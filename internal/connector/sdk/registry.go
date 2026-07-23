package sdk

import (
	"fmt"
	"sync"
)

// Registry holds all active connector implementations.
type Registry struct {
	mu         sync.RWMutex
	connectors map[string]Connector
}

// NewRegistry creates a new empty connector registry.
func NewRegistry() *Registry {
	return &Registry{
		connectors: make(map[string]Connector),
	}
}

// Register adds a new connector to the registry.
func (r *Registry) Register(c Connector) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := c.Name()
	if _, exists := r.connectors[name]; exists {
		return fmt.Errorf("connector %s is already registered", name)
	}

	r.connectors[name] = c
	return nil
}

// Get retrieves a connector by name.
func (r *Registry) Get(name string) (Connector, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, exists := r.connectors[name]
	if !exists {
		return nil, fmt.Errorf("connector %s not found", name)
	}

	return c, nil
}