package template

import (
    "context"
    "text/template"
    "path/filepath"
)

type Generator struct {
    cache  *TemplateCache
    engine *template.Template
}

func NewGenerator() *Generator {
    return &Generator{
        cache: NewTemplateCache(),
    }
}

func (g *Generator) GenerateProject(ctx context.Context, config ProjectConfig) error {
    template, err := g.cache.Get(config.Template)
    if err != nil {
        return err
    }

    return g.engine.ExecuteTemplate(config.OutputDir, template, config.Variables)
}