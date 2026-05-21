package cddverify

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// placeholderPatterns are substrings that mark a receipt value as a
// synthetic/placeholder ref V should NOT try to dereference.
//
// Mirrors the Python predecessor's PLACEHOLDER_PATTERNS so AC behavior is
// preserved (counterfeit rules must reject only real-looking paths).
var placeholderPatterns = []string{
	"git_diff(",
	"sha256:",
	"@sha256:",
	"issue:",
	"@merge_sha_placeholder",
	"@contract_sha_placeholder",
	"@draft",
	"@commit_sha_placeholder",
}

// isFilesystemPath heuristically decides whether s is a real-looking
// filesystem path V should try to dereference, vs. a placeholder/synthetic
// ref.
//
// Mirrors Python's is_filesystem_path exactly so the AC suite's counterfeit
// detection behaves identically.
func isFilesystemPath(s string) bool {
	if s == "" {
		return false
	}
	for _, ph := range placeholderPatterns {
		if strings.Contains(s, ph) {
			return false
		}
	}
	// Strip a trailing @<ref> placeholder if present.
	bare := s
	if i := strings.Index(s, "@"); i >= 0 {
		bare = s[:i]
	}
	if bare == "" {
		return false
	}
	// Real paths start with . or /.
	if strings.HasPrefix(bare, ".") || strings.HasPrefix(bare, "/") {
		return true
	}
	// Bare relative paths with a slash and a file-shaped extension count.
	if strings.Contains(bare, "/") {
		base := bare[strings.LastIndex(bare, "/")+1:]
		if strings.Contains(base, ".") {
			return true
		}
	}
	return false
}

// stripAtRef removes a trailing `@<ref>` placeholder from a path-like
// string (e.g. "foo/bar@sha256:abc" → "foo/bar").
func stripAtRef(s string) string {
	if i := strings.Index(s, "@"); i >= 0 {
		return s[:i]
	}
	return s
}

// isCollapsedMode returns true iff the receipt declares `mode: collapsed`
// at the top level. In collapsed-mode receipts, C1/C2 emit warnings
// instead of failed_predicates (the receipt acknowledges the structural
// degradation; V honestly surfaces it as a warning rather than rejecting
// it as a counterfeit).
//
// Per #389 design (wave-manifest pattern from cycles 375/377/378/388/389).
func isCollapsedMode(r Receipt) bool {
	return mode(r) == "collapsed"
}

// extractActor reads "Actor:" or equivalent metadata-block field from a
// closure-record markdown file's first ~60 lines. Returns "" if not
// found.
//
// Pure (eng/go §2.17): operates on file bytes after read. The os.ReadFile
// call is the IO boundary.
func extractActor(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	prefixes := []string{
		"actor:", "α actor:", "β actor:", "alpha actor:",
		"beta actor:", "author:", "alpha-actor:", "beta-actor:",
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 60 {
		lines = lines[:60]
	}
	for _, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))
		for _, p := range prefixes {
			if strings.HasPrefix(lower, p) {
				// Return the original (non-lowered) tail after the first ":".
				idx := strings.Index(line, ":")
				if idx < 0 {
					return ""
				}
				return strings.TrimSpace(line[idx+1:])
			}
		}
	}
	return ""
}

// parseISO parses an ISO 8601 timestamp. Returns (zero, false) on failure.
//
// Normalizes a trailing "Z" to "+00:00" before parsing.
func parseISO(ts string) (time.Time, bool) {
	if ts == "" {
		return time.Time{}, false
	}
	s := strings.TrimSpace(ts)
	if strings.HasSuffix(s, "Z") {
		s = s[:len(s)-1] + "+00:00"
	}
	for _, layout := range []string{
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05.000-07:00",
		"2006-01-02T15:04:05.000000-07:00",
		time.RFC3339,
		time.RFC3339Nano,
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// ruleC4EvidenceDeref — for each evidence_refs entry that looks like a real
// path, assert the file exists at repoRoot/<path>.
func ruleC4EvidenceDeref(r Receipt, repoRoot string, fail *[]FailedPredicate) {
	er, ok := evidenceRefs(r)
	if !ok {
		return
	}
	// Deterministic key iteration order.
	keys := sortedKeys(er)
	for _, key := range keys {
		value := er[key]
		// Each value is either a single ref or a list of refs.
		var candidates []string
		switch v := value.(type) {
		case string:
			candidates = []string{v}
		case []any:
			for _, item := range v {
				if s, ok := item.(string); ok {
					candidates = append(candidates, s)
				}
			}
		default:
			continue
		}
		for _, ref := range candidates {
			if !isFilesystemPath(ref) {
				continue
			}
			bare := stripAtRef(ref)
			var full string
			if filepath.IsAbs(bare) {
				full = bare
			} else {
				full = filepath.Join(repoRoot, bare)
			}
			if _, err := os.Stat(full); err != nil {
				*fail = append(*fail, FailedPredicate{
					Predicate:   "counterfeit.evidence_ref_unresolved",
					EvidenceRef: fmt.Sprintf("%s: %s", key, ref),
					Diagnostic:  fmt.Sprintf("evidence_ref '%s' points at non-existent path: %s", key, bare),
				})
			}
		}
	}
}

// ruleC1ActorSeparation — α actor ≠ β actor. Fires for CDS receipts only.
//
// When mode: collapsed, the rule emits a warning rather than a
// failed_predicate (degraded review-independence acknowledged by receipt).
func ruleC1ActorSeparation(ctx context.Context, r Receipt, repoRoot string, fail *[]FailedPredicate, warn *[]string) {
	if protocolID(r) != "cnos.cdd.cds.receipt.v1" {
		return
	}
	er, ok := evidenceRefs(r)
	if !ok {
		return
	}
	alphaPath, _ := asString(er, "alpha_closeout")
	betaPath, _ := asString(er, "beta_closeout")
	if alphaPath == "" || betaPath == "" {
		return // CUE already failed; let that diagnostic carry.
	}
	aBare := stripAtRef(alphaPath)
	bBare := stripAtRef(betaPath)
	if aBare == "" || bBare == "" {
		return
	}
	aFull := filepath.Join(repoRoot, aBare)
	bFull := filepath.Join(repoRoot, bBare)
	if _, err := os.Stat(aFull); err != nil {
		return // C4 already flagged it.
	}
	if _, err := os.Stat(bFull); err != nil {
		return
	}
	aActor := extractActor(aFull)
	if aActor == "" {
		aActor = gitLogAuthor(ctx, repoRoot, aBare)
	}
	bActor := extractActor(bFull)
	if bActor == "" {
		bActor = gitLogAuthor(ctx, repoRoot, bBare)
	}
	if aActor == "" || bActor == "" {
		*warn = append(*warn, fmt.Sprintf(
			"counterfeit.actor_separation: could not determine actors for %s or %s; rule skipped.",
			aBare, bBare,
		))
		return
	}
	if strings.EqualFold(strings.TrimSpace(aActor), strings.TrimSpace(bActor)) {
		msg := fmt.Sprintf(
			"α and β closeout authored by same actor: '%s'. "+
				"Structural review-independence requires distinct actors "+
				"(per RECEIPT-VALIDATION.md §Q4 and CCNF kernel doctrine).",
			aActor,
		)
		if isCollapsedMode(r) {
			*warn = append(*warn, fmt.Sprintf("counterfeit.actor_separation (collapsed-mode acknowledged): %s", msg))
		} else {
			*fail = append(*fail, FailedPredicate{
				Predicate:   "counterfeit.actor_separation",
				EvidenceRef: "alpha_closeout / beta_closeout",
				Diagnostic:  msg,
			})
		}
	}
}

// ruleC2VerdictPrecedesMerge — β closeout commit time T_β ≤
// boundary_decision.decided_at T_δ. CDS only. Action ∈ {accept, release}.
func ruleC2VerdictPrecedesMerge(ctx context.Context, r Receipt, repoRoot string, fail *[]FailedPredicate, warn *[]string) {
	if protocolID(r) != "cnos.cdd.cds.receipt.v1" {
		return
	}
	bnd, ok := boundary(r)
	if !ok {
		return
	}
	decidedAt, _ := asString(bnd, "decided_at")
	action, _ := asString(bnd, "action")
	if action != "accept" && action != "release" {
		return // constraint only matters for terminal-accept actions.
	}
	er, ok := evidenceRefs(r)
	if !ok {
		return
	}
	betaPath, _ := asString(er, "beta_closeout")
	if betaPath == "" {
		return
	}
	bBare := stripAtRef(betaPath)
	bFull := filepath.Join(repoRoot, bBare)
	if _, err := os.Stat(bFull); err != nil {
		return // C4 carries it.
	}
	tBetaStr := gitCommitTime(ctx, repoRoot, bBare)
	// Fallback: file mtime when git has no record.
	if tBetaStr == "" {
		fi, err := os.Stat(bFull)
		if err == nil {
			tBetaStr = fi.ModTime().UTC().Format("2006-01-02T15:04:05+00:00")
		}
	}
	tBeta, betaOK := parseISO(tBetaStr)
	tDelta, deltaOK := parseISO(decidedAt)
	if betaOK && deltaOK {
		if tBeta.After(tDelta) {
			msg := fmt.Sprintf(
				"β closeout commit T_β=%s postdates δ decision T_δ=%s. "+
					"β verdict must precede merge per RECEIPT-VALIDATION.md §Q1.",
				tBetaStr, decidedAt,
			)
			if isCollapsedMode(r) {
				*warn = append(*warn, fmt.Sprintf("counterfeit.verdict_precedes_merge (collapsed-mode acknowledged): %s", msg))
			} else {
				*fail = append(*fail, FailedPredicate{
					Predicate:   "counterfeit.verdict_precedes_merge",
					EvidenceRef: fmt.Sprintf("beta_closeout commit time: %s; boundary_decision.decided_at: %s", tBetaStr, decidedAt),
					Diagnostic:  msg,
				})
			}
			return
		}
	}
	// Fallback structural check: closeout exists, non-empty.
	if fi, err := os.Stat(bFull); err == nil && fi.Size() == 0 {
		*fail = append(*fail, FailedPredicate{
			Predicate:   "counterfeit.verdict_precedes_merge",
			EvidenceRef: fmt.Sprintf("beta_closeout: %s", bBare),
			Diagnostic:  "β closeout file is empty; no β verdict to precede δ accept/release.",
		})
	}
}

// ruleC3OverrideDoesNotRewrite — boundary_decision.override.original_validation_verdict.verdict
// must equal validation.verdict. All protocols.
func ruleC3OverrideDoesNotRewrite(r Receipt, fail *[]FailedPredicate) {
	bnd, ok := boundary(r)
	if !ok {
		return
	}
	action, _ := asString(bnd, "action")
	if action != "override" {
		return
	}
	ovd, ok := asMap(bnd, "override")
	if !ok {
		return
	}
	orig, ok := asMap(ovd, "original_validation_verdict")
	if !ok {
		return
	}
	declaredInOverride, hasOverride := orig["verdict"]
	val, ok := validationBlock(r)
	if !ok {
		return
	}
	declaredInValidation, hasVal := val["verdict"]
	if !hasOverride || !hasVal {
		return // CUE has structural responsibility here.
	}
	if !sameAny(declaredInOverride, declaredInValidation) {
		*fail = append(*fail, FailedPredicate{
			Predicate:   "counterfeit.override_does_not_rewrite",
			EvidenceRef: "boundary_decision.override.original_validation_verdict.verdict",
			Diagnostic: fmt.Sprintf(
				"override claims original verdict '%v' but receipt.validation.verdict is '%v'. "+
					"Override never rewrites the ValidationVerdict (RECEIPT-VALIDATION.md §Q4).",
				declaredInOverride, declaredInValidation,
			),
		})
	}
}

// ruleC5ProtocolIDMismatch — protocol_id declares cds/cdr but evidence_refs
// is missing required keys (friendlier diagnostic than raw cue vet).
func ruleC5ProtocolIDMismatch(r Receipt, fail *[]FailedPredicate) {
	pid := protocolID(r)
	er, ok := evidenceRefs(r)
	if !ok {
		// no evidence_refs map at all → CUE will flag; nothing for us.
		return
	}
	switch pid {
	case "cnos.cdd.cds.receipt.v1":
		var missing []string
		for _, k := range CDSRequiredEvidenceRefs {
			if _, present := er[k]; !present {
				missing = append(missing, k)
			}
		}
		if len(missing) > 0 {
			*fail = append(*fail, FailedPredicate{
				Predicate:   "counterfeit.protocol_id_mismatch",
				EvidenceRef: "evidence_refs",
				Diagnostic: fmt.Sprintf(
					"protocol_id declares cnos.cdd.cds.receipt.v1 but missing required CDS evidence_refs keys: %v",
					missing,
				),
			})
		}
	case "cnos.cdd.cdr.receipt.v1":
		var missing []string
		for _, k := range CDRRequiredEvidenceRefs {
			if _, present := er[k]; !present {
				missing = append(missing, k)
			}
		}
		if len(missing) > 0 {
			*fail = append(*fail, FailedPredicate{
				Predicate:   "counterfeit.protocol_id_mismatch",
				EvidenceRef: "evidence_refs",
				Diagnostic: fmt.Sprintf(
					"protocol_id declares cnos.cdd.cdr.receipt.v1 but missing required CDR evidence_refs keys: %v",
					missing,
				),
			})
		}
	}
}

// sameAny compares two interface values for equality. Strings, numbers,
// bools all compare; nil == nil. For complex types we fall back to
// fmt.Sprint, which is enough for our verdict-string compare use case.
func sameAny(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if as, ok := a.(string); ok {
		if bs, ok := b.(string); ok {
			return as == bs
		}
	}
	return fmt.Sprint(a) == fmt.Sprint(b)
}

// sortedKeys returns the keys of m in lex order. Deterministic iteration
// (eng/go §2.13: no map iteration order in rendered output).
func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j-1] > keys[j]; j-- {
			keys[j-1], keys[j] = keys[j], keys[j-1]
		}
	}
	return keys
}
