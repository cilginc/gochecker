package cmd

import (
	"os"

	"github.com/cilginc/gochecker/internal/ui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a .gochecker.yaml file with some examples.",
	Long: ui.Info(`Bootstrap your version monitoring environment. 
This command creates a default configuration file and sets up the 
necessary structure to start tracking upstream software updates.`),
	Example:       ui.Success("gochecker init\n"),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          initFile,
}

func init() {
	rootCmd.AddCommand(initCmd)

}

func initFile(cmd *cobra.Command, args []string) error {
	content := `packages:
  - name: gochecker
    github:
      repo: cilginc/gochecker
    prefix: "v"
`

	targetFile := cfgFile

	if _, err := os.Stat(targetFile); err == nil {
		return ui.CliError("file '%s' already exists, initialization aborted", targetFile)
	}

	err := os.WriteFile(targetFile, []byte(content), 0644)
	if err != nil {
		return ui.CliError("could not write to %s: %v", targetFile, err)
	}

	ui.CliSuccess("Default configuration created successfully at %s", targetFile)
	return nil
}
