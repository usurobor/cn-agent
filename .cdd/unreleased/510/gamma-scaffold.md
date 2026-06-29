# gamma-scaffold — cnos#510 — docs Pass 4B

**cell:** cnos#510 — docs Pass 4B: move low-risk reader surfaces (beta/guides, beta/evidence, beta/lineage, alpha/doctrine) → reader-intent homes
**protocol:** cds
**cycle-branch:** cycle/510
**base-sha:** 6b24ab6710a019ee8754b063d8287bb7ee154a84
**mode:** docs-only (git mv + redirect stubs + link repoint; no code, no version bump)
**collapse-authority:** β-α-collapse-on-δ (PERSONA.md commitment 5 — docs-only surface; AC oracle is mechanical)

---

## Source of truth

| Surface | Path | Status |
|---|---|---|
| Move map + sub-cell order | `.cdd/releases/docs/2026-06-21/508/pass4-move-map.md`, `pass4-subcell-order.md` | Shipped via #509 / 6b24ab6 |
| Golden-bound / do-not-move list | `pass4-golden-impact-map.md`, `pass4-do-not-move.md` | Shipped |
| Target intent-index dirs | `docs/guides/README.md`, `docs/evidence/README.md`, `docs/concepts/README.md` | Shipped (Pass 1) |
| Stub pattern precedent | Pass 2/3 (#506 / PR #507, commit 09414da) | Shipped |

---

## Scope guardrails

**In scope — move exactly:**
- `docs/beta/guides/` → `docs/guides/`
- `docs/beta/evidence/` → `docs/evidence/`
- `docs/beta/lineage/` → `docs/concepts/lineage/`
- `docs/alpha/doctrine/` → `docs/concepts/doctrine/`

**Do NOT touch:**
- `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` (golden-bound)
- Any `*.golden.yml` file
- `.cdd/` (except this unreleased/510/ dir)
- `.github/workflows/`
- Versioned snapshots (`**/X.Y.Z/`)
- Historical PRA/audit records
- `docs/gamma/cdd/`, `docs/alpha/protocol/`, `docs/alpha/agent-runtime/`
- `src/` outside of active link repoints

---

## Per-AC oracle

| AC | Oracle | Check |
|---|---|---|
| AC1 | `git diff --name-status origin/main..HEAD` shows `R` (rename) entries for each moved file | grep `^R` diff |
| AC2 | Every moved file's old path has a `# Moved` stub pointing to new path | `git grep "# Moved" -- docs/beta/guides/ docs/beta/evidence/ docs/beta/lineage/ docs/alpha/doctrine/` |
| AC3 | Active markdown links in docs/ and OPERATOR.md point to new homes | Manual link audit in READMEs and OPERATOR.md |
| AC4 | No diff under `.cdd/` (except unreleased/510/), no version snapshot, no historical record | `git diff --name-only` excludes .cdd/ (except 510/) |
| AC5 | `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` unchanged | `git diff -- docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` empty |
| AC6 | No `*.golden.yml` in diff; no `.github/workflows/` in diff | `git diff --name-only \| grep -E '\.golden\.yml\|\.github/workflows'` returns empty |
| AC7 | `git grep -nE 'docs/(beta/guides\|beta/evidence\|beta/lineage\|alpha/doctrine)/'` returns only stubs + frozen citations | grep audit |
| AC8 | `docs/README.md` portal + intent-index READMEs link to new local homes | inspect updated READMEs |
| AC9 | No bidi/hidden/control chars in moved or changed files | `LC_ALL=C grep -rP '[\x00-\x08\x0b\x0c\x0e-\x1f\x7f\x{200b}-\x{200f}\x{202a}-\x{202e}\x{feff}]'` over changed files |
| AC10 | Moved files' content unchanged except required relative-link repairs | `git diff` of moved files shows only link-path changes + stub content at old paths |

---

## Friction notes

- `docs/guides/README.md`, `docs/evidence/README.md`, `docs/concepts/README.md` already exist as Pass 1 intent-index dirs. Merge content under them; keep their READMEs and fold moved files beneath.
- `docs/guides/` and `docs/evidence/` already have README.md files pointing to the old bundle paths — update those links to local.
- `docs/beta/guides/WRITE-A-SKILL.md` and `docs/beta/guides/DOJO.md` use `../../../src/packages/` relative links — after moving to `docs/guides/`, these must become `../../src/packages/` (one fewer `../` level since `docs/guides/` is shallower than `docs/beta/guides/`).
- `docs/concepts/README.md` already exists with links to old doctrine and lineage paths — update these.
- `OPERATOR.md` at repo root has 4 active navigation links to `docs/beta/guides/` — update to `docs/guides/`.
- `docs/alpha/doctrine/judgment-for-agents/JUDGMENT-FOR-AGENTS.md` has `../../../papers/FOUNDATIONS.md` relative link — depth is preserved after move (`docs/alpha/doctrine/subdir/` → `docs/concepts/doctrine/subdir/`, both 3 levels under docs/), no change needed.
- Absolute GitHub URLs within moved files (essays) are frozen content — leave as-is per no-prose-reflow constraint.
- `docs/alpha/design/WRITER-PACKAGE.md` has many backtick-cited paths to `docs/alpha/doctrine/` — these are design document historical citations, classified as frozen/historical per AC7 exception.
- Governance snapshot files (`docs/beta/governance/3.14.4/GLOSSARY.md`, `docs/beta/governance/GLOSSARY.md`) reference `docs/alpha/doctrine/COHERENCE-SYSTEM.md` — these are snapshot artifacts, classified as frozen/historical per AC7.
- CHANGELOG.md reference to `docs/beta/evidence/rca/` — historical record, frozen.

---

## Implementation sequence

1. git mv all files in four source bundles to targets
2. Create `# Moved` stubs at old file paths
3. Update relative link depths in guides files (WRITE-A-SKILL.md, DOJO.md)
4. Update `docs/guides/README.md` — change links from `../beta/guides/X` → `X`
5. Update `docs/evidence/README.md` — change links from `../beta/evidence/X` → `X`
6. Update `docs/concepts/README.md` — update doctrine and lineage links
7. Update `OPERATOR.md` — change 4 links from `docs/beta/guides/` → `docs/guides/`
8. Run AC oracle checks
9. Commit + push
