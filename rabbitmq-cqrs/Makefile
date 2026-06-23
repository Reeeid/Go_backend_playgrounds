.PHONY: up down setup subscriber publisher

up:
	docker compose up -d

down:
	docker compose down

setup:
	go mod tidy

# ターミナル1で起動（イベント受信・プロジェクション更新）
subscriber:
	go run ./cmd/subscriber

# ターミナル2で実行（イベント発行）
publisher:
	go run ./cmd/publisher
