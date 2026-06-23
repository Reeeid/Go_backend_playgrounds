package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"cqrs_claude/internal/application"
	"cqrs_claude/internal/domain/post"
	"cqrs_claude/internal/infrastructure/eventbus"
	"cqrs_claude/internal/infrastructure/eventstore"
	"cqrs_claude/internal/projection"
	"cqrs_claude/internal/query"
)

func main() {
	// Infrastructure
	es := eventstore.New()
	bus := eventbus.New()
	store := query.NewPostReadStore()

	// Projection
	proj := projection.NewPostProjection(store)
	bus.Subscribe(post.PostCreatedType, proj)
	bus.Subscribe(post.PostUpdatedType, proj)
	bus.Subscribe(post.PostDeletedType, proj)

	// Handlers（1コマンド1Handler）
	createHandler := application.NewCreatePostHandler(es, bus)
	updateHandler := application.NewUpdatePostHandler(es, bus)
	deleteHandler := application.NewDeletePostHandler(es, bus)

	ctx := context.Background()
	authorID := "018f4b3c-2d8e-7f3a-9b5c-1a2b3c4d5e6f"

	// IDは呼び出し側で生成（冪等性のため）
	postID := uuid.New().String()

	// Create
	if err := createHandler.Handle(ctx, application.CreatePostCommand{
		ID:       postID,
		Title:    "はじめてのCQRS+ES",
		Content:  "イベントソーシングの実装例です",
		AuthorID: authorID,
	}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Created ===")
	printPost(store, postID)

	// Update
	if err := updateHandler.Handle(ctx, application.UpdatePostCommand{
		ID:      postID,
		Title:   "更新されたCQRS+ES",
		Content: "内容を更新しました",
		Status:  "published",
	}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Updated ===")
	printPost(store, postID)

	// Delete
	if err := deleteHandler.Handle(ctx, application.DeletePostCommand{
		ID: postID,
	}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("=== Deleted ===")
	if _, err := store.GetByID(postID); err != nil {
		fmt.Println("Post not found (as expected)")
	}
}

func printPost(store *query.PostReadStore, id string) {
	p, err := store.GetByID(id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("  ID:       %s\n", p.ID)
	fmt.Printf("  Title:    %s\n", p.Title)
	fmt.Printf("  Content:  %s\n", p.Content)
	fmt.Printf("  Status:   %s\n", p.Status)
	fmt.Printf("  AuthorID: %s\n", p.AuthorID)
}
