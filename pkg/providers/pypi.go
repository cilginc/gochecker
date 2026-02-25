package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cilginc/gochecker/pkg"
)

type PyPI struct {
	*pkg.PyPI
}

type pypiResponse struct {
	Info struct {
		Version string `json:"version"`
	} `json:"info"`
}

// LatestVersion fetches the latest version from the PyPI API.
func (p *PyPI) LatestVersion(ctx context.Context) (string, error) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", p.Package)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Good practice to set a User-Agent for API requests
	req.Header.Set("User-Agent", "gochecker/1.0 (https://github.com/cilginc/gochecker)")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", pkg.ErrPyPIRequest, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return "", pkg.ErrPyPINotFound
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %d", pkg.ErrPyPIStatus, resp.StatusCode)
	}

	var data pypiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", pkg.ErrPyPIDecode
	}

	if data.Info.Version == "" {
		return "", pkg.ErrPyPIEmptyVersion
	}

	return data.Info.Version, nil
}

