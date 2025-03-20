package ai

type LearningConfig struct {
	LearningRate float64
	BatchSize    int
	Epochs       int
	ModelPath    string
	Features     []FeatureConfig
}

type FeatureConfig struct {
	Name       string
	Weight     float64
	Normalized bool
	Required   bool
}

type BuildConfig struct {
	ProjectType  string
	Dependencies []string
	Environment  map[string]string
	Constraints  BuildConstraints
}

type BuildConstraints struct {
	MaxDuration int64
	MaxMemory   int64
	MaxCPU      int
	Priority    int
}
