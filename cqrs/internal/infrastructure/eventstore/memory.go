package eventstore

import (
	"cqrs_exercise/internal/domain/event"
	"sync"
)

type MemoryEventStore struct {
	mu    sync.RWMutex
	store map[string][]event.Event
}

func New() *MemoryEventStore {
	return &MemoryEventStore{store: make(map[string][]event.Event)}
}

func (s *MemoryEventStore) Save(events []event.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, evt := range events {
		s.store[evt.AggregateID()] = append(s.store[evt.AggregateID()], evt)
	}
	return nil
}

func (s *MemoryEventStore) Load(aggregateID string) ([]event.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.store[aggregateID], nil
}
