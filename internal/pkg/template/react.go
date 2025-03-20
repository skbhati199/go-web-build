package template

type ReactOptions struct {
    ProjectName   string
    UseTypeScript bool
    Features      ReactFeatures
}

type ReactFeatures struct {
    Router       bool
    StateManager string // redux, mobx, zustand
    Styling      string // css-modules, styled-components, tailwind
    Testing      bool
}

func (g *TemplateGenerator) getReactTemplate(useTypeScript bool) *Template {
    if useTypeScript {
        return g.cache.Get("react-typescript")
    }
    return g.cache.Get("react-javascript")
}