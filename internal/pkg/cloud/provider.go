package cloud

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUnsupportedProvider = errors.New("unsupported cloud provider")
	ErrInvalidCredentials  = errors.New("invalid cloud credentials")
	ErrDeploymentFailed    = errors.New("deployment failed")
)

type Provider interface {
	Name() string
	Configure(ctx context.Context, config CloudConfig) error
	Deploy(ctx context.Context, options DeployOptions) error
	Scale(ctx context.Context, options ScaleOptions) error
	Monitor(ctx context.Context) (*Metrics, error)
}

type CloudConfig struct {
	Provider    string
	Region      string
	Credentials Credentials
	Resources   ResourceConfig
}

type Credentials struct {
	AccessKey string
	SecretKey string
	ProjectID string
	TokenPath string
}

type ResourceConfig struct {
	CPU          string
	Memory       string
	AutoScale    bool
	MinInstances int
	MaxInstances int
	Storage      StorageConfig
}

type StorageConfig struct {
	Type            string
	Size            int
	Replicas        int
	Backup          bool
	BackupRetention time.Duration
}

type DeployOptions struct {
	ServiceName string
	Version     string
	Environment string
	Replicas    int
	Port        int
	Env         map[string]string
	Labels      map[string]string
}

type ScaleOptions struct {
	ServiceName string
	Replicas    int
	CPU         int
	Memory      int
	Immediate   bool
}

type Metrics struct {
	CPU          float64
	Memory       float64
	NetworkIn    float64
	NetworkOut   float64
	ResponseTime float64
	ErrorRate    float64
	Uptime       float64
	LastUpdated  time.Time
}

type baseProvider struct {
	name   string
	config CloudConfig
}

func (b *baseProvider) Name() string {
	return b.name
}

func (b *baseProvider) validateConfig(config CloudConfig) error {
	if config.Provider == "" {
		return errors.New("provider type is required")
	}
	if config.Region == "" {
		return errors.New("region is required")
	}
	if config.Resources.MinInstances < 0 || config.Resources.MaxInstances < config.Resources.MinInstances {
		return errors.New("invalid instance configuration")
	}
	return nil
}

func NewCloudProvider(providerType string) (Provider, error) {
	switch providerType {
	case "aws":
		return newAWSProvider(), nil
	// case "gcp":
	// 	return newGCPProvider(), nil
	// case "azure":
	// 	return newAzureProvider(), nil
	default:
		return nil, ErrUnsupportedProvider
	}
}
