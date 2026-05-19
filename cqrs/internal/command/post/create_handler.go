package command

import (
	"context"

	"cqrs_exercise/internal/domain/event"
	"cqrs_exercise/internal/domain/post"
)

type CreatePostHandler struct {
	postRepo   post.PostCommandRepository
	eventStore event.EventStore
	publisher  event.EventPublisher
}

func NewCreatePostHandler(postRepo post.PostCommandRepository, eventStore event.EventStore, publisher event.EventPublisher) *CreatePostHandler {
	return &CreatePostHandler{
		postRepo:   postRepo,
		eventStore: eventStore,
		publisher:  publisher,
	}
}

func (h *CreatePostHandler) Handle(ctx context.Context, command CreatePostCommand) error {
	return nil
}

type UpdatePostHandler struct {
	postRepo   post.PostCommandRepository
	eventStore event.EventStore
	publisher  event.EventPublisher
}

func NewUpdatePostHandler(postRepo post.PostCommandRepository, eventStore event.EventStore, publisher event.EventPublisher) *UpdatePostHandler {
	return &UpdatePostHandler{
		postRepo:   postRepo,
		eventStore: eventStore,
		publisher:  publisher,
	}
}

func (h *UpdatePostHandler) Handle(ctx context.Context, command UpdatePostCommand) error {
	return nil
}

type DeletePostHandler struct {
	postRepo   post.PostCommandRepository
	eventStore event.EventStore
	publisher  event.EventPublisher
}

func NewDeletePostHandler(postRepo post.PostCommandRepository, eventStore event.EventStore, publisher event.EventPublisher) *DeletePostHandler {
	return &DeletePostHandler{
		postRepo:   postRepo,
		eventStore: eventStore,
		publisher:  publisher,
	}
}

func (h *DeletePostHandler) Handle(ctx context.Context, command DeletePostCommand) error {
	return nil
}
