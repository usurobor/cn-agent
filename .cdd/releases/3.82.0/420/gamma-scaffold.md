# γ scaffold — Cycle 420

**Cycle:** [cnos#420](https://github.com/usurobor/cnos/issues/420) — Sub 6 of [cnos#404](https://github.com/usurobor/cnos/issues/404)
**Branch:** `cycle/420` (from `origin/main` @ `f87e6e24`)
**Mode:** docs-only / cleanup (final cycle of the cnos#404 handoff extraction wave)
**Collapse pattern:** γ+α+β collapsed on δ. Commits: `α-420:`, `β-420:`, `γ-420:`.
**Date:** 2026-05-22

## Intent

Final cycle of the cnos#404 handoff extraction wave. Subs 1–5 have landed all v0.1 content into `src/packages/cnos.handoff/`; Sub 6 sweeps residual citations, tightens the package's v0.1-complete state across loader / HANDOFF.md / README.md / extraction-map.md, and emits the close-instruction so the merge auto-closes both #420 and #404.

## Wave progression (post-survey)

| Sub | Issue | Status | Surface | Outcome |
|---|---|---|---|---|
| Sub 1 | [cnos#415](https://github.com/usurobor/cnos/issues/415) | landed | package skeleton + extraction map | bootstrapped |
| Sub 2 | [cnos#416](https://github.com/usurobor/cnos/issues/416) | landed | cross-repo/SKILL.md wholesale move (644 lines) | canonical home flipped |
| Sub 3 | [cnos#417](https://github.com/usurobor/cnos/issues/417) | landed | dispatch/SKILL.md synthesis (γ§2.5 + operator§3a + delta§2) | three sources → one home |
| Sub 4 | [cnos#418](https://github.com/usurobor/cnos/issues/418) | landed | mid-flight/SKILL.md + artifact-channel/SKILL.md | two parallel surfaces |
| Sub 5 | [cnos#419](https://github.com/usurobor/cnos/issues/419) | landed | receipt-stream/SKILL.md (post-release §5.6b extraction) | aggregator doctrine moved |
| **Sub 6** | **[cnos#420](https://github.com/usurobor/cnos/issues/420)** | **this cycle** | **cross-reference cleanup + close tracker** | **wave closes** |

## Surface plan (D1–D7)

### D1 — Residual citation sweep
Audit complete. Findings:
- **cnos.cdd / cnos.cds / cnos.cdr / cnos.handoff:** all canonical-home statements now correctly name `cnos.handoff/skills/handoff/<area>/SKILL.md`. Old-path cites that remain are uniformly framed as "consumer-side", "role-side realization", "pointer", or "runbook (pointer to canonical home)" — i.e. the cdd-side surfaces correctly self-describe as consumers of the new canonical doctrine. AC1 passes structurally as authored.
- **One stale citation discovered:** `.cdd/iterations/INDEX.md` header line 5 reads `See \`cdd/post-release/SKILL.md\` Step 5.6b for the per-finding shape and the aggregator update procedure.` — this should re-point to the canonical wire-format home at `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` (with the cdd-side runbook as the secondary pointer).
- **Historical PRAs in `docs/gamma/cdd/docs/2026-05-*/...POST-RELEASE-ASSESSMENT.md`** name old-path surfaces in their archival prose (e.g. "post-release/SKILL.md §5.6b re-attributes cdd-iteration.md as ε's work product"). These are **immutable historical records** of past cycles and **must not be edited** — they describe what was true at the cycle's pin, not what is currently canonical.

### D2 — HANDOFF.md final pass
Current state already has 5 Landed entries, 0 Forthcoming. Add: "v0.1 complete — cnos#404 wave closed 2026-05-22" version-note tightening; minor wording edits to remove "Sub 6 (cross-reference cleanup + close tracker) remains" residue from the Status line.

### D3 — handoff/SKILL.md loader pass
Drop the "v0.1 status … forthcoming" paragraph in §"Load order"; tighten the §"Sub-skills" intro that says "Forthcoming per Subs 2–5"; rewrite the §"v0.1 caveat" section to a v0.1-complete statement (the five sub-skills are landed; the loader names them as canonical, not as advisory). The "Pending Sub 2 or Sub 3" qualifiers in Steps 1 and the §"Conflict rule" line about "(once `HANDOFF.md` lands per Sub 2 or Sub 3 …)" must come out — HANDOFF.md has landed.

### D4 — README.md v0.1-complete update
Rewrite §"Package Structure" to drop the "**v0.1 skeleton**" + "Pending Subs 2–5" framing. Replace §"Forthcoming surfaces" with §"Surfaces" listing the five Landed sub-skills with their landing issue numbers. Rewrite §"Status" to declare v0.1 complete; cite the #404 wave closure. Repair the broken Quick Start link to `cnos.cdd/skills/cdd/cross-repo/SKILL.md` ("the 644-line primary extraction candidate") — that file is now a 28-line compat pointer, not the canonical surface. Repair the "(until Sub 2 of cnos#404 moves it here)" parenthetical on the cds line.

### D5 — extraction-map.md final-status
Add a top-of-file v0.1-complete note: "**Wave complete (2026-05-22)**: all surfaces migrated; cnos#404 closed via cnos#420." Update §9 ("Cross-reference cleanup + close tracker — Sub 6") with a "Status: v0.1 complete" line and reference cnos#420. Update §11 ("Open questions") to note Q1–Q7 resolved or carried.

### D6 — Compatibility-stub decision
The 28-line stub at `cnos.cdd/skills/cdd/cross-repo/SKILL.md` (carrying its own line-29 self-note "A future cleanup cycle may remove this file entirely once all consumers cite cnos.handoff directly") was authored by Sub 2 (#416). Decision for Sub 6: **preserve as-is**. Rationale: cnos.cdr v0.1 documentation (per the cnos.cdr README.md L17 and cnos.cdr CDR.md L497) still mentions cross-repo via parent-package reference patterns; cnos.cds extraction-map.md §111 and L283 explicitly cite `cnos.cdd/skills/cdd/cross-repo/SKILL.md` (as historical pre-migration prose). External repos (cph, sigma, future c-d-X projects) may also cite the old path; the 28-line stub costs nothing and preserves resolvability. The stub's own "may be removed in a future cleanup cycle" wording is the correct posture; no edit needed. Document the decision in gamma-closeout.md.

### D7 — Final extraction-map sweep
Verify every numbered row §1–§9 has a "Status" field. §1–§6 already carry "Status: v0.1 migrated" lines from prior subs; §7 has "Destination commitment: Open" (deferred-by-design); §8 has "Status: v0.1 migrated"; §9 receives a new "Status: v0.1 complete" line via D5.

### INDEX.md header repair (out-of-scope-of-D1-but-discovered-by-D1)
Re-point line 5 of `.cdd/iterations/INDEX.md` from `cdd/post-release/SKILL.md` to the canonical wire-format home at `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` (with the cdd-side runbook pointer as secondary).

## Implementation contract

| Axis | Pinned value |
|---|---|
| Language | Markdown |
| CLI integration target | None |
| Package scoping | Primary: `src/packages/cnos.handoff/{README.md, skills/handoff/SKILL.md, skills/handoff/HANDOFF.md, docs/extraction-map.md}`. Secondary: `.cdd/iterations/INDEX.md` header repair. Compatibility stub at `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` preserved untouched. |
| Existing-binary disposition | Compatibility stub preserved (D6 decision documented in closeouts) |
| Runtime dependencies | None |
| JSON/wire contract | N/A |
| Backward compat | Compatibility stub preserves any unfound external citations; old-path cites in historical PRAs untouched (archival records). |

## Acceptance criteria

AC1–AC11 per [cnos#420](https://github.com/usurobor/cnos/issues/420). All mechanical greps or file-existence/diff checks. Verified in `self-coherence.md` post-α.

## Hard rules

1. No content moves. Subs 2–5 migrated everything.
2. CCNF kernel in `cnos.cdd/skills/cdd/CDD.md` unchanged (AC9).
3. No new schemas / runtime / cdd-verify changes (AC10).
4. Compatibility stub at `cnos.cdd/skills/cdd/cross-repo/SKILL.md` preserved (per D6 decision).
5. Historical PRAs under `docs/gamma/cdd/docs/2026-05-*/` not edited (archival records).
6. cnos.cdr not modified (no residual cites discovered).

## Non-goals

- Do NOT move any new content.
- Do NOT redesign anything.
- Do NOT delete the compatibility stub.
- Do NOT modify CCNF kernel content in CDD.md.
- Do NOT preempt cnos#405 Track A / Track B work.
- Do NOT edit historical PRAs.

## Operator action on close

Merge `cycle/420` to `main` with `--no-ff` and message including BOTH `Closes #420` AND `Closes #404` so GitHub auto-closes the parent tracker.
