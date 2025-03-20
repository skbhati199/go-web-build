package optimization

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type defaultOptimizer struct {
	config DefaultConfig
}

type DefaultConfig struct {
	OutputDir     string
	MinifyJS      bool
	MinifyCSS     bool
	MinifyHTML    bool
	InlineAssets  bool
	CacheStrategy string
}

func newDefaultOptimizer() FrameworkOptimizer {
	return &defaultOptimizer{
		config: DefaultConfig{
			OutputDir:     "dist",
			MinifyJS:      true,
			MinifyCSS:     true,
			MinifyHTML:    true,
			InlineAssets:  false,
			CacheStrategy: "memory",
		},
	}
}

func (d *defaultOptimizer) Optimize(ctx context.Context, config OptimizationConfig) error {
	fmt.Println("Applying default optimizations...")

	// Apply basic optimizations
	if err := d.minifyAssets(config); err != nil {
		return fmt.Errorf("failed to minify assets: %w", err)
	}

	// Apply compression if enabled
	if config.Performance.Compression.Enable {
		if err := d.compressAssets(config.Performance.Compression); err != nil {
			return fmt.Errorf("failed to compress assets: %w", err)
		}
	}

	// Apply code splitting if enabled
	if config.Performance.CodeSplitting {
		if err := d.splitCode(config); err != nil {
			return fmt.Errorf("failed to split code: %w", err)
		}
	}

	return nil
}

func (d *defaultOptimizer) minifyAssets(config OptimizationConfig) error {
	fmt.Println("Minifying assets...")

	// Example: Minify JavaScript files
	jsFiles, err := filepath.Glob(filepath.Join(d.config.OutputDir, "*.js"))
	if err != nil {
		return fmt.Errorf("failed to find JS files: %w", err)
	}

	for _, file := range jsFiles {
		cmd := exec.Command("uglifyjs", file, "-o", file, "--compress", "--mangle")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to minify JS file %s: %w", file, err)
		}
	}

	// Example: Minify CSS files
	cssFiles, err := filepath.Glob(filepath.Join(d.config.OutputDir, "*.css"))
	if err != nil {
		return fmt.Errorf("failed to find CSS files: %w", err)
	}

	for _, file := range cssFiles {
		cmd := exec.Command("cleancss", "-o", file, file)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to minify CSS file %s: %w", file, err)
		}
	}

	return nil
}

func (d *defaultOptimizer) compressAssets(config CompressionConfig) error {
	fmt.Println("Compressing assets...")

	// Example: Compress files using gzip
	files, err := filepath.Glob(filepath.Join(d.config.OutputDir, "*"))
	if err != nil {
		return fmt.Errorf("failed to find files: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file, ".gz") {
			continue
		}

		cmd := exec.Command("gzip", "-k", "-f", file)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to compress file %s: %w", file, err)
		}
	}

	return nil
}

func (d *defaultOptimizer) splitCode(config OptimizationConfig) error {
	fmt.Println("Splitting code...")

	// Example: Use Webpack for code splitting
	cmd := exec.Command("webpack", "--config", "webpack.config.js", "--mode", "production")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to split code: %w", err)
	}

	return nil
}
