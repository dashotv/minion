package minion

import (
	"fmt"
	"time"

	"github.com/dashotv/grimoire"
	"github.com/pkg/errors"
)

type JobDataStatus string

const (
	JobDataStatusPending  JobDataStatus = "pending"
	JobDataStatusRunning  JobDataStatus = "running"
	JobDataStatusFailed   JobDataStatus = "failed"
	JobDataStatusFinished JobDataStatus = "finished"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type JobData struct {
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`

	Kind string `bson:"kind,omitempty" json:"kind,omitempty"`
	Args string `bson:"args,omitempty" json:"args,omitempty"`

	Status   string            `bson:"status,omitempty" json:"status,omitempty"`
	Attempts []*JobDataAttempt `bson:"attempts,omitempty" json:"attempts,omitempty"`
}

func (d *JobData) AddAttempt(a *JobDataAttempt) int {
	d.Status = a.Status
	d.Attempts = append(d.Attempts, a)
	return len(d.Attempts) - 1
}

func (d *JobData) UpdateAttempt(i int, a *JobDataAttempt) {
	d.Status = a.Status
	d.Attempts[i] = a
}

type JobDataAttempt struct {
	StartedAt  time.Time `bson:"started_at,omitempty" json:"started_at,omitempty"`
	Duration   float64   `bson:"duration,omitempty" json:"duration,omitempty"`
	Status     string    `bson:"status,omitempty" json:"status,omitempty"`
	Error      string    `bson:"error,omitempty" json:"error,omitempty"`
	Stacktrace []string  `bson:"stacktrace,omitempty" json:"stacktrace,omitempty"`
}

func (a *JobDataAttempt) Start() {
	a.StartedAt = time.Now()
	a.Status = "running"
}

func (a *JobDataAttempt) Finish(err error) {
	a.Status = "finished"
	a.Duration = time.Since(a.StartedAt).Seconds()
	if err != nil {
		a.Status = "failed"
		a.Error = err.Error()

		err, ok := errors.Cause(err).(stackTracer)
		if !ok {
			return
		}

		st := err.StackTrace()
		if len(st) > 10 {
			st = st[:10]
		}
		for _, f := range st {
			a.Stacktrace = append(a.Stacktrace, fmt.Sprintf("%+v", f))
		}
	}
}

type Job[T Payload] struct {
	*JobData

	// Args are the arguments for the job.
	Args T
}
