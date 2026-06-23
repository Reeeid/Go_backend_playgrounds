package post

import "cqrs_claude/internal/domain/event"

const (
	PostCreatedType = "PostCreated"
	PostUpdatedType = "PostUpdated"
	PostDeletedType = "PostDeleted"
)

type PostCreatedEvent struct {
	event.BaseEvent
	id       PostID
	title    Title
	content  Content
	authorID AuthorID
}

func (e PostCreatedEvent) ID() PostID        { return e.id }
func (e PostCreatedEvent) Title() Title      { return e.title }
func (e PostCreatedEvent) Content() Content  { return e.content }
func (e PostCreatedEvent) AuthorID() AuthorID { return e.authorID }

type PostUpdatedEvent struct {
	event.BaseEvent
	id      PostID
	title   Title
	content Content
	status  Status
}

func (e PostUpdatedEvent) ID() PostID       { return e.id }
func (e PostUpdatedEvent) Title() Title     { return e.title }
func (e PostUpdatedEvent) Content() Content { return e.content }
func (e PostUpdatedEvent) Status() Status   { return e.status }

type PostDeletedEvent struct {
	event.BaseEvent
	id PostID
}

func (e PostDeletedEvent) ID() PostID { return e.id }
