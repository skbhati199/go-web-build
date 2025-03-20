package template

import (
    "context"
    "path/filepath"
    "text/template"
)

type TemplateGenerator struct {
    cache  *TemplateCache
    config *Config
}

func NewTemplateGenerator(config *Config) *TemplateGenerator {
    return &TemplateGenerator{
        cache:  NewTemplateCache(),
        config: config,
    }
}

func (g *TemplateGenerator) GenerateReactProject(ctx context.Context, opts *ReactOptions) error {
    template := g.getReactTemplate(opts.UseTypeScript)
    return g.generateFromTemplate(ctx, template, opts)
}

func (g *TemplateGenerator) generateFromTemplate(ctx context.Context, tmpl *Template, data interface{}) error {
    return tmpl.Execute(data)
}