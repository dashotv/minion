package minion

import (
	"encoding/json"

	"github.com/pkg/errors"

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

func (m *Minion) enqueueTo(queue string, in Payload) error {
	if in == nil {
		return errors.New("payload is nil")
	}

	args, err := json.Marshal(in)
	if err != nil {
		return errors.Wrap(err, "marshaling job args")
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
		return errors.Wrap(err, "creating job")
	}

	m.notify("job:created", data.ID.Hex(), data.Kind)
	return nil
}
