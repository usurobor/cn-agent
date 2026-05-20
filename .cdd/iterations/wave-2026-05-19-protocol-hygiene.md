# ε iteration — Wave 2026-05-19 Protocol Hygiene (launch-side)

**Wave:** `.cdd/waves/2026-05-19-protocol-hygiene/`
**Wave master tracker:** none (the wave is three independent cycles, not a master+subs shape)
**ε actor:** cnos γ acting as ε per role-collapse rule (`epsilon/SKILL.md §2`); user-delegated this session
**Output type:** launch-side ε iteration. The close-side ε work (consolidating per-cycle findings after the three cycles merge) writes into this same file; this commit lands the launch-side findings only.

## §1 Why this iteration exists

The user delegated "Run epsilon on cph cdd... Launch waves as epsilon launching delta to launch gammas and then reporting cdd iterations if needed" (session 2026-05-19). ε's responsibilities for a wave-launch include:

- Observe whether the cdd protocol's roadmap and surfaces are coherent with the current canonical doctrine (`COHERENCE-CELL.md`, `COHERENCE-CELL-NORMAL-FORM.md`, `ROLES.md`).
- Collect `cdd-*-gap` findings from prior cycle close-outs and from the wave-planning review itself.
- Apply MCA discipline — ship immediate patches when the fix is clear; file `next-MCA` when the pattern is real but the patch shape isn't yet clear; drop only one-off non-patterns.
- Write `cdd-iteration.md` (this file) as the durable record per `epsilon/SKILL.md §1`.

This session reviewed cnos issues #366, #375, #377, #378, #379, #384 + the persona/protocol/project work landed on `claude/file-cnos-cdr-issue-fi9Ld`. Four findings surface.

## §2 Findings dispositioned

### F1: #366 Phase 3 V signature drift (cdd-protocol-gap)

**Source:** ε review of #366 against `COHERENCE-CELL-NORMAL-FORM.md` (which landed at #370 merge `0d9f7498` on 2026-05-17).

**Class:** `cdd-protocol-gap` (doctrine drift between roadmap and ratified normal form).

**Trigger:** ε process-gap check during wave-launch review.

**Description:** #366 Phase 3's deliverable text said `V : Contract × Receipt × EvidenceGraph → ValidationVerdict`. The ratified normal form at `COHERENCE-CELL-NORMAL-FORM.md` (#370) defines `V : Contract × Receipt → ValidationVerdict` — evidence refs are bound *into* the receipt by γ; V dereferences them through the receipt; δ does not read evidence directly. The roadmap text was authored before #370 landed; the doctrine landed; the roadmap text did not update. If Phase 3 had dispatched against the roadmap signature, the implementation would have contradicted the ratified normal form. Avoidable design instability.

**Root cause:** No canonical rule that says "when a normal form / doctrine ratifies, downstream roadmap text must be reconciled before the next cycle in the affected chain dispatches." The roadmap drift class is itself an `cdd-protocol-gap` (see F4 below). For this finding specifically, the rule that would have prevented it is "ε reviews roadmap-vs-doctrine coherence before the next dispatchable phase fires." That's what this session did — but the discipline isn't written down.

**Disposition:** `patch-landed`. #366 body updated this session — Phase 3 V signature corrected to match #370; recursion equation in §"Recursion the implementation must instantiate" rewritten; note added under Phase 2 explaining the lineage.

**Patch:** cnos issue #366 body (modified 2026-05-19 via `mcp__github__issue_write`); see #366 "Update — 2026-05-19" block at top of body for the delta summary.

**Affects:** Phase 3 implementation (now dispatches against the correct signature); Phase 2.5 dispatch (which references the corrected signature); future Phase 6 (ε upscope) and Phase 7 (CDD.md rewrite) — both will read the corrected roadmap.

### F2: #366 Phase 2.5 missing — schemas/cdd/ conflates generic CDD with CDS-specific evidence (cdd-skill-gap + cdd-protocol-gap)

**Source:** ε review of `schemas/cdd/receipt.cue` against the persona/protocol/project enforcement chain (`ROLES.md §4a`, landed this session 2026-05-19).

**Class:** `cdd-skill-gap` (the schema surface is the artifact; the gap is its shape) + `cdd-protocol-gap` (the roadmap missed the phase).

**Trigger:** ε process-gap check during wave-launch review; reinforced by the cnos#376 body modification (this session) which added a CDR receipt schema requiring `claim_refs`, `data_refs`, `method_refs`, `result_refs`, `claim_status`, `limitations`, `reproduction` — fields that have no home in the current `schemas/cdd/`.

**Description:** Phase 2 (#369) shipped `schemas/cdd/receipt.cue` requiring fields like `self_coherence_ref`, `beta_review_ref`, `alpha_closeout_ref`, `beta_closeout_ref`, `gamma_closeout_ref`, `diff_ref`, optional `ci_refs`. These are right for the software-development instantiation (CDS), but they are not generic. A CDR receipt cannot validate against these fields — its required fields are claim/data/method-shaped. Without a generic/domain split, Phase 3's V would either bake CDS-specific requirements into the generic validator (CDR can't use it) or strip CDS-specific requirements to a permissive base (V can't actually validate anything). Neither is acceptable.

**Root cause:** Phase 2's scope didn't distinguish "the cell receptor schema" from "the CDS-specific instantiation of that schema." `COHERENCE-CELL-NORMAL-FORM.md` (#370, shipped 2026-05-17) names the kernel as substrate-independent, but the kernel doesn't have an executable schema home that's separate from CDS until Phase 2.5 ships.

**Disposition:** `patch-landed` (substance: Phase 2.5 inserted into #366 roadmap this session) + `next-MCA` (action: Phase 2.5 cycle to be filed as a child issue at dispatch).

**Patch:** #366 body now carries a Phase 2.5 section with two architectural options ((a) split into `schemas/cdd/` generic + `schemas/cds/` + `schemas/cdr/`; (b) freeze an adapter boundary with `evidence_refs` as an open typed map) and six draft ACs. The dispatch decision picks (a) or (b).

**Affects:** Phase 3 (cannot dispatch until Phase 2.5 lands); cnos#376 CDR receipt schema (depends on Phase 2.5 generic surface); a future cnos.cds bootstrap issue (depends on Phase 2.5 splitting CDS-specific fields into their own home).

**Cross-reference:** `ROLES.md §4a.3` (receipts enforce discipline) requires that V validates each protocol's receipt against that protocol's schema. Phase 2.5 makes this mechanical.

### F3: ε wave-launch surface undocumented (cdd-protocol-gap)

**Source:** ε self-observation while authoring this wave's manifest + dispatch prompts.

**Class:** `cdd-protocol-gap`.

**Trigger:** ε process-gap check during wave-launch.

**Description:** The user's directive "launch waves as epsilon launching delta to launch gammas" describes a real protocol chain: ε produces a wave plan that δ executes by dispatching γ for each cycle. cph used this pattern for the cdr-refactor wave 2026-05-18 (master cph#11; subs cph#12–15) — the wave-closeout, receipt, and ε iteration are all real artifacts. But no canonical cdd skill names the chain. `epsilon/SKILL.md` describes ε's per-finding discipline (ship MCA / file MCI / drop); it does not describe ε's wave-planning role. `operator/SKILL.md` §5.2 describes wave-mode dispatch from δ's side but not the upstream ε contribution. The result: ε wave-launches happen, but the protocol for them is invented per session.

**Root cause:** ε's role surface is finding-triggered; wave-launching wasn't surfaced as a separate ε responsibility when `epsilon/SKILL.md §1` was written. The cph wave demonstrated the pattern (ε iteration documents the wave's protocol findings; δ executes; γ runs cycles), but the pattern is implicit in the cph artifacts, not codified in cnos doctrine.

**Disposition:** `next-MCA`. The fix shape is "add an `epsilon/SKILL.md §1.5` or new section naming ε's wave-planning role, with the artifact set (wave manifest + dispatch prompts + launch-side iteration + close-side iteration) and the chain ε→δ→γ explicit." But this patch needs to land *after* this session's wave runs and produces empirical evidence of what's actually right (rather than committing the protocol on the basis of one session's improvisation).

**Issue filed:** none yet. Fold into Phase 6 (ε upscope) when that dispatches; or file a standalone cnos issue when the post-close ε iteration for this wave produces enough evidence to write the rule.

**First AC for the eventual MCA:** `epsilon/SKILL.md` names ε's wave-launch role; cites this wave as the empirical anchor; defines the artifact set (`.cdd/waves/{wave-id}/{manifest,status,dispatch-prompts}.md`); names the chain ε produces wave plan → δ dispatches γ per cycle → γ runs α/β.

### F4: Roadmap drift class — no canonical rule for "when doctrine ratifies, downstream roadmap text must reconcile" (cdd-protocol-gap)

**Source:** F1's root cause analysis.

**Class:** `cdd-protocol-gap`.

**Trigger:** ε process-gap check during F1 disposition.

**Description:** F1 surfaced a specific drift (`Contract × Receipt × EvidenceGraph` in #366 vs `Contract × Receipt` in #370). The deeper question: what catches this class of drift? When `COHERENCE-CELL-NORMAL-FORM.md` shipped at #370 on 2026-05-17, who was supposed to scan downstream roadmaps and reconcile? No one had it as a step. The drift sat for ~2 days until this session's ε review caught it. If Phase 3 had dispatched in those two days against the stale signature, the wasted work would have been the cost.

The β skill `review/diff-context/SKILL.md §2.1.8 "Authority-surface conflict"` catches authority conflicts *within a single cycle's diff*. It does not catch authority drift that lands cleanly in one cycle and stales doctrine in another cycle's issue body.

**Root cause:** No protocol rule for cross-cycle authority propagation. The closest existing surface is `gamma/SKILL.md` §"Canonical-skill staleness check before each γ phase change" — but that fires for γ within an active cycle, not for ε across cycles. The cph wave-iteration (2026-05-18) had a similar dynamic (rule 3.11b discoverability under wave-mode) — both are "the rule says X; the wave's reality is X'; until someone reviews across cycles, the drift persists."

**Disposition:** `next-MCA`. The fix shape is "ε runs a roadmap-coherence check during wave-launch: for each open roadmap issue (`#366`, `#376`, etc.), grep against the current ratified doctrine for signature/path/contract drift; surface any drift as a finding before the next phase dispatches." This is what this session did manually; codifying it as an ε responsibility makes it routine.

**Issue filed:** none yet. Could fold into Phase 6 (ε upscope) — the upscoped ε surface naturally absorbs this. Or could be its own small issue.

**First AC for the eventual MCA:** `epsilon/SKILL.md` names "roadmap-coherence check" as a wave-launch step; the check grep-greps each open roadmap issue against the current ratified doctrine paths; drift findings are surfaced before any phase that depends on the drifted surface dispatches.

## §3 Wave-launch decisions

Three γ cycles dispatched in this wave (`.cdd/waves/2026-05-19-protocol-hygiene/`):

| Cycle | Issue | Rationale |
|---|---|---|
| 1 | #375 | γ-axis preventive gate; small skill-patch; cleanly closes the rule-3.11b round-trip cost the cph wave (and cnos #369) paid. |
| 2 | #378 | Same axis as #375 but on the β-side discoverability under wave-mode; cph cdr-refactor wave (2026-05-18) is the empirical anchor with 4-of-4 sub uniform pattern. |
| 3 | #377 | Cross-repo coordination protocol; this session surfaced four more empirical anchors (cn-sigma agent-activate-skill; cph coherence-drift-sweep-followup with its own Rule 6 β-skill patch; bootstrap-cdr modified disposition; gait→cph empirical-anchor swap) that all feed into the canonical surface. |

Wave is parallel — no inter-cycle dependencies. δ may dispatch all three γ sessions concurrently or sequentially in any order. Wave manifest at `.cdd/waves/2026-05-19-protocol-hygiene/manifest.md`; dispatch prompts at `.cdd/waves/2026-05-19-protocol-hygiene/dispatch-prompts.md`.

Wave is **not** dispatching Phase 2.5 or Phase 3 from #366. Those wait on the architectural-choice decision (which fold into the same decision-class as cnos#376 AC7). Wave 1 (Phase 2.5 + Phase 3) launches in a subsequent session.

## §4 INDEX update

`cnos:.cdd/iterations/INDEX.md` (if it exists; per the convention named in `cdd/post-release/SKILL.md` §5.6b) should add one row:

```markdown
| Cycle | Issue | Date | Findings | Patches | MCAs | No-patch | Path |
|-------|-------|------|----------|---------|------|----------|------|
| wave-2026-05-19-protocol-hygiene (launch-side) | n/a (wave) | 2026-05-19 | 4 | 2 | 2 | 0 | .cdd/iterations/wave-2026-05-19-protocol-hygiene.md |
```

Patches landed (immediate): F1 (#366 body Phase 3 signature corrected), F2 (#366 body Phase 2.5 inserted).
MCAs filed (next): F3 (ε wave-launch surface — fold into Phase 6 or file standalone), F4 (roadmap drift class — fold into Phase 6 or file standalone).
No-patch decisions: none.

## §5 Open / carried debt

- **Phase 2.5 cycle not yet dispatched.** Roadmap insertion done; the actual schema split is the next dispatch. Architectural choice (a) vs (b) must be made at dispatch. Cross-references cnos#376 AC7 (architecture-choice question) — both decisions live in the same conceptual decision-class (how the generic kernel relates to domain overlays).
- **β-skill Rule 6 patch on `claude/file-cnos-cdr-issue-fi9Ld`.** Single-line semantic addition to `cdd/beta/SKILL.md` from the cph iteration ε-MCA on 2026-05-19. Awaits merge to cnos main. Not blocking this wave — Wave 2's β cycles operate without it if it hasn't merged.
- **Sigma task (3 deliverables) on the same branch.** Operator-side execution: cn-rho repo creation, cn-sigma discipline patch, cph#32 comment. Not blocking this wave.
- **F3 and F4 next-MCAs uncommitted.** ε's protocol-iteration responsibilities (wave-launch, roadmap-coherence check) need codification but the patch shape is best decided after this wave's empirical evidence lands. Carry as named debt; surface in the post-close ε iteration for this wave.

## §6 Next-MCA commitment

After this wave's three cycles close (375, 378, 377), ε returns to this file and:

1. Adds a §"close-side" section to this iteration document with per-cycle findings.
2. Disposes F3 (ε wave-launch surface) — either fold into a Phase 6 sub-cycle issue or file standalone with the wave-empirical-anchor argument.
3. Disposes F4 (roadmap drift class) — same disposition path as F3.
4. Updates `cnos:.cdd/iterations/INDEX.md` with the close-side row.

Until then, the launch-side iteration (this file as committed) is sufficient to dispatch the wave.

Filed by ε on 2026-05-19.
