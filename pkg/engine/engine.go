package engine

import (
	"context"
	"os"
	"sync"

	"github.com/cilginc/gochecker/pkg"
	"github.com/cilginc/gochecker/pkg/providers"
)

// Check retrieves the current status of packages without saving changes.
func Check(ctx context.Context, configPath ...string) ([]pkg.Result, error) {
	_, packages, err := loadConfigAndPackages(configPath...)
	if err != nil {
		return nil, err
	}
	return Run(ctx, packages), nil
}

// Execute checks for updates and persists new versions to the config file.
func Execute(ctx context.Context, configPath ...string) ([]pkg.Result, error) {
	cfg, packages, err := loadConfigAndPackages(configPath...)
	if err != nil {
		return nil, err
	}

	results := Run(ctx, packages)

	// Update the config object with new versions from successful results
	updated := false
	resultMap := make(map[string]pkg.Result)
	for _, r := range results {
		resultMap[r.Name] = r
	}

	for i := range cfg.Packages {
		if r, ok := resultMap[cfg.Packages[i].Name]; ok && r.Error == nil && r.Updated {
			cfg.Packages[i].Version = r.NewVersion
			updated = true
		}
	}

	if updated {
		if err := cfg.SaveVersions(); err != nil {
			return results, err
		}
	}

	return results, nil
}

// Run performs concurrent version checks for the provided packages.
func Run(ctx context.Context, packages []pkg.Package) []pkg.Result {
	var wg sync.WaitGroup
	resultCh := make(chan pkg.Result, len(packages))
	token := os.Getenv(pkg.GITHUB_PAT_TOKEN_ENV_VAR)

	for _, p := range packages {
		wg.Add(1)
		go func(p pkg.Package) {
			defer wg.Done()
			resultCh <- checkPackage(ctx, p, token)
		}(p)
	}

	wg.Wait()
	close(resultCh)

	results := make([]pkg.Result, 0, len(packages))
	for r := range resultCh {
		results = append(results, r)
	}
	return results
}

func loadConfigAndPackages(configPath ...string) (*pkg.Config, []pkg.Package, error) {
	cfg, err := pkg.CheckConfig(configPath...)
	if err != nil {
		return nil, nil, err
	}
	if err := cfg.LoadVersions(); err != nil {
		return nil, nil, err
	}
	return cfg, cfg.Packages, nil
}

func checkPackage(ctx context.Context, p pkg.Package, githubToken string) pkg.Result {
	res := pkg.Result{
		Name:       p.Name,
		OldVersion: p.Version,
	}

	if p.Provider.GitHub != nil {
		provider := &providers.GitHub{
			GitHub: p.Provider.GitHub,
			Token:  githubToken,
		}

		version, err := provider.LatestVersion(ctx)
		if err != nil {
			res.Error = err
			return res
		}

		res.NewVersion = version
		res.Updated = version > p.Version
	}

	return res
}
