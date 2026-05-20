# α self-coherence — cycle #375

Issue: usurobor/cnos#375 — "γ-side pre-dispatch gate for gamma-scaffold.md (rule 3.11b symmetry; cdd-protocol-gap)".

Branch: `cycle/375`.

Class: skill-patch (docs-only disconnect). Wave: 2026-05-19 protocol hygiene.

## Gap

`review/SKILL.md` rule 3.11b is a **binding β-side** gate: when `.cdd/unreleased/{N}/gamma-scaffold.md` is missing on the cycle branch at review time and no `## Protocol exemption` exists in the sub-issue body, β must return REQUEST CHANGES (D-severity, `protocol-compliance`).

The γ-side parallel surface — `gamma/SKILL.md` §2.5 Step 3a (cycle-branch creation) and Step 3b (α dispatch) — did not block dispatch on the artifact's existence. Verified by inspection of `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` at `origin/main` HEAD (`c9017153`): §2.5 Step 3a contained γ-owned branch pre-flight checks (branch non-existence, no stalled `.cdd/unreleased/{N}/` on main) but no scaffold check; §2.5 Step 3b's existing prose proceeded directly from the heading to the polling paragraph without any scaffold gate.

Peer-enumeration (per `gamma/SKILL.md` §2.2a):

- `rg gamma-scaffold src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` at `c9017153` returns no matches.
- `rg gamma-scaffold src/packages/cnos.cdd/skills/cdd/CDD.md` at `c9017153` returns no matches.
- `rg gamma-scaffold src/packages/cnos.cdd/skills/cdd/review/SKILL.md` at `c9017153` returns multiple matches (rule 3.11b and the artifact-completeness gate row in the review template) — confirming the asymmetry: β-side names the artifact, γ-side does not.

The negation "γ-side gate for `gamma-scaffold.md` does not exist" is therefore backed by grep-evidence per the §2.2a precedent.

Cycle #369 paid the round-trip cost empirically: R1 RC (β rule 3.11b D1 against α review-ready at `6835197d` — bc59a1a3) → γ recovery path (a) authoring `gamma-scaffold.md` at `227d2373` → β R2 APPROVED at `4e179db6` → merged at `ff54f2a0`. The cycle was otherwise review-ready at R1 (10 ACs met, 10/11 contract-integrity rows yes; the only "no" row was the missing scaffold).

## AC mapping

### AC1 — Pre-dispatch check named in γ skill or CDD step 3 as a binding gate

**Implementation:** `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` §2.5 Step 3b now opens with a sub-section "Pre-dispatch γ scaffold check (binding gate)". The text reads: *"γ MUST NOT proceed to α dispatch (the prompt-production + δ-routing flow below) until `.cdd/unreleased/{N}/gamma-scaffold.md` exists on `origin/cycle/{N}`. This is a gate, not advisory prose."*

The sub-section is bound to Step 3b (α dispatch) and includes a bash procedure (`git ls-tree -r --name-only origin/cycle/<N> .cdd/unreleased/<N>/gamma-scaffold.md`) that γ runs before producing the α prompt. Canonical surface chosen: `gamma/SKILL.md` (not `CDD.md`); the issue body lists both as valid choices and §2.5 Step 3b is the dispatch flow's actual home.

**Oracle test:** §2.5 Step 3b contains the literal string `.cdd/unreleased/{N}/gamma-scaffold.md` and the literal string "MUST NOT". Binding-gate language present.

**Status:** PASS.

### AC2 — Symmetry with rule 3.11b documented

**Implementation:** The sub-section's "Dual of `review/SKILL.md` rule 3.11b" paragraph names rule 3.11b explicitly and frames the γ-side check as its dual:

> γ-side pre-dispatch enforcement pays the same cost once *per cycle* at scaffold time, before α and β load context. The two gates are symmetric; together they guarantee scaffold presence without either side leaning on the other.

The paragraph also names rule 3.11b's verdict semantics (REQUEST CHANGES, D-severity, `protocol-compliance`) so the cross-reference is concrete, not just a name-drop.

**Oracle test:** the new sub-section contains the literal string `rule 3.11b` and frames the γ-side check as the dual.

**Status:** PASS.

### AC3 — Empirical anchor cited (cycle #369)

**Implementation:** The sub-section's "Empirical anchor — cycle #369" paragraph cites:

- cycle #369
- β R1 D1 (rule 3.11b fire on missing `gamma-scaffold.md`)
- γ recovery commit `227d2373` (with its actual commit message)
- β R2 APPROVED at `4e179db6`
- merge at `ff54f2a0`

SHAs verified locally before being placed in the patch:
- `git log --oneline ff54f2a0` → "γ-369: merge cycle/369 — Phase 2 schema-typing of receptor (CUE schemas + fixtures)"
- `git log --oneline 227d2373` → "γ-369: scaffold — recovery path (a) for β R1 D1 (rule 3.11b)"
- `git log --oneline 4e179db6` → "β-369 R2: APPROVED — D1 closed at 227d2373; pre-merge gate all-green"

**Oracle test:** the sub-section contains the literal string `cycle #369` and the literal SHAs `227d2373` and `ff54f2a0`.

**Status:** PASS.

### AC4 — No CI / runtime / release surface change

**Implementation:** `git diff origin/main..HEAD --stat` shows changes only under:

- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` (the canonical patch surface)
- `.cdd/unreleased/375/*.md` (cycle evidence)

No edits to `CDD.md`, `review/SKILL.md`, `beta/SKILL.md`, `alpha/SKILL.md`, `operator/SKILL.md`, `release/SKILL.md`, any CI workflow under `.github/workflows/`, any validator code under `src/packages/cnos.cdd/cmd/` or `scripts/`, or any other runtime surface.

**Oracle test:** see β-review.md §Artifact completeness — γ runs `git diff origin/main..HEAD --stat` and confirms scope.

**Status:** PASS.

## CDD Trace

| CDD surface | How this cycle conforms |
|---|---|
| `cdd/CDD.md` §1.4 (lifecycle) | γ → α → β → close-out, executed within δ-as-agent under §5.2 wave-mode γ=δ collapse + α=β session-collapse. |
| `gamma/SKILL.md` §2.2a (peer enumeration) | §Gap above cites `rg` evidence backing the negation "γ-side gate does not exist". |
| `gamma/SKILL.md` §2.5 Step 3b (the patched surface) | The patch IS this section; binding gate added with bash procedure, rule-3.11b dual framing, empirical anchor. |
| `review/SKILL.md` §3.11b | Untouched (per non-goals); patch cross-references but does not edit. Exemption surface honored (γ-side defers to the sub-issue `## Protocol exemption` per 3.11b's discoverability rule rather than inventing a parallel γ exemption channel). |
| `cdd/CDD.md` §4.3 (γ-owned branch pre-flight) | Untouched. Branch pre-flight (cycle-branch non-existence, stalled-dir check) is distinct from scaffold pre-flight; the new gate is co-located with α dispatch in Step 3b, not with branch creation in Step 3a, because the artifact must exist by *dispatch time*, not by *branch-creation time* (γ may legitimately push the scaffold to `cycle/{N}` after creating the branch and before producing the α prompt). |
| Wave manifest invariant 3 (`.cdd/waves/2026-05-19-protocol-hygiene/manifest.md`) | "rule 3.11b: γ-scaffold required at dispatch" — this cycle's scaffold `gamma-scaffold.md` was committed at `cfe6c96e` before α implementation. |

## Review-readiness signal

**Branch:** `cycle/375` at the SHA reachable from `git rev-parse cycle/375` after this commit lands.

**State:**

- AC1: PASS — binding gate language in `gamma/SKILL.md` §2.5 Step 3b.
- AC2: PASS — rule 3.11b cross-reference + dual framing.
- AC3: PASS — cycle #369 + SHAs `227d2373` and `ff54f2a0` cited verbatim.
- AC4: PASS — diff scope confirmed under `gamma/SKILL.md` + `.cdd/unreleased/375/` only.

**Fix rounds:** R1 only (β-collapsed self-review runs next; if any AC fails it converts to a fix round here).

**β handoff:** this self-coherence is the β trigger. β-collapsed reviewer reads against AC1–AC4 oracles above and issues verdict in `beta-review.md`.
