package cddverify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExitCode is V's process exit code. Values:
//
//	0 — verdict.Result == PASS
//	1 — verdict.Result == FAIL
//	2 — V itself errored (cue missing, receipt not found, etc.)
type ExitCode int

const (
	ExitPASS     ExitCode = 0
	ExitFAIL     ExitCode = 1
	ExitVError  ExitCode = 2 //nolint:revive  // V-error per RECEIPT-VALIDATION.md
)

// Args describes the parsed cdd-verify command-line surface.
//
// Combines V receipt-validation mode (--receipt) with the legacy ledger
// modes (--unreleased, --all, --version, --pr, --cycle, --triadic).
type Args struct {
	// V receipt-mode flags.
	Receipt        string
	Contract       string
	JSON           bool
	StructuralOnly bool

	// Ledger-mode flags.
	Version       string
	PR            string
	Cycle         string
	Triadic       bool
	All           bool
	Unreleased    bool
	Bundle        string
	RepoRootOver  string
	Exceptions    string

	// Help requested.
	Help bool
}

// ParseArgs parses a cdd-verify argv slice (NOT including the command
// name). Returns the parsed Args and any structural error (unknown flag,
// missing value, etc.).
//
// Mirrors the bash predecessor's flag set so backward compat (AC7) holds.
func ParseArgs(argv []string) (Args, error) {
	var a Args
	a.Bundle = "docs/gamma/cdd" // default per bash predecessor
	i := 0
	for i < len(argv) {
		arg := argv[i]
		switch arg {
		case "-h", "--help":
			a.Help = true
			i++
		case "--receipt":
			v, err := takeValue(argv, &i, "--receipt")
			if err != nil {
				return a, err
			}
			a.Receipt = v
		case "--contract":
			v, err := takeValue(argv, &i, "--contract")
			if err != nil {
				return a, err
			}
			a.Contract = v
		case "--json":
			a.JSON = true
			i++
		case "--structural-only":
			a.StructuralOnly = true
			i++
		case "--version":
			v, err := takeValue(argv, &i, "--version")
			if err != nil {
				return a, err
			}
			a.Version = v
		case "--pr":
			v, err := takeValue(argv, &i, "--pr")
			if err != nil {
				return a, err
			}
			a.PR = v
		case "--cycle":
			v, err := takeValue(argv, &i, "--cycle")
			if err != nil {
				return a, err
			}
			a.Cycle = v
		case "--triadic":
			a.Triadic = true
			i++
		case "--all":
			a.All = true
			i++
		case "--unreleased":
			a.Unreleased = true
			i++
		case "--bundle":
			v, err := takeValue(argv, &i, "--bundle")
			if err != nil {
				return a, err
			}
			a.Bundle = v
		case "--repo-root":
			v, err := takeValue(argv, &i, "--repo-root")
			if err != nil {
				return a, err
			}
			a.RepoRootOver = v
		case "--exceptions":
			v, err := takeValue(argv, &i, "--exceptions")
			if err != nil {
				return a, err
			}
			a.Exceptions = v
		default:
			return a, fmt.Errorf("unknown option: %s", arg)
		}
	}
	return a, nil
}

func takeValue(argv []string, i *int, flag string) (string, error) {
	if *i+1 >= len(argv) {
		return "", fmt.Errorf("%s requires a value", flag)
	}
	v := argv[*i+1]
	*i += 2
	return v, nil
}

// Usage writes the cdd-verify usage banner to w. Mirrors the bash
// predecessor's usage block for backward compat (operator muscle memory).
func Usage(w io.Writer) {
	const usage = `Usage:
  cn cdd-verify --version <ver>              release-scoped cycle (canonical contract)
  cn cdd-verify --pr <number>                PR-scoped cycle (no release)
  cn cdd-verify --version <ver> --triadic    include legacy .cdd/ close-out checks
  cn cdd-verify --version <ver> --cycle <N>  check current cycle-scoped artifacts
  cn cdd-verify --all                        repository-wide ledger check (unreleased + released)
  cn cdd-verify --unreleased                 check only .cdd/unreleased/ directories
  cn cdd-verify --receipt <path>             V verdict (Contract × Receipt → ValidationVerdict)
  cn cdd-verify --receipt <path> --json      emit JSON verdict (machine-consumable by δ)
  cn cdd-verify --receipt <path> --contract <p>  V with explicit contract path

Options:
  --bundle <path>        bundle-relative PRA dir (default: docs/gamma/cdd)
  --repo-root <path>     override repo root (testing)
  --exceptions <path>    legacy exceptions file (YAML format)
  --receipt <path>       dispatch into V validator (Phase 3 of #366)
  --contract <path>      V's contract ref (optional; warns if missing)
  --json                 emit V verdict as JSON instead of prose
  --structural-only      skip filesystem-dereference rules (C1/C2/C4); useful
                         for shape-demonstrating fixtures with placeholder paths
`
	fmt.Fprint(w, usage)
}

// Run is the cdd-verify command entrypoint.
//
// Parses argv, then dispatches to either the V receipt-validation path
// (--receipt) or the legacy ledger path (--version/--pr/--all/--unreleased).
//
// stdout/stderr are injected (eng/go §2.6: side-effect boundaries; no
// fmt.Println onto package globals). Returns the V exit code.
func Run(ctx context.Context, argv []string, stdout, stderr io.Writer) ExitCode {
	args, err := ParseArgs(argv)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		Usage(stderr)
		return ExitVError
	}
	if args.Help {
		Usage(stdout)
		return ExitPASS
	}

	// --- V receipt mode dispatch -------------------------------------------
	if args.Receipt != "" {
		// Disallow combining --receipt with ledger modes (exit-code semantics
		// differ; mixing muddles δ-consumer parsing).
		if args.Version != "" || args.PR != "" || args.Cycle != "" ||
			args.Triadic || args.All || args.Unreleased {
			fmt.Fprintln(stderr, "Error: --receipt cannot be combined with --version/--pr/--cycle/--all/--unreleased/--triadic")
			Usage(stderr)
			return ExitVError
		}
		return runReceiptMode(ctx, args, stdout, stderr)
	}

	// Sanity: --contract, --json, or --structural-only without --receipt.
	if args.Contract != "" || args.JSON || args.StructuralOnly {
		fmt.Fprintln(stderr, "Error: --contract, --json, and --structural-only require --receipt")
		Usage(stderr)
		return ExitVError
	}

	// --- Ledger mode dispatch ----------------------------------------------
	return runLedgerMode(ctx, args, stdout, stderr)
}

// runReceiptMode handles the V validator path (--receipt).
func runReceiptMode(ctx context.Context, args Args, stdout, stderr io.Writer) ExitCode {
	// Repo-root resolution.
	repoRoot := args.RepoRootOver
	if repoRoot == "" {
		repoRoot = gitRepoRoot(ctx)
	}
	if repoRoot == "" {
		cwd, _ := os.Getwd()
		repoRoot = cwd
	}

	// Receipt-path resolution.
	receiptPath := args.Receipt
	if !filepath.IsAbs(receiptPath) {
		abs, err := filepath.Abs(receiptPath)
		if err == nil {
			receiptPath = abs
		}
	}
	if _, err := os.Stat(receiptPath); err != nil {
		fmt.Fprintf(stderr, "V: receipt not found: %s\n", receiptPath)
		return ExitVError
	}
	// Cue availability check.
	if !cueAvailable(ctx) {
		fmt.Fprintln(stderr, "V: `cue` is required but not runnable on PATH.")
		return ExitVError
	}

	verdict, err := Validate(ctx, ValidateOptions{
		ReceiptPath:    receiptPath,
		ContractPath:   args.Contract,
		RepoRoot:       repoRoot,
		StructuralOnly: args.StructuralOnly,
	})
	if err != nil {
		fmt.Fprintf(stderr, "V: internal error: %s\n", err)
		return ExitVError
	}

	if args.JSON {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(verdict); err != nil {
			fmt.Fprintf(stderr, "V: emit json: %s\n", err)
			return ExitVError
		}
	} else {
		fmt.Fprintln(stdout, RenderProse(verdict, args.Receipt))
	}

	if verdict.Result == ResultPASS {
		return ExitPASS
	}
	return ExitFAIL
}

// RenderProse returns a human-readable summary of the verdict.
//
// The receipt-path string echoed back is the operator-supplied path (so
// relative paths render as the operator typed them).
func RenderProse(v ValidationVerdict, receiptPath string) string {
	var b strings.Builder
	icon := "PASS"
	if v.Result != ResultPASS {
		icon = "FAIL"
	}
	fmt.Fprintf(&b, "V verdict on %s: %s\n", receiptPath, icon)
	fmt.Fprintf(&b, "  validator: %s @ %s\n", v.Provenance.ValidatorIdentity, v.Provenance.ValidatorVersion)
	fmt.Fprintf(&b, "  checked_at: %s", v.Provenance.CheckedAt)
	if len(v.FailedPredicates) > 0 {
		fmt.Fprintf(&b, "\n  failed_predicates (%d):", len(v.FailedPredicates))
		for _, fp := range v.FailedPredicates {
			fmt.Fprintf(&b, "\n    - [%s] %s", fp.Predicate, fp.Diagnostic)
			if fp.EvidenceRef != "" {
				fmt.Fprintf(&b, "\n        evidence_ref: %s", fp.EvidenceRef)
			}
		}
	}
	if len(v.Warnings) > 0 {
		fmt.Fprintf(&b, "\n  warnings (%d):", len(v.Warnings))
		for _, w := range v.Warnings {
			fmt.Fprintf(&b, "\n    - %s", w)
		}
	}
	return b.String()
}
