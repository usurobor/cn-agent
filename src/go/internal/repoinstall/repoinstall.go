// Package repoinstall implements `cn repo install` — the base CNOS/CDS
// package installer (cnos#608).
//
// It makes an arbitrary Git repository CDS-ready with one command: it
// resolves a cnos release (or an explicit package index), writes a
// deterministic .cn/deps.json + .cn/deps.lock.json, restores the default
// package set (cnos.core, cnos.cdd, cnos.cds) under .cn/vendor/packages/,
// and ensures .gitignore excludes the vendor tree.
//
// This package reuses the existing lock/restore substrate directly
// (internal/restore) rather than reimplementing SHA-256 verification or
// lockfile generation — see restore.GenerateLockFromIndex and
// restore.Restore. It does not write any agent-hub scaffold (cf.
// internal/hubinit, which cn init uses) and it never touches
// .github/workflows/ in base mode.
//
// --dispatch cds (cnos#610) layers a dispatch-workflow render on top of
// the base install: after the base install completes, it invokes the
// vendored cn-install-wake renderer (cnos#609) against the cds-dispatch
// wake manifest, requiring an explicit caller identity (--agent /
// --workflow-pat-secret / --bot-name / --bot-id) for any non-sigma
// agent. It never pushes to a remote (PR-only, by construction — this
// package never calls git). It ensures the canonical cnos.core labels
// via an in-process call into packages/cnos.core/commands/label-doctor
// (cnos#493); if that mechanism cannot resolve the installing repo's
// target (no git remote) or cannot reach the GitHub API, it surfaces a
// named, actionable error rather than silently skipping the labels
// obligation.
//
// This package is cli/-boundary compliant per eng/go §2.18: all domain
// logic lives here, and cli/cmd_repo_install.go is a thin wrapper.
package repoinstall

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	labeldoctor "github.com/usurobor/cnos/packages/cnos.core/commands/label-doctor"
	"github.com/usurobor/cnos/src/go/internal/binupdate"
	"github.com/usurobor/cnos/src/go/internal/hubsetup"
	"github.com/usurobor/cnos/src/go/internal/pkg"
	"github.com/usurobor/cnos/src/go/internal/restore"
)

// DefaultRepo is the GitHub "owner/repo" slug releases are resolved
// against absent an override (tests point this at a fixture server).
const DefaultRepo = "usurobor/cnos"

const (
	defaultAPIBaseURL  = "https://api.github.com"
	defaultDownloadURL = "https://github.com"

	// manifestSchema/manifestProfile are the exact wire values written to
	// .cn/deps.json. Per the implementation contract these must not change
	// shape — pkg.Manifest is the canonical type (schema "cn.deps.v1").
	manifestSchema  = "cn.deps.v1"
	manifestProfile = "cds"
)

// DefaultPackages is the base package set `cn repo install` pins absent
// an explicit --packages override. Order is significant: it is the order
// written to .cn/deps.json (manifest order is operator-controlled, per
// restore.GenerateLockFromIndex's own doc comment; the lockfile it
// generates is independently sorted for determinism).
var DefaultPackages = []string{"cnos.core", "cnos.cdd", "cnos.cds"}

// httpClientDefault mirrors the timeout used by restore.go/binupdate.go
// (OCaml curl flags: --connect-timeout 10 --max-time 300).
var httpClientDefault = &http.Client{Timeout: 300 * time.Second}

// Args is the parsed flag set for `cn repo install`.
type Args struct {
	Release   string   // "" or "latest" resolves the newest release; anything else pins a tag
	IndexPath string   // --index override: local path or http(s) URL
	Packages  []string // --packages csv override; nil uses DefaultPackages
	Dispatch  string   // "" / "none" (default) or "cds"
	DryRun    bool

	// Agent/WorkflowPatSecret/BotName/BotID are the --dispatch cds
	// identity flags (cnos#610). Names mirror cn-install-wake's own
	// flag names 1:1 (src/packages/cnos.core/commands/install-wake/
	// cn-install-wake) so the same identity vocabulary flows through
	// unchanged. Unused when Dispatch != "cds".
	Agent             string // default: "sigma" when empty
	WorkflowPatSecret string // required for any non-sigma Agent
	BotName           string // overrides the renderer's agent_bot_name() lookup
	BotID             string // overrides the renderer's agent_bot_id() lookup
}

// ParseArgs parses the `cn repo install` flag set. Two-token
// "--flag value" form only, matching the convention in
// internal/cell.ParseFinalizeArgs.
func ParseArgs(argv []string) (Args, error) {
	var a Args
	for i := 0; i < len(argv); i++ {
		switch argv[i] {
		case "--dry-run":
			a.DryRun = true
		case "--release":
			v, err := takeValue(argv, &i, "--release")
			if err != nil {
				return a, err
			}
			a.Release = v
		case "--index":
			v, err := takeValue(argv, &i, "--index")
			if err != nil {
				return a, err
			}
			a.IndexPath = v
		case "--packages":
			v, err := takeValue(argv, &i, "--packages")
			if err != nil {
				return a, err
			}
			a.Packages = splitPackages(v)
		case "--dispatch":
			v, err := takeValue(argv, &i, "--dispatch")
			if err != nil {
				return a, err
			}
			a.Dispatch = v
		case "--agent":
			v, err := takeValue(argv, &i, "--agent")
			if err != nil {
				return a, err
			}
			a.Agent = v
		case "--workflow-pat-secret":
			v, err := takeValue(argv, &i, "--workflow-pat-secret")
			if err != nil {
				return a, err
			}
			a.WorkflowPatSecret = v
		case "--bot-name":
			v, err := takeValue(argv, &i, "--bot-name")
			if err != nil {
				return a, err
			}
			a.BotName = v
		case "--bot-id":
			v, err := takeValue(argv, &i, "--bot-id")
			if err != nil {
				return a, err
			}
			a.BotID = v
		default:
			return a, fmt.Errorf("unknown flag %q", argv[i])
		}
	}
	return a, nil
}

func takeValue(argv []string, i *int, flag string) (string, error) {
	if *i+1 >= len(argv) {
		return "", fmt.Errorf("%s requires a value", flag)
	}
	*i++
	return argv[*i], nil
}

func splitPackages(csv string) []string {
	var out []string
	for _, p := range strings.Split(csv, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// Options carries the runtime configuration for a `cn repo install` run.
type Options struct {
	// RepoRoot is the resolved Git repository root to install into.
	// Required — git-root detection is the caller's responsibility
	// (cli/cmd_repo_install.go, mirroring cli/cmd_cell.go's gitRepoRoot
	// precedent) so this package stays testable with a plain t.TempDir().
	RepoRoot string

	Release   string
	IndexPath string
	Packages  []string
	Dispatch  string
	DryRun    bool

	// Agent/WorkflowPatSecret/BotName/BotID — see Args' field docs.
	Agent             string
	WorkflowPatSecret string
	BotName           string
	BotID             string

	// Repo is the "owner/repo" slug used to resolve releases and
	// construct download URLs. Defaults to DefaultRepo.
	Repo string
	// HTTPClient, APIBaseURL, DownloadBase allow tests to point release
	// resolution and index/tarball fetches at an httptest.Server, mirroring
	// binupdate.Options' equivalent test seams.
	HTTPClient   *http.Client
	APIBaseURL   string
	DownloadBase string

	Stdout io.Writer
	Stderr io.Writer
}

// Result records the outcome of a Run, for tests that want structured
// assertions beyond stdout scraping.
type Result struct {
	// ReleaseTag is the resolved release tag. Empty when an explicit
	// --index value was used with no --release override (index-only flow).
	ReleaseTag string
	Manifest   pkg.Manifest
	DryRun     bool
}

// Run executes `cn repo install` against opts.RepoRoot.
func Run(ctx context.Context, opts Options) (*Result, error) {
	if err := validateDispatch(opts.Dispatch); err != nil {
		fmt.Fprintf(opts.Stderr, "✗ %s\n", err)
		return nil, err
	}
	if opts.RepoRoot == "" {
		return nil, fmt.Errorf("repo install: RepoRoot is required")
	}

	repo := opts.Repo
	if repo == "" {
		repo = DefaultRepo
	}
	apiBase := opts.APIBaseURL
	if apiBase == "" {
		apiBase = defaultAPIBaseURL
	}
	dlBase := opts.DownloadBase
	if dlBase == "" {
		dlBase = defaultDownloadURL
	}
	client := opts.HTTPClient
	if client == nil {
		client = httpClientDefault
	}

	names := opts.Packages
	if len(names) == 0 {
		names = append([]string(nil), DefaultPackages...)
	}

	fmt.Fprintf(opts.Stdout, "→ cn repo install")
	if opts.DryRun {
		fmt.Fprintf(opts.Stdout, " (dry-run) — no files will be written")
	}
	fmt.Fprintln(opts.Stdout)
	fmt.Fprintf(opts.Stdout, "✓ Git repository root: %s\n", opts.RepoRoot)

	idxPath, idx, releaseTag, cleanup, err := resolveIndex(ctx, client, apiBase, dlBase, repo, opts.Release, opts.IndexPath)
	if err != nil {
		fmt.Fprintf(opts.Stderr, "✗ %s\n", err)
		return nil, err
	}
	defer cleanup()

	if releaseTag != "" {
		fmt.Fprintf(opts.Stdout, "✓ Resolved cnos release: %s\n", releaseTag)
	} else {
		fmt.Fprintf(opts.Stdout, "✓ Using package index: %s\n", opts.IndexPath)
	}

	deps, err := resolvePins(idx, names, releaseTag)
	if err != nil {
		fmt.Fprintf(opts.Stderr, "✗ %s\n", err)
		return nil, err
	}

	manifest := pkg.Manifest{Schema: manifestSchema, Profile: manifestProfile, Packages: deps}
	plannedFiles := []string{
		filepath.Join(".cn", "deps.json"),
		filepath.Join(".cn", "deps.lock.json"),
		".gitignore",
	}

	if opts.DryRun {
		printPlan(opts.Stdout, manifest, plannedFiles, releaseTag, opts.IndexPath, opts.Dispatch)
		return &Result{ReleaseTag: releaseTag, Manifest: manifest, DryRun: true}, nil
	}

	if err := applyInstall(ctx, opts, manifest, idxPath); err != nil {
		fmt.Fprintf(opts.Stderr, "✗ %s\n", err)
		return nil, err
	}

	if opts.Dispatch == "cds" {
		// Base install (above) has already completed and written its
		// files (AC1/C1) — the dispatch render layers on top of it, not
		// instead of it. runDispatchCds prints its own identity/PAT-scope
		// stdout lines and, today, always ends by surfacing the cnos#493
		// canonical-label gap as a returned error (AC3) — it does not
		// silently report success while that obligation is unmet.
		if err := runDispatchCds(ctx, opts); err != nil {
			return nil, err
		}
		fmt.Fprintf(opts.Stdout, "Dispatch: cds (agent: %s)\n", resolveDispatchAgent(opts.Agent))
	} else {
		fmt.Fprintf(opts.Stdout, "Dispatch: none (base install only — no .github/workflows/ changes)\n")
	}
	fmt.Fprintf(opts.Stdout, "\n✓ cn repo install complete.\n")

	return &Result{ReleaseTag: releaseTag, Manifest: manifest, DryRun: false}, nil
}

// validateDispatch validates the --dispatch flag value itself (Mock B4).
// It does not gate "cds" on any precondition — cnos#610 wires "cds" to
// the (now-merged, #609) renderer; identity/label preconditions are
// enforced later in Run, per-precondition, with their own named errors
// (see runDispatchCds / ensureCanonicalDispatchLabels) rather than one
// blanket refusal here.
func validateDispatch(d string) error {
	switch d {
	case "", "none", "cds":
		return nil
	default:
		return fmt.Errorf("unknown --dispatch value %q (want \"none\" or \"cds\")", d)
	}
}

// resolveDispatchAgent returns the effective --agent value for
// --dispatch cds, defaulting to "sigma" to match cn-install-wake's own
// default-agent convention (src/packages/cnos.core/commands/install-wake/
// cn-install-wake: "Default agent: sigma ... the only agent with
// substrate bindings today"). This keeps a bare `--dispatch cds` (no
// --agent) working exactly as it always implicitly has, preserving AC5
// backward-compat.
func resolveDispatchAgent(a string) string {
	if a == "" {
		return "sigma"
	}
	return a
}

// dispatchWorkflowPath returns the fixed output path --dispatch cds
// always renders to (AC1/C3): .github/workflows/cnos-cds-dispatch.yml
// under the repo root, matching cnos-cds-dispatch.yml already committed
// on main as the live sigma-bound substrate artifact.
func dispatchWorkflowPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".github", "workflows", "cnos-cds-dispatch.yml")
}

// runDispatchCds renders the cds dispatch workflow after the base
// install has completed successfully. Sequence (per
// .cdd/unreleased/610/gamma-scaffold.md "Surfaces α is expected to
// touch" §1):
//
//  1. Resolve identity; a non-sigma --agent with no --workflow-pat-secret
//     fails HERE, before the renderer is ever invoked — so no partial
//     .github/workflows/cnos-cds-dispatch.yml can exist on this path
//     (AC2/C2, Mock C2 "no partial render").
//  2. Invoke the vendored cn-install-wake renderer (cnos#609) against
//     the cds-dispatch wake manifest, targeting the fixed output path.
//  3. State the workflow-scope PAT requirement + "never pushes to main"
//     fact in stdout (AC4/C6 — this package never calls git/gh; the
//     never-pushes-main property holds by construction, not by a guard).
//  4. Ensure the canonical dispatch labels via the cnos#493
//     label-doctor mechanism (AC3): an in-process call into
//     packages/cnos.core/commands/label-doctor (both packages are
//     go.work-linked, mirroring cmd_issues_fsm.go's cross-module
//     issuesfsm import). Does not silently skip the obligation: if the
//     installing repo's target (owner/repo, resolved from its git
//     "origin" remote) or a GitHub token cannot be resolved, or the
//     GitHub API call itself fails, this surfaces a named, actionable
//     error once identity resolution + the render themselves succeed.
func runDispatchCds(ctx context.Context, opts Options) error {
	agent := resolveDispatchAgent(opts.Agent)

	patSecret := opts.WorkflowPatSecret
	if patSecret == "" {
		if agent != "sigma" {
			err := fmt.Errorf("--workflow-pat-secret is required for --agent %q (no default substrate PAT-secret binding for non-sigma agents); pass --workflow-pat-secret <NAME> naming the GitHub Actions secret that holds this agent's workflow-scoped PAT", agent)
			fmt.Fprintf(opts.Stderr, "✗ %s\n", err)
			return err
		}
		// sigma's default substrate PAT binding, mirroring
		// cn-install-wake's own default (renderer authority; this is
		// just the display value printed below).
		patSecret = "SIGMA_WORKFLOW_PAT"
	}

	rendererPath := filepath.Join(pkg.VendorPath(opts.RepoRoot, "cnos.core"), "commands", "install-wake", "cn-install-wake")
	if _, statErr := os.Stat(rendererPath); statErr != nil {
		err := fmt.Errorf("dispatch renderer not found at %s (the cnos.core package must be installed for --dispatch cds): %w", rendererPath, statErr)
		fmt.Fprintf(opts.Stderr, "✗ %s\n", err)
		return err
	}

	outPath := dispatchWorkflowPath(opts.RepoRoot)

	args := []string{"cds-dispatch", "--out", outPath, "--agent", agent}
	if opts.WorkflowPatSecret != "" {
		args = append(args, "--workflow-pat-secret", opts.WorkflowPatSecret)
	}
	if opts.BotName != "" {
		args = append(args, "--bot-name", opts.BotName)
	}
	if opts.BotID != "" {
		args = append(args, "--bot-id", opts.BotID)
	}

	cmd := exec.CommandContext(ctx, rendererPath, args...)
	cmd.Dir = opts.RepoRoot
	cmd.Stdout = opts.Stdout
	cmd.Stderr = opts.Stderr
	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("render dispatch workflow: %w", err)
		fmt.Fprintf(opts.Stderr, "✗ %s\n", err)
		return err
	}

	fmt.Fprintf(opts.Stdout, "✓ rendered .github/workflows/cnos-cds-dispatch.yml\n")
	fmt.Fprintf(opts.Stdout, "  identity: %s\n", agent)
	fmt.Fprintf(opts.Stdout, "  pat secret: %s\n", patSecret)
	fmt.Fprintf(opts.Stdout, "⚠ dispatch grants scheduled write access on merge. Review before merging.\n")
	fmt.Fprintf(opts.Stdout, "  This changes .github/workflows/ — the installing token needs `workflow` scope.\n")
	fmt.Fprintf(opts.Stdout, "  Dispatch never pushes to main (PR-only).\n")

	if err := ensureCanonicalDispatchLabels(ctx, opts); err != nil {
		fmt.Fprintf(opts.Stderr, "✗ %s\n", err)
		return err
	}
	return nil
}

// ensureCanonicalDispatchLabels ensures every canonical label in
// src/packages/cnos.core/labels.json (the 7 status:* lifecycle labels +
// dispatch:cell) exists on the installing repo (opts.RepoRoot's git
// "origin" remote), with canonical color/description, via the cnos#493
// label-doctor mechanism (packages/cnos.core/commands/label-doctor),
// called in-process — not a subprocess/vendored-binary exec, unlike
// runDispatchCds's cn-install-wake invocation above: label-doctor is
// pure API+diff logic with no templating/file-artifact constraint
// requiring a separate process (see the γ scaffold's Implementation
// contract → Existing-binary disposition row).
//
// labeldoctor.Doctor resolves its own target repo (from opts.RepoRoot's
// git remote) and its own GitHub token ($GITHUB_TOKEN then $GH_TOKEN);
// this function passes RepoRoot/Stdout/Stderr through and otherwise
// leaves resolution to that package, so a caller with no configured git
// remote (as every existing repoinstall_test.go dispatch-cds fixture
// uses — a plain t.TempDir() or a git-init'd repo with no "origin") gets
// a named, actionable error rather than any live network call.
func ensureCanonicalDispatchLabels(ctx context.Context, opts Options) error {
	res, err := labeldoctor.Doctor(ctx, labeldoctor.Options{
		RepoRoot: opts.RepoRoot,
		Stdout:   opts.Stdout,
		Stderr:   opts.Stderr,
	})
	if err != nil {
		return fmt.Errorf("canonical dispatch labels not ensured: %w", err)
	}
	if len(res.Applied) > 0 {
		fmt.Fprintf(opts.Stdout, "✓ label-doctor repaired %d canonical label(s): %s\n", len(res.Applied), strings.Join(res.Applied, ", "))
	} else {
		fmt.Fprintf(opts.Stdout, "✓ label-doctor: all canonical labels already present and matching\n")
	}
	return nil
}

// applyInstall performs the non-dry-run write path: .cn/deps.json,
// .cn/deps.lock.json (via restore.GenerateLockFromIndex), package restore
// (via restore.Restore), and the .gitignore entry (via
// hubsetup.EnsureGitignoreEntry — reused, not duplicated).
func applyInstall(ctx context.Context, opts Options, manifest pkg.Manifest, idxPath string) error {
	cnDir := filepath.Join(opts.RepoRoot, ".cn")
	if err := os.MkdirAll(cnDir, 0755); err != nil {
		return fmt.Errorf("create .cn: %w", err)
	}

	manifestPath := filepath.Join(cnDir, "deps.json")
	if err := writeManifest(manifestPath, manifest); err != nil {
		return err
	}
	fmt.Fprintf(opts.Stdout, "✓ wrote .cn/deps.json\n")

	lockResult, err := restore.GenerateLockFromIndex(opts.RepoRoot, idxPath)
	if err != nil {
		return fmt.Errorf("deps lock: %w", err)
	}
	fmt.Fprintf(opts.Stdout, "✓ wrote .cn/deps.lock.json (%d package(s))\n", lockResult.Count)

	installed, err := restore.Restore(ctx, opts.RepoRoot, idxPath)
	if err != nil {
		return fmt.Errorf("deps restore: %w", err)
	}
	if restore.HasErrors(installed) {
		for _, r := range restore.Errors(installed) {
			fmt.Fprintf(opts.Stderr, "✗ %s@%s: %v\n", r.Name, r.Version, r.Err)
		}
		return fmt.Errorf("deps restore: %d package(s) failed", len(restore.Errors(installed)))
	}
	for _, r := range installed {
		fmt.Fprintf(opts.Stdout, "✓ restored %s@%s\n", r.Name, r.Version)
	}

	if err := hubsetup.EnsureGitignoreEntry(opts.RepoRoot, opts.Stdout); err != nil {
		return fmt.Errorf("gitignore: %w", err)
	}

	return nil
}

// writeManifest writes .cn/deps.json deterministically: stable field
// order (json.MarshalIndent over a struct, not a map), the exact package
// order Run resolved (operator-controlled, per restore.go's own
// determinism comment), and no timestamp field.
func writeManifest(path string, m pkg.Manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal deps.json: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write deps.json: %w", err)
	}
	return nil
}

// printPlan renders the --dry-run report (Mock A). It writes nothing to
// disk — the caller (Run) returns immediately after this call. The
// "Dispatch: none" line for dispatch == "" / "none" is byte-identical to
// the pre-cnos#610 text (pinned backward-compat); dispatch == "cds"
// prints a distinct line so --dry-run --dispatch cds does not falsely
// claim no .github/workflows/ change is planned.
func printPlan(w io.Writer, m pkg.Manifest, plannedFiles []string, releaseTag, indexArg, dispatch string) {
	if releaseTag != "" {
		fmt.Fprintf(w, "  Would fetch package index + tarballs for release %s (%d package(s))\n", releaseTag, len(m.Packages))
	} else {
		fmt.Fprintf(w, "  Would fetch package index from %s (%d package(s))\n", indexArg, len(m.Packages))
	}
	fmt.Fprintf(w, "  Would write .cn/deps.json:\n")
	for _, d := range m.Packages {
		fmt.Fprintf(w, "      %-12s %s\n", d.Name, d.Version)
	}
	fmt.Fprintf(w, "  Would run: (internal) deps lock    → .cn/deps.lock.json\n")
	fmt.Fprintf(w, "  Would run: (internal) deps restore → .cn/vendor/packages/ (%d package(s))\n", len(m.Packages))
	fmt.Fprintf(w, "  Would ensure .gitignore contains: .cn/vendor/\n")
	if dispatch == "cds" {
		fmt.Fprintf(w, "  Would render: .github/workflows/cnos-cds-dispatch.yml (via cn-install-wake)\n")
		fmt.Fprintf(w, "Dispatch: cds (would render .github/workflows/cnos-cds-dispatch.yml)\n\n")
	} else {
		fmt.Fprintf(w, "Dispatch: none (base install only — no .github/workflows/ changes)\n\n")
	}
	fmt.Fprintf(w, "Planned committed diff (%d files):\n", len(plannedFiles))
	for _, f := range plannedFiles {
		fmt.Fprintf(w, "  %s\n", f)
	}
	fmt.Fprintf(w, "\nRun without --dry-run to apply.\n")
}

// --- Index resolution ---

// resolveIndex determines the package index to install from. It returns
// a local file path that restore.GenerateLockFromIndex/restore.Restore can
// read directly, the parsed index (for pin resolution), the resolved
// release tag ("" when an explicit --index value was used with no
// --release override), and a cleanup func removing any temp files this
// call created (always non-nil; safe to defer unconditionally).
func resolveIndex(ctx context.Context, client *http.Client, apiBase, dlBase, repo, release, indexArg string) (indexPath string, idx *pkg.PackageIndex, releaseTag string, cleanup func(), err error) {
	noop := func() {}

	pinFromRelease := ""
	if release != "" && release != "latest" {
		pinFromRelease = release
	}

	if indexArg != "" {
		if isRemoteURL(indexArg) {
			body, ferr := fetchBytes(ctx, client, indexArg)
			if ferr != nil {
				return "", nil, "", noop, fmt.Errorf("fetch package index %s: %w", indexArg, ferr)
			}
			parsed, perr := pkg.ParsePackageIndex(body)
			if perr != nil {
				return "", nil, "", noop, fmt.Errorf("parse package index %s: %w", indexArg, perr)
			}
			rewritten := rewriteRelativeEntries(parsed, indexArg)
			tmpPath, tmpCleanup, werr := writeTempIndex(rewritten)
			if werr != nil {
				return "", nil, "", noop, werr
			}
			return tmpPath, rewritten, pinFromRelease, tmpCleanup, nil
		}

		data, rerr := os.ReadFile(indexArg)
		if rerr != nil {
			return "", nil, "", noop, fmt.Errorf("read package index %s: %w", indexArg, rerr)
		}
		parsed, perr := pkg.ParsePackageIndex(data)
		if perr != nil {
			return "", nil, "", noop, fmt.Errorf("parse package index %s: %w", indexArg, perr)
		}
		// Local path: used directly, exactly like the existing
		// cn build → dist/ → cn deps restore dev workflow — relative
		// tarball URLs are resolved by restore.go against indexPath's
		// own directory, so no rewriting is needed here.
		return indexArg, parsed, pinFromRelease, noop, nil
	}

	// No explicit --index: resolve a release and fetch its index.
	tag := release
	if tag == "" || tag == "latest" {
		rel, rerr := binupdate.FetchLatestRelease(ctx, client, apiBase, repo)
		if rerr != nil {
			return "", nil, "", noop, fmt.Errorf("resolve latest cnos release: %w", rerr)
		}
		tag = rel.Tag
	}

	base := fmt.Sprintf("%s/%s/releases/download/%s/", dlBase, repo, tag)
	indexURL := base + "index.json"
	body, ferr := fetchBytes(ctx, client, indexURL)
	if ferr != nil {
		return "", nil, "", noop, fmt.Errorf("fetch package index for release %s: %w", tag, ferr)
	}
	parsed, perr := pkg.ParsePackageIndex(body)
	if perr != nil {
		return "", nil, "", noop, fmt.Errorf("parse package index for release %s: %w", tag, perr)
	}

	rewritten := rewriteRelativeEntriesFromBase(parsed, base)
	tmpPath, tmpCleanup, werr := writeTempIndex(rewritten)
	if werr != nil {
		return "", nil, "", noop, werr
	}
	return tmpPath, rewritten, tag, tmpCleanup, nil
}

// isRemoteURL returns true if the URL has an http:// or https:// scheme.
// Mirrors restore.go's unexported helper of the same name/behavior — kept
// local rather than exported cross-package because it is a one-line
// string check, not a shared parser (eng/go §2.17 governs parser
// duplication, not trivial scheme checks).
func isRemoteURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// fetchBytes GETs url and returns the response body. Used for both
// index.json fetches (release + explicit --index URL forms).
func fetchBytes(ctx context.Context, client *http.Client, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// rewriteRelativeEntries rewrites relative tarball URLs in idx against
// the directory of indexURL (an http(s) URL the index itself was fetched
// from). Absolute (http/https) entries are left untouched.
func rewriteRelativeEntries(idx *pkg.PackageIndex, indexURL string) *pkg.PackageIndex {
	u, err := url.Parse(indexURL)
	if err != nil {
		// indexURL was already successfully fetched via this exact string,
		// so a parse failure here would be unreachable in practice; fall
		// back to leaving entries untouched rather than panicking.
		return idx
	}
	dir := *u
	dir.Path = path.Dir(u.Path) + "/"
	dir.RawQuery = ""
	dir.Fragment = ""
	return rewriteRelativeEntriesFromBase(idx, dir.String())
}

// rewriteRelativeEntriesFromBase returns a copy of idx with every
// relative (non-http, non-absolute-path) tarball URL prefixed with base.
// Implementation requirement 4 (issue #608): normalize the index into a
// form restore.Restore's existing fetchTarball can consume — after this
// rewrite, every entry is an absolute http(s) URL, so fetchTarball's
// isRemoteURL branch always fires regardless of indexDir.
func rewriteRelativeEntriesFromBase(idx *pkg.PackageIndex, base string) *pkg.PackageIndex {
	out := &pkg.PackageIndex{Schema: idx.Schema, Packages: make(map[string]map[string]pkg.IndexEntry, len(idx.Packages))}
	for name, versions := range idx.Packages {
		outVersions := make(map[string]pkg.IndexEntry, len(versions))
		for v, e := range versions {
			if !isRemoteURL(e.URL) && !filepath.IsAbs(e.URL) {
				e.URL = base + e.URL
			}
			outVersions[v] = e
		}
		out.Packages[name] = outVersions
	}
	return out
}

// writeTempIndex marshals idx to a temp file outside the target repo's
// working tree (so it never appears in `git status`) and returns its
// path plus a cleanup func.
func writeTempIndex(idx *pkg.PackageIndex) (string, func(), error) {
	dir, err := os.MkdirTemp("", "cn-repo-install-index-*")
	if err != nil {
		return "", func() {}, fmt.Errorf("create temp index dir: %w", err)
	}
	cleanup := func() { os.RemoveAll(dir) }

	data, merr := json.MarshalIndent(idx, "", "  ")
	if merr != nil {
		cleanup()
		return "", func() {}, fmt.Errorf("marshal package index: %w", merr)
	}
	p := filepath.Join(dir, "index.json")
	if werr := os.WriteFile(p, data, 0644); werr != nil {
		cleanup()
		return "", func() {}, fmt.Errorf("write temp package index: %w", werr)
	}
	return p, cleanup, nil
}

// --- Pin resolution ---

// resolvePins resolves an exact version for each requested package name
// against idx. When pinVersion is set (a resolved release tag or an
// explicit --release value), every package must be present at exactly
// that version. Otherwise (index-only flow, no --release), a package
// with exactly one version in the index resolves unambiguously to that
// version; zero or multiple versions is an error naming the package.
func resolvePins(idx *pkg.PackageIndex, names []string, pinVersion string) ([]pkg.ManifestDep, error) {
	deps := make([]pkg.ManifestDep, 0, len(names))
	var missing []string
	var ambiguous []string

	for _, name := range names {
		versions := idx.Packages[name]
		switch {
		case pinVersion != "":
			if _, ok := versions[pinVersion]; !ok {
				missing = append(missing, fmt.Sprintf("%s@%s", name, pinVersion))
				continue
			}
			deps = append(deps, pkg.ManifestDep{Name: name, Version: pinVersion})
		case len(versions) == 1:
			for v := range versions {
				deps = append(deps, pkg.ManifestDep{Name: name, Version: v})
			}
		case len(versions) == 0:
			missing = append(missing, name)
		default:
			ambiguous = append(ambiguous, name)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("package(s) not found in index: %s", strings.Join(missing, ", "))
	}
	if len(ambiguous) > 0 {
		return nil, fmt.Errorf("package(s) have multiple versions in index; pass --release to pin one: %s", strings.Join(ambiguous, ", "))
	}
	return deps, nil
}
