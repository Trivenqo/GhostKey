package domain

import "context"

// Identity represents the canonical domain entity
type Identity struct {
	ID       string
	Metadata map[string]interface{}
}

// IdentityRepository defines the interface extension
type IdentityRepository interface {
	FindByID(ctx context.Context, id string) (*Identity, error)
}