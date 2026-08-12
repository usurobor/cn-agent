package cellkernel

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// taskSpec is the empty cell with a task hung on its contract. The kernel
// never interprets the bytes, so any JSON object serves — and the payloads
// below are deliberately NOT CDS issues, since a CDS shape appearing in the
// generic kernel's own corpus would suggest the kernel knows one.
func taskSpec(task string) Spec {
	s := EmptySpec()
	s.Contract.Task = json.RawMessage(task)
	return s
}

// The task rides the frozen contract into the record, and therefore into the
// canonical bytes the ONE scope-lift digest is taken over. No second digest is
// needed and none is added; this test is the assertion that the existing one
// already covers it.
func TestTaskIsCarriedIntoTheRecordAndBoundToTheDigest(t *testing.T) {
	const task = `{"anything":"the kernel never reads this","n":1}`
	cl := mechClosure(t, taskSpec(task))

	if !bytes.Equal(cl.Receipt.Record.Contract.Task, json.RawMessage(task)) {
		t.Fatalf("record dropped the task: %s", cl.Receipt.Record.Contract.Task)
	}
	if !bytes.Contains(cl.Receipt.Record.canonicalBytes(), []byte(`"task":`)) {
		t.Fatal("the canonical bytes carry no task, so the digest cannot bind it")
	}
	if err := VerifyClosure(taskSpec(task).Contract, testMeta(ModeMechanical), cl); err != nil {
		t.Fatalf("closure must self-verify: %v", err)
	}
	// Survives serialization: a verifier reading the emitted JSON sees the same
	// task and recomputes the same digest.
	rt := roundTrip(t, cl)
	if err := VerifyClosure(taskSpec(task).Contract, testMeta(ModeMechanical), rt); err != nil {
		t.Fatalf("round-tripped closure must verify: %v", err)
	}

	// One byte, one different digest. Without this the previous assertions
	// would hold for a task the digest ignored.
	other := mechClosure(t, taskSpec(`{"anything":"a different opaque payload","n":2}`))
	if other.Receipt.ScopeLiftDigest == cl.Receipt.ScopeLiftDigest {
		t.Fatal("a one-byte change to the issue left the scope-lift digest unmoved")
	}
	// ... and the same task must still give the same digest, or the inequality
	// above would be proving nothing but nondeterminism.
	if again := mechClosure(t, taskSpec(task)); again.Receipt.ScopeLiftDigest != cl.Receipt.ScopeLiftDigest {
		t.Fatal("the same task produced two digests")
	}
}

// clone() must copy the task bytes. A shared backing array is a mutable value
// inside a struct the whole design treats as frozen: a seat handed
// AlphaInput.Contract could write through it into what the runtime records.
func TestCloneDoesNotShareTheTaskBackingArray(t *testing.T) {
	original := Contract{ID: "c", Goal: "g", Task: json.RawMessage(`{"id":"AAAA"}`)}
	cp := original.clone()
	if !bytes.Equal(cp.Task, original.Task) {
		t.Fatalf("clone changed the task: %s", cp.Task)
	}
	for i := range cp.Task {
		cp.Task[i] = 'Z'
	}
	if string(original.Task) != `{"id":"AAAA"}` {
		t.Fatalf("writing through the clone reached the original: %s", original.Task)
	}
}

// The kernel's ONLY rules for the opaque slot are structural. Both are typed
// InvalidRecord — integrity, not contract-unmet — because a record whose
// canonical bytes cannot be produced is not a repairable episode.
func TestStructurallyInadmissibleTaskIsAnInvalidRecord(t *testing.T) {
	cases := map[string]json.RawMessage{
		"not JSON": json.RawMessage(`{"id":`),
		"over bound": json.RawMessage(`{"pad":"` +
			strings.Repeat("x", maxOpaqueSlotBytes) + `"}`),
	}
	for name, task := range cases {
		t.Run(name, func(t *testing.T) {
			// The honest path refuses it before alpha runs...
			_, err := RunEpisode(context.Background(), taskSpec(string(task)), testMeta(ModeMechanical),
				WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
			if err == nil {
				t.Fatal("an inadmissible task must not start an episode")
			}
			if !strings.Contains(err.Error(), "contract.task") {
				t.Fatalf("rejected for the wrong reason: %v", err)
			}
			// ...and a record carrying one anyway does not self-verify, which is
			// the half a forger controls.
			rec := EpisodeRecord{
				Canon: RecordCanon, EpisodeID: "ep", Mode: ModeMechanical,
				ResolvedSpec: testMeta(ModeMechanical).ResolvedSpec,
				Contract:     Contract{ID: "c", Task: task},
				Alpha:        StationRecord{ExecutionID: "a", Artifacts: []Artifact{}},
				Beta:         StationRecord{ExecutionID: "b", Artifacts: []Artifact{}},
			}
			var found bool
			for _, f := range validateRecord(rec) {
				if f.Class == InvalidRecord && strings.Contains(f.Detail, "task") {
					found = true
				}
			}
			if !found {
				t.Fatalf("validateRecord did not flag the task: %+v", validateRecord(rec))
			}
		})
	}
}

// An absent task stays absent: `omitempty` keeps it out of the canonical bytes
// entirely, so existing closures without a task keep their digests and every
// non-CDS cell is unaffected. Requiring a task is a protocol's rule, not the
// kernel's.
func TestAbsentTaskIsAdmissibleAndInvisible(t *testing.T) {
	cl := mechClosure(t, EmptySpec())
	if bytes.Contains(cl.Receipt.Record.canonicalBytes(), []byte(`"task"`)) {
		t.Fatalf("an absent task reached the canonical bytes: %s", cl.Receipt.Record.canonicalBytes())
	}
	if cl.Status != Accepted {
		t.Fatalf("a taskless cell must still close: %q", cl.Status)
	}
}
