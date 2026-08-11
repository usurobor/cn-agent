package installpreflight

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// defaultAPIBase is the GitHub REST API root used when Options.APIBaseURL
// is empty.
const defaultAPIBase = "https://api.github.com"

// ghSecret is the GitHub Actions repository-secret wire shape returned
// by GET /repos/{owner}/{repo}/actions/secrets. This endpoint is
// presence-only BY DESIGN on GitHub's side — it never returns a
// secret's value over the API — so this type intentionally carries NO
// Value/value field. AC3's load-bearing property: any future edit that
// adds one would silently start threading secret material through this
// package; the absence here IS the audit surface (see
// preflight_test.go's TestGhListSecrets_NeverDecodesAValueField).
type ghSecret struct {
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ghSecretsPage is GET .../actions/secrets' top-level response shape.
type ghSecretsPage struct {
	TotalCount int        `json:"total_count"`
	Secrets    []ghSecret `json:"secrets"`
}

// ghRequest issues an authenticated request against the GitHub REST
// API and returns the raw response for the caller to interpret.
// Mirrors label-doctor/github.go's ghRequest (itself mirroring
// issues-fsm/fetch.go's ghRequest) exactly: dependency-free net/http,
// no third-party GitHub client, no `gh` CLI shellout. Not imported
// directly — separate go.work module — but follows the same shape by
// design per the γ scaffold's Implementation contract.
func ghRequest(ctx context.Context, method, apiURL, token string, body []byte) (*http.Response, error) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, apiURL, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "cn-install-preflight")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return http.DefaultClient.Do(req)
}

// ghListSecrets returns every repo-level Actions secret's presence-only
// wire shape via GET /repos/{owner}/{repo}/actions/secrets, paginated
// (100 per page) until a short page is returned. It NEVER decodes a
// value field (see ghSecret's doc comment) — this is AC3's mechanical
// oracle: the response is decoded into a type that has nowhere to put
// a secret value even if GitHub ever sent one.
func ghListSecrets(ctx context.Context, apiBase, repo, token string) ([]ghSecret, error) {
	var all []ghSecret
	for page := 1; ; page++ {
		listURL := fmt.Sprintf("%s/repos/%s/actions/secrets?per_page=100&page=%d", apiBase, repo, page)
		resp, err := ghRequest(ctx, http.MethodGet, listURL, token, nil)
		if err != nil {
			return nil, fmt.Errorf("github api list actions secrets: %w", err)
		}
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("github api list actions secrets: read body: %w", readErr)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("github api list actions secrets: HTTP %d: %s", resp.StatusCode, string(body))
		}
		var batch ghSecretsPage
		if err := json.Unmarshal(body, &batch); err != nil {
			return nil, fmt.Errorf("github api list actions secrets: decode: %w", err)
		}
		all = append(all, batch.Secrets...)
		if len(batch.Secrets) < 100 {
			break
		}
	}
	return all, nil
}

// ghRepo is the presence-only subset of GET /repos/{owner}/{repo} this
// package reads: just the authenticated caller's permissions on the
// repo. Per the γ scaffold's Friction note 3, the single-call
// permissions.push field is used in preference to the two-call
// collaborator/permission endpoint (which would first require
// resolving which account the token belongs to via a separate GET
// /user call) — "does the token used for install have push access to
// this repo" is answered directly and in one round trip.
type ghRepo struct {
	Permissions struct {
		Push bool `json:"push"`
	} `json:"permissions"`
}

// ghCheckPushAccess reports whether token has push access to repo via
// GET /repos/{owner}/{repo}'s permissions.push field for the
// authenticated caller.
func ghCheckPushAccess(ctx context.Context, apiBase, repo, token string) (bool, error) {
	repoURL := fmt.Sprintf("%s/repos/%s", apiBase, repo)
	resp, err := ghRequest(ctx, http.MethodGet, repoURL, token, nil)
	if err != nil {
		return false, fmt.Errorf("github api get repo: %w", err)
	}
	body, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if readErr != nil {
		return false, fmt.Errorf("github api get repo: read body: %w", readErr)
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("github api get repo: HTTP %d: %s", resp.StatusCode, string(body))
	}
	var parsed ghRepo
	if err := json.Unmarshal(body, &parsed); err != nil {
		return false, fmt.Errorf("github api get repo: decode: %w", err)
	}
	return parsed.Permissions.Push, nil
}
