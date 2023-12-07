package minion

import (
	"fmt"

	"github.com/pkg/errors"
)

type Runner struct {
	ID     int
	Minion *Minion
	Queue  chan string
}

func (r *Runner) Run() {
	for jobID := range r.Queue {
		err := r.runJob(jobID)
		if err != nil {
			r.Minion.Log.Errorf("minon:runner:error: %s", err)
		}
	}
}

func (r *Runner) runJob(jobID string) error {
	r.Minion.notify("job:load", jobID, "-")

	d := &JobData{}
	err := r.Minion.db.Find(jobID, d)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("finding job: %s", jobID))
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
	r.Minion.notify("job:start", jobID, d.Kind)
	i := d.AddAttempt(attempt)
	err = r.Minion.db.Save(d)
	if err != nil {
		return errors.Wrap(err, "updating job")
	}

	err = job.Work(r.Minion.Context)
	e := errors.Wrap(err, "running job")
	attempt.Finish(e)
	r.Minion.notify("job:finish", jobID, d.Kind)

	d.UpdateAttempt(i, attempt)
	err = r.Minion.db.Save(d)
	if err != nil {
		return errors.Wrap(err, "updating job")
	}

	if e != nil {
		r.Minion.notify("job:fail", jobID, d.Kind)
	} else {
		r.Minion.notify("job:success", jobID, d.Kind)
	}
	return e
}
