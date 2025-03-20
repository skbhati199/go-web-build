package npm

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/sonukumar/go-web-build/internal/pkg/manager"
)

type NPM struct {
	workDir string
}

func New(workDir string) *NPM {
	return &NPM{workDir: workDir}
}

func (n *NPM) Install(ctx context.Context, packages ...string) error {
	args := append([]string{"install"}, packages...)
	cmd := exec.CommandContext(ctx, "npm", args...)
	cmd.Dir = n.workDir
	return cmd.Run()
}

func (n *NPM) Uninstall(ctx context.Context, packages ...string) error {
	args := append([]string{"uninstall"}, packages...)
	cmd := exec.CommandContext(ctx, "npm", args...)
	cmd.Dir = n.workDir
	return cmd.Run()
}

func (n *NPM) Update(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "npm", "update")
	cmd.Dir = n.workDir
	return cmd.Run()
}

func (n *NPM) List(ctx context.Context) ([]manager.Package, error) {
	cmd := exec.CommandContext(ctx, "npm", "list", "--json")
	cmd.Dir = n.workDir

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list packages: %w", err)
	}

	var npmList struct {
		Dependencies map[string]struct {
			Version         string            `json:"version"`
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		} `json:"dependencies"`
	}

	if err := json.Unmarshal(output, &npmList); err != nil {
		return nil, fmt.Errorf("failed to parse npm list: %w", err)
	}

	var packages []manager.Package
	for name, info := range npmList.Dependencies {
		pkg := manager.Package{
			Name:         name,
			Version:      info.Version,
			Dependencies: info.Dependencies,
			DevMode:      len(info.DevDependencies) > 0,
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}
