package presentation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"ghostkey/internal/ownership/application"
)

type OwnershipHandler struct {
	useCase *application.OwnershipUseCase
}

func NewOwnershipHandler(useCase *application.OwnershipUseCase) *OwnershipHandler {
	return &OwnershipHandler{useCase: useCase}
}

func (h *OwnershipHandler) RegisterRoutes(router fiber.Router) {
	ownershipGroup := router.Group("/ownership")
	ownershipGroup.Get("/", h.ListOwnerships)
	ownershipGroup.Post("/", h.AssignOwnership)

	identitiesGroup := router.Group("/identities")
	identitiesGroup.Get("/:id/ownership", h.GetIdentityOwnership)
}

func (h *OwnershipHandler) ListOwnerships(c *fiber.Ctx) error {
	results, err := h.useCase.ListOwnerships(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"data": results})
}

func (h *OwnershipHandler) AssignOwnership(c *fiber.Ctx) error {
	var input application.AssignOwnershipInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}

	if input.IdentityID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "identity_id is required",
		})
	}

	ownership, err := h.useCase.AssignManualOwnership(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Ownership assigned successfully",
		"data":    ownership,
	})
}

func (h *OwnershipHandler) GetIdentityOwnership(c *fiber.Ctx) error {
	idParam := c.Params("id")
	identityID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid identity ID",
		})
	}

	ownership, err := h.useCase.GetIdentityOwnership(c.Context(), identityID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if ownership == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No ownership mapping found for this identity",
		})
	}

	return c.JSON(fiber.Map{"data": ownership})
}