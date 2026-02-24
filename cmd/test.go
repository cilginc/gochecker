package cmd

import (
	"fmt"

	"github.com/cilginc/gochecker/internal/config"
	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:          "test",
	Short:        "Test configuration parsing",
	SilenceUsage: true,
	RunE:         testYAML,
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func testYAML(cmd *cobra.Command, args []string) error {
	cfg, err := config.CheckConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)
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
