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

func (p *Post) Apply(e event.Event) {
	switch e := e.(type) {
	case PostCreatedEvent:
		p.id = e.id
		p.title = e.title
		p.content = e.content
		p.authorID = e.authorID
		p.status = e.status
		p.createdAt = e.OccurredAt()
		p.updatedAt = e.OccurredAt()
		p.version = 1
	case PostUpdatedEvent:
		p.title = e.title
		p.content = e.content
		p.status = e.status
		p.updatedAt = e.OccurredAt()
		p.version++
	case PostDeletedEvent:
		// PostDeletedEventが発生した場合の処理は、コマンドハンドラーで行うため、ここでは何もしない
	}
}

func (p *Post) GetUncommittedEvents() []event.Event {
	return p.uncommittedEvents
}

func (p *Post) ClearUncommittedEvents() {
	p.uncommittedEvents = nil
}
