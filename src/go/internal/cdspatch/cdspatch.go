// Package cdspatch is the CDS-owned `cds.patch` fill: the constructor for a
// patch-producing alpha. Only this fill knows that producing a CDS patch
// takes workspace cognition, a git worktree, and loaded skills — the generic
// runner dispatches a fill id and receives a cellkernel.Alpha, nothing more.
//
// The constructor composes three reusable subsystems and owns none of their
// internals: cellcog constructs the bounded provider adapter (this package
// contains no Claude/Codex argv), cellskill resolves and loads exact skill
// bodies, and cellwork prepares the disposable worktree. What comes back is
// one immutable, provider-neutral PatchAlpha.
package cdspatch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellcog"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
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
// overlay #CDSPatchAlpha pins the same shape for the independent oracle.
type Decl struct {
	Fill      string         `json:"fill"`
	Cognition cellcog.Config `json:"cognition"`
	Workspace WorkspaceDecl  `json:"workspace"`
	Skills    []string       `json:"skills"`
}

type WorkspaceDecl struct {
	Kind    string `json:"kind"` // "git-worktree" is the only kind
	Repo    string `json:"repo"`
	BaseSHA string `json:"base_sha"`
}

// ResolvedDecl is what the closure records for this seat: the declaration
// with skills expanded to ordered canonical refs + content digests, and the
// provider/model that actually held the seat. Deterministic bytes — this is
// the canonical form the scope-lift digest covers.
type ResolvedDecl struct {
	Fill      string          `json:"fill"`
	Cognition cellcog.Config  `json:"cognition"`
	Workspace WorkspaceDecl   `json:"workspace"`
	Skills    []ResolvedSkill `json:"skills"`
}

type ResolvedSkill struct {
	Ref    string `json:"ref"`
	SHA256 string `json:"sha256"`
}

// Factory returns the cds.patch alpha factory, closed over the skill
// resolver it loads bodies from (the one construction-time effect).
func Factory(skills cellskill.Resolver) cellfill.AlphaFactory {
	return func(ctx context.Context, decl json.RawMessage) (cellfill.ConstructedAlpha, error) {
		var d Decl
		if err := cellfill.StrictDecode(decl, &d); err != nil {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: %w", Fill, err)
		}
		if d.Workspace.Kind != "git-worktree" {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: unknown workspace kind %q", Fill, d.Workspace.Kind)
		}
		if d.Workspace.Repo == "" || d.Workspace.BaseSHA == "" {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: workspace needs repo and base_sha", Fill)
		}
		if len(d.Skills) == 0 {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: a patch alpha needs its skills", Fill)
		}

		coder, mode, err := cellcog.New(d.Cognition)
		if err != nil {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: %w", Fill, err)
		}
		loaded, err := cellskill.LoadAll(skills, d.Skills)
		if err != nil {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: %w", Fill, err)
		}
		// Pin the revision now, so the recorded declaration names a commit
		// rather than a moving ref: "resolved" has to mean resolved.
		base, err := cellwork.ResolveBase(ctx, d.Workspace.Repo, d.Workspace.BaseSHA)
		if err != nil {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: %w", Fill, err)
		}

		resolved := ResolvedDecl{
			Fill:      Fill,
			Cognition: d.Cognition,
			Workspace: WorkspaceDecl{Kind: d.Workspace.Kind, Repo: d.Workspace.Repo, BaseSHA: base},
			Skills:    make([]ResolvedSkill, 0, len(loaded)),
		}
		for _, s := range loaded {
			resolved.Skills = append(resolved.Skills, ResolvedSkill{Ref: s.Ref, SHA256: s.SHA256})
		}
		canon, err := json.Marshal(resolved)
		if err != nil {
			return cellfill.ConstructedAlpha{}, fmt.Errorf("fill %q: canonicalize: %w", Fill, err)
		}

		return cellfill.ConstructedAlpha{
			Constructed: cellfill.Constructed{Decl: canon, Mode: cellkernel.ExecutionMode(mode)},
			Seat: PatchAlpha{
				coder:  coder,
				repo:   d.Workspace.Repo,
				base:   base,
				skills: loaded,
			},
		}, nil
	}
}

// PatchAlpha is the provider-neutral patch-producing seat. It materializes a
// disposable worktree at the base, lets the coder change files in it, then
// MEASURES the change as a unified diff. Nothing here trusts the coder's
// account of itself: a coder that claims a sweeping refactor and wrote
// nothing produces no diff, and an episode with no diff cannot satisfy a
// contract requiring one — false completion is unrepresentable, not caught
// late.
type PatchAlpha struct {
	coder  cellcog.Coder
	repo   string
	base   string
	skills []cellskill.Skill
}

func (a PatchAlpha) Produce(ctx context.Context, in cellkernel.AlphaInput) (cellkernel.AlphaOutput, error) {
	if a.coder == nil {
		return cellkernel.AlphaOutput{}, cellcog.ErrNoProvider
	}
	wt, release, err := cellwork.Materialize(ctx, a.repo, a.base)
	if err != nil {
		return cellkernel.AlphaOutput{}, err
	}
	defer release()

	if err := a.coder.Work(ctx, wt.Dir, RenderPrompt(in.Contract, a.skills)); err != nil {
		return cellkernel.AlphaOutput{}, fmt.Errorf("cds.patch: coder %q: %w", a.coder.Name(), err)
	}
	diff, err := wt.Diff(ctx)
	if err != nil {
		return cellkernel.AlphaOutput{}, err
	}

	// base_sha is runtime-computed — the coder never sees it — binding the
	// episode to the exact commit even when the caller named a moving ref.
	artifacts := []cellkernel.ArtifactCandidate{
		{ID: BaseArtifactID, Kind: BaseArtifactKind, Text: wt.BaseSHA},
	}
	if strings.TrimSpace(diff) == "" {
		return cellkernel.AlphaOutput{
			Matter:    cellkernel.Matter{Data: fmt.Sprintf("no change was made to %s at %s", a.repo, wt.BaseSHA)},
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

// RenderPrompt is pure and deterministic over the contract and the LOADED
// skills: the exact skill bodies are injected, not their names — naming a
// skill without its text is not loading it.
func RenderPrompt(c cellkernel.Contract, skills []cellskill.Skill) string {
	var b strings.Builder
	b.WriteString("You are the alpha (producing) seat of a CNOS coherence cell, working on real code.\n")
	b.WriteString("You are in a disposable worktree. Edit the files here to meet the contract.\n\n")
	fmt.Fprintf(&b, "CONTRACT %s\nGOAL: %s\n", c.ID, c.Goal)
	b.WriteString("\nHOW YOUR WORK IS RECORDED\n")
	b.WriteString("Your change is measured as a unified diff of this worktree, not taken from\n")
	b.WriteString("your summary. Anything you do not write to a file does not exist. An empty\n")
	b.WriteString("diff closes the episode as unmet, whatever you say about it.\n")
	b.WriteString("An independent reviewer reads that diff afterwards.\n")
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
