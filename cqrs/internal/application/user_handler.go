package command

import (
	"context"

	"cqrs_exercise/internal/domain/event"
	"cqrs_exercise/internal/domain/user"
)

type CreateUserHandler struct {
	userRepo   user.UserCommandRepository
	eventStore event.EventStore
	publisher  event.EventPublisher
}

func NewCreateUserHandler(userRepo user.UserCommandRepository, eventStore event.EventStore, publisher event.EventPublisher) *CreateUserHandler {
	return &CreateUserHandler{
		userRepo:   userRepo,
		eventStore: eventStore,
		publisher:  publisher,
	}
}

func (h *CreateUserHandler) Handle(ctx context.Context, command CreateUserCommand) error {
	return nil
}

type UpdateUserHandler struct {
	userRepo   user.UserCommandRepository
	eventStore event.EventStore
	publisher  event.EventPublisher
}

func NewUpdateUserHandler(userRepo user.UserCommandRepository, eventStore event.EventStore, publisher event.EventPublisher) *UpdateUserHandler {
	return &UpdateUserHandler{
		userRepo:   userRepo,
		eventStore: eventStore,
		publisher:  publisher,
	}
}

func (h *UpdateUserHandler) Handle(ctx context.Context, command UpdateUserCommand) error {
	return nil
}

type DeleteUserHandler struct {
	userRepo   user.UserCommandRepository
	eventStore event.EventStore
	publisher  event.EventPublisher
}

func NewDeleteUserHandler(userRepo user.UserCommandRepository, eventStore event.EventStore, publisher event.EventPublisher) *DeleteUserHandler {
	return &DeleteUserHandler{
		userRepo:   userRepo,
		eventStore: eventStore,
		publisher:  publisher,
	}
}

func (h *DeleteUserHandler) Handle(ctx context.Context, command DeleteUserCommand) error {
	return nil
}
