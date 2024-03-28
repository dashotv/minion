package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion"
	"github.com/dashotv/minion/database"
)

var mongoURI = "mongodb://localhost:27017"
var mongoDatabase = "minion_development"
var mongoCollection = "jobs"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	min := setupMinion()
	min.Start(ctx)

	exit, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()

	defer func() {
		min.Log.Info("stopping")
	}()

	go func() {
		select {
		case <-time.After(3 * time.Second):
			err := min.Enqueue(&Sleeper{Seconds: 25})
			if err != nil {
				min.Log.Errorf("enqueuing job: %s", err)
			}
		case <-ctx.Done():
			return
		}
	}()
	go func() {
		select {
		case <-time.After(5 * time.Second):
			err := min.Enqueue(&Sleeper{Seconds: 3})
			if err != nil {
				min.Log.Errorf("enqueuing job: %s", err)
			}
		case <-ctx.Done():
			return
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			err := min.Enqueue(&Number{Number: i})
			if err != nil {
				min.Log.Errorf("enqueuing job: %s", err)
			}
		}
		min.Log.Info("jobs queued")
	}()
	go func() {
		select {
		case <-time.After(25 * time.Second):
			for i := 0; i < 100; i++ {
				err := min.Enqueue(&Number{Number: i})
				if err != nil {
					min.Log.Errorf("enqueuing job: %s", err)
				}
			}
			min.Log.Info("jobs queued")
		case <-ctx.Done():
			return
		}
	}()

	select {
	case <-exit.Done():
		fmt.Println("interrupt")
	case <-time.After(120 * time.Second):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("context done")
	}

	min.Stop()
}

func Fatal(format string, err error) {
	fmt.Printf("fatal: %s\n", fmt.Sprintf(format, err))
	os.Exit(1)
}

func setupMinion() *minion.Minion {
	err := resetDatabase()
	if err != nil {
		Fatal("resetting database: %s", err)
	}

	dev, err := zap.NewDevelopment()
	if err != nil {
		Fatal("configuring logger: %s", err)
	}
	log := dev.Sugar()
	cfg := &minion.Config{
		// Debug:       true,
		Concurrency: 5,
		BufferSize:  10,
		Logger:      log,
		DatabaseURI: mongoURI,
		Database:    mongoDatabase,
		Collection:  mongoCollection,
	}

	m, err := minion.New("testing", cfg)
	if err != nil {
		Fatal("creating minion: %s", err)
	}
	m.Queue("number", 3, 10, 1)

	err = minion.Register(m, &Sleeper{})
	if err != nil {
		Fatal("registering worker: %s", err)
	}
	err = minion.RegisterWithQueue(m, &Number{}, "number")
	if err != nil {
		Fatal("registering worker: %s", err)
	}

	m.SubscribeStats(func(s minion.Stats) {
		fmt.Printf("STATS: %+v\n", s)
	})

	return m
}

type Sleeper struct {
	minion.WorkerDefaults[*Sleeper]
	Seconds int
}

func (b *Sleeper) Kind() string { return "Sleeper" }
func (b *Sleeper) Timeout(*minion.Job[*Sleeper]) time.Duration {
	return time.Duration(10) * time.Second
}
func (b *Sleeper) Work(ctx context.Context, job *minion.Job[*Sleeper]) error {
	fmt.Printf("both: sleep %d\n", job.Args.Seconds)
	time.Sleep(time.Duration(job.Args.Seconds) * time.Second)
	fmt.Printf("both: done %d\n", job.Args.Seconds)
	return nil
}

type Number struct {
	minion.WorkerDefaults[*Number]
	Number int
}

func (n *Number) Kind() string { return "number" }
func (n *Number) Work(ctx context.Context, job *minion.Job[*Number]) error {
	i := rand.Intn(5)
	time.Sleep(time.Duration(i) * time.Second)
	fmt.Printf("number: %d %d\n", job.Args.Number, i)
	if i == 4 {
		return fae.New("random error")
	}
	if i == 3 {
		panic("random panic")
	}
	return nil
}

func resetDatabase() error {
	con, err := database.New(mongoURI, mongoDatabase, mongoCollection)
	if err != nil {
		return fae.Wrap(err, "creating database")
	}

	_, err = con.Jobs.Collection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return fae.Wrap(err, "deleting all jobs")
	}

	return nil
}
