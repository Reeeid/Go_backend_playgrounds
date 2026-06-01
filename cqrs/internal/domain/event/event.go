package event

import "time"

type Event interface {
	AggregateID() string
	EventType() string
	OccurredAt() time.Time
}

type EventHandler interface {
	Handle(event Event) error
}

type EventStore interface {
	Save(events []Event) error
	Load(aggregateID string) ([]Event, error)
}

type EventPublisher interface {
	Publish(event Event) error
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handler EventHandler) error
}

type EventBus struct {
	handlers map[string][]EventHandler
}
