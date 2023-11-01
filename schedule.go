package minion

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

// Schedule adds a job to the cron scheduler.
func (m *Minion) Schedule(schedule, name string) (cron.EntryID, error) {
	f, ok := m.Jobs[name]
	if !ok {
		return 0, fmt.Errorf("job not found: %s", name)
	}
	return m.Cron.AddFunc(schedule, func() {
		m.Queue <- &Job{name, f}
	})
}

// Remove removes a job from the cron scheduler.
func (m *Minion) Remove(id cron.EntryID) {
	m.Cron.Remove(id)
}

// Enqueue adds a job to the queue.
func (m *Minion) Enqueue(name string) error {
	f, ok := m.Jobs[name]
	if !ok {
		return fmt.Errorf("job not found: %s", name)
	}
	m.Queue <- &Job{name, f}
	return nil
}

func (m *Minion) Register(name string, f Func) {
	m.Jobs[name] = f
}
