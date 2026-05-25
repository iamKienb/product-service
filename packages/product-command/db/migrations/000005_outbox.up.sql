CREATE TABLE outbox (
    id UUID PRIMARY KEY,
    aggregate_id UUID NOT NULL,   
    aggregate_type TEXT NOT NULL,
    event_type TEXT NOT NULL,      
    payload JSONB NOT NULL DEFAULT '{}',
    partition_key TEXT NOT NULL,
    idempotency_key UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    CONSTRAINT uk_outbox_idempotency UNIQUE (idempotency_key, event_type)
);

CREATE INDEX idx_outbox_created_at ON outbox(created_at);
