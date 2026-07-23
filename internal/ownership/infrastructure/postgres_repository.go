package infrastructure

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Trivenqo/GhostKey/internal/ownership/domain"
)

type PostgresOwnershipRepository struct {
	db *pgxpool.Pool
}

func NewPostgresOwnershipRepository(db *pgxpool.Pool) domain.OwnershipRepository {
	return &PostgresOwnershipRepository{db: db}
}

func (r *PostgresOwnershipRepository) Upsert(ctx context.Context, o *domain.Ownership) error {
	query := `
		INSERT INTO identity_ownership (identity_id, owner_email, team_name, department, mapping_source, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (identity_id) DO UPDATE SET
			owner_email = EXCLUDED.owner_email,
			team_name = EXCLUDED.team_name,
			department = EXCLUDED.department,
			mapping_source = EXCLUDED.mapping_source,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		o.IdentityID,
		o.OwnerEmail,
		o.TeamName,
		o.Department,
		o.MappingSource,
	).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *PostgresOwnershipRepository) GetByIdentityID(ctx context.Context, identityID uuid.UUID) (*domain.Ownership, error) {
	query := `
		SELECT id, identity_id, COALESCE(owner_email, ''), COALESCE(team_name, ''), COALESCE(department, ''), mapping_source, created_at, updated_at
		FROM identity_ownership
		WHERE identity_id = $1
	`
	o := &domain.Ownership{}
	err := r.db.QueryRow(ctx, query, identityID).Scan(
		&o.ID,
		&o.IdentityID,
		&o.OwnerEmail,
		&o.TeamName,
		&o.Department,
		&o.MappingSource,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return o, nil
}

func (r *PostgresOwnershipRepository) ListWithIdentities(ctx context.Context) ([]domain.IdentityOwnershipDTO, error) {
	query := `
		SELECT 
			i.id,
			i.arn,
			i.account_id,
			i.provider,
			i.type,
			COALESCE(o.owner_email, ''),
			COALESCE(o.team_name, ''),
			COALESCE(o.department, ''),
			COALESCE(o.mapping_source, ''),
			(o.id IS NOT NULL) as is_mapped
		FROM identities i
		LEFT JOIN identity_ownership o ON i.id = o.identity_id
		ORDER BY i.created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.IdentityOwnershipDTO
	for rows.Next() {
		var dto domain.IdentityOwnershipDTO
		var source string
		if err := rows.Scan(
			&dto.IdentityID,
			&dto.ARN,
			&dto.AccountID,
			&dto.Provider,
			&dto.Type,
			&dto.OwnerEmail,
			&dto.TeamName,
			&dto.Department,
			&source,
			&dto.IsMapped,
		); err != nil {
			return nil, err
		}
		dto.MappingSource = domain.MappingSource(source)
		results = append(results, dto)
	}

	return results, nil
}