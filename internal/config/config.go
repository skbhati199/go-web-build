package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"environment" validate:"required,oneof=development staging production"`
	Server      ServerConfig   `mapstructure:"server" validate:"required"`
	Build       BuildConfig    `mapstructure:"build" validate:"required"`
	Templates   TemplateConfig `mapstructure:"templates" validate:"required"`
}

type ServerConfig struct {
	Port    int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Host    string `mapstructure:"host" validate:"required"`
	DevMode bool   `mapstructure:"dev_mode"`
}

type BuildConfig struct {
	OutDir    string `mapstructure:"out_dir" validate:"required"`
	SourceMap bool   `mapstructure:"source_map"`
	Minify    bool   `mapstructure:"minify"`
	Cache     bool   `mapstructure:"cache"`
	CacheDir  string `mapstructure:"cache_dir" validate:"required_if=Cache true"`
}

type TemplateConfig struct {
	Directory string `mapstructure:"directory" validate:"required"`
	Cache     bool   `mapstructure:"cache"`
}

func LoadConfig(configPath string, env string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Load base config
	v.AddConfigPath("config")
	v.SetConfigName("default")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading base config: %w", err)
		}
	}

	// Load environment specific config
	if env != "" {
		v.SetConfigName(env)
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("error loading %s config: %w", env, err)
		}
	}

	// Custom config overrides
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("error merging custom config: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("environment", "development")
	v.SetDefault("server.port", 3000)
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.dev_mode", true)
	v.SetDefault("build.out_dir", "dist")
	v.SetDefault("build.source_map", true)
	v.SetDefault("build.minify", true)
	v.SetDefault("build.cache", true)
	v.SetDefault("build.cache_dir", ".cache")
	v.SetDefault("templates.directory", "templates")
	v.SetDefault("templates.cache", true)
}
