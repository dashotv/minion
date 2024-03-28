package minion

import "github.com/dashotv/fae"

type registration struct {
	args        Payload
	factory     factory
	queue       string
	concurrency int
	bufferSize  int
}

func Register[T Payload](m *Minion, worker Worker[T]) error {
	return RegisterWithQueue(m, worker, "default")
}

func RegisterWithQueue[T Payload](m *Minion, worker Worker[T], queue string) error {
	var args T

	kind := args.Kind()
	if _, ok := m.workers[kind]; ok {
		return fae.Errorf("worker already registered for kind: %s", kind)
	}

	m.workers[kind] = registration{
		args:    args,
		factory: &workerFactory[T]{worker: worker},
		queue:   queue,
	}

	return nil
}
