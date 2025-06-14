package services

import (
	log "github.com/sirupsen/logrus"

	"github.com/concurrent-file-synchronizer/eventbus"
)

const (
	worckerCount = 5
)

type SynchronizerService struct {
	eventBus        *eventbus.EventBus
	watcher         *FileWatcher
	indexer         *FileIndexer
	differentiator  *FileDifferentiator
	syncer          *FileSyncer
	progressTracker *ProgressTracker
}

func NewSynchronizerService(localDir, remoteDir string) *SynchronizerService {
	eventBus := eventbus.NewEventBus()
	log.Infof("Creating event bus for Synchronizer Service...")
	watcher, err := NewFileWatcher(localDir, remoteDir, eventBus)
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}
	log.Infof("File watcher created for local directory: %s and remote directory: %s", localDir, remoteDir)
	indexer := NewFileIndexer(eventBus)
	differentiator := NewFileDifferentiator(eventBus)
	syncer := NewFileSyncer(worckerCount, eventBus)
	progressTracker := NewProgressTracker(eventBus)

	return &SynchronizerService{
		eventBus:        eventBus,
		watcher:         watcher,
		indexer:         indexer,
		differentiator:  differentiator,
		syncer:          syncer,
		progressTracker: progressTracker,
	}
}

func (s *SynchronizerService) Start() {
	log.Info("Starting Synchronizer Service...")
	if err := s.watcher.Start(); err != nil {
		log.Fatalf("Failed to start file watcher: %v", err)
	}
	s.indexer.Start()
	s.differentiator.Start()
	s.syncer.Start()
	s.progressTracker.Start()
	log.Info("Synchronizer Service started successfully.")
}

func (s *SynchronizerService) Stop() {
	log.Info("Stopping Synchronizer Service...")
	if s.eventBus != nil {
		s.eventBus.Close()
	}
	if s.watcher != nil {
		s.watcher.Stop()
	}
	log.Info("Synchronizer Service stopped.")
}
