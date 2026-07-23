package application

import (
	"context"

	"github.com/google/uuid"

	"ghostkey/internal/ownership/domain"
)

type AssignOwnershipInput struct {
	IdentityID uuid.UUID `json:"identity_id"`
	OwnerEmail string    `json:"owner_email"`
	TeamName   string    `json:"team_name"`
	Department string    `json:"department"`
}

type OwnershipUseCase struct {
	repo domain.OwnershipRepository
}

func NewOwnershipUseCase(repo domain.OwnershipRepository) *OwnershipUseCase {
	return &OwnershipUseCase{repo: repo}
}

func (uc *OwnershipUseCase) AssignManualOwnership(ctx context.Context, input AssignOwnershipInput) (*domain.Ownership, error) {
	ownership := &domain.Ownership{
		IdentityID:    input.IdentityID,
		OwnerEmail:    input.OwnerEmail,
		TeamName:      input.TeamName,
		Department:    input.Department,
		MappingSource: domain.SourceManual,
	}

	err := uc.repo.Upsert(ctx, ownership)
	if err != nil {
		return nil, err
	}

	return ownership, nil
}

func (uc *OwnershipUseCase) AssignAutoOwnership(ctx context.Context, identityID uuid.UUID, ownerEmail, teamName, department string) error {
	// Skip auto-assignment if already manually mapped
	existing, err := uc.repo.GetByIdentityID(ctx, identityID)
	if err == nil && existing != nil && existing.MappingSource == domain.SourceManual {
		return nil
	}

	ownership := &domain.Ownership{
		IdentityID:    identityID,
		OwnerEmail:    ownerEmail,
		TeamName:      teamName,
		Department:    department,
		MappingSource: domain.SourceAuto,
	}

	return uc.repo.Upsert(ctx, ownership)
}

func (uc *OwnershipUseCase) GetIdentityOwnership(ctx context.Context, identityID uuid.UUID) (*domain.Ownership, error) {
	return uc.repo.GetByIdentityID(ctx, identityID)
}

func (uc *OwnershipUseCase) ListOwnerships(ctx context.Context) ([]domain.IdentityOwnershipDTO, error) {
	return uc.repo.ListWithIdentities(ctx)
}