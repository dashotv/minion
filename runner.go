package minion

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion/database"
)

type Runner struct {
	ID     int
	Minion *Minion
	Queue  *Queue
}

func (r *Runner) Run(ctx context.Context) {
	for jobID := range r.Queue.channel {
		err := r.runJob(ctx, jobID)
		if err != nil {
			m := err.Error()
			if len(m) > 100 {
				m = m[:97] + "..."
			}
			r.Minion.Log.Errorf("runner: %s", err)
		}
	}
}

// runJob runs a job
func (r *Runner) runJob(ctx context.Context, jobID string) (err error) {
	r.Minion.notify("job:load", jobID, "-")

	job, d, err := r.loadJob(jobID)
	if err != nil {
		return fae.Wrap(err, "loading job")
	}

	defer func() {
		if recovery := recover(); recovery != nil {
			err = fae.Errorf("panic (outside of job work): %v\n%s", recovery, string(debug.Stack()))
		}

		if err != nil {
			r.Minion.notify("job:fail", jobID, d.Kind)
		} else {
			r.Minion.notify("job:success", jobID, d.Kind)
		}
	}()

	return r.runJobAttempt(ctx, jobID, d, job)
}

func (r *Runner) loadJob(jobID string) (wrapped, *database.Model, error) {
	d := &database.Model{}
	err := r.Minion.db.Jobs.Find(jobID, d)
	if err != nil {
		return nil, d, fae.Wrap(err, fmt.Sprintf("finding job: %s", jobID))
	}

	w, ok := r.Minion.workers[d.Kind]
	if !ok {
		e := fae.Errorf("worker not found for kind: %s", d.Kind)
		d.Status = string(database.StatusCancelled)
		_ = r.Minion.db.Jobs.Save(d)
		return nil, d, e
	}

	job := w.factory.Create(d)
	err = job.Unmarshal()
	if err != nil {
		return nil, d, fae.Wrap(err, "unmarshaling job")
	}

	return job, d, nil
}

func (r *Runner) runJobAttempt(ctx context.Context, jobID string, d *database.Model, job wrapped) error {
	attempt := &database.Attempt{}
	attempt.Start()
	i := d.AddAttempt(attempt)
	err := r.Minion.db.Jobs.Save(d)
	if err != nil {
		return fae.Wrap(err, "updating job")
	}

	r.Minion.notify("job:start", jobID, d.Kind)
	err = r.runJobWork(ctx, job)
	e := fae.Wrap(err, "running job")
	attempt.Finish(e)
	r.Minion.notify("job:finish", jobID, d.Kind)

	d.UpdateAttempt(i, attempt)
	err = r.Minion.db.Jobs.Save(d)
	if err != nil {
		return fae.Wrap(err, "updating job")
	}

	return e
}

// runJobWork runs the job's Work method, this is separate
// to be able to handle deferred panics without affecting
// the job's attempt status
// we use named return so recover can set the error
func (r *Runner) runJobWork(ctx context.Context, job wrapped) (err error) {
	defer func() {
		if recovery := recover(); recovery != nil {
			err = fae.Errorf("panic: %v", recovery)
		}
	}()

	t := time.Duration(r.Minion.Config.Timeout) * time.Second
	if job.Timeout() > 0 {
		t = job.Timeout()
	}

	select {
	case <-time.After(t):
		err = fae.Errorf("timeout")
	case <-ctx.Done():
		return fae.Errorf("cancelled")
	default:
		err = job.Work(ctx)
	}

	return err
}

// WithTimeout runs a delegate function with a timeout,
//
// Example: Wait for a channel
//
//	if value, ok := WithTimeout(time.Second, func() error {return <- inbox}); ok {
//	    // returned
//	} else {
//	    // didn't return
//	}
//
// Example: To send to a channel
//
//	_, ok := WithTimeout(time.Second, func() error {outbox <- myValue; return nil})
//	if !ok {
//	    // didn't send
//	}
func withTimeout(timeout time.Duration, delegate func() error) (err error, ok bool) {
	ch := make(chan error, 1) // buffered
	go func() {
		defer func() {
			if recovery := recover(); recovery != nil {
				ch <- fae.Errorf("panic: %v", recovery)
			}
		}()
		ch <- delegate()
	}()
	select {
	case err = <-ch:
		return err, true
	case <-time.After(timeout):
		return nil, false
	}
}
