<!-- sections: [Cycle, Verdict, Diff Summary, Decisions, Findings Encountered, Debt, Skills-That-Helped, Outputs] -->
<!-- completed: [Cycle, Verdict, Diff Summary, Decisions, Findings Encountered, Debt, Skills-That-Helped, Outputs] -->

# α Closeout — Cycle #388

**Cycle:** [#388](https://github.com/usurobor/cnos/issues/388) — Phase 2.5 (#366): split generic `schemas/cdd/` from CDS/CDR domain evidence
**Date:** 2026-05-21
**Branch:** `cycle/388`
**Implementation SHA:** `7294b473d849740f9fbae252d267b3396fe27497`
**Base SHA:** `7cdcd8d4e4d17340734d0c9df3786d499655c0b0`
**α identity:** alpha / alpha@cdd.cnos (γ+α+β-collapsed on δ, single Claude Code session)

---

## §Verdict (α-self-statement)

The schema refactor lands cleanly. All 7 ACs are mechanically verified (cue vet + rg + ls). The architectural-choice decision is durably recorded in `schemas/cdd/README.md §"Architectural choice"`. The fixture corpus is correctly re-categorised. The CDS / CDR / generic schemas all compose via CUE unification (`cdd.#Receipt` is the parent definition). `cue vet` runs regression-free.

α signals review-readiness at SHA `7294b473`. β proceeds.

## §Diff Summary

```
 .cdd/unreleased/388/design-notes.md                | 318 ++++++++++++++++++
 .cdd/unreleased/388/gamma-scaffold.md              | 225 +++++++++++++
 .cdd/unreleased/388/self-coherence.md              | 159 +++++++++
 cue.mod/module.cue                                 |  16 +
 schemas/cdd/README.md                              | 357 ++++++/-184 --
 schemas/cdd/contract.cue                           |   3 +/-
 schemas/cdd/fixtures/invalid-fail-no-boundary-*    |  26 +/-
 schemas/cdd/fixtures/invalid-gamma-preflight-*     |  24 +/-
 schemas/cdd/fixtures/invalid-override-masks-*      |  30 +/-
 schemas/cdd/fixtures/valid-generic-receipt.yaml    |  45 +
 schemas/cdd/fixtures/valid-receipt.yaml            |  35 --  (moved to cds/)
 schemas/cdd/receipt.cue                            |  70 +/-
 schemas/cdr/README.md                              | 137 +
 schemas/cdr/fixtures/valid-cdr-receipt.yaml        |  63 +
 schemas/cdr/receipt.cue                            |  74 +
 schemas/cds/README.md                              | 114 +
 schemas/cds/fixtures/valid-receipt.yaml            |  48 +
 schemas/cds/receipt.cue                            |  57 +
 18 files, 1617 insertions(+), 184 deletions(-)
```

## §Decisions (recorded during α work)

1. **Use CUE module + import** rather than inline-mirroring the generic shape in each domain package. Trade-off: adds `cue.mod/module.cue` at repo root. Verified empirically that this does not break existing `tools/validate-skill-frontmatter.sh --self-test`.

2. **Widen `evidence_refs` value-type to a union.** Original primitive `[string]: #EvidenceRef` (string-only) was too narrow for CDR's list-shaped fields. Created `#EvidenceRefValue: #EvidenceRef | [...#EvidenceRef]`. Preserves type discipline.

3. **CDS field names: closure-record shape, not §4a.3 sketch.** Cycle/388 implements what #369 actually shipped (`self_coherence`, `beta_review`, three closeouts, `diff`, optional `ci`) rather than §4a.3's forward-looking names (`artifact_refs`, `test_refs`, `ci_refs`, `diff_ref`, `debt_refs`). Divergence named as known debt in `schemas/cds/README.md`.

4. **CDR list-fields require min-length-one.** Per ROLES §4a.3 ("no data_refs → V rejects research receipt"), the four CDR evidence_refs keys (claim, data, method, result) are encoded as `[_, ...cdd.#EvidenceRef]` to enforce non-emptiness structurally.

5. **Drop `evidence_root_ref` from generic `#Receipt`.** It pointed at a CDS-process-evidence root; CDR has a different (per-wave) root structure. Migrated `evidence_root` into `schemas/cds/`'s `evidence_refs` map (since it's CDS-specific). The generic schema gets `evidence_refs` only — domain-specific keys decide where the root ref lives.

6. **`contract.cue` comment rewritten.** AC2 oracle initially returned 1 hit (a comment in `contract.cue` mentioning `diff_ref`). Rewrote the comment to use generic language (`evidence_refs.diff (per schemas/cds/)`); AC2 now returns 0 hits.

7. **Generic-fixture protocol_id = `cnos.cdd.generic.receipt.v1`.** Used in all four `schemas/cdd/fixtures/` files. Documents this as the "no domain overlay" dispatch identifier in `schemas/cdd/README.md`.

8. **Essay update deferred.** Folding the (a) decision back into `CCNF-AND-TYPED-TRUST.md §"Open questions"` is out of #388's ACs. Deferred to a follow-up doc cycle.

## §Findings Encountered

**F1 — CDR `evidence_refs` shape conflict with original generic primitive.**
- Discovery: First-pass `schemas/cdr/receipt.cue` failed `cue vet` with `conflicting values string and [...#EvidenceRef]` on each of `claim/data/method/result` — because the generic `evidence_refs: [string]: #EvidenceRef` declared a string-value map.
- Resolution: Widened to `#EvidenceRefValue: #EvidenceRef | [...#EvidenceRef]`. Both CDS (singular keys) and CDR (list keys) unify against the generic primitive.
- Lesson: A "typed open primitive" that allows the domain to choose cardinality requires a union type at the generic layer.

**F2 — CUE module file is required for cross-package imports.**
- Discovery: Initial attempt `cue vet schemas/cdd/file.cue schemas/cds/file.cue ...` failed with `found packages "cdd" and "cds" in "/home/user/cnos"`. CUE refuses to mix packages on the command line without a module.
- Resolution: Added `cue.mod/module.cue` at repo root. With it in place, `import "cnos.dev/cnos/schemas/cdd"` from `schemas/cds/receipt.cue` resolves automatically; `cue vet -c -d "#CDSReceipt" schemas/cds/receipt.cue schemas/cds/fixtures/valid-receipt.yaml` works without explicitly listing the cdd files.
- Lesson: Recorded in design-notes §Proposal and `schemas/cdd/README.md §"How to run"`.

**F3 — AC2 false-positive on `contract.cue` comment line.**
- Discovery: `rg "self_coherence_ref|beta_review_ref|alpha_closeout_ref|beta_closeout_ref|gamma_closeout_ref|diff_ref|ci_refs" schemas/cdd/` returned 1 hit at `schemas/cdd/contract.cue:33` — a comment describing the v1 derivation of `diff_ref`.
- Resolution: Rewrote the comment to use generic language (`evidence_refs.diff (per schemas/cds/)`). AC2 now returns 0 hits.
- Lesson: The AC oracle is strict about *occurrence* of the field names, even in comments. Comments referencing "what was in v1" should be rewritten to point at the new location, not preserve old names.

## §Debt

Itemized at `self-coherence.md §Debt`:

1. ROLES §4a.3 CDS sketch vs `schemas/cds/` realized shape divergence — named in `schemas/cds/README.md`.
2. CDR schema is a skeleton; cnos#376 Sub 1 / Sub 2 may refine — named in `schemas/cdr/README.md`.
3. Essay v0.1.0 open-question #1 resolution not folded into essay — deferred follow-up doc cycle.
4. γ skill does not yet emit `protocol_id` automatically — Phase 3 or γ-skill follow-up may address.

## §Skills-That-Helped

- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — L7 output format informed the "Architectural choice" section structure (decision + rationale + design-source citations + leverage / negative leverage). The skill's §2.4 "Proposal" guidance ("start with types") shaped the `#Receipt` / `#CDSReceipt` / `#CDRReceipt` typed-inheritance idiom decision.
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — AC oracle authoring guided the per-AC oracle in self-coherence and beta-review (each oracle is a specific file + specific check + expected outcome).
- `src/packages/cnos.eng/skills/eng/code/SKILL.md` — CUE refactor discipline; CUE inheritance idiom; module file conventions.
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` — design source; the essay's §"Schema direction" was directly cited in the schemas/cdd/README.md decision-record.

## §Outputs

**Schema files (CUE):**
- `cue.mod/module.cue` (new) — repo-level CUE module declaration
- `schemas/cdd/receipt.cue` (modified) — generic `#Receipt`, `#EvidenceRef`, `#EvidenceRefValue`, `#ProtocolGapRef`; CDS field names removed; `protocol_id` added as required field
- `schemas/cdd/contract.cue` (one-line comment fix only)
- `schemas/cds/receipt.cue` (new) — `#CDSReceipt: cdd.#Receipt & {...}` with closure-record evidence keys
- `schemas/cdr/receipt.cue` (new) — `#CDRReceipt: cdd.#Receipt & {...}` with claim/data/method/result lists, `#ClaimStatus`, `#Reproduction`

**READMEs:**
- `schemas/cdd/README.md` (rewritten) — three-package layout; "Architectural choice" decision record (AC1); preserved Scope-Lift Invariant; how-to-run for all three packages
- `schemas/cds/README.md` (new) — CDS package description; §4a.3 vs realized debt
- `schemas/cdr/README.md` (new) — CDR package description; skeleton-not-frozen debt

**Fixtures:**
- `schemas/cdd/fixtures/valid-generic-receipt.yaml` (new) — generic-only receipt
- `schemas/cdd/fixtures/invalid-*.yaml` (×3, modified) — trimmed CDS keys, added protocol_id
- `schemas/cds/fixtures/valid-receipt.yaml` (moved from cdd/, reshaped) — CDS receipt with nested evidence_refs map
- `schemas/cdr/fixtures/valid-cdr-receipt.yaml` (new) — minimal CDR receipt

**Cycle artifacts:**
- `.cdd/unreleased/388/gamma-scaffold.md` — γ pre-dispatch scaffold
- `.cdd/unreleased/388/design-notes.md` — L7 design record
- `.cdd/unreleased/388/self-coherence.md` — α self-coherence with per-AC oracle evidence
- `.cdd/unreleased/388/beta-review.md` — β review with β-α-collapse acknowledgment
- `.cdd/unreleased/388/alpha-closeout.md` — this file
- `.cdd/unreleased/388/beta-closeout.md` — to be authored
- `.cdd/unreleased/388/gamma-closeout.md` — to be authored
