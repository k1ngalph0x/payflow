CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS payflow_payments(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    amount NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL,
    reference VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_payflow_payments_user_id ON payflow_payments(user_id)
CREATE INDEX idx_payflow_merchant_id ON payflow_payments(merchant_id)