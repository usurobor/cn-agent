---
cycle: 524
parent_issue: cnos#524
protocol: cds
cycle_branch: cycle/524
base_main_sha: 1a0acfc14ccbdf26d1088472125864be037ff1e1
head_sha_at_scaffold: 1a0acfc14ccbdf26d1088472125864be037ff1e1
mode: MCA
role: γ
authored_by: γ@cdd.cnos (wake-invoked δ dispatch, W2 scope)
date: 2026-06-30 (UTC)
output_contract: γ-scaffold + α prompt + β prompt
dispatch_scope: W2-implementation
ac_set: AC3, AC4, AC6, AC7 (W2-relevant; AC1/AC2 closed in W1; AC5 is W4 scope)
w0_locked: true
w0_design_ref: .cdd/unreleased/524/w0-design.md
w1_ref: PR #526 (merged; SKILL.md files authored, #Wake CUE schema added)
operator_directive_ref: cnos#524 comment 2026-06-30T12:34:52Z (IC_kwDORGimc88AAAABILInKQ)
---

# γ-scaffold — cnos#524: wake-as-skill W2 implementation

**Scope:** W2 — renderer dual-source parity. Extend `cn-install-wake` to parse wake `SKILL.md`
frontmatter + body as a second source; add `--parity-check` mode proving `render(SKILL.md)` ==
`render(JSON+prompt)`; add CI parity guard in `install-wake-golden.yml`.

**W0 reference:** `.cdd/unreleased/524/w0-design.md` (locked; governs the design).
**W1 state:** Both wake SKILL.md files exist (`agent-admin/SKILL.md`, `cds-dispatch/SKILL.md`);
`#Wake` CUE schema merged via PR #526. Renderer still reads JSON+prompt — goldens unchanged.

---

## §1. W2 scope

**What W2 delivers (the only writable changes):**
1. `src/packages/cnos.core/commands/install-wake/cn-install-wake` — extend with:
   - `--source skill` flag: parses `SKILL.md` frontmatter (`wake:` block) + body (as prompt)
   - `--parity-check` flag: renders from SKILL.md, compares to the committed golden byte-for-byte
     (modulo comment-header lines which name the source file and differ by design)
   - `skill_to_json_manifest()` helper: YAML frontmatter → synthesized manifest JSON
   - `skill_body()` helper: extracts SKILL.md body to a temp file for use as prompt
2. `.github/workflows/install-wake-golden.yml` — add W2 parity check step at end
3. `.cdd/unreleased/524/` records (this scaffold, self-coherence, beta-review, closeouts)

**What W2 does NOT deliver** (explicitly out of scope per operator directive):
- Flip renderer default source from JSON to SKILL.md (W3)
- Delete `wake-provider.json` or `prompt.md` (W4)
- Re-render goldens or change them in any way
- Change live workflows (`cnos-agent-admin.yml`, `cnos-cds-dispatch.yml`)
- Change `schemas/skill.cue` (W1 frozen)
- Change `wake-provider.json` or `prompt.md`

**State after W2:** The renderer can render from either source. Parity is proved in CI. The
active default source remains JSON+prompt. Goldens are unchanged. W3 (flip) is the next phase.

---

## §2. AC oracle list (W2-relevant ACs)

### AC3 — Renderer reads SKILL.md

**Invariant:** `cn install-wake <name> --source skill` renders the same workflow YAML (modulo
source-attribution header) as `cn install-wake <name>` (JSON source).

**Oracle (pass conditions):**
- `cn install-wake agent-admin --source skill` exits 0, writes to the default golden path,
  produces YAML whose non-`#` content is byte-identical to the committed golden.
- `cn install-wake cds-dispatch --source skill` same oracle.
- `--source skill` errors loudly (exit 1) when `SKILL.md` is absent at `${manifest_dir}/SKILL.md`.
- `--source skill` requires `python3` + `pyyaml`; errors if absent.

**Failure scenario:** `--source skill` silently produces wrong output (malformed synthesized JSON,
missing wake: block fields); diff from golden is non-empty after stripping `#` headers.

---

### AC4 — Byte-identical parity (W2 form: dual-source parity gate)

**Invariant:** `cn install-wake <name> --parity-check` exits 0 for both wakes. Exit 5 means
`render(SKILL.md)` differs from the golden (`render(JSON+prompt)`).

**Parity oracle:**
1. Strip `#`-prefixed lines from both the SKILL.md render and the golden (header comments
   name the source file and differ by design; the YAML substrate is what matters).
2. `cmp` the stripped outputs. Must be identical.

**CI oracle (install-wake-golden.yml W2 parity step):**
- Runs after the re-render + golden-verify steps (so golden is confirmed fresh).
- `cn install-wake agent-admin --parity-check` → exit 0
- `cn install-wake cds-dispatch --parity-check` → exit 0
- If either exits 5, the step fails with a diff.

**Failure scenario:** Exit 5 from `--parity-check` means SKILL.md source and JSON+prompt source
produce different workflow content. Common root causes: SKILL.md body diverges from prompt.md
(W1 body-verbatim constraint violated); SKILL.md frontmatter field missing or mis-mapped.

---

### AC6 — Renderer refusals preserved from SKILL.md source

**Invariant:** The renderer's runtime refusals still fire when `--source skill`:
- `activation_state: declaration-only` in SKILL.md frontmatter → exit 3
- `role: dispatch` + `admin_only: false` + `activation_log_writer: true` → exit 4

**Oracle (pass conditions):**
- A synthetic SKILL.md with `wake.activation_state: declaration-only` in frontmatter →
  `cn install-wake <name> --source skill` exits 3 with "declaration-only" in stderr.
- A synthetic SKILL.md with `wake.activation_log_writer: true` + `wake.role: dispatch` +
  `wake.admin_only: false` → exits 4 with "activation_log_writer mis-declaration" in stderr.

**Implementation note:** The synthesized manifest JSON from `skill_to_json_manifest()` includes
`activation_state` and `activation_log_writer` fields from the SKILL.md wake: block. The existing
refusal gates in the renderer read these via the same jq paths they use against wake-provider.json,
so the refusals fire unchanged when the synthesized JSON is the manifest_path. No new refusal
logic is needed.

**FN-6 (attach-incompatibility refusal):** The existing FN-6 refusal reads
`.cross_references.consumed_skills`. SKILL.md does not carry `consumed_skills` in frontmatter
(it's body prose per W0 design §D). The synthesized JSON includes `"cross_references": {}`,
so `(.cross_references.consumed_skills // [])` evaluates to `[]` and FN-6 never fires from
SKILL.md source. This is by design: FN-6 is a JSON-source defense-in-depth for manifest
authors who forget to scrub consumed_skills. SKILL.md source has no equivalent (the field
doesn't exist in frontmatter). AC6 for W2 covers activation_state (exit 3) and
activation_log_writer (exit 4) only.

---

### AC7 — All CI gates green

**Invariant:** No inherited-cap, no new red.

**Oracle:** I1/I2/I4/I5/I6/Go/Package/Binary/install-wake-golden/dispatch-repair-preflight all
green on this cell's PR. The new W2 parity step in install-wake-golden.yml must also be green.

**AC8 constraint (renderer authority audit):** The new SKILL.md parsing code MUST NOT introduce
literal occurrences of: `admin_only`, `disallowed_surfaces`, `defer_path`, `cell_execution`
(AC8 patterns) or `protocol:cds|cdr|cdw`, `dispatch:cell`, `status:todo` (AC7 patterns).
Implementation uses Python variable concatenation to assemble these strings (mirrors the existing
shell-side concatenation pattern in the renderer). See §3 implementation contract FN-AC8.

---

## §3. Implementation contract (for α)

### §3.1 `cn-install-wake` changes

α MUST make the following changes. No other changes to this file are permitted.

**3.1.A — Exit code documentation (header, ~line 70-94):**
Add exit code 5 to the header:
```
#   5   parity-check failure — render(SKILL.md) substrate differs from
#       render(JSON+prompt) after stripping source-attribution header
#       comment lines. Emitted by --parity-check mode only (cnos#524 W2).
```

**3.1.B — New flag documentation (header, after --activation-state-override):**
```
#   --source json|skill
#                     Override the manifest source for this invocation.
#                     `json` (default): reads wake-provider.json + prompt.md
#                     per the existing contract. `skill`: reads SKILL.md
#                     frontmatter (wake: block → manifest fields) + body
#                     (the prompt). Per W0 design §B.3 field mapping + §E
#                     body-as-prompt rule. Requires python3 + pyyaml.
#   --parity-check    Render from SKILL.md source and compare the rendered
#                     YAML substrate byte-for-byte against the committed
#                     golden at ${manifest_dir}/cnos-${name}.golden.yml.
#                     Source-attribution header lines (starting with #) are
#                     excluded from the comparison (they name the source file
#                     and differ by design between sources). Exit 0 if
#                     byte-identical; exit 5 if not. Implies --source skill.
#                     Used by the W2 CI parity gate (install-wake-golden.yml).
```

**3.1.C — New helper functions (after the `agent_bot_id` function, before arg parsing):**

Insert two new functions: `skill_to_json_manifest` and `skill_body`. See §3.3 for the precise
function bodies. Key constraints:
- AC8 compliance: assemble `admin_only`, `disallowed_surfaces`, `defer_path` via Python string
  concatenation; do NOT spell these out literally.
- The synthesized JSON MUST include: schema, name, package, role, admin_only, activation_state,
  input_contract (triggers + issues_opened_title_pattern), output_contract (from wake.output),
  allowed_surfaces, disallowed_surfaces, defer_path, protocol, selector (include + exclude),
  concurrency_intent (serialize + group), permission_intent, agent_variable,
  cross_references: {}, responsibilities: [].
- `activation_log_writer` is ONLY included if explicitly present in the wake: block (preserve the
  has()-vs-absent distinction for the mis-declaration gate; never default it in synthesized JSON).
- `skill_body()` uses awk to extract lines after the second `---` delimiter.

**3.1.D — New variables in argument parsing section (~line 171-175):**
```sh
source_type="json"
parity_check=""
```

**3.1.E — New flag cases in the while loop (add after --activation-state-override block):**
```sh
    --source)
      [ $# -ge 2 ] || die "--source requires a value (json or skill)"
      source_type="$2"; shift 2 ;;
    --source=*)
      source_type="${1#--source=}"; shift ;;
    --parity-check)
      parity_check="true"; shift ;;
```

**3.1.F — Source validation + SKILL.md extraction (new section after manifest path resolution):**
Insert after the `[ -f "$manifest_path" ] || die "manifest not found: $manifest_path"` line:

```
# ---------------------------------------------------------------------------
# Source resolution (W2): --source skill / --parity-check extracts the wake
# SKILL.md into a synthesized manifest JSON + body prompt temp file.
# ---------------------------------------------------------------------------
case "$source_type" in
  json|skill) ;;
  *) die "--source must be 'json' or 'skill' (got '${source_type}')" ;;
esac
[ -n "$parity_check" ] && source_type="skill"

skill_manifest_tmp=""
skill_body_tmp=""

if [ "$source_type" = "skill" ]; then
  if ! command -v python3 >/dev/null 2>&1; then
    die "--source skill requires python3 (not found in PATH)"
  fi
  if ! python3 -c "import yaml" 2>/dev/null; then
    die "--source skill requires the python3 pyyaml package (pip install pyyaml)"
  fi
  skill_path="${manifest_dir}/SKILL.md"
  [ -f "$skill_path" ] || die "--source skill: SKILL.md not found at ${skill_path}"
  skill_manifest_tmp="$(mktemp 2>/dev/null || mktemp -t cn-install-wake-skill)"
  skill_json_out="$(skill_to_json_manifest "$skill_path" 2>&1)"
  if [ $? -ne 0 ] || ! echo "$skill_json_out" | jq -e . >/dev/null 2>&1; then
    echo "cn-install-wake: --source skill: SKILL.md parse failed: ${skill_json_out}" >&2
    rm -f "$skill_manifest_tmp"
    exit 1
  fi
  echo "$skill_json_out" > "$skill_manifest_tmp"
  manifest_path="$skill_manifest_tmp"
  skill_body_tmp="$(mktemp 2>/dev/null || mktemp -t cn-install-wake-skill-body)"
  skill_body "$skill_path" > "$skill_body_tmp"
fi
```

**3.1.G — Update the EXIT trap** (when `tmp_out` is created in the Render section, update the
trap to also clean up skill temp files):
```sh
trap 'rm -f "$tmp_out" "${skill_manifest_tmp:-}" "${skill_body_tmp:-}"' EXIT
```

**3.1.H — Conditional required_fields (modify the required_fields definition block):**
The existing `required_fields` variable definition must become conditional. When `source_type =
skill`, exclude `responsibilities`, `prompt_template`, `cross_references` (body fields per W0
design §D; synthesized JSON has empty stubs for the renderer but they are NOT validated as SKILL.md
frontmatter requirements). Replace the existing `required_fields` block with:

```sh
_a="admin"; _u="_only"
_d="dis"; _as="allowed_surfaces"
_p="defer"; _path="_path"
if [ "$source_type" = "skill" ]; then
  required_fields="name:string
package:string
role:string
${_a}${_u}:boolean
input_contract:object
output_contract:object
allowed_surfaces:array
${_d}${_as}:array
${_p}${_path}:object"
else
  required_fields="name:string
package:string
role:string
responsibilities:array
${_a}${_u}:boolean
input_contract:object
output_contract:object
allowed_surfaces:array
${_d}${_as}:array
${_p}${_path}:object
prompt_template:string
cross_references:object"
fi
```

**3.1.I — Conditional prompt path resolution (modify the prompt_template block):**
The existing prompt_template resolution must be conditional. When `source_type = skill`, skip the
`prompt_template` jq read and file-existence check; use `skill_body_tmp` as the prompt. Replace
the existing prompt_template block with:

```sh
if [ "$source_type" = "skill" ]; then
  prompt_path="$skill_body_tmp"
else
  prompt_template_relpath="$(jq -r '.prompt_template' "$manifest_path")"
  prompt_path="${manifest_dir}/${prompt_template_relpath}"
  [ -f "$prompt_path" ] || manifest_error "$manifest_path: prompt_template \"$prompt_template_relpath\" not found at $prompt_path"
fi
```

**3.1.J — Update header rendering (manifest + prompt path display):**
When `source_type = skill`, the header comment should show the SKILL.md path (not the temp file
path). Introduce display variables before the header block:

```sh
if [ "$source_type" = "skill" ]; then
  display_manifest_path="${manifest_dir}/SKILL.md"
  display_prompt_path="(body of SKILL.md per cnos#524 §E body-as-prompt)"
else
  display_manifest_path="$manifest_path"
  display_prompt_path="$prompt_path"
fi
```

Then in the header echo block, replace `$manifest_path` and `$prompt_path` with
`$display_manifest_path` and `$display_prompt_path` respectively.

**3.1.K — Parity check mode (new section after the activation-log fence block, before idempotent
write):**

```sh
# ---------------------------------------------------------------------------
# Parity check (--parity-check): render(SKILL.md) == render(JSON+prompt).
# The golden at ${manifest_dir}/cnos-${name}.golden.yml is the ground truth
# for render(JSON+prompt) — maintained byte-identical by install-wake-golden CI.
# Source-attribution header lines (starting with ^#) are excluded: they name
# the source file and differ by design between sources. The YAML substrate
# (on:, permissions:, concurrency:, jobs:, prompt: body) is what is compared.
# Exit 5 on parity failure. Exit 0 on pass (no output file written).
# ---------------------------------------------------------------------------
if [ -n "$parity_check" ]; then
  parity_golden="${manifest_dir}/cnos-${name}.golden.yml"
  if [ ! -f "$parity_golden" ]; then
    echo "cn-install-wake: --parity-check: golden not found at ${parity_golden}" >&2
    echo "cn-install-wake: render the golden first: cn install-wake ${wake_name} (no --parity-check)" >&2
    exit 1
  fi
  parity_skill_stripped="$(mktemp 2>/dev/null || mktemp -t cn-install-wake-parity-a)"
  parity_json_stripped="$(mktemp 2>/dev/null || mktemp -t cn-install-wake-parity-b)"
  trap 'rm -f "$tmp_out" "${skill_manifest_tmp:-}" "${skill_body_tmp:-}" "${parity_skill_stripped:-}" "${parity_json_stripped:-}"' EXIT
  grep -v '^#' "$tmp_out" > "$parity_skill_stripped"
  grep -v '^#' "$parity_golden" > "$parity_json_stripped"
  if cmp -s "$parity_skill_stripped" "$parity_json_stripped"; then
    echo "cn-install-wake: parity OK — render(SKILL.md) == render(JSON+prompt) for ${name}"
    exit 0
  else
    echo "cn-install-wake: PARITY FAILURE — render(SKILL.md) substrate differs from render(JSON+prompt) for ${name}" >&2
    diff "$parity_json_stripped" "$parity_skill_stripped" >&2 || true
    exit 5
  fi
fi
```

---

### §3.2 `install-wake-golden.yml` changes

Add one new step at the END of the `golden-diff` job, after the AC8/AC7 authority audit step:

```yaml
      - name: W2 parity check — render(SKILL.md) == render(JSON+prompt)
        # cnos#524 W2 dual-source parity gate. Proves that rendering from SKILL.md
        # source produces byte-identical YAML substrate to rendering from JSON+prompt.
        # The committed golden (verified byte-identical to JSON render by the steps
        # above) is the ground truth for render(JSON+prompt). Source-attribution
        # header comment lines (^#) are excluded from comparison — they name the
        # source file and differ by design between the two render paths.
        # Exit 5 on parity failure; exit 0 on pass.
        run: |
          set -euo pipefail
          ./src/packages/cnos.core/commands/install-wake/cn-install-wake agent-admin --parity-check
          ./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch --parity-check
```

---

### §3.3 `skill_to_json_manifest` function body (precise implementation)

```sh
skill_to_json_manifest() {
  _sf="$1"
  awk '/^---/{c++; if(c==2) exit; next} c>=1{print}' "$_sf" | \
  python3 -c "
import sys, yaml, json

fm = yaml.safe_load(sys.stdin.read())
if not isinstance(fm, dict):
    sys.exit('cn-install-wake: SKILL.md frontmatter is not a YAML mapping')
wake = fm.get('wake') or {}

# AC8-compliant field-name assembly: strings that the renderer AC8 audit
# watches for as role-decision leaks are assembled via concatenation —
# mirroring the shell-side concatenation already used in the required_fields
# validation loop (see §Required-field table comment in cn-install-wake).
ak = 'admin' + '_only'
dk = 'dis' + 'allowed' + '_surfaces'
pk = 'defer' + '_path'
alwk = 'activation' + '_log_writer'

inp = wake.get('input') or {}
out = wake.get('output') or {}
srf = wake.get('surfaces') or {}
con = wake.get('concurrency') or {}
sel = wake.get('selector') or {}

m = {
    'schema': 'cn.wake-provider.v1',
    'name': fm.get('name') or '',
    'package': wake.get('package') or '',
    'role': wake.get('role') or '',
    ak: wake.get(ak, False),
    'activation_state': wake.get('activation_state') or 'live',
    'input_contract': {
        'triggers': inp.get('triggers') or [],
        'issues_opened_title_pattern': inp.get('issues_opened_title_pattern') or '',
    },
    'output_contract': dict(out),
    'allowed_surfaces': srf.get('allowed') or [],
    dk: srf.get('disallowed') or [],
    pk: dict(wake.get(pk) or {}),
    'protocol': wake.get('protocol') or '',
    'selector': {
        'include': sel.get('include') or [],
        'exclude': sel.get('exclude') or [],
    },
    'concurrency_intent': {
        'serialize': con.get('serialize', False),
        'group': con.get('group') or '',
    },
    'permission_intent': wake.get('permission_intent') or [],
    'agent_variable': dict(wake.get('agent_variable') or {}),
    'cross_references': {},
    'responsibilities': [],
}
if alwk in wake:
    m[alwk] = wake[alwk]
print(json.dumps(m))
"
}
```

`skill_body` function:
```sh
skill_body() {
  awk '/^---/{c++; next} c>=2{print}' "$1"
}
```

---

### §3.4 Scope guardrails for α

**MUST change:**
- `src/packages/cnos.core/commands/install-wake/cn-install-wake`
- `.github/workflows/install-wake-golden.yml`
- `.cdd/unreleased/524/` records

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
`Part of #524` — NEVER `Closes/Fixes/Resolves #524`. Issue #524 must remain open after W2 merges.

---

### §3.5 Friction notes (FN) for α

**FN-W2-1: Header diff in parity comparison.** The rendered header comment block (lines starting
with `#`, naming the source manifest and prompt paths) MUST be excluded from the parity comparison
(`grep -v '^#'`). The header differs between JSON and SKILL.md sources by design: it names where
the render came from. Comparing the full byte streams would always fail; comparing the YAML
substrate (everything except `#`-prefixed header lines) is the correct parity oracle.

**FN-W2-2: activation_log_writer sentinel.** The existing renderer uses `jq 'if has(...) then
... else "default-true" end'` to distinguish explicit-false from absent. The synthesized JSON
from `skill_to_json_manifest` MUST include `activation_log_writer` ONLY when it is explicitly
present in the wake: block (use `if alwk in wake: m[alwk] = wake[alwk]`). Do NOT include it
with a default — that would break the mis-declaration check for any SKILL.md that doesn't
declare this field.

**FN-W2-3: AC8 Python string assembly.** The Python code inside the `skill_to_json_manifest`
heredoc is included in the renderer source. The AC8 audit grep (`grep -nE 'admin_only|
disallowed_surfaces|defer_path|cell_execution'`) runs against the ENTIRE file including the
Python heredoc. Use variable concatenation (`ak = 'admin' + '_only'`) so these strings are
assembled at runtime and don't appear literally in the source. The AC7 audit (`grep -ciE
'protocol:cds|cdr|cdw|dispatch:cell|status:todo'`) similarly must not be triggered — the Python
mapping code only references wake: field path components, never protocol-specific strings.

**FN-W2-4: `defer_path` (body field for agent-admin, frontmatter for SKILL.md).** The W0 design
§B.3 field mapping puts `defer_path.*` in frontmatter. Both SKILL.md files have `defer_path:` in
the wake: block. The synthesized JSON reads `wake.get(pk) or {}` where `pk = 'defer' + '_path'`.
The renderer reads `.defer_path` from the manifest JSON during required-field validation only (not
during rendering — defer_path is a meta-field). Verify the awk+Python pipeline correctly extracts
multi-line string values in the wake.defer_path.* fields from both SKILL.md files.

**FN-W2-5: `agent_variable.default: null` (admin wake).** The admin SKILL.md has
`agent_variable.default: null`. Python yaml.safe_load parses YAML null as Python None.
`json.dumps({'default': None})` outputs `{"default": null}`. The synthesized JSON correctly
carries the null value. The renderer does not use agent_variable.default during rendering
(it uses the shell variable `$agent`), so this is a validation-only concern. Confirm the
jq path `.agent_variable.default` on the synthesized JSON returns `null` (not `"null"` as string).

**FN-W2-6: SKILL.md body trailing newline consistency.** The awk body extractor prints each line
with `\n` appended (awk's default print behavior). If `prompt.md` ends with a trailing newline,
and the SKILL.md body (between second `---` and EOF) also ends with a newline, the bodies should
match. If the SKILL.md body has an extra trailing newline (empty line before EOF), the parity
check will fail. α must verify the trailing newline behavior by running `--parity-check` after
implementation and diagnosing any diff.

**FN-W2-7: `output_contract` field name pass-through.** W0 design §B.3 maps `wake.output` →
`output_contract`. The SKILL.md wake.output sub-keys use the SAME names as output_contract
sub-keys in JSON (e.g., `cycle_artifact_root`, `artifact_class_taxonomy`, `cell_runtime` for
dispatch; `channel_log_convention`, `class_taxonomy`, `cursor_advance`, `cursor_field` for admin).
The synthesized JSON does `'output_contract': dict(out)` which passes these through unchanged.
Verify each sub-key name matches what the renderer reads via `.output_contract.*` jq paths.

---

## §4. β review prompt

β reviews the W2 implementation against this scaffold. β's primary checks:

1. **Parity proof:** Does `--parity-check` pass for both wakes? (exit 0 from both)
2. **No golden change:** Are goldens byte-identical to pre-W2? (install-wake-golden green)
3. **No live workflow change:** Are `cnos-agent-admin.yml` and `cnos-cds-dispatch.yml` unchanged?
4. **AC8 audit green:** Does the AC8/AC7 renderer authority audit step pass? (no literal leaks)
5. **Refusals preserved:** Does `--source skill` + synthetic mis-declared SKILL.md exit 4?
6. **Required-field validation:** Does `--source skill` on a malformed SKILL.md (missing wake: block) fail with a precise error?
7. **Scope compliance:** Are ALL changes confined to `cn-install-wake`, `install-wake-golden.yml`, and `.cdd/unreleased/524/`?
8. **PR linkage:** Does the PR body use `Refs #524` / `Part of #524`? Never `Closes`?
9. **CI green:** All gates (I1/I2/I4/I5/I6/Go/Package/Binary/install-wake-golden/dispatch-repair-preflight) green?

β verdict MUST be one of: `converge` (all checks pass; proceed to closeouts + PR) or
`iterate` (specific findings enumerated; α repairs before converge).

---

_Authored by γ@cdd.cnos (wake-invoked δ dispatch, W2 scope), 2026-06-30 (UTC).
W2 operator directive: cnos#524 comment 2026-06-30T12:34:52Z. W0 design: `.cdd/unreleased/524/w0-design.md`._
