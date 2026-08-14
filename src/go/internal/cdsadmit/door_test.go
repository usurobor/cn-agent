// This file is the witness for D6: nothing cognitive runs before admission
// passes. It is an EXTERNAL test package because it drives the real runner —
// internal/cellrun imports this package, so the witness has to sit outside it
// to avoid a cycle. That is also the honest place for it: the property is
// about the runner's order of operations, and asserting it from inside the
// pure admission package would be asserting it about a function that has no
// provider to invoke in the first place.
package cdsadmit_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cdsadmit"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellmethod"
	"github.com/usurobor/cnos/src/go/internal/cellrun"

	"github.com/usurobor/cnos/src/go/internal/celltest"
)

const corpusDir = "../../../../schemas/cds/fixtures/runinput"

// recordingCell is the spec the witness runs. Its seats are the recording
// fills below — the ONLY fills in the registry it is given, so there is no
// other constructor a run could reach.
const recordingCell = `{
  "version": "cnos.cellspec.v0",
  "contract": {"id": "witness", "goal": "record any seat activity"},
  "protocol_id": "cnos.cellkernel.episode-closure.v0",
  "alpha": {"fill": "witness.recording"},
  "beta": {"fill": "witness.recording"}
}`

// ledger records every point at which cognition could begin. Construction is
// recorded as well as invocation, because a fill constructs its provider
// adapter and loads its skill bodies at CONSTRUCTION time: a run that got as
// far as building a seat has already touched the thing this property says it
// must not.
type ledger struct{ events []string }

func (l *ledger) note(what string) { l.events = append(l.events, what) }

type recordingAlpha struct{ l *ledger }

func (a recordingAlpha) Produce(context.Context, cellkernel.AlphaInput) (cellkernel.AlphaOutput, error) {
	a.l.note("alpha.Produce")
	return cellkernel.AlphaOutput{Matter: cellkernel.Matter{Data: "m"}}, nil
}

type recordingBeta struct{ l *ledger }

func (b recordingBeta) Review(context.Context, cellkernel.BetaInput) (cellkernel.BetaOutput, error) {
	b.l.note("beta.Review")
	return cellkernel.BetaOutput{Review: cellkernel.Review{Pass: true, Notes: "n"}}, nil
}

// recordingRegistry carries the REAL door beside the recording fills. The
// property under test is the runner's order of operations, so the door has to
// be the shipped one — a stub door would prove that a stub refuses early.
func recordingRegistry(l *ledger) cellfill.Registry {
	reg := doorlessRegistry(l)
	reg.Door = cdsadmit.Door
	return reg
}

// doorlessRegistry is the same registry with NO admission door: a runtime that
// ships no profile. It is a legitimate registry, not a broken one.
func doorlessRegistry(l *ledger) cellfill.Registry {
	const fill = "witness.recording"
	return cellfill.Registry{
		Alpha: map[string]cellfill.AlphaFill{
			// No declared requirement: this witness is about the RUNNER's order
			// of operations, and a fill that refused a subjectless run would
			// stop several cases below before the door was ever reached.
			fill: {Construct: func(context.Context, json.RawMessage, cellmethod.View) (cellfill.ConstructedAlpha, error) {
				l.note("construct alpha")
				return cellfill.ConstructedAlpha{
					Constructed: cellfill.Constructed{
						Decl: json.RawMessage(`{"fill":"` + fill + `"}`),
						Mode: cellkernel.ModeCognitive,
					},
					Seat: recordingAlpha{l},
				}, nil
			}},
		},
		Beta: map[string]cellfill.BetaFill{
			// No declared requirement here either, and for the same reason as
			// the alpha above.
			fill: {Construct: func(context.Context, json.RawMessage, cellmethod.View) (cellfill.ConstructedBeta, error) {
				l.note("construct beta")
				return cellfill.ConstructedBeta{
					Constructed: cellfill.Constructed{
						Decl: json.RawMessage(`{"fill":"` + fill + `"}`),
						Mode: cellkernel.ModeCognitive,
					},
					Seat: recordingBeta{l},
				}, nil
			}},
		},
	}
}

func runWitness(t *testing.T, l *ledger, inputPath string) (int, string, string) {
	t.Helper()
	return runWith(t, recordingRegistry(l), "--input", inputPath)
}

// runWith drives the real runner over the recording cell with the given extra
// arguments.
func runWith(t *testing.T, reg cellfill.Registry, extra ...string) (int, string, string) {
	t.Helper()
	specPath := filepath.Join(t.TempDir(), "cell.json")
	if err := os.WriteFile(specPath, []byte(recordingCell), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := cellrun.Run(context.Background(), reg,
		append([]string{"--contract", specPath}, extra...),
		strings.NewReader(""), &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

// receiptOf decodes the admission receipt a refusal printed to stdout, and
// insists that a closure was NOT printed: a closure would be a claim that an
// episode happened, and no episode exists before a contract is admitted.
func receiptOf(t *testing.T, stdout string) struct {
	Kind        string `json:"kind"`
	Outcome     string `json:"outcome"`
	InputDigest string `json:"input_digest"`
	Reason      string `json:"reason"`
} {
	t.Helper()
	var receipt struct {
		Kind        string `json:"kind"`
		Outcome     string `json:"outcome"`
		InputDigest string `json:"input_digest"`
		Reason      string `json:"reason"`
	}
	if err := json.Unmarshal([]byte(stdout), &receipt); err != nil {
		t.Fatalf("stdout is not an admission receipt (%v): %s", err, stdout)
	}
	if strings.Contains(stdout, "closure_schema") {
		t.Fatal("a refusal emitted an episode closure")
	}
	return receipt
}

// AC1's second half: each of the three refusals happens with ZERO seat
// activity. The ledger is the assertion — if admission ever moved after the
// registry is dispatched, "construct alpha" would appear and every case below
// would fail.
// ENVELOPE AND PAYLOAD REFUSALS ARE THE SAME PATH. The first three fixtures
// are refused by a payload rule, the last three by the envelope — a wrong
// `kind`, an unknown key, a key spelled with a capital. All six must exit 4
// with a receipt on stdout and zero seat activity.
//
// They did not. The envelope three exited 2 with EMPTY stdout and no receipt,
// because the runner decoded the envelope itself and treated failure as a read
// error, while payload refusals went through the door. A wrong `kind` is
// decisively inadmissible — the same class of fact as a design with no
// approach — and an operator reading exit 2 and silence was told their file
// could not be read. Listing both classes in one table is what makes "the same
// path" checkable rather than asserted.
func TestARefusedRunInputBuildsAndInvokesNoSeat(t *testing.T) {
	cases := map[string]string{
		"malformed issue":    "runinput-malformed-issue.json",
		"malformed design":   "runinput-malformed-design.json",
		"wrong-kind subject": "runinput-wrong-kind-subject.json",
		"bad envelope kind":  "runinput-bad-kind.json",
		"unknown key":        "runinput-unknown-key.json",
		"mixed-case key":     "runinput-mixed-case-key.json",
	}
	for name, fixture := range cases {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(corpusDir, fixture)
			var l ledger
			code, stdout, stderr := runWitness(t, &l, path)

			if len(l.events) != 0 {
				t.Fatalf("a refused run input reached the seats: %v", l.events)
			}
			if code != 4 {
				t.Fatalf("exit = %d, want 4 (run input refused); stderr: %s", code, stderr)
			}
			receipt := receiptOf(t, stdout)
			if receipt.Outcome != "rejected" {
				t.Fatalf("outcome = %q, want rejected", receipt.Outcome)
			}
			if receipt.Kind != cdsadmit.ReceiptKind {
				t.Fatalf("receipt kind = %q, want %q", receipt.Kind, cdsadmit.ReceiptKind)
			}
			if receipt.Reason == "" {
				t.Fatal("a refusal must name its reason")
			}
			// The reason reaches the operator on stderr too, so a caller
			// reading only the exit code and stderr is not left guessing.
			if !strings.Contains(stderr, receipt.Reason) {
				t.Fatalf("stderr does not carry the refusal reason\n stderr: %s\n reason: %s", stderr, receipt.Reason)
			}
			// O3: the receipt names WHICH document it refused. Nothing is
			// frozen into a contract and no closure exists, so this digest is
			// the only record of the artifact the decision was about.
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			sum := sha256.Sum256(raw)
			if want := hex.EncodeToString(sum[:]); receipt.InputDigest != want {
				t.Fatalf("input_digest = %q, want the sha256 of the refused document %q",
					receipt.InputDigest, want)
			}
		})
	}
}

// F3: a payload past the kernel's opaque-slot bound is refused AT THE DOOR,
// with zero constructions. Before this, an 84 KiB issue was admitted, its
// subject was resolved against a real repository, and both seats were built —
// the provider adapter constructed and every skill body read — before the
// kernel refused with `episode malfunction`. The ledger is the assertion: it
// records construction, not just invocation, because construction is where a
// provider adapter and skill loading already happened.
func TestAnOversizePayloadIsRefusedWithZeroConstructions(t *testing.T) {
	// Built on the run input that points at a REAL repository, so that without
	// the bound this document would be admitted, pinned, and carried into
	// construction. Against the committed fixture's `.` the pinning step would
	// fail first for an unrelated reason, and the ledger could never fire — the
	// test would pass whether or not the door held the bound.
	raw, err := os.ReadFile(runInputAgainstARealRepo(t))
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	// One field of the otherwise valid issue, inflated past the bound. The
	// document stays a well-formed, admissible-shaped issue — only its size is
	// wrong, so nothing but the bound can be what refuses it.
	var issue map[string]any
	if err := json.Unmarshal(doc["issue"], &issue); err != nil {
		t.Fatal(err)
	}
	issue["id"] = strings.Repeat("x", cellkernel.MaxOpaqueSlotBytes)
	if doc["issue"], err = json.Marshal(issue); err != nil {
		t.Fatal(err)
	}
	if len(doc["issue"]) <= cellkernel.MaxOpaqueSlotBytes {
		t.Fatalf("the fixture is not oversize (%d bytes)", len(doc["issue"]))
	}
	out, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "oversize.json")
	if err := os.WriteFile(path, out, 0o600); err != nil {
		t.Fatal(err)
	}

	var l ledger
	code, stdout, stderr := runWitness(t, &l, path)
	if len(l.events) != 0 {
		t.Fatalf("an oversize run input reached the seats: %v", l.events)
	}
	if code != 4 {
		t.Fatalf("exit = %d, want 4 (run input refused); stderr: %s", code, stderr)
	}
	receipt := receiptOf(t, stdout)
	if receipt.Outcome != "rejected" {
		t.Fatalf("outcome = %q, want rejected", receipt.Outcome)
	}
	if !strings.Contains(receipt.Reason, "contract slot bound") {
		t.Fatalf("refused for the wrong reason: %s", receipt.Reason)
	}
	// The same document one byte under the bound is admitted, so the refusal
	// above is the BOUND firing and not the inflated issue being malformed.
	issue["id"] = strings.Repeat("x", cellkernel.MaxOpaqueSlotBytes-len(doc["issue"])+len(issue["id"].(string))-1)
	if doc["issue"], err = json.Marshal(issue); err != nil {
		t.Fatal(err)
	}
	if len(doc["issue"]) != cellkernel.MaxOpaqueSlotBytes-1 {
		t.Fatalf("the under-bound fixture is %d bytes, want %d", len(doc["issue"]), cellkernel.MaxOpaqueSlotBytes-1)
	}
	under, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if _, receiptBytes, err := cdsadmit.Door(under); err != nil {
		t.Fatalf("a payload one byte under the bound must be admitted: %v\n%s", err, receiptBytes)
	}
}

// F4: the runner runs with a registry carrying NO door. A cell that acts on
// nothing is still a cell — the bool and stub corpus cells are exactly that —
// and the runner must not require a profile in order to run one.
func TestACellWithNoRunInputRunsWithoutADoor(t *testing.T) {
	var l ledger
	code, stdout, stderr := runWith(t, doorlessRegistry(&l))
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, "closure_schema") {
		t.Fatalf("a doorless run must still emit a closure: %s", stdout)
	}
	want := []string{"construct alpha", "construct beta", "alpha.Produce", "beta.Review"}
	if strings.Join(l.events, ",") != strings.Join(want, ",") {
		t.Fatalf("ledger = %v, want %v", l.events, want)
	}
}

// ...and a doorless runtime handed a run input says so, rather than refusing.
// Refusing would claim a door judged the document and found it inadmissible;
// nothing judged it. The distinction is the difference between "your input is
// wrong" and "this binary admits no input", and only one of them is repairable
// by editing the document.
func TestADoorlessRuntimeGivenAnInputSaysSoRatherThanRefusing(t *testing.T) {
	var l ledger
	code, stdout, stderr := runWith(t, doorlessRegistry(&l),
		"--input", filepath.Join(corpusDir, "valid-run-input.json"))
	if code == 4 {
		t.Fatal("a missing door must not be reported as a refused document")
	}
	if code != 2 {
		t.Fatalf("exit = %d, want 2 (no door); stderr: %s", code, stderr)
	}
	if !strings.Contains(stderr, "no admission door") {
		t.Fatalf("stderr does not name the missing door: %s", stderr)
	}
	if stdout != "" {
		t.Fatalf("a malfunction wrote to stdout: %s", stdout)
	}
	if len(l.events) != 0 {
		t.Fatalf("a doorless run with an input reached the seats: %v", l.events)
	}
}

// runInputAgainstARealRepo rewrites the corpus fixture's subject to point at a
// throwaway repository. The committed fixture names `.`, which is what an
// author writes and what pinning resolves — but pinning really does open the
// repository, so the admitted case needs one that exists.
func runInputAgainstARealRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	celltest.Git(t, repo, "init", "-q", "-b", "main")
	celltest.Git(t, repo, "commit", "-qm", "base", "--allow-empty")
	raw, err := os.ReadFile(filepath.Join(corpusDir, "valid-run-input.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	subject, err := json.Marshal(map[string]string{
		"kind": "git.snapshot/0.1", "repo": repo, "base_sha": "HEAD",
	})
	if err != nil {
		t.Fatal(err)
	}
	doc["subject"] = subject
	out, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "run-input.json")
	if err := os.WriteFile(path, out, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// The ledger is not silent by construction. Without this, every assertion
// above would pass against a registry whose factories were never wired up, or
// against a runner that failed for some unrelated reason before it got
// anywhere — which is exactly the vacuity the ledger exists to rule out.
func TestTheLedgerRecordsAnAdmittedRun(t *testing.T) {
	var l ledger
	code, stdout, stderr := runWitness(t, &l, runInputAgainstARealRepo(t))
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr: %s", code, stderr)
	}
	want := []string{"construct alpha", "construct beta", "alpha.Produce", "beta.Review"}
	if strings.Join(l.events, ",") != strings.Join(want, ",") {
		t.Fatalf("ledger = %v, want %v", l.events, want)
	}
	if !strings.Contains(stdout, "closure_schema") {
		t.Fatalf("an admitted run must emit a closure: %s", stdout)
	}
}
