package post

import (
	"cqrs_exercise/internal/domain/event"
)

type PostCreatedEvent struct {
	event.BaseEvent
	id       PostID
	title    Title
	content  Content
	authorID AuthorID
	status   Status
}

func NewPostCreatedEvent(post *Post) PostCreatedEvent {
	return PostCreatedEvent{
		id:        post.id,
		title:     post.title,
		content:   post.content,
		authorID:  post.authorID,
		status:    post.status,
		BaseEvent: event.NewBaseEvent(post.id.String(), "PostCreatedEvent"),
	}
}

type PostUpdatedEvent struct {
	event.BaseEvent
	id      PostID
	title   Title
	content Content
	status  Status
}

func NewPostUpdatedEvent(post *Post) PostUpdatedEvent {
	return PostUpdatedEvent{
		id:        post.id,
		title:     post.title,
		content:   post.content,
		status:    post.status,
		BaseEvent: event.NewBaseEvent(post.id.String(), "PostUpdatedEvent"),
	}
}

type PostDeletedEvent struct {
	event.BaseEvent
	id PostID
}

func NewPostDeletedEvent(post *Post) PostDeletedEvent {
	return PostDeletedEvent{
		id:        post.id,
		BaseEvent: event.NewBaseEvent(post.id.String(), "PostDeletedEvent"),
	}
}
