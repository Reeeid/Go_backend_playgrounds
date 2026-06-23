// Package projection はイベントを購読してクエリ側の読み取りモデルを構築する。
package projection

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"rabbitmq_cqrs/internal/domain"
)

// UserReadModel はクエリ側の読み取り専用のユーザーモデル。
type UserReadModel struct {
	ID    string
	Name  string
	Email string
}

// UserProjection はイベントを受信して UserReadModel を更新するプロジェクション。
type UserProjection struct {
	mu    sync.RWMutex
	store map[string]*UserReadModel
}

func NewUserProjection() *UserProjection {
	return &UserProjection{store: make(map[string]*UserReadModel)}
}

// Handle はイベントバスから渡されたメッセージを処理する。
func (p *UserProjection) Handle(eventType string, body []byte) error {
	switch eventType {
	case "UserCreatedEvent":
		return p.onUserCreated(body)
	case "UserUpdatedEvent":
		return p.onUserUpdated(body)
	default:
		return nil
	}
}

func (p *UserProjection) onUserCreated(body []byte) error {
	var e domain.UserCreatedEvent
	if err := json.Unmarshal(body, &e); err != nil {
		return fmt.Errorf("unmarshal UserCreatedEvent: %w", err)
	}
	p.mu.Lock()
	p.store[e.ID] = &UserReadModel{ID: e.ID, Name: e.Name, Email: e.Email}
	p.mu.Unlock()
	log.Printf("[projection] UserCreated -> id=%s name=%s", e.ID, e.Name)
	return nil
}

func (p *UserProjection) onUserUpdated(body []byte) error {
	var e domain.UserUpdatedEvent
	if err := json.Unmarshal(body, &e); err != nil {
		return fmt.Errorf("unmarshal UserUpdatedEvent: %w", err)
	}
	p.mu.Lock()
	if m, ok := p.store[e.ID]; ok {
		m.Name = e.Name
		m.Email = e.Email
	}
	p.mu.Unlock()
	log.Printf("[projection] UserUpdated -> id=%s name=%s", e.ID, e.Name)
	return nil
}

// FindByID は ID でユーザーを取得する（クエリ側 API から呼ぶ）。
func (p *UserProjection) FindByID(id string) (*UserReadModel, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	m, ok := p.store[id]
	return m, ok
}

// All は全ユーザーを返す。
func (p *UserProjection) All() []*UserReadModel {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]*UserReadModel, 0, len(p.store))
	for _, m := range p.store {
		out = append(out, m)
	}
	return out
}
