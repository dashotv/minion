package minion

import (
	"github.com/robfig/cron/v3"

	"github.com/dashotv/minion/database"
)

// Schedule adds (and Registers) a job to the cron scheduler.
func (m *Minion) Schedule(schedule string, in Payload) (cron.EntryID, error) {
	return m.cron.AddFunc(schedule, func() {
		m.notify("job:scheduled", "-", in.Kind())
		m.enqueueTo("schedule", in)
	})
}

// ScheduleFunc adds a function to the cron scheduler which is only stored
// in the database if it fails.
func (m *Minion) ScheduleFunc(schedule, name string, f func() error) (cron.EntryID, error) {
	return m.cron.AddFunc(schedule, func() {
		err := f()
		if err != nil {
			m.Log.Error(err)

			data := &database.Model{
				Args:   "{}",
				Kind:   name,
				Status: string(database.StatusFailed),
				Queue:  "schedule_func",
			}

			err = m.db.Jobs.Save(data)
			if err != nil {
				m.Log.Errorf("error saving job: %s", err)
			}
		}
	})
}

// Remove removes a job from the cron scheduler.
func (m *Minion) Remove(id cron.EntryID) {
	m.cron.Remove(id)
}
