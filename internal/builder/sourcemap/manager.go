package sourcemap

import "fmt"

type SourceMapManager struct {
	handler  *SourceMapHandler
	uploader ErrorTracker
	debug    *DebugSymbolGenerator
}

func NewSourceMapManager(storageDir, version string, uploader ErrorTracker) *SourceMapManager {
	return &SourceMapManager{
		handler:  NewSourceMapHandler(storageDir, version, true),
		uploader: uploader,
		debug:    NewDebugSymbolGenerator(),
	}
}

func (m *SourceMapManager) ProcessSourceMap(originalFile, generatedFile string, sourceMap []byte) error {
	// Store the source map
	if err := m.handler.StoreSourceMap(generatedFile, sourceMap); err != nil {
		return fmt.Errorf("failed to store source map: %w", err)
	}

	// Upload to error tracking service
	if err := m.uploader.UploadSourceMap(generatedFile, sourceMap); err != nil {
		return fmt.Errorf("failed to upload source map: %w", err)
	}

	// Generate debug symbols
	if _, err := m.debug.GenerateSymbols(originalFile, generatedFile, sourceMap); err != nil {
		return fmt.Errorf("failed to generate debug symbols: %w", err)
	}

	return nil
}
