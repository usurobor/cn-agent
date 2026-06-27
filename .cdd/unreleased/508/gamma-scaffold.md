---
cycle: 508
branch: cycle/508
main-sha: 62f93b2f97358aa1954dd10cf8ff262a0b897665
issue: "#508"
mode: explore
authored-by: gamma
---

# γ Scaffold — Cycle 508

## 1. Header

- **Issue:** #508 — docs Pass 4A: audit alpha/beta/gamma bundle path dependencies before reader-intent migration
- **Cycle branch:** `cycle/508`
- **Base SHA (main):** `62f93b2f97358aa1954dd10cf8ff262a0b897665`
- **Mode:** `explore` — audit/investigation only; no file moves, no code, no version bump
- **Work shape:** substantial (7 ACs, wide scan surface, golden-bound path analysis)

## 2. Mode

`explore` — The cell's only output is analysis artifacts. No file moves, no stub creation or retirement, no golden re-renders, no `src/` citation repoints, no prose edits, no frontmatter changes. Physical moves (4B–4E) are separate cells filed only after operator review of this audit.

## 3. Surfaces γ expects α to touch

α may only add files. No existing file may be edited, renamed, or deleted. Expected additions:

| Path | Purpose |
|---|---|
| `.cdd/unreleased/508/self-coherence.md` | α's gap declaration, mode, ACs covered, CDD Trace, review-readiness signal |
| `.cdd/unreleased/508/pass4-path-inventory.txt` | AC1 oracle — raw `git grep` output of all `docs/(alpha|beta|gamma)/` refs |
| `.cdd/unreleased/508/pass4-triad-token-inventory.txt` | Supplemental `alpha/|beta/|gamma/` token scan |
| `.cdd/unreleased/508/pass4-golden-inventory.txt` | `find` output listing all `*golden*` / `*.golden.*` files |
| `.cdd/unreleased/508/pass4-classification.md` | AC2 oracle — every inventory line labeled with exactly one taxonomy class |
| `.cdd/unreleased/508/pass4-move-map.md` | AC3 oracle — bundle-by-bundle old → reader-intent path map |
| `.cdd/unreleased/508/pass4-golden-impact-map.md` | AC4 oracle — per-golden-file path impact table |
| `.cdd/unreleased/508/pass4-do-not-move.md` | AC5 oracle — explicit frozen/historical/golden-bound exclusion list |
| `.cdd/unreleased/508/pass4-subcell-order.md` | AC6 oracle — 4B→4C→4D→4E sequencing rationale |
| *(optional)* `docs/` planning location | Mirror of migration map for operator review convenience |

α must NOT create or modify any file outside `.cdd/unreleased/508/` and (optionally) a `docs/` planning location — and even there, only additions are permitted.

## 4. AC oracle approach

### AC1 — Full inventory exists

**α produces:** `pass4-path-inventory.txt` via `git grep -nE 'docs/(alpha|beta|gamma)/' -- src docs .github scripts test tests schemas cue.mod tools`.

**α records:** total line count at production time.

**β verifies independently:** Re-runs the same command on `cycle/508` HEAD, counts lines, confirms count matches α's recorded figure. Any discrepancy is a binding finding.

### AC2 — Every reference classified

**α produces:** `pass4-classification.md` labeling each inventory line with exactly one of: `markdown-link` · `inline-path-citation` · `source/package-doctrine` · `generated/golden-bound` · `frozen/historical` · `test-fixture`.

**α verifies:** Sum of per-class counts equals total inventory line count.

**β verifies independently:** Spot-checks a random sample from each class (minimum 5 per class when class has ≥5 entries); verifies the taxonomy label is defensible for each sampled line; confirms total count matches AC1.

### AC3 — Bundle-by-bundle move map exists

**α produces:** `pass4-move-map.md` mapping each bundle under `docs/alpha/`, `docs/beta/`, `docs/gamma/` to its reader-intent destination (or explicit "stays/defer" note) with per-bundle reference counts drawn from AC1.

**β verifies independently:** Confirms every bundle directory under the three triad roots appears in the map; verifies reference counts are consistent with the inventory; flags any bundle that has references but no map entry.

### AC4 — CI/golden impact map exists

**α produces:** `pass4-golden-impact-map.md` starting from the 2 known golden-bound wake files:
- `cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`
- `cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`

Each golden file is listed with: the specific triad docs paths it cites, the bundle(s) those paths belong to, and the re-render/guard implication if that bundle moves.

**β verifies independently:** Opens each of the 2 named golden files, confirms each triad path they cite appears in the impact map; checks that no golden file found in `pass4-golden-inventory.txt` is omitted if it references a triad path.

### AC5 — Frozen records and snapshots explicitly excluded

**α produces:** `pass4-do-not-move.md` naming `.cdd/` records, version-stamped snapshot paths (`**/3.14.5/` pattern), historical PRA/audit docs, and any path the AC4 impact map flags as golden-bound. Entries named, not guessed.

**β verifies independently:** Confirms `.cdd/` records are listed; confirms at least one version-stamped snapshot path appears (or states none found, with evidence); confirms all AC4 golden-bound paths are also listed in the do-not-move file.

### AC6 — Recommended sub-cell order written before any move

**α produces:** `pass4-subcell-order.md` containing the 4B→4C→4D→4E ordering with rationale, per-cell scope, and explicit reference to the golden-impact map for sequencing decisions.

**β verifies independently:** Confirms the ordering doc cites the AC4 impact map; confirms riskiest/golden-bound bundles are sequenced after the audit (i.e., not in 4B unless confirmed low-risk); flags any ordering that contradicts the impact map.

### AC7 — No file moves in 4A

**α invariant:** zero renames, zero deletions of `docs/alpha|beta|gamma` files.

**β verifies independently:** Runs `git diff --name-status origin/main...HEAD` and confirms: (1) every line starts with `A` (addition) or `M` (modification to a newly-added file); (2) zero lines start with `R` (rename); (3) zero deletions of any `docs/alpha/`, `docs/beta/`, or `docs/gamma/` path appear.

## 5. Expected diff scope

`git diff --name-status origin/main...HEAD` for the completed cycle/508 must show:

- **Only `A` lines** (additions) — all under `.cdd/unreleased/508/` and optionally a `docs/` planning location
- **Zero `R` lines** (renames)
- **Zero deletions** of any `docs/alpha|beta|gamma` path
- The γ-scaffold itself: `.cdd/unreleased/508/gamma-scaffold.md` (already present at cycle open)

Any other status is a protocol violation. β enforcing AC7 is the binding gate.

## 6. Scope guardrails — what α must NOT do

- **No file moves.** Not even "obviously safe" ones.
- **No stub creation or retirement.**
- **No golden re-renders.**
- **No prose edits** to any existing doc.
- **No `src/` citation repoints.**
- **No frontmatter changes** to any existing file.
- **No version bumps.**
- **Do not trust raw grep alone.** Some references are historical/frozen; classify before mapping.
- **Do not invent triad metadata.** Triad (alpha/beta/gamma) is coherence metadata in the framework, not folder authorship claims.
- **Do not open or edit `.cdd/unreleased/508/gamma-scaffold.md`** (this file).

## 7. α Dispatch Prompt

---

**Role:** α  
**Branch:** `cycle/508`  
**Issue:** #508 — docs Pass 4A: audit alpha/beta/gamma bundle path dependencies before reader-intent migration  
**Protocol:** CDS  

### Context

Passes 1–3 of the docs reorg have shipped: the reader-intent portal is live (`docs/README.md`), essays are consolidated to `docs/papers/`, and the two-wakes blocker (#467, #486, #487) is cleared. The `docs/alpha/`, `docs/beta/`, `docs/gamma/` bundles still physically hold content; intent indexes link back into the triad. Pass 4 will physically move those bundles — but only after this audit (4A) maps every path dependency and CI golden impact.

Your job is **audit only**. Produce the five analysis artifacts. Move nothing.

### Skills to load (Tier 3)

- `cnos.core/skills/design` — boundary decisions, source-of-truth clarity
- `eng/tool` — shell discipline for the inventory commands

### Implementation guidance — starting commands

Run these first and capture output verbatim:

```bash
# AC1 primary inventory
git grep -nE 'docs/(alpha|beta|gamma)/' -- \
  src docs .github scripts test tests schemas cue.mod tools > .cdd/unreleased/508/pass4-path-inventory.txt

# Supplemental triad-token scan
git grep -nE 'alpha/|beta/|gamma/' -- \
  src docs .github scripts test tests schemas cue.mod tools > .cdd/unreleased/508/pass4-triad-token-inventory.txt

# Golden-file inventory
find . -name '*golden*' -o -name '*.golden.*' > .cdd/unreleased/508/pass4-golden-inventory.txt
```

After capturing, record the line count of `pass4-path-inventory.txt` (the AC1 total). This number must appear in `self-coherence.md` and `pass4-classification.md`.

### Implementation contract

| Axis | Value |
|---|---|
| Language | bash (shell grep/find for inventory) + markdown (artifact files) |
| CLI integration target | none — explore cell; no CLI changes |
| Package scoping | artifacts under `.cdd/unreleased/508/` only; optionally a `docs/` planning location for the migration map |
| Existing-binary disposition | no binary changes |
| Runtime dependencies | standard `git grep`, `find`, `bash` — no external tools |
| JSON/wire contract preservation | N/A — no wire contracts touched |
| Backward-compat invariant | no files moved, renamed, or deleted; zero `R` lines in `git diff --name-status origin/main...HEAD` (AC7 hard gate) |

### AC1–AC7 oracle list

| AC | Artifact | How α satisfies it | How β verifies it |
|---|---|---|---|
| AC1 | `pass4-path-inventory.txt` | Run the starting command verbatim; record line count | Re-run same command; confirm line count matches |
| AC2 | `pass4-classification.md` | Label every inventory line with exactly one taxonomy class; sum per-class counts = total | Spot-check ≥5 per class; confirm total sum = AC1 count |
| AC3 | `pass4-move-map.md` | Map every bundle under `docs/alpha|beta|gamma` to reader-intent destination or "stays/defer"; include per-bundle ref counts | Confirm all bundles appear; ref counts consistent with inventory |
| AC4 | `pass4-golden-impact-map.md` | Start from 2 known golden files; list every triad path each cites; state re-render implication | Open both golden files; verify all cited paths appear in map |
| AC5 | `pass4-do-not-move.md` | Name `.cdd/` records, version-stamped snapshots, historical PRA/audit docs, all AC4 golden-bound paths | Confirm `.cdd/` listed; confirm AC4 golden-bound paths appear |
| AC6 | `pass4-subcell-order.md` | Write 4B→4C→4D→4E order with rationale; cite impact map for sequencing | Confirm ordering references AC4 map; confirm riskiest bundles not in 4B unless impact-map-confirmed low-risk |
| AC7 | `git diff --name-status` | Make only additions; zero renames | Run `git diff --name-status origin/main...HEAD`; zero `R` lines, zero `docs/alpha|beta|gamma` deletions |

### Classification taxonomy (AC2)

Every line in the inventory must be labeled exactly one of:

- `markdown-link` — a `[text](docs/alpha|beta|gamma/...)` or `](docs/alpha|beta|gamma/...)` reference in a Markdown file
- `inline-path-citation` — a path string mentioned in prose or a comment (not a hyperlink, not a tool call)
- `source/package-doctrine` — a reference in a `src/` file, schema, or package config that encodes a path as a code-level dependency or convention
- `generated/golden-bound` — a reference inside a `*.golden.yml` or other generated/byte-checked artifact
- `frozen/historical` — a reference in a `.cdd/` record, version-stamped snapshot, historical PRA/audit doc, or changelog entry that must not be rewritten
- `test-fixture` — a path reference in a test or fixture file (`test/`, `tests/`, `.github/` workflow that checks a specific path)

A line may appear to match multiple classes; pick the most specific (e.g., a golden file reference is `generated/golden-bound`, not `inline-path-citation`).

### Known golden-bound files (start here for AC4)

Both of the following files are known to reference triad docs paths; verify and record exact paths:

1. `cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`
2. `cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`

Also scan `pass4-golden-inventory.txt` for any additional `*.golden.yml` files that reference `docs/(alpha|beta|gamma)/`.

### Recommended 4B–4E physical-move order for AC6

Do NOT execute these moves. Document the recommended order only:

- **4B (low-risk reader surfaces):** `beta/guides/` → `guides/`, `beta/evidence/` → `evidence/`, `beta/lineage/` → `concepts/lineage/`, `alpha/doctrine/` → `concepts/doctrine/`
- **4C (development process):** `gamma/cdd/` → `development/cdd/`, etc.
- **4D (reference bundles):** `alpha/protocol/`, `alpha/agent-runtime/`, etc.
- **4E (architecture):** `alpha/cognitive-substrate/`, `beta/architecture/`, etc.

Confirm or adjust sequencing based on your impact map. Riskiest/golden-bound bundles must be sequenced after the audit confirms their goldens.

### Non-goals

- No file moves, renames, or deletions (including "obviously safe" ones)
- No stub creation or retirement
- No golden re-renders
- No prose edits to existing docs
- No `src/` citation repoints
- No frontmatter changes
- No version bumps
- No CLI changes

### **Do NOT move or rename any file. AC7 is a hard gate: zero `R` lines in the diff.**

### Self-coherence signal

When all 7 ACs are satisfied and artifacts committed, write `.cdd/unreleased/508/self-coherence.md` with: gap summary, mode (`explore`), AC1 line count, list of artifacts produced, CDD Trace, and your review-readiness signal.

---

## 8. β Dispatch Prompt

---

**Role:** β  
**Branch:** `cycle/508`  
**Issue:** #508 — docs Pass 4A: audit alpha/beta/gamma bundle path dependencies before reader-intent migration  
**Protocol:** CDS  

### Your task

Independently verify that α's artifacts for cycle/508 satisfy AC1–AC7. Do not read α's self-coherence.md reasoning as your source of truth — derive your verdicts by re-running the oracles.

### AC verification checklist

**AC1 — Inventory completeness**

1. Run independently: `git grep -nE 'docs/(alpha|beta|gamma)/' -- src docs .github scripts test tests schemas cue.mod tools | wc -l`
2. Compare your line count to the count recorded in `pass4-path-inventory.txt` (or `self-coherence.md`).
3. If counts diverge by more than 0 (any divergence is a discrepancy): binding finding, severity D (data-integrity).

**AC2 — Classification coverage**

1. Open `pass4-classification.md`.
2. Verify the sum of per-class counts equals the AC1 total.
3. Spot-check: for each of the 6 taxonomy classes, sample at minimum 5 entries (or all entries if fewer than 5 in that class). Verify each sampled entry's taxonomy label is defensible.
4. Classification taxonomy for spot-checking:
   - `markdown-link` — Markdown hyperlink syntax pointing to a triad docs path
   - `inline-path-citation` — prose or comment path string, not a hyperlink
   - `source/package-doctrine` — code-level path dependency in `src/`, schema, or package config
   - `generated/golden-bound` — reference inside a `*.golden.yml` or byte-checked artifact
   - `frozen/historical` — `.cdd/` record, version-stamped snapshot, PRA/audit doc, changelog entry
   - `test-fixture` — path reference in `test/`, `tests/`, or a CI workflow checking a specific path
5. Misclassification of `generated/golden-bound` entries (e.g., classified as `inline-path-citation`) is a C-severity finding.

**AC3 — Move map completeness**

1. List all directories immediately under `docs/alpha/`, `docs/beta/`, `docs/gamma/` in the repo.
2. Confirm every bundle appears in `pass4-move-map.md` with either a reader-intent destination or an explicit "stays/defer" note.
3. Verify per-bundle reference counts in the map are consistent with `pass4-path-inventory.txt`.
4. A bundle with references but no map entry is a D-severity finding.

**AC4 — Golden impact map**

1. Open `cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` — enumerate every `docs/(alpha|beta|gamma)/` path it contains.
2. Open `cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` — enumerate every `docs/(alpha|beta|gamma)/` path it contains.
3. Verify each path appears in `pass4-golden-impact-map.md` with its bundle assignment and re-render implication.
4. From `pass4-golden-inventory.txt`, check any additional `*.golden.yml` or `*.golden.*` file for triad path references; any reference not in the impact map is a D-severity finding.
5. Missing golden-bound path in the impact map is a D-severity finding.

**AC5 — Do-not-move list**

1. Open `pass4-do-not-move.md`.
2. Confirm `.cdd/` records are listed (named, not just "`.cdd/` directory").
3. Confirm version-stamped snapshot paths (e.g., `**/3.14.5/`) appear, or confirm with evidence that none exist in the repo.
4. Confirm all AC4 golden-bound paths also appear in this list.
5. Any AC4 golden-bound path absent from the do-not-move list is a D-severity finding.

**AC6 — Sub-cell order**

1. Open `pass4-subcell-order.md`.
2. Confirm the ordering explicitly cites the AC4 golden-impact map.
3. Confirm riskiest/golden-bound bundles are NOT placed in 4B unless the impact map confirms they carry zero golden-bound references.
4. Ordering that contradicts the impact map (risky bundle early without justification) is a C-severity finding.

**AC7 — Zero renames (hard gate)**

1. Run: `git diff --name-status origin/main...HEAD`
2. Confirm: every line begins with `A` (addition).
3. Confirm: zero lines begin with `R` (rename).
4. Confirm: zero lines show deletion of any path under `docs/alpha/`, `docs/beta/`, or `docs/gamma/`.
5. Any `R` line anywhere in the diff is an RC-blocking finding. Any deletion of a triad docs path is an RC-blocking finding.

### Verdict

Issue your verdict as the standard `#CDSReceipt` typed block with per-AC rows. AC7 violations are automatically RC-blocking regardless of other findings. All D-severity findings are RC-blocking. C-severity findings require explicit disposition before APPROVED.

---

## 9. Friction notes

Known complexities α and β should be aware of:

**Scale.** The issue body cites 684 references to `docs/(alpha|beta|gamma)/` across 221 files. The classification artifact will be large. α should process it systematically (one taxonomy class at a time, or one directory scope at a time) to avoid gaps. β's spot-check must sample deliberately — random sampling across 684 entries with 6 classes means some classes may be undersampled; bias toward the classes with fewer entries.

**Golden-bound paths — highest severity.** Only 2 golden files are currently identified, but `pass4-golden-inventory.txt` may reveal more. Any newly discovered golden-bound triad path that was classified as a lower-severity class (e.g., `markdown-link`) is a misclassification. The golden-bound class takes priority over all others.

**Frozen/historical entries — do not rewrite.** `.cdd/` records and version-stamped snapshots (e.g., entries under paths like `docs/gamma/3.14.5/`) must appear in the do-not-move list. A common error is treating these as ordinary `markdown-link` references because their Markdown syntax is identical; the frozen/historical class is determined by the file's location and purpose, not its syntax.

**Classification taxonomy precision.** The 6 classes are mutually exclusive by priority:
1. `generated/golden-bound` (highest priority — overrides syntax-based classes)
2. `frozen/historical` (overrides syntax-based classes for .cdd/ records and snapshots)
3. `test-fixture` (overrides link-syntax classes when the file is in test/tests/)
4. `source/package-doctrine` (overrides link-syntax classes when the file is in src/)
5. `markdown-link` (applies when none of the above apply and syntax is a link)
6. `inline-path-citation` (lowest priority — fallback for prose/comment references)

**AC7 is absolute.** No "temporary" renames, no "preparatory" moves, no stub creation. The `git diff --name-status` check has no exceptions in a 4A cell.

**Artifact location.** All inventory artifacts must land under `.cdd/unreleased/508/`. Placing them at repo root (e.g., `pass4-path-inventory.txt` at `/`) is an error — β will flag it as a scope violation.
