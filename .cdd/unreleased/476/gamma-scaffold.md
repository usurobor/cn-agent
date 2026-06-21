# γ scaffold — cycle/476

**Issue:** [cnos#476](https://github.com/usurobor/cnos/issues/476) — `cn-wake-install: v0 subcommand to render package-owned wakes into substrate (consume cn.wake-provider.v1)`. Sub 3 of cnos#467 (master tracker `agent/wake-orchestration`); builds on Sub 1 (cnos#468 — label doctrine, merged at `c0048bef`) and Sub 2 (cnos#470 — agent-admin wake provider, merged at `043bf7aa`). Implements the renderer scope of cnos#450 (amended).

**Mode:** **design-and-build** (per #476 mode declaration). v0 renderer + golden fixture + validator; small design call within cycle (form pin + render-target pin per below). γ pins both; α may override only with a recorded structural reason in `self-coherence.md` §Gap.

**Branch:** `cycle/476` from `origin/main` at `fcc5cdb9a533ad86e67524bcf05a33d2b4592e8a` (verified by `git rev-parse origin/main` at γ scaffold time; HEAD of `cycle/476` is identical because γ created the branch from `origin/main` without intervening commits).

**Cycle execution mode:** **bootstrap (pre-dispatch δ/channel path).** The γ-interface session acts as bootstrap-δ; γ/α/β are spawned as separate sub-agents via the Agent tool; `.cdd/unreleased/476/` artifacts are the shared memory across roles (no chat-state continuity). The package-owned dispatch wakes that would normally claim a `dispatch:cell + protocol:cdd + status:todo` cycle do not exist yet — Sub 4 builds the first one. Wave context: Sub 1 (#468) merged, Sub 2 (#470) merged at `043bf7aa` (the wake-provider manifest this sub RENDERS), this is Sub 3 (the renderer that consumes Sub 2's manifest); Subs 4-6 follow. A friction log will be appended to `gamma-closeout.md` capturing renderer-contract learnings — these feed potential amendments to `cnos.core/skills/agent/wake-provider/SKILL.md`, NOT Sub 5 (cnos#472 is a separate dispatch-prompt-template thread).

---

## Form choice (pinned by γ; α may override only with a recorded structural reason)

The issue offers two surfaces for the renderer: (A) Go subcommand in the `cn` binary under `src/go/internal/wake-install/` + dispatch wiring in `src/go/internal/cli/`; or (B) shell command under `cnos.core/commands/install-wake/cn-install-wake` registered in `cnos.core/cn.package.json`'s `commands` map.

γ pins **Option B — shell command at `src/packages/cnos.core/commands/install-wake/cn-install-wake`** registered in `cnos.core/cn.package.json`'s `commands` map as `install-wake` (entrypoint `commands/install-wake/cn-install-wake`).

Rationale:

1. **Sub 2's reservation.** Sub 2's γ scaffold (cycle/470 §"Form choice") explicitly reserved `cn-install-wake` as Sub 3's "execute the install" verb, separating it from the declaration (which lives under `orchestrators/agent-admin/`). Picking the shell-command form here honors that reservation directly: `orchestrators/agent-admin/` = data + prompt (declaration, Sub 2); `commands/install-wake/` = the verb (Sub 3, this cycle).
2. **Sibling-pattern precedent.** `cnos.core/commands/{daily, weekly, save}/` are all shell commands registered in `cn.package.json`'s `commands` map; the discovery + dispatch infrastructure already exists at `src/go/internal/discover/` (`ScanPackageCommands` walks `.cn/vendor/packages/*/cn.package.json` and registers an `ExecCommand` per entry). A new entry costs one row in the manifest + one executable file; no Go-side wiring change.
3. **Lighter v0 surface.** The renderer is essentially: read JSON manifest → read Markdown prompt template → emit YAML to a pinned path (idempotently). Shell with `jq` + heredocs handles this without introducing Go dependencies (`gopkg.in/yaml.v3` or equivalent) or new Go test infrastructure. Idempotence is straightforward: compute the YAML, compare bytes to the existing file, write only if differs.
4. **Faster iteration on the renderer-contract.** This is v0; the friction log will feed amendments to `cnos.core/skills/agent/wake-provider/SKILL.md` (manifest field gaps, substrate-vs-package boundary edges). Shell makes iterating the renderer cheap; the Go form pays its tax up-front for a contract still in flight.
5. **CLI integration cleanly registers.** Per Sub 2's pin rationale, "CLI integration target = None for the declaration"; for THIS cycle (the verb), CLI integration target = shell command registered in `cnos.core/cn.package.json`'s `commands` map. This is the natural symmetric move.

**Trade-offs acknowledged:**

- YAML emission in shell heredoc is more fragile than Go's `yaml.Marshal` (escaping, indentation discipline). Mitigated by the golden-fixture invariant (CI re-renders and byte-diffs); any escape bug is caught mechanically before merge.
- Less testable than Go (no `go test`); α writes a validation script invoked from CI instead. The golden fixture serves as the primary correctness oracle.
- If the v0 contract proves stable and the renderer needs to grow (more substrates, more sophisticated dispatch), a future cycle may port to Go. The shell v0 does NOT block that port — the contract surface is the same.

**If α discovers a structural reason to pick Option A (Go subcommand) instead** — e.g., a `cn`-binary invariant that commands consuming `cn.wake-provider.v1` must be Go-resident; or shell-side YAML emission proves unreproducible across `bash` versions; or path-confinement / package-isolation rules in `discover.ScanPackageCommands` prevent the command from reading from arbitrary installed packages' `orchestrators/` trees — α records the reason in `self-coherence.md` §Gap citing the file + line that forces the override, picks Go form, and adjusts package scoping to `src/go/internal/wake-install/` + `src/go/internal/cli/cmd_wake.go` (paralleling `cmd_doctor.go`, `cmd_dispatch.go`). γ's pin is not absolute.

---

## Render target path (pinned by γ; α may override only with a recorded structural reason)

The issue offers three target paths for the rendered `cnos-agent-admin.yml`: (a) golden-file-only location (cn test fixtures or per-package `*.golden.yml`); (b) `.github/workflows/.rendered/cnos-agent-admin.yml` (parallel directory under `.github/workflows/`); (c) `.github/workflows/cnos-agent-admin.yml` (would activate the new wake as a SECOND production wake — PROHIBITED by operator guardrail "no silent cutover").

γ pins **option (a) — golden file as a per-package sibling at `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`** AND the renderer writes its run-time output at the same path by default when invoked WITHOUT `--out` (i.e., the default install target IS the golden, which makes the golden + render-output paths trivially coincident, so `cn install-wake agent-admin` updates the golden in place; CI re-renders and diffs the working tree against `HEAD` to catch drift).

Rationale:

1. **Operator guardrail "no silent cutover" is binding.** Path (c) would create a second active production wake firing at the same cron slots as `claude-wake.yml`, doubling channel writes and racing the cursor. Prohibited; the cutover is a separate cycle.
2. **Path (b) `.github/workflows/.rendered/` is risky.** GitHub Actions discovers workflow files under `.github/workflows/` recursively; while subdirectories are NOT scanned for workflow definitions today (per current GHA semantics — only files directly under `.github/workflows/`), this behavior is undocumented as a stability contract and could change. The dot-prefix is also a leaky-by-convention signal. Risk of accidental activation if the directory rule changes.
3. **Path (a) is unambiguous.** The rendered file lives under the package source tree, named `*.golden.yml`, in the same directory as the declaration it derives from. CI's check is: re-render with `cn install-wake agent-admin` (or the inner-loop equivalent), `git diff cycle/476 --` exits clean. No path under `.github/workflows/` is created or modified, so AC7's byte-identical invariant on `claude-wake.yml` is automatically satisfied at the directory level (cutover is later).
4. **Tests-fixture variant rejected for v0.** Placing the golden under `src/go/internal/wake-install/testdata/` couples the golden to the Go form (Option A). Per-package sibling is form-agnostic; if α overrides to Go, the golden path remains valid and the Go-side test loads it directly.
5. **`cnos-` prefix preserved.** The rendered file's `name:` field and filename both follow the `cnos-{wake-name}.yml` convention per cnos#467 master and Sub 2's manifest's `superseded_substrate_artifact` semantics. The `.golden.yml` suffix is the v0 marker that this is a fixture, not an active workflow.

**If α discovers a structural reason to pick path (b) instead** — e.g., the operator-review workflow requires the rendered artifact to be visible alongside `claude-wake.yml` for side-by-side comparison; or GitHub's workflow discovery semantics on subdirectories prove to be guaranteed-recursive (operator-confirmed via GHA docs); or the cn renderer's idempotence check requires writing to a substrate-canonical path — α records the reason in `self-coherence.md` §Gap. γ's pin is not absolute.

**Render-target invariant for cutover (later cycle):** when cutover is approved (separate cycle), the operator runs `cn install-wake agent-admin --out .github/workflows/cnos-agent-admin.yml` (or equivalent) to materialize the active workflow, then deletes `claude-wake.yml` in the same commit. The Sub 3 renderer MUST support this `--out` flag so the cutover cycle does not require re-implementing the renderer. (Mechanical proof of `--out` support is added to AC2's oracle below.)

---

## Expected touched surfaces (α's prediction; α records actual in self-coherence §CDD Trace step 6)

| File | Purpose | New / edit |
|---|---|---|
| `src/packages/cnos.core/commands/install-wake/cn-install-wake` | AC1, AC2 — the renderer shell script (γ-pinned form); reads `cn.wake-provider.v1` manifest + paired prompt template; emits GHA YAML idempotently | new |
| `src/packages/cnos.core/cn.package.json` | Add `install-wake` entry to `commands` map (entrypoint `commands/install-wake/cn-install-wake`); bump `version` if cnos.core release rules require it (γ pin: do NOT bump version in this sub; γ records candidate release note in closeout) | edit (commands map only) |
| `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` | AC3, AC6 — golden rendering of the agent-admin provider per γ-pinned render target | new |
| `src/packages/cnos.core/commands/install-wake/README.md` *(optional)* | Cross-references; invocation examples; `--out` flag documentation | optional |
| `.github/workflows/install-wake-golden.yml` *(if not subsumed by an existing CI workflow)* | AC6 — CI mechanism that re-renders the golden + `git diff --exit-code` against `cycle/476` HEAD; runs on every PR touching `cnos.core/orchestrators/` or `cnos.core/commands/install-wake/`. α may instead add a step to an existing workflow if one carries equivalent triggers; α picks one approach and documents in `self-coherence.md` §Debt or §Gap | new OR edit existing |
| `.cdd/unreleased/476/self-coherence.md` | α's review-readiness signal + CDD Trace | new (α writes incrementally per `alpha/SKILL.md` §2.5) |
| `.cdd/unreleased/476/beta-review.md` | β's R1+ verdict + findings | new (β writes per `beta/SKILL.md`) |
| `.cdd/unreleased/476/alpha-closeout.md` | α's post-merge cycle narrative | new (α writes after β merge per re-dispatch) |
| `.cdd/unreleased/476/beta-closeout.md` | β's post-merge release context | new (β writes at merge) |
| `.cdd/unreleased/476/gamma-closeout.md` | γ's cycle closure declaration + friction log (renderer-contract learnings) | new (γ writes after β passes; out of scope for α/β prompts) |
| `.cdd/unreleased/476/gamma-scaffold.md` | this file | (already written) |

**Explicitly NOT touched:**

- `.github/workflows/claude-wake.yml` — AC7 invariant: the existing hand-written wake remains byte-identical on `cycle/476`. Cutover happens in a separate cycle.
- `.github/workflows/cnos-agent-admin.yml` — would activate a second production wake; PROHIBITED.
- `.github/workflows/.rendered/*.yml` — alternative path; γ pin says no.
- Any package other than `cnos.core` — Sub 4 builds the cnos.cdd dispatch wake provider; not this sub.
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` — δ wake-invoked mode is Sub 5.
- `cn-sigma:.cn-sigma/spec/OPERATOR.md` — out of scope; cited from this provider, not edited.
- `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` + `prompt.md` — Sub 2's shipped artifacts; CONSUMED by the renderer, NOT modified.
- `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` — the contract α reads; renderer must conform; friction notes about manifest field gaps go to `gamma-closeout.md` and may seed a follow-up issue to amend, NOT this cycle.

---

## Implementation contract (pinned by δ at dispatch; verified at this scaffold)

| Axis | Pinned value |
|---|---|
| Language | **Shell (POSIX sh / bash)** + Markdown + YAML (the renderer EMITS YAML; the source is shell). No Go (per γ form-pin Option B); if α overrides to Option A, this becomes **Go** + Markdown + YAML and `src/go/` package scoping shifts accordingly. No Python. |
| CLI integration target | **`install-wake` registered in `src/packages/cnos.core/cn.package.json`'s `commands` map** with entrypoint `commands/install-wake/cn-install-wake`; invoked as `cn install-wake <wake-name>` (per `discover.ScanPackageCommands` + `ExecCommand` infrastructure already present in `src/go/internal/discover/`). NO change to `src/go/internal/cli/` (no new Go-side cmd file). |
| Package scoping | `src/packages/cnos.core/` only. New files under `cnos.core/commands/install-wake/` (renderer) and `cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` (golden fixture, per γ-pinned render target). `cnos.core/cn.package.json` edited to add the `install-wake` commands-map entry only. NO edits to any other package; NO edits to `src/go/`. |
| Existing-binary disposition | N/A — `cn` Go binary unchanged. Sibling shell commands (`daily`, `weekly`, `save`) untouched. The new `install-wake` is additive in the `commands` map. |
| Runtime dependencies | Standard POSIX shell utilities (`sh` / `bash`, `cat`, `cut`, `diff`, `mkdir`, `mv`, `cp`, `rm`, `test`, `printf`); `jq` for JSON manifest parsing (commonly available; α confirms whether `jq` is in CI image; if not, α picks a fallback — e.g., a small Python one-liner under a pinned `python3` — but the γ pin prefers `jq` for shell-idiomatic parsing); `sha256sum` (or `shasum -a 256` on macOS, with portability shim) for idempotence. NO new runtime dependencies added to the repo's `go.mod`, `package.json`, or any other dependency manifest. |
| JSON/wire contract | Consumes `cn.wake-provider.v1` (the schema shipped in Sub 2; α MUST read `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` for the field-level contract and `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` for the reference instance). Emits GitHub Actions workflow YAML matching `claude-code-action@v1` invocation shape (`on:` block, `permissions:` block, `concurrency:` block, `jobs.{name}.steps[].uses: actions/checkout@v4`, `jobs.{name}.steps[].uses: anthropics/claude-code-action@v1`, `prompt:` inlined from `prompt_template`). The substrate-rendering target mapping is per `wake-provider/SKILL.md §2.4` "Substrate-rendering target" + §2.5 canonical split table. |
| Backward compat | `.github/workflows/claude-wake.yml` MUST remain byte-identical on `cycle/476`. Mechanical check before signaling review-readiness: `git diff origin/main..HEAD -- .github/workflows/claude-wake.yml \| wc -l` MUST return 0. Additionally, `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` and `prompt.md` MUST remain byte-identical (Sub 2 declarations are CONSUMED, not edited): `git diff origin/main..HEAD -- src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json src/packages/cnos.core/orchestrators/agent-admin/prompt.md \| wc -l` MUST return 0. |

---

## AC mapping (oracles + surfaces; α records actual evidence in `self-coherence.md` §ACs)

| AC | Invariant (per #476) | Oracle (mechanical) | Surface |
|----|----------------------|---------------------|---------|
| AC1 | Renderer form pinned + rationale documented in `gamma-scaffold.md` §"Form choice"; α accepts the pin or records override per the discipline. | This scaffold's §"Form choice" exists; α's `self-coherence.md` §Gap either acknowledges the pin or records override with file + line citing the override reason. β verifies the scaffold §"Form choice" and α's §Gap line are coherent. | `.cdd/unreleased/476/gamma-scaffold.md` (this file); `.cdd/unreleased/476/self-coherence.md` |
| AC2 | Renderer consumes `cn.wake-provider.v1` manifest + paired prompt template at `orchestrators/{name}/wake-provider.json` + `orchestrators/{name}/prompt.md`. Rejects malformed manifests with precise error. Supports `--out <path>` flag for emitting to non-default paths (needed for future cutover). | Positive: `cn install-wake agent-admin` (invoked from a repo with `cnos.core` installed, or via a local-dev shim) succeeds; exits 0; emits the golden output at the γ-pinned path. Negative: invoking against a fabricated manifest with `"schema"` absent fails with stderr matching `/required field "schema"/`; invoking with `"admin_only": "true"` (string instead of bool) fails with `/admin_only must be boolean/`. `--out` flag oracle: `cn install-wake agent-admin --out /tmp/test-out.yml` writes to `/tmp/test-out.yml` and not to the default path. | `cnos.core/commands/install-wake/cn-install-wake`; positive + negative test invocations (α writes a small test script under `commands/install-wake/` or in a CI step). |
| AC3 | Renders cnos.core/agent-admin provider to a syntactically valid GHA workflow YAML at the γ-pinned target path. Output is human-readable (prompt content preserved verbatim), substrate-correct (valid `on:`, `permissions:`, `concurrency:`, `jobs:`, `claude-code-action@v1`), substrate-coupled only at the renderer's responsibility boundary. | `python3 -c "import yaml; yaml.safe_load(open('src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml'))"` exits 0; `grep -c "anthropics/claude-code-action@v1" cnos-agent-admin.golden.yml` ≥ 1; `grep -cE "^on:" cnos-agent-admin.golden.yml` ≥ 1; `grep -cE "^permissions:" cnos-agent-admin.golden.yml` ≥ 1; `grep -cE "^concurrency:" cnos-agent-admin.golden.yml` ≥ 1; `grep -cE "^jobs:" cnos-agent-admin.golden.yml` ≥ 1. | Rendered file at `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`. |
| AC4 | Rendered workflow's `prompt:` field contains the content of `cnos.core/orchestrators/agent-admin/prompt.md`. The admin-only constraint language (per `prompt.md` §"Admin-only boundary" + §"Disallowed surfaces") is present in the rendered workflow. | `grep -c "MUST NOT execute cells under any circumstance" cnos-agent-admin.golden.yml` ≥ 1; `grep -cE "cell_execution\|disallowed_surfaces\|cell execution" cnos-agent-admin.golden.yml` ≥ 1; the prompt template's body character count is within 95-105% of the original (`wc -c prompt.md` baseline; rendered `prompt:` block character count via a small extractor; allows for YAML indentation overhead, forbids silent truncation). | Rendered file content; `prompt.md` original. |
| AC5 | Rendered workflow's `on:` triggers map 1:1 to the manifest's `input_contract.triggers` (`schedule` → `on.schedule` with cron entries; `issues_opened_title_match` → `on.issues` with type filter + title-match conditional). Rendered `permissions:` block respects `output_contract` + `allowed_surfaces` and matches `permission_intent` from the manifest (no permissions beyond the manifest's authorization). | Per-trigger oracle (enumerated, NOT aggregated — see §"Claim-class verification" injected into α prompt below): `for trigger in $(jq -r '.input_contract.triggers[]' wake-provider.json); do grep -q "$trigger" cnos-agent-admin.golden.yml \|\| echo "MISSING: $trigger"; done` MUST emit zero `MISSING:` lines. Per-permission oracle: `for perm in $(jq -r '.permission_intent[]' wake-provider.json); do gha_key=$(echo "$perm" \| sed 's/\./: /'); grep -q "$gha_key" cnos-agent-admin.golden.yml \|\| echo "MISSING: $perm"; done` MUST emit zero `MISSING:` lines. Negative: `grep -oE "^  [a-z_-]+: (write\|read)" cnos-agent-admin.golden.yml \| awk '{print $1}' \| sort -u` MUST equal the manifest's `permission_intent` set (no extras). | Rendered file vs `wake-provider.json`; comparison script α writes or β re-runs. |
| AC6 | Golden fixture committed at the γ-pinned path; idempotence verified. CI mechanism re-renders the golden + diffs against `HEAD` on every PR touching `cnos.core/orchestrators/` or `cnos.core/commands/install-wake/`. Running the renderer twice produces byte-identical output. | Golden file exists (`test -f cnos-agent-admin.golden.yml`); idempotence: `cn install-wake agent-admin && sha256sum cnos-agent-admin.golden.yml > /tmp/r1.sum && cn install-wake agent-admin && sha256sum cnos-agent-admin.golden.yml > /tmp/r2.sum && diff /tmp/r1.sum /tmp/r2.sum` exits 0; CI workflow exists (or step exists in an existing workflow) that runs the equivalent re-render + diff on PR. | Golden file; CI workflow file; idempotence proof in `self-coherence.md` §ACs. |
| AC7 | `.github/workflows/claude-wake.yml` byte-identical on this cycle. γ closeout includes the candidate cutover-issue title/body (γ-only; out of α/β scope). | `git diff origin/main HEAD -- .github/workflows/claude-wake.yml \| wc -l` MUST be 0; γ records cutover candidate in `gamma-closeout.md`. | Mechanical diff; `gamma-closeout.md`. |
| AC8 | Package-vs-renderer authority split preserved (per `wake-provider/SKILL.md §3` "Authority split"). Renderer source contains no admin-only / disallowed-surfaces / defer-path / cell-execution language (those come from the prompt template substitution); package's manifest contains no GHA YAML / runs-on / claude-code-action / GITHUB_TOKEN beyond the existing AC7 carve-out for `relationship_to_substrate`. | Renderer-side audit: `grep -ciE "admin.only\|disallowed_surfaces\|defer.path\|cell_execution" src/packages/cnos.core/commands/install-wake/cn-install-wake` MUST be 0. Package-side audit: `grep -ciE "runs-on\|uses:\|claude-code-action\|GITHUB_TOKEN\|GHA" src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json src/packages/cnos.core/orchestrators/agent-admin/prompt.md \| grep -vE "AC5\|AC7\|relationship_to_substrate"` MUST be ≤ already-audited carve-outs from cnos#470 (β re-runs the same grep that passed at #470 merge and confirms identical count). | Renderer source; Sub 2 declaration files (read-only); audit grep recorded in `beta-review.md`. |

**Mechanical gate (run by α at pre-review-gate row 9 polyglot re-audit and by β at row 3 pre-merge non-destructive merge-test):**

```bash
# Mechanical proof of AC7 (no claude-wake.yml changes):
git diff --name-only origin/main..HEAD -- .github/workflows/claude-wake.yml | wc -l  # MUST be 0

# Mechanical proof of Sub 2 declaration not edited:
git diff --name-only origin/main..HEAD -- src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json src/packages/cnos.core/orchestrators/agent-admin/prompt.md | wc -l  # MUST be 0

# Mechanical proof of scope discipline:
git diff --name-only origin/main..HEAD | grep -vE '^(src/packages/cnos\.core/(commands/install-wake/|orchestrators/agent-admin/cnos-agent-admin\.golden\.yml|cn\.package\.json)|\.cdd/unreleased/476/|\.github/workflows/install-wake-golden\.yml)$' | wc -l  # MUST be 0  (allow CI workflow file if α added new)

# Mechanical proof of AC6 idempotence:
cn install-wake agent-admin
sha256sum src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml | cut -d' ' -f1 > /tmp/r1.sum
cn install-wake agent-admin
sha256sum src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml | cut -d' ' -f1 > /tmp/r2.sum
diff /tmp/r1.sum /tmp/r2.sum  # MUST exit 0

# Mechanical proof of AC8 package-vs-renderer authority split (renderer side):
grep -ciE "admin.only|disallowed_surfaces|defer.path|cell_execution" src/packages/cnos.core/commands/install-wake/cn-install-wake  # MUST be 0

# Mechanical proof of AC8 package-vs-renderer authority split (package side; carve-out for relationship_to_substrate):
grep -ciE "runs-on|uses:|claude-code-action|GITHUB_TOKEN|GHA" src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json src/packages/cnos.core/orchestrators/agent-admin/prompt.md | grep -vE "AC5|AC7|relationship_to_substrate"  # MUST be ≤ existing-as-of-#470-merge count (β re-runs against origin/main at base SHA and confirms identical)
```

---

## α dispatch prompt

(Below is the full self-contained prompt δ feeds to the spawned α session. α has zero shared memory with γ or β. The prompt carries everything α needs.)

---

```
You are α (alpha) for CDD cycle cnos#476.

# Your one job
Implement cnos#476 (`cn-wake-install: v0 subcommand to render package-owned wakes into substrate (consume cn.wake-provider.v1)`). Produce the review-ready artifact set on `cycle/476`. Write `.cdd/unreleased/476/self-coherence.md` incrementally per the α SKILL §2.5 (one section per commit, pushed). Signal review-readiness by appending the `## Review-readiness` section to `self-coherence.md` and pushing. **Do not spawn anyone. Do not dispatch β.** δ (the operator-routed bootstrap session) will dispatch β after you signal review-readiness.

# Environment
- Working directory: `/home/user/cnos` (the cnos repository working tree)
- Branch: **`cycle/476`** (created by γ from `origin/main` at SHA `fcc5cdb9a533ad86e67524bcf05a33d2b4592e8a`)
- Repository: `usurobor/cnos`
- Issue URL: https://github.com/usurobor/cnos/issues/476
- Master tracker URL: https://github.com/usurobor/cnos/issues/467 (read the "Foundational architecture — package-owned wake providers (authoritative)" section)
- Predecessor sibs: cnos#468 (Sub 1, merged c0048bef); cnos#470 (Sub 2, merged 043bf7aa) — the manifest you RENDER

# Git identity
Before any commit, set:
  git config user.email "alpha@cdd.cnos"
  git config user.name "alpha"
Verify with `git config --get user.email`. (Per `operator/SKILL.md §Git identity for role actors` + `alpha/SKILL.md §2.6 row 14`.)

# Cycle execution mode (read before claiming any rule does not apply)
This cycle is executed through the **pre-dispatch δ/channel path**: γ-interface session acts as bootstrap-δ, spawns α (you) via the Agent tool, then will spawn β after your review-readiness signal. The package-owned dispatch wakes that would normally claim this cycle do not exist yet — Sub 3 (this cycle) is part of building them. β is NOT yet spawned; do not poll for β verdicts during your authoring run. δ will run β after you exit on review-readiness. α exits after signaling review-readiness (sequential bounded dispatch per `cnos.cds/skills/cds/CDS.md` §"Field 6: Actor collapse rule"); δ will re-dispatch you after β merges for `alpha-closeout.md`.

# Skills to load (in this order, before any implementation commit)
Tier 1:
1. `src/packages/cnos.cdd/skills/cdd/CDD.md` (if present; else skip)
2. `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (your role)

Lifecycle sub-skills as work requires:
3. `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md`
4. `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` (this is design-and-build — load before any renderer-form decision)

Tier 3 (issue-specific — read these FIRST before any code change):
5. **`src/packages/cnos.core/skills/agent/wake-provider/SKILL.md`** — the contract you are CONSUMING; pay special attention to §3 "Authority split" and §2.5 canonical table; the renderer must not leak into the package's authority and vice versa
6. **`src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json`** — the reference instance you render; READ ONLY; do NOT edit
7. **`src/packages/cnos.core/orchestrators/agent-admin/prompt.md`** — the prompt template you inline; READ ONLY; do NOT edit
8. `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` (cnos#468 merged predecessor; contextual)
9. `src/packages/cnos.core/cn.package.json` (current commands map; you'll add the `install-wake` entry)
10. `src/packages/cnos.core/commands/daily/cn-daily`, `commands/weekly/cn-weekly`, `commands/save/cn-save` — shell-command sibling precedent (environment variables `CN_HUB_PATH`, `CN_PACKAGE_ROOT`, `CN_COMMAND_NAME` exported by `cn dispatch`; follow the same convention)
11. `src/go/internal/discover/` (`ScanPackageCommands` shows how the cn binary discovers shell commands via `cn.package.json`'s `commands` map — your new entry must conform)
12. `.github/workflows/claude-wake.yml` (READ ONLY — the substrate-bound wake you are reproducing the rendering of; AC7 invariant says you do NOT modify this file)
13. `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` (cited from the manifest's `output_contract`; the prompt template references it)

Also re-read the γ scaffold for full context including form-choice rationale, render-target rationale, and AC mapping:
14. `.cdd/unreleased/476/gamma-scaffold.md`

# Implementation contract (pinned by δ; α MUST NOT improvise these axes)

| Axis | Pinned value |
|---|---|
| Language | Shell (POSIX sh / bash; entrypoint is shebang `#!/bin/sh` or `#!/bin/bash` per α's call — daily/weekly/save use `#!/bin/sh`) + Markdown + YAML (emitted, not source). No Go, no Python (unless `jq` is unavailable and α picks a `python3` fallback for JSON parsing — record the choice in §Gap). |
| CLI integration target | `install-wake` registered in `src/packages/cnos.core/cn.package.json`'s `commands` map, entrypoint `commands/install-wake/cn-install-wake`. Invoked as `cn install-wake <wake-name>`. NO change to `src/go/internal/cli/`. |
| Package scoping | `src/packages/cnos.core/` only. New: `commands/install-wake/cn-install-wake` + optional `commands/install-wake/README.md`; `orchestrators/agent-admin/cnos-agent-admin.golden.yml`. Edit: `cn.package.json` `commands` map (one entry added). Optionally: a CI workflow under `.github/workflows/install-wake-golden.yml` for the golden-diff check. |
| Existing-binary disposition | N/A — `cn` Go binary unchanged. Sibling commands untouched. |
| Runtime dependencies | POSIX shell utilities + `jq` (preferred for JSON parsing). If `jq` is not in the cnos CI image, use a `python3` fallback and record in §Gap. NO new deps added to any dep manifest. |
| JSON/wire contract | Consume `cn.wake-provider.v1` per `cnos.core/skills/agent/wake-provider/SKILL.md`. Reject malformed manifests with precise stderr per `wake-provider/SKILL.md §3.5`. Emit GHA YAML per `claude-code-action@v1` invocation shape (use `claude-wake.yml` as the structural reference — but reproduce ONLY the shape that the manifest's fields dictate; do not hardcode claude-wake.yml's specific cron slots or secret names beyond what the manifest authorizes). |
| Backward compat | `.github/workflows/claude-wake.yml` byte-identical; Sub 2 declarations (`orchestrators/agent-admin/wake-provider.json` + `prompt.md`) byte-identical. Mechanical: `git diff origin/main..HEAD -- .github/workflows/claude-wake.yml src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json src/packages/cnos.core/orchestrators/agent-admin/prompt.md | wc -l` MUST be 0. |

α MUST NOT change any pinned axis. If you find a structural reason to change one, surface to γ via `.cdd/unreleased/476/gamma-clarification.md` (write it, commit it on `cycle/476`, and stop) before making any contradicting commit.

# Form choice (γ-pinned; override only with recorded structural reason)
γ pins the renderer as a **shell command at `src/packages/cnos.core/commands/install-wake/cn-install-wake`** registered in `cnos.core/cn.package.json`'s `commands` map as `install-wake`. NOT a Go subcommand. Rationale is in `.cdd/unreleased/476/gamma-scaffold.md` §"Form choice"; read it before authoring. You MAY override to Option A (Go subcommand under `src/go/internal/wake-install/` + `src/go/internal/cli/cmd_wake.go`) ONLY if you discover a structural reason — e.g., `discover.ScanPackageCommands` path-confinement rules prevent the command from reading from arbitrary installed packages' `orchestrators/` trees; or shell-side YAML emission proves unreproducible across `bash` versions; or a `cn`-binary invariant requires consuming-of-`cn.wake-provider.v1` to be Go-resident. If you override, record the reason in `self-coherence.md` §Gap citing the file + line that forces the override.

# Render target (γ-pinned; override only with recorded structural reason)
γ pins the rendered golden at **`src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`** (per-package sibling). The renderer's default output path (when invoked WITHOUT `--out`) IS this path; `cn install-wake agent-admin` updates the golden in place. The renderer MUST also support `--out <path>` for future cutover (when a separate cycle activates the rendered workflow at `.github/workflows/cnos-agent-admin.yml`). Rationale in `gamma-scaffold.md` §"Render target path". You MAY override to `.github/workflows/.rendered/cnos-agent-admin.yml` ONLY with a recorded structural reason. You MUST NOT render to `.github/workflows/cnos-agent-admin.yml` — that activates a second production wake, PROHIBITED.

# Acceptance criteria (verbatim from cnos#476 — AC1 through AC8)

[α MUST read the issue body via `mcp__github__issue_read` with method=get, owner=usurobor, repo=cnos, issue_number=476, OR via the GitHub web URL if MCP github is unavailable. Each AC's Invariant / Oracle / Surface block must be read in full; the γ scaffold's AC mapping table is a navigation aid, not a substitute.]

AC1 — Renderer form pinned + rationale documented (γ scaffold §"Form choice"; α accepts or overrides per discipline)
AC2 — Renderer consumes `cn.wake-provider.v1` manifest; rejects malformed inputs with precise error
AC3 — Renders cnos.core/agent-admin provider to a syntactically valid GHA workflow YAML at the γ-pinned target path
AC4 — Rendered prompt block carries the admin-only boundary verbatim from `prompt.md`
AC5 — Triggers + permissions in the rendered workflow match the manifest contract
AC6 — Golden fixture exists; idempotence verified; CI mechanism in place
AC7 — `.github/workflows/claude-wake.yml` byte-identical on this cycle
AC8 — Package-vs-renderer authority split preserved (renderer source no admin-only / cell-execution / defer-path language; package manifest no GHA-specific YAML beyond #470 carve-outs)

# Expected surfaces (γ scaffold §"Expected touched surfaces" — re-read it; this is summary only)

NEW:
- `src/packages/cnos.core/commands/install-wake/cn-install-wake` (the renderer shell script)
- `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` (golden fixture)
- (Optional) `src/packages/cnos.core/commands/install-wake/README.md`
- (Optional, OR step in existing) `.github/workflows/install-wake-golden.yml` (CI re-render + golden diff)
- `.cdd/unreleased/476/self-coherence.md` (your review-readiness signal)

EDIT:
- `src/packages/cnos.core/cn.package.json` (one new entry in `commands` map: `"install-wake": {"entrypoint": "commands/install-wake/cn-install-wake", "summary": "Render a package-owned wake into the substrate"}`). DO NOT bump `version` — γ pin is "no version bump in this sub"; γ records candidate release note in closeout.

NOT TOUCHED (refusal conditions — if you find yourself editing these, STOP and re-read this prompt):
- `.github/workflows/claude-wake.yml` (AC7 byte-identical invariant)
- `.github/workflows/cnos-agent-admin.yml` (would activate second production wake; PROHIBITED)
- `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` (Sub 2's shipped manifest; CONSUMED, not edited)
- `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` (Sub 2's shipped prompt; INLINED, not edited)
- `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` (the contract you read; friction goes to gamma-closeout, NOT a same-cycle amendment)
- Any package other than `cnos.core` (esp. `cnos.cdd/`, `cnos.cdr/`, `cnos.kata/`)
- `src/go/` (γ-pinned form is shell; if you override to Go, this changes — record in §Gap)
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` (Sub 5)
- `cn-sigma:` anything (different repo)

# Refusal conditions (binding)
α MUST refuse and surface to γ via `.cdd/unreleased/476/gamma-clarification.md` (commit on cycle/476, then stop) if:
- The cycle requires editing `.github/workflows/claude-wake.yml` to satisfy any AC (it does not; AC7 says the opposite).
- The cycle requires rendering to `.github/workflows/cnos-agent-admin.yml` directly (PROHIBITED — would activate a second production wake; cutover is a separate cycle).
- The cycle requires editing the Sub 2 declaration (`orchestrators/agent-admin/wake-provider.json` or `prompt.md`) — those are CONSUMED, not modified. If α finds a manifest field is missing for the renderer's needs, record a friction note in §Debt for γ's closeout (may seed a follow-up to amend `wake-provider/SKILL.md`); IF the missing field BLOCKS rendering, file a `gamma-clarification.md` and stop.
- The cycle requires implementing the CDD dispatch wake (Sub 4), δ wake-invoked-mode (Sub 5), or end-to-end smoke (Sub 6).
- The cycle requires implementing any non-GHA substrate support (deferred; out of scope per operator guardrail).
- The cycle requires implementing any non-claude-code-action backend (deferred; cnos#452 territory).
- A pinned implementation-contract axis appears unsatisfiable from this scope.
- An AC oracle from the issue body appears unsatisfiable from the pinned form / scope.

# Process

1. Set git identity (above).
2. `git fetch origin && git switch cycle/476` (the branch exists — γ created it).
3. Read everything in §"Skills to load" above — in order, especially items 5-7 (the contract you consume + the reference instance).
4. Re-read `.cdd/unreleased/476/gamma-scaffold.md` (this scaffold) in full.
5. Read cnos#476 issue body in full (verbatim ACs).
6. Read cnos#467 master tracker — specifically §"Foundational architecture — package-owned wake providers (authoritative)" and §"Acceptance criteria (wave-level)" AC4.
7. Read cnos#470 (the predecessor sub) — its `.cdd/releases/docs/2026-06-21/470/gamma-scaffold.md` is the scaffold-shape precedent; its `wake-provider/SKILL.md` is the contract.
8. Begin writing `.cdd/unreleased/476/self-coherence.md` incrementally per α SKILL §2.5:
   - §Gap (issue, version/mode = `design-and-build`, form-choice acknowledgment or documented override, render-target acknowledgment or documented override)
   - §Skills (Tier 1/2/3 loaded)
   - Then implement (commits as you go):
     - The renderer shell script (consumes `cn.wake-provider.v1`; emits YAML; idempotent; supports `--out`)
     - The `cn.package.json` `commands` map entry
     - The golden fixture (`cn install-wake agent-admin` produces it; commit the produced output)
     - The CI mechanism (new workflow file OR step in existing) that re-renders + diffs
     - The positive + negative test invocations (AC2 oracles)
   - §ACs (AC-by-AC evidence pointing at the diff; mechanical re-runs of every oracle)
   - §Self-check (did α push ambiguity onto β? is every claim backed by diff evidence?)
   - §Debt (known debt; "none" is acceptable if true; renderer-contract friction notes go HERE if minor, in `gamma-clarification.md` if blocking)
   - §CDD Trace (steps 0-7)
   - §Review-readiness (the signal that exits α to β)
9. Before §Review-readiness: run ALL mechanical gates from the γ scaffold §"AC mapping" → "Mechanical gate" block. ALL must pass.
10. Run α SKILL §2.6 pre-review-gate rows 1-15. Row 9 (polyglot re-audit) for this cycle = shell syntax (`bash -n` on `cn-install-wake`), JSON validity (`jq . cn.package.json`), YAML validity (`python3 -c "import yaml; yaml.safe_load(...)"` on the golden), Markdown grep, AC oracle re-runs.
11. Append §Review-readiness with the implementation SHA (per α SKILL §2.6 "SHA convention for readiness signal" — implementation SHA, NOT readiness-signal-commit's-own SHA).
12. Push to `origin/cycle/476`. Exit. Do NOT poll for β; δ runs β.

# Claim-class verification (BINDING; from cnos#472 friction; auto-injected by γ per cnos#470 PRA F-α-1)

**Any §Self-check claim of form "all X work/resolve/pass" MUST be backed by a per-X mechanical proof table.** Aggregated claims are forbidden. Specifically for this cycle:

- "All manifest fields render correctly" → table with row per `wake-provider.json` required field, column showing the grep / `jq` command that verifies the rendered YAML carries the field's transformed value.
- "All triggers map" → table with row per trigger in `input_contract.triggers` (today: `schedule`, `issues_opened_title_match`), column showing the per-trigger grep against `cnos-agent-admin.golden.yml`.
- "All permissions correct" → table with row per `permission_intent` entry (today: `contents.write`, `issues.write`, `pull_requests.write`, `id_token.write`), column showing the per-permission grep against the rendered `permissions:` block.
- "All ACs pass" → table with row per AC (AC1-AC8), each with the mechanical oracle command + observed result.
- "Idempotence verified" → table showing the two consecutive `sha256sum` outputs side-by-side, with the `diff` exit code.

No "looks good", "all enumerated", "all checked" — every per-item line must be proved. This is cite-§3.13c of `cnos.cdd/skills/cdd/review/SKILL.md` (β SKILL Rule 6 anchors oracle evidence on code, not doc; aggregated claims are doc, per-item proofs are code).

# Output artifacts (you write)
- `.cdd/unreleased/476/self-coherence.md` (incremental, one section per commit)
- Implementation files under `src/packages/cnos.core/` per §"Expected surfaces"

You do NOT write:
- `beta-review.md` (β writes it)
- `alpha-closeout.md` (you write it later, after δ re-dispatches you post-β-merge)
- `gamma-closeout.md` (γ writes it)
- Any `cdd-iteration.md` (γ writes it if findings warrant)
- `RELEASE.md` (γ decides; not in α's scope)

# Commit message guidance
- One commit per self-coherence section: `α-476 self-coherence: §Gap`, `α-476 self-coherence: §Skills`, etc.
- Implementation commits: `α-476: cn-install-wake renderer (AC1+AC2)`, `α-476: cn.package.json install-wake entry`, `α-476: agent-admin golden fixture (AC3+AC6)`, `α-476: CI golden-diff workflow (AC6)`, `α-476: AC2 negative-case test`, etc.
- Push after each commit. The cycle branch is the coordination surface.

# Stop condition
Stop when §Review-readiness is appended to `self-coherence.md` and pushed to `origin/cycle/476`. Do NOT wait for β. Do NOT message anyone. Return a short summary: which commits landed, where review-readiness was signaled, any known debt declared in §Debt, any γ-clarification you needed to file, any form-pin or render-target-pin override you made and the recorded reason. δ will read your summary and dispatch β.
```

---

## β dispatch prompt

(Below is the full self-contained prompt δ feeds to the spawned β session after α signals review-readiness. β has zero shared memory with γ or α.)

---

```
You are β (beta) for CDD cycle cnos#476.

# Your one job
Review α's implementation of cnos#476 (`cn-wake-install: v0 subcommand to render package-owned wakes into substrate (consume cn.wake-provider.v1)`) on `cycle/476`. Write `.cdd/unreleased/476/beta-review.md` with a per-round verdict (R1, R2, ...) and binding findings. If verdict = APPROVE, merge to `main`; if RC, request changes and exit (α will re-dispatch). After merge, write `.cdd/unreleased/476/beta-closeout.md`. **Do not spawn anyone. Do not dispatch γ or α.** δ (the operator-routed bootstrap session) will dispatch α for close-out after merge.

# Environment
- Working directory: `/home/user/cnos`
- Branch under review: **`cycle/476`** (γ created from `origin/main` at SHA `fcc5cdb9a533ad86e67524bcf05a33d2b4592e8a`; α implemented on it)
- Branch to merge into: `main`
- Repository: `usurobor/cnos`
- Issue URL: https://github.com/usurobor/cnos/issues/476
- Master tracker URL: https://github.com/usurobor/cnos/issues/467
- PR URL: {δ will inject the PR URL here at dispatch time if one exists; otherwise β operates branch-direct via `git merge` per β SKILL §"Phase map"}

# Git identity
Before any commit (including the merge commit), set:
  git config user.email "beta@cdd.cnos"
  git config user.name "beta"
Verify with `git config --get user.email`. Per β SKILL §"Pre-merge gate" row 1 — if you use a worktree for merge-testing, use `git config --worktree user.email "beta@cdd.cnos"` to avoid leaking to shared config.

# Cycle execution mode (read before claiming any rule does not apply)
This cycle is executed through the **pre-dispatch δ/channel path**: γ-interface session acts as bootstrap-δ, spawned α (already exited at review-readiness), now spawning β (you). Sequential bounded dispatch per `cnos.cds/skills/cds/CDS.md` §"Field 6: Actor collapse rule". After your merge + `beta-closeout.md`, you exit; δ re-dispatches α for `alpha-closeout.md` and then γ writes `gamma-closeout.md`.

# Skills to load (in this order, before reviewing)
Tier 1:
1. `src/packages/cnos.cdd/skills/cdd/CDD.md` (if present)
2. `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` (your role)
3. `src/packages/cnos.cdd/skills/cdd/review/SKILL.md`
4. `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` (merge mechanics only; this cycle does NOT cut a release)

Tier 3 (issue-specific):
5. **`src/packages/cnos.core/skills/agent/wake-provider/SKILL.md`** — the contract α's renderer must consume; especially §3 "Authority split" and §2.5 canonical table; you audit BOTH the renderer source (for role-decision leakage) AND the Sub 2 manifest (for substrate-emission leakage; this latter audit confirms #470's existing carve-outs were not enlarged)
6. **`src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json`** — the manifest the renderer consumed; READ ONLY; verify byte-identical on cycle/476
7. **`src/packages/cnos.core/orchestrators/agent-admin/prompt.md`** — the prompt template; READ ONLY; verify byte-identical
8. `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` (cnos#468 — contextual)
9. `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` (cited from the manifest; affects rendered output indirectly)
10. `.github/workflows/claude-wake.yml` (AC7 byte-identical invariant — confirm unchanged; AND use as structural reference for what a valid GHA workflow looks like when validating AC3/AC5)
11. `src/packages/cnos.core/cn.package.json` (verify the `install-wake` entry was added correctly)
12. `src/go/internal/discover/` (`ScanPackageCommands` — confirm α's entry would be discovered correctly by the cn binary)

Also re-read for context:
13. `.cdd/unreleased/476/gamma-scaffold.md` (γ's scaffold — AC mapping table is your review checklist; form-choice + render-target rationale are γ's pins; verify which path α took and whether the override (if any) is justified)
14. `.cdd/unreleased/476/self-coherence.md` (α's review-readiness signal — CDD Trace, AC evidence, known debt, claim-class verification tables)

# Issue context (for verbatim ACs)
Read cnos#476 issue body in full via `mcp__github__issue_read` (owner=usurobor, repo=cnos, issue_number=476, method=get). ACs AC1-AC8 are the binding gate. β SKILL Rule 6 anchors oracle evidence on code (the diff), not doc — re-grep α's actual artifacts; do not trust α's claims without verifying against the diff.

# Pre-merge gate (β SKILL §"Pre-merge gate"; binding before `git merge`)
Run all four rows. If any row fails, RC with severity D on that row's failure mode.

Row 1 — Identity truth: `git config user.email` returns `beta@cdd.cnos`. If you used a merge-test worktree, ensure no leak via the `--worktree` discipline.

Row 2 — Canonical-skill freshness: `git fetch --verbose origin main && git rev-parse origin/main` returns SHA matching session-start snapshot. If `origin/main` advanced mid-cycle, re-load Tier-1 skills and re-evaluate.

Row 3 — Non-destructive merge-test (this cycle's specific validators):
- `git worktree add /tmp/cnos-merge-476/wt origin/main && cd /tmp/cnos-merge-476/wt && git config --worktree user.name beta && git config --worktree user.email beta@cdd.cnos && git merge --no-ff --no-commit origin/cycle/476`
- Confirm zero unmerged paths.
- Validators specific to this cycle's surface:
  - `bash -n src/packages/cnos.core/commands/install-wake/cn-install-wake` (shell syntax)
  - `jq . src/packages/cnos.core/cn.package.json > /dev/null` (JSON validity)
  - `python3 -c "import yaml; yaml.safe_load(open('src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml'))"` (YAML validity)
  - `./tools/validate-skill-frontmatter.sh` IF the cycle modifies any SKILL.md frontmatter (it shouldn't — wake-provider/SKILL.md should be byte-identical)
  - `cn-cdd-verify` IF present (for `.cdd/` artifacts)
  - The kata under `src/packages/cnos.kata/katas/{N}/` whose surface the cycle touches — likely none, but enumerate by `ls src/packages/cnos.kata/katas/` and confirm no kata covers `cn install-wake` rendering (if one does, run it)
- Mechanical AC7 check: `git diff origin/main HEAD -- .github/workflows/claude-wake.yml | wc -l` MUST be 0
- Mechanical Sub 2 declaration unchanged: `git diff origin/main HEAD -- src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json src/packages/cnos.core/orchestrators/agent-admin/prompt.md | wc -l` MUST be 0
- Mechanical scope check: `git diff --name-only origin/main HEAD | grep -vE '^(src/packages/cnos\.core/(commands/install-wake/|orchestrators/agent-admin/cnos-agent-admin\.golden\.yml|cn\.package\.json)|\.cdd/unreleased/476/|\.github/workflows/install-wake-golden\.yml)$' | wc -l` MUST be 0 (allow CI workflow file if α added new)
- Mechanical AC8 renderer-side: `grep -ciE 'admin.only|disallowed_surfaces|defer.path|cell_execution' src/packages/cnos.core/commands/install-wake/cn-install-wake` MUST be 0
- Mechanical AC8 package-side carve-out unchanged: `grep -ciE 'runs-on|uses:|claude-code-action|GITHUB_TOKEN|GHA' src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json src/packages/cnos.core/orchestrators/agent-admin/prompt.md` MUST equal the same grep run against origin/main at base SHA (no enlargement)
- Mechanical AC6 idempotence (re-render in the merge-test worktree): from inside the worktree, invoke `./src/packages/cnos.core/commands/install-wake/cn-install-wake agent-admin` (or however α wired the local-dev invocation); `sha256sum src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` → record; re-invoke; re-`sha256sum`; the two outputs MUST be byte-identical AND match the committed file's sha (`sha256sum` against the committed version)
- Tear down worktree.

Row 4 — γ artifact completeness: `git ls-tree -r origin/cycle/476 .cdd/unreleased/476/gamma-scaffold.md` returns a non-empty line.

# Review approach (β SKILL Rule 6 + Rule 7)

For each AC (AC1-AC8 per cnos#476 body):
1. Read the AC's Invariant / Oracle / Surface verbatim from the issue body.
2. Re-grep the cycle's diff (NOT α's `self-coherence.md` claims) to verify the oracle passes against actual artifacts. β SKILL Rule 6: code is evidence, doc is hypothesis.
3. Per Rule 6a: if a γ-scaffold oracle regex is brittle (misses a real literal the code emits), widen β-side and record both the scaffold regex and the widened regex in `beta-review.md`. Do NOT request α to reshape implementation to fit a brittle scaffold regex.
4. Per Rule 6b: if a name overpromises (e.g., the renderer's `--out` flag accepts a path but silently rendrers to default; or a function `validate_manifest` only checks presence, not type), flag as name-overpromise.

Specific β audits for this cycle:
- **Audit the package-vs-renderer authority split (AC8) on BOTH sides.**
  - Renderer side: `grep -ciE 'admin.only|disallowed_surfaces|defer.path|cell_execution' src/packages/cnos.core/commands/install-wake/cn-install-wake` MUST be 0. If non-zero, the renderer is leaking package-authority strings into substrate emission — RC with D-severity finding, classification `architecture` (contract violation).
  - Package side: re-run the existing #470 grep against cycle/476 to confirm no enlargement of substrate-bound strings in `wake-provider.json` or `prompt.md`.
- **Audit the idempotence claim (AC6).** Re-render the golden in the merge-test worktree (above) and byte-compare twice. If the second render differs, the renderer has non-determinism (timestamps, UUIDs, sort-order dependencies, environment-variable leakage). RC with D-severity, classification `correctness`.
- **Audit the golden-fixture invariant (AC6).** Re-render at HEAD and `git diff` against the committed golden. If non-empty, the committed golden is stale OR the renderer is non-deterministic. RC.
- **Audit α's claim-class verification tables in §Self-check.** If α has any "all X work" claim WITHOUT a per-X table, RC with D-severity finding, classification `protocol-compliance` (cites cnos#472 friction). Per-row table absence on a class-claim is itself a finding.
- **Audit α's form-pin or render-target-pin override (if any).** If α overrode γ's pin, verify §Gap cites the specific file + line that structurally forced the override. If the override reason is weak (preference, taste, simplicity-without-evidence), RC with D-severity, classification `protocol-compliance`.

Rule 7 — Implementation-contract coherence (binding):
For each of the 7 axes in the dispatch contract (γ scaffold §"Implementation contract"), confirm the diff conforms. The diff hunks must map onto the pinned rows:
- Language: only `.sh` / extension-less shell + `.json` (cn.package.json edit) + `.md` (optional README, CDD artifacts) + `.yml` (golden + optional CI workflow) added; no `.go`, no `.py` (unless α picked python3 jq-fallback and recorded in §Gap).
- CLI integration target: `cn.package.json`'s `commands` map has new `install-wake` entry; `src/go/internal/cli/` unchanged.
- Package scoping: only `src/packages/cnos.core/commands/install-wake/`, `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`, `src/packages/cnos.core/cn.package.json`, `.cdd/unreleased/476/`, and optionally `.github/workflows/install-wake-golden.yml` touched; nothing else.
- Existing-binary disposition: N/A; verify no changes to `daily/cn-daily`, `weekly/cn-weekly`, `save/cn-save`, or `src/go/`.
- Runtime dependencies: shell + `jq` (or recorded python3 fallback); no new entries in any dep manifest.
- JSON/wire contract: rendered YAML matches the `claude-code-action@v1` invocation shape; consumes `cn.wake-provider.v1` per the contract skill.
- Backward compat: `claude-wake.yml` byte-identical; Sub 2 declarations byte-identical (mechanical).

If any axis diverges without a `.cdd/unreleased/476/gamma-clarification.md` updating it, RC with D-severity finding, classification `implementation-contract`.

# Verdict format

Write `.cdd/unreleased/476/beta-review.md` with this structure:

```
# β review — cycle/476

## R1 ({YYYY-MM-DD HH:MM UTC} — base origin/main: {SHA} — head origin/cycle/476: {SHA})

### Verdict: {APPROVE | REQUEST CHANGES | REJECT}

### Contract Integrity
{Identity, branch, base/head SHAs verified}

### Pre-merge gate
| Row | Result | Notes |
| 1 (Identity truth) | pass/fail | |
| 2 (Canonical-skill freshness) | pass/fail | |
| 3 (Non-destructive merge-test) | pass/fail | each validator named |
| 4 (γ artifact completeness) | pass/fail | |

### AC coverage
| AC | Status | Evidence (re-grepped from diff) | Notes |
| AC1 | pass/fail | {form-pin acknowledged / overridden + reason} | |
| AC2 | pass/fail | {positive case run; negative cases run} | |
| AC3 | pass/fail | {YAML parses; structural greps} | |
| AC4 | pass/fail | {admin-only language present in rendered prompt} | |
| AC5 | pass/fail | {per-trigger oracle; per-permission oracle} | |
| AC6 | pass/fail | {golden exists; idempotence (sha twice); CI mechanism} | |
| AC7 | pass/fail | (mechanical: 0 diff lines on claude-wake.yml) | |
| AC8 | pass/fail | (renderer-side grep = 0; package-side carve-out unchanged) | |

### Authority-split audit (AC8 — explicit)
| Side | Audit | Result |
| Renderer (cn-install-wake) | grep for role-decision strings | 0 / non-zero |
| Package (wake-provider.json + prompt.md) | grep for substrate-emission strings | count == #470 baseline |

### Idempotence audit (AC6 — explicit)
| Step | sha256sum |
| Render 1 | {sha} |
| Render 2 | {sha} |
| Committed golden | {sha} |
| All three equal? | yes/no |

### Claim-class verification audit (cnos#472 friction; binding)
| α §Self-check claim | Per-item table present? | If absent, RC severity D classification protocol-compliance |
| "All manifest fields render correctly" | yes/no | |
| "All triggers map" | yes/no | |
| "All permissions correct" | yes/no | |
| "All ACs pass" | yes/no | |
| "Idempotence verified" | yes/no | |

### Implementation-contract coherence (Rule 7)
| Axis | Status | Evidence |
| Language | pass/fail | |
| CLI integration | pass/fail | |
| Package scoping | pass/fail | |
| Existing-binary disposition | pass/fail | |
| Runtime dependencies | pass/fail | |
| JSON/wire contract | pass/fail | |
| Backward compat | pass/fail | |

### Binding findings (if RC or REJECT)
- F1 (severity D/C/B/A; classification {protocol-compliance|implementation-contract|correctness|architecture|...}): {description; surface; oracle; remediation}
- F2 ...

### Approval line (if APPROVE)
"Approved at head {SHA}; merging via `git merge --no-ff origin/cycle/476` to main."
```

# On APPROVE: merge

If R1 (or RN) verdict is APPROVE and pre-merge gate passes:
1. From a clean checkout on `main`: `git fetch origin && git switch main && git pull --ff-only`.
2. `git merge --no-ff origin/cycle/476 -m "Merge pull request #N from usurobor/cycle/476 — cn-wake-install renderer (v0)"` (or use PR-merge if δ opened one and passed the PR URL).
3. Verify the merge SHA. Push: `git push origin main`.
4. Write `.cdd/unreleased/476/beta-closeout.md` (per β SKILL Rule 5 §"Closure discipline" — review summary, implementation assessment, technical review, process observations; NOT release notes — this cycle does not cut a release; NOT a triage table — γ owns that).
5. Commit `beta-closeout.md` on `main`; push.
6. Exit. Do NOT spawn α or γ for close-out; δ re-dispatches.

# On REQUEST CHANGES: exit cleanly

If R1 (or RN) verdict is RC:
1. Ensure `beta-review.md` has the binding findings clearly enumerated with severities + classifications.
2. Push `beta-review.md` to `origin/cycle/476` (NOT to main).
3. Exit. δ will re-dispatch α to fix findings; α appends a fix-round section to `self-coherence.md` and re-signals review-readiness; δ re-dispatches you for R2.

# Refusal conditions (binding)
β MUST refuse and surface to γ via `.cdd/unreleased/476/gamma-clarification.md` (commit on `cycle/476`, then exit without verdict) if:
- The operator (or anyone) attempts to instruct β directly on implementation, scope, or merge — β communicates only via γ-via-artifact-channel during a cycle (β SKILL Rule 1).
- The harness offers β a separate `claude/{slug}-{rand}`-style branch as the review/merge surface — refuse; review verdict lives on `origin/cycle/476` (β SKILL Rule 1).
- The cycle requires β to author the fix it judges — β does not author α's work (β SKILL Rule 1: independence).
- Approve-with-follow-up is requested — β SKILL Rule 4: no "approve with follow-up" except explicitly named design-scope deferrals filed before merge.

# Stop condition
Stop when (a) merge completed + `beta-closeout.md` written and pushed, OR (b) `beta-review.md` with RC verdict pushed to `cycle/476`. In either case, return a short summary: round number, verdict, merge SHA (if APPROVE), findings (if RC), any γ-clarification you filed. δ reads your summary and either re-dispatches α (RC path) or re-dispatches α for close-out (APPROVE path).
```

---

## Non-goals (verbatim from cnos#476 §"Non-goals" + operator constraints)

From cnos#476:
- **NO production cutover** of `claude-wake.yml`. Cutover is a separate cycle.
- **NO CDD dispatch wake provider** — Sub 4.
- **NO δ wake-invoked mode skill** — Sub 5.
- **NO end-to-end smoke test** — Sub 6.
- **NO renderer support for non-GHA substrates** — deferred until needed.
- **NO renderer support for non-claude-code-action backends** — deferred; cnos#452 territory.
- **NO cnos.cdr / cnos.cdw provider declarations** — per-package follow-ons.
- **NO cohere implementation** — cnos#444.
- **NO claim protocol or cell-runtime mechanics** — cnos#454.
- **NO tag / version bump / RELEASE.md** — γ records candidate release note in closeout; does not cut release.

From operator constraints (binding for this scaffold):
- **MUST NOT** render to `.github/workflows/cnos-agent-admin.yml` (would activate a second production wake at the same cron slots as `claude-wake.yml`; doubles channel writes; races the cursor).
- **MUST NOT** edit `cnos.core/orchestrators/agent-admin/wake-provider.json` or `prompt.md` (Sub 2 shipped these; the renderer CONSUMES them).
- **MUST NOT** edit `.github/workflows/claude-wake.yml`.
- **MUST NOT** amend `wake-provider/SKILL.md` in this cycle (friction notes flow to gamma-closeout and may seed a follow-up issue).
- **NO scope-creep into Sub 4 (cdd-dispatch wake provider), Sub 5 (δ wake-invoked mode), or Sub 6 (end-to-end smoke).**
- **Render-target invariant for cutover:** the renderer MUST support `--out <path>` so a future cutover cycle does NOT require re-implementing the renderer.

## Refusal conditions for α (consolidated; mirrored from α dispatch prompt for visibility)

α MUST refuse (via `gamma-clarification.md` and stop) when:
1. The cycle requires editing `.github/workflows/claude-wake.yml`.
2. The cycle requires rendering to `.github/workflows/cnos-agent-admin.yml` directly.
3. The cycle requires editing the Sub 2 declaration files (`orchestrators/agent-admin/wake-provider.json` or `prompt.md`); friction notes about manifest field gaps go to §Debt unless they BLOCK rendering.
4. The cycle requires authoring CDD dispatch wake (Sub 4), δ wake-invoked-mode (Sub 5), end-to-end smoke (Sub 6).
5. The cycle requires non-GHA substrate support or non-claude-code-action backend support.
6. A pinned implementation-contract axis (γ scaffold table) appears unsatisfiable from this scope.
7. An AC oracle from the issue body appears unsatisfiable from the pinned form / scope.
8. The cycle requires editing any path outside §"Expected touched surfaces" allowed list.

α MAY (without refusal) design within the renderer's behavior — the YAML emission shape (within `claude-code-action@v1` shape constraint), the manifest field-to-YAML mapping algorithm, the error-message wording for AC2 negative cases, the choice of `jq` vs `python3` for JSON parsing, the `--out` flag's exact name/spelling — and override γ's form-pin or render-target-pin (with a structural reason recorded). This is `design-and-build` mode.

---

## Friction-log seed (for γ closeout)

The friction-log in `gamma-closeout.md` will accumulate renderer-contract learnings observed during α work + β review. γ watches for these signals during the cycle (consumed AFTER β APPROVE; not actionable mid-cycle):

1. **Manifest field gaps.** Does `cn.wake-provider.v1` carry every field the renderer needs to produce a substrate-correct workflow? If α surfaces a missing field (e.g., the manifest doesn't carry enough cron-slot information; or a permission_intent value doesn't map cleanly to a GHA permissions key), record the field + the rendering need in `gamma-closeout.md` §Friction; γ may file a follow-up issue to amend `wake-provider/SKILL.md` and bump schema to `cn.wake-provider.v2` if breaking.
2. **Substrate-vs-package boundary fuzz.** Are there cases where α had to make a substrate decision that "felt like" a package decision (or vice versa)? E.g., "should the cron slots be hardcoded in the renderer or driven by a `schedule_intent` in the manifest?" Each such moment is a candidate for `wake-provider/SKILL.md §2.5` table-row addition. Record the fuzz point + the resolution α picked.
3. **Idempotence pitfalls.** Did α encounter any non-determinism source (timestamp injection, UUID generation, environment-variable leakage, sort-order instability, file-mtime-dependent output)? Record the source + the mitigation. If the mitigation requires a contract-level guarantee (e.g., "renderer MUST NOT inject timestamps"), record as a candidate `wake-provider/SKILL.md §3.X` rule addition.
4. **Golden-test mechanism shape.** Did the chosen CI mechanism (new workflow vs step in existing) work cleanly, or did α encounter friction (PR-trigger scoping, secret access from forks, CI runner availability)? Record + recommend.
5. **`--out` flag design.** Does the renderer's `--out` flag carry enough flexibility for the future cutover cycle? E.g., does it support emitting to a directory (renderer derives filename) vs requiring an exact file path? Record + recommend for cutover-cycle prep.
6. **β-side audit friction.** Did β struggle to mechanically verify any AC's oracle? E.g., the per-permission grep relies on a YAML-key transformation that's brittle. Record the brittleness + the widened oracle β used.

Other friction items observed mid-cycle that do NOT seed `wake-provider/SKILL.md` amendments (but may seed other things):
- γ-scaffold prompt-clarity gaps (feed Sub 5 / cnos#472 dispatch-prompt-template work, NOT this cycle).
- δ-routing friction (process; not implementation).
- CDD-side process observations (feed `gamma/SKILL.md` or `alpha/SKILL.md` patches via γ's PRA).

---

## Cycle scope sizing (per #476 — kept whole)

| Factor | Reading | Splitting signal? |
|---|---|---|
| (a) New code surface | 1 renderer shell script + 1 golden fixture + 1 CI workflow step + 1 cn.package.json entry; all in cnos.core | No |
| (b) Cross-module breadth | cnos.core only (commands/ + orchestrators/ + cn.package.json) + optionally one `.github/workflows/install-wake-golden.yml` CI file; no other package, no `src/go/`, no `claude-wake.yml` touch | At edge — γ considered split (design cycle then build) but the form is shell + the design call is the form-pin + render-target-pin which γ has resolved here; α's work is mostly mechanical from this scaffold |
| (c) Lifecycle span | design (form + render-target pinned in scaffold) + impl + golden + CI validator | At edge; whole-cycle judged tractable for shell form |
| (d) MCA preconditions | Sub 2's contract converged; γ's form-pin + render-target-pin in this scaffold resolve the remaining design calls | design-and-build |
| (e) Independent shippability | Yes — Sub 4 consumes the renderer; cutover is a separate cycle | n/a |

**Decision: keep whole (design-and-build).** Per cnos#476 §"Cycle scope sizing" provisional decision; γ confirms in scaffold.

---

Filed by γ@cnos (γ-spawned-as-sub-agent for cycle/476, bootstrap-δ path) on 2026-06-21.
