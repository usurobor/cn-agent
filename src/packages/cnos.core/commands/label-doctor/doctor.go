// Package labeldoctor implements `cn label-doctor` (cnos#493): it audits
// a target GitHub repo's labels against src/packages/cnos.core/
// labels.json (schema cn.labels.v1) and, unless --dry-run, repairs any
// drift — creating missing labels and correcting color/description on
// drifted ones, idempotently.
//
// The audit is fully generic: every entry in the manifest is diffed
// against the live label set by name/color/description, with no
// per-label special-casing (in particular, status:review — the
// originally-reported drift — is not treated any differently than any
// other canonical label).
//
// Design authority: cnos#493, .cdd/unreleased/493/gamma-scaffold.md.
// GitHub REST primitives (github.go) follow the dependency-free
// net/http idiom already established by
// cnos.issues/commands/issues-fsm/fetch.go (no `gh` CLI shellout, no
// third-party GitHub client) — see that file's ghEnsureLabelExists
// (cnos#615) for the closest existing precedent, which this package
// does not import (separate Go module, unexported) but mirrors
// structurally, extended with the list (GET) and drift-repair (PATCH)
// primitives that precedent lacks.
package labeldoctor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Status classifies one canonical label's live-repo state relative to
// labels.json (AC1's audit classification: present-and-matching /
// present-but-drifted / missing).
type Status string

const (
	StatusMatch   Status = "match"
	StatusDrifted Status = "drifted"
	StatusMissing Status = "missing"
)

// Finding is one canonical label's audit result.
type Finding struct {
	Name                 string
	Status               Status
	CanonicalColor       string
	CanonicalDescription string
	LiveColor            string
	LiveDescription      string
}

// NeedsRepair reports whether Finding requires a create or update call.
func (f Finding) NeedsRepair() bool { return f.Status != StatusMatch }

// Audit diffs manifest's canonical labels against the live label set
// (keyed by name), generically: every manifest entry is checked against
// name/color/description with no per-label special-casing. Color
// comparison is case-insensitive (GitHub itself is case-insensitive on
// hex); description comparison is exact.
func Audit(manifest Manifest, live map[string]ghLabel) []Finding {
	findings := make([]Finding, 0, len(manifest.Labels))
	for _, l := range manifest.Labels {
		lv, ok := live[l.Name]
		f := Finding{
			Name:                 l.Name,
			CanonicalColor:       l.Color,
			CanonicalDescription: l.Description,
		}
		switch {
		case !ok:
			f.Status = StatusMissing
		case !strings.EqualFold(lv.Color, l.Color) || lv.Description != l.Description:
			f.Status = StatusDrifted
			f.LiveColor = lv.Color
			f.LiveDescription = lv.Description
		default:
			f.Status = StatusMatch
			f.LiveColor = lv.Color
			f.LiveDescription = lv.Description
		}
		findings = append(findings, f)
	}
	return findings
}

// ErrDrift is returned by Doctor when DryRun is true and at least one
// finding is not StatusMatch — the AC5 CI-guard sentinel: a --dry-run
// invocation must exit nonzero when drift exists against the manifest,
// and nil when the live repo is already fully canonical.
var ErrDrift = errors.New("label-doctor: drift detected against labels.json")

// Options is Doctor's structured configuration — the entry point
// in-process Go callers (repoinstall.go's ensureCanonicalDispatchLabels)
// use directly. Run (cli.go) is the flag-parsed CLI entry point and
// calls into Doctor after resolving flags.
type Options struct {
	// Repo is "owner/repo". If empty, resolved from RepoRoot's git
	// "origin" remote.
	Repo string
	// RepoRoot anchors Repo/LabelsPath resolution when they are empty.
	// Required in that case.
	RepoRoot string
	// Token is the GitHub token. If empty, resolved from $GITHUB_TOKEN
	// then $GH_TOKEN.
	Token string
	// LabelsPath is the canonical manifest path. If empty, resolved via
	// resolveDefaultManifestPath(RepoRoot).
	LabelsPath string
	// DryRun reports drift only; makes no GitHub API mutations.
	DryRun bool
	Stdout io.Writer
	Stderr io.Writer
}

// Result records Doctor's outcome.
type Result struct {
	Repo     string
	Manifest Manifest
	Findings []Finding
	// Applied holds the names of labels actually created/updated. Empty
	// when DryRun is true (report-only) or when no drift existed.
	Applied []string
	// Failed holds the names of labels whose repair was attempted but
	// failed (e.g. GitHub rejecting a manifest value, as opposed to a
	// network/auth failure that aborts Doctor entirely before any
	// attempt). Always empty when DryRun is true. A non-empty Failed
	// does not stop other findings from being attempted — see Doctor's
	// apply-loop doc comment.
	Failed []string
	// LiveLabels is every label name present on Repo at audit time
	// (sorted), independent of manifest membership. Findings/Audit only
	// reports canonical (manifest) labels — a caller that additionally
	// needs "present but not canonical" names (e.g. cnos#656's
	// `cn repo status`, which reports unknown/extra labels as
	// informational) computes that by set-differencing LiveLabels
	// against Manifest.Labels rather than Doctor re-listing them.
	LiveLabels []string
}

// Doctor audits opts.Repo's live labels against opts.LabelsPath (or
// their resolved defaults) and, unless DryRun, repairs every drifted or
// missing entry — idempotently: a second call against unchanged live
// state finds StatusMatch everywhere and performs zero mutating API
// calls (AC4's idempotence oracle).
func Doctor(ctx context.Context, opts Options) (*Result, error) {
	repo := opts.Repo
	if repo == "" {
		r, err := resolveRepoFromGitRemote(ctx, opts.RepoRoot)
		if err != nil {
			return nil, err
		}
		repo = r
	}

	token := opts.Token
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		token = os.Getenv("GH_TOKEN")
	}

	labelsPath := opts.LabelsPath
	if labelsPath == "" {
		p, err := resolveDefaultManifestPath(opts.RepoRoot)
		if err != nil {
			return nil, err
		}
		labelsPath = p
	}

	manifest, err := LoadManifest(labelsPath)
	if err != nil {
		return nil, err
	}

	liveList, err := ghListLabels(ctx, repo, token)
	if err != nil {
		return nil, fmt.Errorf("label-doctor: %w", err)
	}
	live := make(map[string]ghLabel, len(liveList))
	for _, l := range liveList {
		live[l.Name] = l
	}

	findings := Audit(manifest, live)
	liveNames := make([]string, 0, len(liveList))
	for _, l := range liveList {
		liveNames = append(liveNames, l.Name)
	}
	sort.Strings(liveNames)
	res := &Result{Repo: repo, Manifest: manifest, Findings: findings, LiveLabels: liveNames}

	renderReport(stdoutOrDiscard(opts.Stdout), repo, findings, opts.DryRun)

	if opts.DryRun {
		if hasDrift(findings) {
			return res, ErrDrift
		}
		return res, nil
	}

	// One label's repair failing (e.g. GitHub rejecting a manifest value
	// outright — empirically hit live: labels.json's dispatch:cell
	// description is 149 bytes, over GitHub's 100-character label-
	// description cap, so PATCHing it 422s) must not block every OTHER
	// label's repair. Every finding gets its own attempt; errors are
	// collected and joined into a single returned error at the end, so
	// Doctor still tells the caller exactly what succeeded (res.Applied)
	// and what didn't (res.Failed) rather than aborting mid-manifest.
	var errs []error
	for _, f := range findings {
		if !f.NeedsRepair() {
			continue
		}
		var err error
		switch f.Status {
		case StatusMissing:
			err = ghCreateLabel(ctx, repo, token, ghLabel{Name: f.Name, Color: f.CanonicalColor, Description: f.CanonicalDescription})
			if err != nil {
				err = fmt.Errorf("label-doctor: create %q: %w", f.Name, err)
			}
		case StatusDrifted:
			err = ghUpdateLabel(ctx, repo, token, f.Name, f.CanonicalColor, f.CanonicalDescription)
			if err != nil {
				err = fmt.Errorf("label-doctor: repair %q: %w", f.Name, err)
			}
		}
		if err != nil {
			errs = append(errs, err)
			res.Failed = append(res.Failed, f.Name)
			fmt.Fprintf(stdoutOrDiscard(opts.Stdout), "  ✗ %s: %s\n", f.Name, err)
			continue
		}
		res.Applied = append(res.Applied, f.Name)
	}
	if len(errs) > 0 {
		return res, errors.Join(errs...)
	}
	return res, nil
}

func hasDrift(findings []Finding) bool {
	for _, f := range findings {
		if f.NeedsRepair() {
			return true
		}
	}
	return false
}

func stdoutOrDiscard(w io.Writer) io.Writer {
	if w == nil {
		return io.Discard
	}
	return w
}

// renderReport prints the AC1 audit artifact shape: every canonical
// label's classification (match / drifted / missing), plus canonical
// vs. live color/description for anything not matching. This output IS
// the AC1 "audit complete" artifact when the caller is `cn label-doctor`
// run interactively or in CI (the γ scaffold's oracle: "α's own tool
// output ... is the actual artifact this AC gates").
func renderReport(w io.Writer, repo string, findings []Finding, dryRun bool) {
	fmt.Fprintf(w, "label-doctor: auditing %d canonical label(s) against %s\n", len(findings), repo)
	for _, f := range findings {
		switch f.Status {
		case StatusMatch:
			fmt.Fprintf(w, "  match    %-20s color=%s\n", f.Name, f.CanonicalColor)
		case StatusDrifted:
			fmt.Fprintf(w, "  drifted  %-20s color: %s -> %s; description: %q -> %q\n", f.Name, f.LiveColor, f.CanonicalColor, f.LiveDescription, f.CanonicalDescription)
		case StatusMissing:
			fmt.Fprintf(w, "  missing  %-20s color=%s (will be created)\n", f.Name, f.CanonicalColor)
		}
	}
	if dryRun {
		fmt.Fprintln(w, "(dry-run: no mutations performed)")
	}
}
