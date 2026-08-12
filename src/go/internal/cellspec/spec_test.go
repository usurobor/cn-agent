package cellspec

import (
	"context"
	"encoding/json"
	"fmt"
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
    "subject": {"kind": "git.snapshot/0.1", "repo": ".", "base_sha": "$base_sha"},
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

// Holes resolve IN PLACE inside the seat trees AND inside the contract's
// subject: the skill-list hole and the subject's base hole are each replaced
// where they sit. The subject is the one contract field holes reach, because it
// names where this run acts and that is exactly what an invoker supplies.
func TestResolveFillsHolesInPlace(t *testing.T) {
	r, err := mustParse(t).Resolve(map[string]string{"language": "cnos.eng:eng/go", "base_sha": "abc123"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	a := string(r.Alpha)
	for _, want := range []string{`"cnos.eng:eng/go"`, `"fill":"cds.patch"`} {
		if !strings.Contains(a, want) {
			t.Errorf("resolved alpha missing %s: %s", want, a)
		}
	}
	if strings.Contains(a, "$") {
		t.Fatalf("unresolved hole survived: %s", a)
	}
	if got := string(r.Subject); !strings.Contains(got, `"base_sha":"abc123"`) {
		t.Fatalf("the subject hole was not filled: %s", got)
	}
	if strings.Contains(string(r.Subject), "$") {
		t.Fatalf("unresolved hole survived in the subject: %s", r.Subject)
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
	// No subject either: the cdd fills act on nothing, and a spec that declared
	// one would need an adapter this registry does not wire.
	src := strings.Replace(fixture,
		`"subject": {"kind": "git.snapshot/0.1", "repo": ".", "base_sha": "$base_sha"},`, "", 1)
	src = strings.Replace(src,
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

// subjectFixture is the generic-envelope spec used for subject tests: cdd
// seats, so no CDS registry is needed, and the ONE hole sits in the contract's
// subject, which is the position under test.
const subjectFixture = `{
  "version": "cnos.cellspec.v0",
  "contract": {"id": "c1", "goal": "do the thing",
    "subject": {"kind": "opaque.to.this.package/0.1", "at": "$where"}},
  "protocol_id": "cnos.cdd.cds.receipt.v1",
  "params": {"where": {"required": true}},
  "alpha": {"fill": "cdd.bool", "value": "true"},
  "beta": {"fill": "cdd.bool-check"}
}`

func buildSubjectFixture(t *testing.T, src string, reg cellfill.Registry) (cellkernel.Spec, error) {
	t.Helper()
	s, err := Parse([]byte(src))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r, err := s.Resolve(map[string]string{"where": "abc123"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	spec, _, err := r.Build(context.Background(), reg)
	return spec, err
}

// The subject is carried OPAQUELY and pinned by the wired adapter, ONCE. This
// package never learns what a subject is: it hands the resolved bytes over and
// records the bytes it gets back, so a git snapshot and any later kind cost
// this loader the same nothing.
func TestSubjectIsPinnedOnceByTheWiredAdapter(t *testing.T) {
	calls := 0
	reg := cellfill.CddFills()
	reg.PinSubject = func(_ context.Context, subject json.RawMessage) (json.RawMessage, error) {
		calls++
		if !strings.Contains(string(subject), `"at":"abc123"`) {
			return nil, fmt.Errorf("the adapter was handed an unresolved subject: %s", subject)
		}
		return json.RawMessage(`{"kind":"opaque.to.this.package/0.1","at":"pinned"}`), nil
	}

	spec, err := buildSubjectFixture(t, subjectFixture, reg)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if calls != 1 {
		t.Fatalf("the subject was pinned %d times; both stations must receive one pinning", calls)
	}
	if string(spec.Contract.Subject) != `{"kind":"opaque.to.this.package/0.1","at":"pinned"}` {
		t.Fatalf("the kernel contract does not carry the PINNED subject: %s", spec.Contract.Subject)
	}
}

// A declared subject with no adapter fails construction. Running anyway would
// mean each station resolving the subject for itself, which is the divergence
// the contract slot exists to remove; and a spec with no subject must still
// build, or every cell that acts on nothing would need an adapter.
func TestSubjectWithoutAnAdapterFailsAndNoSubjectStillBuilds(t *testing.T) {
	_, err := buildSubjectFixture(t, subjectFixture, cellfill.CddFills())
	if err == nil {
		t.Fatal("a declared subject with no wired adapter must fail construction")
	}
	if !strings.Contains(err.Error(), "no subject adapter") {
		t.Fatalf("failed for the wrong reason: %v", err)
	}

	subjectless := strings.Replace(subjectFixture,
		`,
    "subject": {"kind": "opaque.to.this.package/0.1", "at": "$where"}`, "", 1)
	if subjectless == subjectFixture {
		t.Fatal("the subjectless variant is identical to the fixture, so it proves nothing")
	}
	s, err := Parse([]byte(strings.Replace(subjectless, `"params": {"where": {"required": true}},`, "", 1)))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r, err := s.Resolve(nil)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	spec, _, err := r.Build(context.Background(), cellfill.CddFills())
	if err != nil {
		t.Fatalf("a cell that acts on nothing must still build: %v", err)
	}
	if spec.Contract.Subject != nil {
		t.Fatalf("a subjectless spec grew one: %s", spec.Contract.Subject)
	}
}

// A subject the adapter refuses stops the run at construction, before either
// seat exists — the same discipline that keeps an ill-defined issue from
// reaching a provider.
func TestRefusedSubjectFailsBuild(t *testing.T) {
	reg := cellfill.CddFills()
	reg.PinSubject = func(context.Context, json.RawMessage) (json.RawMessage, error) {
		return nil, fmt.Errorf("no such revision")
	}
	if _, err := buildSubjectFixture(t, subjectFixture, reg); err == nil ||
		!strings.Contains(err.Error(), "no such revision") {
		t.Fatalf("a refused subject must fail the build with its own reason, got %v", err)
	}
}
