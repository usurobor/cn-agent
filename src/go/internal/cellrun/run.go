// Package cellrun is the domain logic behind `cn cell run`: read a serialized
// cell spec and an optional run input, admit the run input, fill the spec's
// parameter holes, run one episode through the cellkernel, and emit the
// terminal episode closure. It owns the IO, parsing, and rendering;
// internal/cli holds only a thin dispatch wrapper (eng/go §2.18).
//
// It owns no GitHub, ref, PR, or custody policy (Pi β #31 C3): --contract and
// --input are local paths or "-" for stdin; there is no network access. The
// emitted closure is a generic cnos.cellkernel.episode-closure.v0 with
// protocol_validated=false — the spec's protocol_id is declared provenance,
// never a validated claim — and re-verifies whole via
// cellkernel.VerifyClosure.
//
// ORDER IS THE PROPERTY (D6). The admission door decodes and judges the run
// input, and the subject is pinned, BEFORE the fill registry's constructors are
// dispatched — so no seat exists, and therefore no provider can be invoked,
// until the run input has been admitted. That is a property of this function's
// shape, not of a policy flag; internal/cdsadmit/door_test.go is its witness.
//
// The door itself arrives IN THE REGISTRY, wired by the composition root beside
// the fills. This package therefore names no protocol: it imports no CDS
// package, holds no admission rule, and cannot read the receipt it prints. A
// registry with no door is a runtime that admits no run input — which is every
// cell that existed before run inputs did, and still runs.
package cellrun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellspec"
	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

// maxContractBytes bounds the serialized spec read from a file or stdin. The
// run input is bounded by the same number: both are hand-authored documents
// read into memory before anything is decided about them, and a second
// constant would be two numbers to justify instead of one.
const maxContractBytes = 1 << 20 // 1 MiB

// Help is the usage text for `cn cell run`.
const Help = `cn cell run - Run one local cell episode through the CCNF kernel

USAGE:
  cn cell run --contract <path|-> [--input <path|->] [--param <name>=<value> ...]

DESCRIPTION:
  Reads a serialized cell spec (JSON) from a file or stdin, fills its typed
  parameter holes from --param flags, and runs a single episode through the
  kernel. Prints a generic episode closure (cnos.cellkernel.episode-closure.v0)
  as JSON to stdout; the whole closure re-verifies via VerifyClosure. The
  spec's protocol_id is declared provenance only (protocol_validated=false).

  --input supplies the per-run contract: a run-input document carrying an
  issue, a design, and a git subject reference. It is admitted structurally
  before any seat is constructed; the subject's base is pinned to an exact
  commit once, and the admitted issue, design and pinned subject are frozen
  into the contract and covered by the closure's one scope-lift digest. A
  refused run input prints an admission receipt and mints no episode.

  Exit status: 0 accepted; 1 non-accepted terminal (needs_repair/degraded/
  rejected); 2 usage error or seat malfunction; 3 simulated (a non-authoritative
  stub run); 4 run input refused at admission.

FLAGS:
  --contract <path|->    serialized cell spec; "-" reads stdin (required, once)
  --input <path|->       run input: issue + design + subject (optional, once)
  --param <name>=<value> fill a declared parameter hole (repeatable, no dups)`

// Run executes one `cn cell run` invocation and returns its exit code. The
// fill registry arrives already assembled from the application composition
// root: this package dispatches it and never learns what a fill needs.
func Run(ctx context.Context, reg cellfill.Registry, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	inv, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "✗ %v\n\n%s\n", err, Help)
		return 2
	}

	// The run input is read, admitted and pinned FIRST. Everything below this
	// block constructs seats; nothing above it can.
	bind, code, err := admit(ctx, reg.Door, inv.input, stdin, stdout)
	if err != nil {
		fmt.Fprintf(stderr, "✗ %v\n", err)
		return code
	}

	data, err := readDoc("contract", inv.contract, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "✗ %v\n", err)
		return 2
	}

	spec, err := cellspec.Parse(data)
	if err != nil {
		fmt.Fprintf(stderr, "✗ %v\n", err)
		return 2
	}

	resolved, err := spec.Resolve(inv.params)
	if err != nil {
		fmt.Fprintf(stderr, "✗ %v\n", err)
		return 2
	}

	kspec, meta, err := resolved.Build(ctx, reg, bind)
	if err != nil {
		fmt.Fprintf(stderr, "✗ %v\n", err)
		return 2
	}

	cl, err := cellkernel.RunEpisode(ctx, kspec, meta)
	if err != nil {
		fmt.Fprintf(stderr, "✗ episode malfunction: %v\n", err)
		return 2
	}

	// Self-check: the emitted closure must independently re-verify whole,
	// against the contract AND invocation metadata this invocation built —
	// never the closure's own.
	if err := cellkernel.VerifyClosure(kspec.Contract, meta, cl); err != nil {
		fmt.Fprintf(stderr, "✗ closure failed self-verification: %v\n", err)
		return 2
	}

	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cl); err != nil {
		fmt.Fprintf(stderr, "✗ encode closure: %v\n", err)
		return 2
	}

	switch cl.Status {
	case cellkernel.Accepted:
		return 0
	case cellkernel.Simulated:
		return 3 // ran, but non-authoritative — never ordinary accepted (Pi D5).
	default:
		return 1 // needs_repair / degraded / rejected.
	}
}

// invocation is the parsed command line. `input` empty means no run input was
// supplied — the pre-run-input behaviour, kept because a cell that acts on
// nothing is still a cell (the bool and stub corpus cells are exactly that).
// Whether a given PROFILE may omit it is admission's rule to state, not this
// parser's.
type invocation struct {
	contract string
	input    string
	params   map[string]string
}

func parseArgs(args []string) (invocation, error) {
	inv := invocation{params: make(map[string]string)}
	contractSet, inputSet := false, false
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--contract", "--input":
			flag := args[i]
			set, dst := &contractSet, &inv.contract
			if flag == "--input" {
				set, dst = &inputSet, &inv.input
			}
			if i+1 >= len(args) {
				return invocation{}, fmt.Errorf("%s requires a path or '-'", flag)
			}
			if *set {
				return invocation{}, fmt.Errorf("%s given more than once", flag)
			}
			i++
			*dst, *set = args[i], true
		case "--param":
			if i+1 >= len(args) {
				return invocation{}, fmt.Errorf("--param requires <name>=<value>")
			}
			i++
			k, v, ok := strings.Cut(args[i], "=")
			if !ok || k == "" {
				return invocation{}, fmt.Errorf("malformed --param %q (want <name>=<value>)", args[i])
			}
			if _, dup := inv.params[k]; dup {
				return invocation{}, fmt.Errorf("--param %q given more than once", k)
			}
			inv.params[k] = v
		default:
			return invocation{}, fmt.Errorf("unknown argument %q", args[i])
		}
	}
	if !contractSet {
		return invocation{}, fmt.Errorf("--contract is required")
	}
	// One stream, one document. Reading both from stdin would silently give
	// the whole stream to whichever read ran first and an empty one to the
	// other, which is a corrupt run reported as a malformed file.
	if inv.contract == "-" && inv.input == "-" {
		return invocation{}, fmt.Errorf("--contract and --input cannot both read stdin")
	}
	return inv, nil
}

// admit reads the run input, dispatches the registry's admission door, then
// pins the admitted subject exactly once (D4) — before any seat exists. It
// returns the binding the contract freezes.
//
// The door arrives in the registry, like a fill. This function names no
// profile: it does not know what an issue is, what makes one admissible, or
// what the receipt it prints says. That is why `cn cell run` stays the generic
// runner — a second protocol is a different door in the registry, not a branch
// here.
//
// An absent --input yields a zero binding and no refusal: a cell with no run
// input is the shape every cell had before run inputs existed. The returned
// code is the exit status for the returned error.
func admit(ctx context.Context, door cellfill.Door, path string, stdin io.Reader, stdout io.Writer) (cellspec.Binding, int, error) {
	if path == "" {
		return cellspec.Binding{}, 0, nil
	}
	if door == nil {
		// Not a refusal: refusing would claim a door judged this document and
		// found it inadmissible. Nothing judged it, and saying so is the whole
		// difference between "your input is wrong" and "this binary admits no
		// input".
		return cellspec.Binding{}, 2, fmt.Errorf("--input given but this runtime carries no admission door")
	}
	raw, err := readDoc("run input", path, stdin)
	if err != nil {
		return cellspec.Binding{}, 2, err
	}

	admitted, receipt, err := door(raw)
	if err != nil {
		// A refusal is a RESULT, not a usage error: the receipt goes to
		// stdout under its own kind, and the exit code is its own. An episode
		// closure is never emitted, because no episode exists before a
		// contract is admitted. The receipt is emitted as the door serialized
		// it — this function never decodes it, so it learns no profile
		// vocabulary in order to print one.
		if errors.Is(err, cellfill.ErrRefused) {
			enc := json.NewEncoder(stdout)
			enc.SetIndent("", "  ")
			if encErr := enc.Encode(receipt); encErr != nil {
				return cellspec.Binding{}, 2, fmt.Errorf("encode admission receipt: %w", encErr)
			}
			return cellspec.Binding{}, 4, err
		}
		return cellspec.Binding{}, 2, err
	}

	// Pin ONCE, here, and never again: what the contract carries from this
	// point is an exact commit, so the two stations cannot each re-resolve a
	// moving name and measure against different trees while the record claims
	// one.
	pinned, err := cellwork.Pin(ctx, admitted.Subject)
	if err != nil {
		return cellspec.Binding{}, 2, err
	}
	return cellspec.Binding{Issue: admitted.Issue, Design: admitted.Design, Subject: pinned}, 0, nil
}

// readDoc reads one bounded document from a path or stdin. The IO wrapper over
// the pure decoders (eng/go §2.17): `what` names the document in diagnostics,
// so a missing spec and a missing run input do not read alike.
func readDoc(what, path string, stdin io.Reader) ([]byte, error) {
	var r io.Reader
	if path == "-" {
		r = stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open %s %q: %w", what, path, err)
		}
		defer f.Close()
		r = f
	}
	data, err := io.ReadAll(io.LimitReader(r, maxContractBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", what, err)
	}
	if len(data) > maxContractBytes {
		return nil, fmt.Errorf("%s exceeds %d bytes", what, maxContractBytes)
	}
	return data, nil
}
