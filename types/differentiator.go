package types

import "os"

type ActionType string

const (
	ActionTypeUpload   ActionType = "UPLOAD"
	ActionTypeDownload ActionType = "DOWNLOAD"
	ActionTypeDelete   ActionType = "DELETE"
)

type FileDifferentiatorEvent struct {
	EventUUID      string
	ActionType     ActionType
	SourceFilePath string
	DestFilePath   string
	SourceFileInfo os.FileInfo
	DestFileInfo   os.FileInfo
}
