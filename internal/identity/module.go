package identity

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/Trivenqo/GhostKey/internal/identity/adapters"
	"github.com/Trivenqo/GhostKey/internal/identity/application"
	"github.com/Trivenqo/GhostKey/internal/identity/ports"
)

// RegisterModule acts as the discovery and wiring module
func RegisterModule(router fiber.Router, db *sql.DB) {
	// Dependency Injection
	repo := &adapters.PostgresIdentityRepository{DB: db}
	useCase := application.NewGetIdentityUseCase(repo)
	handler := ports.NewIdentityHandler(useCase)

	// Routing
	group := router.Group("/identities")
	group.Get("/:id", handler.GetIdentity)
}