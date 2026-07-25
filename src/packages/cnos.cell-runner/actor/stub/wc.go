package stub

import (
	"fmt"
	"os"
	"strings"

	"cnos.dev/cnos/cell-runner/kernel"
)

// wcActor is the Working Cell: telos = artifact. Its α edits the target file;
// its V asserts the contract's mechanical acceptance predicates against the
// edited file.
type wcActor struct{ baseActor }

// NewWC returns a deterministic WC actor.
func NewWC() kernel.Actor { return wcActor{} }

// Alpha appends the required `## Coherence` section to the target file.
func (wcActor) Alpha(rc *kernel.RunContext) error {
	content, err := os.ReadFile(rc.Target)
	if err != nil {
		return fmt.Errorf("read target: %w", err)
	}
	body := string(content)
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	body += "\n" + CoherenceHeading + "\n\n" + coherenceBody
	if err := os.WriteFile(rc.Target, []byte(body), 0o644); err != nil {
		return fmt.Errorf("write target: %w", err)
	}
	rc.Facts.ArtifactPresent = true
	rc.AddEvidence("artifact", rc.Target)
	rc.Log("  α produce: appended %q section to target", CoherenceHeading)
	return nil
}

// Beta approves unless the beta-reject fault is injected.
func (wcActor) Beta(rc *kernel.RunContext) error {
	if rc.FailInjection == "beta-reject" {
		rc.Facts.ReviewRejected = true
		rc.AddEvidence("beta_review", "rejected")
		rc.Log("  β review: REJECT (fault injected: beta-reject)")
		return nil
	}
	rc.Facts.ReviewApproved = true
	rc.AddEvidence("beta_review", "approved")
	rc.Log("  β review: approve (artifact matches the required shape)")
	return nil
}

// Validate mechanically asserts the acceptance predicates: the `## Coherence`
// section is present and well-formed. The v-fail fault forces a rejection.
func (wcActor) Validate(rc *kernel.RunContext) error {
	content, err := os.ReadFile(rc.Target)
	if err != nil {
		return fmt.Errorf("read target for validation: %w", err)
	}
	if rc.FailInjection == "v-fail" {
		setValidation(rc, false, []string{"coherence_section_present (fault injected: v-fail)"})
		rc.Log("  V validate: FAIL (fault injected: v-fail)")
		return nil
	}
	var failed []string
	if !hasCoherenceSection(string(content)) {
		failed = append(failed, "coherence_section_present")
	}
	if !coherenceSectionWellFormed(string(content)) {
		failed = append(failed, "coherence_section_wellformed")
	}
	if len(failed) > 0 {
		setValidation(rc, false, failed)
		rc.Log("  V validate: FAIL (%s)", strings.Join(failed, ", "))
		return nil
	}
	setValidation(rc, true, nil)
	rc.AddEvidence("v_validate", "acceptance predicates satisfied")
	rc.Log("  V validate: PASS (section present and well-formed)")
	return nil
}
