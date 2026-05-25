package port

import (
	"context"
	"product-command-module/internal/domain/shared"
	"time"

	"github.com/google/uuid"
)

type OutboxParam struct {
	AggregateID   uuid.UUID
	AggregateType string
	Events        []shared.DomainEvent
}

type OutboxEvent struct {
	ID             uuid.UUID
	AggregateID    uuid.UUID
	AggregateType  string
	EventType      string
	Payload        interface{}
	PartitionKey   string
	IdempotencyKey uuid.UUID
	CreatedAt      time.Time
}

type OutboxRepository interface {
	SaveOutbox(ctx context.Context, event *OutboxEvent) error
	SaveOutboxBatch(ctx context.Context, event []OutboxEvent) error
}
