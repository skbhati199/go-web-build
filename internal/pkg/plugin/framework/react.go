package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/sonukumar/go-web-build/internal/pkg/optimization"
)

type ReactOptimizer struct {
	config OptimizationConfig
}

type OptimizationConfig struct {
	TreeShaking   bool `json:"treeShaking"`
	CodeSplitting bool `json:"codeSplitting"`
	LazyLoading   bool `json:"lazyLoading"`
	Prefetching   bool `json:"prefetching"`
	OutputDir     string `json:"outputDir"`
	Environment   string `json:"environment"`
}

func (r *ReactOptimizer) Name() string {
	return "react-optimizer"
}

func (r *ReactOptimizer) Version() string {
	return "1.0.0"
}

func (r *ReactOptimizer) Init(ctx context.Context, config json.RawMessage) error {
	if err := json.Unmarshal(config, &r.config); err != nil {
		return fmt.Errorf("failed to parse React optimizer config: %w", err)
	}
	
	// Set defaults if not provided
	if r.config.OutputDir == "" {
		r.config.OutputDir = "dist"
	}
	
	if r.config.Environment == "" {
		r.config.Environment = "production"
	}
	
	return nil
}

func (r *ReactOptimizer) Execute(ctx context.Context, params map[string]interface{}) error {
	// Convert our plugin config to the optimization package config
	optimizerConfig := optimization.OptimizationConfig{
		Framework:   "react",
		Environment: r.config.Environment,
		Performance: optimization.PerformanceConfig{
			TreeShaking:   r.config.TreeShaking,
			CodeSplitting: r.config.CodeSplitting,
			LazyLoading:   r.config.LazyLoading,
			Compression: optimization.CompressionConfig{
				Enable:    true,
				Level:     9,
				Algorithm: "gzip",
			},
		},
	}
	
	// Get the project directory from params or use current directory
	projectDir, ok := params["projectDir"].(string)
	if !ok {
		var err error
		projectDir, err = filepath.Abs(".")
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}
	
	// Create the optimizer and run it
	optimizer := optimization.NewOptimizer("react")
	if err := optimizer.Optimize(ctx, optimizerConfig); err != nil {
		return fmt.Errorf("optimization failed: %w", err)
	}
	
	// Handle prefetching if enabled (not part of the core optimizer)
	if r.config.Prefetching {
		if err := r.enablePrefetching(projectDir); err != nil {
			return fmt.Errorf("failed to enable prefetching: %w", err)
		}
	}
	
	fmt.Println("React optimization completed successfully")
	return nil
}

func (r *ReactOptimizer) Cleanup(ctx context.Context) error {
	// Clean up any temporary files or resources
	return nil
}

func (r *ReactOptimizer) enablePrefetching(projectDir string) error {
	// This is a custom feature not in the core optimizer
	// Implement prefetching for React components
	
	// Example: Add prefetch webpack plugin to the webpack config
	webpackConfigPath := filepath.Join(projectDir, "webpack.config.js")
	
	// Check if webpack config exists
	if _, err := exec.Command("test", "-f", webpackConfigPath).Output(); err != nil {
		return fmt.Errorf("webpack config not found: %w", err)
	}
	
	// Add prefetching plugin (this is a simplified example)
	cmd := exec.Command("npx", "webpack-cli", "--config", webpackConfigPath, "--prefetch")
	cmd.Dir = projectDir
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to enable prefetching: %s: %w", output, err)
	}
	
	return nil
}
