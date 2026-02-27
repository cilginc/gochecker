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

func runUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	paths, err := config.GetConfigPaths(recursive, cfgFile, recursiveDir)
	if err != nil {
		return err
	}

	var allResults []pkg.Result

	for _, path := range paths {
		res, err := engine.Execute(ctx, path)
		if err != nil {
			_ = ui.CliError("Failed to update %s: %v", path, err)
			continue
		}
		allResults = append(allResults, res...)
	}

	if outputFormat != "text" {
		return output.RenderResults(outputFormat, allResults)
	}

	ui.CliInfo("Processing updates for %d packages across all configurations...", len(allResults))
	fmt.Println("--------------------------------------------------")

	updatedCount := 0
	for _, v := range allResults {
		if v.Error != nil {
			fmt.Printf("%s %-15s: %v\n", ui.Warn("✖"), v.Name, v.Error)
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
		ui.CliSuccess("%d packages were updated and synchronized.", updatedCount)
	} else if len(allResults) > 0 {
		ui.CliSuccess("All records are already in sync. No changes needed.")
	} else {
		ui.CliWarn("No packages were found to process.")
	}

	return nil
}
