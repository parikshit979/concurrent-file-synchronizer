package services

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/concurrent-file-synchronizer/eventbus"
	"github.com/concurrent-file-synchronizer/types"
)

type FileIndexer struct {
	eventbus *eventbus.EventBus
}

func NewFileIndexer(eventBus *eventbus.EventBus) *FileIndexer {
	return &FileIndexer{
		eventbus: eventBus,
	}
}

func (fi *FileIndexer) Start() {
	fi.eventbus.Subscribe(eventbus.TopicWatcherEvent, fi.handleEvent)
}

func (fi *FileIndexer) handleEvent(event any) {
	switch e := event.(type) {
	case types.FileWatcherEvent:
		fi.processEvent(&e)
	default:
		log.Error("Received unknown event type in FileWatcherEvent:", e)
	}
}

func (fi *FileIndexer) processEvent(event *types.FileWatcherEvent) {
	log.Infof("Processing FileWatcherEvent: %s for file: %s", event.EventType, event.SourceFilePath)

	fileIndexerEvent := types.FileIndexerEvent{
		EventUUID: event.EventUUID,
		EventType: event.EventType,
		SourceFile: &types.FileDetails{
			FilePath: event.SourceFilePath,
		},
		DestFile: &types.FileDetails{
			FilePath: event.DestFilePath,
		},
	}

	switch event.EventType {
	case types.FileEventTypeCreate, types.FileEventTypeModify:
		fileInfo, err := os.Stat(event.SourceFilePath)
		if err != nil {
			log.Errorf("Failed to get file info for %s: %v", event.SourceFilePath, err)
			return
		}
		fileIndexerEvent.SourceFile.FileInfo = fileInfo

		fileInfo, err = os.Stat(event.DestFilePath)
		if err != nil {
			log.Infof("Failed to get file info for %s: %v", event.DestFilePath, err)
		}
		fileIndexerEvent.DestFile.FileInfo = fileInfo

		fileIndexerEvent.SourceFile.Checksum = "dummy-checksum"
		fileIndexerEvent.DestFile.Checksum = "dummy-checksum"
	}

	log.Infof("FileIndexerEvent created: %+v", fileIndexerEvent)
	fi.eventbus.Publish(eventbus.TopicIndexerEvent, fileIndexerEvent)
}
