package templates

import (
	"encoding/json"
	"fmt"
)

type TemplateRegistry struct {
	Templates map[string]*Template
	BasePath  string
}

type Template struct {
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Path        string          `json:"path"`
	Config      json.RawMessage `json:"config"`
	Files       []string        `json:"files"`
	Variables   map[string]any  `json:"variables"`
	Validations map[string]any  `json:"validations"`
}

func NewRegistry(basePath string) *TemplateRegistry {
	return &TemplateRegistry{
		Templates: make(map[string]*Template),
		BasePath:  basePath,
	}
}

func (r *TemplateRegistry) Register(name string, template *Template) error {
	if _, exists := r.Templates[name]; exists {
		return fmt.Errorf("template %s already registered", name)
	}
	r.Templates[name] = template
	return nil
}
