# β whole-branch review — PR #679 `claude/repo-cleanup-newcomer`

**Reviewed HEAD:** `988ea4e1f28b41a5b49926c2b1c225754a0850ed` (PR base `78c5d467`; whole-PR diff 147 files, +2,481 / −12,554)
**Reviewer:** β (independent, adversarial). Findings derived from the live tree; remediation receipts read only after forming my own view.
**Scope note:** local `main` is stale (merge-base `421399ef` ≠ PR base `78c5d467`), so `git diff main...HEAD` on this checkout is polluted by unrelated main-only code (label-doctor, recursive-cell). Those are **not** part of #679. Reviewed the real PR range `78c5d467...988ea4e1`, which is docs + `.cdd` receipts only. Confirmed against the GitHub API (additions 2,481 / deletions 12,554 / 147 files).

---

## Verdict: REQUEST-CHANGES (narrow) — one residual status-truth miss on a `Status: Current` security spec; everything else is APPROVE-quality.

Headline: the cleanup is genuinely good and the 7 D1 corrections verify honest — **but** `docs/architecture/security/SECURITY-MODEL.md` still presents removed OCaml modules as the current locus of security enforcement, under `Status: Current`, with no archived-runtime banner, **while this same PR added exactly that banner to its two sibling OCaml-citing docs.** That is the flagged D1 defect class, unresolved on one doc. Fix is one line already authored twice in this PR.

**Findings by severity:** 1 blocking (D/M1) · 1 major-but-disclosed/stack-contingent (M2) · 4 nits (N1–N4).

---

## D-level (blocking) — fix before merge

### M1 — SECURITY-MODEL.md claims `Status: Current` but cites deleted OCaml modules as where enforcement lives (status-truth, highest-priority dimension)

- **Where:** `docs/architecture/security/SECURITY-MODEL.md:3` (`**Status:** Current`), `:104` ("The four typed FSMs in `cn_protocol.ml` provide compile-time safety"), and the "Where enforcement lives" table `:151`–`:155` citing `cn_protocol.ml`, `cn_agent.ml`, `cn_mail.ml`, `cn_commands.ml`, `cn_hub.ml`.
- **Evidence:** `find . -name '*.ml'` → **0 files in tree**. `grep sender_transition|receiver_transition|materialize_branch src/go src/packages` → **none**. The named modules do not exist; the OCaml tree is off `main` (branch `legacy/ocaml-thread-reference`). The doc has **no** archived-runtime banner.
- **Why blocking:** This is the exact defect D1 was raised to fix ("no unimplemented/removed thing reads as shipped"), and the PR's AC "Lifecycle status is honest" is checked. A newcomer reading the canonical *security* spec cold is told, as current fact, that compile-time FSM safety and IO/mail/hub enforcement live in `cn_*.ml` — then greps and finds nothing. Trust-destroying, and it's the security doc.
- **Self-inconsistency that removes any benefit-of-the-doubt:** this same PR (`78c5d467...988ea4e1`) **added** the archived-OCaml banner to `docs/architecture/ARCHITECTURE.md:9-15` and `docs/reference/cli/CLI.md:173-175` (verified in the diff). SECURITY-MODEL got only its `Date:` line removed — the banner pass skipped it.
- **Exact fix:** add the same banner already used in ARCHITECTURE.md (copy `docs/architecture/ARCHITECTURE.md:9-15`, adjusting "throughout this document"), e.g. immediately under the Status header:
  > **Runtime note.** The active runtime is the Go `cn` binary (`src/go/`). The `cn_*.ml` modules named below are the archived CN thread reference implementation (branch `legacy/ocaml-thread-reference`); the shipped Go runtime implements the enforcement subset under `src/go/internal/` (`dispatch/`, `cli/cmd_*.go`, `doctor/`). The FSMs and module table describe the reference architecture, not current Go symbols.

  Better still (optional): repoint the enforcement table at the real Go loci (`src/go/internal/cli/cmd_*.go`, `internal/dispatch/`, `internal/doctor/`), which do exist and do implement the "all effects go through `cn`" gatekeeper principle.

With M1 fixed this branch is a clean APPROVE.

---

## Per-dimension findings

### 1. Status-truth / honest claims — **substantially resolved, ONE residual miss (M1)**

I re-ran the falsification greps myself before trusting any label:

- `PROTOCOL.md` — was **falsely `Implemented (v2 — cn_protocol.ml)`**; now `Normative target — not implemented in the Go runtime` + banner. Verified: 0 `.ml` files, `sender_transition|receiver_transition|materialize_branch` absent from `src/go`. **Correct.**
- `BUILD-AND-DIST.md` / `PACKAGE-ARTIFACTS.md` — `Implemented`; cited impl paths **all exist**: `src/go/internal/pkgbuild/build.go`, `cli/cmd_build.go`, `restore/restore.go`, `cli/cmd_deps.go`, `pkg/pkg.go`, `activation/index.go`. **Honest.**
- `MEMORY.md`, `HYBRID-LLM-ROUTING.md`, `MESSAGE-PACKET-TRANSPORT.md`, `HUB-PLACEMENT-MODELS.md`, `THREAD-EVENT-MODEL.md`, `GIT-CN-PACKAGE.md`, `ORCHESTRATORS.md`, `EXTENSION-REGISTRY.md`, `RUNTIME-EXTENSIONS.md` — all now `Draft (proposed — not implemented)` + implementation line; present-tense "defines/replaces" softened to "proposes/would replace". Spot-checked greps (`placement.json`, `cn.packet`, `routing_policy`, `cn.orchestrator`, `cn.extension`) → **none in `src/go`**. **Correct.**
- `POLYGLOT-PACKAGES-AND-PROVIDERS.md` — `Partially implemented`, split honestly (commands ship in `pkg.go`; provider surface is a target). **Honest.**
- `COGNITIVE-SUBSTRATE.md` — `Normative — enforced`; taxonomy is enforced by `pkg.go` content classes + `doctor.go`. **Honest.**

**Missed (M1):** SECURITY-MODEL.md (above). Also a lesser instance —
**N2:** `docs/reference/runtime/AGENT-RUNTIME.md:229` "proposes a native OCaml agent runtime … that replaces OpenClaw", citing `cn_*.ml` throughout with no banner. It is honestly `Status: Draft`, so it does **not** read as shipped — lower severity — but it proposes building on the abandoned OCaml runtime and would benefit from the same banner. Nit, not blocking.

Determination: **the status-truth defect is NOT fully resolved.** The 7 corrections are sound, but the audit scoped itself to the 14 *trimmed* specs and spot-checked only for *presence of a Status header* — so a header-bearing but OCaml-stale doc (SECURITY-MODEL, `Status: Current`) passed the spot-check while still misrepresenting its lifecycle.

### 2. Deletion appropriateness — **clean, one prose orphan (N1)**

Deleted classes are genuine history: `docs/development/plans/*` (16), `kata/runs/*`, `evidence/rca/2026-*`, cycle-logs/critiques, `POST-RELEASE-EPOCH-*`, audits, and issue-keyed `reference/` design records (`DESIGN-227-*`, `SELF-COHERENCE-227`, `CORE-REFACTOR`, `PLAN-174-*`). Spot-checked deleted `docs/reference/legacy/OCAML-THREAD-REFERENCE.md` and `schemas/DESIGN-LLM-SCHEMA*` — no active inbound refs; the OCaml is reachable via the `legacy/ocaml-thread-reference` git branch that surviving docs point to. No load-bearing current-state doc deleted.

**N1 (nit):** `docs/reference/packages/DESIGN-266-dist-out-of-git.md:76,237` still names the now-deleted `DESIGN-227-distribution-pipeline.md` and `SELF-COHERENCE-227.md` as current descriptions of `cn build`. Prose, not a markdown link (so lychee-invisible), but a live doc pointing a reader at deleted files. Fix: drop those two filenames or repoint to the surviving `PACKAGE-AUTHORING.md` / `PACKAGE-SYSTEM.md` already in the same sentence.

### 3. Link integrity — **clean on the reader surface**

Independent offline scan of all 1,393 `**/*.md` (my own resolver, lychee-`--offline`-equivalent): **129 dead internal links, 100% inside `.cdd/`**; **zero** in the active reader surface (`docs/`, root, `src/packages/**`). `lychee.toml` sets `exclude_path = [".cdd"]` by policy (frozen work-cell evidence, audited via `cn cdd verify`, not gated), so the CI I4 gate scope is clean. Deleted-basename sweep found only historical CHANGELOG rows and the N1 prose pointer. Frontmatter/nav indexes resolve.

### 4. Doctrine / portal consistency — **M2 (standalone contradiction, disclosed & stack-contingent)**

`DOCUMENTATION-SYSTEM.md` §5 was in fact substantially rewritten by #679 (v3.13.0 → v3.82.0) and now holds **two paragraphs in tension**:
- `:100` (para 1, NEW): "Released snapshots are **not** kept as folders on `main` — git history is the archive. A document on `main` is always the current one."
- `:102` (para 2, OLD): "A frozen record (a dated design decision, a cycle log, a completed plan) **is left in place** under its dated or `Superseded` header … Only the live reader surface is kept current."

#679's *action* deletes exactly those frozen records — directly contradicting §5 para 2 as shipped in #679. `docs/README.md:10` and the C1-fixed footer `:62` both assert the NEW stance. So standalone-#679 is incoherent: doctrine para 2 says "leave frozen records in place," the PR deletes them. **As a stack with #681** (which carries the §5 reversal) this resolves, and the PR's Known Debt explicitly says "reviewed as a stack with #681; do not merge with an intermediate closure claim." Given that disclosure I do not escalate M2 to an independent blocker, but flag: either land the stack atomically, or pull the §5-para-2 deletion into #679 so #679 is self-coherent.

**N3 (nit):** the C1 commit message claims the new footer "matches … §5." It matches §5 **para 1** but contradicts §5 **para 2** — a half-true closure claim.

### 5. PR-body truth — **accurate**

- "1 cleanup cell, 9 passes — β in passes 1–3 & 5; γ-verify 4 & 6–9" — matches the round-0{1..9} + review-fix receipts. ✓
- Two stat numbers: reader-surface "125 files, +294 / −12,554" and whole-PR "147 files, +2,481 / −12,554". The whole-PR triple **matches the GitHub API exactly** (147 / +2,481 / −12,554). ✓
- **N4 (nit):** title says "−12.3k"; body says −12,554. Reconciles as reader-surface *net* (+294 −12,554 = −12,260 ≈ −12.3k). Not an overclaim, just two different denominators; harmless.
- The D2 lineage correction (was "9 β→α→γ cells", now honest pass structure) and the D1/C1/C2/B1 remediation rows all correspond to real commits (`a7db6dda`, `988ea4e1`, `39d78197`, merge `cd4050ce`). No remaining overclaim except the checked AC "Lifecycle status is honest" which M1 falsifies for one doc.

### 6. CI — **11/11 green on the reviewed HEAD `988ea4e1`**

Confirmed via check-runs API: Go build & test, Package/Binary verification, I1 drift, I2 protocol schema sync, I4 repo-link validation, I5 SKILL frontmatter, I6 CDD ledger, and the #516/#524/#648 dispatch guards — all `success`. No reason to doubt any gate: the changes are docs + receipts (no Go source touched in the PR range), I4's active-docs scope is clean (dimension 3), and I6/ledger passed. Green is credible.

### 7. Newcomer coherence — **strong and honest**

Cold path root `README.md` → `docs/README.md` → specs: identity is clear and front-loaded ("cnos gives AI agents a Git-native home … a model such as Claude or GPT acts as the temporary reasoning body"), the failure-mode/answer table is concrete, and the shipped/planned boundary is now legible thanks to the D1 status labels (`Draft (proposed — not implemented)` vs `Implemented`). The one thing that would trip a careful cold reader is M1: the security spec's OCaml citations.

---

## What I converge on (genuinely good)

- The D1 remediation is real work, not cosmetic: every "Implemented" label I checked cites files that exist; every "Draft (proposed)" I checked falsifies clean in `src/go`; the false `PROTOCOL: Implemented` and the mislabeled `BUILD-AND-DIST: Draft` were both correctly flipped by evidence.
- Link integrity on the reader surface is **zero-dead**, and the `.cdd` exclusion is principled and matches CI.
- Deletions are true history; no current-state load-bearing doc was lost (one recoverable prose pointer, N1).
- PR body stats and lineage are honest and match the API; CI is green on the exact reviewed SHA.
- Newcomer identity is clear, and the shipped-vs-proposed boundary is materially better than before this branch.

The branch is one banner away from clean. Add the archived-runtime note to SECURITY-MODEL.md (M1); confirm the #681 stack lands atomically for §5 (M2); optionally sweep N1–N4.
