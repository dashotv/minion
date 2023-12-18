package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/dashotv/grimoire"

	"github.com/dashotv/minion"
)

func main() {
	min := setupMinion()
	min.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer func() {
		min.Log.Info("stopping")
		db, err := grimoire.New[*minion.Model]("mongodb://localhost:27017", "minion", "jobsExample")
		if err != nil {
			Fatal("creating db: %s", err)
		}
		cur, err := db.Collection.Find(context.Background(), bson.M{"status": "finished", "$expr": bson.M{"$gt": bson.A{bson.M{"$size": "$attempts"}, "1"}}})
		if err != nil {
			Fatal("finding: %s", err)
		}
		defer cur.Close(context.Background())
		for cur.Next(context.Background()) {
			var d *minion.Model
			err := cur.Decode(&d)
			if err != nil {
				Fatal("decoding: %s", err)
			}
			fmt.Printf("%s %s %d\n", d.ID.Hex(), d.Kind, len(d.Attempts))
		}
	}()

	go func() {
		<-time.After(3 * time.Second)
		err := min.Enqueue(&Sleeper{Seconds: 25})
		if err != nil {
			min.Log.Errorf("enqueuing job: %s", err)
		}
	}()
	go func() {
		<-time.After(5 * time.Second)
		err := min.Enqueue(&Sleeper{Seconds: 3})
		if err != nil {
			min.Log.Errorf("enqueuing job: %s", err)
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
		<-time.After(25 * time.Second)
		for i := 0; i < 100; i++ {
			err := min.Enqueue(&Number{Number: i})
			if err != nil {
				min.Log.Errorf("enqueuing job: %s", err)
			}
		}
		min.Log.Info("jobs queued")
	}()

	select {
	case <-c:
		fmt.Println("interrupt")
	case <-time.After(120 * time.Second):
		fmt.Println("done")
	}
}

func Fatal(format string, err error) {
	fmt.Printf("fatal: %s\n", fmt.Sprintf(format, err))
	os.Exit(1)
}

func setupMinion() *minion.Minion {
	ctx := context.Background()

	db, err := grimoire.New[*minion.Model]("mongodb://localhost:27017", "minion", "jobsExample")
	if err != nil {
		Fatal("creating db: %s", err)
	}

	if _, err := db.Collection.DeleteMany(context.Background(), bson.M{}); err != nil {
		Fatal("clearing db: %s", err)
	}

	ctx = context.WithValue(ctx, "db", db)

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
		Database:    "minion",
		Collection:  "jobsExample",
		DatabaseURI: "mongodb://localhost:27017",
	}

	m, err := minion.New(ctx, cfg)
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
	err = minion.Register(m, &ScheduleWorker{})
	if err != nil {
		Fatal("registering worker: %s", err)
	}

	_, err = m.Schedule("*/5 * * * * *", &SchedulePayload{})
	if err != nil {
		Fatal("scheduling worker: %s", err)
	}

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
		return errors.New("random error")
	}
	if i == 3 {
		panic("random panic")
	}
	return nil
}

type SchedulePayload struct{}

func (p *SchedulePayload) Kind() string { return "schedule" }

type ScheduleWorker struct {
	minion.WorkerDefaults[*SchedulePayload]
}

func (s *ScheduleWorker) Work(ctx context.Context, job *minion.Job[*SchedulePayload]) error {
	db := ctx.Value("db").(*grimoire.Store[*minion.Model])

	for _, n := range []string{"default", "schedule", "number"} {
		err := stats(db, n)
		if err != nil {
			return errors.Wrap(err, "stats")
		}
	}

	// 	total, err := db.Count(bson.M{"status": bson.M{"$nin": bson.A{minion.StatusFinished, minion.StatusFailed}}})
	// 	if err != nil {
	// 		return errors.Wrap(err, "counting jobs")
	// 	}
	// 	scheduled, err := db.Count(bson.M{"queue": "schedule", "status": bson.M{"$nin": bson.A{minion.StatusFinished, minion.StatusFailed}}})
	// 	if err != nil {
	// 		return errors.Wrap(err, "counting jobs")
	// 	}
	// 	running, err := db.Count(bson.M{"queue": "default", "status": minion.StatusRunning})
	// 	if err != nil {
	// 		return errors.Wrap(err, "counting jobs")
	// 	}
	//
	// 	a := []string{string(minion.StatusFinished), string(minion.StatusFailed)}
	// 	list, err := db.Query().
	// 		Where("queue", "default").
	// 		In("status", a).
	// 		Limit(-1).
	// 		Run()
	// 	if err != nil {
	// 		return errors.Wrap(err, "querying jobs")
	// 	}
	//
	// 	duration := 0.0
	// 	if len(list) > 0 {
	// 		sum := 0
	// 		for _, item := range list {
	// 			for _, a := range item.Attempts {
	// 				sum += int(a.Duration)
	// 			}
	// 		}
	// 		duration = float64(sum) / float64(len(list))
	// 	}
	//
	// 	fmt.Printf("jobs: %s running %d / scheduled %d / total %d / duration: %5.2f (%d)\n", time.Now(), running, scheduled, total, duration, len(list))
	return nil
}

func stats(db *grimoire.Store[*minion.Model], queue string) error {
	failed, err := db.Count(bson.M{"queue": queue, "status": minion.StatusFailed})
	if err != nil {
		return errors.Wrap(err, "counting jobs")
	}
	running, err := db.Count(bson.M{"queue": queue, "status": minion.StatusRunning})
	if err != nil {
		return errors.Wrap(err, "counting jobs")
	}
	pending, err := db.Count(bson.M{"queue": queue, "status": bson.M{"$nin": bson.A{minion.StatusFinished, minion.StatusFailed}}})
	if err != nil {
		return errors.Wrap(err, "counting jobs")
	}
	finished, err := db.Count(bson.M{"queue": queue, "status": bson.M{"$in": bson.A{minion.StatusFinished, minion.StatusFailed}}})
	if err != nil {
		return errors.Wrap(err, "counting jobs")
	}

	fmt.Printf("stats: %s: %d %d %d %d\n", queue, running, pending, failed, finished)
	return nil
}
