package sourcemap

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

type SourceMapBuilder struct {
	config   *SourceMapConfig
	handler  *SourceMapHandler
	uploader ErrorTracker
	debugger *DebugSymbolGenerator
}

func NewSourceMapBuilder(config *SourceMapConfig) *SourceMapBuilder {
	return &SourceMapBuilder{
		config:   config,
		handler:  NewSourceMapHandler(filepath.Join("build", "sourcemaps"), "v1", config.PrivateStorage),
		uploader: createErrorTracker(config.ErrorTracking),
		debugger: NewDebugSymbolGenerator(),
	}
}

func (b *SourceMapBuilder) GenerateSourceMap(source, filename string) (*SourceMap, error) {
	// Create base source map
	sourceMap := &SourceMap{
		Version:    3,
		Sources:    []string{filename},
		Names:      []string{},
		Mappings:   "",
		File:       filename,
		SourceRoot: b.config.SourceRoot,
	}

	if b.config.IncludeContent {
		sourceMap.SourcesContent = []string{source}
	}

	// Apply optimizations for production
	if b.config.Mode == ProductionMode {
		if err := b.optimizeSourceMap(sourceMap); err != nil {
			return nil, err
		}
	}

	// Handle source map storage
	if err := b.processSourceMap(sourceMap); err != nil {
		return nil, err
	}

	return sourceMap, nil
}

func (b *SourceMapBuilder) processSourceMap(sourceMap *SourceMap) error {
	// Convert to JSON
	data, err := json.Marshal(sourceMap)
	if err != nil {
		return err
	}

	// Handle storage based on type
	switch b.config.Type {
	case InlineType:
		encoded := base64.StdEncoding.EncodeToString(data)
		sourceMap.Url = fmt.Sprintf("data:application/json;base64,%s", encoded)

	case ExternalType, BothType:
		if err := b.handler.StoreSourceMap(sourceMap.File, data); err != nil {
			return err
		}
		sourceMap.Url = fmt.Sprintf("%s.map", sourceMap.File)

		// Upload to error tracking if configured
		if b.uploader != nil {
			if err := b.uploader.UploadSourceMap(sourceMap.File, data); err != nil {
				return err
			}
		}
	}

	// Generate debug symbols
	if _, err := b.debugger.GenerateSymbols(sourceMap.Sources[0], sourceMap.File, data); err != nil {
		return err
	}

	return nil
}

func (b *SourceMapBuilder) optimizeSourceMap(sourceMap *SourceMap) error {
	// Remove source content in production if not needed
	if !b.config.IncludeContent {
		sourceMap.SourcesContent = nil
	}

	// Apply path rewrites
	for _, rewrite := range b.config.PathRewrites {
		for i, source := range sourceMap.Sources {
			sourceMap.Sources[i] = applyPathRewrite(source, rewrite)
		}
	}

	return nil
}

func applyPathRewrite(source string, rewrite PathRewrite) string {
	if rewrite.Pattern == "" {
		return source
	}
	return strings.Replace(source, rewrite.Pattern, rewrite.Replacement, -1)
}
