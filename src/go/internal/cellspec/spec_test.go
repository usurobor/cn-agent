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

// A hole is rejected for ITS OWN fact. Malformed and undeclared are different
// facts, and the earlier version of this test proved neither: it inserted an
// illegal DECLARED PARAMETER KEY, so Parse rejected the key before the seat
// value was ever examined, and the surviving case covered only the
// well-formed undeclared spelling (Pi #59 C1). These drive `spliceValue`
// directly, through a spec whose declarations are all legal.
func TestHolesRejectedForTheirOwnFact(t *testing.T) {
	cases := []struct {
		name string
		hole string
		want string
	}{
		{name: "malformed", hole: `"$bad-name"`, want: `hole "$bad-name" is malformed`},
		{name: "undeclared", hole: `"$nosuchparam"`, want: `undeclared parameter "nosuchparam"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := strings.Replace(fixture, `"$base_sha"`, tc.hole, 1)
			s, err := Parse([]byte(src))
			if err != nil {
				t.Fatalf("the spec's DECLARATIONS are legal; only the seat hole is at issue: %v", err)
			}
			_, err = s.Resolve(map[string]string{"language": "cnos.eng:eng/go", "base_sha": "x"})
			if err == nil {
				t.Fatalf("hole %s must fail resolution", tc.hole)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("rejected for the wrong reason: got %v, want mention of %q", err, tc.want)
			}
		})
	}
}

// An explicitly supplied empty value is a VALUE, not an absence. `--param p=`
// is how a fake's meaningless model is written, so collapsing it into "unset"
// silently substituted a default or reported the parameter missing
// (Pi #59 C2).
func TestExplicitEmptyValueIsPreserved(t *testing.T) {
	const spec = `{
  "version": "cnos.cellspec.v0",
  "contract": {"id": "c1", "goal": "g"},
  "protocol_id": "cnos.cellkernel.episode-closure.v0",
  "params": {
    "model": {"required": false, "default": "a-default"},
    "needed": {"required": true}
  },
  "alpha": {"fill": "cdd.stub", "model": "$model"},
  "beta": {"fill": "cdd.stub"}
}`
	s, err := Parse([]byte(spec))
	if err != nil {
		t.Fatal(err)
	}

	r, err := s.Resolve(map[string]string{"model": "", "needed": "y"})
	if err != nil {
		t.Fatalf("an explicitly supplied empty value must resolve: %v", err)
	}
	if got := string(r.Alpha); !strings.Contains(got, `"model":""`) {
		t.Fatalf("explicit empty was not preserved: %s", got)
	}

	// Omitting it still falls back to the declared default...
	r, err = s.Resolve(map[string]string{"needed": "y"})
	if err != nil {
		t.Fatal(err)
	}
	if got := string(r.Alpha); !strings.Contains(got, `"model":"a-default"`) {
		t.Fatalf("omitted value must take the default: %s", got)
	}

	// ...and an omitted REQUIRED parameter still fails.
	if _, err := s.Resolve(map[string]string{"model": ""}); err == nil {
		t.Fatal("an omitted required parameter must still fail")
	}
}

// The task is carried OPAQUELY: this package must not learn what a CDS issue
// is, so a task that is not one — and a `$hole` string sitting inside one —
// both survive to the kernel byte for byte. Holes are spliced into the seat
// declarations only; a task is authored literally and frozen literally, so
// what a seat admits is what the author wrote.
func TestTaskIsCarriedOpaquelyAndHoleFree(t *testing.T) {
	const task = `{"kind":"not.a.cds.issue","note":"$base_sha stays a literal here"}`
	// Seats swapped for cdd fills so this test needs no CDS registry: whether
	// a task is carried is the GENERIC loader's property, not any fill's.
	src := strings.Replace(fixture,
		`"contract": {"id": "c1", "goal": "do the thing",`,
		`"contract": {"id": "c1", "goal": "do the thing", "task": `+task+`,`, 1)
	src = src[:strings.Index(src, `  "alpha": {`)] +
		`  "alpha": {"fill": "cdd.bool", "value": "true"},
  "beta": {"fill": "cdd.bool-check"}
}`

	s, err := Parse([]byte(src))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r, err := s.Resolve(map[string]string{"language": "cnos.eng:eng/go", "base_sha": "abc123"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	spec, _, err := r.Build(context.Background(), cellfill.CddFills())
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if got := string(spec.Contract.Task); got != task {
		t.Fatalf("task did not reach the kernel unchanged:\n got: %s\nwant: %s", got, task)
	}
}
