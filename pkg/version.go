package pkg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func (cfg *Config) LoadVersions(path ...string) error {
	targetPath := DEFAULT_VERSIONS_FILE
	if len(path) > 0 && path[0] != "" {
		targetPath = path[0]
	}

	// [TODO]: Clean the code here.
	srcName, srcVersion, err := checkSRCINFO()
	if err != nil {
		fmt.Printf("Warning: Could not sync from .SRCINFO: %v\n", err)
	} else {
		for i := range cfg.Packages {
			if cfg.Packages[i].Name == srcName {
				cfg.Packages[i].Version = srcVersion
				break
			}
		}
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

func checkSRCINFO(path ...string) (name string, version string, err error) {
	targetPath := DEFAULT_SRCINFO_FILE
	if len(path) > 0 && path[0] != "" {
		targetPath = path[0]
	}

	file, err := os.Open(targetPath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "pkgbase", "pkgname":
			if name == "" || key == "pkgbase" {
				name = value
			}
		case "pkgver":
			version = value
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	if name == "" || version == "" {
		return "", "", fmt.Errorf(
			"metadata missing in %s (name: %q, ver: %q)",
			targetPath,
			name,
			version,
		)
	}

	return name, version, nil
}
