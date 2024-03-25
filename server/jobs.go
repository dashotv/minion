package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/dashotv/minion"
	"github.com/dashotv/minion/database"
)

func setupJobs(s *Server) error {
	mcfg := &minion.Config{
		Logger:      s.Log.Named("minion"),
		Debug:       s.Config.Debug,
		DatabaseURI: s.Config.MongoURI,
		Database:    s.Config.MongoDatabase,
		Collection:  s.Config.MongoCollection,
	}

	m, err := minion.New("minion", mcfg)
	if err != nil {
		return err
	}

	j := &Jobs{
		Minion:       m,
		Log:          s.Log.Named("jobs"),
		DB:           s.DB,
		keepFinished: s.Config.KeepFinishedJobs,
		keepFailed:   s.Config.KeepFailedJobs,
	}
	if _, err := m.ScheduleFunc("0 0 8 * * *", "jobs_cleanup", j.jobs_cleanup); err != nil {
		return err
	}

	s.Jobs = j
	return nil
}

type Jobs struct {
	Minion *minion.Minion
	Log    *zap.SugaredLogger
	DB     *database.Connector

	keepFinished int
	keepFailed   int
}

func (j *Jobs) Start(ctx context.Context) error {
	return j.Minion.Start(ctx)
}

func (j *Jobs) Stop() error {
	j.Minion.Stop()
	return nil
}

func (j *Jobs) jobs_cleanup() error {
	_, err := j.DB.Jobs.Collection.DeleteMany(context.Background(), bson.M{"status": database.StatusFinished, "updated_at": bson.M{"$lt": time.Now().Add(-time.Hour * time.Duration(j.keepFinished))}})
	if err != nil {
		return fmt.Errorf("cleaning up finished jobs: %w", err)
	}

	_, err = j.DB.Jobs.Collection.DeleteMany(context.Background(), bson.M{"status": database.StatusFailed, "updated_at": bson.M{"$lt": time.Now().Add(-time.Hour * time.Duration(j.keepFailed))}})
	if err != nil {
		return fmt.Errorf("cleaning up failed jobs: %w", err)
	}
	return nil
}
