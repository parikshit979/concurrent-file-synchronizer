package services

import (
	log "github.com/sirupsen/logrus"

	"github.com/concurrent-file-synchronizer/eventbus"
	"github.com/concurrent-file-synchronizer/tasks"
	"github.com/concurrent-file-synchronizer/tasks/task"
	"github.com/concurrent-file-synchronizer/types"
)

type FileSyncer struct {
	eventbus    *eventbus.EventBus
	taskManager *tasks.TaskManager
}

func NewFileSyncer(workerCount int, eventBus *eventbus.EventBus) *FileSyncer {
	return &FileSyncer{
		eventbus:    eventBus,
		taskManager: tasks.NewTaskManager(workerCount, eventBus),
	}
}

func (fs *FileSyncer) Start() {
	fs.eventbus.Subscribe(eventbus.TopicDifferentiatorEvent, fs.handleEvent)
	fs.taskManager.Start()
}

func (fs *FileSyncer) handleEvent(event any) {
	switch e := event.(type) {
	case types.FileDifferentiatorEvent:
		fs.processEvent(&e)
	default:
		log.Error("Received unknown event type in FileDifferentiator:", e)
	}
}

func (fs *FileSyncer) processEvent(event *types.FileDifferentiatorEvent) {
	log.Infof("Processing FileDifferentiatorEvent: %s for file: %s", event.ActionType, event.SourceFilePath)

	syncTask := types.SyncTask{
		EventUUID:  event.EventUUID,
		Action:     event.ActionType,
		SourcePath: event.SourceFilePath,
		DestPath:   event.DestFilePath,
		FileInfo:   event.SourceFileInfo,
	}

	fs.taskManager.Submit(task.NewTask(event.ActionType, syncTask))
	log.Infof("SyncTask submitted for processing: %+v", syncTask)
}
