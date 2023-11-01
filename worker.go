package minion

type Worker struct {
	ID     int
	Queue  chan *Job
	Runner func(workerID int, j *Job)
}

func (w *Worker) Run() {
	for j := range w.Queue {
		w.Runner(w.ID, j)
	}
}
