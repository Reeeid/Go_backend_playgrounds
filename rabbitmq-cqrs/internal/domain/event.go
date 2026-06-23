package domain

import "time"

type Event interface {
	AggregateID() string
	EventType() string
	OccurredAt() time.Time
}

// UserCreatedEvent はユーザー作成時に発生するドメインイベント。
type UserCreatedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (e UserCreatedEvent) AggregateID() string { return e.ID }
func (e UserCreatedEvent) EventType() string   { return "UserCreatedEvent" }
func (e UserCreatedEvent) OccurredAt() time.Time { return e.CreatedAt }

// UserUpdatedEvent はユーザー更新時に発生するドメインイベント。
type UserUpdatedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e UserUpdatedEvent) AggregateID() string   { return e.ID }
func (e UserUpdatedEvent) EventType() string     { return "UserUpdatedEvent" }
func (e UserUpdatedEvent) OccurredAt() time.Time { return e.UpdatedAt }
