CREATE TABLE IF NOT EXISTS payflow_wallet_transactions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES payflow_wallets(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    type VARCHAR(10) NOT NULL, 
    amount NUMERIC(14,2) NOT NULL,
    reference TEXT,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX Idx_wallet_tx_user_id ON payflow_wallet_transactions(user_id);
CREATE INDEX Idx_wallet_tx_wallet_id ON payflow_wallet_transactions(wallet_id);
CREATE INDEX Idx_wallet_tx_created_at ON payflow_wallet_transactions(created_at DESC);