package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sonukumar/go-web-build/internal/builder"
	"github.com/sonukumar/go-web-build/internal/config"
	"github.com/sonukumar/go-web-build/test/utils"
)

type TestBuilder struct {
	config *config.Config
	dir    string
}

func setupTestBuilder(t *testing.T, dir string) *TestBuilder {
	t.Helper()

	cfg := &config.Config{
		Environment: "development",
		Build: config.BuildConfig{
			OutDir: filepath.Join(dir, "dist"),
		},
		Templates: config.TemplateConfig{},
		Server: config.ServerConfig{
			Port:    8080,
			Host:    "localhost",
			DevMode: false,
		},
	}

	return &TestBuilder{
		config: cfg,
		dir:    dir,
	}
}

func (b *TestBuilder) Build(ctx context.Context, configFile string, mode string) error {
	buildOpts := builder.Options{
		Mode:    mode,
		OutDir:  b.config.Build.OutDir,
		Config:  configFile,
		BaseDir: b.dir,
	}

	return builder.New().Build(ctx, buildOpts)
}

func TestBuildProcess(t *testing.T) {
	tempDir := t.TempDir()
	defer utils.CleanupTestDir(t, tempDir)

	tests := []struct {
		name     string
		config   string
		mode     string
		timeout  time.Duration
		validate func(context.Context, string) error
	}{
		{
			name:    "Development Build",
			config:  "dev.config.json",
			mode:    "development",
			timeout: 1 * time.Minute,
			validate: func(ctx context.Context, dir string) error {
				return validateDevBuild(ctx, dir)
			},
		},
		{
			name:    "Production Build",
			config:  "prod.config.json",
			mode:    "production",
			timeout: 2 * time.Minute,
			validate: func(ctx context.Context, dir string) error {
				return validateProdBuild(ctx, dir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			builder := setupTestBuilder(t, tempDir)
			err := builder.Build(ctx, tt.config, tt.mode)
			if err != nil {
				t.Fatalf("build failed: %v", err)
			}

			if err := tt.validate(ctx, tempDir); err != nil {
				t.Errorf("validation failed: %v", err)
			}
		})
	}
}

func validateDevBuild(ctx context.Context, dir string) error {
	expectedFiles := []string{
		"dist/index.html",
		"dist/assets",
		"dist/static",
	}

	for _, file := range expectedFiles {
		if err := utils.WaitForFile(ctx, filepath.Join(dir, file), 30*time.Second); err != nil {
			return errors.New("dev build validation failed: " + err.Error())
		}
	}
	return nil
}

func validateProdBuild(ctx context.Context, dir string) error {
	expectedFiles := []string{
		"dist/index.html",
		"dist/assets",
		"dist/static",
		"dist/sourcemaps",
	}

	for _, file := range expectedFiles {
		if err := utils.WaitForFile(ctx, filepath.Join(dir, file), 30*time.Second); err != nil {
			return errors.New("prod build validation failed: " + err.Error())
		}
	}
	return verifyProdOptimizations(dir)
}

func verifyProdOptimizations(dir string) error {
	// Implement production build verification
	files := []string{
		filepath.Join(dir, "dist/assets/main.min.js"),
		filepath.Join(dir, "dist/static/css/main.min.css"),
		filepath.Join(dir, "dist/sourcemaps/main.js.map"),
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("failed to verify optimization for %s: %v", file, err)
		}

		if info.Size() == 0 {
			return fmt.Errorf("optimized file %s is empty", file)
		}

		// Verify source maps
		if strings.HasSuffix(file, ".map") {
			if err := verifySourceMap(file); err != nil {
				return fmt.Errorf("invalid source map %s: %v", file, err)
			}
		}
	}
	return nil
}

func verifySourceMap(mapFile string) error {
	data, err := os.ReadFile(mapFile)
	if err != nil {
		return err
	}

	var sourceMap struct {
		Version    int      `json:"version"`
		Sources    []string `json:"sources"`
		Mappings   string   `json:"mappings"`
		SourceRoot string   `json:"sourceRoot"`
	}

	if err := json.Unmarshal(data, &sourceMap); err != nil {
		return err
	}

	if sourceMap.Version != 3 {
		return fmt.Errorf("invalid source map version: %d", sourceMap.Version)
	}

	if len(sourceMap.Sources) == 0 {
		return errors.New("source map contains no sources")
	}

	if sourceMap.Mappings == "" {
		return errors.New("source map contains no mappings")
	}

	return nil
}
