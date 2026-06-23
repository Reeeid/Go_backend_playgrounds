package eventstore

import (
	"encoding/json"
	"fmt"
	"time"

	"cqrs_exercise/internal/domain/event"

	"gorm.io/gorm"
)

type eventRecord struct {
	ID               uint      `gorm:"primaryKey;autoIncrement"`
	AggregateID      string    `gorm:"column:aggregate_id;uniqueIndex:idx_agg_ver"`
	AggregateVersion int       `gorm:"column:aggregate_version;uniqueIndex:idx_agg_ver"`
	EventType        string    `gorm:"column:event_type"`
	Payload          []byte    `gorm:"column:payload"`
	OccurredAt       time.Time `gorm:"column:occurred_at"`
}

type PostgresEventStore struct {
	db          *gorm.DB
	deserialize map[string]func([]byte) (event.Event, error)
}

func New(db *gorm.DB, deserialize map[string]func([]byte) (event.Event, error)) *PostgresEventStore {
	return &PostgresEventStore{db: db, deserialize: deserialize}
}

func (s *PostgresEventStore) Save(events []event.Event) error {
	records := make([]eventRecord, 0, len(events))
	for _, e := range events {
		payload, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("serialize event %s: %w", e.EventType(), err)
		}
		records = append(records, eventRecord{
			AggregateID:      e.AggregateID(),
			AggregateVersion: e.AggregateVersion(),
			EventType:        e.EventType(),
			Payload:          payload,
			OccurredAt:       e.OccurredAt(),
		})
	}
	return s.db.Create(&records).Error
}

func (s *PostgresEventStore) Load(aggregateID string) ([]event.Event, error) {
	var records []eventRecord
	err := s.db.
		Where("aggregate_id = ?", aggregateID).
		Order("aggregate_version ASC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}

	events := make([]event.Event, 0, len(records))
	for _, r := range records {
		factory, ok := s.deserialize[r.EventType]
		if !ok {
			return nil, fmt.Errorf("unknown event type: %s", r.EventType)
		}
		e, err := factory(r.Payload)
		if err != nil {
			return nil, fmt.Errorf("deserialize %s: %w", r.EventType, err)
		}
		events = append(events, e)
	}
	return events, nil
}
