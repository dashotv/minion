package minion

type Loggable interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type DefaultLogger struct{}

func (d *DefaultLogger) Infof(format string, args ...interface{})  {}
func (d *DefaultLogger) Errorf(format string, args ...interface{}) {}
