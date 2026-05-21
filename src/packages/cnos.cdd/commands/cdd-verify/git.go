package cddverify

import (
	"context"
	"os/exec"
	"strings"
)

// gitLogAuthor returns the most recent git author email for filePath
// (relative to repoRoot), or "" if the file is not tracked or git is
// unavailable.
//
// Adapter (eng/go §2.6, §3.7): subprocess via exec.CommandContext with
// argv (eng/go §3.10).
func gitLogAuthor(ctx context.Context, repoRoot, filePath string) string {
	cmd := exec.CommandContext(ctx, "git", "log", "-1", "--format=%ae", "--", filePath)
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// gitCommitTime returns ISO 8601 committer time for the most recent
// commit touching filePath, or "" if unavailable.
func gitCommitTime(ctx context.Context, repoRoot, filePath string) string {
	cmd := exec.CommandContext(ctx, "git", "log", "-1", "--format=%cI", "--", filePath)
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// gitRepoRoot returns `git rev-parse --show-toplevel` for cwd, or "" on
// failure. Used when --repo-root is not supplied.
func gitRepoRoot(ctx context.Context) string {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
