# cnos.handoff Extraction Map — v0.1

**Version:** 0.1 (Sub 1 of [cnos#404](https://github.com/usurobor/cnos/issues/404); shipped under [cnos#415](https://github.com/usurobor/cnos/issues/415))
**Date:** 2026-05-22
**Placement:** `src/packages/cnos.handoff/docs/extraction-map.md`
**Audience:** δ operators dispatching Subs 2–6 of [cnos#404](https://github.com/usurobor/cnos/issues/404); reviewers verifying that every handoff-class surface in `cnos.cdd` (and the handoff-class surfaces that migrated into `cnos.cds` during the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave) has a named `cnos.handoff` destination and a named migration sub.
**Scope:** A surface-by-surface migration plan from handoff-class doctrine currently resident in `cnos.cdd` (and in `cnos.cds` per [cnos#411](https://github.com/usurobor/cnos/issues/411)) onto its planned `cnos.handoff` canonical home, with the migration sub-issue named for each surface.

> This document is a **mapping table + commentary**, not a re-statement of CDD / CDS / CDR doctrine. It does not embed source prose; citations carry the source path and the section marker at which the content currently resides. Subs 2–5 of [cnos#404](https://github.com/usurobor/cnos/issues/404) execute the migration against this map; Sub 6 re-points cross-references in cdd / cds / cdr to the new cnos.handoff canonical homes and closes the tracker.

---

## 0. Method

### 0.1 Source pin

Source evidence cited herein refers to the cnos repository's `main` branch at the head of `cycle/415`'s parent commit (`bc29c009` — Merge cycle/414: Add DECREASING-INCOHERENCE.md essay). All line-numbered citations resolve at that pin. The handoff-class surfaces have shifted homes across the [cnos#403](https://github.com/usurobor/cnos/issues/403) wave (some surfaces named in [cnos#404](https://github.com/usurobor/cnos/issues/404)'s opening table moved from `cnos.cdd` into `cnos.cds` during Subs 2–5 of that wave); the rows below record each surface at its **current** canonical home at the source pin, not at the home `cnos#404`'s opening table named at filing time.

### 0.2 cnos.handoff v0.1 target surfaces

The migration lands into three surface classes within `cnos.handoff`:

1. **`skills/handoff/HANDOFF.md`** — the canonical wire-format contract (forthcoming; ships with Sub 2 or Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404), whichever lands the first substantive doctrine). Sections under `HANDOFF.md` will be named `§"Directional cases"`, `§"Wire-format artifacts"`, `§"Consumer protocol"`, etc.
2. **`skills/handoff/<area>/SKILL.md`** — sub-area skill files for operationally-heavy surfaces (cross-repo, dispatch, mid-flight, artifact-channel, receipt-stream). Each sub-area skill is the operational expansion of the `HANDOFF.md` section that names it. The cross-repo case is special: Sub 2 moves the entire 644-line `cnos.cdd/skills/cdd/cross-repo/SKILL.md` wholesale rather than expanding from a `HANDOFF.md` section, because the source is already a complete canonical surface.
3. **`schemas/handoff/`** (optional; deferred decision) — if Sub 3 lifts the implementation-contract 7-axes table from Markdown into a typed JSON schema, this is where it lives. Sub 3 decides; not pre-committed by this map.

### 0.3 Sub naming

The migration sub-issues follow [cnos#404](https://github.com/usurobor/cnos/issues/404)'s wave shape:

- **Sub 2** — "Move cross-repo state machine + bundles into cnos.handoff" (P2; depends on this Sub 1). Owns the entire `cross-repo/SKILL.md` wholesale move (644 lines).
- **Sub 3** — "Move dispatch-prompt template + implementation-contract schema into cnos.handoff" (P2; depends on Sub 1). Owns `cdd/gamma/SKILL.md §2.5` (dispatch prompts + 7-axes table + four-surface mesh) + `cdd/delta/SKILL.md §2` (inward-membrane enrichment doctrine).
- **Sub 4** — "Move intra-cycle handoff mechanisms into cnos.handoff" (P2; depends on Sub 1). Owns the mid-flight `gamma-clarification.md` mechanism (cnos#391 empirical anchor) + the `.cdd/unreleased/{N}/` artifact channel rules (post-#411 canonical in `cnos.cds`). Sub 4 decides migration semantics for the cds-resident piece: move / mirror / cite-only.
- **Sub 5** — "Move cross-cycle receipt-stream into cnos.handoff" (P2; depends on Sub 1). Owns `cdd/post-release/SKILL.md §5.6b` (per-finding shape + `cdd-iteration.md` cadence rule + aggregator update procedure + cross-repo trace bundle). May also absorb `cdd/harness/SKILL.md §5.4` polling primitives if Sub 5 decides they are handoff-class rather than runtime-substrate.
- **Sub 6** — "Cross-reference cleanup + close tracker" (P2; depends on Subs 2–5). Sweeps citations in `cnos.cdd` / `cnos.cds` / `cnos.cdr` (and the docs essays under `docs/gamma/`) to re-point at the new `cnos.handoff` canonical homes. Closes [cnos#404](https://github.com/usurobor/cnos/issues/404).

The Sub-2/Sub-3/Sub-4/Sub-5/Sub-6 split follows the natural cohesion of the surfaces; cross-cutting bullets (e.g., the four-surface mesh which γ template / δ enrichment / α constraint / β verification span Sub 2's cross-repo, Sub 3's dispatch, and Sub 4's mid-flight) are noted in the row's commentary.

### 0.4 Coverage commitment

Per [cnos#415](https://github.com/usurobor/cnos/issues/415) AC5 ("`docs/extraction-map.md` contains tables for each of the 6 surfaces listed in D2"; mechanical: `grep -c "^## " docs/extraction-map.md` ≥ 6), the tables below carry one `##`-headed section per surface family. The six required families are:

1. Cross-repo state machine + bundles (Sub 2)
2. Dispatch-prompt template (Sub 3)
3. Implementation-contract schema (Sub 3)
4. Mid-flight `gamma-clarification.md` mechanism (Sub 4)
5. `.cdd/unreleased/{N}/` artifact channel rules (Sub 4)
6. `cdd-iteration.md` receipt-stream + INDEX.md aggregator (Sub 5)

Two additional handoff surfaces were discovered during authoring (recorded as additional `##`-headed sections; counted toward the ≥ 6 floor):

7. Polling primitives (Sub 5 or deferred; the boundary is contested)
8. Spec-staleness propagation (Sub 4; folded into mid-flight or artifact-channel)

Plus the close-out row:

9. Cross-reference cleanup + close tracker (Sub 6).

### 0.5 Note-column conventions

Each row's note column uses one of:

- **verbatim move** — the source content moves as-is; rename only the file path and update cross-references.
- **rewrite** — the source content needs minor edits (e.g., software-class vocabulary stripped; consumer-protocol-neutral phrasing) before landing.
- **pointer-only** — the destination is a one-line pointer to a canonical home elsewhere (used for surfaces where the canonical home stays put and `cnos.handoff` only provides a loader entrypoint).
- **migration-semantics-undecided** — the destination is committed but the migration mechanism is one of (move / mirror / cite-only) and the dispatching δ for the sub decides. Used for rows where the source is already in `cnos.cds` post-#411.

---

## 1. Cross-repo state machine + bundles (§Cross-repo) — Sub 2

**Source family:** [`cnos.cdd/skills/cdd/cross-repo/SKILL.md`](../../cnos.cdd/skills/cdd/cross-repo/SKILL.md) — entire skill, 644 lines (the largest extraction target in this wave).
**Source content:** Four directional cases (a: inbound proposal; b: outbound iteration trace; c: bilateral iteration; d: operator-pending bundle); canonical path discipline; 8-event STATUS state machine (drafted / submitted / accepted / modified / rejected / landed / withdrawn / revised); transition graph; emitters per event; master/sub `landed` rule; bundle-state phases; bundle file set per case; LINEAGE.md schema per case; feedback-patch format; bundle archival rule (asymmetric: source-side may delete on terminal STATUS; target-side preserved indefinitely); hat-collapse attribution; known protocol edge cases; embedded kata.
**Migration sub:** Sub 2 of [cnos#404](https://github.com/usurobor/cnos/issues/404).
**Destination commitment:** Move wholesale to `cnos.handoff/skills/handoff/cross-repo/SKILL.md` (same filename, same frontmatter shape; `parent: handoff` instead of `parent: cdd`). The source file is deleted from `cnos.cdd` and replaced with a single-line pointer (or removed entirely if no `cnos.cdd` skill needs a cross-repo entrypoint).

| Source (path + §-anchor) | Destination (path under cnos.handoff/) | Sub | Note |
|---|---|---|---|
| `cnos.cdd/skills/cdd/cross-repo/SKILL.md` (entire file; 644 lines) | `cnos.handoff/skills/handoff/cross-repo/SKILL.md` | Sub 2 | Verbatim move; rename `parent: cdd` → `parent: handoff` in frontmatter; update internal cross-references to other cnos.cdd skills (`gamma/SKILL.md`, `post-release/SKILL.md`, `issue/SKILL.md`) to use absolute cross-package paths (`../../../cnos.cdd/skills/cdd/<role>/SKILL.md`). |
| `cnos.cdd/skills/cdd/cross-repo/SKILL.md §2.3` (STATUS state machine; 8-event vocabulary; transition graph; emitters; master/sub rule) | `cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.3` | Sub 2 | Inside the wholesale move; called out as a row because the STATUS state machine is the most-cited piece of cross-repo doctrine and Sub 6's reference sweep depends on the anchor staying stable. |
| `cnos.cdd/skills/cdd/cross-repo/SKILL.md §2.5` (bundle file set per case a/b/c/d) | `cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.5` | Sub 2 | Verbatim. |
| `cnos.cdd/skills/cdd/cross-repo/SKILL.md §2.6` (LINEAGE.md schema per case) | `cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.6` | Sub 2 | Verbatim. |
| `cnos.cdd/skills/cdd/cross-repo/SKILL.md §2.7` (feedback-patch format) | `cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.7` | Sub 2 | Verbatim. |
| `cnos.cdd/skills/cdd/cross-repo/SKILL.md §2.8` (bundle archival rule; asymmetric) | `cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.8` | Sub 2 | Verbatim. |
| `cnos.cdd/skills/cdd/cross-repo/SKILL.md §2.9` (hat-collapse attribution) | `cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.9` | Sub 2 | Verbatim. |

**Open question for Sub 2 dispatcher:** Does the wholesale move stay one file, or does Sub 2 split the skill (e.g., extract `cross-repo/status/SKILL.md` for the STATUS state machine, `cross-repo/bundles/SKILL.md` for the file sets, `cross-repo/lineage/SKILL.md` for the LINEAGE schemas)? Recommendation: wholesale move for v0.1; consider splitting in a post-#404 cycle once consumer pressure reveals which sub-surfaces churn independently.

---

## 2. Dispatch-prompt template (§Dispatch) — Sub 3

**Source family:** [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.5`](../../cnos.cdd/skills/cdd/gamma/SKILL.md) "Steps 3a–5 — Create the cycle branch, dispatch, unblock without leakage" (lines 170–278 at source pin).
**Source content:** γ-owned branch pre-flight gate; pre-dispatch γ scaffold check (binding gate per `gamma-scaffold.md` presence); polling cross-reference into `harness/SKILL.md §5.4`; dispatch prompt formats per role (γ prompt at lines 204–209; α prompt at lines 211–218; β prompt at lines 252–258); implementation-contract section markdown template (the 7-axes table at lines 226–238); prompt rules; spec-staleness propagation sub-section (lines 272–278); the four-surface mesh declaration (γ template ↔ δ enrichment ↔ α constraint ↔ β verification).
**Migration sub:** Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404).
**Destination commitment:** `cnos.handoff/skills/handoff/dispatch/SKILL.md` — new sub-skill file. Splits cleanly: the dispatch-prompt format is wire-format; the γ-owned branch pre-flight gate is cycle-lifecycle (stays in cdd) and is cited from the new dispatch skill rather than duplicated.

| Source (path + §-anchor) | Destination (path under cnos.handoff/) | Sub | Note |
|---|---|---|---|
| `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Dispatch prompts"` (γ / α / β prompt formats) | `cnos.handoff/skills/handoff/dispatch/SKILL.md §"Prompt formats"` | Sub 3 | Verbatim move; the prompt format is wire-format (operator → α/β/γ/δ via cn dispatch). The line-176 γ-owned branch pre-flight (`cycle/{N}` from `origin/main`) stays in `cdd/gamma/SKILL.md` because it is cycle-lifecycle, not wire-format; Sub 3 cites it from the new dispatch skill. |
| `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Implementation contract section"` (the markdown template at lines 226–238) | `cnos.handoff/skills/handoff/dispatch/SKILL.md §"Implementation contract"` | Sub 3 | See row §3 below for the schema (the 7 axes themselves) — this row is the *template injection rule* (γ MUST include the section in the α prompt; γ MUST NOT dispatch with empty / TBD rows). The schema row and template row land together in `dispatch/SKILL.md`. |
| `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Prompt rules"` (point both roles at the issue; include explicit `Branch:` line; use `--json title,body,state,comments`; do not restate the algorithm in the prompt; do not smuggle constraints into chat prose) | `cnos.handoff/skills/handoff/dispatch/SKILL.md §"Prompt rules"` | Sub 3 | Verbatim move. These rules are wire-format discipline (what the prompt envelope must contain). |
| `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Spec-staleness propagation"` (lines 272–278; how γ propagates spec changes to long-lived polling sessions) | `cnos.handoff/skills/handoff/dispatch/SKILL.md §"Spec-staleness propagation"` (or folded into mid-flight per row §5) | Sub 3 / Sub 4 (boundary contested) | Wire-format aspect (gamma-coordination.md write under `.cdd/unreleased/{N}/`) is handoff-class; the cycle-lifecycle aspect (which skill files trigger propagation — `CDD.md`, role skills, `release/SKILL.md`, `review/SKILL.md`) stays in cdd. Recommendation: fold into mid-flight (Sub 4) so all "γ writes to in-flight α/β" mechanisms live in one place. |
| `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "The mesh"` (the four-surface mesh declaration; γ template ↔ δ enrichment ↔ α constraint ↔ β verification) | `cnos.handoff/skills/handoff/dispatch/SKILL.md §"Four-surface mesh"` | Sub 3 | Verbatim move with cross-references re-pointed: the δ enrichment cite (`delta/SKILL.md §2`) becomes a cnos.handoff-internal cite once row §3 also lands in `dispatch/SKILL.md`; the α constraint cite (`alpha/SKILL.md §3.6`) stays a cdd cite (α is consumer-side); the β verification cite (`beta/SKILL.md` Role Rules Rule 7) stays a cdd cite. |
| `cnos.cdd/skills/cdd/operator/SKILL.md §3a` (relocation pointer; "see `delta/SKILL.md §2`") | `cnos.handoff/skills/handoff/dispatch/SKILL.md` (citation only; `§3a` itself is preserved as a redirect in cdd or deleted) | Sub 3 | The §3a anchor is the pre-Phase-4a name; after `cnos#397` the canonical content lives in `cdd/delta/SKILL.md §2` (see row §3). The §3a redirect is preserved or deleted at Sub 3's discretion; Sub 6 sweeps any remaining cites. |

---

## 3. Implementation-contract schema (§Implementation-contract) — Sub 3

**Source family:** [`cnos.cdd/skills/cdd/delta/SKILL.md §2`](../../cnos.cdd/skills/cdd/delta/SKILL.md) "Inward membrane — γ contract → α-ready dispatch (implementation-contract enrichment)" (canonical home post-`cnos#397`); plus the markdown template in [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Implementation contract section"`](../../cnos.cdd/skills/cdd/gamma/SKILL.md) (lines 226–238 at source pin).
**Source content:** The 7 architectural axes (Language; CLI integration target; Package scoping; Existing-binary disposition; Runtime dependencies; JSON/wire contract preservation; Backward-compat invariant); δ's review-before-routing duty; fill-or-escalate paths; the four-surface mesh (γ template / δ enrichment / α constraint / β verification — see row §2 for the mesh declaration); cnos#389/#391/#392/#393 empirical anchors.
**Migration sub:** Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404).
**Destination commitment:** `cnos.handoff/skills/handoff/dispatch/SKILL.md §"Implementation contract"` — co-resident with the prompt-template row §2 because the schema and the template land together. Optionally lifted to `schemas/handoff/implementation-contract.cue` if Sub 3 chooses to type the axes.

| Source (path + §-anchor) | Destination (path under cnos.handoff/) | Sub | Note |
|---|---|---|---|
| `cnos.cdd/skills/cdd/delta/SKILL.md §2` (inward-membrane doctrine; the 7 axes; δ's review-before-routing duty) | `cnos.handoff/skills/handoff/dispatch/SKILL.md §"Implementation contract" → "Schema"` and `→ "Inward-membrane enrichment"` | Sub 3 | Verbatim move; the 7-axis table itself is wire-format; the "δ-as-inward-membrane" role doctrine moves with it because it is the *enforcement mechanism* for the axes (without δ enrichment the axes are empty / TBD on dispatch). Sub 3 decides whether to keep the role doctrine here (one place) or split it (axes in handoff, doctrine in cdd delta with citation). Recommendation: keep together; the axes only function if the enrichment doctrine is co-located. |
| `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Implementation contract section"` template (lines 226–238; the markdown table γ injects into the α prompt) | `cnos.handoff/skills/handoff/dispatch/SKILL.md §"Implementation contract" → "Prompt-injection template"` | Sub 3 | Verbatim move; the template is the wire-format rendering γ writes into the α prompt. (Also covered as a row in §2 above; this row is the schema-level commitment that the template is canonicalized in cnos.handoff.) |
| `cnos.cdd/skills/cdd/delta/SKILL.md §2.1` (the four-surface mesh, δ side) | `cnos.handoff/skills/handoff/dispatch/SKILL.md §"Implementation contract" → "Four-surface mesh"` | Sub 3 | Verbatim; folds with the γ-side mesh row in §2. |
| `cnos.cdd/skills/cdd/delta/SKILL.md §2.3` (empirical anchors: cnos#389 Python-not-Go; cnos#391 wrong package scoping; cnos#392 first δ-pinned cycle; cnos#393 first-class template) | `cnos.handoff/skills/handoff/dispatch/SKILL.md §"Implementation contract" → "Empirical anchor"` | Sub 3 | Verbatim. |
| `cnos.cdd/skills/cdd/alpha/SKILL.md §3.6` (α constraint: do not improvise; surface unpinned rows to γ via `gamma-clarification.md`) — referenced, not migrated | (no destination; stays in cdd; cited from cnos.handoff dispatch skill) | Sub 3 | Pointer-only. The α constraint is consumer-side cycle-lifecycle doctrine (how α reads the implementation contract); the wire format itself is in cnos.handoff. Sub 3 adds the citation; Sub 6 sweeps. |
| `cnos.cdd/skills/cdd/beta/SKILL.md` Role Rules Rule 7 (β verification: verify diff conforms to every pinned axis; non-conformance → RC severity D classification `implementation-contract`) — referenced, not migrated | (no destination; stays in cdd; cited from cnos.handoff dispatch skill) | Sub 3 | Pointer-only. Same reasoning as α constraint — consumer-side cycle-lifecycle doctrine. |
| (Optional) lift the 7 axes from Markdown table into typed schema at `schemas/handoff/implementation-contract.cue` | `schemas/handoff/implementation-contract.cue` | Sub 3 (deferred decision) | Open question Q2 from `gamma-scaffold.md`. If Sub 3 types the axes, the Markdown template references the schema by `$ref`; if Sub 3 leaves them as Markdown, the schema directory is not created. Recommendation: defer to Sub 3 dispatch; if Sub 3 stays B-lite, leave as Markdown. |

---

## 4. Mid-flight `gamma-clarification.md` mechanism (§Mid-flight) — Sub 4

**Source family:** Currently empirical-only — no canonical home. Cited from:
- [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Implementation contract section"`](../../cnos.cdd/skills/cdd/gamma/SKILL.md) ("Mid-cycle contract changes are logged in `.cdd/unreleased/{N}/gamma-clarification.md`"; line 242 at source pin)
- [`cnos.cdd/skills/cdd/cross-repo/SKILL.md`](../../cnos.cdd/skills/cdd/cross-repo/SKILL.md) (cited as a cache-bust mechanism for issue-edit polling)
- [`cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Mid-flight clarification"`](../../cnos.cds/skills/cds/CDS.md) (post-`cnos#411` canonical home for the issue-edit cache-bust aspect)

**Empirical anchor:** [cnos#391](https://github.com/usurobor/cnos/issues/391) rescue, 2026-05-21 — the cycle where the mechanism first crystallized as a γ→in-flight-α channel.
**Source content (synthesized; no canonical home yet):** file path (`.cdd/unreleased/{N}/gamma-clarification.md` on `cycle/{N}`); authoring role (γ or operator); reader role (in-flight α or β); trigger conditions (axis re-pin, scope clarification, AC clarification, blocker resolution); semantics (in-flight α/β reads the file at next polling cycle; counts as a wake-up event); cache-bust function (an issue-body edit propagates only when an in-flight session refreshes its `gh issue view` cache, and a `gamma-clarification.md` write is the explicit signal to refresh).
**Migration sub:** Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404).
**Destination commitment:** `cnos.handoff/skills/handoff/mid-flight/SKILL.md` — new sub-skill, authored from the empirical anchor.

| Source (path + §-anchor) | Destination (path under cnos.handoff/) | Sub | Note |
|---|---|---|---|
| (no canonical home) cnos#391 empirical anchor | `cnos.handoff/skills/handoff/mid-flight/SKILL.md` (entire file; new authoring) | Sub 4 | Rewrite (authoring from empirical anchor; not a verbatim move). Sub 4 distills the post-cnos#391 practice across cycles into a canonical surface. The file path / authoring role / reader role / trigger conditions / semantics / cache-bust function are the section structure. |
| `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5` cite (line 242: "Mid-cycle contract changes are logged in `.cdd/unreleased/{N}/gamma-clarification.md`") | (cite migrates: cdd cite now points at cnos.handoff/skills/handoff/mid-flight/SKILL.md) | Sub 4 (cite update); Sub 6 (sweep) | Pointer-only on the cdd side; Sub 6 sweeps the cite during cross-reference cleanup. |
| `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Mid-flight clarification"` (post-#411 canonical for issue-edit cache-bust aspect) | `cnos.handoff/skills/handoff/mid-flight/SKILL.md §"Cache-bust function"` | Sub 4 | Migration-semantics-undecided: Sub 4 picks among (a) move the cds section into cnos.handoff (writes both packages — re-evaluate the boundary in Sub 4's β-review); (b) keep canonical in cnos.cds and add pointer mirror in cnos.handoff (preserves CDS untouched); (c) cite-only from cnos.handoff loaders (most conservative). Recommendation: (b) — the mid-flight mechanism is wire-format-shaped, so cnos.handoff is the right canonical home, and cnos.cds becomes a pointer. But this is Sub 4's call; flagged here for the dispatching δ. |
| (folded) `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Spec-staleness propagation"` (lines 272–278; the gamma-coordination.md cousin of gamma-clarification.md for long-lived polling sessions) | `cnos.handoff/skills/handoff/mid-flight/SKILL.md §"Spec-staleness propagation"` | Sub 4 | Folded here per row §2's recommendation; all "γ writes to in-flight α/β" mechanisms in one place. |

---

## 5. `.cdd/unreleased/{N}/` artifact channel rules (§Artifact-channel) — Sub 4

**Source family:** Post-`cnos#411` canonical home is [`cnos.cds/skills/cds/CDS.md §"Artifact contract"`](../../cnos.cds/skills/cds/CDS.md) — specifically the sub-sections `→ "Location matrix"` (canonical paths for `self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, `cdd-iteration.md`, `RELEASE.md`, the version-snapshot directory, the PRA, the cross-repo trace dir), `→ "Ownership matrix"` (per-role artifact ownership: α owns self-coherence + alpha-closeout; β owns beta-review + beta-closeout; γ owns gamma-closeout + scaffold + cdd-iteration), `→ "Ordered flow"` (the 13-stage flow), and `→ "Frozen snapshot rule"` (artifacts in `.cdd/unreleased/{N}/` freeze on merge).
**Pre-#411 home:** `cnos.cdd/skills/cdd/CDD.md` (pre-#402; the channel rules were folded into the CCNF kernel before the cnos#403 wave extracted them to cds).
**Source content:** Sequential α→β→γ handoff over the cycle directory; per-role write ownership (α writes alpha-closeout; β writes beta-review then beta-closeout; γ writes gamma-scaffold + gamma-closeout + cdd-iteration); freeze-on-merge semantics; release-time move from `.cdd/unreleased/{N}/` to `.cdd/releases/{X.Y.Z}/{N}/` (or `.cdd/releases/docs/{ISO-date}/{N}/` for docs-only releases).
**Migration sub:** Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404).
**Destination commitment:** `cnos.handoff/skills/handoff/artifact-channel/SKILL.md` — new sub-skill. The boundary between "channel rules" (wire-format; handoff) and "what each artifact contains" (cycle-lifecycle; cds-resident) is the load-bearing decision Sub 4 makes.

| Source (path + §-anchor) | Destination (path under cnos.handoff/) | Sub | Note |
|---|---|---|---|
| `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Location matrix"` (canonical paths under `.cdd/unreleased/{N}/`) | `cnos.handoff/skills/handoff/artifact-channel/SKILL.md §"Location matrix"` | Sub 4 | Migration-semantics-undecided. The path *patterns* (`.cdd/unreleased/{N}/<artifact>.md`) are wire-format; the *specific artifacts* (self-coherence; beta-review; etc.) are cycle-lifecycle. Sub 4 splits the section: path-pattern + frozen-snapshot semantics → cnos.handoff; per-artifact contents → stays in cnos.cds. Recommendation: pointer-only from cnos.handoff; the cds section is the canonical home for the matrix itself, and cnos.handoff cites it. |
| `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Ownership matrix"` (α/β/γ/δ/ε per-artifact ownership) | `cnos.handoff/skills/handoff/artifact-channel/SKILL.md §"Ownership matrix"` (pointer) | Sub 4 | Pointer-only. The ownership matrix is consumer-protocol-specific (CDS-shaped; CDR has different artifacts; future c-d-X may have other artifacts). Wire-format-side cnos.handoff records the *pattern* (every cycle dir has a per-role-owned writer); consumer-side cds records the *specific assignment*. |
| `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Ordered flow"` (13-stage flow: design → contract → plan → tests → code → docs → self-coherence → review → gate → release → observe → assess → close) | `cnos.handoff/skills/handoff/artifact-channel/SKILL.md §"Sequential rule"` (pointer or rewrite) | Sub 4 | Pointer-only or thin rewrite. The 13-stage *ordering* is CDS-specific (research-class CDR's flow differs; future c-d-X may have other flows). cnos.handoff records the *invariant* (one role at a time writes; the predecessor role's artifact is the successor's input); consumer-side cds records the *specific 13-stage flow*. |
| `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Frozen snapshot rule"` (artifacts in `.cdd/unreleased/{N}/` freeze on merge; release-time move discipline) | `cnos.handoff/skills/handoff/artifact-channel/SKILL.md §"Frozen snapshot rule"` | Sub 4 | Verbatim move — this is wire-format (post-merge immutability of the channel directory). Sub 4 decides whether the cds section becomes a pointer (canonical at handoff) or stays as cds-side canonical with a handoff pointer (cite-only). Recommendation: (a) — frozen-snapshot rule is wire-format, move it to handoff; the release-time directory rename is also wire-format. |

**Open question for Sub 4 dispatcher:** The cnos.cds rows above are flagged migration-semantics-undecided. The dispatching δ for Sub 4 picks (a) move; (b) mirror; (c) cite-only **per row** (not as a blanket policy). The Location matrix + Ownership matrix should likely be (c) — consumer-specific. The Frozen snapshot rule + path-pattern invariant should likely be (a) — wire-format. The Ordered flow should likely be (b) — invariant in handoff, specific assignment in cds. This is judgment; the map records the options.

---

## 6. `cdd-iteration.md` receipt-stream + INDEX.md aggregator (§Receipt-stream) — Sub 5

**Source family:** [`cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b`](../../cnos.cdd/skills/cdd/post-release/SKILL.md) "Author `cdd-iteration.md` (when applicable)" (lines 293–340 at source pin) + [`.cdd/iterations/INDEX.md`](../../../../.cdd/iterations/INDEX.md) (the aggregator file itself; one row per cycle).
**Source content:** Cadence rule (`cdd-iteration.md` required only when `protocol_gap_count > 0`); per-finding shape (F1, F2, … with Source / Class / Trigger / Description / Root cause / Disposition fields; sub-fields per disposition: patch-landed / next-MCA / no-patch); aggregator-update procedure (one row per cycle in INDEX.md with Cycle / Issue / Date / Findings / Patches / MCAs / No-patch / Path columns); cross-repo trace bundle (when a finding's disposition is `patch-landed` with `Cross-repo` target, γ creates a bundle at `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` per cross-repo skill); empty-findings courtesy stub allowance (backward compatibility).
**Migration sub:** Sub 5 of [cnos#404](https://github.com/usurobor/cnos/issues/404).
**Destination commitment:** `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` — new sub-skill. The receipt-stream is cross-cycle handoff (one cycle's γ writes; the ε reader at any later cycle consumes); fits cleanly in cnos.handoff.

| Source (path + §-anchor) | Destination (path under cnos.handoff/) | Sub | Note |
|---|---|---|---|
| `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b → "Cadence rule"` (`cdd-iteration.md` required only when `protocol_gap_count > 0`) | `cnos.handoff/skills/handoff/receipt-stream/SKILL.md §"Cadence rule"` | Sub 5 | Verbatim move; the cadence rule is wire-format (when γ MUST write the receipt). Note: the `cdd-iteration.md` filename itself becomes consumer-specific in time (cds has its own `cds-iteration.md` analogue per `cnos.cds/skills/cds/CDS.md §"Closure"`); the wire-format rule is "the cycle's iteration receipt" regardless of filename. |
| `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b → "Per-finding shape"` (lines 305–326: F1, F2, … with Source / Class / Trigger / Description / Root cause / Disposition + sub-fields) | `cnos.handoff/skills/handoff/receipt-stream/SKILL.md §"Per-finding shape"` | Sub 5 | Verbatim move. Could be lifted to typed schema at `schemas/handoff/iteration-finding.cue` if Sub 5 chooses; defer that to Sub 5 dispatch. |
| `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b → "Aggregator update"` (the INDEX.md row format) | `cnos.handoff/skills/handoff/receipt-stream/SKILL.md §"Aggregator update"` | Sub 5 | Verbatim move. |
| `cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b → "Cross-repo trace"` (bundle creation when patch-landed has Cross-repo target) | `cnos.handoff/skills/handoff/receipt-stream/SKILL.md §"Cross-repo trace"` | Sub 5 | Verbatim move; the bundle-creation rule is wire-format (it composes with the cross-repo skill row §1). The cross-repo skill cite migrates from `cnos.cdd/skills/cdd/cross-repo/SKILL.md` to `cnos.handoff/skills/handoff/cross-repo/SKILL.md` (the latter is the Sub 2 destination; Sub 5 lands after Sub 2 so the cite resolves to the new home). |
| `.cdd/iterations/INDEX.md` (the aggregator file itself; 23+ rows at source pin; one per cycle) | (file stays at `.cdd/iterations/INDEX.md`; not migrated) | Sub 5 | The aggregator file is *data*, not doctrine; it does not move. The doctrine (how to update it) moves. The filename will eventually become consumer-specific (the cds analogue would be `.cds/iterations/INDEX.md`); cnos.handoff records the *pattern* (per-protocol-package aggregator at a known path). |
| `cnos.cds/skills/cds/CDS.md §"Closure"` (cds-side cycle-iteration trigger references that cite §5.6b in cdd) | (cite migrates: cds cite now points at cnos.handoff/skills/handoff/receipt-stream/SKILL.md) | Sub 5 (cite update); Sub 6 (sweep) | Pointer-only on the cds side; Sub 6 sweeps. |

---

## 7. Polling primitives (§Polling) — Sub 5 or deferred (boundary contested)

**Source family:** [`cnos.cdd/skills/cdd/harness/SKILL.md §5.4`](../../cnos.cdd/skills/cdd/harness/SKILL.md) (single-named-branch transition loop; `Monitor`-wrapped polling; reachability re-probe; transition-only emission). Post-`cnos#411` canonical doctrine is in [`cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"`](../../cnos.cds/skills/cds/CDS.md); operational realization remains in `cnos.cdd/skills/cdd/harness/SKILL.md §5.4`.
**Source content:** Polling query forms (gh / MCP / git); wake-up mechanism; reachability preflight; transition-only emission discipline; synchronous baseline pull; `git fetch` reliability re-probe.
**Migration sub:** Sub 5 of [cnos#404](https://github.com/usurobor/cnos/issues/404) (or deferred to a post-#404 cycle).
**Destination commitment:** **Open.** The polling primitives are wire-format (γ observes α/β progress on the cycle branch); but the operational realization is runtime-substrate (harness mechanics). Recommendation: defer — leave canonical in cnos.cds + operational in cnos.cdd/harness; do not duplicate in cnos.handoff for v0.1. Re-evaluate when a non-CDD/CDS consumer (CDR or future c-d-X) needs to reuse the primitives.

| Source (path + §-anchor) | Destination (path under cnos.handoff/) | Sub | Note |
|---|---|---|---|
| `cnos.cdd/skills/cdd/harness/SKILL.md §5.4` (operational mechanics) | (no destination; stays in cdd as runtime substrate) | Sub 5 or deferred | Pointer-only. Cite from cnos.handoff loaders if needed. |
| `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"` (post-#411 canonical doctrine) | `cnos.handoff/skills/handoff/polling/SKILL.md` (deferred) | Sub 5 or post-#404 | Migration-semantics-undecided; recommend deferring until a non-CDS consumer surfaces. The doctrine is wire-format-shaped but the consumer count is still one (CDS); design-skill §3.5 ("interfaces belong to consumers") says wait for substitution pressure. |

---

## 8. Spec-staleness propagation (§Spec-staleness) — Sub 4 (folded)

**Source family:** [`cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Spec-staleness propagation"`](../../cnos.cdd/skills/cdd/gamma/SKILL.md) (lines 272–278 at source pin).
**Source content:** Identity-rotation mode (each role invocation loads skills fresh from the filesystem; the next dispatch picks them up); long-lived polling sessions (γ writes a coordination note to `.cdd/unreleased/{N}/gamma-coordination.md` on `cycle/{N}` naming the spec change and affected skill path); when to propagate (changes to `CDD.md`, role skill files, `release/SKILL.md`, `review/SKILL.md` landing on `main` while in-flight); when not to propagate (changes outside the cdd package; doc-only; issues/threads).
**Migration sub:** Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404) — folded into mid-flight (row §4) because the mechanism is a γ-writes-to-in-flight-α/β channel.
**Destination commitment:** `cnos.handoff/skills/handoff/mid-flight/SKILL.md §"Spec-staleness propagation"` (already a row in §4 above; cross-listed here for explicit coverage of the additional discovered surface).

| Source (path + §-anchor) | Destination (path under cnos.handoff/) | Sub | Note |
|---|---|---|---|
| `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Spec-staleness propagation"` (lines 272–278) | `cnos.handoff/skills/handoff/mid-flight/SKILL.md §"Spec-staleness propagation"` | Sub 4 | Verbatim move; the "when to propagate" rule (changes to specific cdd files) becomes consumer-specific (CDD-shaped); the *invariant* (γ writes a coordination note when spec changes mid-flight) is wire-format. Sub 4 splits: invariant in handoff; specific file list in cdd-consumer-side. |

---

## 9. Cross-reference cleanup + close tracker (§Cleanup) — Sub 6

**Source family:** N/A — this is the cleanup sub, not an extraction sub.
**Source content:** All cross-references in `cnos.cdd/`, `cnos.cdr/`, `cnos.cds/`, and `docs/gamma/` essays that cite handoff-class surfaces. Per-citation update: cite the new `cnos.handoff/skills/handoff/<area>/SKILL.md` canonical home rather than the pre-migration cdd / cds path.
**Migration sub:** Sub 6 of [cnos#404](https://github.com/usurobor/cnos/issues/404).
**Destination commitment:** N/A — Sub 6 modifies citations in cdd / cdr / cds / docs; does not author new cnos.handoff surfaces.

| Source (path + §-anchor) | Destination (path under cnos.handoff/) | Sub | Note |
|---|---|---|---|
| All `cross-repo/SKILL.md` cites in `cnos.cdr/skills/cdr/*.md`, `cnos.cds/skills/cds/CDS.md`, `cnos.cdd/skills/cdd/{gamma,post-release,issue}/SKILL.md` | (cites re-pointed at `cnos.handoff/skills/handoff/cross-repo/SKILL.md`) | Sub 6 | Mechanical sweep after Sub 2 lands. |
| All `gamma/SKILL.md §2.5` dispatch-prompt cites in `cnos.cdd/skills/cdd/{alpha,beta,operator,delta,harness}/SKILL.md` | (cites re-pointed at `cnos.handoff/skills/handoff/dispatch/SKILL.md`) | Sub 6 | Mechanical sweep after Sub 3 lands. |
| All `delta/SKILL.md §2` implementation-contract cites in `cnos.cdd/skills/cdd/{alpha,beta,gamma,operator}/SKILL.md` | (cites re-pointed at `cnos.handoff/skills/handoff/dispatch/SKILL.md §"Implementation contract"`) | Sub 6 | Mechanical sweep after Sub 3 lands. |
| All `gamma-clarification.md` cites in `cnos.cdd/`, `cnos.cds/`, `cnos.cdr/` | (cites re-pointed at `cnos.handoff/skills/handoff/mid-flight/SKILL.md`) | Sub 6 | Mechanical sweep after Sub 4 lands. |
| All `post-release/SKILL.md §5.6b` cites + INDEX.md format cites | (cites re-pointed at `cnos.handoff/skills/handoff/receipt-stream/SKILL.md`) | Sub 6 | Mechanical sweep after Sub 5 lands. |
| Close [cnos#404](https://github.com/usurobor/cnos/issues/404) — name the wave complete; cite the six Subs (1–6) and their landed commits | (close tracker comment + STATUS update) | Sub 6 | Wave closeout artifact. |

---

## 10. Coverage verification

The six required surface families per [cnos#415](https://github.com/usurobor/cnos/issues/415) D2 + the two additional discovered surfaces + the close-out row are all covered:

| Required surface family per cnos#415 D2 | Covered in section § |
|---|---|
| 1. Cross-repo state machine + bundles (Sub 2) | §1 |
| 2. Dispatch-prompt template (Sub 3) | §2 |
| 3. Implementation-contract schema (Sub 3) | §3 |
| 4. Mid-flight `gamma-clarification.md` mechanism (Sub 4) | §4 |
| 5. `.cdd/unreleased/{N}/` artifact channel rules (Sub 4) | §5 |
| 6. `cdd-iteration.md` receipt-stream + INDEX.md aggregator (Sub 5) | §6 |
| (additional) Polling primitives (Sub 5 or deferred) | §7 |
| (additional) Spec-staleness propagation (Sub 4 folded) | §8 |
| Sub 6 — cross-reference cleanup + close tracker | §9 |

Mechanical check: `grep -c "^## " src/packages/cnos.handoff/docs/extraction-map.md` returns ≥ 6 (currently 10 `##` headings under §0–§10 in this file).

Sub 6 of [cnos#404](https://github.com/usurobor/cnos/issues/404) will, on completion of Subs 2–5, sweep all handoff-class citations in cdd / cdr / cds / docs against this map: each citation either resolves (content migrated to the named cnos.handoff destination) or updates (citation re-pointed at the cnos.handoff surface named above). No citation is silently dropped.

## 11. Open questions

These are noted for the operator dispatching Subs 2–6 against this map; none are blockers for Sub 1.

- **Q1 (Sub 2 scope).** Does Sub 2 move `cross-repo/SKILL.md` wholesale (644 lines as one file), or split into `cross-repo/{status,bundles,lineage,feedback,archival}/SKILL.md` sub-files? Recommendation: wholesale for v0.1; consider splitting if a post-#404 cycle reveals which sub-surfaces churn independently.
- **Q2 (Sub 3 schema lift).** Does Sub 3 lift the implementation-contract 7-axes table to a typed schema at `schemas/handoff/implementation-contract.cue`, or keep as Markdown table? Recommendation: defer to Sub 3 dispatch; if B-lite, leave as Markdown.
- **Q3 (Sub 4 unification).** Does Sub 4 unify mid-flight + artifact-channel into one `skills/handoff/intra-cycle/SKILL.md`, or keep them split? Recommendation: split per the per-mechanism cohesion; mid-flight is asynchronous γ→in-flight-α; artifact-channel is sequential α→β→γ. Different rates of fire.
- **Q4 (Sub 5 polling absorption).** Does Sub 5 absorb `harness/SKILL.md §5.4` polling primitives, or leave as cnos.cdd runtime substrate? Recommendation: leave for v0.1; re-evaluate when CDR or future c-d-X needs the primitives.
- **Q5 (Sub 4 + Sub 5 cds-source rows).** The rows in §4 and §5 above marked migration-semantics-undecided each have a recommended choice (move / mirror / cite-only) in the row's note column; the dispatching δ for each sub makes the final call per row. The recommendation is not a pre-commitment.
- **Q6 (Sub 6 scope of sweep).** Sub 6's citation sweep covers cnos.cdd / cnos.cdr / cnos.cds skill files. Does it also cover `docs/gamma/essays/*.md` and `docs/gamma/cdd/RATIONALE.md`? Recommendation: yes — any essay that cites a handoff-class surface by old path becomes stale; the sweep should be exhaustive.
- **Q7 (HANDOFF.md authoring sub).** Sub 2 is the largest extraction (cross-repo wholesale); should Sub 2 also land `HANDOFF.md` (the canonical wire-format contract), or should `HANDOFF.md` wait for Sub 3 (dispatch — which is the first surface to span multiple sub-skills and benefits from a shared contract)? Recommendation: Sub 3 lands `HANDOFF.md` because it is the first sub that needs a cross-sub-skill consensus on directional cases. Sub 2 leaves `HANDOFF.md` empty (or absent); Sub 3 lands it.

These are migration-coordination questions, not destination-uncertainty questions. The destinations named in the per-surface tables above are stable commitments for Subs 2–6; the open questions are about *how* the migration executes (one file vs split; schema lift vs Markdown; combine vs separate; sweep scope), not *where* the canonical home lands.

---

**End of extraction map.**
