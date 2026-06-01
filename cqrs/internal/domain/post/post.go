package post

import (
	"cqrs_exercise/internal/domain/event"
	"time"
)

type Post struct {
	uncommittedEvents []event.Event
	id                PostID
	title             Title
	content           Content
	authorID          AuthorID
	status            Status
	createdAt         time.Time
	updatedAt         time.Time
	version           int
}

func NewPost(id PostID, title Title, content Content, authorID AuthorID) *Post {
	p := &Post{}
	p.raise(NewPostCreatedEvent(id, title, content, authorID))
	return p
}

func (p *Post) Apply(e event.Event) {
	switch e := e.(type) {
	case PostCreatedEvent:
		p.id = e.id
		p.title = e.title
		p.content = e.content
		p.authorID = e.authorID
		p.status = e.status
		p.version = 1
	case PostUpdatedEvent:
		p.title = e.title
		p.content = e.content
		p.status = e.status
		p.version++
	case PostDeletedEvent:
		p.status = StatusArchived
		p.version++
	}
}

func (p *Post) GetUncommittedEvents() []event.Event {
	return p.uncommittedEvents
}

func (p *Post) ClearUncommittedEvents() {
	p.uncommittedEvents = nil
}

func (p *Post) raise(e event.Event) {
	p.Apply(e)
	p.uncommittedEvents = append(p.uncommittedEvents, e)
}

func (p *Post) Update(title Title, content Content, status Status) {
	p.raise(NewPostUpdatedEvent(p.id, title, content, status))
}

func (p *Post) Delete() {
	p.raise(NewPostDeletedEvent(p.id))
}

func ConstructPostFromEvents(events []event.Event) *Post {
	p := &Post{}
	for _, e := range events {
		p.Apply(e)
	}
	return p
}
