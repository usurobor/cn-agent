package cellkernel

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"
)

// subjectSpec is the empty cell with a subject hung on its contract. As with
// the task, the payload is deliberately NOT the shape any real adapter uses:
// an adapter-shaped document in the generic kernel's own corpus would suggest
// the kernel knows what a subject is.
func subjectSpec(subject string) Spec {
	s := EmptySpec()
	s.Contract.Subject = json.RawMessage(subject)
	return s
}

// The subject rides the frozen contract into the record and therefore into the
// canonical bytes the ONE scope-lift digest is taken over — the same carriage
// the task already had, which is why no second digest is added.
func TestSubjectIsCarriedIntoTheRecordAndBoundToTheDigest(t *testing.T) {
	const subject = `{"anything":"the kernel never reads this","n":1}`
	cl := mechClosure(t, subjectSpec(subject))

	if !bytes.Equal(cl.Receipt.Record.Contract.Subject, json.RawMessage(subject)) {
		t.Fatalf("record dropped the subject: %s", cl.Receipt.Record.Contract.Subject)
	}
	if !bytes.Contains(cl.Receipt.Record.canonicalBytes(), []byte(`"subject":`)) {
		t.Fatal("the canonical bytes carry no subject, so the digest cannot bind it")
	}
	if err := VerifyClosure(subjectSpec(subject).Contract, testMeta(ModeMechanical), cl); err != nil {
		t.Fatalf("closure must self-verify: %v", err)
	}
	rt := roundTrip(t, cl)
	if err := VerifyClosure(subjectSpec(subject).Contract, testMeta(ModeMechanical), rt); err != nil {
		t.Fatalf("round-tripped closure must verify: %v", err)
	}

	// One byte, one different digest. Without this the assertions above would
	// hold for a subject the digest ignored.
	other := mechClosure(t, subjectSpec(`{"anything":"the kernel never reads this","n":2}`))
	if other.Receipt.ScopeLiftDigest == cl.Receipt.ScopeLiftDigest {
		t.Fatal("a one-byte change to the subject left the scope-lift digest unmoved")
	}
	// ... and the same subject must still give the same digest, or the
	// inequality above proves nothing but nondeterminism.
	if again := mechClosure(t, subjectSpec(subject)); again.Receipt.ScopeLiftDigest != cl.Receipt.ScopeLiftDigest {
		t.Fatal("the same subject produced two digests")
	}
}

// clone() must copy the subject bytes for the reason it copies the task's: a
// shared backing array is a mutable value inside a struct the whole design
// treats as frozen, and here it would be mutable by BOTH seats.
func TestCloneDoesNotShareTheSubjectBackingArray(t *testing.T) {
	original := Contract{ID: "c", Goal: "g", Subject: json.RawMessage(`{"at":"AAAA"}`)}
	cp := original.clone()
	if !bytes.Equal(cp.Subject, original.Subject) {
		t.Fatalf("clone changed the subject: %s", cp.Subject)
	}
	for i := range cp.Subject {
		cp.Subject[i] = 'Z'
	}
	if string(original.Subject) != `{"at":"AAAA"}` {
		t.Fatalf("writing through the clone reached the original: %s", original.Subject)
	}
}

// The kernel's ONLY rules for the opaque subject are structural, and both are
// typed InvalidRecord — integrity, not contract-unmet — because a record whose
// canonical bytes cannot be produced is not a repairable episode.
func TestStructurallyInadmissibleSubjectIsAnInvalidRecord(t *testing.T) {
	cases := map[string]json.RawMessage{
		"not JSON":   json.RawMessage(`{"repo":`),
		"over bound": json.RawMessage(`{"pad":"` + strings.Repeat("x", maxOpaqueSlotBytes) + `"}`),
	}
	for name, subject := range cases {
		t.Run(name, func(t *testing.T) {
			// The honest path refuses it before alpha runs...
			_, err := RunEpisode(context.Background(), subjectSpec(string(subject)), testMeta(ModeMechanical),
				WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
			if err == nil {
				t.Fatal("an inadmissible subject must not start an episode")
			}
			if !strings.Contains(err.Error(), "contract.subject") {
				t.Fatalf("rejected for the wrong reason: %v", err)
			}
			// ...and a record carrying one anyway does not self-verify, which is
			// the half a forger controls.
			rec := EpisodeRecord{
				Canon: RecordCanon, EpisodeID: "ep", Mode: ModeMechanical,
				ResolvedSpec: testMeta(ModeMechanical).ResolvedSpec,
				Contract:     Contract{ID: "c", Subject: subject},
				Alpha:        StationRecord{ExecutionID: "a", Artifacts: []Artifact{}},
				Beta:         StationRecord{ExecutionID: "b", Artifacts: []Artifact{}},
			}
			var found bool
			for _, f := range validateRecord(rec) {
				if f.Class == InvalidRecord && strings.Contains(f.Detail, "subject") {
					found = true
				}
			}
			if !found {
				t.Fatalf("validateRecord did not flag the subject: %+v", validateRecord(rec))
			}
		})
	}
}

// An absent subject stays absent: `omitempty` keeps it out of the canonical
// bytes, so every cell that acts on nothing keeps its digest and its closure.
func TestAbsentSubjectIsAdmissibleAndInvisible(t *testing.T) {
	cl := mechClosure(t, EmptySpec())
	if bytes.Contains(cl.Receipt.Record.canonicalBytes(), []byte(`"subject"`)) {
		t.Fatalf("an absent subject reached the canonical bytes: %s", cl.Receipt.Record.canonicalBytes())
	}
	if cl.Status != Accepted {
		t.Fatalf("a subjectless cell must still close: %q", cl.Status)
	}
}

// The watching seats capture exactly what each station was handed, so "both seats
// received the same subject" can be asserted on bytes rather than assumed from
// the fact that one contract was frozen.
type subjectWatchAlpha struct{ saw *json.RawMessage }

func (r subjectWatchAlpha) Produce(_ context.Context, in AlphaInput) (AlphaOutput, error) {
	*r.saw = in.Contract.Subject
	return AlphaOutput{Matter: Matter{Data: "m"}}, nil
}

type subjectWatchBeta struct{ saw *json.RawMessage }

func (r subjectWatchBeta) Review(_ context.Context, in BetaInput) (BetaOutput, error) {
	*r.saw = in.Contract.Subject
	return BetaOutput{Review: Review{Pass: true, Notes: "n"}}, nil
}

// Both stations receive the same subject bytes, and neither can reach the
// other's copy. Sameness is structural — one frozen contract, cloned per
// station — and this is the executable statement of it.
func TestBothStationsReceiveTheSameSubjectBytes(t *testing.T) {
	const subject = `{"anything":"one frozen value","n":3}`
	var atAlpha, atBeta json.RawMessage
	s := subjectSpec(subject)
	s.Alpha, s.Beta = subjectWatchAlpha{&atAlpha}, subjectWatchBeta{&atBeta}

	cl := mechClosure(t, s)
	if string(atAlpha) != subject {
		t.Fatalf("alpha saw %s, want %s", atAlpha, subject)
	}
	if string(atBeta) != string(atAlpha) {
		t.Fatalf("the stations saw different subjects:\n alpha: %s\n beta:  %s", atAlpha, atBeta)
	}
	if string(cl.Receipt.Record.Contract.Subject) != subject {
		t.Fatalf("the record carries a third value: %s", cl.Receipt.Record.Contract.Subject)
	}
	// Distinct backing arrays: alpha writing through its copy must not reach
	// beta's, or "same bytes" would be one shared mutable buffer.
	if len(atAlpha) > 0 && &atAlpha[0] == &atBeta[0] {
		t.Fatal("the two stations share one subject buffer")
	}
}

// AC1's boundary, checked mechanically rather than asserted in a comment: the
// kernel's own source may not name any subject vocabulary. A kernel that knew
// what `base_sha` was would have stopped carrying opaque bytes, and no test of
// behaviour would notice — the leak is in the vocabulary, so the vocabulary is
// what is measured.
//
// Test files are excluded on purpose: a test that forbids a word has to be
// able to write it down.
func TestKernelNamesNoSubjectVocabulary(t *testing.T) {
	forbidden := []*regexp.Regexp{
		regexp.MustCompile(`git\.snapshot`),
		regexp.MustCompile(`base_sha`),
		regexp.MustCompile(`\brepo(sitory|sitories)?\b`),
		regexp.MustCompile(`\bworktree\b`),
	}
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read the kernel's own package directory: %v", err)
	}
	scanned := 0
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		src, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		scanned++
		for _, re := range forbidden {
			if loc := re.FindIndex(src); loc != nil {
				t.Errorf("%s names %q at byte %d: the kernel carries opaque bytes and must not learn what a subject is",
					name, src[loc[0]:loc[1]], loc[0])
			}
		}
	}
	// A scan that read no files would report a clean boundary for a package it
	// never opened.
	if scanned == 0 {
		t.Fatal("no kernel source files were scanned")
	}
}
