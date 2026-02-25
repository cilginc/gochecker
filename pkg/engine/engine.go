package engine

import (
	"context"
	"os"
	"sync"

	"github.com/cilginc/gochecker/pkg"
	"github.com/cilginc/gochecker/pkg/providers"
)

// [TODO]: DRY clean the code here.
func Execute(ctx context.Context, configPath ...string) ([]pkg.Result, error) {
	cfg, err := pkg.CheckConfig(configPath...)
	if err != nil {
		return nil, err
	}

	if err := cfg.LoadVersions(); err != nil {
		return nil, err
	}

	results := Run(ctx, cfg.Packages)

	resultMap := make(map[string]pkg.Result)
	for _, r := range results {
		resultMap[r.Name] = r
	}

	for i := range cfg.Packages {
		if r, ok := resultMap[cfg.Packages[i].Name]; ok {
			if r.Error == nil && r.Updated {
				cfg.Packages[i].Version = r.NewVersion
			}
		}
	}

	if err := cfg.SaveVersions(); err != nil {
		return results, err
	}

	return results, nil
}

func Run(ctx context.Context, packages []pkg.Package) []pkg.Result {
	var wg sync.WaitGroup
	resultCh := make(chan pkg.Result, len(packages))

	for _, p := range packages {
		wg.Add(1)

		go func(p pkg.Package) {
			defer wg.Done()

			res := pkg.Result{
				Name:       p.Name,
				OldVersion: p.Version,
			}

			if p.Provider.GitHub != nil {
				token := os.Getenv(pkg.GITHUB_PAT_TOKEN_ENV_VAR)
				provider := &providers.GitHub{
					GitHub: p.Provider.GitHub,
					Token:  token,
				}

				version, err := provider.LatestVersion(ctx)
				if err != nil {
					res.Error = err
					resultCh <- res
					return
				}

				res.NewVersion = version
				res.Updated = version > p.Version
			}

			resultCh <- res
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

func Check(ctx context.Context, configPath ...string) ([]pkg.Result, error) {
	cfg, err := pkg.CheckConfig(configPath...)
	if err != nil {
		return nil, err
	}

	if err := cfg.LoadVersions(); err != nil {
		return nil, err
	}

	return Run(ctx, cfg.Packages), nil
}
