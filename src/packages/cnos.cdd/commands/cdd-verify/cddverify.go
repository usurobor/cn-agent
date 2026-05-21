package cddverify

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ValidateOptions parameterizes a V invocation.
type ValidateOptions struct {
	ReceiptPath    string // required
	ContractPath   string // optional
	RepoRoot       string // resolved (caller may pass git rev-parse --show-toplevel)
	StructuralOnly bool   // skip filesystem-dereference rules (C1/C2/C4)
}

// Validate runs V over (contract, receipt) → ValidationVerdict.
//
// Per CCNF: V's signature is Contract × Receipt → ValidationVerdict. The
// contract path is optional in practice; V's predicates are derived from
// the receipt + schemas. The contract reference is recorded in
// provenance.input_refs.contract_ref but a missing contract surfaces only
// as a warning.
//
// StructuralOnly=true runs only CUE structural validation and the
// counterfeit rules that don't require filesystem dereference (C3, C5).
func Validate(ctx context.Context, opts ValidateOptions) (ValidationVerdict, error) {
	var fail []FailedPredicate
	var warn []string

	// Parse receipt.
	receipt, err := ReadReceipt(ctx, opts.ReceiptPath)
	if err != nil {
		return emitVerdict(ResultFAIL,
			[]FailedPredicate{{
				Predicate:   "v.receipt_parse",
				EvidenceRef: opts.ReceiptPath,
				Diagnostic:  fmt.Sprintf("could not parse receipt YAML: %s", err.Error()),
			}},
			nil,
			opts.ReceiptPath, opts.ContractPath, "",
		), nil
	}
	if receipt == nil {
		return emitVerdict(ResultFAIL,
			[]FailedPredicate{{
				Predicate:   "v.receipt_shape",
				EvidenceRef: opts.ReceiptPath,
				Diagnostic:  "receipt is not a YAML mapping at top level",
			}},
			nil,
			opts.ReceiptPath, opts.ContractPath, "",
		), nil
	}

	// Determine evidence_root for provenance.
	var evidenceRoot string
	if er, ok := evidenceRefs(receipt); ok {
		if s, ok := asString(er, "evidence_root"); ok {
			evidenceRoot = stripAtRef(s)
		}
	}

	// Dispatch on protocol_id.
	pid := protocolID(receipt)
	switch {
	case pid == "":
		fail = append(fail, FailedPredicate{
			Predicate:   "v.dispatch.missing_protocol_id",
			EvidenceRef: "receipt.protocol_id",
			Diagnostic:  "receipt has no protocol_id; V cannot dispatch (per #388 README §protocol_id dispatch convention).",
		})
	default:
		entry, ok := DispatchTable[pid]
		if !ok {
			fail = append(fail, FailedPredicate{
				Predicate:   "v.dispatch.unknown_protocol_id",
				EvidenceRef: fmt.Sprintf("receipt.protocol_id: %s", pid),
				Diagnostic:  fmt.Sprintf("unknown protocol_id: %s. Known: %v", pid, KnownProtocols()),
			})
		} else {
			// AC2 / AC3 / AC4 / AC5 — CUE structural validation.
			ok, diagnostics := cueVet(ctx, opts.ReceiptPath, entry, opts.RepoRoot)
			if !ok {
				for _, d := range diagnostics {
					fail = append(fail, FailedPredicate{
						Predicate:   "cue.structural_vet",
						EvidenceRef: fmt.Sprintf("protocol_id: %s (%s)", pid, entry.Definition),
						Diagnostic:  d,
					})
				}
			}
		}
	}

	// AC6 / AC8 — counterfeit-receipt rules.
	// C5 first (friendly diagnostic; runs even when CUE fails). Filesystem-free.
	ruleC5ProtocolIDMismatch(receipt, &fail)
	// C3 — override does not rewrite verdict. Filesystem-free.
	ruleC3OverrideDoesNotRewrite(receipt, &fail)
	if !opts.StructuralOnly {
		// C4 — evidence deref.
		ruleC4EvidenceDeref(receipt, opts.RepoRoot, &fail)
		// C1 — actor separation (CDS only).
		ruleC1ActorSeparation(ctx, receipt, opts.RepoRoot, &fail, &warn)
		// C2 — β verdict precedes merge (CDS only).
		ruleC2VerdictPrecedesMerge(ctx, receipt, opts.RepoRoot, &fail, &warn)
	} else {
		warn = append(warn,
			"V running in --structural-only mode: skipping evidence-ref "+
				"dereference (C4), actor-separation (C1), and "+
				"verdict-precedes-merge (C2) rules.",
		)
	}

	// Contract presence (warn-only).
	if opts.ContractPath != "" {
		cp := opts.ContractPath
		if !filepath.IsAbs(cp) {
			cp = filepath.Join(opts.RepoRoot, opts.ContractPath)
		}
		if _, err := os.Stat(cp); err != nil {
			warn = append(warn, fmt.Sprintf(
				"contract_ref points at non-existent path: %s; V proceeds (contract presence is not load-bearing).",
				opts.ContractPath,
			))
		}
	}

	result := ResultPASS
	if len(fail) > 0 {
		result = ResultFAIL
	}
	return emitVerdict(result, fail, warn, opts.ReceiptPath, opts.ContractPath, evidenceRoot), nil
}

// emitVerdict builds a ValidationVerdict struct with provenance filled in.
func emitVerdict(result string, fail []FailedPredicate, warn []string, receiptPath, contractPath, evidenceRoot string) ValidationVerdict {
	if fail == nil {
		fail = []FailedPredicate{}
	}
	if warn == nil {
		warn = []string{}
	}
	refs := InputRefs{ReceiptRef: receiptPath}
	if contractPath != "" {
		refs.ContractRef = contractPath
	}
	if evidenceRoot != "" {
		refs.EvidenceRootRef = evidenceRoot
	}
	return ValidationVerdict{
		Result:           result,
		FailedPredicates: fail,
		Warnings:         warn,
		Provenance: Provenance{
			ValidatorIdentity: ValidatorIdentity,
			ValidatorVersion:  ValidatorVersion,
			CheckedAt:         time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			InputRefs:         refs,
		},
	}
}

// sortFailedPredicatesByPredicate is a deterministic sort applied where
// the caller wants stable diagnostic ordering. Not used in the default
// path (we preserve the rule-firing order which matches Python), but
// available for tests.
func sortFailedPredicatesByPredicate(fps []FailedPredicate) {
	sort.SliceStable(fps, func(i, j int) bool {
		return fps[i].Predicate < fps[j].Predicate
	})
}
