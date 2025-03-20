package template

import (
	"fmt"
	"path/filepath"

	"github.com/sonukumar/go-web-build/internal/template-engine/engine"
)

type ReactOptions struct {
	ProjectName   string
	UseTypeScript bool
	Features      ReactFeatures
}

type ReactFeatures struct {
	Router       bool
	StateManager string // redux, mobx, zustand
	Styling      string // css-modules, styled-components, tailwind
	Testing      bool
}

type ReactTemplate struct {
	engine *engine.TemplateEngine
	cache  map[string]*engine.TemplateData
}

func NewReactTemplate(templatesDir string) *ReactTemplate {
	return &ReactTemplate{
		engine: engine.NewTemplateEngine(templatesDir),
		cache:  make(map[string]*engine.TemplateData),
	}
}

func (t *ReactTemplate) Generate(opts *ReactOptions) error {
	config := map[string]interface{}{
		"router":       opts.Features.Router,
		"stateManager": opts.Features.StateManager,
		"styling":      opts.Features.Styling,
		"testing":      opts.Features.Testing,
	}

	// Generate base React project
	if err := t.engine.GenerateReactProject(opts.ProjectName, opts.UseTypeScript, config); err != nil {
		return fmt.Errorf("failed to generate React project: %w", err)
	}

	// Add additional features
	if err := t.addFeatures(opts); err != nil {
		return fmt.Errorf("failed to add features: %w", err)
	}

	return nil
}

func (t *ReactTemplate) addFeatures(opts *ReactOptions) error {
	projectPath := filepath.Clean(opts.ProjectName)

	if opts.Features.Router {
		if err := t.engine.GenerateReactRouterProject(projectPath, nil); err != nil {
			return fmt.Errorf("failed to add router: %w", err)
		}
	}

	if opts.Features.StateManager == "redux" {
		if err := t.engine.GenerateReactReduxProject(projectPath, nil); err != nil {
			return fmt.Errorf("failed to add Redux: %w", err)
		}
	}

	if opts.Features.Testing {
		if err := t.engine.AddTestingSetup(projectPath, nil); err != nil {
			return fmt.Errorf("failed to add testing setup: %w", err)
		}
	}

	return nil
}

func (t *ReactTemplate) ValidateOptions(opts *ReactOptions) error {
	if opts.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}

	if opts.Features.StateManager != "" {
		validManagers := map[string]bool{
			"redux":   true,
			"mobx":    true,
			"zustand": true,
		}
		if !validManagers[opts.Features.StateManager] {
			return fmt.Errorf("unsupported state manager: %s", opts.Features.StateManager)
		}
	}

	if opts.Features.Styling != "" {
		validStyles := map[string]bool{
			"css-modules":       true,
			"styled-components": true,
			"tailwind":          true,
		}
		if !validStyles[opts.Features.Styling] {
			return fmt.Errorf("unsupported styling option: %s", opts.Features.Styling)
		}
	}

	return nil
}
