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

func NewCredential(userID UserID, hashPassword HashPassword) *Credential {
	c := &Credential{}
	c.raise(NewUserPasswordCreatedEvent(userID, hashPassword))
	return c
}

func (c *Credential) Apply(e event.Event) {
	switch e := e.(type) {
	case UserPasswordCreatedEvent:
		c.userID = e.userID
		c.password = e.hashPassword
		c.createdAt = e.OccurredAt()
		c.updatedAt = e.OccurredAt()
		c.version = 1
	case UserPasswordUpdatedEvent:
		c.password = e.hashPassword
		c.updatedAt = e.OccurredAt()
		c.version++
	case UserPasswordDeletedEvent:
		c.updatedAt = e.OccurredAt()
		c.version++
	}
}

func (c *Credential) GetUncommitedEvents() []event.Event {
	return c.uncommitedEvents
}

func (c *Credential) ClearUncommitedEvents() {
	c.uncommitedEvents = nil
}

func (c *Credential) raise(e event.Event) {
	c.Apply(e)
	c.uncommitedEvents = append(c.uncommitedEvents, e)
}

func (c *Credential) Update(hashPassword HashPassword) {
	c.raise(NewUserPasswordUpdatedEvent(c.userID, hashPassword))
}

func (c *Credential) Delete() {
	c.raise(NewUserPasswordDeletedEvent(c.userID))
}
