package port

import (
	"context"

	"github.com/Trivenqo/GhostKey/internal/shared/canonical"
)

// EventPublisher defines how the application pushes domain events to the outside world.
type EventPublisher interface {
	PublishIdentityDiscovered(ctx context.Context, identity canonical.Identity) error
}