package serverless

import (
    "context"
)

type Deployer struct {
    provider Provider
    config   Config
}

type Config struct {
    Provider    string
    Region     string
    Functions  []FunctionConfig
    Triggers   []TriggerConfig
}

type FunctionConfig struct {
    Name        string
    Runtime     string
    Memory      int
    Timeout     int
    Environment map[string]string
}

func NewDeployer(config Config) *Deployer {
    return &Deployer{
        provider: selectProvider(config.Provider),
        config:   config,
    }
}

func (d *Deployer) Deploy(ctx context.Context) (*DeploymentResult, error) {
    resources := d.prepareResources()
    return d.provider.Deploy(ctx, resources)
}