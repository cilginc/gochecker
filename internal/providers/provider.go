package providers

import (
	"context"
)

type Provider interface {
	LatestVersion(ctx context.Context) (string, error)
}
