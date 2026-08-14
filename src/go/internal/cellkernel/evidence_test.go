package cellkernel

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// evidenceDir holds the committed closures of episodes that RENTED COGNITION.
// The shared corpus runs the deterministic fake only — a CI job that rented a
// provider would be the provider service this project does not build — so these
// files are the only standing record that the cognitive path ever ran.
//
// A committed artifact nothing checks is decoration. Two things make these
// evidence: each must pass VerifyClosure, the runtime's one scope-lift
// verification boundary; and each value the record carries must MOVE the
// digest, so a closure cannot be edited into a receipt for an episode that did
// not happen.
//
// The check lives here, in the kernel's own package, deliberately. Canonical
// bytes are the kernel's definition, and a checker that recomputed them
// elsewhere would be a second authority free to drift from the one that
// produced them — it would be testing its own copy. An earlier attempt did
// exactly that in Python and failed all six files on encoding/json's HTML
// escaping, which is a difference between two implementations and not a fact
// about any episode.
const evidenceDir = "../../../../docs/architecture/evidence/wcc-0.1"

func evidenceFiles(t *testing.T) []string {
	t.Helper()
	names, err := filepath.Glob(filepath.Join(evidenceDir, "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	// A vanished or renamed directory must not read as a suite that passed —
	// the same vacuity class as a corpus with no fixtures in it.
	if len(names) < 6 {
		t.Fatalf("expected at least 6 committed episodes in %s, found %d", evidenceDir, len(names))
	}
	return names
}

// loadEvidence decodes one committed closure. It re-reads per call rather than
// handing out a shared value: the mutation table below writes through slices
// the record owns, and a clone written once is a clone one careless mutation
// makes silently wrong for every later case.
func loadEvidence(t *testing.T, path string) Closure {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var cl Closure
	if err := json.Unmarshal(raw, &cl); err != nil {
		t.Fatalf("%s: %v", filepath.Base(path), err)
	}
	return cl
}

// otherJSON returns a valid document that is not raw, for mutating an opaque
// slot. Whitespace will NOT do: encoding/json compacts a RawMessage on the way
// out, so a space inserted into a slot is erased before the digest ever sees
// it. A whitespace mutation would have reported "the digest does not bind this
// slot" while proving only that json.Marshal normalizes.
func otherJSON(raw json.RawMessage) json.RawMessage {
	return json.RawMessage("[" + string(raw) + "]")
}

// mutations names every value an episode record carries, paired with a change
// to it that is visible in JSON by construction — append, wrap, or negate.
// The table IS the claim: a value that stops being bound fails here, rather
// than being discovered when a receipt is trusted for something it never
// proved.
var mutations = map[string]func(*EpisodeRecord){
	"canon":                  func(r *EpisodeRecord) { r.Canon += "X" },
	"episode_id":             func(r *EpisodeRecord) { r.EpisodeID += "X" },
	"execution_mode":         func(r *EpisodeRecord) { r.Mode = ModeMechanical },
	"resolved_spec.version":  func(r *EpisodeRecord) { r.ResolvedSpec.Version += "X" },
	"resolved_spec.protocol": func(r *EpisodeRecord) { r.ResolvedSpec.DeclaredProtocol += "X" },
	"resolved_spec.alpha":    func(r *EpisodeRecord) { r.ResolvedSpec.Alpha = otherJSON(r.ResolvedSpec.Alpha) },
	"resolved_spec.beta":     func(r *EpisodeRecord) { r.ResolvedSpec.Beta = otherJSON(r.ResolvedSpec.Beta) },
	"contract.id":            func(r *EpisodeRecord) { r.Contract.ID += "X" },
	"contract.goal":          func(r *EpisodeRecord) { r.Contract.Goal += "X" },
	"contract.issue":         func(r *EpisodeRecord) { r.Contract.Issue = otherJSON(r.Contract.Issue) },
	"contract.design":        func(r *EpisodeRecord) { r.Contract.Design = otherJSON(r.Contract.Design) },
	"contract.subject":       func(r *EpisodeRecord) { r.Contract.Subject = otherJSON(r.Contract.Subject) },
	"contract.required_evidence": func(r *EpisodeRecord) {
		r.Contract.RequiredEvidence = append(r.Contract.RequiredEvidence, RequiredRef{ID: "x", Kind: "x", Producer: RoleAlpha})
	},
	"alpha.execution_id":  func(r *EpisodeRecord) { r.Alpha.ExecutionID += "X" },
	"alpha.artifact.id":   func(r *EpisodeRecord) { r.Alpha.Artifacts[0].ID += "X" },
	"alpha.artifact.kind": func(r *EpisodeRecord) { r.Alpha.Artifacts[0].Kind += "X" },
	"alpha.artifact.text": func(r *EpisodeRecord) { r.Alpha.Artifacts[0].Text += "X" },
	"matter.data":         func(r *EpisodeRecord) { r.Matter.Data += "X" },
	"beta.execution_id":   func(r *EpisodeRecord) { r.Beta.ExecutionID += "X" },
	"review.pass":         func(r *EpisodeRecord) { r.Review.Pass = !r.Review.Pass },
	"review.notes":        func(r *EpisodeRecord) { r.Review.Notes += "X" },
	"review.unit":         func(r *EpisodeRecord) { r.Review.Assessment[0].Unit += "X" },
	"review.disposition":  func(r *EpisodeRecord) { r.Review.Assessment[0].Disposition = DispositionFinding },
	"review.unit.reason":  func(r *EpisodeRecord) { r.Review.Assessment[0].Reason += "X" },
	"review.unit.citations": func(r *EpisodeRecord) {
		r.Review.Assessment[0].Citations = append(r.Review.Assessment[0].Citations, "x")
	},
}

func TestCommittedEpisodesVerifyAndBindEveryValueTheyClaim(t *testing.T) {
	for _, path := range evidenceFiles(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			cl := loadEvidence(t, path)
			rec := cl.Receipt.Record

			// The record's own contract and resolved spec stand in for the
			// parent's trusted values, which makes VerifyClosure's binding
			// checks (2) vacuous HERE and nothing else: the digest recompute
			// and the whole derivation chain — verdict ← V, decision ← δ,
			// status ← lift, repair — are re-derived from scratch. The
			// mutation table below is what covers the binding class.
			meta := RunMeta{ResolvedSpec: rec.ResolvedSpec, ExecutionMode: rec.Mode}
			if err := VerifyClosure(rec.Contract, meta, cl); err != nil {
				t.Fatalf("committed closure does not verify: %v", err)
			}

			// These episodes are the cognitive evidence. A mechanical one
			// filed here would make the directory prove the opposite of what
			// it exists to prove, and would also make the execution_mode
			// mutation below a no-op.
			if rec.Mode != ModeCognitive {
				t.Fatalf("execution_mode = %q; this directory records episodes that rented cognition", rec.Mode)
			}
			if cl.Status != Accepted {
				t.Fatalf("status = %q; every committed episode here closed accepted", cl.Status)
			}

			for label, mutate := range mutations {
				t.Run(label, func(t *testing.T) {
					m := loadEvidence(t, path).Receipt.Record
					mutate(&m)
					if sha256hex(m.canonicalBytes()) == cl.Receipt.ScopeLiftDigest {
						t.Fatalf("the scope-lift digest does not bind %s: changing it left the digest unmoved", label)
					}
				})
			}
		})
	}
}
