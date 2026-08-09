package cellspec

import (
	"context"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

const fixture = `{
  "version": "cnos.cellspec.v0",
  "contract": {"id": "c1", "goal": "do the thing",
    "required_evidence": [
      {"id": "diff", "kind": "diff", "producer": "alpha"},
      {"id": "beta_review", "kind": "review", "producer": "beta"}
    ]},
  "protocol_id": "cnos.cdd.cds.receipt.v1",
  "profile": "stub",
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
}

func TestResolveMissingRequired(t *testing.T) {
	if _, err := mustParse(t).Resolve(nil); err == nil || !strings.Contains(err.Error(), "language") {
		t.Fatalf("want missing-required error, got %v", err)
	}
}

func TestResolveDomainRejectsTypo(t *testing.T) {
	if _, err := mustParse(t).Resolve(map[string]string{"language": "cobol"}); err == nil || !strings.Contains(err.Error(), "domain") {
		t.Fatalf("want domain error, got %v", err)
	}
}

func TestResolveUnknownParam(t *testing.T) {
	if _, err := mustParse(t).Resolve(map[string]string{"language": "go", "bogus": "x"}); err == nil || !strings.Contains(err.Error(), "unknown parameter") {
		t.Fatalf("want unknown-parameter error, got %v", err)
	}
}

// TestParseRejects covers the strict-parse negatives (Pi #32 D5).
func TestParseRejects(t *testing.T) {
	base := `"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"cnos.cdd.receipt.v1","alpha":{"skills":[]},"beta":{"skills":[]}`
	cases := map[string]string{
		"unknown field":     "{" + base + `,"bogus":1}`,
		"trailing data":     "{" + base + "}{}",
		"duplicate key":     `{"version":"cnos.cellspec.v0","version":"x","contract":{"id":"c","goal":"g"},"protocol_id":"cnos.cdd.receipt.v1","alpha":{"skills":[]},"beta":{"skills":[]}}`,
		"missing version":   `{"contract":{"id":"c","goal":"g"},"protocol_id":"cnos.cdd.receipt.v1","alpha":{"skills":[]},"beta":{"skills":[]}}`,
		"wrong version":     `{"version":"v9","contract":{"id":"c","goal":"g"},"protocol_id":"cnos.cdd.receipt.v1","alpha":{"skills":[]},"beta":{"skills":[]}}`,
		"unknown protocol":  "{" + `"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"made.up","alpha":{"skills":[]},"beta":{"skills":[]}` + "}",
		"missing beta":      `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"cnos.cdd.receipt.v1","alpha":{"skills":[]}}`,
		"unknown profile":   "{" + base + `,"profile":"wat"}`,
		"bad param kind":    "{" + base + `,"params":{"p":{"kind":"weird"}}}`,
		"bad evidence prod": `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g","required_evidence":[{"id":"x","kind":"k","producer":"gamma"}]},"protocol_id":"cnos.cdd.receipt.v1","alpha":{"skills":[]},"beta":{"skills":[]}}`,
		"dup evidence id":   `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g","required_evidence":[{"id":"x","kind":"k","producer":"alpha"},{"id":"x","kind":"k","producer":"beta"}]},"protocol_id":"cnos.cdd.receipt.v1","alpha":{"skills":[]},"beta":{"skills":[]}}`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse([]byte(in)); err == nil {
				t.Fatalf("Parse(%s) = nil error, want rejection", name)
			}
		})
	}
}

func TestStubProfileIsSimulated(t *testing.T) {
	r, err := mustParse(t).Resolve(map[string]string{"language": "go"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	kspec, meta, err := r.Build()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if meta.ExecutionMode != cellkernel.ModeStub {
		t.Fatalf("mode = %q, want stub", meta.ExecutionMode)
	}
	env, err := cellkernel.RunEpisode(context.Background(), kspec, cellkernel.WithMeta(meta))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if env.Status != cellkernel.Simulated {
		t.Fatalf("status = %q, want simulated", env.Status)
	}
	if err := cellkernel.VerifyEnvelope(env); err != nil {
		t.Fatalf("envelope must verify: %v", err)
	}
	// The α diff and β beta_review are each bound with the authorized producer.
	var diffAlpha, reviewBeta bool
	for _, e := range env.Receipt.Evidence {
		if e.ID == "diff" && e.Producer == cellkernel.RoleAlpha {
			diffAlpha = true
		}
		if e.ID == "beta_review" && e.Producer == cellkernel.RoleBeta {
			reviewBeta = true
		}
	}
	if !diffAlpha || !reviewBeta {
		t.Fatalf("producer-attributed evidence not bound: %+v", env.Receipt.Evidence)
	}
}

const boolFixture = `{
  "version": "cnos.cellspec.v0",
  "contract": {"id":"cell-bool","goal":"produce bool true",
    "required_evidence":[{"id":"bool","kind":"value","producer":"alpha"}]},
  "protocol_id": "cnos.cellkernel.episode-receipt.v0",
  "profile": "bool",
  "params": {"value": {"kind":"value","required":true,"domain":["true","false"]}},
  "alpha": {"skills": []},
  "beta": {"skills": []}
}`

func TestBoolProfileAcceptedAndUnmet(t *testing.T) {
	s, err := Parse([]byte(boolFixture))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	run := func(value string) cellkernel.Status {
		r, err := s.Resolve(map[string]string{"value": value})
		if err != nil {
			t.Fatalf("resolve(%s): %v", value, err)
		}
		kspec, meta, err := r.Build()
		if err != nil {
			t.Fatalf("build: %v", err)
		}
		if meta.ExecutionMode != cellkernel.ModeMechanical {
			t.Fatalf("mode = %q, want mechanical", meta.ExecutionMode)
		}
		env, err := cellkernel.RunEpisode(context.Background(), kspec, cellkernel.WithMeta(meta))
		if err != nil {
			t.Fatalf("run: %v", err)
		}
		return env.Status
	}
	if got := run("true"); got != cellkernel.Accepted {
		t.Errorf("value=true: status %q, want accepted", got)
	}
	if got := run("false"); got != cellkernel.NeedsRepair {
		t.Errorf("value=false: status %q, want needs_repair", got)
	}
}
