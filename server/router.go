package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.elastic.co/apm/module/apmechov4/v2"
	"go.infratographer.com/x/echox/echozap"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion/database"
	"github.com/dashotv/minion/static"
)

var pagesize = 25

type Router struct {
	Port                int
	ShutdownWaitSeconds int

	DB   *database.Connector
	Echo *echo.Echo
	Log  *zap.SugaredLogger
	Jobs *Jobs
}

// Router creates and registers the routes of the minion package
func setupRouter(s *Server) error {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(echozap.Middleware(s.Log.Named("echo").Desugar()))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       ".",
		Index:      "index.html", // This is the default html page for your SPA
		Browse:     false,
		HTML5:      true,
		Filesystem: http.FS(static.FS),
	})) // https://echo.labstack.com/docs/middleware/static
	e.Use(apmechov4.Middleware())

	r := &Router{Port: s.Config.Port, Echo: e, DB: s.DB, Log: s.Log.Named("router"), Jobs: s.Jobs}
	e.HTTPErrorHandler = r.customHTTPErrorHandler

	g := e.Group("/jobs")
	g.GET("", r.handleList)
	g.GET("/", r.handleList)
	g.POST("", r.handleCreate)
	g.POST("/", r.handleCreate)
	g.GET("/:id", r.handleGet)
	g.PATCH("/:id", r.handlePatch)
	g.PUT("/:id", r.handleUpdate)
	g.DELETE("/:id", r.handleDelete)

	s.Router = r
	return nil
}

func (r *Router) customHTTPErrorHandler(err error, c echo.Context) {
	r.Log.Errorf("handler error: %v", err)
	he, ok := err.(*echo.HTTPError)
	if !ok {
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	code := he.Code
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, map[string]string{"error": "true", "message": he.Error()})
		}
		if err != nil {
			c.Logger().Error(fae.Errorf("error handling error: %w", err))
		}
	}
}

func (r *Router) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := r.Echo.Start(fmt.Sprintf(":%d", r.Port)); err != nil && err != http.ErrServerClosed {
			r.Echo.Logger.Fatal("shutting down the server")
		}
		cancel()
	}()

	<-ctx.Done()
}

func (r *Router) Stop() error {
	// start gracefully shutdown with a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.ShutdownWaitSeconds)*time.Second)
	defer cancel()

	if err := r.Echo.Shutdown(ctx); err != nil {
		return fae.Errorf("error shutting down the server: %w", err)
	}

	return nil
}

func (r *Router) handleList(c echo.Context) error {
	page := QueryParamInt(c, "page", 1)
	limit := QueryParamInt(c, "limit", pagesize)
	client := c.QueryParam("client")
	skip := (page - 1) * limit
	status := c.QueryParam("status")

	stats, err := r.jobStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": err.Error()})
	}

	q := r.DB.Jobs.Query().Limit(limit).Skip(skip).Desc("created_at")

	if client != "" {
		q = q.Where("client", client)
	}
	if status != "" {
		q = q.Where("status", status)
	}

	list, err := q.Run()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, H{"error": false, "stats": stats, "results": list})
}

func (r *Router) handleCreate(c echo.Context) error {
	kind := c.QueryParam("kind")
	if kind == "" {
		return fae.New("missing kind")
	}

	client := c.QueryParam("client")
	if client == "" {
		return fae.New("missing client")
	}

	j := &database.Model{
		Kind:   kind,
		Client: client,
		Args:   "{}",
		Queue:  "default",
		Status: string(database.StatusPending),
	}

	if err := r.DB.Jobs.Save(j); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, H{"error": false})
}

func (r *Router) handleDelete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return fae.New("missing id")
	}
	hard := c.QueryParam("hard") == "true"

	if id == string(database.StatusPending) && !hard {
		filter := bson.M{"status": database.StatusPending}
		if _, err := r.DB.Jobs.Collection.UpdateMany(context.Background(), filter, bson.M{"$set": bson.M{"status": database.StatusCancelled}}); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, H{"error": false})
	} else if id == string(database.StatusFailed) && hard {
		filter := bson.M{"status": database.StatusFailed}
		if _, err := r.DB.Jobs.Collection.UpdateMany(context.Background(), filter, bson.M{"$set": bson.M{"status": database.StatusArchived}}); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, H{"error": false})
	} else if id == string(database.StatusCancelled) && hard {
		filter := bson.M{"status": database.StatusCancelled}
		if _, err := r.DB.Jobs.Collection.UpdateMany(context.Background(), filter, bson.M{"$set": bson.M{"status": database.StatusArchived}}); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, H{"error": false})
	}

	j, err := r.DB.Jobs.Get(id, &database.Model{})
	if err != nil {
		return err
	}

	j.Status = string(database.StatusCancelled)
	if err := r.DB.Jobs.Save(j); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, H{"error": false})
}

func (r *Router) handleGet(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, H{"error": "not implemented"})
}
func (r *Router) handlePatch(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return fae.New("missing id")
	}
	if err := r.Jobs.Minion.Requeue(id); err != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, H{"error": false})
}
func (r *Router) handleUpdate(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, H{"error": "not implemented"})
}
