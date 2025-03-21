package templateengine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/skbhati199/go-web-build/internal/template-engine/engine"
)

func setupTestTemplates(t *testing.T) (string, func()) {
	tempDir := filepath.Join(os.TempDir(), "go-web-build-test")
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test template structure
	templates := map[string]map[string]string{
		"react-typescript": {
			"template.json": `{
				"name": "react-typescript",
				"version": "1.0.0",
				"framework": "react",
				"language": "typescript"
			}`,
			"package.json": `{
				"name": "react-typescript-template",
				"version": "1.0.0"
			}`,
		},
		"vue-javascript": {
			"template.json": `{
				"name": "vue-javascript",
				"version": "1.0.0",
				"framework": "vue",
				"language": "javascript"
			}`,
			"package.json": `{
				"name": "vue-javascript-template",
				"version": "1.0.0"
			}`,
		},
	}

	for templateName, files := range templates {
		templateDir := filepath.Join(tempDir, templateName)
		if err := os.MkdirAll(templateDir, 0755); err != nil {
			t.Fatalf("Failed to create template directory: %v", err)
		}

		for fileName, content := range files {
			filePath := filepath.Join(templateDir, fileName)
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to create template file: %v", err)
			}
		}
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestNewManager(t *testing.T) {
	templatesDir, cleanup := setupTestTemplates(t)
	defer cleanup()

	manager := NewManager(templatesDir)
	if manager == nil {
		t.Fatal("NewManager returned nil")
	}

	if manager.engine == nil {
		t.Error("Template engine is nil")
	}
	if manager.validator == nil {
		t.Error("Template validator is nil")
	}
	if manager.versions == nil {
		t.Error("Version manager is nil")
	}
	if manager.cache == nil {
		t.Error("Template cache is nil")
	}
}

func TestCreateProject(t *testing.T) {
	templatesDir, cleanup := setupTestTemplates(t)
	defer cleanup()

	tests := []struct {
		name      string
		framework string
		language  string
		version   string
		config    map[string]interface{}
		wantErr   bool
	}{
		{
			name:      "Valid React TypeScript Project",
			framework: "react",
			language:  "typescript",
			version:   "1.0.0",
			config:    map[string]interface{}{"strict": true},
			wantErr:   false,
		},
		{
			name:      "Valid Vue JavaScript Project",
			framework: "vue",
			language:  "javascript",
			version:   "1.0.0",
			config:    nil,
			wantErr:   false,
		},
		{
			name:      "Invalid Framework",
			framework: "invalid",
			language:  "typescript",
			version:   "1.0.0",
			config:    nil,
			wantErr:   true,
		},
		{
			name:      "Invalid Version",
			framework: "react",
			language:  "typescript",
			version:   "999.0.0",
			config:    nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager(templatesDir)
			err := manager.CreateProject("test-project", tt.framework, tt.language, tt.version, tt.config)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Verify project creation
				projectPath := filepath.Join(templatesDir, "test-project")
				if _, err := os.Stat(projectPath); os.IsNotExist(err) {
					t.Error("Project directory was not created")
				}
			}
		})
	}
}

func TestTemplateVersioning(t *testing.T) {
	templatesDir, cleanup := setupTestTemplates(t)
	defer cleanup()

	manager := NewManager(templatesDir)

	// Test version caching
	templateName := "react-typescript"
	version := "1.0.0"
	cacheKey := templateName + "@" + version

	// First call should cache the version
	err := manager.CreateProject("test-project", "react", "typescript", version, nil)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Verify cache entry
	if _, ok := manager.cache[cacheKey]; !ok {
		t.Error("Template version was not cached")
	}
}

func TestIsValidFramework(t *testing.T) {
	tests := []struct {
		framework string
		want      bool
	}{
		{"react", true},
		{"vue", true},
		{"next", true},
		{"vite", true},
		{"invalid", false},
		{"REACT", true}, // Test case insensitivity
		{"Vue", true},   // Test case insensitivity
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.framework, func(t *testing.T) {
			if got := isValidFramework(tt.framework); got != tt.want {
				t.Errorf("isValidFramework(%q) = %v, want %v", tt.framework, got, tt.want)
			}
		})
	}
}

func TestTemplateData(t *testing.T) {
	templatesDir, cleanup := setupTestTemplates(t)
	defer cleanup()

	manager := NewManager(templatesDir)

	config := map[string]interface{}{
		"strict": true,
		"port":   3000,
	}

	data := &engine.TemplateData{
		ProjectName:   "test-project",
		Framework:     "react",
		Language:      "typescript",
		Version:       "1.0.0",
		Configuration: config,
	}

	// Verify template data is correctly populated
	if data.ProjectName != "test-project" {
		t.Error("Project name not set correctly")
	}
	if data.Framework != "react" {
		t.Error("Framework not set correctly")
	}
	if data.Language != "typescript" {
		t.Error("Language not set correctly")
	}
	if data.Version != "1.0.0" {
		t.Error("Version not set correctly")
	}
	if data.Configuration["strict"] != true {
		t.Error("Configuration not set correctly")
	}

	err := manager.CreateProject("test-project", "react", "typescript", "1.0.0", config)

	if err != nil {
		t.Errorf("Failed to create project: %v", err)
	}
}
