# γ Clarification — #287

**Cycle:** #287 — γ creates the cycle branch — α and β only check out `cycle/{N}`
**Branch:** `cycle/287`
**Author:** γ (`gamma@cdd.cnos`)
**Issue:** [#287](https://github.com/usurobor/cnos/issues/287)
**Trigger:** β R1 verdict at `045f1f2` raises F1 (D, scope drift) + F2 (C, status truth) predicated on the artifact-fact claim that σ's `70ff2b1` is on `origin/cycle/287` only and not on `origin/main`.

This file transfers a mechanical environment fact under `gamma/SKILL.md` §2.5 ("provide mechanical environment help"). It does not forward β's reasoning to α nor adjudicate β's findings beyond their factual premises.

---

## 1. Decision

F1 and F2 are predicated on β's local view of `origin/main`. As of γ's authoritative fetch at the timestamp below, β's premise does not match origin state. F1 (D, scope drift) and F2 (C, status truth) collapse on a fresh fetch and re-evaluation. F3 (C, identity drift) and F4 (A, §4.1 ambiguity) stand on their own and require α response.

**γ does not authorize σ's `70ff2b1` as a #287 ride-along** — because σ's `70ff2b1` is not on `cycle/287` as an unrequested commit. It is on `origin/main`, and α's rebase brought it transitively into `cycle/287`'s ancestry. There is nothing to authorize.

α may proceed to a fix-round addressing F3 + F4 only. β re-fetches and re-evaluates F1 + F2 against current `origin/main` before R2 verdict.

## 2. Mechanical environment fact

Captured by γ at **2026-04-30 00:32 UTC** in `/home/user/cnos`:

```text
$ git fetch --verbose origin main
POST git-upload-pack (305 bytes)
From http://127.0.0.1:44681/git/usurobor/cnos
 * branch            main       -> FETCH_HEAD
 = [up to date]      main       -> origin/main

$ git rev-parse origin/main
70ff2b1b80e49a30a4d7dddded49a5bd33669b32

$ git merge-base --is-ancestor 70ff2b1 origin/main && echo yes
yes

$ git branch -r --contains 70ff2b1
  origin/cycle/287
  origin/main

$ git merge-base origin/main origin/cycle/287
70ff2b1b80e49a30a4d7dddded49a5bd33669b32

$ git log --pretty="%h %an <%ae> %s" origin/main..origin/cycle/287
045f1f2 beta <beta@cdd.cnos> beta-review(287): R1 Pass 3 — Findings + Regressions + Notes; verdict REQUEST CHANGES
286c5f5 beta <beta@cdd.cnos> beta-review(287): R1 Pass 2 — §2.0 Issue Contract walk
cefb4a8 beta <beta@cdd.cnos> beta-review(287): R1 Pass 1 — verdict header + §2.0.0 Contract Integrity
6f44dbb alpha <alpha@cnos.local> self-coherence(287): §Review-readiness signal — round 1 ready for β
3336fe3 alpha <alpha@cnos.local> self-coherence(287): §Pre-review gate + §Debt update post-rebase
6a897a4 alpha <alpha@cnos.local> skill(cdd): γ creates the cycle branch — α and β only check out cycle/{N}
```

Reading:

- `origin/main` HEAD is `70ff2b1`, not `a8e67b7`.
- `70ff2b1` is reachable from `origin/main`. σ pushed it to `main`, not to `cycle/287`.
- merge-base(main, cycle/287) = `70ff2b1` — every commit up to and including σ's `70ff2b1` is shared between `main` and `cycle/287`. The cycle's diff scope (`origin/main..origin/cycle/287`) is **α's three commits + β's three commits**. σ's commit is **not in the diff scope**.
- α's §Known debt #5 ("`origin/main` advanced with `70ff2b1`") and §Pre-review gate row 1 ("base SHA: `70ff2b1` (`origin/main`)") are factually correct against current origin state.

## 3. Implication for β's R1 findings

| # | Severity | Finding (β R1 wording) | Status after γ clarification |
|---|---|---|---|
| F1 | D | σ commit `70ff2b1` lands an out-of-scope spec change on `cycle/287` | **Collapses.** σ pushed to `main`, not to `cycle/287`. The commit is part of `main`'s history, brought transitively into `cycle/287` by α's rebase per `alpha/SKILL.md` §2.6 row 1. The diff `origin/main..origin/cycle/287` does not include σ's commit. There is no scope drift on `cycle/287` to authorize. β to re-verify on re-fetch and withdraw. |
| F2 | C | α's §Known debt #5 + §Pre-review gate row 1 misidentify `70ff2b1` as `origin/main` | **Collapses.** α's claim matches current origin state. β's local origin/main was stale (at `a8e67b7`, the previous main HEAD before σ's push). β to re-verify on re-fetch and withdraw. |
| F3 | C | α git identity drifts from canonical `alpha@cdd.{project}` | **Stands.** Independent of F1/F2. γ confirms by inspection: all three α commits are signed `alpha <alpha@cnos.local>` instead of `alpha@cdd.cnos`. α addresses in fix-round (rebase commits to re-author with canonical email, or amend identity going forward — α's call). |
| F4 | A | `CDD.md` §4.1 row 2 ambiguous in isolation | **Stands.** Independent of F1/F2. β's polish suggestion (parenthetical or footnote naming γ as creator) is sound. α may ship the wording fix in the same fix-round or defer to a follow-up; α's call. |

## 4. Why this happened — and what it implies for the cycle

β's stale-origin/main view is the same `git fetch --quiet` reliability failure mode that:

- α correctly worked around at rebase time (α's pre-review gate row 1 fired and α rebased onto `70ff2b1`)
- 3.61.0 PRA §4b cycle-iteration named as the avoidable-tooling-failure trigger
- AC8 of this very issue addresses ("**`git fetch` reliability is an explicit dependency** ... on N successive empty polling iterations, do a synchronous reachability re-probe ... surface to operator on failure")

β's R1 polling loop (per `CDD.md` §1.4 β step 3 reference snippet, embedded in the diff α delivered) is built around **transition-only stdout on cycle/287's head SHA**. The loop watches `origin/cycle/287` head transitions to wake β. It does not poll `origin/main` for transitions, so when σ pushed `70ff2b1` to `main` between β's intake fetch (when main was at `a8e67b7`) and β's R1 review, β's local `origin/main` ref was never updated — β had no wake-up event for "main moved" because main is not the polled surface.

This is correct loop design (β polls the cycle branch, not main) but exposes a missing rule: **when β reads `git diff main..cycle/{N}` to compute the review surface, β should re-fetch `origin/main` synchronously immediately before computing the diff** — otherwise the review base may lag `main` by however long β's session has been alive.

This is a §Tracking gap not yet covered by AC8. It is a candidate close-out finding, not a blocker for #287's merge.

## 5. Recommended path forward

1. **β** re-fetches `origin/main` synchronously, re-runs `git rev-parse origin/main` and `git branch -r --contains 70ff2b1`, and confirms γ's mechanical fact. β appends a note to `beta-review.md` recording the re-fetch result and withdrawing F1 + F2 if the fact holds.
2. **α** addresses F3 + F4 in a fix-round per `alpha/SKILL.md` §2.7:
   - F3: configure canonical git identity going forward (`git config user.email "alpha@cdd.cnos"`); decide whether to retroactively re-author the three existing α commits (adds rebase/force-push overhead, breaks β's already-committed reviewed-SHA references) or accept the existing commits as legacy and apply the rule from the next fix-round commit forward. α's call; surface either choice in the fix-round narrative.
   - F4: ship the §4.1 row 2 wording fix in the fix-round, or defer to follow-up; α's call.
3. **β** re-reviews on R2 against the fix-round.
4. **γ** captures the §Tracking gap in §4 above (β's diff-base re-fetch rule) as a close-out finding, to be triaged in `gamma-closeout.md` and folded into either #287's CDD §Tracking text (if α elects to bundle) or a follow-up issue.

## 6. Role-separation note

This file transfers an **artifact fact** (the output of `git rev-parse` and `git branch -r --contains`, captured at γ's authoritative fetch) and the **decision** that γ does not authorize a scope expansion (because there is nothing to authorize). It does not forward β's review reasoning to α and does not author the fix. β re-evaluates F1 + F2 on its own, against current origin state. α addresses F3 + F4 on its own, with the response routed through `self-coherence.md` fix-round per the standard contract.

— γ (`gamma@cdd.cnos`) at 2026-04-30 00:32 UTC
