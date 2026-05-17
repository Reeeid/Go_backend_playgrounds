package event

type EventStore interface {
	Save(events []Event) error
	Load(aggregateID string) ([]Event, error)
}
