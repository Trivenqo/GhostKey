package usecase

import (
	"context"
	"fmt"

	"github.com/Trivenqo/GhostKey/internal/discovery/application/port"
	"github.com/Trivenqo/GhostKey/internal/discovery/domain/repository"
	"github.com/Trivenqo/GhostKey/internal/shared/canonical"
)

// RegisterIdentityUseCase orchestrates storing a newly discovered identity
// and broadcasting its arrival to other modules.
type RegisterIdentityUseCase struct {
	repo      repository.IdentityRepository
	publisher port.EventPublisher
}

func NewRegisterIdentityUseCase(repo repository.IdentityRepository, pub port.EventPublisher) *RegisterIdentityUseCase {
	return &RegisterIdentityUseCase{
		repo:      repo,
		publisher: pub,
	}
}

// Execute runs the discovery registration workflow.
func (uc *RegisterIdentityUseCase) Execute(ctx context.Context, identity canonical.Identity) error {
	// 1. Persist the identity (insert or update)
	if err := uc.repo.Upsert(ctx, identity); err != nil {
		return fmt.Errorf("failed to persist identity: %w", err)
	}

	// 2. Announce discovery to other modules (Ownership/Risk)
	if err := uc.publisher.PublishIdentityDiscovered(ctx, identity); err != nil {
		return fmt.Errorf("failed to publish discovery event: %w", err)
	}

	return nil
}