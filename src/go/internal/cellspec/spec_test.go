package cellspec

import (
	"context"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

// The CDS-shaped fixture: fill-owned seats with holes in constructor-argument
// positions (Pi cds-fill-construction-51).
const fixture = `{
  "version": "cnos.cellspec.v0",
  "contract": {"id": "c1", "goal": "do the thing",
    "required_evidence": [
      {"id": "diff", "kind": "diff", "producer": "alpha"},
      {"id": "beta_review", "kind": "review", "producer": "beta"}
    ]},
  "protocol_id": "cnos.cdd.cds.receipt.v1",
  "params": {
    "language": {"required": true, "domain": ["cnos.eng:eng/go", "cnos.eng:eng/ocaml"]},
    "base_sha": {"required": true}
  },
  "alpha": {
    "fill": "cds.patch",
    "cognition": {"provider": "fake", "model": ""},
    "workspace": {"kind": "git-worktree", "repo": ".", "base_sha": "$base_sha"},
    "skills": ["cnos.eng:eng/code", "cnos.eng:eng/test", "$language", "cnos.eng:eng/write-functional"]
  },
  "beta": {"fill": "cdd.mechanical-unmet"}
}`

func mustParse(t *testing.T) CellSpec {
	t.Helper()
	s, err := Parse([]byte(fixture))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return s
}

// Holes resolve IN PLACE inside the seat trees: the workspace hole and the
// skill-list hole are replaced where they sit.
func TestResolveFillsHolesInPlace(t *testing.T) {
	r, err := mustParse(t).Resolve(map[string]string{"language": "cnos.eng:eng/go", "base_sha": "abc123"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	a := string(r.Alpha)
	for _, want := range []string{`"base_sha":"abc123"`, `"cnos.eng:eng/go"`, `"fill":"cds.patch"`} {
		if !strings.Contains(a, want) {
			t.Errorf("resolved alpha missing %s: %s", want, a)
		}
	}
	if strings.Contains(a, "$") {
		t.Fatalf("unresolved hole survived: %s", a)
	}
}

func TestResolveMissingRequired(t *testing.T) {
	if _, err := mustParse(t).Resolve(map[string]string{"language": "cnos.eng:eng/go"}); err == nil || !strings.Contains(err.Error(), "base_sha") {
		t.Fatalf("want missing-required error, got %v", err)
	}
}

func TestResolveDomainRejectsTypo(t *testing.T) {
	if _, err := mustParse(t).Resolve(map[string]string{"language": "cobol", "base_sha": "x"}); err == nil || !strings.Contains(err.Error(), "domain") {
		t.Fatalf("want domain error, got %v", err)
	}
}

func TestResolveUnknownParamAndHole(t *testing.T) {
	if _, err := mustParse(t).Resolve(map[string]string{"language": "cnos.eng:eng/go", "base_sha": "x", "bogus": "y"}); err == nil {
		t.Fatal("want unknown-parameter error")
	}
	undeclared := strings.Replace(fixture, `"$base_sha"`, `"$undeclared"`, 1)
	s, err := Parse([]byte(undeclared))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, err := s.Resolve(map[string]string{"language": "cnos.eng:eng/go", "base_sha": "x"}); err == nil {
		t.Fatal("a hole referencing an undeclared parameter must fail resolution")
	}
}

// TestParseRejects covers the strict generic-envelope negatives (Pi #32 D5,
// round-5 D2, round-6 D2). Seat interiors belong to fills, not to Parse.
func TestParseRejects(t *testing.T) {
	base := `"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"p","alpha":{"fill":"cdd.stub"},"beta":{"fill":"cdd.stub"}`
	cases := map[string]string{
		"unknown field":     "{" + base + `,"bogus":1}`,
		"trailing data":     "{" + base + "}{}",
		"duplicate key":     `{"version":"cnos.cellspec.v0","version":"x",` + base[len(`"version":"cnos.cellspec.v0",`):] + `}`,
		"missing version":   `{"contract":{"id":"c","goal":"g"},"protocol_id":"p","alpha":{"fill":"f"},"beta":{"fill":"f"}}`,
		"wrong version":     `{"version":"v9","contract":{"id":"c","goal":"g"},"protocol_id":"p","alpha":{"fill":"f"},"beta":{"fill":"f"}}`,
		"missing beta":      `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"p","alpha":{"fill":"f"}}`,
		"missing fill":      `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"p","alpha":{},"beta":{"fill":"f"}}`,
		"empty goal":        `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":""},"protocol_id":"p","alpha":{"fill":"f"},"beta":{"fill":"f"}}`,
		"case-alias key":    `{"version":"bad","Version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"p","alpha":{"fill":"f"},"beta":{"fill":"f"}}`,
		"case-alias nested": `{"version":"cnos.cellspec.v0","contract":{"id":"c","Goal":"g","goal":"g"},"protocol_id":"p","alpha":{"fill":"f"},"beta":{"fill":"f"}}`,
		"null anywhere":     `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"p","alpha":{"fill":"f","skills":null},"beta":{"fill":"f"}}`,
		"unknown param key": "{" + base + `,"params":{"p":{"kind":"weird"}}}`,
		"bad evidence prod": `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g","required_evidence":[{"id":"x","kind":"k","producer":"gamma"}]},"protocol_id":"p","alpha":{"fill":"f"},"beta":{"fill":"f"}}`,
		"dup evidence id":   `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g","required_evidence":[{"id":"x","kind":"k","producer":"alpha"},{"id":"x","kind":"k","producer":"beta"}]},"protocol_id":"p","alpha":{"fill":"f"},"beta":{"fill":"f"}}`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse([]byte(in)); err == nil {
				t.Fatalf("Parse(%s) = nil error, want rejection", name)
			}
		})
	}
}

// --- Build through the fill registry --------------------------------------

const stubCell = `{
  "version": "cnos.cellspec.v0",
  "contract": {"id": "cell-0", "goal": "empty",
    "required_evidence": [
      {"id": "diff", "kind": "diff", "producer": "alpha"},
      {"id": "beta_review", "kind": "review", "producer": "beta"}
    ]},
  "protocol_id": "cnos.cdd.receipt.v1",
  "alpha": {"fill": "cdd.stub"},
  "beta": {"fill": "cdd.stub"}
}`

func buildCell(t *testing.T, src string, params map[string]string) (cellkernel.Spec, cellkernel.RunMeta) {
	t.Helper()
	s, err := Parse([]byte(src))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r, err := s.Resolve(params)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	kspec, meta, err := r.Build(context.Background(), cellfill.CddFills())
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	return kspec, meta
}

// Stub fills make an honestly simulated episode; positional artifacts land on
// their own sides; the closure verifies against the caller's contract+meta.
func TestStubFillsAreSimulated(t *testing.T) {
	kspec, meta := buildCell(t, stubCell, nil)
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
	var diffAlpha, reviewBeta bool
	for _, a := range cl.Receipt.Record.Alpha.Artifacts {
		diffAlpha = diffAlpha || a.ID == "diff"
	}
	for _, a := range cl.Receipt.Record.Beta.Artifacts {
		reviewBeta = reviewBeta || a.ID == "beta_review"
	}
	if !diffAlpha || !reviewBeta {
		t.Fatal("positional artifacts not bound to their sides")
	}
}

const boolCell = `{
  "version": "cnos.cellspec.v0",
  "contract": {"id":"cell-bool","goal":"produce bool true",
    "required_evidence":[{"id":"bool","kind":"value","producer":"alpha"}]},
  "protocol_id": "cnos.cellkernel.episode-closure.v0",
  "params": {"value": {"required":true,"domain":["true","false"]}},
  "alpha": {"fill": "cdd.bool", "value": "$value"},
  "beta": {"fill": "cdd.bool-check"}
}`

// The bool cell has a genuinely mechanical review predicate, so it may accept.
func TestBoolFillAcceptedAndUnmet(t *testing.T) {
	run := func(value string) cellkernel.Status {
		kspec, meta := buildCell(t, boolCell, map[string]string{"value": value})
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

// The loader is fill-blind: an unknown fill fails at Build, before any seat
// or provider is touched, with the same corpus shape CUE rejects.
func TestUnknownFillFailsAtBuild(t *testing.T) {
	src := strings.Replace(stubCell, `"fill": "cdd.stub"}`, `"fill": "no.such"}`, 1)
	s, err := Parse([]byte(src))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r, err := s.Resolve(nil)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if _, _, err := r.Build(context.Background(), cellfill.CddFills()); err == nil || !strings.Contains(err.Error(), "no.such") {
		t.Fatalf("want unknown-fill error, got %v", err)
	}
}

// A fill's arguments are strict: an unknown key inside a seat fails that
// fill's decode even though the generic envelope cannot know the shape.
func TestFillArgumentsAreStrict(t *testing.T) {
	src := strings.Replace(stubCell, `"alpha": {"fill": "cdd.stub"}`, `"alpha": {"fill": "cdd.stub", "Extra": 1}`, 1)
	s, err := Parse([]byte(src))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r, err := s.Resolve(nil)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if _, _, err := r.Build(context.Background(), cellfill.CddFills()); err == nil {
		t.Fatal("an unknown key in a seat declaration must fail the fill decode")
	}
}

// A malformed hole must be rejected FOR ITS OWN REASON, not incidentally. The
// shared corpus cannot prove this: its fixtures also carry a base SHA no
// repository resolves, so worktree construction would return the same exit 2
// even if the hole checks disappeared entirely (Pi #58 C1). Each case here
// pins a DISTINCT diagnostic, so neither can stand in for the other.
func TestMalformedHolesRejectedForTheirOwnReason(t *testing.T) {
	// An illegal identifier is caught at Parse, as a parameter NAME, because a
	// hole spelling and a parameter name are the same grammar.
	declared := strings.Replace(fixture, `"base_sha": {"required": true}`,
		`"base_sha": {"required": true}, "bad-name": {"required": false, "default": "x"}`, 1)
	declared = strings.Replace(declared, `"$base_sha"`, `"$bad-name"`, 1)
	if _, err := Parse([]byte(declared)); err == nil {
		t.Fatal("an illegal parameter identifier must fail Parse")
	} else if !strings.Contains(err.Error(), "is not a legal identifier") {
		t.Fatalf("rejected for the wrong reason: %v", err)
	}

	// An undeclared hole survives Parse and dies at Resolve, which is the only
	// stage that knows which parameters exist.
	undeclared := strings.Replace(fixture, `"$base_sha"`, `"$nosuchparam"`, 1)
	s, err := Parse([]byte(undeclared))
	if err != nil {
		t.Fatalf("a well-formed but undeclared hole must survive Parse: %v", err)
	}
	_, err = s.Resolve(map[string]string{"language": "cnos.eng:eng/go", "base_sha": "x"})
	if err == nil {
		t.Fatal("a hole referencing an undeclared parameter must fail resolution")
	}
	if !strings.Contains(err.Error(), `undeclared parameter "nosuchparam"`) {
		t.Fatalf("rejected for the wrong reason: %v", err)
	}
}
