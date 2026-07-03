# γ scaffold — cycle/558

Issue: [#558](https://github.com/usurobor/cnos/issues/558) — "docs: refresh canonical glossary for current cnos work model"

Mode: **design-and-build** (per issue header — no separate DESIGN.md/PLAN.md exist for this cell; the "design" is small enough to converge inline: which terms to add, which stale refs to strip). γ verified no stable DESIGN.md/PLAN.md pair exists for this gap, so MCA mode is not applicable (per `cnos.cdd/skills/cdd/issue/SKILL.md` MCA preconditions) — `design-and-build` is correct, matching the issue's own `Mode:` header.

Work-shape: **small-change** (single file is the primary edit surface, plus 4 one-line link insertions; 8 ACs, all docs-only, AC8 is a negative/guard AC).

Branch: `cycle/558`

Base SHA: `2d0afca31d0917ead4f3c8b555a780da0c337280` (verified fast-forward descendant of the δ-supplied SHA `284be5693cc162c0fcd6c97fb69f22d9a4b1b5ea`; the two commits in between — `16b6ddfc` agent-admin heartbeat, `2d0afca3` board-map regeneration — are unrelated to this cell's surface).

---

## Surfaces γ expects α to touch

- `docs/reference/governance/GLOSSARY.md` (primary rewrite/refresh target)
- `docs/README.md` (add a glossary link)
- `docs/reference/README.md` (add a glossary link)
- `docs/development/README.md` (add a glossary link)
- `docs/development/issues/TAXONOMY.md` (add a glossary link)

No other files. AC8 ("No code changes") is a hard guard — the diff must not touch anything outside `docs/` (or `.cdd/unreleased/558/` for α/β artifacts).

---

## Peer enumeration (γ/SKILL.md §2.2a)

`ls docs/reference/governance/` shows `GLOSSARY.md` already exists (this is a refresh, not a net-new glossary — consistent with issue non-goal "do not create a second glossary"). `rg -l "glossary" docs/` returns only the four entrypoint files named above plus the glossary itself — no second glossary-shaped file was found; the issue's non-goal is uncontested.

---

## Source-of-truth table

| Claim / surface | Canonical source | Status |
|---|---|---|
| Current glossary content (stale baseline) | `docs/reference/governance/GLOSSARY.md` | Exists, stale — cites `docs/alpha/doctrine/...`, `docs/gamma/ENGINEERING-LEVELS.md`, `packages/cnos.core/...` (missing `src/` prefix), root `threads/` structure, `C_Σ ≥ 0.80` threshold |
| Docs entrypoints to link glossary from | `docs/README.md`, `docs/reference/README.md`, `docs/development/README.md`, `docs/development/issues/TAXONOMY.md` | All exist, read in full by γ; none currently link the glossary |
| "Filesystem organized for readers, not role folders" rule | `docs/README.md` lines 5–12, 45–53 (already states: "The filesystem is organized for readers... α/β/γ triad is no longer a filing taxonomy; it is kept only as a coherence measurement... never as folders") | Shipped — glossary AC4 should restate this rule and cite `docs/README.md`, not invent new wording |
| TSC current pass threshold (glossary says stale `C_Σ ≥ 0.80`) | `docs/reference/protocol/cn/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md` (2026-06-23 errata correcting the normative threshold) | Current threshold is **PASS ≥ 0.75** (Θ default; teams may set stricter). `docs/reference/ctb/SEMANTICS-NOTES.md` uses `Θ` generically. α must cite the errata, not invent a number. |
| Role letters α/β/γ (TSC axes today; also CDD/CDS role actors) | `docs/reference/governance/GLOSSARY.md` (TSC axis def, keep); `ROLES.md` (root) §1 "Generic Role-Scope Ladder", §9 Glossary | α = produces; β = reviews/discriminates; γ = coordinates the α↔β loop |
| δ (role) | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`; `ROLES.md` §1 | Boundary/operator actor — reads V's verdict, records `BoundaryDecision`, operates wave cadence, dispatches every role |
| ε (role) | `src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md`; `ROLES.md` §4b, §9 | Reads cross-cycle receipt streams, surfaces `*-gap` findings, writes protocol-patch proposals (`cdd-iteration.md`) |
| V (validator) | `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` (`V : Contract × Receipt × EvidenceGraph → ValidationVerdict`) | Pure validator predicate emitting PASS/FAIL/WARN; distinct from δ's `BoundaryDecision` |
| cell | `docs/papers/CELL-OF-CELLS.md` §3, §19; `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` (draft, defers to CDD.md for current behavior); `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` (the cell-shaped issue contract) | Bounded work unit: receives a message, closes internal work via α/β/γ triad, validates boundary transmission via V+δ, emits a receipt upward |
| wave | `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` §10 "Wave Coordination"; `docs/development/issues/TAXONOMY.md` line 19 ("wave δ may relabel") | A δ-run sequence of related cycles under one durable, git-committed manifest (`.cdd/waves/{wave-id}/manifest.md`) |
| cell vs wave distinction | `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` §10 vs `issue/SKILL.md` | Cell = one bounded unit of work (one issue, one α/β/γ loop); wave = an ordered sequence of cells under one δ-owned manifest |
| matter | `ROLES.md` §9 Glossary; `docs/papers/CELL-OF-CELLS.md` (`matterₙ := αₙ.produce(...)`) | The deliverable α produces (protocol-specific: code+tests in CDD, prose in CDW, etc.) — not CN-Shell/runtime substrate |
| receipt | `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md`; `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md`; existing GLOSSARY.md "Receipts" entry (CN Shell sense) | Two valid current senses that must both be named: (1) CN Shell runtime receipt for every mutation (`state/receipts/`, already in glossary); (2) the CDD cell receipt γ emits at close-out carrying contract/evidence/validation/boundary blocks, consumed by V then δ. Glossary should keep sense (1) and add sense (2), disambiguated. |
| review request vs review verdict | `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (review-readiness signal, line ~96 "review request"); `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` (RC/APPROVE verdict written to `beta-review.md`) | Review request = α's outbound signal that the cell is ready for β (via `self-coherence.md` + branch state); review verdict = β's RC(request-changes)/APPROVE decision. These are not formally split into two named artifacts today — α should state this plainly rather than inventing false precision. |
| projection vs role-renaming | `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` line 297 | "This is a projection under scope-lift. It is not a flat role-renaming inside a single cell. δₙ is not literally βₙ₊₁; εₙ is not literally γₙ₊₁." — projection = a typed object (decision + receipt-stream observation) crossing a scope boundary; role-renaming = the (rejected) idea that a role at scope n is just relabeled as a different role at scope n+1 |
| trust claim vs coherence witness | `docs/papers/DUMB-MODELS-SMART-CELLS.md` (lines ~17, 39, 120, 348, 552) | Trust claim = the strength of attribution/signature on a receipt; coherence witness = whether the artifacts still describe one coherent unit of work (what TSC measures) — independent, composable properties of a receipt |
| signature vs attestation | `docs/papers/DUMB-MODELS-SMART-CELLS.md` (lines ~257–258, 779–780) | Signature = cryptographic proof the artifact is unaltered/tied to an identity; attestation = an external/third-party witness statement about the event |
| wake / wake-as-skill | `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md`; `schemas/skill.cue` `#Wake` definition | A wake is a typed `SKILL.md`-shaped module (`artifact_class: wake`) rendered into `.golden.yml` GitHub Actions workflows by `cn install-wake` |
| golden | `.golden.yml` files, e.g. `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` (header: "DO NOT EDIT. Rendered by `cn install-wake`") | A rendered, checked-in fixture treated as source-of-truth to diff against — golden-file convention |
| issue (cell-shaped contract) | `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` | Executable work contract: gap, impact, status truth, source of truth, scope, ACs, proof plan, non-goals |
| dispatch:cell / protocol:cds / effort/* | `docs/development/issues/TAXONOMY.md` §"Dispatch and protocol", §"Effort labels" | `dispatch:cell` = issue is executable by CDS/CDD; `protocol:{cds,cdd}` = which concrete protocol; `effort/{S,M,L,XL}` = ordinal size estimate, not priority |
| I4 link-validation oracle | `.github/workflows/build.yml` job `link-check` ("Repo link validation (I4)"), `lychee.toml` | `lychee --offline --no-progress --hidden --config lychee.toml '**/*.md'`, `exclude_path = [".cdd"]` (frozen evidence out of scope); fails CI on any broken link in active docs |

---

## Per-AC oracle list

**AC1 — `docs/reference/governance/GLOSSARY.md` is current.**
Oracle: manual/inspectable — every path cited inside the glossary resolves (`git ls-files` or `test -f` on every backtick-quoted path in the file); no reference to a path that does not exist in the current tree at merge time. Mechanical proxy: the I4 lychee run (AC7) catches dead *links*, but plain-text backtick path citations (not `[]()` links) need a manual `rg` sweep by β — grep every `` `docs/...` ``, `` `src/packages/...` ``, `` `packages/...` `` token in the file and `test -f` each one.

**AC2 — glossary defines at least the N listed terms.**
Oracle: mechanical. `for t in cell wave matter receipt "review request" projection α β γ δ ε V MCA MCI TSC "trust claim" "coherence witness" wake SKILL.md golden issue "dispatch:cell" "effort/*"; do grep -c -F "$t" docs/reference/governance/GLOSSARY.md; done` — every term must resolve to a **defined entry** (a heading or bolded term with a definition paragraph), not just an incidental mention. β checks each term has its own `###`/`**term**` entry with prose, not merely appearing inside another entry's body.

**AC3 — stale references removed/corrected.**
Oracle: mechanical negative grep. `rg -n "docs/alpha|docs/beta|docs/gamma" docs/reference/governance/GLOSSARY.md` → zero hits (except if used as a historical note explicitly marked as retired, which the issue does not ask for). `rg -n "^threads/|root \`threads/\`" GLOSSARY.md` for root-`threads/` claims not otherwise corrected. `rg -n "packages/cnos\.core|packages/cnos\.cdd" docs/reference/governance/GLOSSARY.md` → any hit must have the `src/` prefix (`src/packages/cnos.core/...`) or be flagged. `rg -n "C_Σ ≥ 0.80|0\.80" docs/reference/governance/GLOSSARY.md` → zero hits; replaced with the current `PASS ≥ 0.75` figure, cited to the errata source in the table above.

**AC4 — states the filesystem/reader rule.**
Oracle: inspectable. Glossary contains a statement equivalent to `docs/README.md`'s "The filesystem is organized for readers. The triad is kept for measurement" (or an explicit cite to that doc) plus the role-grammar/measurement-axis framing ("α/β/γ are role grammar and measurement axes, not docs folders"). Positive case: the exact clause or a faithful paraphrase with citation is present. Negative case: the old "Defined in: docs/alpha/doctrine/..." style path attribution without the reader-organization caveat.

**AC5 — distinguishes the 5 term pairs.**
Oracle: inspectable, per-pair. For each pair (review request vs review verdict; trust vs coherence; signature vs attestation; projection vs role-renaming; cell vs wave), the glossary must contain prose that names both terms **and** states how they differ (not just two separate definitions with no explicit contrast). β checks each pair for an explicit "X is not Y" / "X vs Y" sentence, matching the distinctions surfaced in the source-of-truth table above (γ does not require α to invent new distinctions — the citations above already state them).

**AC6 — linked from main docs entrypoints.**
Oracle: mechanical. `grep -l "governance/GLOSSARY" docs/README.md docs/reference/README.md docs/development/README.md docs/development/issues/TAXONOMY.md` → all four files match (accounting for each file's relative path depth to the glossary).

**AC7 — I4 link validation remains green.**
Oracle: mechanical, CI-native. The `link-check` job ("Repo link validation (I4)") in `.github/workflows/build.yml` must report success on the cycle branch's CI run — `gh run list --branch cycle/558 --json name,conclusion | jq` filtered to that job name, or equivalently `lychee --offline --no-progress --hidden --config lychee.toml '**/*.md'` run locally must exit 0. Any new link α adds to the glossary or the four entrypoints (including the new cross-links from AC6) must resolve.

**AC8 — no code changes.**
Oracle: mechanical, negative. `git diff --stat origin/main...cycle/558 -- . ':!docs' ':!.cdd/unreleased/558'` must be empty. Any file outside `docs/` (excluding this cycle's own `.cdd/unreleased/558/` artifact channel) appearing in the diff is a violation.

---

## Scope guardrails (restated from the issue, binding on α)

**In scope:**
- Update `docs/reference/governance/GLOSSARY.md`.
- Remove stale paths and outdated claims.
- Add current terms for CDD, waves, issue taxonomy, wakes, skills, trust, and coherence.
- Link the glossary from `docs/README.md`, `docs/reference/README.md`, `docs/development/README.md`, `docs/development/issues/TAXONOMY.md`.
- Keep definitions short and precise.

**Out of scope (do not touch, do not drift into):**
- Do not create a second glossary.
- Do not change doctrine semantics (i.e. do not redefine what MCA/CAP/TSC/etc. *mean* — only correct stale path citations and thresholds against already-shipped canonical sources; α is transcribing current truth, not inventing new doctrine).
- Do not change issue labels (i.e. do not edit `TAXONOMY.md`'s label definitions themselves — only add the glossary link).
- Do not start Demo 0.
- Do not rewrite all docs repo-wide to use glossary terms in this cycle — only the four named entrypoints get a link.
- No code changes anywhere (AC8) — this is the hard guard; the cycle branch diff must be confined to `docs/` plus this cycle's `.cdd/unreleased/558/` artifacts.

---

## Empirical anchor / friction notes

- **Friction 1 (AC2 term ambiguity):** the issue lists `α, β, γ, δ, ε, V` as terms to define, but the current glossary defines α/β/γ *only* as TSC measurement axes (Pattern/Relation/Exit), not as CDD/CDS role actors. δ, ε, and V are role actors with **no TSC-axis meaning**. α must add role-actor definitions for δ/ε/V (and probably a short role-actor gloss for α/β/γ too, cross-referencing but not replacing the existing TSC-axis entries) without collapsing the two senses into one confusing entry. This is a real design decision left to α's judgment — γ surfaces it here rather than smuggling a prescriptive answer, per `cnos.cdd/skills/cdd/gamma/SKILL.md` §Rule 3.5 (γ transfers artifact facts, not hidden reasoning) — but the source-of-truth table above already resolves the citations α needs.
- **Friction 2 (review request):** "review request" is not a separately named, canonical artifact in the current skills (only "review-readiness signal" and β's "verdict" are named). α should state the distinction honestly (signal vs verdict) rather than inventing a formal "review request" object that doesn't exist in the skills. This keeps AC5 satisfiable without overclaiming (issue/SKILL.md's "overclaiming" failure mode).
- **Friction 3 (TSC threshold correction):** the existing glossary's `C_Σ ≥ 0.80` figure is stale relative to the 2026-06-23 errata in `docs/reference/protocol/cn/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md`, which sets the default `Θ` at **0.75**. This is exactly the "stale TSC threshold" the issue's Problem section flags — α must cite the errata source, not just change the number silently.
- **Friction 4 (receipt double-sense):** "receipt" already has a shipped CN-Shell-runtime meaning in the glossary (`state/receipts/`, one per mutation). The issue's AC2 also wants "receipt" in the CDD-cell sense (γ's close-out receipt, validated by V then δ). Both are current and both should appear, clearly disambiguated (two senses under one heading, or two headings) — not a replacement of one by the other.
- No γ-level ambiguity was found that blocks dispatch; the above are α-facing judgment calls with citations already supplied, not open design questions requiring a δ escalation.

---

## α prompt

```text
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 558 --repo usurobor/cnos --json title,body,state,comments
Refs #558
Branch: cycle/558
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Markdown |
| CLI integration target | N/A |
| Package scoping | `docs/reference/governance/GLOSSARY.md` (primary); `docs/README.md`, `docs/reference/README.md`, `docs/development/README.md`, `docs/development/issues/TAXONOMY.md` (one-line glossary link each) |
| Existing-binary disposition | N/A |
| Runtime dependencies | None |
| JSON/wire contract preservation | N/A |
| Backward-compat invariant | Existing glossary entries that are still current are preserved as-is; stale entries are corrected in place (paths, thresholds); new entries are additive. No doctrine semantics change — only path/terminology/threshold corrections against already-shipped canonical sources and new term additions. |

## Scaffold and AC oracles

γ has filed `.cdd/unreleased/558/gamma-scaffold.md` on this branch. It contains:
- a source-of-truth table (canonical paths for every AC2-listed term, the "filesystem organized for readers" rule, the current TSC threshold correction, and the 5 AC5 term-pair distinctions);
- a per-AC oracle list (AC1–AC8) naming exactly how each AC will be checked;
- scope guardrails restating the issue's in/out-of-scope lists;
- 4 friction notes on judgment calls left open for you (role-actor sense of α/β/γ/δ/ε/V vs TSC-axis sense; "review request" honesty; TSC threshold citation; receipt double-sense).

Read the scaffold before starting. Do not re-derive the source-of-truth citations from scratch — they are already resolved there; cite them directly in the glossary entries you write.

## Acceptance criteria (from the issue body — read the live issue for full text)

- AC1: `docs/reference/governance/GLOSSARY.md` is current.
- AC2: glossary defines at least: cell, wave, matter, receipt, review request, projection, α, β, γ, δ, ε, V, MCA, MCI, TSC, trust claim, coherence witness, wake, SKILL.md, golden, issue, dispatch:cell, effort/*.
- AC3: stale references to old docs/alpha, docs/beta, docs/gamma, root threads/, old package paths (missing `src/` prefix), and the stale TSC threshold are removed or corrected.
- AC4: glossary states the current rule — the filesystem is organized for readers; α/β/γ are role grammar and measurement axes, not docs folders.
- AC5: glossary distinguishes: review request vs review verdict; trust vs coherence; signature vs attestation; projection vs role-renaming; cell vs wave.
- AC6: glossary is linked from docs/README.md, docs/reference/README.md, docs/development/README.md, docs/development/issues/TAXONOMY.md.
- AC7: I4 link validation (`link-check` job, `.github/workflows/build.yml`, `lychee.toml`) remains green.
- AC8: no code changes — diff confined to `docs/` plus this cycle's own `.cdd/unreleased/558/` artifacts.

Follow your role-skill's load order and write your usual self-coherence artifact at `.cdd/unreleased/558/self-coherence.md` before signaling review-readiness.
```

---

## β prompt

```text
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 558 --repo usurobor/cnos --json title,body,state,comments
Refs #558
Branch: cycle/558

γ has filed `.cdd/unreleased/558/gamma-scaffold.md` on this branch — read it before reviewing.
It contains the per-AC oracle list (AC1–AC8), the source-of-truth table (canonical paths for every glossary term and the 5 AC5 term-pair distinctions), and the scope guardrails.

Walk each AC (AC1–AC8) independently against the diff using the oracle named in the scaffold for that AC — do not substitute your own oracle where the scaffold already names a mechanical check (e.g. AC2/AC3/AC6/AC8 are grep-checkable; AC7 is CI-checkable via the `link-check` job). Confirm the implementation contract's 7 axes are honored (Markdown-only, no code, no new package/binary surface) per your Rule 7 (implementation-contract coherence).

Render a verdict per AC: converge or iterate, with findings. Overall verdict: APPROVE only if every AC converges and AC8's negative-diff guard holds.
```

---

**End of scaffold.**
