package pkg

import (
	"errors"
	"net/http"
)

// Package represents a version check target.
type Package struct {
	Name     string   `yaml:"name"    json:"name"`
	Provider Provider `yaml:",inline" json:",inline"`
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

// [TODO]: improve error messages
var (
	ErrInvalidConfig   = errors.New("invalid config")
	ErrUnknownProvider = errors.New("unknown provider")
)

type Client struct {
	workers int
	http    *http.Client
}
