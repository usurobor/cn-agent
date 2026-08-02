# β Review — cnos#691

## §R0 — Independent AC-by-AC verification

**Method:** re-ran every grep/diff command from scratch against `cycle/691` (HEAD `b4f6792`) vs `origin/main` (`737e573`), read the full text of all three docs, and cross-checked #690's issue body (`gh issue view 690`) and the operator-clarification comment for the ranked model, ref topology, and #684/#688 disposition language. Did not take α's `self-coherence.md` claims at face value; every quoted grep output below was independently reproduced.

### AC1 — `MEMORY.md` rewritten to ranked model; triadic framing retired

```
$ grep -n "triadic\|α — episodic\|β — reflective\|γ — working continuity" docs/reference/runtime/MEMORY.md
(no output, exit 1 — zero hits)
```

Read the full file (`docs/reference/runtime/MEMORY.md`, 0.3.0, title "Memory in cnos — ranked model"). Confirmed independently:
- No episodic/reflective/working-continuity ontology anywhere; explicit "Supersedes" line states the v0.2.0 split is retired, not relabeled.
- `r0`/`rN` model stated in "Principle" and "The model (KISS/YAGNI)" §§1–6.
- `reads:` provenance field stated as its own point (5) and present in the entry-schema YAML block.
- "Promotion ≠ rank" stated as its own bolded point (6), with cross-link to the essay.
- Status header distinguishes doctrine-canonical-now ("ratified by #690") from implementation-not-yet-existing ("Sub 2 through Sub 5, future work") — matches the oracle's explicit instruction not to overclaim implementation status.

**Verdict: PASS.**

### AC2 — `AGENT-MEMORY-LOG-STRUCTURED.md` reconciled

Read the full diff (`git diff origin/main..HEAD -- docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md`, +5/-... net small) and the full essay body. Confirmed:
- New "Doctrine cross-reference" callout at the top explicitly names `MEMORY.md` as the concrete canonical doctrine this essay feeds, reconciles the essay's richer `inputs:`/`basis:`/`coherence:` illustrative schema against `MEMORY.md`'s minimal `reads:` field as "the same provenance requirement, not two different structural models," and scopes all of the essay's "activation logs" references to the r0-substrate *role* rather than the specific v0 convention.
- "Related analogues" section now leads with `MEMORY.md` and annotates the `AGENT-ACTIVATION-LOG-v0.md` link as historical/superseded-for-memory.
- Independently grepped every remaining "activation log" mention in the essay body (lines 30, 32, 103, 113, 166, 201, 207, 209, 281, 292, 301, 321, 470, 502) — all are either generic-role descriptions (already covered by the top-of-file scoping note) or the final Related-analogues link (which itself now carries the superseded annotation). None assert the v0 convention is still the current agent-memory mechanism after the top-of-file caveat is read alongside them.
- No contradiction found between the essay's rank vocabulary and `MEMORY.md`'s.

**Verdict: PASS.**

### AC3 — `AGENT-ACTIVATION-LOG-v0.md` narrowed with supersession framing, not deleted

```
$ git diff origin/main..HEAD -- docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md | grep "^-" | grep -v "^---"
-status: field convention
```
Only one line removed from the body/frontmatter (the old `status:` value) — everything else in the diff is additive. Confirmed the file still exists at the same path (210 lines) and still contains §0 ("Writer Locality invariant", line 30) and §0.1 ("Wake-class writer ownership", line 46) with their descriptive bodies fully intact (byte-identical — only one line changed anywhere in the whole file, and it isn't inside §0/§0.1).

Confirmed independently:
- New frontmatter `status:` states "historical convention — superseded for agent-memory purposes by cnos#690 (docs/reference/runtime/MEMORY.md); §0/§0.1 ... remain current," plus a new `superseded_for_memory_by:` list pointing at `MEMORY.md` and `cnos#690`.
- New three-paragraph prose block immediately after the title (before §0) states the box-topology model now governs agent-memory continuity, explicitly carves out §0/§0.1 as *not* superseded (live production mechanics), and frames §1-onward as accurate historical record, not current/future source of truth.
- Scanned the rest of the document (§1 through §10, all headers listed) for any residual "this is the current/only way agent memory persists" claim without the caveat attached — found none; the body describes mechanics factually without re-asserting current-model status anywhere past the top-of-doc notice.

**Verdict: PASS.**

### AC4 — box topology authoritative; #684/#688 named as prior exploration

```
$ grep -n "684\|688" docs/reference/runtime/MEMORY.md docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md
docs/reference/runtime/MEMORY.md:7:**Prior exploration, subsumed (not controlling):** #684 and its held PR #688 explored a direction-based ref scheme...
docs/reference/runtime/MEMORY.md:82:- Issue #684 / PR #688 — prior substrate exploration (direction-based ref scheme); subsumed, not controlling; closed.
```
Matches α's report; zero hits in the other two docs, which is fine — MEMORY.md is the natural single location and the oracle only requires "at least one of the three."

Diffed the ref topology block in `MEMORY.md` against #690's issue body verbatim (fetched via `gh issue view 690`):
```
MEMORY.md:
usurobor/cnos      : refs/heads/sigma/<activation-id>        Sigma-at-cnos  r0
usurobor/bumpt     : refs/heads/sigma/<activation-id>        Sigma-at-bumpt r0
usurobor/tsc       : refs/heads/sigma/<activation-id>        Sigma-at-tsc   r0
usurobor/cmp       : refs/heads/sigma/box, refs/.../cloud    two activations, one repo
usurobor/cn-sigma  : refs/heads/sigma/home                   home's own r0 box (off HEAD)
```
Identical to #690 issue body's "Ref topology" section, character-for-character. Invariants line ("orphan (no `main` ancestry), single writer, fast-forward-only, no force-push, no-delete-while-registered") also matches the issue body verbatim. #684/#688 disposition language ("prior substrate exploration, not the controlling topology... closed as subsumed") is quoted from κ's operator-clarification comment on #690, correctly attributed.

Scope discipline confirmed: MEMORY.md explicitly states "This topology is doctrine ... not yet-implemented runtime state," matching the issue's Sub 2/3/5 phasing.

**Verdict: PASS.**

### AC5 — no residual contradiction across the three docs

- Re-read all three files together. `MEMORY.md` states the ranked model + box topology + promotion-is-separate-event + explicitly marks the v0 activation-log convention superseded-for-memory; the essay cross-links to `MEMORY.md` and scopes its own activation-log mentions; the v0 convention itself carries the supersession notice while keeping §0/§0.1 live. No contradiction found among the three about what the current agent-memory model is.
- Repo-wide triadic grep, reproduced independently:
```
$ grep -rln "triadic" docs/ | sort
docs/architecture/CELL-RUNTIME.md
docs/concepts/doctrine/coherence-for-agents/COHERENCE-FOR-AGENTS.md
docs/development/board/board-data.json
docs/development/board/index.html
docs/papers/CELL-OF-CELLS.md
docs/papers/COHERENCE-SYSTEM.md
docs/papers/CTB-LANGUAGE-SPEC-v0.2-draft.md
docs/papers/CTB-v4.0.0-VISION.md
docs/papers/DECREASING-INCOHERENCE.md
docs/papers/EXECUTABLE-SKILLS.md
docs/papers/FOUNDATIONS.md
docs/papers/SHIPPING-SOFTWARE-AFTER-AI.md
docs/reference/README.md
docs/reference/ctb/README.md
docs/reference/ctb/SEMANTICS-NOTES.md
docs/reference/governance/GLOSSARY.md
docs/reference/runtime/ORCHESTRATORS.md
```
Matches α's list exactly (17 files). Independently pulled `grep -n -i -B1 -A1 "triadic"` context for every one of these 17 files and read each hit: all refer to CTB/TSC's `tri()` triadic-carrier / triadic-self-coherence concept (α pattern / β relation / γ process axes, `[L|C|R]` terms, CDD's triadic lifecycle, CELL-RUNTIME's WC/PC/CC deployment shapes) — a distinct, still-current doctrine unrelated to the retired episodic/reflective/working-continuity memory split. α's claim that these are "an unrelated TSC/CTB triadic carrier concept" is verified true, not just accepted.
- Additionally checked beyond the prescribed "triadic" grep: searched for other docs referencing `MEMORY.md` or the retired memory-class names (`episodic memory`, `reflective memory`, `working continuity`) repo-wide. Found two incidental hits, neither a genuine AC5 violation:
  - `docs/guides/MIGRATION.md` — a pre-existing, unrelated `memory/` → `threads/` migration guide referencing a *different* `MEMORY.md` concept (personal curated-facts file, not `docs/reference/runtime/MEMORY.md`); untouched by and unrelated to this cell's scope.
  - `docs/papers/AGENT-ACTIVATION-LOGS-AND-EVENTUAL-CONSISTENCY.md:45` — uses "working continuity" as an ordinary English phrase describing a sleeping activation, not as a reference to the retired triadic memory-class name; does not assert a competing memory model.
  Neither is introduced or worsened by this cell's diff; both predate it and are out of this cell's explicit scope (three named docs only). Noting as informational, not a blocking finding.
- `git diff --stat origin/main..HEAD` reproduced:
```
.cdd/unreleased/691/gamma-scaffold.md              | 122 +++++++
.cdd/unreleased/691/self-coherence.md              | 113 +++++++
docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md         |   5 +-
.../conventions/AGENT-ACTIVATION-LOG-v0.md         |  11 +-
docs/reference/runtime/MEMORY.md                   | 356 ++++-----------------
5 files changed, 304 insertions(+), 303 deletions(-)
```
Confirmed via a negated pathspec diff (`git diff --stat origin/main..HEAD -- . ':!docs/reference/runtime/MEMORY.md' ':!docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md' ':!docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md' ':!.cdd/unreleased/691'`) that returns empty — nothing outside the three docs and this cell's own `.cdd/unreleased/691/` artifacts changed. No `.go`, no `.cn-sigma/`, no CI, no FSM/label file.

**Verdict: PASS.**

## Summary

All five ACs independently verified against the diff on `cycle/691` (HEAD `b4f6792`), #690's issue body, and κ's operator-clarification comment — not against α's self-report alone. No scope break, no dangling contradiction introduced by this cell, no overclaiming of implementation status. The two incidental repo-wide hits noted under AC5 (`MIGRATION.md`, `AGENT-ACTIVATION-LOGS-AND-EVENTUAL-CONSISTENCY.md`) predate this cell and are outside its explicit three-doc scope; not a basis for iteration.

**verdict: converge**
