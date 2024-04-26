package database

type Status string

const (
	StatusPending   Status = "pending"
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusFailed    Status = "failed"
	StatusFinished  Status = "finished"
	StatusCancelled Status = "cancelled"
	StatusArchived  Status = "archived"
)
