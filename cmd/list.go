package cmd

import (
	"fmt"

	"github.com/cilginc/gochecker/internal/config"
	"github.com/cilginc/gochecker/internal/output"
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
	paths, err := config.GetConfigPaths(recursive, cfgFile, recursiveDir)
	if err != nil {
		return err
	}
	var allPackages []pkg.Package

	for _, cfgPath := range paths {
		cfg, err := pkg.CheckConfig(cfgPath)
		if err != nil {
			_ = ui.CliError("%s: %s", cfgPath, err)
			continue
		}

		if err := cfg.LoadVersions(); err != nil {
			ui.CliWarn("No version history found for %s.", cfgPath)
		}

		allPackages = append(allPackages, cfg.Packages...)
	}

	if outputFormat != "text" {
		return output.RenderPackages(outputFormat, allPackages)
	}

	ui.CliInfo("Tracked Packages (All Configurations):")
	fmt.Println("--------------------------------------------------")

	if len(allPackages) == 0 {
		fmt.Println("No packages found in any configuration.")
		return nil
	}

	for _, p := range allPackages {
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
	fmt.Printf("Total: %d packages monitored across all files.\n", len(allPackages))

	return nil
}
