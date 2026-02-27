package cmd

import (
	"fmt"

	"github.com/cilginc/gochecker/internal/config"
	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Validate the configuration syntax",
	Long: ui.Info(`Perform a syntax check on your configuration file. 
This command ensures the YAML structure is valid and the package 
definitions follow the required schema without performing any live checks.`),
	Example: ui.Success("  gochecker test\n") +
		ui.Success("  gochecker test --config my-config.yaml"),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runTest,
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func runTest(cmd *cobra.Command, args []string) error {
	paths, err := config.GetConfigPaths(recursive, cfgFile, recursiveDir)
	if err != nil {
		return err
	}

	for i, cfgPath := range paths {
		if i > 0 {
			fmt.Println()
		}

		ui.CliInfo("Validating configuration: %s", cfgPath)

		cfg, err := pkg.CheckConfig(cfgPath)
		if err != nil {
			_ = ui.CliError("Configuration invalid: %v (%s)", err, cfgPath)
			continue
		}

		if len(cfg.Packages) == 0 {
			ui.CliWarn("The configuration file is valid but contains no packages. (%s)", cfgPath)
			continue
		}

		ui.CliSuccess("Syntax check passed! Found %d well-formatted packages.", len(cfg.Packages))

		for _, p := range cfg.Packages {
			fmt.Printf("  %s %s\n", ui.Success("✔"), p.Name)
		}
	}

	return nil
}
