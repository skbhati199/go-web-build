package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type TemplateMetadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Framework   string   `json:"framework"`
	Language    string   `json:"language"`
	Files       []string `json:"files"`
	Required    []string `json:"required"`
}

type TemplateValidator struct {
	templatesDir string
}

func NewTemplateValidator(templatesDir string) *TemplateValidator {
	return &TemplateValidator{
		templatesDir: templatesDir,
	}
}

func (v *TemplateValidator) ValidateTemplate(templateName string) error {
	metadata, err := v.loadMetadata(templateName)
	if err != nil {
		return fmt.Errorf("failed to load template metadata: %w", err)
	}

	if err := v.validateFiles(templateName, metadata.Files); err != nil {
		return fmt.Errorf("template file validation failed: %w", err)
	}

	if err := v.validateSyntax(templateName, metadata.Files); err != nil {
		return fmt.Errorf("template syntax validation failed: %w", err)
	}

	return nil
}

func (v *TemplateValidator) loadMetadata(templateName string) (*TemplateMetadata, error) {
	metadataPath := filepath.Join(v.templatesDir, templateName, "template.json")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	var metadata TemplateMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (v *TemplateValidator) validateFiles(templateName string, files []string) error {
	templatePath := filepath.Join(v.templatesDir, templateName)
	for _, file := range files {
		path := filepath.Join(templatePath, file)
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("template file not found: %s", file)
		}
	}
	return nil
}

func (v *TemplateValidator) validateSyntax(templateName string, files []string) error {
	templatePath := filepath.Join(v.templatesDir, templateName)
	for _, file := range files {
		path := filepath.Join(templatePath, file)
		if _, err := os.ReadFile(path); err != nil {
			return fmt.Errorf("failed to read template file %s: %w", file, err)
		}
	}
	return nil
}
