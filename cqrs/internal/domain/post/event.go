package post

import (
	"time"
)

type PostCreatedEvent struct {
	id         PostID
	title      Title
	content    Content
	authorID   AuthorID
	status     Status
	occurredAt time.Time
}

func NewPostCreatedEvent(post *Post) PostCreatedEvent {
	return PostCreatedEvent{
		id:         post.id,
		title:      post.title,
		content:    post.content,
		authorID:   post.authorID,
		status:     post.status,
		occurredAt: time.Now(),
	}
}

type PostUpdatedEvent struct {
	id         PostID
	title      Title
	content    Content
	status     Status
	occurredAt time.Time
}

func NewPostUpdatedEvent(post *Post) PostUpdatedEvent {
	return PostUpdatedEvent{
		id:         post.id,
		title:      post.title,
		content:    post.content,
		status:     post.status,
		occurredAt: time.Now(),
	}
}

type PostDeletedEvent struct {
	id         PostID
	occurredAt time.Time
}

func NewPostDeletedEvent(post *Post) PostDeletedEvent {
	return PostDeletedEvent{
		id:         post.id,
		occurredAt: time.Now(),
	}
}
