# α self-coherence — cycle/426

**Verdict:** All 10 ACs PASS. Review-ready.

α verifies each AC against its declared oracle. Where the oracle has a mechanical command, α runs the command and records the output. Where the oracle is structural, α reads the artifact and confirms the structure.

## AC1: Workflow file exists with correct shape

**Oracle:** `test -f .github/workflows/publish-3.82.0-release.yml`. Contains `on: push: paths: [...self-path...]`, `permissions: contents: write`, 4 named steps (checkout at tag, verify RELEASE.md, softprops/action-gh-release@v1 invocation, self-delete).

**Result:**

```
$ test -f .github/workflows/publish-3.82.0-release.yml && echo present
present
$ grep -E "^(name|on|permissions|jobs):" .github/workflows/publish-3.82.0-release.yml
name: One-shot — publish 3.82.0 GH release
on:
permissions:
jobs:
$ grep -E "^      - (name|uses):" .github/workflows/publish-3.82.0-release.yml
      - uses: actions/checkout@v4
      - name: Verify root RELEASE.md is v3.82.0
      - name: Create GitHub Release
      - name: Self-delete the workflow file
```

The workflow has all four top-level keys present, push trigger path-filtered on its own path, `permissions: contents: write`, and 4 steps (1 `uses:` for the checkout + 3 `name:` for verify/release/self-delete). **PASS.**

## AC2: Workflow checks out at the tag (not main)

**Oracle:** `grep "ref:.*3.82.0" .github/workflows/publish-3.82.0-release.yml` returns ≥ 1 hit.

**Result:**

```
$ grep -c "ref:.*3.82.0" .github/workflows/publish-3.82.0-release.yml
1
$ grep "ref:" .github/workflows/publish-3.82.0-release.yml
          ref: '3.82.0'
```

The checkout pins `ref: '3.82.0'` (the literal tag string, not `main` or `HEAD`). `git ls-remote origin refs/tags/3.82.0` resolves to `fd1d654e8d6361f0db2a6407f19e573a96d1054d` — the canonical v3.82.0 baseline commit, where root `RELEASE.md` is the correct body. **PASS.**

## AC3: Release body uses root RELEASE.md

**Oracle:** `grep "body_path: RELEASE.md" .github/workflows/publish-3.82.0-release.yml` returns 1 hit.

**Result:**

```
$ grep -c "body_path: RELEASE.md" .github/workflows/publish-3.82.0-release.yml
1
$ grep "body_path:" .github/workflows/publish-3.82.0-release.yml
          body_path: RELEASE.md
```

The `softprops/action-gh-release@v1` invocation sources its body from root `RELEASE.md` of the checked-out tree (which is the tag, per AC2). **PASS.**

## AC4: Receipt file has all 6 fields

**Oracle:** `test -f .cdd/unreleased/426/remote-runner-receipt-publish-3.82.0.md`; contains 6 numbered or named field rows; evidence field is post-run-fillable placeholder.

**Result:**

```
$ test -f .cdd/unreleased/426/remote-runner-receipt-publish-3.82.0.md && echo present
present
$ grep -c "^## [1-6]\." .cdd/unreleased/426/remote-runner-receipt-publish-3.82.0.md
6
$ grep "^## [1-6]\." .cdd/unreleased/426/remote-runner-receipt-publish-3.82.0.md
## 1. Who asked for it
## 2. What it will run
## 3. Where it will run
## 4. What authority it has
## 5. Evidence
## 6. Who may accept the result
```

All 6 fields present, each substantively populated:

- §1 — operator + cnos#425 follow-on directive (release-boundary, post-build-failure-diagnosis) cited.
- §2 — 4 actual job steps with specific git/softprops commands listed at command granularity.
- §3 — `ubuntu-latest` GH-hosted ephemeral VM; explicit exclusions (not self-hosted, not third-party); same surface as cnos#425.
- §4 — `GITHUB_TOKEN` with `contents: write` scoped to `usurobor/cnos`; lists what scope grants + does NOT grant; names branch-protection as substrate-enforced.
- §5 — post-run-fillable per hard rule 1; 6-row evidence-shape table declares evidence shape in advance.
- §6 — operator with explicit acceptance criterion (release body must match `.cdd/releases/3.82.0/RELEASE.md`); explicitly forbids agent-side acceptance.

Plus expected-effect (5 numbered steps), failure-modes (7-row table with mitigations including release-already-exists idempotency case), acceptance-criteria (7 numbered + partial/rejection branches), V/δ composition note, Relationship-to-cnos#425 section. **PASS.**

## AC5: Workflow header cites doctrine

**Oracle:** Workflow file top comments cite `BOX-AND-THE-RUNNER.md`, `delta/SKILL.md §8`, and the receipt path.

**Result:**

```
$ head -8 .github/workflows/publish-3.82.0-release.yml
# One-shot remote-runner workflow — publish the 3.82.0 GH release.
#
# Doctrine: docs/gamma/essays/BOX-AND-THE-RUNNER.md
# Skill:    src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §8 (Remote-runner delegation)
# Receipt:  .cdd/unreleased/426/remote-runner-receipt-publish-3.82.0.md
# Issue:    https://github.com/usurobor/cnos/issues/426
#
# Second use of the remote-runner-delegation primitive landed in cnos#425.
```

All three required citations present in lines 3–5, plus the issue link (line 6), plus second-use motivation context (lines 8 onward through line 32). **PASS.**

## AC6: Self-delete present

**Oracle:** Final step removes the workflow file with `git rm` + `git commit` + `git push`.

**Result:**

```
$ grep -A7 "Self-delete the workflow file" .github/workflows/publish-3.82.0-release.yml
      - name: Self-delete the workflow file
        run: |
          git fetch origin main
          git checkout main
          git config user.name "cnos-remote-runner[bot]"
          git config user.email "cnos-remote-runner@cnos.local"
          git rm .github/workflows/publish-3.82.0-release.yml
          git commit -m "release: remove one-shot publish-3.82.0-release workflow (cnos#426)"
          git push origin main
```

Final step contains all three required actions (`git rm`, `git commit`, `git push origin main`) plus two prerequisite worktree-switch actions (`git fetch origin main` + `git checkout main`) needed because step 1 checks out the tag rather than main. The commit message includes the issue ref. **PASS.**

## AC7: No CCNF kernel / CDS / CDR / handoff changes

**Oracle:**
- `git diff origin/main..HEAD -- src/packages/cnos.cdd/skills/cdd/CDD.md` returns 0 lines.
- `git diff origin/main..HEAD -- src/packages/cnos.cds/ src/packages/cnos.cdr/ src/packages/cnos.handoff/` returns 0 lines.
- `git diff origin/main..HEAD -- src/packages/cnos.cdd/` returns 0 lines (delta/SKILL.md unchanged in this cycle).

**Result:**

```
$ git diff origin/main -- src/packages/cnos.cdd/skills/cdd/CDD.md | wc -l
0
$ git diff origin/main -- src/packages/cnos.cds/ src/packages/cnos.cdr/ src/packages/cnos.handoff/ | wc -l
0
$ git diff origin/main -- src/packages/cnos.cdd/ | wc -l
0
```

CDD.md byte-identical; CDS/CDR/handoff packages byte-identical; entire cnos.cdd/ byte-identical (delta/SKILL.md §8 inherited from cnos#425 unchanged). **PASS.**

## AC8: No schemas / runtime / scripts changes

**Oracle:** `git diff origin/main..HEAD -- schemas/ src/go/ src/packages/cnos.core/ src/packages/cnos.eng/ src/packages/cnos.kata/ src/packages/cnos.cdd.kata/ scripts/` returns 0 lines.

**Result:**

```
$ git diff origin/main -- schemas/ src/go/ src/packages/cnos.core/ src/packages/cnos.eng/ src/packages/cnos.kata/ src/packages/cnos.cdd.kata/ scripts/ | wc -l
0
```

0 lines across all 7 excluded paths. **PASS.**

## AC9: No essay or README edits

**Oracle:** `git diff --name-only origin/main..HEAD -- docs/` returns 0 lines.

**Result:**

```
$ git diff --name-only origin/main -- docs/
(empty output)
```

No docs/ files modified. `BOX-AND-THE-RUNNER.md` and `docs/gamma/essays/README.md` inherited unchanged from cnos#425. **PASS.**

## AC10: No release.yml modifications

**Oracle:** `git diff origin/main..HEAD -- .github/workflows/release.yml` returns 0 lines.

**Result:**

```
$ git diff origin/main -- .github/workflows/release.yml | wc -l
0
```

`release.yml` byte-identical to origin/main. The build-job fix is deferred to a separate cycle and filed as cdd-iteration F1 with disposition `next-MCA`. **PASS.**

## Summary

| AC | Verdict | Evidence |
|----|---------|----------|
| AC1: workflow file exists with correct shape | PASS | 4 named/uses steps; all top-level keys present |
| AC2: checks out at the tag, not main | PASS | `ref: '3.82.0'` (resolves to fd1d654e) |
| AC3: body_path RELEASE.md | PASS | 1 hit |
| AC4: receipt has 6 fields filled | PASS | 6 sections populated; §5 post-run-fillable per hard rule 1 |
| AC5: workflow header cites doctrine | PASS | essay + skill + receipt + issue cited in lines 3–6 |
| AC6: self-delete present | PASS | git rm + commit + push (+ worktree switch prereq) |
| AC7: no CCNF kernel / CDS / CDR / handoff / cnos.cdd changes | PASS | 0/0/0 lines |
| AC8: no schemas / runtime / scripts | PASS | 0 lines across 7 excluded paths |
| AC9: no essay or README edits | PASS | 0 files in docs/ |
| AC10: no release.yml modifications | PASS | 0 lines |

All 10 ACs PASS. Bundle is review-ready for β.
