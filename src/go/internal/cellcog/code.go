package cellcog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

// Coder is the second port, deliberately separate from Provider: a seat that
// must CHANGE files needs a capability a text-only Provider does not have, and
// putting that difference in the type rather than in an optional field keeps
// it visible at every call site.
//
// Scope, stated honestly: a Coder is pointed at one directory and granted
// file tools only, and only that directory is ever measured. That is a
// working-directory boundary plus the provider's own workspace rules — NOT an
// OS sandbox. A provider that writes an absolute path elsewhere is not
// stopped by this package; it merely gains nothing, since anything outside
// the worktree is invisible to the diff and cannot become evidence. Real
// containment (namespace or container isolation) belongs under the adapter
// and is not implemented.
type Coder interface {
	Name() string

	// Work carries out the prompt inside dir. Its return value says only
	// whether it ran — what it actually did is measured from the worktree,
	// never taken from its word.
	Work(ctx context.Context, dir, prompt string) error
}

// Artifact ids a code episode produces. The CDS contract's canonical
// diff-first rule requires exactly {id "diff", kind "diff"}.
const (
	DiffArtifactID   = "diff"
	DiffArtifactKind = "diff"
	BaseArtifactID   = "base_sha"
	BaseArtifactKind = "sha"
)

// CodeAlpha is a producing seat that works on real code: it materializes a
// disposable worktree at a base commit, lets a rented coder change files in
// it, and then MEASURES the change as a unified diff.
//
// Nothing here trusts the coder's account of itself. The diff is computed by
// the runtime from the worktree, so a coder that claims a sweeping refactor
// and wrote nothing produces no diff — and an episode with no diff cannot
// satisfy a contract that requires one. False completion is not a review
// failure to catch later; it is unrepresentable here.
type CodeAlpha struct {
	Coder   Coder
	Repo    string // repository to cut the worktree from
	BaseRef string // revision the work starts at; the resolved SHA is bound
	Skills  []string
}

func (a CodeAlpha) Produce(ctx context.Context, in cellkernel.AlphaInput) (cellkernel.AlphaOutput, error) {
	if a.Coder == nil {
		return cellkernel.AlphaOutput{}, ErrNoProvider
	}
	wt, release, err := cellwork.Materialize(ctx, a.Repo, a.BaseRef)
	if err != nil {
		return cellkernel.AlphaOutput{}, err
	}
	defer release()

	if err := a.Coder.Work(ctx, wt.Dir, RenderCodePrompt(in.Contract, a.Skills)); err != nil {
		return cellkernel.AlphaOutput{}, fmt.Errorf("cellcog: coder %q: %w", a.Coder.Name(), err)
	}

	diff, err := wt.Diff(ctx)
	if err != nil {
		return cellkernel.AlphaOutput{}, err
	}

	// A base_sha artifact is runtime-computed, not coder-claimed: the coder
	// never sees it. It binds the episode to the exact commit the work
	// started from, even when the caller named a moving revision.
	artifacts := []cellkernel.ArtifactCandidate{
		{ID: BaseArtifactID, Kind: BaseArtifactKind, Text: wt.BaseSHA},
	}
	if strings.TrimSpace(diff) == "" {
		// No change: emit no diff. The contract that required one closes
		// needs_repair, which is the truth — the goal was not met.
		return cellkernel.AlphaOutput{
			Matter:    cellkernel.Matter{Data: fmt.Sprintf("no change was made to %s at %s", a.Repo, wt.BaseSHA)},
			Artifacts: artifacts,
		}, nil
	}

	// The diff is both what beta reviews (matter) and what V checks
	// (evidence). Same bytes, two roles: beta never receives artifacts, and V
	// never reads matter.
	return cellkernel.AlphaOutput{
		Matter: cellkernel.Matter{Data: diff},
		Artifacts: append(artifacts,
			cellkernel.ArtifactCandidate{ID: DiffArtifactID, Kind: DiffArtifactKind, Text: diff}),
	}, nil
}

// RenderCodePrompt is pure and deterministic. It asks for a change to the
// working copy — not for a report about one — because the report is not what
// gets recorded.
func RenderCodePrompt(c cellkernel.Contract, skills []string) string {
	var b strings.Builder
	b.WriteString("You are the alpha (producing) seat of a CNOS coherence cell, working on real code.\n")
	b.WriteString("You are in a disposable worktree. Edit the files here to meet the contract.\n\n")
	fmt.Fprintf(&b, "CONTRACT %s\nGOAL: %s\n", c.ID, c.Goal)
	if len(skills) > 0 {
		fmt.Fprintf(&b, "SKILLS: %s\n", strings.Join(skills, ", "))
	}
	b.WriteString("\nHOW YOUR WORK IS RECORDED\n")
	b.WriteString("Your change is measured as a unified diff of this worktree, not taken from\n")
	b.WriteString("your summary. Anything you do not write to a file does not exist. An empty\n")
	b.WriteString("diff closes the episode as unmet, whatever you say about it.\n")
	b.WriteString("An independent reviewer reads that diff afterwards.\n")
	return b.String()
}

// FakeCoder makes one deterministic, real change so CI exercises the whole
// substrate — worktree, edit, measured diff — offline and without a model.
// A run behind it is `mechanical`: nothing was rented.
type FakeCoder struct{}

func (FakeCoder) Name() string { return "fake" }

func (FakeCoder) Work(_ context.Context, dir, prompt string) error {
	note := "fake coder: deterministic change, no cognition was rented\n" +
		fmt.Sprintf("prompt bytes: %d\n", len(prompt))
	return os.WriteFile(filepath.Join(dir, "CELL-FAKE-CHANGE.txt"), []byte(note), 0o600)
}
