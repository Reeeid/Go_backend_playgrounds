package user

import (
	"cqrs_exercise/internal/domain/event"
)

type UserCreatedEvent struct {
	event.BaseEvent
	id      UserID
	name    Name
	image   string
	email   Email
	profile Profile
}

// イベントにポインタを返す明確な理由は薄い。不変な記録なので値で返す。
// ただし、イベントのサイズの大きさやコピーコストにもよる
func NewUserCreatedEvent(user *User) UserCreatedEvent {
	return UserCreatedEvent{
		id:        user.id,
		name:      user.name,
		image:     user.image,
		email:     user.email,
		profile:   user.profile,
		BaseEvent: event.NewBaseEvent(user.id.String(), "UserCreatedEvent", 1),
	}
}

type UserUpdatedEvent struct {
	event.BaseEvent
	id      UserID
	name    Name
	image   string
	email   Email
	profile Profile
}

func NewUserUpdatedEvent(user *User) UserUpdatedEvent {
	return UserUpdatedEvent{
		id:        user.id,
		name:      user.name,
		image:     user.image,
		email:     user.email,
		profile:   user.profile,
		BaseEvent: event.NewBaseEvent(user.id.String(), "UserUpdatedEvent"),
	}
}

type UserDeletedEvent struct {
	event.BaseEvent
	id UserID
}

func NewUserDeletedEvent(user *User) UserDeletedEvent {
	return UserDeletedEvent{
		id:        user.id,
		BaseEvent: event.NewBaseEvent(user.id.String(), "UserDeletedEvent"),
	}
}

type UserPasswordSetEvent struct {
	event.BaseEvent
	userID       UserID
	hashPassword HashPassword
}

func NewUserPasswordSetEvent(userID UserID, hashPassword HashPassword) UserPasswordSetEvent {
	return UserPasswordSetEvent{
		userID:       userID,
		hashPassword: hashPassword,
		BaseEvent:    event.NewBaseEvent(userID.String(), "UserPasswordSetEvent"),
	}
}

type UserPasswordUpdatedEvent struct {
	event.BaseEvent
	userID       UserID
	hashPassword HashPassword
}

func NewUserPasswordUpdatedEvent(userID UserID, hashPassword HashPassword) UserPasswordUpdatedEvent {
	return UserPasswordUpdatedEvent{
		userID:       userID,
		hashPassword: hashPassword,
		BaseEvent:    event.NewBaseEvent(userID.String(), "UserPasswordUpdatedEvent"),
	}
}

type UserPasswordDeletedEvent struct {
	event.BaseEvent
	userID UserID
}

func NewUserPasswordDeletedEvent(userID UserID) UserPasswordDeletedEvent {
	return UserPasswordDeletedEvent{
		userID:    userID,
		BaseEvent: event.NewBaseEvent(userID.String(), "UserPasswordDeletedEvent"),
	}
}
