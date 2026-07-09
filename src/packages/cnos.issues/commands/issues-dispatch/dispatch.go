package issuesdispatch

import (
	"context"
	"fmt"
	"io"
	"strings"
)

// Options configures Dispatch. Every GitHub-side dependency is an
// injectable field with a live default, mirroring issues-fsm's
// TerminalOptions/ScanOptions injectable-dependency pattern (terminal.go /
// scan.go) — so the orchestration logic (which of the body edit / label
// flip actually happen, and in what order) is testable without a real
// network call.
type Options struct {
	Repo  string
	Token string
	Issue int

	// Apply gates mutation, mirroring `cn issues fsm`'s evaluate/scan/
	// terminal sub-verbs: without --apply this command is read-only and
	// reports what it would do; with --apply it performs the body edit
	// and/or label flip it determined were needed.
	Apply bool

	// GetIssue fetches the target issue's current body + labels. Defaults
	// to a live GitHub REST GET (liveGetIssue).
	GetIssue func(ctx context.Context, repo string, issue int, token string) (issueState, error)

	// EditBody rewrites the issue body. Defaults to the new
	// ghEditIssueBody primitive (fetch.go) — the one genuinely new HTTP
	// call this cell adds.
	EditBody func(ctx context.Context, repo string, issue int, token, body string) error

	// RemoveLabel / AddLabel are the label-mutation primitives. Default to
	// fetch.go's ghRemoveLabel/ghAddLabel (reused-in-shape, not imported —
	// see fetch.go's doc comments for why this package cannot literally
	// import issues-fsm's unexported functions across the Go-module
	// boundary and instead reproduces their exact idiom).
	RemoveLabel func(ctx context.Context, repo string, issue int, token, label string) error
	AddLabel    func(ctx context.Context, repo string, issue int, token, label string) error
}

func (o *Options) setDefaults() {
	if o.GetIssue == nil {
		o.GetIssue = liveGetIssue
	}
	if o.EditBody == nil {
		o.EditBody = ghEditIssueBody
	}
	if o.RemoveLabel == nil {
		o.RemoveLabel = ghRemoveLabel
	}
	if o.AddLabel == nil {
		o.AddLabel = ghAddLabel
	}
}

// Result is Dispatch's caller-visible outcome for one issue.
type Result struct {
	Issue          int    `json:"issue"`
	PreviousStatus string `json:"previous_status,omitempty"`
	NewStatus      string `json:"new_status,omitempty"`
	BodyChanged    bool   `json:"body_changed"`
	LabelFlipped   bool   `json:"label_flipped"`
	Applied        bool   `json:"applied"`
	NoOp           bool   `json:"no_op"`
	Note           string `json:"note,omitempty"`
}

// statusLabelPrefix is the label-doctrine status namespace (label-doctrine
// /SKILL.md §1.1). currentStatusLabel strips it to get the bare state
// token Dispatch reasons over ("ready", "todo", ...).
const statusLabelPrefix = "status:"

// currentStatusLabel returns the bare state token of the first status:*
// label found in labels, or "" if none is present. label-doctrine pins
// "exactly one status:* label per open issue" as an invariant; this
// function does not itself enforce that invariant (a drifted issue with
// zero or multiple status:* labels is not this primitive's job to
// repair) — it simply reads the first one found, which is enough for
// Dispatch's own gating logic below.
func currentStatusLabel(labels []string) string {
	for _, l := range labels {
		if strings.HasPrefix(l, statusLabelPrefix) {
			return strings.TrimPrefix(l, statusLabelPrefix)
		}
	}
	return ""
}

// Dispatch is the cnos#640 one-verb primitive's engine: authorize issue
// #{Issue} for dispatch by flipping status:ready -> status:todo when
// currently held at status:ready, and by stripping the legacy body-hold
// phrase whenever ReconcileBody finds one — regardless of the issue's
// current status, since a stray hold sentence is never correct once
// labels are the sole source of truth for dispatch readiness (AC2). The
// two actions are computed together and, when both apply, executed body-
// first: an EditBody failure aborts before any label mutation is
// attempted, so this primitive can never itself produce the exact
// contradiction it exists to fix (status:todo label + hold phrase still
// present in the body). A label-flip failure after a successful body edit
// is safe by the same construction — it leaves the body already clean and
// the label at status:ready (not yet authorized), which is a re-runnable,
// non-contradictory state, not the #614/#633 shape.
//
// Idempotency (AC1): running Dispatch twice against the same issue
// produces the same end state both times. The first run (starting at
// status:ready with a hold-phrase body, the common fresh-dispatch case)
// performs both actions; the second run observes status:todo with an
// already-clean body and takes neither action (NoOp==true, zero mutation
// calls) — see dispatch_test.go's idempotency fixture.
func Dispatch(ctx context.Context, opts *Options) (Result, error) {
	opts.setDefaults()
	res := Result{Issue: opts.Issue}

	state, err := opts.GetIssue(ctx, opts.Repo, opts.Issue, opts.Token)
	if err != nil {
		return res, fmt.Errorf("cn issues dispatch: fetch issue #%d: %w", opts.Issue, err)
	}

	status := currentStatusLabel(state.Labels)
	res.PreviousStatus = status
	res.NewStatus = status

	cleanedBody, bodyChanged := ReconcileBody(state.Body)
	res.BodyChanged = bodyChanged

	needsLabelFlip := status == "ready"
	res.LabelFlipped = needsLabelFlip

	if !bodyChanged && !needsLabelFlip {
		res.NoOp = true
		res.Note = "already clean: no hold phrase in body, no status:ready label to flip"
		return res, nil
	}

	if !opts.Apply {
		res.Note = dryRunNote(bodyChanged, needsLabelFlip)
		return res, nil
	}

	// Body first: never flip the label into a state (status:todo) that
	// could still coexist with an un-cleaned hold phrase.
	if bodyChanged {
		if err := opts.EditBody(ctx, opts.Repo, opts.Issue, opts.Token, cleanedBody); err != nil {
			return res, fmt.Errorf("cn issues dispatch: edit body for #%d: %w (label not touched)", opts.Issue, err)
		}
	}

	if needsLabelFlip {
		if err := opts.RemoveLabel(ctx, opts.Repo, opts.Issue, opts.Token, statusLabelPrefix+"ready"); err != nil {
			return res, fmt.Errorf("cn issues dispatch: remove status:ready from #%d: %w", opts.Issue, err)
		}
		if err := opts.AddLabel(ctx, opts.Repo, opts.Issue, opts.Token, statusLabelPrefix+"todo"); err != nil {
			return res, fmt.Errorf("cn issues dispatch: add status:todo to #%d: %w", opts.Issue, err)
		}
		res.NewStatus = "todo"
	}

	res.Applied = true
	res.Note = appliedNote(bodyChanged, needsLabelFlip)
	return res, nil
}

func dryRunNote(bodyChanged, labelFlipped bool) string {
	switch {
	case bodyChanged && labelFlipped:
		return "would clean hold-phrase body AND flip status:ready -> status:todo -- rerun with --apply"
	case bodyChanged:
		return "would clean hold-phrase body (label unchanged) -- rerun with --apply"
	case labelFlipped:
		return "would flip status:ready -> status:todo (body already clean) -- rerun with --apply"
	default:
		return "nothing to do"
	}
}

func appliedNote(bodyChanged, labelFlipped bool) string {
	switch {
	case bodyChanged && labelFlipped:
		return "cleaned hold-phrase body and flipped status:ready -> status:todo"
	case bodyChanged:
		return "cleaned hold-phrase body (label was not status:ready; unchanged)"
	case labelFlipped:
		return "flipped status:ready -> status:todo (body was already clean)"
	default:
		return "nothing to do"
	}
}

// Render writes a short, deterministic summary of res to w, mirroring the
// issues-fsm sub-verbs' plain-text render idiom (renderScan / renderTerminal).
func Render(w io.Writer, res Result) {
	fmt.Fprintf(w, "cn issues dispatch — issue #%d\n", res.Issue)
	fmt.Fprintf(w, "  previous status: %s\n", orNone(res.PreviousStatus))
	fmt.Fprintf(w, "  new status:      %s\n", orNone(res.NewStatus))
	fmt.Fprintf(w, "  body changed:    %v\n", res.BodyChanged)
	fmt.Fprintf(w, "  label flipped:   %v\n", res.LabelFlipped)
	fmt.Fprintf(w, "  applied:         %v\n", res.Applied)
	fmt.Fprintf(w, "  %s\n", res.Note)
}

func orNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}
