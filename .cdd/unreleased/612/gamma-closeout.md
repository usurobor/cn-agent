# γ close-out — cycle #612

**Scope of this artifact.** Process-gap-audit retrospective at the β-converge boundary, written before the cycle-PR is opened. Issue #612 remains `OPEN` at this point. Full release-time closure (`.cdd/unreleased/612/` → `.cdd/releases/{X.Y.Z}/612/`, `RELEASE.md`, branch deletion) is deferred to the actual merge/release step.

---

## Cycle summary

**Issue:** [cnos#612](https://github.com/usurobor/cnos/issues/612) — CLI ergonomics for onboarding, per `docs/development/design/cn-repo-install-MOCKS.md` §Mock F. Sub-issue of the #607 wave; customer evidence from cnos#606 C3 (`tsc`, the first external tenant) and escalated by the operator as release-blocking against #618's next public binary cut.

**What shipped:** `cn --version`/`-V`; global and per-command `--help`/`-h` (including four commands whose rich `Help()` text existed but was previously dead code, never wired to any dispatch path); `cn init` refusing unrecognized flags instead of scaffolding a stray directory from one; `cn help` and `cn status` (the latter found by β, not named in the issue) displaying the space-separated form that actually invokes each command instead of a non-dispatchable internal registry key. All four ACs closed with `mock_parity` evidence (2 match, 2 justified-exceed, 0 missed). An end-to-end CI smoke-test step was added against the real built binary, not just unit tests.

**Review arc:** R0 converge, with one blocking finding (a second instance of the AC4 bug class, in `cn status`) found and fixed within the same round. No iteration back to α was needed.

---

## Process-gap audit

**Was γ's own scaffold complete?** Mostly — the "AC oracle approach" table and "Surfaces touched" table correctly anticipated three of the four fix locations (`main.go`, `dispatch.go`/`cmd_help.go`, `cmd_init.go`) and correctly flagged the dead-`Help()`-method discovery pattern as a likely finding (the scaffold's β prompt explicitly told β to "grep for other call sites" printing `spec.Name` directly — which is exactly how the `cmd_status.go` gap was found). The scaffold did not itself name `cmd_status.go` as a surface, which is the one gap: AC4's own issue text scopes itself to `cn help` only, and the scaffold inherited that scope without independently auditing for other `spec.Name`-printing call sites before dispatch. This is a smaller, less structurally significant version of the same failure mode #608's γ closeout (cnos#616) already named and filed a doctrine fix for — negative/adjacent-surface clauses in an AC's own text don't automatically get an independent oracle unless someone thinks to look. Not filing a new issue for this instance: cnos#616 (open, from #608's closeout) already covers exactly this class of gap in `gamma/SKILL.md`'s scaffold-authoring guidance, and this cycle's own β prompt already demonstrated the mitigating practice (explicitly telling β to grep for sibling surfaces) that caught it here at zero cost (same review round, no iteration). Filing a second issue for the same doctrine gap would be duplicative.

**Cycle-iteration trigger check (`gamma/SKILL.md` §2.8):** review rounds = 1 (not > 2); total findings = 1 blocking + 3 non-blocking (below the ≥10 mechanical-overload floor, and the blocking finding was substantive, not mechanical); no avoidable tooling/environment failure. No trigger fires.

---

## Follow-up issues

No new issues filed. The one process-gap identified above is already covered by the open cnos#616 (filed during #608's closeout); no new debt items surfaced beyond the one non-blocking CI-grep-specificity observation in `beta-review.md`, which was judged not worth a separate tracker (see `beta-closeout.md`).

---

## Triage carryforward

1. `alpha-closeout.md` and `beta-closeout.md` are both present on this branch alongside this file.
2. Full release-time closure has not run — issue #612 stays `OPEN` until the cycle-PR merges.
3. δ still owes: `REVIEW-REQUEST.yml`, opening/confirming the PR (mechanical finalizer, `cn cell finalize`), and requesting `status:in-progress → status:review` via `cn issues fsm evaluate --issue 612 --apply` once the deliverable-evidence preflight (cds-dispatch/SKILL.md §"Closeout integrity preflight", cnos#524) is satisfied.
4. Parent wave #607 is unaffected; #609-#611 (renderer generalization, dispatch install, bootstrap delegation) remain their own sub-issues, untouched by this cycle.

---

**Cycle #612's γ process-gap audit is closed. Full cycle closure is deferred to the merge/release step.**
