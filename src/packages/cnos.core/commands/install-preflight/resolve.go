package installpreflight

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// githubRemotePattern matches both GitHub remote URL forms `git remote
// get-url origin` can return:
//
//	https://github.com/{owner}/{repo}.git  (or without the trailing .git)
//	git@github.com:{owner}/{repo}.git
//
// and captures owner/repo. Duplicated from
// label-doctor/resolve.go's githubRemotePattern (a separate go.work
// module — not importable directly) rather than shared, per the γ
// scaffold's Implementation contract ("model, don't import directly").
var githubRemotePattern = regexp.MustCompile(`github\.com[:/]([^/]+)/(.+?)(\.git)?/?$`)

// resolveRepoFromGitRemote resolves "owner/repo" from repoRoot's git
// "origin" remote. repoRoot may be any directory inside the checkout —
// git itself walks upward to find .git, so this does not require an
// already-resolved repository root. Mirrors
// label-doctor/resolve.go's resolveRepoFromGitRemote exactly.
func resolveRepoFromGitRemote(ctx context.Context, repoRoot string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	if repoRoot != "" {
		cmd.Dir = repoRoot
	}
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("install-preflight: could not resolve target repo (owner/repo): git remote get-url origin: %w", err)
	}
	url := strings.TrimSpace(string(out))
	m := githubRemotePattern.FindStringSubmatch(url)
	if m == nil {
		return "", fmt.Errorf("install-preflight: could not resolve target repo (owner/repo): unrecognized git remote URL %q", url)
	}
	return m[1] + "/" + m[2], nil
}
