# M5 — Issue Authoring Kata

**Class:** method
**Default level target:** L6
**Purpose:** Prove issue authoring produces an executable work contract — gap, source of truth, constraint strata, ACs with negative space, proof plan, non-goals — without ambiguity, overclaiming, or false closure.

## Scenario

Given a bounded change request and surrounding repo state, produce an issue pack under two modes (baseline / CDD) and compare gap framing, source-of-truth alignment, constraint strata, AC quality, and non-goal discipline.

## Required artifacts

- Problem statement (exists / expected / divergence) — 3–5 lines
- Status truth split (shipped / current spec / draft target / planned / not in scope / unknown)
- Source-of-truth table with canonical paths and status labels
- Scope (in scope / out of scope / deferred / blocked by)
- Constraint strata (hard gate / exception-backed / optional-defaulted / validated-if-present / ignored-deferred)
- Acceptance criteria — numbered, independently testable, with invariant / oracle / positive / negative / surface
- Proof plan — invariant / surface / oracle / positive / negative / operator-visible projection / known gap
- Non-goals (mandatory for 3+ ACs)
- Path resolution base (when paths must be validated)
- Tier 3 skills + active design constraints
- Related artifacts with concrete paths

## Scoring

- L5: issue exists, locally executable
- L6: status truth split present; source-of-truth resolves; ACs have negative space; non-goals named; constraint strata stratified; cross-surface projections enumerated
- L7: gap framing changes future work shape; exception ledger names removable debt; proof plan names what v0 cannot prove

## Worked examples

Three drill scenarios that exercise the issue-authoring discipline (`issue/SKILL.md`) on different problem shapes. Each names anti-patterns the baseline tends to ship.

### Example 1 — Schema validation gate

**Given:** A spec defines fields. Current files may violate it. CI does not enforce it. Some legacy files need temporary exceptions.

**Produce an issue that includes:**
current spec path; hard gate fields; exception-backed fields; optional/defaulted fields; enum list; exception shape; path resolution base; positive and negative fixtures; CI job; notification/status aggregation; non-goals.

**Anti-patterns:**
exception example exempts a hard-gate field; CI job added but notify aggregation omitted; schema and shell both own the same validation rule; path resolution base omitted; "verify CI passes" without negative proof.

**Check:** Can an engineer implement without asking which fields are hard failures? Can reviewer reject a missing negative test? Can CI failure tell contributor what to fix?

### Example 2 — README/source-map alignment

**Given:** The project framing changed. Root README is stale. Docs README has broken links. Upstream formal foundation must be linked. Runtime enforcement does not yet exist.

**Produce an issue that includes:**
current framing; desired framing; affected README files; source map layers; exact doc paths; status labels; link-checking expectation; runtime truth caveat; non-goals.

**Anti-patterns:**
making draft CTB the top-level project frame; linking to stale paths; omitting upstream formal docs; claiming runtime enforcement exists; updating root README but not docs README.

**Check:** Do readers see the current system frame? Do links resolve? Do status labels prevent overclaiming?

### Example 3 — Checker against witness theater

**Given:** A spec defines witness fields. No checker enforces them. A well-written close-out can still be fabricated.

**Produce an issue that includes:**
the risk name; required witness fields; what v0 checks structurally; what v0 does not prove semantically; valid and invalid fixtures; future independence checks; non-goals.

**Anti-patterns:**
field presence treated as semantic truth; no invalid fixture; checker claims runtime enforcement; TSC or upstream verifier treated as already integrated.

**Check:** Does at least one class of false closure become mechanically rejectable? Does the issue avoid claiming more than v0 can prove?
