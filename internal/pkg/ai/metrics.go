package ai

import (
	"sync"
)

type MetricsCollector interface {
	Collect() []Metric
	Add(metric Metric)
	Reset()
}

type metricsCollector struct {
	metrics []Metric
	mu      sync.RWMutex
}

func newMetricsCollector() MetricsCollector {
	return &metricsCollector{
		metrics: make([]Metric, 0),
	}
}

func (c *metricsCollector) Collect() []Metric {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metrics := make([]Metric, len(c.metrics))
	copy(metrics, c.metrics)
	return metrics
}

func (c *metricsCollector) Add(metric Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = append(c.metrics, metric)
}

func (c *metricsCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = make([]Metric, 0)
}
