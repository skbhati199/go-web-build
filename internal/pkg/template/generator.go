package template

import (
	"context"
	"fmt"
	"os"

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
	// Ensure project directory exists
	if err := os.MkdirAll(opts.ProjectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Convert ReactOptions to template engine configuration
	config := map[string]interface{}{
		"router":       opts.Features.Router,
		"stateManager": opts.Features.StateManager,
		"styling":      opts.Features.Styling,
		"testing":      opts.Features.Testing,
	}

	// Determine template name based on language
	templateName := "react-javascript"
	if opts.UseTypeScript {
		templateName = "react-typescript"
	}

	// Generate base React project

	// opts.UseTypeScript ? "typescript" : "javascript",
	var language string
	if opts.UseTypeScript {
		language = "typescript"
	} else {
		language = "javascript"
	}

	templateData := &engine.TemplateData{
		ProjectName:   opts.ProjectName,
		Framework:     "react",
		Language:      language,
		Configuration: config,
	}

	if err := g.generateFromTemplate(ctx, templateName, templateData); err != nil {
		return fmt.Errorf("failed to generate React project: %w", err)
	}

	// Add additional features based on options
	if err := g.addFeatures(ctx, opts); err != nil {
		return fmt.Errorf("failed to add features: %w", err)
	}

	return nil
}

func (g *TemplateGenerator) addFeatures(ctx context.Context, opts *ReactOptions) error {
	// projectPath := filepath.Clean(opts.ProjectName)

	if opts.Features.Router {
		var language string
		if opts.UseTypeScript {
			language = "typescript"
		} else {
			language = "javascript"
		}
		routerData := &engine.TemplateData{
			ProjectName: opts.ProjectName,
			Framework:   "react",
			Language:    language,
			Configuration: map[string]interface{}{
				"routerVersion": "6.14.0",
			},
		}

		if err := g.generateFromTemplate(ctx, "react-router", routerData); err != nil {
			return fmt.Errorf("failed to add router: %w", err)
		}
	}

	if opts.Features.StateManager == "redux" {
		var language string
		if opts.UseTypeScript {
			language = "typescript"
		} else {
			language = "javascript"
		}
		reduxData := &engine.TemplateData{
			ProjectName: opts.ProjectName,
			Framework:   "react",
			Language:    language,
			Configuration: map[string]interface{}{
				"reduxVersion":        "4.2.1",
				"reduxToolkitVersion": "1.9.5",
			},
		}

		if err := g.generateFromTemplate(ctx, "react-redux", reduxData); err != nil {
			return fmt.Errorf("failed to add Redux: %w", err)
		}
	}

	if opts.Features.Testing {
		var language string
		if opts.UseTypeScript {
			language = "typescript"
		} else {
			language = "javascript"
		}
		testingData := &engine.TemplateData{
			ProjectName: opts.ProjectName,
			Framework:   "react",
			Language:    language,
			Configuration: map[string]interface{}{
				"testingLibrary": true,
				"jest":           true,
			},
		}

		if err := g.generateFromTemplate(ctx, "react-testing", testingData); err != nil {
			return fmt.Errorf("failed to add testing setup: %w", err)
		}
	}

	// Add build scripts and configuration
	var language string
	if opts.UseTypeScript {
		language = "typescript"
	} else {
		language = "javascript"
	}
	buildData := &engine.TemplateData{
		ProjectName: opts.ProjectName,
		Framework:   "react",
		Language:    language,
		Configuration: map[string]interface{}{
			"optimization": true,
			"minify":       true,
			"sourceMaps":   opts.Features.Styling,
		},
	}

	if err := g.generateFromTemplate(ctx, "react-build", buildData); err != nil {
		return fmt.Errorf("failed to add build configuration: %w", err)
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
