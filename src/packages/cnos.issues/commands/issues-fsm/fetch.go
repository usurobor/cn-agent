package issuesfsm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// githubAPIBase is the GitHub REST API root. It is a package-level var
// (rather than a literal baked into every call site) solely so the
// cnos#569 Phase 2 label-write tests can point it at an httptest server
// and assert on the exact requests issued — the live default is always
// "https://api.github.com". assembleLive's read-only GETs intentionally
// keep their existing literal URLs (byte-identical to Phase 1); only the
// new mutation path below (applyStatusLabel and its helpers) reads this
// var, since that path is the one cnos#569 adds and the one the tests
// need to intercept.
var githubAPIBase = "https://api.github.com"

// assembleLive builds a FactSnapshot for issue N by combining the GitHub
// REST API (workflow-run / PR / label state, via net/http + a GITHUB_TOKEN)
// with local git/filesystem observation (branch / diff / .cdd-artifact
// facts, via os/exec and os), per the issue's fact model. This is the "live"
// path; the "fixture" path (LoadFixture) is what AC3-AC7 actually gate, so
// this path is kept honest but deliberately simple — see the γ scaffold's
// implementation contract.
//
// assembleLive must be run from within a cnos checkout (it shells out to
// git in the current working directory and reads .cdd/unreleased/{N}/
// relative to it).
func assembleLive(ctx context.Context, repo string, issue int, token string) (FactSnapshot, error) {
	snap := FactSnapshot{Issue: issue}

	// --- GitHub REST: labels ---
	if repo != "" {
		var ghIssue struct {
			Labels []struct {
				Name string `json:"name"`
			} `json:"labels"`
		}
		if err := ghGetJSON(ctx, fmt.Sprintf("https://api.github.com/repos/%s/issues/%d", repo, issue), token, &ghIssue); err == nil {
			for _, l := range ghIssue.Labels {
				snap.Labels = append(snap.Labels, l.Name)
			}
		}
		// else: leave Labels empty rather than failing the whole snapshot —
		// the evaluator will report state "" (unlabeled) and the caller can
		// see the gap in the printed facts.
	}

	branch := "cycle/" + strconv.Itoa(issue)

	// --- local git: branch existence + commits beyond base ---
	if err := exec.CommandContext(ctx, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+branch).Run(); err == nil {
		snap.BranchExists = true
		out, err := exec.CommandContext(ctx, "git", "rev-list", "--count", "main.."+branch).Output()
		if err == nil {
			if n, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
				snap.CommitsBeyondBase = n
			}
		}
	}

	// --- local filesystem: .cdd/unreleased/{N}/ artifacts ---
	dir := filepath.Join(".cdd", "unreleased", strconv.Itoa(issue))
	if entries, err := os.ReadDir(dir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			snap.CDDArtifacts = append(snap.CDDArtifacts, e.Name())
			switch e.Name() {
			case "REVIEW-REQUEST.yml":
				snap.ReviewRequestPresent = true
			case "delta-repair.md":
				snap.RepairContractPresent = true
			}
		}
		sort.Strings(snap.CDDArtifacts)
	}
	// A beta-review.md with an explicit "iterate" verdict is also repair
	// context; a lightweight substring check is enough for the live path
	// (the fixture path is what the ACs gate — see package doc).
	if !snap.RepairContractPresent {
		if data, err := os.ReadFile(filepath.Join(dir, "beta-review.md")); err == nil {
			if strings.Contains(string(data), "verdict: iterate") {
				snap.RepairContractPresent = true
			}
		}
	}

	// A gamma-scaffold.md's `cell_kind:` line is the canonical recording
	// point cnos#570 chose for the cell-kind observation seam
	// (CELL-KINDS.md §"Recording point"). This is observation only: it
	// populates the seam that snapshot.go already carries
	// (FactSnapshot.CellKind), but no transition rule in table.go consumes
	// it — TestSeam_CellKindNotEnforced locks that invariant.
	if data, err := os.ReadFile(filepath.Join(dir, "gamma-scaffold.md")); err == nil {
		if kind := parseCellKind(string(data)); kind != "" {
			snap.CellKind.Observed = kind
			snap.CellKind.Source = "cdd_artifact"
		}
	}

	if repo == "" {
		snap.normalizeCellKind()
		return snap, nil
	}

	// --- GitHub REST: PR existence + commit count ---
	owner := repo
	if i := strings.IndexByte(repo, '/'); i > 0 {
		owner = repo[:i]
	}
	var prs []struct {
		Number int `json:"number"`
	}
	prURL := fmt.Sprintf("https://api.github.com/repos/%s/pulls?head=%s:%s&state=all&per_page=1", repo, owner, branch)
	if err := ghGetJSON(ctx, prURL, token, &prs); err == nil && len(prs) > 0 {
		snap.PRExists = true
		var commits []json.RawMessage
		commitsURL := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/commits?per_page=100", repo, prs[0].Number)
		if err := ghGetJSON(ctx, commitsURL, token, &commits); err == nil {
			snap.PRCommitCount = len(commits)
		}
	}

	// --- GitHub REST: most recent workflow run on this branch ---
	var runs struct {
		WorkflowRuns []struct {
			Status     string `json:"status"`
			Conclusion string `json:"conclusion"`
		} `json:"workflow_runs"`
	}
	runsURL := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs?branch=%s&per_page=1", repo, branch)
	if err := ghGetJSON(ctx, runsURL, token, &runs); err == nil && len(runs.WorkflowRuns) > 0 {
		r := runs.WorkflowRuns[0]
		snap.RunState = r.Status // "queued", "in_progress", "completed"
		if r.Status == "completed" {
			snap.RunConclusion = r.Conclusion
			snap.ChecksState = r.Conclusion
		} else {
			snap.ChecksState = "pending"
		}
	}

	snap.normalizeCellKind()
	return snap, nil
}

// cellKindLinePattern matches a `cell_kind:` recording line in either the
// bold-markdown form γ scaffolds use (`**cell_kind:** `implementation` ...`,
// with optional trailing prose after the value) or the plain form
// (`cell_kind: implementation`). Kept intentionally simple — a line-scan /
// single regex, per CELL-KINDS.md §"Recording point" (cnos#570) — this is an
// observation convenience, not a schema-validated parse.
var cellKindLinePattern = regexp.MustCompile("(?m)^\\*{0,2}cell_kind:\\*{0,2}\\s*`?([A-Za-z_][A-Za-z0-9_-]*)")

// parseCellKind scans gamma-scaffold.md content for the first `cell_kind:`
// line and returns the kind token, or "" if none is present. See
// CELL-KINDS.md §"Recording point" for the convention this parses.
func parseCellKind(scaffold string) string {
	m := cellKindLinePattern.FindStringSubmatch(scaffold)
	if m == nil {
		return ""
	}
	return m[1]
}

// ghGetJSON issues an authenticated GET against the GitHub REST API and
// decodes the JSON body into out. Mirrors commands/issues-map/fetch.go's
// dependency-free stdlib approach (net/http, no third-party client).
func ghGetJSON(ctx context.Context, url, token string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "cn-issues-fsm")
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

// --- cnos#569 Phase 2: the label-write primitives ---------------------
//
// Everything above this point (and everything in the rest of the
// package) is read-only, exactly as Phase 1 (cnos#568) shipped it. The
// functions below are the ONLY code in this package that mutate GitHub
// state, and they are reachable only from issuesfsm.go's --apply branch,
// itself gated on Evaluate having already produced outcome=="proposed"
// with a non-empty TargetState (i.e. the transition table's guards
// already passed — see table.go's Rule matching). No third-party GitHub
// client: same dependency-free net/http + auth-header idiom as
// ghGetJSON above.

// ghRequest issues an authenticated non-GET request against the GitHub
// REST API and returns the raw response for the caller to interpret
// (status-code handling differs between add-label and remove-label, in
// particular the 404-tolerance remove-label wants — see
// applyStatusLabel's doc comment).
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
	req.Header.Set("User-Agent", "cn-issues-fsm")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return http.DefaultClient.Do(req)
}

// ghAddLabel adds label to issue (repo "owner/name") via POST
// /repos/{repo}/issues/{issue}/labels. GitHub's add-labels endpoint is
// itself idempotent (re-adding a label the issue already carries is not
// an error), so no pre-check is needed here.
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
// /repos/{repo}/issues/{issue}/labels/{label}. A 404 (label already
// absent) is tolerated, not an error — see applyStatusLabel's doc
// comment for why.
func ghRemoveLabel(ctx context.Context, repo string, issue int, token, label string) error {
	removeURL := fmt.Sprintf("%s/repos/%s/issues/%d/labels/%s", githubAPIBase, repo, issue, url.PathEscape(label))
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

// applyStatusLabel is the single mutation entry point this package
// exposes: it moves issue #{issue}'s status:* label from fromState to
// toState (label names are the CDS "status:" prefix convention
// CurrentState already reads — see snapshot.go). fromState is the
// Decision's CurrentState (may be "" if the issue was unlabeled).
//
// 404-tolerance decision (documented per the γ scaffold's Friction
// notes, since it is an implementation-contract-adjacent call the
// scaffold left to α): the remove-old-label call tolerates a 404
// (already removed). This keeps a manually-recovered or
// partially-applied state idempotent — e.g. an operator who already
// hand-removed status:in-progress, or a prior --apply call that added
// the new label but crashed before removing the old one, must not turn
// a second --apply into a hard failure. The add-new-label call is not
// given the same tolerance for a *conflicting* error (only 200/201 are
// accepted) because GitHub's add-labels endpoint is already idempotent
// for the "label already present" case (see ghAddLabel), so there is no
// analogous benign-404 case to tolerate there.
func applyStatusLabel(ctx context.Context, repo string, issue int, token, fromState, toState string) error {
	if fromState != "" && fromState != toState {
		if err := ghRemoveLabel(ctx, repo, issue, token, "status:"+fromState); err != nil {
			return fmt.Errorf("remove old label status:%s: %w", fromState, err)
		}
	}
	if err := ghAddLabel(ctx, repo, issue, token, "status:"+toState); err != nil {
		return fmt.Errorf("add label status:%s: %w", toState, err)
	}
	return nil
}
