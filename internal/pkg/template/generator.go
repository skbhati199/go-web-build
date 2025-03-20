package template

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/sonukumar/go-web-build/internal/template-engine/engine"
	"github.com/sonukumar/go-web-build/internal/template-engine/version"
)

type TemplateGenerator struct {
	engine       *engine.TemplateEngine
	versions     *version.VersionManager
	templatesDir string
}

type Config struct {
	TemplatesDir string
	Version      string
	CacheTimeout int
}

func NewTemplateGenerator(config *Config) *TemplateGenerator {
	return &TemplateGenerator{
		engine:       engine.NewTemplateEngine(config.TemplatesDir),
		versions:     version.NewVersionManager(config.TemplatesDir),
		templatesDir: config.TemplatesDir,
	}
}

func (g *TemplateGenerator) GenerateReactProject(ctx context.Context, opts *ReactOptions) error {
	// Convert ReactOptions to template engine configuration
	config := map[string]interface{}{
		"router":       opts.Features.Router,
		"stateManager": opts.Features.StateManager,
		"styling":      opts.Features.Styling,
		"testing":      opts.Features.Testing,
	}

	// Generate base React project
	if err := g.engine.GenerateReactProject(opts.ProjectName, opts.UseTypeScript, config); err != nil {
		return fmt.Errorf("failed to generate React project: %w", err)
	}

	// Add additional features based on options
	if err := g.addFeatures(ctx, opts); err != nil {
		return fmt.Errorf("failed to add features: %w", err)
	}

	return nil
}

func (g *TemplateGenerator) addFeatures(ctx context.Context, opts *ReactOptions) error {
	projectPath := filepath.Clean(opts.ProjectName)

	if opts.Features.Router {
		if err := g.engine.GenerateReactRouterProject(projectPath, nil); err != nil {
			return fmt.Errorf("failed to add router: %w", err)
		}
	}

	if opts.Features.StateManager == "redux" {
		if err := g.engine.GenerateReactReduxProject(projectPath, nil); err != nil {
			return fmt.Errorf("failed to add Redux: %w", err)
		}
	}

	if opts.Features.Testing {
		if err := g.engine.AddTestingSetup(projectPath, nil); err != nil {
			return fmt.Errorf("failed to add testing setup: %w", err)
		}
	}

	return nil
}

func (g *TemplateGenerator) generateFromTemplate(ctx context.Context, templateName string, data *engine.TemplateData) error {
	// Get latest version if not specified
	templateVersion, err := g.versions.GetLatestVersion(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template version: %w", err)
	}

	data.Version = templateVersion.Version.String()
	return g.engine.Generate(templateName, data, data.ProjectName)
}
