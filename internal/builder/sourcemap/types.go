package sourcemap

// SourceMap represents the structure of a source map file
type SourceMap struct {
	Version        int      `json:"version"`
	Sources        []string `json:"sources"`
	Names          []string `json:"names"`
	Mappings       string   `json:"mappings"`
	File           string   `json:"file"`
	SourceRoot     string   `json:"sourceRoot,omitempty"`
	SourcesContent []string `json:"sourcesContent,omitempty"`
	Url            string   `json:"-"` // Internal use only, not marshaled to JSON
}
