# α self-coherence — cycle/476

## §Gap

**Issue:** [cnos#476](https://github.com/usurobor/cnos/issues/476) — `cn-wake-install: v0 subcommand to render package-owned wakes into substrate (consume cn.wake-provider.v1)`. Sub 3 of cnos#467; builds on Sub 1 (cnos#468 — merged `c0048bef`) and Sub 2 (cnos#470 — merged `043bf7aa`).

**Version / mode:** `design-and-build`. No version bump in this sub (per γ scaffold pin "no version bump in this sub; γ records candidate release note in closeout"). Renderer + golden fixture + CI validator + AC2 negative tests; small design call within cycle (γ pinned both form + render target in `gamma-scaffold.md`).

**Branch:** `cycle/476` (γ-created from `origin/main@fcc5cdb9a533ad86e67524bcf05a33d2b4592e8a`; γ scaffold at `417541ad`).

**Form-choice acknowledgment.** γ pin (per `gamma-scaffold.md` §"Form choice") is **shell command** at `src/packages/cnos.core/commands/install-wake/cn-install-wake`, registered in `cnos.core/cn.package.json`'s `commands` map as `install-wake`. α **accepts** the pin. No structural reason discovered to override to Option A (Go subcommand). Specifically:

- `discover.ScanPackageCommands` (`src/go/internal/discover/discover.go`) walks `.cn/vendor/packages/*/cn.package.json` and registers shell commands via `ExecCommand`, exporting `CN_HUB_PATH`, `CN_PACKAGE_ROOT`, `CN_COMMAND_NAME`. Path confinement (`filepath.Rel` check) only enforces the entrypoint stays inside the package directory; the entrypoint may then read arbitrary paths (incl. other installed packages' `orchestrators/` trees) at runtime — no structural block.
- `jq` is present in the sandbox (`jq-1.7`) and standard in the cnos CI image (Ubuntu runners ship `jq`). `python3` + `PyYAML` are also present (3.11.15 + PyYAML 6.0.1), so YAML validation does not require new infra.
- YAML emission via shell heredoc is constrained to a fixed-shape template (`claude-code-action@v1` invocation), not arbitrary YAML serialization; escaping risk is bounded by the golden-fixture invariant (CI re-renders + byte-diffs).
- Sub 2's γ scaffold (cycle/470 §"Form choice") explicitly reserved `cn-install-wake` as Sub 3's "execute the install" verb under `commands/install-wake/`; honoring that reservation is the symmetric move.

**Render-target acknowledgment.** γ pin (per `gamma-scaffold.md` §"Render target path") is **`src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`** (per-package sibling). α **accepts** the pin. The renderer's default output path (when invoked WITHOUT `--out`) is this path; `cn install-wake agent-admin` updates the golden in place. The renderer supports `--out <path>` for the future cutover cycle.

**No γ-clarification filed.** No pinned implementation-contract axis appears unsatisfiable; no AC oracle appears unsatisfiable from the pinned form/scope.

**Refusal conditions checked against pre-coding plan:**
- Will NOT edit `.github/workflows/claude-wake.yml` (AC7 byte-identical invariant — confirmed unchanged at session start).
- Will NOT render to `.github/workflows/cnos-agent-admin.yml` (prohibited — would activate second production wake).
- Will NOT edit `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` or `prompt.md` (Sub 2 declarations; CONSUMED, not modified).
- Will NOT edit `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` (the contract; friction-log routes to gamma-closeout).
- Will NOT touch `src/go/` (γ-pinned form is shell).
- Will NOT touch any other package (cnos.cdd, cnos.cdr, cnos.kata).

## §Skills

**Tier 1 (CDD lifecycle):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role contract; §2.5 incremental self-coherence; §2.6 pre-review gate; §3.6 implementation contract is δ's

**Lifecycle sub-skills:**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — issue shape (read for AC interpretation)
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — design-and-build mode (form-pin decision)

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` — **the contract this renderer consumes**; §2.1 required fields; §2.4 substrate-rendering target; §2.5 canonical authority-split table; §3 authority split rules; §3.5 reject malformed declarations; §4.1–§4.6 verification checks
- `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` — reference instance (READ ONLY)
- `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` — prompt template (READ ONLY; inlined verbatim)
- `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` — cnos#468 contextual; renderer must not emit anything that contradicts label doctrine
- `src/packages/cnos.core/cn.package.json` — commands map (one entry added)
- `src/packages/cnos.core/commands/{daily,weekly,save}/cn-*` — shell-command sibling precedent (`#!/bin/sh`; `set -eu`; CN_HUB_PATH/CN_PACKAGE_ROOT/CN_COMMAND_NAME dispatch contract)
- `src/go/internal/discover/discover.go` — `ScanPackageCommands` + `ExecCommand` (confirms entry shape; path confinement; env-var contract)
- `.github/workflows/claude-wake.yml` — substrate-bound wake (READ ONLY; AC7 byte-identical invariant); structural reference for `claude-code-action@v1` invocation shape
- `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` — cited from the manifest's `output_contract.channel_log_convention`

**γ scaffold:** `.cdd/unreleased/476/gamma-scaffold.md` (re-read for form-choice rationale §"Form choice"; render-target rationale §"Render target path"; AC mapping table §"AC mapping"; mechanical gate block; claim-class verification injection).

**Issue body:** read in full via `mcp__github__issue_read` (issue_number=476). AC1–AC8 Invariant/Oracle/Surface blocks parsed verbatim; non-goals checked against scope plan.

**Master tracker:** cnos#467 §"Foundational architecture — package-owned wake providers (authoritative)" + wave-level AC4 — already absorbed via Sub 2 (#470) merge state and γ scaffold's distillation.

