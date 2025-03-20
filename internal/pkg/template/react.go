package template

// ReactOptions defines the configuration options for generating a React project
type ReactOptions struct {
	ProjectName   string
	UseTypeScript bool
	Features      ReactFeatures
}

// ReactFeatures defines the optional features for a React project
type ReactFeatures struct {
	Router       bool
	StateManager string // "redux", "mobx", "context", etc.
	Styling      string // "css", "sass", "styled-components", "tailwind", etc.
	Testing      bool
	Development  bool
}

// NewReactOptions creates a new ReactOptions with default values
func NewReactOptions(projectName string) *ReactOptions {
	return &ReactOptions{
		ProjectName:   projectName,
		UseTypeScript: false,
		Features: ReactFeatures{
			Router:       false,
			StateManager: "",
			Styling:      "css",
			Testing:      true,
			Development:  true,
		},
	}
}

// WithTypeScript enables TypeScript for the React project
func (o *ReactOptions) WithTypeScript() *ReactOptions {
	o.UseTypeScript = true
	return o
}

// WithRouter enables React Router for the project
func (o *ReactOptions) WithRouter() *ReactOptions {
	o.Features.Router = true
	return o
}

// WithRedux enables Redux for state management
func (o *ReactOptions) WithRedux() *ReactOptions {
	o.Features.StateManager = "redux"
	return o
}

// WithStyling sets the styling approach for the project
func (o *ReactOptions) WithStyling(styling string) *ReactOptions {
	o.Features.Styling = styling
	return o
}

// WithTesting enables testing setup for the project
func (o *ReactOptions) WithTesting() *ReactOptions {
	o.Features.Testing = true
	return o
}

// WithoutDevelopment disables development mode features
func (o *ReactOptions) WithoutDevelopment() *ReactOptions {
	o.Features.Development = false
	return o
}
