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

	// RepairContractPresent reports whether a repair contract / prior beta
	// or delta findings are present (e.g. delta-repair.md, or a beta-review
	// with an "iterate" verdict) — the repair-re-entry context required by
	// cnos#516 for a changes -> todo transition.
	RepairContractPresent bool `json:"repair_contract_present"`

	// ChecksState is the most recently observed CI/checks state where
	// available: "", "pending", "success", or "failure".
	ChecksState string `json:"checks_state,omitempty"`
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
	return snap, nil
}
