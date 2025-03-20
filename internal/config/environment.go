package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
	Staging     Environment = "staging"
)

func LoadEnvironmentConfig(env Environment) (*Config, error) {
	configPath := getConfigPath(env)
	config, err := LoadConfig(configPath, string(env))
	if err != nil {
		return nil, fmt.Errorf("failed to load %s config: %w", env, err)
	}

	// Override with environment variables
	applyEnvironmentVariables(config)

	// Validate configuration
	validator := NewValidator()
	if err := validator.Validate(config); err != nil {
		return nil, err
	}

	return config, nil
}

func getConfigPath(env Environment) string {
	configName := fmt.Sprintf("config.%s.yaml", env)
	if customPath := os.Getenv("CONFIG_PATH"); customPath != "" {
		return filepath.Join(customPath, configName)
	}
	return filepath.Join("config", configName)
}

func applyEnvironmentVariables(config *Config) {
	if port := os.Getenv("PORT"); port != "" {
		// Parse and set the port
		finalPort, err := strconv.Atoi(port)
		if err != nil {
			// Handle the error appropriately, e.g., log it
			return
		}
		config.Server.Port = finalPort
	}
	if host := os.Getenv("HOST"); host != "" {
		config.Server.Host = host
	}
	// Add more environment variable overrides as needed
}
