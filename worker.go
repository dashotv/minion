package minion

type Worker struct {
	ID     int
	Queue  chan Runnable
	Runner func(workerID int, job Runnable)
}

func (w *Worker) Run() {
	for j := range w.Queue {
		w.Runner(w.ID, j)
	}
}
