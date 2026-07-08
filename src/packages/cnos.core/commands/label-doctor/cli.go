package labeldoctor

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// Run is the entry point for `cn label-doctor`. It mirrors
// cnos.issues/commands/issues-fsm's Run(ctx, args, stdin, stdout,
// stderr) signature exactly, per the γ scaffold's CLI integration axis,
// so src/go/internal/cli/cmd_label_doctor.go can be as thin as
// cmd_issues_fsm.go (Run delegates entirely to this package).
func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("label-doctor", flag.ContinueOnError)
	fs.SetOutput(stderr)
	repoFlag := fs.String("repo", "", "target repository as owner/name (default: resolved from the current git checkout's \"origin\" remote)")
	tokenFlag := fs.String("token", "", "GitHub token (default: $GITHUB_TOKEN, then $GH_TOKEN)")
	dryRun := fs.Bool("dry-run", false, "report drift only; make no GitHub API mutations. Exits nonzero if any canonical label is missing or drifted (CI-guard mode).")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: cn label-doctor [--repo owner/name] [--token TOKEN] [--dry-run]")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "Audits (and, unless --dry-run, repairs) the target repo's GitHub labels")
		fmt.Fprintln(stderr, "against src/packages/cnos.core/labels.json (schema cn.labels.v1): creates")
		fmt.Fprintln(stderr, "any missing canonical label and corrects color/description on any drifted")
		fmt.Fprintln(stderr, "one. Idempotent: a second run against unchanged live state performs zero")
		fmt.Fprintln(stderr, "mutating API calls.")
		fmt.Fprintln(stderr)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	repoRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("label-doctor: %w", err)
	}

	res, doctorErr := Doctor(ctx, Options{
		Repo:     *repoFlag,
		RepoRoot: repoRoot,
		Token:    *tokenFlag,
		DryRun:   *dryRun,
		Stdout:   stdout,
		Stderr:   stderr,
	})
	if doctorErr == nil {
		return nil
	}
	if errors.Is(doctorErr, ErrDrift) {
		n := 0
		if res != nil {
			for _, f := range res.Findings {
				if f.NeedsRepair() {
					n++
				}
			}
		}
		fmt.Fprintf(stderr, "✗ label-doctor: drift detected (%d finding(s) not canonical) — rerun without --dry-run to repair\n", n)
		return doctorErr
	}
	fmt.Fprintf(stderr, "✗ %s\n", doctorErr)
	return doctorErr
}
