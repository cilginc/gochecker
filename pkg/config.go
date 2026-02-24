package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// CheckConfig parses the YAML configuration from the specified path or ".gochecker.yaml" by default.
func CheckConfig(path ...string) (*Config, error) {
	targetPath := DEFAULT_CONFIG_FILE
	if len(path) > 0 && path[0] != "" {
		targetPath = path[0]
	}

	data, err := os.ReadFile(targetPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", ErrConfigNotFound, targetPath)
		}
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}

	decoder := yaml.NewDecoder(
		bytes.NewReader(data),
		yaml.DisallowUnknownField(),
	)

	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	err = validateConfig(&cfg)
	if err != nil {
		return &cfg, err
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if len(cfg.Packages) == 0 {
		return fmt.Errorf("%w: no packages defined", ErrInvalidConfig)
	}
	return nil
}
