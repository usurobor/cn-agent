package issuesmap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// fetchIssues pulls all open issues for repo (owner/name) from the GitHub
// REST API, following pagination. REST is used over GraphQL for a
// dependency-free stdlib implementation; open-issue volume is small enough
// that per-page fetches are cheap. The Actions GITHUB_TOKEN is sufficient
// (public read); an empty token still works for public repos, subject to
// unauthenticated rate limits.
func fetchIssues(ctx context.Context, repo, token string) ([]ghIssue, error) {
	if repo == "" {
		return nil, fmt.Errorf("cn issues map: --repo (or $GITHUB_REPOSITORY) is required for the live path; use --fixture/--stdin offline")
	}
	client := &http.Client{}
	var all []ghIssue
	for page := 1; page <= 100; page++ {
		u := fmt.Sprintf("https://api.github.com/repos/%s/issues?state=open&per_page=100&page=%d", repo, page)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		req.Header.Set("User-Agent", "cn-issues-map")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("github api request: %w", err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("github api %s: HTTP %d: %s", u, resp.StatusCode, truncate(body, 300))
		}
		var batch []ghIssue
		if err := json.Unmarshal(body, &batch); err != nil {
			return nil, fmt.Errorf("decode issues page %d: %w", page, err)
		}
		if len(batch) == 0 {
			break
		}
		all = append(all, batch...)
		if len(batch) < 100 {
			break
		}
	}
	return all, nil
}

func truncate(b []byte, n int) string {
	s := string(b)
	if len(s) > n {
		return s[:n] + "…"
	}
	return s
}
