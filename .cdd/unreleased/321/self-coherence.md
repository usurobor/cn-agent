## Gap

**Issue:** #321 — core(activate): introduce Kernel/Persona/Operator activation triad and retire templates class

**Version/mode:** MCA — system-shaping change across package layout, pkg code, activate command, docs, and kata.

**Incoherence:** `cn activate` 3.70.0 emits a structurally wrong `## Identity` bucket that conflates kernel, persona, and operator into one unstructured list. On `cn-sigma` (the canonical reference hub), it reports `no identity files found`. The kernel ships under `templates/SOUL.md` as a "template" in both docs and code, entrenching a fork-the-kernel pattern. `pkg.ContentClasses` and `PACKAGE-SYSTEM.md` both list `templates` as a content class. Per-hub identity slots (`Name/Role/Operator`) live in the kernel file, contradicting kernel semantics. `scanPackages` reads only the vendor directory, not `deps.json`. No ordered first-read guidance, no latest-reflection pointer.

**Selected gap:** The triad (Kernel / Persona / Operator) is the correct layering. Each layer owns one canonical path. The kernel ships under `doctrine/` (not `templates/`). The `templates` content class is removed from code and docs. `cn activate` produces three layered sections, three kernel states, three deps states, an ordered `## Read first`, and a latest-reflection pointer.

---

## Skills

**Tier 1:**
- `CDD.md` (canonical lifecycle)
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (α role)

**Tier 2 (loaded per issue):**
- `src/packages/cnos.core/skills/write/SKILL.md` — KERNEL.md reframing and PACKAGE-SYSTEM.md updates
- `src/packages/cnos.core/skills/design/SKILL.md` — layer conflation is a design failure; rename + package code + activate redesign coupled by design

**Tier 3 (issue-named):**
- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — Go implementation of resolvers, `pkg.ContentClasses` edit
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` — sigma-shape / init-only / init+setup fixtures; positive + negative proofs
- `src/packages/cnos.eng/skills/eng/tool/SKILL.md` — stdout/stderr contract preserved
- `src/packages/cnos.eng/skills/eng/ux-cli/SKILL.md` — `## Read first` ordering, empty-state messaging, restore guidance wording

---

## ACs

**AC1: `scanIdentity` and the conflated `## Identity` bucket are removed**
- **Met.** `scanIdentity` function deleted from `activate.go`. `writePrompt` no longer emits `## Identity`. Oracle: `grep -n "scanIdentity\|## Identity\b" src/go/internal/activate/activate.go` → 0 hits. Tests: `TestNoIdentityBucket` across sigma/init-only/init+setup fixtures.

**AC2: Three layered resolvers exist with correct semantics**
- **Met.** `resolveKernel`, `scanPersona`, `scanOperator` each own one layer and one canonical path. Unit tests on each resolver via `TestKernelState_*`, `TestResolversRejectLegacyPaths`. Sigma-shape fixture: kernel resolves to vendored `doctrine/KERNEL.md`, persona to `spec/PERSONA.md`, operator to `spec/OPERATOR.md`. Negative: `TestResolversRejectLegacyPaths` asserts `spec/SOUL.md`, `spec/identity.md`, `agent/identity.md`, `spec/USER.md` not surfaced.

**AC3: Prompt has three layered sections with correct headers**
- **Met.** `writePrompt` emits `## Kernel`, `## Persona`, `## Operator` in that order. `TestSigmaShapeActivatesCorrectly` asserts all three present and in correct ordinal position. `TestInitOnlyFixture_NoKernelReference` confirms init-only also has all three headers. Negative: `TestNoIdentityBucket` confirms `## Identity` absent.

**AC4: Kernel renamed, relocated, and reframed**
- **Met.** `git mv src/packages/cnos.core/templates/SOUL.md → src/packages/cnos.core/doctrine/KERNEL.md`. Title changed to `# Kernel`. Lead description reframes as universal coherence kernel. Per-hub identity slots (`Name/Role/Operator`, "Configure identity through the configure-agent skill") removed; replaced with pointer to `spec/PERSONA.md` and `spec/OPERATOR.md`. Oracle: file exists at new path; old path absent (templates/ deleted); `grep "set by operator\|Configure identity" src/packages/cnos.core/doctrine/KERNEL.md` → 0 hits.

**AC5: `templates/` directory and dead `USER.md` are deleted**
- **Met.** `git rm src/packages/cnos.core/templates/USER.md`. Directory is empty and removed. Oracle: `ls src/packages/cnos.core/templates/` → no such directory.

**AC6: `templates` removed from package content classes — in docs and code**
- **Met.** `pkg.ContentClasses` now has 7 entries; `"templates"` absent. `PACKAGE-SYSTEM.md` §1.1 and §4 list 7 classes; §4.2 (Why templates are a content class) deleted. `hubstatus_test.go` `TestRunContentClassesAllSeven` (was AllEight) updated. `build_test.go` templates fixtures replaced with `doctrine/KERNEL.md`. Oracle: `grep templates src/go/internal/pkg/pkg.go` → 0 hits on ContentClasses; `go test ./internal/pkg/... ./internal/hubstatus/... ./internal/pkgbuild/...` → all pass.

**AC7: Kernel resolution distinguishes three states**
- **Met.** `resolveKernel` returns `{state: "vendored", path, version}` / `{state: "manifest-only"}` / `{state: "none"}`. Tests: `TestKernelState_Vendored`, `TestKernelState_ManifestOnly`, `TestKernelState_None`. Each fixture exercises the corresponding state. Negative: manifest-only does not print a resolved path; no-manifest does not print restore guidance.

**AC8: Deps detection distinguishes three states**
- **Met.** `scanDeps` returns `{state: "restored", packages}` / `{state: "manifest-only"}` / `{state: "none"}`. Tests: `TestDepsState_Restored`, `TestDepsState_ManifestOnly`, `TestDepsState_None`.

**AC9: Ordered `## Read first` section**
- **Met.** `writePrompt` emits `## Read first` with entries: persona (1) → operator (2) → kernel (3) → deps (4) → latest reflection (5, when present). Tests: `TestReadFirstSection_OrderedSigma` (ordinal positions verified), `TestReadFirstSection_InitOnlyOrdered`. Negative: "inspect files directly" footer absent from new prompt.

**AC10: Latest reflection pointer (path only)**
- **Met.** `latestReflection` reads `threads/reflections/daily/`, picks lexically last non-directory entry (YYYY-MM-DD.md → chronologically latest). Tests: `TestLatestReflection_Present` (two daily files, latest path appears), `TestLatestReflection_Empty` (empty dir → no "not present" text, line omitted).

**AC11: Sigma-shape fixture activates correctly**
- **Met.** `TestSigmaShapeActivatesCorrectly`: sigma fixture produces `## Kernel` with vendored path `@3.71.0`, `## Persona` with `spec/PERSONA.md`, `## Operator` with `spec/OPERATOR.md`, `## Read first`. Negatives: no `no identity files found`, no `spec/SOUL.md`.

**AC12a: Init-only fixture activates with explicit absence and no kernel reference**
- **Met.** `TestInitOnlyFixture_NoKernelReference`: kernel = `no kernel reference`, no restore guidance, persona absent explicit, operator absent explicit, `spec/SOUL.md` not surfaced.

**AC12b: Init+setup fixture activates with restore guidance**
- **Met.** `TestInitPlusSetupFixture_RestoreGuidance`: kernel = `dependency manifest declares cnos.core; not restored — run cn deps restore`, no `vendored at`, persona absent, operator absent, `spec/SOUL.md` not surfaced.

**AC13: R5 kata extended**
- **Met.** `src/packages/cnos.kata/katas/R5-activate/kata.md` adds P7 (triad split), P8 (kernel states), P9 (deps states), P10 (read-first ordering), P11 (latest reflection). P1–P6 unchanged.

**AC14: No regression**
- **Met.** `go test ./...` → all 12 packages green (activate, activation, binupdate, cli, discover, doctor, hubinit, hubsetup, hubstatus, pkg, pkgbuild, restore). R5 P1–P6 preserved (no existing test assertions weakened). Stdout/stderr contract preserved (diagnostic arrow still goes to stderr, prompt to stdout). No model invocation.

**AC15: No stale references to old paths in active code or active docs**
- **Met with notes.** `grep -rn "templates/SOUL\|templates/USER" src/ docs/alpha/ docs/beta/ CHANGELOG.md README.md` produces three hits:
  1. `CHANGELOG.md:140,143` — the new 3.71.0 entry I wrote documenting the removal ("was moved from X to Y"). These are intentional release notes, not stale references.
  2. `docs/alpha/cli/SETUP-INSTALLER.md:424` — references `src/agent/templates/SOUL.md`, which is a path in an agent source directory tree (not `cnos.core/templates/SOUL.md`). This is a different path describing a different surface (`cn setup`'s template reading behavior from an agent source repo, not the cnos.core package kernel). The oracle `grep "templates/SOUL"` is broader than the stated invariant "No reference to `cnos.core/templates/SOUL.md`"; this match is not a stale reference to the removed path.
- No active code or active doc references `cnos.core/templates/SOUL.md` or `cnos.core/templates/USER.md`. Historical references in `docs/gamma/cdd/3.24.0/` preserved.
