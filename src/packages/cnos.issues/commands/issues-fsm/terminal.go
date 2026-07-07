package issuesfsm

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
)

// --- cnos#615: terminal-hygiene reconciliation -------------------------
//
// This file implements `cn issues fsm terminal`: a sibling reconciler to
// the #593 recovery scanner (scan.go), NOT a mode of it and NOT a new FSM
// state. scan.go's job is "does this OPEN issue's live-run/matter evidence
// match what the transition table proposes for its CURRENT status:*
// state" -- it is glue around Evaluate(table, snap) (table.go), the same
// engine `evaluate` uses.
//
// This pass has no FSM state to evaluate at all. A CLOSED issue's "next
// step" is not a transitions-table lookup -- it is a direct function of
// GitHub's own state_reason on the Issue resource, a fact this package's
// FSM engine (table.go / transitions.json) has no vocabulary for and does
// not need one for. transitions.json remains byte-unchanged by this file
// (AC4); Evaluate/Table are never called here.
//
// For every CLOSED dispatch:cell+protocol:{P} issue still carrying
// status:*/dispatch:cell/protocol:{P}: strip those labels (reusing
// ghRemoveLabel from fetch.go verbatim) and add the resolution/* label
// GitHub's own state_reason resolves to (reusing ghAddLabel verbatim, plus
// one new small primitive, ghEnsureLabelExists, to create
// resolution/not-planned the first time it is needed).

// TerminalIssueResult is one closed issue's terminal-hygiene outcome
// within one `cn issues fsm terminal` run.
type TerminalIssueResult struct {
	Issue       int      `json:"issue"`
	StateReason string   `json:"state_reason"`
	StaleLabels []string `json:"stale_labels,omitempty"`
	Resolution  string   `json:"resolution,omitempty"`
	// Reconciled reports whether labels were actually mutated (only
	// meaningful with Apply == true).
	Reconciled bool   `json:"reconciled"`
	Note       string `json:"note,omitempty"`
	Error      string `json:"error,omitempty"`
}

// TerminalReport is the full output of one terminal-sweep run.
type TerminalReport struct {
	Protocol string                `json:"protocol"`
	Results  []TerminalIssueResult `json:"results"`
}

// AnyError reports whether any swept issue recorded a per-issue error.
func (r TerminalReport) AnyError() bool {
	for _, res := range r.Results {
		if res.Error != "" {
			return true
		}
	}
	return false
}

func countTerminalErrors(r TerminalReport) int {
	n := 0
	for _, res := range r.Results {
		if res.Error != "" {
			n++
		}
	}
	return n
}

// terminalCandidate is one closed dispatch:cell+protocol:{P} issue as
// returned by ListClosedCandidates: its full current label set and
// GitHub's own close-reason classification, both read in a single list
// call (mirrors liveListActiveIssues's single-query-then-client-filter
// shape in scan.go -- no per-issue follow-up GET is needed).
type terminalCandidate struct {
	Number      int
	Labels      []string
	StateReason string
}

// resolutionForStateReason resolves GitHub's Issue.state_reason to this
// cell's pinned resolution label, per the operator's binding clarifying
// comment on cnos#615 (posted 2026-07-07T01:51:47Z; supersedes the issue
// body's original "not_planned -> resolution/wontfix-or-equivalent"
// phrasing):
//
//	"completed"   -> resolution/completed  (already exists in the repo)
//	"not_planned" -> resolution/not-planned (this cell creates it)
//
// Any other value (in practice: only reachable if GitHub ever returns an
// enum value this pass does not recognize -- "reopened" cannot occur on a
// currently-closed issue, and null/"" means GitHub simply did not record a
// reason) resolves to "" so the caller skips the issue and mutates
// nothing. This pass never guesses a resolution label for an unrecognized
// reason (cnos#368 "never blind-act on missing evidence" doctrine, the
// same discipline scan.go's "blocked" outcome already encodes).
//
// GitHub's state_reason has no "duplicate" enum value at all -- an issue
// closed "as a duplicate" via GitHub's UI reports state_reason
// "not_planned", so it already receives resolution/not-planned under the
// mapping above with no extra branch (see gamma-scaffold.md Friction note
// 2 for the full reachability argument; this is why no dedicated
// "duplicate" case or fixture exists here).
func resolutionForStateReason(stateReason string) string {
	switch stateReason {
	case "completed":
		return "resolution/completed"
	case "not_planned":
		return "resolution/not-planned"
	default:
		return ""
	}
}

// resolutionLabelColor / resolutionLabelDescription pin the
// resolution/not-planned label this cell creates to match
// resolution/completed's live values exactly, per the operator's pinned
// instruction. Verified live via
// `gh api repos/usurobor/cnos/labels/resolution%2Fcompleted` at α
// implementation time (2026-07-07): color "ededed", description null --
// see self-coherence.md §ACs AC2 for the verification transcript.
const (
	resolutionLabelColor       = "ededed"
	resolutionLabelDescription = ""
)

// TerminalOptions configures RunTerminalSweep. Every GitHub-side
// dependency is an injectable field with a live default, mirroring
// ScanOptions's injectable-dependency pattern (scan.go) -- so the
// orchestration logic (which labels get removed/added for which issue) is
// testable without a real network call. Unlike ScanOptions, this pass
// needs no local-git/CI observation and no Table -- it never evaluates an
// FSM state (AC4).
type TerminalOptions struct {
	Repo     string
	Token    string
	Protocol string
	Apply    bool

	// ListClosedCandidates returns every closed dispatch:cell+protocol:
	// {Protocol} issue's current label set and state_reason. Defaults to a
	// live GitHub REST query (liveListClosedCandidates).
	ListClosedCandidates func(ctx context.Context, repo, protocol, token string) ([]terminalCandidate, error)

	// RemoveLabel, AddLabel, and EnsureLabelExists are the label-mutation
	// primitives. Defaults to fetch.go's ghRemoveLabel/ghAddLabel (reused
	// verbatim -- this pass forks no DELETE/POST-labels logic of its own)
	// plus the one new primitive this cell adds, ghEnsureLabelExists.
	RemoveLabel       func(ctx context.Context, repo string, issue int, token, label string) error
	AddLabel          func(ctx context.Context, repo string, issue int, token, label string) error
	EnsureLabelExists func(ctx context.Context, repo, token, name, color, description string) error
}

func (o *TerminalOptions) setDefaults() {
	if o.ListClosedCandidates == nil {
		o.ListClosedCandidates = liveListClosedCandidates
	}
	if o.RemoveLabel == nil {
		o.RemoveLabel = ghRemoveLabel
	}
	if o.AddLabel == nil {
		o.AddLabel = ghAddLabel
	}
	if o.EnsureLabelExists == nil {
		o.EnsureLabelExists = ghEnsureLabelExists
	}
}

// RunTerminalSweep is the terminal-hygiene reconciler's engine (cnos#615):
// for every CLOSED dispatch:cell+protocol:{Protocol} issue still carrying
// status:*/dispatch:cell/protocol:{Protocol}, remove those labels and add
// the resolution/* label GitHub's own state_reason resolves to.
//
// Idempotence (AC3) is structural, not bolted on: once a closed issue's
// dispatch:cell and protocol:{Protocol} labels are both removed,
// ListClosedCandidates's own list-query filter (labels=
// dispatch:cell,protocol:{P}&state=closed) excludes it from every
// subsequent run -- no separate "already processed" marker/check is
// needed (mirrors #593's scanner's idempotence falling out of `todo` not
// being in activeScanStates). An OPEN dispatch:cell+protocol:{P} issue is
// never a candidate at all: the list query's own state=closed filter
// excludes it structurally, the same way scan.go's ListActiveIssues never
// sees a closed issue.
func RunTerminalSweep(ctx context.Context, opts *TerminalOptions) (TerminalReport, error) {
	opts.setDefaults()
	report := TerminalReport{Protocol: opts.Protocol}

	candidates, err := opts.ListClosedCandidates(ctx, opts.Repo, opts.Protocol, opts.Token)
	if err != nil {
		return report, fmt.Errorf("list closed candidates: %w", err)
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Number < candidates[j].Number })

	for _, c := range candidates {
		report.Results = append(report.Results, terminalOne(ctx, opts, c))
	}

	if report.AnyError() {
		return report, fmt.Errorf("cn issues fsm terminal: %d issue(s) reported an error -- see report for detail", countTerminalErrors(report))
	}
	return report, nil
}

// terminalOne reconciles one closed candidate. It never evaluates an FSM
// state and never calls Evaluate/Table -- see this file's package doc.
func terminalOne(ctx context.Context, opts *TerminalOptions, c terminalCandidate) TerminalIssueResult {
	res := TerminalIssueResult{Issue: c.Number, StateReason: c.StateReason}

	// Collect every status:* label present (label-doctrine says at most
	// one, but this pass does not assume exactly one -- it strips
	// whichever status:* labels are actually found), plus dispatch:cell
	// and protocol:{P} -- the three label classes the issue's own wording
	// names as "live selectors" to clear.
	var stale []string
	for _, l := range c.Labels {
		if strings.HasPrefix(l, "status:") {
			stale = append(stale, l)
		}
	}
	stale = append(stale, "dispatch:cell", "protocol:"+opts.Protocol)
	sort.Strings(stale)
	res.StaleLabels = stale

	resolution := resolutionForStateReason(c.StateReason)
	if resolution == "" {
		res.Note = fmt.Sprintf("state_reason %q not recognized; skipped -- never guessing a resolution label", orNoneStr(c.StateReason, "(empty)"))
		return res
	}
	res.Resolution = resolution

	if !opts.Apply {
		res.Note = fmt.Sprintf("would remove %s and add %s -- rerun with --apply", strings.Join(stale, ", "), resolution)
		return res
	}

	for _, l := range stale {
		if err := opts.RemoveLabel(ctx, opts.Repo, c.Number, opts.Token, l); err != nil {
			res.Error = fmt.Sprintf("remove %s: %v", l, err)
			return res
		}
	}
	if err := opts.EnsureLabelExists(ctx, opts.Repo, opts.Token, resolution, resolutionLabelColor, resolutionLabelDescription); err != nil {
		res.Error = fmt.Sprintf("ensure label %s exists: %v", resolution, err)
		return res
	}
	if err := opts.AddLabel(ctx, opts.Repo, c.Number, opts.Token, resolution); err != nil {
		res.Error = fmt.Sprintf("add %s: %v", resolution, err)
		return res
	}
	res.Reconciled = true
	res.Note = fmt.Sprintf("reconciled: removed %s, added %s", strings.Join(stale, ", "), resolution)
	return res
}

// --- live dependency defaults -----------------------------------------

// liveListClosedCandidates queries the GitHub REST API for CLOSED issues
// carrying both dispatch:cell and protocol:{protocol}, returning each
// candidate's full label set and state_reason straight from the single
// list response -- no per-issue follow-up GET is needed. This single
// query IS the AC1 recognizer: it can only ever return issues that are
// closed AND carry both dispatch:cell and protocol:{protocol}; a closed
// issue missing either label, or an open issue regardless of its labels,
// is excluded by the query's own labels=/state= filter and never reaches
// this function's caller at all.
func liveListClosedCandidates(ctx context.Context, repo, protocol, token string) ([]terminalCandidate, error) {
	if repo == "" {
		return nil, fmt.Errorf("repo is required")
	}
	if protocol == "" {
		return nil, fmt.Errorf("protocol is required")
	}
	var issues []struct {
		Number int `json:"number"`
		Labels []struct {
			Name string `json:"name"`
		} `json:"labels"`
		StateReason *string `json:"state_reason"`
	}
	listURL := fmt.Sprintf("%s/repos/%s/issues?labels=dispatch:cell,protocol:%s&state=closed&per_page=100", githubAPIBase, repo, protocol)
	if err := ghGetJSON(ctx, listURL, token, &issues); err != nil {
		return nil, err
	}
	out := make([]terminalCandidate, 0, len(issues))
	for _, is := range issues {
		c := terminalCandidate{Number: is.Number}
		for _, l := range is.Labels {
			c.Labels = append(c.Labels, l.Name)
		}
		if is.StateReason != nil {
			c.StateReason = *is.StateReason
		}
		out = append(out, c)
	}
	return out, nil
}

// renderTerminal writes a stable, deterministic report to w -- one line
// per swept issue plus a summary, mirroring renderScan's plain-text
// convention (scan.go).
func renderTerminal(w io.Writer, r TerminalReport) {
	fmt.Fprintf(w, "cn issues fsm terminal — protocol %s\n\n", orNone(r.Protocol))
	if len(r.Results) == 0 {
		fmt.Fprintln(w, "No closed dispatch:cell issues carrying stale status:*/dispatch:cell/protocol:{P} labels found for this protocol.")
		return
	}
	for _, res := range r.Results {
		fmt.Fprintf(w, "#%d  state_reason=%s\n", res.Issue, orNone(res.StateReason))
		if res.Error != "" {
			fmt.Fprintf(w, "  error: %s\n", res.Error)
			continue
		}
		fmt.Fprintf(w, "  %s\n", res.Note)
		if res.Reconciled {
			fmt.Fprintln(w, "  reconciled: true")
		}
	}
}
