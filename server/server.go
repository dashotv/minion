package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/caarlos0/env/v10"

	"github.com/dashotv/minion/database"
)

type Server struct {
	Config *Config
	Log    *zap.SugaredLogger
	DB     *database.Connector
	Jobs   *Jobs
	Router *Router
}

func setup() (*Server, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	s := &Server{Config: cfg}

	funcs := []func(*Server) error{
		setupLogger,
		setupDatabase,
		setupJobs,
		setupRouter,
	}
	for _, f := range funcs {
		if err := f(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}
