package types

import "os"

type FileDetails struct {
	FilePath string
	FileInfo os.FileInfo
	Checksum string
}

type FileIndexerEvent struct {
	EventUUID  string
	EventType  FileEventType
	SourceFile *FileDetails
	DestFile   *FileDetails
}
