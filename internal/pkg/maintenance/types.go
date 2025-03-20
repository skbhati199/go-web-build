package maintenance

import (
	"context"
	"time"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type PackageInfo struct {
	Current  string `json:"current"`
	Wanted   string `json:"wanted"`
	Latest   string `json:"latest"`
	Location string `json:"location"`
}

type Update struct {
	Package     string
	FromVersion string
	ToVersion   string
	Type        string // major, minor, patch
	Breaking    bool
}

type PerformanceReport struct {
	BuildMetrics    BuildMetrics
	RuntimeMetrics  RuntimeMetrics
	Recommendations []string
	Timestamp       time.Time
}

type BuildMetrics struct {
	Duration   float64
	BundleSize int64
	Chunks     int
}

type RuntimeMetrics struct {
	LoadTime    float64
	FirstPaint  float64
	Interactive float64
}

type MetricsCollector interface {
	Collect(ctx context.Context) (*Metrics, error)
}

type Metrics struct {
	Build    BuildMetrics
	Runtime  RuntimeMetrics
	Warnings []string
}

type VulnerabilityScanner interface {
	Scan(ctx context.Context) ([]Vulnerability, error)
}

type SecurityPatcher interface {
	ApplyPatches(ctx context.Context) error
}

type SecurityNotifier interface {
	NotifyVulnerability(vuln Vulnerability) error
}

type Vulnerability struct {
	ID          string
	Package     string
	Severity    string
	Description string
	FixVersion  string
}
