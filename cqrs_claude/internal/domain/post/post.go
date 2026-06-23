package post

import (
	"errors"

	"cqrs_claude/internal/domain/event"
)

var ErrCannotEditArchivedPost = errors.New("cannot edit archived post")

type Post struct {
	id                PostID
	title             Title
	content           Content
	authorID          AuthorID
	status            Status
	version           int
	uncommittedEvents []event.Event
}

// Apply はイベントを集約の状態に反映する（イベント再生時に使用）
func (p *Post) Apply(evt event.Event) {
	switch e := evt.(type) {
	case PostCreatedEvent:
		p.id = e.id
		p.title = e.title
		p.content = e.content
		p.authorID = e.authorID
		p.status = StatusDraft
	case PostUpdatedEvent:
		p.title = e.title
		p.content = e.content
		p.status = e.status
	case PostDeletedEvent:
		p.status = StatusArchived
	}
	p.version++
}

// raise はイベントを状態に反映しつつ未コミットリストに追加する
func (p *Post) raise(evt event.Event) {
	p.Apply(evt)
	p.uncommittedEvents = append(p.uncommittedEvents, evt)
}

func (p *Post) UncommittedEvents() []event.Event { return p.uncommittedEvents }
func (p *Post) ClearUncommittedEvents()          { p.uncommittedEvents = nil }
func (p *Post) ID() PostID                       { return p.id }
func (p *Post) Version() int                     { return p.version }

// NewPost は新規Postを作成しPostCreatedEventを発行する
func NewPost(id PostID, title Title, content Content, authorID AuthorID) *Post {
	p := &Post{}
	p.raise(PostCreatedEvent{
		BaseEvent: event.NewBaseEvent(id.String(), PostCreatedType),
		id:        id,
		title:     title,
		content:   content,
		authorID:  authorID,
	})
	return p
}

// Update はPostを更新しPostUpdatedEventを発行する
func (p *Post) Update(title Title, content Content, status Status) error {
	if p.status == StatusArchived {
		return ErrCannotEditArchivedPost
	}
	p.raise(PostUpdatedEvent{
		BaseEvent: event.NewBaseEvent(p.id.String(), PostUpdatedType),
		id:        p.id,
		title:     title,
		content:   content,
		status:    status,
	})
	return nil
}

// Delete はPostを削除しPostDeletedEventを発行する
func (p *Post) Delete() {
	p.raise(PostDeletedEvent{
		BaseEvent: event.NewBaseEvent(p.id.String(), PostDeletedType),
		id:        p.id,
	})
}

// ReconstructPost はイベント履歴からPostの状態を再構築する（ESの核心）
func ReconstructPost(events []event.Event) *Post {
	p := &Post{}
	for _, evt := range events {
		p.Apply(evt)
	}
	return p
}
