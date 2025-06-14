package services

import (
	log "github.com/sirupsen/logrus"

	"github.com/concurrent-file-synchronizer/eventbus"
	"github.com/concurrent-file-synchronizer/types"
)

type ProgressTracker struct {
	eventbus *eventbus.EventBus
}

func NewProgressTracker(eventbus *eventbus.EventBus) *ProgressTracker {
	return &ProgressTracker{
		eventbus: eventbus,
	}
}

func (pt *ProgressTracker) Start() {
	pt.eventbus.Subscribe(eventbus.TopicProgressTrackerEvent, pt.handleEvent)
}

func (pt *ProgressTracker) handleEvent(event any) {
	switch e := event.(type) {
	case types.ProgressTrackerEvent:
		pt.processTrackerEvent(&e)
	default:
		log.Error("Received unknown event type in ProgressTrackerEvent:", e)
	}
}

func (pt *ProgressTracker) processTrackerEvent(event *types.ProgressTrackerEvent) {
	log.Infof("Processing ProgressTrackerEvent: %+v for file: %s", event, event.FileName)
}
