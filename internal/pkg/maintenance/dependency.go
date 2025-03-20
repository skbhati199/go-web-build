package maintenance

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
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

func (m *DependencyManager) processUpdates(updates map[string]PackageInfo) []Update {
	var result []Update
	for pkg, info := range updates {
		update := Update{
			Package:     pkg,
			FromVersion: info.Current,
			ToVersion:   info.Latest,
			Type:        determineUpdateType(info.Current, info.Latest),
			Breaking:    isBreakingChange(info.Current, info.Latest),
		}
		result = append(result, update)
		m.logger.Info("Found update", "package", pkg, "from", info.Current, "to", info.Latest)
	}
	return result
}

func determineUpdateType(current, latest string) string {
	currentParts := strings.Split(strings.TrimPrefix(current, "v"), ".")
	latestParts := strings.Split(strings.TrimPrefix(latest, "v"), ".")

	if len(currentParts) != 3 || len(latestParts) != 3 {
		return "unknown"
	}

	if currentParts[0] != latestParts[0] {
		return "major"
	}
	if currentParts[1] != latestParts[1] {
		return "minor"
	}
	return "patch"
}

func isBreakingChange(current, latest string) bool {
	return determineUpdateType(current, latest) == "major"
}

func (m *DependencyManager) UpdateDependencies(ctx context.Context, updates []Update) error {
	for _, update := range updates {
		if update.Breaking {
			m.logger.Info("Skipping breaking change", "package", update.Package)
			continue
		}

		cmd := exec.CommandContext(ctx, "npm", "install", fmt.Sprintf("%s@%s", update.Package, update.ToVersion))
		cmd.Dir = m.projectPath

		if err := cmd.Run(); err != nil {
			m.logger.Error("Failed to update package", "package", update.Package, "error", err)
			return fmt.Errorf("failed to update %s: %w", update.Package, err)
		}

		m.logger.Info("Updated package", "package", update.Package, "version", update.ToVersion)
	}
	return nil
}
