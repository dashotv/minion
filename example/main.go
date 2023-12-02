package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/dashotv/grimoire"

	"github.com/dashotv/minion"
)

func main() {
	min := setupMinion()
	min.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-time.After(3 * time.Second)
		err := min.Enqueue(&Both{Seconds: 25})
		if err != nil {
			log.Fatal("enqueuing job", err)
		}
	}()
	go func() {
		<-time.After(5 * time.Second)
		err := min.Enqueue(&Both{Seconds: 3})
		if err != nil {
			log.Fatal("enqueuing job", err)
		}
	}()

	select {
	case <-c:
		fmt.Println("interrupt")
	case <-time.After(60 * time.Second):
		fmt.Println("done")
	}
}

func setupMinion() *minion.Minion {
	ctx := context.Background()

	db, err := grimoire.New[*minion.JobData]("mongodb://localhost:27017", "minion", "jobs")
	if err != nil {
		log.Fatal("creating db", err)
	}

	ctx = context.WithValue(ctx, "db", db)

	log := zap.NewExample().Sugar()
	cfg := &minion.Config{
		Concurrency: 5,
		Logger:      log,
		Database:    "minion",
		Collection:  "jobs",
		DatabaseURI: "mongodb://localhost:27017",
	}

	m, err := minion.New(ctx, cfg)
	if err != nil {
		log.Fatal("creating minion", err)
	}

	err = minion.Register(m, &Both{})
	if err != nil {
		log.Fatal("registering worker", err)
	}
	err = minion.Register(m, &ScheduleWorker{})
	if err != nil {
		log.Fatal("registering worker", err)
	}

	_, err = m.Schedule("* * * * * *", &SchedulePayload{})
	if err != nil {
		log.Fatal("scheduling worker", err)
	}

	return m
}

type Both struct {
	Seconds int
	minion.WorkerDefaults[*Both]
}

func (b *Both) Kind() string {
	return "both"
}
func (b *Both) Work(ctx context.Context, job *minion.Job[*Both]) error {
	fmt.Printf("both: sleep %d\n", job.Args.Seconds)
	time.Sleep(time.Duration(job.Args.Seconds) * time.Second)
	fmt.Printf("both: done %d\n", job.Args.Seconds)
	return nil
}

type SchedulePayload struct{}

func (p *SchedulePayload) Kind() string {
	return "schedule"
}

type ScheduleWorker struct {
	minion.WorkerDefaults[*SchedulePayload]
}

func (s *ScheduleWorker) Work(ctx context.Context, job *minion.Job[*SchedulePayload]) error {
	db := ctx.Value("db").(*grimoire.Store[*minion.JobData])
	list, err := db.Query().
		Where("status", minion.JobDataStatusRunning).
		Where("kind", "both").
		Asc("created_at").
		Run()
	if err != nil {
		return errors.Wrap(err, "querying jobs")
	}

	fmt.Printf("running jobs: %d\n", len(list))
	for _, item := range list {
		if item.ID == job.ID {
			continue
		}
		fmt.Printf("%s %s %s\n", item.ID.Hex(), item.Kind, item.Args)
	}
	return nil
}
