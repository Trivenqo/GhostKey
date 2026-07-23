CREATE TABLE IF NOT EXISTS identities (
    id VARCHAR(255) PRIMARY KEY,
    provider VARCHAR(100) NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    account_ref VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_used_at TIMESTAMP WITH TIME ZONE,
    
    -- System timestamps
    sys_created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    sys_updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Deduplication constraint: an identity is uniquely identified by where it came from
    CONSTRAINT uq_identities_provider_ext_acct UNIQUE (provider, external_id, account_ref)
);