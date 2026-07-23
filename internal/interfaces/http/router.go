package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/Trivenqo/GhostKey/internal/bootstrap"
	"github.com/Trivenqo/GhostKey/internal/interfaces/http/health"
)

func NewRouter(container *bootstrap.Container) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(recover.New())

	v1 := app.Group("/v1")

	// Register Core Endpoints correctly on the v1 group!
	healthHandler := health.NewHandler(container)
	v1.Get("/healthz", healthHandler.Healthz)
	v1.Get("/readyz", healthHandler.Readyz)

	return app
}