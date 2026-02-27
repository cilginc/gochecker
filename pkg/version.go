package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func (cfg *Config) LoadVersions(path ...string) error {
	targetPath := DEFAULT_VERSIONS_FILE
	if len(path) > 0 && path[0] != "" {
		targetPath = path[0]
	}

	data, err := os.ReadFile(targetPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// No versions file yet that's fine
			// [TODO]: probably not fine
			return nil
		}
		return fmt.Errorf("%w: %v", ErrVersionsRead, err)
	}

	var vf VersionFile
	if err := json.Unmarshal(data, &vf); err != nil {
		return fmt.Errorf("%w: %v", ErrVersionsParse, err)
	}

	for i := range cfg.Packages {
		if v, ok := vf.Packages[cfg.Packages[i].Name]; ok {
			cfg.Packages[i].Version = v
		}
	}

	return nil
}

func (cfg *Config) SaveVersions(path ...string) error {
	targetPath := DEFAULT_VERSIONS_FILE
	if len(path) > 0 && path[0] != "" {
		targetPath = path[0]
	}

	vf := VersionFile{
		Packages: make(map[string]string),
	}

	for _, p := range cfg.Packages {
		if p.Version != "" {
			vf.Packages[p.Name] = p.Version
		}
	}

	if len(vf.Packages) == 0 {
		return nil
	}

	data, err := json.MarshalIndent(vf, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrVersionsParse, err)
	}

	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return fmt.Errorf("%w: %v", ErrVersionsWrite, err)
	}

	return nil
}
