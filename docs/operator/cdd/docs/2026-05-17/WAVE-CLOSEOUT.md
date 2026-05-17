# Wave Closeout — 2026-05-17 — #369 + #370 Parallel γ-as-δ Dispatch

**Operator:** Sigma (cn-sigma hub; usurobor's CDD coherence partner)
**Wave:** parallel dispatch of #369 (Phase 2 schemas) + #370 (Phase 1.5 normal form) as γ-acting-as-δ subprocesses
**Wall clock:** ~50 minutes for both cycles, end-to-end (dispatch → merge → close-out)
**Outcome:** both cycles shipped same-day; #366 Phases 3–7 unblocked; six issues filed

---

## Cycles shipped

### #370 — Phase 1.5: Coherence-cell normal form

- α produced `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` (435 lines; 9 kernel sections in order: Preamble → Kernel → Cell Outcomes → Recursion Modes → Scope-Lift → Two-Layer Separation → Non-goals → Closure → ACs)
- β R1 → REQUEST CHANGES (F1: CDD-Trace heading drift; CI red)
- α R2 → fix at `846800e8` (rename §CDD-Trace → §CDD Trace)
- β R2 → APPROVED (F1 closed; CI green; merge-test PASS)
- γ-acting-as-δ merged at `0d9f7498` (`Closes #370` in subject)
- γ landed step 13a tooling-gap patch (`4a0115d2`): aligned §CDD-Trace → §CDD Trace across `alpha/SKILL.md` + `cn-cdd-verify` comment
- γ close-out at `00cceda1` (cycle dir moved to `.cdd/releases/docs/2026-05-17/370/`)
- Issue closed `2026-05-17T12:35:32Z`
- §9.1 triggers fired: avoidable tooling failure + loaded-skill miss × 2 — all triaged with disposition
- **Next-MCA filed by γ-370: #373** (preventive `--worktree` identity write across role skills)

### #369 — Phase 2: Schemas

- α produced `schemas/cdd/contract.cue`, `schemas/cdd/receipt.cue`, `schemas/cdd/boundary_decision.cue`, `schemas/cdd/README.md` (with §Scope-Lift Invariant), plus 4 fixtures (1 valid + 3 doctrine-load-bearing invalid exercising #367 AC3 + AC6 freezes)
- β R1 → REQUEST CHANGES (D1: rule 3.11b binding — missing `gamma-scaffold.md` on cycle branch)
- γ recovery path (a) at `227d2373` — scaffold authored, single-row re-check primed
- β R2 → APPROVED unconditionally (D1 closed; pre-merge gate all-green)
- γ-acting-as-δ merged at `ff54f2a0` (`Closes #369` in subject; no-ff merge)
- γ close-out at `bbc65e2c` (cycle dir moved + INDEX.md row + cdd-iteration.md)
- Issue closed `2026-05-17T12:40:59Z`
- §9.1 trigger fired: γ-side rule 3.11b discoverability gap
- **Next-MCA filed by γ-369: #375** (γ-side pre-dispatch gate mirroring β-side rule 3.11b)

### #367 — closed post-hoc (Phase 1 cycle from 2026-05-15)

- Caught at session start: #367 merged 2026-05-15 13:54 UTC but `gh issue view 367` returned `state: OPEN`
- Root cause: merge subject was reference `(#367)` not close-keyword `Closes #367`; docs-only disconnect path skipped explicit `gh issue close`
- Closed manually with comment pointing to merge + close-out
- Filed cnos **#368** proposing α merge subject template fix + γ close-out gh-state assertion

---

## Issues filed during this wave

| # | Title | Filer | Class |
|---|---|---|---|
| #368 | CDD close-token gap — docs-only disconnect cycles leave issues OPEN after merge | Operator (Sigma) | cdd-skill-gap |
| #371 | Require observable dispatch for cycle subprocesses (stream-json) — default text mode is operator-blind | Operator (Sigma) | cdd-skill-gap |
| #372 | Roadmap: Operator UX for CDD — observation-first, multi-surface | Operator (Sigma) | operator-ux-roadmap |
| #374 | Build cnos.telegram v0 — CI notifications inbound, operator-to-agent dispatch outbound (claude + codex hook) | Operator (Sigma) | tooling / new-package |
| #373 | Preventive `--worktree` identity write across all role skills when `extensions.worktreeConfig=true` | γ-370 | cdd-skill-gap (recurrence from #301 O8) |
| #375 | γ-side pre-dispatch gate mirroring β-side rule 3.11b (gamma-scaffold pre-dispatch) | γ-369 | cdd-skill-gap (γ-axis) |

**Six total — four operator-filed, two γ-filed in-cycle.**

---

## #366 roadmap status after wave

| Phase | Status | Tracker |
|---|---|---|
| Phase 0 — Doctrine | ✓ shipped | #364 |
| Phase 1 — Validation surface design | ✓ shipped | #367 |
| Phase 1.5 — Normal form | ✓ shipped (this wave) | #370 |
| Phase 2 — Schemas | ✓ shipped (this wave) | #369 |
| Phase 3 — Validator (`cn-cdd-verify` rewrite) | **unblocked, dispatchable** | (to be filed) |
| Phase 4 — δ split | gates on Phase 3; bundles #371 + #373 inputs | (to be filed) |
| Phase 5 — γ shrink | gates on Phase 4; may also consume #375 | (to be filed) |
| Phase 6 — ε relocation | gates on Phase 3; may run parallel to 4–5 | (to be filed) |
| Phase 7 — `CDD.md` rewrite | gates on Phases 3–6 stabilizing | (to be filed) |

**Phase 4 input bundle:**
- #371 — dispatch observability contract (harness emits stream-json by contract)
- #373 — preventive `--worktree` identity write across role skills

**γ-axis preventive (independent of phase gating):**
- #375 — γ pre-dispatch gate

---

## Key learnings

### 1. Parallel γ-as-δ dispatch works cleanly

Two cycles dispatched in parallel as γ-acting-as-δ subprocesses (Phase 4 δ split has not landed; γ assumes δ for now). Independent worktrees + branches. Both ran R1 RC → recovery → R2 APPROVED → merge → close-out without operator intervention. Total wall clock: ~50 minutes for two cycles producing 5 deliverable files + 2 merge commits + 2 close-outs + 2 next-MCA issues. This is the dispatch shape that should be standardized once `cn dispatch` lands (Phase 4).

### 2. Black-box dispatch is operator-hostile

Operator dispatched both cycles using default `claude -p` text-mode. Subprocess logs stayed 0 bytes for the entire 50 minutes. Operator-side observation reduced to `git log` polling on cycle branches as a filesystem proxy for in-process state. Had either γ-subprocess hit an unrecoverable error, operator would not have known until the subprocess exited. The fix is `--output-format stream-json` + JSONL capture — codified in #371 as a skill rule (operator/SKILL.md + gamma/SKILL.md), to be implemented as a harness primitive in Phase 4. Until #371 lands, operator-side standing rule applies.

### 3. γ instances exercise protocol-iteration discipline in-cycle

Both γ-369 and γ-370 detected protocol gaps from their own cycles and filed next-MCAs (#375 and #373 respectively) without operator prompting. γ-370 explicitly traced #373 back to a recurrence from cycle #301 O8 (first surfacing of the `--worktree` identity-leak class). This is exactly the ε-discipline pattern, exercised at γ-scope mid-cycle — γ's `cdd-iteration.md` finding → next-MCA issue file as standing protocol-evolution mechanism. Visible to operator only post-hoc because subprocess logs were blind (see learning 2).

### 4. Auth identity vs commit identity gap

Issues filed via `gh issue create` from the operator's droplet show **usurobor** as the GitHub author (gh CLI auth identity), regardless of which agent (γ-369, γ-370, operator) ran the command. Operator initially assumed #373 was filed by usurobor directly; actually γ-370 filed it. The issue body cites the filer, but the GitHub author field does not. For ε's Phase 6 reading of receipt streams across cycles, this matters: ε reading "filed by" from GitHub API gets the wrong actor. Worth noting in Phase 6 design.

### 5. Operator UX is observation-first, not command-first

Telegram was the initial frame for phone-side operator surface. This session's evidence falsified that framing: 80% of operator pain was observation (multi-stream tracking, decision-moment surfacing, async attention pattern), not input. Captured as #372 (operator UX roadmap) with three-surface architecture: web dashboard (primary) + telegram (alerts + quick commands) + terminal TUI (optional). Telegram-alone was the wrong v0; a focused telegram CI bridge + agent dispatch is the right v0 of #372 Phase 2 — that's #374.

### 6. R1 RC → R2 APPROVED is the default pattern, not the exception

Both cycles in this wave needed a recovery round. Pattern was clean both times:

- #369 R1 RC → γ recovery path (a) — γ writes missing artifact, β R2 single-row re-check
- #370 R1 RC → α R2 — α writes mechanical fix, β R2 verifies F1 closed

Suggests one full recovery round should be expected in cycle scope sizing, not framed as exceptional. The recovery paths are codified in `gamma/SKILL.md` §recovery and `review/SKILL.md` §3.11b binding; they worked as designed.

### 7. Step 13a tooling-gap patches land in-cycle without scope leak

γ-370 detected internal tool drift (cn-cdd-verify checks `§CDD Trace` but α/SKILL.md emitted `§CDD-Trace`) mid-cycle and landed a step 13a patch in the same cycle (`4a0115d2`). This is the §2.5b "documented patch landed inline" path — γ recorded it as a finding+patch in cdd-iteration.md without expanding cycle scope or smuggling unrelated work. The patch was contained, surface-bounded, and tooling-gap-classed. Working as intended.

---

## Open follow-ups

Out of scope for this wave but on operator's plate:

- **#371 / #373 dispatchable** as standalone skill-patch cycles, or as Phase 4 input batch
- **#375 dispatchable** as standalone γ-axis skill patch (small surface; could ship before Phase 4)
- **#372 Phase 1 (dashboard)** — design cycle, larger scope (~v0 web app)
- **#374 (cnos.telegram v0)** — design-and-build cycle, ~400–800 line package
- **Phase 3 cycle** for #366 — `cn-cdd-verify` rewrite against typed receptor (now unblocked by #369)
- **Older stale cycle/* branches** at origin (cycle/308, /327, /358 still listed as remote refs) — possible cleanup pass; pre-session debris, not from this wave

---

## Wave artifacts

| Artifact | Path |
|---|---|
| #370 deliverable | `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` |
| #369 deliverable | `schemas/cdd/` (3 CUE + README + fixtures/) |
| #370 cycle evidence | `.cdd/releases/docs/2026-05-17/370/` |
| #369 cycle evidence | `.cdd/releases/docs/2026-05-17/369/` |
| #370 γ PRA | `docs/gamma/cdd/docs/2026-05-17/POST-RELEASE-ASSESSMENT.md` |
| Operator session reflection | cn-sigma `threads/reflections/daily/20260517.md` |
| This wave closeout | `docs/operator/cdd/docs/2026-05-17/WAVE-CLOSEOUT.md` |
| Wave merge SHAs | #370: `0d9f7498` · #369: `ff54f2a0` · #370 close-out: `00cceda1` · #369 close-out: `bbc65e2c` |
