package dto

import (
	"cqrs_exercise/internal/domain/user"
	"time"
)

type UserCreatedEventDTO struct {
	AggregateID string    `json:"aggregate_id"`
	EventType   string    `json:"event_type"`
	EventTime   time.Time `json:"event_time"`
	Version     int       `json:"version"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Email       string    `json:"email"`
	Profile     string    `json:"profile"`
}

func NewUserCreateEventDTO(u user.UserCreatedEvent) UserCreatedEventDTO {
	return UserCreatedEventDTO{
		AggregateID: u.AggregateID(),
		EventType:   u.EventType(),
		EventTime:   u.OccurredAt(),
		Version:     u.AggregateVersion(),
		UserID:      u.ID().String(),
		Name:        u.Name().String(),
		Image:       u.Image(),
		Email:       u.Email().String(),
		Profile:     u.Profile().String(),
	}
}
