---
cycle: 500
role: alpha
version: R0
mode: design + build
---

# Self-coherence — cnos#500 — cdd/review-return

## §R0 Design decisions

### Command shape

`cn cell return` and `cn cell resume` — registered as `cell-return` and `cell-resume` in the kernel registry using the existing noun-verb dispatch pattern (`dispatch.go` `ResolveCommand` §1 lookup). User-facing: `cn cell return --issue N --verdict V --review path`, `cn cell resume --issue N`.

Rationale: the `cn cell` namespace is explicitly pinned in the implementation contract; the noun-verb pattern is the canonical form in this codebase (`dispatch.go` `GroupMembers`/`PrintGroup`); `return` and `resume` are the most semantically precise verbs for the two operations (return = deliver verdict; resume = re-arm the existing cycle).

### HI contract placement

New file `src/packages/cnos.core/orchestrators/agent-admin/hi-contract.md`. The existing `prompt.md` is the functional wake prompt for the agent-admin wake — it is runtime-loaded by the substrate and must not carry doctrinal contract prose that is not wake-specific. A separate `hi-contract.md` is the correct placement for the HI behavioral contract (which governs all HI sessions, not just the agent-admin wake).

### AC6 enforcement mechanism

Convention-only enforcement with β Rule 7 as the backstop. CI grep of HI authorship signatures in role-owned paths is not feasible: the HI commits do not carry a stable machine-distinguishable authorship signature (the HI may be dispatched under any git identity; the only stable signal is the artifact path + content, which requires semantic parsing). The convention is: prohibited surfaces are named in `hi-contract.md`; β's Rule 7 diff-vs-contract verification catches violations at review time. This is the same enforcement model as the implementation-contract axes. Named as known debt in §Debt.

### δ SKILL.md amendment structure

New `§9.10` subsection in `delta/SKILL.md` for the `resumed-from-changes` wake-invoked mode shape. Existing §9.1–§9.9 preserved intact. The §9.6 carve-out ("status:changes is EXTERNAL") is amended to read: `cn cell return` acts on behalf of the operator's stated verdict — the operator is still the authority; the command is the mechanical translator. This preserves the carve-out's intent while installing the mechanism.

### Cycle sizing judgment

Keep as one cell. Five-factor check:
1. One new Go package (`src/go/internal/cell/`) — bounded; no changes to existing packages
2. Three existing surfaces touched (delta SKILL.md, agent-admin, dispatch-protocol) — all additive amendments
3. Spans design + code + docs — explicitly the declared mode ("design + build")
4. Design pass required and completed above
5. Operator suggested split as optional ("start with the live worker; γ can make that call") — γ kept whole; all 7 ACs are tightly coupled (the HI contract, the schema, the CLI command, and the δ amendment are a single conceptual unit that splits poorly: AC2+AC6+AC7 require the same schema, AC3+AC4 require the same CLI command, AC5 requires understanding both)

---

## §Gap

**Issue:** cnos#500 — cdd/review-return: mechanically route operator iterate verdicts back into an existing cell

**Version/mode:** R0 / design + build

**The gap:** The live CDS dispatch wake handles a cell through `status:review` with no autonomous mechanism to re-engage from `status:review` when the operator's final read returns `iterate-narrowly`. Without this primitive, the HI absorbs operator-final-read corrections inline — crossing role boundaries. This was observed empirically in cycle/497: the HI edited α-owned `self-coherence.md`, rewrote γ-owned `gamma-closeout.md`, and inserted §R1 sections into β/α/γ artifacts, framing the absorption as "δ-direct R1." That absorption was declared a `degraded_recovery: human_interface_applied_operator_patch` in `.cdd/unreleased/497/gamma-closeout.md §5`.

**This cycle installs:**
1. A typed `cn.operator-review.v1` schema for HI to capture operator verdicts
2. `cn cell return --issue N --verdict V --review path` — mechanical `status:review → status:changes` transition
3. `cn cell resume --issue N` — re-arm the existing cycle (branch + artifact directory preserved; R[N+1] appended)
4. δ SKILL.md §9 amendment for the `resumed-from-changes` wake-invoked mode shape
5. HI behavioral contract — explicit prohibited surfaces + MUST NOT language
6. `degraded_recovery` declaration schema as first-class convention

---

## §Skills

**Tier 1 — CDD lifecycle:**
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (this role)
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` (amended §9; §9.6 carve-out; §9.10 new subsection)

**Tier 2 — always-applicable eng:**
- `src/packages/cnos.core/skills/write/SKILL.md` (governing question, one-file-one-job, front-load the point)

**Tier 3 — issue-specific:**
- `src/packages/cnos.handoff/skills/handoff/dispatch/SKILL.md` (implementation contract; 7-axis schema; four-surface mesh)

**Active constraints from loaded skills:**
- write/SKILL.md §3.2 — one file, one governing question (applied to hi-contract.md and operator-review/SKILL.md)
- write/SKILL.md §3.4 — front-load the point (applied to all authored Markdown)
- alpha/SKILL.md §2.5 — incremental self-coherence write discipline (one section per commit)
- alpha/SKILL.md §2.3 — peer enumeration before closure claims
- alpha/SKILL.md §3.6 — implementation contract is δ's; α MUST NOT improvise
