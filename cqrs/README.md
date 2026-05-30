# CQRS + Event Sourcing

## アーキテクチャフロー

```mermaid
graph TD
    subgraph Write
        CMD[Command] --> CH[CommandHandler]
        CH -->|1. Load| ES[(EventStore)]
        ES -->|events| RC[ReconstructPost/User]
        RC --> AGG[Aggregate Post/User]
        AGG -->|raise| UE[uncommittedEvents]
        UE -->|2. Save| ES
        UE -->|3. Publish| EB[EventBus]
    end

    subgraph Read
        EB -->|Handle| PROJ[Projection]
        PROJ -->|update| RS[(ReadStore)]
        QH[QueryHandler] -->|ref| RS
        QRY[Query] --> QH
    end
```

## クラス図

```mermaid
classDiagram
    class Event {
        <<interface>>
        +AggregateID()
        +EventType()
        +OccurredAt()
    }
    class BaseEvent {
        +AggregateID
        -eventType
        -occurredAt
    }
    class EventStore {
        <<interface>>
        +Save()
        +Load()
    }
    class EventBus {
        <<interface>>
        +Publish()
        +Subscribe()
    }
    class EventHandler {
        <<interface>>
        +Handle()
    }
    class Post {
        -id
        -title
        -content
        -authorID
        -status
        -version
        +Apply()
        +Update()
        +Delete()
        +UncommittedEvents()
        +ClearUncommittedEvents()
    }
    class PostCreatedEvent {
        -id
        -title
        -content
        -authorID
    }
    class PostUpdatedEvent {
        -id
        -title
        -status
    }
    class PostDeletedEvent {
        -id
    }
    class User {
        -id
        -name
        -email
        -version
        +Apply()
        +Update()
        +Delete()
    }
    class Credential {
        -userID
        -password
        -createdAt
        -updatedAt
        +Update()
        +Delete()
    }
    class CreatePostHandler {
        -eventStore
        -publisher
        +Handle()
    }
    class UpdatePostHandler {
        -eventStore
        -publisher
        +Handle()
    }
    class DeletePostHandler {
        -eventStore
        -publisher
        +Handle()
    }
    class MemoryEventStore {
        -store
        +Save()
        +Load()
    }
    class MemoryEventBus {
        -handlers
        +Publish()
        +Subscribe()
    }
    class PostProjection {
        -store
        +Handle()
    }
    class PostReadStore {
        -store
        +Save()
        +GetByID()
        +Delete()
    }
    class PostReadModel {
        +ID
        +Title
        +Content
        +AuthorID
        +Status
    }

    BaseEvent ..|> Event
    PostCreatedEvent --|> BaseEvent
    PostUpdatedEvent --|> BaseEvent
    PostDeletedEvent --|> BaseEvent
    Post --> PostCreatedEvent : raises
    Post --> PostUpdatedEvent : raises
    Post --> PostDeletedEvent : raises
    Credential --> User : belongs to
    CreatePostHandler ..|> EventHandler
    UpdatePostHandler ..|> EventHandler
    DeletePostHandler ..|> EventHandler
    MemoryEventStore ..|> EventStore
    MemoryEventBus ..|> EventBus
    PostProjection ..|> EventHandler
    CreatePostHandler --> EventStore : uses
    CreatePostHandler --> EventBus : uses
    PostProjection --> PostReadStore : updates
```
