package cmd

import (
	"context"
	"fmt"

	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg/engine"
	"github.com/spf13/cobra"
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "github",
	Long:  `github`,
	RunE:  runThis,
}

func init() {
	rootCmd.AddCommand(githubCmd)
}

func runThis(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	results, err := engine.Execute(ctx)
	if err != nil {
		return err
	}

	for _, r := range results {
		if r.Error != nil {
			ui.CliError("%s: %v", r.Name, r.Error)
			continue
		}

		if r.Updated {
			fmt.Println(ui.Title(r.Name), ui.Success(r.NewVersion))
		} else {
			fmt.Println(ui.Title(r.Name), ui.Info("up-to-date"))
		}
	}

	return nil
}
