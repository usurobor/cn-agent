# cdd-iteration — Cycle 423

**Cycle:** [cnos#423](https://github.com/usurobor/cnos/issues/423) — Build-fix for cnos#416 cross-repo stub missing triggers
**Date:** 2026-05-23
**Mode:** docs-only / build-fix (P0)
**Authority:** `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` (canonical home of cdd-iteration shape; per-finding shape; cadence rule; INDEX row format).

## §0 Findings

`protocol_gap_count: 1`. One binding `cdd-skill-gap` finding (loaded-skill miss). One patch landed this cycle (the α-423 frontmatter fix). One `next-MCA` disposition recorded for future operator selection. Zero `no-patch` dispositions.

## §1 The finding

### Class

`cdd-skill-gap` — loaded-skill miss.

### Trigger

The activation validator at `src/go/internal/activation/Validate` emits `IssueEmptyTriggers` for any SKILL.md whose frontmatter has empty `triggers:`. The test `TestValidate_RealCorpus_NoEmptyTriggers` (at `src/go/internal/activation/index_test.go:312`) walks the live package corpus and asserts that no `IssueEmptyTriggers` issues surface. Post-merge of cycle/422 (which closed v3.82.0), CI went red on `main`:

```
--- FAIL: TestValidate_RealCorpus_NoEmptyTriggers (0.30s)
    index_test.go:312: empty-triggers against real corpus: package cnos.cdd: skill cdd/cross-repo has no triggers
    index_test.go:319: real-corpus validation: 85 skills copied, 27 issues surfaced
FAIL    github.com/usurobor/cnos/src/go/internal/activation
```

The offending skill is the cnos#416-authored compatibility stub at `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` — a 28-line pointer to the new canonical home at `cnos.handoff/skills/handoff/cross-repo/SKILL.md`. Its frontmatter declares `artifact_class: pointer` / `status: moved` / `canonical: cnos.handoff/skills/handoff/cross-repo/SKILL.md` but omits `triggers:`.

### Description

cnos#416's dispatch brief was thorough on the canonical-side authoring (a 643-line handoff SKILL.md), citation re-pointing across five cdd-side skills, and extraction-map updates. But it did not enumerate the activation validator's "every SKILL.md frontmatter must ship with non-empty `triggers:`" requirement. The agent authored a structurally-correct deprecation stub — pointer narrative in the body; canonical / status / artifact_class fields in frontmatter — but with no `triggers:` field. The cnos#416 self-coherence pass-set did not catch this because the cycle's ACs were extraction-mechanics (lines moved, citation re-points landed) rather than "all packaged skills pass `internal/activation/Validate`". CI on `main` post-cycle/422 merge surfaced the gap as a hard failure that this cycle (cnos#423) now patches.

### Root cause

**Dispatch-brief incomplete.** The cnos#416 dispatch brief did not load `internal/activation/Validate` requirements as a Tier 3 reference. Specifically, the brief did not surface:

1. The contract that **every** packaged SKILL.md must have non-empty `triggers:` (no exemption for `artifact_class: pointer` skills).
2. The test that enforces this contract against the live corpus.
3. The pattern that any deprecation stub must therefore carry `triggers:` even if the stub is purely a backward-compat breadcrumb with no operational role.

This is a `cdd-skill-gap` (the stub-authoring pattern documentation did not load the activation/Validate contract) classified by trigger as a **loaded-skill miss** (the loaded skills — `cnos.handoff/skills/handoff/extraction-map.md`, the dispatch template, the citation-repoint pattern — did not include activation/Validate as a precondition).

### Disposition

**`next-MCA`.** Two architecturally-coherent remediation paths exist; operator picks at next-MCA dispatch:

**Option A — Patch the dispatch template / stub-authoring pattern documentation:**
- Modify the cdd dispatch template (or `cnos.cdd/skills/cdd/activation/SKILL.md`, or wherever the stub-authoring pattern is documented) to require non-empty `triggers:` on all SKILL.md files including pointers.
- This treats the issue as a documentation gap: the contract is sound (every loaded skill needs triggers); future stub-authoring dispatches must ship triggers.
- Pros: zero runtime risk; preserves the validator's strict contract.
- Cons: relies on every future stub author remembering the rule.

**Option B — Patch `internal/activation/Validate` to exempt `artifact_class: pointer` skills:**
- Modify the validator at `src/go/internal/activation/validate.go` to skip the empty-triggers check when `artifact_class: pointer` is set.
- This treats the issue as a contract gap: pointer skills are structurally never loaded by trigger keyword (they're navigation aids, not operational surfaces), so the empty-triggers requirement doesn't apply.
- Pros: architecturally cleaner; encodes the semantic that pointers are not trigger-addressable.
- Cons: more invasive (Go code change + test update + potential cascading exemptions for other `artifact_class` values).

The operator picks **A or B (or both)** at the next-MCA dispatch.

### First AC for the eventual MCA

> Every SKILL.md frontmatter ships with non-empty `triggers:` OR `internal/activation/Validate` skips `artifact_class: pointer` skills from the empty-triggers check.

Falsifiable mechanically (run the corpus survey + the test); concrete (names both remediation paths).

## §2 Trigger assessment (per ε §1 / receipt-stream cadence rules)

- **Review churn:** R1 APPROVED on first pass; no fix rounds. **No additional trigger fired.**
- **Mechanical overload:** P0 build-fix; one-file edit; mechanical ratio not applicable. **No trigger fired.**
- **Avoidable tooling / environment failure:** the underlying failure (CI red on main post-v3.82.0 merge) IS the primary trigger; the build-fix cycle exists to resolve it.
- **Loaded skill failed to prevent a finding:** **trigger fired**. cnos#416's loaded skill set (extraction-map, dispatch template, citation-repoint pattern) did not include activation/Validate's contract. This is the binding finding recorded in §1.

## §3 Wave + release context

This cycle is a **post-release build-fix** following the v3.82.0 release boundary (cycle/422). The release was merged but the tag has not been pushed because CI is red on main. Cycle/423 unblocks the tag push:

- After merge, `go test ./...` is green on `main`.
- δ may push the `v3.82.0` tag via `scripts/release.sh`.
- The post-tag pause-protocol-evolution Stop condition (per cycle/422's RELEASE.md) remains in effect, with one open operator decision: pick the next-MCA disposition (Option A vs B vs both) from §1.

## INDEX update

Adds the following row to `.cdd/iterations/INDEX.md`:

```
| 423 | #423 | 2026-05-23 | 1 | 1 | 1 | 0 | .cdd/unreleased/423/cdd-iteration.md |
```

Filed by γ on 2026-05-23.
