package minion

import (
	"context"
	"encoding/json"

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

	return &Minion{
		Context:       ctx,
		Concurrency:   cfg.Concurrency,
		Log:           cfg.Logger,
		Debug:         cfg.Debug,
		db:            db,
		queue:         make(chan string, cfg.Concurrency*cfg.Concurrency),
		cron:          cron.New(cron.WithSeconds()),
		workers:       make(map[string]workerInfo),
		notifications: make(chan *Notification),
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
	Concurrency int
	Logger      *zap.SugaredLogger
	Database    string
	Collection  string
	DatabaseURI string
	Debug       bool
}

type Minion struct {
	Debug         bool
	Context       context.Context
	Concurrency   int
	Log           *zap.SugaredLogger
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
	if m.Debug {
		m.Subscribe(m.debug)
	}

	for w := 0; w < m.Concurrency; w++ {
		runner := &Runner{
			ID:     w,
			Minion: m,
		}
		go runner.Run()
	}

	go func() {
		m.cron.Start()
	}()

	go func() {
		m.listen()
	}()

	return nil
}

func (m *Minion) Enqueue(in Payload) error {
	args, err := json.Marshal(in)
	if err != nil {
		return errors.Wrap(err, "marshaling job args")
	}

	data := &JobData{
		Args: string(args),
		Kind: in.Kind(),
	}

	err = m.db.Save(data)
	if err != nil {
		return errors.Wrap(err, "creating job")
	}

	m.queue <- data.ID.Hex()
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
