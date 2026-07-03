---
issue: 568
cycle_branch: cycle/568
role: gamma
---

# γ closeout — issue #568

## Cycle summary

cds/issues FSM Phase 1 — read-only issue-state reconciler. Shipped: a
declarative transition table (`src/packages/cnos.cds/skills/cds/fsm/transitions.json`,
package-owned by `cnos.cds`) consumed by a generic, CDS-agnostic evaluator
engine (`src/packages/cnos.issues/commands/issues-fsm/`, package-owned by
`cnos.issues`, new Go module mirroring the `issues-map` co-location + kernel
shim pattern); a compiled-in `cn issues fsm evaluate --issue {N}` CLI
(`--fixture` offline path + live GitHub REST/`git` assembly) that prints
facts + decision with zero mutation; 9 fixtures covering the AC3–AC7
failure modes (empty review, dead-in-progress with/without matter, repair
re-entry gating, idempotence); a new CI step wired into `build.yml` (AC9).
25 files changed, all inside the declared scope — no protocol/prompt files
touched, no label-write code path introduced, no unrelated files modified.

**R0-only convergence.** α implemented all 10 ACs and wrote
`self-coherence.md §R0` with per-AC file:line evidence, including three
honestly-flagged partial-verification notes (live-path test coverage,
`repair_contract_present` heuristic on the live path, I5/cue unavailable
locally). β independently re-derived every AC — reading the table and
engine directly, building the binary and driving the CLI against fixtures
itself, grepping for label-write signatures itself, diffing the protected
protocol files itself — and returned `verdict: converge` at R0 with zero
AC violations and no scope creep. No iteration round was needed; the
scaffold's AC→oracle→surface table gave β a complete, mechanically
checkable list to walk, and α's implementation matched it on the first
pass.

## Process friction

**Scaffold sufficiency: high.** The γ scaffold's implementation contract
(package boundary, CLI resolution mechanics, dispatch-boundary rule,
stdlib-only constraint, exact file paths for the transition table and the
new module) left α very little room to improvise the load-bearing axes.
β's review credits this directly — it was able to check "does the engine
avoid switching on CDS state names" and "is the table pure JSON" as
grep-provable oracles because the scaffold specified both the check and
where to look for it, rather than leaving AC1's abstract "declarative
source of truth" wording for α/β to individually interpret.

**No gaps β had to fill.** β's review shows no instance of β having to
invent an oracle the scaffold or issue didn't supply — every AC1–AC10
check maps to an existing scaffold row or issue AC text. The one place β
did original analysis beyond the checklist (AC5's rule-ordering proof —
confirming first-match-wins semantics in `Evaluate` mean the δ-recovery
rule is genuinely checked before the requeue fallback, not just that the
fixture happens to produce the right output) was warranted scrutiny on
the highest-stakes AC (the cnos#368 regression guard), not a scaffold
gap.

**Environment constraint, not a cycle gap.** Both α and β independently
noted the same two environment limits (no CI runner available to either
of them; `cue` unavailable for I5 skill-frontmatter-check). Both handled
it the same way — ran the closest available local equivalent and stated
plainly what wasn't verified rather than papering over it. This is
symmetric with prior cycles' documented pattern (e.g. cycle/514's dune/opam
note) and needs no cycle-specific action.

## Triage

| Item | Source | Disposition |
|---|---|---|
| `RunConclusion` captured in `FactSnapshot`, printed in the facts block, but no guard function in `guardFuncs` reads it | β R0, non-blocking observation | Not an AC violation — AC5's oracle as written doesn't require a success/failure distinction on dead runs. Confirmed independently: `RunConclusion` has zero references outside `snapshot.go`'s struct field and its doc comment. |
| `checks_passing` guard function is defined (`table.go`) and documented in `transitions.json`'s guard glossary, but never referenced by any `all_true`/`any_true`/`all_false` list in any transition rule | Found during this closeout audit, extending β's observation | Same disposition as above — defined but currently inert. Confirmed by grep: `checks_passing` appears exactly once in `transitions.json` (the guard-name glossary at line 28), never inside a `rules[]` condition. |
| `cdd_artifacts_present` guard function likewise defined but unreferenced by any rule | Found during this closeout audit | Same pattern as the two above. |
| Live-path `fetch.go` has no dedicated unit test; `repair_contract_present` on the live path is a substring heuristic, not a full β/δ-findings parse | α R0 honesty note, independently reconfirmed by β | Accepted as mirroring existing `issues-map` precedent's test-coverage shape, not a new gap this cycle introduced. Fixture path (what AC6 actually gates) is unaffected. |
| I5 (skill-frontmatter-check), I2, I4, install-wake-golden, dispatch-repair-preflight, dispatch-closeout-integrity not re-run in an actual CI runner by either α or β | Both R0 artifacts | Structurally inferred safe (no changed file falls inside those gates' known scope) but not empirically observed by either role. Genuinely unresolved until CI runs on the PR — δ should treat CI as the closing check for these gates, not this closeout. |

**Cleanly closed:** AC1–AC10 all independently verified twice (α + β) against
running code, not prose. Phase-1 boundary (AC8: zero mutation, zero
protocol-file diff, zero new status labels) holds under direct grep and
diff, confirmed by both roles separately. No scope creep — β's triple-dot
merge-base diff confirms exactly the declared file set changed.

**Genuinely unresolved (not this cycle's job to resolve):** the CI-runner
gates listed above will only be empirically confirmed when the PR's actual
CI runs — that is δ's/the PR pipeline's job, not a gap in this cycle's
work. The three captured-but-unconsumed guard predicates
(`RunConclusion`/`checks_passing`/`cdd_artifacts_present`) are the one item
with a forward-looking action: **recommend δ file a follow-up issue (or at
minimum leave an explicit note on #569) for Phase 2 (#569) awareness** —
these three facts are already assembled into every `FactSnapshot` and
already load-bearing in the CLI's printed output, but no Phase-1 rule
consumes them. Phase 2 grants the FSM authority to apply labels; before
that flip, someone should deliberately decide whether these three
predicates are intentional headroom (e.g. distinguishing a *failed* dead
run from a *succeeded* one before proposing δ-recovery, or requiring CI to
be green before proposing a state advance) or genuinely dead code to prune.
This is a design decision, not a Phase-1 defect — γ is not filing the
issue, only recommending it for δ's judgment.

## Signoff

R0 converge, no iteration required. cnos#568 ready for δ to land the
closeout triad and open the PR.
