package minion

type Reportable interface {
	Report(t ReportType, name string, workerID int) error
}
type ReporterFunc func(ReportType, string, int) error
type ReportType int

const (
	ReportableUnknown ReportType = iota
	ReportableStart
	ReportableFinish
	ReportableError
	ReportableDuration
)

func (m *Minion) Report(t ReportType, name string, workerID int) error {
	return m.Reporter.Report(t, name, workerID)
}

type DefaultReporter struct{}

func (r *DefaultReporter) Report(t ReportType, name string, workerID int) error {
	return nil
}
