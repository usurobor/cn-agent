// Package cellwork materializes the working copy a producing seat needs and
// measures what it did with it.
//
// This is the matter substrate seam (G1), and it lives OUTSIDE the kernel on
// purpose: the kernel owns no git, no repo and no branch, so a worktree is
// prepared here, handed to a seat as a plain directory, and reduced back to
// one artifact of text — a unified diff — before anything crosses into an
// episode.
//
// The load-bearing property is that the runtime measures the change instead
// of believing a claim about it. A seat may say whatever it likes in its
// answer; the diff this package computes is what the record carries, and a
// seat that changed nothing produces no diff to carry.
package cellwork

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	gitTimeout   = 2 * time.Minute
	maxDiffBytes = 1 << 20 // matches the kernel's matter bound
)

// Worktree is a disposable checkout of a repository at one base commit.
// Nothing outside it is writable by the seat that receives it.
type Worktree struct {
	Dir     string // absolute path the seat may modify
	Repo    string // repository the worktree was cut from
	BaseSHA string // resolved commit the work starts from
}

// ResolveBase pins a revision to its exact commit SHA. Construction calls it
// so a declaration recorded as "resolved" names a commit and not a moving
// name like HEAD; the episode then materializes at that pinned SHA.
func ResolveBase(ctx context.Context, repo, base string) (string, error) {
	repoAbs, err := repoPath(repo)
	if err != nil {
		return "", err
	}
	sha, err := git(ctx, repoAbs, "rev-parse", "--verify", base+"^{commit}")
	if err != nil {
		return "", fmt.Errorf("cellwork: base %q does not resolve in %s: %w", base, repoAbs, err)
	}
	return strings.TrimSpace(sha), nil
}

// RepoRoot reports the top level of the repository containing dir. The
// composition root uses it to anchor package resolution when no hub is
// present (the cnos#593 fallback), so nothing resolves relative to the
// process working directory.
func RepoRoot(ctx context.Context, dir string) (string, error) {
	out, err := git(ctx, dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("cellwork: no git repository at %s: %w", dir, err)
	}
	return strings.TrimSpace(out), nil
}

func repoPath(repo string) (string, error) {
	abs, err := filepath.Abs(repo)
	if err != nil {
		return "", fmt.Errorf("cellwork: resolve repo path: %w", err)
	}
	if _, err := os.Stat(filepath.Join(abs, ".git")); err != nil {
		return "", fmt.Errorf("cellwork: %s is not a git repository: %w", abs, err)
	}
	return abs, nil
}

// Materialize cuts a detached worktree of repo at base and returns it with a
// release function. base may be any revision git resolves; the resolved SHA is
// recorded so the episode binds the exact commit, not a moving name.
//
// The caller must always call release, which removes the worktree and the
// seat's ability to touch anything.
func Materialize(ctx context.Context, repo, base string) (Worktree, func(), error) {
	repoAbs, err := repoPath(repo)
	if err != nil {
		return Worktree{}, nil, err
	}
	sha, err := ResolveBase(ctx, repoAbs, base)
	if err != nil {
		return Worktree{}, nil, err
	}

	dir, err := os.MkdirTemp("", "cnos-cell-worktree-")
	if err != nil {
		return Worktree{}, nil, fmt.Errorf("cellwork: create worktree dir: %w", err)
	}
	// git insists on creating the leaf itself.
	wtDir := filepath.Join(dir, "wt")
	if _, err := git(ctx, repoAbs, "worktree", "add", "--detach", wtDir, sha); err != nil {
		os.RemoveAll(dir)
		return Worktree{}, nil, fmt.Errorf("cellwork: add worktree at %s: %w", sha, err)
	}

	release := func() {
		// Best effort, in order: git forgets the worktree, then the bytes go.
		rmCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), gitTimeout)
		defer cancel()
		_, _ = git(rmCtx, repoAbs, "worktree", "remove", "--force", wtDir)
		os.RemoveAll(dir)
	}
	return Worktree{Dir: wtDir, Repo: repoAbs, BaseSHA: sha}, release, nil
}

// Diff reduces whatever the seat did to one unified diff against the base,
// including files it created. An empty result means the seat changed nothing —
// the caller must not manufacture evidence from it.
func (w Worktree) Diff(ctx context.Context) (string, error) {
	// Staging everything is what makes new files visible to `diff`; the index
	// belongs to this disposable worktree alone.
	if _, err := git(ctx, w.Dir, "add", "-A"); err != nil {
		return "", fmt.Errorf("cellwork: stage worktree: %w", err)
	}
	out, err := git(ctx, w.Dir, "diff", "--cached", "--no-color")
	if err != nil {
		return "", fmt.Errorf("cellwork: compute diff: %w", err)
	}
	if len(out) > maxDiffBytes {
		return "", fmt.Errorf("cellwork: diff exceeds %d bytes", maxDiffBytes)
	}
	return out, nil
}

// git runs one git command in dir and returns stdout.
func git(ctx context.Context, dir string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, gitTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", errors.New("git " + strings.Join(args, " ") + ": " + msg)
	}
	return stdout.String(), nil
}
