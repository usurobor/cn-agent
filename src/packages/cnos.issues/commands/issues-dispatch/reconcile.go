// Package issuesdispatch implements `cn issues dispatch --issue N`: the
// cnos#640 one-verb dispatch primitive. It authorizes a design-first issue
// for dispatch by (a) flipping its status:ready label to status:todo when
// the issue is currently held at status:ready, and (b) stripping the
// legacy "Not dispatched — status:ready ... dispatch on explicit operator
// authorization" body-hold prose when present, so a human/κ never again
// has to hand-edit both surfaces separately and risk leaving them
// contradictory (the cnos#614 -> cnos#633 recurrence this cell fixes).
//
// This command is operator/human-invoked only. No dispatch wake or
// scan/reconciler process calls it on its own initiative — see
// dispatch-protocol/SKILL.md §1.2 "Human dispatch gate", which this
// primitive operationalizes rather than replaces: the status:ready ->
// status:todo transition remains the operator's authorization event, this
// command simply makes performing that event (plus the body cleanup that
// must accompany it) a single caller-visible operation instead of two
// separately-remembered hand edits.
package issuesdispatch

import "regexp"

// holdPhrasePattern matches the legacy design-first body-hold convention
// named in cnos#640 (the exact #614/#633 shape) and reconstructed in the
// γ scaffold's oracle: a (usually bold) "Not dispatched" lead-in, an em-
// or hyphen-dash, arbitrary prose in between (the current status token,
// "Design-first;", etc.), ending at the fixed closing clause "dispatch on
// explicit operator authorization." No canonical doc ever defined this
// convention (γ confirmed via repo-wide grep — see gamma-scaffold.md §2
// row 1); the pattern below is reverse-engineered from the two live
// occurrences (#614, #633) and the issue's own body quote, not derived
// from a spec. It intentionally does NOT match #618's different phrasing
// ("Filing-only — operator holds dispatch.") — #618 is explicitly out of
// scope for this cell (gamma-scaffold.md §5); a future occurrence of that
// or any other phrasing is a new pattern-family to add explicitly, not
// something this primitive should guess at.
//
// (?is) — case-insensitive, dot matches newline, so the non-greedy body
// between the lead-in and the closing clause may itself span lines (the
// live fixtures wrap the sentence across the "status:X" and "Design-
// first;" clauses on one paragraph, but nothing guarantees that never
// wraps).
var holdPhrasePattern = regexp.MustCompile(`(?is)\**Not dispatched\**\s*[—-]\s*.*?dispatch on explicit operator authorization\.`)

// trailingSeparatorPattern trims a trailing markdown horizontal rule that
// commonly precedes the hold-phrase paragraph in the live fixtures, so
// ReconcileBody does not leave a dangling "---" as the new last line of
// the body once the hold paragraph itself is removed and trailing
// whitespace has already been trimmed. Anchored at the absolute end of
// string (no (?m)) — this is applied only after trimTrailingBlank, so the
// separator, if present, is always the literal tail of the string.
var trailingSeparatorPattern = regexp.MustCompile(`\n{1,2}-{3,}$`)

// ReconcileBody strips the hold-phrase paragraph from body if present and
// returns the cleaned body plus whether a change was made. When body
// carries no hold phrase, ReconcileBody returns body unchanged and
// changed=false — the AC1 "safe no-op" oracle for an already-clean issue:
// this function must never rewrite a body that has nothing to clean, so a
// caller can use changed to decide whether a body-edit HTTP call is even
// worth making.
//
// This is a pure function (no I/O) so it is directly unit-testable
// without any GitHub double.
func ReconcileBody(body string) (cleaned string, changed bool) {
	loc := holdPhrasePattern.FindStringIndex(body)
	if loc == nil {
		return body, false
	}
	cleaned = body[:loc[0]] + body[loc[1]:]
	// Trim trailing whitespace left by the removal, then fold a
	// now-orphaned trailing "---" separator (if the hold paragraph was
	// preceded by one) and trim again, so the cleaned body never ends in
	// a dangling rule or blank lines.
	cleaned = trimTrailingBlank(cleaned)
	cleaned = trailingSeparatorPattern.ReplaceAllString(cleaned, "")
	cleaned = trimTrailingBlank(cleaned)
	return cleaned, true
}

// trimTrailingBlank removes trailing whitespace/newlines, mirroring
// strings.TrimRight(s, " \t\n\r") without importing strings solely for
// this one call site's exact semantics (kept explicit for readability
// next to the two regexps above).
func trimTrailingBlank(s string) string {
	end := len(s)
	for end > 0 {
		c := s[end-1]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			end--
			continue
		}
		break
	}
	return s[:end]
}
