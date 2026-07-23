package discovery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Trivenqo/GhostKey/internal/bootstrap"
	"github.com/Trivenqo/GhostKey/internal/discovery/application/usecase"
	"github.com/Trivenqo/GhostKey/internal/discovery/infrastructure/postgres"
	"github.com/Trivenqo/GhostKey/internal/discovery/interfaces/http"
)

// Register wires up the entire Discovery module and attaches its HTTP routes.
func Register(container *bootstrap.Container, router fiber.Router) {
	// 1. Infrastructure Layer
	repo := postgres.NewIdentityRepository(container.DB)

	// 2. Application Layer
	listUseCase := usecase.NewListIdentitiesUseCase(repo)

	// 3. Interfaces Layer (HTTP)
	handler := http.NewHandler(listUseCase)

	// 4. Register Routes
	discoveryGroup := router.Group("/identities")
	discoveryGroup.Get("/", handler.ListIdentities)
}