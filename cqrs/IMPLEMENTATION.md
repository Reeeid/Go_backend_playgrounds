# CQRS + Event Sourcing 実装順

## 前提の整理

### EventStore の選択肢

| | PostgreSQL | DynamoDB |
|---|---|---|
| 楽観的排他制御 | UNIQUE制約 | Condition Expression（ビルトイン）|
| 集約の全イベント取得 | WHERE + ORDER BY | Partition Key = aggregateID, Sort Key = version |
| 向き不向き | 既存インフラに合わせるとき | スケール・スキーマレス重視 |

DynamoDB は Partition Key + Sort Key の組み合わせが EventStore と相性がいい。  
PostgreSQL でも問題なく作れる。どちらか決めてから実装に入る。

---

## 実装順

### Step 1 — Event インターフェースを固める

```go
// domain/event/event.go
type Event interface {
    AggregateID()      string
    AggregateVersion() int\
    EventType()        string
    OccurredAt()       time.Time
}
```

`AggregateVersion` がないと DB に保存したとき順序も排他制御もできない。

---

### Step 2 — ドメインイベントを直す（フィールドを公開）

```go
// domain/user/events.go
type UserCreatedEvent struct {
    ID         string    // 大文字に
    Name       string
    Email      string
    Version    int
    EventTime time.Time
}

func (e UserCreatedEvent) AggregateID()      string    { return e.ID }
func (e UserCreatedEvent) AggregateVersion() int       { return e.Version }
func (e UserCreatedEvent) EventType()        string    { return "UserCreatedEvent" }
func (e UserCreatedEvent) OccurredAt()       time.Time { return e.OccurredAt }
```

フィールドを公開することで `json.Marshal` がそのまま使える。  
Value Object（UserID, Name 等）は string/int に展開してイベントに持たせる。

---

### Step 3 — Aggregate を修正

```go
// domain/user/user.go
func (u *User) raise(e event.Event) {
    u.Apply(e)
    u.uncommittedEvents = append(u.uncommittedEvents, e)
}
```

`version` カウントは `Apply` 内でインクリメント。  
`NewUser` で最初のイベントを version=1 で生成。

---

### Step 4 — EventStore を実装

**PostgreSQL の場合**

```go
// infrastructure/eventstore/postgres.go
type eventRecord struct {
    AggregateID      string `gorm:"uniqueIndex:idx_agg_ver"`
    AggregateVersion int    `gorm:"uniqueIndex:idx_agg_ver"` // ← 重複でエラー = 排他制御
    EventType        string
    Payload          []byte
    OccurredAt       time.Time
}
```

**DynamoDB の場合**

```
Table: events
  PK: aggregate_id  (string)
  SK: version       (number)
  event_type, payload, occurred_at
```

Condition Expression `attribute_not_exists(version)` で楽観的排他制御。

---

### Step 5 — デシリアライズのレジストリを作る

```go
// infrastructure/eventstore/registry.go
var Registry = map[string]func([]byte) (event.Event, error){
    "UserCreatedEvent": func(b []byte) (event.Event, error) {
        var e user.UserCreatedEvent
        return e, json.Unmarshal(b, &e)
    },
    "UserUpdatedEvent": func(b []byte) (event.Event, error) { ... },
}
```

`Load` 時に `event_type` でこのマップを引いて型を復元する。

---

### Step 6 — EventBus を整える

```go
// infrastructure/eventbus/memory.go
type MemoryEventBus struct { ... }

func (b *MemoryEventBus) Publish(evt event.Event) error { ... }
func (b *MemoryEventBus) Subscribe(eventType string, handler event.EventHandler) { ... }
```

インターフェースの `Subscribe` の戻り値（error の有無）を実装と合わせる。

---

### Step 7 — Command Handler を実装

```go
// application/user_handler.go
func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) error {
    u := user.NewUser(...)                          // 1. Aggregate 生成
    events := u.GetUncommittedEvents()
    if err := h.eventStore.Save(events); err != nil { // 2. EventStore に保存
        return err
    }
    for _, e := range events {
        h.eventBus.Publish(e)                       // 3. EventBus で通知
    }
    u.ClearUncommittedEvents()
    return nil
}
```

Update/Delete は `Load → ConstructFromEvents → 操作 → Save` の流れ。

---

### Step 8 — Projection を実装

EventBus の subscriber として登録し、イベントを受け取って Read Model（PostgreSQL の通常テーブルや Redis 等）に書き込む。

```go
// projection/user_projection.go
type UserProjection struct { db *gorm.DB }

func (p *UserProjection) Handle(evt event.Event) error {
    switch e := evt.(type) {
    case user.UserCreatedEvent:
        return p.db.Create(&UserView{ID: e.ID, Name: e.Name}).Error
    }
    return nil
}
```

---

### Step 9 — Query Handler を実装

Read Model から普通に SELECT するだけ。EventStore は使わない。

```go
// query/user_query.go
func (h *UserQueryHandler) GetByID(id string) (*UserView, error) {
    var v UserView
    return &v, h.db.Where("id = ?", id).First(&v).Error
}
```

---

### Step 10 — main.go で wire up

```go
db, _ := db.NewDB(dsn)
store := eventstore.New(db, eventstore.Registry)
bus   := eventbus.New()

projection := &projection.UserProjection{DB: db}
bus.Subscribe("UserCreatedEvent", projection)

createUserHandler := application.NewCreateUserHandler(store, bus)
```

---

## ファイル構成（最終形）

```
internal/
├── domain/
│   ├── event/
│   │   └── event.go          # Event / EventStore / EventBus インターフェース
│   ├── user/
│   │   ├── user.go           # Aggregate
│   │   ├── events.go         # イベント定義（フィールド公開）
│   │   ├── values.go
│   │   └── repository.go     # UserQueryRepository のみ
│   └── post/
│       └── ...
├── application/
│   ├── handler.go            # CommandHandler インターフェース
│   ├── user_command.go
│   ├── user_handler.go
│   ├── post_command.go
│   └── post_handler.go
├── projection/
│   ├── user_projection.go
│   └── post_projection.go
├── query/
│   ├── user_query.go
│   └── post_query.go
└── infrastructure/
    ├── db/postgres.go
    ├── eventstore/
    │   ├── postgres.go       # or dynamodb.go
    │   └── registry.go       # デシリアライズ用マップ
    └── eventbus/
        └── memory.go
```
