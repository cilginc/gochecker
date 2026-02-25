package cmd

import (
	"fmt"

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
	ui.CliInfo("Validating configuration: %s", cfgFile)

	cfg, err := pkg.CheckConfig(cfgFile)
	if err != nil {
		return ui.CliError("Configuration invalid: %v", err)
	}

	if len(cfg.Packages) == 0 {
		ui.CliWarn("The configuration file is valid but contains no packages.")
		return nil
	}

	ui.CliSuccess("Syntax check passed! Found %d well-formatted packages.", len(cfg.Packages))

	for _, p := range cfg.Packages {
		fmt.Printf("  %s %s\n", ui.Success("✔"), p.Name)
	}

	return nil
}
