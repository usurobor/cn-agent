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
	"runtime"
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

// tarEntry is a single file to write into a fixture tarball, with an
// explicit mode — unlike makeRepoInstallTarGz above (always 0644), the
// dispatch fixture below needs the real cn-install-wake renderer
// script's executable bit preserved so the installed cn binary can
// exec it once restore.Restore extracts it.
type tarEntry struct {
	name string
	mode int64
	data []byte
}

func makeRepoInstallTarGzEntries(t *testing.T, entries []tarEntry) ([]byte, string) {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	for _, e := range entries {
		hdr := &tar.Header{Name: e.name, Mode: e.mode, Size: int64(len(e.data))}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(e.data); err != nil {
			t.Fatal(err)
		}
	}
	tw.Close()
	gw.Close()
	data := buf.Bytes()
	h := sha256.Sum256(data)
	return data, hex.EncodeToString(h[:])
}

// writeCliDispatchFixtureIndex builds a package index + two tarballs:
// cnos.core (carrying the REAL commands/install-wake/cn-install-wake
// renderer script, executable bit preserved) and cnos.cds (carrying the
// REAL orchestrators/cds-dispatch/SKILL.md this cycle edited) — read
// live from the repository source tree via repoGoModuleRoot's sibling
// resolver, so this CLI-level test exercises the actual shipped
// renderer + prose, not a synthetic stand-in. Mirrors
// repoinstall_test.go's writeDispatchFixtureIndex (test-fixture
// duplication across packages, not production parser duplication — per
// eng/go §2.17, already the convention this file's own
// makeRepoInstallTarGz follows relative to repoinstall_test.go's
// makeTarGz).
func writeCliDispatchFixtureIndex(t *testing.T) string {
	t.Helper()
	// repoGoModuleRoot resolves <repoRoot>/src/go; go up two more levels
	// for <repoRoot> itself (src/go -> src -> repoRoot).
	repoRoot := filepath.Join(repoGoModuleRoot(t), "..", "..")
	rendererScript, err := os.ReadFile(filepath.Join(repoRoot, "src", "packages", "cnos.core", "commands", "install-wake", "cn-install-wake"))
	if err != nil {
		t.Fatalf("read renderer script: %v", err)
	}
	skillMD, err := os.ReadFile(filepath.Join(repoRoot, "src", "packages", "cnos.cds", "orchestrators", "cds-dispatch", "SKILL.md"))
	if err != nil {
		t.Fatalf("read cds-dispatch SKILL.md: %v", err)
	}

	dir := t.TempDir()

	coreTar, coreSHA := makeRepoInstallTarGzEntries(t, []tarEntry{
		{name: "cn.package.json", mode: 0644, data: []byte(`{"name": "cnos.core", "version": "9.9.9"}`)},
		{name: "commands/install-wake/cn-install-wake", mode: 0755, data: rendererScript},
	})
	coreTarName := "cnos.core-9.9.9.tar.gz"
	if err := os.WriteFile(filepath.Join(dir, coreTarName), coreTar, 0644); err != nil {
		t.Fatal(err)
	}

	cdsTar, cdsSHA := makeRepoInstallTarGzEntries(t, []tarEntry{
		{name: "cn.package.json", mode: 0644, data: []byte(`{"name": "cnos.cds", "version": "9.9.9"}`)},
		{name: "orchestrators/cds-dispatch/SKILL.md", mode: 0644, data: skillMD},
	})
	cdsTarName := "cnos.cds-9.9.9.tar.gz"
	if err := os.WriteFile(filepath.Join(dir, cdsTarName), cdsTar, 0644); err != nil {
		t.Fatal(err)
	}

	idx := pkg.PackageIndex{
		Schema: "cn.package-index.v1",
		Packages: map[string]map[string]pkg.IndexEntry{
			"cnos.core": {"9.9.9": {URL: coreTarName, SHA256: coreSHA}},
			"cnos.cds":  {"9.9.9": {URL: cdsTarName, SHA256: cdsSHA}},
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

// cnos#610: this test previously asserted --dispatch cds always failed
// through the CLI, naming #609. That guard is exactly what this cell
// removes (see TestRepoInstall_DispatchCds_IdentityFlagsWireThrough
// below for the new positive contract). The fixture here
// (writeFixtureIndex) is a minimal single-package manifest with no
// actual renderer script content, so --dispatch cds still fails today
// — but for a different, still-genuine reason: the renderer binary
// itself was never vendored. No partial .github/workflows/ may exist
// either way.
func TestRepoInstall_DispatchCds_RendererNotVendored_CliWiring(t *testing.T) {
	indexPath := writeFixtureIndex(t, "cnos.core", "9.9.9")

	repoDir := t.TempDir()
	initGitRepo(t, repoDir)
	t.Chdir(repoDir)

	_, stderr, err := runRepoInstall(t, []string{"--index", indexPath, "--packages", "cnos.core", "--dispatch", "cds"})
	if err == nil {
		t.Fatal("expected --dispatch cds to fail when the renderer script is not vendored")
	}
	if !strings.Contains(stderr, "dispatch renderer not found") {
		t.Errorf("stderr should name the missing renderer, got: %q", stderr)
	}
	if _, statErr := os.Stat(filepath.Join(repoDir, ".github")); !os.IsNotExist(statErr) {
		t.Error(".github must not exist after --dispatch cds fails")
	}
}

// AC2/C2: --agent acme with no --workflow-pat-secret fails early through
// the full CLI wiring, before any renderer invocation — no partial
// .github/workflows/ directory.
func TestRepoInstall_DispatchCds_MissingIdentity_CliWiring(t *testing.T) {
	indexPath := writeFixtureIndex(t, "cnos.core", "9.9.9")

	repoDir := t.TempDir()
	initGitRepo(t, repoDir)
	t.Chdir(repoDir)

	_, stderr, err := runRepoInstall(t, []string{"--index", indexPath, "--packages", "cnos.core", "--dispatch", "cds", "--agent", "acme"})
	if err == nil {
		t.Fatal("expected --dispatch cds --agent acme (no --workflow-pat-secret) to fail")
	}
	if !strings.Contains(stderr, "--workflow-pat-secret") {
		t.Errorf("stderr should name the missing flag, got: %q", stderr)
	}
	if _, statErr := os.Stat(filepath.Join(repoDir, ".github")); !os.IsNotExist(statErr) {
		t.Error(".github must not exist after a missing-identity failure")
	}
}

// AC1/C1/C3 + AC4/C4/C6 + AC5-positive, through the FULL CLI flag
// surface (--agent / --workflow-pat-secret / --bot-name / --bot-id),
// proving cmd_repo_install.go's Args→Options pass-through actually
// wires these four new flags through to repoinstall.Run, not just that
// repoinstall.Run itself accepts them (repoinstall_test.go already
// covers that at the package-internal level). Builds the real cn
// binary and drives it as a subprocess against a fixture vendoring the
// REAL cn-install-wake renderer + the REAL cds-dispatch/SKILL.md.
func TestRepoInstall_DispatchCds_IdentityFlagsWireThrough(t *testing.T) {
	binPath := buildCnBinary(t)
	indexPath := writeCliDispatchFixtureIndex(t)

	repoDir := t.TempDir()
	initGitRepo(t, repoDir)

	cmd := exec.Command(binPath, "repo", "install",
		"--index", indexPath,
		"--packages", "cnos.core,cnos.cds",
		"--dispatch", "cds",
		"--agent", "acme",
		"--workflow-pat-secret", "ACME_WORKFLOW_PAT",
		"--bot-name", "acme-bot",
		"--bot-id", "12345678",
	)
	cmd.Dir = repoDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()

	// AC3's cnos#493 label gap is still expected today (see
	// repoinstall_test.go's TestRun_DispatchCds_RendersWorkflow_
	// ThenSurfacesLabelGap for why this is the correct current
	// behavior) — what THIS test proves is that the CLI flags actually
	// reached the renderer, which only happens if Args→Options
	// pass-through works.
	if runErr == nil {
		t.Fatalf("expected the cnos#493 label-gap error, got success\nstdout: %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "cnos#493") {
		t.Errorf("stderr should carry the cnos#493 diagnostic, got: %q", stderr.String())
	}

	workflowPath := filepath.Join(repoDir, ".github", "workflows", "cnos-cds-dispatch.yml")
	data, statErr := os.ReadFile(workflowPath)
	if statErr != nil {
		t.Fatalf("expected rendered workflow at %s: %v\nstdout: %s\nstderr: %s", workflowPath, statErr, stdout.String(), stderr.String())
	}
	content := string(data)

	// The four CLI flags must have reached the render (proves Args→
	// Options wiring, not just that repoinstall.Run accepts them).
	for _, want := range []string{"ACME_WORKFLOW_PAT", "acme-bot", "12345678", "cds-dispatch-acme"} {
		if !strings.Contains(content, want) {
			t.Errorf("rendered workflow missing identity value %q (CLI flag did not reach the renderer):\n%s", want, content)
		}
	}

	// AC4/C4 + AC5 positive: zero sigma leak for this non-sigma agent.
	for _, leak := range []string{"SIGMA_WORKFLOW_PAT", "41898282", "cds-dispatch-sigma", "today: `sigma`", "agent-admin-sigma", "cn-sigma:"} {
		if strings.Contains(content, leak) {
			t.Errorf("acme render (via full CLI wiring) leaks %q", leak)
		}
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

// --- real dispatch-path regression (cnos#608 β R0 finding) ---
//
// Every test above drives RepoInstallCmd.Run directly via
// runRepoInstall, which hand-constructs an Invocation{HubPath: ""}. That
// bypasses main.go's discoverHub() entirely — the function that
// populates inv.HubPath in the real `cn` binary by walking upward from
// cwd looking for ANY ancestor .cn/ directory, unbounded by the current
// git repository's root and without ever checking that the found
// directory is itself a git repository. β's R0 review found that
// cmd_repo_install.go used to trust a non-empty inv.HubPath as the
// install root outright, skipping the git-repo-root check entirely —
// and that no test in the original diff could have caught it, because
// the test harness never exercises the inv.HubPath != "" branch through
// a faithful discoverHub()-shaped path.
//
// The tests below close that gap by building the actual cn binary and
// invoking it as a subprocess — main.go's real discoverHub() runs, not a
// reimplementation of it — against a cwd that (a) is not a git repo at
// all and (b) has an unrelated ancestor .cn/ directory that discoverHub
// would find.

// buildCnBinary builds the real `cn` binary once per test (from the
// module root, mirroring the `(cd src/go && go build -o $CN_BIN
// ./cmd/cn)` convention already used by the shell test harnesses under
// src/packages/cnos.cdd/commands/cdd-verify/) and returns its path.
func buildCnBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	binPath := filepath.Join(dir, "cn")
	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/cn")
	cmd.Dir = repoGoModuleRoot(t)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build ./cmd/cn: %v\n%s", err, out)
	}
	return binPath
}

// repoGoModuleRoot resolves the src/go module root from this test file's
// own path (runtime.Caller), independent of the process's current
// working directory — several sibling tests in this file call
// t.Chdir, and this helper must not depend on ordering relative to
// those.
func repoGoModuleRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	// thisFile: .../src/go/internal/cli/cmd_repo_install_test.go
	return filepath.Join(filepath.Dir(thisFile), "..", "..")
}

// AC1 (real dispatch path): cwd is NOT inside any git repository at all,
// but an unrelated ancestor directory happens to contain a .cn/
// directory (e.g. a different project, a stale hub, a parent
// workspace). main.go's discoverHub() will find that ancestor and
// populate inv.HubPath with it. cn repo install must still fail with
// AC1's exact error — it must never silently treat the discovered
// ancestor as the install root, and must not touch it.
func TestRepoInstall_RealDispatch_AncestorHubOutsideGitRepo_FailsClearly(t *testing.T) {
	binPath := buildCnBinary(t)

	// Layout:
	//   outer/.cn/                        <- unrelated ancestor "hub"; not a git repo itself
	//   outer/nested/not-a-git-repo/      <- cwd; not inside any git repo
	outer := t.TempDir()
	ancestorHub := filepath.Join(outer, ".cn")
	if err := os.MkdirAll(ancestorHub, 0755); err != nil {
		t.Fatal(err)
	}
	cwd := filepath.Join(outer, "nested", "not-a-git-repo")
	if err := os.MkdirAll(cwd, 0755); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binPath, "repo", "install", "--dry-run")
	cmd.Dir = cwd
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()

	if runErr == nil {
		t.Fatalf("expected `cn repo install` to fail outside a git repository (ancestor .cn found at %s); stdout:\n%s\nstderr:\n%s", ancestorHub, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "✗ cn repo install must be run inside a Git repository.") {
		t.Errorf("stderr = %q, want the exact AC1 git-repo-required message", stderr.String())
	}

	// The ancestor "hub" must be completely untouched — no silent
	// walk-up install, no write of any kind under outer/ (beyond the
	// two directories this test itself created).
	hubEntries, rerr := os.ReadDir(ancestorHub)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if len(hubEntries) != 0 {
		t.Errorf("ancestor .cn/ must not be written to; found: %v", hubEntries)
	}
	outerEntries, rerr := os.ReadDir(outer)
	if rerr != nil {
		t.Fatal(rerr)
	}
	for _, e := range outerEntries {
		if e.Name() != ".cn" && e.Name() != "nested" {
			t.Errorf("unexpected entry written under the ancestor directory: %s", e.Name())
		}
	}
}

// AC1 (real dispatch path, control case): the same ancestor-.cn/ layout,
// but cwd IS inside its own git repository (nested under the ancestor
// hub, mirroring a monorepo-of-repos layout). discoverHub() still finds
// the unrelated ancestor .cn/ first (it walks up from cwd and stops at
// the first .cn/ found, which sits above the inner git repo's root), so
// this must also resolve to the INNER repo's root — never the
// ancestor's — and must succeed there, not silently install into the
// ancestor.
func TestRepoInstall_RealDispatch_AncestorHubWithNestedGitRepo_UsesInnerRoot(t *testing.T) {
	binPath := buildCnBinary(t)
	indexPath := writeFixtureIndex(t, "cnos.core", "9.9.9")

	outer := t.TempDir()
	ancestorHub := filepath.Join(outer, ".cn")
	if err := os.MkdirAll(ancestorHub, 0755); err != nil {
		t.Fatal(err)
	}
	innerRepo := filepath.Join(outer, "nested", "actual-git-repo")
	if err := os.MkdirAll(innerRepo, 0755); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, innerRepo)

	cmd := exec.Command(binPath, "repo", "install", "--index", indexPath, "--packages", "cnos.core")
	cmd.Dir = innerRepo
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("cn repo install: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "Git repository root: "+innerRepo) {
		t.Errorf("stdout should report the inner repo as the resolved root, got:\n%s", stdout.String())
	}

	if _, statErr := os.Stat(filepath.Join(innerRepo, ".cn", "deps.json")); statErr != nil {
		t.Errorf("expected .cn/deps.json to be written under the inner repo root: %v", statErr)
	}

	ancestorEntries, rerr := os.ReadDir(ancestorHub)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if len(ancestorEntries) != 0 {
		t.Errorf("ancestor .cn/ must not be written to; found: %v", ancestorEntries)
	}
}
