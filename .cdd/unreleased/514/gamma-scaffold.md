# gamma-scaffold — cnos#514

**Issue:** cnos#514
**Mode:** design-and-build (substantial cycle; 4D migration series)
**Cycle branch:** cycle/514
**Base SHA:** a8b13b71b7bef699fe6b9b955c0b96e56d1ea093
**Scaffolded:** 2026-06-28

---

## 1. Status truth table

| Pass | Issue | Status |
|---|---|---|
| 4A audit | #509 | merged |
| 4B | #511 | merged |
| 4C | #513 | merged |
| **4D (this cycle)** | **#514** | **in progress** |

Prior passes (#509, #511, #513) are merged to main. This cycle has no upstream dependency — 4D may proceed immediately.

---

## 2. Surfaces α will touch

### 2.1. Doc bundle moves (10 bundles)

| Source bundle | Destination | Frozen subdirs to leave in place |
|---|---|---|
| `docs/alpha/protocol/` | `docs/reference/protocol/cn/` | none |
| `docs/alpha/agent-runtime/` | `docs/reference/runtime/` | `3.10.0/`, `3.14.0/`, `3.7.0/`, `3.8.0/` |
| `docs/alpha/runtime-extensions/` | `docs/reference/runtime/extensions/` | `1.0.6/` |
| `docs/alpha/package-system/` | `docs/reference/packages/` | none |
| `docs/alpha/cli/` | `docs/reference/cli/` | none |
| `docs/alpha/ctb/` | `docs/reference/ctb/` | none |
| `docs/alpha/schemas/` | `docs/reference/schemas/` | none (only 2 files: peers.schema.json, protocol-contract.json) |
| `docs/beta/schema/` | `docs/reference/schemas/` | `3.14.4/` |
| `docs/alpha/security/` | `docs/architecture/security/` | none |
| `docs/alpha/cognitive-substrate/` | `docs/architecture/cognitive-substrate/` | none |

**Active files confirmed at docs/alpha/agent-runtime/ root (move these):**
AGENT-RUNTIME.md, CAA.md, CONFIGURE-AGENT.md, CORE-REFACTOR.md, GIT-CN-PACKAGE.md, GO-KERNEL-COMMANDS.md, HYBRID-LLM-ROUTING.md, MEMORY.md, ORCHESTRATORS.md, PLAN-174-orchestrator-runtime.md, POLYGLOT-PACKAGES-AND-PROVIDERS.md, PROVIDER-CONTRACT-v1.md, README.md, RUNTIME-CONTRACT-v2.md

**Active files confirmed at docs/alpha/runtime-extensions/ root (move these):**
README.md, RUNTIME-EXTENSIONS.md

**Active files confirmed at docs/beta/schema/ root (merge into docs/reference/schemas/):**
DESIGN-LLM-SCHEMA.md, README.md

**Pre-cleared schema collision check (AC10):** docs/alpha/schemas/ contains `peers.schema.json` + `protocol-contract.json`; docs/beta/schema/ (active files only) contains `DESIGN-LLM-SCHEMA.md` + `README.md`. No basename collision. Confirmed safe to merge into `docs/reference/schemas/`.

### 2.2. Path repairs (NOT content moves — doc-comment lines only)

| File | Line | Old path | New path |
|---|---|---|---|
| `.github/workflows/build.yml` | 223 | `docs/alpha/schemas/protocol-contract.json` | `docs/reference/schemas/protocol-contract.json` |
| `test/cmd/cn_contract_test.ml` | 3 | `docs/alpha/schemas/protocol-contract.json` | `docs/reference/schemas/protocol-contract.json` |
| `test/cmd/cn_traceability_test.ml` | 5 | `docs/alpha/schemas/protocol-contract.json` | `docs/reference/schemas/protocol-contract.json` |
| `src/ocaml/cmd/cn_workflow.ml` | 4 | `docs/alpha/agent-runtime/ORCHESTRATORS.md` | `docs/reference/runtime/ORCHESTRATORS.md` |
| `src/ocaml/transport/cn_packet.ml` | 8 | `docs/alpha/protocol/MESSAGE-PACKET-TRANSPORT.md` | `docs/reference/protocol/cn/MESSAGE-PACKET-TRANSPORT.md` |

**Note:** `test/cmd/cn_contract_test.ml:16` and `test/cmd/cn_traceability_test.ml:36,46` contain runtime fixture path strings (`"protocol-contract.json"` with no directory prefix) — these reference the adjacent `test/cmd/protocol-contract.json` fixture file which STAYS in place. α must NOT repoint those strings.

**Note:** `src/ocaml/cmd/cn_workflow.ml:20` contains a second reference to `ORCHESTRATORS.md` in prose context — α must assess whether this also needs repointing as an active citation.

### 2.3. Do-NOT-touch files

- `test/cmd/protocol-contract.json` — sibling fixture; stays in place; not a docs/ path
- All frozen/versioned snapshot subdirs within moved bundles (listed in §2.1)
- All `.cdd/` evidence and receipt files
- All wake goldens and golden-bound files
- `docs/alpha/architecture/`, `docs/alpha/design/`, `docs/alpha/doctrine/`, `docs/alpha/essays/`, `docs/alpha/vision/` — outside named 10 bundles; not in scope

---

## 3. AC oracle — per-AC verification approach

| AC | Text | Verification approach |
|---|---|---|
| AC1 | Only 10 named 4D bundles move | `git diff --name-status main...HEAD` shows only paths under the 10 named source bundles (plus build.yml and OCaml repairs). Any path outside scope = STOP. |
| AC2 | `git mv` for history | `git log --follow --diff-filter=R` on moved files shows rename events. Spot-check at least one file per bundle. |
| AC3 | Moved-notice stubs at every old active path | Verify stub exists at every old path (redirect, not a content copy). Every old active file should resolve to a stub citing the new location. |
| AC4 | Frozen content untouched | `git diff --name-only main...HEAD \| grep -E '3\\.10\\.0|3\\.14\\.0|3\\.7\\.0|3\\.8\\.0|1\\.0\\.6|3\\.14\\.4'` returns empty. `.cdd/` files unchanged. |
| AC5 | Active Markdown links repointed | `git grep -nE 'docs/alpha/(protocol\|agent-runtime\|runtime-extensions\|package-system\|cli\|ctb\|schemas\|security\|cognitive-substrate)\|docs/beta/schema' -- docs src .github test` after repair returns only frozen/historical/stub contexts. |
| AC6 | Source/path/workflow/test/OCaml citations repointed | Verify build.yml:223, cn_contract_test.ml:3, cn_traceability_test.ml:5, cn_workflow.ml:4, cn_packet.ml:8 all show new paths. |
| AC7 | Relative links recomputed for new depth | Spot-check files at changed nesting depth (e.g. a file moving from depth 3 to depth 4) to confirm `../`-relative links remain valid. |
| AC8 | build.yml changes path-only | `git diff main...HEAD -- .github/workflows/build.yml` shows exactly one line changed (line 223, path string only). Any other job/matrix/dependency/cache change = STOP. |
| AC9 | OCaml gate | `dune build && dune runtest` (or `opam exec -- dune build && opam exec -- dune runtest`) remains green or fails identically to known baseline. Any NEW OCaml failure = STOP. |
| AC10 | Schema collision gate | Pre-cleared (confirmed: no basename collision between alpha/schemas and beta/schema active files). Verify at merge that `docs/reference/schemas/` contains peers.schema.json, protocol-contract.json, DESIGN-LLM-SCHEMA.md, README.md with no filename conflict. |
| AC11 | Golden protection | `git diff --name-only main...HEAD \| grep -Ei 'golden\|snapshot\|fixture'` returns empty (fixture = `test/cmd/protocol-contract.json` is unchanged). |
| AC12 | Do-not-touch integrity | Any do-not-touch file touched for active-link liveness is listed in receipt as δ-accepted override. |
| AC13 | Stale-reference classification | After all moves, the absolute stale sweep AND relative stale sweep both return only: redirect stubs / frozen historical contexts / deferred explicit / external literal citations. No unclassified active stale refs remain. |
| AC14 | No prose edits | `git diff main...HEAD -- docs/` shows only: file additions (new locations), stub additions (old locations), and link-repair diffs. No prose/semantic reflow. |
| AC15 | Hidden/bidi hygiene | `git diff main...HEAD` and CDD receipt files pass a hidden/bidi/object-replacement sweep. |
| AC16 | Required checks recorded in self-coherence.md | name-status proof; absolute stale grep; relative stale grep; do-not-touch proof; golden no-change proof; workflow diff proof; OCaml dune result; link-check result; honest β/δ repair record — all present. |
| AC17 | Merge posture | Inherited reds (I4/I5/I6) accepted only if proven identical to known baseline. Any new red in Go/OCaml/workflow/schema(I2)/package-verification/link-integrity/golden-protection = STOP. |

---

## 4. Scope guardrails

### 4.1. Frozen subdirs — do not move these

- `docs/alpha/agent-runtime/3.10.0/`
- `docs/alpha/agent-runtime/3.14.0/`
- `docs/alpha/agent-runtime/3.7.0/`
- `docs/alpha/agent-runtime/3.8.0/`
- `docs/alpha/runtime-extensions/1.0.6/`
- `docs/beta/schema/3.14.4/`

### 4.2. Do-not-move bundles (outside named 10)

- `docs/alpha/architecture/`
- `docs/alpha/design/`
- `docs/alpha/doctrine/`
- `docs/alpha/essays/`
- `docs/alpha/vision/`
- Any other `docs/alpha/` or `docs/beta/` path not in the 10-bundle move map

### 4.3. Stop conditions

- `build.yml` change beyond the single path string at line 223
- New OCaml test failure (dune runtest)
- New Go build/test failure (if Go source is touched — expected: none)
- Wake golden change or golden-bound file change
- Schema basename collision in `docs/reference/schemas/` (pre-cleared but re-check at merge)
- Active stale refs that cannot be repaired without prose edits or do-not-touch overrides not pre-authorized
- Frozen snapshot edits required
- Prose/semantic edits in moved docs
- Any doc outside the named 10 bundles requiring a move

---

## 5. Friction notes and pre-verified facts

### 5.1. Schema collision — pre-cleared (AC10)

Confirmed at scaffold time:
- `docs/alpha/schemas/` active files: `peers.schema.json`, `protocol-contract.json`
- `docs/beta/schema/` active files: `DESIGN-LLM-SCHEMA.md`, `README.md`
- No basename collision. Safe merge into `docs/reference/schemas/`.
- `docs/beta/schema/3.14.4/` is frozen; stays in place.

### 5.2. runtime-extensions subfolder collision-avoidance

`docs/alpha/runtime-extensions/` moves to `docs/reference/runtime/extensions/` (a subdirectory of `docs/reference/runtime/`, not a sibling). `docs/alpha/agent-runtime/` moves to `docs/reference/runtime/` (parent). The `1.0.6/` frozen subdir inside runtime-extensions stays at its source path (`docs/alpha/runtime-extensions/1.0.6/`). α must NOT move `1.0.6/` to the new location — only the active files `README.md` and `RUNTIME-EXTENSIONS.md` move.

### 5.3. build.yml path-only constraint (AC8, I2 gate)

`build.yml:223` is currently GREEN. The `protocol-contract-check` job runs `diff docs/alpha/schemas/protocol-contract.json test/cmd/protocol-contract.json`. After the move, line 223 must read `diff docs/reference/schemas/protocol-contract.json test/cmd/protocol-contract.json`. That is the sole permitted change. Any other `build.yml` modification is STOP-class.

### 5.4. OCaml doc-comment repair scope

The 4 OCaml file repairs are doc-comment / string-literal citation repairs only. The executable logic is unchanged. `dune build && dune runtest` must remain green (or fail identically to any known inherited baseline).

### 5.5. Reference and architecture destination dirs exist

`docs/reference/` and `docs/architecture/` both exist at HEAD (confirmed: each contains `README.md`). α does not need to create them; α creates subdirectories within them.

### 5.6. Alpha-surface citation volume

`docs/alpha/agent-runtime` has approximately 116 active refs (issue body). This is the broadest active-surface of any pass. α must run the stale sweep after all moves before declaring review-readiness.

### 5.7. Relative link depth change

Files moving from `docs/alpha/<bundle>/FILE.md` (depth 3 from repo root) to `docs/reference/<bundle>/FILE.md` (depth 3) or `docs/reference/runtime/FILE.md` (depth 3) typically maintain the same depth. However, `docs/alpha/runtime-extensions/` moves to `docs/reference/runtime/extensions/` (depth 4), and `docs/alpha/protocol/` moves to `docs/reference/protocol/cn/` (depth 4). Files at those new deeper paths need relative link adjustments if they cross-link to other docs.

---

## 6. α dispatch prompt

```
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 514 --json title,body,state,comments
Branch: cycle/514
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | N/A — docs migration + YAML path repair + OCaml doc-comment repair only |
| CLI integration target | N/A |
| Package scoping | docs/ tree; .github/workflows/build.yml (path string at line 223 only); src/ocaml/cmd/cn_workflow.ml, src/ocaml/transport/cn_packet.ml, test/cmd/cn_contract_test.ml, test/cmd/cn_traceability_test.ml (doc-comments / comment-string citations only) |
| Existing-binary disposition | N/A |
| Runtime dependencies | None beyond git; dune/opam for OCaml verification |
| JSON/wire contract preservation | protocol-contract.json moves docs path only; no JSON content change; test/cmd/protocol-contract.json fixture stays untouched |
| Backward-compat invariant | All moved active paths must have redirect stubs at old locations; frozen version-stamped snapshots stay in place |

## Pre-verified facts α may use without re-deriving

- Base SHA: a8b13b71b7bef699fe6b9b955c0b96e56d1ea093
- Schema collision (AC10): pre-cleared — no basename collision between alpha/schemas and beta/schema active files
- build.yml:223 is currently GREEN (I2); single path-string change permitted
- Frozen subdirs confirmed: agent-runtime/3.10.0/, 3.14.0/, 3.7.0/, 3.8.0/; runtime-extensions/1.0.6/; beta/schema/3.14.4/
- docs/reference/ and docs/architecture/ destination dirs exist (each has README.md)
- runtime-extensions → docs/reference/runtime/extensions/ (depth 4; deeper than source)
- protocol/ → docs/reference/protocol/cn/ (depth 4; deeper than source)
- test/cmd/protocol-contract.json fixture: DO NOT TOUCH
```

---

## 7. β review prompt

```
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 514 --json title,body,state,comments
Branch: cycle/514

## Review orientation

This is a docs-migration + path-repair cycle (Pass 4D). The implementation contract is in
.cdd/unreleased/514/gamma-scaffold.md §6. Key verification gates (per AC oracle in §3 of
that file):

- AC8 (STOP-class): build.yml diff must show exactly one changed line (line 223, path string only).
- AC9 (STOP-class): dune build && dune runtest must be green or match known inherited baseline.
- AC11 (STOP-class): no golden/snapshot/fixture changes (test/cmd/protocol-contract.json untouched).
- AC13: stale sweep must return only classified refs (stub / frozen / deferred / external).
- AC16: all required checks must appear in self-coherence.md.
- AC17: no new reds in any gate (Go/OCaml/workflow/schema/package-verification/link-integrity/golden).

Run the required checks from the issue body §"Required checks" in order and record results.
```

---

## 8. Expected diff scope

- 10 bundle directories moved (active files only; frozen subdirs stay)
- Old active paths replaced by redirect stubs
- 1 YAML line changed (build.yml:223)
- 5 OCaml/test doc-comment lines changed (4 files)
- Active Markdown links across the repo repointed to new paths
- Relative links inside moved files recomputed where depth changed

Net file count estimate: ~50–80 file renames (git mv) + ~10–20 stub additions + ~1–5 link-repair-only edits in external docs.

---

*gamma-scaffold authored by γ (R0 phase, cycle/514). δ dispatches α and β; γ does not spawn roles directly.*
