# beta-review.md — cnos#684 (β, R0)

## Scope of this review

Independently re-derived against `gamma-scaffold.md`'s "AC oracle list," the
β prompt embedded in that scaffold, and `gamma-clarification.md`'s two
confirmed scope calls — not trusted from `self-coherence.md`'s own AC-by-AC
narrative. Base SHA (`origin/main`, re-fetched synchronously this session):
`7856c8c` — matches the base α's self-coherence.md records post-rebase, no
drift. Head SHA (`origin/cycle/684`): `91b39dc`. Cross-checked α's claims
against the issue's own pinned body + both operator comments via
`gh issue view 684 --repo usurobor/cnos --json body,comments`, not only
against γ's paraphrase.

## Scope guardrails (hard gate) — independently verified

- `git diff origin/main...origin/cycle/684 --stat` → exactly 8 files: the
  two convention docs (`AGENT-ACTIVATION-CHANNEL-v1.md` new,
  `AGENT-ACTIVATION-LOG-v0.md` 6-line pointer edit) plus 6 files under
  `.cdd/unreleased/684/` (`CLAIM-REQUEST.yml`, `dry-run-migration-plan.md`,
  `gamma-clarification.md`, `gamma-scaffold.md`, `self-coherence.md`,
  `verify-channel-reconstruction.sh`). **No path under `.cn-sigma/**`, no
  path under `.github/workflows/*`, no path under `docs/development/board/**`**
  (each independently confirmed via `git diff ... --stat -- <path>` →
  empty output for all three).
- `git tag -l` → empty. No tag created or moved.
- `git for-each-ref refs/heads/channels` → empty. No orphan ref created.
- `git log origin/main..origin/cycle/684 --oneline` → 14 commits, all
  additive doc/artifact commits (γ scaffold, δ clarification, α's 4
  artifact commits, α's incremental self-coherence commits). No commit
  subject or content indicates a force-push-shaped operation; the migration
  plan and verification script are described, not invoked (confirmed by
  reading both files — see AC6/AC7 below).
- **Verdict: guardrail gate PASSES.** No D-severity scope violation.

## AC1 — registration schema declares repo+ref+cursor for both directions; bare-SHA cursor rejected

**MET.** `AGENT-ACTIVATION-CHANNEL-v1.md` §3.1 (L70–79) shows the shared
endpoint registry with `repo:`/`ref:` for both `foreign_to_home` and
`home_to_foreign`. §3.2 (L83–93) is a table naming, per direction, the
cursor owner and the qualified `{repo, ref, sha}` location — independently
re-read, not merely grepped. L92 states the falsifiable rule verbatim: "A
cursor entry with a bare SHA and no `repo:`/`ref:` pair fails validation."
§7's entry format (L156–157) uses the same qualified `{repo, ref, sha}`
shape for `cursor_in`/`cursor_out`, consistent with §3.2 rather than
reverting to v0's bare `<agent>@<sha>` form. No code path or example in the
doc accepts a lone SHA as a complete cursor record.

## AC2 — attach contract fetch-reader-ref / write-and-ff-writer-ref; writer never pushes other stream; reader never writes

**MET.** §4 (L104–121) states the six-step loop explicitly (fetch reader ref
→ read cursor→tip → do work → append to own writer stream → commit+ff →
record cursor), followed by two explicit ruling-out sentences: "A writer
never pushes the other stream" and "A reader never writes," each with a
one-line justification tied to specific step numbers. The per-channel
writer/reader table (L118–121) names both roles for both refs concretely
(Sigma-at-cnos vs Sigma-at-home). Matches the oracle's check exactly.

## AC3 — orphan-ref invariants enforced; non-ff or cross-writer push rejected

**MET.** §5 (L127–134) is a 6-row table. Independently mapped each row
against the issue's Proposal §4 list (orphaned start; single writer; ff-only;
no force-push; no delete-while-registered; cursor-commit-reachable) — 1:1
correspondence confirmed, no invariant dropped or merged. Each row names a
concrete enforcement mechanism (GitHub branch-protection settings for
invariants 2–5; a CI/validation check for invariants 1 and 6, since native
branch protection has no "must be orphaned" or "cursor reachable" primitive)
— concrete enough for δ/operator to configure without re-deriving from
prose, per the oracle's bar.

## AC4 — promotion boundary documented

**MET.** §6 (L138–144) contains a "Promotion boundary" section with the
required MUST-NOT binding language verbatim: "Absent promotion, a channel
entry MUST NOT silently govern current behavior." Names the promotion
targets explicitly (issue, CDD contract/design doc, decision record, code
change on `main`) and cites Kernel §2.1's no-silent-drops principle by name,
matching the oracle's "binding language, not descriptive prose" bar.

## AC5 — rename landed; scoped per the confirmed γ/δ split

**MET at the confirmed scope.** Independently re-read `gamma-clarification.md`
§1 (settled: AC5's physical component structurally conflicts with the
operator's own no-`git rm`-of-live-paths guardrail; γ's design-landed/
physical-strip split is the only reading consistent with both, confirmed by
δ, not relitigated here) and cross-checked it against the actual artifact.
`AGENT-ACTIVATION-CHANNEL-v1.md` §8 "Rename status" (L174–182) states the
split explicitly: "Design-landed (this cycle...)" naming §0/§2/§3/§4/§5/§6 as
fully specified, vs. "Physical-strip deferred (operator-executed, AC7)"
naming the forbidden `git rm` action and pointing to
`dry-run-migration-plan.md` as where it is *planned*, not executed. This is
the split claim, not a claim of full physical completion — `self-coherence.md`'s
AC5 row (L46–47) makes the same claim in the same words. Confirmed:
`.cn-sigma/logs/**` is untouched by this diff (guardrail check above), so
the "design-landed, not physically landed" framing is accurate, not
aspirational.

## AC6 — content-preservation verification specified, not executed

**MET at the confirmed scope-down.** Independently re-read
`gamma-clarification.md` §2 (settled: `delta/SKILL.md` §9.12 is standing
doctrine, not a per-cycle call; AC6 is correctly scoped to "procedure
specified, not executed"). `.cdd/unreleased/684/verify-channel-reconstruction.sh`
exists, is a complete runnable procedure (blob-SHA set-membership check via
`git rev-list --objects` + `comm -23`, no bespoke hashing needed — sound
approach), and its header (L1–24) states plainly why this cell does not run
it, citing both the sparse-checkout exclusion and §9.12 doctrine. Confirmed
independently: `git sparse-checkout list` on this checkout → `/*` then
`!/.cn-sigma`; `.cn-sigma` absent from the working tree. `bash -n
verify-channel-reconstruction.sh` → syntax OK (re-run this session,
independently of α's claimed run). Script logic never hardcodes or reads a
`.cn-sigma/` path in *this* repo — it takes checkout paths as CLI arguments
(`MAIN_CHECKOUT`/`ORPHAN_CHECKOUT`), so authoring it required no read of
`.cn-sigma/` content, consistent with the doctrine boundary. Correctly
sequenced: `dry-run-migration-plan.md` §2 "Sequencing" step 1 names this
script as the binding precondition gate before the strip step (§3.3) may
run — confirmed by reading that section directly, not taking α's citation
on faith. `self-coherence.md`'s AC6 claim (L49–50) is "procedure specified,
not executed" — matches, not an executed-evidence claim.

## AC7 — dry-run-only rewrite plan; no destructive action taken

**MET.** `.cdd/unreleased/684/dry-run-migration-plan.md`, all sections
independently re-read:
- §1 names the target glob `.cn-sigma/logs/**` and explicitly excludes
  `docs/development/board/**` with rationale (cites `328451d4`'s debounce
  and #681) — states exclusion, not inclusion, satisfying the β prompt's
  explicit check on this point.
- §3.1 names `git filter-repo --analyze` as the operator's mechanism for
  obtaining the actual commit count / commits-that-become-empty at
  execution time, deliberately not hardcoding the operator's approximate
  ~1024 figure as authoritative.
- §3.2/§3.3 give the `filter-repo`-style command plan for relocate and
  strip, both explicitly framed as running only against disposable clones.
- §4 covers all three named blast-radius items — tags `3.82.x` (§4.1,
  remap-vs-preserve framed as an explicit operator decision point), open
  PRs (§4.2, snapshot + commit-map rebase + PR-freeze recommendation), and
  #682's CDD-by-SHA custody (§4.3, named as a hard sequencing dependency,
  explicitly not modifying #682's own scope) — each with a stated
  remap/re-anchor strategy.
- §5 names the rollback ref `refs/heads/backup/main-pre-cnos684-rewrite`
  with naming rationale.
- §6 restates the operator-executed, separately-gated framing and gives a
  self-verifiable claim, independently re-confirmed this session: `git tag
  -l` empty, `git diff origin/main...origin/cycle/684 --stat` touches only
  doc/plan files, no `.cn-sigma/**`, no ref/tag mutation.

## WRITER_LOCALITY_VIOLATION (negative case)

**MET.** §3.1's shared registry block (L70–79) independently re-read: no
`cursor:` field on either the `foreign_to_home` or `home_to_foreign` entry,
matching comment 1's amended YAML shape exactly (cross-checked against the
issue's actual comment body via `gh issue view`, not γ's paraphrase alone).
§3.3 (L94–98) names the trigger condition concretely (a side attempting to
write the other side's read-cursor record, with two concrete examples) and
the required response (reject the write, surface `WRITER_LOCALITY_VIOLATION`
as a named error/rejection state).

## `AGENT-ACTIVATION-LOG-v0.md` — not gutted, superseded-by pointer added

**MET.** File still exists at 206 lines, full v0 content intact (§0 through
§10 all present, independently read in full). Only change is a 6-line diff:
new frontmatter field `superseded_by:` plus a blockquote note (L16–17)
directly under the frontmatter pointing to v1, stating the file "remains the
historical record" and is "not deleted or gutted." `AGENT-ACTIVATION-CHANNEL-v1.md`'s
own frontmatter (`supersedes:` field, L6) and prose (§0, §10 "What changed
from v0") reference v0 without duplicating or silently conflicting with it —
v0's §0/§0.1 (Writer Locality, wake-class ownership) are explicitly named in
v1 §10 as "unchanged," not restated as if newly invented.

## Cross-check against the raw issue (not γ's paraphrase)

Read the issue body and both operator comments directly via `gh issue view`.
Confirmed: AC1–AC5 as pinned in the issue body match γ's oracle-table
restatement; the cursor-ownership amendment (comment
@2026-08-01T06:59:04Z) and the AC6/AC7 + guardrail comment
(@2026-08-01T07:50:17Z) match γ's source-of-truth table and the scaffold's
own quoted excerpts verbatim — no paraphrase drift found between the raw
issue text and what γ's scaffold (and α's artifacts) treat as binding.

## Findings

None. Every AC oracle check, the WRITER_LOCALITY_VIOLATION negative case,
the AC5/AC6 scope-split language, and every scope guardrail were
independently re-derived against the actual committed files and the raw
issue/comment text — all pass. `self-coherence.md`'s claims match the
artifacts it cites; no doc-vs-code (or doc-vs-doc) drift found.

## Verdict

`verdict: converge`

All 8 oracle rows (AC1–AC7 + WRITER_LOCALITY_VIOLATION) are satisfied by the
actual committed artifacts, independently re-verified against file content,
not self-coherence.md's narrative. The AC5/AC6 scope-splits match the
γ-scaffolded and δ-confirmed interpretation exactly, stated as splits/
scope-downs rather than claims of full completion. The scope guardrails
(no `.cn-sigma/**`, no `.github/workflows/*`, no `docs/development/board/**`,
no ref/tag mutation) hold with zero exceptions — the hard gate passes
cleanly. `AGENT-ACTIVATION-LOG-v0.md` is preserved, not gutted. Recommend
proceeding to merge.
