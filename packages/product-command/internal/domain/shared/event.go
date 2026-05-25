package shared

type DomainEvent interface {
	EventName() string
	IntegrationPayload() map[string]interface{}
}

type EventEntity struct {
	events []DomainEvent
}

func (e *EventEntity) AddEvent(event DomainEvent) {
	e.events = append(e.events, event)
}

func (e *EventEntity) Flush() []DomainEvent {
	cloned := make([]DomainEvent, len(e.events))
	copy(cloned, e.events)
	e.events = nil
	return cloned
}

func (e *EventEntity) ClearEvent() {
	e.events = nil
}
