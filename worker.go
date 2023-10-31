package minion

import "go.uber.org/zap"

type Worker struct {
	ID     int
	Log    *zap.SugaredLogger
	Queue  chan Runnable
	Runner func(workerID int, job Runnable)
}

func (w *Worker) Run() {
	for j := range w.Queue {
		w.Runner(w.ID, j)
	}
}
