package labeldoctor

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
// and captures owner/repo. This utility does not exist anywhere in
// src/go/internal today (confirmed absent by the γ scaffold's
// grep -rn "git remote\|ParseGitHubRemote\|RepoFromRemote" across
// src/go/internal — no hits); repoinstall.Options.Repo is a different
// value entirely (the source of cnos releases, not the installing
// target's owner/repo — see the γ scaffold's Implementation contract →
// Runtime dependencies row).
var githubRemotePattern = regexp.MustCompile(`github\.com[:/]([^/]+)/(.+?)(\.git)?/?$`)

// resolveRepoFromGitRemote resolves "owner/repo" from repoRoot's git
// "origin" remote. repoRoot may be any directory inside the checkout —
// git itself walks upward to find .git, so this does not require an
// already-resolved repository root.
func resolveRepoFromGitRemote(ctx context.Context, repoRoot string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	if repoRoot != "" {
		cmd.Dir = repoRoot
	}
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("label-doctor: could not resolve target repo (owner/repo): git remote get-url origin: %w", err)
	}
	url := strings.TrimSpace(string(out))
	m := githubRemotePattern.FindStringSubmatch(url)
	if m == nil {
		return "", fmt.Errorf("label-doctor: could not resolve target repo (owner/repo): unrecognized git remote URL %q", url)
	}
	return m[1] + "/" + m[2], nil
}
