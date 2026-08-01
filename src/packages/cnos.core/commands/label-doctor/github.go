package labeldoctor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// githubAPIBase is the GitHub REST API root. It is a package-level var
// (rather than a literal baked into every call site) solely so this
// package's tests can point it at an httptest server and assert on the
// exact requests issued — the live default is always
// "https://api.github.com". Mirrors cnos.issues/commands/issues-fsm/
// fetch.go's githubAPIBase idiom exactly.
var githubAPIBase = "https://api.github.com"

// ghLabel is the GitHub REST label wire shape (a strict subset of
// GitHub's label object — id/url/default are irrelevant to the audit).
type ghLabel struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// ghRequest issues an authenticated request against the GitHub REST API
// and returns the raw response for the caller to interpret. Mirrors
// issues-fsm/fetch.go's ghRequest exactly (dependency-free net/http, no
// third-party GitHub client).
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
	req.Header.Set("User-Agent", "cn-label-doctor")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return http.DefaultClient.Do(req)
}

// ghListLabels returns every repo-level label via GET
// /repos/{repo}/labels, paginated (100 per page) until a short page is
// returned. This is the "list-repo-labels GET" primitive the γ scaffold
// names as missing from the issues-fsm precedent (that package only has
// ghEnsureLabelExists, a create-only POST).
func ghListLabels(ctx context.Context, repo, token string) ([]ghLabel, error) {
	var all []ghLabel
	for page := 1; ; page++ {
		listURL := fmt.Sprintf("%s/repos/%s/labels?per_page=100&page=%d", githubAPIBase, repo, page)
		resp, err := ghRequest(ctx, http.MethodGet, listURL, token, nil)
		if err != nil {
			return nil, fmt.Errorf("github api list labels: %w", err)
		}
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("github api list labels: read body: %w", readErr)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("github api list labels: HTTP %d: %s", resp.StatusCode, string(body))
		}
		var batch []ghLabel
		if err := json.Unmarshal(body, &batch); err != nil {
			return nil, fmt.Errorf("github api list labels: decode: %w", err)
		}
		all = append(all, batch...)
		if len(batch) < 100 {
			break
		}
	}
	return all, nil
}

// ghCreateLabel creates a repo-level label via POST /repos/{repo}/labels.
// A 422 whose error body names code "already_exists" is tolerated as a
// no-op — this keeps the primitive idempotent-by-construction even if
// the audit's own live-state read raced with an external creator,
// mirroring issues-fsm/fetch.go's ghEnsureLabelExists (cnos#615) idiom.
// Unlike that precedent (which tolerates ANY 422 unconditionally), this
// function inspects the error body's code first: GitHub returns 422 for
// OTHER validation failures too (e.g. a description over 100
// characters — found live against a real canonical-manifest value,
// src/packages/cnos.core/labels.json's dispatch:cell description is 149
// bytes), and treating every 422 as a harmless "already exists" no-op
// would silently swallow those as false successes.
func ghCreateLabel(ctx context.Context, repo, token string, l ghLabel) error {
	payload, err := json.Marshal(l)
	if err != nil {
		return err
	}
	createURL := fmt.Sprintf("%s/repos/%s/labels", githubAPIBase, repo)
	resp, err := ghRequest(ctx, http.MethodPost, createURL, token, payload)
	if err != nil {
		return fmt.Errorf("github api create label %q: %w", l.Name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnprocessableEntity {
		b, _ := io.ReadAll(resp.Body)
		if isAlreadyExistsError(b) {
			return nil
		}
		return fmt.Errorf("github api create label %q: HTTP 422: %s", l.Name, string(b))
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github api create label %q: HTTP %d: %s", l.Name, resp.StatusCode, string(b))
	}
	return nil
}

// isAlreadyExistsError reports whether a GitHub 422 response body's
// errors[] array names code "already_exists" (the create-label
// endpoint's shape for "a label with this name already exists" —
// distinct from other 422 causes such as an over-length description or
// an invalid color).
func isAlreadyExistsError(body []byte) bool {
	var parsed struct {
		Errors []struct {
			Code string `json:"code"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return false
	}
	for _, e := range parsed.Errors {
		if e.Code == "already_exists" {
			return true
		}
	}
	return false
}

// ghUpdateLabel repairs an existing label's color/description via PATCH
// /repos/{repo}/labels/{name}. This is the "update-drifted-label PATCH"
// primitive the γ scaffold names as genuinely new — no existing
// precedent in this repo covers repair-drift-on-existing (only create,
// via ghEnsureLabelExists). The label's name is never renamed by this
// call (new_name is intentionally omitted): drift repair here is
// color/description only, matching AC2's scope.
func ghUpdateLabel(ctx context.Context, repo, token, name, color, description string) error {
	payload, err := json.Marshal(struct {
		Color       string `json:"color"`
		Description string `json:"description"`
	}{Color: color, Description: description})
	if err != nil {
		return err
	}
	updateURL := fmt.Sprintf("%s/repos/%s/labels/%s", githubAPIBase, repo, url.PathEscape(name))
	resp, err := ghRequest(ctx, http.MethodPatch, updateURL, token, payload)
	if err != nil {
		return fmt.Errorf("github api update label %q: %w", name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github api update label %q: HTTP %d: %s", name, resp.StatusCode, string(b))
	}
	return nil
}
