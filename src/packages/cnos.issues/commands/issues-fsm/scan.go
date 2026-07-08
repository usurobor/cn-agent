package issuesfsm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// activeScanStates is the set of CDS states the recovery scanner reconciles:
// only states where a cell can be silently stranded mid-cycle by a dead
// workflow run. status:todo/ready await a claim or operator "go"; status:changes
// awaits operator repair; status:blocked awaits operator resolution -- none of
// these "strand" the way in-progress/review do when the run driving them
// dies, so the scanner does not touch them (cnos#593 scope).
var activeScanStates = map[string]bool{
	"in-progress": true,
	"review":      true,
}

// ScanIssueResult is one issue's reconciliation outcome within a scan run.
type ScanIssueResult struct {
	Issue        int    `json:"issue"`
	CurrentState string `json:"current_state"`
	Outcome      string `json:"outcome"`
	Action       string `json:"action"`
	TargetState  string `json:"target_state,omitempty"`
	// Reconciled reports whether a status-label transition was applied
	// (only meaningful with Apply == true).
	Reconciled bool `json:"reconciled"`
	// Finalized reports whether `cn cell finalize` was invoked for this
	// issue this run (cases 3/4: matter checkpointing).
	Finalized bool `json:"finalized"`
	// Commented reports whether an issue comment was posted recording the
	// reconciliation (never posted for a no-op or an already-checkpointed
	// dead-with-matter state -- see RunScan's idempotence discipline).
	Commented bool   `json:"commented,omitempty"`
	Note      string `json:"note,omitempty"`
	Error     string `json:"error,omitempty"`
}

// ScanReport is the full output of one scan run.
type ScanReport struct {
	Protocol string            `json:"protocol"`
	Results  []ScanIssueResult `json:"results"`
}

// AnyError reports whether any scanned issue recorded a per-issue error.
func (r ScanReport) AnyError() bool {
	for _, res := range r.Results {
		if res.Error != "" {
			return true
		}
	}
	return false
}

// ScanOptions configures RunScan. Every dependency the live path needs is an
// injectable field with a live default -- mirroring cell.Finalizer's
// RunGH/RunGHJSON injection pattern (src/go/internal/cell/cell.go) -- so the
// orchestration logic (which issue gets which action) is testable without a
// real subprocess or network call.
type ScanOptions struct {
	Repo     string
	Token    string
	Protocol string
	Apply    bool
	Table    *Table

	// ListActiveIssues returns the open dispatch:cell+protocol:{Protocol}
	// issue numbers currently at status:in-progress or status:review.
	// Defaults to a live GitHub REST query.
	ListActiveIssues func(ctx context.Context, repo, protocol, token string) ([]int, error)

	// AssembleFacts observes one issue's FactSnapshot. Defaults to
	// assembleLive (fetch.go) -- the exact same live observation path
	// `cn issues fsm evaluate` uses, so a scanned issue's facts are read
	// identically to how a single-issue evaluate call would read them.
	AssembleFacts func(ctx context.Context, repo string, issue int, token string) (FactSnapshot, error)

	// RunFinalize invokes the mechanical checkpoint+PR-open finalizer
	// (cnos#591) for one issue and returns its captured stdout (used only
	// to detect whether a NEW draft PR was opened this call, for the
	// idempotent-comment decision below) and any error.
	//
	// Defaults to self-re-exec: the currently running `cn` binary invoked
	// as `cn cell finalize --issue N`. The finalizer lives in a different
	// Go module (src/go/internal/cell, part of the src/go module) that this
	// package (a separate module under src/packages/cnos.issues/) cannot
	// import directly -- Go's internal-package visibility rule scopes
	// `.../src/go/internal/cell` to importers rooted under `.../src/go/`.
	// Subprocess invocation of the same compiled binary reuses cnos#591's
	// finalizer verbatim instead of duplicating its checkpoint/PR logic.
	RunFinalize func(ctx context.Context, issue int) (stdout string, err error)

	// PostComment posts one issue comment. Defaults to a live GitHub REST
	// call. Injectable so tests can assert exactly when (and how often) a
	// comment is posted without a real network call.
	PostComment func(ctx context.Context, repo string, issue int, body string) error
}

func (o *ScanOptions) setDefaults() {
	if o.ListActiveIssues == nil {
		o.ListActiveIssues = liveListActiveIssues
	}
	if o.AssembleFacts == nil {
		o.AssembleFacts = assembleLive
	}
	if o.RunFinalize == nil {
		o.RunFinalize = liveRunFinalize
	}
	if o.PostComment == nil {
		token := o.Token
		o.PostComment = func(ctx context.Context, repo string, issue int, body string) error {
			return livePostComment(ctx, repo, issue, token, body)
		}
	}
}

// RunScan is the recovery scanner's engine (cnos#593 Sub C): for every open
// dispatch:cell+protocol:{Protocol} issue at status:in-progress or
// status:review, observe its facts, evaluate them against Table, and
// reconcile per the matched rule -- exactly the per-state behavior
// `cn issues fsm evaluate` already implements (table.go / transitions.json),
// scanned across every active issue instead of the one issue an operator
// names.
//
// RunScan performs no cognition of its own: every action it takes traces to
// a rule the transition table already proposes for `evaluate`. The two
// things this function adds beyond a single `evaluate` call are (a)
// iteration across every active issue and (b) invoking the mechanical
// finalizer (cnos#591) when a "propose_delta_recovery" decision is matched,
// so a dead run's matter is checkpointed into a draft PR without a human
// having to notice the strand first.
//
// cnos#630: once matter is checkpointed (a PR exists) and no
// REVIEW-REQUEST.yml has followed, the table proposes
// "propose_status_todo_with_matter" (in-progress -> todo) instead of
// re-matching propose_delta_recovery forever -- this is the mechanical exit
// from the "partial-matter in-progress wedge": the cell re-enters the
// status:todo claim queue with its branch/PR PRESERVED (never deleted, the
// cnos#368 protection restated for the reconciler's own action), so a fresh
// claim resumes from the existing matter instead of leaving the cell
// invisible to the claim queue forever. That decision flows through the
// same generic "proposed with a TargetState" reconciliation branch below
// that Cases 2/5 already use -- no separate code path.
//
// No status-label mutation happens anywhere in this function except through
// applyStatusLabel -- the same primitive `evaluate --apply` uses. RunScan
// never calls gh issue edit directly and never invents a transition the
// table did not already propose (AC8).
func RunScan(ctx context.Context, opts *ScanOptions) (ScanReport, error) {
	opts.setDefaults()
	report := ScanReport{Protocol: opts.Protocol}

	issues, err := opts.ListActiveIssues(ctx, opts.Repo, opts.Protocol, opts.Token)
	if err != nil {
		return report, fmt.Errorf("list active issues: %w", err)
	}
	sort.Ints(issues)

	for _, issue := range issues {
		report.Results = append(report.Results, scanOne(ctx, opts, issue))
	}

	if report.AnyError() {
		return report, fmt.Errorf("cn issues fsm scan: %d issue(s) reported an error -- see report for detail", countScanErrors(report))
	}
	return report, nil
}

func countScanErrors(r ScanReport) int {
	n := 0
	for _, res := range r.Results {
		if res.Error != "" {
			n++
		}
	}
	return n
}

func scanOne(ctx context.Context, opts *ScanOptions, issue int) ScanIssueResult {
	res := ScanIssueResult{Issue: issue}

	snap, err := opts.AssembleFacts(ctx, opts.Repo, issue, opts.Token)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	snap.Issue = issue

	dec, err := Evaluate(opts.Table, snap)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	res.CurrentState = dec.CurrentState
	res.Outcome = dec.Outcome
	res.Action = dec.Action
	res.TargetState = dec.TargetState

	if !activeScanStates[dec.CurrentState] {
		// Defensive re-check: ListActiveIssues already filtered to
		// in-progress/review, but a concurrent transition between listing
		// and observing this issue is possible (e.g. a dispatch claim
		// firing raced this scan). Re-verify rather than trust the stale
		// listing -- never act on a state this scanner does not own.
		res.Note = "current state is not an active-scan state (in-progress/review); skipped"
		return res
	}

	switch {
	case dec.Outcome == "valid":
		// Case 1 (live run active) and any other steady "nothing to do"
		// state: no-op, no comment.
		res.Note = dec.Reason

	case dec.Outcome == "blocked":
		// Case 6 (stale status:review with missing proof) and any other
		// blocked-in-place decision: report, never mutate, never comment.
		// Commenting on every scan tick against an unchanged blocked state
		// would violate AC9 idempotence (duplicate comments accumulating
		// every ~15 minutes); the blocked reason is operator-visible via
		// this report / the workflow log instead.
		res.Note = dec.Reason

	case dec.Outcome == "proposed" && dec.Action == "propose_delta_recovery":
		// Case 3: matter (commits) exists but is not yet checkpointed into
		// a PR, with no active run. Checkpoint via the cnos#591 finalizer.
		// delta-recovery has no TargetState because the correct next step
		// is "matter is now checkpointed, awaiting REVIEW-REQUEST.yml or
		// an operator/δ decision," never a blind status move (the
		// cnos#368 protection this same rule already encodes for
		// `evaluate`).
		//
		// cnos#630: once a PR exists, transitions.json's new
		// "propose_status_todo_with_matter" rule matches BEFORE this one
		// (see transitions.json's in-progress rule ordering), so this case
		// is reached only pre-checkpoint (branch commits, no PR yet) --
		// the "PR already existed" else-branch below is now a defensive
		// fallback for an observe-vs-apply race (facts observed with
		// pr_exists==false, but a PR appeared between observation and this
		// finalize call), not the steady-state path it used to be.
		stdout, ferr := opts.RunFinalize(ctx, issue)
		if ferr != nil {
			res.Error = fmt.Sprintf("finalize: %v", ferr)
			break
		}
		res.Finalized = true
		if strings.Contains(stdout, "opened draft PR for") {
			res.Note = "dead run with matter: no PR existed; checkpointed via cn cell finalize and opened a draft PR. Status remains in-progress -- awaiting REVIEW-REQUEST.yml or operator/δ follow-up."
			body := fmt.Sprintf(
				"cds recovery scan: dead run detected on #%d (no active workflow run; matter present) with no pull request. "+
					"Checkpointed via `cn cell finalize` and opened a draft PR for `cycle/%d`. "+
					"Status remains `status:in-progress` -- the branch/PR is preserved so the next dispatched round resumes from it rather than starting over.",
				issue, issue)
			if opts.Apply {
				if cerr := opts.PostComment(ctx, opts.Repo, issue, body); cerr != nil {
					res.Error = fmt.Sprintf("post comment: %v", cerr)
					break
				}
				res.Commented = true
			}
		} else {
			// A PR (or the branch/commits it depends on) already existed
			// at finalize time despite this tick's facts observing
			// pr_exists==false (observe-vs-apply race -- see the case
			// comment above); cn cell finalize's own idempotent no-op path
			// ran. Nothing new happened this tick -- do not re-comment.
			// The next scan tick will re-observe pr_exists==true and
			// route through transitions.json's propose_status_todo_with_matter
			// rule (the generic "proposed with a TargetState" case below),
			// which is the cnos#630 mechanical exit from this state.
			res.Note = "dead run with matter: PR (or commits) already checkpointed; cn cell finalize confirmed idempotent no-op. Status remains in-progress -- awaiting the next scan tick's propose_status_todo_with_matter reconciliation."
		}

	case dec.Outcome == "proposed" && dec.TargetState != "":
		// Cases 2 and 5, plus cnos#630's propose_status_todo_with_matter:
		// dead-in-progress-no-matter -> todo, PR+REVIEW-REQUEST.yml
		// present -> review, or dead-in-progress-with-checkpointed-PR ->
		// todo (matter preserved, cnos#630 AC1/AC2). Apply through the
		// SAME applyStatusLabel primitive `evaluate --apply` uses -- the
		// FSM (this package) remains the only label-writer; scan never
		// bypasses it with a direct gh issue edit. The comment posted
		// below (built from dec.Reason, the transitions.json rule's own
		// prose) is the cnos#630 AC4 audit note: every mechanical
		// reversion this branch performs is documented on the issue so a
		// subsequent claim/scan pass does not misread it as an
		// unexplained orphan.
		if !opts.Apply {
			res.Note = fmt.Sprintf("would reconcile %s (%s) -- rerun with --apply", dec.EnabledTransition, dec.Reason)
			break
		}
		if err := applyStatusLabel(ctx, opts.Repo, issue, opts.Token, dec.CurrentState, dec.TargetState); err != nil {
			res.Error = fmt.Sprintf("apply %s: %v", dec.EnabledTransition, err)
			break
		}
		res.Reconciled = true
		res.Note = fmt.Sprintf("reconciled %s (%s)", dec.EnabledTransition, dec.Reason)
		body := fmt.Sprintf(
			"cds recovery scan: reconciled #%d -- %s.\n\n%s",
			issue, dec.EnabledTransition, dec.Reason)
		if cerr := opts.PostComment(ctx, opts.Repo, issue, body); cerr != nil {
			res.Error = fmt.Sprintf("post comment: %v", cerr)
			break
		}
		res.Commented = true

	default:
		res.Note = dec.Reason
	}

	return res
}

// --- live dependency defaults -----------------------------------------

// liveListActiveIssues queries the GitHub REST API for open issues carrying
// both dispatch:cell and protocol:{protocol}, then filters client-side to
// those whose status:* label is in activeScanStates (GitHub's issues-list
// `labels` query param is an AND across all given label names; it cannot
// itself express an OR across two different status labels).
func liveListActiveIssues(ctx context.Context, repo, protocol, token string) ([]int, error) {
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
	}
	listURL := fmt.Sprintf("%s/repos/%s/issues?labels=dispatch:cell,protocol:%s&state=open&per_page=100", githubAPIBase, repo, protocol)
	if err := ghGetJSON(ctx, listURL, token, &issues); err != nil {
		return nil, err
	}
	var out []int
	for _, is := range issues {
		state := ""
		for _, l := range is.Labels {
			if strings.HasPrefix(l.Name, "status:") {
				state = strings.TrimPrefix(l.Name, "status:")
				break
			}
		}
		if activeScanStates[state] {
			out = append(out, is.Number)
		}
	}
	return out, nil
}

// liveRunFinalize self-re-execs the currently running `cn` binary as
// `cn cell finalize --issue N`, capturing combined stdout+stderr so the
// caller can distinguish "opened a new draft PR" from "idempotent no-op"
// (see cell.Finalizer's stdout markers in src/go/internal/cell/cell.go).
func liveRunFinalize(ctx context.Context, issue int) (string, error) {
	self, err := os.Executable()
	if err != nil || self == "" {
		self = os.Args[0]
	}
	cmd := exec.CommandContext(ctx, self, "cell", "finalize", "--issue", strconv.Itoa(issue))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	runErr := cmd.Run()
	return buf.String(), runErr
}

// livePostComment posts one issue comment via the GitHub REST API, reusing
// the same authenticated-request idiom as ghRequest/ghAddLabel in fetch.go.
func livePostComment(ctx context.Context, repo string, issue int, token, body string) error {
	payload, err := json.Marshal(struct {
		Body string `json:"body"`
	}{Body: body})
	if err != nil {
		return err
	}
	commentURL := fmt.Sprintf("%s/repos/%s/issues/%d/comments", githubAPIBase, repo, issue)
	resp, err := ghRequest(ctx, http.MethodPost, commentURL, token, payload)
	if err != nil {
		return fmt.Errorf("github api post comment: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github api post comment: HTTP %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// renderScan writes a stable, deterministic report to w -- one line per
// scanned issue plus a summary, mirroring Decision.Render's plain-text
// convention.
func renderScan(w io.Writer, r ScanReport) {
	fmt.Fprintf(w, "cn issues fsm scan — protocol %s\n\n", orNone(r.Protocol))
	if len(r.Results) == 0 {
		fmt.Fprintln(w, "No active (status:in-progress / status:review) dispatch:cell issues found for this protocol.")
		return
	}
	for _, res := range r.Results {
		fmt.Fprintf(w, "#%d  state=%s  outcome=%s  action=%s\n", res.Issue, orNone(res.CurrentState), orNone(res.Outcome), orNoneStr(res.Action, "none"))
		if res.Error != "" {
			fmt.Fprintf(w, "  error: %s\n", res.Error)
			continue
		}
		fmt.Fprintf(w, "  %s\n", res.Note)
		if res.Reconciled {
			fmt.Fprintln(w, "  reconciled: true")
		}
		if res.Finalized {
			fmt.Fprintln(w, "  finalized: true")
		}
		if res.Commented {
			fmt.Fprintln(w, "  commented: true")
		}
	}
}
