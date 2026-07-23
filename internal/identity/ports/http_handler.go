package ports

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Trivenqo/GhostKey/internal/identity/application"
)

// IdentityResponse maps the domain entity to a public structure
type IdentityResponse struct {
	ID       string                 `json:"id"`
	Metadata map[string]interface{} `json:"metadata"`
}

type IdentityHandler struct {
	useCase *application.GetIdentityUseCase
}

func NewIdentityHandler(useCase *application.GetIdentityUseCase) *IdentityHandler {
	return &IdentityHandler{useCase: useCase}
}

func (h *IdentityHandler) GetIdentity(c *fiber.Ctx) error {
	// Extract dynamic path parameter
	id := c.Params("id")

	query := application.GetIdentityQuery{ID: id}
	
	// Execute use case
	identity, err := h.useCase.Execute(c.Context(), query)
	if err != nil {
		if err.Error() == "identity not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not Found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	// DTO Mapping
	resp := IdentityResponse{
		ID:       identity.ID,
		Metadata: identity.Metadata,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}