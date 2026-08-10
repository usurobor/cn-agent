// Package cellrun is the domain logic behind `cn cell run`: read a serialized
// cell spec, fill its parameter holes, run one episode through the cellkernel,
// and emit the terminal episode closure. It owns the IO, parsing, and
// rendering; internal/cli holds only a thin dispatch wrapper (eng/go §2.18).
//
// It owns no GitHub, ref, PR, or custody policy (Pi β #31 C3): --contract is a
// local path or "-" for stdin; there is no network access. The emitted closure
// is a generic cnos.cellkernel.episode-closure.v0 with protocol_validated=false
// — the spec's protocol_id is declared provenance, never a validated claim — and
// re-verifies whole via cellkernel.VerifyClosure.
package cellrun

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellspec"
)

// maxContractBytes bounds the serialized spec read from a file or stdin.
const maxContractBytes = 1 << 20 // 1 MiB

// Help is the usage text for `cn cell run`.
const Help = `cn cell run - Run one local cell episode through the CCNF kernel

USAGE:
  cn cell run --contract <path|-> [--param <name>=<value> ...]

DESCRIPTION:
  Reads a serialized cell spec (JSON) from a file or stdin, fills its typed
  parameter holes from --param flags, and runs a single episode through the
  kernel. Prints a generic episode closure (cnos.cellkernel.episode-closure.v0)
  as JSON to stdout; the whole closure re-verifies via VerifyClosure. The
  spec's protocol_id is declared provenance only (protocol_validated=false).

  Exit status: 0 accepted; 1 non-accepted terminal (needs_repair/degraded/
  rejected); 2 usage error or seat malfunction; 3 simulated (a non-authoritative
  stub run).

FLAGS:
  --contract <path|->    serialized cell spec; "-" reads stdin (required, once)
  --param <name>=<value> fill a declared parameter hole (repeatable, no dups)`

// Run executes one `cn cell run` invocation and returns its exit code.
func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	contractPath, params, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "✗ %v\n\n%s\n", err, Help)
		return 2
	}

	data, err := readContract(contractPath, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "✗ %v\n", err)
		return 2
	}

	spec, err := cellspec.Parse(data)
	if err != nil {
		fmt.Fprintf(stderr, "✗ %v\n", err)
		return 2
	}

	resolved, err := spec.Resolve(params)
	if err != nil {
		fmt.Fprintf(stderr, "✗ %v\n", err)
		return 2
	}

	kspec, meta, err := resolved.Build()
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
	// against the contract this invocation built — never the closure's own.
	if err := cellkernel.VerifyClosure(kspec.Contract, cl); err != nil {
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

func parseArgs(args []string) (contract string, params map[string]string, err error) {
	params = make(map[string]string)
	contractSet := false
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--contract":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--contract requires a path or '-'")
			}
			if contractSet {
				return "", nil, fmt.Errorf("--contract given more than once")
			}
			i++
			contract = args[i]
			contractSet = true
		case "--param":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--param requires <name>=<value>")
			}
			i++
			k, v, ok := strings.Cut(args[i], "=")
			if !ok || k == "" {
				return "", nil, fmt.Errorf("malformed --param %q (want <name>=<value>)", args[i])
			}
			if _, dup := params[k]; dup {
				return "", nil, fmt.Errorf("--param %q given more than once", k)
			}
			params[k] = v
		default:
			return "", nil, fmt.Errorf("unknown argument %q", args[i])
		}
	}
	if !contractSet {
		return "", nil, fmt.Errorf("--contract is required")
	}
	return contract, params, nil
}

func readContract(path string, stdin io.Reader) ([]byte, error) {
	var r io.Reader
	if path == "-" {
		r = stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open contract %q: %w", path, err)
		}
		defer f.Close()
		r = f
	}
	data, err := io.ReadAll(io.LimitReader(r, maxContractBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read contract: %w", err)
	}
	if len(data) > maxContractBytes {
		return nil, fmt.Errorf("contract exceeds %d bytes", maxContractBytes)
	}
	return data, nil
}
