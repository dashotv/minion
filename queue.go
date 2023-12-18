package minion

type Queue struct {
	Name        string
	Concurrency int
	BufferSize  int
	Interval    int
	channel     chan string
}

func (q *Queue) Full() bool {
	return channelBufferFull(q.channel)
}

func (q *Queue) Remaining() int {
	return channelBufferRemaining(q.channel)
}

// Queue adds a new queue to Minion.
func (m *Minion) Queue(name string, concurrency, buffersize, interval int) {
	if concurrency == 0 {
		concurrency = m.Config.Concurrency
	}
	if buffersize == 0 {
		buffersize = m.Config.BufferSize
	}
	if interval == 0 {
		interval = m.Config.PollingInterval
	}

	m.queues[name] = &Queue{name, concurrency, buffersize, interval, make(chan string, buffersize)}
}
