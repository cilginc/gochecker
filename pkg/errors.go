package pkg

import "errors"

var (
	ErrInvalidConfig   = errors.New("invalid config")
	ErrUnknownProvider = errors.New("unknown provider")
	ErrConfigNotFound  = errors.New("config file not found")
	ErrConfigRead      = errors.New("failed to read configuration file")
	ErrVersionsRead    = errors.New("failed to read versions file")
	ErrVersionsWrite   = errors.New("failed to write versions file")
	ErrVersionsParse   = errors.New("failed to parse versions file")
)

var (
	// GitHub provider errors
	ErrGitHubRequest   = errors.New("github request failed")
	ErrGitHubRateLimit = errors.New("github rate limit exceeded")
	ErrGitHubForbidden = errors.New("github access forbidden")
	ErrGitHubStatus    = errors.New("github unexpected status code")
	ErrGitHubDecode    = errors.New("github response decode failed")
)

var (
	// PyPI provider errors
	ErrPyPIRequest      = errors.New("pypi request failed")
	ErrPyPINotFound     = errors.New("pypi package not found")
	ErrPyPIStatus       = errors.New("pypi unexpected status code")
	ErrPyPIDecode       = errors.New("pypi response decode failed")
	ErrPyPIEmptyVersion = errors.New("pypi returned an empty version")
)

var (
	// OCI provider errors
	ErrOCIRequest  = errors.New("oci request failed")
	ErrOCINotFound = errors.New("oci image not found")
	ErrOCIDecode   = errors.New("oci response decode failed")
)

var (
	// AUR provider errors
	ErrAURRequest  = errors.New("aur request failed")
	ErrAURNotFound = errors.New("aur package not found")
	ErrAURDecode   = errors.New("aur response decode failed")
)

var (
	// Git provider errors
	ErrGitRequest        = errors.New("git remote request failed")
	ErrGitBranchNotFound = errors.New("git branch not found")
	ErrGitNoTags         = errors.New("no tags found in repository")
)
