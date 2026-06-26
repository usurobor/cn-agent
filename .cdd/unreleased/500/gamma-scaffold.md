# γ-scaffold — cycle/500 — cnos#500

**Mode:** design + build
**Cycle branch:** cycle/500
**Main SHA at scaffold time:** 3095fa2b44145490c8e5241bd347165a53ace827
**Protocol:** cds

---

## Source-of-truth table

Derived from issue body §"Source of truth" and related-doctrine references.

| Claim / surface | Canonical path | Status |
|---|---|---|
| Agent-admin contract: HI never executes cells; cell-shaped work defers to dispatch wake | `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` | Shipped; respected for outbound direction only |
| δ contract; named failure mode "invisible meddling" | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` | Shipped; violated by cycle/497 "δ-direct" framing |
| Live wake claim sequence + lifecycle transitions; does NOT cover operator-iterate re-entry | `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` | Shipped; gap confirmed: no `resumed-from-changes` shape |
| `status:review` → `status:changes` external-rejection semantics; no mechanical implementation that resumes the same cell | `src/packages/cnos.core/labels.json` + dispatch-protocol skill | Shipped doctrine; missing mechanical implementation |
| cycle/497 R1 empirical witness for the gap | `.cdd/unreleased/497/operator-review.md` + cycle/497 commit `dd819f00` | Live evidence on main |
| δ wake-invoked mode §9 — current contract (covers outbound through `status:review`; carves out `status:changes` as external) | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9` | Shipped; §9.6 explicitly names `status:changes` as external/out-of-scope |
| Existing `cn` binary CLI structure | `src/go/internal/cli/` | Shipped; `cmd_*.go` pattern; no `cmd_cell.go` exists |
| `src/go/internal/cell/` package | (does not exist) | Gap — new package required |
| Operator-review typed schema | `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md` | Gap — new skill required (this cycle's AC1) |
| HI behavioral contract (explicit prohibited surfaces) | `src/packages/cnos.core/orchestrators/agent-admin/` | Shipped; needs amendment to list prohibited role-owned surfaces explicitly (this cycle's AC6) |
| `degraded_recovery` declaration schema | (does not exist) | Gap — new schema/convention required (this cycle's AC7) |

---

## AC oracle approach

### AC1 — Typed `cn.operator-review.v1` schema

**Evidence α must produce:**
- `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md` exists and defines required fields: `schema`, `issue`, `pr`, `verdict` (enum: `converge`/`iterate`/`reject`/`clarify`), `reviewer`, `captured_by`, `captured_at`, `findings[]` (each with `id`, `surface`, `problem`, `expected_change`).
- An example artifact in the skill or adjacent fixture that conforms to the schema.
- Confirm that `.cdd/unreleased/497/operator-review.md` (the empirical motivator) conforms to the schema as stated.

**How β independently verifies:**
- Read the SKILL.md schema definition; independently confirm field list completeness.
- Check that the 497 `operator-review.md` frontmatter matches every required field.
- No oracle tool exists yet — β applies structural schema review: required fields present, no missing field, enum values cover the issue's four verdict types.

**Pass:** `operator-review/SKILL.md` exists with all required fields listed, example conforms, 497 artifact conforms. No required field absent.

---

### AC2 — HI records verdict without editing role-owned matter

**Evidence α must produce:**
- HI behavioral contract documented — either as an amendment to `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` or as a new HI-contract section.
- Named failure mode "invisible meddling" from `delta/SKILL.md` explicitly cited in the contract.
- HI artifact attribution (`captured_by: gamma-interface (HI)`) distinguished from α/β/γ artifact attribution.
- Prohibited surfaces explicitly listed (not implied): `self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, `gamma-scaffold.md`.

**How β independently verifies:**
- Read the agent-admin contract diff; confirm the prohibited-surfaces list is explicit and not buried in prose.
- Check that the delta SKILL.md "invisible meddling" failure mode is cited by path + anchor, not paraphrased.
- Negative check: agent-admin prompt must NOT say "HI should avoid" — must say "HI MUST NOT author or amend."

**Pass:** explicit MUST NOT language; named prohibited surfaces list; delta SKILL.md citation present; `captured_by` attribution distinguished from role artifacts.

---

### AC3 — Mechanical `status:review` → `status:changes` transition command

**Evidence α must produce:**
- `cn cell return` command (or equivalent per design decision) implemented in Go at `src/go/internal/cli/cmd_cell.go` + `src/go/internal/cell/` package.
- On `--verdict iterate` or `--verdict reject`: removes `status:review` label and adds `status:changes` label on the GitHub issue.
- Smoke fixture (test or documented manual invocation sequence) that demonstrates the transition on a cell at `status:review`.

**How β independently verifies:**
- Read the Go implementation diff; confirm label mutation calls are correct (`status:review` removal + `status:changes` addition).
- Run or read the smoke fixture; confirm it passes.
- Check that `converge` verdict does NOT apply the transition (only `iterate`/`reject` do).

**Pass:** label transition is mechanical (not manual), command compiles, smoke fixture demonstrates the correct before/after label state.

**Note:** This AC has the highest mechanical clarity of all 7 — the oracle is a label-state diff before/after the command. Clear pass/fail.

---

### AC4 — Existing cycle resumed (not replaced)

**Evidence α must produce:**
- `cn cell resume --issue {N}` (or equivalent) does not create a new cycle branch or new `.cdd/unreleased/{N}/` directory.
- Existing artifacts on `cycle/{N}` and in `.cdd/unreleased/{N}/` persist unchanged.
- R[N+1] section is appended to the appropriate artifact (e.g., `self-coherence.md §R[N+1]`) rather than replacing it.
- δ SKILL.md §9 amendment names this as the `resumed-from-changes` shape with R[N+1] increment.

**How β independently verifies:**
- Smoke test or documented verification: before/after artifact state on a test cycle shows no files deleted, branch unchanged.
- Read the Go implementation: confirm no directory-creation or branch-creation call in the resume path.
- Read the δ SKILL.md §9 amendment; confirm `resumed-from-changes` shape is defined with R[N+1] increment documented.

**Pass:** no new directory; branch unchanged; existing artifacts present; R[N+1] increment documented in δ skill amendment.

---

### AC5 — δ routes R[N+1] in order

**Evidence α must produce:**
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9` amended (or extended) to define the `resumed-from-changes` wake-invoked mode shape.
- Amendment names: (1) input includes operator-review artifact + prior round's artifacts + R[N+1] increment; (2) routing sequence: α R[N+1] → β R[N+1] → γ closeout amendment; (3) the resumed-cell contract distinctions from first-claim contract.

**How β independently verifies:**
- Read the δ SKILL.md diff; confirm the `resumed-from-changes` shape is a named variant (not buried in prose) with the three routing-sequence steps explicit.
- Confirm the amendment does not conflict with existing §9.3 routing sequence for first-claim cells.
- Confirm the carve-out that `status:changes` is external (§9.6) is preserved or explicitly superseded by this amendment.

**Pass:** `resumed-from-changes` shape named; R[N+1] routing sequence explicit; no contradictions with existing §9 contract.

**Design-pass judgment required:** α must decide whether to add a new §9.N subsection or amend existing sections. The structure is α's design decision; the content is specified by the AC.

---

### AC6 — HI cannot substitute for dispatched roles

**Evidence α must produce:**
- HI contract document (per AC2 surface) explicitly lists prohibited artifact surfaces.
- The list is exhaustive: `self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, `gamma-scaffold.md`.
- CI guard OR convention enforcement named: either a CI grep check for HI authorship signatures in role-owned paths, OR a documented convention that β's Rule 7 catches violations at review time (TBD per design pass).

**How β independently verifies:**
- Read the contract document; confirm each prohibited surface is named.
- Confirm the enforcement mechanism (CI guard or convention) is either implemented or named as a named non-goal with justification.
- Negative: contract must not say "HI should generally avoid" — it must use MUST NOT language.

**Pass:** all six surfaces named; MUST NOT language; enforcement mechanism present or explicitly named as future-TBD with reason.

**Design-pass judgment required:** whether CI enforcement is feasible in this cycle. The issue body says "(TBD per design pass)" — α must decide and document.

---

### AC7 — Declared bootstrap-exception escape hatch

**Evidence α must produce:**
- A `degraded_recovery` declaration schema or convention exists (YAML frontmatter convention, CUE schema, or Markdown section template).
- Required fields in the declaration: `reason`, `status` (at minimum), plus `recovery_type: human_interface_applied_operator_patch`.
- cycle/497's `gamma-closeout.md` carries the declaration as the first canonical witness (either α verifies this is already present from the cycle/497 recovery, or α adds it as part of this cycle's scope — issue body is ambiguous on who adds it).
- A grep/structured check can detect the declaration's presence.

**How β independently verifies:**
- `grep -rn "degraded_recovery" .cdd/` returns at least the 497 gamma-closeout hit.
- Read the schema/convention document; confirm required fields are present.
- Confirm detection command is documented (issue says "detectable by a grep / structured check").

**Pass:** schema exists; 497 gamma-closeout carries the declaration; grep detects it.

**Friction note:** The issue body says cycle/497's gamma-closeout should carry the declaration "post-this-issue's R1 absorption pass" — but cycle/497 is already closed (or nearly so). α must determine the current state of `.cdd/unreleased/497/gamma-closeout.md` and decide whether this AC requires an amendment to 497 artifacts or only the schema definition + 497 as the stated first witness (which may already satisfy the oracle if the 497 closeout already has the declaration language).

---

## Scope guardrails

**In scope:**
- `cn.operator-review.v1` typed schema skill at `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md`
- `cn cell return --issue {N} --verdict {V} --review {path}` and `cn cell resume --issue {N}` CLI commands in `src/go/internal/cli/cmd_cell.go` + `src/go/internal/cell/` package
- `status:review` → `status:changes` mechanical label transition on `iterate`/`reject` verdict
- δ SKILL.md §9 amendment for `resumed-from-changes` shape
- HI behavioral contract amendment in `src/packages/cnos.core/orchestrators/agent-admin/`
- `degraded_recovery` declaration schema/convention
- Smoke fixture for AC3/AC4

**Out of scope / non-goals:**
- Full Go-side runtime ownership of routing + round increment (broader mechanical-orchestration migration)
- Replacing claude-code-action as substrate carrier (cnos#452)
- Revisiting actor-collapse rule semantics (CDS doctrine; this cycle uses the rule, not redefines it)
- Per-protocol review-return variants (this is the generic CDD primitive)
- Notification / dispatch-summary implementation (cnos#495 Sub 2)
- Private-repo policy work

---

## Friction notes

1. **Design-pass judgment on command shape.** The issue explicitly marks exact flag names and subcommand names as α's design decision. α must resolve `cn cell return` vs alternative forms and document the decision in `self-coherence.md`. β must not flag the design decision itself as a finding — only flag deviations from the chosen shape.

2. **AC5/AC7 cross-dependency with δ SKILL.md.** AC5 requires amending `delta/SKILL.md §9`. That skill file is ~500 lines and has strict section-ordering (§9.1 through §9.8). α must use the large-file authoring rule (section-by-section; section-manifest). The amendment must land in the correct subsection without disrupting the existing §9 contract.

3. **AC7 / cycle/497 closeout state.** The issue body says 497's gamma-closeout should carry the `degraded_recovery` declaration. At scaffold time, `.cdd/unreleased/497/gamma-closeout.md` is present on `main`. α must check whether the declaration is already present (the 497 recovery sequence was still in progress at issue-filing time). If already present, AC7 only requires schema definition. If absent, α must patch 497's gamma-closeout — which is a write to `main`, not to `cycle/500`. This cross-branch write requires explicit operator clarification or a design decision: does AC7 require patching 497, or does it only require the schema + 497 as the named first-witness?

4. **HI contract placement.** The issue names `src/packages/cnos.core/orchestrators/agent-admin/` but the existing prompt.md there is a functional prompt, not a behavioral contract document. α's design pass must decide: amend the existing `prompt.md` (adding a MUST NOT section), or create a new `hi-contract.md` or equivalent. Either is acceptable; the decision must be documented.

5. **AC6 CI enforcement ambiguity.** The issue marks enforcement as "TBD per design pass." A grep-based CI guard checking for HI authorship signatures in role-owned artifact paths is feasible but requires knowing what signatures to grep. α must either implement a concrete guard or explicitly name "convention-only enforcement with β Rule 7 as the backstop" and document this as a known gap.

6. **Five-factor cycle sizing.** This cycle introduces (a) one new Go package (`src/go/internal/cell/`), (b) touches 3+ existing surfaces (delta SKILL.md, agent-admin, dispatch-protocol), (c) spans design + code + docs, (d) requires a design pass before implementation. The issue body already carried a "split into 500A/500B/500C" suggestion from the operator. α's design pass should explicitly revisit this and produce a written split-or-keep justification. If α decides to keep whole, the justification must cite the five-factor check.

7. **`status:changes` is currently an external transition only.** δ SKILL.md §9.6 explicitly carves out `status:changes` as external (operator/planner authority). This cycle installs the mechanism for the runtime to *initiate* the `status:review → status:changes` transition (AC3). The δ SKILL.md amendment (AC5) must reconcile this: the amendment should clarify that the `cn cell return` command acts *on behalf of the operator's stated verdict* — the operator is still the authority; the command is the mechanical translator of that authority. Failing to make this distinction would conflict with the existing §9.6 carve-out.

---

## Dispatch prompts

### α prompt

You are α. Project: cnos (protocol:cds).
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 500 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/500
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, src/packages/cnos.handoff/skills/handoff/dispatch/SKILL.md

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go (for `cn cell return/resume` CLI command) + Markdown (for operator-review schema skill, HI behavioral contract, δ-skill amendment, bootstrap-exception schema) |
| CLI integration target | New `cn cell` subcommand family: `cn cell return --issue {N} --verdict {V} --review {path}` and `cn cell resume --issue {N}`; exact flags and subcommand names are α's design decision during the design pass, constrained to the `cn cell {subcommand}` namespace |
| Package scoping | Go: `src/go/internal/cli/cmd_cell.go` (new command file) + `src/go/internal/cell/` (new package for cell operations); Markdown: `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md` (new schema skill) + amendment to `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9` (resumed-from-changes shape) + HI behavioral contract in `src/packages/cnos.core/orchestrators/agent-admin/` |
| Existing-binary disposition | Preserve `cn` binary; `cn cell` is a new additive subcommand group; no existing `cn` commands modified |
| Runtime dependencies | None beyond existing (`gh` CLI is already an assumed runtime dep; no new language runtimes) |
| JSON/wire contract preservation | N/A — all artifacts are new; existing `status:review → status:changes` label transition follows existing `cnos.core/labels.json` label doctrine unchanged |
| Backward-compat invariant | All existing cell lifecycle transitions (status:todo → status:in-progress → status:review → status:changes → status:todo) preserved; `status:changes` semantics unchanged; `cn cell return/resume` is additive; existing δ wake-invoked-mode contract extended not replaced |

### β prompt

You are β. Project: cnos (protocol:cds).
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 500 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/500
