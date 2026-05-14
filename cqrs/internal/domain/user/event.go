package user

import (
	"time"
)

type UserCreatedEvent struct {
	id         UserID
	name       Name
	image      Image
	email      Email
	profile    Profile
	occurredAt time.Time
}

// イベントにポインタを返す明確な理由は薄い。不変な記録なので値で返す。
// ただし、イベントのサイズの大きさやコピーコストにもよる
func NewUserCreatedEvent(user *User) UserCreatedEvent {
	return UserCreatedEvent{
		id:         user.id,
		name:       user.name,
		image:      user.image,
		email:      user.email,
		profile:    user.profile,
		occurredAt: time.Now(),
	}
}

type UserUpdatedEvent struct {
	id         UserID
	name       Name
	image      Image
	email      Email
	profile    Profile
	occurredAt time.Time
}

func NewUserUpdatedEvent(user *User) UserUpdatedEvent {
	return UserUpdatedEvent{
		id:         user.id,
		name:       user.name,
		image:      user.image,
		email:      user.email,
		profile:    user.profile,
		occurredAt: time.Now(),
	}
}

type UserDeletedEvent struct {
	id         UserID
	occurredAt time.Time
}

func NewUserDeletedEvent(user *User) UserDeletedEvent {
	return UserDeletedEvent{
		id:         user.id,
		occurredAt: time.Now(),
	}
}

type UserPasswordSetEvent struct {
	userID       UserID
	hashPassword HashPassword
	occurredAt   time.Time
}

func NewUserPasswordSetEvent(userID UserID, hashPassword HashPassword) UserPasswordSetEvent {
	return UserPasswordSetEvent{
		userID:       userID,
		hashPassword: hashPassword,
		occurredAt:   time.Now(),
	}
}

type UserPasswordUpdatedEvent struct {
	userID       UserID
	hashPassword HashPassword
	occurredAt   time.Time
}

func NewUserPasswordUpdatedEvent(userID UserID, hashPassword HashPassword) UserPasswordUpdatedEvent {
	return UserPasswordUpdatedEvent{
		userID:       userID,
		hashPassword: hashPassword,
		occurredAt:   time.Now(),
	}
}

type UserPasswordDeletedEvent struct {
	userID     UserID
	occurredAt time.Time
}
