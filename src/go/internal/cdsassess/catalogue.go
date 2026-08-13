package cdsassess

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cdsissue"
	"github.com/usurobor/cnos/src/go/internal/cellcheck"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

// The two checker units. They are named constants because three things must
// agree on the spelling — the catalogue, the forcing rules, and the prompt —
// and a literal repeated in three places is three chances to disagree about
// which obligation was disposed of.
const (
	UnitMatterNonEmpty = "check:matter-nonempty"
	UnitProjectVerify  = "check:project-verify"
)

// UnitKind separates the obligations a SEAT decides from the ones the RUNTIME
// decides. It is not decoration: the difference is the whole of who is allowed
// to answer, and Reconcile enforces it.
type UnitKind string

const (
	KindAcceptance UnitKind = "acceptance"
	KindCheck      UnitKind = "check"
)

// Unit is one catalogue entry.
//
// Statement and Verification are the criterion's own words, carried on the
// unit so the prompt renders the catalogue from ONE value rather than walking
// the issue again beside it — two renderings of one list are two lists.
//
// Forced is the disposition the runtime already decided, or nil when the unit
// is the seat's to decide. It is on the unit rather than in a second map so
// that "which units are the seat's" cannot be answered differently by two
// pieces of code.
type Unit struct {
	ID           string
	Kind         UnitKind
	Statement    string
	Verification string
	Forced       *Forced
}

// Forced is a runtime-decided disposition and the runtime's reason for it. The
// reason travels with the disposition because a forced unit's reason is not
// the seat's to write either: the seat did not decide it, so it cannot explain
// it.
type Forced struct {
	Disposition cellkernel.Disposition
	Reason      string
}

// Catalogue is the ordered, closed list of obligations one episode is assessed
// against. Order is part of it: Reconcile requires the answer to arrive in
// this order, so "the catalogue" is one sequence and not a set that two readers
// could enumerate differently.
type Catalogue struct {
	Units []Unit
}

// Build derives the catalogue from the ADMITTED issue and nothing else:
//
//	one unit per issue.acceptance[].id, in the issue's order
//	check:matter-nonempty
//	check:project-verify
//
// NOT ONE UNIT PER SKILL. A skill is a body of guidance, not a decidable
// obligation: a seat asked "does this patch satisfy eng/go" has no criterion
// and no verification route, so it must either guess — which is the failure
// this whole cycle exists to end — or return `unverified` for every skill unit,
// which makes coverage noise out of the thing coverage is for. That is not
// hypothetical; a skill-ref catalogue on the donor branch is what produced a
// confident false finding. An acceptance criterion carries its verification
// route BY CONSTRUCTION, because cdsissue.Admit refuses one that does not. The
// skills stay in the methodology, where they shape HOW a unit is judged rather
// than pretending to be units.
//
// Pure: an issue in, a catalogue out. The forced dispositions arrive later,
// from Decide, because they depend on what the episode measured and the
// catalogue's SHAPE must not.
func Build(iss cdsissue.Issue) Catalogue {
	units := make([]Unit, 0, len(iss.Acceptance)+2)
	for _, c := range iss.Acceptance {
		units = append(units, Unit{
			ID:           c.ID,
			Kind:         KindAcceptance,
			Statement:    c.Statement,
			Verification: c.Verification,
		})
	}
	return Catalogue{Units: append(units,
		Unit{ID: UnitMatterNonEmpty, Kind: KindCheck},
		Unit{ID: UnitProjectVerify, Kind: KindCheck},
	)}
}

// Decide returns the catalogue with its two check units forced from what the
// runtime measured. Pure, and total: an observation that never ran is still an
// observation, so every path through the seat produces the same shape.
//
// The mapping is fixed and the seat has no say in it:
//
//	the matter carries a reviewable change  → pass, else finding
//	checker pass / fail / unavailable       → pass / finding / unverified
//
// `unavailable → unverified` is the one that matters. A recipe that could not
// run proves nothing about the candidate, so calling it a finding would blame
// a candidate for a missing toolchain, and calling it a pass would report
// evidence nobody has.
func (c Catalogue) Decide(m cellkernel.Matter, obs cellcheck.Observation) Catalogue {
	out := Catalogue{Units: make([]Unit, 0, len(c.Units))}
	for _, u := range c.Units {
		switch u.ID {
		case UnitMatterNonEmpty:
			if fault := MatterFault(m); fault != "" {
				u.Forced = &Forced{Disposition: cellkernel.DispositionFinding, Reason: fault}
			} else {
				u.Forced = &Forced{
					Disposition: cellkernel.DispositionPass,
					Reason:      fmt.Sprintf("the matter carries a unified diff of %d bytes", len(m.Data)),
				}
			}
		case UnitProjectVerify:
			u.Forced = checkerForced(obs)
		}
		out.Units = append(out.Units, u)
	}
	return out
}

// checkerForced maps one observation onto a disposition, naming the recipe and
// the step that decided it. The step is in the reason because "the tests
// failed" and "it did not compile" are different next work, and a reason that
// said only "project-verify failed" would make a reader open the closure to
// find out which.
func checkerForced(obs cellcheck.Observation) *Forced {
	var where string
	for _, s := range obs.Steps {
		if s.Status != cellcheck.Pass {
			where = fmt.Sprintf(" at step %q (exit %d): %s", s.Name, s.Exit, strings.TrimSpace(s.Tail))
			break
		}
	}
	switch obs.Status {
	case cellcheck.Pass:
		return &Forced{
			Disposition: cellkernel.DispositionPass,
			Reason:      fmt.Sprintf("%s passed every step against the reconstructed candidate", obs.Recipe),
		}
	case cellcheck.Fail:
		return &Forced{
			Disposition: cellkernel.DispositionFinding,
			Reason:      fmt.Sprintf("%s failed against the reconstructed candidate%s", obs.Recipe, where),
		}
	default:
		return &Forced{
			Disposition: cellkernel.DispositionUnverified,
			Reason: fmt.Sprintf("%s could not be run against the candidate%s, so nothing was measured",
				obs.Recipe, where),
		}
	}
}

// MatterFault reports why the matter carries nothing reviewable, or "" if it
// does.
//
// The unit is called `check:matter-nonempty` and this predicate is very
// slightly wider than that name: a non-blank matter that is not a unified diff
// also fails it. An id is a stable label rather than a specification, and the
// two cases are the same fact — `cds.patch` measures its product as a diff, so
// a matter with no `diff --git ` file header is a change by no spelling this
// profile can read. The REASON says which rule fired, so nothing is hidden
// behind the label.
//
// It is structural, not a diff parser: the seat reads the diff, this only
// decides whether there is one to read. The alternative — sending everything to
// a provider and trusting it to notice — is what produced a passing verdict on
// the sentence "no change was made to …".
func MatterFault(m cellkernel.Matter) string {
	switch {
	case strings.TrimSpace(m.Data) == "":
		return "the matter is empty, so there is no change to judge"
	case !strings.Contains("\n"+m.Data, "\ndiff --git "):
		return "the matter carries no unified diff (no `diff --git ` file header), so there is no change to judge"
	}
	return ""
}

// Assessment is the answer shape: exactly one disposition per catalogue unit.
// It is what the provider is constrained to by AnswerSchema and what the
// deterministic judge returns, so both arrive at Reconcile as the same value.
type Assessment struct {
	Units []cellkernel.UnitResult `json:"units"`
}

// Reconcile is the RUNTIME'S half of the verdict: it checks the answer against
// the catalogue and returns what the record carries, or a fault.
//
// A fault is not a failing review. Every condition below means the seat did not
// review — it malfunctioned, or it answered a question nobody asked — and the
// only honest thing to do with such an answer is to refuse it. Turning any of
// them into `finding` would put a judgement in the record that nobody made.
//
// viewComplete is whether the seat saw the whole candidate. When it did not,
// an acceptance `pass` is DOWNGRADED to `unverified` rather than trusted: on a
// partial view a pass cannot be told apart from a pass on evidence that was
// never shown. The seat is told this in its prompt; the downgrade is what makes
// the instruction more than a request.
func Reconcile(c Catalogue, a Assessment, viewComplete bool) ([]cellkernel.UnitResult, error) {
	answered := make(map[string]cellkernel.UnitResult, len(a.Units))
	for _, u := range a.Units {
		if _, dup := answered[u.Unit]; dup {
			return nil, fmt.Errorf("the assessment disposes of unit %q more than once; one obligation has one disposition", u.Unit)
		}
		answered[u.Unit] = u
	}
	// The vocabulary and the reason rule belong HERE, where the seat's answer is
	// judged, not at the kernel's seal. Both are things this seat got wrong, and
	// a rule enforced only downstream surfaces as an episode malfunction naming
	// neither the fill nor the provider that produced it — a true failure
	// attributed to the wrong component.
	for _, u := range a.Units {
		switch u.Disposition {
		case cellkernel.DispositionPass, cellkernel.DispositionFinding, cellkernel.DispositionUnverified:
		default:
			return nil, fmt.Errorf("the assessment reports disposition %q for unit %q; the vocabulary is pass, finding or unverified",
				u.Disposition, u.Unit)
		}
		if u.Disposition != cellkernel.DispositionPass && strings.TrimSpace(u.Reason) == "" {
			return nil, fmt.Errorf("unit %q is %q with no reason; a judgement that withholds its reason is not a review",
				u.Unit, u.Disposition)
		}
	}
	known := make(map[string]bool, len(c.Units))
	for _, u := range c.Units {
		known[u.ID] = true
	}
	// UNKNOWN before MISSING: a seat that invented a unit and dropped a real one
	// has done two different things wrong, and reporting only the absence would
	// hide the invention.
	for _, u := range a.Units {
		if !known[u.Unit] {
			return nil, fmt.Errorf("the assessment disposes of unit %q, which is not in the catalogue", u.Unit)
		}
	}
	for _, u := range c.Units {
		if _, ok := answered[u.ID]; !ok {
			return nil, fmt.Errorf("the assessment does not dispose of catalogue unit %q; coverage must be exact", u.ID)
		}
	}
	// Same set, so the only remaining coverage question is ORDER. It is checked
	// because the catalogue is a sequence: a reader comparing two episodes'
	// assessments reads them positionally, and an answer that permuted the units
	// would compare against the wrong obligations.
	if len(a.Units) != len(c.Units) {
		return nil, fmt.Errorf("the assessment carries %d dispositions for %d catalogue units", len(a.Units), len(c.Units))
	}
	for i, u := range c.Units {
		if a.Units[i].Unit != u.ID {
			return nil, fmt.Errorf("the assessment is out of catalogue order: position %d disposes of %q, the catalogue's is %q",
				i, a.Units[i].Unit, u.ID)
		}
	}

	out := make([]cellkernel.UnitResult, 0, len(c.Units))
	for i, u := range c.Units {
		got := a.Units[i]
		if u.Forced != nil {
			// The runtime measured this one. A seat that disagrees has claimed
			// authority it does not have — over a fact it was TOLD in its prompt
			// — so this is a malfunction and not a difference of opinion.
			if got.Disposition != u.Forced.Disposition {
				return nil, fmt.Errorf(
					"the assessment reports %q for unit %q, which the runtime measured as %q; a measured unit is not the seat's to decide",
					got.Disposition, u.ID, u.Forced.Disposition)
			}
			// The runtime's reason, not the seat's: the seat did not decide this
			// unit, so it is not the one that can explain it. Citations survive —
			// a pointer into the matter is the seat's own work.
			out = append(out, cellkernel.UnitResult{
				Unit:        u.ID,
				Disposition: u.Forced.Disposition,
				Reason:      u.Forced.Reason,
				Citations:   got.Citations,
			})
			continue
		}
		if !viewComplete && got.Disposition == cellkernel.DispositionPass {
			out = append(out, cellkernel.UnitResult{
				Unit:        u.ID,
				Disposition: cellkernel.DispositionUnverified,
				Reason: "the reconstructed view was incomplete, so a pass cannot be told apart from a pass on " +
					"content that was never shown; the seat's reason was: " + strings.TrimSpace(got.Reason),
				Citations: got.Citations,
			})
			continue
		}
		out = append(out, got)
	}
	return out, nil
}

// Pass reports whether every unit passed. It is the ONE derivation of the
// review's boolean from the dispositions, and the kernel's V re-derives the
// same fact independently from the record — so a seat cannot report findings
// and a passing verdict, whichever of the two it got wrong.
func Pass(units []cellkernel.UnitResult) bool {
	for _, u := range units {
		if u.Disposition != cellkernel.DispositionPass {
			return false
		}
	}
	return true
}

// Summary is the review's notes: a mechanical count, in catalogue order, of
// what was decided. It states no judgement of its own — the judgement is the
// unit list, and notes that paraphrased it would be a second account of the
// same thing.
func Summary(units []cellkernel.UnitResult) string {
	counts := map[cellkernel.Disposition]int{}
	for _, u := range units {
		counts[u.Disposition]++
	}
	return fmt.Sprintf("%d catalogue units: %d pass, %d finding, %d unverified",
		len(units),
		counts[cellkernel.DispositionPass],
		counts[cellkernel.DispositionFinding],
		counts[cellkernel.DispositionUnverified])
}

// AnswerSchema is the JSON Schema the provider's answer is constrained to. The
// unit ids are pinned as an `enum` and the length is pinned to the catalogue's,
// so the shape a provider can emit at all is already close to the shape
// Reconcile requires — but the schema is a convenience, never the authority:
// Reconcile re-checks coverage on the decoded value, because a provider that
// ignored or partially honoured the constraint must not thereby pass.
func AnswerSchema(c Catalogue) json.RawMessage {
	ids := make([]string, 0, len(c.Units))
	for _, u := range c.Units {
		ids = append(ids, u.ID)
	}
	enum, _ := json.Marshal(ids)
	n := len(c.Units)
	return json.RawMessage(fmt.Sprintf(`{"type":"object","additionalProperties":false,`+
		`"required":["units"],"properties":{"units":{"type":"array","minItems":%d,"maxItems":%d,`+
		`"items":{"type":"object","additionalProperties":false,"required":["unit","disposition","reason"],`+
		`"properties":{"unit":{"type":"string","enum":%s},`+
		`"disposition":{"type":"string","enum":["pass","finding","unverified"]},`+
		`"reason":{"type":"string"},`+
		`"citations":{"type":"array","items":{"type":"string"}}}}}}}`,
		n, n, enum))
}
