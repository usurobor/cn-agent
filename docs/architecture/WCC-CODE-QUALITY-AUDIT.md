# WCC 0.1 — code quality audit

**Status: all eight findings closed.** Q1–Q3 were fixed by the cell itself, in
four rented episodes whose closures are committed under
[`evidence/wcc-0.1/`](evidence/wcc-0.1/). Q4–Q8 were fixed by hand. The
findings and their measurements are kept below as written; what closed each one
is recorded in [Disposition](#disposition) at the end.

Read against `eng/go`, `eng/write-functional`, `eng/evolve` (L7), and ordinary
industry practice. Scope: the cell packages added or reshaped on
`claude/wcc-0.1` — `cellkernel`, `cellspec`, `cellfill`, `cellfills`,
`cellrun`, `cellcog`, `cellskill`, `cellwork`, `cellmethod`, `cellcheck`,
`cellinput`, `cdsissue`, `cdsdesign`, `cdsadmit`, `cdspatch`, `cdsassess`.

Not a correctness review — the adversarial suite and four review rounds cover
that. This is about **repetition, signal-to-noise, and unnecessary
complication**. Every finding below is measured, not impressionistic, and each
names the smallest change that closes it.

Baseline measurements (non-test lines):

| Package | lines | code | comment | comment % |
|---|---|---|---|---|
| cellkernel | 1250 | 796 | 326 | 26% |
| cellspec | 504 | 344 | 135 | 26% |
| cellfill | 598 | 366 | 186 | 31% |
| cellfills | 75 | 29 | 41 | **54%** |
| cellrun | 291 | 191 | 75 | 25% |
| cellcog | 538 | 242 | 262 | **48%** |
| cellwork | 598 | 347 | 216 | 36% |
| cellmethod | 229 | 120 | 90 | 39% |
| cellcheck | 354 | 186 | 148 | **41%** |
| cellinput | 90 | 38 | 46 | **51%** |
| cdsissue | 289 | 182 | 89 | 30% |
| cdsdesign | 140 | 78 | 53 | 37% |
| cdsadmit | 261 | 95 | 152 | **58%** |
| cdspatch | 269 | 134 | 120 | **44%** |
| cdsassess | 795 | 496 | 254 | 32% |

Test-to-source ratio 0.55 (9,166 test lines against 16,705 source lines).

---

## Q1 — HIGH. Two fill constructors are the same function written twice.

`cdspatch.Factory` and `cdsassess.Factory` are structurally identical: same
steps, same order, same error wrapping, differing only in the concrete
Constructed type and the expected projection role. Diffing them with comments
and blanks stripped gives 30 lines each and a difference of only type names and
two message strings.

`exactShape` is worse — `cdspatch` and `cdsassess` versions are byte-identical
apart from the label:

```go
func exactShape(decl json.RawMessage) error {
	if err := cellfill.OnlyKeys(decl, "cds.patch", "fill", "cognition"); err != nil { ... }
	if cog, ok := cellfill.Field(decl, "cognition"); ok {
		if err := cellfill.OnlyKeys(cog, "cds.patch.cognition", "provider", "model"); err != nil { ... }
	}
	return nil
}
```

`Decl` is declared identically in both packages. `ResolvedDecl` differs only in
which comment it carries.

**Why it matters beyond tidiness.** Every future fill copies this block, and a
rule added to one copy is a rule missing from the others. The `NeedsSubject`
gap already showed what an unread duplicate costs.

**Smallest change.** One helper in `cellfill` that decodes a seat declaration
whose shape is `{fill, cognition}` and validates the projection role, returning
the decoded declaration and the constructed cognition. Each fill then states
only what is genuinely its own: which cognition port it wants and what it does
with the projection. Do NOT build a generic fill framework — this is one
function replacing one duplicated block.

## Q2 — HIGH. Three implementations of "bound this output".

`cellcog.boundedBuffer`, `cellcheck.tailBuffer`, and `cellwork`'s bounded `git`
reader all solve the same problem, and two `tail` functions with different
signatures exist (`cellcog/command.go:99` takes a limit, `cellcheck:349` does
not). They differ in behaviour, and the differences are accidents rather than
decisions: `boundedBuffer` keeps the HEAD and reports truncation, `tailBuffer`
keeps the TAIL and prefixes a marker.

**Smallest change.** One internal package (or one file in `cellfill`) with a
head-keeping and a tail-keeping bounded writer and a single `tail(s, n)`.
Callers pick the policy; nobody re-implements the mechanism. Keep the
`String()`-side truncation marker, which is load-bearing.

## Q3 — MEDIUM. Comment ratio is inverted in five packages.

`cdsadmit` 58%, `cellfills` 54%, `cellinput` 51%, `cellcog` 48%, `cdspatch`
44%. `cellcog/provider_claude.go` carries a **66-line comment block** on a
single argv function.

Much of this is genuinely load-bearing — the "why", the rejected alternative,
the measured evidence. That is the house style and it should stay. But at 50%+
the ratio has passed the point where a reader can find the code, and three
distinguishable things are mixed together:

1. rationale a maintainer needs (keep, in the code);
2. incident history — "an earlier version claimed X; measurement falsified it"
   (belongs in the commit that fixed it and in the memory box, not permanently
   in the header);
3. restatement of what the next line plainly does (delete).

**Smallest change.** Per package, move class 2 out and delete class 3. Target
under 35% for the packages listed, without deleting a single "why". `eng/evolve
§3.8`: when a claim needs qualifying a third time, delete it — the same applies
to a paragraph that has been narrowed twice.

## Q4 — MEDIUM. Four functions over 80 lines, each doing several things.

| lines | function |
|---|---|
| 92 | `cdsassess.Reconcile` |
| 91 | `cellspec.Resolved.Build` |
| 82 | `cdsadmit.admit` |
| 81 | `cellkernel.RunEpisode` |

`Reconcile` is the clearest: it validates the answer's vocabulary, checks
coverage four ways, applies forced dispositions, and downgrades on an
incomplete view. Those are four separable predicates over the same two values,
and each is independently testable today only through the whole function.

`Build` does lookup, requirement checking, methodology loading, two
constructions and metadata assembly.

**Smallest change.** Extract named predicates where the extraction is obvious
and leave the orchestration flat. `RunEpisode` is the fixed protocol and its
length is the protocol's; leave it. This is not a call to hit a line count — it
is that `Reconcile`'s four rules deserve four names.

## Q5 — MEDIUM. `cellcheck.Run`'s step list is a slice of anonymous structs.

The recipe is `[]struct{ name string; argv []string }` with a stop-at-first-
failure loop, and `format` is then bolted on after the loop because it needs
different failure logic. So the recipe is really "three uniform steps plus one
special one", written as if it were four uniform steps.

**Smallest change.** Give a step a named type with its own classifier, so
`format`'s different rule is a property of that step rather than code outside
the loop. This also makes adding `cue vet` later a data change rather than a
control-flow change.

## Q6 — LOW. Test helpers duplicated across packages.

`testRepo` is written three times; `contractFor` twice; `writeJSON` twice;
`writeSkill` twice. Each is a handful of lines, but they drift — the three
`testRepo`s already differ in what they return.

**Smallest change.** One `internal/celltest` package for the git-fixture and
JSON-fixture helpers used by more than one package. Do not move
single-consumer helpers; a test helper used once belongs beside its test.

## Q7 — LOW. `Observation.Candidate` is a field nothing fills.

Documented as unfilled with a test asserting `Run` never invents one. That was
correct while no caller could supply it; `cdsassess` now materializes the
candidate and could. Either fill it or delete it — a permanently empty field in
a serialized shape teaches readers to expect it and gives them nothing.

## Q8 — LOW. Error strings restate their package prefix inconsistently.

`cellwork:`, `cellcheck:`, `cds issue:`, `fill "cds.patch":`, `cell spec:`.
Five conventions. The wrapping is already unambiguous in most paths, so several
messages carry the same word twice by the time they reach an operator.

**Smallest change.** One convention: the package prefixes at the boundary where
it is first wrapped, and inner errors do not repeat it.

---

## Not findings

- **The test-to-source ratio (0.55) is right**, not high. The mutation-witness
  discipline is why every defect in four review rounds was caught by an oracle
  that failed.
- **The dual Go/CUE authority is duplication on purpose.** It has caught real
  divergences; it stays.
- **`cellkernel`'s length** is the fixed protocol plus its validation. It reads
  as one thing because it is one thing.

## Acceptance for this work

Each finding names its own check. Additionally, and non-negotiably: the full
suite, the corpus, and the adversarial suite must all still pass, and no
mutation witness may be weakened to make a refactor land. A quality change that
costs a witness is not a quality change.

---

## Disposition

| | closed by | how |
|---|---|---|
| Q1 | the cell (`self-q1.json`) | one `cellfill.SeatDecl` helper; the two fills state only what is their own |
| Q2 | the cell (`self-q2.json`) | one `cellbound` writer with head- and tail-keeping policies, one `Tail` |
| Q3 | the cell (`self-q3a`, `self-q3b`) | incident history moved to the commits; restatement deleted; every "why" kept |
| Q4 | by hand | `Reconcile` 92→27 behind four named rules; `Build` 91→42; `admit`'s slot loop is `slotFault` |
| Q5 | by hand | `step` is a named type with its own `plan` and `classify` |
| Q6 | by hand | `internal/celltest` — `Git`, `Repo`, `Skill`; the seven copies had already drifted |
| Q7 | by hand | `Observation.Candidate` deleted; the identity it would carry is already in the digest |
| Q8 | by hand | one identity, at the package boundary; the rule is stated in `cellfill`'s header |

Measured after (non-test lines), against the same packages as the baseline:

| Package | lines | code | comment | comment % |
|---|---|---|---|---|
| cellkernel | 1250 | 796 | 326 | 26% |
| cellspec | 525 | 356 | 141 | 27% |
| cellfill | 690 | 397 | 244 | 35% |
| cellfills | 51 | 29 | 17 | 33% |
| cellrun | 291 | 191 | 75 | 26% |
| cellcog | 378 | 217 | 130 | **34%** (was 48%) |
| cellwork | 576 | 329 | 214 | 37% |
| cellmethod | 229 | 120 | 90 | 39% |
| cellcheck | 361 | 180 | 160 | **44%** (was 41%) |
| cellinput | 66 | 38 | 22 | **33%** (was 51%) |
| cdsissue | 289 | 182 | 89 | 31% |
| cdsdesign | 140 | 78 | 53 | 38% |
| cdsadmit | 182 | 104 | 62 | **34%** (was 58%) |
| cdspatch | 186 | 111 | 61 | **33%** (was 44%) |
| cdsassess | 808 | 494 | 267 | 33% |

All five packages Q3 named are now at or under 35%. **`cellcheck` went the other
way — 41% to 44%** — and it was not in Q3's list. Q5 and Q7 removed code from it
and added the reasons for both changes; the ratio is arithmetic, not new
padding. It is left as measured rather than adjusted to look like progress.

Two packages are new since the baseline: `cellbound` (105 lines, Q2's one
bounded writer) and `celltest` (73 lines, Q6's fixtures).

**A correction to this document's own measurement.** The baseline reported a
test-to-source ratio of "0.55 (9,166 test lines against 16,705 source lines)".
Those two numbers are not the same measurement: the numerator counted the cell
packages' tests and the denominator counted all of `internal/`. Measured
consistently over the cell packages alone, the ratio is 9,479 / 7,470 = **1.27**
— the mutation-witness discipline costs more than the original figure implied,
and the "not a finding" verdict on it stands for the same reason at the true
number.
