# γ scaffold — cycle/393

**Issue:** cnos#393 — δ-as-architect: implementation-contract enrichment at dispatch (α/β/γ/δ skill patches; Phase 4 input)
**Mode:** design-and-build, γ+α+β-collapsed-on-δ (per breadth-2026-05-12 precedent and the operator-pinned dispatch)
**Branch:** `cycle/393` from `origin/main` (base SHA: `e531dba0`)
**Priority:** P1
**Parent:** cnos#366 (Phase 4 design input)

## Surfaces γ expects α to touch (4-patch surface)

1. **Patch A** — `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`
   - Add rule "Implementation contract is δ's, not α's" under `## 3. Rules` (current rules 3.1–3.5; this becomes 3.6).
   - α MUST NOT change 7 contract axes pinned by δ at dispatch; α MUST surface to γ/δ before coding if any axis is unpinned.
   - Cross-references Patch B (β verifies), Patch C (γ template names the axes), Patch D (δ inward enrichment).
   - Empirical anchors: cnos#389 (Python-not-Go), cnos#391 (wrong package scoping + separate binary).

2. **Patch B** — `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md`
   - Add `### 7. Implementation-contract coherence` under `## Role Rules` (current rules 1–6; this becomes Rule 7).
   - β verifies the 7 axes conform to the contract pinned at dispatch; behavior-only ACs are necessary but not sufficient.
   - Non-conformance → REQUEST CHANGES, severity D, classification `implementation-contract`.
   - Cross-references Patch A, Patch C, Patch D.
   - Empirical anchors: cnos#389 R1 + cnos#391 R1 behavior-only APPROVE failures.

3. **Patch C** — `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md`
   - Extend §2.5 Step 3b dispatch-prompt template with required `## Implementation contract` section enumerating 7 axes:
     - Language
     - CLI integration target
     - Package scoping
     - Existing-binary disposition
     - Runtime dependencies
     - JSON/wire contract preservation
     - Backward-compat invariant
   - Rule: γ MUST NOT dispatch with empty / "TBD" rows; γ asks δ if a value is undecidable.
   - Cross-references Patch D (δ owns inward enrichment of the template), Patch A (α MUST NOT improvise), Patch B (β verifies).
   - Empirical anchors: cnos#389 + cnos#391.

4. **Patch D** — `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md`
   - Add a new top-level section "δ as inward membrane: implementation-contract enrichment at dispatch."
   - δ is a two-sided membrane: **outward** (boundary decision on receipts; existing scope) + **inward** (implementation-contract enrichment of γ's dispatch contracts going to α).
   - At dispatch time, δ reviews γ's `## Implementation contract` section, fills unpopulated rows, blocks dispatch if undecidable.
   - This is the design-prerequisite for Phase 4 (δ split); Phase 4 will eventually move this content to `delta/SKILL.md`.
   - Cross-references Patch C (γ template), Patch A (α constraint), Patch B (β verification).
   - Empirical anchors: cnos#389 + cnos#391.

## AC oracle per AC1–AC7

| AC | Oracle | Surface |
|---|---|---|
| AC1 | `rg "α MUST NOT change the implementation-contract\|Implementation contract is δ" src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` ≥1 hit; anchors cite #389 + #391 | `alpha/SKILL.md` |
| AC2 | `rg "Implementation-contract coherence\|Rule 7" src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` ≥1 hit; anchors cite #389 R1 + #391 R1 behavior-only APPROVE failures | `beta/SKILL.md` |
| AC3 | `rg "## Implementation contract" src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` ≥1 hit; section enumerates all 7 axes; rule "γ MUST NOT dispatch with empty rows" named | `gamma/SKILL.md` §2.5 |
| AC4 | `rg "inward membrane\|implementation-contract enrichment" src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` ≥1 hit; names δ as two-sided; references Patch C template | `operator/SKILL.md` |
| AC5 | cnos#366 Phase 4 body explicitly names δ as two-sided membrane (inward + outward); absorbs cnos#393 as Phase 4 input | cnos#366 body |
| AC6 | Cross-reference mesh: α cites γ template + β + δ; β cites α + γ template + δ; γ cites α + β + δ; δ cites γ template + α + β. Bidirectional via grep. | All 4 patches |
| AC7 | Each patch cites cnos#389 + cnos#391 by issue number as empirical anchors. | All 4 patches |

## Cross-skill referential mesh plan (AC6)

Bidirectional mesh — each patch cites the other three by name + path:

```
            ┌─────────────────────────────┐
            │      Patch C (γ)            │
            │   "## Implementation        │
            │    contract" template       │
            └────────┬───────────┬────────┘
                     │           │
       cites template│           │cites template
                     ▼           ▼
┌──────────────────────┐   ┌──────────────────────┐
│  Patch A (α)         │◄──┤  Patch D (δ-inward)  │
│  "α MUST NOT         │   │  "implementation-    │
│  improvise"          │   │  contract enrichment"│
└─────┬──────┬─────────┘   └─────┬───────────────┘
      │      │                    │
      │      └────── cites α ─────┤
      │                            │
      └── cites β ──► ┌────────────┴────┐
                      │  Patch B (β)    │
                      │  "Rule 7 IC     │
                      │   coherence"    │
                      └─────────────────┘
```

Mesh is **unidirectional pairwise** to avoid hard circular reasoning: each patch cites the *artifact* of the other three (the rule, the template section, the inward-membrane section), not the patch's own justification surface. Citations are "see X for the role-side enforcement" rather than "as proved in X." No patch is logically dependent on another for its rule to make sense; each is locally self-justifying via the empirical anchors (#389, #391). The mesh is for **discoverability**: a future α/β/γ session loading any one finds the others.

## Empirical anchor (per §2.2a)

This cycle's framing rests on cnos#389 + cnos#391 + cnos#392 cdd-iteration:

- **cnos#389** (Python-not-Go): α implemented V in Python despite the repo being Go-native. The α SKILL did not name "language" as a δ-pinned axis; α had room to improvise; β's AC oracles were behavior-only ("does V validate? does it reject counterfeit receipts?") and did not surface the implementation-contract drift.
- **cnos#391** (wrong package scoping + separate binary): α placed the Go port in a separate binary at the wrong package path. Same root cause: α improvised on package scoping and CLI integration; β APPROVED on behavior-only ACs.
- **cnos#392** (Phase 3 remediation v2): δ pinned the 7-axis implementation contract at dispatch as an ad-hoc operator action. The cycle succeeded specifically *because* the contract was pinned; F1–F4 in cnos#392's `cdd-iteration.md` forecast the four patches this cycle ships.

§Gap is empirical (not "X does not exist"); the surface this cycle adds is genuinely additive (Patch A adds a new rule 3.6; Patch B adds a new Rule 7; Patch C extends an existing template; Patch D adds a new operator section), and `rg "Implementation contract" src/packages/cnos.cdd/skills/cdd/{alpha,beta,gamma,operator}/SKILL.md` returns 0 hits today.

## Expected diff scope

- 4 SKILL.md files modified, additive only (no rule renumbering of existing rules; no rewrite of existing template rows; no removal of existing operator sections).
- 1 GitHub issue body edited (cnos#366 Phase 4 section).
- `.cdd/unreleased/393/` populated per CDD §1.4 + closure gate (gamma-scaffold, design-notes, self-coherence, beta-review, cdd-iteration, alpha-closeout, beta-closeout, gamma-closeout) plus INDEX.md row.
- Backward compat: all existing rules (α 3.1–3.5, β 1–6, γ §2.5 dispatch-prompt template existing rows, δ existing sections) preserved verbatim.

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Markdown (skill files with YAML frontmatter) |
| CLI integration target | N/A (skill files; not commands) |
| Package scoping | `src/packages/cnos.cdd/skills/cdd/{alpha,beta,gamma,operator}/SKILL.md` (4 files) |
| Existing-binary disposition | N/A (existing skill files; additive patches; preserve existing rules) |
| Runtime dependencies | None |
| JSON/wire contract | N/A |
| Backward compat | Existing rules (1–6 for β; 3.1–3.5 for α; existing γ §2.5 template; existing operator/δ sections) preserved; new content additive |

Eating our own dog food — this scaffold is dispatched *with* a fully populated `## Implementation contract` section, which is exactly what Patch C will mandate going forward.
