// Package eventbus は RabbitMQ をバックエンドとしたイベントバスを提供する。
//
// CQRS での役割:
//   コマンド側 (Publisher) がイベントを exchange に送信し、
//   クエリ側  (Subscriber) がキューにバインドして受信・プロジェクションを更新する。
//
//   [Command Handler]
//        │ Publish("UserCreatedEvent", payload)
//        ▼
//   [RabbitMQ Exchange: "domain.events"  type=topic]
//        │  routing key = "UserCreatedEvent"
//        ▼
//   [Queue: "user.projection"]  ←  binding: "UserCreatedEvent", "UserUpdatedEvent"
//        │
//        ▼
//   [Query Side: Projection Handler]
package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const exchangeName = "domain.events"

// RabbitMQBus はイベントのパブリッシュとサブスクライブを担う。
type RabbitMQBus struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func New(url string) (*RabbitMQBus, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("amqp dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}
	// topic exchange: routing key でイベント種別ごとに振り分ける
	if err := ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("declare exchange: %w", err)
	}
	return &RabbitMQBus{conn: conn, ch: ch}, nil
}

func (b *RabbitMQBus) Close() {
	b.ch.Close()
	b.conn.Close()
}

// Publish はイベントを JSON にシリアライズして exchange に送る。
// routing key にイベント型名を使うことで、サブスクライバーが絞り込める。
func (b *RabbitMQBus) Publish(ctx context.Context, eventType string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return b.ch.PublishWithContext(ctx, exchangeName, eventType, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

// MessageHandler はキューから受信したメッセージを処理する関数型。
type MessageHandler func(eventType string, body []byte) error

// Subscribe はキューを作成して exchange にバインドし、メッセージをハンドラーに渡す。
// patterns は routing key パターン（例: "UserCreatedEvent", "User.*"）。
func (b *RabbitMQBus) Subscribe(queueName string, patterns []string, handler MessageHandler) error {
	q, err := b.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}
	for _, pattern := range patterns {
		if err := b.ch.QueueBind(q.Name, pattern, exchangeName, false, nil); err != nil {
			return fmt.Errorf("bind %s: %w", pattern, err)
		}
	}
	msgs, err := b.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.RoutingKey, msg.Body); err != nil {
				log.Printf("[eventbus] handler error (key=%s): %v — nack", msg.RoutingKey, err)
				msg.Nack(false, true) // requeue
				continue
			}
			msg.Ack(false)
		}
	}()

	log.Printf("[eventbus] subscribed queue=%s patterns=%v", queueName, patterns)
	return nil
}
