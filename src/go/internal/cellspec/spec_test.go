package cellspec

import (
	"context"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

const fixture = `{
  "contract": {"id": "c1", "goal": "do the thing",
    "required_evidence": [{"id": "diff", "kind": "diff"}]},
  "protocol_id": "cnos.cdd.cds.receipt.v1",
  "params": {
    "language": {"kind": "skill", "required": true, "domain": ["go", "ocaml"]},
    "style": {"kind": "skill", "required": false, "default": "functional"}
  },
  "alpha": {"skills": ["eng", "$language", "$style"]},
  "beta": {"skills": ["$language", "cds-review"]}
}`

func mustParse(t *testing.T) CellSpec {
	t.Helper()
	s, err := Parse([]byte(fixture))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return s
}

func TestResolveFillsAndSplices(t *testing.T) {
	r, err := mustParse(t).Resolve(map[string]string{"language": "go"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if got := strings.Join(r.AlphaSkills, ","); got != "eng,go,functional" {
		t.Errorf("alpha skills = %q, want eng,go,functional", got)
	}
	if got := strings.Join(r.BetaSkills, ","); got != "go,cds-review" {
		t.Errorf("beta skills = %q, want go,cds-review", got)
	}
	if r.Params["style"] != "functional" {
		t.Errorf("style default not applied: %q", r.Params["style"])
	}
}

func TestResolveMissingRequired(t *testing.T) {
	_, err := mustParse(t).Resolve(nil)
	if err == nil || !strings.Contains(err.Error(), "language") {
		t.Fatalf("want missing-required error naming language, got %v", err)
	}
}

func TestResolveDomainRejectsTypo(t *testing.T) {
	_, err := mustParse(t).Resolve(map[string]string{"language": "cobol"})
	if err == nil || !strings.Contains(err.Error(), "domain") {
		t.Fatalf("want domain error, got %v", err)
	}
}

func TestResolveUnknownParam(t *testing.T) {
	_, err := mustParse(t).Resolve(map[string]string{"language": "go", "bogus": "x"})
	if err == nil || !strings.Contains(err.Error(), "unknown parameter") {
		t.Fatalf("want unknown-parameter error, got %v", err)
	}
}

func TestOptionalUnfilledSkillDropped(t *testing.T) {
	// A spec whose optional param has no default: the $ref is dropped, not error.
	s := mustParse(t)
	s.Params["style"] = ParamSpec{Kind: "skill"} // optional, no default
	r, err := s.Resolve(map[string]string{"language": "ocaml"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if got := strings.Join(r.AlphaSkills, ","); got != "eng,ocaml" {
		t.Errorf("alpha skills = %q, want eng,ocaml (style dropped)", got)
	}
}

func TestParseRejectsUnknownField(t *testing.T) {
	_, err := Parse([]byte(`{"contract":{"id":"c","goal":"g"},"protocol_id":"p","alpha":{"skills":[]},"beta":{"skills":[]},"bogus":1}`))
	if err == nil {
		t.Fatal("want error for unknown top-level field")
	}
}

func TestKernelSpecRunsToAccepted(t *testing.T) {
	r, err := mustParse(t).Resolve(map[string]string{"language": "go"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	res, err := cellkernel.RunEpisode(context.Background(), r.KernelSpec())
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Status != cellkernel.Accepted {
		t.Fatalf("status = %q, want accepted", res.Status)
	}
	// The required diff evidence was bound by the stub α.
	found := false
	for _, e := range res.Receipt.EvidenceRefs {
		if e.ID == "diff" {
			found = true
		}
	}
	if !found {
		t.Fatal("required diff evidence not bound")
	}
}
