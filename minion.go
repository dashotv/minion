package minion

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Minion struct {
	Concurrency int
	Queue       chan Runnable
	Cron        *cron.Cron
	Log         *zap.SugaredLogger
}

type MinionFunc func(id int, log *zap.SugaredLogger) error

func NewMinion(concurrency int, log *zap.SugaredLogger) *Minion {
	return &Minion{
		Concurrency: concurrency,
		Log:         log,
		Queue:       make(chan Runnable, concurrency*concurrency),
		Cron:        cron.New(cron.WithSeconds()),
	}
}

func (m *Minion) Start() error {
	m.Log.Infof("starting minion (concurrency=%d)...", m.Concurrency)

	for w := 0; w < m.Concurrency; w++ {
		worker := &Worker{
			ID:     w,
			Log:    m.Log.Named(fmt.Sprintf("worker=%d", w)),
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

func (m *Minion) Schedule(schedule string, job Runnable) (cron.EntryID, error) {
	return m.Cron.AddFunc(schedule, func() {
		m.Queue <- job
	})
}

func (m *Minion) Remove(id cron.EntryID) {
	m.Cron.Remove(id)
}

func (m *Minion) Enqueue(job Runnable) error {
	m.Queue <- job
	return nil
}

func (m *Minion) run(workerID int, job Runnable) {
	id := ""
	if idable, ok := job.(Idable); ok {
		id = fmt.Sprintf("[%s]", idable.GetID())
	}

	name := ""
	if nameable, ok := job.(Nameable); ok {
		name = fmt.Sprintf("'%s'", nameable.GetName())
	}

	reporter, reportable := job.(Reportable)
	start := time.Now()
	m.Log.Infof("job:%s%s: starting", id, name)

	if reportable {
		err := reporter.Report(ReportableStart, workerID, m)
		if err != nil {
			m.Log.Infof("job:%s%s: starting, failed to report: %s", id, name, err)
			return
		}
	}

	err := job.Run(workerID, m)
	if err != nil {
		m.Log.Infof("job:%s%s: failing", id, name)
		if reportable {
			err := reporter.Report(ReportableError, workerID, m)
			if err != nil {
				m.Log.Infof("job:%s%s: failing, also failed to report: %s", id, name, err)
				return
			}
		}
		return
	}

	m.Log.Infof("job:%s%s: finishing", id, name)
	if reportable {
		err := reporter.Report(ReportableFinish, workerID, m)
		if err != nil {
			m.Log.Infof("job:%s%s: finishing, failed to report: %s", id, name, err)
			return
		}
	}

	diff := time.Since(start)
	m.Log.Infof("job:%s%s: duration: %s", id, name, diff)
	if reportable {
		err := reporter.Report(ReportableDuration, workerID, m)
		if err != nil {
			m.Log.Infof("job:%s%s: duration, failed to report: %s", id, name, err)
			return
		}
	}
}

// func (m *Minion) Add(name string, f MinionFunc) error {
// 	// 	j := &MinionJob{
// 	// 		Name: name,
// 	// 	}
// 	//
// 	// 	err := db.MinionJob.Save(j)
// 	// 	if err != nil {
// 	// 		return errors.Wrap(err, "failed to save minion job")
// 	// 	}
//
// 	mf := func(id int, log *zap.SugaredLogger) error {
// 		log.Infof("starting %s: %s", name, j.ID.Hex())
// 		err := f(id, log)
//
// 		j.ProcessedAt = time.Now()
// 		if err != nil {
// 			log.Errorf("processing %s: %s: %s", name, j.ID.Hex(), err)
// 			j.Error = errors.Wrap(err, "failed to run minion job").Error()
// 		}
//
// 		err = db.MinionJob.Update(j)
// 		if err != nil {
// 			log.Errorf("error %s: %s: %s", name, j.ID.Hex(), err)
// 			return errors.Wrap(err, "failed to save minion job")
// 		}
//
// 		log.Infof("finished %s: %s", name, j.ID.Hex())
// 		return nil
// 	}
//
// 	m.Queue <- &Job{ID: j.ID.Hex(), Func: mf}
// 	return nil
// }
//
// func (m *Minion) AddCron(spec, name string, f MinionFunc) error {
// 	_, err := m.Cron.AddFunc(spec, func() {
// 		minion.Add(name, func(id int, log *zap.SugaredLogger) error {
// 			return f(id, log)
// 		})
// 	})
//
// 	if err != nil {
// 		return errors.Wrap(err, "adding cron function")
// 	}
//
// 	return nil
// }
