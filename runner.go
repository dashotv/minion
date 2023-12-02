package minion

import (
	"github.com/pkg/errors"
)

type Runner struct {
	ID     int
	Minion *Minion
}

func (r *Runner) Run() {
	for jobID := range r.Minion.queue {
		r.runJob(jobID)
	}
}

func (r *Runner) runJob(jobID string) error {
	d := &JobData{}
	err := r.Minion.db.Find(jobID, d)
	if err != nil {
		return errors.Wrap(err, "finding job")
	}

	w, ok := r.Minion.workers[d.Kind]
	if !ok {
		return errors.Errorf("worker not found for kind: %s", d.Kind)
	}

	job := w.factory.Create(d)
	err = job.Unmarshal()
	if err != nil {
		return errors.Wrap(err, "unmarshaling job")
	}

	attempt := &JobDataAttempt{}
	attempt.Start()
	i := d.AddAttempt(attempt)
	err = r.Minion.db.Save(d)
	if err != nil {
		return errors.Wrap(err, "updating job")
	}

	err = job.Work(r.Minion.Context)
	e := errors.Wrap(err, "running job")
	attempt.Finish(e)
	if err != nil {
		return e
	}

	d.UpdateAttempt(i, attempt)
	err = r.Minion.db.Save(d)
	if err != nil {
		return errors.Wrap(err, "updating job")
	}

	return e
}
