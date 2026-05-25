package outbox

import (
	"context"
	"fmt"
	"product-command-module/internal/application/port"
)

func (r *outboxRepository) SaveOutbox(ctx context.Context, event *port.OutboxEvent) error {
	if err := r.getQuerier(ctx).SaveOutbox(ctx, toInfraOutbox(event)); err != nil {
		return fmt.Errorf("infra: save outbox: %w", err)
	}

	return nil
}

func (r *outboxRepository) SaveOutboxBatch(ctx context.Context, events []port.OutboxEvent) error {
	if err := r.getQuerier(ctx).SaveOutboxBatch(ctx, toInfraOutboxBatch(events)); err != nil {
		return fmt.Errorf("infra: save outbox bath: %w", err)
	}

	return nil
}
