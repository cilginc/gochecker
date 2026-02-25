package pkg

// Package represents a version check target.
type Package struct {
	Name     string   `yaml:"name"    json:"name"`
	Version  string   `yaml:"version" json:"version"`
	Provider Provider `yaml:",inline" json:",inline"`

	// Global Options
	Prefix string `yaml:"prefix,omitempty" json:"prefix,omitempty"`
}

type VersionFile struct {
	Packages map[string]string `json:"packages"`
}

// Provider defines which upstream source is used.
type Provider struct {
	GitHub *GitHub `yaml:"github,omitempty" json:"github,omitempty"`
}

// GitHub provider configuration.
type GitHub struct {
	Repo   string `yaml:"repo"   json:"repo"`
	Branch string `yaml:"branch" json:"branch"`
}

type Config struct {
	Packages []Package `yaml:"packages" json:"packages"`
}

type Result struct {
	Name       string
	OldVersion string
	NewVersion string
	Updated    bool
	Error      error
}

const DEFAULT_CONFIG_FILE = ".gochecker.yaml"
const DEFAULT_VERSIONS_FILE = ".gochecker-lock.json"

const GITHUB_PAT_TOKEN_ENV_VAR = "GITHUB_TOKEN"
