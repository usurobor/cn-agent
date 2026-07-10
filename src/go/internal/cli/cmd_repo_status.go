package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/repostatus"
)

const repoStatusHelp = `cn repo status - Report this repo's CNOS/CDS install state and drift

USAGE:
  cn repo status [flags]

DESCRIPTION:
  Reports source channel/release, per-package desired-vs-locked-vs-vendored
  versions, the dispatch workflow's presence/tier/drift, canonical label
  status, orphan vendored packages, locally-edited managed files, and
  whether a newer release is available (cnos#656, Phase 1 of the cnos#655
  'cn repo' lifecycle wave).

  This is repo-scoped: "is THIS repo's CNOS install in sync?" — distinct
  from the existing global 'cn status' (binary/hub health at the current
  agent hub). Run 'cn status' for hub/binary health; run 'cn repo status'
  for repo-install drift.

  Read-only: never writes .cn/repo.state.json or any other file, with or
  without any flag below. When the ledger is absent (pre-cnos#656 installs),
  reports as much as it can from deps.json/deps.lock.json/vendor/the live
  workflow file directly, with dispatch drift degraded to "unknown" (no
  ledger sha to classify against).

FLAGS:
  --json    Emit the full report as machine-readable JSON (schema
            cn.repo.status.v1) instead of the human-readable summary.
  --check   Exit 1 if any drift was detected (packages out of sync, dispatch
            workflow drifted, canonical labels drifted, orphan vendored
            packages, or locally-edited managed files); exit 0 otherwise.
            Without --check, exit code is always 0 regardless of drift —
            drift is reported, not fatal.

EXAMPLES:
  cn repo status
  cn repo status --json
  cn repo status --check

EXIT CODES:
  0  Success (report printed; with --check, no drift detected)
  1  Error (no .cn/deps.json — this repo has no 'cn repo install' state),
     or with --check, drift was detected
`

// RepoStatusCmd implements "cn repo status" (cnos#656). NeedsHub is false —
// this is repo-scoped (resolved via the git repository root), not agent-
// hub-scoped, mirroring RepoInstallCmd's own A5 rationale.
type RepoStatusCmd struct{}

func (c *RepoStatusCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:     "repo-status",
		Summary:  "Report this repo's CNOS/CDS install state and drift",
		Source:   SourceKernel,
		Tier:     TierKernel,
		NeedsHub: false,
	}
}

func (c *RepoStatusCmd) Help() string {
	return repoStatusHelp
}

func (c *RepoStatusCmd) Run(ctx context.Context, inv Invocation) error {
	jsonOut := false
	checkMode := false
	for _, a := range inv.Args {
		switch a {
		case "--help", "-h":
			fmt.Fprint(inv.Stdout, repoStatusHelp)
			return nil
		case "--json":
			jsonOut = true
		case "--check":
			checkMode = true
		default:
			err := fmt.Errorf("unknown flag: %s", a)
			fmt.Fprintf(inv.Stderr, "✗ cn repo status: %s\n\n", err)
			fmt.Fprint(inv.Stderr, repoStatusHelp)
			return err
		}
	}

	// Repo-root resolution: git root, never inv.HubPath — same rationale
	// as cmd_repo_install.go (A5's repo-scope-vs-hub-scope distinction).
	root, gerr := gitRepoRoot(ctx)
	if gerr != nil {
		fmt.Fprintf(inv.Stderr, "✗ cn repo status must be run inside a Git repository.\n")
		return fmt.Errorf("repo status: not inside a git repository: %w", gerr)
	}

	st, err := repostatus.Run(ctx, repostatus.Options{RepoRoot: root})
	if err != nil {
		fmt.Fprintf(inv.Stderr, "✗ %s\n", err)
		return err
	}

	if jsonOut {
		data, merr := json.MarshalIndent(st, "", "  ")
		if merr != nil {
			return fmt.Errorf("repo status: marshal json: %w", merr)
		}
		fmt.Fprintln(inv.Stdout, string(data))
	} else {
		renderStatus(inv.Stdout, root, st)
	}

	if checkMode && st.Drift {
		return fmt.Errorf("repo status: drift detected")
	}
	return nil
}

func renderStatus(w io.Writer, repoRoot string, st *repostatus.Status) {
	fmt.Fprintf(w, "→ cn repo status — %s\n", repoRoot)

	if st.Ledger.Present {
		if st.Source.Release != "" {
			fmt.Fprintf(w, "✓ source: %s @ %s (%s)\n", st.Source.Channel, st.Source.Release, st.Source.Index)
		} else {
			fmt.Fprintf(w, "✓ source: %s (%s)\n", st.Source.Channel, st.Source.Index)
		}
	} else {
		fmt.Fprintf(w, "○ source: unknown (no .cn/repo.state.json ledger)\n")
	}

	allInSync := true
	for _, p := range st.Packages {
		if !p.InSync {
			allInSync = false
		}
	}
	if allInSync {
		fmt.Fprintf(w, "✓ deps.json / deps.lock.json / vendor: in sync (%d package(s))\n", len(st.Packages))
		for _, p := range st.Packages {
			fmt.Fprintf(w, "    %-12s %s\n", p.Name, p.Desired)
		}
	} else {
		fmt.Fprintf(w, "⚠ deps.json / deps.lock.json / vendor: out of sync\n")
		for _, p := range st.Packages {
			if !p.InSync {
				fmt.Fprintf(w, "    %s: desired=%s locked=%s vendored=%s\n", p.Name, p.Desired, p.Locked, p.Vendored)
			}
		}
	}

	if st.Dispatch.Present {
		switch st.Dispatch.Drift {
		case repostatus.DriftMatchesLedger:
			fmt.Fprintf(w, "✓ dispatch: %s (%s tier) — %s matches ledger\n", st.Dispatch.ID, st.Dispatch.Tier, st.Dispatch.Path)
		case repostatus.DriftRendererMoved:
			fmt.Fprintf(w, "⚠ dispatch: %s — %s differs from ledger (renderer_moved: matches a fresh re-render; run `cn repo repair`)\n", st.Dispatch.ID, st.Dispatch.Path)
		case repostatus.DriftUserEdit:
			fmt.Fprintf(w, "✗ dispatch: %s — %s differs from ledger (user_edit)\n", st.Dispatch.ID, st.Dispatch.Path)
		case repostatus.DriftUnclassified:
			fmt.Fprintf(w, "✗ dispatch: %s — %s differs from ledger (could not classify further: re-render comparison failed)\n", st.Dispatch.ID, st.Dispatch.Path)
		default:
			fmt.Fprintf(w, "○ dispatch: %s present, drift unknown (no ledger record)\n", st.Dispatch.Path)
		}
	} else if st.Dispatch.Drift == repostatus.DriftRemoved {
		fmt.Fprintf(w, "⚠ dispatch: %s is recorded in the ledger but no longer exists on disk (removed)\n", st.Dispatch.Path)
	} else {
		fmt.Fprintf(w, "✓ dispatch: none\n")
	}

	switch st.Labels.Status {
	case repostatus.LabelsOK:
		fmt.Fprintf(w, "✓ labels: ok\n")
	case repostatus.LabelsDrifted:
		var parts []string
		if len(st.Labels.Missing) > 0 {
			parts = append(parts, "missing: "+strings.Join(st.Labels.Missing, ", "))
		}
		if len(st.Labels.Drifted) > 0 {
			parts = append(parts, "drifted: "+strings.Join(st.Labels.Drifted, ", "))
		}
		fmt.Fprintf(w, "✗ labels: drifted (%s)\n", strings.Join(parts, "; "))
	default:
		fmt.Fprintf(w, "○ labels: unknown (could not resolve repo target or GitHub token)\n")
	}
	if len(st.Labels.Unknown) > 0 {
		fmt.Fprintf(w, "  unknown (non-canonical, informational): %s\n", strings.Join(st.Labels.Unknown, ", "))
	}

	if len(st.OrphanPackages) > 0 {
		fmt.Fprintf(w, "⚠ orphan vendored packages: %s\n", strings.Join(st.OrphanPackages, ", "))
	} else {
		fmt.Fprintf(w, "✓ orphan vendored packages: none\n")
	}

	if len(st.LocalEdits) > 0 {
		fmt.Fprintf(w, "⚠ locally-edited managed files:\n")
		for _, e := range st.LocalEdits {
			fmt.Fprintf(w, "    %s (%s)\n", e.Path, e.Classification)
		}
	} else {
		fmt.Fprintf(w, "✓ locally-edited managed files: none\n")
	}

	if st.UpdateAvailable.Checked {
		if st.UpdateAvailable.Available {
			fmt.Fprintf(w, "  update available: yes — %s (run `cn repo update`)\n", st.UpdateAvailable.Release)
		} else {
			fmt.Fprintf(w, "  update available: no (%s is the latest resolved release)\n", st.Source.Release)
		}
	} else {
		fmt.Fprintf(w, "  update available: unknown (network check skipped or failed)\n")
	}

	fmt.Fprintln(w)
	if st.Drift {
		fmt.Fprintf(w, "Drift detected (see ⚠/✗ rows above).\n")
	} else {
		fmt.Fprintf(w, "No drift.\n")
	}
}
