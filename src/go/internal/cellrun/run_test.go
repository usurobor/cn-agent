package cellrun

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

const boolSpecJSON = `{"version":"cnos.cellspec.v0",` +
	`"contract":{"id":"cell-bool","goal":"b","required_evidence":[{"id":"bool","kind":"value","producer":"alpha"}]},` +
	`"protocol_id":"cnos.cellkernel.episode-receipt.v0","profile":"bool",` +
	`"params":{"value":{"kind":"value","required":true,"domain":["true","false"]}},` +
	`"alpha":{"skills":[]},"beta":{"skills":[]}}`

// run drives Run with buffered stdio and returns the exit code + streams.
func run(stdin string, args ...string) (code int, stdout, stderr string) {
	var out, errb bytes.Buffer
	code = Run(context.Background(), args, strings.NewReader(stdin), &out, &errb)
	return code, out.String(), errb.String()
}

func TestAcceptedFromStdin(t *testing.T) {
	code, stdout, stderr := run(boolSpecJSON, "--contract", "-", "--param", "value=true")
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr)
	}
	if stderr != "" {
		t.Errorf("stderr not empty: %q", stderr)
	}
	// stdout must be exactly one valid receipt that re-verifies.
	dec := json.NewDecoder(strings.NewReader(stdout))
	var rc cellkernel.Receipt
	if err := dec.Decode(&rc); err != nil {
		t.Fatalf("stdout not a receipt: %v", err)
	}
	if dec.More() {
		t.Fatal("stdout carried more than one JSON value")
	}
	if err := cellkernel.VerifyReceipt(rc); err != nil {
		t.Fatalf("emitted receipt does not verify: %v", err)
	}
}

func TestExitCodes(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "spec.json")
	if err := os.WriteFile(tmp, []byte(boolSpecJSON), 0o600); err != nil {
		t.Fatal(err)
	}
	big := "{" + strings.Repeat(" ", maxContractBytes+10) + "}"

	cases := []struct {
		name  string
		stdin string
		args  []string
		want  int
	}{
		{"file accepted", "", []string{"--contract", tmp, "--param", "value=true"}, 0},
		{"needs_repair exit 1", boolSpecJSON, []string{"--contract", "-", "--param", "value=false"}, 1},
		{"malformed json", "{not json", []string{"--contract", "-"}, 2},
		{"unknown arg", "", []string{"--bogus"}, 2},
		{"missing contract", "", []string{"--param", "value=true"}, 2},
		{"dup param", boolSpecJSON, []string{"--contract", "-", "--param", "value=true", "--param", "value=false"}, 2},
		{"dup contract", boolSpecJSON, []string{"--contract", "-", "--contract", "-"}, 2},
		{"unknown protocol", `{"version":"cnos.cellspec.v0","contract":{"id":"c","goal":"g"},"protocol_id":"made.up","profile":"stub","alpha":{"skills":[]},"beta":{"skills":[]}}`, []string{"--contract", "-"}, 2},
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
