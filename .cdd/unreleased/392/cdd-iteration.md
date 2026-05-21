# cdd-iteration — cycle/392

**Issue:** cnos#392
**Date:** 2026-05-21
**Cycle:** 392
**Mode:** design-and-build, γ+α+β-collapsed-on-δ
**Findings:** 4 (F1–F4, all operator-named at dispatch)
**Patches:** 0 (all four findings dispositioned to next-MCA; ε will fold them into a coordinated cdd patch)
**MCAs:** 4
**No-patch:** 0

Per `cdd/post-release/SKILL.md` §5.6b, this file records protocol-gap and skill-gap findings produced during the cycle.

---

## F1 — `cdd-protocol-gap`: α cannot decide implementation language/architecture

### Finding
α cannot decide implementation language/architecture. The α role is to execute against an issue's ACs; when an issue under-specifies architectural choices (language, CLI integration target, package location), α must either improvise (regression-prone) or surface (slow). Cycle #389 collapsed at exactly this: α picked Python despite the repo being Go-native.

### Root cause
The α SKILL doctrine does not explicitly forbid α from making implementation-language and architectural choices. The δ role at dispatch was not explicitly recognized as the "architect" tier responsible for pinning these choices upstream. cycles #389 and #391 are both consequences of this protocol gap.

### Evidence
- cnos#389 closeouts and PR — V shipped in Python, contradicting cnos's Go-native v3→v4 direction.
- cnos#391 dispatch and closeout — Go port attempted but with under-specified scope (package location, CLI integration); operator caught both during scaffold review.
- This cycle (#392) succeeded specifically because δ-as-architect pinned 7 axes at dispatch. The contrast with #391 is the test.

### Disposition: `next-MCA`
Patch `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` with an explicit rule:

> **α MUST NOT change implementation language or architecture.** Implementation language, CLI integration target, package scoping, and other axes of the implementation contract are pinned by δ-as-architect at dispatch. When the issue under-specifies these, α MUST surface (not improvise).

ε will fold this into the coordinated MCA covering F1–F4 + the Phase 4 δ-split scope refinement.

---

## F2 — `cdd-protocol-gap`: β must check implementation-contract coherence

### Finding
β review under β-α-collapse approved wrong-architecture implementations in cycles #389 (Python on a Go-native repo) and #391 (wrong location, wrong CLI) because the oracle sub-checks tested behavior (does V work?), not implementation contract (is it built the right way?). The structural review-independence collapse acknowledged in mode: collapsed receipts compounds this — there's no second pair of eyes catching the architectural drift.

### Root cause
The β SKILL doctrine does not explicitly require β to verify implementation-contract conformance before APPROVE. Behavior-only oracles let architectural drift through.

### Evidence
- cnos#389 β-α-collapse review: AC oracles passed; receipt valid; Python architecture survived to merge.
- cnos#391 β review: scope under-specified at dispatch; rescoped after the fact.
- This cycle (#392) explicitly required β to walk the 7-axis implementation-contract table — and that catch-net is what α's β-α-collapse review needs.

### Disposition: `next-MCA`
Patch `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` adding **Rule 7 — Implementation-contract coherence**:

> Before APPROVE, β MUST verify the cycle's implementation conforms to the implementation contract pinned at dispatch (language, CLI surface, package scoping, runtime deps, JSON/contract preservation, backward-compat invariant). When the issue body contains an `## Implementation contract` section, β walks each axis and confirms conformance with concrete file:line evidence.

ε folds into the coordinated MCA.

---

## F3 — `cdd-skill-gap`: γ dispatch-prompt template lacks an implementation-contract section

### Finding
γ's dispatch-prompt template (the prompt operator/δ writes to launch a cycle) does not have a standard `## Implementation contract` section. cycles #389 and #391 dispatch prompts under-specified architectural decisions; #392's dispatch fixed this by adding the pinned-axes table inline, but the template itself doesn't formalize the pattern.

### Root cause
`cnos.cdd/skills/cdd/gamma/SKILL.md` §2.5 (dispatch-prompt format) does not enumerate implementation-contract axes. Without a template slot, dispatch authors don't know to fill it.

### Disposition: `next-MCA`
Patch `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` §2.5 adding an `## Implementation contract` section to the dispatch-prompt format. The section enumerates:
- Language
- CLI integration target
- Package scoping
- Existing-binary disposition
- Runtime deps
- JSON/contract preservation
- Backward-compat invariant

(matching the 7 axes from this cycle's issue body).

ε folds into the coordinated MCA.

---

## F4 — `cdd-protocol-gap` (NEW since #391): δ is a two-sided membrane

### Finding
δ is documented in `operator/SKILL.md` as a one-sided membrane — outward (verdict-on-receipt); operator-facing trust surface. But cycle #392's dispatch made visible an **inward** δ: δ-as-architect, the role that fills the `## Implementation contract` section in γ's dispatch prompts. This inward δ is structurally distinct from the outward δ and is what F1–F3 are pointing at.

### Root cause
`src/packages/cnos.cdd/skills/operator/SKILL.md` documents only outward δ (verdict-on-receipt at merge boundary). The inward δ (implementation-contract enrichment at dispatch) is not named, not skilled, and not protocol-codified.

### Disposition: `next-MCA`
Folds into **Phase 4 (δ split)** scope refinement: the Phase 4 issue captures δ as a two-sided membrane explicitly:
- **outward δ (current):** verdict-on-receipt at merge.
- **inward δ (new):** implementation-contract enrichment at dispatch.

Forthcoming standalone tracking issue captures this; ε will reference it from the coordinated MCA.

---

## Aggregator-update note (for ε)

INDEX.md row to add:

```
| 392 | #392 | 2026-05-21 | 4 | 0 | 4 | 0 | .cdd/unreleased/392/cdd-iteration.md |
```

(Findings: 4 — F1, F2, F3, F4. Patches: 0 in this cycle. MCAs: 4 deferred to ε's coordinated patch. No-patch: 0.)
