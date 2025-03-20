package maintenance

import (
	"context"
	"fmt"
	"time"
)

type MaintenanceScheduler struct {
	dependency  *DependencyManager
	performance *PerformanceMonitor
	interval    time.Duration
	logger      Logger
}

func NewMaintenanceScheduler(dm *DependencyManager, pm *PerformanceMonitor, logger Logger) *MaintenanceScheduler {
	return &MaintenanceScheduler{
		dependency:  dm,
		performance: pm,
		interval:    24 * time.Hour,
		logger:      logger,
	}
}

func (s *MaintenanceScheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run initial maintenance
	if err := s.runMaintenance(ctx); err != nil {
		s.logError("Initial maintenance failed", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.runMaintenance(ctx); err != nil {
				s.logError("Scheduled maintenance failed", err)
			}
		}
	}
}

func (s *MaintenanceScheduler) runMaintenance(ctx context.Context) error {
	// Check for dependency updates
	updates, err := s.dependency.CheckUpdates(ctx)
	if err != nil {
		return fmt.Errorf("dependency check failed: %w", err)
	}

	if len(updates) > 0 {
		s.logger.Info("Found dependency updates", "count", len(updates))
		if err := s.dependency.ApplySecurityPatches(ctx); err != nil {
			return fmt.Errorf("security patch failed: %w", err)
		}
	}

	// Run performance analysis
	report, err := s.performance.AnalyzePerformance(ctx)
	if err != nil {
		return fmt.Errorf("performance analysis failed: %w", err)
	}

	if len(report.Recommendations) > 0 {
		s.logger.Info("Performance recommendations available", "count", len(report.Recommendations))
	}

	return nil
}

func (s *MaintenanceScheduler) SetInterval(duration time.Duration) {
	if duration > 0 {
		s.interval = duration
	}
}

func (s *MaintenanceScheduler) logError(msg string, err error) {
	s.logger.Error(msg, "error", err)
}
