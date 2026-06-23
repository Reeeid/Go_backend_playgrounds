# CQRS + Event Sourcing 実装例

## アーキテクチャ概要

```
Command → Application → Aggregate → EventStore
                                  → EventBus → Projection → ReadStore
Query                                                      → QueryHandler → ReadStore
```

## 各層の責務

### domain/event
- `Event` インターフェイス（マーカー + 共通メソッド）
- `BaseEvent` 埋め込み構造体（全イベント共通フィールド）
- `EventStore` / `EventPublisher` インターフェイス定義

### domain/post
- `Post` 集約（`apply()` + `uncommittedEvents` でESを実現）
- ドメインイベント（`PostCreatedEvent` 等）
- Value Objects（`Title`, `Content` 等）

### application
- Command定義（`CreatePostCommand` 等）
- Command Handler（集約操作 → EventStore保存 → EventBus発行）

### infrastructure
- `MemoryEventStore` : イベントをインメモリに保存・再生
- `MemoryEventBus` : イベントを購読者に配信

### projection
- `PostProjection` : イベントを受け取りRead Modelを更新

### query
- `PostReadStore` : クエリ用の非正規化モデル（読み取り専用）

---

## ESの核心：イベント再生による状態復元

```go
// DBには現在の状態ではなくイベント履歴が保存される
events := [PostCreatedEvent, PostUpdatedEvent, PostPublishedEvent]

// イベントを順番に適用して現在の状態を復元
func ReconstructPost(events []event.Event) *Post {
    p := &Post{}
    for _, evt := range events {
        p.Apply(evt) // 各イベントで状態を変化させる
    }
    return p
}
```

## CQRSの核心：Write側とRead側の分離

```
Write側: Post集約 + EventStore（正規化されたイベント）
Read側:  PostReadModel（クエリに最適化された非正規化データ）
```

EventBusがProjectionを通じてRead Modelを同期する。

---

## 実行

```bash
go run ./cmd/main.go
```
