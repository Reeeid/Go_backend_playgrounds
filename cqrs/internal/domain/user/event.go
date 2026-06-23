package user

import "time"

type UserCreatedEvent struct {
	aggregateID string
	eventType   string
	eventTime   time.Time
	version     int
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
	return e.eventTime
}
func (e UserCreatedEvent) AggregateVersion() int {
	return e.version
}
func (e UserCreatedEvent) ID() UserID {
	return e.id
}
func (e UserCreatedEvent) Name() Name {
	return e.name
}
func (e UserCreatedEvent) Image() string {
	return e.image
}
func (e UserCreatedEvent) Email() Email {
	return e.email
}
func (e UserCreatedEvent) Profile() Profile {
	return e.profile
}

func NewUserCreatedEvent(id UserID, name Name, image string, email Email, profile Profile) UserCreatedEvent {
	return UserCreatedEvent{
		aggregateID: id.String(),
		eventType:   "UserCreatedEvent",
		eventTime:   time.Now(),
		version:     1,
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
	eventTime   time.Time
	version     int
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
	return e.eventTime
}
func (e UserUpdatedEvent) AggregateVersion() int {
	return e.version
}
func (e UserUpdatedEvent) ID() UserID {
	return e.id
}
func (e UserUpdatedEvent) Name() Name {
	return e.name
}
func (e UserUpdatedEvent) Image() string {
	return e.image
}
func (e UserUpdatedEvent) Email() Email {
	return e.email
}
func (e UserUpdatedEvent) Profile() Profile {
	return e.profile
}

func NewUserUpdatedEvent(id UserID, version int, name Name, image string, email Email, profile Profile) UserUpdatedEvent {
	return UserUpdatedEvent{
		aggregateID: id.String(),
		eventType:   "UserUpdatedEvent",
		eventTime:   time.Now(),
		version:     version,
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
	eventTime   time.Time
	version     int
	id          UserID
}

func (e UserDeletedEvent) AggregateID() string {
	return e.aggregateID
}
func (e UserDeletedEvent) EventType() string {
	return e.eventType
}
func (e UserDeletedEvent) OccurredAt() time.Time {
	return e.eventTime
}
func (e UserDeletedEvent) AggregateVersion() int {
	return e.version
}
func (e UserDeletedEvent) ID() UserID {
	return e.id
}

func NewUserDeletedEvent(id UserID, version int) UserDeletedEvent {
	return UserDeletedEvent{
		aggregateID: id.String(),
		eventType:   "UserDeletedEvent",
		eventTime:   time.Now(),
		version:     version,
		id:          id,
	}
}

type UserPasswordCreatedEvent struct {
	aggregateID  string
	eventType    string
	eventTime    time.Time
	version      int
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
	return e.eventTime
}
func (e UserPasswordCreatedEvent) AggregateVersion() int {
	return e.version
}
func (e UserPasswordCreatedEvent) UserID() UserID {
	return e.userID
}
func (e UserPasswordCreatedEvent) HashPassword() HashPassword {
	return e.hashPassword
}

func NewUserPasswordCreatedEvent(userID UserID, hashPassword HashPassword) UserPasswordCreatedEvent {
	return UserPasswordCreatedEvent{
		aggregateID:  userID.String(),
		eventType:    "UserPasswordSetEvent",
		eventTime:    time.Now(),
		version:      1,
		userID:       userID,
		hashPassword: hashPassword,
	}
}

type UserPasswordUpdatedEvent struct {
	aggregateID  string
	eventType    string
	eventTime    time.Time
	version      int
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
	return e.eventTime
}
func (e UserPasswordUpdatedEvent) AggregateVersion() int {
	return e.version
}
func (e UserPasswordUpdatedEvent) UserID() UserID {
	return e.userID
}
func (e UserPasswordUpdatedEvent) HashPassword() HashPassword {
	return e.hashPassword
}

func NewUserPasswordUpdatedEvent(userID UserID, version int, hashPassword HashPassword) UserPasswordUpdatedEvent {
	return UserPasswordUpdatedEvent{
		aggregateID:  userID.String(),
		eventType:    "UserPasswordUpdatedEvent",
		eventTime:    time.Now(),
		version:      version,
		userID:       userID,
		hashPassword: hashPassword,
	}
}

type UserPasswordDeletedEvent struct {
	aggregateID string
	eventType   string
	eventTime   time.Time
	version     int
	userID      UserID
}

func (e UserPasswordDeletedEvent) AggregateID() string {
	return e.aggregateID
}
func (e UserPasswordDeletedEvent) EventType() string {
	return e.eventType
}
func (e UserPasswordDeletedEvent) OccurredAt() time.Time {
	return e.eventTime
}
func (e UserPasswordDeletedEvent) AggregateVersion() int {
	return e.version
}
func (e UserPasswordDeletedEvent) UserID() UserID {
	return e.userID
}

func NewUserPasswordDeletedEvent(userID UserID, version int) UserPasswordDeletedEvent {
	return UserPasswordDeletedEvent{
		aggregateID: userID.String(),
		eventType:   "UserPasswordDeletedEvent",
		eventTime:   time.Now(),
		version:     version,
		userID:      userID,
	}
}
