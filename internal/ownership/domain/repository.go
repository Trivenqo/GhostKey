package domain

import (
	"context"
)

type OwnershipRepository interface {
	Upsert(ctx context.Context, ownership *Ownership) error
	GetByIdentityID(ctx context.Context, identityID string) (*Ownership, error) // Changed to string
	ListWithIdentities(ctx context.Context) ([]IdentityOwnershipDTO, error)
}