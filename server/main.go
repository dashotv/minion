package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/dotenv-org/godotenvvault"
)

func main() {
	if err := godotenvvault.Load(); err != nil {
		fmt.Printf("failed to load env: %s\n", err)
		os.Exit(1)
	}

	s, err := setup()
	if err != nil {
		fmt.Printf("setup error: %s\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exit, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()

	s.Log.Infof("starting (port=%d)", s.Config.Port)
	go s.Router.Start(ctx)
	go s.Jobs.Start(ctx)

	select {
	case <-exit.Done():
		s.Log.Warnf("interrupted")
	case <-ctx.Done():
		s.Log.Warnf("canceled")
	}

	if err := s.Router.Stop(); err != nil {
		s.Log.Fatalf("error stopping server: %s\n", err)
	}

	if err := s.Jobs.Stop(); err != nil {
		s.Log.Fatalf("error stopping jobs: %s\n", err)
	}
}
