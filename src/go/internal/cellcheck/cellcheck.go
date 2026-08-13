// Package cellcheck runs ONE closed, compiled-in checker recipe against a
// candidate directory and reports what happened as a typed observation.
//
// The recipe is closed on purpose. A cell cannot supply commands, arguments,
// or an environment: a checker whose argv came from a declaration would let the
// thing being judged choose the judge, and "the tests pass" would mean whatever
// the candidate said it meant. What a cell may do is decline to be measured by
// this recipe; it may not redefine it.
//
// The observation is not evidence. Its consumer is the CDS assessing fill,
// which runs the recipe against the candidate it reconstructed and lets the
// result FORCE one catalogue unit — a `fail` becomes a finding and an
// `unavailable` becomes `unverified`, with no cognitive seat able to override
// either. What the observation is not is an artifact: it does not enter the
// record as evidence, and nothing here mints one.
//
// It runs against a DIRECTORY and knows nothing about views, matters, or
// subjects. The directory must be the root of the candidate repository — the
// recipe's paths are repo-root-relative and the changed-path query asks git
// about that root.
package cellcheck

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// RecipeID names this recipe. It is part of the observation because a reader
// must be able to tell WHICH closed recipe produced a status, not merely that
// something ran.
const RecipeID = "cnos.project-verify.v0"

// maxTailBytes bounds the output one step reports. A step's failure is
// explained by its last lines; carrying an unbounded build log would put an
// unbounded value into a value that travels.
const maxTailBytes = 4 << 10

// Status is the closed outcome set. `unavailable` is NOT a kind of failure: a
// step that could not start proves nothing about the candidate, and collapsing
// it into `fail` would blame a candidate for a missing binary.
type Status string

const (
	Pass        Status = "pass"
	Fail        Status = "fail"
	Unavailable Status = "unavailable"
)

// Step is one recipe step and what it did.
type Step struct {
	Name   string
	Status Status
	// Exit is the process exit code; -1 when the step never ran. It is NOT the
	// status: `format` fails with exit 0, because `gofmt -l` reports by listing.
	Exit int
	// Tail is the bounded tail of the step's combined stdout+stderr.
	Tail string
}

// Observation is the whole run.
type Observation struct {
	Recipe string
	// Candidate identifies what the recipe ran against. Run LEAVES IT EMPTY and
	// that is deliberate: this package is handed a directory, and a digest it
	// invented from that directory would not be the identity a caller means by
	// "the candidate". The caller that materialized the directory from a view
	// knows that digest and assigns it; until such a caller exists the field is
	// empty in every observation this package produces.
	Candidate string
	// Steps are the steps that RAN, in order. The list is short when the run
	// stopped early — see Run.
	Steps  []Step
	Status Status
}

// recipe is the closed step list, in order. `go -C src/go` because the module
// root is a fixed subdirectory of the candidate repository; the recipe is this
// project's, not a general one.
//
// `-count=1` on the test step because a cached PASS is not a result: without it
// a candidate whose tests were never re-run against its own change would report
// green from another tree's cache.
//
// `cue vet` is deliberately ABSENT. The corpus gate needs a `cue` binary that
// nothing in this runtime declares or installs, and a missing tool must be
// `unavailable` rather than a silent skip — so a step that would be unavailable
// on every ordinary machine would make every observation unavailable, and no
// run could accept. Adding it waits on the binary being declared.
var recipe = []struct {
	name string
	argv []string
}{
	{"build", []string{"go", "-C", "src/go", "build", "./..."}},
	{"vet", []string{"go", "-C", "src/go", "vet", "./..."}},
	{"test", []string{"go", "-C", "src/go", "test", "-count=1", "./..."}},
}

// Run executes the recipe against `dir` and returns what it observed. It never
// returns an error: an observation is its own outcome, and a checker that could
// report "I failed to tell you" through a second channel would have two ways to
// say `unavailable`.
//
// It STOPS at the first step that does not pass. Running `go test` after a
// compile error reports one defect twice and costs the whole test suite's
// runtime; the step list records exactly how far the recipe got, so where it
// stopped is read off the observation rather than guessed.
func Run(ctx context.Context, dir, baseSHA string) Observation {
	obs := Observation{Recipe: RecipeID, Status: Pass}
	for _, s := range recipe {
		st, _ := run(ctx, dir, s.name, s.argv)
		obs.Steps = append(obs.Steps, st)
		if st.Status != Pass {
			obs.Status = st.Status
			return obs
		}
	}
	st := formatStep(ctx, dir, baseSHA)
	obs.Steps = append(obs.Steps, st)
	obs.Status = st.Status
	return obs
}

// formatStep checks gofmt over ONLY the .go paths the candidate changed.
//
// Repo-wide would be red for every candidate regardless of its patch: `gofmt -l
// src/go` lists 18 files on a clean tree of this repository today. That is
// measured, not hypothetical — a repo-wide step would fail on the base commit
// itself, so no candidate could ever reach `pass` and the whole recipe would
// carry no information about any change.
func formatStep(ctx context.Context, dir, baseSHA string) Step {
	paths, err := changedGoFiles(ctx, dir, baseSHA)
	if err != nil {
		// The changed set is what makes this step meaningful. Falling back to
		// repo-wide, or to "nothing changed", would each turn an unanswerable
		// question into an answer.
		return Step{Name: "format", Status: Unavailable, Exit: -1, Tail: tail(err.Error())}
	}
	if len(paths) == 0 {
		return Step{Name: "format", Status: Pass, Tail: "no changed .go paths"}
	}
	st, out := run(ctx, dir, "format", append([]string{"gofmt", "-l"}, paths...))
	// `gofmt -l` REPORTS BY LISTING: it exits 0 whether or not it found
	// unformatted files. Trusting the exit code here would make this step pass
	// for every candidate, which is exactly the vacuous guard the step exists
	// to avoid. A non-empty listing is the failure.
	if st.Status == Pass && strings.TrimSpace(out) != "" {
		st.Status = Fail
	}
	return st
}

// changedGoFiles is the .go paths the candidate changed relative to the PINNED
// BASE, in path order.
//
// Against the base, not against HEAD, and the distinction is the same one
// cellwork.Diff already carries: a seat has a shell and therefore git, so it
// may commit its own work. Once it does, `git status` is clean, and a
// HEAD-relative changed set is empty — this step would then report "no changed
// .go paths" and PASS with an unformatted file in the change. A gate that goes
// green because the seat tidied up is worse than no gate.
//
// Deleted paths and rename sources are dropped — there is nothing there to
// format — and only regular files survive, so a changed symlink cannot be
// handed to gofmt as if it were source. Paths git reports are relative to the
// REPOSITORY ROOT, which is not necessarily `dir`, so the root is asked for
// rather than assumed: resolving it wrongly would drop every path and pass.
func changedGoFiles(ctx context.Context, dir, baseSHA string) ([]string, error) {
	if baseSHA == "" {
		// Without a base there is no "what the candidate changed", and the two
		// available fallbacks — repo-wide, or nothing — are both answers to a
		// question that cannot be answered here.
		return nil, fmt.Errorf("cellcheck: no pinned base to measure the candidate against")
	}
	root, err := gitOut(ctx, dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("cellcheck: locate the candidate's repository root: %w", err)
	}
	// The recipe's other steps run `go -C src/go ...` relative to `dir`, so the
	// whole recipe already means "dir is the candidate repository root". git
	// reports paths relative to the ROOT, so a `dir` below it would drop every
	// path and this step would pass from having looked in the wrong place.
	// Refuse instead: an unanswerable question must not return an answer.
	if same, err := sameDir(root, dir); err != nil {
		return nil, fmt.Errorf("cellcheck: compare the candidate directory with its repository root: %w", err)
	} else if !same {
		return nil, fmt.Errorf("cellcheck: %s is not the candidate repository root (%s); the recipe measures a whole repository", dir, root)
	}
	// Tracked changes against the base, plus paths the candidate created and
	// never staged. `git diff` alone cannot see an untracked file, and a new
	// unformatted source file is exactly the case this step must catch.
	tracked, err := gitOut(ctx, dir, "diff", "--name-only", "-z", "--diff-filter=ACMRTUXB", baseSHA)
	if err != nil {
		return nil, fmt.Errorf("cellcheck: list the candidate's changed paths: %w", err)
	}
	untracked, err := gitOut(ctx, dir, "ls-files", "-z", "--others", "--exclude-standard")
	if err != nil {
		return nil, fmt.Errorf("cellcheck: list the candidate's new paths: %w", err)
	}

	seen := make(map[string]bool)
	var paths []string
	for _, out := range []string{tracked, untracked} {
		// `-z`, not the default: without it git QUOTES paths containing spaces
		// or non-ASCII, and the quoted name is not the name on disk.
		for _, path := range strings.Split(strings.TrimSuffix(out, "\x00"), "\x00") {
			if path == "" || !strings.HasSuffix(path, ".go") || seen[path] {
				continue
			}
			info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(path)))
			switch {
			case os.IsNotExist(err):
				continue // deleted, or a rename source: nothing to format
			case err != nil:
				// Any OTHER stat failure is a question this step cannot answer.
				// Skipping it would turn "I could not look" into "nothing
				// changed", which reads as a pass.
				return nil, fmt.Errorf("cellcheck: stat candidate path %s: %w", path, err)
			case !info.Mode().IsRegular():
				continue
			}
			seen[path] = true
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	return paths, nil
}

// sameDir reports whether two paths name the same directory, resolving symlinks
// so that a temp dir reached through /tmp and through /private/tmp compares
// equal rather than looking like two places.
func sameDir(a, b string) (bool, error) {
	ra, err := filepath.EvalSymlinks(a)
	if err != nil {
		return false, err
	}
	rb, err := filepath.EvalSymlinks(b)
	if err != nil {
		return false, err
	}
	return ra == rb, nil
}

// gitOut runs one git command in dir and returns its stdout.
func gitOut(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}

// run executes one step and classifies the result. It returns the step and the
// step's UNTRUNCATED combined output, because one caller decides on the output
// itself rather than on the exit code.
func run(ctx context.Context, dir, name string, argv []string) (Step, string) {
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = dir
	// The environment is inherited: `go` needs its cache and toolchain
	// configuration, and a hand-built environment would measure a build the
	// project does not otherwise perform.
	// Bounded as the child writes, not after it exits: `go test ./...` on a
	// loud candidate can emit hundreds of megabytes, and a limit applied to a
	// fully buffered result is not a limit — the same rule cellwork.git states
	// for the same reason. The TAIL is what a reader needs, so the buffer keeps
	// the last bytes and drops the rest.
	buf := &tailBuffer{max: maxStepOutputBytes}
	cmd.Stdout, cmd.Stderr = buf, buf
	err := cmd.Run()
	out := buf.String()
	st := Step{Name: name, Tail: tail(string(out))}
	switch {
	case ctx.Err() != nil:
		// A cancelled or timed-out step was killed; the non-zero exit that
		// leaves behind is the signal, not the candidate's verdict.
		st.Status, st.Exit = Unavailable, -1
		if st.Tail == "" {
			st.Tail = ctx.Err().Error()
		}
	case err == nil:
		st.Status, st.Exit = Pass, 0
	default:
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			// It ran and said no. That is about the candidate.
			st.Status, st.Exit = Fail, ee.ExitCode()
		} else {
			// It never started: binary missing, directory unreadable. A
			// candidate must not be blamed for that.
			//
			// The discriminator is exactly "did the process start", and it is
			// narrower than "is this the machine's fault": a step that starts
			// and then dies on an unwritable build cache or an undownloadable
			// toolchain exits non-zero, so it lands in Fail above. Those are
			// machine conditions reported against the candidate, and this
			// predicate cannot separate them. Stated rather than implied,
			// because the sentence above it would otherwise promise more than
			// exec.ExitError can deliver.
			st.Status, st.Exit = Unavailable, -1
			if st.Tail == "" {
				st.Tail = tail(err.Error())
			}
		}
	}
	return st, string(out)
}

// maxStepOutputBytes bounds what one step's output may cost in memory. It is
// larger than maxTailBytes because a reader wants the tail while a diagnostic
// sometimes needs a little more context than the tail alone.
const maxStepOutputBytes = 1 << 20

// tailBuffer keeps the LAST max bytes written to it and discards the rest as
// they arrive, so a child that never stops talking costs a bounded amount of
// memory rather than all of it. It never reports a short write: a bound must
// fail the STEP here, not kill the child with a broken pipe and turn a loud
// test suite into an unavailable checker.
type tailBuffer struct {
	max  int
	buf  []byte
	over bool
}

func (b *tailBuffer) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	if len(b.buf) > b.max {
		b.over = true
		b.buf = b.buf[len(b.buf)-b.max:]
	}
	return len(p), nil
}

func (b *tailBuffer) String() string {
	if b.over {
		return "...[truncated]...\n" + string(b.buf)
	}
	return string(b.buf)
}

// tail keeps the LAST maxTailBytes: a failing build says what went wrong at
// the end, and the truncation is announced rather than silent.
func tail(s string) string {
	if len(s) <= maxTailBytes {
		return s
	}
	return "...[truncated]...\n" + s[len(s)-maxTailBytes:]
}
