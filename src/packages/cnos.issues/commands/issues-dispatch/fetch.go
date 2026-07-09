package issuesdispatch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// githubAPIBase is the GitHub REST API root, kept as a package-level var
// (not a literal baked into every call site) so tests can point it at an
// httptest server and assert on the exact requests issued — mirrors
// issues-fsm/fetch.go's githubAPIBase idiom exactly (cnos#640 γ scaffold
// §2's "reuse the HTTP primitives shape" instruction). The live default is
// always "https://api.github.com".
var githubAPIBase = "https://api.github.com"

// ghGetJSON issues an authenticated GET against the GitHub REST API and
// decodes the JSON body into out. Verbatim shape reuse of
// issues-fsm/fetch.go's ghGetJSON — issues-dispatch is its own Go module
// (own go.mod, per the implementation contract), so the function itself
// cannot be imported across the package boundary; this reproduces the
// same dependency-free net/http idiom rather than forking a different one.
func ghGetJSON(ctx context.Context, url, token string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "cn-issues-dispatch")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api %s: HTTP %d", url, resp.StatusCode)
	}
	return json.Unmarshal(body, out)
}

// ghRequest issues an authenticated non-GET request against the GitHub
// REST API and returns the raw response for the caller to interpret.
// Verbatim shape reuse of issues-fsm/fetch.go's ghRequest (see ghGetJSON's
// doc comment for why this is reproduced rather than imported).
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
	req.Header.Set("User-Agent", "cn-issues-dispatch")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return http.DefaultClient.Do(req)
}

// ghAddLabel adds label to issue via POST /repos/{repo}/issues/{issue}/labels.
// Verbatim shape reuse of issues-fsm/fetch.go's ghAddLabel.
func ghAddLabel(ctx context.Context, repo string, issue int, token, label string) error {
	payload, err := json.Marshal(struct {
		Labels []string `json:"labels"`
	}{Labels: []string{label}})
	if err != nil {
		return err
	}
	addURL := fmt.Sprintf("%s/repos/%s/issues/%d/labels", githubAPIBase, repo, issue)
	resp, err := ghRequest(ctx, http.MethodPost, addURL, token, payload)
	if err != nil {
		return fmt.Errorf("github api add label %q: %w", label, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github api add label %q: HTTP %d: %s", label, resp.StatusCode, string(b))
	}
	return nil
}

// ghRemoveLabel removes label from issue via DELETE
// /repos/{repo}/issues/{issue}/labels/{label}. A 404 (label already absent)
// is tolerated, not an error. Verbatim shape reuse of issues-fsm/fetch.go's
// ghRemoveLabel.
func ghRemoveLabel(ctx context.Context, repo string, issue int, token, label string) error {
	removeURL := fmt.Sprintf("%s/repos/%s/issues/%d/labels/%s", githubAPIBase, repo, issue, label)
	resp, err := ghRequest(ctx, http.MethodDelete, removeURL, token, nil)
	if err != nil {
		return fmt.Errorf("github api remove label %q: %w", label, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github api remove label %q: HTTP %d: %s", label, resp.StatusCode, string(b))
	}
	return nil
}

// ghEditIssueBody is the one genuinely new HTTP primitive this cell adds
// (γ scaffold §2 / §3 implementation contract row "runtime dependencies" —
// no existing precedent in this repo). It rewrites an issue's body via
// PATCH /repos/{repo}/issues/{issue}, modeled directly on ghAddLabel's
// shape above: marshal a small payload, call ghRequest, check the status
// code, surface the response body on failure.
func ghEditIssueBody(ctx context.Context, repo string, issue int, token, body string) error {
	payload, err := json.Marshal(struct {
		Body string `json:"body"`
	}{Body: body})
	if err != nil {
		return err
	}
	editURL := fmt.Sprintf("%s/repos/%s/issues/%d", githubAPIBase, repo, issue)
	resp, err := ghRequest(ctx, http.MethodPatch, editURL, token, payload)
	if err != nil {
		return fmt.Errorf("github api edit issue body: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github api edit issue body: HTTP %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// issueState is what Dispatch needs to observe about the target issue: its
// current body and its full label set (from which the current status:*
// label is derived — see currentStatusLabel in dispatch.go).
type issueState struct {
	Body   string
	Labels []string
}

// liveGetIssue fetches the target issue's body + labels via GET
// /repos/{repo}/issues/{issue} — the read half of the primitive, mirroring
// assembleLive's single-GET label read in issues-fsm/fetch.go but also
// pulling body (fsm's read never needed body).
func liveGetIssue(ctx context.Context, repo string, issue int, token string) (issueState, error) {
	var raw struct {
		Body   string `json:"body"`
		Labels []struct {
			Name string `json:"name"`
		} `json:"labels"`
	}
	url := fmt.Sprintf("%s/repos/%s/issues/%d", githubAPIBase, repo, issue)
	if err := ghGetJSON(ctx, url, token, &raw); err != nil {
		return issueState{}, err
	}
	st := issueState{Body: raw.Body}
	for _, l := range raw.Labels {
		st.Labels = append(st.Labels, l.Name)
	}
	return st, nil
}
