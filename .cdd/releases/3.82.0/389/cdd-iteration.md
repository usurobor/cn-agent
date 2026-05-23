# cdd-iteration — Cycle #389 (Phase 3 of #366)

**Cycle:** #389 — implement V (`cn-cdd-verify` rewrite — Contract × Receipt → ValidationVerdict)
**Merge:** [FILL IN AT MERGE TIME]
**Closed:** 2026-05-21
**Mode:** design-and-build (γ+α+β-collapsed-on-δ per `breadth-2026-05-12` wave manifest precedent; β-α-collapse acknowledged in `beta-review.md` and on the cycle's own CDS receipt via `mode: collapsed`)
**Rounds:** R1 APPROVE (no fix-round)
**ACs:** 8/8 PASS (37/37 oracle sub-checks in `tests/cdd/test_cn_cdd_validate_receipt.sh`)

Per `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b`: cycle's `cdd-*-gap` findings recorded so ε can aggregate without re-parsing prose. Written even when findings count would be zero (closure-gate hygiene per #388's ε iteration close-side; not repeated here).

## §1 Findings dispositioned

### F1 — `cdd-protocol-gap`: receipt-level `mode` field not yet typed in schemas

- **Source:** cycle 389 β-review §Findings F0 (surfaced and addressed in-cycle as a V-implementation feature)
- **Class:** `cdd-protocol-gap`
- **Trigger:** ε process-gap check during cycle close-out
- **Description:** Cycle #389 introduces a receipt-level `mode: "collapsed"` field. V's counterfeit rules C1 (actor_separation) and C2 (verdict_precedes_merge) honor it: structural review-independence violations emit warnings rather than failed_predicates when this field is declared. The field is structurally accepted today via `#Receipt`'s open-extension (`...`) but is not pinned in `schemas/cdd/receipt.cue`. A typed `mode?: "collapsed" | "distributed"` field would: (1) document the field as load-bearing in the schema itself; (2) let CUE reject unknown values (e.g., `mode: "partial"`); (3) make the field discoverable to future receipt authors.
- **Root cause:** the field was discovered as a design need *during* the build phase of #389 when V's strict counterfeit rules conflicted with the wave-manifest β-α-collapse mode. Adding the field to V was the right in-cycle move; pinning it in the schema would have expanded scope.
- **Disposition:** `next-MCA`. Patch shape: add `mode?: "collapsed" | "distributed"` to `#Receipt` in `schemas/cdd/receipt.cue`; document under `schemas/cdd/README.md`. ~30 lines.
- **Issue filed:** none yet. Suggest a small standalone cycle, or fold into Phase 5 γ-shrink (which will touch the receipt-authoring surface anyway).
- **First AC for the eventual MCA:** `schemas/cdd/receipt.cue` declares `mode?: "collapsed" | "distributed"` as a typed enum; `cue vet -c -d '#Receipt' <fixture-with-mode>` exits 0; `cue vet -c -d '#Receipt' <fixture-with-mode: invalid>` exits 1.

### F2 — `cdd-tooling-gap`: V does not yet resolve `git_diff(...)` synthetic refs

- **Source:** cycle 389 design-notes §Risks + alpha-closeout §Debt item D1
- **Class:** `cdd-tooling-gap`
- **Trigger:** ε process-gap check during cycle close-out
- **Description:** V's `is_filesystem_path` filter excludes `git_diff(...)` synthetic refs (and `sha256:`, `issue:#`, `@*_placeholder`) from C4 dereference. They pass through silently as non-paths. Future scope-drift detection — V checking that the actual cycle diff respects the contract's `non_goals` — would need a `git_diff(origin/main..merge_sha)` resolver.
- **Disposition:** `next-MCA`. Out of scope for Phase 3 (V's hard ACs do not require diff resolution). Best landed alongside Phase 4 δ-split (which is the surface that benefits from scope-drift detection) or Phase 5 γ-shrink.
- **Issue filed:** none yet.
- **First AC for the eventual MCA:** V's C-rule set includes a `git_diff_resolves` predicate that runs `git diff <base>..<merge>` and asserts that all modified files match at least one entry in the contract's `scope.in_scope` (and no entry in `scope.out_of_scope` or `non_goals`).

### F3 — design-finding (within-cycle, addressed): `mode: collapsed` field

- **Source:** cycle 389 β-review §Findings F0
- **Class:** none (within-cycle engineering refinement; surfaced and resolved before merge)
- **Description:** First V run against `schemas/cds/fixtures/valid-receipt.yaml` triggered C1 and C2 because cycle 369's α and β closeouts shared an actor (`gamma@cdd.cnos`) and the synthetic `decided_at` timestamp predated the actual β commit. This is correct behavior under strict rules but wrong under the established γ+α+β-collapsed-on-δ mode. Resolved by adding the `mode: collapsed` receipt-level field (described in F1 above as the structural follow-up).
- **Disposition:** `drop` (within-cycle refinement; recorded for posterity). The follow-on F1 captures the future schema-formalization debt.

## §2 No-findings observations (informational)

- 37/37 AC oracle sub-checks PASS in `tests/cdd/test_cn_cdd_validate_receipt.sh` on first run. No mechanical-overload signal.
- `cue vet` works cleanly across all fixtures (valid → exit 0; invalid → exit 1 with structural diagnostic).
- Backward compatibility with existing `cn cdd-verify --version/--pr/--all/--unreleased/--triadic/--cycle` modes preserved (verified via BC.a/b in the harness).
- The cycle's own deliverables can be V-validated post-merge: a CDS receipt for #389 will declare `mode: collapsed` and dereference its own `.cdd/unreleased/389/*.md` evidence cleanly.

## §3 Trigger assessment (per `gamma/SKILL.md` §2.8 table)

| Trigger | Fire condition | Fired? | ε note |
|---|---|---|---|
| Review churn | review rounds > 2 | **No** | R1 APPROVE on first pass. |
| Mechanical overload | mechanical ratio > 20% AND findings ≥ 10 | **No** | 3 findings; F1+F2 named for follow-up; F3 in-cycle refinement. |
| Avoidable tooling/environment failure | environment blocked the cycle | **No** | `cue vet`, `python3`, PyYAML, `jq`, `yq`, `git` all available; no harness pattern fired. |
| Loaded-skill miss | a loaded skill should have prevented a finding | **No** | F1 emerged from V's first contact with the wave-manifest precedent; the convention is fresh from this cycle. F2 was scoped out by design. F3 was caught and addressed before merge. |

## §4 INDEX update

Add to `.cdd/iterations/INDEX.md` (columns: Cycle | Issue | Date | Findings | Patches | MCAs | No-patch | Path):

```
| 389 | #389 | 2026-05-21 | 3 | 0 | 2 | 1 | .cdd/unreleased/389/cdd-iteration.md |
```

Findings=3 (F1, F2, F3). Patches=0 (no within-cycle iteration-finding patches — F3's `mode: collapsed` resolution is a design-feature addition, not a `cdd-iteration` patch). MCAs=2 (F1, F2 → next-MCA). No-patch=1 (F3 → drop).
