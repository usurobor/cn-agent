# α design notes — cycle/393

**Issue:** cnos#393
**Mode:** design-and-build, γ+α+β-collapsed-on-δ
**Goal:** 4 coordinated skill patches forming a coherent referential mesh, plus a body edit to cnos#366.

## §Design half — the doctrine being captured

The doctrine is: **δ is a two-sided membrane.** CCNF (`COHERENCE-CELL-NORMAL-FORM.md`) and the RECEIPT-VALIDATION design name δ as *the boundary*, but they describe δ-as-outward only: receipt + verdict → boundary decision at scope `n+1`. The complementary inward face — γ's contract → α-ready dispatch — has been done implicitly by the operator at every cycle and has only been formalized as a per-cycle ad-hoc operator action (cnos#392 was the first cycle where δ explicitly pinned the 7-axis contract at dispatch).

The four patches transform this from operator-implicit to skill-explicit, distributed coherently across the four role surfaces:

```
       γ writes the contract template (Patch C; surfaces 7 axes)
                          │
                          ▼
       δ enriches at dispatch (Patch D; inward-membrane function)
                          │
                          ▼
       α MUST NOT improvise (Patch A; constrains implementation)
                          │
                          ▼
       β verifies conformance (Patch B; binding gate at review)
```

Each patch is independently coherent (it makes sense in isolation) but the mesh ensures a future session loading any one is pointed at the others. The mesh is unidirectional pairwise via "see X for the role-side enforcement" links — no circular reasoning required.

## §Build half — the 4 patches, patch by patch

### Patch A — `alpha/SKILL.md` insertion plan

**Insertion point:** under `## 3. Rules`, after current `### 3.5. Keep role boundaries clean` (line ~341–345 in current file). New subsection `### 3.6. Implementation contract is δ's, not α's`.

**Note on numbering:** The issue body says "Rule 8 (or whatever next number — check current numbering)." α's current Rules section is numbered 3.1–3.5 (one-period style), not "Rule N" as β uses. Adding 3.6 matches α's existing convention. Cross-referencing in other patches will use the human-readable name "Rule 3.6: Implementation contract is δ's, not α's" so the reference is stable even if α's section numbering shifts.

**Content sketch:**

```markdown
### 3.6. Implementation contract is δ's, not α's

α MUST NOT change the implementation-contract axes pinned by δ at dispatch
time: language, CLI integration target, package scoping, runtime
dependencies, existing-binary disposition, JSON/wire contract preservation,
backward-compat invariants. If any of these is unpinned in the dispatch
prompt's `## Implementation contract` section, α MUST surface to γ/δ
before coding; α MUST NOT improvise.

The 7 axes form the "implementation contract" — the architectural shape
the cycle ships, distinct from the behavioral ACs the cycle satisfies.
γ writes the contract values at dispatch (`gamma/SKILL.md` §2.5 Step 3b
`## Implementation contract` template); δ ratifies and enriches if γ
under-specified (`operator/SKILL.md` §"δ as inward membrane");
β verifies α's implementation conforms before APPROVE
(`beta/SKILL.md` §Role Rules Rule 7: Implementation-contract coherence).

Empirical anchor: cnos#389 (α implemented V in Python despite cnos being
Go-native; the α SKILL did not name "language" as a δ-pinned axis;
α had room to improvise) and cnos#391 (α placed the Go port in a separate
binary at the wrong package path; α improvised on package scoping and
CLI integration). In both cycles β's behavior-only AC oracles APPROVE-d
without catching the implementation-contract drift; cnos#392 was the
first cycle where δ pinned the contract at dispatch as an ad-hoc
operator action, and the cycle succeeded specifically because of it.

- ❌ "The dispatch didn't say which package to use; I'll pick `pkg/foo/`"
- ❌ "Python's faster for this; I'll switch from Go"
- ✅ "Implementation contract row 3 (package scoping) is unpinned; surfacing to γ via the artifact channel before coding"
```

### Patch B — `beta/SKILL.md` insertion plan

**Insertion point:** under `## Role Rules`, after current `### 6. Anchor oracle evidence on code, not doc` (line ~146–162 in current file). New subsection `### 7. Implementation-contract coherence`.

**Content sketch:**

```markdown
### 7. Implementation-contract coherence

Before APPROVE, β verifies the cycle's implementation conforms to the
implementation contract pinned at dispatch (the 7 axes in
`gamma/SKILL.md` §2.5 Step 3b `## Implementation contract` template:
language, CLI integration target, package scoping, runtime dependencies,
existing-binary disposition, JSON/wire contract preservation,
backward-compat invariants).

Behavior-only ACs (does V validate? does the parser accept the new
shape?) are necessary but not sufficient. Implementation-contract
conformance is a binding gate parallel to Rule 6's code-first oracle
anchoring: where Rule 6 prevents doc-vs-code drift, Rule 7 prevents
behavior-vs-shape drift.

Verification: for each axis pinned in the dispatch prompt's
`## Implementation contract` section, β confirms the diff conforms.
If the diff diverges from any pinned axis without an explicit γ
clarification updating the contract, the verdict is REQUEST CHANGES
with a D-severity finding classified `implementation-contract`.

α owns the constraint (`alpha/SKILL.md` §3.6: α MUST NOT improvise
implementation contract). δ owns enrichment at dispatch
(`operator/SKILL.md` §"δ as inward membrane"). β owns verification.

Empirical anchor: cnos#389 R1 and cnos#391 R1 both APPROVE-d on
behavior-only ACs. Neither caught the implementation-contract drift
(Python-not-Go in #389; wrong package + separate binary in #391).
Adding this rule makes the implementation-contract surface a
first-class review gate, dual to the α-side constraint in §3.6
and the γ-side template in §2.5.

- ❌ "AC oracles passed; APPROVE." (without checking the 7 axes)
- ❌ "Implementation works; the package path is a minor detail."
- ✅ "AC oracles pass; implementation contract row 1 (language=Go) confirmed; row 3 (package scoping = `src/go/internal/cdd/...`) confirmed; APPROVE."
- ✅ "AC oracles pass but row 1 (language) drifted to Python; RC with D-severity finding, classification `implementation-contract`."
```

### Patch C — `gamma/SKILL.md` §2.5 insertion plan

**Insertion point:** under `## 2.5 ... #### Step 3b — Subscribe and dispatch α and β`, immediately after the existing `**α prompt (γ produces, δ dispatches):**` code block (line ~338) and before the `**β prompt (γ produces, δ dispatches):**` block. The template addition appends a required `## Implementation contract` section that γ writes into both α and β prompts.

**Concrete approach:** Rather than mutating the existing α-prompt code block (which is a minimal canonical form), add an explicit subsection titled `##### Implementation contract section (required)` between the α and β prompt examples. This keeps the existing prompt skeleton intact (backward compat) and adds the new γ obligation as a clearly-labeled extension.

**Content sketch:**

```markdown
##### Implementation contract section (required for α prompt)

γ MUST include a `## Implementation contract` section in the α prompt
(and reference it from the β prompt) enumerating the 7 architectural
axes that pin the cycle's shape:

\`\`\`markdown
## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | <Go ; Python ; TypeScript ; Markdown ; ...> |
| CLI integration target | <`cn` subcommand ; standalone binary ; library only ; N/A> |
| Package scoping | <e.g. `src/go/internal/cdd/...` ; `src/packages/<pkg>/...` ; N/A> |
| Existing-binary disposition | <preserve ; replace ; deprecate ; N/A> |
| Runtime dependencies | <list or "None"> |
| JSON/wire contract preservation | <preserve as-is ; explicit shape change with migration ; N/A> |
| Backward-compat invariant | <e.g. "existing rules preserved; new content additive" ; "breaking change documented in §Migration"> |
\`\`\`

γ writes the contract values per repo conventions or escalates to δ
for ratification before dispatch. **γ MUST NOT dispatch with empty /
"TBD" rows.** If γ doesn't know a value, γ asks δ — the inward
membrane function (`operator/SKILL.md` §"δ as inward membrane")
exists exactly for this enrichment. δ fills the row or blocks dispatch
if the row is genuinely undecidable (escalate to operator-as-human).

α's role-side constraint (α MUST NOT improvise) is in
`alpha/SKILL.md` §3.6. β's review-side gate (verify conformance
before APPROVE) is in `beta/SKILL.md` §Role Rules Rule 7. The four
surfaces form one coherent mesh: γ template (this section), δ inward
enrichment, α MUST NOT improvise, β verifies.

Empirical anchor: cnos#389 (Python-not-Go) and cnos#391 (wrong
package scoping + separate binary) both shipped wrong-shape
implementations because γ's dispatch prompt under-specified the
contract and α improvised. cnos#392 was the first cycle where δ
pinned the contract at dispatch as an ad-hoc operator action;
cnos#393 (this patch) makes the contract a first-class γ obligation.

- ❌ Dispatch with rows left empty or marked "TBD"
- ❌ Dispatch with implicit conventions ("everyone knows we use Go")
- ✅ Every row populated explicitly; δ escalation logged in
  `gamma-clarification.md` when γ asked δ to fill a row
```

### Patch D — `operator/SKILL.md` insertion plan

**Insertion point:** new top-level section between current `## 3. Gate` and `## 4. Override`. Title: `## 3a. δ as inward membrane: implementation-contract enrichment at dispatch`. Using `3a` rather than promoting it to a new top-level number minimizes downstream renumber-drift on §4 / §5 / §6 / §7 references throughout the file.

**Note on Phase 4:** Per the issue, "Patch D ... is the design-prerequisite for Phase 4 (δ split) which will eventually move this to `delta/SKILL.md`." So this section is intentionally bounded as a design-prerequisite anchor; it names the doctrine, references the γ template, and notes the Phase 4 move-target. Phase 4 (cnos#366) takes responsibility for relocating it to `delta/SKILL.md` when the split lands.

**Content sketch:**

```markdown
## 3a. δ as inward membrane: implementation-contract enrichment at dispatch

`COHERENCE-CELL-NORMAL-FORM.md` names δ as the cell's boundary — the
actor that receives the receipt and verdict and records a boundary
decision (`accept`, `release`, `override`, `reject`, `repair_dispatch`).
That is the **outward-facing membrane**: receipt + V verdict →
parent-scope boundary decision. The current `## 3. Gate` section
above describes this surface.

δ is also a two-sided membrane. The complementary face is **inward**:
γ's protocol-level contract (gap, ACs, oracle, evidence) →
α-ready dispatch enriched with the implementation-contract axes
α needs to execute coherently.

\`\`\`text
δ as two-sided membrane:

  outward:  receipt + V verdict → parent-scope boundary decision  (§3 above)
  inward:   γ contract → α-ready dispatch                          (this section)
            (implementation contract enriched here)
\`\`\`

**The inward function.** γ writes the protocol-level contract per
`gamma/SKILL.md` §2.5 Step 3b, including the
`## Implementation contract` section enumerating the 7 axes (language,
CLI integration target, package scoping, existing-binary disposition,
runtime dependencies, JSON/wire contract preservation, backward-compat
invariant). γ populates the rows from repo conventions and the
issue body. **δ reviews γ's dispatch prompt before routing it to α
and ensures every row is populated.** If γ left a row unpopulated
or marked "TBD," δ fills it per repo conventions, or — if the row
is genuinely undecidable (e.g. the language choice is part of the
cycle's design question, not its execution shape) — blocks dispatch
and escalates to operator-as-human.

**Why this is δ's, not γ's, surface.** γ writes what the cycle is
*for* (gap, ACs, oracle). δ writes what the cycle's output is
*shaped like* (language, package, integration target). The two
contracts are distinct: γ's is protocol-level; δ's is
implementation-level. Mixing them produced cnos#389
(α improvised language) and cnos#391 (α improvised package
scoping + binary disposition). cnos#392 was the first cycle
where δ pinned the implementation contract at dispatch; the
cycle succeeded specifically because of it.

**The mesh.** This section is the δ side of a four-surface mesh:

- γ template:    `gamma/SKILL.md` §2.5 Step 3b
                 (`## Implementation contract` template — the 7 axes)
- α constraint:  `alpha/SKILL.md` §3.6
                 ("Implementation contract is δ's, not α's")
- β verification:`beta/SKILL.md` §Role Rules Rule 7
                 ("Implementation-contract coherence")
- δ enrichment:  this section (inward-membrane function)

**Phase 4 (δ split) — relocation target.** This section is a
**design-prerequisite anchor** for Phase 4 of cnos#366 (δ split:
`operator/SKILL.md` → `delta/SKILL.md` + harness substrate).
When Phase 4 lands, this section moves to `delta/SKILL.md` as
part of the membrane-policy surface. The two-sided framing is
the design input Phase 4 absorbs; the relocation itself happens
in Phase 4's cycle, not here.

Empirical anchor: cnos#389 (Python-not-Go) and cnos#391 (wrong
package scoping + separate binary) are the failure-mode evidence
that motivates this section. cnos#393 (this patch) makes the
inward function doctrine; Phase 4 implements δ-inward in
`delta/SKILL.md`.

- ❌ δ routes γ's α prompt with rows blank — "α can figure it out"
- ❌ δ fills rows by guessing — no consultation with γ on intent
- ✅ δ reviews γ's `## Implementation contract` section row-by-row;
  enriches per repo conventions; escalates the row if undecidable
- ✅ δ logs any enrichment in `.cdd/unreleased/{N}/gamma-clarification.md`
  (or a δ-specific channel if Phase 4 has carved one) so the contract
  trail is auditable
```

## cnos#366 Phase 4 body update plan (AC5)

The current Phase 4 section already lists inputs (#371, #373, #384, #375). The minimal-change update adds:

- One new bullet under "Inputs surfaced by parallel cycles ...": `#393 — δ-as-architect: implementation-contract enrichment at dispatch. Names δ as explicitly two-sided membrane (outward + inward) and ships the precursor skill patches (α/β/γ/operator) that Phase 4 absorbs.`
- One sentence appended to the Phase 4 deliverable bullet for `delta/SKILL.md`: `Includes the inward-membrane surface from #393 (implementation-contract enrichment at dispatch).`
- One sentence at the top of the Phase 4 narrative: `δ is **two-sided** — outward (boundary decision on receipts; current scope) AND inward (implementation-contract enrichment of γ's contracts going to α). Phase 4 must carve both surfaces explicitly.`

Backward compat: existing Phase 4 prose preserved; additions are additive.

## §Cross-reference plan (AC6 evidence)

Each patch cites the other three by name + path (unidirectional pairwise; no circular reasoning). The citations are for **discoverability** (a future α/β/γ session loading any one finds the others), not for justification (each rule is locally self-justifying via the empirical anchors).

- Patch A (α §3.6) cites: γ template (Patch C), β Rule 7 (Patch B), δ inward (Patch D), #389, #391.
- Patch B (β Rule 7) cites: γ template (Patch C), α §3.6 (Patch A), δ inward (Patch D), #389 R1, #391 R1.
- Patch C (γ §2.5 §§) cites: α §3.6 (Patch A), β Rule 7 (Patch B), δ inward (Patch D), #389, #391, #392.
- Patch D (δ inward) cites: γ template (Patch C), α §3.6 (Patch A), β Rule 7 (Patch B), #389, #391, #392.

`rg "alpha/SKILL.md §3.6\|beta/SKILL.md.*Rule 7\|gamma/SKILL.md §2.5\|operator/SKILL.md.*inward membrane"` against the diff hunks shows the bidirectional structure at review time.
