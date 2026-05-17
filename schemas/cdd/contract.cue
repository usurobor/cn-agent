// see schemas/cdd/README.md §Scope-Lift Invariant.
//
// Schema ID: cnos.cdd.contract.v1
//
// Types the cycle's accountability surface — what V validates the receipt
// against. The contract is sourced from the issue body at the SHA committed
// to in the issue's `## Source of truth` table (see
// src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md §Q5 / §Input
// contract). The receipt is computed; the contract is authored. This cycle
// types refs; Phase 3 implements derive_receipt and V.
//
// v1 is the *schema-artifact* version, not a CDD protocol version.

package cdd

// #Contract is V's first input — the accountability surface for the cycle.
// Fields are reference-shaped: V reads refs and resolves them against the
// evidence graph rather than carrying inline copies (RECEIPT-VALIDATION.md
// §Input contract).
#Contract: {
	// contract_id: stable identifier for the contract artifact, scoped to
	// the cycle. Typically `cnos.cdd.contract.v1:#{issue}`.
	contract_id: string

	// issue_ref: a ref to the authoring substrate — the issue body. Typically
	// of the form `issue:#{N}` or a GitHub URL.
	issue_ref: string

	// branch: the cycle branch name (e.g. `cycle/369`).
	branch: string

	// base_sha: the merge-base on origin/main from which the cycle branched.
	// The receipt's diff_ref resolves to git_diff(base_sha..merge_sha).
	base_sha: string

	// contract_sha: optional pin of the issue-body SHA the contract is read
	// at — when the issue body is edited mid-cycle, contract_sha records
	// which version of the body V is checking against. RECEIPT-VALIDATION.md
	// §Input contract names this as the SHA committed to in `## Source of
	// truth`.
	contract_sha?: string

	// scope: in/out/deferred surfaces named by the issue.
	scope: {
		in_scope: [...string]
		out_of_scope: [...string]
		deferred?: [...string]
	}

	// acceptance_oracle_refs: refs into the contract that resolve to the
	// executable predicates V evaluates. Each entry typically names one AC's
	// oracle (e.g. `issue:#369#AC1.oracle`). The predicate set is
	// contract-derived; V does not invent predicates.
	acceptance_oracle_refs: [...string]

	// non_goals: surfaces the cycle must not touch. V can use these to flag
	// scope drift in the diff.
	non_goals: [...string]

	// source_of_truth_refs: refs to canonical artifacts the contract cites.
	// Used by V when an AC oracle requires resolution against a canonical
	// document (e.g. `RECEIPT-VALIDATION.md` for #369).
	source_of_truth_refs: [...string]
}
