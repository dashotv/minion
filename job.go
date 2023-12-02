package minion

/* this implementation is heavily inspired by
   https://github.com/riverqueue/river

   I had built something similar, then I saw their announcment
   and it filled in the gaps for me.
*/

import (
	"context"
	"encoding/json"
)

type workerInfo struct {
	jobArgs Payload
	factory factory
}

type Payload interface {
	Kind() string
}

type WorkerDefaults[T Payload] struct{}
type Worker[T Payload] interface {
	Work(ctx context.Context, job *Job[T]) error
}

type wrapped interface {
	Unmarshal() error
	Work(ctx context.Context) error
}

type factory interface {
	Create(data *JobData) wrapped
}

type workerFactory[T Payload] struct {
	worker Worker[T]
}

func (f *workerFactory[T]) Create(data *JobData) wrapped {
	return &wrappedWorker[T]{data: data, worker: f.worker}
}

type wrappedWorker[T Payload] struct {
	job    *Job[T]
	data   *JobData
	worker Worker[T]
}

func (w *wrappedWorker[T]) Work(ctx context.Context) error {
	return w.worker.Work(ctx, w.job)
}

func (w *wrappedWorker[T]) Unmarshal() error {
	w.job = &Job[T]{
		JobData: w.data,
	}
	return json.Unmarshal([]byte(w.data.Args), &w.job.Args)
}
