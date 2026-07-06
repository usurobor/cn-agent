# γ closeout — cycle/609

**Issue:** [cnos#609](https://github.com/usurobor/cnos/issues/609) — cds-install Sub 2: generalize `cn install wake` identity (agent/PAT/bot) — successor to #549.

**Branch:** `cycle/609` (base `main@612a96d6`).

**Mode:** substantial (renderer + CI-test + one-line prompt-consistency fix). No `.cdd/releases/` disconnect performed yet — that happens at merge time per `release/SKILL.md` §2.5a.

**Cycle execution mode:** wake-invoked-δ (per `cnos.cdd/skills/cdd/delta/SKILL.md` §9), driven by the `cds-dispatch` wake firing that claimed issue #609. δ performed the γ-scaffold investigation and authoring directly (reading the renderer, the design-oracle doc, and the existing CI-test patterns), then dispatched α and β as independent `Agent`-tool sub-sessions per §9.3 (α does not spawn β; β does not spawn α; no chat-state shared between them beyond the cycle branch + `.cdd/unreleased/609/` artifacts).

**Rounds:** 1 (R0 converge — α's R0 implementation passed β's R0 review with zero findings; no iteration needed).

**Filed by:** δ (cds-dispatch wake-invoked cycle; gamma-scaffold + gamma-closeout authored directly by the wake-invoked-δ session; α/β dispatched as separate agent sessions).

---

## Cycle summary

α generalized `cn-install-wake`'s hardcoded sigma-only substrate identity bindings into caller-supplied flags, per the pinned design oracle at `docs/development/design/cn-repo-install-MOCKS.md` §Mock C + §Mock E:

1. **`src/packages/cnos.core/commands/install-wake/cn-install-wake`** — new `--workflow-pat-secret`, `--bot-name`, `--bot-id` flags; fail-early identity resolution for non-sigma agents (Mock C2 — dies before any `--out` write); the three literal `SIGMA_WORKFLOW_PAT` occurrences replaced with the resolved `${workflow_pat_secret}`; the "Mechanical recovery scanner" and "Mechanical checkpoint + PR finalizer" steps branch on `$agent`: sigma keeps the exact `cd src/go && go build` self-build (byte-identical), any other agent gets a tenant-portable `install.sh`-based acquisition (Mock E2).
2. **`src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`** — one-line parameterization-consistency fix: the wake's own self-referential `cds-dispatch-sigma` mention (in the "Disallowed surfaces" bullet) now reads `cds-dispatch-{agent}`, matching the three other `{agent}`-substituted occurrences of the same value already in the file. The adjacent `agent-admin-sigma` (a different wake's identity) and a dated historical incident citation in the same bullet were deliberately left untouched — see gamma-scaffold.md Friction note 2.
3. **`.github/workflows/install-wake-golden.yml`** — five new CI steps: AC1 positive (acme render carries all four caller-supplied identity values), AC1 negative (Mock C2 fail-early, no partial `--out` file), AC2 two-layer zero-leak audit (a synthetic fixture manifest isolates the renderer's own leak-freedom from any package's legacy prompt prose, plus a scoped grep on the real acme cds-dispatch render for the renderer-controlled tokens), and E2/E4 (tenant acquisition present for non-sigma, `go build` self-build preserved for sigma).

4 commits on `cycle/609` (gamma-scaffold + α ×3 + β review); α/β closeouts filed alongside this one. Zero lines touched in `src/go/**`, `docs/development/design/cn-repo-install-MOCKS.md`, either committed golden (`cnos-cds-dispatch.golden.yml`, `cnos-agent-admin.golden.yml`), or the live `.github/workflows/cnos-cds-dispatch.yml` — all confirmed byte-identical post-change by both α (self-coherence.md) and β (beta-review.md) independently re-rendering and diffing.

---

## Closure declaration (pending PR merge)

All four in-scope ACs (AC1, AC2/Mock-C4, AC3/Mock-C5, gates) plus the two tenant-portability invariants (Mock E2, E4) PASS as independently verified by α (self-coherence.md §R0) and re-verified from scratch by β (beta-review.md §R0, `verdict: converge`, zero findings).

| AC / invariant | Status | Primary surface |
|---|---|---|
| AC1 — configurable identity (Mock C1) | PASS | `cn-install-wake` new flags + identity-resolution block |
| AC1 negative / Mock C2 — fail-early, no partial render | PASS | `cn-install-wake` `die()` ordering (verified by β to fire before any `tmp_out`/`out_path` write) |
| AC2 / Mock C4 — zero sigma leak | PASS | `install-wake-golden.yml` synthetic-fixture + scoped-grep steps; `cds-dispatch/SKILL.md` self-reference fix |
| AC3 / Mock C5 / E4 — sigma byte-identical | PASS | existing golden-diff + idempotence steps (regression oracle); zero diff confirmed independently by α and β |
| E2 — tenant acquisition, no `src/go` build | PASS | `cn-install-wake` non-sigma branch of the scanner/finalizer steps |
| AC4 — gates green | PASS (install-wake-golden surface); no `src/go/**` touched so no separate Go-gate run needed for this cycle | CI + local verification by α and β |

Full per-invariant evidence is in [α's self-coherence.md §AC verification](./self-coherence.md) and [β's beta-review.md §R0 table](./beta-review.md); γ does not re-derive the evidence here.

`mock_parity` (per issue #609's "Parity requirement"):

```
mock_parity:
  - id: C2
    invariant: "missing --workflow-pat-secret/identity fails early, no partial render"
    verified_by: [alpha, beta]
    missed: 0
  - id: C4
    invariant: "zero sigma/SIGMA_WORKFLOW_PAT/41898282 leak for agent != sigma"
    verified_by: [alpha, beta]
    missed: 0
  - id: C5
    invariant: "--agent sigma reproduces committed golden byte-for-byte"
    verified_by: [alpha, beta]
    missed: 0
  - id: E2
    invariant: "tenant (non-sigma) render has no cd src/go / go build ./cmd/cn"
    verified_by: [alpha, beta]
    missed: 0
  - id: E4
    invariant: "--agent sigma still self-builds via go build, byte-for-byte"
    verified_by: [alpha, beta]
    missed: 0
```

---

## Process-gap audit (γ triage)

| ID | Finding | Source | CAP class | Disposition |
|---|---|---|---|---|
| **F-γ-1** | The tenant-install temp path is hardcoded to `/tmp/cnos-cn-609` (cycle-numbered), and the scanner/finalizer steps' install logic is written once and reused rather than duplicated. If a future cycle touches only one of the two steps, the shared block could drift. | α self-coherence.md §Debt | cdd-skill-gap (minor; renderer-internal naming convention) | **Drop.** The cycle-numbered path is cosmetic (any fixed, collision-free path works — GitHub Actions runners are single-use); renaming it later is a zero-risk one-line change if it's ever noticed as odd. Not worth a follow-up issue. |
| **F-γ-2** | `--git-user-name`/`--git-user-email` flags requested in a later issue comment were deliberately not implemented — the pinned design oracle (Mock C/E) does not include them, and no existing renderer output emits a `git config user.*` step to parameterize. | gamma-scaffold.md Friction note 1 | project scope-clarification (not a defect) | **Issue (follow-up, not filed in this cycle).** If the operator confirms this is still wanted, it should land as its own sub (a real behavior addition: a new conditionally-rendered `git config` step), not be folded into this cycle's "generalize existing hardcoded values" scope. Recording here so a future triage pass sees the deliberate omission and its reasoning rather than mistaking it for an oversight. |
| **F-γ-3** | AC2's zero-leak oracle could not be run as a single unconditional `grep -i sigma` against the live `cds-dispatch` wake, because that wake's prompt body carries two out-of-scope literal "sigma" strings (a different wake's concurrency-group name; a historical incident citation) that this cycle's scope explicitly excludes from editing. | gamma-scaffold.md Friction note 2 | cdd-skill-gap (oracle-design, not implementation) | **Drop, with note for future wake-provider work.** The two-layer split (synthetic fixture + scoped grep) is the correct oracle shape for this exact class of problem (renderer generalization inside a package whose own prose has historical/cross-wake references); no process change needed. If a future cycle wants a stricter single-command oracle, that would require either scrubbing all package-owned wake prompts of every substring match (a much larger, out-of-scope prose-hygiene pass) or accepting the split-oracle shape as the standing pattern — γ recommends the latter. |

No findings required a patch within this cycle's commit set; all three are either cosmetic, deliberate-scope documentation, or oracle-design notes for future reference.

---

## What remains

- Open the cycle-PR scoped to #609 (β approval δ requests the review transition against).
- Write `REVIEW-REQUEST.yml` per `dispatch-protocol/SKILL.md` §2.4 / `cnos#569` Phase 2.
- Request `status:in-progress → status:review` via `cn issues fsm evaluate --issue 609 --apply`.
- Post the PR-URL comment on issue #609 (δ's `status:review` return token per `delta/SKILL.md` §9.6).

This cycle does not merge itself — merge + release-boundary actions are the external reviewer's / operator's authority per `dispatch-protocol/SKILL.md` §"Lifecycle transitions" (the dispatch wake never transitions out of `status:review`).
