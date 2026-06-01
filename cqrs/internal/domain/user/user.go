package user

import (
	"cqrs_exercise/internal/domain/event"
	"time"
)

type User struct {
	uncommittedEvents []event.Event
	id                UserID
	name              Name
	image             string
	email             Email
	profile           Profile
	createdAt         time.Time
	updatedAt         time.Time
	version           int
}

func NewUser(id UserID, name Name, image string, email Email, profile Profile) *User {
	u := &User{}
	u.raise(NewUserCreatedEvent(id, name, image, email, profile))
	return u
}

func (u *User) Apply(event event.Event) {
	switch e := event.(type) {
	case UserCreatedEvent:
		u.id = e.id
		u.name = e.name
		u.image = e.image
		u.email = e.email
		u.profile = e.profile
		u.updatedAt = e.OccurredAt()
		u.version = 1
	case UserUpdatedEvent:
		u.name = e.name
		u.image = e.image
		u.email = e.email
		u.profile = e.profile
		u.updatedAt = e.OccurredAt()
		u.version++
	case UserDeletedEvent:
		u.updatedAt = e.OccurredAt()
		u.version++
	}
}

func (u *User) GetUncommittedEvents() []event.Event {
	return u.uncommittedEvents
}

func (u *User) ClearUncommittedEvents() {
	u.uncommittedEvents = nil
}

func (u *User) raise(e event.Event) {
	u.Apply(e)
	u.uncommittedEvents = append(u.uncommittedEvents, e)
}

func (u *User) Update(name Name, image string, email Email, profile Profile) {
	u.raise(NewUserUpdatedEvent(u.id, name, image, email, profile))
}

func (u *User) Delete() {
	u.raise(NewUserDeletedEvent(u.id))
}

func ConstructUserFromEvent(events []event.Event) *User {
	u := &User{}
	for _, e := range events {
		u.Apply(e)
	}
	return u
}
