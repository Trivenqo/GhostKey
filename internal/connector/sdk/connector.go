package sdk

import (
	"context"
	"time"

	"github.com/Trivenqo/GhostKey/internal/shared/canonical"
)

// RawRecord represents an untyped JSON payload directly from a provider API.
type RawRecord map[string]interface{}

// Credentials represents the authentication material needed for a connector.
type Credentials map[string]string

// Cursor tracks pagination and incremental polling state.
type Cursor struct {
	Token     string
	UpdatedAt time.Time
}

// Page represents a single page of results from a provider API.
type Page struct {
	Items      []RawRecord
	NextCursor Cursor
	HasMore    bool
}

// HealthStatus represents the current operational state of a connector.
type HealthStatus struct {
	Healthy    bool
	LastSyncAt time.Time
	ErrorRate  float64
	Message    string
}

// Connector is the interface every provider (AWS, GitHub, etc.) must implement.
type Connector interface {
	Name() string
	Authenticate(ctx context.Context, creds Credentials) error
	Discover(ctx context.Context, cursor Cursor) (Page, error)
	SupportsWebhook() bool
	HandleWebhook(ctx context.Context, payload []byte) ([]canonical.Identity, error)
	HealthCheck(ctx context.Context) HealthStatus
}

// Normalizer converts provider-specific raw JSON into our canonical model.
type Normalizer interface {
	Normalize(raw RawRecord) (canonical.Identity, error)
}