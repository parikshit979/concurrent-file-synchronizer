package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/concurrent-file-synchronizer/eventbus"
	"github.com/concurrent-file-synchronizer/types"
)

type FileWatcher struct {
	watcher   *fsnotify.Watcher
	localDir  string
	remoteDir string
	eventBus  *eventbus.EventBus
}

func NewFileWatcher(localDir, remoteDir string, eventBus *eventbus.EventBus) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		watcher:   watcher,
		localDir:  localDir,
		remoteDir: remoteDir,
		eventBus:  eventBus,
	}

	// Start the watcher for the local directory
	if err := fw.watcher.Add(localDir); err != nil {
		return nil, err
	}

	return fw, nil
}

func (fw *FileWatcher) Start() error {
	log.Infof("Setting up a watcher for directory: %s", fw.localDir)

	if _, err := os.Stat(fw.localDir); os.IsNotExist(err) {
		log.Infof("Directory '%s' does not exist, creating it.", fw.localDir)
		if err := os.Mkdir(fw.localDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	go func() {
		for {
			select {
			case event, ok := <-fw.watcher.Events:
				if !ok {
					return
				}
				fw.handleEvent(event)
			case err, ok := <-fw.watcher.Errors:
				if !ok {
					return
				}
				log.Infof("Watcher Error: %v", err)
			}
		}
	}()

	return nil
}

func (fw *FileWatcher) handleEvent(event fsnotify.Event) {
	log.Infof("Received watcher event: %+v", event)

	fileEvent := types.FileWatcherEvent{
		EventUUID:      uuid.New().String(),
		SourceFilePath: event.Name,
		// TODO: Remove fw.localDir from the path
		DestFilePath: fw.remoteDir + "/" + filepath.Base(event.Name),
	}
	if event.Has(fsnotify.Create) {
		fileEvent.EventType = types.FileEventTypeCreate
	} else if event.Has(fsnotify.Write) {
		fileEvent.EventType = types.FileEventTypeModify
	} else if event.Has(fsnotify.Remove) {
		fileEvent.EventType = types.FileEventTypeDelete
	}

	log.Infof("FileWatcherEvent created: %+v", fileEvent)
	fw.eventBus.Publish(eventbus.TopicWatcherEvent, fileEvent)
}

func (fw *FileWatcher) Stop() {
	if err := fw.watcher.Close(); err != nil {
		log.Infof("Error closing watcher: %v", err)
	}

	log.Infof("FileWatcher closed.")
}
