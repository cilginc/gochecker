package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/cilginc/gochecker/pkg"
)

type GitHub struct {
	*pkg.GitHub
	Token string
}

func (g *GitHub) LatestVersion(ctx context.Context) (string, error) {
	baseURL := "https://api.github.com"
	if g.Host != "" {
		baseURL = fmt.Sprintf("https://%s/api/v3", g.Host)
	}

	var endpoint string
	switch {
	case g.UseLatestRelease:
		endpoint = fmt.Sprintf("%s/repos/%s/releases/latest", baseURL, g.Repo)
		return g.fetchSingleTag(ctx, endpoint, "tag_name")

	case g.UseMaxRelease || g.IncludePrerelease:
		endpoint = fmt.Sprintf("%s/repos/%s/releases", baseURL, g.Repo)
		return g.fetchAndSort(ctx, endpoint, "tag_name", g.IncludePrerelease)

	case g.UseLatestTag || g.UseMaxTag:
		endpoint = fmt.Sprintf("%s/repos/%s/tags", baseURL, g.Repo)
		return g.fetchAndSort(ctx, endpoint, "name", false)

	default:
		endpoint = fmt.Sprintf("%s/repos/%s/releases/latest", baseURL, g.Repo)
		return g.fetchSingleTag(ctx, endpoint, "tag_name")
	}
}

func (g *GitHub) fetchSingleTag(ctx context.Context, url, key string) (string, error) {
	resp, err := g.makeRequest(ctx, url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("%w: %v", pkg.ErrGitHubDecode, err)
	}

	if val, ok := data[key].(string); ok {
		return val, nil
	}
	return "", fmt.Errorf("%w: version key not found", pkg.ErrGitHubDecode)
}

func (g *GitHub) fetchAndSort(
	ctx context.Context,
	url, key string,
	includePre bool,
) (string, error) {
	resp, err := g.makeRequest(ctx, url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	var items []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return "", fmt.Errorf("%w: %v", pkg.ErrGitHubDecode, err)
	}

	if len(items) == 0 {
		return "", fmt.Errorf("no versions found")
	}

	var versions []string
	for _, item := range items {
		if key == "tag_name" && !includePre {
			if isPre, ok := item["prerelease"].(bool); ok && isPre {
				continue
			}
		}

		if v, ok := item[key].(string); ok {
			versions = append(versions, v)
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no valid versions found after filtering")
	}

	// [TODO]: Maybe use Masterminds/semver later.
	sort.Slice(versions, func(i, j int) bool {
		return versions[i] > versions[j]
	})

	return versions[0], nil
}

func (g *GitHub) makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", pkg.ErrGitHubRequest, err)
	}

	if g.Token != "" {
		req.Header.Set("Authorization", "token "+g.Token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", pkg.ErrGitHubRequest, err)
	}

	if resp.StatusCode == http.StatusForbidden {
		if resp.Header.Get("X-Ratelimit-Remaining") == "0" {
			unixTimeStr := resp.Header.Get("X-Ratelimit-Reset")
			sec, _ := strconv.ParseInt(unixTimeStr, 10, 64)
			resetTime := time.Unix(sec, 0).Format(time.RFC1123)
			return nil, fmt.Errorf("%w: retry at %s", pkg.ErrGitHubRateLimit, resetTime)
		}
		return nil, pkg.ErrGitHubForbidden
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", pkg.ErrGitHubStatus, resp.StatusCode)
	}

	return resp, nil
}
