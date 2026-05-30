CREATE TYPE plans_tier_enum AS ENUM ('bronze', 'silver', 'gold', 'platinum');

CREATE TABLE IF NOT EXISTS plans
(
    id                UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name              VARCHAR(100)    NOT NULL,
    provider          VARCHAR(70)     NOT NULL,
    tier              plans_tier_enum NOT NULL,
    monthly_premium   DECIMAL(10, 2)  NOT NULL,
    deductible        DECIMAL(10, 2),
    out_of_pocket_max DECIMAL(10, 2),
    state             VARCHAR(60)     NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_plans_state ON plans(state);