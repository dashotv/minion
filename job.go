package minion

type ReportableType int

const (
	ReportableStart ReportableType = iota
	ReportableFinish
	ReportableError
	ReportableDuration
)

type Idable interface {
	GetID() string
}

type Nameable interface {
	GetName() string
}

type Runnable interface {
	Run(workerID int, minion *Minion) error
}

type Reportable interface {
	Report(event ReportableType, workerID int, minion *Minion) error
}

type Job struct {
	ID   string
	Name string
}

func (j *Job) GetID() string {
	return j.ID
}

func (j *Job) GetName() string {
	return j.Name
}

func (j *Job) Run(workerID int, minion *Minion) error {
	return nil
}

func (j *Job) Report(event string, workerID int, minion *Minion) error {
	return nil
}
