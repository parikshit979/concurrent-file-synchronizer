package types

import "os"

type SyncTask struct {
	EventUUID  string
	Action     ActionType
	SourcePath string
	DestPath   string
	FileInfo   os.FileInfo
}
