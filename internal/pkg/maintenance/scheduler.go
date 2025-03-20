package maintenance

import (
    "context"
    "time"
)

type MaintenanceScheduler struct {
    dependency  *DependencyManager
    performance *PerformanceMonitor
    interval    time.Duration
}

func NewMaintenanceScheduler(dm *DependencyManager, pm *PerformanceMonitor) *MaintenanceScheduler {
    return &MaintenanceScheduler{
        dependency:  dm,
        performance: pm,
        interval:    24 * time.Hour,
    }
}

func (s *MaintenanceScheduler) Start(ctx context.Context) error {
    ticker := time.NewTicker(s.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            if err := s.runMaintenance(ctx); err != nil {
                s.logError(err)
            }
        }
    }
}