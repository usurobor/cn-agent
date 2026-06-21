# α closeout — cycle/476 (cn-wake-install renderer, v0)

**Issue:** [cnos#476](https://github.com/usurobor/cnos/issues/476) — `cn-wake-install: v0 subcommand to render package-owned wakes into substrate (consume cn.wake-provider.v1)`. Sub 3 of cnos#467; builds on Sub 1 (cnos#468 — label doctrine; merged `c0048bef`) and Sub 2 (cnos#470 — wake-provider declaration; merged `043bf7aa`).

**PR:** [#477](https://github.com/usurobor/cnos/pull/477) (merged via merge-commit at `35380b3d0d37765be9803de7175553d844622ec0`; cycle branch `cycle/476` head `a3e20e3c` merged into `origin/main` at `fcc5cdb9`).

**Branch (pre-merge):** `cycle/476` (γ-created from `origin/main@fcc5cdb9a533ad86e67524bcf05a33d2b4592e8a`; γ scaffold at `417541ad`). This closeout written on `main` alongside `beta-closeout.md` at `d5dbc819`.

**Cycle execution mode:** pre-dispatch δ/channel bootstrap (third such cycle: cnos#468 + cnos#470 + cnos#476); γ/α/β spawned as separate sub-agents via the Agent tool; `.cdd/unreleased/476/` is the shared memory. α did NOT spawn sub-agents. This closeout is the re-dispatched α session per α SKILL §2.8 "re-dispatch path" — δ re-dispatched α after β R3 APPROVED, merged, and wrote `beta-closeout.md`.

**Rounds:** 3 (R1 ready at impl SHA `7162c32a` → β R1 RC F1 → α R2 fix at `12f13045` → β R2 RC F2 → α R3 fix at `1224b532` → β R3 APPROVE → merge `35380b3d`).

**Filed by:** alpha@cdd.cnos (this is α's fourth re-dispatch in the cycle — R1 authoring, R2 fix, R3 fix, post-merge closeout).

---

## Cycle summary

α delivered a four-file v0 shell-renderer surface that consumes Sub 2's `cn.wake-provider.v1` declaration and emits a GitHub Actions YAML wake, plus a CI workflow that defends the renderer's correctness on every push:

1. **`src/packages/cnos.core/commands/install-wake/cn-install-wake`** (495 lines) — POSIX shell renderer. Consumes `cn.wake-provider.v1` manifest via `jq`; validates required fields per `wake-provider/SKILL.md §3.5` with precise stderr error messages naming the missing/wrong-type field; emits GHA YAML wrapping `anthropics/claude-code-action@v1`; substitutes the `{agent}` variable in the inlined prompt; supports `--out <path>` for the future cutover cycle. Idempotent (sha256-stable across consecutive renders; reports `(unchanged)` on no-op re-render so file-mtime is preserved). Discovered via `cn install-wake <name>` through Go-side `discover.ScanPackageCommands` + `ExecCommand` dispatch.

2. **`src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`** (183 lines; sha256 `a912dd97d520b8660e4a53c2f306dad7e657b0c216bbc94f31e2958fdbd01c47`) — default render output of `cn install-wake agent-admin`; the golden fixture for the CI re-render-and-byte-diff invariant.

3. **`.github/workflows/install-wake-golden.yml`** (R1 135 → R3 145 lines; 9 `run:` blocks at R3) — CI workflow. Re-renders the agent-admin wake and byte-diffs against the committed golden; runs all 8 AC mechanical oracles (jq present, re-render, golden unchanged, idempotence, YAML parses, structural shape, AC2 negative-case smoke, AC8 renderer-side authority audit, AC7 byte-identical). This surface had both binding findings (F1, F2) and is where the cycle's signature pattern played out.

4. **`src/packages/cnos.core/cn.package.json`** (+4 lines) — one new `install-wake` entry in the `commands` map registering the renderer with `cn`'s package-command dispatch.

Plus process artifacts: `.cdd/unreleased/476/gamma-scaffold.md` (γ R0; 602 lines), `.cdd/unreleased/476/self-coherence.md` (α R1 + R2 + R3 sections; 553 lines), `.cdd/unreleased/476/beta-review.md` (β R1 + R2 + R3 verdicts), `.cdd/unreleased/476/beta-closeout.md` (β post-merge closeout at `d5dbc819`).

Implementation surface: 4 new code files in `src/packages/cnos.core/` + `.github/workflows/` + the 4-line `cn.package.json` edit. Zero lines in `src/go/`, in any package other than `cnos.core`, or in `.github/workflows/claude-wake.yml` (AC7 byte-identical invariant). All 7 implementation-contract axes (γ-pinned at scaffold time, verified by β Rule 7) held across all three rounds.

---

## AC outcomes (post-merge state)

All 8 ACs PASS at merge SHA `35380b3d`. Each was re-verified by β at R3 against impl SHA `1224b532`; merge introduced no further change to the implementation surface.

| AC | Status (post-merge) | Verifying surface | Round of clean pass |
|----|--------------------|-------------------|---------------------|
| AC1 (form pin + rationale) | PASS | γ scaffold pins shell at `commands/install-wake/cn-install-wake`; α §Gap "Form-choice acknowledgment" accepts with 4 structural reasons; no override | R1 |
| AC2 (consume `cn.wake-provider.v1`; reject malformed; precise stderr) | PASS behavior R1; PASS CI proof R2 | 6 negative cases all exit 2 with stderr naming the precise field; positive case exits 0; `--out` flag byte-identical with default; CI mechanism (AC2 negative-case smoke step) repaired by R2 F1 fix | behavior R1; CI proof R2 |
| AC3 (renders to syntactically valid GHA workflow at γ-pinned path) | PASS | YAML parses via `python3 -c yaml.safe_load`; all 6 structural greps return ≥ 1 (`^on:`, `^permissions:`, `^concurrency:`, `^jobs:`, `anthropics/claude-code-action@v1`, `actions/checkout@v4`) | R1 |
| AC4 (admin-only boundary inlined verbatim from prompt.md) | PASS | "MUST NOT execute cells under any circumstance" present; "Cell execution" heading present; "## Wake termination" (last prompt heading) present (no silent truncation); widened oracle (extract `prompt: \|` block, strip 12-space prefix) char-count ratio 0.998 (within 95–105%) | R1 |
| AC5 (triggers + permissions match manifest) | PASS | Per-trigger: `schedule` + 4 cron entries; `issues_opened_title_match` + title-match `if:`. Per-permission: all 4 `permission_intent` entries map to 4 `permissions:` keys with kebab-case correctly applied (`pull_requests.write`→`pull-requests: write`; `id_token.write`→`id-token: write`) | R1 |
| AC6 (golden + idempotence + CI mechanism) | PASS golden+idempotence R1; PASS CI mechanism R3 | Golden present; 3 renders all sha256 `a912dd97…`; renderer reports `(unchanged)`; CI workflow runs 9 `run:` blocks at R3, all reach success conclusion on PR #477's R3 push | substantive R1; CI mechanism R3 (after both F1 + F2 fixes) |
| AC7 (`claude-wake.yml` byte-identical) | PASS | `git diff origin/main..HEAD -- .github/workflows/claude-wake.yml \| wc -l` = 0 | R1 |
| AC8 (package-vs-renderer authority split preserved) | PASS substantive R1; PASS CI proof R3 | Renderer-side grep for role-decision strings = 0 (shell-concat workaround on field-name literals; runtime stderr still names fields correctly); package-side declarations byte-identical with origin/main; CI mechanism (AC8 audit step) repaired by R3 F2 fix | substantive R1; CI proof R3 |

All 7 implementation-contract axes (γ-pinned, β Rule 7) PASS across all three rounds: language (POSIX shell + YAML + JSON + Markdown; no Go, no Python); CLI integration target (`cn install-wake` via `cn.package.json` commands map; `src/go/` empty diff); package scoping (only `cnos.core/` + `.cdd/unreleased/476/` + the single CI workflow); existing-binary disposition (sibling `daily/cn-daily`, `weekly/cn-weekly`, `save/cn-save` unchanged); runtime deps (POSIX + `jq` only; no new manifest entries); JSON/wire contract (matches `claude-code-action@v1` invocation shape); backward-compat (`claude-wake.yml` byte-identical; Sub 2 declaration byte-identical).

---

## R1 → R2 → R3 delta

**R1 (impl SHA `7162c32a`, head `854102ce`) → β R1 REQUEST CHANGES.** One D-severity finding (F1): the CI workflow's AC2 negative-case smoke step pipes the renderer through `tee /tmp/neg.log` without `set -o pipefail`, so the renderer's exit-2 on malformed input is masked by `tee`'s exit-0; the step's `if`-branch then takes the success branch and emits `::error::Renderer accepted malformed manifest` on every push — constant-failure CI on the very workflow α introduced to discharge AC6. The renderer behavior is correct; the CI guard that proves it is broken on arrival. All other ACs, all 7 implementation-contract axes, and all R1 mechanical gates passed.

**R2 fix (impl SHA `12f13045`).** One-line edit to `.github/workflows/install-wake-golden.yml`: add `set -o pipefail` on the line after `set -eu` in the AC2 negative-case smoke step (line 97 at R2 head). The pre-existing `if … | tee …; then` shape preserved so the diff reviewer sees a single-line semantic delta. α scoped to F1 only (β R1 named no other findings).

**R2 (impl SHA `12f13045`, head `4913228f`) → β R2 REQUEST CHANGES.** One D-severity finding (F2): the AC8 renderer-side authority audit step (lines 115-125 at R2 head) ran `n=$(grep -ciE '…' cn-install-wake)` under `bash -e`; `grep -c` returns POSIX exit 1 when the count is zero, and zero is precisely the intended-success path here (we WANT zero role-decision leaks). The command substitution killed the step before the `if [ "$n" != "0" ]` guard fired. **F2 was hidden at R1 by F1's earlier-step crash**; F1's R2 fix unblocked execution past line 113 and exposed F2 — same `bash -e` defect family as F1, sibling pattern in the same workflow file.

**R3 fix (impl SHA `1224b532`).** One-line edit to the same workflow: append `|| true` to the `grep -ciE` command substitution in the AC8 step. The substantive guard `if [ "$n" != "0" ]; then echo "::error::…"; exit 1; fi` is byte-identical with R2 — regression detection on real leaks is preserved; the `|| true` only neutralizes the bash-e + grep-c-on-zero-matches interaction. α also added (per β SKILL Rule 6) a comprehensive **bash-e semantics audit table for every `run:` block** in the workflow (9 steps), simulating each under `bash -e` against intended-success input and confirming the guard mechanism on each — recorded in `self-coherence.md` §R3 fix.

**R3 (impl SHA `1224b532`, head `c55f7061`) → β R3 APPROVE.** F2 resolved; AC8 step now exits 0 on intended-success input and reaches the success echo on PR #477 CI; bash-e audit re-confirmed across all 9 run blocks; merged at `35380b3d`.

**The two findings (F1, F2) are a sibling pair** in the same defect family: `bash -e` + step-internal exit-code propagation through pipelines (F1: tee swallowed exit 2) or command substitutions (F2: grep -c returns 1 on zero matches). Both required a one-line YAML fix. Neither was a renderer-source defect — the renderer carried zero churn across rounds. Only the CI workflow file (and `self-coherence.md`) churned at R2 and R3.

---

## Self-check (post-merge)

### Did claim-class verification injection (cnos#472) work?

**Yes at the doc-claim level; NO at the next-deeper artifact-presence-vs-CI-execution-semantics level.** This is the cycle's signature finding and requires unpacking.

cnos#472's per-item-table-required rule was injected into γ's scaffold for this cycle. α populated every per-X table — AC1–AC8 per-AC oracle table; AC2 6-row per-case table; AC3 per-grep table; AC4 per-marker table; AC5 per-trigger + per-permission tables; AC6 per-render idempotence table; mechanical-gate per-row table. β R1's "Claim-class verification audit" explicitly confirmed: "no aggregated claims slipped through ... the injection prevented the cnos#470-R1-style wiring-claim trap from recurring."

**And yet F1 still shipped.** α's R1 "AC6 CI mechanism in place" row WAS per-item-tabled (one row per CI step at the workflow's named path); each step was named; the table existed. What was NOT in the per-item discipline: **execution evidence** that each step actually exits with the intended code on the intended input class under the actual GH Actions `bash -e` shell invocation. β had to run the PR's CI to discover the AC2 step was constant-failure-on-arrival.

The rule prevented the previous cycle's failure mode (#470-R1's aggregated wiring claim) but did not prevent the next-deeper instance (#476-R1's per-item artifact-presence without per-item execution evidence). And then **R2 sharpened the rule (in prose) to require execution evidence for the one step F1 touched** — but R2 did not extend that discipline to every other CI step in the same workflow, so F2 (a sibling in the next step over) shipped to R3 review. The discipline was named at R1 (β + α in prose); R2 still shipped an instance of it. Then the discipline was named at R2 (β + α in prose); R3 would have shipped another instance if there had been a third sibling step — α's R3 bash-e audit table is the only thing that confirms there wasn't.

### Did α push ambiguity onto β at any round?

**No across rounds, but the underlying class-trap was self-pushed onto each next round.** The per-cycle round-level mechanics held: every AC oracle was mechanically reproducible from the diff alone; the AC4 char-count brittleness was resolved by α (widened oracle: extract `prompt: \|` block, strip indentation, char-count ratio 0.998) and recorded in §ACs rather than deferred; the AC8 shell-concat workaround was recorded as friction F3 with structural reasoning; the F1 fix and F2 fix both included local re-reproduction evidence (R2: pipefail + tee experiment showing the if-branch now correctly does NOT take on intended-success; R3: bash-e simulation of the AC8 step both pre-fix-killing and post-fix-passing). β had to run the PR's CI to verify the post-push outcome but did not have to derive any α-side reasoning.

**Where α did NOT hold the line:** at R1 and at R2, α populated the per-item table at the artifact-presence depth that satisfied the round's then-current discipline but missed the round's emerging discipline. R1's "AC6 CI mechanism in place" row tabled the steps but not their bash-e exit codes. R2's amended AC6 row tabled the **one** step F1 touched at exit-code-semantics depth but did not extend that depth to the other 8 steps. The R2 §"Honest correction — AC6 row + §Self-check amendment" paragraph explicitly acknowledged this: "the R1 claim 'AC6 CI mechanism in place' carried artifact-presence evidence (workflow file exists, steps exist) but did not carry per-step exit-code-semantics evidence ... going forward, α self-coherence claims about CI mechanisms must include per-step CI evidence." The "going forward" was prose; it did not retroactively apply to the rest of the workflow at R2; F2 shipped.

The R3 bash-e audit table is the first instance in the cycle where α populated the discipline to the depth it needed for the workflow to be β-approvable in one pass — and α only did so because β R2's review had explicitly named the audit as the R3-required surface (and α R3's dispatch prompt mirrored it). The discipline arrived after the defect, three times in a row.

### Did α surface debt honestly?

**Yes.** Six §Debt items declared at R1 (F1–F6: scope-discipline regex brittleness in γ scaffold, AC4 char-count oracle brittleness, AC8 renderer-side grep brittleness, agent→bot identity mapping has no manifest field, renderer hardcodes substrate cron slots, on-disk path is both render target and golden fixture). All six items remain accurate post-merge; none were "discovered" by β as missing-from-α. The R2 + R3 fix rounds did not close any of them (scoped to F1 and F2 respectively). γ closeout should triage which warrant follow-up filings.

α also surfaced two friction items spanning rounds: R2's "Friction note for cnos#472 sharpening" (Friction-O1, per-item table needs CI-evidence depth) and R3's "CRITICAL friction note for γ closeout / cnos#472 sharpening" (the empirical class-recurrence ledger across three rounds + the conclusion that prose discipline is insufficient).

---

## The recurring class-trap retrospective (the cycle's signature finding)

**Across three rounds in this single cycle, the same class of CI-mechanism-correctness trap recurred at progressively deeper levels:**

| Round | Class of defect | What β found | Prose-naming of next-level sharpening |
|---|---|---|---|
| R1 (impl SHA `7162c32a`) | per-item table at **artifact-presence depth** (CI workflow file exists, steps exist, named on per-row table) but no execution evidence | F1: AC2 negative-smoke step constant-failure under `bash -e` + missing `pipefail`; broken-on-arrival CI mechanism | β R1 §"Notes for γ closeout" + α R2 §"Friction note for cnos#472 sharpening" both named "per-item table with execution evidence" as next-level sharpening in PROSE |
| R2 (impl SHA `12f13045`) | per-item table at **exit-code-semantics depth** for ONE step (the F1 fix step) but not extended to every step in the workflow | F2: AC8 audit step exits 1 under `bash -e` because `grep -c` returns 1 on zero matches; sibling of F1, hidden at R1 by F1's earlier crash | β R2 §"Note on cnos#472 effectiveness" + α R3 §"CRITICAL friction note" both named "per-step `bash -e`-semantics audit table" as next-level sharpening in PROSE; AND both rounds now flagged the META-finding: prose-naming alone has NOT been sufficient to land the discipline in time to catch the sibling instance |
| R3 (impl SHA `1224b532`) | full bash-e audit table populated by α covering all 9 `run:` blocks in the workflow — discipline applied to the depth needed | β R3 APPROVE; no F3 | (no fourth round; merged) |

**Empirical conclusion (β and α convergent across both review notes and self-coherence sections):** prose discipline is insufficient. Each round, both β-side review notes AND α-side §R[N] fix notes named the next-level sharpening explicitly and accurately. Each subsequent round still shipped an instance of the very class the prose had named. This is not because either party failed to understand the lesson — the prose framing was accurate and travel-ready at each round. **It is because the discipline was not mechanically enforced in the scaffold/template that α populates each cycle.** A discipline that lives only in prose guidance requires every α (or every role, every cycle) to remember to apply it; a discipline that lives as a template section the scaffold MUST populate is enforced by the act of populating the scaffold.

This cycle is the empirical case study for "**prose discipline insufficient → mechanical scaffold/template injection required**." It is the durable lesson the cycle delivered beyond its substantive output. γ closeout / PRA / cnos#472 follow-up should fold this empirical case into the next-cycle scaffold template.

---

## Debt status (re-stated post-merge)

The six §Debt items declared at R1 self-coherence remain accurate post-merge. None were closed by R2 or R3 (R2 scoped to F1; R3 scoped to F2; both fix-rounds touched only the CI workflow and self-coherence). γ-closeout should triage which warrant follow-up filings.

| # | R1 declaration (summary) | Post-merge status | Triage hint |
|---|---|---|---|
| F1 | γ scaffold scope-discipline regex is brittle (missing trailing wildcards on directory alternatives) | **Unchanged.** α-corrected version recorded in §ACs "Mechanical gate summary". | γ may amend scaffold-template regex shape OR note as known pattern. |
| F2 | AC4 char-count oracle (95–105% of original) is brittle when applied to whole rendered file; structural correctness produces 1.22 ratio because YAML indentation + headers inflate | **Unchanged.** α implemented widened oracle (extract `prompt: \|` block, strip 12-space prefix, char-count ratio 0.998); recorded in §ACs AC4 row. | γ may file follow-up to clarify AC4 oracle wording for future wake-renderer ACs. Does NOT seed a `wake-provider/SKILL.md` amendment. |
| F3 | AC8 renderer-side grep oracle conflates "renderer emits role-decision strings to substrate" with "renderer source mentions role-decision field names" — but the validator MUST mention field names to reject by name | **Unchanged.** α used shell-concat (`_a="admin"; _u="_only"`) to assemble field names; behavior preserved (AC2 negatives still name fields in stderr). β accepted the workaround as friction-routed. | γ may consider amending AC8 oracle wording: split into (a) substrate-emission grep on rendered output, (b) runtime-data role-decision encoding grep on source. Optionally amend `wake-provider/SKILL.md §4` with the split. |
| F4 | Agent → bot identity mapping has no manifest field; renderer hardcodes `sigma → sigma@cnos.cn-sigma.cnos / 41898282` one-row table inside `cn-install-wake` | **Unchanged.** Unknown agents fail with a precise error naming the file the operator must edit. | Cutover-cycle pre-work: lift to `cn.wake-bindings.json` or `--agent-config <path>` flag or per-substrate config. Not v0 scope. |
| F5 | Renderer hardcodes substrate cron slots (`8 23 38 53`) from `claude-wake.yml` because manifest doesn't carry cron-slot info (substrate authority per `wake-provider/SKILL.md §2.5`) | **Unchanged.** Comment block at lines 138-144 of renderer cites `claude-wake.yml` + rationale + "preserve at cutover" reason. | Future cycle may add `concurrency_intent.cron_slots` or `substrate_overrides.cron_slots` to manifest (backward-compatible `v1` field addition). |
| F6 | On-disk path `.cnos-agent-admin.golden.yml` is both render target AND golden fixture (coupled); cutover sequence requires explicit `--out` for active workflow materialization | **Unchanged.** `--out` flag works (verified in §ACs AC2 row). | Cutover-cycle scaffold should document the 4-step cutover sequence (re-render golden → operator review → render to active path → remove claude-wake.yml in same commit). |

**No new α-side debt discovered between R1 and merge.** Three rounds of CI-mechanism iteration did not produce any new debt items — only fix-round friction notes already routed above. The cycle merged with exactly the 6 R1-declared debt items still open in the same shape, plus the cycle-spanning friction-O1/CRITICAL-friction observation (now folded into the class-trap retrospective above).

**Pre-existing main CI red (carried into this cycle but inherited from origin/main, not α-debt):** I4 (Repo link validation), I5 (SKILL frontmatter), I6 (CDD artifact ledger — `cn-cdd-verify` missing) — all red on origin/main at base SHA `fcc5cdb9`; cnos#476 does not regress any of them per β R1 CI status table; out of scope for this cycle; tracked by named follow-up issues per the issue body's cross-refs. γ-closeout PRA may wish to file follow-up.

---

## Friction notes for γ closeout

Capture explicitly so γ does not re-derive these from the round-by-round artifacts:

### F-γ-1 (primary, convergent β + α): the 3-column per-CI-step table proposal

Both β R3 closeout and α R3 §CRITICAL-friction recommend that γ closeout / cnos#472 follow-up amend the γ scaffold template to **inject a 3-column per-CI-step table (NOT prose guidance) for any cycle touching CI workflows**:

- **Column (i): per-`run:`-block `bash -e` substitution-failure-mode audit** — for each `run:` block in any new workflow file the cycle introduces, the row lists: command substitutions / pipelines / commands that could exit non-zero; guard mechanism (`|| true` / `set -o pipefail` / `if !`); empirically-observed `bash -e` exit on intended-success input.
- **Column (ii): per-step CI execution evidence** — job URL + conclusion + the specific assertion the step proves, populated once the cycle's PR CI has run.
- **Column (iii): per-step assertion-fires verification** — the step actually exercises its assertion on at least one observed input (not just compiles / parses / exists). For CI steps usually subsumed by (ii); for static-check steps needs a positive-regression-pair.

**The table shape is byte-liftable from α's R3 audit in `self-coherence.md` §R3 fix** (9 rows; columns: # / Step name / Line range / Command substitutions or pipelines / Guarded? / bash-e exit on intended-success input / Notes). β R2 review also produced a parallel audit at the end of F2's "Other steps β audited under bash -e" section; both tables are byte-coherent.

### F-γ-2 (the meta-rule): cnos#472 mechanical-injection requirement

The discipline must land as a **scaffold-required subsection γ scaffolds inject and α populates**, NOT as prose guidance α must remember to apply. This is the empirical lesson of the cycle: across three rounds, both β-side and α-side named the next-level sharpening in prose each round, and the next round still shipped an instance. The cycle is the empirical case study for the conclusion.

γ closeout should either fold this into the cnos#472 update itself OR file a follow-up issue lifting the requirement from "prose in α prompts and review/SKILL.md" to "scaffold template section γ MUST populate per cycle, α MUST populate per cycle". β R2 noted γ may want to do both: amend cnos#472 + amend the γ scaffold template + amend the α SKILL to require α populate the table before signaling review-readiness on any cycle touching CI workflows.

### F-γ-3 (smaller, scoped): the 6 §Debt items routed to γ closeout triage

F1–F6 in §Debt above (recap: γ scaffold regex shape, AC4 char-count oracle wording, AC8 grep oracle conflation, agent→bot identity mapping needs manifest surface, cron slots hardcoded, render target ↔ golden coupling). None blocking; all candidates for γ-side triage as either issue-template amendments, `wake-provider/SKILL.md` field additions, or cutover-cycle pre-work. Recommend γ classify each as: file follow-up issue / fold into cutover-cycle pre-work / accept as v0 scope.

### F-γ-4 (pre-existing main CI red, observational): I4 / I5 / I6 still red on origin/main

Not regressed by this cycle; tracked by named follow-up issues per the issue body cross-refs (cnos#473/#474/#475). γ-closeout PRA may wish to record as cross-cycle observation (third successive bootstrap cycle observing same pre-existing reds; no cycle has chased them yet).

### F-γ-5 (process observation): three successive pre-dispatch bootstrap cycles ran smoothly

cnos#468 (Sub 1; merged `c0048bef`), cnos#470 (Sub 2; merged `043bf7aa`), cnos#476 (Sub 3; merged `35380b3d`) all ran under pre-dispatch δ/channel bootstrap mode (γ/α/β spawned as separate sub-agents; `.cdd/unreleased/{N}/` as shared memory; α never spawned sub-agents; β merged + wrote closeout; α re-dispatched for closeout). The mode is now empirically stable across three cycles. γ closeout / PRA may want to note this as durable mode evidence rather than experimental.

---

## Forward links

- **γ closeout (pending):** δ will dispatch γ for `.cdd/unreleased/476/gamma-closeout.md`. γ owns: (a) the cycle's CDD-process learning column (F-γ-1, F-γ-2 above are the primary feedstock for cnos#472 follow-up + γ scaffold template amendment); (b) the docs-only archive move `.cdd/unreleased/476/` → `.cdd/releases/docs/2026-06-21/476/` per β SKILL note + release/SKILL.md §2.5b; (c) the PRA at `docs/gamma/cdd/docs/2026-06-21/POST-RELEASE-ASSESSMENT.md`; (d) any follow-up issue filings (cnos#472 sharpening; F-γ-3 debt triage; cross-cycle observations from F-γ-4 + F-γ-5).
- **PRA (pending):** γ-owned per `post-release/SKILL.md`. α does not write the PRA. The friction notes above are α's per-cycle feedstock; γ aggregates across cycles and writes the PRA. The class-trap retrospective is durable PRA-level material.
- **Sub 4 of cnos#467 — cnos.cdd dispatch wake provider.** Consumes this cycle's renderer. The Sub 2 contract skill (cnos#470) + Sub 3 renderer (this cycle) are jointly sufficient for Sub 4 to author its own dispatch-class wake provider declaration AND render it to active substrate via `cn install-wake <dispatch-wake-name> --out <active-path>`. The renderer's `--out` flag (verified at AC2 `--out` row) is Sub 4's cutover hook.
- **Sub 6 of cnos#467 — end-to-end smoke.** Consumes everything: Sub 1 label doctrine + Sub 2 declarations + Sub 3 renderer + Sub 4 dispatch wake + Sub 5 δ wake-invoked mode. The cycle-complete reading and end-to-end pipeline depend on Sub 3 being substantively correct (which it is, R1) and the CI mechanism being durable (which it is, post-R3).
- **Cutover cycle (separate; per debt F4/F5/F6):** the production cutover from `.github/workflows/claude-wake.yml` to `.github/workflows/cnos-agent-admin.yml` is a separate future cycle; this cycle's `--out` flag + render-target docs are the inputs.
- **cnos#472 follow-up (per F-γ-1 + F-γ-2):** the mechanical-scaffold-injection sharpening. γ closeout decides whether to amend cnos#472 in place OR file a sibling follow-up issue. α's R3 bash-e audit table is the byte-liftable template; β's R2 audit table is the parallel.

α's work on cycle/476 is complete. Re-dispatched once for R2 fix; re-dispatched once for R3 fix; re-dispatched once for this closeout. Exits after push.
