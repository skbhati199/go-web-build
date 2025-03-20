package variables

import (
	"fmt"
	"regexp"
	"strings"
)

type Variable struct {
	Name         string
	DefaultValue string
	Description  string
	Required     bool
	Validator    func(string) error
}

type VariableManager struct {
	variables map[string]*Variable
	patterns  map[string]*regexp.Regexp
}

func NewVariableManager() *VariableManager {
	return &VariableManager{
		variables: make(map[string]*Variable),
		patterns: map[string]*regexp.Regexp{
			"standard": regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`),
			"default":  regexp.MustCompile(`\{\{\s*(\w+)\s*\|\|\s*([^}]+)\s*\}\}`),
			"function": regexp.MustCompile(`\{\{\s*(\w+)\((.*?)\)\s*\}\}`),
		},
	}
}

func (m *VariableManager) RegisterVariable(name, defaultValue, description string, required bool, validator func(string) error) {
	m.variables[name] = &Variable{
		Name:         name,
		DefaultValue: defaultValue,
		Description:  description,
		Required:     required,
		Validator:    validator,
	}
}

func (m *VariableManager) GetVariable(name string) (*Variable, bool) {
	v, exists := m.variables[name]
	return v, exists
}

func (m *VariableManager) SetVariable(name string, value interface{}) {
	if v, exists := m.variables[name]; exists {
		// Update existing variable's default value
		v.DefaultValue = fmt.Sprintf("%v", value)
	} else {
		// Create new variable
		m.variables[name] = &Variable{
			Name:         name,
			DefaultValue: fmt.Sprintf("%v", value),
			Description:  "Dynamically set variable",
			Required:     false,
		}
	}
}

func (m *VariableManager) ProcessContent(content string, values map[string]string) (string, error) {
	// Process standard variables
	content = m.patterns["standard"].ReplaceAllStringFunc(content, func(match string) string {
		key := m.patterns["standard"].FindStringSubmatch(match)[1]
		if value := values[key]; value != "" {
			return value
		}
		if v, exists := m.variables[key]; exists {
			return v.DefaultValue
		}
		return match
	})

	// Process variables with default values
	content = m.patterns["default"].ReplaceAllStringFunc(content, func(match string) string {
		parts := m.patterns["default"].FindStringSubmatch(match)
		key, defaultValue := parts[1], parts[2]
		if value := values[key]; value != "" {
			return value
		}
		if v, exists := m.variables[key]; exists && v.DefaultValue != "" {
			return v.DefaultValue
		}
		return strings.Trim(defaultValue, `"'`)
	})

	// Process function calls
	content = m.patterns["function"].ReplaceAllStringFunc(content, func(match string) string {
		parts := m.patterns["function"].FindStringSubmatch(match)
		funcName, args := parts[1], strings.Split(parts[2], ",")
		return m.executeFunction(funcName, args)
	})

	// Validate required variables
	for name, variable := range m.variables {
		value := values[name]
		if value == "" {
			value = variable.DefaultValue
		}

		if variable.Required && value == "" {
			return "", fmt.Errorf("required variable %s is not set", name)
		}

		if variable.Validator != nil {
			if err := variable.Validator(value); err != nil {
				return "", fmt.Errorf("validation failed for variable %s: %w", name, err)
			}
		}
	}

	return content, nil
}

func (m *VariableManager) executeFunction(name string, args []string) string {
	switch name {
	case "uppercase":
		if len(args) > 0 {
			return strings.ToUpper(strings.TrimSpace(args[0]))
		}
	case "lowercase":
		if len(args) > 0 {
			return strings.ToLower(strings.TrimSpace(args[0]))
		}
	case "capitalize":
		if len(args) > 0 {
			s := strings.TrimSpace(args[0])
			if len(s) > 0 {
				return strings.ToUpper(s[:1]) + s[1:]
			}
		}
	}
	return ""
}

// Add this new method
func (m *VariableManager) GetVariables() map[string]string {
	values := make(map[string]string)
	for name, v := range m.variables {
		values[name] = v.DefaultValue
	}
	return values
}
