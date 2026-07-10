// Package repostatus implements `cn repo status` (cnos#656, Phase 1 of the
// cnos#655 `cn repo` lifecycle wave): a read-only report of a repo's
// installed CNOS/CDS state versus its desired/locked/ledger state.
//
// Discipline: this package NEVER writes to opts.RepoRoot (design doc P3 —
// "cn repo status is read-only by default"). It may fetch a release tag
// over the network for the update-available check, and it may invoke the
// vendored cn-install-wake renderer into a throwaway temp file to compute
// a fresh-render comparison sha (A2 drift classification) — neither of
// those touches opts.RepoRoot itself.
//
// This is scoped distinctly from the existing global `cn status`
// (internal/hubstatus — binary/hub health, operates on the agent hub at
// inv.HubPath) per the design doc's A5: `cn repo status` answers "is THIS
// REPO's CNOS install in sync?", resolved via the git repository root,
// never via inv.HubPath.
package repostatus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	labeldoctor "github.com/usurobor/cnos/packages/cnos.core/commands/label-doctor"
	"github.com/usurobor/cnos/src/go/internal/binupdate"
	"github.com/usurobor/cnos/src/go/internal/dispatchrender"
	"github.com/usurobor/cnos/src/go/internal/pkg"
	"github.com/usurobor/cnos/src/go/internal/repostate"
	"github.com/usurobor/cnos/src/go/internal/restore"
)

// StatusSchema is the fixed schema identifier `--json` output declares.
const StatusSchema = "cn.repo.status.v1"

// DefaultRepo mirrors repoinstall.DefaultRepo (kept as an independent
// literal, not an import, to avoid a status->install dependency for a
// single stable constant).
const DefaultRepo = "usurobor/cnos"

// Drift classification values (design doc A2). DriftRemoved is this
// package's own addition beyond A2's three-way split: the ledger records
// a workflow that no longer exists on disk at all (distinct from "never
// had one," which reports Present: false with no Drift value).
const (
	DriftMatchesLedger = "matches_ledger"
	DriftUserEdit      = "user_edit"
	DriftRendererMoved = "renderer_moved"
	DriftUnknown       = "unknown"
	DriftRemoved       = "removed"
)

// Label status values.
const (
	LabelsOK      = "ok"
	LabelsDrifted = "drifted"
	LabelsUnknown = "unknown"
)

type PackageStatus struct {
	Name     string `json:"name"`
	Desired  string `json:"desired,omitempty"`
	Locked   string `json:"locked,omitempty"`
	Vendored string `json:"vendored,omitempty"`
	InSync   bool   `json:"in_sync"`
}

type DispatchStatus struct {
	Present bool   `json:"present"`
	Tier    string `json:"tier,omitempty"`
	ID      string `json:"id,omitempty"`
	Path    string `json:"path,omitempty"`
	// Drift is "" when Present is false; otherwise one of
	// matches_ledger / user_edit / renderer_moved / unknown.
	Drift string `json:"drift,omitempty"`
}

type LabelStatus struct {
	// Status is "ok" (no missing/drifted canonical labels), "drifted"
	// (at least one missing or color/description-drifted canonical
	// label), or "unknown" (repo target or GitHub token could not be
	// resolved — informational, never treated as drift by IsDrift).
	Status  string   `json:"status"`
	Missing []string `json:"missing,omitempty"`
	// Drifted names canonical labels present with a wrong color/
	// description. Exceeds the design doc's illustrative
	// {status,missing,unknown} shape (mock_parity row A4 — safe:
	// additional, non-conflicting information).
	Drifted []string `json:"drifted,omitempty"`
	// Unknown names labels present on the repo but not in the canonical
	// manifest — informational only, never contributes to Status or to
	// the top-level Drift verdict.
	Unknown []string `json:"unknown,omitempty"`
}

type LocalEdit struct {
	Path string `json:"path"`
	// Classification is "user_edit" (content differs from the ledger
	// sha256) or "removed" (the managed file no longer exists on disk).
	Classification string `json:"classification"`
}

type UpdateAvailable struct {
	// Checked reports whether the network resolution actually ran;
	// false means "could not determine" (offline, no token, etc.) —
	// distinct from Available=false, which means "checked, and current".
	Checked   bool   `json:"checked"`
	Available bool   `json:"available"`
	Release   string `json:"release,omitempty"`
}

type LedgerInfo struct {
	Present       bool `json:"present"`
	Reconstructed bool `json:"reconstructed"`
}

// Status is the full `cn repo status` report — the `--json` shape.
type Status struct {
	Schema          string           `json:"schema"`
	Source          repostate.Source `json:"source"`
	Packages        []PackageStatus  `json:"packages"`
	Dispatch        DispatchStatus   `json:"dispatch"`
	Labels          LabelStatus      `json:"labels"`
	OrphanPackages  []string         `json:"orphan_packages"`
	LocalEdits      []LocalEdit      `json:"local_edits"`
	UpdateAvailable UpdateAvailable  `json:"update_available"`
	Ledger          LedgerInfo       `json:"ledger"`
	Drift           bool             `json:"drift"`
}

// Options carries Run's configuration. Every field beyond RepoRoot is
// optional and defaults to production behavior (live GitHub API calls for
// labels/update-available) — tests override HTTPClient/APIBaseURL/
// DownloadBase/Repo/Token/SkipNetwork the same way repoinstall's test
// seams work.
type Options struct {
	RepoRoot string

	// Repo/Token override label-doctor's and the update-available
	// check's target repo/token resolution. Empty means "resolve from
	// RepoRoot's git origin remote" (labels) / DefaultRepo (update
	// check).
	Repo  string
	Token string

	HTTPClient   *http.Client
	APIBaseURL   string
	DownloadBase string

	// SkipNetwork disables the update-available check entirely (tests
	// that don't want any live/mock HTTP call for that informational
	// field). Label lookups are unaffected — they already degrade to
	// LabelsUnknown on any resolution failure.
	SkipNetwork bool
}

// Run produces a Status report for opts.RepoRoot. Never writes to
// opts.RepoRoot (P3).
func Run(ctx context.Context, opts Options) (*Status, error) {
	if opts.RepoRoot == "" {
		return nil, fmt.Errorf("repo status: RepoRoot is required")
	}

	manifestPath := filepath.Join(opts.RepoRoot, ".cn", "deps.json")
	manifest, err := restore.ReadManifest(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("repo status: %s not found or unreadable — this repo has no `cn repo install` state (%w)", manifestPath, err)
	}

	lockPath := filepath.Join(opts.RepoRoot, ".cn", "deps.lock.json")
	lockfile, lockErr := restore.ReadLockfile(lockPath)
	locked := map[string]string{}
	if lockErr == nil {
		for _, dep := range lockfile.Packages {
			locked[dep.Name] = dep.Version
		}
	}

	statePath := filepath.Join(opts.RepoRoot, ".cn", "repo.state.json")
	ledger, ledgerErr := readLedger(statePath)
	ledgerInfo := LedgerInfo{Present: ledgerErr == nil, Reconstructed: ledgerErr != nil}

	packages, orphans := packageStatuses(opts.RepoRoot, manifest, locked, lockErr == nil)

	dispatch := dispatchStatus(ctx, opts.RepoRoot, ledger)
	labels := labelStatus(ctx, opts)
	localEdits := localEditStatus(opts.RepoRoot, ledger)
	update := updateAvailableStatus(ctx, opts, ledger)

	drift := false
	for _, p := range packages {
		if !p.InSync {
			drift = true
		}
	}
	if dispatch.Drift == DriftUserEdit || dispatch.Drift == DriftRendererMoved || dispatch.Drift == DriftRemoved {
		drift = true
	}
	if labels.Status == LabelsDrifted {
		drift = true
	}
	if len(orphans) > 0 || len(localEdits) > 0 {
		drift = true
	}

	source := repostate.Source{}
	if ledger != nil {
		source = ledger.Source
	}

	return &Status{
		Schema:          StatusSchema,
		Source:          source,
		Packages:        packages,
		Dispatch:        dispatch,
		Labels:          labels,
		OrphanPackages:  orphans,
		LocalEdits:      localEdits,
		UpdateAvailable: update,
		Ledger:          ledgerInfo,
		Drift:           drift,
	}, nil
}

func readLedger(path string) (*repostate.RepoState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return repostate.Parse(data)
}

// packageStatuses reports desired/locked/vendored per manifest package,
// plus orphan vendored directories (materialized but absent from the
// lockfile — design doc's orphan-package invariant). Orphan detection
// only runs when lockPresent — with no lockfile at all there is nothing
// authoritative to call a vendored dir "extraneous" against.
func packageStatuses(repoRoot string, manifest *pkg.Manifest, locked map[string]string, lockPresent bool) ([]PackageStatus, []string) {
	statuses := make([]PackageStatus, 0, len(manifest.Packages))
	for _, dep := range manifest.Packages {
		vendored := ""
		if installed, err := restore.ReadInstalledManifest(pkg.VendorPath(repoRoot, dep.Name)); err == nil {
			vendored = installed.Version
		}
		lockedVersion := locked[dep.Name]
		inSync := dep.Version != "" && dep.Version == lockedVersion && lockedVersion == vendored
		statuses = append(statuses, PackageStatus{
			Name: dep.Name, Desired: dep.Version, Locked: lockedVersion, Vendored: vendored, InSync: inSync,
		})
	}

	var orphans []string
	if lockPresent {
		vendorRoot := filepath.Join(repoRoot, ".cn", "vendor", "packages")
		entries, err := os.ReadDir(vendorRoot)
		if err == nil {
			for _, e := range entries {
				if !e.IsDir() {
					continue
				}
				if _, ok := locked[e.Name()]; !ok {
					orphans = append(orphans, e.Name())
				}
			}
		}
	}
	sort.Strings(orphans)
	return statuses, orphans
}

func sha256File(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// dispatchStatus reports dispatch-workflow presence + A2 drift
// classification. When the ledger has no workflow record, it falls back
// to direct file presence (A3 degraded reporting — no ledger sha to
// anchor a three-way classification, so drift is DriftUnknown).
func dispatchStatus(ctx context.Context, repoRoot string, ledger *repostate.RepoState) DispatchStatus {
	const workflowRelPath = ".github/workflows/cnos-cds-dispatch.yml"

	var wf *repostate.ManagedFile
	if ledger != nil {
		wf = ledger.FindManagedFile(workflowRelPath)
	}

	livePath := filepath.Join(repoRoot, ".github", "workflows", "cnos-cds-dispatch.yml")
	liveSHA, liveErr := sha256File(livePath)

	if wf == nil {
		if liveErr != nil {
			return DispatchStatus{Present: false}
		}
		// A3: workflow exists but no ledger record to classify it against.
		return DispatchStatus{Present: true, ID: "cnos-cds-dispatch", Path: workflowRelPath, Drift: DriftUnknown}
	}

	if liveErr != nil {
		// The ledger records a workflow that no longer exists on disk —
		// distinct from "never had one" (wf == nil above): Present is
		// false (nothing is actually active), but Tier/ID/Path are still
		// surfaced for context, and Drift names what happened so this
		// doesn't get silently treated as "no dispatch, no drift".
		return DispatchStatus{Present: false, Tier: wf.Tier, ID: wf.ID, Path: wf.Path, Drift: DriftRemoved}
	}

	status := DispatchStatus{Present: true, Tier: wf.Tier, ID: wf.ID, Path: wf.Path}
	if liveSHA == wf.SHA256 {
		status.Drift = DriftMatchesLedger
		return status
	}

	// Differs from ledger — classify user_edit vs renderer_moved by
	// re-rendering fresh into a temp file with the SAME render contract
	// the ledger recorded (P2) and comparing.
	tmpDir, err := os.MkdirTemp("", "cn-repo-status-render-*")
	if err != nil {
		status.Drift = DriftUnknown
		return status
	}
	defer os.RemoveAll(tmpDir)
	tmpOut := filepath.Join(tmpDir, "workflow.yml")

	renderErr := dispatchrender.Render(ctx, dispatchrender.Options{
		RepoRoot:          repoRoot,
		RendererPackage:   wf.RendererPackage,
		Tier:              wf.Tier,
		Agent:             wf.Agent,
		WorkflowPatSecret: wf.WorkflowPatSecret,
		BotName:           wf.BotName,
		BotID:             wf.BotID,
		OutPath:           tmpOut,
	})
	if renderErr != nil {
		status.Drift = DriftUnknown
		return status
	}
	freshSHA, err := sha256File(tmpOut)
	if err != nil {
		status.Drift = DriftUnknown
		return status
	}
	if freshSHA == liveSHA {
		status.Drift = DriftRendererMoved
	} else {
		status.Drift = DriftUserEdit
	}
	return status
}

// localEditStatus reports non-workflow managed_files whose live content
// differs from the ledger sha256 (workflow drift is reported via
// DispatchStatus.Drift instead — not duplicated here, matching the design
// mocks doc's Mock C example).
func localEditStatus(repoRoot string, ledger *repostate.RepoState) []LocalEdit {
	if ledger == nil {
		return nil
	}
	var edits []LocalEdit
	for _, mf := range ledger.ManagedFiles {
		if mf.Kind == repostate.KindWorkflow {
			continue
		}
		sha, err := sha256File(filepath.Join(repoRoot, mf.Path))
		if err != nil {
			edits = append(edits, LocalEdit{Path: mf.Path, Classification: "removed"})
			continue
		}
		if sha != mf.SHA256 {
			edits = append(edits, LocalEdit{Path: mf.Path, Classification: "user_edit"})
		}
	}
	sort.Slice(edits, func(i, j int) bool { return edits[i].Path < edits[j].Path })
	return edits
}

// labelStatus audits canonical labels read-only (labeldoctor.Doctor with
// DryRun: true never mutates GitHub). Degrades to LabelsUnknown (not an
// error) when the repo target or a GitHub token cannot be resolved —
// status must never fail outright just because labels can't be checked.
func labelStatus(ctx context.Context, opts Options) LabelStatus {
	res, err := labeldoctor.Doctor(ctx, labeldoctor.Options{
		RepoRoot: opts.RepoRoot,
		Repo:     opts.Repo,
		Token:    opts.Token,
		DryRun:   true,
		Stdout:   io.Discard,
		Stderr:   io.Discard,
	})
	if err != nil && err != labeldoctor.ErrDrift {
		return LabelStatus{Status: LabelsUnknown}
	}
	if res == nil {
		return LabelStatus{Status: LabelsUnknown}
	}

	var missing, drifted []string
	canonical := map[string]bool{}
	for _, f := range res.Findings {
		canonical[f.Name] = true
		switch f.Status {
		case labeldoctor.StatusMissing:
			missing = append(missing, f.Name)
		case labeldoctor.StatusDrifted:
			drifted = append(drifted, f.Name)
		}
	}
	var unknown []string
	for _, name := range res.LiveLabels {
		if !canonical[name] {
			unknown = append(unknown, name)
		}
	}

	status := LabelsOK
	if len(missing) > 0 || len(drifted) > 0 {
		status = LabelsDrifted
	}
	return LabelStatus{Status: status, Missing: missing, Drifted: drifted, Unknown: unknown}
}

// updateAvailableStatus best-effort checks whether a newer release than
// the ledger's recorded source.release exists. Degrades to Checked: false
// on any resolution failure (offline, rate-limited, no ledger/no release
// recorded) — this is informational only and never contributes to Drift.
func updateAvailableStatus(ctx context.Context, opts Options, ledger *repostate.RepoState) UpdateAvailable {
	if opts.SkipNetwork || ledger == nil || ledger.Source.Release == "" {
		return UpdateAvailable{Checked: false}
	}
	repo := opts.Repo
	if repo == "" {
		repo = DefaultRepo
	}
	apiBase := opts.APIBaseURL
	if apiBase == "" {
		apiBase = "https://api.github.com"
	}
	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	rel, err := binupdate.FetchLatestRelease(ctx, client, apiBase, repo)
	if err != nil {
		return UpdateAvailable{Checked: false}
	}
	return UpdateAvailable{
		Checked:   true,
		Available: rel.Tag != ledger.Source.Release,
		Release:   rel.Tag,
	}
}
