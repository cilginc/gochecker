package providers

import (
	"context"
	"fmt"
	"sort"

	"github.com/cilginc/gochecker/pkg"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type OCI struct {
	*pkg.OCI
}

type ociTagResponse struct {
	Tags []string `json:"tags"`
}

// [TODO]: We probably don't need a library for this.
func (o *OCI) LatestVersion(ctx context.Context) (string, error) {
	repo, err := name.NewRepository(o.Image)
	if err != nil {
		return "", fmt.Errorf("%w: %v", pkg.ErrInvalidConfig, err)
	}

	tags, err := remote.List(repo, remote.WithContext(ctx))
	if err != nil {
		return "", fmt.Errorf("%w: %v", pkg.ErrOCIRequest, err)
	}

	if len(tags) == 0 {
		return "", pkg.ErrOCINotFound
	}

	sort.Strings(tags)
	return tags[len(tags)-1], nil
}
