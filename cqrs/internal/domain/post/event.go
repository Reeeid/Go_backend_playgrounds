package post

import (
	"time"
)

type PostCreatedEvent struct {
	aggregateID string
	eventTime   time.Time
	eventType   string
	version     int
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
	return e.eventTime
}
func (e PostCreatedEvent) AggregateVersion() int {
	return e.version
}
func (e PostCreatedEvent) ID() PostID {
	return e.id
}
func (e PostCreatedEvent) Title() Title {
	return e.title
}
func (e PostCreatedEvent) Content() Content {
	return e.content
}
func (e PostCreatedEvent) AuthorID() AuthorID {
	return e.authorID
}
func (e PostCreatedEvent) Status() Status {
	return e.status
}

func NewPostCreatedEvent(id PostID, title Title, content Content, authorID AuthorID) PostCreatedEvent {
	return PostCreatedEvent{
		aggregateID: id.String(),
		eventType:   "PostCreatedEvent",
		eventTime:   time.Now(),
		version:     1,
		id:          id,
		title:       title,
		content:     content,
		authorID:    authorID,
		status:      StatusDraft,
	}
}

type PostUpdatedEvent struct {
	aggregateID string
	eventTime   time.Time
	eventType   string
	version     int
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
	return e.eventTime
}
func (e PostUpdatedEvent) AggregateVersion() int {
	return e.version
}
func (e PostUpdatedEvent) ID() PostID {
	return e.id
}
func (e PostUpdatedEvent) Title() Title {
	return e.title
}
func (e PostUpdatedEvent) Content() Content {
	return e.content
}
func (e PostUpdatedEvent) Status() Status {
	return e.status
}

func NewPostUpdatedEvent(id PostID, version int, title Title, content Content, status Status) PostUpdatedEvent {
	return PostUpdatedEvent{
		aggregateID: id.String(),
		eventType:   "PostUpdatedEvent",
		eventTime:   time.Now(),
		version:     version,
		id:          id,
		title:       title,
		content:     content,
		status:      status,
	}
}

type PostDeletedEvent struct {
	aggregateID string
	EventTime   time.Time
	eventType   string
	version     int
	id          PostID
}

func (e PostDeletedEvent) AggregateID() string {
	return e.aggregateID
}
func (e PostDeletedEvent) EventType() string {
	return e.eventType
}
func (e PostDeletedEvent) OccurredAt() time.Time {
	return e.EventTime
}
func (e PostDeletedEvent) AggregateVersion() int {
	return e.version
}
func (e PostDeletedEvent) ID() PostID {
	return e.id
}

func NewPostDeletedEvent(id PostID, version int) PostDeletedEvent {
	return PostDeletedEvent{
		aggregateID: id.String(),
		eventType:   "PostDeletedEvent",
		EventTime:   time.Now(),
		version:     version,
		id:          id,
	}
}
