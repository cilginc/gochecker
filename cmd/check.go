package cmd

import (
	"context"
	"fmt"

	"github.com/cilginc/gochecker/internal/config"
	"github.com/cilginc/gochecker/internal/output"
	"github.com/cilginc/gochecker/internal/ui"
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

func runCheck(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	showNewOnly, _ := cmd.Flags().GetBool("new")

	paths, err := config.GetConfigPaths(recursive, cfgFile, recursiveDir)
	if err != nil {
		return err
	}

	for i, cfgPath := range paths {
		if i > 0 {
			fmt.Println()
		}

		results, err := engine.Check(ctx, cfgPath)
		if err != nil {
			_ = ui.CliError("%s: %s", cfgPath, err)
			continue
		}

		if outputFormat != "text" {
			if err := output.RenderResults(outputFormat, results); err != nil {
				return err
			}
			continue
		}

		ui.CliInfo("Scanning for updates using configuration: %s", cfgPath)

		foundUpdate := false

		for _, r := range results {

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
			continue
		}

		ui.CliSuccess("Check completed!")
	}

	return nil
}
