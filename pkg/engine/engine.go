package engine

import (
	"context"
	"os"
	"strings"
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

	var rawVersion string
	var err error

	switch {
	case p.Provider.GitHub != nil:
		provider := &providers.GitHub{
			GitHub: p.Provider.GitHub,
			Token:  githubToken,
		}
		rawVersion, err = provider.LatestVersion(ctx)

	case p.Provider.PyPI != nil:
		provider := &providers.PyPI{
			PyPI: p.Provider.PyPI,
		}
		rawVersion, err = provider.LatestVersion(ctx)

	case p.Provider.OCI != nil:
		provider := &providers.OCI{
			OCI: p.Provider.OCI,
		}
		rawVersion, err = provider.LatestVersion(ctx)

	case p.Provider.AUR != nil:
		provider := &providers.AUR{
			AUR: p.Provider.AUR,
		}
		rawVersion, err = provider.LatestVersion(ctx)

	default:
		res.Error = pkg.ErrUnknownProvider
		return res
	}

	if err != nil {
		res.Error = err
		return res
	}

	finalVersion := processVersion(rawVersion, p)
	res.NewVersion = finalVersion
	res.Updated = (finalVersion != p.Version)

	return res
}

func processVersion(v string, p pkg.Package) string {
	if p.Prefix != "" {
		v = strings.TrimPrefix(v, p.Prefix)
	}

	return v
}
