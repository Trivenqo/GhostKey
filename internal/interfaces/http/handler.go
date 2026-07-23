package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/Trivenqo/GhostKey/internal/discovery/application/usecase"
	"github.com/Trivenqo/GhostKey/internal/shared/canonical"
)

// IdentityResponse is the DTO (Data Transfer Object) for the API.
// It ensures we don't accidentally leak internal domain structures.
type IdentityResponse struct {
	ID          string            `json:"id"`
	Provider    string            `json:"provider"`
	Type        string            `json:"type"`
	DisplayName string            `json:"display_name"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   string            `json:"created_at"`
}

type Handler struct {
	listUseCase *usecase.ListIdentitiesUseCase
}

func NewHandler(listUseCase *usecase.ListIdentitiesUseCase) *Handler {
	return &Handler{
		listUseCase: listUseCase,
	}
}

// ListIdentities handles GET /v1/identities
func (h *Handler) ListIdentities(c *fiber.Ctx) error {
	// 1. Parse Query Parameters
	provider := c.Query("provider", "")
	
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	query := usecase.ListIdentitiesQuery{
		Provider: provider,
		Limit:    limit,
		Offset:   offset,
	}

	// 2. Call the Application Layer (Business Logic)
	identities, err := h.listUseCase.Execute(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 3. Map Domain Entities to HTTP DTOs
	var response []IdentityResponse
	for _, id := range identities {
		response = append(response, IdentityResponse{
			ID:          id.ID,
			Provider:    id.ExternalRef.Provider,
			Type:        string(id.Type),
			DisplayName: id.DisplayName,
			Metadata:    id.Metadata,
			CreatedAt:   id.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	// Ensure we return an empty array [] instead of null if there are no results
	if response == nil {
		response = make([]IdentityResponse, 0)
	}

	return c.JSON(fiber.Map{
		"data": response,
		"meta": fiber.Map{
			"limit":  limit,
			"offset": offset,
		},
	})
}