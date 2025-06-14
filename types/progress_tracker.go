package types

type ProgressStatus string

const (
	ProgressStatusInProgress ProgressStatus = "INPROGRESS"
	ProgressStatusCompleted  ProgressStatus = "COMPLETED"
	ProgressStatusFailed     ProgressStatus = "FAILED"
	ProgressStatusCancelled  ProgressStatus = "CANCELLED"
)

type ProgressTrackerEvent struct {
	EventUUID string
	FileName  string
	TotalSize int64
	BytesDone int64
	Status    ProgressStatus
	Error     string
}
