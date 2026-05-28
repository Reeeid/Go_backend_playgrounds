package user

import (
	"cqrs_exercise/internal/domain/event"
	"time"
)

type Credential struct {
	uncommitedEvents []event.Event
	userID           UserID
	password         HashPassword
	updatedAt        time.Time
	createdAt        time.Time
	version          int
}
