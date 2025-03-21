package e2e

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/skbhati199/go-web-build/test/utils"
)

func TestCLICommands(t *testing.T) {
	tempDir := t.TempDir()
	defer utils.CleanupTestDir(t, tempDir)

	tests := []struct {
		name    string
		command []string
		verify  func(t *testing.T, output string, dir string)
		timeout time.Duration
	}{
		{
			name:    "Create React Project",
			command: []string{"create", "react-app", "--framework", "react", "--template", "typescript"},
			timeout: 2 * time.Minute,
			verify: func(t *testing.T, output string, dir string) {
				verifyReactProject(t, filepath.Join(dir, "react-app"))
			},
		},
		{
			name:    "Build Development",
			command: []string{"build", "--mode", "development"},
			timeout: 1 * time.Minute,
			verify: func(t *testing.T, output string, dir string) {
				verifyDevBuild(t, dir)
			},
		},
		{
			name:    "Build Production",
			command: []string{"build", "--mode", "production"},
			timeout: 2 * time.Minute,
			verify: func(t *testing.T, output string, dir string) {
				verifyProdBuild(t, dir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("go-web-build", tt.command...)
			cmd.Dir = tempDir

			// Set timeout context
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			// Replace incorrect context cancellation with proper implementation
			go func() {
				<-ctx.Done()
				if cmd.Process != nil {
					cmd.Process.Kill()
				}
			}()

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("command failed: %v\noutput: %s", err, output)
			}

			tt.verify(t, string(output), tempDir)
		})
	}
}

func verifyReactProject(t *testing.T, projectDir string) {
	t.Helper()

	expectedFiles := []string{
		"package.json",
		"tsconfig.json",
		"src/index.tsx",
		"public/index.html",
	}

	utils.VerifyBuildArtifacts(t, projectDir, expectedFiles)
	verifyPackageContent(t, projectDir)
}

func verifyPackageContent(t *testing.T, projectDir string) {
	t.Helper()

	packageJSON, err := ioutil.ReadFile(filepath.Join(projectDir, "package.json"))
	if err != nil {
		t.Fatalf("failed to read package.json: %v", err)
	}

	var pkg struct {
		Name         string            `json:"name"`
		Dependencies map[string]string `json:"dependencies"`
		Scripts      map[string]string `json:"scripts"`
	}

	if err := json.Unmarshal(packageJSON, &pkg); err != nil {
		t.Fatalf("failed to parse package.json: %v", err)
	}

	requiredDeps := []string{"react", "react-dom", "typescript"}
	for _, dep := range requiredDeps {
		if _, ok := pkg.Dependencies[dep]; !ok {
			t.Errorf("missing required dependency: %s", dep)
		}
	}

	requiredScripts := []string{"start", "build", "test"}
	for _, script := range requiredScripts {
		if _, ok := pkg.Scripts[script]; !ok {
			t.Errorf("missing required script: %s", script)
		}
	}
}

func verifyDevBuild(t *testing.T, dir string) {
	t.Helper()

	expectedFiles := []string{
		"dist/index.html",
		"dist/assets/main.js",
		"dist/static/css/main.css",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, file := range expectedFiles {
		if err := utils.WaitForFile(ctx, filepath.Join(dir, file), 5*time.Second); err != nil {
			t.Errorf("failed waiting for file %s: %v", file, err)
		}
	}

	utils.VerifyBuildArtifacts(t, dir, expectedFiles)
	verifyDevBuildContent(t, dir)
}

func verifyDevBuildContent(t *testing.T, dir string) {
	t.Helper()

	mainJS := filepath.Join(dir, "dist/assets/main.js")
	content, err := ioutil.ReadFile(mainJS)
	if err != nil {
		t.Fatalf("failed to read main.js: %v", err)
	}

	if len(content) == 0 {
		t.Error("main.js is empty")
	}

	if !strings.Contains(string(content), "sourceMappingURL") {
		t.Error("source map reference missing in development build")
	}
}

func verifyProdBuild(t *testing.T, dir string) {
	t.Helper()

	expectedFiles := []string{
		"dist/index.html",
		"dist/assets/main.min.js",
		"dist/static/css/main.min.css",
		"dist/sourcemaps/main.js.map",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, file := range expectedFiles {
		if err := utils.WaitForFile(ctx, filepath.Join(dir, file), 5*time.Second); err != nil {
			t.Errorf("failed waiting for file %s: %v", file, err)
		}
	}

	utils.VerifyBuildArtifacts(t, dir, expectedFiles)
	verifyProdBuildOptimizations(t, dir)
}

func verifyProdBuildOptimizations(t *testing.T, dir string) {
	t.Helper()

	mainJS := filepath.Join(dir, "dist/assets/main.min.js")
	content, err := ioutil.ReadFile(mainJS)
	if err != nil {
		t.Fatalf("failed to read main.min.js: %v", err)
	}

	if len(content) == 0 {
		t.Error("main.min.js is empty")
	}

	// Check for minification
	if strings.Count(string(content), "\n") > 10 {
		t.Error("production build JS doesn't appear to be properly minified")
	}

	// Verify source map
	sourceMap := filepath.Join(dir, "dist/sourcemaps/main.js.map")
	if _, err := os.Stat(sourceMap); err != nil {
		t.Errorf("source map file missing: %v", err)
	}
}
