# β Review — Cycle #344-B

**Reviewer:** β (`beta@cdd.cnos`)
**Branch:** `cycle/344-b`
**Base SHA:** `fba356e1`
**Implementation SHA:** `81813964`
**Issue:** #344 (Cycle B ACs: B.AC1–B.AC5)
**Review skill:** `review/SKILL.md` + sub-skills

---

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | α self-coherence correctly marks cycle as substantial, provisional close-out pending β. No draft-as-current conflation. |
| Canonical sources/paths verified | yes | All file paths in self-coherence.md map to actual files in the diff. `activation/SKILL.md §10.1–§10.2` is the governing spec; paths resolve correctly. |
| Scope/non-goals consistent | yes | Cycle B is purely additive (new `activation/templates/` directory). No existing surfaces modified. Non-goals (other transports, tsc adoption) not violated. |
| Constraint strata consistent | yes | No hard gates violated. Secret names follow `CDD_` namespace. Template files use `${{ secrets.CDD_... }}` throughout. |
| Exceptions field-specific/reasoned | yes | Four named debt items in §Debt, each with scope and rationale. No blanket exceptions to hard gates. |
| Path resolution base explicit | yes | All template paths are caller-relative; README instructions give explicit copy-target paths (`.github/cdd/notify.sh`, `.github/workflows/`). |
| Proof shape adequate | yes | This is a template/docs cycle. No runtime behavior claims; no checker claims. `bash -n` and `yaml.safe_load` are adequate for syntax gates. B.AC4 word-count oracle is adequate for the prose constraint. |
| Cross-surface projections updated | n/a | No existing status surfaces, CI jobs, schema, or runtime commands were changed. New templates are additive. `activation/SKILL.md §10.2` already names `templates/telegram-notifier/` as the Cycle B deliverable. |
| No witness theater / false closure | yes | No enforcement claim is made. Templates are declared as reference implementations (copy-paste starting points), not enforced infrastructure. §Debt is honest about heuristic approximations. |
| PR body matches branch files | n/a | No PR opened; review dispatched directly from issue. Self-coherence.md pre-review gate enumerates 8 files; diff shows 8 files matching the enumeration. |

**Phase 1 result:** All rows yes/n/a. Proceed to Phase 2.

---

## §2.0 Issue Contract

### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| B.AC1 | Telegram notifier works end-to-end (script + workflow, 4 events, no hardcoded tokens) | yes | **MET** | `notify.sh` 133 lines, `case` dispatch on 4 events, `${{ secrets.CDD_... }}` in all workflow env blocks. No hardcoded tokens anywhere. `bash -n` passes. |
| B.AC2 | Artifact validator catches missing files (cdd-artifact-validate.yml) | yes | **MET** | `cdd-artifact-validate.yml` has explicit `exit 1` + named error if `validate-release-gate.sh` not found. Triggers on `cycle/**` and `main`. YAML validates. |
| B.AC3 | Workflow templates are copy-pasteable (no tenant-specific hardcoding) | yes | **MET** | All YAML uses `${{ secrets.CDD_TELEGRAM_BOT_TOKEN }}` / `${{ secrets.CDD_TELEGRAM_CHAT_ID }}`. `# Customize:` comments mark the only tenant-variable locations. No repo names, user names, or hardcoded paths. |
| B.AC4 | README walkthrough ≤300 words (bot registration → secret setup → wiring) | partial | **CONTRACT MISMATCH — see Finding #1** | Issue AC names `cnos:cdd/activation/templates/README.md`. That file is 171 words and does NOT contain bot registration or secret setup walkthrough — it is a directory overview that links to `telegram-notifier/README.md`. The walkthrough itself lives in `telegram-notifier/README.md` (266 words, ≤300 ✓). |
| B.AC5 | Notifier honors transport-agnostic event vocabulary (event names match activation/SKILL.md §10) | yes | **MET** | `case "$EVENT" in` handles exactly `cycle-open`, `beta-verdict`, `cycle-rc`, `cycle-merge` — matching `activation/SKILL.md §10.1` table rows 1–4. Unknown-event path exits 0 with warning. |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `activation/templates/telegram-notifier/notify.sh` | yes | present | New file — primary deliverable |
| `activation/templates/telegram-notifier/cdd-notify.yml` | yes | present | New file |
| `activation/templates/telegram-notifier/README.md` | yes | present | New file |
| `activation/templates/github-actions/cdd-artifact-validate.yml` | yes | present | New file |
| `activation/templates/github-actions/cdd-cycle-on-merge.yml` | yes | present | New file |
| `activation/templates/README.md` | yes | present | New file |
| `.cdd/unreleased/344-b/self-coherence.md` | yes | present | CDD artifact |
| `.cdd/unreleased/344-b/alpha-closeout.md` | yes | present | CDD artifact (provisional) |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `self-coherence.md` | yes | yes | Complete through step 7 of CDD Trace. Review-readiness signal present. |
| `alpha-closeout.md` | yes | yes | Marked provisional per α/SKILL.md §2.8 bounded-dispatch fallback. Acceptable. |
| `beta-review.md` | yes | this file | Written incrementally per review/SKILL.md. |
| `beta-closeout.md` | conditional | not yet | Written after verdict on this file. |

### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| `CDD.md` 3.15.0 | Cycle dispatch | yes | yes | Lifecycle followed; CDD Trace through step 7. |
| `alpha/SKILL.md` | α role | yes | yes | Review-readiness signal populated with pre-review gate table. |
| `activation/SKILL.md` | B.AC5, B.AC1 | yes | yes | §10.1 event vocabulary used directly in `case` dispatch; §10.2 adapter contract implemented. |
| Shell conventions | Tier 2 | yes | yes | `bash -n` run; `set -euo pipefail`; exit codes 0/1 per contract. |
| YAML conventions | Tier 2 | yes | yes | `yaml.safe_load` validated all three YAML files. |
| Markdown authoring | Tier 2 | yes | yes | Word counts checked (266, 171). |

---

## Phase 2b: Diff and Context Inspection

**2.1.1 Structural closure:** No structural-prevention claims made. Template scripts have all input sources documented and handled (secrets absent, unknown event, unset CDD_EVENT — all handled with `exit 0` + warning).

**2.1.2 Multi-format parity:** Event names are consistent across `notify.sh` `case` dispatch, `cdd-notify.yml` CDD_EVENT values, README description, and `activation/SKILL.md §10.1` table. No format divergence.

**2.1.3 Snapshot consistency:** No snapshot tests in this project; n/a.

**2.1.4 Stale-path validation:** All files are new additions. No renamed/moved files. No existing consumers to check.

**2.1.5 Branch naming:** `cycle/344-b` follows convention. Branch matches issue dispatch.

**2.1.6 Execution timeline:** `notify.sh` is called as a subprocess from GitHub Actions step. Secrets are injected as environment variables. Python3 used for JSON encoding — dependency on `ubuntu-latest` runner (declared as Debt item 2; acceptable for reference impl).

**2.1.7 Derivation vs validation:** No single-source-of-truth generation claim; n/a.

**2.1.8 Authority-surface conflicts:**
- `activation/SKILL.md §10.2` contract item 3: "Returns exit code 0 on successful delivery, non-zero on failure." `notify.sh` correctly returns exit 1 on curl failure and on Telegram API non-ok response, and exit 0 on successful delivery or graceful skip. Consistent.
- Self-coherence §ACs line for B.AC4 cites `activation/templates/telegram-notifier/README.md` but issue AC names `activation/templates/README.md` — authority conflict noted (Finding #1).

**2.1.9 Module-truth audit:** New directory; no pre-existing module assumptions to audit.

**2.1.10 Contract-implementation confinement:** `notify.sh` handles all 4 canonical events plus empty-event and unknown-event paths. No unclaimed input silently accepted. The `if: >` condition in `notify-beta-verdict` job has an operator-precedence defect (Finding #2).

**2.1.11 Architecture leverage:** Templates are the right abstraction. The tenant-deployment model (copy-paste reference implementation) is appropriate for the scope. Higher-leverage boundary (e.g., vendored skill enforcement) is deferred to Cycle C per the issue.

**2.1.12 Process overhead:** Templates add copy-pasteable CI workflows. The `cdd-artifact-validate.yml` requires `scripts/validate-release-gate.sh` to exist in the tenant repo — the README mentions this indirectly but does not prominently warn tenant about this dependency. This is a B-level observation (Finding #3).

**2.1.13 Design constraints:** No project design constraints document loaded; n/a.

## Phase 2c: Architecture Check

| Check | Result | Notes |
|---|---|---|
| Reason to change preserved | yes | New `templates/` directory has one reason to change: reference implementations evolve. No existing modules touched. |
| Policy above detail preserved | yes | `activation/SKILL.md §10` remains the policy; `templates/telegram-notifier/` remains the detail/example. |
| Interfaces remain truthful | yes | `notify.sh` implements all 5 §10.2 contract items. The interface promises only what all implementations must support. |
| Registry model remains unified | n/a | No registry model changes. |
| Source/artifact/installed boundary preserved | yes | Templates are source artifacts in cnos; tenants copy them to their repos. The source/copy boundary is explicit in every README. |
| Runtime surfaces remain distinct | n/a | No runtime surfaces in cnos itself — templates operate in tenant repos. |
| Degraded paths visible and testable | yes | `notify.sh` has explicit warn-and-exit-0 paths for all degraded states. Logged to stdout for CI capture. |

---

## Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|
| 1 | B.AC4 walkthrough is in `telegram-notifier/README.md`, not in `activation/templates/README.md` as issue AC names. Self-coherence cites the sub-README (correctly describing the work done) but the issue AC names the top-level README. The top-level README (171 words) is a directory overview, not a walkthrough. | Issue B.AC4: "cnos:cdd/activation/templates/README.md walks bot registration…"; `templates/README.md` lines 1–37 contains no bot registration or secret setup prose; self-coherence.md §ACs B.AC4 evidence cites `activation/templates/telegram-notifier/README.md` | C | contract / honest-claim (3.13b) |
| 2 | `cdd-notify.yml` `notify-beta-verdict` job `if:` expression has operator-precedence defect: `A && B && C \|\| D` evaluates as `(A && B && C) \|\| D`, meaning any push to any branch (including main) that modifies `beta-review.md` will trigger a notification. Intended logic is `(A && B) && (C \|\| D)`. | `cdd-notify.yml` lines 60–64: `startsWith(...) && created != true && contains(...added...) \|\| contains(...modified...)` — the final `\|\| D` clause is not parenthesized with the cycle-branch guard | B | mechanical / judgment |
| 3 | `cdd-artifact-validate.yml` depends on `scripts/validate-release-gate.sh` existing in the tenant repo, but this prerequisite is not mentioned in `activation/templates/README.md` or in the workflow file header's "Customize" comment. A tenant copying the file without the script will get a clear error, but the dependency is a surprise. | `cdd-artifact-validate.yml` lines 31–38; `activation/templates/README.md` contains no mention of `validate-release-gate.sh`; the workflow header comment mentions only `VALIDATE_MODE` under "Customize" | A | judgment |

## Regressions Required (D-level only)

None — no D-level findings.

---

## Notes

- **CI state:** No CI runners configured for cycle branches in this environment (per self-coherence §Row 10 note). All mechanical gates (`bash -n`, `yaml.safe_load`, word counts, no-token grep) verified locally in this review pass. Branch CI state is effectively unverifiable; approval marked provisional.
- **B.AC4 interpretation:** The issue AC text is unambiguous (`cnos:cdd/activation/templates/README.md`). However, the walkthrough content exists and meets the ≤300 word constraint in the sub-README. The gap is structural (wrong file), not substantive (content exists). Finding #1 is severity C.
- **Finding #2 fix:** Wrap the last two `contains()` clauses in parentheses: `(contains(...added..., 'beta-review.md') || contains(...modified..., 'beta-review.md'))` and rejoin with `&&` to the preceding conditions.
- **Finding #3 fix:** Add one line to `activation/templates/README.md` under the `cdd-artifact-validate.yml` row: "Requires `scripts/validate-release-gate.sh` in your repository (see cnos `scripts/validate-release-gate.sh` for reference)."

---

## Phase 3: Verdict

**Verdict:** REQUEST CHANGES

**Round:** 1
**Fixed this round:** n/a (round 1)
**Branch CI state:** provisional (no CI runners available for cycle branches; all mechanical gates verified locally)
**Merge instruction:** Do not merge until all findings below are fixed on `cycle/344-b`, then re-dispatch β for round 2.

### Required fixes before merge

**Finding #1 (C — contract/honest-claim):** B.AC4 names `cnos:cdd/activation/templates/README.md` as the file that "walks bot registration, secret setup, first activation, end-to-end ≤300 words." The top-level `templates/README.md` currently contains a directory overview (no walkthrough content). Fix: either move the walkthrough content into `templates/README.md` (and update or keep the sub-README), or update the issue AC evidence in self-coherence.md to accurately reflect the structural decision (walkthrough lives in `telegram-notifier/README.md` with a pointer from the top-level README).

The preferred fix is the structural one: the issue AC was written before the two-level README structure was designed. `self-coherence.md §ACs B.AC4` correctly identified what α built; the issue AC wording was under-specified. α should append to `self-coherence.md` with a fix-round note naming the resolution: either confirm that the sub-README satisfies B.AC4 and update the top-level README to include the walkthrough summary prose (≤300 words total), or document the structural decision explicitly.

**Finding #2 (B — mechanical/judgment):** `cdd-notify.yml` `notify-beta-verdict` job `if:` expression evaluates as `(A && B && C) || D` due to `&&`/`||` precedence, firing on main-branch pushes that modify `beta-review.md`. Fix the expression to:

```yaml
if: >
  startsWith(github.ref, 'refs/heads/cycle/') &&
  github.event.created != true &&
  (contains(join(github.event.commits.*.added, ','), 'beta-review.md') ||
   contains(join(github.event.commits.*.modified, ','), 'beta-review.md'))
```

**Finding #3 (A — judgment):** Add the `validate-release-gate.sh` prerequisite to `activation/templates/README.md` under the `cdd-artifact-validate.yml` table row. One sentence is sufficient.

### Closing the search space

No remaining D-level finding was found in the relevant contract. All three findings are correctable without design decisions outside issue scope — no issue deferral required. The implementation is structurally sound: shell syntax valid, YAML valid, event vocabulary matches `activation/SKILL.md §10.1`, no hardcoded tokens, all adapter contract items implemented. Round 2 narrows to the three findings above.
