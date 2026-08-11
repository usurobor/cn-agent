# γ closeout — cnos#706

## Process-gap audit

Cycle ran γ scaffold → α R0 → β R0 with no fix round; verdict converged on the first review pass. Friction worth carrying forward:

1. **The "two independent sigma bindings" friction note (scaffold §Friction 4) was real and correctly caught.** The scaffold explicitly flagged that `cn-install-wake` carries two separate sigma-only call sites — the cosmetic `agent_bot_name()`/`agent_bot_id()` identity table (AC9's target) and a distinct sigma-only default for the PAT secret name itself (AC6's target) — roughly 550 lines apart in the same file, either of which could be fixed while missing the other. α fixed both; β independently confirmed both were gone from every named surface (renderer, `repoinstall.go`, `cmd_repo_install.go`, live workflow, both golden fixtures) rather than stopping at the first hit. This is a good example of a scaffold friction note doing its job: naming a plausible half-fix before it happened, not after. Worth keeping this pattern (explicit "N distinct bindings, don't fix one and miss the other" notes) in future scaffolds whenever a rename/deletion touches more than one call site in the same file.
2. **The per-AC oracle list held up well against independent re-derivation.** β re-derived every AC's oracle from the issue's own consolidated-spec comment rather than the scaffold's paraphrase, and arrived at the same pass/fail reading on all 10 ACs with no oracle-wording disputes. No refinement needed here for future cycles — the oracle list format (mechanical grep/test-name pairs, not prose) is working.
3. **Cycle-scope-sizing check (five-factor, at-edge 8–10 AC band) correctly kept the cycle whole.** Two factors fired (cross-module breadth, independent shippability) but the "keep whole" decision held up in practice: the single coherent diff (preflight + rename + bot-deletion + docs) reviewed and converged in one round with no cross-PR merge-conflict risk. No evidence this cycle that splitting would have gone better; the sizing heuristic's judgment call reads as correct in hindsight.
4. **One process gap, not a defect:** the scaffold's own §Friction note 5 (mode is `design-and-build`, not MCA, because no `DESIGN.md` exists at a stable `docs/{tier}/{bundle}/{X.Y.Z}/DESIGN.md` path — the converged design lives only in an issue-comment thread) is a recurring shape worth flagging generally: converged in-thread designs are common and this scaffold handled the mode-labeling correctly, but a future cycle touching this same governing gap (dispatch-workflow secret/identity model) would benefit from the design being filed at a stable path rather than re-read from an issue thread each time. Not urgent — no action needed this cycle — but noted as a pattern.

No process-gap here rose to the level of an immediate MCA patch or a project MCI filing; all four points above are observations for future scaffold/dispatch authors, not defects requiring a skill/spec change this cycle.

## Required pre-merge operator action

**Before merging the PR for cnos#706, the operator MUST confirm the GitHub Actions repository secret `CN_DISPATCH_PAT` exists on `usurobor/cnos`** — holding the same PAT value that `SIGMA_WORKFLOW_PAT` currently holds, or a fresh equivalent PAT with the same scopes (Contents + Issues + Pull requests + Workflows = write).

**Why this is not optional:** this cycle's diff renames the secret every live workflow references — `.github/workflows/cnos-cds-dispatch.yml` (checkout `token:`, the mechanical-recovery-scanner's `GH_TOKEN`, the `claude-code-action`'s `github_token:`, the finalizer's `GH_TOKEN`) and `cnos-agent-admin.golden.yml` — from `secrets.SIGMA_WORKFLOW_PAT` to `secrets.CN_DISPATCH_PAT`. **GitHub evaluates a reference to a nonexistent secret as an empty string, not a parse error.** If `CN_DISPATCH_PAT` is not provisioned at the moment this merges to `main`, the next scheduled firing of the live `cnos-cds-dispatch.yml` and `cnos-agent-admin.yml` wakes will run with an empty token in place of a real one, degrading or breaking checkout auth, the recovery scanner, the `claude-code-action` step, and the finalizer — see `beta-review.md` §R0 Finding F1 for the full trace, including β's confirmation that the current `main`-branch workflow (still on the old secret name) is firing normally today, which is exactly the healthy state at risk.

**Two options, name one before merging:**
- **(a)** Create `CN_DISPATCH_PAT` as a new repository secret with the same value `SIGMA_WORKFLOW_PAT` currently holds (GitHub secrets cannot be renamed in place — a new secret must be created), *before* the merge lands.
- **(b)** Sequence the merge so there is no window where the renamed workflow YAML is live on `main` but the new secret is not yet provisioned — e.g., provision the secret and merge in immediate succession with no scheduled firing able to land in between.

**This is not a code-level blocker and could not have been closed by α's diff.** Per AC3 and this cycle's scope guardrails, secret *values* are never handled, collected, or provisioned by this cycle's code by design — presence is checked by name only, never by value. Provisioning the actual secret is infrastructure/operator action outside version control and outside this diff's reach, full stop. β could not confirm secret presence itself (`gh api repos/usurobor/cnos/actions/secrets` returned `403 Resource not accessible by personal access token`) — this is a genuine unknown that only the operator/δ can close.

## Deferred / follow-up

- **Remove `SIGMA_WORKFLOW_PAT` from the repo's actual secrets** once `CN_DISPATCH_PAT` is confirmed working end-to-end (i.e., after at least one successful post-merge scheduled firing of both wakes under the new name). Small cleanup, not blocking — the old secret becomes dead weight once the rename is confirmed live, but leaving it in place briefly is harmless and arguably a useful rollback fallback during the transition window.
- **F2 (LOW, non-blocking, per `beta-review.md`):** `docs/guides/INSTALL-CDS.md`'s Tier 3 runbook rotation-cadence guidance doesn't call out that rotating `CN_DISPATCH_PAT` on the cnos repo itself carries the same secret-must-exist-before-workflow-references-it care this cycle's F1 describes — currently only implied, not stated. Doc nit; can be picked up in a future docs pass, does not block this cycle's merge.

## Deliverable evidence (δ, cnos#524 closeout-integrity preflight)

```
deliverable_evidence:
  pr: "#708 (cycle/706 -> main)"
  head_sha: "d34d7643252391fcac17da48ccc6aa35eec7e5bc"
  base_sha: "7f249ddbb50f230d5d41287b6554ab17b5a1d1d5"
  commits_beyond_base: 15
  closeout_artifacts: [gamma-scaffold.md, self-coherence.md, beta-review.md, alpha-closeout.md, beta-closeout.md, gamma-closeout.md]
```

δ confirmed via `cn issues fsm evaluate --issue 706` (read-only) before requesting the transition: `pr_exists: true`, `commits_beyond_base: 15`, `review_request_present: true`, all six required closeout artifacts present. PR #708 is open, non-draft, and its description carries the F1 pre-merge operator-action warning verbatim.
