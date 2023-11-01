package minion

type Worker struct {
	ID     int
	Queue  chan Func
	Runner func(workerID int, f Func)
}

func (w *Worker) Run() {
	for j := range w.Queue {
		w.Runner(w.ID, j)
	}
}
