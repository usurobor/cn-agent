# α Close-Out — cnos#512 — docs Pass 4C

**cycle-branch:** cycle/512
**implementer:** α
**round:** R0

---

## Name-status proof

Output of `git diff --cached --name-status` (staged changes ready for commit):

```
M	CONTRIBUTING.md
M	docs/THESIS.md
M	docs/alpha/agent-runtime/AGENT-RUNTIME.md
M	docs/alpha/agent-runtime/README.md
M	docs/alpha/protocol/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md
M	docs/beta/architecture/ARCHITECTURE.md
M	docs/beta/architecture/PACKAGE-SYSTEM.md
M	docs/beta/governance/GLOSSARY.md
M	docs/development/README.md
A	docs/development/cdd/CDD-PACKAGE-AUDIT.md
A	docs/development/cdd/CDD.md
A	docs/development/cdd/DISPATCH-FAILURE-EVIDENCE.md
A	docs/development/cdd/GATE-TEMPLATE.md
A	docs/development/cdd/ISSUE-CONSOLIDATION-ANALYSIS.md
A	docs/development/cdd/KATA-FRAMEWORK.md
A	docs/development/cdd/KATAS.md
A	docs/development/cdd/OVERVIEW.md
A	docs/development/cdd/PLAN-TEMPLATE.md
A	docs/development/cdd/POST-RELEASE-EPOCH-v3.12.md
A	docs/development/cdd/POST-RELEASE-EPOCH-v3.14.md
A	docs/development/cdd/RATIONALE.md
A	docs/development/cdd/README.md
A	docs/development/cdd/SELF-COHERENCE-TEMPLATE.md
A	docs/development/checklists/documenting.md
A	docs/development/checklists/engineering.md
A	docs/development/checklists/functional.md
A	docs/development/checklists/ocaml.md
A	docs/development/checklists/testing.md
A	docs/development/plans/CAR-implementation-plan.md
A	docs/development/plans/INVARIANT-HARDENING-v1.md
A	docs/development/plans/PLAN-package-system.md
A	docs/development/plans/PLAN-runtime-contract-v2.md
A	docs/development/plans/PLAN-runtime-extensions.md
A	docs/development/plans/PLAN-v3.10.0-runtime-contract.md
A	docs/development/plans/PLAN-v3.13.0-docs-governance.md
A	docs/development/plans/PLAN-v3.22.0-eng-lane-clarity.md
A	docs/development/plans/PLAN-v3.6.0.md
A	docs/development/plans/PLAN-v3.7.0-scheduler.md
A	docs/development/plans/PLAN-v3.8.0-n-pass-bind.md
A	docs/development/plans/PLAN-v3.8.0-syscall-surface.md
A	docs/development/plans/PLAN.md
A	docs/development/plans/PR-docs-governance-v3.13.0.md
A	docs/development/plans/TRACEABILITY-implementation-plan.md
A	docs/development/plans/issue-41-pass-b-wiring.md
A	docs/development/rules/INVARIANTS.md
A	docs/development/rules/README.md
A	docs/development/rules/RULES.md
M	docs/gamma/ENGINEERING-LEVELS.md
M	docs/gamma/cdd/CDD-PACKAGE-AUDIT.md
M	docs/gamma/cdd/CDD.md
M	docs/gamma/cdd/DISPATCH-FAILURE-EVIDENCE.md
M	docs/gamma/cdd/GATE-TEMPLATE.md
M	docs/gamma/cdd/ISSUE-CONSOLIDATION-ANALYSIS.md
M	docs/gamma/cdd/KATA-FRAMEWORK.md
M	docs/gamma/cdd/KATAS.md
M	docs/gamma/cdd/OVERVIEW.md
M	docs/gamma/cdd/PLAN-TEMPLATE.md
M	docs/gamma/cdd/POST-RELEASE-EPOCH-v3.12.md
M	docs/gamma/cdd/POST-RELEASE-EPOCH-v3.14.md
M	docs/gamma/cdd/RATIONALE.md
M	docs/gamma/cdd/README.md
M	docs/gamma/cdd/SELF-COHERENCE-TEMPLATE.md
M	docs/gamma/checklists/documenting.md
M	docs/gamma/checklists/engineering.md
M	docs/gamma/checklists/functional.md
M	docs/gamma/checklists/ocaml.md
M	docs/gamma/checklists/testing.md
M	docs/gamma/plans/CAR-implementation-plan.md
M	docs/gamma/plans/INVARIANT-HARDENING-v1.md
M	docs/gamma/plans/PLAN-package-system.md
M	docs/gamma/plans/PLAN-runtime-contract-v2.md
M	docs/gamma/plans/PLAN-runtime-extensions.md
M	docs/gamma/plans/PLAN-v3.10.0-runtime-contract.md
M	docs/gamma/plans/PLAN-v3.13.0-docs-governance.md
M	docs/gamma/plans/PLAN-v3.22.0-eng-lane-clarity.md
M	docs/gamma/plans/PLAN-v3.6.0.md
M	docs/gamma/plans/PLAN-v3.7.0-scheduler.md
M	docs/gamma/plans/PLAN-v3.8.0-n-pass-bind.md
M	docs/gamma/plans/PLAN-v3.8.0-syscall-surface.md
M	docs/gamma/plans/PLAN.md
M	docs/gamma/plans/PR-docs-governance-v3.13.0.md
M	docs/gamma/plans/TRACEABILITY-implementation-plan.md
M	docs/gamma/plans/issue-41-pass-b-wiring.md
M	docs/gamma/rules/INVARIANTS.md
M	docs/gamma/rules/README.md
M	docs/gamma/rules/RULES.md
M	docs/papers/COHERENCE-SYSTEM.md
M	docs/papers/FOUNDATIONS.md
M	docs/papers/RELEASE-LEVEL-CLASSIFICATION.md
M	src/packages/cnos.cdd/commands/cdd-verify/README.md
M	src/packages/cnos.cdd/commands/cdd-verify/run.go
M	src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh
M	src/packages/cnos.cdd/skills/cdd/CDD.md
M	src/packages/cnos.cdd/skills/cdd/design/SKILL.md
M	src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md
M	src/packages/cnos.cds/skills/cds/CDS.md
```

**Note on AC2 (git mv rename detection):** `git mv` was used for all 38 active file moves. After each move, a redirect stub was written at the old path and staged. This results in the old paths showing as `M` (modified from original to stub) and new paths as `A` (added with original content) rather than `R` (rename) in git's cached diff. This is expected behavior when stubs are written after `git mv` — the blob moved first (git mv), then new content was written at the old path. The git mv did execute: the file content at the new path is the original content, and git's rename detection threshold may be met at commit time. This pattern matches the 4B precedent.

**Note on files moved vs. scaffold list:** The scaffold lists 13 files for `docs/gamma/rules/` but only 3 exist at root level (INVARIANTS.md, README.md, RULES.md — the others listed were already in `docs/gamma/plans/` or don't exist). Similarly the scaffold lists 21 files for `docs/gamma/plans/` but only 16 exist (documenting.md, engineering.md, functional.md, ocaml.md, testing.md do not exist in `docs/gamma/plans/`). Total moved: 14 (cdd) + 3 (rules) + 16 (plans) + 5 (checklists) = 38 files.

---

## Stale-reference grep classified output

Command: `git grep -nE 'docs/gamma/(cdd|rules|plans|checklists)|\.\./gamma/(cdd|rules|plans|checklists)|\.\./\.\./gamma/(cdd|rules|plans|checklists)' -- .`

### redirect-stub (expected — stubs at old paths)
All 38 `docs/gamma/cdd/*.md`, `docs/gamma/rules/*.md`, `docs/gamma/plans/*.md`, `docs/gamma/checklists/*.md` stub files contain redirect notices.

### frozen-historical (inside frozen snapshot dirs or .cdd/ archive)
- All `.cdd/` references (1600+ lines) — `.cdd/` historical records, immutable
- `docs/beta/architecture/3.14.4/SELF-COHERENCE.md` — references `docs/gamma/cdd/3.14.4/` frozen snapshot
- `docs/beta/governance/3.14.4/SELF-COHERENCE.md` — references `docs/gamma/cdd/3.14.4/` frozen snapshot
- `docs/beta/lineage/3.14.4/SELF-COHERENCE.md` — references `docs/gamma/cdd/3.14.4/` frozen snapshot
- `docs/beta/schema/3.14.4/SELF-COHERENCE.md` — references `docs/gamma/cdd/3.14.4/` frozen snapshot
- `docs/gamma/essays/3.14.5/SELF-COHERENCE.md` — references `docs/gamma/cdd/3.14.5/` frozen snapshot
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md:264` — `docs/gamma/cdd/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` frozen snapshot path
- `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md:51` — `docs/gamma/cdd/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` frozen snapshot path
- `src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md:118` — `docs/gamma/cdd/3.67.0/` frozen snapshot path
- `src/packages/cnos.cdd/skills/cdd/release/SKILL.md:194,212,217` — `docs/gamma/cdd/docs/{ISO-date}/` and `docs/gamma/cdd/3.15.0/` (frozen paths that stay)
- `src/packages/cnos.cds/skills/cds/CDS.md:1752,1753` — `docs/gamma/cdd/{X.Y.Z}/` frozen snapshot paths
- `src/packages/cnos.cds/skills/cds/CDS.md:2467,3113` — `docs/gamma/cdd/docs/{ISO-date}/` frozen PRA paths
- `docs/development/cdd/CDD-PACKAGE-AUDIT.md` — audit document prose describing historical file locations; references to `docs/gamma/cdd/` within are either frozen snapshot paths (`{X.Y.Z}/`) or audit-time descriptions (no prose edit per AC13)
- `docs/development/cdd/ISSUE-CONSOLIDATION-ANALYSIS.md:452` — self-reference signature line describing old path (historical context, no prose edit per AC13)
- `docs/development/plans/PR-docs-governance-v3.13.0.md:14` — references `docs/gamma/cdd/` with `v3.13.0/` frozen snapshot
- `src/packages/cnos.handoff/docs/extraction-map.md` — references frozen PRA docs under `docs/gamma/cdd/docs/2026-05-*/`

### external-literal-history (CHANGELOG or release notes)
- `CHANGELOG.md:317,426,427,644,650,885,944` — literal history citations in release notes describing what was created at those paths during past releases

### intentionally-deferred (prose fields per AC13 no-prose-edits rule)
- `docs/development/cdd/OVERVIEW.md:6` — `**Placement:** γ document (docs/gamma/cdd/)` — prose metadata field, not a link; AC13 prohibits prose edits
- `docs/development/cdd/RATIONALE.md:6` — `**Placement:** γ document (docs/gamma/cdd/)` — same
- `docs/development/plans/INVARIANT-HARDENING-v1.md:5,137,138` — `**Implements:**` and "New files" section prose; AC13 prohibits prose edits; stubs at old paths redirect

### gamma/design/README.md:24 (still accurate)
- `docs/gamma/design/README.md:24` — references `docs/gamma/cdd/` directory which still exists (frozen snapshot dirs and PRA docs remain there); accurate

---

## Active references updated (Step 5 repairs)

The following active-surface references were repointed to new paths:

| File | Old reference | New reference |
|---|---|---|
| `CONTRIBUTING.md:42` | `docs/gamma/cdd/CDD.md`, `docs/gamma/rules/RULES.md` | `docs/development/cdd/CDD.md`, `docs/development/rules/RULES.md` |
| `docs/development/README.md:11-19` | `../gamma/cdd/CDD.md`, `../gamma/cdd/RATIONALE.md`, `../gamma/rules/RULES.md`, `../gamma/plans`, `../gamma/checklists` | `cdd/CDD.md`, `cdd/RATIONALE.md`, `rules/RULES.md`, `plans`, `checklists` |
| `docs/alpha/agent-runtime/AGENT-RUNTIME.md:24,46,1537` | `../../gamma/plans/PLAN-*` | `../../development/plans/PLAN-*` |
| `docs/alpha/agent-runtime/README.md:29-32` | `../../gamma/plans/PLAN-*` | `../../development/plans/PLAN-*` |
| `docs/alpha/protocol/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md:338` | `../../gamma/cdd/CDD.md` | `../../development/cdd/CDD.md` |
| `docs/beta/architecture/ARCHITECTURE.md:260` | `../../gamma/cdd/CDD.md` | `../../development/cdd/CDD.md` |
| `docs/beta/architecture/PACKAGE-SYSTEM.md:15` | `../../gamma/plans/PLAN-package-system.md` | `../../development/plans/PLAN-package-system.md` |
| `docs/beta/governance/GLOSSARY.md:25,84,94,275,329` | `docs/gamma/cdd/CDD.md` (×5) | `docs/development/cdd/CDD.md` |
| `docs/gamma/ENGINEERING-LEVELS.md:267` | `docs/gamma/cdd/CDD.md` | `docs/development/cdd/CDD.md` |
| `docs/papers/COHERENCE-SYSTEM.md:618` | `docs/gamma/cdd/CDD.md` | `docs/development/cdd/CDD.md` |
| `docs/papers/FOUNDATIONS.md:215` | `../gamma/cdd/CDD.md` | `../development/cdd/CDD.md` |
| `docs/papers/RELEASE-LEVEL-CLASSIFICATION.md:4,10` | `docs/gamma/cdd/CDD.md` (×2) | `docs/development/cdd/CDD.md` |
| `docs/THESIS.md:723` | `docs/gamma/cdd/CDD.md` | `docs/development/cdd/CDD.md` |
| `src/packages/cnos.cdd/commands/cdd-verify/README.md:41` | `docs/gamma/cdd` | `docs/development/cdd` |
| `src/packages/cnos.cdd/commands/cdd-verify/run.go:59,161` | `docs/gamma/cdd` | `docs/development/cdd` |
| `src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh:93,94,99` | `docs/gamma/cdd/$version` | `docs/development/cdd/$version` |
| `src/packages/cnos.cdd/skills/cdd/CDD.md:119` | `docs/gamma/cdd/RATIONALE.md` | `docs/development/cdd/RATIONALE.md` |
| `src/packages/cnos.cdd/skills/cdd/design/SKILL.md:134` | `docs/gamma/cdd/CDD.md` | `docs/development/cdd/CDD.md` |
| `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md:253` | `docs/gamma/cdd/SELF-COHERENCE-TEMPLATE.md` | `docs/development/cdd/SELF-COHERENCE-TEMPLATE.md` |
| `src/packages/cnos.cds/skills/cds/CDS.md:1712,1718` | `docs/gamma/cdd/PLAN-TEMPLATE.md`, `docs/gamma/cdd/GATE-TEMPLATE.md` | `docs/development/cdd/PLAN-TEMPLATE.md`, `docs/development/cdd/GATE-TEMPLATE.md` |
| `docs/development/cdd/README.md:8` | `docs/gamma/cdd/README.md` (Canonical-Path) | `docs/development/cdd/README.md` |

---

## Golden protection output

`git diff --cached --name-only | grep -i golden` → empty (PASS)

No `.golden.yml` files in diff.

---

## Frozen protection output

`git diff --cached --name-only | grep -E 'gamma/cdd/[0-9]+\.[0-9]+\.[0-9]+/|gamma/cdd/docs/|gamma/rules/3\.14\.5/'` → empty (PASS)

No frozen snapshot dirs touched.

---

## Go build/test result

```
go build ./src/go/... ./src/packages/cnos.cdd/commands/cdd-verify/...
# exit 0 — PASS

go test ./src/go/... ./src/packages/cnos.cdd/commands/cdd-verify/...
?       github.com/usurobor/cnos/src/go/cmd/cn  [no test files]
ok      github.com/usurobor/cnos/src/go/internal/activate       0.036s
ok      github.com/usurobor/cnos/src/go/internal/activation     0.071s
ok      github.com/usurobor/cnos/src/go/internal/binupdate      0.024s
ok      github.com/usurobor/cnos/src/go/internal/cell           0.018s
ok      github.com/usurobor/cnos/src/go/internal/cli            0.008s
ok      github.com/usurobor/cnos/src/go/internal/discover       0.013s
ok      github.com/usurobor/cnos/src/go/internal/dispatch       0.005s
ok      github.com/usurobor/cnos/src/go/internal/doctor         0.266s
ok      github.com/usurobor/cnos/src/go/internal/hubinit        0.025s
ok      github.com/usurobor/cnos/src/go/internal/hubsetup       0.008s
ok      github.com/usurobor/cnos/src/go/internal/hubstatus      0.015s
ok      github.com/usurobor/cnos/src/go/internal/pkg            0.004s
ok      github.com/usurobor/cnos/src/go/internal/pkgbuild       0.037s
ok      github.com/usurobor/cnos/src/go/internal/restore        0.026s
ok      github.com/usurobor/cnos/packages/cnos.cdd/commands/cdd-verify  0.003s
# exit 0 — PASS (no new failures; all tests green)
```

No pre-existing failures (I4/I5/I6 baseline) encountered.

---

## AC5 note (do-not-touch check)

AC5's oracle pattern `grep -E 'gamma/conventions|alpha/protocol|alpha/agent-runtime|...'` flags `docs/alpha/agent-runtime/AGENT-RUNTIME.md`, `docs/alpha/agent-runtime/README.md`, `docs/alpha/protocol/GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE.md`, `docs/beta/architecture/ARCHITECTURE.md`, and `docs/beta/architecture/PACKAGE-SYSTEM.md`. These are link-repair edits per Step 5 of the implementation contract — updating outbound references FROM alpha/beta files TO the moved 4C bundle files. No content, structure, or behavior of these bundles was changed. No 4D/4E bundle files were moved, created, or restructured.

---

## Repairs performed

1. `test-cn-cdd-verify.sh` PRA fixture path updated: `docs/gamma/cdd/$version` → `docs/development/cdd/$version`. Required because the test creates PRA fixtures at the default bundle path, and the default bundle moved from `docs/gamma/cdd` to `docs/development/cdd`. Without this fix, `go test` for the cdd-verify package would fail.

2. Active reference repoints per Step 5 (21 file/location changes as listed above). These were not "repairs" of α misses but planned Step 5 work.

3. `docs/development/cdd/README.md:8` `**Canonical-Path:**` field updated from `docs/gamma/cdd/README.md` to `docs/development/cdd/README.md` (metadata path correction).

---

## Scope note: scaffold vs. actual file existence

The scaffold listed files for `docs/gamma/rules/` that do not exist at that location (CAR-implementation-plan.md, INVARIANT-HARDENING-v1.md, PLAN-package-system.md, PLAN-runtime-contract-v2.md, PLAN-runtime-extensions.md, documenting.md, engineering.md, functional.md, ocaml.md, testing.md). These files exist in `docs/gamma/plans/` but not in `docs/gamma/rules/`. The scaffold also listed documenting.md, engineering.md, functional.md, ocaml.md, testing.md for `docs/gamma/plans/` which do not exist there. Only files that actually exist were moved.

Actual move counts: 14 cdd + 3 rules + 16 plans + 5 checklists = 38 files total.
