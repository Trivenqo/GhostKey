package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Trivenqo/GhostKey/internal/shared/canonical"
)

type IdentityRepository struct {
	pool *pgxpool.Pool
}

func NewIdentityRepository(pool *pgxpool.Pool) *IdentityRepository {
	return &IdentityRepository{pool: pool}
}

func (r *IdentityRepository) Upsert(ctx context.Context, identity canonical.Identity) error {
	// Convert the metadata map to JSONB for Postgres
	metadataJSON, err := json.Marshal(identity.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO identities (
			id, provider, external_id, account_ref, type, display_name, metadata, created_at, last_used_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		ON CONFLICT (provider, external_id, account_ref) 
		DO UPDATE SET 
			type = EXCLUDED.type,
			display_name = EXCLUDED.display_name,
			metadata = EXCLUDED.metadata,
			last_used_at = EXCLUDED.last_used_at,
			sys_updated_at = NOW();
	`

	_, err = r.pool.Exec(ctx, query,
		identity.ID,
		identity.ExternalRef.Provider,
		identity.ExternalRef.ExternalID,
		identity.ExternalRef.AccountRef,
		string(identity.Type),
		identity.DisplayName,
		metadataJSON,
		identity.CreatedAt,
		identity.LastUsedAt,
	)

	if err != nil {
		return fmt.Errorf("postgres upsert failed: %w", err)
	}

	return nil
}

func (r *IdentityRepository) List(ctx context.Context, providerFilter string, limit, offset int) ([]canonical.Identity, error) {
	query := `
		SELECT id, provider, external_id, account_ref, type, display_name, metadata, created_at, last_used_at
		FROM identities
		WHERE ($1 = '' OR provider = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, providerFilter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("postgres list failed: %w", err)
	}
	defer rows.Close()

	var identities []canonical.Identity

	for rows.Next() {
		var id, provider, extID, acctRef, typ, displayName string
		var metadataJSON []byte
		var createdAt time.Time
		var lastUsedAt *time.Time

		err := rows.Scan(&id, &provider, &extID, &acctRef, &typ, &displayName, &metadataJSON, &createdAt, &lastUsedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var metadata map[string]string
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		identities = append(identities, canonical.Identity{
			ID: id,
			ExternalRef: canonical.ProviderRef{
				Provider:   provider,
				ExternalID: extID,
				AccountRef: acctRef,
			},
			Type:        canonical.IdentityType(typ),
			DisplayName: displayName,
			Metadata:    metadata,
			CreatedAt:   createdAt,
			LastUsedAt:  lastUsedAt,
		})
	}

	return identities, nil
}