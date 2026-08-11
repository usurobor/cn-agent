package cdsreview

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
)

var testSkills = []string{"cnos.cdd:cdd/review", "cnos.eng:eng/go"}

func skillTree(t *testing.T, refs ...string) cellskill.Tree {
	t.Helper()
	root := t.TempDir()
	for _, ref := range refs {
		pkg, path, _ := strings.Cut(ref, ":")
		dir := filepath.Join(root, pkg, "skills", filepath.FromSlash(path))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# body of "+ref+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return cellskill.Tree{Root: root}
}

func declJSON(provider string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(
		`{"fill":"cds.review","cognition":{"provider":%q,"model":""},`+
			`"skills":["cnos.cdd:cdd/review","cnos.eng:eng/go"]}`, provider))
}

var reviewContract = cellkernel.Contract{ID: "c", Goal: "add a CONTRIBUTING guide"}

// stubAnswerer returns exactly what a test wants back, so the decode boundary
// can be driven without renting anything.
type stubAnswerer struct {
	raw string
	err error
}

func (stubAnswerer) Name() string { return "stub" }
func (s stubAnswerer) Answer(context.Context, string, json.RawMessage) (json.RawMessage, error) {
	if s.err != nil {
		return nil, s.err
	}
	return json.RawMessage(s.raw), nil
}

// The reviewing seat holds NO workspace: not a repo, not a directory, not a
// worktree. Independence is structural — it cannot reach alpha's workspace
// because it was never given one and has no tool to look with.
func TestSeatHoldsNoWorkspace(t *testing.T) {
	c, err := Factory(skillTree(t, testSkills...))(context.Background(), declJSON("fake"))
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	seat, ok := c.Seat.(ReviewBeta)
	if !ok {
		t.Fatalf("seat is %T, want ReviewBeta", c.Seat)
	}
	// Whole-value check rather than field-by-field: a workspace added later in
	// any form fails here without anyone remembering to extend the list.
	if got := fmt.Sprintf("%+v", struct {
		Answerer any
		Skills   int
	}{seat.answerer, len(seat.skills)}); !strings.Contains(got, "Skills:2") {
		t.Fatalf("unexpected seat shape: %s", got)
	}
	var decl map[string]json.RawMessage
	if err := json.Unmarshal(c.Decl, &decl); err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"workspace", "repo", "base_sha"} {
		if _, present := decl[forbidden]; present {
			t.Errorf("resolved reviewer declaration must not carry %q: %s", forbidden, c.Decl)
		}
	}
}

// A workspace key is rejected FOR BEING A WORKSPACE KEY. The shared corpus
// cannot witness this: without an installed hub every spec fails earlier on
// skill loading, so an exit code alone would prove nothing.
func TestWorkspaceKeyRejectedForItsOwnReason(t *testing.T) {
	decl := json.RawMessage(`{"fill":"cds.review","cognition":{"provider":"fake","model":""},` +
		`"workspace":{"kind":"git-worktree","repo":".","base_sha":"x"},` +
		`"skills":["cnos.cdd:cdd/review"]}`)
	_, err := Factory(skillTree(t, testSkills...))(context.Background(), decl)
	if err == nil {
		t.Fatal("a reviewer declaring a workspace must fail construction")
	}
	if !strings.Contains(err.Error(), "workspace") {
		t.Fatalf("rejected for the wrong reason: %v", err)
	}
}

// A verdict that cannot be read is NOT a failing verdict. Defaulting an
// unreadable answer to pass:false would fabricate a judgement nobody made —
// the same sin as fabricating completion, from the other direction.
func TestUnreadableVerdictFailsRatherThanDefaulting(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{name: "unknown field", raw: `{"pass":false,"notes":"n","verdict":"extra"}`, want: "schema"},
		{name: "wrong type", raw: `{"pass":"no","notes":"n"}`, want: "schema"},
		{name: "not an object", raw: `"fail"`, want: "schema"},
		{name: "empty notes", raw: `{"pass":false,"notes":"   "}`, want: "no notes"},
		{name: "passing with no reason", raw: `{"pass":true,"notes":""}`, want: "no notes"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			seat := ReviewBeta{answerer: stubAnswerer{raw: tc.raw}}
			_, err := seat.Review(context.Background(), cellkernel.BetaInput{
				Contract: reviewContract,
				Matter:   cellkernel.Matter{Data: "diff --git a/x b/x\n"},
			})
			if err == nil {
				t.Fatalf("verdict %s must fail, not become a review", tc.raw)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("failed for the wrong reason: got %v, want mention of %q", err, tc.want)
			}
		})
	}
}

// A well-formed verdict passes through unchanged, in both directions, and
// carries no artifacts — evidence is gamma/V's channel, never beta's.
func TestWellFormedVerdictIsCarried(t *testing.T) {
	for _, pass := range []bool{true, false} {
		raw := fmt.Sprintf(`{"pass":%t,"notes":"because"}`, pass)
		out, err := ReviewBeta{answerer: stubAnswerer{raw: raw}}.Review(
			context.Background(),
			cellkernel.BetaInput{Contract: reviewContract, Matter: cellkernel.Matter{Data: "d"}})
		if err != nil {
			t.Fatalf("%s: %v", raw, err)
		}
		if out.Review.Pass != pass || out.Review.Notes != "because" {
			t.Fatalf("verdict altered in transit: %+v", out.Review)
		}
		if len(out.Artifacts) != 0 {
			t.Fatalf("a reviewing seat must emit no artifacts, got %d", len(out.Artifacts))
		}
	}
}

// The prompt IS the seat's whole world. If the contract goal or the matter
// were missing from it, the reviewer would be judging something else.
func TestPromptCarriesContractAndMatter(t *testing.T) {
	skills := []cellskill.Skill{{Ref: "cnos.cdd:cdd/review", Body: "REVIEW-SKILL-BODY"}}
	got := RenderPrompt(reviewContract, cellkernel.Matter{Data: "THE-MATTER"}, skills)
	for _, want := range []string{"add a CONTRIBUTING guide", "THE-MATTER", "REVIEW-SKILL-BODY", "cnos.cdd:cdd/review"} {
		if !strings.Contains(got, want) {
			t.Errorf("prompt is missing %q", want)
		}
	}
}

// Construction fails closed, and the resolved declaration records digests —
// naming a skill is not loading it.
func TestConstructionFailsClosedAndRecordsDigests(t *testing.T) {
	tree := skillTree(t, testSkills...)
	for _, tc := range []struct {
		name string
		decl string
		want string
	}{
		{"no skills", `{"fill":"cds.review","cognition":{"provider":"fake","model":""},"skills":[]}`, "needs its skills"},
		{"unknown provider", `{"fill":"cds.review","cognition":{"provider":"clyde","model":"m"},"skills":["cnos.cdd:cdd/review"]}`, "unknown provider"},
		{"fake with a model", `{"fill":"cds.review","cognition":{"provider":"fake","model":"m"},"skills":["cnos.cdd:cdd/review"]}`, "takes no model"},
		{"uninstalled skill", `{"fill":"cds.review","cognition":{"provider":"fake","model":""},"skills":["nope:x"]}`, "not installed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Factory(tree)(context.Background(), json.RawMessage(tc.decl)); err == nil {
				t.Fatal("must fail construction")
			} else if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("wrong reason: %v", err)
			}
		})
	}

	c, err := Factory(tree)(context.Background(), declJSON("fake"))
	if err != nil {
		t.Fatal(err)
	}
	var rd ResolvedDecl
	if err := json.Unmarshal(c.Decl, &rd); err != nil {
		t.Fatal(err)
	}
	if len(rd.Skills) != 2 {
		t.Fatalf("want 2 resolved skills, got %d", len(rd.Skills))
	}
	for i, s := range rd.Skills {
		if s.Ref != testSkills[i] {
			t.Errorf("skill order changed: got %q want %q", s.Ref, testSkills[i])
		}
		if len(s.SHA256) != 64 {
			t.Errorf("skill %q carries no content digest", s.Ref)
		}
	}
	if c.Mode != cellkernel.ExecutionMode("mechanical") {
		t.Errorf("a fake reviewer must be mechanical, got %q", c.Mode)
	}
}

// The deterministic fake never passes. A fake that returned pass:true would be
// false completion wearing a reviewer's clothes.
func TestFakeReviewerNeverPasses(t *testing.T) {
	c, err := Factory(skillTree(t, testSkills...))(context.Background(), declJSON("fake"))
	if err != nil {
		t.Fatal(err)
	}
	out, err := c.Seat.Review(context.Background(), cellkernel.BetaInput{
		Contract: reviewContract,
		Matter:   cellkernel.Matter{Data: "a very convincing diff"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Review.Pass {
		t.Fatal("the deterministic fake must never pass a review it did not perform")
	}
}

var _ cellfill.BetaFactory = Factory(cellskill.Tree{})
