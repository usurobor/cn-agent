package stub

import (
	"fmt"

	"cnos.dev/cnos/cell-runner/kernel"
	"cnos.dev/cnos/cell-runner/model"
)

// pcActor is the Planning Cell: telos = relation_graph. Its α emits a wave
// (one WC contract) from the measurement; its V asserts the wave is well-formed.
type pcActor struct{ baseActor }

// NewPC returns a deterministic PC actor.
func NewPC() kernel.Actor { return pcActor{} }

// Alpha emits a wave carrying exactly one WC contract derived from the CM: a
// contract to add the missing `## Coherence` section to the target.
func (pcActor) Alpha(rc *kernel.RunContext) error {
	if rc.CM == nil {
		return fmt.Errorf("PC planning requires a CM but none was provided")
	}
	wc := model.Contract{
		ID:          "wc-1",
		Telos:       "artifact",
		Target:      rc.CM.Target,
		Instruction: fmt.Sprintf("add the %q section to the target to clear the measured defect", CoherenceHeading),
		Acceptance:  []string{"coherence_section_present", "coherence_section_wellformed"},
	}
	rc.Wave = &model.Wave{
		ID:        "wave-1",
		FromCMRef: rc.CM.MeasuredAtRef,
		Contracts: []model.Contract{wc},
	}
	rc.Facts.ArtifactPresent = true
	rc.AddEvidence("artifact", rc.Ref+"/wave")
	rc.Log("  α produce: planned wave of %d WC contract(s) from the CM", len(rc.Wave.Contracts))
	return nil
}

// Beta approves the planned wave.
func (pcActor) Beta(rc *kernel.RunContext) error {
	rc.Facts.ReviewApproved = true
	rc.AddEvidence("beta_review", "approved")
	rc.Log("  β review: approve (wave targets the measured defect)")
	return nil
}

// Validate asserts the wave is well-formed: exactly one contract, telos artifact.
func (pcActor) Validate(rc *kernel.RunContext) error {
	var failed []string
	if rc.Wave == nil || len(rc.Wave.Contracts) != 1 {
		failed = append(failed, "wave_has_exactly_one_contract")
	} else if rc.Wave.Contracts[0].Telos != "artifact" {
		failed = append(failed, "wave_contract_telos_is_artifact")
	}
	if len(failed) > 0 {
		setValidation(rc, false, failed)
		rc.Log("  V validate: FAIL (%v)", failed)
		return nil
	}
	setValidation(rc, true, nil)
	rc.AddEvidence("v_validate", "wave well-formed")
	rc.Log("  V validate: PASS (wave well-formed: one WC contract, telos=artifact)")
	return nil
}
