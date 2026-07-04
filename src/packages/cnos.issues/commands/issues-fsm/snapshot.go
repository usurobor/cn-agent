package issuesfsm

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// FactSnapshot is the evaluator's entire input: an explicit, hidden-inference-
// free observation of one issue's state. Every field here is named directly
// in the issue's "Fact model" section (cnos#568). Transition logic (table.go)
// consumes a FactSnapshot; it never re-derives facts by guessing.
//
// Two assembly paths produce a FactSnapshot: (a) live, from the GitHub REST
// API + local git (fetch.go); (b) fixture, from a JSON file with exactly this
// shape (LoadFixture below) — used by --fixture and by the hermetic tests
// covering AC3-AC7.
type FactSnapshot struct {
	// Issue is the issue number this snapshot describes.
	Issue int `json:"issue"`

	// Labels is the issue's current label set, verbatim (e.g. "status:review",
	// "dispatch:cell", "protocol:cds"). CurrentState derives the CDS status
	// from the "status:" prefix, mirroring the convention already used by
	// `cn issues map` (issuesmap.go's toRecord).
	Labels []string `json:"labels"`

	// RunState is the most recently observed state of the issue's dispatch
	// workflow run: "" (none observed), "queued", "in_progress", or
	// "completed".
	RunState string `json:"run_state"`

	// RunConclusion is set when RunState == "completed": "success",
	// "failure", "cancelled", etc. (GitHub Actions run conclusion values).
	RunConclusion string `json:"run_conclusion,omitempty"`

	// BranchExists reports whether cycle/{issue} exists.
	BranchExists bool `json:"branch_exists"`

	// CommitsBeyondBase is the number of commits on cycle/{issue} beyond its
	// base (main). Zero means the branch (if it exists) has an empty diff.
	CommitsBeyondBase int `json:"commits_beyond_base"`

	// PRExists reports whether a pull request for the issue exists.
	PRExists bool `json:"pr_exists"`

	// PRCommitCount is the number of commits on the PR beyond its base.
	PRCommitCount int `json:"pr_commit_count"`

	// CDDArtifacts lists the filenames present under .cdd/unreleased/{issue}/
	// (e.g. "gamma-scaffold.md", "self-coherence.md", "delta-repair.md").
	CDDArtifacts []string `json:"cdd_artifacts,omitempty"`

	// ReviewRequestPresent reports whether REVIEW-REQUEST.yml is present
	// under .cdd/unreleased/{issue}/.
	ReviewRequestPresent bool `json:"review_request_present"`

	// ClaimRequestPresent reports whether CLAIM-REQUEST.yml is present
	// under .cdd/unreleased/{issue}/ (cnos#575 AC1). Mirrors the
	// ReviewRequestPresent / REVIEW-REQUEST.yml pattern: an explicit marker
	// the dispatch wake writes before requesting todo -> in-progress via
	// `cn issues fsm evaluate --apply`, so an incidental/offline evaluate
	// call against a status:todo issue never itself proposes a claim.
	ClaimRequestPresent bool `json:"claim_request_present,omitempty"`

	// BlockRequestPresent reports whether BLOCK-REQUEST.yml is present
	// under .cdd/unreleased/{issue}/ (cnos#575 AC2) — the explicit
	// STOP/escalation evidence a hard-block transition requires. Same
	// marker-file shape as ReviewRequestPresent/ClaimRequestPresent: a
	// deterministic, observable presence check rather than comment-text
	// parsing.
	BlockRequestPresent bool `json:"block_request_present,omitempty"`

	// ReleaseRequestPresent reports whether RELEASE-REQUEST.yml is present
	// under .cdd/unreleased/{issue}/ (cnos#575 AC3) — the worker's explicit
	// request to release its claim pre-work (in-progress -> todo). Distinct
	// from the pre-existing dead-run reconciliation rules (table.go's
	// "in-progress" rules, all_false: [run_active]): those fire on a sweep
	// finding an ABANDONED issue; this fires synchronously on the
	// requesting wake's own still-active run releasing its claim before
	// producing matter (cnos#368 protection — never blind-requeue if matter
	// already exists).
	ReleaseRequestPresent bool `json:"release_request_present,omitempty"`

	// RepairContractPresent reports whether a repair contract / prior beta
	// or delta findings are present (e.g. delta-repair.md, or a beta-review
	// with an "iterate" verdict) — the repair-re-entry context required by
	// cnos#516 for a changes -> todo transition.
	RepairContractPresent bool `json:"repair_contract_present"`

	// ChecksState is the most recently observed CI/checks state where
	// available: "", "pending", "success", or "failure".
	ChecksState string `json:"checks_state,omitempty"`

	// CellKind is the optional cell-kind seam (cnos#568 operator note;
	// cnos#570 defines the canonical taxonomy). Phase 1 OBSERVES it but no
	// transition rule (table.go) consumes it — it is a forward-compat seam,
	// not enforcement, so its presence or absence cannot change any Phase-1
	// decision. When nothing is observed, normalizeCellKind records
	// source="absent" and defaults the kind to "implementation" (the current
	// CDS dispatch cell kind), so cnos#570 can plug the taxonomy in later
	// without redesigning this fact model.
	CellKind CellKind `json:"cell_kind"`
}

// CellKind is the observation half of the cell-kind seam. The FSM fact model
// must be able to OBSERVE a cell's kind regardless of where cnos#570 later
// decides to record it (issue body, cell.yml, or receipt) — this struct is
// that observation slot. Phase 1 never sets Observed (no source is parsed
// yet); it only records the defaulting.
type CellKind struct {
	// Observed is the cell_kind read from repo state, or "" if none found.
	Observed string `json:"observed,omitempty"`
	// Source records where Observed came from: "absent", "issue_body",
	// "cdd_artifact", or "inferred".
	Source string `json:"source,omitempty"`
	// DefaultedTo is the kind Phase 1 falls back to when Observed is "".
	DefaultedTo string `json:"defaulted_to,omitempty"`
}

// normalizeCellKind fills the seam's defaults explicitly (rather than
// silently): when nothing was observed, Phase 1 records source="absent" and
// defaults the kind to "implementation". Idempotent — safe to call on any
// assembled snapshot. No transition rule reads these fields.
func (f *FactSnapshot) normalizeCellKind() {
	if f.CellKind.Observed == "" {
		if f.CellKind.Source == "" {
			f.CellKind.Source = "absent"
		}
		if f.CellKind.DefaultedTo == "" {
			f.CellKind.DefaultedTo = "implementation"
		}
	}
}

// CurrentState derives the CDS status from f.Labels by the "status:" prefix
// convention (the same convention `cn issues map`'s toRecord uses). Returns
// "" if no status:* label is present.
func CurrentState(f FactSnapshot) string {
	for _, l := range f.Labels {
		if strings.HasPrefix(l, "status:") {
			return strings.TrimPrefix(l, "status:")
		}
	}
	return ""
}

// LoadFixture reads a FactSnapshot from a JSON file — the offline path used
// by --fixture and by every AC3-AC7 test.
func LoadFixture(path string) (FactSnapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FactSnapshot{}, fmt.Errorf("read fixture: %w", err)
	}
	var snap FactSnapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return FactSnapshot{}, fmt.Errorf("parse fixture JSON: %w", err)
	}
	snap.normalizeCellKind()
	return snap, nil
}
