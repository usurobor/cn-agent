// Package cellrun is the domain logic behind `cn cell run`: read a serialized
// cell spec, fill its parameter holes, run one episode through the cellkernel,
// and emit a generic episode receipt. It owns the IO, parsing, and rendering;
// internal/cli holds only a thin dispatch wrapper (eng/go §2.18).
//
// It owns no GitHub, ref, PR, or custody policy (Pi β #31 C3): --contract is a
// local path or "-" for stdin; there is no network access. The emitted receipt
// is a generic cnos.cellkernel.episode-receipt.v0 with protocol_validated=false
// — the spec's protocol_id is declared provenance, never a validated claim.
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
  kernel. Prints a generic episode receipt (cnos.cellkernel.episode-receipt.v0)
  as JSON to stdout. The spec's protocol_id is carried as declared provenance
  only (protocol_validated=false in v0).

  Exit status: 0 accepted; 1 non-accepted terminal or needs_repair; 2 usage
  error or seat malfunction.

FLAGS:
  --contract <path|->    serialized cell spec; "-" reads stdin (required, once)
  --param <name>=<value> fill a declared parameter hole (repeatable, no dups)`

// Run executes one `cn cell run` invocation and returns its exit code:
// 0 accepted, 1 non-accepted terminal / needs_repair, 2 usage error or seat
// malfunction.
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

	kspec, meta, mode, err := resolved.Build()
	if err != nil {
		fmt.Fprintf(stderr, "✗ %v\n", err)
		return 2
	}

	res, err := cellkernel.RunEpisode(ctx, kspec, cellkernel.WithMeta(meta))
	if err != nil {
		fmt.Fprintf(stderr, "✗ episode malfunction: %v\n", err)
		return 2
	}

	// Self-check: the emitted receipt must independently re-verify (Pi #33 D2).
	if err := cellkernel.VerifyReceipt(res.Receipt); err != nil {
		fmt.Fprintf(stderr, "✗ receipt failed self-verification: %v\n", err)
		return 2
	}

	if err := writeReceipt(stdout, mode, res); err != nil {
		fmt.Fprintf(stderr, "✗ %v\n", err)
		return 2
	}

	if res.Status == cellkernel.Accepted {
		return 0
	}
	return 1
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

// receiptOutput is the CLI's structured generic episode receipt: the kernel's
// self-verifying receipt plus the runner's mode/verdict framing. It vets against
// schemas/cdd/episode-receipt.cue.
type receiptOutput struct {
	ReceiptSchema      string                    `json:"receipt_schema"`
	ProtocolValidated  bool                      `json:"protocol_validated"`
	ExecutionMode      string                    `json:"execution_mode"`
	Status             string                    `json:"status"`
	Decision           string                    `json:"decision"`
	Verdict            cellkernel.Verdict        `json:"verdict"`
	cellkernel.Receipt                           // embedded self-verifying receipt
	Repair             *cellkernel.RepairRequest `json:"repair,omitempty"`
}

func writeReceipt(w io.Writer, mode string, res cellkernel.EpisodeResult) error {
	out := receiptOutput{
		ReceiptSchema:     cellspec.EpisodeReceiptSchema,
		ProtocolValidated: false, // v0 runs no protocol-specific validation.
		ExecutionMode:     mode,
		Status:            string(res.Status),
		Decision:          string(res.Decision),
		Verdict:           res.Verdict,
		Receipt:           res.Receipt,
		Repair:            res.Repair,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("encode receipt: %w", err)
	}
	return nil
}
