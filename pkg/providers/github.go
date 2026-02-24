package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cilginc/gochecker/pkg"
)

type GitHub struct {
	*pkg.GitHub
	Token string
}

// Reads the GITHUB_TOKEN environment variable for authentication.
// Any manually provided token will override the value fetched from the environment.
func (g *GitHub) LatestVersion(ctx context.Context) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", g.Repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("%w: %v", pkg.ErrGitHubRequest, err)
	}

	if g.Token != "" {
		req.Header.Set("Authorization", "token "+g.Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", pkg.ErrGitHubRequest, err)
	}
	defer resp.Body.Close()

	// Rate limit handling
	if resp.StatusCode == http.StatusForbidden {
		if resp.Header.Get("X-Ratelimit-Remaining") == "0" {
			unixTimeStr := resp.Header.Get("X-Ratelimit-Reset")

			sec, _ := strconv.ParseInt(unixTimeStr, 10, 64)
			resetTime := time.Unix(sec, 0).Format(time.RFC1123)

			return "", fmt.Errorf("%w: retry at %s",
				pkg.ErrGitHubRateLimit,
				resetTime,
			)
		}
		return "", pkg.ErrGitHubForbidden
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %d",
			pkg.ErrGitHubStatus,
			resp.StatusCode,
		)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("%w: %v",
			pkg.ErrGitHubDecode,
			err,
		)
	}

	return release.TagName, nil
}
