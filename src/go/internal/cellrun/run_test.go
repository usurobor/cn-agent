package cellrun

import (
	"bytes"
	"context"
	"encoding/json"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellfills"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
	"github.com/usurobor/cnos/src/go/internal/cellspec"
)

const boolSpecJSON = `{"version":"cnos.cellspec.v0",` +
	`"contract":{"id":"cell-bool","goal":"b","required_evidence":[{"id":"bool","kind":"value","producer":"alpha"}]},` +
	`"protocol_id":"cnos.cellkernel.episode-closure.v0",` +
	`"params":{"value":{"required":true,"domain":["true","false"]}},` +
	`"alpha":{"fill":"cdd.bool","value":"$value"},"beta":{"fill":"cdd.bool-check"}}`

// testRegistry is the assembled registry with an empty skill tree: these
// cases exercise generic dispatch, not skill loading.
func testRegistry() cellfill.Registry {
	return cellfills.With(cellskill.Tree{Root: "/nonexistent"})
}

func run(stdin string, args ...string) (code int, stdout, stderr string) {
	var out, errb bytes.Buffer
	code = Run(context.Background(), testRegistry(), args, strings.NewReader(stdin), &out, &errb)
	return code, out.String(), errb.String()
}

// parseExpected derives the trusted contract and invocation metadata from
// boolSpecJSON via the same loader path — independently of any emitted closure.
func parseExpected(t *testing.T) (cellkernel.Contract, cellkernel.RunMeta) {
	t.Helper()
	s, err := cellspec.Parse([]byte(boolSpecJSON))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	r, err := s.Resolve(map[string]string{"value": "true"})
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	kspec, meta, err := r.Build(context.Background(), testRegistry(), cellspec.Binding{})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	return kspec.Contract, meta
}

func TestAcceptedFromStdin(t *testing.T) {
	code, stdout, stderr := run(boolSpecJSON, "--contract", "-", "--param", "value=true")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr)
	}
	if stderr != "" {
		t.Errorf("stderr not empty: %q", stderr)
	}
	// stdout must be exactly one valid closure that re-verifies whole.
	dec := json.NewDecoder(strings.NewReader(stdout))
	var cl cellkernel.Closure
	if err := dec.Decode(&cl); err != nil {
		t.Fatalf("stdout not a closure: %v", err)
	}
	var extra json.RawMessage
	if err := dec.Decode(&extra); err != io.EOF {
		t.Fatalf("stdout carried trailing data (want io.EOF, got %v)", err)
	}
	// Re-verify against the contract the test itself trusts — parsed and
	// resolved independently of the emitted closure.
	expected, meta := parseExpected(t)
	if err := cellkernel.VerifyClosure(expected, meta, cl); err != nil {
		t.Fatalf("emitted closure does not verify: %v", err)
	}
}

// F4, checked mechanically rather than asserted in a comment: this package's
// own source may not import a protocol package. The admission door arrives in
// the registry, so a runner that named `cdsadmit` would have taken back the
// coupling the registry exists to remove — and no test of behaviour would
// notice, because the CDS door is the one this binary happens to wire in. The
// leak is in the import graph, so the import graph is what is measured.
//
// Test files are excluded: a test that forbids an import has to be able to
// drive the real door through one.
func TestTheRunnerNamesNoProtocolPackage(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read the runner's own package directory: %v", err)
	}
	scanned := 0
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(token.NewFileSet(), name, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		scanned++
		for _, imp := range f.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if strings.Contains(path, "/internal/cds") {
				t.Errorf("%s imports %s: the runner dispatches a door, it does not name one", name, path)
			}
		}
	}
	// A scan that read no files would report a clean boundary for a package it
	// never opened.
	if scanned == 0 {
		t.Fatal("no runner source files were scanned")
	}
}

func TestExitCodes(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "spec.json")
	if err := os.WriteFile(tmp, []byte(boolSpecJSON), 0o600); err != nil {
		t.Fatal(err)
	}
	stubSpec := `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},` +
		`"protocol_id":"cnos.cellkernel.episode-closure.v0",` +
		`"alpha":{"fill":"cdd.stub"},"beta":{"fill":"cdd.stub"}}`
	big := "{" + strings.Repeat(" ", maxContractBytes+10) + "}"

	cases := []struct {
		name  string
		stdin string
		args  []string
		want  int
	}{
		{"accepted", "", []string{"--contract", tmp, "--param", "value=true"}, 0},
		{"needs_repair", boolSpecJSON, []string{"--contract", "-", "--param", "value=false"}, 1},
		{"simulated stub", stubSpec, []string{"--contract", "-"}, 3},
		{"malformed json", "{not json", []string{"--contract", "-"}, 2},
		{"unknown arg", "", []string{"--bogus"}, 2},
		{"missing contract", "", []string{"--param", "value=true"}, 2},
		{"missing fill", `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"cnos.cellkernel.episode-closure.v0","alpha":{},"beta":{"fill":"cdd.stub"}}`, []string{"--contract", "-"}, 2},
		{"unknown fill", `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"cnos.cellkernel.episode-closure.v0","alpha":{"fill":"no.such"},"beta":{"fill":"cdd.stub"}}`, []string{"--contract", "-"}, 2},
		{"dup param", boolSpecJSON, []string{"--contract", "-", "--param", "value=true", "--param", "value=false"}, 2},
		{"dup contract", boolSpecJSON, []string{"--contract", "-", "--contract", "-"}, 2},
		{"opaque protocol runs", `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"made.up","alpha":{"fill":"cdd.stub"},"beta":{"fill":"cdd.stub"}}`, []string{"--contract", "-"}, 3},
		{"trailing brace", boolSpecJSON + "}", []string{"--contract", "-", "--param", "value=true"}, 2},
		{"oversize contract", big, []string{"--contract", "-"}, 2},
		{"dup input", boolSpecJSON, []string{"--contract", "-", "--input", "a", "--input", "b"}, 2},
		{"input without value", boolSpecJSON, []string{"--contract", "-", "--input"}, 2},
		// One stream cannot carry two documents; taking the whole of stdin for
		// whichever read ran first would be a corrupt run reported as a
		// malformed file.
		{"both read stdin", boolSpecJSON, []string{"--contract", "-", "--input", "-"}, 2},
		{"missing input file", boolSpecJSON, []string{"--contract", "-", "--input", "/nonexistent/run-input.json"}, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, stdout, _ := run(tc.stdin, tc.args...)
			if code != tc.want {
				t.Fatalf("exit=%d, want %d", code, tc.want)
			}
			if tc.want == 2 && stdout != "" {
				t.Errorf("error path wrote to stdout: %q", stdout)
			}
		})
	}
}

// The shipped `cds.patch` cell, run with no --input, over the SHIPPED registry.
// It is refused before its alpha is constructed, and the proof is testRegistry's
// skill root: it does not exist, so construction could not have succeeded — an
// error mentioning a skill would say construction had begun, and an "episode
// malfunction" would say a station had. The run is refused by the declaration
// alone, which is why removing a skill from an installed hub cannot change what
// an operator sees here.
func TestASubjectlessPatchCellIsRefusedBeforeConstruction(t *testing.T) {
	spec, err := os.ReadFile("../../../../schemas/cds/fixtures/code-cell-spec.json")
	if err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := run(string(spec), "--contract", "-",
		"--param", "language=cnos.eng:eng/go", "--param", "provider=fake")

	if code != 2 {
		t.Fatalf("exit = %d, want 2; stderr: %s", code, stderr)
	}
	if stdout != "" {
		t.Fatalf("a refused run wrote to stdout: %s", stdout)
	}
	for _, want := range []string{`"cds.patch"`, "contract.subject", "--input"} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("the refusal does not name %q: %s", want, stderr)
		}
	}
	// Construction never began. Both markers below belong to work that happens
	// strictly after the check: skill loading is inside the constructor, and
	// "episode malfunction" is printed only for a failure raised by a station.
	for _, forbidden := range []string{"skill", "episode malfunction"} {
		if strings.Contains(stderr, forbidden) {
			t.Fatalf("the refusal came after construction (%q): %s", forbidden, stderr)
		}
	}
	// The cell really is the one whose alpha declares the requirement, and the
	// arguments really are otherwise complete — or the assertions above would
	// hold for a spec rejected for some unrelated reason.
	if !strings.Contains(string(spec), `"fill": "cds.patch"`) {
		t.Fatalf("the fixture is not the cds.patch cell: %s", spec)
	}
}
