package application

import (
	"context"

	"github.com/google/uuid"

	"cqrs_claude/internal/domain/event"
	"cqrs_claude/internal/domain/post"
)

// CommandHandler は全コマンドハンドラーの共通インターフェイス
type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

// saveAndPublish はイベントの保存と発行をまとめたヘルパー
func saveAndPublish(es event.EventStore, pub event.EventPublisher, events []event.Event) error {
	if err := es.Save(events); err != nil {
		return err
	}
	for _, evt := range events {
		if err := pub.Publish(evt); err != nil {
			return err
		}
	}
	return nil
}

// --- Create ---

type CreatePostCommand struct {
	ID       string // 呼び出し側で事前生成（冪等性のため）
	Title    string
	Content  string
	AuthorID string
}

type CreatePostHandler struct {
	eventStore event.EventStore
	publisher  event.EventPublisher
}

func NewCreatePostHandler(es event.EventStore, pub event.EventPublisher) *CreatePostHandler {
	return &CreatePostHandler{eventStore: es, publisher: pub}
}

func (h *CreatePostHandler) Handle(ctx context.Context, cmd CreatePostCommand) error {
	postUUID, err := uuid.Parse(cmd.ID)
	if err != nil {
		return err
	}
	postID, err := post.NewPostIDFromUUID(postUUID)
	if err != nil {
		return err
	}
	title, err := post.NewTitle(cmd.Title)
	if err != nil {
		return err
	}
	content, err := post.NewContent(cmd.Content)
	if err != nil {
		return err
	}
	authorUUID, err := uuid.Parse(cmd.AuthorID)
	if err != nil {
		return err
	}
	authorID, err := post.NewAuthorID(authorUUID)
	if err != nil {
		return err
	}

	p := post.NewPost(postID, title, content, authorID)
	if err := saveAndPublish(h.eventStore, h.publisher, p.UncommittedEvents()); err != nil {
		return err
	}
	p.ClearUncommittedEvents()
	return nil
}

// --- Update ---

type UpdatePostCommand struct {
	ID      string
	Title   string
	Content string
	Status  string
}

type UpdatePostHandler struct {
	eventStore event.EventStore
	publisher  event.EventPublisher
}

func NewUpdatePostHandler(es event.EventStore, pub event.EventPublisher) *UpdatePostHandler {
	return &UpdatePostHandler{eventStore: es, publisher: pub}
}

func (h *UpdatePostHandler) Handle(ctx context.Context, cmd UpdatePostCommand) error {
	events, err := h.eventStore.Load(cmd.ID)
	if err != nil {
		return err
	}
	p := post.ReconstructPost(events)

	title, err := post.NewTitle(cmd.Title)
	if err != nil {
		return err
	}
	content, err := post.NewContent(cmd.Content)
	if err != nil {
		return err
	}
	status, err := post.NewStatus(cmd.Status)
	if err != nil {
		return err
	}

	if err := p.Update(title, content, status); err != nil {
		return err
	}
	if err := saveAndPublish(h.eventStore, h.publisher, p.UncommittedEvents()); err != nil {
		return err
	}
	p.ClearUncommittedEvents()
	return nil
}

// --- Delete ---

type DeletePostCommand struct {
	ID string
}

type DeletePostHandler struct {
	eventStore event.EventStore
	publisher  event.EventPublisher
}

func NewDeletePostHandler(es event.EventStore, pub event.EventPublisher) *DeletePostHandler {
	return &DeletePostHandler{eventStore: es, publisher: pub}
}

func (h *DeletePostHandler) Handle(ctx context.Context, cmd DeletePostCommand) error {
	events, err := h.eventStore.Load(cmd.ID)
	if err != nil {
		return err
	}
	p := post.ReconstructPost(events)
	p.Delete()

	if err := saveAndPublish(h.eventStore, h.publisher, p.UncommittedEvents()); err != nil {
		return err
	}
	p.ClearUncommittedEvents()
	return nil
}
