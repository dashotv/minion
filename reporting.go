package minion

type ReporterFunc func(ReportType, string, int) error
type ReportType int

const (
	ReportableStart ReportType = iota
	ReportableFinish
	ReportableError
	ReportableDuration
)

func (m *Minion) Report(t ReportType, name string, workerID int) error {
	if m.Reporter == nil {
		return nil
	}
	return m.Reporter(t, name, workerID)
}
