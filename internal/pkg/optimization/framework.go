package optimization

import (
	"context"
)

type FrameworkOptimizer interface {
	Optimize(ctx context.Context, config OptimizationConfig) error
}

type OptimizationConfig struct {
	Framework   string
	Environment string
	Features    []string
	Performance PerformanceConfig
}

type PerformanceConfig struct {
	TreeShaking   bool
	CodeSplitting bool
	LazyLoading   bool
	Compression   CompressionConfig
}

type CompressionConfig struct {
	Enable    bool
	Level     int
	Algorithm string
}

func NewOptimizer(framework string) FrameworkOptimizer {
	switch framework {
	case "react":
		return newReactOptimizer()
	default:
		return newDefaultOptimizer()
	}
}
