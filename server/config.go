package main

import (
	"fmt"
	"os"

	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"

	"github.com/dashotv/minion/database"
)

type Config struct {
	Debug  bool   `env:"DEBUG" default:"false"`
	Port   int    `env:"PORT" default:"9010"`
	Logger string `env:"LOGGER" default:"dev"`

	MongoURI        string `env:"MONGO_URI" default:"mongodb://localhost:27017"`
	MongoDatabase   string `env:"MONGO_DATABASE" default:"minion_development"`
	MongoCollection string `env:"MONGO_COLLECTION" default:"jobs"`

	ShutdownWaitSeconds int `env:"SHUTDOWN_WAIT_SECONDS" default:"5"`
	KeepFinishedJobs    int `env:"KEEP_FINISHED_JOBS" default:"2"` // hours
	KeepFailedJobs      int `env:"KEEP_FAILED_JOBS" default:"48"`  // hours
}

func setupLogger(s *Server) error {
	switch s.Config.Logger {
	case "dev":
		isTTY := term.IsTerminal(int(os.Stderr.Fd()))
		verbosity := 1
		logStdoutWriter := zapcore.Lock(os.Stderr)
		log := zap.New(zapcore.NewCore(logging.NewEncoder(verbosity, isTTY), logStdoutWriter, zapcore.DebugLevel))
		s.Log = log.Sugar().Named("server")
		return nil
	case "release":
		zapcfg := zap.NewProductionConfig()
		zapcfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		log, err := zap.NewProduction()
		s.Log = log.Sugar().Named("server")
		return err
	default:
		return fmt.Errorf("unknown logger: %s", s.Config.Logger)
	}
}

func setupDatabase(s *Server) error {
	con, err := database.New(s.Config.MongoURI, s.Config.MongoDatabase, s.Config.MongoCollection)
	if err != nil {
		return fmt.Errorf("creating job store: %w", err)
	}
	s.DB = con
	return nil
}
