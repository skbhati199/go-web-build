package templateengine // Update package name

import (
	"fmt"
	"strings"

	"github.com/sonukumar/go-web-build/internal/template-engine/engine"
	"github.com/sonukumar/go-web-build/internal/template-engine/validation"
	"github.com/sonukumar/go-web-build/internal/template-engine/version"
)

type Manager struct {
	engine    *engine.TemplateEngine
	validator *validation.TemplateValidator
	versions  *version.VersionManager
	cache     map[string]*version.TemplateVersion
}

func NewManager(templatesDir string) *Manager {
	return &Manager{
		engine:    engine.NewTemplateEngine(templatesDir),
		validator: validation.NewTemplateValidator(templatesDir),
		versions:  version.NewVersionManager(templatesDir),
		cache:     make(map[string]*version.TemplateVersion),
	}
}

func (m *Manager) CreateProject(name, framework, language, ver string, config map[string]interface{}) error {
	if !isValidFramework(framework) {
		return fmt.Errorf("unsupported framework: %s", framework)
	}

	templateName := fmt.Sprintf("%s-%s", strings.ToLower(framework), strings.ToLower(language))

	if err := m.validator.ValidateTemplate(templateName); err != nil {
		return fmt.Errorf("invalid template combination: %w", err)
	}

	var templateVersion *version.TemplateVersion
	var err error

	// Check cache first
	cacheKey := fmt.Sprintf("%s@%s", templateName, ver)
	if cachedVersion, ok := m.cache[cacheKey]; ok {
		templateVersion = cachedVersion
	} else {
		if ver != "" {
			templateVersion, err = m.versions.GetVersion(templateName, ver)
		} else {
			templateVersion, err = m.versions.GetLatestVersion(templateName)
		}
		if err != nil {
			return fmt.Errorf("failed to get template version: %w", err)
		}
		// Cache the template version
		m.cache[cacheKey] = templateVersion
	}

	data := &engine.TemplateData{
		ProjectName:   name,
		Framework:     framework,
		Language:      language,
		Version:       templateVersion.Version.String(),
		Configuration: config,
	}

	return m.engine.Generate(templateName, data, name)
}

func isValidFramework(framework string) bool {
	validFrameworks := map[string]bool{
		"react": true,
		"vue":   true,
		"next":  true,
		"vite":  true,
	}
	return validFrameworks[framework]
}
