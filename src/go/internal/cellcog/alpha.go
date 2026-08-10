package cellcog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

// ResponseSchema is the envelope a rented alpha must answer with. It is
// pinned so the prompt and the parser cannot drift apart.
const ResponseSchema = `{"matter": "<what you produced, as text>", "artifacts": [{"id": "<id>", "kind": "<kind>", "text": "<content>"}]}`

// Alpha is a rented producing seat: it renders the frozen contract into a
// prompt, asks a provider, and parses the answer into kernel candidates.
type Alpha struct {
	Provider Provider
	Skills   []string
}

// Produce implements cellkernel.Alpha. A provider failure or an unparseable
// answer is an error — the seat produced nothing, so there is no matter to
// review and the episode does not close. Retrying is a repair driver's job
// (Case 4), deliberately not this seat's.
func (a Alpha) Produce(ctx context.Context, in cellkernel.AlphaInput) (cellkernel.AlphaOutput, error) {
	if a.Provider == nil {
		return cellkernel.AlphaOutput{}, ErrNoProvider
	}
	raw, err := a.Provider.Complete(ctx, RenderAlphaPrompt(in.Contract, a.Skills))
	if err != nil {
		return cellkernel.AlphaOutput{}, fmt.Errorf("cellcog: provider %q: %w", a.Provider.Name(), err)
	}
	out, err := ParseAlphaResponse([]byte(raw))
	if err != nil {
		return cellkernel.AlphaOutput{}, fmt.Errorf("cellcog: provider %q: %w", a.Provider.Name(), err)
	}
	return out, nil
}

// RenderAlphaPrompt is pure and deterministic: the same contract and skills
// always render the same prompt, so a run is reproducible up to the provider.
func RenderAlphaPrompt(c cellkernel.Contract, skills []string) string {
	var b strings.Builder
	b.WriteString("You are the alpha (producing) seat of a CNOS coherence cell.\n")
	b.WriteString("You do the work; an independent reviewer and a mechanical verifier judge it.\n\n")
	fmt.Fprintf(&b, "CONTRACT %s\nGOAL: %s\n", c.ID, c.Goal)

	if len(skills) > 0 {
		fmt.Fprintf(&b, "SKILLS: %s\n", strings.Join(skills, ", "))
	}

	required := requiredFrom(c, cellkernel.RoleAlpha)
	if len(required) > 0 {
		b.WriteString("\nREQUIRED EVIDENCE — your answer must carry one artifact per line below,\nwith exactly these ids and kinds:\n")
		for _, r := range required {
			fmt.Fprintf(&b, "  - id %q, kind %q\n", r.ID, r.Kind)
		}
	} else {
		b.WriteString("\nREQUIRED EVIDENCE: none; return an empty artifacts list.\n")
	}

	b.WriteString("\nANSWER FORMAT — return one JSON object and nothing else:\n")
	b.WriteString(ResponseSchema)
	b.WriteString("\n\nNo prose outside the JSON. Every field is required.\n")
	b.WriteString("`matter` is what a reviewer reads to judge whether the goal was met.\n")
	return b.String()
}

func requiredFrom(c cellkernel.Contract, role cellkernel.Role) []cellkernel.RequiredRef {
	out := make([]cellkernel.RequiredRef, 0, len(c.RequiredEvidence))
	for _, r := range c.RequiredEvidence {
		if r.Producer == role {
			out = append(out, r)
		}
	}
	return out
}

// alphaResponse mirrors ResponseSchema.
type alphaResponse struct {
	Matter    string             `json:"matter"`
	Artifacts []responseArtifact `json:"artifacts"`
}

type responseArtifact struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Text string `json:"text"`
}

// ParseAlphaResponse turns a provider's raw answer into kernel candidates.
// Pure: bytes in, data out (eng/go §2.17).
//
// This parser is deliberately less paranoid than cellspec's: a cell spec is
// an authority document where a smuggled key could launder an invocation
// claim, whereas this is untrusted candidate data that carries no authority
// at all — the kernel bounds and UTF-8-checks it at sealing, and V judges
// evidence positionally. Strictness here only has to stop a malformed answer
// from being mistaken for a well-formed one.
func ParseAlphaResponse(data []byte) (cellkernel.AlphaOutput, error) {
	dec := json.NewDecoder(bytes.NewReader(unfence(data)))
	dec.DisallowUnknownFields()
	var r alphaResponse
	if err := dec.Decode(&r); err != nil {
		return cellkernel.AlphaOutput{}, fmt.Errorf("provider answer is not the required JSON envelope: %w", err)
	}
	var extra json.RawMessage
	if err := dec.Decode(&extra); err != io.EOF {
		return cellkernel.AlphaOutput{}, fmt.Errorf("provider answer carried trailing data after the JSON object")
	}
	if strings.TrimSpace(r.Matter) == "" {
		return cellkernel.AlphaOutput{}, fmt.Errorf("provider answer has empty matter")
	}

	cands := make([]cellkernel.ArtifactCandidate, 0, len(r.Artifacts))
	for _, a := range r.Artifacts {
		cands = append(cands, cellkernel.ArtifactCandidate{ID: a.ID, Kind: a.Kind, Text: a.Text})
	}
	return cellkernel.AlphaOutput{Matter: cellkernel.Matter{Data: r.Matter}, Artifacts: cands}, nil
}

// unfence strips one markdown code fence around the answer. This is the ONE
// tolerated normalization: models routinely fence JSON they were asked to
// return bare, and rejecting that would fail the seam for a formatting habit
// rather than a content fault. Anything beyond a single fence is a malformed
// answer, not something to guess at.
func unfence(data []byte) []byte {
	t := bytes.TrimSpace(data)
	if !bytes.HasPrefix(t, []byte("```")) {
		return t
	}
	if i := bytes.IndexByte(t, '\n'); i >= 0 {
		t = t[i+1:]
	}
	return bytes.TrimSpace(bytes.TrimSuffix(bytes.TrimSpace(t), []byte("```")))
}

// MatterBeta is the mechanical reviewer paired with a rented alpha while beta
// is still mechanical (Case 2). It re-reads the matter projection and passes
// iff it carries content.
//
// This is a deliberately weak review, and says so in its notes: a mechanical
// seat cannot judge arbitrary prose. It is not tautological — an alpha that
// answers with whitespace fails — but real judgement arrives with a rented
// beta (Case 3). Required evidence is never beta's business: V checks it.
type MatterBeta struct{}

func (MatterBeta) Review(_ context.Context, in cellkernel.BetaInput) (cellkernel.BetaOutput, error) {
	if strings.TrimSpace(in.Matter.Data) == "" {
		return cellkernel.BetaOutput{Review: cellkernel.Review{
			Pass:  false,
			Notes: "mechanical review: matter is empty, the goal cannot have been met",
		}}, nil
	}
	return cellkernel.BetaOutput{Review: cellkernel.Review{
		Pass: true,
		Notes: fmt.Sprintf("mechanical review only: matter carries %d bytes of content; "+
			"a mechanical seat cannot judge whether it meets the goal", len(in.Matter.Data)),
	}}, nil
}
