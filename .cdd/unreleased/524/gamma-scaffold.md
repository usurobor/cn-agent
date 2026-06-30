---
cycle: 524
parent_issue: cnos#524
protocol: cds
cycle_branch: cycle/524
base_main_sha: e275d9acecf08705389d6e846d42819f2d68f7f8
head_sha_at_scaffold: e275d9acecf08705389d6e846d42819f2d68f7f8
mode: MCA
role: γ
authored_by: γ@cdd.cnos (wake-invoked δ dispatch, W3 scope)
date: 2026-06-30 (UTC)
output_contract: γ-scaffold + α prompt + β prompt
dispatch_scope: W3-implementation
ac_set: AC3, AC4, AC6, AC7 (W3-relevant)
w0_locked: true
w0_design_ref: .cdd/unreleased/524/w0-design.md
w1_ref: PR #526 (merged; SKILL.md files authored, #Wake CUE schema added)
w2_ref: PR #527 (merged; dual-source parity gate, --parity-check mode, CI step)
operator_directive_ref: cnos#524 comment 2026-06-30T14:37:36Z (IC_kwDORGimc88AAAABIMOirw)
---

# γ-scaffold — cnos#524: wake-as-skill W3 implementation

**Scope:** W3 — renderer source-flip. Flip `cn-install-wake` default source of truth from
`wake-provider.json + prompt.md` to `SKILL.md`. Invert `--parity-check` semantics: it now
proves render(JSON+prompt) == render(SKILL.md) (rather than the W2 direction).

**W0 reference:** `.cdd/unreleased/524/w0-design.md` (locked; governs the design).
**W1 state:** Both wake SKILL.md files exist; `#Wake` CUE schema in `schemas/skill.cue` (PR #526).
**W2 state:** `cn-install-wake` can parse SKILL.md via `--source skill`; parity proved byte-for-byte
including headers; CI parity guard green (PR #527).

---

## §1. W3 scope

**What W3 delivers (the only writable changes):**
1. `src/packages/cnos.core/commands/install-wake/cn-install-wake` — flip default source:
   - `source_type="json"` → `source_type="skill"` (default now reads SKILL.md)
   - `--parity-check` now implies `--source json` (not `--source skill`)
   - Update `--source`, `--parity-check`, and exit-5 doc comments
   - Update parity check success/failure messages
2. `.github/workflows/install-wake-golden.yml` — update W2 parity check step comment/name
   to reflect W3 semantics (now proves render(JSON+prompt) == render(SKILL.md))
3. `.cdd/unreleased/524/` records (this scaffold, self-coherence, beta-review, closeouts)

**What W3 does NOT deliver** (explicitly out of scope per operator directive):
- Delete `wake-provider.json` or `prompt.md` (W4)
- Re-render goldens or change them in any way
- Change live workflows (`cnos-agent-admin.yml`, `cnos-cds-dispatch.yml`)
- Change `schemas/skill.cue`
- Change `wake-provider.json` or `prompt.md`
- Change the SKILL.md files for either wake

**State after W3:** Normal render (`cn install-wake <name>`) reads SKILL.md. JSON+prompt are
still present; parity mode (`--parity-check`) proves they still produce byte-identical output.
Goldens are unchanged (W2 proved byte-identity). W4 (delete JSON+prompt) is the final phase.

---

## §2. AC oracle list (W3-relevant ACs)

### AC3 — Renderer reads SKILL.md (now DEFAULT path)

**Invariant:** After W3, `cn install-wake <name>` (no `--source` flag) reads SKILL.md by
default and produces output byte-identical to the committed golden.

**Oracle (pass conditions):**
- `cn install-wake agent-admin` exits 0, golden unchanged. (CI "Re-render" steps)
- `cn install-wake cds-dispatch` exits 0, golden unchanged. (CI "Re-render" steps)
- `git diff --exit-code` on both goldens after re-render: clean. (CI "Verify goldens unchanged")
- Negative proof: mutating `wake-provider.json` or `prompt.md` in a temp copy does NOT
  change the normal render output. (AC narrative in self-coherence §R0)

**Failure scenario:** Flip to SKILL.md default causes a field to be read differently
(e.g. Python YAML parsing handles a field differently than jq), producing different YAML.
The "Verify goldens unchanged" step fails with a diff.

---

### AC4 — Byte-identical goldens (W3 form: source-flip is a no-op on output)

**Invariant:** W3 source-flip produces zero bytes of change in committed goldens.
The byte-identity oracle applies at W3: rendering from SKILL.md (new default) must produce
output byte-identical to the committed golden (which was itself byte-identical to JSON+prompt
per W2 parity proof).

**Oracle:** `install-wake-golden` CI: "Re-render" → "Verify goldens unchanged" → `git diff`
returns clean for both `cnos-agent-admin.golden.yml` and `cnos-cds-dispatch.golden.yml`.

**Failure scenario:** The golden diff is non-empty after the source flip. This means the W2
parity proof did not hold for some subtle case (e.g. trailing newline difference not caught by
W2's full-byte comparison). STOP and surface; do not re-commit the golden.

---

### AC6 — Renderer refusals preserved from SKILL.md source (now the DEFAULT)

**Invariant:** The renderer's runtime refusals still fire when the default source (SKILL.md)
is used:
- `activation_state: declaration-only` in frontmatter → exit 3
- `role: dispatch + admin_only: false + activation_log_writer: true` → exit 4

**Oracle:** CI "AC5 — declaration-only refusal" and "AC4 (cycle/496) — renderer refusal on
mis-declaration" steps continue to pass. These smokes use `--manifest` with JSON test fixtures,
so they are unaffected by the default source flip (they explicitly supply a manifest file,
bypassing the source selection logic). After W3, the AC5 step for the live `cds-dispatch`
manifest (which has `activation_state: live` in SKILL.md) still renders successfully.

**Parity gate (W3 form):** After flipping `--parity-check` to imply `--source json`, the CI
"W3 parity check" step renders from JSON+prompt and compares against the SKILL.md-rendered
golden. If JSON+prompt still produces byte-identical output, AC6 is satisfied (the refusal
gates were always in the renderer, not in the source path; changing which source is default
does not affect their behavior).

---

### AC7 — All CI gates green

**Invariant:** No inherited-cap, no new red.

**Oracle:** I1/I2/I4/I5/I6/Go/Package/Binary/install-wake-golden/dispatch-repair-preflight all
green on this cell's PR. The "W3 parity check" step must replace the "W2 parity check" step
in install-wake-golden.yml and pass (exit 0 from both `--parity-check` invocations).

**AC8 constraint:** The flip from `source_type="json"` to `source_type="skill"` is a single
string change. It does not add any literal role-decision strings. The AC8/AC7 renderer
authority audit must remain green.

---

## §3. Implementation contract (for α)

### §3.1 `cn-install-wake` changes (minimal — 6 targeted edits)

**3.1.A — Flip default source (line ~298):**
```sh
# BEFORE:
source_type="json"
# AFTER:
source_type="skill"
```

**3.1.B — Invert `--parity-check` source implication (line ~416):**
```sh
# BEFORE:
[ -n "$parity_check" ] && source_type="skill"
# AFTER:
[ -n "$parity_check" ] && source_type="json"
```

**3.1.C — Update `--source` doc comment in header (~lines 55-62):**
Change `` `json` (default) `` to `` `skill` (default) `` and reorder so `skill` is described
first. `` `json` `` becomes the non-default path description.

**3.1.D — Update `--parity-check` doc comment in header (~lines 63-75):**
Change from "Render from SKILL.md source..." to "Render from JSON+prompt source
(wake-provider.json + prompt.md) and compare FULL byte-for-byte against the committed golden
(now produced from SKILL.md by default). Proves render(JSON+prompt) == render(SKILL.md).
Exit 0 if byte-identical; exit 5 if not. Implies --source json."

**3.1.E — Update exit code 5 doc comment (~lines 100-106):**
Change "render(SKILL.md) differs from render(JSON+prompt)" to
"render(JSON+prompt) differs from render(SKILL.md)".

**3.1.F — Update parity-check success/failure messages (in the --parity-check block, ~lines 1125-1129):**
```sh
# BEFORE (success):
echo "cn-install-wake: parity OK — render(SKILL.md) == render(JSON+prompt) for ${name} (full byte identity, headers included)"
# AFTER:
echo "cn-install-wake: parity OK — render(JSON+prompt) == render(SKILL.md) for ${name} (full byte identity, headers included)"

# BEFORE (failure):
echo "cn-install-wake: PARITY FAILURE — render(SKILL.md) differs from render(JSON+prompt) for ${name} (full byte comparison, headers included)" >&2
# AFTER:
echo "cn-install-wake: PARITY FAILURE — render(JSON+prompt) differs from render(SKILL.md) for ${name} (full byte comparison, headers included)" >&2
```

### §3.2 `install-wake-golden.yml` changes (one step update)

Update the "W2 parity check" step:
- Rename step to: `W3 parity check — render(JSON+prompt) == render(SKILL.md)`
- Update the step comment to reflect W3 semantics: `--parity-check` now implies `--source json`;
  the committed goldens are now produced from SKILL.md (new default); parity proves
  render(JSON+prompt) is still byte-identical to render(SKILL.md).

### §3.3 Scope guardrails for α

**MUST change:**
- `src/packages/cnos.core/commands/install-wake/cn-install-wake` (6 targeted edits per §3.1)
- `.github/workflows/install-wake-golden.yml` (one step rename/comment update)
- `.cdd/unreleased/524/` records (this scaffold + self-coherence + beta-review + closeouts)

**MUST NOT change:**
- `.github/workflows/cnos-agent-admin.yml`
- `.github/workflows/cnos-cds-dispatch.yml`
- `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`
- `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`
- `wake-provider.json` (either wake)
- `prompt.md` (either wake)
- `schemas/skill.cue`
- `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md`
- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`
- Any `*.golden.yml` file
- `src/go/` (no Go changes)

**PR linkage invariant:** Every commit message and the PR body MUST use `Refs #524` or
`Part of #524` — NEVER `Closes/Fixes/Resolves/close/fix/resolve #524` (including negated forms
like "does not close #524"). Issue #524 must remain open after W3 merges.

### §3.4 Friction notes (FN) for α

**FN-W3-1: --parity-check flow with new default.** After the flip, `--parity-check` sets
`source_type="json"`. The code path then checks `if [ "$source_type" = "skill" ]; then` for
extracting SKILL.md — this block is now SKIPPED (source is json). The existing JSON path
(`prompt_template_relpath = jq -r '.prompt_template' ...`) runs as normal. The parity comparison
is render(JSON+prompt) vs golden. This is correct.

**FN-W3-2: `display_prompt_path` in skill default.** When `source_type="skill"` (new default),
`json_prompt_path` is derived from `jq -r '.prompt_template' "$json_manifest_path"`. This reads
`wake-provider.json`'s `prompt_template` field to get the display path for the header comment.
The header still shows `prompt.md` path regardless of SKILL.md being the actual source — this
is the stable-header design from W2 (the header is source-stable so parity check can do full
byte comparison). No change needed; this is already correct in the current code.

**FN-W3-3: AC8 audit.** The flip changes only string literals `"json"` and `"skill"`. Neither
is a role-decision string. The AC8/AC7 audit must remain green.

**FN-W3-4: Negative proof method.** To prove the normal render path reads SKILL.md (not JSON),
α should demonstrate in `self-coherence.md §R0` that the `source_type="skill"` branch is taken
for both wakes in a fresh `cn install-wake` invocation, via code inspection. A mechanical
demonstration (mutating JSON in a temp copy and confirming unchanged output) is the gold
standard but may require a shell script in the CI; for W3 the code inspection plus the parity
check prove it sufficiently.

---

## §4. β review prompt

β reviews the W3 implementation against this scaffold. β's primary checks:

1. **Goldens unchanged:** `git diff --exit-code` on both golden files after re-render — clean?
2. **Normal render reads SKILL.md:** Source flip confirmed in code (`source_type="skill"`) +
   `install-wake-golden` CI re-render steps pass with unchanged goldens.
3. **Parity check inverted:** `--parity-check` now implies `--source json`; CI "W3 parity check"
   step passes (exit 0) for both wakes.
4. **No golden/live-workflow change:** Diff excludes
   `cnos-agent-admin.golden.yml`, `cnos-cds-dispatch.golden.yml`,
   `.github/workflows/cnos-agent-admin.yml`, `.github/workflows/cnos-cds-dispatch.yml`.
5. **AC8 audit green:** Renderer authority audit passes (no literal role-decision strings added).
6. **Scope compliance:** Changes confined to `cn-install-wake`, `install-wake-golden.yml`,
   `.cdd/unreleased/524/`. No other files touched.
7. **PR linkage:** PR body uses `Refs #524` / `Part of #524`. Never `Closes/Fixes/Resolves`.
8. **CI green:** I1/I2/I4/I5/I6/Go/Package/Binary/install-wake-golden/dispatch-repair-preflight.
9. **Comments updated:** `--source`, `--parity-check`, exit-5 doc, parity messages, CI step.

β verdict MUST be one of: `converge` (all checks pass; proceed to closeouts + PR) or
`iterate` (specific findings enumerated; α repairs before converge).

---

_Authored by γ@cdd.cnos (wake-invoked δ dispatch, W3 scope), 2026-06-30 (UTC).
W3 operator directive: cnos#524 comment 2026-06-30T14:37:36Z. W0 design: `.cdd/unreleased/524/w0-design.md`.
W2 baseline: PR #527 merged at e275d9ac._
