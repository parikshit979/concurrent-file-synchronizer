package eventbus

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

// Topic defines a string identifier for an event topic.
type Topic string

const (
	TopicWatcherEvent         Topic = "WatcherEvent"
	TopicIndexerEvent         Topic = "IndexerEvent"
	TopicDifferentiatorEvent  Topic = "DifferentiatorEvent"
	TopicSyncerEvent          Topic = "SyncerEvent"
	TopicProgressTrackerEvent Topic = "ProgressTrackerEvent"
	TopicTaskErrorEvent       Topic = "TaskErrorEvent"
)

// EventBus provides a simple publish-subscribe mechanism using Go channels.
type EventBus struct {
	subscribers map[Topic][]chan any
	mu          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[Topic][]chan any),
	}
}

// Subscribe allows a component to listen for events on a given topic.
// It returns a channel where events will be received.
func (eb *EventBus) Subscribe(topic Topic, handler func(data any)) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// Create a new channel for this specific subscription
	ch := make(chan any, 100) // Buffered channel to prevent blocking publisher

	eb.subscribers[topic] = append(eb.subscribers[topic], ch)
	log.Infof("EventBus: New subscription to topic: %s", topic)

	// Start a goroutine to continuously read from this channel and call the handler
	go func() {
		for data := range ch {
			handler(data)
		}
	}()
}

// Publish sends an event to all subscribers of a specific topic.
func (eb *EventBus) Publish(topic Topic, data any) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if channels, ok := eb.subscribers[topic]; ok {
		for _, ch := range channels {
			select {
			case ch <- data:
				// Event sent successfully
			default:
				log.Infof("WARNING: EventBus: Subscriber channel for topic %s is full. Dropping event.", topic)
			}
		}
	}
}

// Close closes all subscriber channels.
func (eb *EventBus) Close() {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	for _, channels := range eb.subscribers {
		for _, ch := range channels {
			close(ch)
		}
	}
	eb.subscribers = make(map[Topic][]chan any)
	log.Infoln("EventBus: All subscriber channels closed.")
}
