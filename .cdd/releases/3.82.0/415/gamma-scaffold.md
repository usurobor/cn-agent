# γ scaffold — cycle/415

**Issue:** [cnos#415](https://github.com/usurobor/cnos/issues/415) — Sub 1 of [#404](https://github.com/usurobor/cnos/issues/404): Bootstrap `cnos.handoff` package skeleton + extraction map
**Branch:** `cycle/415` (from `origin/main` at `bc29c009` — Merge cycle/414)
**Mode:** small-change — docs-only + package skeleton + extraction map; no runtime/harness/validator touch
**Collapse:** β-α-collapse-on-δ; dispatched as `γ+α+β collapsed on δ` per docs-class cycle pattern (cnos.cds §"Field 6: Actor collapse rule")

## Surfaces γ expects α to touch

| Path | Author |
|---|---|
| `src/packages/cnos.handoff/cn.package.json` | new file, verbatim shape from #415 |
| `src/packages/cnos.handoff/README.md` | new, ≥ 50 lines; mirrors cnos.cds README section shape |
| `src/packages/cnos.handoff/skills/handoff/SKILL.md` | new, mirrors cnos.cds/cnos.cdr loader-skill frontmatter |
| `src/packages/cnos.handoff/skills/handoff/.gitkeep` | placeholder; HANDOFF.md omitted (Subs 2–5 land it) |
| `src/packages/cnos.handoff/docs/extraction-map.md` | new; ≥ 6 surface-family tables; load-bearing artifact for Subs 2–6 |
| `.cdd/unreleased/415/{gamma-scaffold,self-coherence,beta-review,alpha-closeout,beta-closeout,gamma-closeout,cdd-iteration}.md` | cycle close-outs |
| `.cdd/iterations/INDEX.md` | append one row |

## AC oracle approach

All 9 ACs are mechanical or read-checks:

- **AC1** — file presence + line count: `test -f <path>`; `wc -l README.md` ≥ 50.
- **AC2** — `cn` package discovery: `cnos.handoff` discoverable as v1 package; no `cnos.core` changes required.
- **AC3** — read-check on README section shape and content claims.
- **AC4** — read-check on SKILL.md frontmatter (name/description/artifact_class/governing_question/triggers/calls) and section shape (load order / rule / cross-protocol).
- **AC5** — `grep -c "^## " docs/extraction-map.md` ≥ 6; each row's destination resolves under `cnos.handoff/`; each row's sub field names Sub 2 / 3 / 4 / 5.
- **AC6** — three `git diff origin/main..HEAD -- <pkg>` invocations: 0 lines for cnos.cdd, cnos.cdr, cnos.cds.
- **AC7** — `test -d schemas/ccnf-o/` returns nonzero.
- **AC8** — `git diff origin/main..HEAD -- src/packages/cnos.cdd/commands/cdd-verify/` returns 0 lines.
- **AC9** — `git diff origin/main..HEAD -- src/go/` returns 0 lines.

## Empirical anchor (§2.2a peer enumeration)

Structural precedent: `cnos.cds` v0.1 bootstrap at cnos#406 (Sub 1 of cnos#403; merged 2026-05-22; commit `378a54f0`). The `cnos.cds/` package — `cn.package.json` (140 bytes; same schema), `README.md` (64 lines), `skills/cds/SKILL.md` (120 lines; same frontmatter shape), `docs/extraction-map.md` (≥ 12 surface tables) — is the verbatim template. cnos#415 sits in the same structural slot in the cnos#404 wave that cnos#406 occupied in the cnos#403 wave.

Secondary precedent: `cnos.cdr` v0.1 at cnos#376 (the original three-package split that motivated #403 and #404).

## Expected diff scope

- ~7 new files under `src/packages/cnos.handoff/` (manifest, README, loader-skill, .gitkeep, extraction-map; plus the two skill dir placeholders if needed)
- ~7 new files under `.cdd/unreleased/415/` (scaffold, self-coherence, beta-review, alpha-closeout, beta-closeout, gamma-closeout, cdd-iteration courtesy stub)
- 1 line added to `.cdd/iterations/INDEX.md`
- **Zero lines changed** in `src/packages/cnos.cdd/`, `src/packages/cnos.cdr/`, `src/packages/cnos.cds/`, `src/go/`, or `schemas/ccnf-o/` (last must remain non-existent).

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Markdown skills + JSON manifest |
| CLI integration target | None new; package loads via existing `cn` package-discovery |
| Package scoping | `src/packages/cnos.handoff/` (new; peer of cnos.cdd / cnos.cdr / cnos.cds) |
| Package name | `cnos.handoff` (pinned; **not** `coordination` / `relay` / `dispatch`) |
| Existing-binary disposition | N/A |
| Runtime dependencies | None |
| JSON/wire contract | `cn.package.json` schema `cn.package.v1`, version `0.1.0`, engines `cnos >= 3.81.0` |
| Backward-compat invariant | cnos.cdd / cnos.cdr / cnos.cds **untouched**; all extraction-target surfaces stay in their current homes until Subs 2–5 move them |

## Authoring judgment calls (recorded for β read-check)

### J1: HANDOFF.md present-or-absent

**Decision:** Omit `HANDOFF.md` for v0.1; use `.gitkeep` placeholder under `skills/handoff/` to keep the directory tracked.

**Reasoning:** The cnos.cds Sub 1 precedent (cnos#406) shipped without `CDS.md`; Sub 2 (cnos#407) landed it. cnos#415 mirrors that exactly. A 1-paragraph stub would (a) duplicate the "forthcoming" disclaimer already in SKILL.md and README.md, (b) require Sub 2/3/4/5 to overwrite rather than create, (c) tempt a sub-author to treat the stub as a constraint. `.gitkeep` is the cleaner option; the SKILL.md `calls:` block already names the forthcoming surfaces as advisory targets.

### J2: Extraction-map source paths after post-#411 reality

**Discovery:** Post-#403 wave closure (cnos#411), several extraction-target surfaces named in cnos#415's D2 are no longer in `cnos.cdd` — they migrated to `cnos.cds/skills/cds/CDS.md`:
- "`.cdd/unreleased/{N}/` artifact channel rules" → now `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Location matrix" + "Ordered flow" + "Frozen snapshot rule"`
- "cycle-state evidence / mid-flight clarification / cross-repo proposals" → now `cnos.cds/skills/cds/CDS.md §"Coordination surfaces"` (consolidated under §3 of the cds extraction map)
- The pre-#402 `CDD.md` artifact-channel references no longer exist; `CDD.md` now only carries the kernel + pointers into cds.

**Decision:** Record each extraction-map row's source path at its **current** canonical home (which is sometimes `cnos.cds` after #411, not `cnos.cdd`). For rows where the canonical home is `cnos.cds`, flag the migration-semantics question in the row's note column: Sub 4/Sub 5 must decide whether to (a) move the surface from cnos.cds into cnos.handoff (writes both packages — re-evaluate the boundary), (b) keep canonical in cnos.cds and add a pointer-only mirror in cnos.handoff, or (c) re-scope (leave canonical in cnos.cds; do not duplicate; cite from cnos.handoff loaders).

**Why this is the right move now:** cnos#415's job is to ship a *credible map* — not to pre-commit Subs 4/5 to one of (a)/(b)/(c). Recording the post-#411 reality in the note column gives the dispatching δ for each subsequent sub the live evidence they need; pretending the surface is still in cnos.cdd would propagate a stale citation that β-review on Sub 4 would flag immediately.

**The pure-cnos.cdd-source rows are unchanged:** the cross-repo state machine (Sub 2), the dispatch-prompt template + implementation-contract schema in `cdd/gamma/SKILL.md §2.5` + `cdd/delta/SKILL.md §2` (Sub 3), the receipt-stream in `cdd/post-release/SKILL.md §5.6b` (Sub 5), the mid-flight `gamma-clarification.md` mechanism (Sub 4; still empirical-only with no canonical home) all retain their cnos.cdd sources.

### J3: Operator/SKILL.md §3a deprecation

**Discovery:** `cnos.cdd/skills/cdd/operator/SKILL.md §3a` is now a relocation-pointer only — the implementation-contract enrichment doctrine moved to `delta/SKILL.md §2` per cnos#397 (Phase 4a of cnos#366).

**Decision:** Source path in the extraction map for the inward-membrane row reads `cnos.cdd/skills/cdd/delta/SKILL.md §2` (canonical) with a parenthetical note that `operator/SKILL.md §3a` is the pre-Phase-4a anchor name preserved as a redirect. Sub 3 reads `delta/SKILL.md §2` for the canonical content.

### J4: Additional handoff surfaces discovered

While reading the cdd skill files I noted these additional handoff-class surfaces (each recorded as an extraction-map row):

- `cnos.cdd/skills/cdd/harness/SKILL.md §5.4` — polling primitives (gh / MCP / git query forms; wake-up; reachability re-probe; transition-only emission; synchronous baseline pull). Post-#411 the doctrine is canonical in `cnos.cds §"Coordination surfaces" → "Polling primitives"` with operational realization left in `cdd/harness/SKILL.md`. Handoff-class (the primitives are how γ observes α/β progress on the cycle branch). → Sub 5 or Sub 6 (the agent dispatching Sub 5 decides whether to absorb it).
- `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5` "Spec-staleness propagation" sub-bullet — handoff doctrine for cross-activation skill-loading-staleness. → folded into Sub 4 (artifact channel + mid-flight family).
- `cnos.cdd/skills/cdd/cross-repo/SKILL.md §2.5` — bundle file set per case (the case-(a)/(b)/(c)/(d) directional bundle shapes); already inside the Sub 2 wholesale move, so it is not a separate row.

## Open-questions reservoir (for Subs 2–6 to decide)

These are not Sub 1 blockers; they are flagged in the extraction-map Open Questions section:

- **Q1:** Does Sub 2 move `cross-repo/SKILL.md` wholesale (644 lines) or split it (e.g., keep the bundle file set + LINEAGE schemas in cnos.handoff; keep the case-(d) operator-pending workflow in cnos.cdd)?
- **Q2:** Does Sub 3 lift the implementation-contract 7-axes table to a typed JSON schema in `schemas/handoff/`, or keep it as Markdown table?
- **Q3:** Does Sub 4 unify the mid-flight + artifact-channel skills into one (`skills/handoff/intra-cycle/SKILL.md`) or keep them split? The split mirrors cnos#391 + α→β→γ as two distinct mechanisms; the unification mirrors "intra-cycle handoff" as one family.
- **Q4:** Does Sub 5 absorb `harness/SKILL.md §5.4` polling primitives, or leave them as cnos.cdd runtime substrate? (`§5.4` is operationally tied to the harness; the doctrine — query / wake-up / reachability — is handoff-class. The boundary is genuinely contested.)
- **Q5:** Sub 6's scope — cross-reference cleanup in cdd / cds / cdr to re-point at the new cnos.handoff homes. Each canonical surface that moves needs a citation sweep; the sweep is mechanical but the count depends on how Subs 2–5 land.

## Scaffold complete

γ obligations met per `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5`:
- Pre-dispatch scaffold present on `cycle/415` (this file).
- Implementation contract pinned.
- Empirical anchor named.
- Surface map and AC oracle approach declared.
- Mid-cycle changes will be logged in `gamma-clarification.md` if any axis re-pins.

α may proceed with skeleton authoring.
