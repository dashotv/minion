package minion

import (
	"context"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/dashotv/grimoire"
)

type Minion struct {
	Config  *Config
	Context context.Context
	Log     *zap.SugaredLogger

	queues        map[string]*Queue
	notifications chan *Notification
	workers       map[string]registration
	db            *grimoire.Store[*Model]
	cron          *cron.Cron
	subs          []func(*Notification)
	listening     bool
}

type Config struct {
	Concurrency     int
	BufferSize      int
	PollingInterval int
	Timeout         int
	RetryCanceled   bool
	Logger          *zap.SugaredLogger
	Database        string
	Collection      string
	DatabaseURI     string
	Debug           bool
}

func New(ctx context.Context, cfg *Config) (*Minion, error) {
	db, err := grimoire.New[*Model](cfg.DatabaseURI, cfg.Database, cfg.Collection)
	if err != nil {
		return nil, errors.Wrap(err, "creating job store")
	}

	if cfg.Concurrency == 0 {
		cfg.Concurrency = 1
	}
	if cfg.BufferSize == 0 {
		cfg.BufferSize = 100
	}
	if cfg.PollingInterval == 0 {
		cfg.PollingInterval = 1
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 600 // 10 minutes
	}

	queues := map[string]*Queue{
		"default":  {"default", cfg.Concurrency, cfg.BufferSize, cfg.PollingInterval, make(chan string, cfg.BufferSize)},
		"schedule": {"schedule", cfg.Concurrency, cfg.BufferSize, 1, make(chan string, cfg.BufferSize)},
	}

	return &Minion{
		Context:       ctx,
		Config:        cfg,
		Log:           cfg.Logger,
		db:            db,
		queues:        queues,
		notifications: make(chan *Notification, cfg.BufferSize*cfg.BufferSize),
		cron:          cron.New(cron.WithSeconds()),
		workers:       make(map[string]registration),
		subs:          []func(*Notification){},
	}, nil
}

func (m *Minion) Start() error {
	// m.Log.Infof("starting minion (concurrency=%d/%d)...", m.Concurrency, m.Concurrency*m.Concurrency)
	if m.Config.Debug {
		m.Subscribe(m.debug)
	}

	for _, queue := range m.queues {
		for w := 0; w < queue.Concurrency; w++ {
			runner := &Runner{
				ID:     w,
				Minion: m,
				Queue:  queue,
			}
			go runner.Run()
		}
		p := &Producer{Minion: m, Queue: queue}
		p.Run()
	}

	go func() {
		m.cron.Start()
	}()

	go func() {
		if len(m.subs) > 0 {
			m.listen()
		}
	}()

	go func() {
		if m.Config.RetryCanceled {
			res, err := m.db.Collection.UpdateMany(m.Context, bson.M{"status": "canceled"}, bson.M{"$set": bson.M{"status": "pending"}})
			if err != nil {
				m.Log.Errorf("querying canceled jobs: %s", err)
				return
			}
			if res.ModifiedCount > 0 {
				m.Log.Infof("resuming %d canceled jobs", res.ModifiedCount)
			}
		}
	}()

	return nil
}

// Schedule adds (and Registers) a job to the cron scheduler.
func (m *Minion) Schedule(schedule string, in Payload) (cron.EntryID, error) {
	return m.cron.AddFunc(schedule, func() {
		m.notify("job:scheduled", "-", in.Kind())
		m.enqueueTo("schedule", in)
	})
}

// Remove removes a job from the cron scheduler.
func (m *Minion) Remove(id cron.EntryID) {
	m.cron.Remove(id)
}
