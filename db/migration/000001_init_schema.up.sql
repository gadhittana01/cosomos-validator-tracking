CREATE TABLE IF NOT EXISTS delegation_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    validator_address TEXT NOT NULL,
    delegator_address TEXT NOT NULL,
    amount_uatom BIGINT NOT NULL,
    change_uatom BIGINT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS daily_aggregates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    validator_address TEXT NOT NULL,
    delegator_address TEXT NOT NULL,
    date DATE NOT NULL,
    total_amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);