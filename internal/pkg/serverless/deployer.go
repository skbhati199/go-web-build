package serverless

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/skbhati199/go-web-build/internal/template-engine/engine"
	"github.com/skbhati199/go-web-build/internal/template-engine/validation"
)

// Provider defines the interface for serverless deployment providers
type Provider interface {
	Configure(ctx context.Context, config ProviderConfig) error
	Deploy(ctx context.Context, resources []Resource) (*DeploymentResult, error)
	Remove(ctx context.Context, name string) error
	GetStatus(ctx context.Context, name string) (*DeploymentStatus, error)
}

// ProviderConfig contains configuration for serverless providers
type ProviderConfig struct {
	Type         string
	Region       string
	Credentials  map[string]string
	FunctionOpts FunctionOptions
}

// FunctionOptions contains configuration options for serverless functions
type FunctionOptions struct {
	Runtime     string
	MemorySize  int
	Timeout     int
	Environment map[string]string
}

// Config contains the deployer configuration
type Config struct {
	ProviderConfig ProviderConfig
	Functions      []FunctionConfig
	Triggers       []TriggerConfig
	TemplatesDir   string
}

// FunctionConfig defines a serverless function
type FunctionConfig struct {
	Name        string
	Runtime     string
	Memory      int
	Timeout     int
	Environment map[string]string
	Template    string // Template name for the function
	Handler     string
	CodePath    string
}

// TriggerConfig defines a trigger for a serverless function
type TriggerConfig struct {
	Type       string // http, schedule, queue, etc.
	Function   string
	Properties map[string]interface{}
}

// Resource represents a deployable serverless resource
type Resource struct {
	Type       string
	Name       string
	Properties map[string]interface{}
}

// DeploymentResult contains the result of a deployment
type DeploymentResult struct {
	Resources []Resource
	Endpoint  string
	Version   string
	CreatedAt time.Time
}

// DeploymentStatus represents the current status of a deployment
type DeploymentStatus struct {
	State     string
	Resources []ResourceStatus
	LastError string
	UpdatedAt time.Time
}

// ResourceStatus represents the status of a deployed resource
type ResourceStatus struct {
	Name   string
	Type   string
	State  string
	URL    string
	Errors []string
}

// Deployer handles serverless deployments
type Deployer struct {
	provider  Provider
	config    Config
	template  *engine.TemplateEngine
	validator *validation.TemplateValidator
}

// NewDeployer creates a new serverless deployer
func NewDeployer(config Config) (*Deployer, error) {
	provider, err := NewProvider(config.ProviderConfig.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize provider: %w", err)
	}

	if err := provider.Configure(context.Background(), config.ProviderConfig); err != nil {
		return nil, fmt.Errorf("failed to configure provider: %w", err)
	}

	return &Deployer{
		provider:  provider,
		config:    config,
		template:  engine.NewTemplateEngine(config.TemplatesDir),
		validator: validation.NewTemplateValidator(config.TemplatesDir),
	}, nil
}

// Deploy deploys serverless resources
func (d *Deployer) Deploy(ctx context.Context) (*DeploymentResult, error) {
	// Validate and process templates first
	if err := d.processTemplates(ctx); err != nil {
		return nil, fmt.Errorf("template processing failed: %w", err)
	}

	resources := d.prepareResources()
	return d.provider.Deploy(ctx, resources)
}

// Remove removes a deployed serverless application
func (d *Deployer) Remove(ctx context.Context, name string) error {
	return d.provider.Remove(ctx, name)
}

// GetStatus gets the status of a deployed serverless application
func (d *Deployer) GetStatus(ctx context.Context, name string) (*DeploymentStatus, error) {
	return d.provider.GetStatus(ctx, name)
}

// processTemplates processes function templates
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

// prepareResources prepares resources for deployment
func (d *Deployer) prepareResources() []Resource {
	resources := make([]Resource, 0)

	// Create resources for each function
	for _, function := range d.config.Functions {
		functionResource := Resource{
			Type: "Function",
			Name: function.Name,
			Properties: map[string]interface{}{
				"runtime":     function.Runtime,
				"memory":      function.Memory,
				"timeout":     function.Timeout,
				"environment": function.Environment,
				"handler":     function.Handler,
				"codePath":    function.CodePath,
			},
		}
		resources = append(resources, functionResource)
	}

	// Create resources for each trigger
	for _, trigger := range d.config.Triggers {
		triggerResource := Resource{
			Type: fmt.Sprintf("%sTrigger", trigger.Type),
			Name: fmt.Sprintf("%s-%s-trigger", trigger.Function, trigger.Type),
			Properties: map[string]interface{}{
				"function":   trigger.Function,
				"properties": trigger.Properties,
			},
		}
		resources = append(resources, triggerResource)
	}

	return resources
}

// NewProvider creates a new provider based on the provider type
func NewProvider(providerType string) (Provider, error) {
	switch providerType {
	case "aws":
		return NewAWSProvider(), nil
	// case "gcp":
	// 	return NewGCPProvider(), nil
	// case "azure":
	// 	return NewAzureProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}
