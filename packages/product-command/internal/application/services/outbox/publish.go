package outbox

import (
	"context"
	"product-command-module/internal/application/port"
	"time"

	"github.com/google/uuid"
)

func (s *outboxService) Publish(ctx context.Context, param port.OutboxParam) error {
	events := param.Events
	if len(events) == 0 {
		return nil
	}

	messages := make([]port.OutboxEvent, 0, len(events))
	now := time.Now().UTC()

	for _, event := range events {
		messages = append(messages, port.OutboxEvent{
			ID:             uuid.New(),
			AggregateID:    param.AggregateID,
			AggregateType:  param.AggregateType,
			EventType:      event.EventName(),
			Payload:        event.IntegrationPayload(),
			PartitionKey:   param.AggregateID.String(),
			IdempotencyKey: uuid.New(),
			CreatedAt:      now,
		})
	}

	return s.outboxRepo.SaveOutboxBatch(ctx, messages)
}
