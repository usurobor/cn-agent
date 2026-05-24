# α self-coherence — cycle/425

**Verdict:** All 11 ACs PASS. Review-ready.

α verifies each AC against its declared oracle. Where the oracle has a mechanical command, α runs the command and records the output. Where the oracle is structural, α reads the artifact and confirms the structure.

## AC1: Essay file exists at canonical path; line count 200–400

**Oracle:** `test -f docs/gamma/essays/BOX-AND-THE-RUNNER.md`; line count 200–400.

**Result:**

```
$ wc -l docs/gamma/essays/BOX-AND-THE-RUNNER.md
288 docs/gamma/essays/BOX-AND-THE-RUNNER.md
```

288 ∈ [200, 400]. **PASS.**

## AC2: Required H2 sections present (≥ 9)

**Oracle:** `grep -c "^## " docs/gamma/essays/BOX-AND-THE-RUNNER.md` returns ≥ 9.

**Result:**

```
$ grep -c "^## " docs/gamma/essays/BOX-AND-THE-RUNNER.md
10
```

The 10 H2 sections, in order:

```
## Remote-Runner Delegation as an Effect Surface   (subtitle)
## Point
## Progress note
## What I realized
## What this means for cnos
## What not to celebrate
## Memory to keep
## The 6-field receipt convention
## First use: cnos#425 v3.82.0 tag retarget
## References
```

10 ≥ 9. **PASS.**

## AC3: Critical content present

**Oracle:** `grep -c "remote-runner delegation\|effect surface\|6-field receipt\|agent boundary" docs/gamma/essays/BOX-AND-THE-RUNNER.md` returns ≥ 5.

**Result:**

```
$ grep -c "remote-runner delegation\|effect surface\|6-field receipt\|agent boundary" docs/gamma/essays/BOX-AND-THE-RUNNER.md
17
```

The four key phrases all appear in their canonical forms; in addition, the seed's load-bearing claims appear verbatim:

- "sandbox that guards only the local shell is not guarding the agent" — present in §"What I realized" line 1.
- "the agent boundary is not the place where the model sits. It is the full path from intention to effect." — present in §"Memory to keep" and restated for emphasis.
- "Any artifact that can make another system execute is an effect surface." — present in §"What this means for cnos" as the canonical rule statement.
- 6-field list (who asked / what runs / where / authority / evidence / who may accept) — present as code block in §"What this means for cnos" and as the §"The 6-field receipt convention" template.

17 ≥ 5. **PASS.**

## AC4: README updated (Document Map + Reading Order)

**Oracle:** `docs/gamma/essays/README.md` Document Map has a new row for `BOX-AND-THE-RUNNER.md`; Reading Order has an entry for it.

**Result:**

Document Map row appended after `CELL-OF-CELLS.md`:

```
| [BOX-AND-THE-RUNNER.md](./BOX-AND-THE-RUNNER.md) | Essay (DRAFT v0.1.0) | The Box and the Runner: remote-runner delegation as an effect surface; the agent boundary is the full path from intention to effect; 6-field receipt convention for artifacts that cause another body to execute |
```

Reading Order item 7 added:

```
7. **[BOX-AND-THE-RUNNER.md](./BOX-AND-THE-RUNNER.md)** — remote-runner delegation as an effect surface; the boundary discipline for artifacts that cause another body to execute
```

```
$ grep -c BOX-AND-THE-RUNNER docs/gamma/essays/README.md
2
```

2 ≥ 2 (Document Map row + Reading Order item). **PASS.**

## AC5: Skill patch present + line count preserved

**Oracle:** `grep -c "remote-runner\|BOX-AND-THE-RUNNER" src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` returns ≥ 2; file line count ≥ 95% of pre-cycle (319 × 0.95 ≈ 303).

**Result:**

```
$ grep -c "remote-runner\|BOX-AND-THE-RUNNER" src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
12
$ wc -l src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
372 src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
```

12 ≥ 2; 372 ≥ 303. The §8 section adds 53 lines without modifying any prior content (verified by `git diff origin/main -- src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` showing pure additions in lines 320–372). **PASS.**

## AC6: Workflow file shape correct

**Oracle:** File exists with named triggers + permissions + 3 substantive steps (tag force-move + push) + self-delete step.

**Result:**

```
$ test -f .github/workflows/repoint-3.82.0.yml && echo present
present
$ grep -E "^(name|on|permissions|jobs):" .github/workflows/repoint-3.82.0.yml
name: One-shot — repoint 3.82.0 tag
on:
permissions:
jobs:
$ grep -c "^      - name:" .github/workflows/repoint-3.82.0.yml
4
```

The workflow has:

- `name:` — "One-shot — repoint 3.82.0 tag"
- `on: push: branches: [main], paths: ['.github/workflows/repoint-3.82.0.yml']` — push trigger path-filtered on its own path
- `permissions: contents: write` — minimal scope for the operations
- `jobs: repoint:` — single job on `ubuntu-latest`
- 4 named steps: checkout / configure git identity / force-move tag (with the load-bearing `git tag -f 3.82.0 fd1d654e` + `git push --force origin 3.82.0`) / self-delete (with `git rm` + `git commit` + `git push origin main`)

The structure matches the YAML shape pinned in cnos#425 verbatim modulo header comments. **PASS.**

## AC7: Receipt file has all 6 fields

**Oracle:** `test -f .cdd/unreleased/425/remote-runner-receipt-3.82.0.md`; contains all 6 named field rows.

**Result:**

```
$ test -f .cdd/unreleased/425/remote-runner-receipt-3.82.0.md && echo present
present
$ grep -c "^## [1-6]\." .cdd/unreleased/425/remote-runner-receipt-3.82.0.md
6
$ grep "^## [1-6]\." .cdd/unreleased/425/remote-runner-receipt-3.82.0.md
## 1. Who asked for it
## 2. What it will run
## 3. Where it will run
## 4. What authority it has
## 5. Evidence
## 6. Who may accept the result
```

All 6 fields present, each substantively populated:

- §1 — operator + 2 dispatch directives cited verbatim.
- §2 — 4 actual job steps with specific git/gh commands listed at command granularity.
- §3 — `ubuntu-latest` GH-hosted ephemeral VM; explicit exclusions (not self-hosted, not third-party).
- §4 — `GITHUB_TOKEN` with `contents: write` scoped to `usurobor/cnos`; lists what scope grants + does NOT grant; names branch-protection as substrate-enforced.
- §5 — post-run-fillable per hard rule 2; 5-row evidence-shape table declares evidence shape in advance.
- §6 — operator with explicit acceptance criterion (release body must match `.cdd/releases/3.82.0/RELEASE.md`); explicitly forbids agent-side acceptance.

Plus expected-effect (5 numbered steps in order), failure-modes (6-row table with mitigations), acceptance-criteria (6 numbered + partial/rejection branches), V/δ composition note. **PASS.**

## AC8: cnos.cdd kernel + other packages untouched

**Oracle:**
- `git diff origin/main..HEAD -- src/packages/cnos.cdd/skills/cdd/CDD.md` returns 0 lines.
- `git diff origin/main..HEAD -- src/packages/cnos.cds/ src/packages/cnos.cdr/ src/packages/cnos.handoff/` returns 0 lines.
- Only `delta/SKILL.md` touched in `cnos.cdd/skills/`.

**Result:**

```
$ git diff origin/main -- src/packages/cnos.cdd/skills/cdd/CDD.md | wc -l
0
$ git diff origin/main -- src/packages/cnos.cds/ src/packages/cnos.cdr/ src/packages/cnos.handoff/ | wc -l
0
$ git diff --name-only origin/main -- src/packages/cnos.cdd/skills/
src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
```

CDD.md byte-identical to origin/main; CDS/CDR/handoff packages byte-identical; only delta/SKILL.md modified in cdd/skills/. **PASS.**

## AC9: No mojibake; Greek letters used natively where applicable

**Oracle:** Greek letters used natively where applicable.

**Result:** The essay's subject is workflows/runners/effect surfaces, not the cell algorithm; it does not require Greek letters in the prose, and α did not introduce Latin transliterations of cell-algorithm symbols (no "alpha"/"beta"/"gamma"/"delta"/"epsilon" written out where Greek would be the canonical form). The skill section (`delta/SKILL.md §8`) uses Greek (δ, α, β, γ) natively as inherited from the surrounding file (which uses Greek throughout — verified by the §8 text mentioning δ-as-role, the α/β/γ producer/reviewer/closer triad, and the V/δ composition).

Mojibake check (defensive):

```
$ grep -cE "Î[±²³´µ]|â¡|matterâ|Î±â" docs/gamma/essays/BOX-AND-THE-RUNNER.md src/packages/cnos.cdd/skills/cdd/delta/SKILL.md .github/workflows/repoint-3.82.0.yml .cdd/unreleased/425/remote-runner-receipt-3.82.0.md
0
0
0
0
```

No mojibake across any of the 4 substantive artifacts. **PASS.**

## AC10: Only BOX-AND-THE-RUNNER + README in essays

**Oracle:** `git diff --name-only origin/main..HEAD -- docs/gamma/essays/` returns exactly BOX-AND-THE-RUNNER.md (added) + README.md (modified).

**Result:**

```
$ git diff --name-only origin/main -- docs/gamma/essays/
docs/gamma/essays/README.md
docs/gamma/essays/BOX-AND-THE-RUNNER.md
```

Exactly 2 files, both in `docs/gamma/essays/`. Verification of byte-identity for the 6 prior essays (5 in `docs/gamma/essays/` + 1 in `docs/essays/`):

```
$ git diff origin/main -- docs/gamma/essays/STATELESS-AGENCY.md docs/gamma/essays/EXECUTABLE-SKILLS.md docs/gamma/essays/COHERENCE-MUST-BE-FREE.md docs/gamma/essays/CCNF-AND-TYPED-TRUST.md docs/gamma/essays/DECREASING-INCOHERENCE.md docs/gamma/essays/CELL-OF-CELLS.md docs/essays/agent-first.md
(0 lines)
$ git diff origin/main -- docs/alpha/essays/
(0 lines)
```

**PASS.**

## AC11: No schemas / runtime / scripts changes

**Oracle:** `git diff origin/main..HEAD -- schemas/ src/go/ src/packages/cnos.core/ src/packages/cnos.eng/ src/packages/cnos.kata/ src/packages/cnos.cdd.kata/ scripts/` returns 0 lines.

**Result:**

```
$ git diff origin/main -- schemas/ src/go/ src/packages/cnos.core/ src/packages/cnos.eng/ src/packages/cnos.kata/ src/packages/cnos.cdd.kata/ scripts/ | wc -l
0
```

0 lines across all 7 excluded paths. **PASS.**

## Summary

| AC | Verdict |
|----|---------|
| AC1: file exists, line count 200–400 | PASS (288) |
| AC2: ≥ 9 H2 sections | PASS (10) |
| AC3: ≥ 5 critical-content hits | PASS (17) |
| AC4: README Document Map + Reading Order | PASS (2 entries) |
| AC5: skill patch ≥ 2 grep hits + ≥ 95% line count | PASS (12; 372 ≥ 303) |
| AC6: workflow shape correct | PASS (named triggers + permissions + 3 substantive steps + self-delete) |
| AC7: receipt has 6 fields | PASS (all 6 populated; §5 post-run-fillable per hard rule) |
| AC8: kernel + cds/cdr/handoff untouched; only delta/SKILL.md in cdd/skills | PASS (0/0; only delta/SKILL.md) |
| AC9: no mojibake; Greek natively where applicable | PASS (0 mojibake; Greek inherited from delta/SKILL.md) |
| AC10: only BOX-AND-THE-RUNNER + README in essays | PASS (exactly 2 files; other essays byte-identical) |
| AC11: no schemas/runtime/scripts | PASS (0 lines across 7 excluded paths) |

All 11 ACs PASS. Bundle is review-ready for β.
