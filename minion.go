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
		cfg.PollingInterval = 1
	}

	queues := map[string]chan string{
		"default":  make(chan string, cfg.BufferSize),
		"schedule": make(chan string, cfg.BufferSize),
	}

	return &Minion{
		Context:       ctx,
		Config:        cfg,
		Log:           cfg.Logger,
		db:            db,
		queues:        queues,
		notifications: make(chan *Notification, cfg.BufferSize*cfg.BufferSize),
		cron:          cron.New(cron.WithSeconds()),
		workers:       make(map[string]workerInfo),
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

	queues        map[string]chan string
	notifications chan *Notification
	workers       map[string]workerInfo
	db            *grimoire.Store[*JobData]
	cron          *cron.Cron
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
			Queue:  m.queues["default"],
		}
		go runner.Run()
	}
	for w := 0; w < m.Config.Concurrency; w++ {
		runner := &Runner{
			ID:     w,
			Minion: m,
			Queue:  m.queues["schedule"],
		}
		go runner.Run()
	}

	go func() {
		m.producer("default", m.Config.PollingInterval)
	}()

	go func() {
		m.producer("schedule", 1)
	}()

	go func() {
		m.cron.Start()
	}()

	go func() {
		if len(m.subs) > 0 {
			m.listen()
		}
	}()

	return nil
}

func (m *Minion) Enqueue(in Payload) error {
	return m.enqueueTo("default", in)
}

func (m *Minion) enqueueTo(queue string, in Payload) error {
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
		Queue:  queue,
	}

	err = m.db.Save(data)
	if err != nil {
		return errors.Wrap(err, "creating job")
	}

	m.notify("job:created", data.ID.Hex(), data.Kind)
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

func (m *Minion) Subscribe(f func(*Notification)) {
	m.subs = append(m.subs, f)
}

func (m *Minion) producer(queue string, interval int) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		if channelBufferFull(m.queues[queue]) {
			continue
		}

		i := channelBufferRemaining(m.queues[queue])
		list, err := m.db.Query().Where("queue", queue).Where("status", "pending").Asc("created_at").Limit(i).Run()
		if err != nil {
			m.Log.Errorf("querying pending jobs: %s", err)
		}

		for _, j := range list {
			j.Status = "queued"
			err = m.db.Save(j)
			if err != nil {
				m.Log.Errorf("updating job: %s", err)
			}

			m.notify("job:queued", j.ID.Hex(), j.Kind)
			m.queues[queue] <- j.ID.Hex()
		}
	}
}

func (m *Minion) debug(n *Notification) {
	m.Log.Debugf("event=%s job=%s kind=%s", n.Event, n.JobID, n.Kind)
}

func (m *Minion) notify(event string, jobID string, kind string) {
	if !m.listening {
		// m.Log.Warnf("no listeners for notification: %s", event)
		return
	}
	if channelBufferFull(m.notifications) {
		m.Log.Debugf("notification buffer full: %d", cap(m.notifications))
	}
	m.notifications <- &Notification{event, jobID, kind}
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
	Kind  string
}
