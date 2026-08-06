// Package installpreflight implements the operator-prerequisite
// presence checks for `cn repo install --dispatch cds` (cnos#706): it
// verifies the required GitHub Actions repository secrets exist BY
// NAME (never reading a value) and that the installing token has push
// access to the target repo — both via presence-only GitHub REST
// calls, so the CLI never receives or handles a secret's value.
//
// This is genuinely new surface: no prior mechanism in this repo lists
// GitHub Actions secrets or checks repo-permission presence (confirmed
// absent by cnos#706's γ scaffold via
// `grep -rn "actions/secrets\|collaborators.*permission"` across
// src/go and src/packages at scaffold time). It mirrors label-doctor's
// (cnos#493) dependency-free net/http idiom — see
// label-doctor/github.go's ghRequest and label-doctor/resolve.go's
// resolveRepoFromGitRemote — rather than inventing a third pattern.
// Those functions are not imported directly (this package lives in its
// own go.work module, exactly like label-doctor does relative to
// issues-fsm's own ghRequest precedent); the shape is mirrored, not
// shared, per the γ scaffold's Implementation contract.
//
// AC3's load-bearing property: GitHub's
// GET /repos/{owner}/{repo}/actions/secrets endpoint is presence-only
// by design on GitHub's own side — it returns name/created_at/
// updated_at and NEVER a value — so ghSecret (github.go) intentionally
// carries no Value/value field. This package can never read, log, or
// forward a secret's value, by construction, not merely by discipline.
package installpreflight

import (
	"context"
	"fmt"
	"os"
)

// Options configures Check.
type Options struct {
	// RepoRoot anchors git-remote based Repo resolution when Repo is
	// empty. Required in that case.
	RepoRoot string
	// Repo is "owner/repo" — the repo being installed INTO. This is
	// distinct from repoinstall.Options.Repo, which names the cnos
	// release SOURCE, not the installing target (see repoinstall.go's
	// runPreflight doc comment). If empty, resolved from RepoRoot's
	// git "origin" remote (mirrors label-doctor.Doctor()'s Repo
	// resolution exactly).
	Repo string
	// Token is the GitHub token used for both checks. If empty,
	// resolved from $GITHUB_TOKEN then $GH_TOKEN — mirrors
	// label-doctor.Doctor()'s Token resolution order exactly.
	Token string
	// APIBaseURL overrides the GitHub REST API root. Empty defaults to
	// https://api.github.com. Tests point this at an httptest.Server.
	APIBaseURL string
	// SecretNames are the GitHub Actions repository-secret names that
	// must be present. May be empty (e.g. the PAT-free engine tier —
	// see repoinstall.go's requiredSecretNames — needs no repo
	// secret), in which case only push access is checked.
	SecretNames []string
}

// Result records Check's presence-only findings. It never carries a
// secret VALUE, by construction (see package doc comment).
type Result struct {
	// Repo is the resolved "owner/repo" the checks ran against.
	Repo string
	// Present holds every requested secret name found present, in
	// SecretNames order.
	Present []string
	// Missing holds every requested secret name NOT found present, in
	// SecretNames order.
	Missing []string
	// PushAccess reports whether Token has push access to Repo.
	PushAccess bool
}

// Ready reports whether every requested secret is present AND the
// token has push access — the "operator prerequisites satisfied, ok
// to proceed" gate (AC1/AC4).
func (r *Result) Ready() bool {
	return len(r.Missing) == 0 && r.PushAccess
}

// Check runs the two presence-only GitHub REST checks (AC3) and
// returns a Result. It issues only read-only GETs — against
// /repos/{owner}/{repo}/actions/secrets (secret NAMES only, never a
// value) and /repos/{owner}/{repo} (the authenticated caller's
// permissions.push field) — and mutates nothing.
func Check(ctx context.Context, opts Options) (*Result, error) {
	repo := opts.Repo
	if repo == "" {
		r, err := resolveRepoFromGitRemote(ctx, opts.RepoRoot)
		if err != nil {
			return nil, err
		}
		repo = r
	}

	token := opts.Token
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		token = os.Getenv("GH_TOKEN")
	}

	apiBase := opts.APIBaseURL
	if apiBase == "" {
		apiBase = defaultAPIBase
	}

	live, err := ghListSecrets(ctx, apiBase, repo, token)
	if err != nil {
		return nil, fmt.Errorf("install-preflight: %w", err)
	}
	present := make(map[string]bool, len(live))
	for _, s := range live {
		present[s.Name] = true
	}

	res := &Result{Repo: repo}
	for _, name := range opts.SecretNames {
		if present[name] {
			res.Present = append(res.Present, name)
		} else {
			res.Missing = append(res.Missing, name)
		}
	}

	pushAccess, err := ghCheckPushAccess(ctx, apiBase, repo, token)
	if err != nil {
		return nil, fmt.Errorf("install-preflight: %w", err)
	}
	res.PushAccess = pushAccess

	return res, nil
}
