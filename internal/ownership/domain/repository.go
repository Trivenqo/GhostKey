package domain

import (
	"context"

	"github.com/google/uuid"
)

type OwnershipRepository interface {
	Upsert(ctx context.Context, ownership *Ownership) error
	GetByIdentityID(ctx context.Context, identityID uuid.UUID) (*Ownership, error)
	ListWithIdentities(ctx context.Context) ([]IdentityOwnershipDTO, error)
}