package plugin

import (
	"context"
	"encoding/json"
)

type Plugin interface {
	Name() string
	Version() string
	Init(ctx context.Context, config json.RawMessage) error
	Execute(ctx context.Context, params map[string]interface{}) error
	Cleanup(ctx context.Context) error
}

type PluginManager struct {
	plugins map[string]Plugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

func (pm *PluginManager) Register(p Plugin) error {
	pm.plugins[p.Name()] = p
	return nil
}

func (pm *PluginManager) Execute(ctx context.Context, name string, params map[string]interface{}) error {
	plugin, exists := pm.plugins[name]
	if !exists {
		return ErrPluginNotFound
	}
	return plugin.Execute(ctx, params)
}