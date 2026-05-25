package outbox

import (
	"encoding/json"
	"product-command-module/db/repository"
	"product-command-module/internal/application/port"

	"github.com/iamKienb/go-core/postgres/conv"
	"github.com/jackc/pgx/v5/pgtype"
)

func toInfraOutbox(e *port.OutboxEvent) repository.SaveOutboxParams {
	payload, _ := json.Marshal(e.Payload)
	return repository.SaveOutboxParams{
		ID:             conv.UUID(e.ID),
		AggregateID:    conv.UUID(e.AggregateID),
		AggregateType:  e.AggregateType,
		EventType:      e.EventType,
		Payload:        payload,
		PartitionKey:   e.PartitionKey,
		IdempotencyKey: conv.UUID(e.IdempotencyKey),
		CreatedAt:      conv.TimeStampZ(&e.CreatedAt),
	}
}

func toInfraOutboxBatch(events []port.OutboxEvent) repository.SaveOutboxBatchParams {
	n := len(events)
	params := repository.SaveOutboxBatchParams{
		Ids:             make([]pgtype.UUID, n),
		AggregateIds:    make([]pgtype.UUID, n),
		AggregateTypes:  make([]string, n),
		EventTypes:      make([]string, n),
		Payloads:        make([][]byte, n),
		PartitionKeys:   make([]string, n),
		IdempotencyKeys: make([]pgtype.UUID, n),
		CreatedAts:      make([]pgtype.Timestamptz, n),
	}

	for i, e := range events {
		payload, err := json.Marshal(e.Payload)
		if err != nil {
			payload = []byte("{}")
		}

		params.Ids[i] = conv.UUID(e.ID)
		params.AggregateIds[i] = conv.UUID(e.AggregateID)
		params.AggregateTypes[i] = e.AggregateType
		params.EventTypes[i] = e.EventType
		params.Payloads[i] = payload
		params.PartitionKeys[i] = e.PartitionKey
		params.IdempotencyKeys[i] = conv.UUID(e.IdempotencyKey)
		params.CreatedAts[i] = conv.TimeStampZ(&e.CreatedAt)
	}

	return params
}
