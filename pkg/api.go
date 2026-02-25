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
	PyPI   *PyPI   `yaml:"pypi,omitempty"   json:"pypi,omitempty"`
	OCI    *OCI    `yaml:"oci,omitempty"    json:"oci,omitempty"`
	AUR    *AUR    `yaml:"aur,omitempty"    json:"aur,omitempty"`
}

type Result struct {
	Name       string
	OldVersion string
	NewVersion string
	Updated    bool
	Error      error
}

type Config struct {
	Packages []Package `yaml:"packages" json:"packages"`
}

const DEFAULT_CONFIG_FILE = ".gochecker.yaml"
const DEFAULT_VERSIONS_FILE = ".gochecker-lock.json"

const GITHUB_PAT_TOKEN_ENV_VAR = "GITHUB_TOKEN"

type GitHub struct {
	// Required: "owner/repo"
	Repo string `yaml:"repo" json:"repo"`

	// Optional
	Branch string `yaml:"branch,omitempty" json:"branch,omitempty"`
	Path   string `yaml:"path,omitempty"   json:"path,omitempty"`

	// For GitHub Enterprise (example: github.example.com)
	Host string `yaml:"host,omitempty" json:"host,omitempty"`

	// Release options
	UseLatestRelease  bool `yaml:"use_latest_release,omitempty"  json:"use_latest_release,omitempty"`
	UseMaxRelease     bool `yaml:"use_max_release,omitempty"     json:"use_max_release,omitempty"`
	UseReleaseName    bool `yaml:"use_release_name,omitempty"    json:"use_release_name,omitempty"`
	IncludePrerelease bool `yaml:"include_prereleases,omitempty" json:"include_prereleases,omitempty"`

	// Tag options
	UseLatestTag bool   `yaml:"use_latest_tag,omitempty" json:"use_latest_tag,omitempty"`
	UseMaxTag    bool   `yaml:"use_max_tag,omitempty"    json:"use_max_tag,omitempty"`
	Query        string `yaml:"query,omitempty"          json:"query,omitempty"`
}

type PyPI struct {
	Package string `yaml:"package" json:"package"`
}

type OCI struct {
	Image string `yaml:"image" json:"image"`
}

type AUR struct {
	Package string `yaml:"package"                 json:"package"`
	// Strips the pkgrel. (example: 1.0.2-1 -> 1.0.2)
	StripRelease bool `yaml:"strip_release,omitempty" json:"strip_release,omitempty"`
}
