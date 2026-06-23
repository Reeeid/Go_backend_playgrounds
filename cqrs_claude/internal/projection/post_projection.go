package projection

import (
	"cqrs_claude/internal/domain/event"
	"cqrs_claude/internal/domain/post"
	"cqrs_claude/internal/query"
)

// PostProjection はイベントを受け取りPostReadStoreを更新する
type PostProjection struct {
	store *query.PostReadStore
}

func NewPostProjection(store *query.PostReadStore) *PostProjection {
	return &PostProjection{store: store}
}

func (p *PostProjection) Handle(evt event.Event) error {
	switch e := evt.(type) {
	case post.PostCreatedEvent:
		p.store.Save(query.PostReadModel{
			ID:       e.ID().String(),
			Title:    e.Title().String(),
			Content:  e.Content().String(),
			AuthorID: e.AuthorID().String(),
			Status:   string(post.StatusDraft),
		})
	case post.PostUpdatedEvent:
		rm, err := p.store.GetByID(e.ID().String())
		if err != nil {
			return err
		}
		rm.Title = e.Title().String()
		rm.Content = e.Content().String()
		rm.Status = string(e.Status())
		p.store.Save(rm)
	case post.PostDeletedEvent:
		p.store.Delete(e.ID().String())
	}
	return nil
}
