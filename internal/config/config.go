package config

import (
	"bytes"
	"os"

	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg"
	"github.com/goccy/go-yaml"
)

// LoadConfig parses the YAML configuration from the specified path or ".gochecker.yaml" by default.
func LoadConfig(path ...string) (*pkg.Config, error) {
	targetPath := ".gochecker.yaml"
	if len(path) > 0 && path[0] != "" {
		targetPath = path[0]
	}

	data, err := os.ReadFile(targetPath)
	if err != nil {
		return nil, ui.CliError("failed to read %s: %v", path, err)
	}

	decoder := yaml.NewDecoder(
		bytes.NewReader(data),
		yaml.DisallowUnknownField(),
	)

	var cfg pkg.Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, ui.CliError("invalid %s: %v", path, err)
	}

	err = validateConfig(&cfg)
	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}

func validateConfig(cfg *pkg.Config) error {
	if len(cfg.Packages) == 0 {
		return ui.CliError("no packages defined in config")
	}
	return nil
}
