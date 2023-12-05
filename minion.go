package minion

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/dashotv/grimoire"
)

func New(ctx context.Context, cfg *Config) (*Minion, error) {
	db, err := grimoire.New[*JobData](cfg.DatabaseURI, cfg.Database, cfg.Collection)
	if err != nil {
		return nil, errors.Wrap(err, "creating job store")
	}

	if cfg.BufferSize == 0 {
		cfg.BufferSize = 100
	}
	if cfg.PollingInterval == 0 {
		cfg.PollingInterval = 5
	}

	return &Minion{
		Context:       ctx,
		Config:        cfg,
		db:            db,
		queue:         make(chan string, cfg.BufferSize),
		cron:          cron.New(cron.WithSeconds()),
		workers:       make(map[string]workerInfo),
		notifications: make(chan *Notification, cfg.BufferSize),
		subs:          []func(*Notification){},
	}, nil
}

func Register[T Payload](m *Minion, worker Worker[T]) error {
	var args T

	kind := args.Kind()
	if _, ok := m.workers[kind]; ok {
		return errors.Errorf("worker already registered for kind: %s", kind)
	}

	m.workers[kind] = workerInfo{
		jobArgs: args,
		factory: &workerFactory[T]{worker: worker},
	}

	return nil
}

type Config struct {
	Concurrency     int
	BufferSize      int
	PollingInterval int
	Logger          *zap.SugaredLogger
	Database        string
	Collection      string
	DatabaseURI     string
	Debug           bool
}

type Minion struct {
	Config  *Config
	Context context.Context
	Log     *zap.SugaredLogger

	workers       map[string]workerInfo
	queue         chan string
	db            *grimoire.Store[*JobData]
	cron          *cron.Cron
	notifications chan *Notification
	subs          []func(*Notification)
	listening     bool
}

func (m *Minion) Start() error {
	// m.Log.Infof("starting minion (concurrency=%d/%d)...", m.Concurrency, m.Concurrency*m.Concurrency)
	if m.Config.Debug {
		m.Subscribe(m.debug)
	}

	for w := 0; w < m.Config.Concurrency; w++ {
		runner := &Runner{
			ID:     w,
			Minion: m,
		}
		go runner.Run()
	}

	go func() {
		m.producer()
	}()

	go func() {
		m.cron.Start()
	}()

	go func() {
		m.listen()
	}()

	return nil
}

func (m *Minion) Enqueue(in Payload) error {
	if in == nil {
		return errors.New("payload is nil")
	}

	args, err := json.Marshal(in)
	if err != nil {
		return errors.Wrap(err, "marshaling job args")
	}

	data := &JobData{
		Args:   string(args),
		Kind:   in.Kind(),
		Status: "pending",
	}

	err = m.db.Save(data)
	if err != nil {
		return errors.Wrap(err, "creating job")
	}

	return nil
}

// Schedule adds (and Registers) a job to the cron scheduler.
func (m *Minion) Schedule(schedule string, in Payload) (cron.EntryID, error) {
	return m.cron.AddFunc(schedule, func() {
		m.Enqueue(in)
	})
}

// Remove removes a job from the cron scheduler.
func (m *Minion) Remove(id cron.EntryID) {
	m.cron.Remove(id)
}

func (m *Minion) Subscribe(f func(*Notification)) {
	m.subs = append(m.subs, f)
}

func (m *Minion) producer() {
	for {
		time.Sleep(time.Duration(m.Config.PollingInterval) * time.Second)

		if len(m.queue) == cap(m.queue) {
			continue
		}

		i := cap(m.queue) - len(m.queue)
		list, err := m.db.Query().Where("status", "pending").Limit(i).Run()
		if err != nil {
			m.Log.Errorf("querying pending jobs: %s", err)
		}

		for _, j := range list {
			m.queue <- j.ID.Hex()
		}
	}
}

func (m *Minion) debug(n *Notification) {
	m.Log.Debugf("event=%s job=%s", n.Event, n.JobID)
}

func (m *Minion) notify(event string, jobID string) {
	if !m.listening {
		m.Log.Warnf("no listeners for notification: %s", event)
		return
	}
	m.notifications <- &Notification{event, jobID}
}

func (m *Minion) listen() {
	m.listening = true
	for n := range m.notifications {
		for _, s := range m.subs {
			s(n)
		}
	}
	m.listening = false
}

type Notification struct {
	Event string
	JobID string
}
