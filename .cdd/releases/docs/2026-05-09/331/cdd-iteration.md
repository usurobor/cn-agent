---
retroactive: true
cycle: 331
issue: "#331"
pr: "#332"
merge_sha: "315e5292"
reconstructed_by: "#335"
reconstructed_date: "2026-05-09"
---

# CDD Iteration — Cycle #331

> **RETROACTIVE RECONSTRUCTION** — This artifact was not produced during the cycle.
> Reconstructed post-merge from git history by cycle #335 (issue #335).
> Original merge: `315e5292` on 2026-05-08. Reconstruction date: 2026-05-09.
>
> Irony note: This cycle introduced Step 5.6b (the requirement to produce `cdd-iteration.md` for cycles with `cdd-*-gap` findings) and did not produce this file for itself. Cycle #335 produces it retroactively.

This file lists the six `cdd-*-gap` findings addressed by cycle #331's patches. Each finding was identified in the `usurobor/tsc` cross-repo supercycle retrospective and proposed as a patch to `cnos:src/packages/cnos.cdd/`. All dispositions are `patch-landed` since the patches are already on `cnos:main`.

---

### F1: honest-claim verification rule missing from review skill

- **Source:** γ-triage / PRA §3 — cross-repo supercycle retrospective in `usurobor/tsc`
- **Class:** `cdd-protocol-gap`
- **Trigger:** γ process-gap check — review skill lacked a rule requiring doc claims to match artifacts
- **Description:** `cdd/review/SKILL.md` had no explicit rule requiring that every factual claim in a reviewed document trace to a live artifact. Reviewers could approve "validation performed" claims without verifying the artifact existed.
- **Root cause:** The review skill focused on structural correctness (ACs present, scope defined) but not on the honesty of factual assertions within the reviewed artifact. As cycles grew more complex, self-coherence files began making claims about artifacts that did not yet exist.
- **Disposition:** `patch-landed`
- **Patch:** `1a4f1966`
- **Affects:** `cdd/review/SKILL.md` §3.13

---

### F2: mode declaration and MCA preconditions missing from issue skill

- **Source:** γ-triage / PRA §3 — cross-repo supercycle retrospective in `usurobor/tsc`
- **Class:** `cdd-protocol-gap`
- **Trigger:** γ process-gap check — issues lacked explicit mode headers, causing disconnect confusion
- **Description:** `cdd/issue/SKILL.md` had no requirement for a mode declaration (e.g., `design-and-build`, `docs-only`, `fix`) or MCA precondition gates. Issues defaulted to the full triadic protocol even when a docs-only or small-change disconnect was appropriate.
- **Root cause:** The issue template was designed for full triadic cycles. As the protocol matured, docs-only and small-change paths emerged but were not reflected in the issue template, leaving γ to decide the path ad hoc at dispatch time.
- **Disposition:** `patch-landed`
- **Patch:** `68c18d23`
- **Affects:** `cdd/issue/SKILL.md` mode declaration and MCA preconditions

---

### F3: docs-only disconnect path missing from release skill

- **Source:** γ-triage / PRA §3 — cross-repo supercycle retrospective in `usurobor/tsc`
- **Class:** `cdd-protocol-gap`
- **Trigger:** γ process-gap check — no protocol path existed for cycles that ship only protocol artifacts, with no version bump
- **Description:** `cdd/release/SKILL.md` only specified the tagged-release disconnect (§2.5a: move to `.cdd/releases/{X.Y.Z}/{N}/`). Cycles that ship only documentation, protocol artifacts, or assessments had no canonical path — they either were forced to fake a version bump or left cycle directories in `.cdd/unreleased/` indefinitely.
- **Root cause:** The original release skill was designed when all cycles produced code and version bumps. The triadic protocol's evolution produced docs-only cycles (PRAs, skill patches, retroactive close-outs) that don't warrant version bumps but still need a disconnect.
- **Disposition:** `patch-landed`
- **Patch:** `10f71e7f`
- **Affects:** `cdd/release/SKILL.md` §2.5b

---

### F4: review-rounds and finding-class metrics missing from CDD tracking

- **Source:** γ-triage / PRA §3 — cross-repo supercycle retrospective in `usurobor/tsc`
- **Class:** `cdd-metric-gap`
- **Trigger:** γ process-gap check — no structured metrics existed to detect review churn or mechanical-finding accumulation
- **Description:** CDD had no standardized tracking for review-round counts or finding classification (mechanical / wiring / honest-claim / judgment / contract). Without these metrics, §9.1 triggers (e.g., "mechanical ratio > 20%") could not be computed, and the CHANGELOG ledger lacked a Rounds column.
- **Root cause:** The triadic protocol evolved empirically. Metrics were discussed in individual PRAs but not standardized. The result was that "review churn" could only be estimated from narrative, not measured from structured data.
- **Disposition:** `patch-landed`
- **Patch:** `878caf62`
- **Affects:** `CDD.md` §Tracking, `cdd/post-release/SKILL.md` §5.5, `CHANGELOG.md` ledger header

---

### F5: honest-grading rubric missing from release skill

- **Source:** γ-triage / PRA §3 — cross-repo supercycle retrospective in `usurobor/tsc`
- **Class:** `cdd-skill-gap`
- **Trigger:** γ process-gap check — TSC grades in the CHANGELOG were subjective and not reproducible
- **Description:** `cdd/release/SKILL.md` called for "honest TSC grades" in the CHANGELOG but provided no rubric defining what each grade means. Different cycles applied the letter grades inconsistently, and grade inflation was common (many cycles rated A+/A even when drift items existed).
- **Root cause:** The grading instruction was aspirational ("honest grades") without an operationalizable definition. Without a rubric, γ defaulted to optimistic grades, and the CHANGELOG signal degraded.
- **Disposition:** `patch-landed`
- **Patch:** `b27fc15b`
- **Affects:** `cdd/release/SKILL.md` §3.8

---

### F6: cdd-iteration.md artifact and INDEX.md aggregator not specified in protocol

- **Source:** γ-triage / PRA §3 — cross-repo supercycle retrospective in `usurobor/tsc`
- **Class:** `cdd-protocol-gap`
- **Trigger:** γ process-gap check — cdd-self-improvement findings were buried in PRA §3 prose with no structured form
- **Description:** When a cycle produced `cdd-*-gap` findings, those findings were recorded only as prose in PRA §3. They could not be aggregated, the protocol could not measure its own learning rate, and cross-repo work lost traceability. No canonical artifact type or aggregation file existed.
- **Root cause:** CDD's self-improvement loop was designed for system-level cycles, not for cdd-about-cdd cycles. The absence of a first-class artifact type meant cdd-self-improvement was treated as commentary rather than structured work.
- **Disposition:** `patch-landed`
- **Patch:** `d3c3b3ca`
- **Affects:** `cdd/post-release/SKILL.md` Step 5.6b, `cdd/gamma/SKILL.md` §2.10 row 14
