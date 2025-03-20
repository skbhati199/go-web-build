package serverless

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/sonukumar/go-web-build/internal/template-engine/engine"
	"github.com/sonukumar/go-web-build/internal/template-engine/validation"
)

type Deployer struct {
	provider  Provider
	config    Config
	template  *engine.TemplateEngine
	validator *validation.TemplateValidator
}

type Config struct {
	Provider     string
	Region       string
	Functions    []FunctionConfig
	Triggers     []TriggerConfig
	TemplatesDir string
}

type FunctionConfig struct {
	Name        string
	Runtime     string
	Memory      int
	Timeout     int
	Environment map[string]string
	Template    string // Template name for the function
}

type TriggerConfig struct {
	Type       string
	Function   string
	Properties map[string]interface{}
}

type DeploymentResult struct {
	Resources []Resource
	Endpoint  string
	Version   string
}

func NewDeployer(config Config) *Deployer {
	return &Deployer{
		provider:  selectProvider(config.Provider),
		config:    config,
		template:  engine.NewTemplateEngine(config.TemplatesDir),
		validator: validation.NewTemplateValidator(config.TemplatesDir),
	}
}

func (d *Deployer) Deploy(ctx context.Context) (*DeploymentResult, error) {
	// Validate and process templates first
	if err := d.processTemplates(ctx); err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	resources := d.prepareResources()
	return d.provider.Deploy(ctx, resources)
}

func (d *Deployer) processTemplates(ctx context.Context) error {
	for _, function := range d.config.Functions {
		if function.Template != "" {
			if err := d.validator.ValidateTemplate(function.Template); err != nil {
				return fmt.Errorf("invalid template for function %s: %w", function.Name, err)
			}

			data := &engine.TemplateData{
				ProjectName: function.Name,
				Framework:   "serverless",
				Language:    function.Runtime,
				Version:     "1.0.0",
				Configuration: map[string]interface{}{
					"memory":      function.Memory,
					"timeout":     function.Timeout,
					"environment": function.Environment,
				},
			}

			outputDir := filepath.Join("build", "functions", function.Name)
			if err := d.template.Generate(function.Template, data, outputDir); err != nil {
				return fmt.Errorf("failed to generate function %s: %w", function.Name, err)
			}
		}
	}
	return nil
}

func (d *Deployer) prepareResources() []Resource {
	// Implementation for resource preparation
	return nil
}

type Resource struct {
	Type       string
	Name       string
	Properties map[string]interface{}
}

type Provider interface {
	Deploy(context.Context, []Resource) (*DeploymentResult, error)
}

func selectProvider(providerType string) Provider {
	// Provider selection implementation
	return nil
}
