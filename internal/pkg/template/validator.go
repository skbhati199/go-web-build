package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type TemplateValidator struct {
	requiredFiles map[string][]string
	requiredDeps  map[string][]string
}

func NewTemplateValidator() *TemplateValidator {
	return &TemplateValidator{
		requiredFiles: map[string][]string{
			"react-typescript": {
				"tsconfig.json",
				"src/index.tsx",
				"package.json",
				"public/index.html",
				"src/App.tsx",
			},
			"react-javascript": {
				"src/index.js",
				"package.json",
				"public/index.html",
				"src/App.js",
			},
		},
		requiredDeps: map[string][]string{
			"react-typescript": {
				"react",
				"react-dom",
				"typescript",
				"@types/react",
				"@types/react-dom",
			},
			"react-javascript": {
				"react",
				"react-dom",
			},
		},
	}
}

func (v *TemplateValidator) Validate(templateName string, path string) error {
	if err := v.validateTemplate(templateName); err != nil {
		return err
	}

	if err := v.validateFiles(templateName, path); err != nil {
		return err
	}

	if err := v.validatePackageJSON(templateName, path); err != nil {
		return err
	}

	return nil
}

func (v *TemplateValidator) validateTemplate(templateName string) error {
	if _, exists := v.requiredFiles[templateName]; !exists {
		return fmt.Errorf("unknown template: %s", templateName)
	}
	return nil
}

func (v *TemplateValidator) validateFiles(templateName, path string) error {
	files := v.requiredFiles[templateName]
	for _, file := range files {
		fullPath := filepath.Join(path, file)
		if !fileExists(fullPath) {
			return fmt.Errorf("required file missing: %s", file)
		}
	}
	return nil
}

func (v *TemplateValidator) validatePackageJSON(templateName, path string) error {
	packageJSON := filepath.Join(path, "package.json")
	content, err := os.ReadFile(packageJSON)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %v", err)
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(content, &pkg); err != nil {
		return fmt.Errorf("invalid package.json format: %v", err)
	}

	requiredDeps := v.requiredDeps[templateName]
	for _, dep := range requiredDeps {
		if _, ok := pkg.Dependencies[dep]; !ok {
			if _, ok := pkg.DevDependencies[dep]; !ok {
				return fmt.Errorf("missing required dependency: %s", dep)
			}
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
