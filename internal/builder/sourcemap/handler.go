package sourcemap

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

type SourceMapHandler struct {
	storageDir string
	version    string
	isPrivate  bool
}

func NewSourceMapHandler(storageDir, version string, isPrivate bool) *SourceMapHandler {
	return &SourceMapHandler{
		storageDir: storageDir,
		version:    version,
		isPrivate:  isPrivate,
	}
}

func (h *SourceMapHandler) CompressSourceMap(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(data); err != nil {
		return nil, fmt.Errorf("failed to compress source map: %w", err)
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (h *SourceMapHandler) StoreSourceMap(filename string, data []byte) error {
	compressed, err := h.CompressSourceMap(data)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(data)
	versionedPath := filepath.Join(h.storageDir, h.version, base64.URLEncoding.EncodeToString(hash[:]))

	if err := os.MkdirAll(filepath.Dir(versionedPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.WriteFile(versionedPath, compressed, 0644)
}
