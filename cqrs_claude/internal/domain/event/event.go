package event

import "time"

// Event はすべてのドメインイベントのマーカーインターフェイス
type Event interface {
	isEvent()
	AggregateID() string
	EventType() string
	OccurredAt() time.Time
}

// BaseEvent は全イベント共通フィールド
type BaseEvent struct {
	aggregateID string
	eventType   string
	occurredAt  time.Time
}

func NewBaseEvent(aggregateID, eventType string) BaseEvent {
	return BaseEvent{
		aggregateID: aggregateID,
		eventType:   eventType,
		occurredAt:  time.Now(),
	}
}

func (e BaseEvent) isEvent()              {}
func (e BaseEvent) AggregateID() string   { return e.aggregateID }
func (e BaseEvent) EventType() string     { return e.eventType }
func (e BaseEvent) OccurredAt() time.Time { return e.occurredAt }

// EventStore はイベントの永続化・再生を担当
type EventStore interface {
	Save(events []Event) error
	Load(aggregateID string) ([]Event, error)
}

// EventHandler は単一イベントの処理を担当（Projection等）
type EventHandler interface {
	Handle(evt Event) error
}

// EventPublisher はイベントを購読者に配信
type EventPublisher interface {
	Publish(evt Event) error
	Subscribe(eventType string, handler EventHandler)
}
