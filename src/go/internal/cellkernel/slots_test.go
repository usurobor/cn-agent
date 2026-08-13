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

// The three opaque slots obey identical rules, so they are tested through one
// table rather than three near-identical files. Each entry is the slot's name
// in the canonical JSON, a setter, and a reader.
//
// The payloads below are deliberately NOT CDS issues, designs or git
// snapshots: a domain-shaped document in the generic kernel's own corpus would
// suggest the kernel knows one.
var slots = []struct {
	name string
	set  func(*Contract, json.RawMessage)
	get  func(Contract) json.RawMessage
}{
	{"issue", func(c *Contract, r json.RawMessage) { c.Issue = r }, func(c Contract) json.RawMessage { return c.Issue }},
	{"design", func(c *Contract, r json.RawMessage) { c.Design = r }, func(c Contract) json.RawMessage { return c.Design }},
	{"subject", func(c *Contract, r json.RawMessage) { c.Subject = r }, func(c Contract) json.RawMessage { return c.Subject }},
}

// slotSpec is the empty cell with one opaque slot hung on its contract.
func slotSpec(set func(*Contract, json.RawMessage), payload string) Spec {
	s := EmptySpec()
	set(&s.Contract, json.RawMessage(payload))
	return s
}

// AC3. Each slot rides the frozen contract into the record, and therefore into
// the canonical bytes the ONE scope-lift digest is taken over. No second
// digest is needed and none is added; this is the assertion that the existing
// one already covers them.
func TestEachSlotIsCarriedIntoTheRecordAndBoundToTheDigest(t *testing.T) {
	for _, slot := range slots {
		t.Run(slot.name, func(t *testing.T) {
			const payload = `{"anything":"the kernel never reads this","n":1}`
			spec := func() Spec { return slotSpec(slot.set, payload) }
			cl := mechClosure(t, spec())

			if !bytes.Equal(slot.get(cl.Receipt.Record.Contract), json.RawMessage(payload)) {
				t.Fatalf("record dropped the %s: %s", slot.name, slot.get(cl.Receipt.Record.Contract))
			}
			if !bytes.Contains(cl.Receipt.Record.canonicalBytes(), []byte(`"`+slot.name+`":`)) {
				t.Fatalf("the canonical bytes carry no %s, so the digest cannot bind it", slot.name)
			}
			if err := VerifyClosure(spec().Contract, testMeta(ModeMechanical), cl); err != nil {
				t.Fatalf("closure must self-verify: %v", err)
			}
			// Survives serialization: a verifier reading the emitted JSON sees
			// the same payload and recomputes the same digest.
			rt := roundTrip(t, cl)
			if err := VerifyClosure(spec().Contract, testMeta(ModeMechanical), rt); err != nil {
				t.Fatalf("round-tripped closure must verify: %v", err)
			}

			// One byte, one different digest. Without this the assertions above
			// would hold for a payload the digest ignored.
			other := mechClosure(t, slotSpec(slot.set, `{"anything":"the kernel never reads this","n":2}`))
			if other.Receipt.ScopeLiftDigest == cl.Receipt.ScopeLiftDigest {
				t.Fatalf("a one-byte change to the %s left the scope-lift digest unmoved", slot.name)
			}
			// ...and the same payload must still give the same digest, or the
			// inequality above proves nothing but nondeterminism.
			if again := mechClosure(t, spec()); again.Receipt.ScopeLiftDigest != cl.Receipt.ScopeLiftDigest {
				t.Fatalf("the same %s produced two digests", slot.name)
			}
		})
	}
}

// The three slots are digested as three DISTINCT positions. Without this, one
// slot's bytes appearing in another's position would move the digest just as
// much, and "a change to the design moves the digest" would be true of a
// record that had confused design with issue.
func TestTheSlotsAreDistinctPositionsInTheDigest(t *testing.T) {
	const payload = `{"anything":"one value in three places"}`
	seen := make(map[string]string)
	for _, slot := range slots {
		cl := mechClosure(t, slotSpec(slot.set, payload))
		d := cl.Receipt.ScopeLiftDigest
		if prev, dup := seen[d]; dup {
			t.Fatalf("the %s and the %s slot digest identically", prev, slot.name)
		}
		seen[d] = slot.name
	}
}

// clone() must copy every slot's bytes. A shared backing array is a mutable
// value inside a struct the whole design treats as frozen: a seat handed
// AlphaInput.Contract could write through it into what the runtime records.
func TestCloneDoesNotShareASlotBackingArray(t *testing.T) {
	for _, slot := range slots {
		t.Run(slot.name, func(t *testing.T) {
			const payload = `{"id":"AAAA"}`
			original := Contract{ID: "c", Goal: "g"}
			slot.set(&original, json.RawMessage(payload))
			cp := original.clone()
			if !bytes.Equal(slot.get(cp), slot.get(original)) {
				t.Fatalf("clone changed the %s: %s", slot.name, slot.get(cp))
			}
			raw := slot.get(cp)
			for i := range raw {
				raw[i] = 'Z'
			}
			if string(slot.get(original)) != payload {
				t.Fatalf("writing through the clone reached the original: %s", slot.get(original))
			}
		})
	}
}

// The kernel's ONLY rules for an opaque slot are structural. All are typed
// InvalidRecord — integrity, not contract-unmet — because a record whose
// canonical bytes cannot be produced, or which no reader downstream can accept,
// is not a repairable episode.
//
// The non-object cases are the SECOND authority's rule made this one's: the
// closure schema declares all three slots `{...}`, so before this a scalar slot
// self-verified here and then failed `cue vet` with a type conflict. A slot
// carries one tagged value; a bare scalar is not a tagged anything.
func TestStructurallyInadmissibleSlotIsAnInvalidRecord(t *testing.T) {
	payloads := map[string]json.RawMessage{
		"not JSON":   json.RawMessage(`{"id":`),
		"over bound": json.RawMessage(`{"pad":"` + strings.Repeat("x", MaxOpaqueSlotBytes) + `"}`),
		"number":     json.RawMessage(`42`),
		"string":     json.RawMessage(`"a string"`),
		"bool":       json.RawMessage(`true`),
		"null":       json.RawMessage(`null`),
		"array":      json.RawMessage(`[{"id":"x"}]`),
	}
	for _, slot := range slots {
		for name, payload := range payloads {
			t.Run(slot.name+"/"+name, func(t *testing.T) {
				// The honest path refuses it before alpha runs...
				_, err := RunEpisode(context.Background(), slotSpec(slot.set, string(payload)),
					testMeta(ModeMechanical), WithIDSource(seqIDs{"ep-t", "alpha-t", "beta-t"}))
				if err == nil {
					t.Fatal("an inadmissible slot must not start an episode")
				}
				if !strings.Contains(err.Error(), "contract."+slot.name) {
					t.Fatalf("rejected for the wrong reason: %v", err)
				}
				// ...and a record carrying one anyway does not self-verify,
				// which is the half a forger controls.
				c := Contract{ID: "c"}
				slot.set(&c, payload)
				rec := EpisodeRecord{
					Canon: RecordCanon, EpisodeID: "ep", Mode: ModeMechanical,
					ResolvedSpec: testMeta(ModeMechanical).ResolvedSpec,
					Contract:     c,
					Alpha:        StationRecord{ExecutionID: "a", Artifacts: []Artifact{}},
					Beta:         StationRecord{ExecutionID: "b", Artifacts: []Artifact{}},
				}
				var found bool
				for _, f := range validateRecord(rec) {
					if f.Class == InvalidRecord && strings.Contains(f.Detail, slot.name) {
						found = true
					}
				}
				if !found {
					t.Fatalf("validateRecord did not flag the %s: %+v", slot.name, validateRecord(rec))
				}
			})
		}
	}
}

// Absent slots stay absent: `omitempty` keeps them out of the canonical bytes
// entirely, so existing closures keep their digests and every cell that
// carries no run input is unaffected. Requiring a slot is a protocol's rule,
// not the kernel's.
func TestAbsentSlotsAreAdmissibleAndInvisible(t *testing.T) {
	cl := mechClosure(t, EmptySpec())
	canon := cl.Receipt.Record.canonicalBytes()
	for _, slot := range slots {
		if bytes.Contains(canon, []byte(`"`+slot.name+`"`)) {
			t.Fatalf("an absent %s reached the canonical bytes: %s", slot.name, canon)
		}
	}
	if cl.Status != Accepted {
		t.Fatalf("a cell with no run input must still close: %q", cl.Status)
	}
}

// The watching seats capture exactly what each station was handed, so "both
// seats received the same subject" can be asserted on bytes rather than
// assumed from the fact that one contract was frozen.
type slotWatchAlpha struct{ saw *Contract }

func (r slotWatchAlpha) Produce(_ context.Context, in AlphaInput) (AlphaOutput, error) {
	*r.saw = in.Contract
	return AlphaOutput{Matter: Matter{Data: "m"}}, nil
}

type slotWatchBeta struct{ saw *Contract }

func (r slotWatchBeta) Review(_ context.Context, in BetaInput) (BetaOutput, error) {
	*r.saw = in.Contract
	return BetaOutput{Review: Review{Pass: true, Notes: "n"}}, nil
}

// D4's second half: both stations receive the SAME pinned bytes, and neither
// can reach the other's copy. Sameness is structural — one frozen contract,
// cloned per station — and this is the executable statement of it.
func TestBothStationsReceiveTheSameSlotBytes(t *testing.T) {
	const payload = `{"anything":"one frozen value","n":3}`
	for _, slot := range slots {
		t.Run(slot.name, func(t *testing.T) {
			var atAlpha, atBeta Contract
			s := slotSpec(slot.set, payload)
			s.Alpha, s.Beta = slotWatchAlpha{&atAlpha}, slotWatchBeta{&atBeta}

			cl := mechClosure(t, s)
			gotAlpha, gotBeta := slot.get(atAlpha), slot.get(atBeta)
			if string(gotAlpha) != payload {
				t.Fatalf("alpha saw %s, want %s", gotAlpha, payload)
			}
			if string(gotBeta) != string(gotAlpha) {
				t.Fatalf("the stations saw different %s values:\n alpha: %s\n beta:  %s",
					slot.name, gotAlpha, gotBeta)
			}
			if string(slot.get(cl.Receipt.Record.Contract)) != payload {
				t.Fatalf("the record carries a third value: %s", slot.get(cl.Receipt.Record.Contract))
			}
			// Distinct backing arrays: alpha writing through its copy must not
			// reach beta's, or "same bytes" would be one shared mutable buffer.
			if len(gotAlpha) > 0 && &gotAlpha[0] == &gotBeta[0] {
				t.Fatalf("the two stations share one %s buffer", slot.name)
			}
		})
	}
}

// The boundary checked mechanically rather than asserted in a comment: the
// kernel's own source may not name any domain vocabulary belonging to what
// rides its opaque slots. A kernel that knew what `base_sha` or `acceptance`
// was would have stopped carrying opaque bytes, and no test of behaviour would
// notice — the leak is in the vocabulary, so the vocabulary is what is
// measured.
//
// Test files are excluded on purpose: a test that forbids a word has to be
// able to write it down.
func TestKernelNamesNoSlotVocabulary(t *testing.T) {
	forbidden := []*regexp.Regexp{
		regexp.MustCompile(`git\.snapshot`),
		regexp.MustCompile(`base_sha`),
		regexp.MustCompile(`\brepo(sitory|sitories)?\b`),
		regexp.MustCompile(`\bworktree\b`),
		regexp.MustCompile(`\bacceptance\b`),
		regexp.MustCompile(`\binvariants\b`),
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
				t.Errorf("%s names %q at byte %d: the kernel carries opaque bytes and must not learn what rides them",
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
