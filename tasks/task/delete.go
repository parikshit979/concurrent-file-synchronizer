package task

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/concurrent-file-synchronizer/eventbus"
	"github.com/concurrent-file-synchronizer/types"
)

type DeleteTask struct {
	taskOpt types.SyncTask
}

func NewDeleteTask(taskOpt types.SyncTask) *DeleteTask {
	return &DeleteTask{
		taskOpt: taskOpt,
	}
}

func (t *DeleteTask) Execute(ctx context.Context, eventBus *eventbus.EventBus) error {
	log.Infof("Starting delete task for file: %s", t.taskOpt.DestPath)

	if err := os.RemoveAll(t.taskOpt.DestPath); err != nil {
		return fmt.Errorf("failed to delete path %s: %v", t.taskOpt.DestPath, err)
	}

	progress := types.ProgressTrackerEvent{
		EventUUID: t.taskOpt.EventUUID,
		FileName:  t.taskOpt.DestPath,
		Status:    types.ProgressStatusCompleted,
	}
	eventBus.Publish(eventbus.TopicProgressTrackerEvent, progress)
	return nil
}
