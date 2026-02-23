package cmd

import (
	"os"
	"regexp"

	"github.com/cilginc/gochecker/internal/ui"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gochecker",
	Short: "A brief description of your application",
	Long:  ui.Info("I still don't know"),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	for _, arg := range os.Args {
		if arg == "--no-color" {
			ui.DisableColor()
			break
		}
	}

	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")

	colorFlags := func(s string) string {
		re := regexp.MustCompile(`(?m)(^\s+-[a-zA-Z], --[a-zA-Z0-9-]+|^\s+--[a-zA-Z0-9-]+)`)
		return re.ReplaceAllStringFunc(s, func(m string) string {
			return ui.Info(m)
		})
	}

	cobra.AddTemplateFuncs(map[string]interface{}{
		"yellow":     ui.Warn,
		"cyan":       ui.Info,
		"green":      ui.Success,
		"bold":       ui.Title,
		"colorFlags": colorFlags,
	})

	rootCmd.SetHelpTemplate(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{yellow "Usage:"}}
  {{if .Runnable}}{{bold .UseLine}}{{end}}{{if .HasAvailableSubCommands}}{{cyan .CommandPath}} [command]{{end}}

{{if .HasAvailableSubCommands}}{{yellow "Available Commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{green (rpad .Name .NamePadding)}} {{.Short}}{{end}}{{end}}

{{end}}{{if .HasAvailableLocalFlags}}{{yellow "Flags:"}}
{{.LocalFlags.FlagUsages | colorFlags | trimTrailingWhitespaces}}

{{end}}{{if .HasAvailableInheritedFlags}}{{yellow "Global Flags:"}}
{{.InheritedFlags.FlagUsages | colorFlags | trimTrailingWhitespaces}}

{{end}}{{if .HasHelpSubCommands}}{{yellow "Additional help topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{green (rpad .CommandPath .CommandPathPadding)}} {{.Short}}{{end}}{{end}}

{{end}}{{if .HasAvailableSubCommands}}{{cyan "Use"}} "{{.CommandPath}} [command] --help" {{cyan "for more information about a command."}}{{end}}
`)
}
