# self-coherence — cycle #609

**manifest:** sections = [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]
**completed:** [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]

---

## Gap

**Issue:** [cnos#609](https://github.com/usurobor/cnos/issues/609) — "cds-install Sub 2: generalize `cn install wake` identity (agent/PAT/bot) — successor to #549".

**Mode:** `substantial` (per γ's scaffold `.cdd/unreleased/609/gamma-scaffold.md` — multi-file renderer change + golden/test additions).

**Binding design surface:** `docs/development/design/cn-repo-install-MOCKS.md` (read-only oracle) — Mock C (`cn repo install --dispatch cds` for a non-sigma agent) and Mock E (tenant-portable dispatch wake acquisition). γ's scaffold resolved the issue-body AC1 oracle as byte-identical to Mock C's console example, and pinned this cycle's scope to Mock C's invariants C2/C4/C5 (C1/C3/C6 belong to `cn repo install`, issue #610) plus Mock E's E2/E4.

**Gap being closed:** `cn-install-wake` (the shell renderer that emits `cnos-cds-dispatch.golden.yml` / `.github/workflows/cnos-cds-dispatch.yml`) hard-bound every dispatch-wake render to `sigma` — the literal `SIGMA_WORKFLOW_PAT` secret name (×3), `bot_name`/`bot_id` via a `case "$1" in sigma) ... ; *) die` table with no override, and a `cds-dispatch-sigma` self-reference in the live wake's own prompt body (`cnos.cds/orchestrators/cds-dispatch/SKILL.md`, which every other occurrence of the same concept already substitutes via `{agent}`). A non-sigma agent could not render a usable dispatch wake at all, and even if the identity tables were extended, the wake's own prompt would still leak the literal string `cds-dispatch-sigma` regardless of `--agent`. Separately (Mock E, cnos#606 dogfooding), the rendered "Mechanical recovery scanner" and "Mechanical checkpoint + PR finalizer" steps assumed `src/go/` exists in the installing repo (`cd src/go && go build ...`), which is false for any tenant repo that only vendors cnos packages.

**What closes it:**
- Three new `cn-install-wake` flags — `--workflow-pat-secret <NAME>`, `--bot-name <name>`, `--bot-id <id>` — each supporting both `--flag value` and `--flag=value` forms, matching the existing `--agent`/`--out` style.
- `workflow_pat_secret` resolves to: the flag value; else `SIGMA_WORKFLOW_PAT` when `--agent` is `sigma`; else `die` **before** the render/write phase runs (Mock C2 fail-early — `out_path` is never touched until the final `cp` at the very end of the script, so any `die` earlier leaves no partial `--out` file).
- `bot_name`/`bot_id` resolve to the flag values if given, else fall back to the existing `agent_bot_name()`/`agent_bot_id()` tables (unchanged — they already `die` on unknown non-sigma agents).
- All three literal `SIGMA_WORKFLOW_PAT` occurrences in the render body (checkout token, `claude-code-action` `github_token`, finalizer `GH_TOKEN`) now interpolate `${workflow_pat_secret}` — byte-identical output for `--agent sigma` (default), since `workflow_pat_secret` resolves to the literal string `SIGMA_WORKFLOW_PAT` in that case.
- The "Mechanical recovery scanner" and "Mechanical checkpoint + PR finalizer" steps branch on `$agent`: `sigma` keeps the exact `(cd src/go && go build -o /tmp/cn-scan|cn-finalize ./cmd/cn)` lines (byte-for-byte, per Mock E4); any other agent gets a new "Install cn (tenant-portable — no src/go build)" step that acquires the prebuilt `cn` binary via the repo's existing `install.sh` (`curl -fsSL .../install.sh | BIN_DIR=/tmp/cnos-cn-609 sh`) and both mechanical steps then invoke `/tmp/cnos-cn-609/cn` directly.
- One mechanical, non-behavioral fix in `cnos.cds/orchestrators/cds-dispatch/SKILL.md`'s "Disallowed surfaces" section: the one literal `cds-dispatch-sigma` occurrence (describing this wake's own concurrency group) is now `cds-dispatch-{agent}`, matching the same value's `{agent}`-substituted form used three other times in the same file. No other line touched (the adjacent `agent-admin-sigma` — a different wake's identity — and the historical `cn-sigma:.cn-sigma/logs/20260624.md` citation are untouched, per γ's Friction note 2).
- New CI steps in `.github/workflows/install-wake-golden.yml`: AC1 positive/negative, AC2 two-layer zero-leak (synthetic fixture + scoped live-wake grep), E2/E4.

**Base SHA:** `612a96d6a5cac69486a48b50ca8432daf3bdaac2` (per γ's scaffold; `origin/main` at scaffold time). Verified unchanged at implementation time — no rebase was needed this round.

---

## Skills

**Tier 1 (lifecycle):** `CDD.md` + `cnos.cdd/skills/cdd/alpha/SKILL.md` (this role's execution surface). No separate `plan/SKILL.md` invocation: γ's scaffold already carried a concrete per-flag resolution table, an AC-to-invariant mapping, and an expected-diff-scope list — functioning as the plan. `design/SKILL.md` not separately invoked: the design surface (`cn-repo-install-MOCKS.md`) is pinned, read-only, and pre-existing (landed via cycle/608) — this cycle implements against it, it does not author new design content.

**Tier 2/3:** none loaded beyond the base lifecycle skills — this cycle is a POSIX shell renderer + one prose-string edit + CI YAML, with no Go, no new dependency, and no package-specific Tier-3 skill named by the issue or scaffold. The renderer's own header comment (§"Authority split") and the wake-provider contract it implements (`cnos.core/skills/agent/wake-provider/SKILL.md` §2.5/§3) governed the substrate-vs-package authority split already established by prior cycles (476, 485, 496, 591, 593) — this cycle extends that same split (identity binding is renderer/substrate authority) rather than introducing a new one.

---

## ACs

All commands below were run against the implementation SHA (current branch HEAD before this file's own commit — see §Review-readiness for the exact value).

### AC1 — configurable identity (Mock C), positive + negative (Mock C2)

**Positive** — `cn-install-wake cds-dispatch --agent acme --workflow-pat-secret ACME_WORKFLOW_PAT --bot-name acme-bot --bot-id 12345678 --out /tmp/...` renders successfully (exit 0) and the output contains `ACME_WORKFLOW_PAT`, `acme-bot`, `12345678`, and `cds-dispatch-acme`. Verified locally (exact command CI's new "AC1 (cnos#609) — configurable identity, positive (acme render)" step runs) — all four needles present.

**Negative (Mock C2)** — `cn-install-wake cds-dispatch --agent acme --out /tmp/cnos609-ac1-neg-out.yml` (no other identity flags) exits non-zero (`rc=1`), the `--out` file is **not** created, and stderr names the missing flag (`--workflow-pat-secret is required for agent 'acme' ...`). Verified locally — all three conditions held.

### AC2 — zero sigma leak (Mock C4), two-layer

**Layer (a), synthetic fixture:** a minimal `role: dispatch` SKILL.md fixture (mirrors the AC5/AC2-negative fixture-authoring pattern already in `install-wake-golden.yml` — inline heredoc, `mktemp -d`) with `protocol`/`selector.include` populated, rendered via `--manifest <fixture>` with `--agent acme --workflow-pat-secret ACME_WORKFLOW_PAT --bot-name acme-bot --bot-id 12345678`. `grep -inE 'sigma|SIGMA_WORKFLOW_PAT|41898282'` on the output returned zero matches. Verified locally.

**Layer (b), scoped grep on the real cds-dispatch acme render:** the AC1-positive acme render, scanned for the four renderer-controlled leak tokens only (`SIGMA_WORKFLOW_PAT`, `bot_name: "sigma`, `bot_id: "41898282"`, `cds-dispatch-sigma`) — zero matches. This is the oracle that confirms the SKILL.md `{agent}`-substitution fix actually closes the one real self-reference gap (an unscoped `grep -i sigma` against the live wake's full prompt body would never reach zero, per γ's Friction note 2 — `agent-admin-sigma` and the dated historical citation are legitimate package-owned prose, out of this cycle's scope). Verified locally.

### AC3 (Mock C5) — sigma backward-compat, byte-identical

`./cn-install-wake cds-dispatch` and `./cn-install-wake agent-admin` (default, no new flags) both report `(unchanged)` against the committed goldens; `git diff --stat` on both golden files and on `.github/workflows/cnos-cds-dispatch.yml` shows zero changes. `sha256sum` of the live `.github/workflows/cnos-cds-dispatch.yml` and the per-package golden are identical (`3dee3d1574...`), matching the pre-existing "Verify live cds-dispatch workflow matches golden (sha256)" CI step's invariant. Verified locally, before writing this file.

### AC4 — gates green

`install-wake-golden.yml`'s pre-existing steps (goldens unchanged, idempotence ×2, YAML-parses, substrate-shape ×2, AC5 declaration-only refusal, AC2-negative malformed-SKILL.md, AC4 write-fence ×5, AC8/AC7 renderer-authority audit) were all re-run locally against the changed renderer and pass unchanged (see §Self-check for the exact commands/output). No Go code was touched (scope guardrail), so `go build`/`go test` were not re-run as part of this cycle's own change — CI will still run them as part of the repo-wide gate set, unaffected by this diff.

### E2 — tenant-portable acquisition, no `cd src/go` / `go build ./cmd/cn`

The acme render (non-sigma agent) contains no `cd src/go` and no `go build ./cmd/cn`; it contains an `install.sh`-based acquisition line (`curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | BIN_DIR=/tmp/cnos-cn-609 sh`) in a new "Install cn (tenant-portable — no src/go build)" step, and both the "Mechanical recovery scanner" and "Mechanical checkpoint + PR finalizer" steps invoke `/tmp/cnos-cn-609/cn` directly. Verified locally.

### E4 — sigma still self-builds, byte-identical

The sigma-default `cnos-cds-dispatch.golden.yml` (already re-rendered by the "Re-render cds-dispatch wake" CI step earlier in the same job) still contains `(cd src/go && go build -o /tmp/cn-scan ./cmd/cn)` and `-o /tmp/cn-finalize`. Verified locally via `grep -qF` against the golden file already on disk — no new render needed, consistent with γ's scaffold's AC3/E4 oracle note (this is the same regression gate as AC3).

---

## Self-check

Every AC above is backed by a command actually run this round, reproduced here (trimmed):

```
$ ./cn-install-wake agent-admin
cn-install-wake: agent-admin → .../cnos-agent-admin.golden.yml (unchanged)
$ ./cn-install-wake cds-dispatch
cn-install-wake: cds-dispatch → .../cnos-cds-dispatch.golden.yml (unchanged)
$ git diff --stat -- .../cnos-cds-dispatch.golden.yml .../cnos-agent-admin.golden.yml .github/workflows/cnos-cds-dispatch.yml
(no output — zero diff)
$ sha256sum .github/workflows/cnos-cds-dispatch.yml .../cnos-cds-dispatch.golden.yml
3dee3d15741786d4419f2c45ba97dddbc7962d0ba5f4269dafa83421660017d1  (both — identical)

$ ./cn-install-wake cds-dispatch --agent acme --workflow-pat-secret ACME_WORKFLOW_PAT \
    --bot-name acme-bot --bot-id 12345678 --out /tmp/cnos609-acme-render.yml
cn-install-wake: cds-dispatch → /tmp/cnos609-acme-render.yml (rendered)
# all 4 needles (ACME_WORKFLOW_PAT, acme-bot, 12345678, cds-dispatch-acme) present — AC1 positive OK

$ ./cn-install-wake cds-dispatch --agent acme --out /tmp/cnos609-ac1-neg-out.yml
cn-install-wake: --workflow-pat-secret is required for agent 'acme' (...)
rc=1; /tmp/cnos609-ac1-neg-out.yml does not exist — AC1 negative (Mock C2) OK

# AC2 synthetic fixture (role:dispatch, protocol+selector) rendered with acme + full identity:
grep -inE 'sigma|SIGMA_WORKFLOW_PAT|41898282' → zero matches — AC2 layer (a) OK

# AC2 scoped grep on the acme render from AC1-positive:
grep -qF 'SIGMA_WORKFLOW_PAT' / 'bot_name: "sigma' / 'bot_id: "41898282"' / 'cds-dispatch-sigma'
→ zero matches — AC2 layer (b) OK

# E2: acme render has no 'cd src/go', no 'go build ./cmd/cn', has 'install.sh' — OK
# E4: golden still has 'cd src/go && go build -o /tmp/cn-scan ./cmd/cn' and '-o /tmp/cn-finalize' — OK

$ sh -n src/packages/cnos.core/commands/install-wake/cn-install-wake  → OK (syntax check)
$ python3 -c "import yaml; yaml.safe_load(open('.github/workflows/install-wake-golden.yml'))"  → parses, 25 steps
```

Also re-ran, unchanged, the pre-existing AC5 (declaration-only refusal), AC2-negative (malformed SKILL.md), and the AC8/AC7 renderer-authority-leak audit (`admin_only|disallowed_surfaces|defer_path|cell_execution` and `protocol:cds|cdr|cdw|dispatch:cell|status:todo` greps against the renderer source) — both return zero leaks; my new code introduces none of these strings.

I did not run `go build`/`go test` — no file under `src/go/**` was touched (scope guardrail explicitly named in γ's scaffold), so there is no new Go surface to test.

---

## Debt

1. **Fixed tenant-acquisition path name (`/tmp/cnos-cn-609`) is cycle-numbered.** γ's scaffold explicitly allowed "e.g. `/tmp/cnos-cn-609` or similar — pick one consistent name," but a cycle-numbered temp path baked into a general-purpose renderer reads oddly to a future reader. A follow-up could rename it to something acquisition-purpose-named (e.g. `/tmp/cnos-cn`) with no behavior change — cosmetic only, not a correctness gap, and out of scope to bikeshed further this cycle per γ's explicit allowance.
2. **The tenant-acquisition install.sh step is duplicated per-branch** (once in the scanner block, once in the finalizer block would have been the naive approach) — I instead render exactly one "Install cn" step (inside the scanner's non-sigma branch) that both the scanner and finalizer steps then reuse via the same job/runner filesystem. This relies on the "Mechanical recovery scanner" and "Mechanical checkpoint + PR finalizer" blocks staying co-gated on the same `role == "dispatch"` condition (true today — both read the same `$role` variable) and rendering in the same job. If a future cycle ever de-couples their gating conditions, the finalizer's `/tmp/cnos-cn-609/cn` invocation would need its own acquisition step. Flagging for whoever touches this gating next; not a bug against this cycle's scope (Mock E's own example uses the identical one-install-many-invocations shape).
3. **No `--git-user-name`/`--git-user-email` flags** — explicitly excluded per γ's scaffold Friction note 1 (not in the pinned Mock C/E design surface); not implemented, per explicit instruction not to re-litigate this.
4. **`cn repo install`'s consumption of these new flags is issue #610's job**, not this cycle's — this cycle only lands the renderer-side flags; no repo-installer wiring exists yet to invoke them end-to-end from `cn repo install --dispatch cds`.

---

## CDD Trace

| Surface | Change | AC / invariant |
|---|---|---|
| `src/packages/cnos.core/commands/install-wake/cn-install-wake` | usage-comment doc block for 3 new flags | doc-only, no functional AC |
| same file, arg-parsing block | `--workflow-pat-secret`/`--bot-name`/`--bot-id` (both `--flag value` and `--flag=value` forms) | AC1 |
| same file, "Resolve substrate bindings" section | `workflow_pat_secret` fail-early resolution; `bot_name`/`bot_id` override-then-fallback | AC1, Mock C2 |
| same file, checkout/claude-code-action/finalizer steps | `SIGMA_WORKFLOW_PAT` → `${workflow_pat_secret}` (×3) | AC1, AC2, AC3/Mock C5 |
| same file, "Mechanical recovery scanner" + "Mechanical checkpoint + PR finalizer" blocks | agent-branched: sigma unchanged; non-sigma gets "Install cn" step + `/tmp/cnos-cn-609/cn` invocation | E2, E4 |
| `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` | one literal `cds-dispatch-sigma` → `cds-dispatch-{agent}` in "Disallowed surfaces" | AC2 layer (b) |
| `.github/workflows/install-wake-golden.yml` | 5 new CI steps: AC1 positive, AC1 negative, AC2 synthetic-fixture, AC2 scoped-grep, E2/E4 | AC1, AC2, E2, E4 |

**Mock parity (per `cn-repo-install-MOCKS.md`'s receipt-parity contract, scoped to this sub's owned rows per γ's scaffold):**

| ID | Expectation | Observed | Evidence | Verdict | How |
|---|---|---|---|---|---|
| C2 | Missing identity → fails early, no partial render | `rc=1`, `/tmp/cnos609-ac1-neg-out.yml` not created, stderr names `--workflow-pat-secret` | Local run, this file's §Self-check | match | Resolution happens before any `out_path` write; `out_path` is only touched by the final `cp` at the end of the script. |
| C4 | Zero sigma leak for agent ≠ sigma | Two-layer grep (fixture + scoped live-wake) both zero matches | Local run, this file's §ACs / §Self-check | match | `${workflow_pat_secret}`/`${bot_name}`/`${bot_id}` substitution + the one SKILL.md `{agent}` fix close both the renderer-authority and the one prompt-body self-reference leak. |
| C5 | `--agent sigma` byte-identical to committed golden | `(unchanged)` on both goldens; sha256 match on live vs. golden `cnos-cds-dispatch.yml` | Local run, this file's §AC3 | match | New flags only take effect when explicitly passed or agent ≠ sigma; default path is untouched. |
| E2 | Tenant render has no `cd src/go`/`go build ./cmd/cn`; uses `install.sh` | Confirmed absent/present respectively | Local run, this file's §E2 | match | Acme-branch of both mechanical steps replaced with `install.sh`-based acquisition. |
| E4 | sigma still self-builds, byte-identical | Golden still contains both `go build -o /tmp/cn-scan`/`-o /tmp/cn-finalize` lines | Local run, this file's §E4 | match | sigma branch of both mechanical steps is untouched, verbatim. |

`missed: 0` across this sub's owned rows (C2, C4, C5, E2, E4). C1/C3/C6 (Mock C) and E1/E3 (Mock E) are `cn repo install`'s concern (issue #610), not this renderer's, per γ's scaffold — not scored here.

---

## Review-readiness

**Base SHA:** `612a96d6a5cac69486a48b50ca8432daf3bdaac2` (`origin/main` at scaffold time — confirmed still current; no rebase needed this round).

**Implementation SHA:** run `git log -1 --format='%H' origin/cycle/609` after this file's commit lands, to resolve the exact HEAD (per alpha/SKILL.md §2.6 SHA convention — not naming a value this very commit would advance past).

**Changed files:** `src/packages/cnos.core/commands/install-wake/cn-install-wake`, `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`, `.github/workflows/install-wake-golden.yml`, `.cdd/unreleased/609/self-coherence.md`.

**Verification:** every AC/invariant claim above (C2, C4, C5, E2, E4) was exercised locally by running the exact shell the new CI steps run (not merely writing the CI YAML and trusting it); `sh -n` syntax-checked the renderer; the new CI YAML was parsed with PyYAML (25 steps, parses clean) after fixing two step names that YAML mis-parsed as nested mappings (unquoted `name:` values containing `: ` — renamed to comma-separated phrasing, matching every other existing step name's colon-free convention). Both protected goldens and the live `cnos-cds-dispatch.yml` are confirmed byte-identical (`git diff --stat` empty, sha256 match) — no `src/go/**` file was touched.

**Known debt:** see §Debt above (4 items, none blocking — a cosmetic naming note, a documented single-install-step reuse dependency, two explicitly out-of-scope items already named by γ's scaffold).

**This cycle is ready for β review.**
