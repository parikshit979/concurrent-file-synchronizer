package tasks

import (
	"github.com/concurrent-file-synchronizer/eventbus"
	"github.com/concurrent-file-synchronizer/tasks/task"
)

type TaskManager struct {
	workerPool *WorkerPool
}

func NewTaskManager(workerCount int, eventBus *eventbus.EventBus) *TaskManager {
	return &TaskManager{
		workerPool: NewWorkerPool(workerCount, eventBus),
	}
}

func (tm *TaskManager) Start() {
	tm.workerPool.Start()
}

func (tm *TaskManager) Submit(task task.Task) {
	if task == nil {
		return
	}
	tm.workerPool.Submit(task)
}

func (tm *TaskManager) Stop() {
	tm.workerPool.Stop()
}
