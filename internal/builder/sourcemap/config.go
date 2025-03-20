package sourcemap

type SourceMapMode string
type SourceMapType string

const (
	DevelopmentMode SourceMapMode = "development"
	ProductionMode  SourceMapMode = "production"

	InlineType   SourceMapType = "inline"
	ExternalType SourceMapType = "external"
	BothType     SourceMapType = "both"
)

type SourceMapConfig struct {
	Mode           SourceMapMode       `json:"mode"`
	Type           SourceMapType       `json:"type"`
	IncludeContent bool                `json:"includeContent"`
	SourceRoot     string              `json:"sourceRoot"`
	PathRewrites   []PathRewrite       `json:"pathRewrites"`
	Compression    bool                `json:"compression"`
	PrivateStorage bool                `json:"privateStorage"`
	ErrorTracking  ErrorTrackingConfig `json:"errorTracking"`
}

type PathRewrite struct {
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`
}

type ErrorTrackingConfig struct {
	Provider   string `json:"provider"` // sentry, rollbar
	ProjectID  string `json:"projectId"`
	AuthToken  string `json:"authToken"`
	SourceRoot string `json:"sourceRoot"`
}
