package application

import (
	"context"
	"errors"
	"github.com/Trivenqo/GhostKey/internal/identity/domain"
)

// GetIdentityQuery encapsulates the input parameters
type GetIdentityQuery struct {
	ID string
}

type GetIdentityUseCase struct {
	repo domain.IdentityRepository
}

func NewGetIdentityUseCase(repo domain.IdentityRepository) *GetIdentityUseCase {
	return &GetIdentityUseCase{repo: repo}
}

func (uc *GetIdentityUseCase) Execute(ctx context.Context, q GetIdentityQuery) (*domain.Identity, error) {
	// Core business validation rule
	if q.ID == "" {
		return nil, errors.New("identity ID cannot be empty")
	}

	// Delegate to repository adapter
	return uc.repo.FindByID(ctx, q.ID)
}