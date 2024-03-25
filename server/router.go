package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.infratographer.com/x/echox/echozap"

	"github.com/dashotv/minion/database"
	"github.com/dashotv/minion/static"
)

var pagesize = 25

type Router struct {
	Port                int
	ShutdownWaitSeconds int

	DB   *database.Connector
	Echo *echo.Echo
}

// Router creates and registers the routes of the minion package
func setupRouter(s *Server) error {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(echozap.Middleware(s.Log.Named("router").Desugar()))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       ".",
		Index:      "index.html", // This is the default html page for your SPA
		Browse:     false,
		HTML5:      true,
		Filesystem: http.FS(static.FS),
	})) // https://echo.labstack.com/docs/middleware/static

	r := &Router{Port: s.Config.Port, Echo: e, DB: s.DB}

	e.GET("/jobs", r.handleList)
	e.POST("/jobs", r.handleCreate)
	e.GET("/jobs/:id", r.handleGet)
	e.PATCH("/jobs/:id", r.handlePatch)
	e.PUT("/jobs/:id", r.handleUpdate)
	e.DELETE("/jobs/:id", r.handleDelete)

	s.Router = r
	return nil
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
		return fmt.Errorf("error shutting down the server: %w", err)
	}

	return nil
}

func (r *Router) handleList(c echo.Context) error {
	page := QueryParamInt(c, "page", 1)
	limit := QueryParamInt(c, "limit", pagesize)
	skip := (page - 1) * limit

	stats, err := r.jobStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": err.Error()})
	}

	list, err := r.DB.Jobs.Query().Limit(limit).Skip(skip).Desc("created_at").Run()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, H{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, H{"error": false, "stats": stats, "results": list})
}
func (r *Router) handleCreate(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, H{"error": "not implemented"})
}
func (r *Router) handleGet(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, H{"error": "not implemented"})
}
func (r *Router) handlePatch(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, H{"error": "not implemented"})
}
func (r *Router) handleUpdate(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, H{"error": "not implemented"})
}
func (r *Router) handleDelete(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, H{"error": "not implemented"})
}
