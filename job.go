package minion

type Job struct {
	Name    string
	Func    Func
	Payload interface{}
}
