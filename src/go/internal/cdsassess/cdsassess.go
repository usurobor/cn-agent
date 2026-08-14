// Package cdsassess is the CDS-owned `cds.assess` fill: the constructor for a
// reviewing beta that decides every declared obligation and can say it could
// not decide one.
//
// It is the counterpart to `cds.patch`, and the asymmetry between them is the
// point. `cds.patch` produces, so it needs a workspace and file tools, and what
// it did is MEASURED afterwards. `cds.assess` judges, so it needs neither: its
// canonical input is `(contract, matter)` and its whole product is a verdict.
// It declares no workspace, is offered no tools, and returns a value the
// runtime reconciles rather than a change the runtime measures.
//
// THE SEAT IS NOT THE ONLY JUDGE, and that is what makes this more than a
// prompt. Two of the catalogue's units are decided by the runtime — whether
// the matter carries a reviewable change, and what the closed checker recipe
// observed against the reconstructed candidate — and a cognitive answer that
// contradicts either is a fault rather than a verdict. What cognition is asked
// for is exactly the part no mechanical check can supply: whether each
// acceptance criterion is met.
//
// INDEPENDENCE IS STRUCTURAL, NOT PROMISED. The seat receives only what
// cellkernel.BetaInput carries — a frozen contract copy and the sealed matter
// projection. It cannot reach alpha's worktree, artifacts, session or
// transcript, because it is never given them and has no tool with which to
// look. The runtime DOES cut a throwaway checkout from `(contract.subject,
// matter)` and reduce it to a bounded value which the seat is handed inside its
// prompt; production can therefore affect what the reviewer sees only by
// changing the matter, which is exactly what the reviewer is judging. The
// independence claim rests on the channel being closed, which this runtime can
// demonstrate, and not on a containment this substrate cannot.
package cdsassess

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cdsissue"
	"github.com/usurobor/cnos/src/go/internal/cellcheck"
	"github.com/usurobor/cnos/src/go/internal/cellcog"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellmethod"
	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

// Fill is the tag this package registers under.
const Fill = "cds.assess"

// The authored declaration is cellfill.SeatDecl, the closed {fill, cognition}
// shape this seat shares with `cds.patch`. Note what is absent from it: no
// workspace, because a reviewer that could open the worktree would not be
// reviewing the matter it was handed; and no skills, because the cell declares
// one methodology bundle and this seat receives a projection of it, exactly as
// `cds.patch` does. Neither key is admitted by the shared decoder, so those two
// absences are one rule rather than two copies of one.

// ResolvedDecl is what the closure records for this seat: the provider and the
// REQUESTED MODEL SELECTOR, and the identity of the projection that held it.
// The same "selector, not observed identity" caveat that applies to alpha
// applies here — nothing observes what actually served the request.
type ResolvedDecl struct {
	Fill        string              `json:"fill"`
	Cognition   cellcog.Config      `json:"cognition"`
	Methodology cellmethod.Recorded `json:"methodology"`
}

// Factory returns the cds.assess beta factory. It closes over nothing: the
// methodology arrives as an argument and the subject arrives on the contract,
// so there is no second place either could come from.
func Factory() cellfill.BetaFactory {
	return func(_ context.Context, decl json.RawMessage, method cellmethod.View) (cellfill.ConstructedBeta, error) {
		// The shared decode: the closed {fill, cognition} key language and the
		// projection-role check, identical for both CDS seats and owned by
		// neither. What remains here is this fill's own — the two refusals in
		// this seat's words, and the ANSWERING port, which is the one thing
		// about its cognition that differs from the producing side.
		//
		// The fill states its own requirement, and only the fill can: this seat
		// judges real code against a methodology, and there is nothing to hold
		// the candidate to without one.
		d, err := cellfill.AdmitSeatDecl(decl, Fill, cellmethod.RoleAdversarial, method, cellfill.SeatRefusal{
			NoMethodology: "an assessing beta needs the cell's methodology, and this cell declares none",
			WrongRole:     "an assessing seat takes the adversarial projection",
		})
		if err != nil {
			return cellfill.ConstructedBeta{}, err
		}

		answerer, mode, err := cellcog.NewAnswerer(d.Cognition)
		var seat judge
		switch {
		case errors.Is(err, cellcog.ErrNoDeterministicAnswer):
			// The deterministic path. cellcog cannot fabricate a judgement
			// because it owns no verdict vocabulary; this package does, so the
			// refusal is written here in the language the catalogue speaks.
			seat = refusingJudge{}
		case err != nil:
			return cellfill.ConstructedBeta{}, fmt.Errorf("fill %q: %w", Fill, err)
		default:
			seat = rentedJudge{answerer: answerer}
		}

		resolved := ResolvedDecl{Fill: Fill, Cognition: d.Cognition, Methodology: method.Recorded()}
		canon, err := json.Marshal(resolved)
		if err != nil {
			return cellfill.ConstructedBeta{}, fmt.Errorf("fill %q: canonicalize: %w", Fill, err)
		}
		return cellfill.ConstructedBeta{
			Constructed: cellfill.Constructed{Decl: canon, Mode: cellkernel.ExecutionMode(mode)},
			Seat:        AssessBeta{judge: seat, method: method},
		}, nil
	}
}

// AssessBeta is the reviewing seat. Its ONLY fields are the judge it consults
// and the projection it judges by: no repository, no directory, no worktree
// handle, no producer session, no transcript, no alpha result. That is asserted
// by a test over this struct's fields, because a promise about what a seat
// cannot reach is worth exactly as much as the field list behind it.
type AssessBeta struct {
	judge  judge
	method cellmethod.View
}

func (b AssessBeta) Review(ctx context.Context, in cellkernel.BetaInput) (cellkernel.BetaOutput, error) {
	if b.judge == nil {
		return cellkernel.BetaOutput{}, cellcog.ErrNoProvider
	}
	// The issue is admitted by the SAME predicate alpha's contract went through.
	// A reviewer given criteria it cannot read would be back to judging
	// plausibility, which is what the catalogue exists to replace.
	issue, err := cdsissue.Admit(in.Contract.Issue)
	if err != nil {
		return cellkernel.BetaOutput{}, fmt.Errorf("cds.assess: %w", err)
	}
	subject, err := cellwork.AdmitSubject(in.Contract.Subject)
	if err != nil {
		return cellkernel.BetaOutput{}, fmt.Errorf("cds.assess: %w", err)
	}

	view, obs, viewComplete, err := b.candidate(ctx, subject, in.Matter)
	if err != nil {
		return cellkernel.BetaOutput{}, err
	}
	cat := Build(issue).Decide(in.Matter, obs)

	answer, err := b.judge.Judge(ctx, RenderPrompt(in.Contract, issue, in.Matter, view, obs, cat, b.method), cat)
	if err != nil {
		return cellkernel.BetaOutput{}, fmt.Errorf("cds.assess: judge %q: %w", b.judge.Name(), err)
	}
	units, err := Reconcile(cat, answer, viewComplete)
	if err != nil {
		// A malformed or non-covering answer is a MALFUNCTION, not a verdict.
		// Inventing `pass:false` from an unreadable answer would fabricate a
		// judgement nobody made — the same reason a missing verdict is an error
		// rather than a failing one.
		return cellkernel.BetaOutput{}, fmt.Errorf("cds.assess: judge %q did not review: %w", b.judge.Name(), err)
	}
	// Artifacts are gamma/V's channel, never beta's: the assessment is the whole
	// product of this seat.
	return cellkernel.BetaOutput{Review: cellkernel.Review{
		Pass:       Pass(units),
		Notes:      Summary(units),
		Assessment: units,
	}}, nil
}

// candidate derives the reviewing seat's evidence from `(contract.subject,
// matter)` and nothing else: a throwaway checkout at the pinned base with the
// matter applied, reduced to a bounded value, with the closed checker recipe
// run against that same directory before it is released.
//
// ONE materialization serves both. Reconstructing for the value and then
// materializing again for the checker would measure two trees and report one,
// which is the class of defect the single pinned subject exists to prevent.
//
// A matter with nothing reviewable in it short-circuits: there is no patch to
// apply, so no candidate exists, and the observation says the recipe never ran
// rather than blaming the candidate for a change that was never made. A
// reconstruction that FAILS on a matter that did carry a diff is a different
// thing — the runtime could not build a view it should have been able to build
// — and surfaces as an episode malfunction rather than as a verdict.
func (b AssessBeta) candidate(ctx context.Context, subject cellwork.Subject, m cellkernel.Matter) (cellwork.View, cellcheck.Observation, bool, error) {
	notRun := cellcheck.Observation{Recipe: cellcheck.RecipeID, Status: cellcheck.Unavailable}
	if MatterFault(m) != "" {
		return cellwork.View{}, notRun, false, nil
	}
	obs := notRun
	view, err := cellwork.Reconstruct(ctx, subject.Repo, subject.BaseSHA, m.Data, func(dir string) {
		obs = cellcheck.Run(ctx, dir, subject.BaseSHA)
	})
	if err != nil {
		return cellwork.View{}, notRun, false, fmt.Errorf("cds.assess: reconstruct the candidate view: %w", err)
	}
	return view, obs, !view.Truncated, nil
}

// judge is the cognition seam this fill needs, declared HERE because this is
// where it is consumed (eng/go §2.3). It is narrower than cellcog.Answerer on
// purpose: what this seat needs is an Assessment over a Catalogue, and the two
// implementations below differ in whether anything was rented — not in what
// they are asked for.
type judge interface {
	Name() string
	Judge(ctx context.Context, prompt string, c Catalogue) (Assessment, error)
}

// rentedJudge asks a provider, constrained to the catalogue's own schema.
type rentedJudge struct{ answerer cellcog.Answerer }

func (j rentedJudge) Name() string { return j.answerer.Name() }

func (j rentedJudge) Judge(ctx context.Context, prompt string, c Catalogue) (Assessment, error) {
	raw, err := j.answerer.Answer(ctx, prompt, AnswerSchema(c))
	if err != nil {
		return Assessment{}, err
	}
	// Strict, because a malformed verdict is not a failing verdict. A seat whose
	// answer cannot be read has not reviewed.
	var a Assessment
	if err := cellfill.StrictDecode(raw, &a); err != nil {
		return Assessment{}, fmt.Errorf("the answer does not match the requested schema: %w", err)
	}
	return a, nil
}

// refusingJudge is the deterministic seat: it rents nothing, so it forms no
// judgement, and it says so per unit rather than passing or failing.
//
// It never passes an acceptance unit, and that is the whole of its honesty. A
// fake that returned `pass` would be a false completion wearing a reviewer's
// clothes, and one that guessed from the prompt would be inventing judgement it
// does not have. For the two units the RUNTIME decided it returns exactly what
// the runtime decided — not because it agrees, but because those units are not
// its to answer, and an answer that contradicted them would be the same fault
// a rented seat's contradiction is.
type refusingJudge struct{}

func (refusingJudge) Name() string { return "no-cognition" }

func (refusingJudge) Judge(_ context.Context, _ string, c Catalogue) (Assessment, error) {
	units := make([]cellkernel.UnitResult, 0, len(c.Units))
	for _, u := range c.Units {
		if u.Forced != nil {
			units = append(units, cellkernel.UnitResult{
				Unit: u.ID, Disposition: u.Forced.Disposition, Reason: u.Forced.Reason,
			})
			continue
		}
		units = append(units, cellkernel.UnitResult{
			Unit:        u.ID,
			Disposition: cellkernel.DispositionUnverified,
			Reason:      "no cognition was rented, so this criterion was not decided by anyone",
		})
	}
	return Assessment{Units: units}, nil
}

// RenderPrompt builds the reviewing seat's entire world: the contract it judges
// against, the admitted issue whose criteria it decides, the matter, the
// runtime-derived view of what that matter produces, what the closed checker
// observed, the catalogue it must dispose of, and the adversarial projection of
// the cell's methodology. Nothing else reaches the seat — in particular no
// tool, so the view is a bounded VALUE rather than a directory it could open.
//
// The issue block is cdsissue.Render — the same call the producing seat's
// contract went through. β is told exactly what α was told it must satisfy,
// which is the property that makes verification cheaper than production.
func RenderPrompt(c cellkernel.Contract, issue cdsissue.Issue, m cellkernel.Matter, view cellwork.View,
	obs cellcheck.Observation, cat Catalogue, method cellmethod.View) string {
	var b strings.Builder
	b.WriteString("You are the beta (assessing) seat of a CNOS coherence cell.\n")
	b.WriteString("You did not produce this work and you cannot see the workspace it came\n")
	b.WriteString("from. You have no tools. Judge ONLY what is below.\n\n")
	fmt.Fprintf(&b, "CONTRACT %s\nGOAL: %s\n\n", c.ID, c.Goal)
	b.WriteString(cdsissue.Render(issue))

	b.WriteString("\nHOW TO JUDGE\n")
	b.WriteString("Return exactly one disposition for EVERY unit of the catalogue below, in\n")
	b.WriteString("the catalogue's order, and no others. Each is one of:\n")
	b.WriteString("  pass       — the obligation is met, and you can point at why\n")
	b.WriteString("  finding    — you checked and it is not met\n")
	b.WriteString("  unverified — you could not decide it from what you were shown\n")
	b.WriteString("`unverified` is a real answer and the right one whenever the evidence you\n")
	b.WriteString("need is absent. Do not pass an obligation you could not check: absence of\n")
	b.WriteString("evidence is not evidence, and a confident pass on something you did not\n")
	b.WriteString("see is the single worst thing this seat can do.\n")
	b.WriteString("Every finding and every unverified needs a reason. A judgement without a\n")
	b.WriteString("reason is not review.\n")
	b.WriteString("Decide each acceptance unit by its criterion's stated verification route.\n")
	b.WriteString("A change that is real, large or well written but does not meet the goal is\n")
	b.WriteString("a finding.\n")
	b.WriteString("Where a criterion is about the STATE of a file rather than about the change\n")
	b.WriteString("to it, decide it from the reconstructed view below, not from the hunks: a\n")
	b.WriteString("diff shows the lines around a change and nothing else, so a symbol you\n")
	b.WriteString("cannot see in a hunk may well be present in the file.\n")
	b.WriteString("The two `check:` units were already measured by the runtime and their\n")
	b.WriteString("dispositions are stated in the catalogue. Repeat them exactly. They are not\n")
	b.WriteString("yours to decide, and an answer that differs from them is discarded whole.\n")

	fmt.Fprintf(&b, "\n===== MATTER (%d bytes) =====\n", len(m.Data))
	b.WriteString(m.Data)
	if !strings.HasSuffix(m.Data, "\n") {
		b.WriteString("\n")
	}
	b.WriteString(renderView(view))
	b.WriteString(renderObservation(obs))
	b.WriteString(renderCatalogue(cat))
	fmt.Fprintf(&b, "\n(methodology bundle sha256 %s)\n", method.SHA256)
	b.WriteString(method.Text)
	return b.String()
}

// renderView writes the reconstructed view into the prompt, and says exactly
// what it is. The provenance sentence is load-bearing: a reviewer that read
// this block as the producer's account of its own work would be back to judging
// a claim, which is the thing the reconstruction removes.
//
// Omitted content is named as omitted rather than shown as absent, and a
// truncated view is stated at the top, because a reviewer must be able to tell
// "I checked and it is not there" from "I was not shown it".
func renderView(v cellwork.View) string {
	var b strings.Builder
	fmt.Fprintf(&b, "\n===== RECONSTRUCTED VIEW (%d files) =====\n", len(v.Files))
	b.WriteString("This is the post-application content of every file the matter above\n")
	b.WriteString("touches. The RUNTIME built it by checking out the contract's pinned base\n")
	b.WriteString("and applying that same patch to a throwaway copy. It is not the\n")
	b.WriteString("producer's account of anything, and nothing else was consulted.\n")
	if len(v.Files) == 0 {
		b.WriteString("\nTHERE IS NO VIEW: no candidate was reconstructed, because the matter\n")
		b.WriteString("carries no change to apply. Every unit whose evidence would have come\n")
		b.WriteString("from the candidate is unverified.\n")
	}
	if v.Truncated {
		b.WriteString("\nTHIS VIEW IS INCOMPLETE: it hit its size bound, so some file content\n")
		b.WriteString("below is omitted. A unit whose verification needs content you were not\n")
		b.WriteString("shown is `unverified`, never `pass`.\n")
	}
	for _, f := range v.Files {
		fmt.Fprintf(&b, "\n--- %s (%s", f.Path, f.Status)
		if f.From != "" {
			fmt.Fprintf(&b, " from %s", f.From)
		}
		b.WriteString(") ---\n")
		switch {
		case f.Status == cellwork.FileDeleted:
			b.WriteString("(deleted by this change; it has no content after the change)\n")
		case f.Binary:
			b.WriteString("(binary; its content is not shown)\n")
		case f.Symlink:
			fmt.Fprintf(&b, "(symbolic link to %q; that target path IS its content, and the\n"+
				"runtime did not follow it — nothing at the other end was read)\n", f.Content)
		case f.Omitted:
			b.WriteString("(content omitted: this view hit its size bound)\n")
		case f.Content == "":
			b.WriteString("(the file is empty after the change)\n")
		default:
			b.WriteString(f.Content)
			if !strings.HasSuffix(f.Content, "\n") {
				b.WriteString("\n")
			}
		}
	}
	return b.String()
}

// renderObservation states what the closed recipe did, step by step. The seat
// is shown it rather than asked to reproduce it: a cognitive seat cannot run
// anything, so a checker result it was not shown would be a fact it could only
// guess at — and one it WAS shown is a fact it must not contradict.
func renderObservation(obs cellcheck.Observation) string {
	var b strings.Builder
	fmt.Fprintf(&b, "\n===== CHECKER OBSERVATION (%s) =====\n", obs.Recipe)
	b.WriteString("The runtime ran this closed recipe against the reconstructed candidate\n")
	b.WriteString("above. The recipe is fixed and neither seat can change it. It stops at\n")
	b.WriteString("the first step that does not pass.\n")
	fmt.Fprintf(&b, "overall: %s\n", obs.Status)
	if len(obs.Steps) == 0 {
		b.WriteString("(no step ran)\n")
	}
	for _, s := range obs.Steps {
		fmt.Fprintf(&b, "\n--- step %s: %s (exit %d) ---\n", s.Name, s.Status, s.Exit)
		if tail := strings.TrimSpace(s.Tail); tail != "" {
			b.WriteString(tail)
			b.WriteString("\n")
		}
	}
	return b.String()
}

// renderCatalogue writes the ordered obligations, each acceptance unit with the
// verification route its criterion declared and each check unit with the
// disposition the runtime already measured.
func renderCatalogue(c Catalogue) string {
	var b strings.Builder
	fmt.Fprintf(&b, "\n===== ASSESSMENT CATALOGUE (%d units, in order) =====\n", len(c.Units))
	for i, u := range c.Units {
		fmt.Fprintf(&b, "\n%d. %s [%s]\n", i+1, u.ID, u.Kind)
		if u.Statement != "" {
			fmt.Fprintf(&b, "   %s\n", u.Statement)
		}
		if u.Verification != "" {
			fmt.Fprintf(&b, "   verified by: %s\n", u.Verification)
		}
		if u.Forced != nil {
			fmt.Fprintf(&b, "   ALREADY MEASURED BY THE RUNTIME: %s — %s\n", u.Forced.Disposition, u.Forced.Reason)
		}
	}
	return b.String()
}
