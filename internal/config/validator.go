package config

import (
	"fmt"
	"strings"
)

type Validator struct {
	errors []string
}

func NewValidator() *Validator {
	return &Validator{
		errors: make([]string, 0),
	}
}

func (v *Validator) Validate(cfg *Config) error {
	v.validateEnvironment(cfg.Environment)
	v.validateServer(cfg.Server)
	v.validateBuild(cfg.Build)
	v.validateTemplates(cfg.Templates)

	if len(v.errors) > 0 {
		return fmt.Errorf("configuration validation failed:\n%s", strings.Join(v.errors, "\n"))
	}
	return nil
}

func (v *Validator) validateServer(server ServerConfig) {
	if server.Port < 1 || server.Port > 65535 {
		v.errors = append(v.errors, "invalid port number")
	}
	if server.Host == "" {
		v.errors = append(v.errors, "host is required")
	}
}

func (v *Validator) validateBuild(build BuildConfig) {
	if build.OutDir == "" {
		v.errors = append(v.errors, "build output directory is required")
	}
	if build.Cache && build.CacheDir == "" {
		v.errors = append(v.errors, "cache directory is required when cache is enabled")
	}
}

func (v *Validator) validateTemplates(templates TemplateConfig) {
	if templates.Directory == "" {
		v.errors = append(v.errors, "templates directory is required")
	}
}

func (v *Validator) validateEnvironment(env string) {
	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
	}
	if !validEnvs[env] {
		v.errors = append(v.errors, "invalid environment, must be one of: development, staging, production")
	}
}
