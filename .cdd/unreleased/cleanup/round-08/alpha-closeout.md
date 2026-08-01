# α Closeout — Round 08 (docs-only TRIM: strip CDD design-cycle apparatus from canonical specs)

**Role:** α (implementer). **Branch:** `claude/repo-cleanup-newcomer`. **Base HEAD:** `fa75ed94`.
**Mandate:** de-historicize the non-dot reader surface to current-state only; git history is the archive. Dotdirs exempt. Standard: cnos write skill.
**Contract:** γ Round-07 §2/§5 — trim (not delete) the design-cycle apparatus from 6 canonical reference/architecture specs + sweep header stamps on 5 siblings. No commit, no PR.
**Result:** 11 docs edited, all under `docs/` (non-dot reader surface). `git diff --stat`: 37 insertions / 571 deletions. No defer path touched.

Apparatus classes removed across the 6 targets: provenance header stamps (`Issue:`/`Mode:`/`Active Skills:`/`Engineering Level:`), `## Impact Graph`, `## Acceptance Criteria`, `## CDD Trace`, `## File Changes` (cycle create/edit/delete plans), `## Known Debt` (per-item judged), and `## Problem` gap-framing reframed to a current-state `## Purpose`. Kept in every doc: governed markers (`Version:`/`Status:`/`Doc-Class:`/`Owns:`/`Parent:`), the full spec body, and design-rationale sections (Constraints, Challenged Assumption, Decision/Proposal, Leverage, Alternatives Considered, Prior Art, Non-goals, Migration, Process Cost, Tuning Strategy).

---

## Per-doc detail (the 6 targets)

### 1. `docs/reference/runtime/MEMORY.md`
- **Removed stamps:** `**Issue:** #100`, `**Mode:** MCA`, `**Active Skills:** …`, `**Engineering Level:** L7`. Kept `**Version:** 0.2.0`.
- **Reframed:** `## Problem` (the "Today: … / What is missing / current incoherence 1–4" gap-narration) → a 1-paragraph `## Purpose`. The current-state runtime facts it duplicated already live in `## Current Evidence` (kept), so no fact lost.
- **Removed sections:** `## File Changes` (RUNTIME-CONTRACT/cn_runtime_contract/AGENT-RUNTIME edit plan + "No new files in v1"); `## Acceptance Criteria` (7 `[ ]`); `## Known Debt` (3 items, all pure future-work — restore-map "may become useful later", recall/search tooling "worth packaging later", zone-name "tightening pass"); `## CDD Trace` (Step 0–6 table).
- **Anchor repair:** `## Related` referenced removed content — dropped the "(this design supersedes the original spec direction)" narration and rewrote "#156 — Attached hubs (AC9: …)" → "(agent memory stays in hub, tagged by workspace)", removing the now-dangling AC9 pointer. Doc-to-doc links (AGENT-NETWORK.md, HUB-PLACEMENT-MODELS.md) kept.
- **Preserved:** Decision (triadic α/β/γ model), Constraints, Challenged Assumption, Current Evidence, full Proposal (§1–6), Leverage/Negative Leverage, Alternatives Considered. No spec substance lost.

### 2. `docs/reference/runtime/HYBRID-LLM-ROUTING.md`
- **Removed stamps:** `**Issue:** TBD`, `**Mode:** MCA`, `**Active Skills:** …`, `**Engineering Level:** L7`. Kept `**Version:** 0.1.0` + `**Status:** Draft`.
- **Reframed:** `## Problem` → `## Purpose` (kept the two-shapes rationale as "shapes the policy avoids"; dropped the "cnos is moving toward…" gap setup).
- **Removed sections:** `## Impact Graph` (Downstream/Upstream code-impact + Authority relationships); `## File Changes` (future edit plan); `## Acceptance Criteria` (AC1–AC10 `[ ]`); `## Known Debt` (4 items, all future tuning). `## CDD Trace`.
- **Substance-preservation edit:** the one Impact-Graph authority nuance not stated in the body — "commands/orchestrators may declare routing hints, but the runtime makes the final decision" — was folded into Proposal §1 ("Skills may shape the task, and commands or orchestrators may declare routing hints, but the body makes the final provider decision"), so removing Impact Graph loses no authority fact.
- **Preserved:** Constraints, Challenged Assumption, full Proposal §1–8 (placement, routing model, task classes, thresholds, dependency depth, receipts, runtime-contract JSON, failure policy), Leverage/Negative Leverage, Alternatives, Process Cost, Non-goals, **Tuning Strategy** (kept — operational tuning approach, not enumerated apparatus).

### 3. `docs/reference/runtime/POLYGLOT-PACKAGES-AND-PROVIDERS.md`
- **Removed stamps:** `**Issue:** #170`, `**Mode:** MCA`, `**Active Skills:** …`, `**Engineering Level:** L7`. Kept `**Version:** 0.2.0`.
- **Reframed:** `## Problem` → `## Purpose` (kept the Go/Rust/packages shape breakdown + two-shapes-avoided list; dropped "cnos is converging…" gap setup).
- **Removed sections:** `## File Changes`; `## Acceptance Criteria` (8 plain bullets — criteria for accepting the cycle); `## Known Debt` (4 items, all future work — provider protocol spec, index format "later", first migration "future work", surfacing "follow-up"); `## CDD Trace`; and the trailing "If you want, the next useful artifact is…" design-cycle chatter line.
- **Preserved:** Decision (core/runtime/authority rules), Constraints, Challenged Assumption, Prior Art, full Proposal §1–9 (package substrate, kernel-as-package, commands vs providers, command + provider execution contracts incl. Go descriptor code + extension JSON, A2A split, artifact/platform targeting, runtime-contract implications, migration phases), Leverage/Negative Leverage, Alternatives. No `## Impact Graph` existed here.

### 4. `docs/reference/packages/PACKAGE-ARTIFACTS.md`
- **Retitled:** `# Design: Package Artifact Distribution…` → `# Package Artifact Distribution and Command Content Class` (dropped the historicizing `Design:` prefix).
- **Removed stamps:** `**Issue:** #167`, `**Mode:** MCA`, `**Active Skills:** …`, `**Engineering Level:** L7`. Kept `**Version:** 0.5.0`.
- **Reframed:** `## Problem` (two gaps + point-in-time incident evidence: #155 sandbox git failure, `cn_deps.ml` 27-refs, #162) → a 2-bullet `## Purpose`. The "why" rationale is preserved in the kept `## Challenged Assumption` (two replaced assumptions); the incident forensics are git-history archive.
- **Removed sections:** `## Impact Graph` (Downstream/Upstream/Copies-and-authority — authority facts already stated in Proposal §2/§3/§6 and Challenged Assumption); `## File Changes` (Create/Edit/Delete src+scripts+workflow plan); `## Acceptance Criteria` (AC1–AC10 `[ ]`); `## CDD Trace`.
- **Known Debt → Limitations (KEPT):** all 4 items are genuine current limitations a newcomer needs — "Package index is a flat file; no versioned API or caching layer", "No third-party package hosting story yet", "No package signing beyond SHA-256 checksums", "Version solving is exact-match only (no ranges)". Renamed section to `## Limitations`, items unchanged (already current-state phrasing).
- **Preserved:** Constraints, Challenged Assumption, Design Decisions 1–3, full Proposal §1–8 (artifact layout, loading conventions, skill activation/encapsulation, CDD worked example, index, lockfile, restore flow, dev override, commands class, discovery precedence, help/doctor), Leverage/Negative Leverage, Alternatives, Process Cost, Non-goals, **Migration** (kept — carries current fact "one hub exists (sigma)" + cutover decision).

### 5. `docs/reference/protocol/cn/MESSAGE-PACKET-TRANSPORT.md`
- **Removed stamps:** `**Issue:** #150`, `**Mode:** MCA`, `**Active Skills:** …`, `**Engineering Level:** L7`. Kept `**Version:** 3.31.0`.
- **Reframed:** `## Problem` → `## Purpose` + `### Failure modes it defends against`. The 8 failure modes are the protocol's threat model (current-state) — kept verbatim. Dropped the point-in-time Pi/Sigma equivocation incident and the `get_branch_files` mechanism narration (the mechanism is restated in the kept §"Why This Is Bulletproof Against the Current Failure Class").
- **Removed sections:** `## Impact Graph` (Downstream/Upstream tables + Copies-and-authority — every authority row is restated in Proposal §1/§4/§5, so no authority fact lost); `## 15. Acceptance Criteria` (phased `[ ]`).
- **Known Debt → Limitations (KEPT, judged):** renamed `## 16. Known Debt` → `## 15. Limitations`; kept 3 genuine current gaps (dedup index GC not yet designed; attachments/multipart not supported; packet schema `cn.packet.v2` not yet designed); dropped the 1 pure migration-process item ("legacy non-packet branches may need temporary compat during migration" — the kept §14 Migration Plan already owns rollout).
- **Renumbered** trailing sections to stay contiguous: 16 Summary → **16 Summary** (after 14 Migration Plan, 15 Limitations). No in-doc `§15/16/17` cross-refs existed (grep clean), so renumber breaks nothing.
- **Preserved:** Constraints, Challenged Assumption, full Proposal §1–13 (packet schema OCaml types, layout, transport adapters, envelope rules, signature model, dedup/equivocation, 9-step validation pipeline, materialization, send-side, error handling, traceability, Prior Art mapping, bulletproofing), §14 Migration Plan, §16 Summary.

### 6. `docs/architecture/HUB-PLACEMENT-MODELS.md`
- **Removed stamps:** `**Issue:** #156`, `**Mode:** MCA`, `**Active Skills:** …`. Kept `**Version:** 1.0.0` + the descriptive subtitle `## Standalone and Attached Hubs for Sandboxed Agents`.
- **Reframed:** `## Problem` (sandboxed-agent gap, hub-root=workspace-root incoherence, three problem classes, "why raw submodule isn't enough") → `## Purpose` stating the two roots + two modes and keeping the three problem classes as "problems the model prevents". Root-ownership lists are restated in §4 Root Semantics (kept).
- **Removed sections:** `## Impact Graph` (Downstream/Upstream/Copies-and-authority — restated in §2/§4); `## File Changes`; `## Acceptance Criteria` (10 `[ ]`).
- **Known Debt → Limitations (KEPT, judged):** renamed `## Known Debt` → `## Limitations`; kept 2 genuine limitations (multi-workspace attachments not supported; CI-specific backend quirks not handled); dropped the 1 process item ("migration of old single-root assumptions needs implementation sequencing").
- **Preserved:** Constraints (incl. Challenged assumption), full Proposal §1–10 (core decision, placement manifest JSON, attached backends, root semantics, execution model, runtime-contract impact incl. YAML, package/sync/doctor semantics, coherence rationale), Leverage/Negative Leverage, Alternatives, Process Cost, Non-goals.

---

## Sibling header-stamp sweep (stamp-only; no restructuring)

Stripped only unambiguous process-provenance stamps; left all spec content and governed markers:

- `docs/reference/packages/BUILD-AND-DIST.md` — removed `**Issue:** #219, #186`. Kept `Version`/`Status: Draft`.
- `docs/reference/protocol/cn/PROTOCOL.md` — removed `**Author:** usurobor (aka Axiom)`, `**Contributors:** Sigma`, `**Date:** 2026-02-11`. **Kept** `**Status:** Implemented (v2 — cn_protocol.ml matches this spec)` — it is a current-state code-correspondence marker, not provenance; and sibling guidance is stamp-only, no restructuring.
- `docs/reference/protocol/cn/THREAD-EVENT-MODEL.md` — removed `**Issue:** #153`. Kept `Version`/`Status`/`Purpose`/`Related`.
- `docs/reference/runtime/GIT-CN-PACKAGE.md` — removed `**Issue:** #218`. Kept `Version`/`Status`/`Parent`.
- `docs/reference/runtime/ORCHESTRATORS.md` — removed `**Issue:** #170`. Kept `Version`/`Status`/`Doc-Class`/`Canonical-Path`/`Owns`/`Does-Not-Own`. `## 1. Problem` left intact (no restructuring on siblings).

---

## Link integrity / falsification gate

- **No fragment links** to any of the 6 docs anywhere in the repo (`git grep '<basename>.md#'` → 0). **No anchor links** to any removed heading (`#acceptance-criteria|#cdd-trace|#impact-graph|#known-debt|#file-changes|#problem` across `docs/` → 0). No heading I removed was a link target.
- **No dangling prose refs**: `see (acceptance criteria|cdd trace|impact graph|known debt|file changes)` → 0; residual `AC[0-9]` prose refs in the 6 docs → 0 (the one MEMORY→#156 AC9 pointer was repaired). MESSAGE-PACKET `§15/16/17` cross-refs → 0 before renumber.
- **Residual-apparatus grep** over all 11 edited docs for `^**(Issue|Mode|Active Skills|Engineering Level|Author|Contributors|Date):` and `^## (Acceptance Criteria|CDD Trace|Impact Graph|Known Debt|File Changes)` → **clean (0 hits)**.
- **Coherence re-read** of all 6 targets: each reads header → `## Purpose` → spec body → clean tail; no doubled/orphaned `---`, no "see section above" dangles.
- **Defer-path integrity:** `git status --short` = 11 `M` under `docs/` only. No `.cdd/` (except this receipt), `.cn-sigma/`, `src/`, `schemas/`, `scripts/`, `tests/`, `.github/`, `install.sh`, `cn.json`, no root `RELEASE/CHANGELOG/ROLES/OPERATOR/SUSTAINABILITY`, no `docs/development/board/`, none of the 6 code-coupled do-not-delete docs, and NOT the CDD `design` template. Untouched.

## Judgment calls (documented)
- **`## File Changes` removed** though not literally in γ's enumerated list: it is a cycle create/edit/delete plan with zero spec substance, structurally identical to `## CDD Trace`, and "reflect only current state" forbids a future-edit plan on the reader surface. Low risk, no substance lost.
- **Rationale sections kept** (Challenged Assumption, Leverage, Alternatives, Prior Art, Non-goals, Migration, Process Cost, Tuning Strategy): these explain what each system IS / why it is shaped that way — legitimate spec content under the write skill, not point-in-time process apparatus. γ's contract named a bounded apparatus set; I stripped exactly that set (+ File Changes) and preserved everything else, honoring "when in doubt, PRESERVE substance."
- **Known Debt judged per-item:** removed entirely from MEMORY / HYBRID / POLYGLOT (every item pure future-work); kept as `## Limitations` in PACKAGE-ARTIFACTS (4), MESSAGE-PACKET (3), HUB-PLACEMENT (2) — genuine current limits a newcomer must know; dropped only the per-doc migration-sequencing/process items.

## Convergence read
**CONVERGED (docs surface).** This was the terminal reader-surface matter γ identified: the whole-doc historical-record class was exhausted in R7, and this round trims the last residual class — the within-doc CDD design-cycle apparatus on the 6 canonical specs. All 6 now read as current-state specs; the 5 sibling stamps are cleared. Every claim validates against the tree (greps clean, no defer path touched, no dangling links). No other outstanding docs-only class remains. The remaining agenda is entirely the **deferred code pass** (γ R7 §5): the 6 code-coupled do-not-delete docs + their src/schemas/scripts refs, `DESIGN-266` deletion, root `RELEASE/CHANGELOG` gate, root relocations, `docs/development/board/`, and the doctrine reconciliation (`DOCUMENTATION-SYSTEM.md §Frozen history` + the CDD `design` template that prescribes this apparatus at promotion time). None of that is α-docs scope.
