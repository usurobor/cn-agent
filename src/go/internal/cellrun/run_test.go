package cellrun

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellspec"
)

const boolSpecJSON = `{"version":"cnos.cellspec.v0",` +
	`"contract":{"id":"cell-bool","goal":"b","required_evidence":[{"id":"bool","kind":"value","producer":"alpha"}]},` +
	`"protocol_id":"cnos.cellkernel.episode-closure.v0",` +
	`"params":{"value":{"kind":"value","required":true,"domain":["true","false"]}},` +
	`"alpha":{"fill":"cdd.bool","value":"$value"},"beta":{"fill":"cdd.bool-check"}}`

func run(stdin string, args ...string) (code int, stdout, stderr string) {
	var out, errb bytes.Buffer
	code = Run(context.Background(), args, strings.NewReader(stdin), &out, &errb)
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
	kspec, meta, err := r.Build(registry())
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
