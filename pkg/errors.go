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
