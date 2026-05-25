package user

import "cqrs_exercise/internal/domain/event"

type User struct {
	id      UserID
	name    Name
	image   string
	email   Email
	profile Profile
}

func (u *User) Apply(event event.Event) {
	switch e := event.(type) {
	case UserCreatedEvent:
		u.id = e.id
		u.name = e.name
		u.image = e.image
		u.email = e.email
		u.profile = e.profile
	case UserUpdatedEvent:
		u.name = e.name
		u.image = e.image
		u.email = e.email
		u.profile = e.profile
	case UserDeletedEvent:
		// UserDeletedEventが発生した場合の処理は、コマンドハンドラーで行うため、ここでは何もしない
	}
}

func (u *User) GetUncommittedEvents() []event.Event {
	// Userエンティティはイベントを保持しないため、常に空のスライスを返す
	return nil
}

func (u *User) ClearUncommittedEvents() {
	// Userエンティティはイベントを保持しないため、何もしない
}
