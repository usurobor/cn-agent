# α close-out — cycle #618

**Issue:** [usurobor/cnos#618](https://github.com/usurobor/cnos/issues/618) — release-infra: FSM-capable `cn` tooling channel
**Branch:** `cycle/618` @ `5cddfbaa` (γ process-gap closeout commit; HEAD at authoring time)
**Round:** R0 (α implemented, β independently reviewed, converged on first pass — no RC round)
**Authored by:** α, wake-invoked closeout dispatch, no memory of the implementation session

## Scope note (binding on how to read this document)

This is written per `alpha/SKILL.md` §2.8's close-out mechanism, but outside the
standard re-dispatch path that section describes: that path assumes β has already
merged (`beta-closeout.md` exists, there is a merge SHA, `gamma/SKILL.md` §2.10
closure has run). At authoring time the cycle has **not** merged — `beta-review.md`
records `verdict: converge` but no PR has been opened, and `gamma-closeout.md`
(pre-existing on this branch, authored under the same pre-merge constraint) is
explicit that it is not a full §2.10 closure declaration either. Per this session's
own wake instructions (commit and push to `cycle/618`; do not touch code, spawn
roles, open a PR, or touch labels), this document is the α-side implementation
retrospective closest in spirit to §2.8's "provisional close-out fallback," authored
post-convergence rather than at review-readiness time since β's verdict is already
known and can be incorporated directly rather than guessed at.

**Voice, per §2.8:** factual observations and patterns only — this document does not
recommend dispositions or triage priority; that is γ's role and `gamma-closeout.md`
already performed it for the process-gap findings this cycle surfaced.

## 1. What was hard or ambiguous

- **AC4 was the only AC that was not a mechanical instruction-following exercise.**
  AC1–AC3 and AC2 had a scaffold-supplied mechanism (workflow_dispatch tag input,
  `CNOS_CHANNEL` selector) that only needed correct implementation. AC4 named a
  historical failure (`fatal: could not read Username ... terminal prompts
  disabled`) without pre-selecting a fix mechanism — the scaffold explicitly left
  "persist-credentials, an explicit token, job-level permissions, or something else"
  open, and flagged (friction note 4) that the failure might not even reproduce on
  a fresh attempt. This was the one place genuine engineering judgment was required
  rather than execution of a already-made decision, and it required going back to a
  prior cycle's (#429) own RCA of the same failure class rather than trusting the
  scaffold's framing alone.
- **Distinguishing "fixed" from "hardened against a plausible mechanism" for AC4.**
  The historical failure could not be reproduced live in this environment (no real
  tag push), and a prior cycle had already concluded the same failure class was
  "likely transient/GitHub-side, not a structural permissions defect." Writing an
  honest AC4 disposition meant holding two things at once: a real, non-cosmetic
  change was made (`fetch-tags: true` removes the specific raw-SHA-to-tag-ref fetch
  shape seen in the one failing log this repo has) *and* that change cannot be
  proven to be **the** fix for what may be a non-deterministic GitHub-side flake.
  The temptation in a closure-oriented artifact is to round this up to "fixed";
  the scaffold's own friction note 4 and cnos#429's prior RCA were what kept the
  self-coherence.md AC4 section from doing that.
- **The AC2 oracle command in the scaffold had a shell semantics bug** (`VAR=val
  cmd1 | cmd2` does not export `VAR` into `cmd2`), which meant the *canonical*
  verification instruction for AC2, taken literally, would not exercise the
  mechanism it was meant to test. Recognizing this required reading the oracle
  command character-by-character against POSIX pipeline semantics rather than
  executing it as given — the ambiguity was not in the implementation, it was in
  the acceptance-criterion instrument itself.

## 2. What I would do differently

*(This section is written on best-effort reconstruction from the artifacts left on
the branch — `self-coherence.md`'s Self-check and Debt sections — since this
closeout session has no first-hand memory of the implementation session. It is
scoped to what the artifacts themselves indicate as friction, not to new
second-guessing invented at closeout time.)*

- `self-coherence.md` records that AC1/AC3's dispatch-round-trip and AC2's live
  tooling-tag install were never executable in the implementation sandbox, and
  discloses this as the single largest residual item in §Debt item 1 rather than
  waiting for β to surface it. The artifact trail does not show a point where this
  ceiling was discovered late or caused rework — it was named up front, in the
  ACs' verification-method note, before any AC evidence was written. Nothing in the
  retained record suggests a different sequencing would have surfaced this earlier
  or more cheaply.
- The one place a different approach might have changed the outcome: `install.sh`'s
  header-comment fix for the AC2 oracle bug (documenting `curl ... | CNOS_CHANNEL=tooling
  sh`, env var on the piped `sh`'s side) treats the scaffold's oracle bug as
  something to route around in the shipped artifact's own documentation, rather
  than as something to flag back to γ for a scaffold correction before build started.
  Given the bug was caught during implementation rather than before it, routing
  around it in `install.sh`'s own docs (and disclosing it explicitly in both
  self-coherence.md and this closeout) was the available option; whether flagging
  it back to the scaffold author before proceeding would have been cheaper is not
  determinable from the retained artifacts alone.

## 3. Did the scaffold give enough to work with?

Yes, on the dimensions that matter for a design-and-build cycle: γ's scaffold
(`gamma-scaffold.md`) supplied a concrete per-AC oracle list, a source-of-truth
table, explicit scope guardrails (cnos#607/repo-installer, semver ranges, hosted
registry, D14-held feature version — all named as non-goals), a 7-axis
implementation contract, and six friction notes recording judgment calls already
made (tooling-channel mechanism, tag-naming convention, AC4's reproduction caveat,
AC2's accidental-pass trap). `self-coherence.md`'s own header states the scaffold's
central design decision was followed "as-is, no deviation," and β independently
re-derived the same AC1/AC3 mechanism from the diff and the real
`action-gh-release@v1` source without needing to consult the scaffold's reasoning
first — which is evidence the scaffold's decision was concrete enough to be
independently reconstructible, not merely asserted.

The one place the scaffold under-delivered was mechanical rather than substantive:
the literal shell snippet in the AC2 oracle (§1 above) had a real bug. This is
notable specifically because it originated in the scaffold-authoring role's own
artifact, not in the issue body or in the implementation — the scaffold is not
automatically immune to the class of shell pitfall an implementation is expected to
catch elsewhere.

## 4. Did β's review catch anything I missed?

`beta-review.md` records `verdict: converge` with **zero blocking findings**. β's
own §Checks I ran myself section shows it independently rebuilt its own driver
harness for AC2 rather than reusing the one described in self-coherence.md, fetched
`action-gh-release`'s real source independently rather than trusting the
self-coherence.md's characterization of it, and ran a live network probe
(`curl -fsSI .../releases/latest`) that self-coherence.md itself did not perform. On
every AC, β's independent verdict matches the self-coherence.md claim, and on every
point where β could have found a gap (the AC2 oracle bug, the AC4 transient-failure
caveat, the commit-author-identity mismatch), β's review explicitly states it
*independently re-derived* the same conclusion rather than accepting α's
self-disclosure — i.e., β corroborated rather than corrected.

β's review carries two non-blocking observations, both explicitly scoped by β
itself as not AC1–4 defects and not requiring an iteration round:

- **F1** — an operator dispatching `workflow_dispatch` with `smoke-only=false` and
  no `tag` input still hits `action-gh-release`'s own generic
  `"⚠️ GitHub Releases requires a tag"` error rather than a cnos#618-specific
  friendlier message. β itself notes this is a pre-existing rough edge outside
  AC3's oracle scope (which requires an explicit tag), not a regression introduced
  this cycle.
- **F2** — a process note, not a code finding: β held `gh` credentials sufficient to
  perform a real `workflow_dispatch` run and a real tag push against the production
  repo (which would have closed the live-verification gap named in
  self-coherence.md §Debt item 1), and chose not to, reasoning that neither the
  scaffold nor its dispatch instructions explicitly authorized a production
  side-effecting action. This surfaced a gap in the cycle's own coordination
  surfaces — no scaffold field or dispatch convention states whether such an action
  is ever authorized — that neither self-coherence.md nor the scaffold anticipated.
  `gamma-closeout.md` (already on this branch) treats this as Finding B and records
  it as a candidate cross-cycle follow-up.

## 5. Debt and residual risk — as disclosed on the branch, restated here for the
   closeout record

These are carried forward from `self-coherence.md` §Debt, independently confirmed
by β's review, and restated here without softening or upgrading their certainty:

1. **AC4's fix is disclosed as "possibly transient," not proven.** The best
   available evidence — this cycle's own log analysis of run 26340314834 plus
   cnos#429's prior RCA of the identical failure class — concludes the original
   checkout-auth failure was *likely* transient/GitHub-side rather than a static
   workflow misconfiguration. `fetch-tags: true` removes the specific
   raw-SHA-to-tag-ref fetch request shape implicated in the one failing log this
   repo has, and the explicit `persist-credentials`/`token` pins are genuine
   hardening against future silent drift — but if the failure is truly a rare
   GitHub-backend flake, no client-side YAML change can guarantee it will not recur.
   This is stated in `self-coherence.md` §Debt item 2 as a disclosed limit, not a
   proven fix, and β's independent review (item 6) reached the same conclusion by
   its own re-derivation rather than accepting the claim.
2. **The AC2 oracle command as literally written in the scaffold has a real
   env-propagation bug** (`VAR=val curl ... | sh` does not carry `VAR` past the
   pipe into the piped `sh`), documented in `self-coherence.md` §AC2 and §Debt item
   3, in `install.sh`'s own header comment, and independently re-confirmed by β
   (beta-review.md item 3's fourth bullet, from first-principles POSIX pipeline
   semantics rather than trusting the self-disclosure). This is a defect in the
   scaffold's canonical oracle text, not in the shipped diff; anyone who later
   copy-pastes the scaffold's literal oracle command should be aware it would
   silently exercise the default (stable) path rather than the tooling channel.
3. **Commit-author identity on this branch (`sigma@cnos.cn-sigma.cnos`) does not
   match the per-role `alpha@cdd.cnos` pattern named in `alpha/SKILL.md` §2.6 row
   14.** `self-coherence.md` §Debt item 6 discloses this as the harness/session
   identity for this wake-invoked run rather than a per-role identity, consistent
   across every commit on this branch (`b9ae2222`, `88a856f7`, `bd7d2c49`,
   `15cb8fd5`, `6212665a`, `a15dbf3a`, `92471fff`, `5cddfbaa`, and this closeout's
   own commit). β independently confirmed the same observation without treating it
   as a cnos#618-scoped defect. `gamma-closeout.md` (Finding C) judged this an
   agent/harness-level limitation outside this repository's control rather than a
   project-level follow-up.
4. **No live GitHub Actions execution occurred at any point in this cycle** — not
   during implementation, not during review, and not during this closeout. AC1's
   dispatch-and-inspect round trip, AC2's install against a real published tooling
   tag plus the `3.82.0` negative control, AC3's tag-targeted publish confirmation,
   and AC4's real tag-push-reaches-`publish`-green confirmation all remain
   unverified by live execution. Every AC's disposition rests on local mechanical
   checks (shellcheck, `sh -n`, `yamllint`, a Python YAML parse, driver-script runs
   against stubbed and, in β's case, live-network `curl`) plus direct reading of
   `action-gh-release@v1`'s own source and historical run logs — not on an actual
   dispatch or tag push against `usurobor/cnos`. Both roles named this ceiling
   explicitly rather than allowing it to be discovered later.

## Closing state of this session

Per this wake's explicit scope: this document is authored, committed, and pushed to
`cycle/618`. No code was touched, no other CDD role was spawned, no PR was opened,
and no labels were touched. Full cycle closure (merge, `RELEASE.md`, the
`.cdd/unreleased/618/` → `.cdd/releases/{X.Y.Z}/618/` move, branch cleanup, and the
§2.10 closure declaration) remains pending merge and is out of scope for this
session, consistent with `gamma-closeout.md`'s own scope note on this same branch.
