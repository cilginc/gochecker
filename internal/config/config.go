package config

import (
	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg"
)

func CheckConfig(path ...string) (*pkg.Config, error) {
	cfg, err := pkg.CheckConfig()
	if err != nil {
		ui.CliError("%v", err)
		return nil, err
	}
	return cfg, nil
}
