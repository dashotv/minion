package minion

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Minion struct {
	Concurrency int
	Queue       chan Func
	Cron        *cron.Cron
	Log         Loggable
	Reporter    func(ReportableType, int, *Minion) error
}

type Func func() error

type ReportableType int

const (
	ReportableStart ReportableType = iota
	ReportableFinish
	ReportableError
	ReportableDuration
)

func New(concurrency int) *Minion {
	return &Minion{
		Concurrency: concurrency,
		Log:         &DefaultLogger{},
		Queue:       make(chan Func, concurrency*concurrency),
		Cron:        cron.New(cron.WithSeconds()),
	}
}

func (m *Minion) WithLogger(log Loggable) *Minion {
	m.Log = log
	return m
}
func (m *Minion) WithReporter(reporter func(ReportableType, int, *Minion) error) *Minion {
	m.Reporter = reporter
	return m
}

// Start starts the minion.
func (m *Minion) Start() error {
	m.Log.Infof("starting minion (concurrency=%d/%d)...", m.Concurrency, m.Concurrency*m.Concurrency)

	for w := 0; w < m.Concurrency; w++ {
		worker := &Worker{
			ID:     w,
			Queue:  m.Queue,
			Runner: m.run,
		}
		go worker.Run()
	}

	go func() {
		m.Cron.Start()
	}()

	return nil
}

// Schedule adds a job to the cron scheduler.
func (m *Minion) Schedule(schedule string, f Func) (cron.EntryID, error) {
	return m.Cron.AddFunc(schedule, func() {
		m.Queue <- f
	})
}

// Remove removes a job from the cron scheduler.
func (m *Minion) Remove(id cron.EntryID) {
	m.Cron.Remove(id)
}

// Enqueue adds a job to the queue.
func (m *Minion) Enqueue(f Func) {
	m.Queue <- f
}

func (m *Minion) Report(t ReportableType, workerID int) error {
	if m.Reporter == nil {
		return nil
	}
	return m.Reporter(t, workerID, m)
}

func (m *Minion) run(workerID int, f Func) {
	head := fmt.Sprintf("worker=%d", workerID)

	start := time.Now()
	m.Log.Infof("%s: starting", head)

	err := m.Report(ReportableStart, workerID)
	if err != nil {
		m.Log.Errorf("%s: starting, failed to report: %s", head, err)
		return
	}

	err = f()
	if err != nil {
		m.Log.Errorf("%s: failing: %s", head, err)
		rerr := m.Report(ReportableError, workerID)
		if rerr != nil {
			m.Log.Errorf("%s: failing, also failed to report: %s", head, rerr)
			return
		}
		return
	}

	diff := time.Since(start)
	m.Log.Infof("%s: duration: %s", head, diff)
	err = m.Report(ReportableDuration, workerID)
	if err != nil {
		m.Log.Errorf("%s: duration, failed to report: %s", head, err)
		return
	}

	m.Log.Infof("%s: finishing", head)
	err = m.Report(ReportableFinish, workerID)
	if err != nil {
		m.Log.Errorf("%s: finishing, failed to report: %s", head, err)
		return
	}
}
