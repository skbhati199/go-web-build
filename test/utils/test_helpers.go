package utils

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func WaitForFile(ctx context.Context, path string, timeout time.Duration) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	deadline := time.Now().Add(timeout)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if _, err := os.Stat(path); err == nil {
				return nil
			}
			if time.Now().After(deadline) {
				return errors.New("timeout waiting for file: " + path)
			}
		}
	}
}

func VerifyFileExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", path)
		return
	}
	if err != nil {
		t.Errorf("error checking file %s: %v", path, err)
		return
	}
	if info.Size() == 0 {
		t.Errorf("file %s exists but is empty", path)
	}
}

func VerifyBuildArtifacts(t *testing.T, buildDir string, files []string) {
	t.Helper()
	for _, file := range files {
		path := filepath.Join(buildDir, file)
		VerifyFileExists(t, path)
	}
}

func CleanupTestDir(t *testing.T, dir string) {
	t.Helper()
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("failed to cleanup test directory: %v", err)
	}
}
