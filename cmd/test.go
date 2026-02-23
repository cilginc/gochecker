package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:          "test",
	Short:        "Test configuration parsing",
	SilenceUsage: true,
	RunE:         checkYAML,
}

func init() {
	rootCmd.AddCommand(testCmd)
}

// [TODO]: make this into pkg/
func checkYAML(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig(".gochecker.yaml")
	if err != nil {
		return err
	}

	if err := validateConfig(cfg); err != nil {
		return err
	}

	printConfig(cfg)
	return nil
}

func loadConfig(path string) (*pkg.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, cliError("failed to read %s: %v", path, err)
	}

	decoder := yaml.NewDecoder(
		bytes.NewReader(data),
		yaml.DisallowUnknownField(),
	)

	var cfg pkg.Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, cliError("invalid %s: %v", path, err)
	}

	return &cfg, nil
}

func validateConfig(cfg *pkg.Config) error {
	if len(cfg.Packages) == 0 {
		return cliError("no packages defined in config")
	}
	return nil
}

func printConfig(cfg *pkg.Config) {
	for _, p := range cfg.Packages {
		fmt.Println(ui.Title("Package:"), ui.Info(p.Name))

		if p.Provider.GitHub != nil {
			fmt.Println("  ", ui.Success("✔ Provider:"), "GitHub")
			fmt.Println("  ", ui.Info("Repo:"), p.Provider.GitHub.Repo)
			fmt.Println()
			continue
		}

		fmt.Println("  ", ui.Warn("⚠ You need to specify a valid provider."))
		fmt.Println()
	}
}

func cliError(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return fmt.Errorf("%s %s", ui.Err("✖ Error:"), msg)
}
