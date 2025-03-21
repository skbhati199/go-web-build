package templates

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// Update how templates are registered to match your directory structure
func (r *TemplateRegistry) LoadTemplates() error {
	// Instead of looking for templates like "react-typescript"
	// Look for them in nested directories like "react/typescript"

	// Example implementation:
	frameworks, err := os.ReadDir(r.BasePath)
	if err != nil {
		return err
	}

	for _, framework := range frameworks {
		if !framework.IsDir() {
			continue
		}

		variants, err := os.ReadDir(filepath.Join(r.BasePath, framework.Name()))
		if err != nil {
			continue
		}

		for _, variant := range variants {
			if !variant.IsDir() {
				continue
			}

			templatePath := filepath.Join(r.BasePath, framework.Name(), variant.Name())
			// Load template from this path...

			// Example:
			template := &Template{
				Name:    fmt.Sprintf("%s/%s", framework.Name(), variant.Name()),
				Version: "1.0.0",
				Path:    templatePath,
				// ... other fields
			}

			if err := r.Register(template.Name, template); err != nil {
				return err
			}
		}
	}

	return nil
}
