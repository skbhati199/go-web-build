package ai

type OptimizedConfig struct {
	CacheStrategy   string
	BuildParameters map[string]interface{}
	ResourceLimits  ResourceLimits
	Metadata        BuildMetadata
}

type BuildMetadata struct {
	Confidence      float64
	Recommendations []string
	Warnings        []string
	ModelVersion    string
	LastOptimized   int64
}
