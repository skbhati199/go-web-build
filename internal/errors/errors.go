package errors

import (
	"fmt"
)

type ErrorType string

const (
	ErrorTypeValidation ErrorType = "VALIDATION"
	ErrorTypeConfig     ErrorType = "CONFIG"
	ErrorTypeTemplate   ErrorType = "TEMPLATE"
	ErrorTypeBuild      ErrorType = "BUILD"
	ErrorTypeSystem     ErrorType = "SYSTEM"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Type, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewValidationError(msg string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: msg,
		Err:     err,
	}
}

func NewConfigError(msg string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeConfig,
		Message: msg,
		Err:     err,
	}
}

func NewTemplateError(msg string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeTemplate,
		Message: msg,
		Err:     err,
	}
}

func NewBuildError(msg string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeBuild,
		Message: msg,
		Err:     err,
	}
}

func NewSystemError(msg string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeSystem,
		Message: msg,
		Err:     err,
	}
}
