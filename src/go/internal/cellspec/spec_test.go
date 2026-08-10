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
		"missing beta":      `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"cnos.cdd.receipt.v1","alpha":{"skills":[]}}`,
		"unknown profile":   "{" + base + `,"profile":"wat"}`,
		"bad param kind":    "{" + base + `,"params":{"p":{"kind":"weird"}}}`,
		"bad evidence prod": `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g","required_evidence":[{"id":"x","kind":"k","producer":"gamma"}]},"protocol_id":"cnos.cdd.receipt.v1","alpha":{"skills":[]},"beta":{"skills":[]}}`,
		"dup evidence id":   `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g","required_evidence":[{"id":"x","kind":"k","producer":"alpha"},{"id":"x","kind":"k","producer":"beta"}]},"protocol_id":"cnos.cdd.receipt.v1","alpha":{"skills":[]},"beta":{"skills":[]}}`,
		"empty goal":        `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":""},"protocol_id":"cnos.cdd.receipt.v1","profile":"stub","alpha":{"skills":[]},"beta":{"skills":[]}}`,
		"missing skills":    `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"cnos.cdd.receipt.v1","profile":"stub","alpha":{},"beta":{"skills":[]}}`,
		"case-alias key":    `{"version":"bad","Version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"cnos.cdd.receipt.v1","profile":"stub","alpha":{"skills":[]},"beta":{"skills":[]}}`,
		"case-alias nested": `{"version":"cnos.cellspec.v0","contract":{"id":"c","Goal":"g","goal":"g"},"protocol_id":"cnos.cdd.receipt.v1","profile":"stub","alpha":{"skills":[]},"beta":{"skills":[]}}`,
		"null skills":       `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"cnos.cdd.receipt.v1","profile":"stub","alpha":{"skills":null},"beta":{"skills":[]}}`,
		"null evidence":     `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g","required_evidence":null},"protocol_id":"cnos.cdd.receipt.v1","profile":"stub","alpha":{"skills":[]},"beta":{"skills":[]}}`,
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
	cl, err := cellkernel.RunEpisode(context.Background(), kspec, meta)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cl.Status != cellkernel.Simulated {
		t.Fatalf("status = %q, want simulated", cl.Status)
	}
	if err := cellkernel.VerifyClosure(kspec.Contract, meta, cl); err != nil {
		t.Fatalf("closure must verify: %v", err)
	}
	// Positional ownership: the diff sits under Alpha, beta_review under Beta.
	var diffAlpha, reviewBeta bool
	for _, a := range cl.Receipt.Record.Alpha.Artifacts {
		if a.ID == "diff" {
			diffAlpha = true
		}
	}
	for _, a := range cl.Receipt.Record.Beta.Artifacts {
		if a.ID == "beta_review" {
			reviewBeta = true
		}
	}
	if !diffAlpha || !reviewBeta {
		t.Fatalf("positional artifacts not bound: alpha=%+v beta=%+v",
			cl.Receipt.Record.Alpha.Artifacts, cl.Receipt.Record.Beta.Artifacts)
	}
}

const boolFixture = `{
  "version": "cnos.cellspec.v0",
  "contract": {"id":"cell-bool","goal":"produce bool true",
    "required_evidence":[{"id":"bool","kind":"value","producer":"alpha"}]},
  "protocol_id": "cnos.cellkernel.episode-closure.v0",
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
		cl, err := cellkernel.RunEpisode(context.Background(), kspec, meta)
		if err != nil {
			t.Fatalf("run: %v", err)
		}
		return cl.Status
	}
	if got := run("true"); got != cellkernel.Accepted {
		t.Errorf("value=true: status %q, want accepted", got)
	}
	if got := run("false"); got != cellkernel.NeedsRepair {
		t.Errorf("value=false: status %q, want needs_repair", got)
	}
}

// --- Cognitive profile (Phase 3, Case 2) ----------------------------------

const cognitiveFixture = `{
  "version": "cnos.cellspec.v0",
  "contract": {"id":"cell-cog","goal":"answer the question"},
  "protocol_id": "cnos.cellkernel.episode-closure.v0",
  "profile": "cognitive",
  "params": {"provider": {"kind":"value","required":true,"domain":["claude","fake"]}},
  "alpha": {"skills": ["eng"]},
  "beta": {"skills": []}
}`

// The provider decides the mode, because the mode must tell the truth about
// how the work was produced: only a provider that really rents cognition may
// run `cognitive`; the deterministic fake is `mechanical`.
func TestCognitiveProfileModeFollowsProvider(t *testing.T) {
	s, err := Parse([]byte(cognitiveFixture))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	for provider, want := range map[string]cellkernel.ExecutionMode{
		"fake":   cellkernel.ModeMechanical,
		"claude": cellkernel.ModeCognitive,
	} {
		t.Run(provider, func(t *testing.T) {
			r, err := s.Resolve(map[string]string{"provider": provider})
			if err != nil {
				t.Fatalf("resolve: %v", err)
			}
			_, meta, err := r.Build()
			if err != nil {
				t.Fatalf("build: %v", err)
			}
			if meta.ExecutionMode != want {
				t.Fatalf("provider %q: mode = %q, want %q", provider, meta.ExecutionMode, want)
			}
			if meta.ResolvedSpec.Params["provider"] != provider {
				t.Fatal("the provider that held the seat is not disclosed in the record")
			}
		})
	}
}

// A cognitive episode closes through the real loader path and self-verifies.
func TestCognitiveProfileClosesWithFake(t *testing.T) {
	s, err := Parse([]byte(cognitiveFixture))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r, err := s.Resolve(map[string]string{"provider": "fake"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	kspec, meta, err := r.Build()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	cl, err := cellkernel.RunEpisode(context.Background(), kspec, meta)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cl.Status != cellkernel.Accepted {
		t.Fatalf("status = %q, want accepted", cl.Status)
	}
	if err := cellkernel.VerifyClosure(kspec.Contract, meta, cl); err != nil {
		t.Fatalf("closure must verify: %v", err)
	}
}

func TestCognitiveProfileRejectsBadProvider(t *testing.T) {
	s, err := Parse([]byte(cognitiveFixture))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	// A typo must fail resolution against the closed domain...
	if _, err := s.Resolve(map[string]string{"provider": "clyde"}); err == nil {
		t.Fatal("provider outside the declared domain must fail resolution")
	}
	// ...and an undeclared provider must fail the build even if a spec widens
	// its own domain.
	wide := strings.Replace(cognitiveFixture, `["claude","fake"]`, `["clyde"]`, 1)
	ws, err := Parse([]byte(wide))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r, err := ws.Resolve(map[string]string{"provider": "clyde"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if _, _, err := r.Build(); err == nil {
		t.Fatal("unknown provider must fail the build")
	}
}

// The profile must declare the hole that steers it.
func TestCognitiveProfileRequiresProviderParam(t *testing.T) {
	noParam := `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},` +
		`"protocol_id":"p","profile":"cognitive","alpha":{"skills":[]},"beta":{"skills":[]}}`
	if _, err := Parse([]byte(noParam)); err == nil {
		t.Fatal("cognitive profile without a provider parameter must be rejected")
	}
}
