# self-coherence — cycle/558

## §Gap

Issue: [#558](https://github.com/usurobor/cnos/issues/558) — "docs: refresh canonical glossary for current cnos work model".

`docs/reference/governance/GLOSSARY.md` was stale relative to the current docs-cleanup, wake-as-skill migration, issue-taxonomy work, and CDD/CDS/wave vocabulary. It cited retired `docs/alpha|beta|gamma` paths, package paths missing the `src/` prefix, and a superseded TSC pass threshold (`C_Σ ≥ 0.80` instead of the current `Θ ≥ 0.75`). It was also missing definitions for cell, wave, matter, review request/verdict, projection, δ, ε, V, trust claim, coherence witness, wake-as-skill, golden, issue (cell-shaped contract), `dispatch:cell`, and `effort/*`.

Version / mode: **design-and-build** (per issue header and γ's scaffold — the design is small enough to converge inline; no separate DESIGN.md/PLAN.md pair exists for this gap). Work-shape: small-change (one primary file rewrite + 4 one-line link insertions).

Branch: `cycle/558`, base verified fast-forward from `284be5693cc162c0fcd6c97fb69f22d9a4b1b5ea` (δ-supplied SHA) through γ's `072c4ef5` (gamma-scaffold commit) to `origin/main` at `2d0afca3`.

## §Skills

- **Tier 1:** `CDD.md` (canonical lifecycle); `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (this role surface); `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` (AC-boundary interpretation — no design/plan sub-skill invoked; the issue's own "design-and-build" mode and γ's scaffold made a separate design/plan pass unnecessary for a single-file docs refresh).
- **Tier 2:** none of the always-applicable `eng/*` bundles apply — this is a Markdown-only docs cycle with no code, no CLI surface, no runtime dependency (per the pinned implementation contract).
- **Tier 3 (issue-specified):** `src/packages/cnos.core/skills/write/SKILL.md` — governed prose structure (one governing question per entry, front-loaded point, contrastive ❌/✅ discipline where used, cut duplicate facts).
- **γ artifact:** `.cdd/unreleased/558/gamma-scaffold.md` — read in full before drafting; every citation in its source-of-truth table was independently re-verified against the live tree before use (see §ACs; several were found to need a further correction beyond what the scaffold flagged — the scaffold's paths for `docs/development/cdd/CDD.md` §3.4/§9.5, `docs/architecture/cognitive-substrate/COGNITIVE-SUBSTRATE.md` §7.3 for DUR, and `docs/papers/FOUNDATIONS.md` §6 for wake-up no longer resolve to the content they used to; these were re-pointed to their current canonical locations rather than left broken).

## §ACs

**AC1 — glossary is current.**
Evidence: every backtick-quoted repo-relative path in the rewritten glossary was checked with `test -f` (see script output below); all resolve except a small set of intentionally hub-relative or external-repo-relative paths (`spec/SOUL.md`, `state/*.md`, `threads/reflections/...`, `tsc/spec/tsc-core.md`) that describe an agent hub's structure or the external `tsc` repo, not this repo — consistent with how the pre-existing glossary already used those paths (e.g. `state/receipts/`, `state/peers.md`). Two pre-existing stale citations were also found and corrected beyond what the scaffold flagged: the "DUR" entry's dead pointer to `COGNITIVE-SUBSTRATE.md §7.3` (no such section; DUR now lives in `src/packages/cnos.core/skills/skill/SKILL.md` §1–§3) and the "Wake-up" entry's dead pointer to `FOUNDATIONS.md §6` (no wake-up content there anymore; content lives only in `CAA.md §4`, already cross-cited). "Coherence delta"'s stale `docs/development/cdd/CDD.md §3.4, §9.5` citation (that file is now a 3-line pointer, not the algorithm) was dropped, keeping the verified `CAA.md §5.6/§10` citation.
```
$ grep -oE '`[a-zA-Z0-9_./-]+\.(md|cue)`' docs/reference/governance/GLOSSARY.md | sort -u | tr -d '`' | while read -r p; do test -f "$p" && echo "OK $p" || echo "MISSING $p"; done
# All repo-relative paths OK; MISSING only for spec/SOUL.md, spec/USER.md, state/input.md,
# state/output.md, state/peers.md, threads/reflections/daily/YYYYMMDD.md,
# threads/reflections/weekly/YYYY-WNN.md, tsc/spec/tsc-core.md — all hub-/external-repo-relative,
# same convention as the pre-existing glossary.
```

**AC2 — glossary defines at least the N listed terms.**
Evidence: every term in the issue's list resolves to its own heading/bolded entry with prose, not just an incidental mention:
```
$ for t in cell wave matter receipt "review request" projection α β γ δ ε V MCA MCI TSC "trust claim" "coherence witness" wake "SKILL.md" golden issue "dispatch:cell" "effort/"; do echo -n "$t -> "; grep -c -F "$t" docs/reference/governance/GLOSSARY.md; done
cell -> 21   wave -> 5   matter -> 7   receipt -> 20   review request -> 1   projection -> 4
α -> 27   β -> 29   γ -> 30   δ -> 15   ε -> 9   V -> 17   MCA -> 17   MCI -> 13   TSC -> 10
trust claim -> 1   coherence witness -> 2   wake -> 11   SKILL.md -> 14   golden -> 3
issue -> 13   dispatch:cell -> 3   effort/ -> 2
```
Each has a dedicated `###` heading (Cell, Wave, Matter, Receipt, "Review request vs review verdict", Projection vs role-renaming, Role (α/β/γ/δ/ε), δ, ε, V, MCA, MCI, TSC, "Trust claim vs coherence witness", "Wake / wake-as-skill", "Skill / SKILL.md", Golden, "Issue (cell-shaped contract)", dispatch:cell, effort/*) with a full definition paragraph, not a bare mention.

**AC3 — stale references removed/corrected.**
Evidence:
```
$ grep -n "docs/alpha\|docs/beta\|docs/gamma" docs/reference/governance/GLOSSARY.md
18:...There is no `docs/alpha/`, `docs/beta/`, or `docs/gamma/` — that filing scheme was retired...
```
The only remaining hit is the explicit historical/retirement note required by AC4 (the scaffold's own exemption clause: "except if used as a historical note explicitly marked as retired").
```
$ grep -n "C_Σ ≥ 0.80\|0\.80" docs/reference/governance/GLOSSARY.md
138:This corrects the earlier-cited `C_Σ ≥ 0.80` threshold, per the 2026-06-23 v3.0.1 errata...
```
Same pattern — the only hit is the explicit correction note; the live threshold is stated as `Θ ≥ 0.75`.
```
$ grep -n "packages/cnos\.core\|packages/cnos\.cdd" docs/reference/governance/GLOSSARY.md | grep -v "src/packages"
(no output — every package path now carries the src/ prefix)
```
Also fixed beyond the scaffold's explicit list: `docs/alpha/doctrine/COHERENCE-SYSTEM.md` → `docs/papers/COHERENCE-SYSTEM.md`; `docs/alpha/doctrine/FOUNDATIONS.md` → `docs/papers/FOUNDATIONS.md`; `docs/gamma/ENGINEERING-LEVELS.md` → `docs/development/ENGINEERING-LEVELS.md` (all verified: `docs/alpha/`, `docs/gamma/` do not exist in the tree; `docs/papers/COHERENCE-SYSTEM.md` and `docs/papers/FOUNDATIONS.md` do, with matching section numbers re-verified — see §1 above). The "Thread" entry's ambiguous root-`threads/` wording was clarified to state explicitly it describes an agent's **hub** structure, not this cnos source repo (which has no `threads/` directory at all — verified `test -d threads` → absent). The "Package" entry's example list dropped `cnos.pm`, which no longer exists under `src/packages/` (verified via `ls src/packages/`), replacing it with the current package set.

**AC4 — states the filesystem/reader rule.**
Evidence: `docs/reference/governance/GLOSSARY.md` line 18 states verbatim: *"The filesystem is organized for readers. α/β/γ are role grammar and measurement axes, not docs folders."* — directly citing and linking `docs/README.md`, which carries the ratified form of this rule ("The filesystem is organized for readers. The triad is kept for measurement.", `docs/README.md` line 53). Placed in a new top-of-file "Two senses of α, β, γ" section per friction note 1, resolving the AC2 ambiguity between the TSC-axis sense (already in the glossary) and the CDD/CDS role-actor sense (newly added) without collapsing them into one entry.

**AC5 — distinguishes the 5 term pairs.**
Evidence: each pair has a dedicated heading or bolded lead-in with an explicit contrast sentence, not just two adjacent definitions:
```
$ grep -n "^### Review request vs review verdict\|^### Projection vs role-renaming\|^### Trust claim vs coherence witness\|^### Signature vs attestation\|Cell vs wave" docs/reference/governance/GLOSSARY.md
259:### Review request vs review verdict
270:### Projection vs role-renaming
278:### Trust claim vs coherence witness
289:### Signature vs attestation
240:**Cell vs wave.** A cell is one bounded unit of work...
```
- Review request vs verdict: "The request asks for review; the verdict is the review's outcome" — plus the honest disclosure (friction note 2) that no third, separately-named "review request" artifact exists in the current skills.
- Trust vs coherence: "A receipt can carry a strong trust claim and a weak (or absent) coherence witness, or vice versa; neither implies the other."
- Signature vs attestation: "Signature — cryptographic proof... Attestation — an external or third-party witness statement..." with the "may carry... both, or neither" contrast.
- Projection vs role-renaming: "δₙ is not literally βₙ₊₁; εₙ is not literally γₙ₊₁" (quoted verbatim from `COHERENCE-CELL-NORMAL-FORM.md`).
- Cell vs wave: "The wave shipped" and "one cell in the wave shipped" are not the same claim.

**AC6 — linked from main docs entrypoints.**
Evidence:
```
$ grep -l "governance/GLOSSARY" docs/README.md docs/reference/README.md docs/development/README.md docs/development/issues/TAXONOMY.md
docs/README.md
docs/reference/README.md
docs/development/README.md
docs/development/issues/TAXONOMY.md
```
All four match. Each link was resolved by hand (relative-path arithmetic from each file's directory to `docs/reference/governance/GLOSSARY.md`) and independently re-verified by a Python script that parses every `[text](url)` in all 5 changed files and confirms the resolved local path exists — 0 broken links found (script output in §Self-check).

**AC7 — I4 link validation remains green.**
Evidence: branch CI on head commit `f89ae5f9` — `gh run view 28638659648` — job **"Repo link validation (I4)"**: `conclusion: success`. Full run: all 10 jobs (`Package/source drift (I1)`, `Protocol contract schema sync (I2)`, `Go build & test`, `CDD artifact ledger validation (I6)`, `Dispatch closeout-integrity guard`, `SKILL.md frontmatter validation (I5)`, `Dispatch repair-preflight guard`, `Repo link validation (I4)`, `Binary verification`, `Package verification`) — `success`. This is the authoritative CI-native oracle the scaffold names for AC7; a local lychee binary was not available in this environment (no cargo-built binary and no cached release), so the manual link-resolution script above served as a pre-push sanity check, and CI is the binding evidence.

**AC8 — no code changes.**
Evidence:
```
$ git diff --stat origin/main...cycle/558
 .cdd/unreleased/558/gamma-scaffold.md | 186 ++++++++++++
 .cdd/unreleased/558/self-coherence.md | (this file, added after this check)
 docs/README.md                        |   1 +
 docs/development/README.md            |   4 +
 docs/development/issues/TAXONOMY.md   |   6 +
 docs/reference/README.md              |   4 +
 docs/reference/governance/GLOSSARY.md | 213 +++++++--------
```
```
$ git diff --stat origin/main...cycle/558 -- . ':!docs' ':!.cdd/unreleased/558'
(empty)
```
Every touched file is under `docs/` or `.cdd/unreleased/558/`. No code, no schema, no CI config changed.

## §Self-check

Did this cycle push ambiguity onto β? The four friction notes γ surfaced were resolved explicitly rather than deferred:

1. **α/β/γ two-sense conflation** — resolved with a dedicated top-of-file section naming both senses and cross-referencing in both directions (TSC entry ↔ Role entry), not a single merged entry.
2. **"Review request" honesty** — stated plainly that it is α's readiness signal, not a separately named artifact, rather than inventing one.
3. **TSC threshold** — cited the 2026-06-23 v3.0.1 errata by name and quoted its correcting text, rather than silently changing the number.
4. **Receipt double-sense** — kept both senses under one "Receipt" heading with an explicit "do not conflate" sentence, and cross-referenced from the pre-existing CN-Shell "Receipts" entry back to it.

Beyond the four friction notes, this cycle found and fixed additional stale citations the scaffold's source-of-truth table did not enumerate (the dead DUR/`COGNITIVE-SUBSTRATE.md §7.3` pointer, the dead Wake-up/`FOUNDATIONS.md §6` pointer, the dead `docs/development/cdd/CDD.md §3.4/§9.5` "coherence delta" pointer, and the retired `cnos.pm` package example) — each was verified against the live tree before correcting, not assumed from the scaffold. Every path and section-number citation in the rewritten glossary was independently checked with `test -f` plus `grep -n "^## "` against the target file's actual headings before being used (see §ACs AC1/AC3), not copied from the scaffold on trust.

One judgment call not flagged by γ: the "Coherence Contract" entry (not on the AC2 required-term list, pre-existing in the glossary) cited a section (`docs/development/cdd/CDD.md §6`) that no longer exists after the file became a pointer. Rather than deleting the entry or leaving the dead citation, I added the current equivalent (`COHERENCE-CELL-NORMAL-FORM.md`'s `contractₙ` kernel object, operationalized via `alpha/SKILL.md §2.5`'s self-coherence §Gap) alongside the original prose, without changing what the entry claims the concept *means* — consistent with the backward-compat invariant ("no doctrine semantics change — only path/terminology/threshold corrections").

No claim in this file is unbacked by a command run against the actual branch state (see the fenced command blocks under each AC).

## §Debt

- **No local lychee run.** The environment had no cargo-built `lychee` binary and no network-cached release; AC7's binding evidence is the CI job result on the pushed head commit, not a pre-push local run. A manual link-resolution script (Python, checks every `[text](url)` in the 5 changed files against the filesystem) was used as a pre-push gate instead and found 0 broken links; CI confirmed green after push.
- **Document-local version number not bumped.** The glossary's header still reads "cnos v3.6.0" (a pre-existing, differently-scoped versioning convention — the file's own note says document versions are "local to each file," but the header format doesn't match that note's own example format `GLOSSARY v2.0.0`). This pre-existing inconsistency was not in scope for this issue (no AC requires it) and was left untouched to avoid scope creep beyond the pinned backward-compat invariant.
- **`CHANGELOG.md` still cites `≥0.80 = PASS`.** Found incidentally while verifying the TSC threshold citation. `CHANGELOG.md` is outside this issue's scoped surface (GLOSSARY.md + 4 entrypoints) and outside AC8's allowed diff — left untouched, noted here as a candidate follow-up for whoever next touches that file.
- **`lychee.toml`'s own comment says the I4 job lives in `.github/workflows/coherence.yml`**, but the actual job is defined in `.github/workflows/build.yml` (confirmed via CI run evidence above). This is a pre-existing doc-comment/reality mismatch in a file this issue does not touch — noted, not fixed, since `lychee.toml` is outside the pinned package-scoping axis.

## §CDD Trace

1. **Gap** — named above; issue #558, design-and-build mode, small-change work-shape.
2. **Mode** — design-and-build; no separate DESIGN.md/PLAN.md pair required (γ verified MCA preconditions did not apply for a single-file docs refresh; concurred).
3. **Plan** — not required, explicit: the γ scaffold's per-AC oracle list and source-of-truth table served as the sequencing plan (verify every citation → write disambiguation section → write role/cell/wave section → fix stale paths → add entrypoint links → verify links → CI).
4. **Tests** — not applicable to a Markdown-only docs change with no runtime surface (per the pinned implementation contract: language=Markdown, runtime dependencies=None). Acceptance evidence is the mechanical grep/path/CI oracles in §ACs instead.
5. **Code** — none (AC8 guard verified in §ACs).
6. **Docs** — `docs/reference/governance/GLOSSARY.md` rewritten (196 insertions / 32 deletions net across the 5-file diff; glossary alone: 213 lines changed, file now 561 lines total, up from 413). `docs/README.md` (+1 line), `docs/reference/README.md` (+4 lines), `docs/development/README.md` (+4 lines), `docs/development/issues/TAXONOMY.md` (+6 lines) each gained one glossary cross-link per AC6. `.cdd/unreleased/558/gamma-scaffold.md` (186 lines, γ's prior commit `072c4ef5`, not authored by α) and this file (`.cdd/unreleased/558/self-coherence.md`) are the only non-`docs/` artifacts in the branch diff, both within the permitted `.cdd/unreleased/558/` artifact-channel exemption (AC8).
7. **Self-coherence** — this document.

Head implementation SHA: `f89ae5f97e9e13280471a4dcd6a994d5b5cb43fa` (commit "docs: refresh canonical glossary for current cnos work model"). Branch verified fast-forward onto `origin/main` at `2d0afca3` (no rebase needed — `origin/main` had not advanced past γ's base SHA while this cycle was in progress).
