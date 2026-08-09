package cellkernel

import (
	"errors"
	"strings"
)

// VerifyReceipt mechanically re-verifies a receipt from its serialized content
// alone — the check a parent or another runtime runs on a receipt it received,
// with no access to the original EpisodeRecord (Pi #33 D2). It recomputes every
// internal hash, the content-addressed evidence, id uniqueness, and required-
// evidence presence + producer authority. It does not re-judge β.
//
// It returns nil when the receipt is internally self-consistent and its evidence
// re-hashes; otherwise an error naming every failure.
func VerifyReceipt(rc Receipt) error {
	var failed []string
	add := func(cond bool, msg string) {
		if cond {
			failed = append(failed, msg)
		}
	}

	add(rc.EpisodeID == "", "missing episode_id")
	add(rc.AlphaExecutionID == "", "missing alpha_execution_id")
	add(rc.BetaExecutionID == "", "missing beta_execution_id")
	add(rc.PolicyID == "", "missing beta_input_policy_id")
	add(rc.BetaInputHash == "", "missing beta_input_hash")

	add(hashJSON(rc.Contract) != rc.ContractHash, "contract hash mismatch")
	add(hashJSON(rc.Matter) != rc.MatterHash, "matter hash mismatch")
	add(hashJSON(rc.Review) != rc.ReviewHash, "review hash mismatch")
	add(hashJSON(rc.Evidence) != rc.EvidenceHash, "evidence hash mismatch")

	failed = append(failed, checkEvidence(rc.Contract, rc.Evidence)...)

	if len(failed) > 0 {
		return errors.New("receipt verification failed: " + strings.Join(failed, "; "))
	}
	return nil
}
