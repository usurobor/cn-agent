# α self-coherence — cycle/397

**Issue:** cnos#397 — Phase 4a of #366
**Cycle branch:** cycle/397
**Identity:** delta@cdd.cnos (γ+α+β-collapsed-on-δ; β-α-collapse acknowledged)
**Mode:** design-and-build

## §1 AC oracle results (mechanical)

### AC1 — `delta/SKILL.md` exists with standard skill frontmatter

`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` exists. Frontmatter contains:
- `name: delta` ✓
- `description: …` (one-line δ-role-skill description) ✓
- `artifact_class: skill` ✓
- `governing_question: …` (names the two-sided membrane question) ✓
- `parent: cdd` ✓
- `triggers: [boundary, membrane, override, receipt, verdict, gate, dispatch-enrichment, implementation-contract]` ✓
- `scope: role-local` ✓

Additional frontmatter (matching `operator/SKILL.md` template): `kata_surface: embedded`, `visibility: internal`, `inputs`, `outputs`, `requires`, `calls`. All present. **AC1 PASS.**

### AC2 — Two-sided membrane explicit

`delta/SKILL.md` has:
- `## 1. Outward membrane — receipt + verdict → boundary decision` (lines 65-138)
- `## 2. Inward membrane — γ contract → α-ready dispatch (implementation-contract enrichment)` (lines 141-183)

Each side cites the other surfaces:
- §1 cites: `RECEIPT-VALIDATION.md` (V verdict + boundary decision composition; ordering rule); `operator/SKILL.md` (mechanics location); `COHERENCE-CELL-NORMAL-FORM.md` (cell outcomes).
- §2 cites: `gamma/SKILL.md` §2.5 Step 3b (γ template); `alpha/SKILL.md` §3.6 (α constraint); `beta/SKILL.md` Rule 7 (β verification); cnos#393 (precursor); cnos#389/#391/#392 (empirical anchors).

The doctrine box at the top of `delta/SKILL.md` makes the two-sided framing explicit: outward = receipt + V verdict → BoundaryDecision; inward = γ contract → α-ready dispatch. **AC2 PASS.**

### AC3 — Operator §3a content fully moved (no duplication)

`operator/SKILL.md` §3a is replaced by a one-paragraph redirect pointing at `delta/SKILL.md` §2. The substance (7 axes, fill-or-escalate paths, four-surface mesh, empirical anchors) lives only in `delta/SKILL.md`.

`rg "δ-inward-membrane|inward membrane" src/packages/cnos.cdd/skills/cdd/` returns hits in:
- `delta/SKILL.md` — substance (the doctrine)
- `operator/SKILL.md` — header-only (line 252, the section title kept for in-repo navigation) + the redirect prose (line 254) naming the doctrine to direct readers to delta/

This conforms to the AC3 spec ("`operator/SKILL.md` §3a is replaced with a one-line cross-reference + redirect"): the section title is retained for backwards-compatible navigation (readers searching the operator file for "δ as inward membrane" find the redirect immediately). **AC3 PASS.**

### AC4 — Override semantics + verdict-vs-decision distinction preserved

`delta/SKILL.md` §3 "Override — degraded boundary action":

- Line 189: "Override is a degraded boundary action; it is never a form of validity, and it never rewrites V's verdict."
- §3.1 "Override does not rewrite the ValidationVerdict" — three explicit failure modes (never rewrites; never substitutes for PASS; never emits OVERRIDE-PASS).
- §3.2 "Downstream-consumer detection rule (binding biconditional, per cnos#367)" — the biconditional `receipt is degraded ⇔ boundary.override != null`.
- §3.5 "ValidationVerdict vs BoundaryDecision — the structural distinction" — table + three explicit constraints.

Cites `RECEIPT-VALIDATION.md` §Q4 and §"ValidationVerdict vs BoundaryDecision" as the cnos#367 freeze source. **AC4 PASS.**

### AC5 — Cross-references updated

References that cite `operator/SKILL.md` for δ-role content (specifically §3a and the inward-membrane doctrine) updated to `delta/SKILL.md`:

| File | Line (before) | Update |
|---|---|---|
| `gamma/SKILL.md` | 363 | `operator/SKILL.md` §3a → `delta/SKILL.md` §2 |
| `gamma/SKILL.md` | 370 | `operator/SKILL.md` §3a → `delta/SKILL.md` §2 |
| `alpha/SKILL.md` | 355 | `operator/SKILL.md` §3a → `delta/SKILL.md` §2 |
| `beta/SKILL.md` | 175 | `operator/SKILL.md` §3a → `delta/SKILL.md` §2 |

Total: 4 cross-references updated. Oracle: `rg "operator/SKILL.md.*§3a|operator/SKILL.md.*inward membrane" src/packages/cnos.cdd/skills/cdd/` returns hits only in `delta/SKILL.md` (the §2.2 "Phase 4a landing note" historical reference to where the doctrine *was* anchored before this cycle). No external skill cites `operator/SKILL.md` for δ-role content anymore. **AC5 PASS.**

**References intentionally left at `operator/SKILL.md`** (these are operator-as-coordinator + harness + release-effector citations, NOT δ-role; Phase 4b/4c will move when those phases ship):
- `post-release/SKILL.md` line 39 (`operator/SKILL.md §3.4` — release mechanics)
- `release/SKILL.md` lines 35, 238 (`§3.4` — release mechanics)
- `activation/SKILL.md` multiple refs (identity, dispatch configs, §3.4 mechanics, §5.2/§5.3)
- `review/SKILL.md` line 141 (§5.2 wave-mode = dispatch config)
- `CDD.md` multiple refs (identity, §3.3, §3.4, timeout recovery)
- `COHERENCE-CELL.md` doctrinal predictions
- `RECEIPT-VALIDATION.md` and `COHERENCE-CELL-NORMAL-FORM.md` Phase 4 target references (historical)

### AC6 — `operator/SKILL.md` retains operator-as-coordinator content

`operator/SKILL.md` after edits still contains:
- Frontmatter, Core Principle ✓
- Algorithm §1 (dispatch loop), §2 (Wait), §5 (Dispatch configurations), §6 (What operator does NOT do), §7 (Cycle lifecycle), §8 (Timeout recovery), §9 (Embedded Kata), §10 (Wave Coordination) ✓
- §3 retained as "Gate" with mechanics + redirect to delta/ for role-policy ✓
- §3.4 release-cut runbook (`scripts/release.sh` invocation, manual-tag prohibition, tag-message generation) ✓
- γ=δ collapse content in §5.2 ✓
- Dispatch coordination prose throughout ✓

Cross-reference to `delta/SKILL.md` added in: Core-Principle area (implicit via the "see delta/SKILL.md for the role-skill surface" framing in delta itself), §3 (top-of-section redirect note), §3a (replaced by redirect), §4 (replaced by redirect), algorithm step 7. **AC6 PASS.**

### AC7 — No effector mechanics in this cycle

`delta/SKILL.md`:
- `rg "scripts/release.sh"` returns hits at lines 94 (authority-naming: "the `scripts/release.sh` runbook lives in `operator/SKILL.md` §3.4 (Phase 4c relocates)") and 338 (authority-naming: "Phase 4c — release-effector mechanics (pending)"). Both point at mechanics locations; neither contains mechanics.
- `rg "gh run list"` returns ZERO hits.
- `rg "tag|release|deploy"` returns 32 hits, all in authority-naming context: BoundaryDecision enumerants (`release` as a δ decision), role-policy paragraphs (δ as sole tag-author, blocks until CI green), table entries (Tag push + release | δ-as-role authority), override examples ("Override: release CI red on…"), Phase 4b/4c cross-references.

No runbook prose. The `scripts/release.sh` step list, CI polling commands, branch cleanup runbook stay in `operator/SKILL.md` §3.4. **AC7 PASS.**

## §2 Implementation-contract conformance (per cnos#393 Rule 7)

| Axis | Pinned | Conformance |
|---|---|---|
| Language | Markdown | All edits are .md only. ✓ |
| CLI integration target | N/A | No CLI added or edited. ✓ |
| Package scoping | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` (new) + edits to `cdd/operator/SKILL.md` | New file at exact path; operator/SKILL.md at exact path; gamma/, alpha/, beta/ edits within `cdd/` scope. ✓ |
| Existing-binary disposition | operator/SKILL.md keeps operator-as-coordinator + dispatch role; boundary-policy + δ-inward-membrane sections move to delta/SKILL.md; harness mechanics → Phase 4b; release effector → Phase 4c. | Implemented exactly as pinned. ✓ |
| Runtime dependencies | None | No deps added. ✓ |
| JSON/wire contract | N/A | No JSON edited. ✓ |
| Backward compat | All existing operator/SKILL.md refs from other skills updated to delta/ where appropriate; operator/SKILL.md continues to exist for the operator-as-coordinator surface. | All 4 δ-role-content refs (γ×2, α×1, β×1) updated; operator/SKILL.md retained with §1, §2, §3 (mechanics), §5–§10. ✓ |

## §3 Content-loss audit (β rigor per cnos#393 Rule 7)

Verifying no δ-content drops from operator/SKILL.md without landing in delta/SKILL.md:

| operator/SKILL.md content removed | Lands in delta/SKILL.md |
|---|---|
| §3.1 external-actions table + framing as δ-held | §1.1 (table + framing as "δ-as-role holds authority") |
| §3.2 execute-on-request-not-observation prose + examples | §1.2 (verbatim semantics, with cross-link to delta-§3 for override) |
| §3.3 report-completion prose + examples | §1.3 (verbatim) |
| §3.4 role-policy paragraphs ("triad's work not complete until tagged"; "δ blocks release until CI green; owns recovery on red") | §1.1 + §3.3 (override escape hatch) |
| §3.5 the-tag-is-the-signal prose | §1.4 (verbatim) |
| §3a entire (two-sided framing, 7 axes, mesh, empirical anchors) | §2 (entire — verbatim with Phase 4a landing-note update) |
| §4.1 when-to-override criteria | §3.3 (criteria preserved + extended with the CI-red escape hatch) |
| §4.2 override-protocol three-step ("what / why / new state") + bad/good examples | §3.4 (three-step + bad/good examples preserved + override-block field requirements added per cnos#367) |
| §4.3 not-for-taste discipline | §3.3 (preserved verbatim semantics) |

**Net additions in delta/SKILL.md** (not in operator/SKILL.md before, lifted from `RECEIPT-VALIDATION.md` per AC4 binding):
- §1.5 "BoundaryDecision: the five outcomes δ records" — the five enumerants per CCNF.
- §3.1 "Override does not rewrite the ValidationVerdict" — three explicit failure modes per cnos#367 §Q4.
- §3.2 "Downstream-consumer detection rule" — the biconditional per cnos#367 §Q4.
- §3.5 "ValidationVerdict vs BoundaryDecision" table per cnos#367 §"ValidationVerdict vs BoundaryDecision".
- §4 "Composition with V" — composition rule and outcomes per `RECEIPT-VALIDATION.md` §Validation Interface.
- §5 "What δ-as-role does NOT do" — role-boundary constraints per CCNF.
- §6 "Cross-references and relationships" — discoverability surface.
- §7 "Phase 4 of cnos#366 — what this cycle ships and what remains" — phase tracking.

All net additions are AC-mandated synthesis from existing cited surfaces (cnos#367 freeze + CCNF + operator §3a). No invented content; only relocation + AC-binding consolidation.

## §4 Refusal-condition check (post-build)

- ✅ No ambiguous-surface sections: §3 of operator/SKILL.md was successfully split into role-policy (moved to delta/) and mechanics (retained in operator/); the split was unambiguous because §3.4 paragraphs clearly separated role-doctrine ("triad's work not complete until tagged") from runbook ("`scripts/release.sh` invocation, CI polling, branch cleanup").
- ✅ No parallel cycle conflict: branch list shows cycle/397 fresh; cycle/398/399 unstarted (per branch list).
- ✅ No cross-reference required content change: all 4 updates were path-only swaps.

## §5 Review-readiness

α (collapsed onto δ-as-agent) signals review-readiness. All ACs PASS. Implementation-contract conformance verified. No content drops. β-collapsed-on-δ review follows.
