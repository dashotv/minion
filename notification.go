package minion

func (m *Minion) Subscribe(f func(*Notification)) {
	m.subs = append(m.subs, f)
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

func (m *Minion) listen() {
	m.listening = true
	for n := range m.notifications {
		for _, s := range m.subs {
			s(n)
		}
	}
	m.listening = false
}
