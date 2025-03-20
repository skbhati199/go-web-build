package maintenance

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type SecurityManager struct {
	scanner  VulnerabilityScanner
	patcher  SecurityPatcher
	notifier SecurityNotifier
	logger   Logger
}

func NewSecurityManager(scanner VulnerabilityScanner, patcher SecurityPatcher, notifier SecurityNotifier, logger Logger) *SecurityManager {
	return &SecurityManager{
		scanner:  scanner,
		patcher:  patcher,
		notifier: notifier,
		logger:   logger,
	}
}

func (s *SecurityManager) ScanVulnerabilities(ctx context.Context) error {
	// Run Go vulnerability check
	goCmd := exec.CommandContext(ctx, "govulncheck", "./...")
	_, err := goCmd.CombinedOutput()
	if err != nil {
		s.logger.Error("Go vulnerability check failed", "error", err)
		return fmt.Errorf("go vulnerability check failed: %w", err)
	}

	// Run npm audit
	npmCmd := exec.CommandContext(ctx, "npm", "audit", "--json")
	_, err = npmCmd.CombinedOutput()
	if err != nil {
		s.logger.Error("NPM audit failed", "error", err)
		return fmt.Errorf("npm audit failed: %w", err)
	}

	// Process vulnerabilities through scanner
	vulns, err := s.scanner.Scan(ctx)
	if err != nil {
		return fmt.Errorf("vulnerability scan failed: %w", err)
	}

	// Notify about found vulnerabilities
	for _, vuln := range vulns {
		if err := s.notifier.NotifyVulnerability(vuln); err != nil {
			s.logger.Error("Failed to notify vulnerability", "id", vuln.ID, "error", err)
		}
	}

	return s.processVulnerabilities(vulns)
}

func (s *SecurityManager) ApplySecurityPatches(ctx context.Context) error {
	// Apply npm security fixes
	npmCmd := exec.CommandContext(ctx, "npm", "audit", "fix")
	if output, err := npmCmd.CombinedOutput(); err != nil {
		s.logger.Error("NPM security fix failed", "error", err, "output", string(output))
		return fmt.Errorf("npm security fix failed: %w", err)
	}

	// Apply patches through patcher
	if err := s.patcher.ApplyPatches(ctx); err != nil {
		return fmt.Errorf("failed to apply security patches: %w", err)
	}

	s.logger.Info("Security patches applied successfully")
	return nil
}

func (s *SecurityManager) processVulnerabilities(vulns []Vulnerability) error {
	var criticalVulns []Vulnerability
	var highVulns []Vulnerability

	for _, vuln := range vulns {
		switch strings.ToLower(vuln.Severity) {
		case "critical":
			criticalVulns = append(criticalVulns, vuln)
		case "high":
			highVulns = append(highVulns, vuln)
		}
	}

	if len(criticalVulns) > 0 {
		s.logger.Error("Critical vulnerabilities found", "count", len(criticalVulns))
		return fmt.Errorf("found %d critical vulnerabilities", len(criticalVulns))
	}

	if len(highVulns) > 0 {
		s.logger.Info("High severity vulnerabilities found", "count", len(highVulns))
	}

	return nil
}
