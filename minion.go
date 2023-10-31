package minion

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Minion struct {
	Concurrency int
	Queue       chan Runnable
	Cron        *cron.Cron
	Log         Loggable
}

func New(concurrency int) *Minion {
	return &Minion{
		Concurrency: concurrency,
		Log:         &DefaultLogger{},
		Queue:       make(chan Runnable, concurrency*concurrency),
		Cron:        cron.New(cron.WithSeconds()),
	}
}

func (m *Minion) WithLogger(log Loggable) *Minion {
	m.Log = log
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
func (m *Minion) Schedule(schedule string, job Runnable) (cron.EntryID, error) {
	return m.Cron.AddFunc(schedule, func() {
		m.Queue <- job
	})
}

// Remove removes a job from the cron scheduler.
func (m *Minion) Remove(id cron.EntryID) {
	m.Cron.Remove(id)
}

// Enqueue adds a job to the queue.
func (m *Minion) Enqueue(job Runnable) error {
	m.Queue <- job
	return nil
}

func (m *Minion) run(workerID int, job Runnable) {
	head := fmt.Sprintf("worker=%d", workerID)
	if nameable, ok := job.(Nameable); ok {
		head += fmt.Sprintf(" job=[%s](%s)", nameable.GetID(), nameable.GetName())
	}

	reporter, reportable := job.(Reportable)
	start := time.Now()
	m.Log.Infof("%s: starting", head)

	if reportable {
		err := reporter.Report(ReportableStart, workerID, m)
		if err != nil {
			m.Log.Infof("%s: starting, failed to report: %s", head, err)
			return
		}
	}

	err := job.Run(workerID, m)
	if err != nil {
		m.Log.Infof("%s: failing", head)
		if reportable {
			err := reporter.Report(ReportableError, workerID, m)
			if err != nil {
				m.Log.Infof("%s: failing, also failed to report: %s", head, err)
				return
			}
		}
		return
	}

	m.Log.Infof("%s: finishing", head)
	if reportable {
		err := reporter.Report(ReportableFinish, workerID, m)
		if err != nil {
			m.Log.Infof("%s: finishing, failed to report: %s", head, err)
			return
		}
	}

	diff := time.Since(start)
	m.Log.Infof("%s: duration: %s", head, diff)
	if reportable {
		err := reporter.Report(ReportableDuration, workerID, m)
		if err != nil {
			m.Log.Infof("%s: duration, failed to report: %s", head, err)
			return
		}
	}
}
