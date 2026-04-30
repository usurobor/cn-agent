**Verdict:** APPROVED

**Round:** R1
**Fixed this round:** n/a (first review)
**origin/main SHA:** `50dac4cda4fca4592bc3d827cf5a6ffa62936e22`
**Branch head SHA:** `8c25c48b9f9b8ebf4d6307b187fcfe5482100346`
**Branch CI state:** provisional — docs-only; CI runs on push to main
**Merge instruction:** `git merge --no-ff origin/cycle/312 -m "Closes #312: docs — document agent activation (README + OPERATOR + SETUP-INSTALLER)"` into main

---

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | All three files now distinguish current operation from target runtime |
| Canonical sources/paths verified | yes | All relative links resolve (see §Link Inspection) |
| Scope/non-goals consistent | yes | Each commit touches exactly one file; no sibling-file edits |
| Constraint strata consistent | yes | Hard gates met for all three children |
| Exceptions field-specific/reasoned | n/a | No exceptions claimed |
| Path resolution base explicit | yes | Relative links verified from each file's own directory |
| Proof shape adequate | yes | Manual inspection + grep; AC evidence tied to commit SHAs |
| Cross-surface projections updated | yes | README, OPERATOR, SETUP-INSTALLER form a coherent set |
| No witness theater / false closure | yes | α provided concrete commit SHAs per AC |
| Self-coherence matches branch files | yes | All evidence commit SHAs verified in git log |

---

## §2.0 Issue Contract

### AC Coverage — #313 README.md

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| AC1 | Activation step after health-check block | yes | met | `### Activate the agent` section added after `cn doctor`/`cn status` block; commit `4ba4cc1c` |
| AC2 | Distinguishes hub preparation from activation | yes | met | Opens with "`cn setup` prepares the hub. The hub becomes an active coherent agent when a capable model is pointed at it." Two distinct answers present |
| AC3 | Safe placeholder, no angle brackets, renders cleanly | yes | met | Example uses `HUB_DIR`; single fenced `sh` block; no `cn agent` or angle-bracket placeholder |

### AC Coverage — #314 OPERATOR.md

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| AC1 | First content section is "Current operation" | yes | met | `## Current operation` is the first non-frontmatter heading; commit `32063b19` |
| AC2 | Daemon/systemd under explicit target/planned section | yes | met | `## Target runtime (planned)` heading introduced; daemon/systemd content moved under it |
| AC3 | No unannotated `cn agent --daemon`/systemd outside target section | yes | met | Two occurrences survive outside target section; both annotated: line 146 `(target — when daemon is running)`, line 232 `(planned)` |
| AC4 | Cross-references resolve and are target-labeled | yes | met | See §Link Inspection below |

### AC Coverage — #315 SETUP-INSTALLER.md

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| AC1 | Distinguishes install / setup / activation / target service wiring | yes | met | `## 5. Phases` table names all four phases with status; commit `0a592d15` |
| AC2 | Current activation path named explicitly | yes | met | Phases table "activate" row names Claude UI / Claude Code CLI; paragraph below confirms |
| AC3 | Daemon/service wiring labeled target/planned | partial | met with residual | Operator-instruction content fully labeled; pre-existing context sections (§1 Problem, §2 Goals, §6 Architecture §5.2, §10 Implementation plan, §11 Open questions) carry unlabeled `systemd` references — these are design rationale, not operator instructions; see Notes |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `README.md` | yes | ✓ | Child #313 commit `4ba4cc1c` |
| `OPERATOR.md` | yes | ✓ | Child #314 commit `32063b19` |
| `docs/alpha/cli/SETUP-INSTALLER.md` | yes | ✓ | Child #315 commit `0a592d15` |
| Sibling files (non-goal) | no edits | ✓ | Confirmed per-commit file list |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `self-coherence.md` | yes | yes | Full gap/skills/ACs/self-check/debt/trace; review-readiness R1 |
| `beta-review.md` | yes | in progress | This file |
| Tests | no | n/a | Docs-only; explicitly stated |
| Code diff | no | n/a | Docs-only; no runtime implementation |

### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| `cnos.core/skills/write` | #313, #314, #315 Tier 3 | yes (α) | yes | Prose follows status-truth discipline; front-loads the point |
| `cnos.core/skills/design` | #313, #314, #315 Tier 3 | yes (α) | yes | Phase boundaries (install/setup/activation/target) preserved across all three files |
| `cdd/review` | #313, #314, #315 Tier 3 | yes (α) | yes | No future-as-present; status truth maintained |

---

## §Link Inspection (required per #314 AC4)

### README.md — all relative links

| Link | Resolves? |
|------|-----------|
| `./docs/alpha/agent-runtime/GO-KERNEL-COMMANDS.md` | ✓ |
| `./docs/alpha/architecture/INVARIANTS.md` | ✓ |
| `./docs/alpha/agent-runtime/AGENT-RUNTIME.md` | ✓ |
| `./docs/README.md` | ✓ |
| `./docs/THESIS.md` | ✓ |
| `./docs/alpha/essays/COHERENCE-SYSTEM.md` | ✓ |
| `./docs/alpha/essays/FOUNDATIONS.md` | ✓ |
| `./docs/alpha/essays/MANIFESTO.md` | ✓ |
| `./docs/alpha/doctrine/README.md` | ✓ |
| `./docs/alpha/doctrine/coherence-for-agents/COHERENCE-FOR-AGENTS.md` | ✓ |
| `./docs/alpha/doctrine/ethics-for-agents/ETHICS-FOR-AGENTS.md` | ✓ |
| `./docs/alpha/doctrine/judgment-for-agents/JUDGMENT-FOR-AGENTS.md` | ✓ |
| `./docs/alpha/doctrine/inheritance-for-agents/INHERITANCE-FOR-AGENTS.md` | ✓ |
| `./docs/alpha/protocol/WHITEPAPER.md` | ✓ |
| `./docs/alpha/protocol/PROTOCOL.md` | ✓ |
| `./docs/alpha/security/SECURITY-MODEL.md` | ✓ |
| `./docs/alpha/agent-runtime/CAA.md` | ✓ |
| `./docs/beta/architecture/ARCHITECTURE.md` | ✓ |
| `./docs/gamma/cdd/CDD.md` | ✓ |
| `./docs/gamma/ENGINEERING-LEVELS.md` | ✓ |
| `./CHANGELOG.md` | ✓ |
| `./docs/alpha/ctb/README.md` | ✓ |
| `./docs/alpha/ctb/LANGUAGE-SPEC.md` | ✓ |
| `./docs/alpha/ctb/LANGUAGE-SPEC-v0.2-draft.md` | ✓ |
| `./docs/alpha/ctb/SEMANTICS-NOTES.md` | ✓ |
| `./docs/alpha/ctb/CTB-v4.0.0-VISION.md` | ✓ |
| `./CONTRIBUTING.md` | ✓ |
| `./docs/beta/SUSTAINABILITY.md` | ✓ |
| `./LICENSE` | ✓ |

No new links were added or changed in the README.md diff; all pre-existing links verified green.

### OPERATOR.md — all relative links

| Link | Resolves? | Notes |
|------|-----------|-------|
| `README.md` | ✓ | |
| `docs/alpha/cli/SETUP-INSTALLER.md` | ✓ | |
| `docs/beta/guides/AUTOMATION.md` (×2) | ✓ | Target runtime references; prose marks as target |
| `docs/beta/guides/BUILD-RELEASE.md` | ✓ | |
| `docs/alpha/package-system/PACKAGE-SYSTEM.md` | ✓ | |
| `docs/beta/guides/HANDSHAKE.md` | ✓ | |
| `docs/beta/guides/TROUBLESHOOTING.md` | ✓ | |
| `docs/alpha/cli/CLI.md` | ✓ | |

No new links added by the OPERATOR.md diff; all pre-existing links verified green. Cross-references to `AUTOMATION.md` survive in the Target runtime section with prose marking them as target/planned — AC4 met.

### SETUP-INSTALLER.md — all relative links

| Link | Resolves? | Notes |
|------|-----------|-------|
| `../agent-runtime/AGENT-RUNTIME.md` → `docs/alpha/agent-runtime/AGENT-RUNTIME.md` | ✓ | |
| `../security/SECURITY-MODEL.md` → `docs/alpha/security/SECURITY-MODEL.md` | ✓ | |

No new links added by the SETUP-INSTALLER.md diff; all pre-existing links verified green.

---

## Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|
| — | No findings | — | — | — |

---

## Notes

**N1 — SETUP-INSTALLER.md: pre-existing context-section `systemd` mentions not individually labeled**

Five `systemd` references survive outside individually-labeled target/planned sections:
- Line 16 (§1 Problem): `write systemd units (or cron scripts)` — describes historical manual setup pain
- Line 39 (§2 Goals): `Portable: works without systemd` — design portability goal
- Line 94 (§6 Architecture §5.2): `systemd unit (optional)` — artifact table entry (acknowledged in α's §Debt)
- Line 719 (§10 Implementation plan): Step C, daemon/systemd detection
- Line 732 (§11 Open questions): non-systemd platform support

These are in design-doc context sections (Problem, Goals, Architecture, Implementation plan, Open questions), not operator instructions. None claim `cn agent --daemon` or systemd as current shipped operation. The §5 Phases table is the authoritative status-truth frame; these sections are design rationale governed by that context. No operator misdirection is present. Strict AC3 oracle (each occurrence labeled) technically fires, but per review §3.5 (no phantom blockers), these don't constitute demonstrable incoherence for an operator or reader.

The whole document is `Status: Draft` and these sections describe design motivation and planning. A follow-on pass to add inline labels to these context sections would be appropriate editorial hygiene but is not required for merge.

**N2 — OPERATOR.md: `cn-<name>` angle-bracket token in pre-existing code blocks**

`systemctl start cn-<name>` and `systemctl stop cn-<name>` appear in the Target runtime section's Start/Stop code block. These use `<name>` as a template variable indicating the hub name. This is pre-existing content (unchanged by the diff); it is inside the Target runtime section; no new angle-bracket placeholders were introduced by this cycle. The issue's "HUB_DIR token" constraint applies to the new activation examples, which correctly use `HUB_DIR`. Not a finding.

**N3 — `cn logs` in OPERATOR.md not labeled as target/planned**

`cn logs` (used in §1 Observing) is not in the README's shipped command list. α noted this in the debt section. Not in scope for #312 children; labeling `cn logs` as target/planned is a reasonable follow-on but was not part of the AC surface for this cycle. Not a finding.

---

## Regressions Required (D-level only)

None.

---

## Approval statement

R1 closes the search space on:
- **README.md (#313):** activation step present, correct placeholder, clean rendering
- **OPERATOR.md (#314):** current operation first, daemon/systemd content in labeled target section, all cross-references resolve
- **SETUP-INSTALLER.md (#315):** §5 Phases table distinguishes all four phases, activation path explicit, operator-instruction-level daemon content labeled target
- **Non-goals:** no sibling-file edits; no runtime implementation; no deletion of target content
- **Links:** all relative links in all three files resolve against `origin/cycle/312`

No blocker was found in the relevant contract. The branch is coherent and ready to merge.
