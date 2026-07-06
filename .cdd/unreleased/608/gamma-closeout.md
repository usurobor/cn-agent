# γ close-out — cycle #608

**Scope of this artifact.** This is γ's process-gap-audit retrospective at the β-converge boundary of the wake-invoked routing sequence (`delta/SKILL.md` §9.3 step 5 — δ dispatches γ once more after β's `verdict: converge` lands). It is written and pushed to `origin/cycle/608` **before** the cycle-PR is opened / merged; issue #608 is correctly still `OPEN` at this point (confirmed: `gh issue view 608 --json state` → `OPEN`). It is **not** the full `gamma/SKILL.md` §2.10 release-time closure declaration — `RELEASE.md` authorship, the `.cdd/unreleased/608/` → `.cdd/releases/{X.Y.Z}/608/` move, remote-branch deletion, hub-memory update, and the issue-close-state assertion all happen at actual release/merge time, which has not occurred yet. See §Triage carryforward below for what a later γ invocation (post-merge) still owes.

---

## Cycle summary

**Issue:** [cnos#608](https://github.com/usurobor/cnos/issues/608) — "cds-install Sub 1 (Cn=1 / PR 1): implement `cn repo install` — base installer (dispatch none)". Sub 1 of the wave tracked at #607 (first-class CDS repo installer); PR2 (#609, renderer generalization), PR3 (#610, dispatch install), PR4 (#611, bootstrap delegation) remain out of scope and undispatched.

**What shipped:**
- A new kernel command, `cn repo install`, registered via the existing noun-verb resolver (`RepoInstallCmd`, `Spec().Name == "repo-install"`, `Source: SourceKernel`, `NeedsHub: false`) — runs in a plain `git init`-only repo before any `.cn/` hub exists.
- `internal/repoinstall/` (new package) + `internal/cli/cmd_repo_install.go` (thin dispatch wrapper), reusing `internal/restore`'s existing `GenerateLockFromIndex`/`Restore` directly (zero reimplementation of lock generation, SHA-256 verification, or tar extraction).
- Deterministic `.cn/deps.json` (`cn.deps.v1`) + `.cn/deps.lock.json` (`cn.lock.v2`, SHA-256-pinned, exact versions) + name-based `.cn/vendor/packages/<name>/` restore for the default triple (`cnos.core`, `cnos.cdd`, `cnos.cds`).
- `--release latest|<tag>`, `--index <path-or-url>`, `--dry-run`; idempotent reruns; a hard dispatch guard (`--dispatch cds` fails explicitly, writes nothing, until #609 lands); no agent-hub scaffold (no `spec/SOUL.md`/`agent/`/`threads/`/`state/`, cnos#606 C4); no `.github/workflows/` touch in base mode.
- `docs/development/design/cn-repo-install-MOCKS.md` landed verbatim from the operator-reviewed design branch (`claude/cds-install-guide-2ka54j`), and `docs/guides/INSTALL-CDS.md` rewritten so `cn repo install` is the canonical Layer-1 path.
- All AC1–AC11 closed with PASS-backed evidence; `mock_parity` block 10/10 matched (A1–A5, B1–B4, E1), 0 missed, 0 exceeded.

**Review arc:** R0 (`iterate`) → R1 fix → R1 re-review (`converge`). One blocking finding total across the whole cycle (β R0 Finding 1); the rest of both review rounds is confirmation, not correction. CI green at every reviewed HEAD; final head `7bb75914` (β R1 converge) — `Build` conclusion `success`.

---

## Process-gap audit

### The R0 → iterate → R1 → converge loop

β's R0 blocking finding: `RepoInstallCmd.Run` used `inv.HubPath` (populated by `main.go`'s unbounded upward `discoverHub()` walk, with no verification the found directory is even a git repository) as the install root whenever non-empty — silently bypassing `gitRepoRoot` and installing into an unrelated ancestor `.cn/` directory even when cwd was not inside any git repository at all. This directly contradicted AC1's own text: *"Fails with a clear error if not at a git repo root (does not silently walk up or scaffold)."* β reproduced it against the built binary, not against the test suite or the self-coherence report.

**Is this a scaffold gap, review-loop-working-as-intended, or something else? Both of the first two, honestly — not either/or:**

1. **It is a real, specific scaffold gap.** γ's own scaffold (`.cdd/unreleased/608/gamma-scaffold.md`) reproduced AC1's negative clause verbatim in its "Mock A" fallback table (`A1 | Fails with a clear error if not at a git repo root (does not silently walk up or scaffold).`) — the exact language was already sitting in the document γ authored. But the same scaffold's "AC oracle approach" table, the section α was actually told to build tests against, gave AC1 a concrete check naming only the positive path ("assert exit 0 and the three-file diff exists"). The negative half of the AC's own text never got translated into a required concrete check. This is an internal inconsistency inside γ's own artifact, not something α or β introduced.
2. **α's test design compounded it.** The negative test α wrote (`TestRepoInstall_NotAGitRepo_FailsClearly`) satisfied AC1's literal text via a hand-constructed `Invocation{HubPath: ""}` — bypassing `main.go`'s real `discoverHub()` walk entirely. A test that stubs around the exact boundary an invariant concerns will pass regardless of whether the boundary itself is correct.
3. **The review loop caught it exactly as designed, at normal cost.** β's independent re-verification discipline (`beta/SKILL.md` Role Rule 6 — anchor oracle evidence on code, not on a report; build the real binary, don't trust the test suite's framing) is precisely the safety net that exists for this class of gap: a test that satisfies the letter of an AC while missing its substance. One round (R0 → R1 → converge) is the designed cost of catching this kind of thing, not evidence of process failure. The fix itself was a net deletion of the buggy conditional (not new complexity), and β independently reverted it in an isolated worktree to confirm the new tests actually fail against the pre-fix code before converging — a level of rigor that closes the loop cleanly.

**No formal `gamma/SKILL.md` §2.8 cycle-iteration trigger fired:** review rounds = 1 (not > 2); total findings across both rounds ≈ 7 (1 blocking + 6 non-blocking observations, below the mechanical-overload trigger's ≥10 floor and not majority-mechanical anyway); no avoidable tooling/environment failure (this was a substantive design/implementation bug, not a flake); `beta/SKILL.md`'s own loaded skill (Role Rule 6) is what caught the gap, so no loaded skill *failed* to prevent it on β's side.

**§2.9 independent γ process-gap check — answer is yes, one concrete gap exists:** `gamma/SKILL.md`'s scaffold-authoring guidance (§2.5 Step 3b) names an "AC oracle approach" table as a required scaffold artifact but gives no rule for what happens when an AC's own text carries a negative/prohibition clause. That gap is real, generalizable beyond this cycle, and not something I am patching directly in this closeout — `gamma/SKILL.md` is canonical CDD doctrine and a same-session edit during closeout would ship an unreviewed change to every future cycle's scaffold convention. Instead: **filed [cnos#616](https://github.com/usurobor/cnos/issues/616)** (concrete next MCA — first action named: add a rule to §2.5 Step 3b requiring negative-invariant clauses to get their own concrete oracle check, with dispatch-boundary invariants requiring a built-binary/subprocess check rather than a hand-constructed harness).

### γ's own arithmetic slip (self-critical accounting)

γ's scaffold text claimed "9 rows total for this sub" while separately naming the row set as A1–A5, B1–B4, E1 — which sums to 10, not 9. α caught this, used the correct 10-row set, and flagged the discrepancy rather than silently dropping a row or silently correcting γ's prose. β independently confirmed the 10-row accounting was correct and the disposition non-blocking. Naming this here rather than omitting it: γ's own artifact had an arithmetic error in this cycle, caught downstream at no real cost (non-blocking, one sentence of disclosure in both `self-coherence.md` and `beta-review.md`), but it is the kind of small quality lapse worth being honest about in a retrospective that also just flagged a real oracle-precision gap in the same document.

---

## Follow-up issues

Checked existing trackers before filing anything (`gh issue list --search`) to avoid duplicating what's already tracked:

| α-disclosed debt item (`self-coherence.md` §Debt) | Disposition |
|---|---|
| 1. `cn --version` exits "Unknown command" — now actively referenced by the new `docs/guides/INSTALL-CDS.md` verification step | **Already tracked** — [cnos#612](https://github.com/usurobor/cnos/issues/612) AC1 targets exactly this. No new issue; added a cross-reference comment noting #608's docs now depend on it, raising urgency. |
| 2. `cn cdd verify`'s `checkSmallChangeArtifacts` (`ledger.go` `forUnreleased=false`) hard-fails an in-progress small-change-classified `self-coherence.md`, unlike the triadic path's lenient warn — forced α to switch this cycle's section headers from `§`-prefixed to bare form mid-authorship | **Already tracked** — [cnos#577](https://github.com/usurobor/cnos/issues/577) (classifyCycleType misclassifies triadic cycles as small-change before `beta-review.md` lands). No new issue; added a cross-reference comment naming #608 as a second empirical occurrence. |
| 3. `--index <URL>` + explicit `--release <tag>` combination untested | **Closed within the cycle** — `TestRun_IndexHTTPURL_WithExplicitReleaseTag` added in R1, closing β's R0 non-blocking observation 4. No issue needed. |
| 4. `docs/guides/README.md` not updated with a nav link | **Not a gap** — the issue's own §Docs text explicitly defers this to #611 (bootstrap delegation). No issue needed. |
| 5. No lockfile-level concurrency guard (simultaneous `cn repo install` invocations could race on `.cn/deps.json`/`.cn/deps.lock.json`/vendor-tree writes) | **Not filed.** Inherited from `internal/restore`'s pre-existing lack of any lock/mutex primitive — a whole-substrate architectural gap, not specific to #608, with no concrete evidence of an actual race (this is normally a single-operator, single-invocation onboarding command). Revisit if a real concurrent-invocation failure is ever reported; filing now would be speculative. |
| 6. Windows is not a supported target | **Not filed.** A deliberate, pre-existing platform-scope boundary (`binupdate.go` already restricts platform binaries to linux/macOS) — not new debt from this cycle, and no product decision to add Windows support exists yet. Nothing actionable to file without that prior decision. |
| 7. Transient CI red on `cdd-artifact-check` during incremental authoring | **Same root cause as item 2** — see #577 cross-reference above. |

**New issue filed:** [cnos#616](https://github.com/usurobor/cnos/issues/616) — "cdd/gamma: scaffold AC-oracle-approach must propagate negative-invariant clauses with dispatch-boundary precision (cdd-protocol-gap, cycle #608)". This is the concrete next MCA for the process-gap audit's scaffold-gap finding above (labels: `P2`, `area/cdd`, `kind/skill`, `protocol:cdd`).

**Total new issues this cycle: 1** (#616). Two candidate debts were already covered by existing open issues (#577, #612) and cross-referenced rather than duplicated; two were reasoned through and explicitly not filed (concurrency locking, Windows support); two were resolved or scoped-out within the cycle itself.

---

## Triage carryforward

For the next reader (δ, or a resumed γ):

1. **`alpha-closeout.md` and `beta-closeout.md` do not yet exist on `origin/cycle/608`.** Per `delta/SKILL.md` §9.3 step 5 and §9.5's converge-boundary artifact set, α and β each still owe their own closeout artifacts before the full converge boundary is complete. δ must dispatch both roles for their close-out prompts if not already in flight.
2. **Full `gamma/SKILL.md` §2.10 closure has not run and should not run yet.** No `RELEASE.md`, no `.cdd/unreleased/608/` → `.cdd/releases/{X.Y.Z}/608/` move, no remote-branch deletion, no hub-memory update, no issue-close assertion. Issue #608 is correctly `OPEN`. These are release-time steps that belong to a later γ invocation once the cycle-PR is reviewed and actually merged to `main` (δ opens the cycle-PR and requests `status:review` per §9.6 after this closeout and α/β's closeouts are all present).
3. **Parent wave #607** is unaffected by this cycle's completion — #609 (renderer generalization), #610 (dispatch install), #611 (bootstrap delegation) remain filing-only/operator-gated per the wave's own sequencing; nothing here changes that.
4. **cnos#616** (new, filed this closeout) is a standalone `gamma/SKILL.md` doctrine patch, unassigned — whoever picks up the next CDD-tooling cycle should consider it alongside #577 (a related but distinct classifier gap) and #375 (the earlier rule-3.11b scaffold-gate precedent this issue's shape follows).
5. **CI is green** at the converge HEAD `7bb7591461bdfbc30ee1291a192101336256bbc7` (`Build` → `success`, confirmed via `gh run list --branch cycle/608`).

---

**Cycle #608's γ process-gap audit is closed. Full cycle closure (release-time §2.10 checklist) is deferred to the merge/release step, not this converge-boundary invocation.**
