package plugin

import (
	"context"
	"plugin"
)

type Plugin interface {
	Name() string
	Version() string
	Initialize(ctx context.Context, config map[string]interface{}) error
	Execute(ctx context.Context) error
}

type Manager struct {
	plugins map[string]Plugin
}

func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]Plugin),
	}
}

func (m *Manager) LoadPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}

	symbol, err := p.Lookup("New")
	if err != nil {
		return err
	}

	plugin := symbol.(func() Plugin)()
	m.plugins[plugin.Name()] = plugin
	return nil
}