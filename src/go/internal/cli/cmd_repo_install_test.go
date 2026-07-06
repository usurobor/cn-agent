package cli

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/pkg"
)

// --- fixtures ---

func makeRepoInstallTarGz(t *testing.T, files map[string]string) ([]byte, string) {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	for name, content := range files {
		hdr := &tar.Header{Name: name, Mode: 0644, Size: int64(len(content))}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	tw.Close()
	gw.Close()
	data := buf.Bytes()
	h := sha256.Sum256(data)
	return data, hex.EncodeToString(h[:])
}

// writeFixtureIndex writes a single-package, single-version index.json +
// tarball to a temp dir (relative URL, like `cn build` produces) and
// returns the index path.
func writeFixtureIndex(t *testing.T, name, version string) string {
	t.Helper()
	dir := t.TempDir()
	manifest := `{"name": "` + name + `", "version": "` + version + `"}`
	tarData, sha := makeRepoInstallTarGz(t, map[string]string{"cn.package.json": manifest})
	tarName := name + "-" + version + ".tar.gz"
	if err := os.WriteFile(filepath.Join(dir, tarName), tarData, 0644); err != nil {
		t.Fatal(err)
	}
	idx := pkg.PackageIndex{
		Schema: "cn.package-index.v1",
		Packages: map[string]map[string]pkg.IndexEntry{
			name: {version: {URL: tarName, SHA256: sha}},
		},
	}
	indexPath := filepath.Join(dir, "index.json")
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		t.Fatal(err)
	}
	return indexPath
}

// initGitRepo runs `git init` (+ minimal identity config) in dir.
func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init", "-q")
	runGit(t, dir, "config", "user.email", "test@cnos.test")
	runGit(t, dir, "config", "user.name", "cnos-test")
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

func runRepoInstall(t *testing.T, args []string) (string, string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmd := &RepoInstallCmd{}
	err := cmd.Run(context.Background(), Invocation{
		Args:   args,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	return stdout.String(), stderr.String(), err
}

// --- dispatch wiring ---

// AC1 (dispatch shape): "cn repo install" resolves via the kernel
// registry's noun-verb resolution to RepoInstallCmd, the same mechanism
// already used by "cn cell finalize" / "cn issues fsm".
func TestRepoInstall_ResolvesViaNounVerb(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&RepoInstallCmd{})

	res := ResolveCommand(reg, []string{"repo", "install", "--dry-run"})
	if res.Command == nil {
		t.Fatal("expected 'repo install' to resolve to a command")
	}
	if res.Command.Spec().Name != "repo-install" {
		t.Errorf("resolved command = %q, want repo-install", res.Command.Spec().Name)
	}
	if got := res.Command.Spec(); got.Source != SourceKernel || got.Tier != TierKernel || got.NeedsHub {
		t.Errorf("spec = %+v, want kernel/kernel-tier/NeedsHub=false", got)
	}
}

func TestRepoInstall_HelpFlag(t *testing.T) {
	stdout, _, err := runRepoInstall(t, []string{"--help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"cn repo install", "--release", "--index", "--packages", "--dispatch", "--dry-run"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("--help output missing %q:\n%s", want, stdout)
		}
	}
}

func TestRepoInstall_UnknownFlag(t *testing.T) {
	_, stderr, err := runRepoInstall(t, []string{"--bogus"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
	if !strings.Contains(stderr, "--bogus") {
		t.Errorf("stderr should name the unknown flag, got: %q", stderr)
	}
}

// AC1: fails with a clear error if not at a git repo root — does not
// silently walk up or scaffold anything.
func TestRepoInstall_NotAGitRepo_FailsClearly(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	_, stderr, err := runRepoInstall(t, []string{"--dry-run"})
	if err == nil {
		t.Fatal("expected error outside a git repository")
	}
	if !strings.Contains(stderr, "✗ cn repo install must be run inside a Git repository.") {
		t.Errorf("stderr = %q, want the exact git-repo-required message", stderr)
	}
	entries, rerr := os.ReadDir(dir)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if len(entries) != 0 {
		t.Errorf("must not scaffold anything when not a git repo; found: %v", entries)
	}
}

// AC1 end-to-end: a fresh `git init`-only repo, no .cn/ hub, no vendor
// packages, no prior state — cn repo install succeeds and the three-file
// diff exists.
func TestRepoInstall_FreshGitRepo_EndToEnd(t *testing.T) {
	indexPath := writeFixtureIndex(t, "cnos.core", "9.9.9")

	repoDir := t.TempDir()
	initGitRepo(t, repoDir)
	t.Chdir(repoDir)

	stdout, stderr, err := runRepoInstall(t, []string{
		"--index", indexPath,
		"--packages", "cnos.core",
	})
	if err != nil {
		t.Fatalf("Run: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}

	for _, f := range []string{
		filepath.Join(".cn", "deps.json"),
		filepath.Join(".cn", "deps.lock.json"),
		".gitignore",
	} {
		if _, statErr := os.Stat(filepath.Join(repoDir, f)); statErr != nil {
			t.Errorf("expected %s to exist: %v", f, statErr)
		}
	}
	if _, statErr := os.Stat(filepath.Join(repoDir, ".cn", "vendor", "packages", "cnos.core", "cn.package.json")); statErr != nil {
		t.Errorf("cnos.core not restored: %v", statErr)
	}
}

// AC5: cn repo install; git diff --exit-code; cn repo install; git diff
// --exit-code — both must exit 0, and (the stronger form) git status
// --porcelain is empty after committing once and re-running.
func TestRepoInstall_Idempotent_NoDiffOnSecondRun(t *testing.T) {
	indexPath := writeFixtureIndex(t, "cnos.core", "9.9.9")

	repoDir := t.TempDir()
	initGitRepo(t, repoDir)
	t.Chdir(repoDir)

	args := []string{"--index", indexPath, "--packages", "cnos.core"}

	if _, stderr, err := runRepoInstall(t, args); err != nil {
		t.Fatalf("first install: %v\n%s", err, stderr)
	}
	runGit(t, repoDir, "diff", "--exit-code")

	runGit(t, repoDir, "add", ".cn/deps.json", ".cn/deps.lock.json", ".gitignore")
	runGit(t, repoDir, "commit", "-q", "-m", "install cds base layer")

	if _, stderr, err := runRepoInstall(t, args); err != nil {
		t.Fatalf("second install: %v\n%s", err, stderr)
	}
	runGit(t, repoDir, "diff", "--exit-code")

	status := runGit(t, repoDir, "status", "--porcelain")
	if strings.TrimSpace(status) != "" {
		t.Errorf("git status --porcelain not empty after second install:\n%s", status)
	}
}

// AC6: --dry-run writes nothing; git status --porcelain stays empty.
func TestRepoInstall_DryRun_GitStatusStaysClean(t *testing.T) {
	indexPath := writeFixtureIndex(t, "cnos.core", "9.9.9")

	repoDir := t.TempDir()
	initGitRepo(t, repoDir)
	t.Chdir(repoDir)

	stdout, stderr, err := runRepoInstall(t, []string{"--index", indexPath, "--packages", "cnos.core", "--dry-run"})
	if err != nil {
		t.Fatalf("Run: %v\n%s", err, stderr)
	}
	for _, want := range []string{".cn/deps.json", ".cn/deps.lock.json", ".gitignore"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("dry-run stdout missing %q:\n%s", want, stdout)
		}
	}

	status := runGit(t, repoDir, "status", "--porcelain")
	if strings.TrimSpace(status) != "" {
		t.Errorf("git status --porcelain not empty after --dry-run:\n%s", status)
	}
}

// AC9: --dispatch cds fails explicitly through the full CLI wiring, with
// no partial .github/workflows/cnos-cds-dispatch.yml.
func TestRepoInstall_DispatchCds_CliWiring(t *testing.T) {
	indexPath := writeFixtureIndex(t, "cnos.core", "9.9.9")

	repoDir := t.TempDir()
	initGitRepo(t, repoDir)
	t.Chdir(repoDir)

	_, stderr, err := runRepoInstall(t, []string{"--index", indexPath, "--dispatch", "cds"})
	if err == nil {
		t.Fatal("expected --dispatch cds to fail")
	}
	if !strings.Contains(stderr, "#609") {
		t.Errorf("stderr should name #609, got: %q", stderr)
	}
	if _, statErr := os.Stat(filepath.Join(repoDir, ".github")); !os.IsNotExist(statErr) {
		t.Error(".github must not exist after --dispatch cds fails")
	}
}

// AC10 regression guard: cn repo install must not scaffold the agent-hub
// tree that cn init writes (spec/SOUL.md, agent/, threads/, state/).
func TestRepoInstall_NoAgentHubScaffold(t *testing.T) {
	indexPath := writeFixtureIndex(t, "cnos.core", "9.9.9")

	repoDir := t.TempDir()
	initGitRepo(t, repoDir)
	t.Chdir(repoDir)

	if _, stderr, err := runRepoInstall(t, []string{"--index", indexPath, "--packages", "cnos.core"}); err != nil {
		t.Fatalf("Run: %v\n%s", err, stderr)
	}
	for _, p := range []string{"spec", "agent", "threads", "state"} {
		if _, statErr := os.Stat(filepath.Join(repoDir, p)); !os.IsNotExist(statErr) {
			t.Errorf("agent-hub path %q must not exist after cn repo install", p)
		}
	}
}
