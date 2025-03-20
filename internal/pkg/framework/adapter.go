package framework

import (
    "context"
)

type Adapter interface {
    Name() string
    Version() string
    Initialize(ctx context.Context, config Config) error
    GenerateProject(ctx context.Context, options ProjectOptions) error
    BuildProject(ctx context.Context, mode BuildMode) error
}

type Config struct {
    FrameworkType    string
    TemplateVersion  string
    Features         []string
    CloudNative      bool
    BuildOptions     BuildOptions
}

type BuildOptions struct {
    Optimization     bool
    Serverless      bool
    CloudProvider   string
    CacheStrategy   string
}

func NewFrameworkAdapter(frameworkType string) (Adapter, error) {
    switch frameworkType {
    case "react":
        return newReactAdapter(), nil
    case "vue":
        return newVueAdapter(), nil
    case "svelte":
        return newSvelteAdapter(), nil
    case "angular":
        return newAngularAdapter(), nil
    default:
        return nil, ErrUnsupportedFramework
    }
}