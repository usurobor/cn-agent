# β close-out — cnos#630

## Scope note

Written at the δ §9.5 **converge boundary**, alongside `alpha-closeout.md`
and `gamma-closeout.md`. This is the review-side retrospective the §9.5
converge row requires, not the full post-merge β close-out
`beta/SKILL.md` §"Close-out" describes (that flow's `[Review Summary,
Implementation Assessment, Technical Review, Process Observations,
Release Notes]` shape presumes β has already merged and carried the
release boundary through tag/deploy — that has not happened yet on this
cycle; δ still owns opening the cycle-PR and requesting
`status:review`, per §9.6). No merge, tag, or release action is claimed
or implied by this file.

## Review approach

R0 was a single independent-review pass against α's implementation, not
a re-statement of `self-coherence.md`. The governing rule stated in
`beta-review.md` §R0: "I did not take α's self-coherence.md claims on
trust. Every AC below was re-derived from the diff and re-run from the
shell, independently of the narrative." Concretely, that meant:

- Reading `transitions.json`'s new rule (lines 184–192) directly rather
  than via `self-coherence.md`'s paraphrase, and independently
  positioning it relative to the surrounding `run_active` and
  `propose_delta_recovery` rules.
- Independently diffing `TestAC630_WedgePreFixRuleReproducesStrand`'s
  inline pre-diff `Table` literal against `git diff origin/main...
  cycle/630`'s own removed lines, to confirm the "red" side of AC2's
  red/green pair is a faithful, unmodified copy of the old rule shape —
  not a strawman that would trivially "fail" regardless of what the old
  code actually did.
- Reading both `terminal_test.go`-class test bodies for the redirected
  and newly added tests in full (not sampling by name), including
  confirming `TestAC630_TodoWithExistingBranchAndPRResumesNotDeferred`'s
  assertions match its name.
- Re-running every gate α claimed to have run locally (`go build`/`go
  vet`/`go test -v`/`go test -race`/`cn cdd verify --unreleased`) from a
  clean shell rather than accepting the reported pass/fail counts — and
  independently reproducing the same numbers (77 `--- PASS`, 0 `---
  FAIL` for the package suite).
- Going one step further than α could from its sandbox: α disclosed as
  debt that "branch CI is not independently confirmed green on GitHub
  Actions for this head commit" (no Actions runner available in α's
  environment). β does have `gh` access and used it —
  `gh api repos/usurobor/cnos/commits/a63bc9d2.../check-runs` against
  the actual review head SHA — and confirmed all 11 named checks
  `success`.
- Independently rebuilding the `cn-install-wake` renderer from source
  and re-running it against the checked-out `cycle/630` tree to confirm
  `.github/workflows/cnos-cds-dispatch.yml` reproduces byte-for-byte,
  rather than trusting the sha256 comparison α's `self-coherence.md`
  performed between two already-committed files (a weaker check — two
  files matching each other says nothing about whether either was
  hand-edited in a way the renderer would also happen to reproduce).
- Independently re-diffing `src/packages/cnos.core/labels.json` between
  `origin/main` and `cycle/630` to confirm byte-identity (no new
  `status:*` label), rather than accepting α's statement that the file
  was untouched.

## Verdict confirmation

**`verdict: converge`** (`beta-review.md` §Verdict), unchanged from R0.
All six ACs pass on independent re-derivation; no correctness, scope, or
guardrail finding survived review (`beta-review.md` §Findings: "None
blocking"). Two items are flagged as explicit non-blocking observations,
not findings:

- The §9.11 doctrine addition has no live wake-invoked-δ empirical
  witness yet — self-disclosed by α in `self-coherence.md` §Debt,
  independently re-confirmed by β as accurate and adequately scoped,
  mirroring §9.10's own first-use history (cycle/497 was also a
  bootstrap-exception witness). A legitimate scope boundary a single
  cycle cannot manufacture, not a gap this cycle should have closed.
- α's self-disclosed "branch CI is not independently confirmed green
  on GitHub Actions" gap was closed by β's independent `gh api` pull
  (see "Review approach" above) — noted for γ's closeout record, not as
  a finding requiring action, since the underlying claim (CI is green)
  turned out to be true.

This is a first-pass (`run_class: first_pass`) converge at R0 — no
`iterate` round occurred, no RC was returned, and no re-dispatch of α
was needed.

## Review-process learnings worth carrying forward

- **Independently re-running `go test ./...` from a clean shell, rather
  than trusting α's pasted output, is the single highest-leverage
  review step available whenever the numbers are re-derivable
  mechanically.** Nothing was wrong here, but the step is what
  establishes that rather than assuming α's paste is accurate — and it
  is the same discipline that caught the redirected-test rationale gap
  class in prior cycles (per the general β doctrine, not specific to
  this cycle).
- **Live CI re-verification via `gh api`/`gh run view`, independent of
  a sandboxed α's local-only equivalents, is the correct compensating
  step whenever α's own environment lacks an Actions runner.** α was
  explicit and honest about the exact gates it could not run and
  declared them as disclosed debt rather than claiming green by
  omission — that disclosure is what made it straightforward for β to
  know precisely which claim needed independent closure. This is the
  same pattern `.cdd/unreleased/615/beta-closeout.md` already named as
  the highest-leverage step in that cycle; cnos#630 is a second,
  independent confirmation of the same pattern, which is worth noting
  explicitly as a *recurring* observation now (two occurrences), not
  just a one-off. If a third cycle repeats this same shape, that is
  the point at which `beta/SKILL.md` should consider naming "pull live
  CI via `gh api` for any environment-disclosed CI gap" as an explicit
  binding step in the pre-merge gate, rather than an emergent good
  practice independently rediscovered each time.
- **Independently reproducing a rendered-artifact claim from source
  (rebuilding the renderer and re-running it), rather than accepting a
  sha256 comparison between two already-committed files, closes a real
  verification gap a naive "files match" check would miss** — two
  files can be byte-identical to each other while both having been
  hand-edited identically, or the renderer itself could be silently
  broken in a way that happens to reproduce the committed output. This
  cycle's β pass is the first time this stronger check was performed
  (as far as `.cdd/unreleased/615/beta-closeout.md`'s own review-
  process-learnings section did not name it); worth carrying forward
  as a concrete technique for any future cycle touching a
  renderer-produced artifact (`cn-install-wake` or otherwise).
- No process friction from the review side is being carried forward as
  a follow-up beyond the two observations above — this was a clean R0
  converge with no RC, no iteration, and no gate ambiguity that needed
  γ/δ escalation.
