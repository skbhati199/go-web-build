package integration

import (
	"context"
	"sync"
)

type Integration interface {
	Connect(ctx context.Context, config map[string]interface{}) error
	Execute(ctx context.Context, params map[string]interface{}) error
	Cleanup(ctx context.Context) error
}

type IntegrationManager struct {
	integrations map[string]Integration
	mu           sync.RWMutex
}

func NewIntegrationManager() *IntegrationManager {
	return &IntegrationManager{
		integrations: make(map[string]Integration),
	}
}

func (m *IntegrationManager) Register(name string, integration Integration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.integrations[name] = integration
}
