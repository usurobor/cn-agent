# α close-out — cycle #656

## Summary

Implemented cnos#656 (Phase 1 of the cnos#655 `cn repo` lifecycle wave):
`.cn/repo.state.json` managed-surface ledger (schema `cn.repo.state.v1`,
checked-in CUE schema + CI validation gate), `cn repo install` writing it
idempotently, and a new read-only `cn repo status` command. Full detail in
`self-coherence.md`; this file is the terse process record.

## What shipped

- `schemas/repo_state.cue` + `schemas/fixtures/repo-state/{valid,invalid}/`
  + `scripts/ci/validate-repo-state.sh` + a new CI job
  (`repo-state-schema-check` in `.github/workflows/build.yml`).
- `internal/repostate` (new package) — pure types + deterministic
  `Marshal()`.
- `internal/repoinstall` — `cn repo install` writes/updates the ledger;
  survives a `--dispatch cds` label-doctor failure when the render itself
  succeeded (does NOT survive a stale-file-from-a-prior-run
  misidentification — see the bug found/fixed below).
- `internal/dispatchrender` (new package, extracted from
  `repoinstall.runDispatchCds`) — single renderer-invocation authority
  shared by install and status.
- `internal/repostatus` + `cli.RepoStatusCmd` (`cn repo status`, `--json`,
  `--check`) — read-only; A2 drift classification via a real fresh-render
  comparison.
- `docs/development/design/cn-repo-status-MOCKS.md` (mock-first, A6).
- A small additive change to `label-doctor` (`Result.LiveLabels`) needed
  for "unknown/non-canonical" label reporting.

## Process notes

- `run_class: first_pass` — clean claim, no prior rejection.
- No γ→α→β iteration round in the traditional sense: this cycle was run
  by a single interactive session acting as δ, routing γ/α/β itself. The
  "iteration" that did happen was α's own pre-review critical re-read
  (caught the `DriftRemoved` gap) plus a deliberately-dispatched
  independent adversarial review subagent (caught two more real bugs —
  see `beta-review.md`) — both folded into this same pass rather than
  requiring a separate β-reject/α-repair round.
- Genuine end-to-end proof gathered, not just unit tests: a real
  `cn build` → real local package index → real `cn repo install` → real
  `.cn/repo.state.json` validated against the CUE schema via `cue vet`
  (not a hand-written fixture standing in for real output).
- Existing-behavior regression check: `internal/repoinstall`'s and
  `internal/cli`'s full pre-existing suites pass unmodified except one
  test updated to also stage the new `.cn/repo.state.json` file (a
  legitimate consequence of adding a new install-time artifact, not a
  weakened assertion).

## Evidence

```
$ go build ./... && go vet ./... && go test ./...
(all packages ok)

$ cd src/packages/cnos.core/commands/label-doctor && go build ./... && go vet ./... && go test ./...
ok

$ ./scripts/ci/validate-repo-state.sh --self-test
✓ all checks passed.
```

Manual end-to-end smoke (real build → install → status → hand-edit →
status → --check exit codes → no-write verification) — see
`self-coherence.md` Self-check section for the full transcript summary.
