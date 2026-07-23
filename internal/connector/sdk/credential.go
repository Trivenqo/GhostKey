package sdk

import (
	"context"
)

// CredentialManager defines how the SDK retrieves authentication secrets
// for a specific connector instance.
type CredentialManager interface {
	Get(ctx context.Context, connectorName string) (Credentials, error)
}