// Package cdspatch is the CDS-owned `cds.patch` fill: the constructor for a
// patch-producing alpha. Only this fill knows that producing a CDS patch
// takes workspace cognition and a git worktree — the generic runner dispatches
// a fill id and receives a cellkernel.Alpha, nothing more.
//
// The constructor composes two reusable subsystems and owns none of their
// internals: cellcog constructs the bounded provider adapter (this package
// contains no provider argv at all), and cellwork prepares the disposable
// worktree. What comes back is one immutable, provider-neutral PatchAlpha.
//
// THE SEAT DOES NOT DECLARE ITS SKILLS. It receives the cell's constructive
// methodology projection and records the bundle digest it was held to. The
// `skills` key is gone from this fill in both authorities — the Go decoder
// below and the closed CUE overlay — because two lists of obligations, one on
// the cell and one on the seat, are two lists that drift with nothing able to
// notice.
//
// THE SEAT DOES NOT NAME A REPOSITORY. The declaration carries no workspace:
// the repository and the base commit come from the pinned contract subject,
// read at Produce with cellwork.AdmitSubject. This is not tidying. While the
// fill resolved its own repo and base, a run could carry two repository
// declarations that disagreed — the record's contract.subject naming one and
// its resolved_spec.alpha naming another — and the closure self-verified,
// because each was internally consistent. Even authored identically they were
// two ResolveBase calls at different times, so a moving HEAD could pin twice.
// One source, read once, and the disagreement is unrepresentable rather than
// caught.
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

// Decl is the complete, closed shape of a cds.patch alpha declaration. The
// fill owns this strictness: unknown keys fail StrictDecode, and the CUE
// overlays #CDSPatchAlphaAuthored and #CDSPatchAlphaResolved pin the same
// shapes for the independent oracle: authored admits the structurally
// possible forms, and resolution plus this constructor validate the selected
// provider/model combination.
type Decl struct {
	Fill      string         `json:"fill"`
	Cognition cellcog.Config `json:"cognition"`
}

// ResolvedDecl is what the closure records for this seat: the declaration
// with skills expanded to ordered canonical refs + content digests, and the
// provider plus the REQUESTED MODEL SELECTOR. Deterministic bytes — this is
// the canonical form the scope-lift digest covers.
//
// "Selector", precisely (Pi #57 B2): the recorded model is what the cell
// asked for, not an independently observed immutable model identity. A
// provider may remap a selector to another model — the Claude CLI says so on
// stderr when it does — and nothing here observes or verifies what actually
// served the request. Recording the served identity would need the provider
// to report it and the runtime to check it; neither exists yet, so the field
// claims only what it is.
type ResolvedDecl struct {
	Fill      string         `json:"fill"`
	Cognition cellcog.Config `json:"cognition"`
	// Methodology is what the seat RECEIVED, not what it chose: this fill has
	// no skills key and cannot state one. Recording the projection's role and
	// the bundle digest is how a reader can ask whether the seat was held to
	// the methodology the cell declared — the refs and their body digests live
	// once, on the bundle, and are not copied per seat.
	Methodology cellmethod.Recorded `json:"methodology"`
}

// Factory returns the cds.patch alpha factory.
//
// It closes over nothing. While it held a skill resolver, this fill LOADED its
// own skills from its own list — a second methodology beside the cell's, with
// its own digest and no statement that the two were meant to agree. The
// constructive projection now arrives as an argument, so what holds this seat
// is the cell's one bundle and the fill can only refuse it, never replace it.
func Factory() cellfill.AlphaFactory {
	return func(ctx context.Context, decl json.RawMessage, method cellmethod.View) (cellfill.ConstructedAlpha, error) {
		var d Decl
		if err := exactShape(decl); err != nil {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: %w", Fill, err)
		}
		if err := cellfill.StrictDecode(decl, &d); err != nil {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: %w", Fill, err)
		}
		// The fill states its own requirement, and only the fill can: a patch
		// alpha writes real code and there is nothing to hold it to without a
		// methodology. A cell declaring none is legitimate for other fills, so
		// this is refused here rather than by the loader for everyone.
		if method.Empty() {
			return cellfill.ConstructedAlpha{}, fmt.Errorf(
				"fill %q: a patch alpha needs the cell's methodology, and this cell declares none", Fill)
		}
		if method.Role != cellmethod.RoleConstructive {
			// A producing seat handed the adversarial projection would be told to
			// falsify the work it is about to do. Nothing constructs it that way
			// today; this refuses the wiring mistake rather than trusting that.
			return cellfill.ConstructedAlpha{}, fmt.Errorf(
				"fill %q: a producing seat takes the constructive projection, got %q", Fill, method.Role)
		}

		coder, mode, err := cellcog.New(d.Cognition)
		if err != nil {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: %w", Fill, err)
		}
		// Nothing here pins a revision and nothing here loads a skill. The
		// subject was pinned once and the methodology was loaded once, both
		// before this constructor ran; this declaration records what it
		// actually selects — the provider — and what it was handed.
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

// exactShape states this fill's closed key language explicitly, at each of
// its three object shapes. encoding/json matches field names
// case-insensitively even with DisallowUnknownFields, so `Fill`, `Cognition`
// or a nested `Provider` would otherwise decode in Go while the closed CUE
// overlay rejects them. The fill owns this, not the generic runner.
// `workspace` is absent from the allowed set, and that absence is the whole
// point of the deletion: a declaration that still names a repository is now
// refused by name here and by the closed CUE overlay, so the second source
// cannot come back as a tolerated extra key.
func exactShape(decl json.RawMessage) error {
	if err := cellfill.OnlyKeys(decl, "cds.patch", "fill", "cognition"); err != nil {
		return err
	}
	if cog, ok := cellfill.Field(decl, "cognition"); ok {
		if err := cellfill.OnlyKeys(cog, "cds.patch.cognition", "provider", "model"); err != nil {
			return err
		}
	}
	return nil
}

// PatchAlpha is the provider-neutral patch-producing seat. It materializes a
// disposable worktree at the base, lets the coder change files in it, then
// MEASURES the change as a unified diff. Nothing here trusts the coder's
// account of itself: a coder that claims a sweeping refactor and wrote
// nothing produces no diff, and an episode with no diff cannot satisfy a
// contract requiring one — false completion is unrepresentable, not caught
// late.
//
// It holds no repository and no base. Both are read from the contract's
// pinned subject on every Produce, so what the episode acts on and what the
// record says it acted on are one value that was resolved once.
type PatchAlpha struct {
	coder  cellcog.Coder
	method cellmethod.View
}

func (a PatchAlpha) Produce(ctx context.Context, in cellkernel.AlphaInput) (cellkernel.AlphaOutput, error) {
	if a.coder == nil {
		return cellkernel.AlphaOutput{}, cellcog.ErrNoProvider
	}
	// AdmitSubject, not ParseSubject: a station requires a base already pinned
	// to an exact commit. Re-resolving a moving name here is precisely how the
	// two stations could measure against different trees while one record
	// claimed one. A cds.patch cell therefore cannot run without a run input —
	// an absent subject fails here rather than silently defaulting to some
	// repository nobody named.
	subject, err := cellwork.AdmitSubject(in.Contract.Subject)
	if err != nil {
		return cellkernel.AlphaOutput{}, fmt.Errorf("cds.patch: %w", err)
	}
	// Admitted here as well as at the door, and not because the door is
	// distrusted: this seat reads the FROZEN contract, so admitting the bytes it
	// was actually handed is what makes the issue it renders the issue the
	// episode recorded. It also fails before the worktree is cut.
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

	// base_sha is MEASURED from the materialized worktree — the coder never
	// sees it — so it is what the episode actually stood on, not a copy of what
	// the contract asked for. Equal to contract.subject.base_sha whenever the
	// subject is honest; the point is that a reader can compare them.
	artifacts := []cellkernel.ArtifactCandidate{
		{ID: BaseArtifactID, Kind: BaseArtifactKind, Text: wt.BaseSHA},
	}
	if strings.TrimSpace(diff) == "" {
		return cellkernel.AlphaOutput{
			Matter:    cellkernel.Matter{Data: fmt.Sprintf("no change was made to %s at %s", subject.Repo, wt.BaseSHA)},
			Artifacts: artifacts,
		}, nil
	}
	// The diff is both what beta reviews (matter) and what V checks
	// (evidence). Same bytes, two roles: beta never receives artifacts, V
	// never reads matter.
	return cellkernel.AlphaOutput{
		Matter: cellkernel.Matter{Data: diff},
		Artifacts: append(artifacts,
			cellkernel.ArtifactCandidate{ID: DiffArtifactID, Kind: DiffArtifactKind, Text: diff}),
	}, nil
}

// RenderPrompt is pure and deterministic over the contract and the CONSTRUCTIVE
// PROJECTION: the projection's text carries the exact skill bodies the cell's
// bundle loaded, so what the seat reads is the cell's methodology verbatim.
// This function no longer renders skills itself — it appends a view it was
// handed, which is what stops the seat having a second opinion about what it is
// held to.
func RenderPrompt(c cellkernel.Contract, issue cdsissue.Issue, method cellmethod.View) string {
	var b strings.Builder
	b.WriteString("You are the alpha (producing) seat of a CNOS coherence cell, working on real code.\n")
	b.WriteString("You are in a disposable worktree. Edit the files here to meet the contract.\n\n")
	fmt.Fprintf(&b, "CONTRACT %s\nGOAL: %s\n\n", c.ID, c.Goal)
	// The ISSUE, through the same cdsissue.Render the assessing seat uses. The
	// goal line is one sentence and the acceptance criteria are what the work
	// is actually judged against — a seat given only the goal writes against a
	// summary of its contract and is then reviewed against the contract. Both
	// seats read the same frozen bytes through the same function, so "they were
	// told the same thing" is a property of there being one renderer.
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
