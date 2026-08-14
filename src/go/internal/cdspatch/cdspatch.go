// Package cdspatch is the CDS-owned `cds.patch` fill: the constructor for a
// patch-producing alpha. Only this fill knows that producing a CDS patch takes
// workspace cognition and a git worktree — the runner dispatches a fill id and
// receives a cellkernel.Alpha. cellcog builds the bounded provider adapter (no
// provider argv lives in this package) and cellwork prepares the worktree.
//
// THE SEAT DOES NOT DECLARE ITS SKILLS. It receives the cell's constructive
// projection and records the bundle digest it was held to, because two lists of
// obligations — one on the cell, one on the seat — drift with nothing able to
// notice. The `skills` key is gone from both authorities, Go and CUE.
//
// THE SEAT DOES NOT NAME A REPOSITORY: repository and base come from the pinned
// contract subject, read at Produce. One source read once makes two disagreeing
// repository declarations unrepresentable (CDS-CELL-MIGRATION.md).
package cdspatch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cdsissue"
	"github.com/usurobor/cnos/src/go/internal/cellcog"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellmethod"
	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

// Fill is the tag this package registers under.
const Fill = "cds.patch"

// Artifact identities a patch episode produces. The CDS canonical diff-first
// rule requires exactly {id "diff", kind "diff"}.
const (
	DiffArtifactID   = "diff"
	DiffArtifactKind = "diff"
	BaseArtifactID   = "base_sha"
	BaseArtifactKind = "sha"
)

// The authored declaration is cellfill.SeatDecl: the closed {fill, cognition}
// shape, shared with `cds.assess` because a second copy would be a second key
// language. The CUE overlays pin the same shapes for the independent oracle.

// ResolvedDecl is what the closure records for this seat: the provider and the
// REQUESTED MODEL SELECTOR, as deterministic bytes — the canonical form the
// scope-lift digest covers. "Selector", precisely (Pi #57 B2): a provider may
// remap a selector to another model, the Claude CLI says so on stderr when it
// does, and nothing here verifies what actually served the request.
type ResolvedDecl struct {
	Fill      string         `json:"fill"`
	Cognition cellcog.Config `json:"cognition"`
	// Methodology is what the seat RECEIVED, not what it chose: the projection's role
	// and bundle digest are how a reader asks whether the seat was held to the cell's
	// methodology. Refs and body digests live once, on the bundle.
	Methodology cellmethod.Recorded `json:"methodology"`
}

// Factory closes over nothing: the constructive projection arrives as an argument,
// so the fill can refuse the cell's one bundle but never replace it with its own.
func Factory() cellfill.AlphaFactory {
	return func(ctx context.Context, decl json.RawMessage, method cellmethod.View) (cellfill.ConstructedAlpha, error) {
		// The shared decode owns the key language and the projection-role check every
		// seat of this shape needs; what stays here is genuinely this fill's. Only the
		// fill can state that a patch alpha needs a methodology — a cell declaring none
		// is legitimate for other fills, so this refuses here, not in the loader.
		d, err := cellfill.AdmitSeatDecl(decl, Fill, cellmethod.RoleConstructive, method, cellfill.SeatRefusal{
			NoMethodology: "a patch alpha needs the cell's methodology, and this cell declares none",
			WrongRole:     "a producing seat takes the constructive projection",
		})
		if err != nil {
			return cellfill.ConstructedAlpha{}, err
		}

		coder, mode, err := cellcog.New(d.Cognition)
		if err != nil {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: %w", Fill, err)
		}
		// Nothing here pins a revision or loads a skill: both happened once before this
		// constructor ran, so the declaration records what it selects and receives.
		resolved := ResolvedDecl{
			Fill:        Fill,
			Cognition:   d.Cognition,
			Methodology: method.Recorded(),
		}
		canon, err := json.Marshal(resolved)
		if err != nil {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: canonicalize: %w", Fill, err)
		}

		return cellfill.ConstructedAlpha{
			Constructed: cellfill.Constructed{Decl: canon, Mode: cellkernel.ExecutionMode(mode)},
			Seat:        PatchAlpha{coder: coder, method: method},
		}, nil
	}
}

// PatchAlpha is the provider-neutral patch-producing seat: it materializes a
// disposable worktree at the base, lets the coder change files in it, then
// MEASURES the change as a unified diff. Nothing trusts the coder's account of
// itself — one that claims a sweeping refactor and wrote nothing produces no diff,
// and an episode with no diff cannot satisfy a contract requiring one, so false
// completion is unrepresentable rather than caught late. It holds no repository
// and no base: both are read from the pinned subject on every Produce, so the
// episode acts on the one value the record says it acted on.
type PatchAlpha struct {
	coder  cellcog.Coder
	method cellmethod.View
}

func (a PatchAlpha) Produce(ctx context.Context, in cellkernel.AlphaInput) (cellkernel.AlphaOutput, error) {
	if a.coder == nil {
		return cellkernel.AlphaOutput{}, cellcog.ErrNoProvider
	}
	// AdmitSubject, not ParseSubject: a station requires a base already pinned, since
	// re-resolving a moving name is how two stations could measure different trees
	// while one record claimed one; an absent subject fails rather than defaulting.
	subject, err := cellwork.AdmitSubject(in.Contract.Subject)
	if err != nil {
		return cellkernel.AlphaOutput{}, fmt.Errorf("cds.patch: %w", err)
	}
	// Admitted here as well as at the door, not because the door is distrusted: this
	// seat reads the FROZEN contract, so admitting the bytes it was handed makes the
	// issue it renders the issue the episode recorded, and fails before a worktree.
	issue, err := cdsissue.Admit(in.Contract.Issue)
	if err != nil {
		return cellkernel.AlphaOutput{}, fmt.Errorf("cds.patch: %w", err)
	}
	wt, release, err := cellwork.Materialize(ctx, subject.Repo, subject.BaseSHA)
	if err != nil {
		return cellkernel.AlphaOutput{}, err
	}
	defer release()

	if err := a.coder.Work(ctx, wt.Dir, RenderPrompt(in.Contract, issue, a.method)); err != nil {
		return cellkernel.AlphaOutput{}, fmt.Errorf("cds.patch: coder %q: %w", a.coder.Name(), err)
	}
	diff, err := wt.Diff(ctx)
	if err != nil {
		return cellkernel.AlphaOutput{}, err
	}

	// base_sha is MEASURED from the materialized worktree — the coder never sees it
	// — so it is what the episode stood on, not a copy of what the contract asked
	// for. Equal to contract.subject.base_sha when honest; a reader can compare.
	artifacts := []cellkernel.ArtifactCandidate{
		{ID: BaseArtifactID, Kind: BaseArtifactKind, Text: wt.BaseSHA},
	}
	if strings.TrimSpace(diff) == "" {
		return cellkernel.AlphaOutput{
			Matter:    cellkernel.Matter{Data: fmt.Sprintf("no change was made to %s at %s", subject.Repo, wt.BaseSHA)},
			Artifacts: artifacts,
		}, nil
	}
	// The diff is both what beta reviews (matter) and what V checks (evidence):
	// same bytes, two roles — beta never receives artifacts, V never reads matter.
	return cellkernel.AlphaOutput{
		Matter: cellkernel.Matter{Data: diff},
		Artifacts: append(artifacts,
			cellkernel.ArtifactCandidate{ID: DiffArtifactID, Kind: DiffArtifactKind, Text: diff}),
	}, nil
}

// RenderPrompt is pure over the contract and the CONSTRUCTIVE PROJECTION, whose
// text carries the exact skill bodies the cell's bundle loaded. It appends that
// view rather than rendering skills, so the seat holds no second opinion.
func RenderPrompt(c cellkernel.Contract, issue cdsissue.Issue, method cellmethod.View) string {
	var b strings.Builder
	b.WriteString("You are the alpha (producing) seat of a CNOS coherence cell, working on real code.\n")
	b.WriteString("You are in a disposable worktree. Edit the files here to meet the contract.\n\n")
	fmt.Fprintf(&b, "CONTRACT %s\nGOAL: %s\n\n", c.ID, c.Goal)
	// The ISSUE, through the same cdsissue.Render the assessing seat uses: a seat
	// given only the one-sentence goal writes against a summary of its contract and
	// is then reviewed against the contract — one renderer over the same bytes.
	b.WriteString(cdsissue.Render(issue))
	b.WriteString("\nHOW YOUR WORK IS RECORDED\n")
	b.WriteString("Your change is measured as a unified diff of this worktree, not taken from\n")
	b.WriteString("your summary. Anything you do not write to a file does not exist. An empty\n")
	b.WriteString("diff closes the episode as unmet, whatever you say about it.\n")
	b.WriteString("An independent reviewer reads that diff afterwards.\n")
	fmt.Fprintf(&b, "\n(methodology bundle sha256 %s)\n", method.SHA256)
	b.WriteString(method.Text)
	return b.String()
}
