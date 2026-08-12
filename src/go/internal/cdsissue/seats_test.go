// This is an EXTERNAL test package on purpose. The property under test spans
// both CDS seats, and both of them import cdsissue — an internal test could
// not reach them without an import cycle, and putting the test in one seat's
// package would make it look like that seat's property. It is neither seat's:
// it is the issue's.
package cdsissue_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cdsissue"
	"github.com/usurobor/cnos/src/go/internal/cdspatch"
	"github.com/usurobor/cnos/src/go/internal/cdsreview"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

const oneIssue = `{
	"kind": "cnos.cds.issue.v0",
	"id": "one-issue",
	"problem": {
		"exists": "PROBLEM-EXISTS",
		"expected": "PROBLEM-EXPECTED",
		"diverges": "PROBLEM-DIVERGES"
	},
	"sources": [{"claim": "SOURCE-CLAIM", "path": "SOURCE-PATH"}],
	"scope": {"in": ["SCOPE-IN"], "out": ["SCOPE-OUT"]},
	"acceptance": [
		{"id": "AC1", "statement": "STATEMENT-ONE", "verification": "VERIFICATION-ONE"},
		{"id": "AC2", "statement": "STATEMENT-TWO", "verification": "VERIFICATION-TWO"}
	]
}`

// AC4: one contract, two prompts, one issue block — byte for byte.
//
// Sameness is structural: both seats receive frozen copies of one contract and
// both call cdsissue.Render. This test is what makes that claim falsifiable —
// a seat that formatted the issue itself would drift from the other the first
// time either was edited, and the drift would show up as a reviewer judging
// against criteria alpha was never given.
func TestBothSeatsAreShownTheSameIssueBlock(t *testing.T) {
	contract := cellkernel.Contract{
		ID:   "cds-one",
		Goal: "the goal line, which is not the specification",
		Task: json.RawMessage(oneIssue),
	}
	issue, err := cdsissue.Admit(contract.Task)
	if err != nil {
		t.Fatalf("admit: %v", err)
	}
	skills := []cellskill.Skill{{Ref: "cnos.eng:eng/go", Body: "SKILL-BODY"}}

	alpha := cdspatch.RenderPrompt(contract, issue, skills)
	beta := cdsreview.RenderPrompt(contract, issue, cellkernel.Matter{Data: "diff --git a/x b/x\n"},
		cellwork.View{Files: []cellwork.FileState{{Path: "x", Status: cellwork.FileAdded, Content: "x\n"}}}, skills)

	block := cdsissue.Render(issue)
	if strings.Count(alpha, block) != 1 || strings.Count(beta, block) != 1 {
		t.Fatalf("the rendered issue does not appear exactly once in each prompt (alpha %d, beta %d)",
			strings.Count(alpha, block), strings.Count(beta, block))
	}
	// Compare the two prompts to EACH OTHER, not just each to the renderer:
	// the claim is about what the seats see.
	ai, bi := strings.Index(alpha, block), strings.Index(beta, block)
	if alpha[ai:ai+len(block)] != beta[bi:bi+len(block)] {
		t.Fatal("the issue blocks the two seats see are not byte-identical")
	}

	// Every criterion's verification route reaches both seats. This is the
	// property that makes verification simpler than production: beta is told
	// exactly what alpha was told it must satisfy.
	for _, c := range issue.Acceptance {
		for name, prompt := range map[string]string{"alpha": alpha, "beta": beta} {
			for _, want := range []string{c.ID, c.Statement, c.Verification} {
				if !strings.Contains(prompt, want) {
					t.Errorf("%s prompt is missing %q", name, want)
				}
			}
		}
	}
}
