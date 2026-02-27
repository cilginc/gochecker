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

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Sync local version records",
	Long: ui.Info(`Synchronize your local version database with the current configuration. 
This command ensures that your version tracking file is up-to-date 
with the entries defined in your configuration without checking upstream.`),
	Example:       ui.Success("  gochecker update"),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

// [TODO]: Clean the code here.
// We can add engine.execute which want []string paths.
func runUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
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

	if err := allConfig.SaveVersions(); err != nil {
		return fmt.Errorf("failed to save updated versions: %w", err)
	}

	if outputFormat != "text" {
		return output.RenderResults(outputFormat, result)
	}

	fmt.Println("--------------------------------------------------")

	updatedCount := 0
	for _, v := range result {
		if v.Error != nil {
			fmt.Printf("%s %-15s: %v\n", ui.Warn("✖"), v.Name, ui.Err(v.Error))
			continue
		}

		if v.Updated {
			fmt.Printf("%s %-15s: %s -> %s\n",
				ui.Success("↑"),
				v.Name,
				ui.Warn(v.OldVersion),
				ui.Success(v.NewVersion),
			)
			updatedCount++
		} else {
			fmt.Printf("%s %-15s: %s %s\n",
				ui.Info("-"),
				v.Name,
				v.OldVersion,
				ui.Info("(synced)"),
			)
		}
	}

	fmt.Println("--------------------------------------------------")

	if updatedCount > 0 {
		ui.CliSuccess("%d packages successfully updated and synchronized.", updatedCount)
	} else {
		ui.CliSuccess("All packages are already up to date.")
	}

	if len(processErrors) > 0 {
		ui.CliWarn(
			"\nNote: %d configuration file(s) were skipped due to errors.",
			len(processErrors),
		)
	}

	return nil
}
