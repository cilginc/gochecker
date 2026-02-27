package cmd

import (
	"context"
	"fmt"

	"github.com/cilginc/gochecker/internal/config"
	"github.com/cilginc/gochecker/internal/output"
	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg"
	"github.com/cilginc/gochecker/pkg/engine"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for new versions of tracked packages",
	Long: ui.Info(`Scan all software sources defined in your configuration file 
to detect new versions. It compares upstream versions with your local 
records and highlights updates.`),
	Example: ui.Success("  gochecker check\n") +
		ui.Success("  gochecker check --config custom.yaml --output json"),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolP("new", "n", false,
		"Only display packages that have a newer version available")
}

// [TODO]: Clean the code here. 
func runCheck(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	showNewOnly, _ := cmd.Flags().GetBool("new")

	paths, err := config.GetConfigPaths(recursive, cfgFile, recursiveDir)
	if err != nil {
		return ui.CliError("failed to get config paths: %w", err)
	}

	if len(paths) == 0 {
		ui.CliWarn("No configuration files found.")
		return nil
	}

	allConfig := &pkg.Config{
		Packages: []pkg.Package{},
	}

	var processErrors []error

	for _, path := range paths {
		cfg, err := pkg.CheckConfig(path)
		if err != nil {
			errMsg := fmt.Errorf("[%s] could not be read: %v", path, err)
			_ = ui.CliError("%v", errMsg)
			processErrors = append(processErrors, errMsg)
			continue
		}
		allConfig.Packages = append(allConfig.Packages, cfg.Packages...)
	}

	if len(allConfig.Packages) == 0 {
		if len(processErrors) > 0 {
			return fmt.Errorf(
				"failed to process any packages; %d files had errors",
				len(processErrors),
			)
		}
		ui.CliWarn("No packages found to process.")
		return nil
	}

	if err := allConfig.LoadVersions(); err != nil {
		return fmt.Errorf("failed to load version history: %w", err)
	}

	ui.CliInfo("Checking updates for %d packages...", len(allConfig.Packages))
	result := engine.Run(ctx, allConfig.Packages)

	if outputFormat != "text" {
		return output.RenderResults(outputFormat, result)
	}

	ui.CliInfo("Scanning for updates using configuration: %s", cfgFile)

	foundUpdate := false

	for _, r := range result {
		if r.Error != nil {
			_ = ui.CliError("%s: %v", r.Name, r.Error)
			continue
		}

		if r.Updated {
			foundUpdate = true
			ui.CliInfo("%s: %s → %s",
				r.Name,
				r.OldVersion,
				r.NewVersion,
			)
		} else if !showNewOnly {
			ui.CliSuccess("%s: up to date (%s)",
				r.Name,
				r.OldVersion,
			)
		}
	}

	if showNewOnly && !foundUpdate {
		ui.CliSuccess("No updates found.")
		return nil
	}

	ui.CliSuccess("Check completed!")
	return nil
}
