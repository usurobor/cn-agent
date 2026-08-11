package stub

import (
	"fmt"
	"os"

	"cnos.dev/cnos/cell-runner/kernel"
	"cnos.dev/cnos/cell-runner/model"
)

// ccActor is the Coherence Cell: telos = judgment. Its α IS the Measure step
// (it produces the CM object); its δ emits the judgment the loop routes on.
type ccActor struct{ baseActor }

// NewCC returns a deterministic CC actor.
func NewCC() kernel.Actor { return ccActor{} }

// Alpha is the CC's Measure step: it mechanically measures the target and
// produces the CM object. A missing `## Coherence` section is one defect.
func (ccActor) Alpha(rc *kernel.RunContext) error {
	content, err := os.ReadFile(rc.Target)
	if err != nil {
		return fmt.Errorf("read target for measurement: %w", err)
	}
	defects := []model.Defect{}
	if !hasCoherenceSection(string(content)) {
		defects = append(defects, model.Defect{
			ID:     "missing-required-section",
			Kind:   "missing-required-section",
			Detail: fmt.Sprintf("required %q section is absent from the target", CoherenceHeading),
		})
	}
	score := 1.0
	if len(defects) > 0 {
		score = 0.0
	}
	rc.CM = &model.CM{
		Target:     rc.Target,
		AlphaScore: score,
		Defects:    defects,
		Provenance: model.Provenance{
			MeasuredBy: "cc.stub.measure",
			Method:     "mechanical-required-section-check",
			CMVersion:  "v0",
		},
		MeasuredAtRef: rc.Ref + "/cm",
	}
	rc.Facts.ArtifactPresent = true
	rc.AddEvidence("cm", rc.CM.MeasuredAtRef)
	rc.Log("  α measure: CM produced — defects=%d alpha_score=%.2f (incoherence=%d)",
		len(defects), score, rc.CM.Incoherence())
	return nil
}

// Beta approves the measurement.
func (ccActor) Beta(rc *kernel.RunContext) error {
	rc.Facts.ReviewApproved = true
	rc.AddEvidence("beta_review", "approved")
	rc.Log("  β review: approve (measurement is well-formed)")
	return nil
}

// Validate asserts the CM is structurally well-formed.
func (ccActor) Validate(rc *kernel.RunContext) error {
	var failed []string
	if rc.CM == nil || rc.CM.Target == "" {
		failed = append(failed, "cm_has_target")
	}
	if rc.CM != nil && rc.CM.Provenance.MeasuredBy == "" {
		failed = append(failed, "cm_has_provenance")
	}
	if len(failed) > 0 {
		setValidation(rc, false, failed)
		rc.Log("  V validate: FAIL (%v)", failed)
		return nil
	}
	setValidation(rc, true, nil)
	rc.AddEvidence("v_validate", "CM well-formed")
	rc.Log("  V validate: PASS (CM well-formed)")
	return nil
}

// DeltaAccept finalizes the receipt (via the shared accept path) AND emits the
// CC's judgment computed from the CM: no defects → coherent, else
// request_planning.
func (ccActor) DeltaAccept(rc *kernel.RunContext) error {
	if err := acceptReceipt(rc); err != nil {
		return err
	}
	verdict := "coherent"
	if rc.CM != nil && rc.CM.Incoherence() > 0 {
		verdict = "request_planning"
	}
	rc.Judgment = &model.Judgment{
		Verdict:      verdict,
		CMRef:        rc.CM.MeasuredAtRef,
		Target:       rc.CM.Target,
		DecidedAtRef: rc.Ref + "/judgment",
	}
	rc.AddEvidence("judgment", rc.Judgment.DecidedAtRef)
	rc.Log("  δ judgment: %s", verdict)
	return nil
}
