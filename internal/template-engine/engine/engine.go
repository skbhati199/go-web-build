package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"time"

	"github.com/sonukumar/go-web-build/internal/template-engine/cache"
	"github.com/sonukumar/go-web-build/internal/template-engine/variables"
)

type TemplateEngine struct {
	templatesDir string
	cache        *cache.TemplateCache
	variables    *variables.VariableManager
}

func NewTemplateEngine(templatesDir string) *TemplateEngine {
	return &TemplateEngine{
		templatesDir: templatesDir,
		cache:        cache.NewTemplateCache(24 * time.Hour), // Cache for 24 hours
		variables:    variables.NewVariableManager(),
	}
}

func (e *TemplateEngine) RegisterVariable(name, defaultValue, description string, required bool, validator func(string) error) {
	e.variables.RegisterVariable(name, defaultValue, description, required, validator)
}

func (e *TemplateEngine) Generate(templateName string, data *TemplateData, outputDir string) error {
	// Set up variables
	e.variables.SetVariable("projectName", data.ProjectName)
	e.variables.SetVariable("framework", data.Framework)
	e.variables.SetVariable("language", data.Language)
	for k, v := range data.Configuration {
		e.variables.SetVariable(k, v)
	}

	// Process template with variables
	if err := e.processTemplateWithVariables(templateName, outputDir); err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	return nil
}

func (e *TemplateEngine) processTemplateWithVariables(templateName, outputDir string) error {
	files, err := e.getTemplateFiles(templateName)
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		// Get values map from registered variables
		values := e.variables.GetVariables()

		processedContent, err := e.variables.ProcessContent(string(content), values)
		if err != nil {
			return err
		}

		outputPath := filepath.Join(outputDir, filepath.Base(file))
		if err := os.WriteFile(outputPath, []byte(processedContent), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (e *TemplateEngine) loadTemplate(name string) (*template.Template, error) {
	// Check cache first
	if cached, ok := e.cache.Get(name); ok {
		return cached.(*template.Template), nil
	}

	templatePath := filepath.Join(e.templatesDir, name)
	tmpl, err := template.ParseGlob(filepath.Join(templatePath, "*.tmpl"))
	if err != nil {
		return nil, err
	}

	// Cache the template
	e.cache.Set(name, tmpl)
	return tmpl, nil
}

func (e *TemplateEngine) getTemplateFiles(templateName string) ([]string, error) {
	var files []string
	templatePath := filepath.Join(e.templatesDir, templateName)

	err := filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".tmpl" {
			relPath, err := filepath.Rel(templatePath, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})

	return files, err
}

type TemplateData struct {
	ProjectName   string
	Framework     string
	Language      string
	Dependencies  []string
	Version       string
	Configuration map[string]interface{}
}

// Add this new method
func (e *TemplateEngine) GenerateReactReduxProject(name string, config map[string]interface{}) error {
	data := &TemplateData{
		ProjectName:   name,
		Framework:     "react-redux",
		Language:      "javascript",
		Version:       "1.0.0",
		Configuration: config,
	}

	// Add React-Redux specific variables
	e.variables.SetVariable("reactVersion", "18.2.0")
	e.variables.SetVariable("reduxVersion", "4.2.1")
	e.variables.SetVariable("reduxToolkitVersion", "1.9.5")
	e.variables.SetVariable("reactReduxVersion", "8.1.1")

	return e.Generate("react-redux", data, name)
}
func (e *TemplateEngine) GenerateReactRouterProject(name string, config map[string]interface{}) error {
	data := &TemplateData{
		ProjectName:   name,
		Framework:     "react-router",
		Language:      "javascript",
		Version:       "1.0.0",
		Configuration: config,
	}

	// Add React Router specific variables
	e.variables.SetVariable("reactVersion", "18.2.0")
	e.variables.SetVariable("routerVersion", "6.14.0")

	return e.Generate("react-router", data, name)
}
func (e *TemplateEngine) GenerateReactProject(name string, isTypeScript bool, config map[string]interface{}) error {
	language := "javascript"
	if isTypeScript {
		language = "typescript"
	}

	data := &TemplateData{
		ProjectName:   name,
		Framework:     "react",
		Language:      language,
		Version:       "1.0.0",
		Configuration: config,
	}

	// Add React-specific variables
	e.variables.SetVariable("reactVersion", "18.2.0")
	e.variables.SetVariable("useTypeScript", isTypeScript)

	templateName := fmt.Sprintf("react-%s", language)
	return e.Generate(templateName, data, name)
}
func (e *TemplateEngine) AddTestingSetup(projectPath string, config map[string]interface{}) error {
	data := &TemplateData{
		ProjectName:   filepath.Base(projectPath),
		Framework:     "react-testing",
		Version:       "1.0.0",
		Configuration: config,
	}

	// Add testing-specific variables
	e.variables.SetVariable("jestVersion", "29.5.0")
	e.variables.SetVariable("rtlVersion", "14.0.0")

	return e.Generate("react-testing", data, projectPath)
}
