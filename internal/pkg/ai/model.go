package ai

type Model interface {
	Predict(metrics []Metric) Prediction
	Train(data []TrainingData) error
	Evaluate() ModelMetrics
}

type Prediction struct {
	CacheStrategy string
	BuildParams   map[string]interface{}
	Resources     ResourceLimits
	Confidence    float64
}

type ResourceLimits struct {
	CPU    int
	Memory int64
	Disk   int64
}

type Metric struct {
	Name   string
	Value  float64
	Labels map[string]string
}

type TrainingData struct {
	Features    []float64
	Labels      []string
	Timestamp   int64
	Performance BuildPerformance
}

type BuildPerformance struct {
	Duration      float64
	SuccessRate   float64
	ResourceUsage ResourceUsage
}

type ResourceUsage struct {
	CPUUsage    float64
	MemoryUsage float64
	DiskIO      float64
}

type ModelMetrics struct {
	Accuracy    float64
	Precision   float64
	Recall      float64
	F1Score     float64
	LastUpdated int64
}
