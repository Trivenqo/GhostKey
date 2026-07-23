package domain

import (
	"context"
	"github.com/stretchr/testify/mock"
)

// MockIdentityRepository is a mock implementation of IdentityRepository
type MockIdentityRepository struct {
	mock.Mock
}

// FindByID provides a mock function with given fields: ctx, id
func (m *MockIdentityRepository) FindByID(ctx context.Context, id string) (*Identity, error) {
	args := m.Called(ctx, id)
	
	var identity *Identity
	if args.Get(0) != nil {
		identity = args.Get(0).(*Identity)
	}
	
	return identity, args.Error(1)
}