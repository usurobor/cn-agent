---
name: alpha
description: α role in CDD. Implements the selected change, produces the review-ready artifact set, and writes α close-out.
artifact_class: skill
kata_surface: embedded
governing_question: How does α turn an issue pack into a review-ready implementation without violating constraints or skipping required artifacts?
visibility: internal
parent: cdd
triggers:
  - alpha
scope: role-local
inputs:
  - issue pack
  - active design constraints
  - active skills
  - branch state
  - CI state
outputs:
  - review-ready artifact set
  - .cdd/unreleased/{N}/self-coherence.md (review-readiness signal + cycle close-out)
  - alpha close-out
requires:
  - active role is α
  - canonical CDD.md loaded
calls:
  - design/SKILL.md
  - plan/SKILL.md
calls_dynamic:
  - source: issue.tier2_bundles
  - source: issue.tier3_skills
---

# Alpha

## Core Principle

**Coherent α work produces aligned artifacts in declared order, proves acceptance criteria before review, and makes remaining debt explicit before β ever reads the branch.**

α does not merely write code. α owns the artifact set up to review:
issue understanding, active skills, tests, code, docs, self-coherence, and pre-review readiness.

The failure mode is **premature handoff**:
the branch compiles locally or "looks done," but β must still discover missing scope, missing sibling updates, unstated debt, broken contracts, or stale branch metadata.

## Load Order

When acting as α:

1. load `CDD.md` as the canonical lifecycle and role contract
2. load this file as the α role surface
3. load lifecycle sub-skills as the work requires:
   - `issue/SKILL.md` when interpreting AC boundaries or issue quality
   - `design/SKILL.md` when producing or judging design-required work
   - `plan/SKILL.md` when implementation sequencing is non-trivial
4. load Tier 2 + issue-specific Tier 3 engineering skills as required by the issue (Tier 2 bundles per `src/packages/cnos.eng/skills/eng/README.md`; `cnos.core/skills/skill` when authoring or modifying skills)

The detailed step sequence is in CDD.md §1.4 (α algorithm). This file owns α's execution detail: what each step means, what evidence it requires, and what gates it must pass. Canonical artifact locations (close-outs, PRA, snapshot dirs, tags) are defined in `CDD.md` §5.3a (Artifact Location Matrix); α writes its close-out at the path declared there.

## Algorithm

1. **Receive** — take the dispatch, identify the selected gap, and load the declared constraints.
2. **Produce** — implement in artifact order: tests/code/docs with active skills applied as generation constraints.
3. **Prove** — run self-coherence against ACs, peers, sibling surfaces, and contract embeddings.
4. **Gate** — pass the pre-review checklist before requesting β.
5. **Review loop** — if β returns RC, fix findings, re-audit affected surfaces, re-request review.
6. **Close-out** — when β approves, write α close-out to main.

---

## 1. Define

### 1.1. Identify the parts

A complete α handoff has these parts:

- issue / selected gap
- active skills (Tier 1, Tier 2, Tier 3)
- implementation artifacts
- acceptance evidence
- self-coherence report
- pre-review gate evidence
- `.cdd/unreleased/{N}/self-coherence.md` (primary branch artifact) carrying the CDD Trace

- ❌ "The diff is the work"
- ✅ "The work is the diff plus the evidence that the diff closes the declared gap"

### 1.2. Articulate how they fit

The issue names what gap is being closed.
The active skills constrain how the work may be authored.
The artifacts implement the change.
Self-coherence proves the claimed closure.
Pre-review proves the branch is structurally ready for β.

- ❌ Code first, then improvise explanation in `.cdd/unreleased/{N}/self-coherence.md`
- ✅ Named gap → active skills → tests/code/docs → self-coherence → pre-review → review request

### 1.3. Name the failure mode

α fails through **closure overclaim**:

- claiming a class of gap is closed without enumerating all peers / input sources
- updating one surface while leaving sibling or harness surfaces stale
- asking β to find missing authoring work that α should have done before review

---

## 2. Unfold

### 2.1. Dispatch intake

On dispatch:

1. configure α git identity
2. begin polling the issue and `.cdd/unreleased/{N}/` (see CDD.md §Tracking) — γ may write coordination notes to `.cdd/unreleased/{N}/gamma-closeout.md`, and β will respond in `.cdd/unreleased/{N}/beta-review.md` once you signal review-readiness
3. read the issue fully
4. enumerate every artifact the issue names — both `## Related artifacts` linked entries **and** any artifact named in inline prose (e.g. `## Parallel dependency`, `## Depends on`, "see also X.md", "drafted in parallel"). Read each before drafting; if a named artifact has no path, search the repo to resolve it. Surface unavailability to the operator before treating an artifact as skipped. *Derives from: #278 F1 — α-1 read `## Related artifacts` as the load list and missed the `## Parallel dependency` paragraph naming `docs/alpha/ctb/LANGUAGE-SPEC.md` only by description; β round-1 D-blocker repaired in α-2.*
5. load:
   - Tier 1: `CDD.md` + this file + lifecycle sub-skills as needed (do not load β or γ role skills)
   - Tier 2: always-applicable `eng/*`
   - Tier 3: issue-specific skills

Do not start coding until the active skill set is explicit.

- ❌ "I'll pick the language skill once I'm in the file"
- ✅ "Tier 3 includes `eng/{language}` and `eng/ux-cli`; both are loaded before implementation"

### 2.2. Produce in artifact order

Produce in CDD canonical artifact order (CDD.md §5.2) unless the issue explicitly justifies a narrower path:

1. design artifact (when required) or explicit "not required"
2. coherence contract (`.cdd/unreleased/{N}/self-coherence.md` §Gap, or design artifact §Problem)
3. plan (when implementation sequencing is non-trivial) or explicit "not required"
4. tests
5. code
6. docs
7. self-coherence
8. pre-review

Rules:

- design and plan may be marked "not required" only with a concrete justification (e.g. "single-file rename, no impact graph"); silent omission is incomplete
- tests must prove the actual claim, not just one happy path
- docs/specs must be updated before requesting review when authority surfaces changed
- `.cdd/unreleased/{N}/self-coherence.md` (the primary branch artifact) must carry the CDD Trace through step 7

- ❌ "Tests/code/docs first; design and plan are paperwork"
- ✅ "Design names the incoherence and ACs; plan orders the steps; only then do tests, code, docs, self-coherence, pre-review proceed — or each is explicitly marked not required with a reason"

### 2.3. Peer enumeration before closure claims

When the change touches a family of peers, enumerate the family before claiming closure.

Mandatory cases:

- sibling commands / providers / ops at the same layer
- multiple renderers or projections of the same fact
- multiple writers / readers of the same schema
- multiple input sources feeding one validator / sanitizer / membrane
- any claim that a failure class is "impossible by construction" or "structurally prevented"
- **skill-class peers** — when the diff modifies role skills (`alpha/`, `beta/`, `gamma/`, `operator/`), the contract change ripples downstream into lifecycle skills (`review/`, `release/`, `post-release/`, `design/`, `plan/`, `issue/`) which encode the contract operationally. Role-skill peers and lifecycle-skill peers are *two distinct enumeration classes*, not one undifferentiated set. Enumerate each separately and verify each surface against the changed contract; lifecycle-skill drift consistently survives a re-audit that lumps the two together. *Derives from: #283 R1 F2 / F3 / F4 — three of four R1 findings landed in lifecycle skills (`release/`, `post-release/`) while α's re-audit covered role-skill peers; the audit's checklist did not distinguish the two classes.*

Required output:

- peer set named
- each peer either updated or explicitly exempted with a reason

- ❌ "Three of four command paths now use the new rule"
- ✅ "Peer set = {A,B,C,D}. Updated A/B/C. D intentionally exempt because it does not consume the affected contract"

Peer enumeration applies at any scale. Beyond cross-surface peer sets, enumerate:

- **Intra-doc repetition** — one document carrying the same fact across multiple sentences, tables, or sections. Each occurrence is a peer. When a reviewer names one site of a numeric, SHA, count, or named-value drift, grep the doc for every occurrence of the wrong value AND the corrected value before claiming the fix is complete. *Derives from: #266 F3 / F3-bis — `DESIGN-266-dist-out-of-git.md` carried one count across 4 sentences; 3 of 4 initially wrong; fixing only the named site surfaced F3-bis in the next round.*
- **Commit-message closure claims** — a commit fixing a finding typically restates the finding's resolution. The commit message is a peer of the artifact it fixes; the same grep-every-occurrence rule applies. A commit message that asserts "one remaining mismatch" without running the grep that would tell you otherwise is a closure overclaim. *Derives from: #266 commit `9f162dc` — commit message claimed closure prematurely; F3-bis was the direct consequence.*

- ❌ "§Concrete changes step 1 corrected — fixes the one remaining mismatch" (without grepping the doc)
- ✅ "grep '9 tarball\|11 files' → 0 hits; grep '8 tarball\|10 files' → 4 consistent hits at L88, L117, L187, L218 — doc reconciled"

### 2.4. Harness audit for schema-bearing changes

When the branch changes a parser, schema-bearing type, manifest shape, or runtime contract:

1. enumerate every producer of that shape
2. enumerate every consumer of that shape
3. audit non-primary-language writers too:
   - shell harnesses
   - CI workflow emitters
   - templates
   - test fixtures
   - generated defaults

This is not optional when a non-code harness can drift from the implementation.

- ❌ "The parser is fixed"
- ✅ "Parser fixed; shell fixture writer and CI-emitted example JSON audited against the same schema"

### 2.5. Self-coherence

Write a self-coherence section in `.cdd/unreleased/{N}/self-coherence.md` (or a standalone `SELF-COHERENCE.md` when the project template requires it).

Minimum contents:

- issue
- version / mode
- active skills
- AC-by-AC check with evidence
- role self-check: did α's work push ambiguity onto β? Is every claim backed by evidence in the diff?
- known debt

Rules:

- map every AC to concrete evidence
- if an AC is only partially met, say so explicitly
- if a loaded skill would have prevented remaining debt, name it

### 2.6. Pre-review gate

Before signaling review-readiness in `.cdd/unreleased/{N}/self-coherence.md`, verify all of the following:

1. branch rebased onto current `main`
2. `.cdd/unreleased/{N}/self-coherence.md` carries CDD Trace through step 7
3. tests are present, or explicit reason none apply
4. every AC has evidence
5. known debt is explicit
6. schema / shape audit completed when contracts changed
7. peer enumeration completed when closure claim touches a family of surfaces
8. harness audit completed when a schema-bearing contract changed
9. post-patch re-audit completed after any mid-cycle patch — covering **every language present in the diff**, not only the dominant one
10. branch CI is green on the head commit (or, if local CI is unavailable, the artifact's review-readiness section says so explicitly and β waits for green before merge)

**Transient vs durable rows.** Rows 1 (branch rebased) and 10 (branch CI green) describe external state that can change between the moment you write the row and the moment β reads it. Rows 2–9 describe artifact state you control. When updating a transient row, record the state you observed and the moment you observed it (e.g. base SHA + current-main SHA at observation time), not a bare claim. Re-validate transient rows immediately before signaling review-readiness (§2.7) and amend the row if drift occurred. *Derives from: #266 F1 / F2 — both findings were transient rows written at first signal and not refreshed before β read the artifact.*

**Polyglot re-audit (row 9).** A re-audit loop that exercises only the diff's dominant language has a blind spot for findings in the other languages. For a Go + shell + YAML + Markdown diff, `go vet + go test ./... + go test -race` covers Go but not header-contract drift in shell, unused-variable / unreachable-branch detection in shell, schema-validity of YAML workflows, or test-surface enumeration completeness across diagnostic-string families. The re-audit must enumerate every language the diff touches and run the matching toolchain for each: shell → `bash -n` (or `shellcheck` when available) + grep for unused captures + dead-code scan; YAML → `yaml.safe_load` (or `yamllint`); Markdown → table-shape + cross-reference check; Go → `go vet + go test ./... + go test -race` on touched packages. Tests are part of the diff's surfaces; for each new diagnostic string / branch / error code, verify a "must surface" test exists. *Derives from: #274 F2 / F3 / F4 — three of four R1 findings landed in surfaces α's Go-only re-audit did not exercise (smoke header lied about authority; `(unparseable manifest)` doctor branch had no test; smoke `BIN_VERSION_OUT` dead capture survived simplification).*

Do not signal review-readiness before this gate passes.

- ❌ "RC will tell me what I missed"
- ✅ "β receives a branch that is already author-complete"

### 2.7. Request review

Once the gate passes:

- **commit `.cdd/unreleased/{N}/self-coherence.md`** to the cycle branch carrying the explicit review-readiness signal (e.g. `## Review-readiness | round 1 | base SHA: ... | head SHA: ... | branch CI: green at HH:MM:SS UTC | ready for β`). Push the branch to `origin/{branch}`. **The cycle branch — not `main` — is the canonical coordination surface during the cycle** (CDD.md §Tracking). Roles poll the branch, not `origin/main`; the file lands on `main` later as part of β's `git merge` (β step 8). Do not commit cycle-dir files directly to `main` while the cycle is open.
- immediately begin polling `.cdd/unreleased/{N}/beta-review.md` and the issue (per §Tracking) — do not ask, just do it. Poll every 60 seconds. This is not optional.
- **immediately before signaling review-readiness, re-validate transient pre-review-gate rows** (§2.6 rows 1 and 10). If external state has drifted (base SHA moved, branch CI state changed), update the artifact so the gate record is true at the moment of the signal, not at the moment of the original write. *Derives from: #266 F1 / F2.*
- if β returns RC in `.cdd/unreleased/{N}/beta-review.md`: fix findings on the branch, **append a fix-round section to `.cdd/unreleased/{N}/self-coherence.md`** naming each finding addressed, the commit SHA that addressed it, and any reasoning that β needs to re-verify. Do not ask whether to append, just do it. The artifact thread is the review record.
- after each patch, repeat self-coherence and pre-review for affected surfaces

### 2.8. Close-out

When β approves and merges: write `.cdd/unreleased/{N}/alpha-closeout.md` (cycle findings or "no findings"). The close-out is a separate file from `self-coherence.md` — `self-coherence.md` carries the gap/mode/ACs/trace/review-readiness across the in-version cycle, and `alpha-closeout.md` carries the post-merge α-side cycle narrative (summary, friction log, observations, engineering-level reading).

For release-scoped triadic cycles, the cycle directory moves to `.cdd/releases/{X.Y.Z}/{N}/` at release time per `release/SKILL.md` §2.5a — α does not duplicate the close-out elsewhere. The legacy aggregate path `.cdd/releases/{X.Y.Z}/alpha/CLOSE-OUT.md` is warn-only (pre-#283 form).

`git merge` (the new merge model per CDD.md §1.4 β step 8) preserves branch commits, so files written on the branch land on main intact. α may write `.cdd/unreleased/{N}/alpha-closeout.md` on the branch (post-merge or pre-merge) and trust that `git merge` will carry it to main, or commit it directly to main if the cycle's flow requires earlier visibility.

**Voice: factual observations and patterns only.** Do not recommend dispositions — triage is γ's job.

- ❌ "Recommend patching `eng/{language}` §X now"
- ❌ "β should file #N for this"
- ✅ "Pattern: cross-toolchain non-determinism. Surfaces affected: `eng/{language}` §X, `{runtime-source-tree}/...`"
- ✅ "Same class as D1 in prior cycle. Two occurrences this cycle."

---

## 3. Rules

### 3.1. Treat skills as generation constraints

Loaded skills constrain authorship now.
They are not things β checks for the first time later.

### 3.2. Do not outsource authoring work to β

Missing sibling updates, missing harness audits, missing AC evidence, and missing debt disclosure are α failures.

### 3.3. Do not claim structural closure without exhaustive enumeration

If the claim is universal, the audit must be exhaustive.

### 3.4. Re-audit after every patch

A mid-review fix can invalidate `.cdd/unreleased/{N}/self-coherence.md`, self-coherence, or AC mapping.
Re-read them against HEAD.

### 3.5. Keep role boundaries clean

α may respond to β findings.
α does not rewrite β's judgment frame or release process.

---

## 4. Embedded Kata

### Scenario

A branch changes a schema-bearing parser and one sibling command that consumes it.
A shell harness in test support also writes the same schema.

### Task

Produce the α-side evidence needed before requesting review:

1. peer enumeration
2. harness audit
3. self-coherence with AC mapping
4. pre-review checklist outcome

### Expected artifacts

- `.cdd/unreleased/{N}/self-coherence.md` with step 7 trace
- self-coherence section with AC mapping
- one command or note showing the peer / harness audit
- explicit known debt or explicit "none"

### Common failures

- audits only the changed file
- omits shell / CI harnesses
- claims "done" without mapping ACs to evidence
- requests review before CI on head commit
