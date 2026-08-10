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
	gitTimeout     = 2 * time.Minute
	maxDiffBytes   = 1 << 20 // matches the kernel's matter bound
	maxRefBytes    = 8 << 10 // shas, paths, and git's own chatter
	maxStderrBytes = 8 << 10 // diagnostics only
	waitDelay      = 2 * time.Second
)

// Worktree is a disposable checkout of a repository at one base commit. It is
// the only place whose changes are measured — a seat is pointed at it, not
// confined to it by the operating system, so the guarantee is evidential
// rather than physical: what happens elsewhere never becomes evidence.
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
	sha, err := git(ctx, repoAbs, maxRefBytes, "rev-parse", "--verify", base+"^{commit}")
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
	out, err := git(ctx, dir, maxRefBytes, "rev-parse", "--show-toplevel")
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
	if _, err := git(ctx, repoAbs, maxRefBytes, "worktree", "add", "--detach", wtDir, sha); err != nil {
		os.RemoveAll(dir)
		return Worktree{}, nil, fmt.Errorf("cellwork: add worktree at %s: %w", sha, err)
	}

	release := func() {
		// Best effort, in order: git forgets the worktree, then the bytes go.
		// Failures are not surfaced: cleanup is housekeeping, and an episode's
		// truth does not depend on it.
		rmCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), gitTimeout)
		defer cancel()
		_, _ = git(rmCtx, repoAbs, maxRefBytes, "worktree", "remove", "--force", wtDir)
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
	if _, err := git(ctx, w.Dir, maxRefBytes, "add", "-A"); err != nil {
		return "", fmt.Errorf("cellwork: stage worktree: %w", err)
	}
	out, err := git(ctx, w.Dir, maxDiffBytes, "diff", "--cached", "--no-color")
	if err != nil {
		return "", fmt.Errorf("cellwork: compute diff: %w", err)
	}
	return out, nil
}

// git runs one git command in dir and returns stdout, capturing at most max
// bytes. The bound is applied AS THE CHILD WRITES, not after it exits: a
// repository can produce a diff far larger than memory, and a limit checked
// on a fully buffered result is not a limit.
func git(ctx context.Context, dir string, max int, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, gitTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	stdout := &boundedBuffer{max: max}
	stderr := &boundedBuffer{max: maxStderrBytes}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	cmd.WaitDelay = waitDelay
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", errors.New("git " + strings.Join(args, " ") + ": " + msg)
	}
	if stdout.truncated {
		return "", fmt.Errorf("git %s produced more than %d bytes", strings.Join(args, " "), max)
	}
	return stdout.String(), nil
}

// boundedBuffer captures at most max bytes and remembers that it had to stop.
// It never reports a short write, so the bound fails the command here rather
// than killing the child mid-stream with a broken pipe. (Deliberately the same
// small pattern the provider adapters use; the two adapters share no
// dependency worth creating for fifteen lines.)
type boundedBuffer struct {
	max       int
	buf       bytes.Buffer
	truncated bool
}

func (b *boundedBuffer) Write(p []byte) (int, error) {
	switch room := b.max - b.buf.Len(); {
	case room >= len(p):
		b.buf.Write(p)
	case room > 0:
		b.buf.Write(p[:room])
		b.truncated = true
	default:
		b.truncated = true
	}
	return len(p), nil
}

func (b *boundedBuffer) String() string { return b.buf.String() }
