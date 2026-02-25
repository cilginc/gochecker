package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cilginc/gochecker/pkg"
)

type AUR struct {
	*pkg.AUR
}

type aurResponse struct {
	ResultCount int `json:"resultcount"`
	Results     []struct {
		Version string `json:"Version"`
	} `json:"results"`
}

func (a *AUR) LatestVersion(ctx context.Context) (string, error) {
	url := fmt.Sprintf("https://aur.archlinux.org/rpc/v5/info/%s", a.Package)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", pkg.ErrAURRequest, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status %d", pkg.ErrAURRequest, resp.StatusCode)
	}

	var data aurResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", pkg.ErrAURDecode
	}

	if data.ResultCount == 0 || len(data.Results) == 0 {
		return "", pkg.ErrAURNotFound
	}

	version := data.Results[0].Version

	if a.StripRelease {
		if idx := strings.Index(version, "-"); idx != -1 {
			version = version[:idx]
		}
	}

	return version, nil
}
