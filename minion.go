package minion

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Minion struct {
	Concurrency int
	Queue       chan *Job
	Cron        *cron.Cron
	Log         Loggable
	Jobs        map[string]Func
	Reporter    Reportable
}

type Func func() error

func New(concurrency int) *Minion {
	return &Minion{
		Concurrency: concurrency,
		Log:         &DefaultLogger{},
		Queue:       make(chan *Job, concurrency*concurrency),
		Cron:        cron.New(cron.WithSeconds()),
		Jobs:        make(map[string]Func),
		Reporter:    &DefaultReporter{},
	}
}

func (m *Minion) WithLogger(log Loggable) *Minion {
	m.Log = log
	return m
}
func (m *Minion) WithReporter(reporter Reportable) *Minion {
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

func (m *Minion) run(workerID int, j *Job) {
	head := fmt.Sprintf("worker=%d name=%s", workerID, j.Name)

	start := time.Now()
	m.Log.Infof("%s: starting", head)

	err := m.Report(ReportableStart, j.Name, workerID)
	if err != nil {
		m.Log.Errorf("%s: starting, failed to report: %s", head, err)
	}

	err = j.Func()
	if err != nil {
		m.Log.Errorf("%s: failing: %s", head, err)
		rerr := m.Report(ReportableError, j.Name, workerID)
		if rerr != nil {
			m.Log.Errorf("%s: failing, also failed to report: %s", head, rerr)
		}
	}

	diff := time.Since(start)
	m.Log.Infof("%s: duration: %s", head, diff)
	err = m.Report(ReportableDuration, j.Name, workerID)
	if err != nil {
		m.Log.Errorf("%s: duration, failed to report: %s", head, err)
	}

	m.Log.Infof("%s: finishing", head)
	err = m.Report(ReportableFinish, j.Name, workerID)
	if err != nil {
		m.Log.Errorf("%s: finishing, failed to report: %s", head, err)
	}
}
