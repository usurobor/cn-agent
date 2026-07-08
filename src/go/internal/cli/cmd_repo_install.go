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
  separate, explicit opt-in (cnos#610): after the base install, it
  renders .github/workflows/cnos-cds-dispatch.yml via the cnos.core wake
  renderer (cnos#609), requires an explicit caller identity, is PR-only
  (this command never pushes to main), and requires the installing
  token to hold workflow scope. It also ensures the canonical cnos.core
  labels via label-doctor (cnos#493): it audits the installing repo's
  GitHub labels against cnos.core's labels.json and creates/repairs any
  missing or drifted one. This requires the installing repo to have a
  resolvable "origin" git remote and a GitHub token
  ($GITHUB_TOKEN/$GH_TOKEN) with permission to manage labels; if either
  is unavailable, --dispatch cds still renders the workflow but exits
  nonzero naming the label-doctor failure, and labels must be applied
  manually (e.g. via "cn label doctor").

FLAGS:
  --release V     "latest" (default) or a pinned release tag (e.g. 3.82.0)
  --index P       Package index path or http(s) URL — overrides --release
                   (local/offline/test workflows; skips release resolution)
  --packages CSV  Comma-separated package names (default: cnos.core,cnos.cdd,cnos.cds)
  --dispatch D    "none" (default) — base install only.
                  "cds" — also renders the dispatch workflow (see above).
  --agent NAME    Caller identity for --dispatch cds (default: sigma).
                  Any non-sigma agent requires --workflow-pat-secret.
  --workflow-pat-secret NAME
                  GitHub Actions secret name holding the dispatch agent's
                  workflow-scoped PAT. Required for --dispatch cds with a
                  non-sigma --agent (sigma defaults to SIGMA_WORKFLOW_PAT).
  --bot-name NAME Overrides the rendered claude-code-action bot_name input
                  (--dispatch cds only).
  --bot-id ID     Overrides the rendered claude-code-action bot_id input
                  (--dispatch cds only).
  --dry-run       Report the intended writes/fetches/restores; write nothing

EXAMPLES:
  cn repo install
  cn repo install --release latest --packages cnos.core,cnos.cdd,cnos.cds --dispatch none
  cn repo install --dry-run
  cn repo install --index ./dist/packages/index.json --packages cnos.core,cnos.cdd,cnos.cds
  cn repo install --dispatch cds --agent acme --workflow-pat-secret ACME_WORKFLOW_PAT \
    --bot-name acme-bot --bot-id 12345678

EXIT CODES:
  0  Success (installed, or --dry-run reported the plan)
  1  Error (not a git repo, package/index resolution failure, restore
     failure, or a --dispatch cds precondition: missing identity, or the
     cnos#493 label-doctor mechanism failing to resolve the installing
     repo's target/token or to reach the GitHub API)

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

	// Repo-root resolution: ALWAYS resolve via the actual git repository
	// root (git rev-parse --show-toplevel, wrapped by gitRepoRoot in
	// cmd_cell.go) — never via inv.HubPath. inv.HubPath is populated by
	// main.go's discoverHub(), which walks upward from cwd looking for
	// ANY ancestor .cn/ directory, with no bound at the current git
	// repository's root and no check that the found directory is even a
	// git repository at all. Trusting inv.HubPath here would let an
	// unrelated ancestor .cn/ (a different project, a stale hub, a
	// parent workspace, a monorepo-of-repos layout) silently hijack the
	// install root — exactly the scenario `cn repo install` must guard
	// against, since its entire premise is running in a repo that has no
	// .cn/ yet (β R0 finding, cnos#608). Always ask git directly; this
	// also makes the "not inside a git repository" failure unconditional
	// instead of only reachable when inv.HubPath happens to be empty.
	root, gerr := gitRepoRoot(ctx)
	if gerr != nil {
		fmt.Fprintf(inv.Stderr, "✗ cn repo install must be run inside a Git repository.\n")
		return fmt.Errorf("repo install: not inside a git repository: %w", gerr)
	}
	repoRoot := root

	opts := repoinstall.Options{
		RepoRoot:          repoRoot,
		Release:           args.Release,
		IndexPath:         args.IndexPath,
		Packages:          args.Packages,
		Dispatch:          args.Dispatch,
		DryRun:            args.DryRun,
		Agent:             args.Agent,
		WorkflowPatSecret: args.WorkflowPatSecret,
		BotName:           args.BotName,
		BotID:             args.BotID,
		Stdout:            inv.Stdout,
		Stderr:            inv.Stderr,
	}

	_, err = repoinstall.Run(ctx, opts)
	return err
}
