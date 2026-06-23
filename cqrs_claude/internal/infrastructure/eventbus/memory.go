package eventbus

import (
	"sync"

	"cqrs_claude/internal/domain/event"
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
	defer b.mu.RUnlock()
	for _, h := range b.handlers[evt.EventType()] {
		if err := h.Handle(evt); err != nil {
			return err
		}
	}
	return nil
}

func (b *MemoryEventBus) Subscribe(eventType string, handler event.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}
