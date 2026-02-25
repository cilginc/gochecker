package providers

import (
	"context"
	"fmt"
	"sort"

	"github.com/cilginc/gochecker/pkg"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/storage/memory"
)

type Git struct {
	*pkg.Git
}

func (g *Git) LatestVersion(ctx context.Context) (string, error) {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{g.URL},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("%w: %v", pkg.ErrGitRequest, err)
	}

	if g.UseCommit {
		target := g.Branch
		if target == "" {
			target = "master"
		}

		for _, ref := range refs {
			if ref.Name().Short() == target {
				return ref.Hash().String(), nil
			}
		}
		return "", pkg.ErrGitBranchNotFound
	}

	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, ref.Name().Short())
		}
	}

	if len(tags) == 0 {
		return "", pkg.ErrGitNoTags
	}

	sort.Strings(tags)
	return tags[len(tags)-1], nil
}
