package maintenance

import (
	"context"
	"fmt"
	"time"
)

type PerformanceMonitor struct {
	metrics    MetricsCollector
	thresholds map[string]float64
	logger     Logger
}

func NewPerformanceMonitor(metrics MetricsCollector, logger Logger) *PerformanceMonitor {
	return &PerformanceMonitor{
		metrics: metrics,
		logger:  logger,
		thresholds: map[string]float64{
			"buildTime":   120,  // seconds
			"bundleSize":  1024, // KB
			"loadTime":    3,    // seconds
			"interactive": 5,    // seconds
			"firstPaint":  2,    // seconds
		},
	}
}

func (m *PerformanceMonitor) AnalyzePerformance(ctx context.Context) (*PerformanceReport, error) {
	metrics, err := m.metrics.Collect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect metrics: %w", err)
	}

	recommendations := m.generateRecommendations(metrics)
	m.logPerformanceIssues(metrics, recommendations)

	return &PerformanceReport{
		BuildMetrics:    metrics.Build,
		RuntimeMetrics:  metrics.Runtime,
		Recommendations: recommendations,
		Timestamp:       time.Now(),
	}, nil
}

func (m *PerformanceMonitor) generateRecommendations(metrics *Metrics) []string {
	var recommendations []string

	// Check build metrics
	if metrics.Build.Duration > m.thresholds["buildTime"] {
		recommendations = append(recommendations,
			fmt.Sprintf("Build time (%.2fs) exceeds threshold (%.2fs). Consider optimizing build configuration.",
				metrics.Build.Duration, m.thresholds["buildTime"]))
	}

	if float64(metrics.Build.BundleSize) > m.thresholds["bundleSize"] {
		recommendations = append(recommendations,
			fmt.Sprintf("Bundle size (%.2f KB) exceeds threshold (%.2f KB). Consider code splitting or tree shaking.",
				float64(metrics.Build.BundleSize), m.thresholds["bundleSize"]))
	}

	// Check runtime metrics
	if metrics.Runtime.LoadTime > m.thresholds["loadTime"] {
		recommendations = append(recommendations,
			fmt.Sprintf("Page load time (%.2fs) exceeds threshold (%.2fs). Consider lazy loading or optimizing assets.",
				metrics.Runtime.LoadTime, m.thresholds["loadTime"]))
	}

	if metrics.Runtime.FirstPaint > m.thresholds["firstPaint"] {
		recommendations = append(recommendations,
			fmt.Sprintf("First paint time (%.2fs) exceeds threshold (%.2fs). Consider optimizing critical rendering path.",
				metrics.Runtime.FirstPaint, m.thresholds["firstPaint"]))
	}

	if metrics.Runtime.Interactive > m.thresholds["interactive"] {
		recommendations = append(recommendations,
			fmt.Sprintf("Time to interactive (%.2fs) exceeds threshold (%.2fs). Consider reducing JavaScript bundle size.",
				metrics.Runtime.Interactive, m.thresholds["interactive"]))
	}

	return recommendations
}

func (m *PerformanceMonitor) logPerformanceIssues(metrics *Metrics, recommendations []string) {
	if len(recommendations) > 0 {
		m.logger.Info("Performance issues detected", "count", len(recommendations))
		for _, rec := range recommendations {
			m.logger.Info("Performance recommendation", "message", rec)
		}
	}

	// Log warnings if any
	for _, warning := range metrics.Warnings {
		m.logger.Info("Performance warning", "message", warning)
	}
}

func (m *PerformanceMonitor) UpdateThresholds(newThresholds map[string]float64) {
	for key, value := range newThresholds {
		if value > 0 {
			m.thresholds[key] = value
		}
	}
}
