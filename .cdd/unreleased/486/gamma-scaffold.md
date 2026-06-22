---
cycle: 486
parent_issue: cnos#486
master_tracker: cnos#467 (Sub 5B of wake-orchestration wave)
cycle_branch: cycle/486
authored_by: γ@cdd.cnos (bootstrap-δ via δ-interface session)
date: 2026-06-22 (UTC)
base_main_sha: 950730c74985864537696ec45ebf0023fde16b97
delta_skill_base_sha256_at_base: pinned-in-§5-AC8 (α captures + records in self-coherence §R0)
predecessor_sub_5a: cnos#485 / PR #488 / merged at base sha `950730c7`
successor_sub_5c: cnos#487 (deferred; not yet authorized)
---

# γ-scaffold — cnos#486 (cdd/delta: define dispatch-wake-invoked δ mode amendment)

## 1. Frontmatter / header

| Field | Value |
|---|---|
| Cycle | 486 |
| Parent issue | [cnos#486](https://github.com/usurobor/cnos/issues/486) |
| Master tracker | [cnos#467](https://github.com/usurobor/cnos/issues/467) (Sub 5B of the wake-orchestration wave) |
| Cycle branch | `cycle/486` |
| Branched from | `main` at `950730c7` (post-Sub-5A merge state; PR #488 merged) |
| Mode | `design-and-build` (single SKILL.md amendment; no runtime activation; no smoke) |
| Authored by | γ@cdd.cnos (bootstrap-δ via δ-interface session) |
| Date | 2026-06-22 (UTC) |

## 2. Goal

Amend `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` with a new top-level section — **"Dispatch-wake-invoked mode"** (or equivalent; α decides the precise heading per §6 below) — that lifts the empirical bootstrap-δ pattern (γ-interface drives, spawns γ/α/β as Agent sub-sessions; observed in cycles 470, 476, 485) into doctrine. The amendment defines, in the substrate-agnostic role-skill register:

- **The input contract** — what δ receives when a package-owned dispatch wake (e.g. cds-dispatch per cnos#483) hands it a claimed cell after the cnos#454 §2.2 claim sequence.
- **The δ routing sequence** in wake-invoked mode — how γ is spawned and what context γ needs; how α is spawned and what context α needs; how β is spawned and what context β needs; the operator-directed invariant *δ dispatches every role; γ does not spawn α/β* (which the current bootstrap-δ implicitly violates).
- **R[N] iteration token discipline** — how each round's verdict surfaces to the wake (and through the wake, to the substrate's observability layer) without dictating the substrate encoding.
- **The per-round artifact contract** under `.cdd/unreleased/{N}/` — what MUST be on the cycle branch at each R[N] boundary for the round to be considered complete; this is the framework-side restatement of cnos#483's `output_contract.artifact_class_taxonomy`.
- **Return tokens to the wake** — `status:review` after β converge, `status:blocked` with reason on hard block, claim release on race / invalid claim state (mapped to cnos#454 §2.6 drift handling); plus the cycle-PR URL when applicable.
- **v0 substrate constraints** — substrate timeout horizon (named by class, not by GHA-default 360 min number); per-protocol concurrency group (cnos#454 §2.3 layer 1); single-claim-per-firing (cnos#454 §3.2). Named as constraints δ honors, NOT substrate emission.
- **Empirical anchor citations** — cycles 470, 476, 485 as bootstrap-δ runs whose observed friction informs this contract (3-round class-trap recurrence per cnos#478; first R0-converge under mechanical-injection discipline at cycle/485; branch-as-shared-state pattern as the empirical witness for the no-chat-state discipline).

After this cycle ships, cnos#483's `cds-dispatch/prompt.md` forward-reference "δ's wake-invoked mode is the contract that lands in cnos#467 Sub 5" resolves to a landed citation; cnos#487 (Sub 5C; flip cnos#483's `activation_state` from `declaration-only` to `live` + render + commit + smoke) has a concrete, testable contract for what δ does after the wake's claim sequence hands it a cell.

Cite: cnos#486 (this issue body), cnos#467 (master tracker — wave architecture; doctrine-correction header at top is authoritative: `cnos.cdd` = framework; `cnos.cds` = concrete software protocol; `cds-dispatch` = wake; `protocol:cds` = qualifier).

## 3. Branch name

`cycle/486`. γ has created this branch from `main` at `950730c7` and pushed `gamma-scaffold.md`. α implements on this branch and pushes. β reviews on this branch.

## 4. Touched files (expected)

### α MUST modify

- **`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`** — add a new top-level section (or equivalent extension) named "Dispatch-wake-invoked mode" (α picks the precise heading; recommended placement: after §1 "Outward membrane" and §2 "Inward membrane", before §3 "Override" — i.e. a new §3 or §8 numbered slot, with the existing §3–§8 renumbered if α chooses pre-override placement). The amendment must satisfy all 9 ACs in §5 below.

### α MUST create

- **`.cdd/unreleased/486/self-coherence.md`** — α's per-AC verification table + §R0 (initial) section. Future §R[N] sections appended as β iterates. Sections (per `alpha/SKILL.md` §2.5): §Gap / §Skills / §ACs / §Self-check / §Debt / §CDD Trace / §Review-readiness. Plus an explicit `§Design` section naming the amendment's heading choice + placement + the substrate-agnosticism rationale; plus a `§Friction notes` section if α surfaces new observations beyond γ's.

### α MAY modify (optional; α decides + documents in §Design)

None expected. The amendment is doctrinal (skill-internal), and the surface is bounded by AC8's substrate-agnosticism gate. Specifically α should **NOT** open cross-skill cleanup PRs in this cycle (e.g. amending `gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`, `harness/SKILL.md`, `release-effector/SKILL.md`, `operator/SKILL.md`); those are P1 recommendations from the cycle/485 closeouts (β-skill amendment cluster) that the operator may bundle into a separate β-skill cycle. **If α wants to add an internal cross-reference from `delta/SKILL.md` to the wake-provider / dispatch-protocol skills, that is in-scope; surface in §Design.**

### α MUST NOT touch

- Any other `cnos.cdd/skills/cdd/*/SKILL.md` file (`gamma/`, `alpha/`, `beta/`, `operator/`, `harness/`, `release-effector/`, `post-release/`, `release/`, `review/`, `issue/`, etc.). This cycle is **δ-only**. The β-skill / γ-skill / harness-skill amendments named in cycle/485 closeouts §3 triage table (T1, T2, T11, T12, T13, T14) are separate cycles; bundling them here violates cycle/486's scope.
- `src/packages/cnos.cdd/skills/cdd/SKILL.md` (framework overview) — read-only reference; do not amend.
- `src/packages/cnos.cdd/skills/cdd/CDD.md` (canonical algorithm) — read-only reference; do not amend.
- `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` (cnos#454 / PR #466) — out of scope; the δ amendment CITES this skill but does not modify it.
- `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` (cnos#470) — out of scope; the δ amendment CITES this skill (especially §3.3 substrate-leakage rule + §3.9 dispatch fields + §3.10 activation_state) but does not modify it.
- `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` (cnos#468) — out of scope; cite-only.
- `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` (cnos#483) — out of scope; cited as the reference dispatch wake; Sub 5C / cnos#487 flips `activation_state`.
- `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` (cnos#483) — out of scope; cited.
- `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` (cnos#485) — out of scope; cited as the rendered golden the wake would emit.
- `src/packages/cnos.core/commands/install-wake/cn-install-wake` (cnos#476 + cnos#485) — renderer; out of scope.
- `.github/workflows/` — no substrate emission. **AC8 substrate-agnosticism gate.**
- `.github/workflows/cnos-cds-dispatch.yml` (does not exist yet; Sub 5C lands it) — explicitly out of scope.
- `.cdd/iterations/INDEX.md`, `.cdd/unreleased/485/*`, `.cdd/releases/docs/2026-06-21/*` — read-only reference; do not modify.

## 5. AC-by-AC oracle

The 9 ACs are quoted in summary form from cnos#486's "Acceptance criteria" section; full text in the issue body. For each, α produces verifiable oracle output; β re-runs each oracle on cycle/486 HEAD before issuing a verdict.

The `<delta-skill>` placeholder below is `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`.

### AC1 — §"Dispatch-wake-invoked mode" section added

**Statement (cnos#486):** `cnos.cdd/skills/cdd/delta/SKILL.md` carries a new top-level section (or extension to an existing one) explicitly named "wake-invoked mode" or equivalent. cnos#483's cds-dispatch prompt's forward-reference (cnos#467 Sub 5) now points at a landed citation.

**Oracle (α-runnable; β re-runs):**
```bash
# Grep for the wake-invoked heading; expect ≥1 hit in a markdown section heading.
grep -nE '^#+.*wake-invoked' src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
# Verify forward-reference resolution: the cds-dispatch prompt cites this section
# (path-relative or text-citation; either is acceptable). At minimum a grep for the
# cycle's anchor string returns a hit on the amended delta SKILL:
grep -nE 'wake-invoked|dispatch-wake-invoked' src/packages/cnos.cdd/skills/cdd/delta/SKILL.md | head -5
```

β verifies empirically: the section heading is a real top-level (`##` or higher) section, not buried in a bullet list or sub-paragraph. The heading must be grep-discoverable from cnos#483's prompt.md's existing forward-reference language without requiring further index updates (the forward-reference language in `cds-dispatch/prompt.md` line ~73 reads: *"δ's wake-invoked mode is the contract that lands in cnos#467 Sub 5"* — β confirms the amended SKILL satisfies that reference).

### AC2 — Input contract named

**Statement (cnos#486):** The new section enumerates the inputs δ receives from the wake: claimed-issue-number, protocol identifier, current-main-SHA, wake-run-id, package-runtime context (path to the concrete protocol's skills + the package's own selection/lifecycle skills). Each named input is listed with a one-sentence semantics description.

**Oracle:**
```bash
f=src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
# The 5 named inputs each appear at least once in the new section. The grep
# checks the input names (not exact wording; α may use synonyms but a reviewer
# must be able to map each named input to its presence):
for term in 'claimed.issue\|issue.number\|cell.issue' \
            'protocol.identifier\|protocol qualifier\|protocol:cds\|protocol:.P' \
            'current.main\|main SHA\|head.sha\|head commit' \
            'wake.run.id\|run id\|substrate run' \
            'package.runtime\|concrete protocol.*skill\|cnos\.cds/skills\|cnos\.cdr/skills'; do
  grep -ciE "$term" "$f" || echo "AC2: missing $term"
done
# All five must return ≥1 (β confirms each grep hit is in the wake-invoked
# section, not elsewhere in the file).
```

β verifies empirically: each input appears with a one-sentence semantics description (not just listed as a name). E.g. "claimed-issue-number: the GitHub issue number the wake claimed at §2.2 step 4 of cnos#454; δ reads this issue's body as the cell-shaped specification per `cdd/issue/SKILL.md`" — not just the bare word.

### AC3 — Routing sequence named

**Statement (cnos#486):** The new section names how δ spawns γ → α → β in the wake-firing environment, including what context each role needs to do its job without consulting hidden state. Each role's spawn-context is listed; the sequence is grep-able as a discrete step list; cycle/470 + cycle/476 are cited as the empirical case.

**Oracle:**
```bash
f=src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
# A discrete step list (numbered or bulleted) names the three role-spawn steps:
grep -nE '^\s*[0-9]+\.\s+(spawn|dispatch|route|invoke)\s+(γ|gamma|α|alpha|β|beta)' "$f" \
  || grep -nE '^\s*[-*]\s+(spawn|dispatch|route|invoke)\s+(γ|gamma|α|alpha|β|beta)' "$f"
# Each role's spawn-context is enumerated (the amended section names what γ gets,
# what α gets, what β gets):
grep -ciE 'γ.*receives|γ.*context|γ.*spawn-context|gamma.*receives|gamma.*context' "$f"
grep -ciE 'α.*receives|α.*context|α.*spawn-context|alpha.*receives|alpha.*context' "$f"
grep -ciE 'β.*receives|β.*context|β.*spawn-context|beta.*receives|beta.*context' "$f"
# Empirical anchor citations: cycle/470 and cycle/476 named with path:
grep -nE 'cycle/470|cycles?/470|releases/docs/2026-06-21/470' "$f"
grep -nE 'cycle/476|cycles?/476|releases/docs/2026-06-21/476' "$f"
```

β verifies empirically: the sequence is a discrete, grep-able step list (numbered ideally); each role spawn names the *minimum context* the role needs (scaffold + branch state + source-of-truth files; nothing else) — i.e. the "no chat state" discipline is encoded as the contract, not as a footnote. The operator directive "δ dispatches every role; γ does not spawn α/β" must surface as the routing sequence's invariant: γ scaffolds and exits; δ dispatches α; α implements and exits; δ dispatches β; β reviews and exits; δ iterates by re-dispatching α (or γ if a scaffold-side ambiguity surfaces).

### AC4 — Iteration token discipline

**Statement (cnos#486):** The new section names how R[N] tokens surface to the wake (issue comments naming `R[N] complete; β verdict: iterate|converge`; or branch-commit messages; or `.cdd/unreleased/{N}/*.R{N}.json` artifacts — the choice is α/γ's design call, but the doctrine names the requirement that R[N] is wake-observable). The discipline is grep-able; cycle/470 + cycle/476 empirical observation of R[N] boundaries is cited.

**Oracle:**
```bash
f=src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
# R[N] discipline named (grep for the R[N] token + iteration mechanic):
grep -nE 'R\[N\]|R\{N\}|round N|R0|R1' "$f" | head -10
# Iteration mechanic named (one of: issue comments / branch-commit messages /
# .cdd/ artifacts — the section enumerates the choices and pins the contract):
grep -ciE 'wake.observable|wake observe|surface.*wake|return.*wake|signal.*wake' "$f"
# Cycle/470 and cycle/476 empirical iteration observations cited:
grep -nE '(cycle/470|cycle/476).*(round|iterate|R[12])' "$f" \
  || grep -nE '(cycle/470|cycle/476)' "$f"
```

β verifies empirically: the discipline names a concrete mechanism (the choice; α/γ have latitude per the issue body) AND the invariant (R[N] is wake-observable; the wake doesn't need to read chat-state to know iteration progress). The empirical anchor names cycle/470's R1 finding (broken relative-link path; substantive ambiguity) and cycle/476's R1+R2 findings (CI bash-e class-traps; mechanical class) as the friction the discipline absorbs.

### AC5 — Artifact set per round

**Statement (cnos#486):** The new section names what `.cdd/unreleased/{N}/` MUST contain at each R[N] boundary: at R0 → `gamma-scaffold.md`; at R[N>0] → `self-coherence.md §R[N]` + `beta-review.md §R[N]`; at converge → all closeout artifacts (`alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, optional PRA). The per-round artifact set is enumerated; matches the cnos#483 cds-dispatch manifest's `output_contract.artifact_class_taxonomy`.

**Oracle:**
```bash
f=src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
# The seven canonical artifact classes from cnos#483 manifest's
# output_contract.artifact_class_taxonomy:
for cls in gamma-scaffold self-coherence beta-review alpha-closeout \
           beta-closeout gamma-closeout post-release-assessment; do
  grep -qF "$cls" "$f" || echo "AC5: missing artifact class $cls"
done
# R0 boundary names gamma-scaffold; R[N>0] boundary names self-coherence + beta-review:
grep -nE 'R0.*gamma-scaffold|gamma-scaffold.*R0' "$f"
grep -nE 'R\[N.*\].*self-coherence|self-coherence.*§R\[N' "$f"
grep -nE 'R\[N.*\].*beta-review|beta-review.*§R\[N' "$f"
# Converge boundary names the three closeouts (PRA optional):
grep -ciE 'converge.*closeout|closeout.*converge|on converge.*alpha-closeout' "$f"
```

β verifies empirically: the artifact taxonomy matches the cnos#483 manifest (re-read `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` → `output_contract.artifact_class_taxonomy`); no artifact is dropped, no extra artifact is invented. The "PRA optional" framing is explicit (per the manifest's `artifact_class_notes`: "PRA (post-release-assessment) is scope-dependent — only cycles with explicit retrospective value carry it").

### AC6 — Return token defined

**Statement (cnos#486):** The new section names what δ writes back to the wake when the cycle reaches a stable state: `status:review` after β converge + cycle-PR URL; `status:blocked` + a reason on hard block; release (`status:in-progress → status:todo` + race comment) per cnos#454's serialized-claim-guard semantics. The return tokens are listed; the wake-observable mechanism (label state + comment) is named.

**Oracle:**
```bash
f=src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
# The three named return tokens each surface:
grep -ciE 'status:review.*β converge|β converge.*status:review|converge.*status:review' "$f"
grep -ciE 'status:blocked.*reason|hard block.*status:blocked|blocked.*reason comment' "$f"
grep -ciE 'release.*claim|claim release|status:in-progress.*status:todo.*race' "$f"
# The cycle-PR URL surface is named on converge:
grep -ciE 'cycle.PR|cycle-PR|pull request.*url|PR url' "$f"
# Wake-observable mechanism (label transitions + comments per cnos#454):
grep -nE 'cnos#454|dispatch-protocol' "$f"
```

β verifies empirically: the three return tokens map to the four lifecycle transitions cnos#454 §2.4 + cnos#483's `responsibilities` field §5 enumerate (claim; β converge; hard block; release-back-to-queue). The amendment must NOT invent new lifecycle transitions; it MUST consume the dispatch-protocol's transitions and name which one δ triggers under which condition. The status:changes transition is explicitly EXTERNAL (operator/planner) and out of δ's wake-invoked authority — the amendment names this carve-out so a future reader doesn't conflate β iteration (internal; cell stays status:in-progress) with external rejection (operator-applied status:changes).

### AC7 — Honest about substrate constraints

**Statement (cnos#486):** The new section names the v0 constraints of the wake firing's environment: substrate timeout; concurrency group serialization per cnos#454 §2.3 layer 1; single-claim-per-firing per cnos#454 §3.2. The v0 simplifying assumption ("cell completes or reaches a state-stable boundary within one substrate firing") is named.

**Oracle:**
```bash
f=src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
# Substrate timeout horizon named as a CLASS, not a specific number (AC8 will
# bite if α writes "360 min" or "GitHub Actions default 360"):
grep -ciE 'substrate timeout|wake firing.*time|substrate.*horizon|single.firing.*horizon' "$f"
# Concurrency group serialization per cnos#454 §2.3:
grep -ciE 'concurrency.*serial|serialize.*per.protocol|per.protocol concurrency|cnos#454.*§2.3' "$f"
# Single-claim-per-firing per cnos#454 §3.2:
grep -ciE 'single.claim.per.firing|one.claim.per.firing|cnos#454.*§3.2' "$f"
# v0 simplifying assumption named:
grep -ciE 'state.stable boundary|reaches.*stable.*within|cell completes.*within.*firing' "$f"
```

β verifies empirically: the constraints are named as constraints δ HONORS (i.e. inputs to δ's behavior), NOT substrate emission. The wording must avoid:
- Substrate-specific timeout numbers (e.g. "360 min", "60 min" — α can describe the class as "the substrate's per-firing timeout" without quoting a number; if α cites a number, it must be in a "for the v0 GitHub Actions substrate today" carve-out that AC8's grep allows — but cleaner is to name it abstractly).
- GitHub Actions YAML syntax (`concurrency: cancel-in-progress: false`, etc.) — even when discussing the concurrency-group constraint, α describes it abstractly: "the substrate's per-protocol concurrency group ensures at most one wake firing per repository × wake claims at a time" (no YAML).

The v0 assumption is the load-bearing one: "the wake fires once and the cell completes (or reaches a state-stable boundary) within one substrate firing". This is named explicitly (cnos#486 issue body quotes it).

### AC8 — No substrate emission

**Statement (cnos#486):** The δ skill amendment MUST NOT emit GitHub Actions YAML, runner names, secret-name bindings, or `${{ }}` interpolations — same substrate-agnosticism gate the wake-provider skill enforces per §3.3. `grep -ciE 'github|workflow|yaml|GITHUB_TOKEN|runs-on|claude-code-action' delta/SKILL.md` returns 0 (or matches only descriptive references in a "relationship to substrate" carve-out).

**Oracle:**
```bash
f=src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
# Hard-fail grep: substrate tokens MUST NOT appear in the new section. Note
# this grep also picks up pre-existing references in the file (e.g. the existing
# §8 BOX-AND-THE-RUNNER section already cites GitHub Actions YAML in workflow-
# artifact contexts). The TARGET is that the NEW dispatch-wake-invoked section
# adds zero new substrate-leakage matches. β interprets:
n_new=$(git diff main...HEAD -- "$f" | grep '^+' | grep -ciE \
  'github|workflow|yaml|GITHUB_TOKEN|runs-on|claude-code-action|\$\{\{' || true)
# Acceptable matches in the diff:
#  - Citations of cnos#454, cnos#467, cnos#470, cnos#483 (issue refs — not substrate)
#  - Citations of other SKILL.md files at known paths (e.g. wake-provider/SKILL.md §3.3)
#  - A clearly-delimited "relationship to substrate" carve-out paragraph that
#    NAMES substrate tokens descriptively (i.e. "this skill is substrate-agnostic;
#    when the substrate is GitHub Actions, the wake firing is mapped to a workflow
#    run — but the δ contract names no GHA syntax")
#
# Forbidden matches:
#  - Any `${{ ... }}` interpolation
#  - Any `runs-on: ubuntu-latest` or similar
#  - Any `claude-code-action@v1` substrate-action citation
#  - Any GHA YAML block (e.g. job-level `if:`, `concurrency:` block syntax)
#  - Any secret-name binding (e.g. `secrets.GITHUB_TOKEN`, `secrets.PAT`)

# Strict count (β surfaces as a finding if n_new > 0 except in the explicit
# "relationship to substrate" descriptive carve-out, which β can identify
# by reading the diff context — typically a single paragraph at the bottom of
# the new section explicitly framed as "what the substrate looks like under
# v0 GitHub Actions"):
echo "AC8 new substrate-leakage matches in diff: $n_new"
```

**β review point (paired with cnos#470 wake-provider/SKILL.md §3.3 substrate-leakage rule):** β reads the diff line by line and confirms that any substrate token appears only in:
- A discrete, single, clearly-marked "relationship to substrate (descriptive only)" paragraph, OR
- Citations of issue numbers / SKILL.md path references (not substrate emission).

If α writes substrate YAML in a code block, that is a finding. If α writes `${{ ... }}` interpolations anywhere, that is a finding. If α writes "the wake claim runs on `ubuntu-latest`", that is a finding (α should write "the wake claim runs in the substrate-provided environment per cnos#454 §2.2 step 6"). If α writes "δ writes a `secrets.GITHUB_TOKEN`-authorized issue comment", that is a finding (α should write "δ writes an authenticated issue comment per cnos#454 §2.2 step 4; the substrate provides the authentication carrier").

### AC9 — Bootstrap-friction empirical case cited

**Statement (cnos#486):** The amendment cites cycle/470 + cycle/476 (bootstrap-δ; γ-interface spawning γ/α/β as Agent sub-sessions) as the empirical anchor for the contract; not as a permanent solution, but as the empirical observation set that informs what context each role needs. `grep` finds `cycle/470` or `cycle/476` in the amended SKILL with a citation to `.cdd/releases/docs/2026-06-21/`.

**Oracle:**
```bash
f=src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
# Cycle/470 cited (or its release-docs path):
grep -nE 'cycle/470|releases/docs/2026-06-21/470' "$f"
# Cycle/476 cited (or its release-docs path):
grep -nE 'cycle/476|releases/docs/2026-06-21/476' "$f"
# Cycle/485 cited (the third bootstrap-δ run; this scaffold's predecessor;
# γ-side bonus: not required by cnos#486 but recommended for completeness
# — α surfaces in §Design whether to include or not):
grep -nE 'cycle/485|releases/docs/.+/485|\.cdd/unreleased/485' "$f"
# "bootstrap-δ" framing surfaces (the empirical pattern is named as the
# observation set, NOT as the destination architecture):
grep -ciE 'bootstrap.δ|bootstrap.delta|bootstrap-δ pattern' "$f"
# The destination framing — "wake-invoked mode is what runs in production
# when cds-dispatch claims a real cell" — is also named, so a future reader
# doesn't conflate the empirical observation with the architectural target:
grep -ciE 'wake.invoked.*production|production.*wake.invoked|bootstrap.δ.*informs.*destination|empirical.*not.*destination' "$f"
```

β verifies empirically: cycles 470 and 476 are cited *as empirical observations* (i.e. "the 3-round class-trap recurrence in cycle/476 demonstrates that without R[N]-token discipline, the wake cannot distinguish a CI mechanism flake from a substantive review iteration"), not as architectural prescriptions. The cycle/485 citation (if α includes it) names it as "the first R0-converge under mechanical-injection discipline; empirical evidence that the per-CI-step audit absorbs the bash-e class-trap class". If α omits cycle/485, β does not surface a finding (cnos#486 names 470 + 476 only; 485 is γ's bonus recommendation).

## 6. α implementation prompt

You are α for `cycle/486` — implementer for cnos#486 (cdd/delta: define dispatch-wake-invoked δ mode amendment; Sub 5B of cnos#467 wake-orchestration wave).

This is a fresh Agent session. You have no prior context with the parent session. **Read this entire scaffold + the cited source files before touching the SKILL.** δ (the parent session) has already created `cycle/486` from `main@950730c7` and pushed `.cdd/unreleased/486/gamma-scaffold.md` (this file). Your scaffold is the full input; β reviews you on the cycle branch state alone.

**Working directory:** `/home/user/cnos`. Confirm `git rev-parse --abbrev-ref HEAD` returns `cycle/486` before acting.

**Git identity:** before any commit, set `git config user.email "alpha@cdd.cnos"` and `git config user.name "alpha"` per `alpha/SKILL.md` §2.6 row 14. Verify with `git config --get user.email`.

**Inputs you load before writing any text:**

1. **This scaffold** (`.cdd/unreleased/486/gamma-scaffold.md`) — your ACs, oracles, scope guardrails, source-of-truth list.
2. **The cnos#486 issue body** — re-read for cross-verification; do NOT trust this scaffold as the only source of truth. Use `mcp__github__issue_read` with method=get, owner=usurobor, repo=cnos, issue_number=486.
3. **The 12 source-of-truth files in §10 of this scaffold** — read each before drafting the amendment.

**Scope (in / out — mirrored from cnos#486 + the operator directive that spawned γ):**

*In scope:*
- Amend `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` with a new section named "Dispatch-wake-invoked mode" (or equivalent — you pick the precise heading; γ recommends but does not pin).
- Define the dispatch-wake input contract (the 5 named inputs).
- Define δ's routing sequence in wake-invoked mode (how γ/α/β are spawned + what each role's context is).
- Define R[N] iteration token discipline (wake-observable; mechanism enumerated; one pinned for v0).
- Define the per-R[N] artifact contract under `.cdd/unreleased/{N}/`.
- Define return tokens (status:review / status:blocked / claim release) mapped to cnos#454 §2.4 + §2.6 lifecycle transitions.
- Name v0 substrate constraints honestly (timeout class, concurrency, single-claim-per-firing) WITHOUT substrate emission.
- Cite cycle/470 + cycle/476 as empirical anchors (cycle/485 optional — recommended).
- Surface in §Design which heading you picked and why; which iteration-token mechanism you pinned (issue comments vs branch commits vs `.cdd/` JSON artifacts vs hybrid) and why.

*Out of scope (DO NOT touch — these are gated by 5C or separate cycles; touching them collapses sub boundaries):*
- Renderer changes (`cn-install-wake` — cnos#476 + cnos#485 landed).
- Production `.github/workflows/cnos-cds-dispatch.yml` rendering / commit (Sub 5C / cnos#487).
- Flipping `cnos.cds/orchestrators/cds-dispatch/wake-provider.json` `activation_state` (Sub 5C).
- Real `protocol:cds` smoke cell (Sub 5C).
- CDR / CDW dispatch wakes (future packages).
- NIM / OpenAI / alternative substrate carriers (future waves).
- β-skill / γ-skill / harness-skill / operator-skill amendments (cycle/485 closeout triage T1, T2, T11–T14 — bundled into a separate cycle).
- Substrate emission anywhere in `delta/SKILL.md` (AC8 hard gate).

**Operator-named guardrails (these are review points for β; honor them in your implementation; surface each in §Design):**

**OG-1 — Substrate neutrality.** Per cnos#486 AC8: the δ skill MUST NOT emit GitHub Actions YAML, `runs-on`, `GITHUB_TOKEN`, `claude-code-action`, workflow syntax, or `${{ }}` interpolations. The substrate-leakage rule (per `wake-provider/SKILL.md §3.3`) applies to cnos.cdd's δ skill the same way it applies to wake-provider declarations: substrate decisions are renderer authority (cnos#476 / cnos#485), not framework authority. Surface as an explicit β review point in §Design + write an α-side AC in your `self-coherence.md` §R0: "grep for substrate tokens in the diff returns 0 new lines except in descriptive 'relationship to substrate' contexts". Verify with: `git diff main...HEAD -- src/packages/cnos.cdd/skills/cdd/delta/SKILL.md | grep '^+' | grep -ciE 'github|workflow|yaml|GITHUB_TOKEN|runs-on|claude-code-action|\$\{\{'` returns ≤ a small number (γ-pinned threshold: ≤ 3, all in a single explicit "relationship to substrate" descriptive paragraph; β verifies each match is in that paragraph, not in semantic content).

**OG-2 — Sharp output contract.** Per operator directive: "`delta/SKILL.md` should not become the whole runtime manual. It should define δ's wake-invoked role contract. The renderer, substrate workflow, and live smoke remain #487 / 5C." Resist the temptation to specify production-substrate concerns; resist the temptation to specify renderer behavior; resist the temptation to specify the cell-runtime mechanics that already live elsewhere (γ scaffold authoring is `gamma/SKILL.md`; α implementation is `alpha/SKILL.md`; β review is `beta/SKILL.md`; substrate emission is `wake-provider/SKILL.md` + `wake-template/SKILL.md` per cnos#450 + cnos#470). The amendment names what δ does in wake-invoked mode; it CITES the surrounding role contracts without restating them.

Concretely: if you find yourself writing more than ~50 lines on "what γ does after δ spawns it", you're restating `gamma/SKILL.md`. The amendment should name *which sections of `gamma/SKILL.md` δ invokes* (e.g. "γ executes the Step 3 pre-dispatch scaffold check per `gamma/SKILL.md` §2.5"), not duplicate them. Same for α and β. **Surface the line count for the new section in §Design**; if it exceeds ~300 lines of new content (excluding the citation table), β will surface as a scope-creep finding.

**OG-3 — Bootstrap-δ ≠ wake-invoked δ.** The current bootstrap-δ pattern (γ-interface drives via current Claude Code session; spawns γ/α/β as Agent sub-sessions; γ implicitly violates the "δ dispatches every role" rule because γ-the-driver is also δ-the-orchestrator in the empirical cycles) is the EMPIRICAL CASE that informs the contract — but the wake-invoked mode is what runs in production when cds-dispatch claims a real cell. The δ amendment names what the *production* mode looks like, citing bootstrap-δ runs as the input source but NOT as the destination architecture. The "δ dispatches every role; γ does not spawn α/β" invariant (per operator directive) must surface in the amendment as the production routing sequence's invariant.

**OG-4 — Empirical citation discipline.** When you cite cycle/470, cycle/476, or cycle/485, cite them as empirical observations with a specific finding (e.g. "cycle/470 R1 surfaced a broken relative-link path in α's prompt template — γ-scaffold-side substantive ambiguity") rather than as bare references. The citations are the contract's empirical anchor; β reads them as evidence the contract is grounded, not as decoration.

**Output expectations:**

- **Amended SKILL.md** at `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` with the new section satisfying AC1–AC9.
- **`self-coherence.md`** at `.cdd/unreleased/486/self-coherence.md` with (per `alpha/SKILL.md` §2.5 incremental-write discipline — commit each section separately):
  - **§Gap** — issue, version/mode (`design-and-build`), the new-section heading choice + placement decision.
  - **§Skills** — active Tier 1 (CDD.md + alpha/SKILL.md), Tier 2 (none required — this is a doctrinal amendment, no code/eng surfaces touched), Tier 3 (`delta/SKILL.md` as the target; `dispatch-protocol/SKILL.md` + `wake-provider/SKILL.md` + the four role skills as references; `issue/SKILL.md` for cell-shape understanding).
  - **§ACs** — per-AC verification table (9 rows: AC# / oracle command / actual output / pass-fail).
  - **§Self-check** — did α push ambiguity onto β? Is every claim in the new section backed by a concrete citation? Did α resist the OG-2 scope-creep temptation? Honest accounting of any soft-passes.
  - **§Debt** — known debt (e.g. if you omitted cycle/485 from citations, name it; if you used a substrate token in a descriptive paragraph, name where + why it's the relationship-to-substrate carve-out; if you decided to renumber existing §§ in `delta/SKILL.md`, name what you renumbered and why).
  - **§Design** — your section heading choice + placement decision; the iteration-token mechanism you pinned + why; the routing-sequence shape you chose (numbered list / table / step-named paragraphs); the substrate-agnosticism rationale (which words you considered + rejected; which "relationship to substrate" carve-out paragraph you wrote + why). Plus your three explicit α-side ACs:
    - **α-AC-α1 (OG-1 verification):** grep diff for substrate tokens; result documented.
    - **α-AC-α2 (OG-2 line-count discipline):** new-section line count documented + judged against the ~300-line target.
    - **α-AC-α3 (OG-3 production-routing invariant):** the "δ dispatches every role; γ does not spawn α/β" sentence appears verbatim (or with equivalent semantic content) in the routing sequence section.
  - **§CDD Trace** through step 7 per `alpha/SKILL.md` §2.5.
  - **§Friction notes** — α-side meta-observations for future δ wake-invoked mode iteration (Sub 5C and beyond). If you discover during writing that γ's scaffold under-specified something, name it here.
  - **§R0 review-ready signal** — appended last as a separate commit; e.g. `## R0 | base SHA: 950730c7 | implementation SHA: <SHA-of-last-implementation-commit-BEFORE-this-signal-commit> | ready for β`. Per `alpha/SKILL.md` §2.7 "SHA convention for readiness signal": name the implementation SHA (the last impl commit), NOT the HEAD-at-write-time (which advances).

**Iteration discipline:**

- α uses `§R0`, `§R1`, `§R2`, ... section headers per round in self-coherence.md.
- Each §R[N] documents what β found in R[N-1] (cite finding IDs by reference; do not duplicate β's review prose) and what you fixed (commit SHAs + brief description).
- α does NOT change any §R[N-1] content retrospectively (append-only).
- β iterates by appending `§R[N]` sections to `.cdd/unreleased/486/beta-review.md` with verdict (converge or iterate) + findings (F1, F2, …) numbered per cycle, not per round (so F4 in R2 is the 4th finding overall, not the 4th finding in R2 — this matches the cycle/485 β-review convention).
- Iterate until β's verdict is converge. Per `beta/SKILL.md` Role Rule 4, no "approve with follow-up" — converge requires zero open findings.

**When you're done with R0:** commit + push to `cycle/486` + append the review-ready signal to self-coherence.md as a separate commit + report the head SHA. δ (parent session) reads the signal and dispatches β.

**Test plan (run locally before signaling review-ready):** see §9 of this scaffold. Record each command's exit code + relevant output in `self-coherence.md §ACs`.

## 7. β review prompt

You are β for `cycle/486` — reviewer for cnos#486 (cdd/delta: define dispatch-wake-invoked δ mode amendment; Sub 5B of cnos#467 wake-orchestration wave).

This is a fresh Agent session. You have no prior context. δ (the parent session) dispatches you only after α signals R[N] is review-ready (see α's `.cdd/unreleased/486/self-coherence.md §R[N] review-ready signal`).

**Working directory:** `/home/user/cnos`. The branch is `cycle/486`. Confirm `git rev-parse --abbrev-ref HEAD` returns `cycle/486`.

**Git identity:** before any commit, set `git config user.email "beta@cdd.cnos"` and `git config user.name "beta"` per `beta/SKILL.md` pre-merge gate row 1. Verify with `git config --get user.email`.

**Inputs you load:**

1. **This scaffold** (`.cdd/unreleased/486/gamma-scaffold.md`) — the canonical oracle list and scope guardrails.
2. **α's `.cdd/unreleased/486/self-coherence.md`** — α's per-AC verification + chosen heading + design rationale + any prior §R[N] iterations.
3. **The cnos#486 issue body** — re-read for cross-verification; do not trust the scaffold as the only oracle source. Use `mcp__github__issue_read`.
4. **The 12 source-of-truth files (§10 of this scaffold).**
5. **The cycle branch HEAD diff:** `git diff main...cycle/486 -- src/packages/cnos.cdd/skills/cdd/delta/SKILL.md .cdd/unreleased/486/`.

**Review checklist (per-AC; populate this table in `.cdd/unreleased/486/beta-review.md` §R[N]):**

| AC | Oracle | Pass / Fail | Notes / Finding ID if fail |
|---|---|---|---|
| AC1 | §"wake-invoked mode" section heading present; cnos#483's forward-reference resolves (§5 AC1 oracle) | | |
| AC2 | 5 named inputs enumerated with one-sentence semantics each (§5 AC2 oracle) | | |
| AC3 | discrete step list naming γ/α/β spawn sequence + per-role context; cycle/470 + cycle/476 cited; "δ dispatches every role" invariant present (§5 AC3 oracle + OG-3) | | |
| AC4 | R[N] discipline named; iteration mechanism enumerated + one pinned for v0; empirical anchor citations (§5 AC4 oracle) | | |
| AC5 | per-round artifact set matches cnos#483's `output_contract.artifact_class_taxonomy` (§5 AC5 oracle) | | |
| AC6 | three return tokens defined; mapped to cnos#454 §2.4 lifecycle transitions; status:changes carve-out explicit (§5 AC6 oracle) | | |
| AC7 | v0 substrate constraints named honestly (timeout class, concurrency, single-claim) (§5 AC7 oracle) | | |
| AC8 | **substrate-agnosticism**: `git diff` grep for substrate tokens returns ≤ γ-pinned threshold (3), and every match is in the explicit "relationship to substrate" carve-out paragraph (§5 AC8 oracle + OG-1) | | |
| AC9 | cycle/470 + cycle/476 cited as empirical observations with specific finding language; bootstrap-δ ≠ wake-invoked-δ framing surfaces (§5 AC9 oracle + OG-3) | | |

**Operator-named guardrails (review explicitly; each is a finding if violated):**

- **OG-1 (substrate neutrality; paired with AC8):** Read the diff of `delta/SKILL.md`. Run the grep on the `+` lines:
  ```bash
  git diff main...HEAD -- src/packages/cnos.cdd/skills/cdd/delta/SKILL.md \
    | grep '^+' | grep -ciE 'github|workflow|yaml|GITHUB_TOKEN|runs-on|claude-code-action|\$\{\{'
  ```
  The result MUST be ≤ 3, and every match MUST be in a clearly-delimited "relationship to substrate (descriptive only)" paragraph (β reads each match's context to confirm). Matches in semantic content (code blocks; routing-sequence steps; input-contract definitions; artifact-contract definitions; return-token definitions; constraint enumeration) are findings (D-severity, classification `substrate-leakage` — parallel to `beta/SKILL.md` Role Rule 7 `implementation-contract` classification). Matches in citations of issue numbers (cnos#454, cnos#467, cnos#470, cnos#483, cnos#476, cnos#485, etc.) and SKILL.md path references are acceptable.

- **OG-2 (sharp output contract; paired with α's α-AC-α2):** Read α's §Design line-count claim. Independently verify with `git diff main...HEAD --stat -- src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` and `git diff main...HEAD -- src/packages/cnos.cdd/skills/cdd/delta/SKILL.md | grep '^+' | grep -v '^+++' | wc -l`. If the new-section line count exceeds ~300 lines, surface as a scope-creep finding (C-severity unless content is genuinely role-contract surface; α may justify in §Design — read the justification). The amendment should be the role-contract definition, not the runtime manual.

- **OG-3 (bootstrap-δ ≠ wake-invoked-δ; paired with α-AC-α3):** Read the amendment's routing-sequence section. Confirm the production invariant ("δ dispatches every role; γ does not spawn α/β") surfaces as the production routing sequence's binding rule. The amendment must NOT simply restate the bootstrap-δ pattern as if it were the production architecture. β's test: read the routing sequence and ask "does this contract require the bootstrap-δ session pattern (γ-interface-as-driver spawning Agent sub-sessions), or does it work equally well when invoked by a substrate-firing claude-code-action run that has no parent session?" If the answer is "it requires bootstrap-δ", that is a finding (D-severity, classification `production-routing-leak`).

- **OG-4 (empirical citation discipline; paired with AC9):** For each cycle citation (470, 476, optionally 485), read the citation's context. Each citation MUST name a specific empirical finding (a sentence or clause that summarizes what the cycle observed and how it informs the contract) — NOT just a bare reference. Bare references are findings (C-severity, classification `empirical-citation-weak`).

**Doctrine consistency check (B-severity findings if violated):**

- The amendment cites `cnos.cdd` only in framework contexts (the generic cell-runtime framework; the γ/α/β/δ role contracts). Per cnos#467 doctrine-correction header + PR #480: `cnos.cdd` is NOT the dispatch owner. Any text in the new section claiming "cnos.cdd owns the cds-dispatch wake" or "cnos.cdd's protocol qualifier is `protocol:cdd`" or "cdd-dispatch is the wake" is a finding.
- The amendment cites `cnos.cds` for the concrete software protocol; `cds-dispatch` for the wake; `protocol:cds` for the qualifier. Any drift to `protocol:cdd` / `cdd-dispatch` / `cnos.cdd` (in dispatch contexts) is a finding.

**Non-goal verification:**

Run `git diff main...HEAD --name-only`. Expected: exactly 2 files modified:
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`
- `.cdd/unreleased/486/self-coherence.md` (α-authored) + the scaffold path γ committed first.

If any other file is touched, surface as a finding (the "α MUST NOT touch" list in §4 above is the canonical out-of-scope set). The α-may-touch-optionally row in §4 names "internal cross-reference from `delta/SKILL.md` to wake-provider / dispatch-protocol skills" as in-scope; if α added such cross-references, β verifies they're inside the new section (not scattered through unrelated existing sections).

**NOTE: cycle/486 does NOT touch `.github/workflows/`.** The per-CI-step bash-e audit table (mandatory for cycle/485 because that cycle added `run:` blocks to `install-wake-golden.yml`) **DOES NOT APPLY** to cycle/486. Per γ-scaffold §6.1 of cycle/485's gamma-closeout: "the per-CI-step audit applies only when CI changes; include it only when the cycle touches CI." cycle/486's diff is δ-skill-only + α's self-coherence; no CI surface area. **β: do NOT include a per-CI-step audit table in your review. The audit's absence is correct, not an omission.**

**CI evidence:** there is no canonical green-signal CI workflow for a docs-only / skill-amendment cycle (`install-wake-golden` filters on `src/packages/cnos.cds/orchestrators/**` and `src/packages/cnos.core/commands/install-wake/**` — neither covers `delta/SKILL.md`). The wave-inventory inherited-cap CI failures (I4 / I5 / I6) likely still surface on the PR; β confirms they are not cycle/486-introduced by checking the same checks ran red on cycle/485's PR #488 (per cycle/485's beta-review §CI evidence). If a new red appears on cycle/486's PR that did NOT appear on PR #488, surface as a finding.

**Verdict:** `verdict: converge` or `verdict: iterate`. Per β-closeout cycle/485 §6 (recommendation pulled forward): **do NOT manufacture findings to look thorough.** If the work converges cleanly, write a clean converge verdict.

**Output:** `.cdd/unreleased/486/beta-review.md` with:
- **§R[N]** section header (R0 first time; R1 if α had to iterate).
- The 9-row review checklist table above.
- The four OG verifications (OG-1 through OG-4) each named explicitly with pass/fail.
- The doctrine consistency check + non-goal verification results.
- **Verdict line:** `verdict: converge` OR `verdict: iterate`.
- **Findings** numbered `F1`, `F2`, … (only when verdict = iterate; converge → zero findings; do not pad).

**Verdict semantics:**

- **converge** → α and γ proceed to closeout (α writes `alpha-closeout.md`, β writes `beta-closeout.md`, γ writes `gamma-closeout.md`; PR ships; cnos#467 Sub 5B box ticks).
- **iterate** → α addresses findings, increments to §R[N+1] in self-coherence.md, β re-reviews. γ does NOT re-author the scaffold unless a finding surfaces a scaffold-side ambiguity (in which case δ asks γ to amend).

## 8. Non-goals (mirror of §4 "α MUST NOT touch" + the cnos#486 issue's out-of-scope list, restated explicitly for α and β reference)

This is a **DOCTRINE amendment**, not a runtime implementation. No substrate work, no smoke, no production activation. Specifically:

- **Renderer extensions** — Sub 5A landed (cycle/485 / PR #488 merged at `950730c7`); not in scope for 5B.
- **Production `.github/workflows/cnos-cds-dispatch.yml`** — Sub 5C (cnos#487). Not rendered, not committed, not authored in this cycle.
- **`cnos.cds/orchestrators/cds-dispatch/wake-provider.json` `activation_state` flip** — Sub 5C (cnos#487). Manifest stays `declaration-only` through this cycle.
- **Real `protocol:cds` smoke cell** — Sub 5C (cnos#487).
- **CDR / CDW dispatch wakes** — future packages, separate waves (per cnos#467 §"Out of scope, named separately").
- **NIM / OpenAI / alternative substrate carriers** — separate substrate-replacement wave, post-cnos#467 (per cnos#467 §"Out of scope, longer horizon").
- **Multi-protocol δ invocation** (protocol:cdr / protocol:cdw) — when those concrete protocol packages land their own dispatch wakes (per cnos#486 issue body §"Out of scope, longer horizon").
- **Hot-swapping the agent-execution carrier** (claude-code-action → other) — separate substrate-replacement wave.
- **δ scaling to long-running cells** (>30 min wake firings; cron-staggered continuations) — future enhancement; v0 assumes the wake fires once and the cell completes within one substrate firing (per cnos#486 issue body).
- **β-skill / γ-skill / harness-skill / operator-skill amendments** — cycle/485 closeout triage items T1, T2, T11–T14. Bundle these into a separate β-skill amendment cycle; do NOT smuggle them into cycle/486.
- **Substrate emission of any kind in `delta/SKILL.md`** — AC8 hard gate. Substrate-leakage carve-out is one descriptive paragraph; semantic content remains substrate-agnostic.

## 9. Test plan (α's expected local-verification commands before signaling review-ready)

α runs each of these and records exit codes + output excerpts in `.cdd/unreleased/486/self-coherence.md §ACs`:

```bash
# (0) Confirm cycle branch + git identity:
git rev-parse --abbrev-ref HEAD   # expect: cycle/486
git config --get user.email       # expect: alpha@cdd.cnos

# (1) AC1 — wake-invoked heading present:
grep -nE '^#+.*wake-invoked' src/packages/cnos.cdd/skills/cdd/delta/SKILL.md

# (2) AC2 — 5 named inputs (each grep returns ≥1):
f=src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
for term in 'claimed.issue\|issue.number' \
            'protocol.identifier\|protocol qualifier' \
            'current.main\|main SHA\|head.sha' \
            'wake.run.id\|run id' \
            'package.runtime\|concrete protocol.*skill'; do
  count=$(grep -ciE "$term" "$f")
  echo "AC2 [$term]: $count"
done

# (3) AC3 — routing sequence + cycle citations:
grep -nE '^\s*[0-9]+\.\s+(spawn|dispatch|route|invoke)\s+(γ|gamma|α|alpha|β|beta)' "$f"
grep -nE 'cycle/470|cycle/476' "$f"

# (4) AC4 — R[N] discipline + iteration mechanism:
grep -nE 'R\[N\]|R0|R1' "$f" | head -10
grep -ciE 'wake.observable|surface.*wake|signal.*wake' "$f"

# (5) AC5 — artifact taxonomy match:
for cls in gamma-scaffold self-coherence beta-review alpha-closeout \
           beta-closeout gamma-closeout post-release-assessment; do
  grep -qF "$cls" "$f" && echo "AC5 [$cls]: PASS" || echo "AC5 [$cls]: FAIL"
done

# (6) AC6 — return tokens + lifecycle mapping:
grep -ciE 'status:review' "$f"
grep -ciE 'status:blocked' "$f"
grep -ciE 'status:in-progress.*status:todo|claim release' "$f"
grep -nE 'cnos#454' "$f"

# (7) AC7 — substrate constraints honestly named (without substrate emission):
grep -ciE 'substrate timeout|wake firing.*time|per.firing.*horizon' "$f"
grep -ciE 'concurrency.*serial|per.protocol concurrency' "$f"
grep -ciE 'single.claim.per.firing|one.claim.per.firing' "$f"

# (8) AC8 — substrate-agnosticism (HARD GATE):
n_new=$(git diff main...HEAD -- "$f" | grep '^+' | grep -v '^+++' | \
        grep -ciE 'github|workflow|yaml|GITHUB_TOKEN|runs-on|claude-code-action|\$\{\{' || true)
echo "AC8 new substrate-leakage matches in diff: $n_new"
# Acceptable: $n_new ≤ 3, all in a single explicit "relationship to substrate" carve-out.
# Forbidden: $n_new > 3, or matches in semantic content (β reads context).

# (9) AC9 — empirical citations:
grep -nE 'cycle/470|releases/docs/.+/470' "$f"
grep -nE 'cycle/476|releases/docs/.+/476' "$f"
grep -ciE 'bootstrap.δ|bootstrap.delta' "$f"

# (10) Doctrine consistency — no stale forms:
grep -nE 'protocol:cdd|cdd-dispatch' "$f" \
  | grep -v 'NOT|forbidden\|stale\|superseded\|doctrine.correction' \
  | head -5
# Expect: no hits, OR every hit is in a "stale form prohibited" descriptive context.

# (11) Doctrine consistency — cnos.cdd only in framework contexts:
grep -nE 'cnos\.cdd.*dispatch.*wake|cnos\.cdd.*owns.*dispatch' "$f"
# Expect: no hits. cnos.cdd is the framework; cnos.cds owns cds-dispatch.

# (12) New-section line count (OG-2 sanity check):
git diff main...HEAD -- "$f" | grep '^+' | grep -v '^+++' | wc -l
# Expect: ≤ ~300 (target). If higher, document justification in §Design.
```

β re-runs each of these on cycle/486 HEAD as part of the review.

## 10. Source-of-truth list

All paths are absolute on the local cycle/486 checkout. Base SHA = `950730c7` (main HEAD when cycle/486 branched). Online URLs follow `https://github.com/usurobor/cnos/blob/950730c7/<path>` convention.

| # | Reference | Path / URL | Authority |
|---|---|---|---|
| 1 | cnos#486 (this issue) | https://github.com/usurobor/cnos/issues/486 | The 9 ACs + in-scope/out-of-scope + source-of-truth references. Canonical. |
| 2 | cnos#467 master tracker (doctrine-correction header at top is authoritative) | https://github.com/usurobor/cnos/issues/467 | Wave architecture + sub-issue plan + doctrine-correction (cnos.cdd = framework; cnos.cds = concrete software; cds-dispatch = wake; protocol:cds = qualifier). Stale references in the body's content are SUPERSEDED by the correction header + PR #480. |
| 3 | dispatch-protocol skill (cnos#454 / PR #466) | `/home/user/cnos/src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` | The 3-label selector + serialized claim guard + 4-event lifecycle. δ's wake-invoked mode RECEIVES the claimed cell AFTER this skill's claim sequence (§2.2 step 6 "Launch"). |
| 4 | cds-dispatch wake-provider manifest (cnos#483) | `/home/user/cnos/src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` | The reference dispatch-shape manifest. `activation_state: declaration-only` on main. `output_contract.artifact_class_taxonomy` is the canonical artifact-class set the amendment matches (AC5). |
| 5 | cds-dispatch prompt template (cnos#483) | `/home/user/cnos/src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` | The inlined prompt body the wake renders. Lines ~59-74 forward-reference δ's wake-invoked mode (line 73: "δ's wake-invoked mode is the contract that lands in cnos#467 Sub 5"). The amendment LANDS this forward-reference (AC1). |
| 6 | δ role contract (current — pre-amendment) | `/home/user/cnos/src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` | The skill α amends. Read fully before drafting. Existing sections: §1 outward membrane; §2 inward membrane (dispatch enrichment); §3 override; §4 composition with V; §5 what δ does not do; §6 cross-references; §7 Phase 4 wrap-up; §8 remote-runner delegation. The new section likely lands as §3 (pre-override) or §8/§9 (post-existing-sections) — α decides + documents in §Design. |
| 7 | γ role contract | `/home/user/cnos/src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` | δ routes γ in wake-invoked mode. The amendment names WHICH sections of `gamma/SKILL.md` δ invokes (esp. §2.5 Step 3a branch creation + Step 3b scaffold + dispatch prompts; §2.7 close-out triage). Do NOT restate gamma/SKILL.md content. |
| 8 | α role contract | `/home/user/cnos/src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` | δ routes α in wake-invoked mode. The amendment names WHICH sections (esp. §2.1 dispatch intake; §2.5 self-coherence; §2.6 pre-review gate; §2.7 request review; §2.8 close-out). Do NOT restate. |
| 9 | β role contract | `/home/user/cnos/src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` | δ routes β in wake-invoked mode. The amendment names WHICH sections (esp. Phase map; Pre-merge gate; Role Rules 1-7). Do NOT restate. |
| 10 | CDD framework overview | `/home/user/cnos/src/packages/cnos.cdd/skills/cdd/SKILL.md` | The cell-runtime framework's public-loader entrypoint. Cite as the framework's overview; do NOT amend. |
| 11 | Cell-shaped issue contract | `/home/user/cnos/src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` | The issue-shape δ reads when receiving a claimed cell (input contract; AC2). Skim for the Minimal output pattern section to understand what δ expects in a cell's issue body. |
| 12 | wake-provider skill (cnos#470) — substrate-leakage rule §3.3 + dispatch fields §3.9 + activation-state §3.10 | `/home/user/cnos/src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` | The substrate-agnosticism gate's parent doctrine. Cite §3.3 in the amendment's "relationship to substrate" carve-out paragraph as the doctrinal source for AC8. |
| Empirical anchors |
| E1 | cycle/470 (Sub 2; agent-admin wake-provider manifest) | `/home/user/cnos/.cdd/releases/docs/2026-06-21/470/` | First bootstrap-δ run in this wave. R1+R2 (2 rounds). R1 finding F1: broken relative-link path in α's prompt template — substantive ambiguity, not bash-e mechanics. Cite as the empirical case where γ-scaffold-side substantive precision is necessary; the wake-invoked mode's input contract should pin sha/path invariants to filesystem state (cycle/485's FN-1 verification discipline applies here too). |
| E2 | cycle/476 (Sub 3; cn-install-wake renderer v0) | `/home/user/cnos/.cdd/releases/docs/2026-06-21/476/` | Second bootstrap-δ run. R1+R2+R3 (3 rounds — the heaviest in the wave). F1: missing `set -o pipefail` in CI step. F2 (sibling, unmasked by F1 fix): `grep -c` returns exit 1 on zero matches under `bash -e` — killed AC8 audit step on intended-success input. Led to **cnos#478 mechanical scaffold injection** (per-CI-step audit table format requirement). Cite as the empirical case for R[N] iteration discipline being wake-observable: without that, the wake can't distinguish a CI mechanism flake from a substantive review iteration. |
| E3 | cycle/485 (Sub 5A; renderer extension for dispatch shape) | `/home/user/cnos/.cdd/unreleased/485/` | Third bootstrap-δ run. R0 only (converged first round, zero findings). First R0-converge in the wave; empirical evidence that the cnos#478 mechanical-injection discipline absorbs the cycle/476 class-trap. Recommended (not required by cnos#486) citation — name as "the first R0-converge under mechanical-injection discipline; empirical witness that the per-CI-step audit format absorbs the bash-e class-trap class". |

## 11. Friction notes for future δ dispatch-prompt template

These are γ's meta-observations while scaffolding cycle/486; the parent (δ) session reads them; the amendment α writes consumes them. They are γ-side observations to surface for future δ-scaffold-template evolution, not findings for α to address.

**FN-1. The issue body had enough detail for γ to scaffold without operator clarification.** All 9 ACs are independently testable; the scope guardrails are explicit; cross-references resolve to landed sources. The one place γ had to make a judgment call was AC1's "section name" — cnos#486 says "wake-invoked mode" or "equivalent"; γ pinned "Dispatch-wake-invoked mode" as a recommendation in this scaffold but explicitly left the precise heading choice to α (per the operator-named "α decides, but per γ scaffold's guidance" framing). β can flag this as a finding if the chosen heading is misleading; γ doesn't expect that.

**FN-2. All cited sources are present on main at base sha `950730c7`.** No forward-references except the explicit Sub 5C deferrals, which are correctly named as out-of-scope. The doctrine-correction header at top of cnos#467 is the operative one; the body's stale `cnos.cdd` / `protocol:cdd` / `cdd-dispatch` references are explicitly superseded per PR #480.

**FN-3. The bootstrap-δ → wake-invoked-δ mapping is the central conceptual move.** γ-scaffold §6 OG-3 names this; α-AC-α3 verifies it; β-OG-3 reviews it. The empirical pattern is γ-interface-as-driver spawning Agent sub-sessions (which γ-the-scaffolder is using *right now* to author this file — bootstrap-δ in practice). The destination pattern is δ-as-substrate-firing-agent in a claude-code-action run with no parent session — the branch + .cdd/unreleased/{N}/ tree is the entire input. The amendment names the destination pattern, citing the empirical pattern as input. **This is the load-bearing distinction; if α writes the amendment as a description of bootstrap-δ, the amendment fails OG-3.**

**FN-4. cycle/485's closeouts named "branch-as-shared-state handshake" as a wake-invoked-mode design candidate** (per γ-485 §6.4 + α-485 §7 fourth bullet + β-485 §9 rec #4). Three independent role-closeouts agreed: γ-finishes → δ-dispatches-α via the branch + `.cdd/unreleased/{N}/` tree; α-R[N]-ready → δ-dispatches-β via explicit signal in `self-coherence.md`. **γ recommends but does NOT mandate this shape for the amendment** — α has design latitude per OG-2 (sharp output contract). The cycle/485 closeout's specific recommendation is T15 in the triage table; α can adopt it or refine it. β reviews α's chosen handshake shape against the wake-observable requirement (AC4).

**FN-5. Per-CI-step audit table does NOT apply to this cycle.** Per γ-485 §6.1 (recommendation pulled forward): "include the per-CI-step audit only when the cycle touches CI." cycle/486 touches `delta/SKILL.md` + `.cdd/unreleased/486/self-coherence.md` only — no CI surface area. β's review checklist (§7) explicitly carves this out. α should NOT populate a per-CI-step audit in self-coherence.md (would be scope-creep).

**FN-6. The "verify-cited-sha against filesystem state" discipline (cycle/485 FN-1 chain) applies to this cycle.** γ verified the base main SHA (`950730c7`) by reading the working tree; α inherits + re-verifies in self-coherence §Gap; β re-verifies at PR-review time. This is cycle/485's empirical pattern; cycle/486 continues it. The amendment may want to codify "δ verifies sha-pinned invariants against current filesystem state before dispatching γ" as part of the input-contract section (T4 from cycle/485 triage table); γ recommends this for α to consider, but doesn't mandate.

**FN-7. Sub-section ordering / placement inside `delta/SKILL.md` is α's call.** Three reasonable placements (γ enumerates the trade-offs; α decides + documents in §Design):
- **(a) After §2 inward membrane, before §3 override** — places wake-invoked mode as a parallel-to-override "third mode" of δ behavior (alongside outward + inward membranes). Renumbers existing §3-§8 to §4-§9.
- **(b) At the end as a new §8/§9** — places wake-invoked mode as an extension after all existing content, preserving section numbers. Simpler diff; α likely picks this unless there's a strong organizational reason for (a) or (c).
- **(c) As a §3 sub-section of an existing membrane** — e.g. under §1 outward (since wake-invoked is "what δ does when invoked from outside") or §2 inward (since the wake-invocation IS an inward call). Less clean structurally; γ doesn't recommend.

Recommendation: (b) for simplicity unless α has a strong reason to renumber. Surface decision in §Design.

**FN-8. Iteration-token mechanism choice — γ recommends a hybrid.** Per cnos#486 AC4, α has three candidate mechanisms (issue comments / branch-commit messages / `.cdd/unreleased/{N}/` JSON artifacts). γ recommends — but does NOT mandate — a **hybrid pin** for v0:
- **Branch-state primary** (the cycle branch + .cdd/unreleased/{N}/ tree IS the observable iteration state) — this aligns with the cycle/485 closeouts' "branch-as-shared-state" recommendation.
- **Issue comments secondary** (for human-observable convergence signals: claim, β converge, hard block, race release — these match cnos#454 §2.4 lifecycle transitions exactly).
- **`.cdd/unreleased/{N}/*.R{N}.json` not pinned for v0** — too much machinery; the branch-state + section-headers in self-coherence.md / beta-review.md provide the same observability at lower complexity.

α has latitude to refine or replace this; surface in §Design.

**FN-9. The "δ dispatches every role; γ does not spawn α/β" operator directive is a doctrine point that bootstrap-δ implicitly violates.** Today's γ-interface-as-driver session spawns α and β as Agent sub-sessions. That's empirically how cycles 470/476/485 ran. The amendment names this honestly: bootstrap-δ collapses γ-the-coordinator with δ-the-orchestrator (γ-interface is wearing both hats). The production wake-invoked mode separates them: γ scaffolds and exits; δ dispatches α; α implements and exits; δ dispatches β; β reviews and exits; δ iterates. The cycle branch + `.cdd/unreleased/{N}/` tree carries all the state; no chat-state across spawns. **α surfaces this distinction explicitly in the routing-sequence section** — it's both the load-bearing conceptual move AND the production invariant β reviews for at OG-3.

**FN-10. The amendment's word/line count target is a soft target, not a hard gate.** γ pinned ~300 lines of new content (OG-2 + α-AC-α2). If the amendment genuinely needs more (e.g. because α writes detailed step-by-step routing-sequence prose), α justifies in §Design. β reviews the justification rather than the raw count. The intent is "don't restate `gamma/SKILL.md` / `alpha/SKILL.md` / `beta/SKILL.md` content"; the line count is a heuristic. α-AC-α2 captures the measurement; β-OG-2 reviews the justification.

**FN-11. The closeout triad as the bootstrap-δ norm (cycle/485 question).** cycle/485 was the first bootstrap-δ cycle with the full α/β/γ closeout triad (prior cycles 470/476 shipped α-closeout only). cycle/486 continues the triad pattern per operator directive (the scaffold's task list names γ-closeout authoring as the cycle's terminal γ deliverable). γ does NOT flip this to "the norm" unilaterally; the operator's call. Surface as a friction note for the future δ wake-invoked mode skill amendment (it may want to pin the triad as the default cycle-complete artifact set; or leave it dispatch-prompt-driven as today). The amendment can either:
- **(i) pin the triad** as the wake's converge artifact set (matches cnos#483's `output_contract.artifact_class_taxonomy` directly), OR
- **(ii) leave the triad as scope-dependent** (PRA is already scope-dependent per cnos#483 `artifact_class_notes`; α may extend "PRA is optional" to "the triad's optional surfaces are PRA + γ-closeout's third-role formal report").

γ recommends (i) for clarity. α decides; surface in §Design.

**FN-12. Scaffold word/line count.** Approximately ~9.5K words / ~700 lines (heavier than cycle/485's ~6.5K because cycle/486's amendment surface is doctrinal and needs more explicit OG framing + FN observations). The α + β prompts in §6 and §7 are self-contained per the bootstrap-δ "no chat state" discipline. γ-template-amendment work may want to consider whether the OG-framing sections + FN-template should be hoisted into a shared γ-template skill (per cycle/485 γ FN-6 / γ FN-8 / α §5 / β FN-β-4). Out of scope for this cycle; noted for the parent session to triage forward.
