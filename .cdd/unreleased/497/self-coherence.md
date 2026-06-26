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

**Action taken:** Analyzed both models against the 7 open questions; chose Model B on structural grounds (ledger-vs-namespace distinction; write-locality orthogonality; zero migration cost; wave-master single-root / constant-root-count discovery; mechanical-orchestration simplicity).

**Review-readiness signal:** READY — all 6 ACs pass; no implementation work in diff; decision fully reasoned.

---

## §R1 — α takes ownership of operator-finding adoption

cycle: 497
role: alpha
round: R1
input: `.cdd/unreleased/497/operator-review.md` (HI's typed translation of the human-operator final-read; 6 findings O1–O6)
patch_proposal_under_review: `dd819f00` (HI-authored textual edits to §R0; reframed by operator as operator-supplied patch proposal, not α's work until α inspects and takes ownership)

### Framing

The human operator's final-read returned `iterate (narrowly)`. The verdict left Model B intact and named six findings outside β R0's mechanically-scoped AC oracle. The γ-interface, lacking a mechanical review-return primitive (cnos#500 records the missing `cn cell return` / `cn cell resume`), absorbed the corrections inline as commit `dd819f00` and framed the findings as "narrow mechanical text fixes." That framing was wrong. The operator's ruling reframes `dd819f00` as an **operator-supplied patch proposal** and dispatches a proper role pass. This §R1 is α's substantive engagement with the findings.

The corrections are not mechanical. Each one required semantic judgment of a different class:

- O1 is an asymptotic-complexity correctness call (the original text confused root-count with receipt-count complexity)
- O2 is an architectural anchoring decision (where does canonical protocol identity live in the typed-runtime trajectory)
- O3 is citation accuracy against the live canonical filesystem (semi-mechanical but doctrinally important — a decision artifact must not carry knowingly-uncertain canonical paths)
- O4 is CDD doctrine application (cell closure ≠ boundary acceptance per δ rule)
- O5 is CDS actor-collapse-rule completeness (collapsed configurations must be declared, not just acknowledged)
- O6 is γ doctrine triage requirement (every finding gets explicit disposition; "noted" is not a disposition)

α's job in R1 is to inspect each `dd819f00` patch against §R0 and the operator's `expected_change`, then **adopt** the patch as α's correction, **adjust** it, or **reject** it with explicit reason. Adopting is not rubber-stamping — it is α's affirmative agreement that the patch text is correct as written for §R0.

### Per-finding decisions

#### O1 — Wave-master complexity claim (Argument 4 + Q6) — ADOPT (with completeness extension)

**Patch in §R0 (`dd819f00`):** Argument 4 header changed from "Wave-master resolution stays O(1)" → "Wave-master single-root discovery (constant root count)"; body rewritten to state `O(R)` receipt-walk complexity under either model and `O(P + R)` under Model A; explicit acknowledgement that the earlier R0 wording conflated root-count with receipt-count complexity. Q6 mirrored: single-root walk; `O(R)` either way; Model A is `O(P + R)`; "advantage of Model B is **constant root count**, not constant-time receipt discovery."

**α verdict: adopt.** The patch text is correct. The original `O(1)` claim was an asymptotic-complexity error: a walk over `.cdd/unreleased/` is necessarily proportional to the number of receipts (you have to read each cycle's directory to find its receipts). The structural argument never depended on the receipt-walk being constant-time; it depended on the *root count* being constant (one stable root regardless of how many concrete protocols exist). The patch separates the two concerns cleanly.

**α adjustment:** The HI's patch fixed Argument 4 and Q6 but missed the residue in §R0's `CDD Trace` Action-taken bullet (line 116), which still cited "wave-master O(1) resolution" among the structural grounds. α corrects that bullet to "wave-master single-root / constant-root-count discovery" for full consistency. (This is a minor completeness extension, not a substantive disagreement with the HI's patch text.)

#### O2 — Canonical protocol identity anchored to typed `protocol_id` (Verdict + Argument 1 + Q1) — ADOPT

**Patch in §R0 (`dd819f00`):** Verdict paragraph rewritten to anchor canonical protocol identity to the typed receipt's `protocol_id` field (CDS pin: `cnos.cdd.cds.receipt.v1`); issue labels, `gamma-scaffold.md` metadata, and PR body metadata reclassified as operational supporting surfaces, not canonical discriminators. Argument 1 amended with the same anchor: typed `protocol_id` is the canonical discriminator the Go verifier dispatches on. Q1 amended to close with the explicit rule: "Path = ledger owner; `protocol_id` = concrete protocol."

**α verdict: adopt.** This is the materially-strongest correction of the six. The R0 text framed protocol identity as a property of operational surfaces (labels, PR body, scaffold metadata) — those *do* corroborate, but they are not the typed-runtime canonical. The generic CDD receipt schema already carries `protocol_id` specifically so the verifier can dispatch validation mechanically; pinning the architectural anchor there matches the Go-orchestration trajectory (the long-term `cn cdd verify` runtime does not infer protocol from prose or directory names). The decision artifact is stronger after this anchor: "path = ledger owner; `protocol_id` = concrete protocol" is the typed-identity rule the trajectory needs. **α agrees with the HI's patch text verbatim.**

#### O3 — Stale canonical paths in cross-references — ADOPT

**Patch in §R0 (`dd819f00`):** Cross-references list line 1 changed from `src/packages/cnos.cdd/skills/release/SKILL.md` → `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` (the `cdd/` subdirectory restored). Line 2 changed from `cn-cdd-verify (wherever it lives) — inspects .cdd/; unchanged` → `src/packages/cnos.cdd/commands/cdd-verify/ (operator invocation: cn cdd verify) — inspects .cdd/ as unified ledger; validates typed receipts by protocol_id; unchanged`.

**α verdict: adopt.** Path verification confirms both citations land at real directories on the current tree (`src/packages/cnos.cdd/skills/cdd/release/SKILL.md` exists; `src/packages/cnos.cdd/commands/cdd-verify/` exists with `cddverify.go`, `dispatch.go`, `verdict.go`, `README.md`, etc.). The R0 wording "wherever it lives" is the disqualifying language — a decision artifact about canonical ownership cannot itself carry knowingly-uncertain canonical paths. The HI's patch text also augments the verifier description with the typed-receipt-dispatch detail; that augmentation is on-thesis with O2's anchor and α adopts it.

#### O4 — Closure wording conflated cell closure with boundary acceptance — NOT IN α'S MATTER (γ owns)

**Patch in `dd819f00`:** Affects `gamma-closeout.md` only — rewrote "Cycle 497 is closed" to "Cell 497 is closed and awaiting operator boundary decision; cycle acceptance occurs when operator merges PR #499 and closes cnos#497."

**α verdict: not α's call.** The textual surface is in `gamma-closeout.md`, which is γ's matter. α records here only that the operator's finding (cell closure ≠ boundary acceptance per CDS δ doctrine) is correct on its face and that the patch text in `dd819f00` correctly applies that doctrine; γ's R1 takes ownership of whether to adopt as-written. α's `self-coherence.md` §R0 text does not contain the offending wording, so no §R0 edit is required for O4.

#### O5 — Actor-collapse / configuration-floor declaration missing — NOT IN α'S MATTER (γ owns)

**Patch in `dd819f00`:** Adds a new "Dispatch configuration" section to `gamma-closeout.md` declaring collapsed mode + rationale + configuration floor + calibrated success claim.

**α verdict: not α's call.** Same reasoning as O4: γ owns the closeout's per-role-doctrine completeness. α notes that the operator's expected_change matches the CDS actor-collapse-rule's requirement (collapse acknowledged in scaffold + declared in closeout with configuration-floor consequence) and that scaffold §"Cell mode" already acknowledges the first half; γ R1 owns the closeout side. `self-coherence.md` §R0 carries no actor-collapse declaration responsibility (that is gamma-closeout's role), so no §R0 edit is required for O5.

#### O6 — γ process-gap audit table needed explicit dispositions — NOT IN α'S MATTER (γ owns)

**Patch in `dd819f00`:** Reformulates the `gamma-closeout.md` process-gap table with Type + Disposition + Reason columns; both findings classified `no-patch` with reasons; declares `protocol_gap_count: 0`; declares `cdd-iteration.md not required`.

**α verdict: not α's call.** γ doctrine triage of γ-surfaced findings is γ's role. The operator's expected_change is doctrinally correct — γ skill requires explicit disposition of every finding; "mental note" is not a disposition. γ R1 owns the closeout edit. `self-coherence.md` §R0 carries no γ-findings audit (the closeout owns it), so no §R0 edit is required for O6.

### α R1 ownership statement

After inspecting `dd819f00`'s textual edits to `self-coherence.md` §R0 (O1, O2, O3) against §R0's pre-patch text and the operator's expected_change in `operator-review.md`:

- The §R0 corrections for O1, O2, and O3 (as patched by the HI in `dd819f00`) are **adopted as α's R1 corrections**. The patch text becomes α's content. α adds the one completeness extension noted under O1 (CDD Trace bullet wording).
- The §R0 edits required to fix O4, O5, O6 are zero — those findings surface in other artifacts (γ's closeout) and are γ's matter to take ownership of.
- α's R0 substantive decision (Model B; the five structural arguments; the Q1–Q7 answers) is unchanged. The corrections sharpen precision and anchor the typed-identity rule; they do not alter the verdict.

The decision artifact at `.cdd/unreleased/497/self-coherence.md` is now α's work after this R1 pass: §R0 text reflects α's adopted corrections; this §R1 is α's own analysis; the HI-authored §R1 in `dd819f00` is replaced by this section.

### Retrospective on the operator-finding class

The HI's `dd819f00` commit message framed the findings as "narrow mechanical text fixes." That framing was wrong and α records here why:

- O1 required asymptotic-complexity reasoning (recognizing that `O(1)` is the wrong bound for a walk that touches each receipt, and that the load-bearing claim is constant *root count*, not constant *work*).
- O2 required architectural-anchor judgment (deciding that the typed-runtime `protocol_id` field, not the operational corroborating surfaces, is where canonical identity belongs — a Go-trajectory call).
- O3 required canonical-filesystem audit (verifying real paths exist; recognizing that "wherever it lives" is the disqualifying language).
- O4 required CDD doctrine application (the δ skill's named separation of cell closure from boundary acceptance).
- O5 required CDS rule application (actor-collapse-rule completeness; cycle/487 not weakened framing).
- O6 required γ doctrine application (explicit-triage-per-finding rule).

Each is a semantic correction with a defensible doctrinal anchor. Framing them as mechanical text fixes obscured the per-finding reasoning and contributed to the boundary violation (the HI absorbed inline because each individual fix "looked small," when the *aggregate* required a proper role pass).

### α R1 commit-readiness signal

READY. α has taken ownership of the §R0 textual edits in `dd819f00` (O1, O2, O3 adopted with O1's completeness extension applied) and replaced the HI-authored §R1 with this analysis. The decision artifact is α's work. β R1 (independent review of the corrected decision) and γ R1 (closeouts + degraded_recovery declaration for the `dd819f00` overstep) run in parallel and are not dependencies of this α R1 pass.

— α@cdd.cnos, cycle/497 R1, 2026-06-26 (UTC)
