package minion

/* this implementation is heavily inspired by
   https://github.com/riverqueue/river

   I had built something similar, then I saw their announcment
   and it filled in the gaps for me.
*/

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dashotv/minion/database"
)

type Job[T Payload] struct {
	*database.Model

	// Args are the arguments for the job.
	Args T
}

type Payload interface {
	Kind() string
}

// WorkerDefaults is a helper struct that implements the Worker interface
// with default values.  Embed this struct in your worker to get the
// default behavior.
type WorkerDefaults[T Payload] struct{}

// Timeout returns the timeout for the job, override this method to
// set a timeout specific to this job, otherwise the default timeout
// will be used.
func (w *WorkerDefaults[T]) Timeout(*Job[T]) time.Duration { return 0 }

// Worker is the interface that must be implemented by all workers.
type Worker[T Payload] interface {
	Timeout(*Job[T]) time.Duration
	Work(ctx context.Context, job *Job[T]) error
}

type wrapped interface {
	Unmarshal() error
	Timeout() time.Duration
	Work(ctx context.Context) error
}

type wrappedWorker[T Payload] struct {
	job    *Job[T]
	data   *database.Model
	worker Worker[T]
}

func (w *wrappedWorker[T]) Work(ctx context.Context) error {
	return w.worker.Work(ctx, w.job)
}

func (w *wrappedWorker[T]) Timeout() time.Duration {
	return w.worker.Timeout(w.job)
}

func (w *wrappedWorker[T]) Unmarshal() error {
	w.job = &Job[T]{
		Model: w.data,
	}
	return json.Unmarshal([]byte(w.data.Args), &w.job.Args)
}

type factory interface {
	Create(data *database.Model) wrapped
}

type workerFactory[T Payload] struct {
	worker Worker[T]
}

func (f *workerFactory[T]) Create(data *database.Model) wrapped {
	return &wrappedWorker[T]{data: data, worker: f.worker}
}
