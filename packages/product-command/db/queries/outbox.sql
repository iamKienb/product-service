-- name: SaveOutbox :exec
INSERT INTO outbox (
    id,
    aggregate_id,
    aggregate_type,
    event_type,
    payload,
    partition_key,
    idempotency_key,
    created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);


-- name: SaveOutboxBatch :exec
INSERT INTO outbox (
    id,
    aggregate_id,
    aggregate_type,
    event_type,
    payload,
    partition_key,
    idempotency_key,
    created_at
)
SELECT 
    unnest(@ids::uuid[]),
    unnest(@aggregate_ids::uuid[]),
    unnest(@aggregate_types::text[]),
    unnest(@event_types::text[]),
    unnest(@payloads::jsonb[]),
    unnest(@partition_keys::text[]),
    unnest(@idempotency_keys::uuid[]),
    unnest(@created_ats::timestamptz[]);

