# self-coherence — cnos#497

cycle: 497
role: alpha
round: R0

---

## §R0 — Decision artifact

### Verdict: Model B — `.cdd/` is the CDD framework's receipt ledger

The project adopts **Model B**. `.cdd/` is not a package namespace for cnos.cdd; it is the CDD framework's receipt ledger — the unified surface where the cell-runtime framework records evidence regardless of which concrete protocol authored the cell. **Canonical protocol identity is the typed receipt's `protocol_id` field** (for CDS: `cnos.cdd.cds.receipt.v1`). Issue labels, `gamma-scaffold.md` metadata, and PR body metadata are operational supporting surfaces, not the receipt's canonical protocol discriminator. The path identifies the ledger owner (`.cdd/` = CDD framework); the typed `protocol_id` identifies the concrete protocol (`cnos.cdd.cds.receipt.v1` = CDS-shaped receipt).

**This issue closes with this decision artifact. 497B and 497C never file.**

---

### Structural arguments

**Argument 1 — The ledger-vs-namespace distinction.**

`.cdd/` names the CDD framework's *ledger surface*, not a namespace for cnos.cdd-package-owned artifacts. The analogy holds: `.git/` is git's metadata directory for all commits and branches, not a namespace owned by one branch. `.cdd/` is the receipt ledger for all CDD cells regardless of which concrete protocol produced them. **Canonical protocol identity lives in the receipt's typed `protocol_id` field** (per the generic CDD receipt schema; CDS receipts pin to `cnos.cdd.cds.receipt.v1`), which the Go verifier uses to dispatch protocol-specific validation. Issue labels, `gamma-scaffold.md` metadata, and PR metadata are useful operational corroboration but not the receipt's canonical discriminator. Splitting the path does not add information the typed receipt doesn't already carry; it adds coordination cost for no observability gain.

**Argument 2 — cnos#496's write-locality principle is orthogonal to the ledger question.**

cnos#496 established that agent identity logs (`.cn-sigma/logs/`) are write-local to the admin wake. The dispatch wake does NOT write to `.cn-sigma/`; it DOES write to `.cdd/unreleased/{N}/`. Write-locality applies to *agent identity surfaces*; it does not imply the *artifact ledger* must be split per protocol. The sibling framing in the issue is correct: "package-isolation for *writes* (not for the ledger) is settled" (cnos#496). This issue answers the ledger question *independently* — and the ledger question has a clean answer: the ledger is framework-owned, not protocol-namespace-split.

**Argument 3 — Zero migration cost; all current tooling unchanged.**

Model B preserves the current behavior exactly. The release skill, `cn-cdd-verify`, dispatch-protocol skill, γ/α/β/δ skills, and the cds-dispatch `wake-provider.json` (`output_contract.cycle_artifact_root: .cdd/unreleased/{N}/`) all remain unchanged. The discomfort that motivated this issue — "a `protocol:cds` cell landing artifacts at `.cdd/unreleased/491/` reads as namespace conflation" — is a naming-mismatch-without-substance. The cell already carries `protocol:cds` as its routing label; the directory path does not need to repeat that information.

**Argument 4 — Wave-master single-root discovery (constant root count).**

A wave-master tracker (e.g. cnos#467 or cnos#495) finds all child receipts with a single walk over `.cdd/unreleased/`. The complexity of receipt enumeration itself is `O(R)` (proportional to the number of receipts) under either model; this argument is about *root-count* complexity, not receipt-count complexity. Under Model B: **one stable root**, constant root-count regardless of how many concrete protocols exist. Under Model A: `O(P + R)` (P = number of protocol roots; R = number of receipts) — every wave-master, every release pass, every verifier invocation has to enumerate roots before walking receipts. As the number of concrete protocols grows, the root-enumeration overhead grows with it. Model B keeps discovery simpler and tooling more robust: one root to walk; no protocol-root mapping table; no per-root enumeration. (The earlier R0 wording "O(1) discovery" conflated root-count complexity with receipt-count complexity; corrected to "single-root discovery / constant-root-count" per operator-final-read R1.)

**Argument 5 — Mechanical-orchestration trajectory simplicity.**

When `cn dispatch claim` becomes Go (per the mechanical-vs-LLM partition trajectory cnos#495 established), it writes to `.cdd/unreleased/{N}/` always — a constant path independent of the claimed cell's protocol label. No derived-path computation, no protocol-to-directory mapping table, no renderer changes needed. Model A would require `cn dispatch claim` to either (a) derive the path from the protocol label at claim time or (b) carry a per-protocol path table; both add complexity for no structural gain.

---

### Answers to the 7 open questions

**Q1 — What owns the receipt ledger?**

`.cdd/` is the CDD framework's receipt ledger. The cnos.cdd framework package owns the ledger surface. Concrete protocol packages (cnos.cds, cnos.cdr, future cnos.cdw) write receipt evidence *into* the framework's ledger — they are writers, not owners, of the ledger directory. The distinction: the framework *defines* the receipt contract and *hosts* the ledger; the concrete protocols *produce* receipts that conform to that contract. This is not a namespace conflict; it is a tiered ownership model: framework owns the ledger structure, protocols own the content they contribute. **Canonical protocol identity lives in the receipt's typed `protocol_id` field**, not in the directory path. For CDS receipts: `protocol_id: cnos.cdd.cds.receipt.v1`. The Go verifier dispatches protocol-specific validation by `protocol_id`. Path = ledger owner; `protocol_id` = concrete protocol.

**Q2 — What does release do with `.cds/unreleased/`?**

`.cds/unreleased/` does not exist under Model B. The release skill reads from `.cdd/unreleased/{N}/` exactly as today. No behavior change. The release skill lifts cells from `.cdd/unreleased/{N}/` → `.cdd/releases/{X.Y.Z}/{N}/` (versioned) or `.cdd/releases/docs/{ISO-date}/{N}/` (docs-only disconnects) regardless of which protocol authored the cell.

**Q3 — Does `.cdd/releases/` still exist?**

Yes, `.cdd/releases/` remains as-is. The versioned release tree (`.cdd/releases/{X.Y.Z}/{N}/`) and the docs-only disconnect destination (`.cdd/releases/docs/{ISO-date}/{N}/`) are unchanged under Model B. No split, no per-protocol subdirectory.

**Q4 — Where do docs-only disconnects move?**

Docs-only disconnects remain at `.cdd/releases/docs/{ISO-date}/{N}/`. No change. Protocol identity is carried inside the receipt at that path.

**Q5 — What does `cn-cdd-verify` verify?**

`cn-cdd-verify` continues inspecting `.cdd/` as a unified ledger. No rename (no `cn-cds-verify`); no per-package verifier; no multi-tree enumeration. The verifier reads the receipt at `.cdd/unreleased/{N}/` or `.cdd/releases/.../{N}/` and validates against the framework's receipt contract. The protocol qualifier inside the receipt (if the verifier needs to route protocol-specific checks) is a runtime value, not a directory-path variable.

**Q6 — How does a package-agnostic parent find child receipts?**

Under Model B: single-root walk over `.cdd/unreleased/` listing all `{N}/` directories. To find *which* cells belong to a given wave, the wave-master's cross-reference list in the issue is the primary source; the `.cdd/unreleased/{N}/gamma-scaffold.md` file's `cycle:` field is the secondary confirmation; the receipt's typed `protocol_id` is the canonical protocol-routing discriminator (see Q1). Single root → no protocol-root enumeration; receipt-count complexity is `O(R)` either way. Under Model A: `O(P + R)` (P protocol-root enumerations before receipt walks). The advantage of Model B is **constant root count**, not constant-time receipt discovery. This is already how wave-master tracking works today (cnos#467 listed sub-issues by number; the cells under `.cdd/unreleased/{485,486,487,491,496}/` are findable by number).

**Q7 — How does this interact with the mechanical-orchestration trajectory?**

`cn dispatch claim` (when implemented in Go) writes to `.cdd/unreleased/{N}/` always. The path is a framework constant, not derived from the claimed issue's protocol label. This is simpler and more robust than Model A's derived-path logic. The protocol label (`protocol:cds`) is used for dispatch routing (which wake claims the cell) but not for determining the artifact path (which is framework-governed). The two concerns are orthogonal.

---

### Cross-references to affected surfaces

Under Model B, **no surfaces are changed** by this decision. The decision confirms the current implicit model as the explicit declared model.

Named for completeness (surfaces that would have changed under Model A, and are confirmed stable under Model B):
- `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` — reads `.cdd/unreleased/`; unchanged
- `src/packages/cnos.cdd/commands/cdd-verify/` (operator invocation: `cn cdd verify`) — inspects `.cdd/` as unified ledger; validates typed receipts by `protocol_id`; unchanged
- `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` — references `.cdd/unreleased/` via the output contract; unchanged
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.5 per-R[N] artifact contract — references `.cdd/unreleased/{N}/`; unchanged
- `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` `output_contract.cycle_artifact_root` = `.cdd/unreleased/{N}/`; unchanged
- All γ/α/β/δ skills referencing `.cdd/unreleased/` — unchanged

---

### Follow-up issues

**497B and 497C never file.** Model B is adopted; the original "package-namespace rename" proposal is retracted. No implementation issues follow from this decision.

If future observations surface a structural reason to revisit this decision (e.g. protocol-specific ledger isolation becomes necessary at scale; a concrete protocol's cells need a separate verifier due to schema divergence), that should be filed as a new issue with that evidence.

---

## AC verification

| AC | Result | Evidence |
|---|---|---|
| AC1: Decision artifact exists | PASS | This file at `.cdd/unreleased/497/self-coherence.md` |
| AC2: All 7 open questions answered | PASS | Q1–Q7 answered above under "Answers to the 7 open questions" |
| AC3: Model chosen and reasoned | PASS | Verdict: Model B; 5 structural arguments given |
| AC4: Affected surfaces named | PASS | "No surfaces changed" declared; affected surfaces listed and confirmed stable |
| AC5: Follow-up issues identified | PASS | "497B and 497C never file" declared |
| AC6: No implementation work in this cycle | PASS | No diff outside `.cdd/unreleased/497/` |

---

## CDD Trace

**Gap:** The project used `.cdd/` implicitly as the CDD framework's receipt ledger but had never explicitly decided whether this was Model A (package-namespace) or Model B (CDD receipt ledger). The discomfort from cycle/487 surfaced this as a naming-mismatch-without-substance that needed a formal decision before any implementation work.

**Mode:** design / decision (docs-only; the artifact IS the deliverable)

**Action taken:** Analyzed both models against the 7 open questions; chose Model B on structural grounds (ledger-vs-namespace distinction; write-locality orthogonality; zero migration cost; wave-master O(1) resolution; mechanical-orchestration simplicity).

**Review-readiness signal:** READY — all 6 ACs pass; no implementation work in diff; decision fully reasoned.

---

## §R1 — Operator-final-read absorbed (δ-direct; T-486-7 pattern)

**R0 verdict:** β converge (`beta-review.md §R0`).

**Operator-final-read on PR #499 (T-486-12 P1 defense-in-depth):** iterate narrowly. Architectural verdict (Model B) unchanged; six precision/closure issues found that β R0 did not surface. Absorbed inline as δ-direct R1 (no β re-spawn; the corrections are narrow mechanical text fixes; pattern's 4th empirical sample after cycle/486 R1, cycle/496 R1, cycle/496 R2).

### R1 corrections applied

| # | Issue | Correction |
|---|---|---|
| 1 | Argument 4 + Q6 claimed "Model B keeps discovery O(1)" — conflates root-count complexity with receipt-count complexity | Replaced with "single-root discovery / constant-root-count discovery"; complexity stated honestly as `O(R)` over receipts under either model + `O(P + R)` under Model A due to protocol-root enumeration overhead |
| 2 | Verdict + Argument 1 + Q1 framed protocol identity as "issue labels + gamma-scaffold metadata + PR body" — not the canonical typed identity | Anchored canonical protocol identity to the typed receipt's `protocol_id` field; for CDS: `cnos.cdd.cds.receipt.v1`; issue labels / gamma-scaffold metadata / PR metadata named as operational supporting surfaces, not canonical discriminators; the Go verifier dispatches validation by `protocol_id` |
| 3 | Cross-references cited stale/vague canonical paths: `src/packages/cnos.cdd/skills/release/SKILL.md` (wrong; missing the `cdd/` subdirectory) and `cn-cdd-verify (wherever it lives)` (vague) | Corrected to `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` and `src/packages/cnos.cdd/commands/cdd-verify/` (operator invocation: `cn cdd verify`) |
| 4 | `gamma-closeout.md` said "Cycle 497 is closed" while issue is `status:review` — conflates cell closure with boundary acceptance | Rewrote to "Cell 497 is closed and awaiting operator boundary decision. β verdict: converge. Issue state: status:review. Cycle acceptance occurs when the operator merges PR #499 and closes cnos#497." (CDS doctrine separates cell closure from boundary acceptance) |
| 5 | Actor-collapse / configuration-floor not explicitly declared in `gamma-closeout.md` (γ + α + β collapsed on δ was acknowledged in scaffold but not in closeout) | Added "Dispatch configuration" section to `gamma-closeout.md` declaring: collapsed mode is permitted for docs/design-decision cells; configuration floor — cycle does not claim fully independent α/β actor separation; success claim calibrated — this cycle validates live wake claiming + issue lifecycle + artifact production + PR creation + operator-final-read, but does NOT newly prove independent γ/α/β execution contexts |
| 6 | γ process-gap audit table had "mental notes" without canonical dispositions — γ doctrine requires explicit triage of every finding | Reformulated table with explicit Type + Disposition + Reason columns; both findings classified as `no-patch` (with reasons); declared `protocol_gap_count: 0`; declared `cdd-iteration.md not required` |

### R1 verdict

The decision artifact's substance (Model B; 5 structural arguments; Q1–Q7 answered) survives intact. The corrections are precision (complexity claim), canonical-identity anchoring, citation accuracy, doctrinal closure-vocabulary, configuration-floor honesty, and finding-disposition explicitness. None of these affect the verdict. **R1 is the absorption commit; the architectural decision is unchanged.**

### Operating-model observation

This is the **first live-wake-claimed cycle** under the corrected mouthpiece model (per cnos#496 γ-closeout §2.2 operating-model correction). The cycle executed end-to-end autonomously: live wake claimed → δ wake-invoked mode → γ/α/β collapsed on δ (permitted for design-only cells) → artifacts landed → PR opened → status:review. Operator-final-read on PR caught what β's R0 walk methodologically could not — same pattern that drove T-486-12's P1 promotion. **R1 is δ-direct because the corrections are narrow mechanical text fixes; the live wake has no autonomous mechanism to re-engage from `status:review` so γ-interface absorbs the iterate inline.**

### T-486-7 pattern sample-size = 4

| Cycle | δ-direct R[N] | Operator iterate trigger |
|---|---|---|
| cycle/486 | R1 | δ §9.5 off-by-one |
| cycle/496 | R1 | HEAD@{1} multi-commit baseline escape |
| cycle/496 | R2 | --diff-filter=ACMR delete escape |
| **cycle/497** | **R1** | 6 precision/closure issues (O(1) miswording, protocol_id anchoring, stale paths, closure wording, configuration-floor, finding dispositions) |

The pattern is doctrinally robust: when an operator-final-read iterate is narrow-mechanical-correctness with unambiguous path forward, δ-direct absorbs efficiently without α/β respawn.
