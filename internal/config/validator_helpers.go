package config

import (
	"os"
	"path/filepath"
)

func isValidDirectory(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isValidPort(port int) bool {
	return port > 0 && port < 65536
}

func isAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}
