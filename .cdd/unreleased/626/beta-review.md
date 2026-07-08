# beta-review.md — cnos#626 (β, R0)

## Scope of this review

Independently re-derived (not trusted from self-coherence.md) against
`gamma-scaffold.md`'s "Per-AC oracle list" and the β prompt's "Independent
AC walk." Base SHA: `86042ec5be4b5fb45b213c27dfcf635958f60aac`. Reviewed
commits on `cycle/626`: `51c2c1c` (α implementation) + `87c9ed1`
(self-coherence.md). All commands below were re-run by β directly against
the checked-out `cycle/626` tree, not copy-pasted from α's claims.

## AC1 — doctrine naming the cell/substrate boundary

**MET.** Read `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.12
"cell/substrate identity boundary (cnos#626)" (new, lines 607–624, sibling
to §9.10/§9.11, correctly positioned after §9.9 "Cross-references" at the
end of §9). It explicitly states, in its own opening sentence: *"the cell
carries no substrate identity and no hub-state scope. The wake/substrate
owns identity, channel logs, and reporting; δ and the roles it dispatches
(γ/α/β) do not."* Both halves of the required statement are present
verbatim, not implied.

It cross-references §9.2's five-input contract by name (verified §9.2
exists at line 395 and actually lists exactly the five named inputs:
claimed-issue-number, protocol identifier, current main SHA, wake run id,
package-runtime context — confirmed by direct grep, not taken on faith)
and names the Persona/Operator/hub-state omission as intentional in its
own subsection heading ("§9.2's silence is intentional, not an
oversight"), stating: "the absence is doctrine, not debt." This is an
explicit statement, not a vague inference — passes the bar the β prompt
set for distinguishing a real pass from an implied one.

## AC2 — dispatch cell-worker no longer instructed to run full six-step activation

**MET (first disjunct).** Independently ran:

```
grep -n "Activate per.*activate/SKILL.md" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md
```

→ 0 hits. The old instruction ("Activate per activate/SKILL.md: Kernel →
CA skills → Persona → Operator → hub state → identity confirmation") is
gone, replaced by doctrine stating cell-execution cognition loads Kernel +
CA skills only, does not load Persona/Operator/hub-state, and does not
read `.cn-{agent}/` or thread surfaces. The "and/or" wording of AC2's
issue text means the first disjunct (no full activation) alone is a
legitimate pass; the second disjunct (checkout-scope exclusion) is
correctly not claimed here — it is explicitly deferred to AC3.

Grepped for "Persona"/"Operator"/"hub state" in the section as instructed:
these strings do appear, but only in negation ("does NOT load Persona,
Operator, or hub state") and in the unmodified reference sentence
describing what the admin wake still does. No instruction to load them
for cell execution remains. Confirmed the preserved "You do NOT attach to
a channel like the admin wake does" sentence is retained verbatim, not
duplicated or reworded.

Confirmed `activate/SKILL.md` itself is untouched:
`git diff --stat 86042ec5be4b5fb45b213c27dfcf635958f60aac..HEAD -- src/packages/cnos.core/skills/agent/activate/SKILL.md`
→ no output. The verbatim quote α's new prose pulls from that file
("Items 1–2 are the soul... Items 3–4 are the identity... The split is
constitutive") was independently checked against the actual file content
(line 137) — it is an accurate, non-fabricated citation.

**Rendered-mirror regeneration — independently reproduced, not just
diffed.** β re-ran the actual renderer locally on the checked-out branch
(`./src/packages/cnos.core/commands/install-wake/cn-install-wake
cds-dispatch`) and confirmed:
- the golden fixture re-render is reported "unchanged" by the tool itself
  and `git status --porcelain` shows a clean tree afterward — the
  committed golden is a genuine renderer output, not a hand-edit;
- `diff .github/workflows/cnos-cds-dispatch.yml
  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`
  → no output (byte-identical);
- `sha256sum` of both files matches:
  `8f6897164af05bcb1fc405490c45c8721c74c5c7f301555b00ca802813e76ff1` —
  identical to the hash α recorded in self-coherence.md, independently
  reproduced rather than trusted.

This mechanically reproduces exactly what
`.github/workflows/install-wake-golden.yml`'s "Re-render cds-dispatch
wake" + "Verify goldens unchanged" + "Verify live cds-dispatch workflow
matches golden (sha256)" steps do — confirmed by reading that workflow
file directly, not inferring its logic. The CI gate would pass on this
branch.

α's disclosed discrepancy (the renderer is a standalone shell script,
`src/packages/cnos.core/commands/install-wake/cn-install-wake`, not a `cn`
Go-binary subcommand, despite the scaffold/α-prompt text saying `cn
install-wake ...`) is independently confirmed accurate: `file` on that
path reports "POSIX shell script"; `install-wake-golden.yml`'s own steps
invoke the same script path α used. This is a scaffold-prose/tooling
naming mismatch, correctly flagged by α rather than silently smoothed
over, and does not affect correctness.

## AC3/AC4 — capability removal + write-fence retirement

**Correctly NOT implemented — confirmed honest STOP, not a silent drop.**

- No checkout-isolation / sparse-checkout mechanism was added. Full
  6-file diff (`git diff --name-only
  86042ec5be4b5fb45b213c27dfcf635958f60aac..HEAD`) is: two `.cdd/unreleased/626/`
  doctrine artifacts, the two SKILL.md prose files, and the two
  regenerated YAML mirrors. No checkout-step region, no new manifest
  field, nothing resembling candidate move (1) or (2).
- Write-fence step independently verified byte-identical to `main`,
  including line numbers (the edit above it is a single-line-for-
  single-line YAML block-scalar swap, so no line-number drift occurred
  either): `diff <(sed -n '355,445p' <base-workflow>) <(sed -n
  '355,445p' <head-workflow>)` → no output. `dispatch_activation_log_write_violation`
  still appears at the same lines (358, 385, 431, 435) in both base and
  head.
- `self-coherence.md`'s "## AC3/AC4 — doctrine + plan, deferred to
  operator" section exists (confirmed present, ~61 lines) and is
  substantive, not a one-line punt: it restates both candidate mechanisms
  from the scaffold (sparse-checkout/worktree isolation; relocating
  `.cn-sigma/` out of the repo) with their actual tradeoffs (breaking
  live-wake codebase/PR/finalizer access; requiring new infra/permissions
  the operator would need to authorize), and names four concrete,
  non-generic follow-on requirements (a regression matrix over every live
  cell type's checkout-path needs; a manifest-surface decision; a
  live-fire validation window; the explicit AC3-before-AC4 ordering
  constraint). This matches the operator's design-first bar treating
  doctrine+plan as a first-class deliverable, not filler.

## AC5 — no regression

**MET at the level verifiable from this diff.**
- Zero Go files changed: `git diff --name-only
  86042ec5be4b5fb45b213c27dfcf635958f60aac..HEAD | grep -E '\.go$'` → no
  output (independently re-run).
- `src/packages/cnos.cds/skills/cds/fsm/transitions.json` untouched
  (independently re-run, no output).
- `check-dispatch-repair-preflight.sh` / `check-dispatch-closeout-integrity.sh`
  do not appear anywhere in the changed-file list (independently re-run,
  no output).
- `install-wake-golden` gate logic independently reproduced locally (see
  AC2 above) and would pass.

## AC6 — no new taxonomy, no Demo 0

**MET.** `git diff 86042ec5be4b5fb45b213c27dfcf635958f60aac..HEAD | grep
-iE "cell_kind|status:[a-z-]+"` returns only prose *mentions* inside the
`.cdd/unreleased/626/` doctrine artifacts (γ's oracle table and α's
self-coherence checklist referring to existing labels/values for
verification purposes) — no new `cell_kind` value, no new `status:*`
label definition, no FSM/transitions.json change, no Demo-0-shaped
artifact anywhere in the diff.

## Guardrail verification (independent)

- **#633 / `transitions.json`:** untouched (`git diff --stat` empty, as
  above). No evidence anywhere in the diff or self-coherence.md that
  α found or claimed a requirement to touch it — correctly left alone.
- **Write-fence step:** byte-identical to `main` at the same line range,
  confirmed above.
- **`.cn-sigma/` contents:** `git diff --stat
  86042ec5be4b5fb45b213c27dfcf635958f60aac..HEAD -- .cn-sigma/` → no
  output. Untouched.
- **`activate/SKILL.md`:** untouched, confirmed above under AC2.

## Findings

None. Every claim in α's self-coherence.md that β independently
re-derived (grep results, diff emptiness, sha256 hashes, renderer
re-invocation, write-fence line-range byte-identity) matched. The one
disclosed discrepancy (shell-script renderer vs. the scaffold's `cn
install-wake` prose) is accurate and appropriately surfaced rather than
hidden — not a finding against α, a correct disclosure.

The bounded/unbounded split itself holds up under independent scrutiny:
move (3) (drop full activation, Kernel+CA-skills only) is the only one of
the issue's three candidate moves achievable without touching live-wake
checkout shape or requiring new infra, and α executed exactly that move
and no more — no over-forcing of AC3/AC4, no silent dropping of them
either (the required doctrine+plan section is real and substantive).

## Verdict

`verdict: converge`

AC1 and AC2 (first disjunct) are correctly implemented and independently
verified byte-for-byte against both rendered mirrors and the live CI
gate's logic. AC3/AC4 are correctly and honestly deferred to the operator
with a substantive doctrine+plan writeup, per the operator's own
design-first STOP condition — this is the designed outcome of a
supervised design-first cell, not a partial failure. AC5/AC6 hold as
pure-regression guardrails (zero Go changes, zero FSM changes, zero new
taxonomy). All named guardrails (#633, write-fence, `.cn-sigma/`,
`activate/SKILL.md`) are confirmed untouched by independent inspection,
not by trusting α's own claims.

---

# beta-review.md — cnos#626 (β, R1 — AC3/AC4 continuation)

## Scope of this review

α's R1 round (AC3 implementation + AC4 deliberate deferral). Independent
re-derivation, not a read of α's self-coherence.md narrative alone.

## AC3 — capability removal

Independently verified, not trusted from α's report:

1. **Renderer scoping.** Read `cn-install-wake` directly. `role` (enum-
   validated `admin|dispatch|observer`) and `agent` (default `"sigma"`)
   are plain script variables in scope, unchanged, at the point the new
   block is emitted. The new block is correctly gated
   `if [ "$role" = "dispatch" ]`.
2. **Rendered both real manifests** (`cds-dispatch`, `agent-admin`)
   independently: dispatch gets the sparse-checkout block with the
   correct patterns (`/*`, `!/.cn-sigma`) + `sparse-checkout-cone-mode:
   false`; admin gets none. Confirmed the admin golden fixture re-renders
   byte-identical to its committed state (diff, not assumption).
3. **Non-vacuity check on the new Go test.** Ran
   `TestDispatchRenderer_SparseCheckoutExcludesAgentHub` — passes.
   Independently sabotaged a scratch copy of the renderer (changed the
   gate to `role = "dispatch-DISABLED"`) and re-ran: the test fails as
   expected. This is real evidence the test detects a regression, not a
   tautological check of its own fixture.
4. **Real-git-mechanism proof, repeated independently.** Cloned the
   actual repo (not the Go test's synthetic fixture) and ran
   `git sparse-checkout set --no-cone '/*' '!/.cn-sigma'`; diffed
   `git ls-files` (full checkout) against the resulting working tree's
   file list. The only difference: all `.cn-sigma/*` entries (43 files)
   absent. Nothing else affected — no unintended exclusion.
5. **Other `role: dispatch` manifests.** Grepped
   `src/packages/*/orchestrators/*/SKILL.md` for `role:`. Only one other
   exists (`test-fixtures/log-writer-misdeclaration/SKILL.md`),
   deliberately mis-declared so the renderer refuses at exit 4 before
   reaching checkout-emission — never renders a workflow, correctly out
   of scope, not a gap.
6. **No leak via another step.** Confirmed only one checkout step exists
   in the rendered workflow; the write-fence's own `git fetch` step only
   updates refs/objects (no working-tree materialization), and the
   mechanical-recovery-scanner step talks to the GitHub API only — neither
   can reintroduce `.cn-sigma` into the working tree.
7. **Full test suite.** `go test ./...` (src/go) plus the three
   standalone module roots — all green, `go vet` clean.

## AC4 — write-fence retirement correctly NOT attempted

Confirmed the `dispatch_activation_log_write_violation` fence step and
its detection logic are byte-identical aside from the line-offset shift
caused by the new block's insertion point. No fence code, message text,
or `if: always()` condition touched. This matches — and correctly
respects — the operator's explicit gate ("retire ONLY AFTER AC3 is
proven in-cell") and R0's own follow-on-cell prerequisite list, which
named a live-fire validation window as a precondition AC4 doesn't yet
have (a workflow-checkout change cannot be observed running for real
before the PR that authors it merges — structurally true, not an
excuse). α's self-coherence.md names this honestly rather than silently
attempting the retirement anyway; that restraint is itself a review point
in α's favor, not a gap.

## Regression-matrix check (R0's prerequisite 1)

Independently re-ran α's grep sweep for `.cn-sigma`/`.cn-{agent}` across
`SKILL.md` files under `cnos.cdd`/`cnos.cds`/`cnos.core`. Confirms: zero
references in `gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`,
`cds/SKILL.md`, `issue/SKILL.md` — the actual set of skills cell-execution
cognition loads per `delta/SKILL.md` §9.12. The only references are in
`delta/SKILL.md` §9.12 itself (the doctrine statement), the wake prompt
(`cds-dispatch/SKILL.md`), the admin wake manifest, and agent-identity-
layer skills cells never load. This is real evidence against the failure
mode R0 flagged, not an assumption.

## Findings

None. Zero blocking or non-blocking findings survived independent
re-derivation.

## Verdict

**`verdict: converge`.** AC3 is correctly and completely implemented and
evidenced for this bounded round; AC4 is correctly and honestly deferred,
consistent with both the operator's explicit gate and R0's own stated
prerequisites. Recommend the closeout carry forward R0's own
recommendation: treat the next 1-2 real `cds-dispatch` firings post-merge
as the live-fire validation window before any future cell attempts AC4.
