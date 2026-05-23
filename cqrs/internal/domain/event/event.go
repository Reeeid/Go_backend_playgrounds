package event

import "time"

type Event interface {
	isEvent()
}

type BaseEvent struct {
	AggregateID string
	eventType   string
	occurredAt  time.Time
	version     int
}

func (e BaseEvent) isEvent() {}

func NewBaseEvent(aggregateID string, eventType string) BaseEvent {
	return BaseEvent{
		AggregateID: aggregateID,
		eventType:   eventType,
		occurredAt:  time.Now(),
		version:     1,
	}
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
