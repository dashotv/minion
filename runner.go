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

func (r *Runner) loadJob(jobID string) (wrapped, *JobData, error) {

	d := &JobData{}
	err := r.Minion.db.Find(jobID, d)
	if err != nil {
		return nil, d, errors.Wrap(err, fmt.Sprintf("finding job: %s", jobID))
	}

	w, ok := r.Minion.workers[d.Kind]
	if !ok {
		e := errors.Errorf("worker not found for kind: %s", d.Kind)
		d.Status = "cancelled"
		_ = r.Minion.db.Save(d)
		return nil, d, e
	}

	job := w.factory.Create(d)
	err = job.Unmarshal()
	if err != nil {
		return nil, d, errors.Wrap(err, "unmarshaling job")
	}

	return job, d, nil
}

func (r *Runner) runJobAttempt(jobID string, d *JobData, job wrapped) error {
	attempt := &JobDataAttempt{}
	attempt.Start()
	i := d.AddAttempt(attempt)
	err := r.Minion.db.Save(d)
	if err != nil {
		return errors.Wrap(err, "updating job")
	}

	r.Minion.notify("job:start", jobID, d.Kind)
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

// runJob runs a job
// we use named returns so recover can set the error
func (r *Runner) runJob(jobID string) (err error) {
	r.Minion.notify("job:load", jobID, "-")

	job, d, err := r.loadJob(jobID)
	if err != nil {
		return errors.Wrap(err, "loading job")
	}

	defer func() {
		if recovery := recover(); recovery != nil {
			err = errors.Errorf("running job: panic: %v", recovery)
		}
		if err != nil {
			r.Minion.notify("job:fail", jobID, d.Kind)
		} else {
			r.Minion.notify("job:success", jobID, d.Kind)
		}
	}()

	err = r.runJobAttempt(jobID, d, job)
	return err
}
