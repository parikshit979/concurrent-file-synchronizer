package types

type FileEventType string

const (
	FileEventTypeCreate FileEventType = "CREATE"
	FileEventTypeModify FileEventType = "MODIFY"
	FileEventTypeDelete FileEventType = "DELETE"
)

type FileWatcherEvent struct {
	EventUUID       string
	EventType       FileEventType
	SourceFilePath  string
	DestFilePath    string
	IsRenamed       bool
	RenamedFileName string
}
