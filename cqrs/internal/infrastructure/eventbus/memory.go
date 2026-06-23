package eventbus

import (
	"fmt"
	"sync"

	"cqrs_exercise/internal/domain/event"
)

type MemoryEventBus struct {
	mu       sync.RWMutex
	handlers map[string][]event.EventHandler
}

func New() *MemoryEventBus {
	return &MemoryEventBus{handlers: make(map[string][]event.EventHandler)}
}

func (b *MemoryEventBus) Publish(evt event.Event) error {
	b.mu.RLock()
	handlers := make([]event.EventHandler, len(b.handlers[evt.EventType()]))
	copy(handlers, b.handlers[evt.EventType()])
	b.mu.RUnlock()

	for _, h := range handlers {
		if err := h.Handle(evt); err != nil {
			return fmt.Errorf("eventbus: handler failed for %s: %w", evt.EventType(), err)
		}
	}
	return nil
}

func (b *MemoryEventBus) Subscribe(eventType string, handler event.EventHandler) error {
	if handler == nil {
		return fmt.Errorf("eventbus: handler must not be nil")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

func (b *MemoryEventBus) Unsubscribe(eventType string, handler event.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	hs := b.handlers[eventType]
	for i, h := range hs {
		if h == handler {
			b.handlers[eventType] = append(hs[:i], hs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("eventbus: handler not found for event type %s", eventType)
}
