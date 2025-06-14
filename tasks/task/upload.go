package task

import (
	"context"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/concurrent-file-synchronizer/eventbus"
	"github.com/concurrent-file-synchronizer/types"
)

type UploadTask struct {
	taskOpt types.SyncTask
}

func NewUploadTask(taskOpt types.SyncTask) *UploadTask {
	return &UploadTask{
		taskOpt: taskOpt,
	}
}

func (t *UploadTask) Execute(ctx context.Context, eventBus *eventbus.EventBus) error {
	log.Infof("Starting upload task for file: %s", t.taskOpt.SourcePath)

	sourceFile, err := os.Open(t.taskOpt.SourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	if err := os.MkdirAll(filepath.Dir(t.taskOpt.DestPath), os.ModePerm); err != nil {
		return err
	}

	destFile, err := os.Create(t.taskOpt.DestPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	bytesDone, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	progress := types.ProgressTrackerEvent{
		EventUUID: t.taskOpt.EventUUID,
		FileName:  t.taskOpt.DestPath,
		Status:    types.ProgressStatusCompleted,
		BytesDone: bytesDone,
		TotalSize: t.taskOpt.FileInfo.Size(),
	}
	eventBus.Publish(eventbus.TopicProgressTrackerEvent, progress)
	return nil
}
