package minion

import (
	"encoding/json"

	"github.com/dashotv/fae"
	"github.com/dashotv/minion/database"
)

func (m *Minion) Enqueue(in Payload) error {
	queue := "default"

	reg := m.workers[in.Kind()]
	if reg.queue != "" {
		queue = reg.queue
	}

	return m.enqueueTo(queue, in)
}

func (m *Minion) Requeue(jobID string) error {
	job := &database.Model{}
	err := m.db.Jobs.Find(jobID, job)
	if err != nil {
		return fae.Wrap(err, "finding job")
	}

	job.Status = string(database.StatusPending)
	err = m.db.Jobs.Save(job)
	if err != nil {
		return fae.Wrap(err, "updating job")
	}

	m.notify("job:queued", job.ID.Hex(), job.Kind)
	return nil
}

func (m *Minion) enqueueTo(queue string, in Payload) error {
	if in == nil {
		return fae.New("payload is nil")
	}

	args, err := json.Marshal(in)
	if err != nil {
		return fae.Wrap(err, "marshaling job args")
	}

	data := &database.Model{
		Client: m.Client,
		Args:   string(args),
		Kind:   in.Kind(),
		Status: string(database.StatusPending),
		Queue:  queue,
	}

	err = m.db.Jobs.Save(data)
	if err != nil {
		return fae.Wrap(err, "creating job")
	}

	m.notify("job:created", data.ID.Hex(), data.Kind)
	return nil
}
