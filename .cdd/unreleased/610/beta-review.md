# β review — cnos#610

**Verdict:** REQUEST CHANGES (wake-invoked vocabulary: **iterate**)
**R1 note:** superseded — see `§R1` below for the current, final verdict (**converge**) after
independent re-verification of all 5 findings. This §R0 record is left as-is per the
append-don't-rewrite convention.

**Round:** §R0 (first β pass on this cell; γ's `gamma-scaffold.md` is this cell's R0 artifact, α's `self-coherence.md` is a single R0 pass with no prior `§R[N]` history — per `cnos.cdd/skills/cdd/delta/SKILL.md` §9.5 R0 artifact contract)
**Fixed this round:** n/a (first review pass)
**Branch CI state:** **red** — `Build` workflow fails at review SHA (see §CI status)
**Review base:** `origin/main` = `f7e9aaad34793dbea80c603e315e0ecc0760fdfa` (re-fetched synchronously immediately before computing this diff base; matches the SHA α's self-coherence.md §Review-readiness row 1 records — no drift since α's rebase)
**Review head:** `origin/cycle/610` = `b1406642ac0cdb3609f217de38626d789d50a716`
**Merge instruction:** none this round — do not merge. α fixes the findings below on `cycle/610`, appends `§R1` to `self-coherence.md`, and signals β for the next review pass.

---

## §CI status (rule 3.10)

```
$ gh run list --branch cycle/610 --json status,conclusion,workflowName,headSha --limit 5
[[headSha=b140664] install-wake golden: success]
[[headSha=b140664] Build: failure]   <-- review SHA
[[headSha=b140664] Build: failure]
```

`Build` (required workflow, includes `go build`/`go vet`/`go test` + the `cn cdd verify --unreleased` I6 gate + `test-fixtures.sh`) is **red** at the review SHA. Root cause isolated below (Finding 1). `install-wake golden` is green. Per rule 3.10 this alone forces the verdict to REQUEST CHANGES/RC regardless of the substantive findings below.

---

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | partial | Issue/design-doc status claims (cnos#493 open, cnos#609 merged) verified live via `gh issue view` — accurate. Test-count claim in self-coherence.md is inaccurate (Finding 4). |
| Canonical sources/paths verified | yes | `docs/development/design/cn-repo-install-MOCKS.md` Mock C/E, `cn-install-wake`, `cds-dispatch/SKILL.md` all resolve as cited. |
| Scope/non-goals consistent | yes | Non-goals (no cnos#493 impl, no renderer identity-flag changes) honored; the "no sigma default" non-goal phrase is ambiguous but the implemented default matches the pre-existing renderer's own convention and AC5's own backward-compat requirement — not a demonstrable violation (see Notes). |
| Constraint strata consistent | yes | Tier 1/2/3 skills named and applied per §Skills. |
| Exceptions field-specific/reasoned | n/a | No `.cdd/exceptions.yml` entry for #610 (confirmed via grep); none claimed. |
| Path resolution base explicit | yes | `RepoRoot`-relative paths throughout; `pkg.VendorPath` resolution verified in code + e2e. |
| Proof shape adequate | no | Per-AC oracle evidence (tests + e2e transcripts) is present and independently reproduced, but the issue-mandated `mock_parity` structured block is entirely absent (Finding 2). |
| Cross-surface projections updated | yes | `docs/guides/INSTALL-CDS.md` and `cmd_repo_install.go` help text both updated consistently with the new behavior (see Finding 3 for the process-side gap on this same file). |
| No witness theater / false closure | partial | AC evidence is real and independently reproduced (not vacuous) — but the quoted test-count figures don't match a re-run (Finding 4). |
| PR body matches branch files | yes | PR #620 body is the mechanical finalizer's draft-checkpoint stub ("Refs #610 ... Draft — not yet review-ready"), makes no claim to verify against. |
| γ artifacts present (gamma-scaffold.md) | yes | `.cdd/unreleased/610/gamma-scaffold.md` present on `origin/cycle/610`, read in full. Rule 3.11b satisfied. |

---

## §2.0 Issue Contract

### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| AC1 | dispatch render (Mock C1/C3): writes exactly `.github/workflows/cnos-cds-dispatch.yml` after a successful base install; base artifacts present | yes | **MATCH** | Independently reproduced: real `cn` binary, real `cn build`-produced index, scratch repo. `find .cn .github` shows `.cn/deps.json`, `.cn/deps.lock.json`, `.github/workflows/cnos-cds-dispatch.yml` all present after `--dispatch cds --agent acme --workflow-pat-secret ACME_WORKFLOW_PAT --bot-name acme-bot --bot-id 12345678`. `go test ./internal/repoinstall/... -run TestRun_DispatchCds_RendersWorkflow_ThenSurfacesLabelGap -v` → PASS. |
| AC2 | explicit identity (Mock C2): missing `--workflow-pat-secret`/agent/bot → fail early, nonzero exit, no partial workflow file | yes | **MATCH** | Independently reproduced: `cn repo install --dispatch cds --agent acme` (no `--workflow-pat-secret`) → exit=1, stderr names the missing flag, `.github` does not exist (`ls: cannot access '.github'`), base install artifacts (`.cn/deps.json`) already present. Sigma default (no `--agent`) does **not** trip the gate — reproduced via `TestRun_DispatchCds_SigmaDefault_NoIdentityFlagsRequired`. |
| AC3 | labels ensured: cnos#493 mechanism absence → actionable error, not silent skip | yes | **MATCH** | `cnos#493` confirmed OPEN (`gh issue view 493 --json state` → `OPEN`). No label-install mechanism exists anywhere in `src/go` or `src/packages/cnos.core/commands/` (repo-wide grep, re-verified independently). Reproduced e2e: after a successful acme render, `Run` still exits 1 naming `cnos#493` in both `err.Error()` and stderr. α's judgment call on the AC1/AC3 relationship is **accepted** — see §Judgment calls below. |
| AC4 | no substrate-binding leak + PR-only (Mock C4/C6): no `SIGMA_WORKFLOW_PAT`/`41898282`/`sigma@…`/`cds-dispatch-sigma` in a non-sigma render; never pushes `main`; states `workflow`-scope PAT requirement | yes | **MATCH** | `grep -inE 'sigma\|SIGMA_WORKFLOW_PAT\|41898282' .github/workflows/cnos-cds-dispatch.yml` on the independently-produced acme render → 0 hits. `grep -rn 'exec.Command("git"\|git push' src/go/internal/repoinstall/*.go` (non-test) → 0 hits; `repoinstall.go` never imports git tooling, only `os/exec` for the vendored renderer script. stdout carries `"needs \`workflow\` scope"` and `"never pushes to main"` — confirmed in the live e2e transcript. |
| AC5 | tenant prose-clean render: no `"today: sigma"` / `agent-admin-sigma` / `cn-sigma:` in a non-sigma render; `--agent sigma` still renders correctly (compat preserved); negative oracle non-vacuous | yes | **MATCH**, with an accepted documented deviation | `grep -n 'today: `sigma`\|agent-admin-sigma\|cn-sigma:' .github/workflows/cnos-cds-dispatch.yml` on the acme render → 0 hits (independently reproduced). `--agent sigma` render sha256 (`822bb9ec...c8f8b91`) matches both the committed golden and the live `.github/workflows/cnos-cds-dispatch.yml` at repo root — independently recomputed, all three identical. Negative-case oracle (`TestDispatchRenderer_ProseLeakGrep_CatchesPreFixSigmaPhrasing`) independently re-run → PASS (proves the grep is a real detector, not vacuous). α's judgment call on prose-clean overriding the pinned byte-identical-to-*pre-cycle*-golden floor is **accepted** — see §Judgment calls below. Minor prose-quality regression noted as Finding 5 (non-blocking). |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `docs/guides/INSTALL-CDS.md` | yes | consistent with code, but process gap | Content correctly describes the new behavior (verified against actual CLI output); not named in δ's pinned package-scoping row — Finding 3. |
| `src/go/internal/cli/cmd_repo_install.go` help/usage text | yes | consistent with code | Independently diffed; exit-code table, flag docs, and troubleshooting table entries all match observed behavior. |
| `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` | yes | consistent, one prose nit | Two named leaks (line 101, line 296) both resolved; line-296 relocation to the non-rendered `## Responsibilities (body reference)` appendix independently confirmed (the live rendered YAML does not contain the item-9 text). Line-101 fix produces slightly redundant phrasing for non-sigma agents — Finding 5. |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `gamma-scaffold.md` | yes (§5.1) | yes | Read in full; rule 3.11b satisfied. |
| `self-coherence.md` | yes | yes | Present, but section headers use a non-canonical `## §X` form that fails `cn cdd verify`'s literal section match — Finding 1. Missing the issue-mandated `mock_parity` block — Finding 2. |
| `alpha-clarification-needed.md` | situational | yes | Documents exactly the two interpretive calls this review addresses. |
| `beta-review.md` | yes (this file) | yes (being authored now) | — |
| `alpha-closeout.md` / `beta-closeout.md` / `gamma-closeout.md` | only at converge | not present | Correct for an R0/iterate round per `delta/SKILL.md` §9.5 — these land only after a `verdict: converge`. |

### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| `cnos.core/skills/write/SKILL.md` | prose edits (SKILL.md, help text) | yes (α claims) | mostly | Front-load/one-fix-per-leak discipline mostly held; line-101 fix is a minor exception (Finding 5). |
| `docs/development/design/cn-repo-install-MOCKS.md` (Mock C + E) | design source of truth | yes | yes | Mock C1–C6 all independently re-verified; Mock E3 (finalizer invokes installed `cn`, no `src/go`) independently re-verified through the full `cn repo install --dispatch cds` call path (not just the renderer's own golden) — `grep` on the acme render shows `Install cn (tenant-portable — no src/go build)` + `/tmp/cnos-cn-609/cn cell finalize`, and no `cd src/go`/`go build ./cmd/cn`. |

---

## §Judgment calls (α's two documented interpretive calls)

Both calls in `alpha-clarification-needed.md` are **accepted as-implemented**; no override requested this round.

**1. AC1/AC3 relationship — does `Run()` return non-nil whenever the cnos#493 label gap is present, even after an otherwise-successful render?**

Accepted. AC2's oracle text says "nonzero exit" explicitly; AC3's oracle uses the distinct phrase "an actionable error is returned (not a silent skip)," and the issue's own Proof plan lists "missing labels" as a negative case parallel to "missing identity" — both readings support treating the label gap as a real failure exit, not a warning-only path. Returning `nil` while only printing a `⚠` warning would make cnos#493's absence invisible to any script/CI wrapping this command, which directly contradicts "not a silent skip." Mock C's console transcript predates the AC3 requirement (γ's own scaffold and α's clarification note both establish this) and shows no exit-code signal either way, so it neither confirms nor refutes the reading — it is stale relative to AC3, not authoritative over it. Independently reproduced: the acme e2e run exits 1 naming `cnos#493` even though the render itself (file existence, content, no leaks) fully succeeds.

**2. AC5 prose-clean vs. the pinned byte-identical-to-golden backward-compat floor, for the one un-templatable line-296 leak**

Accepted. The issue's own AC5 text says "renders correctly (compatibility preserved)" — textually weaker than "byte-identical." γ's own scaffold flagged this exact leak as impossible to template honestly and offered the same two options α chose between. α's choice (move the specific historical citation to the non-rendered `## Responsibilities (body reference)` appendix) is independently verified: the live rendered `.github/workflows/cnos-cds-dispatch.yml` does not contain the item-9 incident-record text at all, confirming the appendix-boundary claim is real, not asserted. `--agent sigma`'s render is byte-identical to *this cycle's own* committed golden (sha256 `822bb9ec...c8f8b91`, independently recomputed and matching both the golden and the live workflow file) — the artifact CI actually gates against was updated in the same commit, so no real backward-compat break exists for any caller re-running the install. This resolves the tension in favor of the issue's primary directive (tenant-honest prose) without an observable compatibility cost.

---

## Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|
| 1 | `Build` CI workflow is red at the review SHA. Root cause: `self-coherence.md`'s section headers use the non-canonical `## §Gap` / `## §Skills` / `## §ACs` / `## §CDD Trace` form. Before `beta-review.md` exists, `cn cdd verify`'s `classifyCycleType` (`cdd-verify/ledger.go:424-446`) classifies #610 as `"small-change"` (only `self-coherence.md` present) and applies the **non-lenient** section check (`validateSections(..., forUnreleased=false)`), which requires the exact literal headers `## Gap`, `## Skills`/`## Mode`, `## ACs`/`## AC Coverage`, `## CDD Trace` (`sectionPresent` does an exact-line/prefix match, not fuzzy) — none of which match the `§`-prefixed form, producing `missing required sections: Gap Skills/Mode ACs/AC Coverage CDD Trace` as a hard FAIL. This is **not purely transient**: even once `beta-review.md` lands and the cycle re-classifies as `"triadic"` (lenient during the unreleased window), the *same* non-lenient check re-applies verbatim at release time (`checkTriadicArtifacts` also calls `validateSections(..., forUnreleased=false)`), so the defect will resurface when δ moves this cycle to `.cdd/releases/{X.Y.Z}/610/`. Cycle #608's own `self-coherence.md` explicitly documents having hit and fixed this identical failure mode ("Section headers use the bare form (`## Gap` not `## §Gap`) ... required for `cn cdd verify`'s I6 gate to actually pass") — #610 did not follow that established precedent. | `gh run view 28837384470 --log-failed` (Build, headSha=b140664, conclusion=failure, `## Summary: 181 passed, 1 failed, 119 warnings`); `grep -n "^## §" .cdd/unreleased/610/self-coherence.md`; `src/packages/cnos.cdd/commands/cdd-verify/ledger.go` lines 424-446, 494-512, 545-587; `.cdd/unreleased/608/self-coherence.md`'s own deviation note. | D | mechanical, contract |
| 2 | Issue-mandated `mock_parity` block is entirely absent from `.cdd/unreleased/610/`. The issue's "Parity requirement" section states plainly: "Closeout MUST carry `mock_parity` rows for C1–C6 with `missed: 0`, plus an explicit AC5 tenant-prose-clean row." The design doc's "Receipt parity contract" (`docs/development/design/cn-repo-install-MOCKS.md`) calls this "a hard gate: any row with `verdict: miss` blocks convergence to `status:review`," and states "a missing ID is itself a `miss`." No `mock_parity:` YAML block exists anywhere under `.cdd/unreleased/610/` — confirmed by direct grep. Sibling cycles for the *same* design doc carry this block in the equivalent artifact: #608's `self-coherence.md` (α-authored, same role) embeds a full `mock_parity` block for A1-A5/B1-B4/E1; #609's `gamma-closeout.md` does the same. γ's own β-dispatch prompt in `gamma-scaffold.md` (item 7) explicitly told β to check for exactly this. | `grep -rn "mock_parity" .cdd/unreleased/610/` → 0 hits (outside `gamma-scaffold.md`'s own prose reference to the requirement); `docs/development/design/cn-repo-install-MOCKS.md` §"Receipt parity contract"; `.cdd/unreleased/608/self-coherence.md:661-` (precedent block); `.cdd/unreleased/610/gamma-scaffold.md` §"β prompt" item 7. | D | contract, judgment |
| 3 | Rule 7 (implementation-contract coherence) drift: `docs/guides/INSTALL-CDS.md` is touched but is not named in δ's pinned "Package scoping" row (which lists only `internal/repoinstall/repoinstall.go`+`_test.go`, `internal/cli/cmd_repo_install.go`+`_test.go`, and the `cds-dispatch` SKILL.md + golden). α self-disclosed this in `self-coherence.md` §CDD Trace and judged it in-bounds peer-enumeration, but no `.cdd/unreleased/610/gamma-clarification.md` updating the pinned row exists on the branch, which Rule 7 requires for any diff outside a pinned axis. The content of the fix itself is correct and arguably necessary (the pre-fix doc text was itself a stale/false claim that would otherwise be a rule-3.13 honest-claim violation) — the gap is purely procedural (missing γ ratification), not substantive, and does not require reverting the doc fix. | `git diff origin/main..HEAD -- docs/guides/INSTALL-CDS.md`; `gamma-scaffold.md` §"Implementation contract" row 3 ("Package scoping"); `find .cdd/unreleased/610 -name gamma-clarification.md` → not found. | C | contract |
| 4 | Honest-claim / reproducibility (rule 3.13a): `self-coherence.md` §Review-readiness states "26 tests in internal/repoinstall, 39 in internal/cli" and "65 test functions, 65 PASS, 0 FAIL" at the cited HEAD `a5ed0f2`. Independently re-run at that exact SHA (via a throwaway worktree) and at the final HEAD: `go test ./internal/repoinstall/... -v \| grep -c '^--- PASS'` → 27 (not 26); `go test ./internal/cli/... -v \| grep -c '^--- PASS'` → 45 (not 39); total 72, 0 FAIL (all passing, consistent with the claimed "0 FAIL", but the specific counts are wrong by 1 and 6 respectively). | `go test ./internal/repoinstall/... ./internal/cli/... -v` run independently at both `a5ed0f2` and `b140664` (identical counts at both); baseline count on `origin/main` (23 + 43 = 66) plus the diff's net-new tests (+4/+2 = +6) arithmetically confirms 72 is correct and 65 is not. | C | honest-claim |
| 5 | Prose-quality regression in the AC5 line-101 fix (non-blocking): `"You substrate-execute as `{agent}` (today: `{agent}`; future: per-package bot accounts per cnos#449 follow-up)."` renders for any non-sigma agent as `"...as acme (today: acme; future: ...)"` — grammatically redundant/self-referential, and it silently drops the original sentence's actual meaning ("today only sigma has bindings; other agents' bot accounts are a future item"). Not a leak (grep-clean; AC5's binding oracle still passes) and not blocking, but a readability regression relative to γ's own suggested alternatives ("dropping the parenthetical entirely, or rephrasing agent-neutrally"). | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md:101`; independently rendered via `cn repo install --dispatch cds --agent acme ...` e2e run. | B | judgment |

## Regressions Required (D-level only)

- **Finding 1** — positive: `./cn cdd verify --unreleased --exceptions .cdd/exceptions.yml` exits 0 (or only warns, not fails) on `cycle/610` after the header rename. Negative: the identical command on the *current* HEAD (`b140664`) fails with `missing required sections: Gap Skills/Mode ACs/AC Coverage CDD Trace` (already reproduced above) — regression pair is the before/after of this exact command against the two commits.
- **Finding 2** — positive: `.cdd/unreleased/610/self-coherence.md` contains a `mock_parity` block with `summary.missed: 0` covering C1–C6 + the AC5 row, each `evidence` field citing a real test name or transcript already present in §ACs. Negative: the current HEAD has no such block (already reproduced above via grep).

## Notes

- **"No sigma default" non-goal phrasing** — the issue's Non-goals list "no sigma default" alongside "no direct main push"/"no autonomous write." Read literally this could mean "the `--agent` flag must never silently default to sigma," which the implementation does not honor (`resolveDispatchAgent` defaults empty `--agent` to `"sigma"`). However: (a) this matches the pre-existing `cn-install-wake` renderer's own established default-agent convention (not invented by α), (b) AC5's own positive case ("`--agent sigma` still renders correctly, compatibility preserved") presumes a sigma-identified path continues to work, and the wave-level non-goals doc's parallel phrase is "sigma-only default" (about the *system* not being sigma-exclusive, not about this one flag's blank-value behavior). Not raised as a blocking finding — no demonstrable incoherence, and blocking it would be a phantom-blocker under rule 3.5. Flagged here for γ/operator awareness in case the stricter reading was intended.
- Mock E3 (tenant-portable finalizer, from the issue's scope-extension comment) was independently re-verified through the *installer's own call path* (not just the renderer's direct golden), per γ's scaffold instruction — confirmed present (`Install cn (tenant-portable — no src/go build)` + `cn cell finalize` via the installed binary, no `cd src/go`/`go build ./cmd/cn`).
- All five findings above are mechanically small (header rename; add one YAML block; one paperwork artifact; correct two numbers; optional prose tightening) — none require re-deriving evidence already collected in `self-coherence.md` §ACs. Expect a fast R1.

---

## §R1 — independent re-verification of α's fix round

**Verdict: CONVERGE.**

**Round:** §R1 (second β pass on this cell)
**Fixed this round:** all 5 §R0 findings (2×D, 2×C, 1×B)
**Review base:** `origin/main` = `f7e9aaad34793dbea80c603e315e0ecc0760fdfa` (re-checked — unchanged since §R0)
**Review head:** `origin/cycle/610` = `387b01fbde25fcc100e862ccafed29332ab00e39`
**Branch CI state:** **green** — independently checked via `gh run list --branch cycle/610 --json status,conclusion,workflowName,headSha --limit 6` at this exact head: both `Build` and `install-wake golden` report `status=completed, conclusion=success` at `headSha=387b01f...`. (At the start of this review the `Build` run for this head was still `in_progress`; it was polled to completion rather than assumed.)

Every finding below was re-derived independently — a fresh `cn` binary was built from this HEAD (`go build -o /tmp/cn-beta-r1 ./cmd/cn`), a fresh `dist/packages/index.json` via `cn build`, and fresh scratch git repos for every e2e run. Nothing here is copied from α's transcript.

### Finding 1 (D) — non-canonical headers → CI red — **FIXED, confirmed**

`grep -n "^## " self-coherence.md` shows exactly the canonical set: `## Gap`, `## Skills`, `## ACs`, `## Self-check`, `## Debt`, `## CDD Trace`, `## Review-readiness` (plus the new `## R1 — ...` subsection, which is prose, not a required-section header). No `## §...` form remains.

`./cn cdd verify --unreleased --exceptions .cdd/exceptions.yml` from repo root at this HEAD → `## Summary: 184 passed, 0 failed, 121 warnings (305 total)`, exit 0.

Also independently reproduced the **non-lenient/small-change path** α cites (the path that will re-apply at release time): copied the repo to a scratch dir, moved `beta-review.md` out of `.cdd/unreleased/610/`, re-ran the same binary → `Checking small-change cycle #610` / `✅ self-coherence.md (small-change #610)` / `✅ self-coherence.md sections — basic section validation passed`, no "missing required sections" line, exit 0. Confirms the fix holds under both classification paths, not just the currently-lenient triadic one.

Beyond local reproduction: a real CI run exists at this exact head and is green (see Branch CI state above) — this is the actual gate, not just a simulated one.

### Finding 2 (D) — missing `mock_parity` block — **FIXED, confirmed**

`mock_parity` block present under §CDD Trace with rows `C1, C2, C3, C4, C5, C6, AC5` and `summary: {matched: 7, exceeded: 0, missed: 0, exceed_justified: true}`.

Spot-checked 3 rows (not just 2) against real evidence, not fabrication:
- **C1**: cites `TestRun_DispatchCds_RendersWorkflow_ThenSurfacesLabelGap` — confirmed this function exists at `internal/repoinstall/repoinstall_test.go:951` and is the same test cited in §ACs AC1.
- **C2**: cites `TestRun_DispatchCds_MissingIdentity_FailsEarlyNoPartialWrite` (`repoinstall_test.go:1070`) + `TestRepoInstall_DispatchCds_MissingIdentity_CliWiring` (`internal/cli/cmd_repo_install_test.go:383`) — both confirmed to exist and match §ACs AC2's own citations.
- **AC5** row: cites `TestDispatchRenderer_ProseLeakGrep_CatchesPreFixSigmaPhrasing` — confirmed at `repoinstall_test.go:1114`, matching §ACs AC5's negative-oracle citation.

No row cites a test or transcript that isn't already present in §ACs above it in the same file.

### Finding 3 (C) — missing `gamma-clarification.md` — **FIXED, confirmed**

`.cdd/unreleased/610/gamma-clarification.md` exists, names `docs/guides/INSTALL-CDS.md` (commit `a9c194f`), gives the real peer-enumeration reasoning (three sibling surfaces stating the same stale "fails explicitly" claim: help text, `SKILL.md` prose, and this guide), and formally ratifies retroactively adding the file to γ's pinned "Package scoping" row. Disposition is "Ratified," no content revert requested — consistent with §R0's own finding that the doc content itself was already correct.

### Finding 4 (C) — wrong test counts — **FIXED, confirmed**

Independently re-ran (not copied from α's numbers):
```
$ go test ./internal/repoinstall/... -v | grep -c '^--- PASS'   → 27
$ go test ./internal/cli/... -v | grep -c '^--- PASS'           → 45
$ go test ./internal/repoinstall/... ./internal/cli/... -v 2>&1 | grep -c '^--- FAIL'  → 0
```
27 + 45 = 72, matching the corrected figures now in §Review-readiness rows 3 and 13 exactly. The R0 finding's own arithmetic (66 baseline + 6 net-new = 72) is consistent with this independent count.

### Finding 5 (B) — redundant line-101 prose — **FIXED, confirmed**

Rendered the SKILL.md directly through `cn-install-wake` (not through the test suite, to get an unmediated look) for both a non-sigma agent and the sigma default:

- `--agent acme ...` render: `substrate-execute as \`acme\` (bot-account bindings are supplied via --agent/--workflow-pat-secret/--bot-name/--bot-id; only the default agent carries a built-in binding today, any other agent supplies its own explicitly; future: per-package bot accounts per cnos#449 follow-up).` — no "(today: acme)" tautology, true and non-confusing.
- `grep -inE 'sigma|SIGMA_WORKFLOW_PAT|41898282'` on that acme render → 0 hits. α's claim that the literal word "sigma" was deliberately avoided in the new prose (to not regress the AC4 blanket-grep invariant) is **verified true**, not just claimed — confirmed by direct grep on a freshly-rendered file, not by reading the source text alone.
- `--agent sigma` (default) render sha256: `b80906ec2dba15dcdf3a5852f3b72ed79b8cfa8119bf773d95209ba9ad07fd53`. This matches, byte-for-byte, both the committed golden (`cnos-cds-dispatch.golden.yml`) and the live checked-in `.github/workflows/cnos-cds-dispatch.yml` — all three files independently re-hashed and identical. Matches the new SHA α cites in §R1 exactly.
- `git diff 3fb8e47..387b01f` for `SKILL.md`/golden/live-workflow shows exactly one line changed per file (the line-101 sentence), nothing else touched.

### Nothing else regressed

Full AC1–AC5 e2e re-run from scratch (fresh `cn` binary, fresh `cn build`-produced index, fresh scratch git repos, no reuse of any prior artifact):
- **AC1/AC3**: `--dispatch cds --agent acme --workflow-pat-secret ... --bot-name ... --bot-id ...` → `.cn/deps.json`/`.cn/deps.lock.json` written, then `.github/workflows/cnos-cds-dispatch.yml` rendered, then exit 1 naming `cnos#493` — ordering and behavior unchanged.
- **AC2**: same command minus `--workflow-pat-secret` → exit 1, stderr names the missing flag, `.github` does not exist, `.cn/deps.json` does exist — unchanged.
- **AC4/AC5**: sigma-family grep and prose-leak grep both 0 hits on the acme render; `grep` for `git push`/`exec.Command("git"` in `repoinstall/*.go` (non-test) → 0 hits — unchanged.

`git diff --stat 3fb8e47..387b01f` shows exactly 5 files touched: `gamma-clarification.md` (new), `self-coherence.md`, `.github/workflows/cnos-cds-dispatch.yml`, `cds-dispatch/SKILL.md`, `cnos-cds-dispatch.golden.yml` — precisely the 5 findings, no scope creep, no new files or code paths beyond what the findings required.

### Conclusion

All 5 findings are genuinely fixed, independently re-derived, with no regressions and no scope creep in the fix diff. **This is the final review verdict for this cell.** Per the wake-invoked-mode protocol, δ next dispatches γ+α+β for closeouts (β is not responsible for closeouts in this pass, and does not merge to main — merge happens later via draft PR #620 after external review).

**Final HEAD:** `387b01fbde25fcc100e862ccafed29332ab00e39`
