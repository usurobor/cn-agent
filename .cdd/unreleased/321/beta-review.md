## §2.0.0 Contract Integrity

**Round:** 1
**origin/main SHA (review base):** f9843317d6be65b3b66741e0349f8f8df88d7630
**cycle/321 head SHA:** c8206956353f227383a5f888173a31d2968d989f

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | Shipped (3.70.0) vs draft target vs non-goals clearly separated throughout issue and self-coherence.md. No claim of runtime enforcement that doesn't exist. |
| Canonical sources/paths verified | yes | `pkg.ContentClasses` at `pkg.go:131–138`; `activate.go`; `PACKAGE-SYSTEM.md`; `doctrine/KERNEL.md` — all verified against branch state. |
| Scope/non-goals consistent | yes | `cn init`, `cn setup`, sigma rename, wizard, live-model dogfood all explicitly non-goals. No out-of-scope behavior in diff. |
| Constraint strata consistent | yes | Hard gates (three sections, no `## Identity`, `pkg.ContentClasses` 7 entries, `templates/` absent, KERNEL.md no per-hub slots) all met in diff. Exception-backed and optional fields correctly classified. |
| Exceptions field-specific/reasoned | n/a | No exception-backed fields in this cycle's output. |
| Path resolution base explicit | yes | vendored kernel path = `.cn/vendor/packages/cnos.core/doctrine/KERNEL.md`; persona = `spec/PERSONA.md`; operator = `spec/OPERATOR.md`; all hub-root-relative. Explicit in issue and code. |
| Proof shape adequate | yes | Each AC has invariant, oracle, positive case, negative case. Unit + integration tests in `activate_test.go`. Filesystem inspection for templates/ deletion. Grep oracles for stale-ref check. |
| Cross-surface projections updated | yes | `pkg.ContentClasses` (code) ↔ `PACKAGE-SYSTEM.md` (docs) both updated to 7 classes. `hubstatus_test.go` and `build_test.go` updated. R5 kata extended. CHANGELOG entry added. |
| No witness theater / false closure | yes | Tests run and verify behavior. `go test ./...` 12/12 green (verified non-cached). No claim of machine enforcement beyond what tests demonstrate. |
| PR body matches branch files | n/a | CDD does not use PRs. `self-coherence.md` accurately describes branch state at head `c8206956`. |

---

## §2.0 Issue Contract Walk

### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| AC1 | `scanIdentity` and `## Identity` bucket removed | yes | met | `scanIdentity` deleted from `activate.go`. `writePrompt` emits no `## Identity`. `grep` oracle verifies 0 hits. `TestNoIdentityBucket` across three fixtures. |
| AC2 | Three layered resolvers with correct semantics | yes | met | `resolveKernel`, `scanPersona`, `scanOperator` each own one canonical path. `TestResolversRejectLegacyPaths` verifies `spec/SOUL.md`, `spec/identity.md`, `agent/identity.md`, `spec/USER.md` absent from output. |
| AC3 | Prompt has three layered sections in correct order | yes | met | `writePrompt` emits `## Kernel`, `## Persona`, `## Operator`. `TestSigmaShapeActivatesCorrectly` asserts all three present and in correct ordinal position. |
| AC4 | Kernel renamed, relocated, and reframed | yes | met | `git mv templates/SOUL.md → doctrine/KERNEL.md`. Title `# Kernel`, lead description reframed. Per-hub slots (`Name/Role/Operator`, configure-agent framing) removed. `grep "set by operator\|Configure identity"` → 0 hits confirmed by β inspection. |
| AC5 | `templates/` directory and `USER.md` deleted | yes | met | `templates/` absent from `src/packages/cnos.core/`. `USER.md` removed. Filesystem verified. |
| AC6 | `templates` removed from content classes — code and docs | yes | met | `pkg.ContentClasses` 7 entries, `"templates"` absent (comment explains why). `PACKAGE-SYSTEM.md` §1.1 and §4 list 7 classes. `hubstatus_test.go` `TestRunContentClassesAllSeven` updated. `build_test.go` fixtures updated. |
| AC7 | Kernel resolution distinguishes three states | yes | met | `resolveKernel` returns `{state: "vendored"/"manifest-only"/"none"}`. `TestKernelState_Vendored/ManifestOnly/None` cover each. |
| AC8 | Deps detection distinguishes three states | yes | met | `scanDeps` returns `{state: "restored"/"manifest-only"/"none"}`. `TestDepsState_*` cover each. |
| AC9 | Ordered `## Read first` section | yes | met | `writePrompt` emits `## Read first` persona(1)→operator(2)→kernel(3)→deps(4)→reflection(5). `TestReadFirstSection_OrderedSigma` verifies ordinal positions. `TestReadFirstSection_InitOnlyOrdered` verifies ordering for absent layers. |
| AC10 | Latest reflection pointer (path only) | yes | met | `latestReflection` picks lexically-last non-directory entry in `threads/reflections/daily/`. `TestLatestReflection_Present` verifies latest (not older) file. `TestLatestReflection_Empty` verifies omission when empty. |
| AC11 | Sigma-shape fixture activates correctly | yes | met | `TestSigmaShapeActivatesCorrectly` asserts `## Kernel` vendored at `doctrine/KERNEL.md@3.71.0`, `## Persona` with `spec/PERSONA.md`, `## Operator` with `spec/OPERATOR.md`. Negatives: no `no identity files found`, no `spec/SOUL.md`. |
| AC12a | Init-only fixture: no kernel reference | yes | met | `TestInitOnlyFixture_NoKernelReference`: kernel = `no kernel reference`; no restore guidance; PERSONA/OPERATOR absent explicit; `spec/SOUL.md` not surfaced. |
| AC12b | Init+setup fixture: restore guidance | yes | met | `TestInitPlusSetupFixture_RestoreGuidance`: kernel = manifest-only state with restore guidance; no `vendored at`; PERSONA/OPERATOR absent; `spec/SOUL.md` not surfaced. |
| AC13 | R5 kata extended P7–P11 | yes | met | `R5-activate/kata.md` adds P7 (triad split), P8 (kernel states ×3), P9 (deps states ×3), P10 (read-first ordering), P11 (reflection ×2). P1–P6 unchanged. Inputs section updated with fixture shapes. |
| AC14 | No regression | yes | met | `go test ./...` 12/12 green non-cached. R5 P1–P6 assertions unchanged. Stdout/stderr contract preserved. No model invocation. |
| AC15 | No stale references in active surfaces | yes | met with noted nuance | `grep -rn "templates/SOUL\|templates/USER"` hits: CHANGELOG:140,143 (intentional release notes), `docs/alpha/cli/SETUP-INSTALLER.md:424` references `src/agent/templates/SOUL.md` (different path: agent source tree, not cnos.core package). Issue oracle specifies `cnos.core/templates/SOUL.md` and `cnos.core/templates/USER.md` — neither match in active surfaces. Assessment: AC15 fully met. |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `src/go/internal/activate/activate.go` | yes | updated | Resolvers, `writePrompt`, all per spec |
| `src/go/internal/activate/activate_test.go` | yes | updated | New fixtures and 20+ new test functions |
| `src/go/internal/pkg/pkg.go` | yes | updated | `ContentClasses` reduced to 7, comment added |
| `src/go/internal/pkg/pkg_test.go` | n/a | no change needed | `pkg_test.go` not in diff; `pkgbuild/build_test.go` and `hubstatus_test.go` updated per AC6 |
| `src/go/internal/hubstatus/hubstatus_test.go` | yes | updated | `TestRunContentClassesAllSeven` (was AllEight) |
| `src/go/internal/pkgbuild/build_test.go` | yes | updated | templates→doctrine fixtures in two tests |
| `src/packages/cnos.core/doctrine/KERNEL.md` | yes | created (via mv) | Kernel-only content, no per-hub slots |
| `src/packages/cnos.core/templates/USER.md` | yes | deleted | |
| `docs/alpha/package-system/PACKAGE-SYSTEM.md` | yes | updated | §1.1 table, §4.1 count, §4.2 removed, old §4.3→§4.2, §4.4 N=6→N=7, history row added |
| `src/packages/cnos.kata/katas/R5-activate/kata.md` | yes | updated | P7–P11 added |
| `CHANGELOG.md` | yes | updated | 3.71.0 detailed section (Added/Changed/Fixed/Removed) |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `self-coherence.md` | yes | yes | Complete: gap, mode, all skills, AC mapping (1–15), CDD Trace through step 7a, debt (4 items), pre-review gate (11 rows), review-readiness signal |
| `beta-review.md` | yes | in progress | This document |
| `alpha-closeout.md` | post-merge | pending | α writes after merge |
| `beta-closeout.md` | post-merge | pending | β writes after merge |
| Bootstrap version dir | no | n/a | Small-change bootstrap exemption claimed; this cycle adds to 3.71.0 release, not a new version dir cycle |

### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| CDD.md | always | yes | yes | Lifecycle followed throughout |
| `alpha/SKILL.md` | always | yes | yes | Pre-review gate completed (11 rows), review-readiness signal present |
| `cnos.core/skills/write/SKILL.md` | issue | yes | yes | KERNEL.md reframing, PACKAGE-SYSTEM.md updates — prose is concise and non-redundant |
| `cnos.core/skills/design/SKILL.md` | issue | yes | yes | Layer conflation is a design failure: coupled rename + code + activate redesign shipped as one unit. One-file-one-layer principle applied. |
| `eng/go/SKILL.md` | issue | yes | yes | `resolveKernel`/`scanPersona`/`scanOperator`/`scanDeps`/`latestReflection` follow Parse/IO split; `pkg.ContentClasses` edit correct |
| `eng/test/SKILL.md` | issue | yes | yes | sigma-shape / init-only / init+setup / kernel-state×3 / deps-state×3 / read-first×2 / reflection×2 / legacy-rejection / no-identity-bucket×3 fixtures; positive + negative per AC |
| `eng/tool/SKILL.md` | issue | yes | yes | stdout-only prompt, stderr diagnostics — `TestRunPositive_StdoutOnly` verifies the split |
| `eng/ux-cli/SKILL.md` | issue | yes | yes | `## Read first` ordering, explicit absence messaging, restore guidance wording |

---

## §2.1 Diff and Context Inspection

### 2.1.1 Structural closure

Three resolvers, three exclusive states each. No input source bypasses state detection:

- `resolveKernel`: checks vendored KERNEL.md first; falls back to deps.json cnos.core declaration; falls back to "none". Ordering correct — vendored is the authoritative state. No gap where both vendored and manifest-only could simultaneously surface.
- `scanPersona`, `scanOperator`: each reads exactly one path, returns empty string on absence. No fallback chain, no legacy paths accepted.
- `scanDeps`: reads `deps.json` existence first; then counts vendor/packages/ entries. Three exclusive states. `TestResolversRejectLegacyPaths` verifies none of `spec/SOUL.md`, `spec/identity.md`, `agent/identity.md`, `spec/USER.md` surface under any layered section.
- `latestReflection`: `os.ReadDir` is sorted by name; reverse iteration picks the lexically-last non-directory. For `YYYY-MM-DD.md` naming, lexical == chronological. Correct.

### 2.1.2 Multi-format semantic parity

Key claim about 7 content classes appears in:
- `pkg.ContentClasses` (code): 7 entries, `"templates"` absent ✓
- `PACKAGE-SYSTEM.md` §1.1 table: 7 rows ✓
- `PACKAGE-SYSTEM.md` "All seven content classes": ✓
- `PACKAGE-SYSTEM.md` §4.1 "Current scale: 7": ✓
- `PACKAGE-SYSTEM.md` §4.2 history row for v3.71.0: ✓
- `hubstatus_test.go` `TestRunContentClassesAllSeven` want string: 7 entries ✓
- All agree. No format divergence.

### 2.1.3 Snapshot consistency

No snapshot/expect tests (ppx_expect not in use in Go packages). Go tests use direct string assertion. Not applicable.

### 2.1.4 Stale-path validation

Old paths audited:
- `cnos.core/templates/SOUL.md`: only hits are CHANGELOG release notes (intentional) and `SETUP-INSTALLER.md:424` (`src/agent/templates/SOUL.md` — different surface, not cnos.core). ✓
- `cnos.core/templates/USER.md`: only hits are CHANGELOG release notes. ✓
- `scanIdentity` function: comment in test at line 407 (`// --- AC1: scanIdentity / ## Identity removed ---`). No callers, no definition. ✓
- `scanPackages` function: fully absent from codebase. ✓
- `## Packages and skills` section header: replaced by `## Dependencies` in `writePrompt`. ✓
- `## Identity` section header: fully removed. ✓

### 2.1.5 Branch naming

Branch `cycle/321` — canonical format per CDD.md §4.2. ✓

### 2.1.6 Execution timeline

`activate.go` runs in a single process, reads filesystem at execution time. All state reads are synchronous. No cross-process boundary. Not applicable.

### 2.1.7 Derivation vs validation

`pkg.ContentClasses` is the single source of truth for the class list. `pkgbuild.FindContentClasses` uses it dynamically — no duplication. Removing `"templates"` from the slice propagates automatically to all consumers. ✓

### 2.1.8 Authority-surface conflict

Checked agreement across:
- Code authority (`pkg.ContentClasses`): 7 entries, no templates ✓
- Docs authority (`PACKAGE-SYSTEM.md`): 7 classes, no templates ✓
- Test authority (`hubstatus_test.go`, `build_test.go`): updated to 7 ✓
- Kernel content authority (`doctrine/KERNEL.md`): no per-hub slots, no configure-agent framing ✓
- Issue ACs ↔ self-coherence.md ↔ diff: consistent throughout ✓

No authority-surface conflict found.

### 2.1.9 Module-truth audit

Touched module: `activate.go`. Checked for same kind of assumption (`## Identity` / flat bag pattern): entirely absent. The old `writePrompt` signature `(w, absPath, cfg, identity, packages, memory, threads []string)` is replaced by `(w, absPath, cfg, kernel kernelState, persona, operator string, deps depsState, latest string, memory, threads []string)`. No remaining caller passes the old signature (compilation would fail). ✓

Checked `pkg.go` module for other uses of the removed `"templates"` string: comment explains why it's intentionally absent. No other ContentClasses assignment found. ✓

### 2.1.10 Contract-implementation confinement

- `resolveKernel`: vendor check uses `os.Stat` on exactly one path. Returns on first match. Cannot accept `spec/SOUL.md` (different directory). ✓
- `scanPersona`: reads `filepath.Join(hubPath, "spec", "PERSONA.md")` only. No fallback. ✓
- `scanOperator`: reads `filepath.Join(hubPath, "spec", "OPERATOR.md")` only. No fallback. ✓

### 2.1.11 Architecture leverage

This change makes the fork-the-kernel anti-pattern structurally impossible: the kernel is now read from the vendored package, not from `spec/`. The `## Identity` bucket that conflated three orthogonal layers is eliminated. The `templates` content class that entrenched per-hub editing semantics is removed at both code and docs level. L6 system-safe; potentially L7 in that the entire templates/identity-conflation friction class is eliminated.

### 2.1.12 Process overhead

R5 kata P7–P11 additions are verifiable and operator-facing, documenting the new triad behavior with concrete fixture descriptions. Reasonable overhead for a behavioral change of this scope.

### 2.1.13 Design constraints check

Issue's active design constraints verified:
| Constraint | Status |
|---|---|
| One file, one layer | ✓ — each resolver owns one canonical path, no fallbacks |
| Kernel is referenced, not forked | ✓ — reads from vendored package; `spec/` not consulted for kernel |
| Doctrine is doctrine | ✓ — `doctrine/KERNEL.md`, not `templates/` |
| Code and docs must agree | ✓ — both updated, tests verify |
| Sigma is the canonical reference hub | ✓ — sigma-shape fixture is the primary positive test |
| Activate stays prompt-only | ✓ — no model invocation, `TestRunPositive_NoModelInvocation` verifies |
| Status truth in the prompt | ✓ — all three sections emit explicit absence rather than silent omission |
| Boundary discipline | ✓ — `cn init`, `cn setup`, sigma repo untouched |
| Naming is symmetric | ✓ — KERNEL/PERSONA/OPERATOR in file names, section headers, and resolver names |

---

## Architecture Check

| Check | Result | Notes |
|---|---|---|
| Reason to change preserved | yes | `activate.go` owns activation prompt generation only. `pkg.go` owns content class definitions only. No conflated concerns introduced. |
| Policy above detail preserved | yes | Kernel (foundational doctrine) lives in `doctrine/`; per-hub identity lives in `spec/`. Policy layer explicit in KERNEL.md body. |
| Interfaces remain truthful | yes | `resolveKernel`, `scanPersona`, `scanOperator` each return exactly what they promise: one state per layer, no overclaiming. |
| Registry model remains unified | n/a | No registry/descriptor pattern in this change. |
| Source/artifact/installed boundary preserved | yes | `cnos.core/doctrine/KERNEL.md` (source) → dist tarball (artifact) → `.cn/vendor/packages/cnos.core/doctrine/KERNEL.md` (installed) → `cn activate` reads from installed. Flow explicit and consistent with PACKAGE-SYSTEM.md §2 diagram (updated). |
| Runtime surfaces remain distinct | yes | `activate.go` (prompt generation), `pkg.go` (content class enumeration), `pkgbuild/` (build-time discovery) remain separate. No surface conflation. |
| Degraded paths visible and testable | yes | All three kernel states explicitly named. All three deps states explicitly named. Empty-state messaging names the missing surface. `TestKernelState_None` / `TestDepsState_None` verify absence paths. No silent fallback. |

---

## Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|
| — | No findings | All ACs met, no authority conflicts, no stale paths in active surfaces, all tests green | — | — |
