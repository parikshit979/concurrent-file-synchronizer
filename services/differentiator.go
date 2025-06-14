package services

import (
	log "github.com/sirupsen/logrus"

	"github.com/concurrent-file-synchronizer/eventbus"
	"github.com/concurrent-file-synchronizer/types"
)

type FileDifferentiator struct {
	eventbus *eventbus.EventBus
}

func NewFileDifferentiator(eventBus *eventbus.EventBus) *FileDifferentiator {
	return &FileDifferentiator{
		eventbus: eventBus,
	}
}

func (fd *FileDifferentiator) Start() {
	fd.eventbus.Subscribe(eventbus.TopicIndexerEvent, fd.handleEvent)
}

func (fd *FileDifferentiator) handleEvent(event any) {
	switch e := event.(type) {
	case types.FileIndexerEvent:
		fd.processEvent(&e)
	default:
		log.Error("Received unknown event type in FileIndexerEvent:", e)
	}
}

func (fd *FileDifferentiator) processEvent(event *types.FileIndexerEvent) {
	log.Infof("Processing FileindexerEvent: %s for file: %s", event.EventType, event.SourceFile.FilePath)

	fileDifferentiatorEvent := types.FileDifferentiatorEvent{
		EventUUID:      event.EventUUID,
		SourceFilePath: event.SourceFile.FilePath,
		DestFilePath:   event.DestFile.FilePath,
		SourceFileInfo: event.SourceFile.FileInfo,
		DestFileInfo:   event.DestFile.FileInfo,
	}

	switch event.EventType {
	case types.FileEventTypeCreate, types.FileEventTypeModify:
		if event.SourceFile.FileInfo != event.DestFile.FileInfo || event.SourceFile.Checksum != event.DestFile.Checksum {
			fileDifferentiatorEvent.ActionType = types.ActionTypeUpload
		}
	case types.FileEventTypeDelete:
		fileDifferentiatorEvent.ActionType = types.ActionTypeDelete
	}

	log.Infof("FileDifferentiatorEvent created: %+v", fileDifferentiatorEvent)
	fd.eventbus.Publish(eventbus.TopicDifferentiatorEvent, fileDifferentiatorEvent)
}
