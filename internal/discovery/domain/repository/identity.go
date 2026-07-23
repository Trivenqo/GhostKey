package repository

import (
	"context"

	"github.com/Trivenqo/GhostKey/internal/shared/canonical"
)

type IdentityRepository interface {
	Upsert(ctx context.Context, identity canonical.Identity) error
	
	// NEW: List retrieves a paginated list of identities, optionally filtered by provider.
	List(ctx context.Context, providerFilter string, limit, offset int) ([]canonical.Identity, error)
}