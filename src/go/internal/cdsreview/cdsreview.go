// Package cdsreview is the CDS-owned `cds.review` fill: the constructor for a
// rented reviewing beta. It is the Case-3 counterpart to `cds.patch`, and the
// asymmetry between them is the point.
//
// `cds.patch` produces, so it needs a workspace and file tools, and what it
// did is MEASURED afterwards. `cds.review` judges, so it needs neither: its
// canonical input is `(contract, matter)` and its whole product is a verdict.
// It therefore declares no workspace, is offered no tools, and returns a
// value the runtime decodes rather than a change the runtime measures.
//
// Independence is structural, not promised. The seat receives only what
// `cellkernel.BetaInput` carries — a frozen contract copy and the sealed
// matter projection. It cannot reach alpha's worktree, artifacts or internals,
// because it is never given them and has no tool with which to look.
package cdsreview

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cdsissue"
	"github.com/usurobor/cnos/src/go/internal/cellcog"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
)

// Fill is the tag this package registers under.
const Fill = "cds.review"

// Decl is the complete, closed shape of a cds.review beta declaration. Note
// what is absent: there is no workspace, because a reviewer that could open
// the worktree would not be reviewing the matter it was handed.
type Decl struct {
	Fill      string         `json:"fill"`
	Cognition cellcog.Config `json:"cognition"`
	Skills    []string       `json:"skills"`
}

// ResolvedDecl is what the closure records for this seat: the declaration with
// skills expanded to ordered canonical refs + content digests, and the
// provider plus the requested model selector — the same "selector, not
// observed identity" caveat that applies to alpha applies here.
type ResolvedDecl struct {
	Fill      string          `json:"fill"`
	Cognition cellcog.Config  `json:"cognition"`
	Skills    []ResolvedSkill `json:"skills"`
}

type ResolvedSkill struct {
	Ref    string `json:"ref"`
	SHA256 string `json:"sha256"`
}

// VerdictSchema is the shape the reviewer's answer must satisfy. The provider
// constrains its output to this, so the verdict is DECODED rather than parsed
// hopefully out of prose. It mirrors cellkernel.Review exactly; a reviewer
// cannot report a field the kernel has no place for.
var VerdictSchema = json.RawMessage(`{"type":"object",` +
	`"properties":{"pass":{"type":"boolean"},"notes":{"type":"string"}},` +
	`"required":["pass","notes"],"additionalProperties":false}`)

// Factory returns the cds.review beta factory, closed over the skill resolver
// it loads bodies from (the one construction-time effect).
func Factory(skills cellskill.Resolver) cellfill.BetaFactory {
	return func(_ context.Context, decl json.RawMessage) (cellfill.ConstructedBeta, error) {
		var d Decl
		if err := exactShape(decl); err != nil {
			return cellfill.ConstructedBeta{}, fmt.Errorf("fill %q: %w", Fill, err)
		}
		if err := cellfill.StrictDecode(decl, &d); err != nil {
			return cellfill.ConstructedBeta{}, fmt.Errorf("fill %q: %w", Fill, err)
		}
		if len(d.Skills) == 0 {
			return cellfill.ConstructedBeta{}, fmt.Errorf("fill %q: a reviewing beta needs its skills", Fill)
		}

		answerer, mode, err := cellcog.NewAnswerer(d.Cognition)
		if err != nil {
			return cellfill.ConstructedBeta{}, fmt.Errorf("fill %q: %w", Fill, err)
		}
		loaded, err := cellskill.LoadAll(skills, d.Skills)
		if err != nil {
			return cellfill.ConstructedBeta{}, fmt.Errorf("fill %q: %w", Fill, err)
		}

		resolved := ResolvedDecl{
			Fill:      Fill,
			Cognition: d.Cognition,
			Skills:    make([]ResolvedSkill, 0, len(loaded)),
		}
		for _, s := range loaded {
			resolved.Skills = append(resolved.Skills, ResolvedSkill{Ref: s.Ref, SHA256: s.SHA256})
		}
		canon, err := json.Marshal(resolved)
		if err != nil {
			return cellfill.ConstructedBeta{}, fmt.Errorf("fill %q: canonicalize: %w", Fill, err)
		}

		return cellfill.ConstructedBeta{
			Constructed: cellfill.Constructed{Decl: canon, Mode: cellkernel.ExecutionMode(mode)},
			Seat:        ReviewBeta{answerer: answerer, skills: loaded},
		}, nil
	}
}

// exactShape states this fill's closed key language explicitly. encoding/json
// matches field names case-insensitively even with DisallowUnknownFields, so
// `Fill` or a nested `Provider` would otherwise decode in Go while the closed
// CUE overlay rejects them.
func exactShape(decl json.RawMessage) error {
	if err := cellfill.OnlyKeys(decl, "cds.review", "fill", "cognition", "skills"); err != nil {
		return err
	}
	if cog, ok := cellfill.Field(decl, "cognition"); ok {
		if err := cellfill.OnlyKeys(cog, "cds.review.cognition", "provider", "model"); err != nil {
			return err
		}
	}
	return nil
}

// ReviewBeta is the rented reviewing seat. It holds an answerer and the loaded
// skill bodies; it holds no repository, no directory and no artifacts.
type ReviewBeta struct {
	answerer cellcog.Answerer
	skills   []cellskill.Skill
}

func (b ReviewBeta) Review(ctx context.Context, in cellkernel.BetaInput) (cellkernel.BetaOutput, error) {
	if b.answerer == nil {
		return cellkernel.BetaOutput{}, cellcog.ErrNoProvider
	}
	// The contract's issue is admitted before any cognition is rented, and by
	// the same predicate alpha used. A reviewer given criteria it cannot read
	// would be back to judging plausibility.
	issue, err := cdsissue.Admit(in.Contract.Task)
	if err != nil {
		return cellkernel.BetaOutput{}, fmt.Errorf("cds.review: %w", err)
	}
	if fault := matterFault(in.Matter); fault != "" {
		// Not an error: the episode ran, and its outcome is that nothing
		// reviewable reached the reviewer. Returning pass:false without renting
		// anything is the honest verdict — a reviewer handed no diff has not
		// passed it, and asking a provider to judge `cds.patch`'s "no change
		// was made to ..." sentence buys an opinion about a sentence.
		//
		// The verdict shape is unchanged this cycle: `unverified` as a distinct
		// outcome from `judged and failed` is the next cycle's work, so today
		// the distinction lives in the notes.
		return cellkernel.BetaOutput{Review: cellkernel.Review{Pass: false, Notes: fault}}, nil
	}
	raw, err := b.answerer.Answer(ctx, RenderPrompt(in.Contract, issue, in.Matter, b.skills), VerdictSchema)
	if err != nil {
		return cellkernel.BetaOutput{}, fmt.Errorf("cds.review: answerer %q: %w", b.answerer.Name(), err)
	}
	// Strict, because a malformed verdict is not a failing verdict. A reviewer
	// whose answer cannot be read has not reviewed, and inventing `pass:false`
	// from an unreadable answer would fabricate a judgement nobody made.
	var v cellkernel.Review
	if err := cellfill.StrictDecode(raw, &v); err != nil {
		return cellkernel.BetaOutput{}, fmt.Errorf("cds.review: verdict does not match the requested schema: %w", err)
	}
	if strings.TrimSpace(v.Notes) == "" {
		return cellkernel.BetaOutput{}, fmt.Errorf("cds.review: verdict carries no notes; a judgement must state its reason")
	}
	// Artifacts are gamma/V's channel, never beta's — the verdict is the whole
	// product of this seat.
	return cellkernel.BetaOutput{Review: v}, nil
}

// matterFault reports why the matter cannot be reviewed, or "" if it can.
//
// `cds.patch`'s product is a unified diff, so this is the whole admission
// rule: at least one `diff --git ` file header. It is a structural check, not
// a diff parser — the reviewer reads the diff, this only decides whether there
// is one to read. The alternative, sending everything to a provider and
// trusting it to notice, is what produced a passing verdict on the sentence
// "no change was made to ...".
func matterFault(m cellkernel.Matter) string {
	switch {
	case strings.TrimSpace(m.Data) == "":
		return "not reviewed: the matter is empty, so there is no change to judge"
	case !strings.Contains("\n"+m.Data, "\ndiff --git "):
		return "not reviewed: the matter carries no unified diff (no `diff --git ` file header), " +
			"so there is no change to judge"
	}
	return ""
}

// RenderPrompt builds the reviewer's entire world: the contract it judges
// against, the issue whose criteria it decides, the matter it judges, and the
// skill bodies it judges by. Nothing else reaches the seat.
//
// The issue block is cdsissue.Render — the same call the producing seat makes.
// β is told exactly what α was told it must satisfy, which is the property
// that makes verification cheaper than production.
func RenderPrompt(c cellkernel.Contract, issue cdsissue.Issue, m cellkernel.Matter, skills []cellskill.Skill) string {
	var b strings.Builder
	b.WriteString("You are the beta (reviewing) seat of a CNOS coherence cell.\n")
	b.WriteString("You did not produce this work and you cannot see the workspace it came\n")
	b.WriteString("from. Judge ONLY the matter below against the contract below.\n\n")
	fmt.Fprintf(&b, "CONTRACT %s\nGOAL: %s\n\n", c.ID, c.Goal)
	b.WriteString(cdsissue.Render(issue))
	b.WriteString("\nHOW TO JUDGE\n")
	b.WriteString("Decide each acceptance criterion above by its stated verification route.\n")
	b.WriteString("Pass only if the matter actually meets the goal. A change that is real,\n")
	b.WriteString("large or well written but does not meet the goal FAILS. Absence of\n")
	b.WriteString("evidence is failure, not doubt: you cannot run anything, so a claim you\n")
	b.WriteString("cannot see supported in the matter is unsupported.\n")
	b.WriteString("State the reason in your notes; a verdict without a reason is not review.\n")
	fmt.Fprintf(&b, "\n===== MATTER (%d bytes) =====\n", len(m.Data))
	b.WriteString(m.Data)
	if !strings.HasSuffix(m.Data, "\n") {
		b.WriteString("\n")
	}
	for _, s := range skills {
		sum := sha256.Sum256([]byte(s.Body))
		fmt.Fprintf(&b, "\n===== SKILL %s (sha256 %s) =====\n", s.Ref, hex.EncodeToString(sum[:8]))
		b.WriteString(s.Body)
		if !strings.HasSuffix(s.Body, "\n") {
			b.WriteString("\n")
		}
	}
	return b.String()
}
