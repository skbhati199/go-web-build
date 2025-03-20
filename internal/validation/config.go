package validation

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type ConfigValidator struct {
	validator *validator.Validate
}

func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		validator: validator.New(),
	}
}

func (v *ConfigValidator) ValidateConfig(cfg interface{}) error {
	if err := v.validator.Struct(cfg); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return fmt.Errorf("validation failed: %w", err)
		}
		return fmt.Errorf("invalid config structure: %w", err)
	}
	return nil
}
