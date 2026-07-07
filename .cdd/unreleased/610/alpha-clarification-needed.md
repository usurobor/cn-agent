# α clarification note — cnos#610

Two genuine ambiguities surfaced during implementation that neither the issue body nor
`gamma-scaffold.md` fully resolves. Both are resolved with a documented best-effort
assumption below (not a stall) — flagging them here so β/δ/operator can override the
assumption if it is wrong, rather than discovering the gap silently.

## 1. Does AC3's "actionable error" mean `Run()` returns non-nil (process exits nonzero),
or a printed warning while the command still exits 0?

- AC2 says "fail early, **nonzero exit**." AC3 says "an actionable error is **returned**"
  — a different phrase, not repeating "nonzero exit."
- Mock C's console transcript (predates the AC3 label requirement) shows a fully
  "successful"-looking run with only a `⚠` warning, never a nonzero-exit indication.
- The issue's own Proof plan lists "negative: missing identity, missing labels" as
  parallel negative cases, which reads as though missing labels should also produce a
  failing (nonzero) run, like missing identity does.

**Assumption made:** `Run()` returns a non-nil error whenever `--dispatch cds` completes
its render and the cnos#493 label mechanism is (as of today) absent — i.e., every current
`--dispatch cds` invocation ends in a nonzero exit, even on an otherwise "successful"
render, until cnos#493 ships. This reading treats AC1's "after a successful run" as
describing the base-install-and-render steps succeeding, not the whole `Run()` call
returning nil. See `self-coherence.md` §ACs (AC1/AC3 relationship) for the full reasoning.

**If this is wrong:** the fix is narrow — change `runDispatchCds`'s final step in
`src/go/internal/repoinstall/repoinstall.go` to print the `ensureCanonicalDispatchLabels`
error as a `⚠` stdout warning instead of returning it, and adjust the corresponding test
assertions (`err == nil` instead of `err != nil` on the label-gap path).

## 2. Does AC5's "prose-clean" requirement override the pinned byte-identical-to-golden
backward-compat floor, when the two conflict for one specific leak?

- The δ-pinned implementation contract states: "`--agent sigma` render must remain
  byte-identical to the committed golden ... this is the hard backward-compat floor."
- `gamma-scaffold.md` itself flags that the line-296 leak (the dated
  `cn-sigma:.cn-sigma/logs/20260624.md` incident citation) "cannot be `{agent}`-templated
  honestly" — templating it would fabricate the same specific historical incident for
  every tenant, which is false for anyone but sigma.
- These two statements are in direct tension for this one leak (not the line-101 leak,
  which IS byte-identical-preserving via `{agent}` templating).

**Assumption made:** prose-clean (AC5, the issue's Scope-level operator directive) wins
for this one leak; the specific dated citation was moved to the "Responsibilities (body
reference)" appendix (not rendered into any agent's prompt) rather than kept in the
rendered body for sigma only. This makes `--agent sigma`'s render NOT byte-identical to
the **pre-cycle** golden (it IS byte-identical to this cycle's own re-render, which is
committed alongside). The issue's own AC5 text says "renders correctly (compatibility
preserved)" — not literally "byte-identical" — which is the weaker, satisfied claim.

**If this is wrong:** the fix is to keep the line-296 sentence in the rendered body
unconditionally (restoring pre-cycle byte-identity for sigma), and instead scope AC5's
positive case narrowly to "no MORE sigma leaks than were already known and accepted" —
i.e. treat cnos#610 as closing only the two leaks the issue names as its Scope's
literal example strings (`"today: sigma"`, `agent-admin-sigma`, `cn-sigma:` paths) while
accepting that the specific historical citation is a third, pre-existing, harder-to-close
leak deferred to a follow-up. I judged the honest-prose reading more consistent with the
issue's own Non-goals ("no sigma default" as a design principle) and with γ's explicit
framing of this exact leak as un-templatable, but this is δ/operator's call to ratify.
