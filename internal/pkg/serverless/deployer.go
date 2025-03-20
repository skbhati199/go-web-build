package serverless

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/sonukumar/go-web-build/internal/pkg/cloud"
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

// Fix the typo in DeploymentResult
type DeploymentResult struct {
    Resources []Resource
    Endpoint  string
    Version   string
}

// Add Provider interface
type Provider interface {
    Deploy(context.Context, []Resource) (*DeploymentResult, error)
}

func NewDeployer(config Config) (*Deployer, error) {
    provider, err := cloud.NewCloudProvider(config.Provider)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize cloud provider: %w", err)
    }

    return &Deployer{
        provider:  provider,
        config:    config,
        template:  engine.NewTemplateEngine(config.TemplatesDir),
        validator: validation.NewTemplateValidator(config.TemplatesDir),
    }, nil
}

func (d *Deployer) Deploy(ctx context.Context) (*DeploymentResult, error) {
    // Validate and process templates first
    if err := d.processTemplates(ctx); err != nil {
        return nil, fmt.Errorf("template processing failed: %w", err)
    }

    // Prepare cloud configuration
    cloudConfig := cloud.CloudConfig{
        Provider: d.config.Provider,
        Region:   d.config.Region,
        Resources: cloud.ResourceConfig{
            Compute: cloud.ComputeConfig{
                Type: "lambda",
            },
        },
    }

    // Configure cloud provider
    if err := d.provider.Configure(ctx, cloudConfig); err != nil {
        return nil, fmt.Errorf("failed to configure cloud provider: %w", err)
    }

    // Prepare deployment options
    deployOptions := cloud.DeployOptions{
        ServiceName: "serverless-app",
        Version:     "1.0.0",
        Environment: "production",
        Labels: map[string]string{
            "deployment-type": "serverless",
        },
    }

    // Deploy using cloud provider
    if err := d.provider.Deploy(ctx, deployOptions); err != nil {
        return nil, fmt.Errorf("deployment failed: %w", err)
    }

    resources := d.prepareResources()
    result := &DeploymentResult{
        Resources: resources,
        Endpoint:  fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod", 
            deployOptions.ServiceName, d.config.Region),
        Version:   deployOptions.Version,
    }

    return result, nil
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
				"handler":     fmt.Sprintf("build/functions/%s/index.handler", function.Name),
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
