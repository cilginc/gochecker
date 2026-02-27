package config

import (
	"os"
	"path/filepath"

	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg"
)

func CheckConfig(path ...string) (*pkg.Config, error) {
	cfg, err := pkg.CheckConfig()
	if err != nil {
		_ = ui.CliError("%v", err)
		return nil, err
	}
	return cfg, nil
}

func FindConfigFiles(dir string) ([]string, error) {
	var configs []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && info.Name() == pkg.DEFAULT_CONFIG_FILE {
			configs = append(configs, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(configs) == 0 {
		return nil, ui.CliError("no %s files found in %s", pkg.DEFAULT_CONFIG_FILE, dir)
	}

	return configs, nil
}

func GetConfigPaths(recursive bool, format string, dir string) ([]string, error) {
	if !recursive {
		return []string{format}, nil
	}

	if dir == "" {
		dir = "."
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, ui.CliError("invalid directory: %v", err)
	}

	info, err := os.Stat(absDir)
	if err != nil {
		return nil, ui.CliError("cannot access directory: %v", err)
	}
	if !info.IsDir() {
		return nil, ui.CliError("%s is not a directory", absDir)
	}

	return FindConfigFiles(absDir)
}
