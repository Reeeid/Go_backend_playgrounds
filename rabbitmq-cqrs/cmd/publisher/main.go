// publisher はコマンド側のデモ。
// ユーザー作成・更新コマンドを受け付けてドメインイベントを RabbitMQ に Publish する。
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"rabbitmq_cqrs/internal/domain"
	"rabbitmq_cqrs/internal/eventbus"
)

const amqpURL = "amqp://guest:guest@localhost:5672/"

func main() {
	bus, err := eventbus.New(amqpURL)
	if err != nil {
		log.Fatalf("connect rabbitmq: %v", err)
	}
	defer bus.Close()

	ctx := context.Background()

	// ── コマンド1: ユーザー作成 ──────────────────────────────
	aliceID := uuid.NewString()
	created := domain.UserCreatedEvent{
		ID:        aliceID,
		Name:      "Alice",
		Email:     "alice@example.com",
		CreatedAt: time.Now(),
	}
	if err := bus.Publish(ctx, created.EventType(), created); err != nil {
		log.Fatalf("publish UserCreatedEvent: %v", err)
	}
	fmt.Printf("[publisher] UserCreatedEvent: id=%s\n", aliceID)

	time.Sleep(500 * time.Millisecond)

	// ── コマンド2: ユーザー更新 ──────────────────────────────
	updated := domain.UserUpdatedEvent{
		ID:        aliceID,
		Name:      "Alice Updated",
		Email:     "alice.updated@example.com",
		UpdatedAt: time.Now(),
	}
	if err := bus.Publish(ctx, updated.EventType(), updated); err != nil {
		log.Fatalf("publish UserUpdatedEvent: %v", err)
	}
	fmt.Printf("[publisher] UserUpdatedEvent: id=%s\n", aliceID)

	// ── コマンド3: 別ユーザー作成 ────────────────────────────
	bobID := uuid.NewString()
	bobCreated := domain.UserCreatedEvent{
		ID:        bobID,
		Name:      "Bob",
		Email:     "bob@example.com",
		CreatedAt: time.Now(),
	}
	if err := bus.Publish(ctx, bobCreated.EventType(), bobCreated); err != nil {
		log.Fatalf("publish UserCreatedEvent: %v", err)
	}
	fmt.Printf("[publisher] UserCreatedEvent: id=%s\n", bobID)

	fmt.Println("[publisher] done — check subscriber output")
}
