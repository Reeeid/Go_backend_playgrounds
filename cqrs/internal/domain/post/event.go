package post

import (
	"time"
)

type PostCreatedEvent struct {
	aggregateID string
	occurredAt  time.Time
	eventType   string
	id          PostID
	title       Title
	content     Content
	authorID    AuthorID
	status      Status
}

func (e PostCreatedEvent) AggregateID() string {
	return e.aggregateID
}
func (e PostCreatedEvent) EventType() string {
	return e.eventType
}
func (e PostCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func NewPostCreatedEvent(id PostID, title Title, content Content, authorID AuthorID) PostCreatedEvent {
	return PostCreatedEvent{
		aggregateID: id.String(),
		eventType:   "PostCreatedEvent",
		occurredAt:  time.Now(),
		id:          id,
		title:       title,
		content:     content,
		authorID:    authorID,
		status:      StatusDraft,
	}
}

type PostUpdatedEvent struct {
	aggregateID string
	occurredAt  time.Time
	eventType   string
	id          PostID
	title       Title
	content     Content
	status      Status
}

func (e PostUpdatedEvent) AggregateID() string {
	return e.aggregateID
}
func (e PostUpdatedEvent) EventType() string {
	return e.eventType
}
func (e PostUpdatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func NewPostUpdatedEvent(id PostID, title Title, content Content, status Status) PostUpdatedEvent {
	return PostUpdatedEvent{
		aggregateID: id.String(),
		eventType:   "PostUpdatedEvent",
		occurredAt:  time.Now(),
		id:          id,
		title:       title,
		content:     content,
		status:      status,
	}
}

type PostDeletedEvent struct {
	aggregateID string
	occurredAt  time.Time
	eventType   string
	id          PostID
}

func (e PostDeletedEvent) AggregateID() string {
	return e.aggregateID
}
func (e PostDeletedEvent) EventType() string {
	return e.eventType
}
func (e PostDeletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func NewPostDeletedEvent(id PostID) PostDeletedEvent {
	return PostDeletedEvent{
		aggregateID: id.String(),
		eventType:   "PostDeletedEvent",
		occurredAt:  time.Now(),
		id:          id,
	}
}
