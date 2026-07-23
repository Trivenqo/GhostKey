package health

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Trivenqo/GhostKey/internal/bootstrap"
)

type Handler struct {
	container *bootstrap.Container
}

func NewHandler(c *bootstrap.Container) *Handler {
	return &Handler{container: c}
}

// Liveness check
func (h *Handler) Healthz(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// Readiness check (verifies infrastructure connections)
func (h *Handler) Readyz(c *fiber.Ctx) error {
	ctx := c.Context()

	if err := h.container.DB.Ping(ctx); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "database unavailable"})
	}

	if err := h.container.Redis.Ping(ctx).Err(); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "redis unavailable"})
	}

	if err := h.container.Kafka.Ping(ctx); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "kafka unavailable"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ready"})
}