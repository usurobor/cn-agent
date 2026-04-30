# α close-out — #301

## Cycle summary

| | |
|---|---|
| Issue | #301 — *infra(ci): CUE-based skill frontmatter validation in coherence CI*. P2. Filed by γ; substantial cycle, single CTB v0.1 surface (skill-module frontmatter shape). Closes on `b483f36`. |
| Release | `3.63.0` (`6300081`). Minor bump (new CI gate I5; new authority surfaces under `schemas/`; no runtime/binary/package-source change). Per β O9, tag push deferred to δ (env 403). |
| Branch | `claude/cnos-alpha-tier3-skills-MuE2P` (legacy `{agent}/{slug}-{rand}` shape — pre-#287 `cycle/{N}` convention). Branch shape was set by the harness at α dispatch; α did not select it. |
| Dispatch | γ → α: project `cnos`, issue #301, Tier 3 skills `cnos.core/skills/design`, `eng/tool`, `eng/test`, `eng/ux-cli`, `eng/process-economics` (conditional, did not fire as a design question). γ named work-shape as substantial in the dispatch comment (`4347101265`). |
| Work shape | Substantial. Multi-surface: CUE schema + extraction script + CI job + exception list + fixtures + schemas/README.md + 46 SKILL.md frontmatter backfills + 1 YAML normalisation pair + 1 stale-`calls` removal + 1 in-touched-file prose rename. 60 files in the merge; +2,612 / −58 net (rough). |
| Tier 2 skills applied | `eng/document` (durable docs in `schemas/README.md`, machine-readable diagnostics format), `eng/code` (script soundness, `set -euo pipefail`, no silent fallback). |
| Tier 3 skills applied | `cnos.core/skills/design` (script-vs-schema-vs-filesystem source-of-truth split; CUE owns shape/type/enum, shell owns discovery/extraction/exception/calls-fs; per AC2 surface boundary). `eng/tool` (fail-fast prereq checks for `cue` and `jq`, NO_COLOR + non-TTY honored, exit codes 0/1/2, idempotent). `eng/test` (positive/negative proof via fixtures including a calls-target broken-by-construction fixture; explicit invariants). `eng/ux-cli` (one finding per line in `path :: field :: rule :: reason :: fix` shape, ✓/✗ symbols, exit-code-2 for prereq distinct from exit-1 for validation failure). |
| MCA selection | The issue named 9 ACs with explicit field stratification (hard-gate / spec-required-exception-backed / optional-defaulted / reserved / unknown-package-local). MCA scope was largely pre-decided by γ; α's job was to land the named contract coherently across schema, script, exception list, fixtures, CI job, schemas docs, plus the hard-gate frontmatter backfills and YAML normalisation surfaced *during* implementation. The L6 leverage move is the new CI gate plus the explicit debt ledger — turning frontmatter drift from a silent class into a CI-failure class with a per-file shrinkable ledger. |
| Review rounds | **2 numbered rounds + 1 β-side re-evaluation pass** (β O1 — round-2 Pass 2 against fresh `origin/main = 9d6a0fa` after cycle #287 v3.62.0 advanced main during the review window). R1 REQUEST CHANGES + 3 findings (1 C-judgment, 2 A-mechanical). R2 APPROVED. R2 Pass 2 confirmed APPROVAL stands; verdict unchanged because F1–F3 are content-level and immune to base advance. |
| Findings against α | 3 total: F1 (C, judgment, peer-enumeration miss on in-touched-file body prose), F2 (A, mechanical, intra-doc numeric drift "45" three sites), F3 (A, mechanical, recursively self-stale "head SHA = readiness commit SHA" convention). Resolved on-branch in two commits: `171188e` (the fixes) + `55642db` (the appendix to self-coherence.md). |
| Mechanical ratio | 2/3 = 67%. N=3 well below `eng/process-economics` ten-finding floor; not an automatic process-debt trigger. |
| Local verification | At every review-round head: `./tools/validate-skill-frontmatter.sh` → 56/56 SKILL.md no findings; `./tools/validate-skill-frontmatter.sh --self-test` → 3 positive + 4 negative behave as expected; delimiter survey across all 56 SKILL.md → 0 malformed. CUE `v0.13.2`, `jq 1.7`. |
| Concurrent main activity | Cycle #287 (v3.62.0) shipped to `main` *during* β's review window, advancing `origin/main` from `a8e67b7` → `9d6a0fa`. File-overlap with the cycle branch was a single line (`cdd/gamma/SKILL.md`, disjoint regions); auto-merge clean. β re-evaluated against fresh main in Pass 2; F1–F3 stable, no new findings. After β's Pass 2 verdict, sigma's `4a0f678` clarified merge as β authority (not δ); β proceeded to merge under that doctrine (`b483f36`). |

## Findings (α-side observations)

### O1 — Peer-enumeration miss inside a single file: frontmatter vs body

**What happened.** During implementation α discovered that `release/SKILL.md` declared `calls: - writing` as a static call target. The path "writing" did not resolve under the package skill root; α removed the entry. β's R1 F1 then named two body-prose lines in the same file (L103, L217) that still referenced "the writing skill". The frontmatter and the body were peers; α enumerated only the frontmatter.

**Pattern.** A single file is its own peer set when it carries the same fact in two surfaces (frontmatter + body, table + prose, header + footer). The intra-doc grep rule in `cdd/alpha/SKILL.md` §2.3 — "grep every occurrence of the wrong value AND the corrected value" — applies *within* a file as well as across files. α's peer enumeration treated frontmatter as the only surface; the rule's full statement covers prose surfaces in the same file.

### O2 — Intra-doc numeric drift "45 → 43"

**What happened.** α drafted `self-coherence.md` with `45 entries` for `schemas/skill-exceptions.json` before the script-driven generation produced the final count of 43. Three sites in the artifact carried the stale value (L27, L53, L144). β's R1 F2 named all three.

**Pattern.** Any numeric or value derived from generated data must be re-derived against the actual generated data immediately before the readiness signal, not at draft time. The risk class is the same as #266 F3-bis (one document carrying a count across multiple sites; fix-only-one is a closure overclaim). The rule is already in `cdd/alpha/SKILL.md` §2.3 ("Intra-doc repetition"); α drafted before applying it.

### O3 — Recursively self-stale "head SHA = readiness commit SHA" convention

**What happened.** α's first-round readiness signal named `head SHA = the SHA of the commit that lands this readiness file`. The SHA cannot be known before the commit exists. α refreshed it once (`8adfd44 → 6e5ce21`) and the act of refreshing advanced HEAD by one commit again (`→ ed5f218`), leaving the new value still one commit behind. β's R1 F3 surfaced the loop. α's R2 fix renamed the field to "implementation SHA" (a stable earlier commit) and added a "SHA convention" paragraph documenting that branch HEAD is carried by the polling protocol, not by the inline value. β's β-side O3 names the same root.

**Pattern.** An artifact that tries to name "this commit" by SHA inline is recursively self-stale by construction. Two stable conventions exist: name an earlier commit whose SHA is fixed (implementation SHA), or omit the SHA and let the polling protocol carry HEAD. Either makes the artifact self-consistent on first write. The pattern recurs in any artifact that includes its own commit reference.

### O4 — Non-goal vs AC strict gate: a forced reading

**What happened.** The issue body named "Rewriting non-conformant skills to add missing fields (separate work)" as a non-goal. AC1 named `name`, `description`, `governing_question`, `triggers`, `scope` as hard-gate (no exception path). AC4 named "CI passes on main with exceptions in place" *and* "Hard-gate fields cannot be excepted; only spec-required-but-exception-backed fields are eligible." Of 56 SKILL.md, 30 lacked `description`, 33 lacked `governing_question`, 43 lacked `scope`. The three constraints together strictly imply: every skill must carry all five hard-gate fields on green-CI day. The non-goal could only be reconciled by reading "rewriting" narrowly: do *not* add `inputs`/`outputs`/`artifact_class`/`kata_surface` (the exception-backed four), but *do* add the hard-gate ones. α resolved on this reading and authored 106 hard-gate field additions across 43 SKILL.md.

**Pattern.** When a non-goal and an AC strict gate appear contradictory, the consistent reading is usually that the non-goal scopes a *subset* of the touched surface and the AC carries the harder constraint. The implementer must explicitly state the resolved reading; silent resolution leaves the operator with a non-goal that looks broken in retrospect. α did not state the resolution back to γ; the operator could only verify it from the diff.

### O5 — Strict parser surfaces latent YAML ambiguity

**What happened.** `cdd/design/SKILL.md` had `description: Produce the design artifact for a substantial change: incoherence, ...` (unquoted colon-space mid-value); `cdd/gamma/SKILL.md` had `- delta gate results (observable via git: tags, branch state)` (the same shape inside parens). PyYAML rejected both with `mapping values are not allowed`; CUE rejected both with `mapping values are not allowed in this context`. The frontmatter had been on `main` for cycles without anything noticing because no checked-in tool parsed it strictly. α normalised both lines (quoted them) in scope.

**Pattern.** Adding a strict parser to a tree of unverified YAML surfaces a small population of latent-ambiguity findings on first run. They are not *introduced* by the cycle; they are *exposed* by it. The strict tool earns its keep by making them visible, even when the immediate cost is a one-line normalisation per finding. Treat them as frontmatter-format debt, fix in scope (or except), and document that the validator now catches the class going forward.

### O6 — Cross-package `calls` is a missing v0 primitive, not a per-skill defect

**What happened.** `release/SKILL.md` declared `calls: - writing`. No `writing/SKILL.md` exists under cdd's package skill root; the closest valid target is `cnos.core/skills/write/SKILL.md`, in a different package. LANGUAGE-SPEC §2.4.1 specifies `calls` as intra-package, resolved against the package skill root; cross-package edges have no notation in v0. α removed the entry. The prose body still references the write skill (out-of-band coordination, which is what produced F1).

**Pattern.** When a static `calls` entry would only resolve cross-package, the entry's existence is evidence that the language is missing a primitive (cross-package call edges), not that the entry is locally wrong. The v0 fix is removal because v0 CI must be green; the durable debt is "v0.x or v1 LANGUAGE-SPEC §2.4 should declare cross-package call notation." That debt is outside #301's scope but recurs anywhere a skill genuinely composes across packages.

### O7 — Compound sub-skill names with slashes

**What happened.** Cycle #304 (merged the day before this α dispatch) introduced `cnos.cdd/skills/cdd/review/contract/SKILL.md` and `.../review/implementation/SKILL.md` with `name: review/contract` and `name: review/implementation`. α's first schema draft used `^[a-z][a-z0-9-]*$` for `name` (the conservative reading of "identifier"). The schema rejected both. α relaxed to `^[a-z][a-z0-9_/-]*$`. LANGUAGE-SPEC §2.1 describes `name` only as "identifier"; no explicit grammar.

**Pattern.** When a spec leaves a primitive's grammar open and the practice has produced shapes the strict reading would reject, the schema must either accommodate the practice or force the spec to forbid the shape. The v0 schema chose accommodation; tightening would require renaming `review/contract` and `review/implementation` first.

### O8 — Cycle-dir poll-surface advances at release

**What happened.** α was polling `.cdd/unreleased/301/beta-review.md` for β's round-2 verdict. β had already approved (`acfa0cf`), Pass-2'd (`a477b1d`), merged (`b483f36`), released as 3.63.0 (`6300081`), and written `.cdd/releases/3.63.0/301/beta-closeout.md`. The operator had to surface to α that the cycle directory had moved per `release/SKILL.md` §2.5a. α's polling protocol named only the unreleased path; the release event moved the canonical surface and α's loop did not follow.

**Pattern.** During in-version work the canonical poll path is `.cdd/unreleased/{N}/`; after release it is `.cdd/releases/{version}/{N}/`. α's polling spec should account for the release-time move so the post-merge close-out hand-off doesn't depend on operator nudge. The same gap likely exists in γ's polling for cycles γ tracks across the release boundary.

### O9 — Mass-mechanical SKILL.md backfill landed in one commit with the substantive contract change

**What happened.** α landed all eight surfaces of the cycle (CUE schema, validator script, exceptions JSON, fixtures, CI job + notify aggregation, schemas/README.md, 43 SKILL.md hard-gate backfills, 2 YAML normalisations, 1 stale-`calls` removal) in a single implementation commit `8adfd44`. β's review evaluated 60 files as one diff. The mechanical 43-skill backfill is structurally distinct from the substantive contract surfaces (each backfill adds the same 3 fields with file-specific values; no design judgment per file), and β's commit-shape discipline (split incremental writes per σ's `70ff2b1`) would suggest a per-class commit boundary.

**Pattern.** When one cycle ships both a substantive contract and a uniform mechanical mass-edit driven by that contract's strict gate, splitting the mechanical edit into its own commit makes the review-diff narrower and the contract change auditable on its own. α did not split. β's review absorbed both. No finding emerged from the conflation, but the diff-shape is logged here for γ's PRA judgment on whether to encode the split as a discipline.

### O10 — Frontmatter regression class survives downstream of the merge

**What happened.** During this α close-out session (post-merge, post-release), the harness surfaces system-reminder notices that 28 SKILL.md files under `cnos.core/skills/agent/`, `cnos.core/skills/ops/`, `cnos.core/skills/{audit,compose,design,naturalize,skill,write}/`, and `cnos.eng/skills/eng/{code,document,evolve,follow-up,go,ocaml,performance-reliability,process-economics,rca,ship,skill,test,tool,typescript,ux-cli,write-functional}/` are now back to a frontmatter shape that lacks the hard-gate fields α added in this cycle (e.g. `agent-ops/SKILL.md` is once again only `name: agent-ops` + `triggers: [...]`). If the I5 schema on `origin/main` is unchanged from `8adfd44`, the next CI run will fail on every reverted file.

α did not investigate further: post-merge edits to skills outside this cycle's contract are outside α's role contract for #301 (per `cdd/alpha/SKILL.md` §3.5 "α does not rewrite β's judgment frame or release process"). The observation is logged here as factual for γ's triage; the I5 gate either catches the regression on next CI (which would be evidence the cycle's leverage is working) or it does not (which would be evidence of a follow-up gap to investigate).

**Pattern.** A machine-checkable contract shipped in cycle N can be silently regressed by edits in cycle N+1 if those edits are made without running the validator locally. The I5 gate is the catch, not the prevention; the editor's own pre-commit discipline is the prevention. The pattern is generic to every CI-enforced contract, not specific to skill frontmatter.

## Friction log

### F-α-1 — Issue body internal contradiction (resolved at intake)

The non-goal vs hard-gate-AC tension named in O4 was the cycle's largest implementation-time decision. α spent meaningful intake time picking a consistent reading. The reading α resolved on (narrow non-goal, hard-gate fields *are* in scope) was correct in retrospect — β's review accepted the resulting backfills without finding — but α did not state the resolved reading back to γ before implementation. The right shape would have been an explicit clarification appended to the issue or surfaced in the dispatch acknowledgement.

### F-α-2 — CUE not pre-installed on the dev environment

CUE was not on the dev machine's `$PATH`. α fetched `v0.13.2` (the version intended to pin in CI) from upstream, staged it under `/usr/local/bin/cue`, and proceeded. The cycle's CI uses `cue-lang/setup-cue@v1.0.1`, so the runner installs it; only the dev loop needed the manual stage. Friction was small (one-time, ~30s) but logged because future α work in this surface area depends on the same setup.

### F-α-3 — YAML normalisation surfaced two latent ambiguities α had to fix in scope

The strict CUE parser rejected `cdd/design/SKILL.md`'s `description:` line and `cdd/gamma/SKILL.md`'s ninth `inputs:` entry (both unquoted-colon-space). α had to choose: fix in scope, or add the files to the exception list. Adding to the exception list was not viable because the reject was a YAML parse failure, not a missing-field — exception entries waive *missing* fields, not malformed YAML. α normalised both lines (one-quote each). The friction was conceptual: realising the validator's failure mode for malformed YAML was orthogonal to the exception path.

### F-α-4 — Polling outlived the cycle's coordination surface (O8 root)

α's polling loop was hand-coded against `.cdd/unreleased/301/beta-review.md` and did not transition when the cycle dir moved to `.cdd/releases/3.63.0/301/` at release. The operator had to surface that β had already approved, merged, released, and written the post-release close-out. α was not blocked on a missing artifact; α was polling the wrong path. The fix is in the protocol (poll union of `.cdd/unreleased/{N}/` and `.cdd/releases/*/{N}/`, or follow the cycle-dir-move event), not in α's loop logic.

### F-α-5 — Mass mechanical edit consumed disproportionate review-time budget

The 43-skill hard-gate backfill is a mechanical mass-edit (each addition is the same three-field-shape with file-specific values from the body). α authored each Edit individually because each file has a slightly different existing frontmatter shape (some have `triggers: [a,b]` flow style, some have block style; some have `triggers:` after `name:`, some after other fields). A purely scripted edit would have collided with the schema's hard-gate requirement that descriptions and governing_questions be authored from each skill's actual content (not generic templates). The friction was real: ~1/3 of α's implementation wall-clock went to the backfill. β O9 implicitly raises this as a process question (commit-shape discipline); O9 above logs it from α's side.

## Cycle-level engineering reading

**Diff level: L6.** The cycle introduces a machine-checkable contract for SKILL.md frontmatter, enforces it at CI time, and exposes a per-file shrinkable debt ledger. The contract is the single source of truth for shape; the script owns the orthogonal mechanism layer; the exception list is the visible degraded path; the README is the operator-facing summary; the fixtures are the regression suite. Each surface has one reason to change. The CI gate makes drift visible at push time, not at debugging time. Not L7: the gate is the catch, not the prevention — frontmatter drift can still happen pre-CI; the validator surfaces it on the next push, but does not eliminate the friction class entirely. β's β-side reading agrees: L6.

**Leverage achieved (positive):**

- **One source of truth for skill-module shape.** Before this cycle, the spec (`LANGUAGE-SPEC.md` §2 / §10 / §11) declared a contract that no tool checked; 30+ of 56 skills were missing fields the spec called required. After this cycle, `schemas/skill.cue` is the machine-checkable form of that contract; `cue vet` decides conformance.
- **Frontmatter drift becomes a CI-failure class.** New skills added in future cycles cannot ship without hard-gate fields without the I5 job failing. The class moves from "silent and accumulates" to "loud and CI-blocking on first push."
- **Debt is enumerable and shrinkable.** `schemas/skill-exceptions.json` is 43 entries; each entry maps to one specific skill and one or more specific missing-field fixes. The `spec_ref` field on each entry names the LANGUAGE-SPEC section to author from. Future α work to shrink the list is one-skill-one-fix.
- **Latent YAML ambiguity surfaces.** Two pre-existing unquoted-colon-space ambiguities (design + gamma) are normalised. Future ambiguity of the same class will fail the validator.

**Leverage not achieved (negative / scope):**

- **Body prose drift not caught.** β O4 names three live "writing" → "write" rename sites in body prose (`CDD.md` L611 + L785, `eng/skills/eng/README.md` L165). The validator's scope is frontmatter; body content is out of scope by design.
- **Cross-package `calls` still has no notation.** Per O6, the v0 schema's `calls`-resolution algorithm is intra-package only. The release skill's prose loads the writing skill across package boundaries; the static call edge cannot represent it.
- **CTB v0.2 fields not validated.** Per the issue's non-goals, this cycle is v0.1 only. v0.2 promotion (#289) ships its own validator surface (#303 `ctb-check` v0).

**Substantive change vs ritual.** The cycle's ritual cost is moderate: schemas/README.md (one new operator-facing doc), one new CI job, one exception list (43 entries to shrink), 43 skill-frontmatter backfills (mechanical). The substantive change is the contract itself: every SKILL.md going forward must carry hard-gate fields, and every exception must be field-specific with a `reason` and `spec_ref`. The substantive change earns the ritual.

## Cycle invariants observed

- **Dyad-plus-coordinator preserved.** α owned implementation through review; β owned review → merge → release → β close-out; γ owns the PRA. No role overran.
- **Cycle-branch-as-canonical-coordination-surface preserved through review.** All cycle artifacts (`self-coherence.md`, `beta-review.md`, both close-outs, β's verdict commits) lived on `claude/cnos-alpha-tier3-skills-MuE2P` until merge; β's β-side artifacts (`beta-review.md` round-1 verdict commit `40af6c0`, round-2 `acfa0cf`, Pass-2 `a477b1d`) committed to that branch.
- **α never wrote review artifacts; β never wrote implementation artifacts.** Round-2 fixes were α-authored (`171188e`, `55642db`); the merge / release / β-closeout commits were β-authored (per β O8 with the `beta-merge-test` identity drift on `b483f36` + `6300081`, but the role contract held — α did not produce those commits).
- **Self-coherence carried the CDD Trace through step 7 across both rounds.** Round-1 trace at `self-coherence.md` L1–186; round-2 fix-round appendix at L209–238; readiness signal at L188 (round 1) and L239 (round 2).
- **Pre-review gate row 10 (CI green) deferred to δ explicitly.** β's environment cannot reach `api.github.com`; β proxied via local merge-tree validator pass and risk-assessed the existing CI jobs (β O6). The deferral was named in β's R2 verdict and re-confirmed in Pass-2.
- **Branch never rebased after merge-base.** `merge-base` stayed at `a8e67b7` across rounds; α rebased onto β's R1 verdict commit `40af6c0` for round 2. β's R2 Pass-2 audit confirmed the merge into fresh `origin/main = 9d6a0fa` was conflict-free.
