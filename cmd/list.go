package cmd

import (
	"fmt"

	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracked packages",
	Long: ui.Info(`Display a comprehensive list of all software packages defined 
in your configuration. This shows the package names and their 
currently recorded versions without performing an upstream check.`),
	Example: ui.Success("  gochecker list\n") +
		ui.Success("  gochecker list --config custom.yaml"),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := pkg.CheckConfig(cfgFile)
	if err != nil {
		return ui.CliError("%s", err)
	}

	if err := cfg.LoadVersions(); err != nil {
		ui.CliWarn("No version history found. Showing names only.")
	}

	ui.CliInfo("Tracked Packages in %s:", cfgFile)
	fmt.Println("--------------------------------------------------")

	if len(cfg.Packages) == 0 {
		fmt.Println("No packages found in the configuration.")
		return nil
	}

	for _, p := range cfg.Packages {
		version := p.Version
		if version == "" {
			version = ui.Warn("no version recorded")
		}

		fmt.Printf("%s %-20s %s\n",
			ui.Info("•"),
			ui.Title(p.Name),
			ui.Success(version))
	}

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Total: %d packages monitored.\n", len(cfg.Packages))

	return nil
}
