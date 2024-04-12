package minion

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion/database"
)

type Minion struct {
	Client string
	Config *Config
	Log    *zap.SugaredLogger

	queues        map[string]*Queue
	notifications chan *Notification
	workers       map[string]registration
	db            *database.Connector
	cron          *cron.Cron
	subs          []func(*Notification)
	listening     bool

	statsEntry cron.EntryID
	statsSubs  []func(Stats)

	cancel context.CancelFunc
}

type Config struct {
	Concurrency     int
	BufferSize      int
	PollingInterval int
	Timeout         int

	Router bool
	Port   int

	Logger *zap.SugaredLogger

	Database    string
	Collection  string
	DatabaseURI string

	RetryCanceled       bool
	ShutdownWaitSeconds int
	Debug               bool
}

func New(client string, cfg *Config) (*Minion, error) {
	db, err := database.New(cfg.DatabaseURI, cfg.Database, cfg.Collection)
	if err != nil {
		return nil, fae.Errorf("creating database: %w", err)
	}

	if cfg.Concurrency == 0 {
		cfg.Concurrency = 5
	}
	if cfg.BufferSize == 0 {
		cfg.BufferSize = 10
	}
	if cfg.PollingInterval == 0 {
		cfg.PollingInterval = 1
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 600 // 10 minutes
	}
	if cfg.ShutdownWaitSeconds == 0 {
		cfg.ShutdownWaitSeconds = 5
	}

	queues := map[string]*Queue{
		"default":  {"default", cfg.Concurrency, cfg.BufferSize, cfg.PollingInterval, make(chan string, cfg.BufferSize)},
		"schedule": {"schedule", cfg.Concurrency, cfg.BufferSize, 1, make(chan string, cfg.BufferSize)},
	}

	return &Minion{
		Client:        client,
		Config:        cfg,
		Log:           cfg.Logger,
		db:            db,
		queues:        queues,
		notifications: make(chan *Notification, cfg.BufferSize*cfg.BufferSize),
		cron:          cron.New(cron.WithSeconds()),
		workers:       make(map[string]registration),
		subs:          []func(*Notification){},
		cancel:        nil,
	}, nil
}

func (m *Minion) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	// m.Log.Infof("starting minion (concurrency=%d/%d)...", m.Concurrency, m.Concurrency*m.Concurrency)
	if m.Config.Debug {
		m.Subscribe(m.debug)
	}

	if err := m.db.UpdateAbandonedJobs(ctx, m.Client); err != nil {
		return fae.Errorf("updating abandoned jobs: %w", err)
	}

	if m.Config.RetryCanceled {
		count, err := m.db.UpdateCancelledJobs(ctx, m.Client)
		if err != nil {
			return fae.Errorf("updating cancelled jobs: %w", err)
		}
		if count > 0 {
			m.Log.Infof("resuming %d cancelled jobs", count)
		}
	}

	for _, queue := range m.queues {
		for w := 0; w < queue.Concurrency; w++ {
			runner := &Runner{
				ID:     w,
				Minion: m,
				Queue:  queue,
			}
			go runner.Run(ctx)
		}

		p := &Producer{Minion: m, Queue: queue}
		p.Run(ctx)
	}

	go func() {
		m.cron.Start()
	}()

	go func() {
		if len(m.subs) > 0 {
			m.listen(ctx)
		}
	}()

	go func() {
		if len(m.statsSubs) == 0 {
			return
		}
		for {
			select {
			case <-time.After(time.Duration(1 * time.Second)):
				m.stats(ctx)
			case <-ctx.Done():
				m.Log.Debugf("minion shutting down")
				return
			}
		}
	}()

	return nil
}

func (m *Minion) Stop() {
	if m.cancel != nil {
		m.cancel()
		<-time.After(time.Duration(m.Config.ShutdownWaitSeconds) * time.Second)
	}
}
