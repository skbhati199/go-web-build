package maintenance

import (
	"context"
	"encoding/json"
	"os/exec"
)

type DependencyManager struct {
	projectPath string
	logger      Logger
}

func NewDependencyManager(projectPath string, logger Logger) *DependencyManager {
	return &DependencyManager{
		projectPath: projectPath,
		logger:      logger,
	}
}

func (m *DependencyManager) CheckUpdates(ctx context.Context) ([]Update, error) {
	cmd := exec.CommandContext(ctx, "npm", "outdated", "--json")
	cmd.Dir = m.projectPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var updates map[string]PackageInfo
	if err := json.Unmarshal(output, &updates); err != nil {
		return nil, err
	}

	return m.processUpdates(updates), nil
}

func (m *DependencyManager) ApplySecurityPatches(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "npm", "audit", "fix")
	cmd.Dir = m.projectPath
	return cmd.Run()
}
