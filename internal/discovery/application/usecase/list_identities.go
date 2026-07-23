package usecase

import (
	"context"
	"fmt"

	"github.com/Trivenqo/GhostKey/internal/discovery/domain/repository"
	"github.com/Trivenqo/GhostKey/internal/shared/canonical"
)

type ListIdentitiesQuery struct {
	Provider string
	Limit    int
	Offset   int
}

type ListIdentitiesUseCase struct {
	repo repository.IdentityRepository
}

func NewListIdentitiesUseCase(repo repository.IdentityRepository) *ListIdentitiesUseCase {
	return &ListIdentitiesUseCase{repo: repo}
}

func (uc *ListIdentitiesUseCase) Execute(ctx context.Context, query ListIdentitiesQuery) ([]canonical.Identity, error) {
	// Apply safe defaults
	if query.Limit <= 0 || query.Limit > 100 {
		query.Limit = 50
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	identities, err := uc.repo.List(ctx, query.Provider, query.Limit, query.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list identities: %w", err)
	}

	return identities, nil
}