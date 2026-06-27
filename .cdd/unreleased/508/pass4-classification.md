---
cycle: 508
artifact: pass4-classification.md
role: alpha
ac: AC2
ac1-total: 663
---

# Pass 4A — Reference Classification (AC2)

**AC1 total:** 663 lines in `pass4-path-inventory.txt`

## Taxonomy (priority order, highest first)

| Priority | Class | Rule |
|---|---|---|
| 1 | `generated/golden-bound` | Reference inside a `*.golden.yml` or byte-checked generated artifact |
| 2 | `frozen/historical` | Reference in a `.cdd/` record, version-stamped snapshot (`**/X.Y.Z/`), historical PRA/audit doc, or changelog entry |
| 3 | `test-fixture` | Path reference in `test/`, `tests/`, `.github/` workflow, or a test script (`test-*.sh`) |
| 4 | `source/package-doctrine` | Code-level path dependency in `src/` (non-golden, non-test) |
| 5 | `markdown-link` | Markdown hyperlink syntax `[text](url)` where url contains the triad path |
| 6 | `inline-path-citation` | Prose or comment path string; not a hyperlink (fallback) |

**Note:** A `.cdd/` record is `frozen/historical` because of its PATH, not its syntax. A version-stamped snapshot directory (e.g. `docs/gamma/cdd/3.14.5/`, `docs/beta/architecture/3.14.4/`, `docs/alpha/agent-runtime/3.7.0/`) is frozen by path. `generated/golden-bound` takes priority over `source/package-doctrine` for `*.golden.yml` files.

---

## Per-Class Counts (sum = 663) ✅

| Class | Count | Source |
|---|---|---|
| `generated/golden-bound` | 3 | `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` (1 line) + `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` (2 lines) |
| `frozen/historical` | 352 | All `docs/**/{X.Y.Z}/` version-stamped snapshot paths |
| `test-fixture` | 10 | `.github/workflows/` (4) + `test/` (3) + `src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh` (3) |
| `source/package-doctrine` | 59 | `src/` files (65 total) minus 3 golden-bound minus 3 test-fixture scripts |
| `markdown-link` | 31 | Non-versioned `docs/` and `schemas/` lines with `](` hyperlink syntax |
| `inline-path-citation` | 208 | Non-versioned `docs/` and `schemas/` lines without hyperlink syntax (prose/comment/frontmatter citations) |
| **TOTAL** | **663** | ✅ matches AC1 |

---

## Class 1 — `generated/golden-bound` (3 lines)

These lines come from `*.golden.yml` files which are rendered/byte-checked artifacts. They take the highest priority class regardless of syntax.

```
src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml:225:            - **Channel log convention (read-only for dispatch):** [`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`](../../../../../docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md) — cited to document the admin/dispatch writer-locality split; the dispatch wake does NOT write channel entries.
src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml:82:            4. **Status reporting.** Append today's channel log entry to `.cn-sigma/logs/$(date -u +%Y%m%d).md` per the entry format in [`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`](../../../../../docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md) §5. [...]
src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml:168:            - **Channel log convention:** [`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`](../../../../../docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md) — entry format, cursor mechanics, class taxonomy.
```

All 3 lines reference: `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`

---

## Class 2 — `frozen/historical` (352 lines)

All lines from files inside version-stamped snapshot directories. These directories contain historical snapshots and PRA records that must not be rewritten.

**Versioned snapshot directories found (partial list — see pass4-do-not-move.md for full enumeration):**

`docs/alpha/agent-runtime/` snapshots: `3.7.0/`, `3.8.0/`, `3.10.0/`, `3.14.0/`
`docs/alpha/runtime-extensions/` snapshots: `1.0.6/`
`docs/beta/architecture/` snapshots: `3.14.4/`
`docs/beta/governance/` snapshots: `3.14.4/`
`docs/beta/lineage/` snapshots: `3.14.4/`
`docs/beta/schema/` snapshots: `3.14.4/`
`docs/gamma/cdd/` snapshots: 65+ versions (3.13.0 through 3.81.0)
`docs/gamma/essays/` snapshots: `3.14.5/`
`docs/gamma/rules/` snapshots: `3.14.5/`

**Classification rule:** The frozen/historical label applies because the FILE containing the reference is inside a versioned path. The reference syntax (backtick, markdown link, prose) is irrelevant.

**Top contributing files (lines from versioned paths):**

| Source File | Lines |
|---|---|
| `docs/gamma/cdd/3.14.7/CDD.md` | 14 |
| `docs/gamma/cdd/3.14.6/CDD.md` | 10 |
| `docs/gamma/cdd/3.14.5/CDD.md` | 10 |
| `docs/gamma/cdd/3.14.4/CDD.md` | 10 |
| `docs/gamma/cdd/3.14.1/CDD.md` | 10 |
| `docs/gamma/cdd/3.13.0/CDD.md` | 10 |
| `docs/beta/governance/3.14.4/GLOSSARY.md` | 13 |
| `docs/beta/governance/3.14.4/DOCUMENTATION-SYSTEM.md` | 12 |
| `docs/gamma/cdd/3.39.0/POST-RELEASE-ASSESSMENT.md` | 9 |
| `docs/gamma/cdd/3.59.0/POST-RELEASE-ASSESSMENT.md` | 8 |
| *(many more versioned snapshot files)* | 270 |

---

## Class 3 — `test-fixture` (10 lines)

### `.github/workflows/` (4 lines)

```
.github/workflows/build.yml:223:        run: diff docs/alpha/schemas/protocol-contract.json test/cmd/protocol-contract.json
.github/workflows/cnos-agent-admin.yml:82:            4. **Status reporting.** [...] per the entry format in [`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`](../../../../../docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md) §5. [...]
.github/workflows/cnos-agent-admin.yml:168:            - **Channel log convention:** [`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`](../../../../../docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md) [...]
.github/workflows/cnos-cds-dispatch.yml:225:            - **Channel log convention (read-only for dispatch):** [`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`](../../../../../docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md) [...]
```

**Note:** `cnos-agent-admin.yml` and `cnos-cds-dispatch.yml` embed large prompt strings; the triad-path references are inside those embedded prompts. The governing file is a `.github/workflows/` file (classified as test-fixture, not golden-bound — golden-bound applies to the separately-tracked `*.golden.yml` rendered artifacts, not the raw workflow source).

### `test/` (3 lines)

```
test/cmd/cn_contract_test.ml:3:    Verifies that code constants match docs/alpha/schemas/protocol-contract.json.
test/cmd/cn_sandbox_test.ml:59:  show_normalize "docs/alpha/";
test/cmd/cn_traceability_test.ml:5:    contract (docs/alpha/schemas/protocol-contract.json).
```

### `src/` test scripts (3 lines)

```
src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh:93:  mkdir -p "$repo/docs/gamma/cdd/$version"
src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh:94:  cat > "$repo/docs/gamma/cdd/$version/POST-RELEASE-ASSESSMENT.md" <<EOF
src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh:99:  git -C "$repo" add "docs/gamma/cdd/$version/POST-RELEASE-ASSESSMENT.md" >/dev/null
```

---

## Class 4 — `source/package-doctrine` (59 lines)

All `src/` files (65 total) excluding the 3 `*.golden.yml` lines (→ generated/golden-bound) and the 3 `test-cn-cdd-verify.sh` lines (→ test-fixture). These encode triad-docs paths as code-level constants, skill references, or architectural comments.

**Top contributing files:**

| Source File | Lines |
|---|---|
| `src/packages/cnos.cds/skills/cds/CDS.md` | 8 |
| `src/packages/cnos.core/skills/agent/attach/SKILL.md` | 5 |
| `src/packages/cnos.cdd/skills/cdd/` (multiple files) | ~20 |
| `src/packages/cnos.cdd/commands/cdd-verify/README.md` + `run.go` | 3 |
| `src/go/`, `src/ocaml/` code files | 6 |
| *(remaining skill/design docs in src/)* | ~17 |

**Representative lines:**

```
src/go/cmd/cn/main.go:73:	// See docs/alpha/DESIGN-CONSTRAINTS.md §3.1.
src/go/internal/cli/command.go:7:// Design authority: docs/alpha/agent-runtime/GO-KERNEL-COMMANDS.md
src/go/internal/pkg/pkg.go:126:// (installed-package display), per docs/alpha/package-system/PACKAGE-SYSTEM.md
src/packages/cnos.cdd/commands/cdd-verify/run.go:59:	a.Bundle = "docs/gamma/cdd" // default per bash predecessor
src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md:264:Then write the post-release assessment per `post-release/SKILL.md` at `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` (CDD package: `docs/gamma/cdd/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`).
src/packages/cnos.cds/skills/cds/CDS.md:... (8 lines referencing docs/alpha|beta|gamma paths)
src/packages/cnos.core/skills/agent/attach/SKILL.md:... (5 lines referencing docs/ paths)
```

---

## Class 5 — `markdown-link` (31 lines)

Lines in non-versioned `docs/` or `schemas/` files that use Markdown hyperlink syntax `[text](url)` where the URL contains a triad docs path. These break into:

- 22 lines using GitHub absolute URLs: `[text](https://github.com/usurobor/cnos/blob/main/docs/(alpha|beta|gamma)/...)`
- 9 lines using relative URLs: `[text](../docs/(alpha|beta|gamma)/...)` or `[text](../../(alpha|beta|gamma)/...)`

**Representative lines:**

```
docs/alpha/doctrine/coherence-for-agents/COHERENCE-FOR-AGENTS.md:13:- [COHERENCE-SYSTEM](https://github.com/usurobor/cnos/blob/main/docs/alpha/doctrine/COHERENCE-SYSTEM.md) — [...]
docs/alpha/doctrine/ethics-for-agents/ETHICS-FOR-AGENTS.md:20:- [Glossary](https://github.com/usurobor/cnos/blob/main/docs/beta/governance/GLOSSARY.md) — vocabulary authority
docs/beta/architecture/PACKAGE-SYSTEM.md:7:> **[docs/alpha/package-system/PACKAGE-SYSTEM.md](../../alpha/package-system/PACKAGE-SYSTEM.md)**
docs/papers/AGENT-COMMS-FUTURES-KISS.md:362:- `Agent Activation Log Convention v0`: [`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`](../gamma/conventions/AGENT-ACTIVATION-LOG-v0.md)
schemas/README.md:11:[CTB LANGUAGE-SPEC.md](../docs/alpha/ctb/LANGUAGE-SPEC.md) §2 and §11.
schemas/cdd/README.md:126:- [`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md`](../../docs/gamma/essays/CCNF-AND-TYPED-TRUST.md)
```

---

## Class 6 — `inline-path-citation` (208 lines)

All remaining lines in non-versioned `docs/` and `schemas/` files that mention a triad docs path as a prose string, backtick code span, comment, or frontmatter field — not a Markdown hyperlink. This is the largest non-frozen class.

**Top contributing files:**

| Source File | Lines |
|---|---|
| `docs/gamma/cdd/CDD-PACKAGE-AUDIT.md` | 27 |
| `docs/alpha/design/WRITER-PACKAGE.md` | 26 |
| `docs/beta/governance/GLOSSARY.md` | 14 |
| `docs/beta/governance/DOCUMENTATION-SYSTEM.md` | 12 |
| `docs/alpha/package-system/DESIGN-266-dist-out-of-git.md` | 9 |
| `docs/alpha/agent-runtime/CORE-REFACTOR.md` | 7 |
| `docs/gamma/plans/` (multiple files) | ~15 |
| *(remaining non-versioned docs files)* | ~98 |

**Representative lines:**

```
docs/THESIS.md:703:- `docs/alpha/doctrine/FOUNDATIONS.md`
docs/alpha/agent-runtime/CONFIGURE-AGENT.md:6:**Canonical-Path:** docs/alpha/agent-runtime/CONFIGURE-AGENT.md
docs/alpha/agent-runtime/CORE-REFACTOR.md:443:- `docs/alpha/agent-runtime/RUNTIME-CONTRACT-v2.md`
docs/gamma/cdd/CDD-PACKAGE-AUDIT.md:... (27 lines — self-references to triad paths in audit doc)
schemas/README.md:183:- `docs/alpha/ctb/LANGUAGE-SPEC.md` — v0.1 spec governing field shape
```

---

## Classification Verification

| Check | Result |
|---|---|
| Sum of all per-class counts | 3 + 352 + 10 + 59 + 31 + 208 = **663** ✅ |
| Matches AC1 total | **663** ✅ |
| Priority order respected | `generated/golden-bound` > `frozen/historical` > `test-fixture` > `source/package-doctrine` > `markdown-link` > `inline-path-citation` ✅ |
| No line classified by syntax when higher-priority rule applies | ✅ (versioned-path files are `frozen/historical` regardless of syntax; golden.yml files are `generated/golden-bound` regardless of syntax) |
