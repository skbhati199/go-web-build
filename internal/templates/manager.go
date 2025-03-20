package templates

import (
	"fmt"
	"sync"
)

type Manager struct {
	registry *TemplateRegistry
	cache    sync.Map
}

func NewManager(basePath string) *Manager {
	return &Manager{
		registry: NewRegistry(basePath),
	}
}

func (m *Manager) Load(name string) (*Template, error) {
	// Check cache first
	if tmpl, ok := m.cache.Load(name); ok {
		return tmpl.(*Template), nil
	}

	// Load from registry
	tmpl, ok := m.registry.Templates[name]
	if !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}

	// Cache the template
	m.cache.Store(name, tmpl)
	return tmpl, nil
}
