package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Trivenqo/GhostKey/internal/identity/domain"
)

type PostgresIdentityRepository struct {
	DB *sql.DB
}

func (r *PostgresIdentityRepository) FindByID(ctx context.Context, id string) (*domain.Identity, error) {
	query := `SELECT id, metadata FROM identities WHERE id = $1`
	row := r.DB.QueryRowContext(ctx, query, id)

	var identity domain.Identity
	var metadataJSON []byte

	// Safely scan the row
	err := row.Scan(&identity.ID, &metadataJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("identity not found")
		}
		return nil, err
	}

	// Parse JSON metadata back to map
	if len(metadataJSON) > 0 {
		err = json.Unmarshal(metadataJSON, &identity.Metadata)
		if err != nil {
			return nil, err
		}
	}

	return &identity, nil
}