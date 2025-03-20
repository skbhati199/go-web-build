package builder

import (
	"context"
	"path/filepath"

	"github.com/sonukumar/go-web-build/internal/builder/sourcemap"
)

type Builder struct {
	sourceMapBuilder *sourcemap.SourceMapBuilder
}

type Options struct {
	Mode    string
	OutDir  string
	Config  string
	BaseDir string
}

func New() *Builder {
	return &Builder{
		sourceMapBuilder: sourcemap.NewSourceMapBuilder(&sourcemap.SourceMapConfig{
			Type:           sourcemap.ExternalType,
			Mode:           sourcemap.ProductionMode,
			IncludeContent: true,
			SourceRoot:     "/src",
		}),
	}
}

func (b *Builder) Build(ctx context.Context, opts Options) error {
	// Setup build directory
	buildDir := filepath.Join(opts.BaseDir, opts.OutDir)

	// Configure build based on mode
	if opts.Mode == "production" {
		return b.buildProduction(ctx, buildDir, opts)
	}
	return b.buildDevelopment(ctx, buildDir, opts)
}

func (b *Builder) buildDevelopment(ctx context.Context, buildDir string, opts Options) error {
	// Development build implementation
	return nil
}

func (b *Builder) buildProduction(ctx context.Context, buildDir string, opts Options) error {
	// Production build implementation
	return nil
}
