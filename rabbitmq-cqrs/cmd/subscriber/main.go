// subscriber はクエリ側のデモ。
// RabbitMQ からイベントを受信してプロジェクション（読み取りモデル）を更新し続ける。
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rabbitmq_cqrs/internal/eventbus"
	"rabbitmq_cqrs/internal/projection"
)

const amqpURL = "amqp://guest:guest@localhost:5672/"

func main() {
	bus, err := eventbus.New(amqpURL)
	if err != nil {
		log.Fatalf("connect rabbitmq: %v", err)
	}
	defer bus.Close()

	userProj := projection.NewUserProjection()

	// UserCreatedEvent と UserUpdatedEvent を同一キューで受信
	if err := bus.Subscribe(
		"user.projection",                              // キュー名
		[]string{"UserCreatedEvent", "UserUpdatedEvent"}, // routing key パターン
		userProj.Handle,                               // ハンドラー
	); err != nil {
		log.Fatalf("subscribe: %v", err)
	}

	fmt.Println("[subscriber] waiting for events (Ctrl+C to stop)")

	// 5秒ごとに現在のプロジェクション状態を表示
	go func() {
		for {
			time.Sleep(5 * time.Second)
			users := userProj.All()
			fmt.Printf("\n── Read Model (users: %d) ──\n", len(users))
			for _, u := range users {
				fmt.Printf("  id=%-36s name=%s email=%s\n", u.ID, u.Name, u.Email)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\n[subscriber] shutting down")
}
