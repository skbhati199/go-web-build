package template

import (
    "errors"
    "path/filepath"
)

type TemplateValidator struct {
    requiredFiles map[string][]string
}

func NewTemplateValidator() *TemplateValidator {
    return &TemplateValidator{
        requiredFiles: map[string][]string{
            "react-typescript": {
                "tsconfig.json",
                "src/index.tsx",
                "package.json",
            },
            "react-javascript": {
                "src/index.js",
                "package.json",
            },
        },
    }
}

func (v *TemplateValidator) Validate(templateName string, path string) error {
    files, exists := v.requiredFiles[templateName]
    if !exists {
        return errors.New("unknown template")
    }

    for _, file := range files {
        if !fileExists(filepath.Join(path, file)) {
            return errors.New("template validation failed")
        }
    }
    return nil
}