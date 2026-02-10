CREATE TABLE IF NOT EXISTS payflow_idempotency_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key TEXT NOT NULL,
    user_id UUID NOT NULL,
    payment_reference VARCHAR(100),
    request_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (idempotency_key, user_id)
);
