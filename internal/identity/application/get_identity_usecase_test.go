package application

import (
	"context"
	"errors"
	"testing"
	"github.com/Trivenqo/GhostKey/internal/identity/domain"

	"github.com/stretchr/testify/assert"
)

func TestGetIdentityUseCase_Execute_Success(t *testing.T) {
	mockRepo := new(domain.MockIdentityRepository)
	useCase := NewGetIdentityUseCase(mockRepo)

	expectedIdentity := &domain.Identity{
		ID:       "123",
		Metadata: map[string]interface{}{"role": "admin"},
	}

	// Tell the mock what to return when called
	mockRepo.On("FindByID", context.Background(), "123").Return(expectedIdentity, nil)

	query := GetIdentityQuery{ID: "123"}
	result, err := useCase.Execute(context.Background(), query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "123", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetIdentityUseCase_Execute_EmptyID(t *testing.T) {
	mockRepo := new(domain.MockIdentityRepository)
	useCase := NewGetIdentityUseCase(mockRepo)

	query := GetIdentityQuery{ID: ""}
	result, err := useCase.Execute(context.Background(), query)

	assert.Error(t, err)
	assert.Equal(t, "identity ID cannot be empty", err.Error())
	assert.Nil(t, result)
	// We do not expect FindByID to be called if validation fails
	mockRepo.AssertNotCalled(t, "FindByID")
}

func TestGetIdentityUseCase_Execute_NotFound(t *testing.T) {
	mockRepo := new(domain.MockIdentityRepository)
	useCase := NewGetIdentityUseCase(mockRepo)

	mockRepo.On("FindByID", context.Background(), "999").Return(nil, errors.New("identity not found"))

	query := GetIdentityQuery{ID: "999"}
	result, err := useCase.Execute(context.Background(), query)

	assert.Error(t, err)
	assert.Equal(t, "identity not found", err.Error())
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}