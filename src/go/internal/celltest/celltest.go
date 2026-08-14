// Package celltest holds the fixtures that more than one cell package's tests
// need. It is imported only from _test files, so nothing here reaches a built
// binary.
//
// It exists because these were written seven times. The git runner in
// particular was copied into every package that needed a repository, and the
// copies had already drifted: three `testRepo`s differed in what they returned,
// and one trimmed its output where another did not. A fixture that differs by
// accident makes two tests measure two things while reading as one.
//
// What belongs here is the MECHANISM several packages share, not every helper a
// test uses. A fixture with one consumer belongs beside its test, where a
// reader can see what it sets up without leaving the file — so the shapes built
// on top of Git (an empty-commit repo, a repo with a base tree and a change on
// top) stay in the packages that need them.
package celltest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Git runs one git command in dir and returns its trimmed combined output,
// failing the test if git does.
//
// The identity is fixed rather than inherited: a machine with no configured
// user cannot commit, and a test that passes only where the developer happens
// to have a .gitconfig is a test that reports on the machine.
func Git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

// Repo builds a one-commit git repository carrying a README, and returns its
// path and the commit it pinned. The commit is RETURNED rather than left for
// the caller to resolve: every consumer measures something against a pinned
// base, and resolving `HEAD` later would pin whatever the test did next.
func Repo(t *testing.T) (dir, head string) {
	t.Helper()
	dir = t.TempDir()
	Git(t, dir, "init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("base\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	Git(t, dir, "add", "-A")
	Git(t, dir, "commit", "-qm", "base")
	return dir, Git(t, dir, "rev-parse", "HEAD")
}

// Skill writes one SKILL.md into an installed package tree rooted at root, at
// the path `<pkg>/skills/<path>/SKILL.md` that a `pkg:path` ref resolves to.
func Skill(t *testing.T, root, ref, body string) {
	t.Helper()
	pkg, path, _ := strings.Cut(ref, ":")
	dir := filepath.Join(root, pkg, "skills", filepath.FromSlash(path))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}
