package issuesdispatch

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
)

// Run is the entry point for `cn issues dispatch`. There is exactly one
// sub-verb ("dispatch" itself has no further noun — the issue is named by
// --issue), mirroring `cn issues fsm evaluate`'s single-issue shape rather
// than `scan`/`terminal`'s protocol-wide sweep shape: this primitive is a
// human/κ-invoked authorization action for one named issue, never a
// reconciler that walks a queue on its own initiative (γ scaffold §3
// "operator/human-invoked only" pin, AC4).
func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("issues dispatch", flag.ContinueOnError)
	fs.SetOutput(stderr)
	issue := fs.Int("issue", 0, "issue number to dispatch (required)")
	repo := fs.String("repo", "", "target repository as owner/name (default: $GITHUB_REPOSITORY)")
	token := fs.String("token", "", "GitHub token (default: $GITHUB_TOKEN, then $GH_TOKEN)")
	apply := fs.Bool("apply", false, "apply the dispatch: flips status:ready -> status:todo when the issue is currently held at status:ready, and/or strips the legacy body-hold phrase when present. Off by default -- without --apply this command is read-only (reports what it would do).")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: cn issues dispatch --issue N [--apply] [--repo owner/name]")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "Authorizes issue #N for dispatch in one operation: flips status:ready ->")
		fmt.Fprintln(stderr, "status:todo when currently held at status:ready, and strips the legacy")
		fmt.Fprintln(stderr, "\"Not dispatched -- status:ready ... dispatch on explicit operator")
		fmt.Fprintln(stderr, "authorization\" body-hold phrase when present -- so a human/kappa never")
		fmt.Fprintln(stderr, "again has to hand-edit both the label and the body separately and risk")
		fmt.Fprintln(stderr, "leaving them contradictory (cnos#640, recurrence of cnos#614/#633).")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "Operator/human-invoked only: no dispatch wake or scan/reconciler process")
		fmt.Fprintln(stderr, "calls this command on its own initiative.")
		fmt.Fprintln(stderr)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *issue == 0 {
		return fmt.Errorf("cn issues dispatch: --issue is required")
	}

	r := *repo
	if r == "" {
		r = os.Getenv("GITHUB_REPOSITORY")
	}
	tk := *token
	if tk == "" {
		tk = os.Getenv("GITHUB_TOKEN")
	}
	if tk == "" {
		tk = os.Getenv("GH_TOKEN")
	}
	if r == "" {
		return fmt.Errorf("cn issues dispatch: --repo (or $GITHUB_REPOSITORY) is required")
	}

	opts := &Options{
		Repo:  r,
		Token: tk,
		Issue: *issue,
		Apply: *apply,
	}
	res, err := Dispatch(ctx, opts)
	Render(stdout, res)
	return err
}
