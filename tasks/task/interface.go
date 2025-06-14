package task

import (
	"context"

	"github.com/concurrent-file-synchronizer/eventbus"
	"github.com/concurrent-file-synchronizer/types"
)

type Task interface {
	Execute(ctx context.Context, eventBus *eventbus.EventBus) error
}

func NewTask(taskType types.ActionType, taskOpt types.SyncTask) Task {
	switch taskType {
	case types.ActionTypeUpload:
		return NewUploadTask(taskOpt)
	case types.ActionTypeDelete:
		return NewDeleteTask(taskOpt)
	default:
		return nil
	}
}
