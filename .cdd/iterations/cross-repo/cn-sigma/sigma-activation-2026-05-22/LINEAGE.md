# Cross-Repo Lineage: cnos → cn-sigma (sigma-activation-2026-05-22)

This bundle is **case (d.2) — operator-pending; target-repo-exists-but-unreachable** per `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md §2.1`. The cnos session that authored it has MCP scope restricted to `usurobor/cnos` and cannot write to `usurobor/cn-sigma`; the bundle stages the patches in cnos, the operator applies them to cn-sigma. The lifecycle is `drafted (operator-pending) → operator-applied → archive` per `cross-repo/SKILL.md §2.8` (no STATUS state machine applies; the disposition is recorded in this LINEAGE).

## Purpose

Stage the **Sigma activation** delivery for `usurobor/cn-sigma`: feedback patches against `spec/PERSONA.md` and `spec/OPERATOR.md` that codify the **persona-level learnings** from the cnos#366 → #403 arc — those that belong to Sigma's engineering-persona doctrine regardless of which protocol (cdd / cds / cdr-engineering-side) or project (cnos / cph / future) Sigma operates in — plus a companion empirical-anchor doc mapping the bootstrap arc to the 10 named rules with cycle citations.

The patches PROPOSE additions to cn-sigma's `spec/PERSONA.md` and `spec/OPERATOR.md`; the cn-sigma operator (per the existing `Continuity rule` in cn-sigma's PERSONA.md) decides on application. The bundle structure IS the approval mechanism: this is not a unilateral edit, it is a draft for review.

This sub completes the phase per operator directive ("we'll activate sigma once we're done with the phase"; #403 closed at `378a54f0` on 2026-05-22).

## Source

- **Trigger:** Operator directive after #403 (cnos.cds v0.1 bootstrap wave) closed at `378a54f0`. The cnos#366 → #403 arc surfaced persona-level learnings that have outgrown the protocol-package surface — they belong to Sigma (the engineering persona) as standing doctrine.
- **Reference doctrine:**
  - `ROLES.md §4a` (five-layer enforcement chain: persona / operator contract / protocol overlay / project binding / receipt / validator)
  - `ROLES.md §4a.1` (discipline-profile required persona section: six fields)
  - `ROLES.md §4a.2` (loss-function distinction: engineering = artifact-improvement-under-repairable-feedback; research = truth-preservation-under-uncertainty)
  - `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` (CCNF kernel + typed trust surfaces + recursive realization; design substrate for the arc)
  - cnos#393 issue body (δ-as-two-sided-membrane; outward + inward; implementation-contract enrichment at dispatch)
- **Sibling exemplars:**
  - `cnos:.cdd/iterations/cross-repo/cn-sigma/discipline-section-2026-05-19/` (case d.2 precedent; the patch shape this bundle mirrors; appends `## Discipline` to cn-sigma/spec/PERSONA.md via `@@ EOF @@` hunk)
  - `cnos:.cdd/iterations/cross-repo/cn-sigma/agent-activate-skill/` (case a.1 precedent; LINEAGE.md shape, drafted→accepted direct-acceptance variant)
- **Empirical arc anchored:** cnos#366, #376, #391, #392, #393, #402, #403, #404, #405 (10 named rules cross-referenced in `docs/empirical-anchor-cnos-bootstrap-arc.md`).

## Bundle contents

| File | Purpose |
|---|---|
| `LINEAGE.md` | this file; case (d.2) schema per `cross-repo/SKILL.md §2.6` case (d); first paragraph names the case |
| `PERSONA-additions.patch` | unified diff against `cn-sigma:spec/PERSONA.md` adding `## Engineering-persona protocol commitments` (6 numbered rules) after the existing `## Discipline` section; uses `@@ EOF @@` append-trick per discipline-section-2026-05-19 precedent; cn-sigma operator may relocate per Continuity rule |
| `OPERATOR-additions.patch` | unified diff against `cn-sigma:spec/OPERATOR.md` adding `## CDD wave-execution pattern (engineering-persona operations)` (3 sub-rules) after the existing `## CDD role assignment` section; uses `@@ EOF @@` append-trick (operator may relocate to immediately after `## CDD role assignment` per Continuity rule) |
| `PERSONA-discipline-receipt-additions.patch` | unified diff against `cn-sigma:spec/PERSONA.md` augmenting the existing `Receipt requirements` field under `## Discipline` (shipped by discipline-section-2026-05-19) with anti-gaming guardrails (3 named attacks + TSC-observation rule); appends a new sub-block `### Receipt requirements — anti-gaming guardrails` to the end of PERSONA.md via `@@ EOF @@` (operator may relocate inline under `Receipt requirements`) |
| `docs/empirical-anchor-cnos-bootstrap-arc.md` | ~300–400-line companion doc; one section per named rule (6 persona + 3 operator + 1 discipline-augmentation = 10 total) with cycles cited and one-line lessons; final section names what's next for Sigma (#404 handoff extraction; #405 CCNF-O + TSC steering; v1 deep-role rewrites deferred) |

`FEEDBACK.patch` is **not** included. Per case (d.2) (`cross-repo/SKILL.md §2.5` case-d row), no STATUS state machine applies and no FEEDBACK.patch is required unless the cn-sigma operator explicitly wants a STATUS-event emission in any cnos-side STATUS ledger. The current operator posture (Continuity rule + drafted/applied/archive lifecycle) does not require it.

## Application

The cnos session cannot write to `usurobor/cn-sigma` (MCP scope restricted to `usurobor/cnos`). The operator applies the three patches from cn-sigma's repo root.

```
cd usurobor/cn-sigma

# Apply in D2 → D4 → D3 order (D2 and D4 both anchor on `## Continuity rule`
# in PERSONA.md; D2 first ensures correct stacking. D3 is independent of
# PERSONA-side patches; order vs D3 does not matter).
#
# All three patches use --recount because hunk-header line numbers are
# placeholders (the cnos session does not know cn-sigma's exact line
# counts); --recount regenerates line numbers from context-match against
# the trailing-context anchor lines.

# D2 — adds engineering-persona protocol commitments to PERSONA.md
git apply --recount <path-to-this-bundle>/PERSONA-additions.patch

# D4 — augments Discipline / Receipt requirements with anti-gaming guardrails in PERSONA.md
git apply --recount <path-to-this-bundle>/PERSONA-discipline-receipt-additions.patch

# D3 — adds CDD wave-execution pattern to OPERATOR.md
git apply --recount <path-to-this-bundle>/OPERATOR-additions.patch

git add spec/PERSONA.md spec/OPERATOR.md
git commit -m "persona+operator: Sigma activation (cnos#413; 6 persona + 3 operator + 1 discipline-augmentation)"
git push
```

The patches use placeholder hunk-header line numbers (`@@ -999,1 +999,N @@`) because the cnos session that authored the bundle does not have access to cn-sigma's exact line counts. `git apply --recount` regenerates the hunk-header line numbers from the trailing-context anchor line (`## Continuity rule` for D2 + D4; `## Durable preferences only` for D3), which the operator's pre-flight reads confirm are stable terminal section headings in their respective files. Without `--recount`, vanilla `git apply` rejects the patches as "corrupt" because the placeholder line numbers do not match the target's line count.

If `## Continuity rule` or `## Durable preferences only` content lines have changed since the pre-flight read (2026-05-22), `patch -p1 --fuzz=3 < <patch-name>` is the recommended fallback for context drift.

Each patch's free-prose header notes that the cn-sigma operator may relocate the new block to a more natural insertion point per the Continuity rule (e.g. immediately after the section the header text declares as the logical neighbor — `## Discipline` for D2 + D4; `## CDD role assignment` for D3). Content placement is the operator's decision; content existence is what the patches contribute.

If the operator chooses to apply only a subset (e.g. D2 + D4 but defer D3 pending further review), the bundle remains `drafted (operator-pending)` for the unapplied patches and the LINEAGE is updated post-application to record the partial-disposition.

## Verification after application

Run from `usurobor/cn-sigma` repo root:

```
# D2 — six persona commitments landed
rg "^## Engineering-persona protocol commitments" spec/PERSONA.md
# expect: 1 match

rg "^[0-9]+\. \*\*" spec/PERSONA.md | grep -E "two-sided membrane|γ-clarification|five-layer|loss function|β-α-collapse|Cross-protocol routing"
# expect: 6 matches (one per rule)

# D3 — three operator sub-rules landed
rg "^## CDD wave-execution pattern" spec/OPERATOR.md
# expect: 1 match

rg "Wave dispatch shape|B-lite extraction rule|Implementation contract pinned at dispatch" spec/OPERATOR.md
# expect: 3 matches

# D4 — anti-gaming guardrails landed in/near Discipline section
rg "simplify-away-truth|avoid-hard-refactors|tiny-only-shipments|TSC measurement is observation" spec/PERSONA.md
# expect: 4 matches

# Cross-references resolve
rg "ROLES.md §4a|cnos#393|cnos#405" spec/PERSONA.md spec/OPERATOR.md
# expect: ≥3 hits across the two files

# Empirical-anchor doc cited if operator chose to land it in cn-sigma/docs/
# (the doc is included in this bundle; whether cn-sigma adopts it under docs/ is operator choice)
test -f docs/empirical-anchor-cnos-bootstrap-arc.md && echo "anchor doc landed in cn-sigma"
# expect: either present (if operator copied) or absent (if anchor stays in cnos as a permanent reference)
```

The anchor doc (`docs/empirical-anchor-cnos-bootstrap-arc.md`) lives inside this bundle. The cn-sigma operator may copy it into `cn-sigma:docs/` if they want a permanent record on cn-sigma's side; alternatively, this bundle remains the permanent anchor on cnos (case d.2 bundles preserved indefinitely as audit artifacts per `cross-repo/SKILL.md §2.8.2`-analog reading — though case (d) is not formally covered by §2.8.2, the audit value is the same).

## What is NOT in scope here

This bundle is the **operator-pending case (d.2) feedback channel**. It does not:

- Modify cn-sigma directly. MCP scope restriction; the operator effects the apply.
- Modify any cnos package (`cnos.cdd`, `cnos.cdr`, `cnos.cds`, `cnos.core`, etc.). Hard rule per cnos#413; `git diff origin/main..HEAD -- src/packages/` returns 0 lines.
- Rewrite existing cn-sigma sections (`## Identity`, `## Voice`, `## Memory discipline`, `## Conduct (Sigma-specific)`, `## Discipline`, `## Continuity rule` per operator pre-flight reads). Only ADDS new content; the patches respect cn-sigma's `Continuity rule` requirement for intentional governance of persona evolution.
- File a STATUS event in any cnos-side ledger. Case (d.2) carries no STATUS state machine.
- Bootstrap any new cnos package, propose a roadmap, or define a new validator. The 10 rules codify existing observations into persona standing doctrine; they do not introduce new functionality.
- Auto-apply at the cnos side. The cnos session cannot write to cn-sigma; the operator gate is the application step.
- Touch cn-rho, cph, or any other counterpart repo. Sigma-scoped only.
- Implement CCNF-O or TSC integration. Those are #405; this bundle cites the §4.7 anti-gaming guardrails as authority for D4 but does not implement the steering layer.
- Implement the cnos.handoff / coordination extraction. That is #404; this bundle is unaffected by its disposition.

## Disposition

**Status: drafted (operator-pending).**

Application gate: the cn-sigma operator running `git apply` on each patch (or a documented subset, per Continuity rule). Until then, the bundle reports as `open` in any cnos-side bundle-state inventory per `cross-repo/SKILL.md §2.4` case (d) row (`drafted (operator-pending)` is a phase-name synonym for `open`).

Once the operator confirms application to cn-sigma, this bundle moves to `archived/` per `cross-repo/SKILL.md §2.8.1` case (d) archival rule. The cnos-side mirror's audit value warrants preservation indefinitely (analogous to §2.8.2 target-side mirror preservation); whether cnos archives this bundle or retains it as a permanent anchor is a follow-on decision documented at archival time.

Filed by γ@cnos on 2026-05-22 (cycle/413; γ+α+β-collapsed-on-δ per `cross-repo/SKILL.md §2.9`-analog hat-collapse — though case (d) bundles are not formally bilateral, the operator-as-δ pinned the implementation contract at dispatch; the actor playing γ collapsed α and β onto δ for this docs-only authoring cycle; α=β remains prohibited for substantive code work).
