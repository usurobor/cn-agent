package cli

import (
	"context"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/repoinstall"
)

const repoInstallHelp = `cn repo install - Install the base CNOS/CDS package set into this repository

USAGE:
  cn repo install [flags]

DESCRIPTION:
  Makes an arbitrary Git repository CDS-ready with one command (cnos#608).
  Resolves a cnos release (or an explicit package index), writes a
  deterministic .cn/deps.json + .cn/deps.lock.json, restores the default
  package set (cnos.core, cnos.cdd, cnos.cds) under .cn/vendor/packages/,
  and ensures .gitignore excludes the vendor tree.

  This is a kernel (built-in) command — it runs before any hub/package
  state exists, and before any .cn/ directory is present.

  Base install never writes .github/workflows/ and never requires a
  workflow/agent PAT. Autonomous dispatch install (--dispatch cds) is a
  separate, explicit opt-in gated on renderer generalization (#609); until
  that lands, --dispatch cds fails explicitly and writes nothing.

FLAGS:
  --release V     "latest" (default) or a pinned release tag (e.g. 3.82.0)
  --index P       Package index path or http(s) URL — overrides --release
                   (local/offline/test workflows; skips release resolution)
  --packages CSV  Comma-separated package names (default: cnos.core,cnos.cdd,cnos.cds)
  --dispatch D    "none" (default) — base install only.
                  "cds" — NOT implemented until #609; fails explicitly.
  --dry-run       Report the intended writes/fetches/restores; write nothing

EXAMPLES:
  cn repo install
  cn repo install --release latest --packages cnos.core,cnos.cdd,cnos.cds --dispatch none
  cn repo install --dry-run
  cn repo install --index ./dist/packages/index.json --packages cnos.core,cnos.cdd,cnos.cds

EXIT CODES:
  0  Success (installed, or --dry-run reported the plan)
  1  Error (not a git repo, package/index resolution failure, restore
     failure, or --dispatch cds)

INVARIANTS:
  - .cn/deps.json / .cn/deps.lock.json: deterministic (stable order, no
    timestamps), byte-identical across repeated runs with the same inputs.
  - .cn/vendor/packages/<name>/: name-based path (not <name>@<version>/).
  - Idempotent: a second run produces no further git diff.
  - No agent-hub scaffold is written (no spec/SOUL.md, agent/, threads/,
    state/) — contrast with 'cn init'.
`

// RepoInstallCmd implements "cn repo install" (cnos#608) — the base
// CNOS/CDS package installer. Resolves via the noun-verb dispatcher
// (cli.ResolveCommand builds the "repo"+"-"+"install" = "repo-install"
// lookup key), the same mechanism already used by "cn cell finalize" /
// "cn issues fsm".
//
// NeedsHub is false: this command must run before any .cn/ hub exists —
// that is the entire point of a repo installer.
//
// Domain logic lives in internal/repoinstall (eng/go §2.18); this file
// owns only flag parsing + git-root resolution + dispatch wiring.
type RepoInstallCmd struct{}

func (c *RepoInstallCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:     "repo-install",
		Summary:  "Install the base CNOS/CDS package set into this repository",
		Source:   SourceKernel,
		Tier:     TierKernel,
		NeedsHub: false,
	}
}

func (c *RepoInstallCmd) Help() string {
	return repoInstallHelp
}

func (c *RepoInstallCmd) Run(ctx context.Context, inv Invocation) error {
	for _, a := range inv.Args {
		if a == "--help" || a == "-h" {
			fmt.Fprint(inv.Stdout, repoInstallHelp)
			return nil
		}
	}

	args, err := repoinstall.ParseArgs(inv.Args)
	if err != nil {
		fmt.Fprintf(inv.Stderr, "✗ cn repo install: %s\n\n", err)
		fmt.Fprint(inv.Stderr, repoInstallHelp)
		return err
	}

	// Repo-root resolution: reuse the existing gitRepoRoot helper
	// (cmd_cell.go) rather than duplicating a `git rev-parse
	// --show-toplevel` wrapper. inv.HubPath is only populated when
	// main.go's discoverHub found an existing .cn/ directory — the
	// common case for a fresh repo is an empty HubPath, since this
	// command's entire purpose is to run before .cn/ exists.
	repoRoot := inv.HubPath
	if repoRoot == "" {
		root, gerr := gitRepoRoot(ctx)
		if gerr != nil {
			fmt.Fprintf(inv.Stderr, "✗ cn repo install must be run inside a Git repository.\n")
			return fmt.Errorf("repo install: not inside a git repository: %w", gerr)
		}
		repoRoot = root
	}

	opts := repoinstall.Options{
		RepoRoot:  repoRoot,
		Release:   args.Release,
		IndexPath: args.IndexPath,
		Packages:  args.Packages,
		Dispatch:  args.Dispatch,
		DryRun:    args.DryRun,
		Stdout:    inv.Stdout,
		Stderr:    inv.Stderr,
	}

	_, err = repoinstall.Run(ctx, opts)
	return err
}
