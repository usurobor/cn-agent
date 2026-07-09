# γ close-out — cycle #618

**Issue:** [usurobor/cnos#618](https://github.com/usurobor/cnos/issues/618) — release-infra: FSM-capable `cn` tooling channel
**Branch:** `cycle/618` @ `92471fff` (β R0 converge commit; HEAD at authoring time)
**Round:** R0 (no iteration — α implemented, β independently reviewed, converged on first pass)
**Issue state at authoring time:** `OPEN` (no PR yet; cycle not merged — see §Scope note below)

## Scope note (binding on how to read this document)

This is **not** a full `gamma/SKILL.md` §2.10 closure declaration. The cycle has not
merged (`gh pr list --head cycle/618` returns empty; `gh issue view 618 --json state`
returns `OPEN`), so `alpha-closeout.md`/`beta-closeout.md` do not exist, there is no
merge SHA to verify CI against, `RELEASE.md` has not been written, and the
`.cdd/unreleased/618/` → `.cdd/releases/{X.Y.Z}/618/` move has not happened. Per this
wake's explicit instructions, this document is scoped to the §2.7 (steps 8–9,
triage-close-outs-explicitly) and §2.9 (independent γ process-gap check) obligations
applied at the point the cycle actually is — R0 converged, pre-merge — not to the
15-point closure gate, which is δ/merge-time territory and out of scope for this
session. No code was touched, no other role was spawned, no PR was opened, no labels
were touched.

## 1. Did the scaffold's design decision hold up?

**Yes, with no drift.** The scaffold's central decision (friction note 1: implement
"tooling-channel" as a GitHub prerelease published via an explicit `workflow_dispatch`
`tag` input, not a `git push` of a matching semver tag) was adopted by α verbatim —
`self-coherence.md`'s own header states it "followed as-is, no deviation" — and
independently re-derived by β from the diff and the real `action-gh-release@v1`
source, not merely accepted on α's say-so. Neither role proposed an alternative
mechanism (e.g. a real tag-push channel, a separate release feed, a hosted index).
The corollary decision (friction note 2: because prereleases are excluded from
`/releases/latest`, `install.sh`'s default path is structurally unchanged and the
channel selector must be additive/opt-in) was verified empirically by both roles via
independently-built driver harnesses against stubbed `curl`, not merely asserted —
β rebuilt its own harness rather than reusing α's, and both converged on the same
result (default path byte-identical to pre-cycle behavior). The tag-naming convention
(friction note 5, `tooling-<UTC-date>-<short-sha>`) was adopted with the exact stated
rationale (non-collision with the `[0-9]*.[0-9]*.[0-9]*`/`v*` push-trigger globs) and
not re-litigated.

No part of the scaffold's design surface drifted during implementation or review.

## 2. Friction notes — did they materialize as real problems?

| # | Friction note | Materialized as a real problem? | Evidence |
|---|---|---|---|
| 1 | tooling-channel = prerelease via `workflow_dispatch` tag input, not tag-push | No — adopted cleanly, the separation from AC4's tag-push path worked exactly as the note predicted | self-coherence.md header; beta-review.md AC1/AC3 rows |
| 2 | default `install.sh` path structurally unchanged, opt-in only | No — confirmed byte-identical by two independently-built driver harnesses | self-coherence.md §AC2; beta-review.md item 3 |
| 3 | `1ebdb86a` (repo-installer fix) not on `main`, must not be merged here | No — confirmed absent from the diff by both roles (`git diff --stat`, path grep) | self-coherence.md Self-check; beta-review.md item 1 |
| 4 | AC4's historical auth failure may not reproduce (intermittent); α must confirm before diagnosing | **Yes, exactly as anticipated** — α confirmed the most recent real tag-push (`3.82.2`, run 28966706260) did *not* reproduce a checkout-auth failure, and root-caused the original failure against #429's own prior RCA rather than pattern-matching a plausible fix. This is the friction note doing its job: it prevented an overclaimed "we fixed the auth bug" narrative and produced instead an honest "hardening against a plausible mechanism, not a proven fix for what may be a GitHub-side flake" disposition — carried through unchanged into β's independent verdict | self-coherence.md §AC4 + §Debt item 2; beta-review.md item 6 |
| 5 | tag-naming convention `tooling-<date>-<sha>` | No — adopted as specified | self-coherence.md §AC2 |
| 6 | `install.sh` already "passes" AC2 by accident via the hand-patched `3.82.3` release; β must not accept that as evidence | **Partially — the note's warning was necessary and used.** β's own live network check (`curl -fsSI .../releases/latest` → redirects to `3.82.3`) confirms the accidental-pass condition is still live on the real repo today. β explicitly did not treat it as evidence and instead verified the *new* mechanism via its own stubbed-`curl` driver. The friction note correctly anticipated a real trap and the gate held | beta-review.md item 3, item 3's live-network sub-bullet |

Net: the scaffold's friction notes were not boilerplate — notes 4 and 6 named traps
that were live and would otherwise have invited overclaiming or a vacuous pass, and
both roles' behavior shows the notes were read and actively applied, not just filed.

## 3. New process/doctrine gaps this cycle surfaced

### Finding A — the scaffold's own AC2 oracle command has a shell env-propagation bug

The gamma-scaffold.md AC2 oracle (not the issue body — the issue body's own AC2 text
is prose-only and doesn't carry a literal shell command) reads:

```sh
CNOS_CHANNEL=tooling curl -fsSL https://raw.githubusercontent.com/.../install.sh | BIN_DIR="$tmp/bin" sh
```

`VAR=val cmd1 | cmd2` only exports `VAR` into `cmd1`'s environment (here `curl`'s),
not into `cmd2` (the piped `sh`) — so as literally written, this oracle would
silently exercise the *default* (stable) path instead of the tooling channel, and
could pass vacuously without ever touching the code this cycle added. α caught this
self-disclosed (self-coherence.md §AC2, §Debt item 3); β independently re-derived the
same conclusion from first-principles POSIX pipeline semantics rather than trusting
α's claim (beta-review.md item 3's fourth bullet). No harm reached the shipped
artifact — `install.sh`'s own header comment (added this cycle) documents the correct
invocation shape — but the bug originated in **γ's own scaffold**, the artifact this
role is directly responsible for, and nothing in the scaffold-authoring gate
(`gamma/SKILL.md` §2.5's pre-dispatch scaffold check) currently requires a literal
shell oracle snippet to be dry-run before being committed as canonical AC-oracle
text. This is a loaded-skill-miss in the §2.9 sense: a predictable, mechanical shell
pitfall slipped through because no gate checks it, and it had to be caught twice
(once by α, redundantly by β) rather than never occurring.

**Disposition:** not landed as a skill patch in this session (out of scope per this
wake's explicit "do not touch code" / narrow-scope instruction — a `gamma/SKILL.md`
edit is a legitimate γ action in the normal lifecycle but is outside what this
specific closeout session was asked to do). Recommended as a concrete next MCA below,
not silently dropped.

### Finding B — no scaffold/dispatch convention for authorizing live, production-side-effecting verification

β disclosed (beta-review.md F2) that its session held `gh` credentials sufficient to
perform a real `workflow_dispatch` run and a real tag push against
`usurobor/cnos` — closing α's own disclosed residual-evidence gap (self-coherence.md
§Debt item 1) — but chose not to, reasoning that neither the scaffold nor its
dispatch instructions explicitly authorized a real, production-side-effecting action,
only "whatever local checks are actually runnable." This was the right call in the
absence of explicit authorization, but it means AC1's and AC3's actual
"assets appear at the tag" / "release lands at the given tag" oracle claims, and
AC4's "a real tag push reaches `publish` green," remain **unverified by live
execution** at the point this cycle converged — closed only by static/source-level
reasoning (reading `action-gh-release`'s real TypeScript source, historical log
analysis, prior-cycle RCA cross-reference). Both α and β were explicit and honest
about this ceiling rather than overclaiming; the gap is structural to the sandbox,
not a quality failure by either role.

This is not specific to #618: **every future release-infra cell** that needs to
prove a `workflow_dispatch` publish or a tag-push publish actually works will hit the
same ceiling, because α/β sessions in this environment cannot dispatch live GitHub
Actions runs or push real tags against the production repo. Right now there is no
scaffold field, dispatch-prompt convention, or CDD.md/gamma-SKILL.md doctrine that
says: (a) whether a role is ever permitted to perform such an action, (b) who signs
off if so (operator? δ?), and (c) what the fallback verification bar is when live
execution is withheld (this cycle's answer — direct source-reading of the third-party
action's own code plus historical-log RCA — was strong, but was invented ad hoc by
α/β this cycle rather than following a named doctrine).

**Disposition:** recommended as a follow-up issue below (see §Follow-up
recommendation) rather than resolved in this session — this is squarely a
process/doctrine gap spanning multiple future cycles, not a one-line skill patch.

### Finding C — commit-author-identity debt (recurring, but likely out of this repo's control)

α self-disclosed (self-coherence.md §Debt item 6) that all of this cycle's commits
on `origin/cycle/618` carry the identity `sigma@cnos.cn-sigma.cnos` rather than the
per-role `alpha@cdd.cnos` pattern named in `alpha/SKILL.md` §2.6 row 14, and that
this session's tooling does not permit local git-config changes — i.e. this is a
harness/session-identity limitation, not a per-cycle authoring choice. β independently
confirmed the same observation (beta-review.md, final paragraph) without treating it
as a cnos#618-scoped defect. This is now visible across at least three consecutive
commits on this branch (`b9ae2222`, `88a856f7`, and every α/β commit since), so it is
recurring rather than a one-off — but the fix (per-role git identity provisioning)
lives in the harness/wake-invocation layer, not in this repository's code or CDD
skill files, so a cnos GH issue is likely the wrong instrument. This is judged **not**
to warrant a project-level follow-up issue in this repo; see disposition below.

## 4. Close-out triage table

| Finding | Source | Type | Disposition | Artifact / commit |
|---|---|---|---|---|
| Scaffold design decision (prerelease-via-dispatch-tag) held with no drift | this closeout, §1 | process (positive) | drop — no action needed, confirms doctrine works | — |
| Friction notes 4 and 6 correctly pre-empted real traps | this closeout, §2 | process (positive) | drop — reusable pattern, captured in §Learning below | — |
| **Finding A** — scaffold's own AC2 oracle has an env-propagation shell bug | self-coherence.md §Debt-3; beta-review.md item 3 | process/skill-gap (γ scaffold-authoring gate) | **project MCI** — recommend follow-up issue (not filed this session; see §Follow-up recommendation) | `.cdd/unreleased/618/gamma-scaffold.md` AC2 oracle |
| **Finding B** — no doctrine for authorizing live production-side-effecting verification in release-infra cells | self-coherence.md §Debt-1; beta-review.md F2 | process/doctrine gap (structural, cross-cycle) | **project MCI** — recommend follow-up issue (not filed this session; see §Follow-up recommendation) | N/A — doctrine gap, no single file |
| **Finding C** — commit-author identity is session/harness identity, not per-role | self-coherence.md §Debt-6; beta-review.md | agent/harness MCI | **agent MCI** — track via hub memory / operator note, not a project GH issue (outside this repo's control) | commits `b9ae2222`, `88a856f7`, `bd7d2c49`, `15cb8fd5`, `6212665a`, `a15dbf3a`, `92471fff` |
| β's F1 (no-tag-supplied `workflow_dispatch` UX rough edge) | beta-review.md F1 | one-off UX polish, explicitly out of AC3's scope | **drop** — β itself scoped this as non-blocking, pre-existing, and not a cnos#618 regression; revisit only if operators report real confusion | — |
| `install.sh` unauthenticated GitHub API call, no pagination | self-coherence.md §Debt-4 | one-off forward-looking debt, explicitly flagged as acceptable at current volume | **drop** — α already disclosed and scoped this correctly; not a defect against AC1–AC4 | `install.sh` |
| `docs/guides/BUILD-RELEASE.md` stale OCaml-era content not rewritten | self-coherence.md §Debt-5 | one-off, explicitly out of scope per scaffold | **drop** — scaffold explicitly said α is not obligated to fix the rest of that doc | `docs/guides/BUILD-RELEASE.md` |

## 5. Cycle-iteration trigger assessment (`CDD.md` step 10 / `gamma/SKILL.md` §2.8)

| Trigger | Fired? | Basis |
|---|---|---|
| Review churn (rounds > 2) | No | Single round — R0 converged directly, no RC round |
| Mechanical overload (ratio > 20% and total findings ≥ 10) | No | β recorded 0 blocking findings and 2 non-blocking informational notes total — far under the 10-finding floor |
| Avoidable tooling/environment failure | No — the live-dispatch ceiling (Finding B) is a **structural** sandbox constraint disclosed and worked around honestly by both roles, not a failure that blocked the cycle or that a simple guardrail would have prevented this cycle. It is nonetheless flagged as a genuine cross-cycle process gap in §3 Finding B / §Follow-up recommendation, independent of whether it counts as a "trigger" here | self-coherence.md §Debt-1; beta-review.md F2 |
| Loaded-skill miss | **Yes** — Finding A (γ's own scaffold-authoring gate does not check oracle shell-snippet correctness before commit) is exactly this trigger: a predictable, mechanical error a gate could have prevented | §3 Finding A above |

Per §2.8, the fired "loaded-skill miss" trigger requires an explicit disposition: no
patch was landed in this session (scope-restricted, see §Scope note); a concrete next
MCA is named below (§Follow-up recommendation, item 1) rather than left implicit.

## 6. Independent γ process-gap check (`gamma/SKILL.md` §2.9)

- **Did this cycle reveal a recurring friction?** Yes — two: (1) Finding A, a
  mechanical oracle-authoring pitfall that is generic to any future scaffold
  embedding a literal shell one-liner, not specific to #618; (2) Finding B, the
  live-verification ceiling that will recur in every future release-infra cell.
- **Was any gate too weak or too vague?** Yes — the pre-dispatch scaffold gate
  (`gamma/SKILL.md` §2.5, "Pre-dispatch γ scaffold check") verifies the scaffold
  *exists* and carries the required sections, but has no check on the *correctness*
  of embedded oracle shell snippets. That gap is exactly what let Finding A through.
- **Did a role skill fail to prevent a predictable error?** Yes, per Finding A —
  `git config`/pipeline env-var scoping is a well-known POSIX pitfall; nothing in the
  scaffold-authoring path checks for it mechanically today, so it was caught only by
  the diligence of α and β, not by a gate.
- **Did coordination burden show a better mechanical path?** Yes, per Finding B — β
  had to make an ad hoc judgment call about whether it was authorized to perform a
  live production side effect; a named doctrine (scaffold field or dispatch-prompt
  convention) would remove that judgment call from every future release-infra β.

Since the answer is yes on all four questions and no patch was landed in this
scope-restricted session, both are committed as concrete next MCAs (§Follow-up
recommendation) rather than left as unactioned observations.

## 7. Follow-up recommendation

**Explicit statement: yes, follow-up issues are warranted** — for Findings A and B.
This session is not filing them (per the "do not spawn other roles... recommend if
uncertain" framing of this wake's instructions, and because issue-filing/selection is
ordinarily a γ Step-1/2 activity this session was not dispatched to perform); the
recommendation is left here for the operator or a future γ selection pass to file and
prioritize.

1. **Recommended issue — "γ scaffold gate: verify oracle shell snippets before
   commit."** Scope: extend `gamma/SKILL.md` §2.5's pre-dispatch scaffold check (or
   `issue/SKILL.md`'s oracle-authoring guidance) so that any literal shell command
   embedded in a scaffold or issue body as canonical AC-oracle text is dry-run (even
   against a stub) before being committed, specifically targeting the class of bug
   found here (env vars scoped across a pipe boundary, and other classic POSIX
   pitfalls). Small — likely a doc/skill-only change, no code. Priority: low—this
   cycle's instance did no damage (self-caught, then independently re-confirmed), but
   the pattern is generic and will recur in any future scaffold with a shell oracle.

2. **Recommended issue — "Release-infra cells: doctrine for live production-side-
   effecting verification."** Scope: name, in `CDD.md` or `gamma/SKILL.md`, (a)
   whether α/β sessions are ever authorized to perform a real `workflow_dispatch` run
   or real tag push against the production repo during verification, (b) who grants
   that authorization if so (a scaffold field? an explicit operator sign-off?), and
   (c) what the required fallback verification bar is when live execution is
   withheld (this cycle's ad hoc bar — direct third-party source-reading +
   historical-log RCA cross-reference — worked well and is a candidate reusable
   pattern to name explicitly, see §Learning below). This affects every future
   release-infra cell (cnos#607's installer wave and beyond), not just #618. Priority:
   medium — no defect resulted this cycle because both roles defaulted to the
   conservative, honest choice, but the judgment call currently falls on every
   individual β session rather than being resolved once as doctrine.

Finding C (commit-author identity) is **not** recommended as a cnos project issue —
it is a harness/session-identity provisioning limitation outside this repository's
control, consistently disclosed rather than hidden, and better tracked via hub memory
or an operator-facing note than a GitHub issue against this repo's code or skills.

## Learning (`CELL-KINDS.md` §"Mandatory terminal learning section")

```yaml
learning:
  observations:
    - "Scaffold friction notes (esp. #4 AC4-reproduction caveat, #6 accidental-pass
       warning) were read and actively applied by both α and β, not just filed —
       both notes named live traps that would otherwise have produced overclaiming
       or a vacuous pass."
    - "A shell env-propagation bug (VAR=val curl | sh not exporting VAR to the piped
       sh) originated in γ's own scaffold's literal AC2 oracle, not in the issue body
       or in α's implementation — the scaffold-authoring role is not immune to the
       mechanical errors it expects α/β to catch in code."
    - "β held gh credentials sufficient to perform a real, production-side-effecting
       workflow_dispatch/tag-push verification and chose not to use them absent
       explicit authorization — a defensible default, but one currently decided ad
       hoc per-session rather than by named doctrine."
  process_deltas:
    - "Scaffold-authoring gate (gamma/SKILL.md §2.5) should verify literal shell
       oracle snippets are dry-run before commit, not just that the scaffold's
       sections exist."
    - "CDD.md or gamma/SKILL.md should name whether/when a role may perform a live,
       production-side-effecting verification action, and what the fallback
       evidence bar is when it is withheld."
  reusable_patterns:
    - "When live execution is withheld, the strongest available substitute is
       fetching and reading the actual third-party dependency's own source (here,
       softprops/action-gh-release@v1's real TypeScript) rather than assuming its
       documented behavior — both α and β did this independently and it produced
       verifiable, falsifiable evidence rather than speculation."
    - "Cross-referencing a prior cycle's own RCA of the identical failure class
       (cnos#429's #429-F1 investigation, reused directly for this cycle's AC4) is a
       cheap, high-value verification step that avoided re-deriving a conclusion
       from scratch and avoided overclaiming a fix for what may be a transient
       GitHub-side flake."
  followups:
    - issue: "recommended, not filed — γ scaffold gate: verify oracle shell snippets
              before commit (see §Follow-up recommendation item 1)"
    - issue: "recommended, not filed — release-infra cells: doctrine for live
              production-side-effecting verification (see §Follow-up recommendation
              item 2)"
  operator_burden:
    - "None this cycle — R0 converged with zero blocking findings and no operator
       decision points; both recommended follow-ups are process debt for a future
       cycle, not something that required operator intervention now."
```

## Closing state of this session

Per the explicit scope of this wake: this document is authored, committed, and
pushed to `cycle/618`. No code was touched, no other CDD role was spawned, no PR was
opened, and no labels were touched. Full cycle closure (issue-state assertion,
`alpha-closeout.md`/`beta-closeout.md`, `RELEASE.md`, cycle-directory move, branch
cleanup, and the closure declaration) remains pending merge and is out of scope for
this session.
