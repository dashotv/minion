package minion

import "context"

func (m *Minion) Subscribe(f func(*Notification)) {
	m.subs = append(m.subs, f) // TODO: should handle panices here
}

type Notification struct {
	Event string
	JobID string
	Kind  string
}

func (m *Minion) debug(n *Notification) {
	m.Log.Debugf("event=%s job=%s kind=%s", n.Event, n.JobID, n.Kind)
}

func (m *Minion) notify(event string, jobID string, kind string) {
	if !m.listening {
		// m.Log.Warnf("no listeners for notification: %s", event)
		return
	}
	if channelBufferFull(m.notifications) {
		m.Log.Debugf("notification buffer full: %d", cap(m.notifications))
	}
	m.notifications <- &Notification{event, jobID, kind}
}

func (m *Minion) listen(ctx context.Context) {
	m.listening = true
	for {
		select {
		case n := <-m.notifications:
			for _, s := range m.subs {
				s(n)
			}
		case <-ctx.Done():
			m.listening = false
			return
		}
	}
}
