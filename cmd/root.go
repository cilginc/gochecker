package cmd

import (
	"os"
	"regexp"

	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg"
	"github.com/spf13/cobra"
)

var (
	cfgFile      string = pkg.DEFAULT_CONFIG_FILE
	versionFile  string = pkg.DEFAULT_VERSIONS_FILE
	outputFormat string = "text"
	noColor      bool   = false
	recursive    bool   = false
	recursiveDir string
)

var rootCmd = &cobra.Command{
	Use:   "gochecker",
	Short: "A fast and modern version checker for software packages",
	Long: ui.Info(`Gochecker is a high-performance version monitoring tool inspired by nvchecker.
It allows you to track upstream versions from various sources like GitHub, AUR, 
and custom web pages using Go's powerful concurrency model.`),
}

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

	f := rootCmd.PersistentFlags()

	// --no-color
	f.BoolVar(
		&noColor,
		"no-color",
		false,
		"Disable colorized output for better compatibility with non-interactive shells",
	)

	// --config / -c
	f.StringVarP(
		&cfgFile,
		"config",
		"c",
		pkg.DEFAULT_CONFIG_FILE,
		"Path to the configuration file that defines software sources and tracking rules",
	)

	// --version-file
	f.StringVar(
		&versionFile,
		"version-file",
		pkg.DEFAULT_VERSIONS_FILE,
		"Path to the file where currently detected versions are stored",
	)

	// --output / -o
	f.StringVarP(
		&outputFormat,
		"output",
		"o",
		"text",
		"Set the output format to display results; supported: 'text', 'json', 'yaml'",
	)

	// --recursive / -r
	f.BoolVarP(
		&recursive,
		"recursive",
		"r",
		false,
		"Recursively scan a directory for configuration files",
	)

	// --dir / -d
	f.StringVarP(
		&recursiveDir,
		"dir",
		"d",
		".",
		"Directory to scan when using --recursive (defaults to current directory)",
	)

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
