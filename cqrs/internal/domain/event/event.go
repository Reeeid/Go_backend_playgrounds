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

func (e BaseEvent) isEvent()              {}
func (e BaseEvent) OccurredAt() time.Time { return e.occurredAt }
func (e BaseEvent) EventType() string     { return e.eventType }
func (e BaseEvent) Version() int          { return e.version }

func NewBaseEvent(aggregateID string, eventType string, version int) BaseEvent {
	return BaseEvent{
		AggregateID: aggregateID,
		eventType:   eventType,
		occurredAt:  time.Now(),
		version:     version,
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
