package cloud

import (
    "context"
)

type Provider interface {
    Name() string
    Configure(ctx context.Context, config CloudConfig) error
    Deploy(ctx context.Context, options DeployOptions) error
    Scale(ctx context.Context, options ScaleOptions) error
    Monitor(ctx context.Context) (*Metrics, error)
}

type CloudConfig struct {
    Provider     string
    Region      string
    Credentials Credentials
    Resources   ResourceConfig
}

type ResourceConfig struct {
    CPU          string
    Memory       string
    AutoScale    bool
    MinInstances int
    MaxInstances int
}

func NewCloudProvider(providerType string) (Provider, error) {
    switch providerType {
    case "aws":
        return newAWSProvider(), nil
    case "gcp":
        return newGCPProvider(), nil
    case "azure":
        return newAzureProvider(), nil
    default:
        return nil, ErrUnsupportedProvider
    }
}