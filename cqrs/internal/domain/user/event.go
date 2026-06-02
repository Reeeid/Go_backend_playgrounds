package user

import "time"

type UserCreatedEvent struct {
	aggregateID string
	eventType   string
	occurredAt  time.Time
	id          UserID
	name        Name
	image       string
	email       Email
	profile     Profile
}

func (e UserCreatedEvent) AggregateID() string {
	return e.aggregateID
}
func (e UserCreatedEvent) EventType() string {
	return e.eventType
}
func (e UserCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func NewUserCreatedEvent(id UserID, name Name, image string, email Email, profile Profile) UserCreatedEvent {
	return UserCreatedEvent{
		aggregateID: id.String(),
		eventType:   "UserCreatedEvent",
		occurredAt:  time.Now(),
		id:          id,
		name:        name,
		image:       image,
		email:       email,
		profile:     profile,
	}
}

type UserUpdatedEvent struct {
	aggregateID string
	eventType   string
	occurredAt  time.Time
	id          UserID
	name        Name
	image       string
	email       Email
	profile     Profile
}

func (e UserUpdatedEvent) AggregateID() string {
	return e.aggregateID
}
func (e UserUpdatedEvent) EventType() string {
	return e.eventType
}
func (e UserUpdatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func NewUserUpdatedEvent(id UserID, name Name, image string, email Email, profile Profile) UserUpdatedEvent {
	return UserUpdatedEvent{
		aggregateID: id.String(),
		eventType:   "UserUpdatedEvent",
		occurredAt:  time.Now(),
		id:          id,
		name:        name,
		image:       image,
		email:       email,
		profile:     profile,
	}
}

type UserDeletedEvent struct {
	aggregateID string
	eventType   string
	occurredAt  time.Time
	id          UserID
}

func (e UserDeletedEvent) AggregateID() string {
	return e.aggregateID
}
func (e UserDeletedEvent) EventType() string {
	return e.eventType
}
func (e UserDeletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func NewUserDeletedEvent(id UserID) UserDeletedEvent {
	return UserDeletedEvent{
		aggregateID: id.String(),
		eventType:   "UserDeletedEvent",
		occurredAt:  time.Now(),
		id:          id,
	}
}

type UserPasswordCreatedEvent struct {
	aggregateID  string
	eventType    string
	occurredAt   time.Time
	userID       UserID
	hashPassword HashPassword
}

func (e UserPasswordCreatedEvent) AggregateID() string {
	return e.aggregateID
}
func (e UserPasswordCreatedEvent) EventType() string {
	return e.eventType
}
func (e UserPasswordCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func NewUserPasswordCreatedEvent(userID UserID, hashPassword HashPassword) UserPasswordCreatedEvent {
	return UserPasswordCreatedEvent{
		aggregateID:  userID.String(),
		eventType:    "UserPasswordSetEvent",
		occurredAt:   time.Now(),
		userID:       userID,
		hashPassword: hashPassword,
	}
}

type UserPasswordUpdatedEvent struct {
	aggregateID  string
	eventType    string
	occurredAt   time.Time
	userID       UserID
	hashPassword HashPassword
}

func (e UserPasswordUpdatedEvent) AggregateID() string {
	return e.aggregateID
}
func (e UserPasswordUpdatedEvent) EventType() string {
	return e.eventType
}
func (e UserPasswordUpdatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func NewUserPasswordUpdatedEvent(userID UserID, hashPassword HashPassword) UserPasswordUpdatedEvent {
	return UserPasswordUpdatedEvent{
		aggregateID:  userID.String(),
		eventType:    "UserPasswordUpdatedEvent",
		occurredAt:   time.Now(),
		userID:       userID,
		hashPassword: hashPassword,
	}
}

type UserPasswordDeletedEvent struct {
	aggregateID string
	eventType   string
	occurredAt  time.Time
	userID      UserID
}

func (e UserPasswordDeletedEvent) AggregateID() string {
	return e.aggregateID
}
func (e UserPasswordDeletedEvent) EventType() string {
	return e.eventType
}
func (e UserPasswordDeletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func NewUserPasswordDeletedEvent(userID UserID) UserPasswordDeletedEvent {
	return UserPasswordDeletedEvent{
		aggregateID: userID.String(),
		eventType:   "UserPasswordDeletedEvent",
		occurredAt:  time.Now(),
		userID:      userID,
	}
}
