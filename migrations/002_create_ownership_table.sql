CREATE TABLE IF NOT EXISTS identity_ownership (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identity_id UUID NOT NULL UNIQUE REFERENCES identities(id) ON DELETE CASCADE,
    owner_email VARCHAR(255),
    team_name VARCHAR(255),
    department VARCHAR(255),
    mapping_source VARCHAR(50) NOT NULL DEFAULT 'MANUAL', -- 'AUTO' or 'MANUAL'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ownership_identity_id ON identity_ownership(identity_id);
CREATE INDEX IF NOT EXISTS idx_ownership_team_name ON identity_ownership(team_name);