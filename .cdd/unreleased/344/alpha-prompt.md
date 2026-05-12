# α Dispatch — Cycle #344 (Cycle A)

```
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 344 --json title,body,state,comments
Branch: cycle/344
Tier 3 skills:
  - src/packages/cnos.core/skills/write/SKILL.md
  - src/packages/cnos.core/skills/skill/SKILL.md
```

## Context for δ

- Allowed tools: `--allowedTools "Read,Write,Bash"`
- Git identity for α commits on this cycle: `Alpha <alpha@cdd.cnos>`
- Mode: docs-only (design-and-build; design converged in issue body — 24 sections specified, 37 open questions each with a recommendation).
- Branch: `cycle/344` — γ created it; α checks out, never creates.
- On completion α sets `status: review-ready` in `.cdd/unreleased/344/self-coherence.md` and commits + pushes to `cycle/344`.

## Cycle A scope — what α must implement

**This is Cycle A of a 3-cycle meta-issue (#344). α implements Cycle A only.** Cycles B (reference notifier impl) and C (tsc adoption) are out of scope.

### Primary deliverable: `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md`

New file. 24 sections (§1–§24). Target length: 600–1100 lines.

**Required sections (must all be present, each ≥3 sentences):**

- `## §1 Purpose` — bootstrap an existing repo to operate under cdd
- `## §2 Pre-conditions` — repo exists, primary branch, host supports labels + secrets + CI (GitHub assumed; mark portability surface)
- `## §3 .cdd/ scaffold` — exact directory tree, one-line purpose per directory
- `## §4 Version pin` — `.cdd/CDD-VERSION` format (issue open question #1: SHA + tag, one per line)
- `## §5 Labels` — minimum label set: `cdd`, `mca`, `P0`–`P3`; project may add more
- `## §6 Branch convention` — `cycle/{N}` numbering rule, who allocates N
- `## §7 Identity convention` — references `operator/SKILL.md §"Git identity for role actors"` (cycle #343 landed form: `{role}@{project}.cdd.cnos`; cnos elision `{role}@cdd.cnos`)
- `## §8 Dispatch declaration` — `.cdd/DISPATCH` file: §5.1 vs §5.2 per `operator/SKILL.md §5` (cycle #342 dispatch configurations); affects honest-grading floor per `release/SKILL.md §3.8`
- `## §9 CI integration` — minimum CI surface: artifact validation (file presence per §12 template), test runs, spec runs, project-specific progressions; reference `scripts/validate-release-gate.sh` (cycle #339 mechanical pre-merge gate) as a canonical CI artifact check example
- `## §10 Notification integration` — generic transport-agnostic interface: event vocabulary (≥4 events: cycle-open, β-verdict, RC, merge), adapter contract; Telegram as reference adapter pointing to Cycle B; contract must not bake in Telegram assumptions (A.AC3)
- `## §11 Secrets` — exact secret names (`CDD_TELEGRAM_BOT_TOKEN`, `CDD_TELEGRAM_CHAT_ID`), where they live (GitHub repo secrets), explicit prohibition on committing tokens; labeled table (A.AC4); no example tokens in the file
- `## §12 Cycle-README template` — minimum file set per `releases/{version}/{N}/`: self-coherence.md, alpha-closeout.md, beta-review.md, gamma-closeout.md, cdd-iteration.md
- `## §13 Cross-repo trace bundle init` — pre-create `.cdd/iterations/cross-repo/` with README naming `{target}/{slug}/` convention (issue OQ #9: pre-create; #10: nested form); when to create a bundle; what files belong
- `## §14 Honest-claim manifest convention` — per-cycle `claims.md` at `releases/{version}/{N}/claims.md` (issue OQ #12: separate file); categories per `review/SKILL.md §3.13`: reproducibility, source-of-truth alignment, wiring; no empty manifests allowed (OQ #14)
- `## §15 MCA registry` — `.cdd/MCAs/INDEX.md` table of in-flight MCAs + per-MCA `.cdd/MCAs/{slug}/README.md` (OQ #15); close criteria (OQ #16); auto-suggest trigger (OQ #17: CI alert, not auto-create)
- `## §16 Skill bundle pull/sync` — vendored copy under `.cdd/skills/` (OQ #18); refresh on `.cdd/CDD-VERSION` bump only (OQ #19); integrity SHA check in CI (OQ #20)
- `## §17 Pre-commit / pre-push hooks` — opt-in symlink installation (OQ #21); structure-only locally (OQ #22); block on malformed structure, warn on missing optional (OQ #23)
- `## §18 Per-cycle dispatch override` — `releases/{version}/{N}/DISPATCH-OVERRIDE` with section reference + one-line reason (OQ #24); visible in γ-closeout TSC Grades (OQ #25)
- `## §19 Operator handle registry` — `.cdd/OPERATORS` table: handle, identity-email, role-scopes-permitted, added-cycle (OQ #26); warn first cycle, block thereafter (OQ #27); forward-only history (OQ #28)
- `## §20 Close-out SLA` — 24h default, repo-overridable in `.cdd/CADENCE` (OQ #29); alert only in Cycle B, block-on-breach in follow-on (OQ #30); merge-to-main is visibility (OQ #31)
- `## §21 Release cadence declaration` — `.cdd/CADENCE` enum: `versioned` / `rolling-docs` / `mixed` (OQ #32); prescribe both directory layouts (OQ #33); tagging discipline per cadence (OQ #34)
- `## §22 cdd-iteration cadence` — every cycle writes `cdd-iteration.md` (OQ #35); reference `ROLES.md §1` row 5 and `epsilon/SKILL.md` (cycle #345 landed ε role); cdd-iteration.md is ε's work product; auto-spawn MCA at N=5 same-axis findings in consecutive cycles (OQ #36); reuse review/SKILL.md severities + `info` (OQ #37)
- `## §23 Activation checklist` — numbered, 18–24 steps (one per content section §3–§22), each ≤2 sentences with a verification; no forward dependencies; runnable in order (A.AC5)
- `## §24 Verification` — canonical one-command check confirming activation succeeded

**Oracles:**
- `rg '^## §' src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` → 24 hits
- `wc -l src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` → 600–1100
- No `TODO` or `tbd` markers in the file

### Secondary deliverable: cross-references (A.AC2)

Insert ≤3-line pointer to `activation/SKILL.md` in each of:

1. `src/packages/cnos.cdd/skills/cdd/CDD.md` — in the protocol overview section; name activation as the bootstrap entry point for new tenants
2. `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — at §1 or the first-operator-setup context; pointer for first-time operators
3. `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` — note that activation findings (§22 / cdd-iteration) flow into the cdd-iteration cadence described in post-release

**Oracle:** `rg 'cdd/activation/SKILL.md' src/packages/cnos.cdd/skills/cdd/` → ≥3 hits; no broken relative links

### Tertiary deliverable: cnos self-activation marker files (A.AC6)

cnos predates this skill but must pass its own §24 verification at cycle merge. α creates the following files in cnos, populated per the skill prescription in the sections α is simultaneously authoring:

1. **`.cdd/CDD-VERSION`** — pin the cnos commit SHA currently at HEAD (and tag `3.74.0` if that is the latest tag); one SHA per line, tag on a second line if present
2. **`.cdd/DISPATCH`** — single-line record of the dispatch configuration this repo uses; consult `operator/SKILL.md §5` for the canonical form; record which configuration cnos cycles currently run under (§5.2 for γ-dispatch sessions in Claude Code; α should note the configuration observed for this cycle)
3. **`.cdd/CADENCE`** — `rolling-docs` (cnos uses date-anchored docs-only releases for cdd cycles, not version-tagged engine releases)
4. **`.cdd/MCAs/INDEX.md`** — table of in-flight MCAs in cnos; populate from open issues with the `mca` label (or mark empty if none); columns: slug, originating-cycle, target-close-cycle, owner, status
5. **`.cdd/OPERATORS`** — table of authorized δ handles for cnos; columns: handle, identity-email, role-scopes-permitted, added-cycle; populate from recent cycle close-out commit trailers or leave with `alpha@cdd.cnos` / `beta@cdd.cnos` / `gamma@cdd.cnos` as the known role actors
6. **`.cdd/skills/`** — cnos is itself the canonical skill source, so the "vendored copy" is recursive. Create `.cdd/skills/README.md` declaring this: "cnos is the origin repository for this skill bundle. The canonical skill location is `src/packages/cnos.cdd/skills/cdd/`. A vendored copy is not required because this repo IS the source. Tenant repos vendor from the cnos commit SHA pinned in `.cdd/CDD-VERSION`."

**Note on A.AC6 honest-claim invariant:** cnos must pass its own `activation/SKILL.md §24` verification check at the time β reviews. α documents any gaps discovered during self-application in `self-coherence.md §Known Debt`.

**Oracle:** the five marker files above exist in cnos main at merge time

## 37 open questions — α decides all

The issue body (§Open questions) carries 37 questions across 8 groups, each with a recommendation. α must:

1. Read every question and its recommendation
2. Decide whether to follow the recommendation or deviate (deviation requires one-sentence rationale)
3. Document every decision as a numbered table in `self-coherence.md §Open-Question Decisions` before signaling review-readiness

β will validate that all 37 decisions are recorded and that the skill text implements the chosen answer for each.

Decisions shape the skill text — e.g., OQ #1 (`.cdd/CDD-VERSION` format) directly determines what §4 says; OQ #2 (`.cdd/DISPATCH` format) determines §8; etc. There is no separate "decision doc" — the decisions live in self-coherence.md and the skill text is their implementation.

## Recently-landed context α must reference

These cycles landed before #344 dispatched; the activation skill references them rather than re-deriving:

- **cycle #343** — `operator/SKILL.md §"Git identity for role actors"`: canonical `{role}@{project}.cdd.cnos` form (cnos elision: `{role}@cdd.cnos`). The activation skill's §7 (identity convention) references this section directly.
- **cycle #342** — `operator/SKILL.md §5` (dispatch configurations: §5.1 canonical multi-session, §5.2 single-session δ-as-γ). The activation skill's §8 (dispatch declaration) references §5.
- **cycle #345** — `ROLES.md` at repo root; `epsilon/SKILL.md`. The activation skill's §22 (cdd-iteration cadence) must reference the ε role: `cdd-iteration.md` is ε's work product per `ROLES.md §1` row 5 and `epsilon/SKILL.md §1`.
- **cycle #339** — `scripts/validate-release-gate.sh` (mechanical pre-merge gate). The activation skill's §9 (CI integration) references this script as the canonical `.cdd/` artifact validation CI example.

## Non-goals (Cycle A — do not touch)

- `cdd/activation/templates/` — Cycle B only (reference notifier, GitHub Actions templates)
- tsc-side activation files — Cycle C only
- Any code changes, `VERSION`, `CHANGELOG.md`, `cn.json`, `scripts/`
- Hosts other than GitHub — note the portability surface in §2, but do not author adapters
- Any changes to `alpha/SKILL.md`, `beta/SKILL.md`, `gamma/SKILL.md`, `release/SKILL.md`, `review/SKILL.md` — cross-references go into CDD.md, operator/SKILL.md, post-release/SKILL.md only

## Dispatch note

δ: dispatch α with `claude -p` against the cnos repo on branch `cycle/344`. α loads its SKILL.md first and reads the full issue via `gh issue view 344 --json title,body,state,comments`. The branch already exists on origin — α does not create it. This cycle has 6 ACs but 37 open questions; α's primary context pressure is the decision table in self-coherence.md. Write and push self-coherence.md in incremental commits — do not accumulate all work in a single session without pushing.
