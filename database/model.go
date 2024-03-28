package database

import (
	"time"

	"github.com/dashotv/fae"
	"github.com/dashotv/grimoire"
)

type stackTracer interface {
	StackTrace() []string
}

type Model struct {
	grimoire.Document `bson:",inline"` // includes default model settings
	//ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	//CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	//UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`

	Client string `bson:"client" json:"client" grimoire:"index"`
	Kind   string `bson:"kind" json:"kind" grimoire:"index"`
	Args   string `bson:"args,omitempty" json:"args,omitempty"`
	Queue  string `bson:"queue,omitempty" json:"queue,omitempty"`

	Status   string     `bson:"status,omitempty" json:"status,omitempty" grimoire:"index"`
	Attempts []*Attempt `bson:"attempts,omitempty" json:"attempts,omitempty"`
}

func (d *Model) AddAttempt(a *Attempt) int {
	d.Status = a.Status
	d.Attempts = append(d.Attempts, a)
	return len(d.Attempts) - 1
}

func (d *Model) UpdateAttempt(i int, a *Attempt) {
	d.Status = a.Status
	d.Attempts[i] = a
}

type Attempt struct {
	StartedAt  time.Time `bson:"started_at,omitempty" json:"started_at,omitempty"`
	Duration   float64   `bson:"duration,omitempty" json:"duration,omitempty"`
	Status     string    `bson:"status,omitempty" json:"status,omitempty"`
	Error      string    `bson:"error,omitempty" json:"error,omitempty"`
	Stacktrace []string  `bson:"stacktrace,omitempty" json:"stacktrace,omitempty"`
}

func (a *Attempt) Start() {
	a.StartedAt = time.Now()
	a.Status = string(StatusRunning)
}

func (a *Attempt) Finish(err error) {
	a.Status = string(StatusFinished)
	a.Duration = time.Since(a.StartedAt).Seconds()
	if err != nil {
		a.Status = string(StatusFailed)
		a.Error = fae.Cause(err).Error()

		st := fae.StackTrace(err)
		st = st[1:] // remove the first entry, it's the error
		if len(st) > 10 {
			st = st[:10]
		}
		for _, f := range st {
			a.Stacktrace = append(a.Stacktrace, f)
		}
	}
}
