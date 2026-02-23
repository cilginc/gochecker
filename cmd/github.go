package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/cilginc/gochecker/internal/providers"
	"github.com/cilginc/gochecker/internal/ui"
	"github.com/cilginc/gochecker/pkg"
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
	cfg, err := loadConfig(".gochecker.yaml")
	if err != nil {
		return err
	}

	if err := validateConfig(cfg); err != nil {
		return err
	}

	ctx := context.Background()
	var wg sync.WaitGroup

	for _, p := range cfg.Packages {
		wg.Add(1)

		go func(pkg pkg.Package) {
			defer wg.Done()

			if pkg.Provider.GitHub != nil {
				provider := &providers.GitHub{
					GitHub: pkg.Provider.GitHub,
				}

				version, err := provider.LatestVersion(ctx)
				if err != nil {
					fmt.Printf("%s: %v\n", ui.Err("Error"), err)
					return
				}

				fmt.Println(ui.Title(pkg.Name), ui.Info(version))
			}
		}(p)
	}

	wg.Wait()
	return nil
}
