package ports

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"github.com/Trivenqo/GhostKey/internal/identity/application"
	"github.com/Trivenqo/GhostKey/internal/identity/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIdentityHandler_GetIdentity_Success(t *testing.T) {
	// Setup Mock & UseCase
	mockRepo := new(domain.MockIdentityRepository)
	useCase := application.NewGetIdentityUseCase(mockRepo)
	handler := NewIdentityHandler(useCase)

	// Setup Fiber App
	app := fiber.New()
	app.Get("/v1/identities/:id", handler.GetIdentity)

	// Mock expectations
	expectedIdentity := &domain.Identity{
		ID:       "123",
		Metadata: map[string]interface{}{"role": "admin"},
	}
	mockRepo.On("FindByID", mock.Anything, "123").Return(expectedIdentity, nil)

	// Execute Request
	req := httptest.NewRequest("GET", "/v1/identities/123", nil)
	resp, err := app.Test(req, -1)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseBody IdentityResponse
	json.NewDecoder(resp.Body).Decode(&responseBody)

	assert.Equal(t, "123", responseBody.ID)
	assert.Equal(t, "admin", responseBody.Metadata["role"])
	mockRepo.AssertExpectations(t)
}

func TestIdentityHandler_GetIdentity_NotFound(t *testing.T) {
	mockRepo := new(domain.MockIdentityRepository)
	useCase := application.NewGetIdentityUseCase(mockRepo)
	handler := NewIdentityHandler(useCase)

	app := fiber.New()
	app.Get("/v1/identities/:id", handler.GetIdentity)

	mockRepo.On("FindByID", mock.Anything, "999").Return(nil, errors.New("identity not found"))

	req := httptest.NewRequest("GET", "/v1/identities/999", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}