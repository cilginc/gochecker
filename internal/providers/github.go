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

func (g *GitHub) LatestVersion(ctx context.Context) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", g.Repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	if g.Token != "" {
		req.Header.Set("Authorization", "token "+g.Token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		if resp.Header.Get("X-Ratelimit-Remaining") == "0" {
			unixTimeStr := resp.Header.Get("X-Ratelimit-Reset")

			sec, _ := strconv.ParseInt(unixTimeStr, 10, 64)
			resetTime := time.Unix(sec, 0).Format("15:04:05 MST")

			return "", fmt.Errorf("github api limit exceeded, retry at %s", resetTime)
		}
		return "", fmt.Errorf("access forbidden")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("api returned %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}
