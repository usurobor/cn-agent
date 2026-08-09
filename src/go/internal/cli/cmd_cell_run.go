package cli

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

// CellRunCmd implements "cn cell run" — the GitHub-free local runner. It reads a
// serialized cell spec (a compiled main.cell or a hand-authored contract file),
// fills its parameter holes from --param flags, runs one episode through the
// cellkernel, and prints the terminal receipt as JSON. Exit status: 0 accepted,
// 1 non-accepted terminal/needs_repair, 2 malfunction/usage error.
//
// It owns no GitHub, ref, PR, or custody policy (Pi β #31 C3): --contract is a
// local path or "-" for stdin; there are no network reads.
type CellRunCmd struct{}

func (c *CellRunCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:    "cell-run",
		Summary: "Run one local cell episode from a serialized spec (no GitHub/network)",
		Source:  SourceKernel,
		Tier:    TierKernel,
		// NeedsHub false: runs against an arbitrary local contract path/stdin.
	}
}

func (c *CellRunCmd) Help() string {
	return `cn cell run - Run one local cell episode through the CCNF kernel

USAGE:
  cn cell run --contract <path|-> [--param <name>=<value> ...]

DESCRIPTION:
  Reads a serialized cell spec (JSON) from a file or stdin, fills its typed
  parameter holes from --param flags, and runs a single episode through the
  kernel with stub alpha/beta (no cognition yet). Prints the terminal receipt
  as JSON to stdout.

  Exit status: 0 accepted; 1 non-accepted terminal or needs_repair; 2 usage
  error or seat malfunction.

FLAGS:
  --contract <path|->    serialized cell spec; "-" reads stdin
  --param <name>=<value> fill a declared parameter hole (repeatable)`
}

func (c *CellRunCmd) Run(ctx context.Context, inv Invocation) error {
	contractPath, params, err := parseCellRunArgs(inv.Args)
	if err != nil {
		fmt.Fprintf(inv.Stderr, "✗ %v\n\n%s\n", err, c.Help())
		return &CellRunExit{Code: 2}
	}

	data, err := readContract(contractPath, inv.Stdin)
	if err != nil {
		fmt.Fprintf(inv.Stderr, "✗ %v\n", err)
		return &CellRunExit{Code: 2}
	}

	spec, err := cellspec.Parse(data)
	if err != nil {
		fmt.Fprintf(inv.Stderr, "✗ %v\n", err)
		return &CellRunExit{Code: 2}
	}

	resolved, err := spec.Resolve(params)
	if err != nil {
		fmt.Fprintf(inv.Stderr, "✗ %v\n", err)
		return &CellRunExit{Code: 2}
	}

	res, err := cellkernel.RunEpisode(ctx, resolved.KernelSpec())
	if err != nil {
		fmt.Fprintf(inv.Stderr, "✗ episode malfunction: %v\n", err)
		return &CellRunExit{Code: 2}
	}

	if err := writeReceipt(inv.Stdout, spec, resolved, res); err != nil {
		fmt.Fprintf(inv.Stderr, "✗ %v\n", err)
		return &CellRunExit{Code: 2}
	}

	if res.Status == cellkernel.Accepted {
		return nil
	}
	return &CellRunExit{Code: 1}
}

func parseCellRunArgs(args []string) (contract string, params map[string]string, err error) {
	params = make(map[string]string)
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--contract":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--contract requires a path or '-'")
			}
			i++
			contract = args[i]
		case "--param":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--param requires <name>=<value>")
			}
			i++
			k, v, ok := strings.Cut(args[i], "=")
			if !ok || k == "" {
				return "", nil, fmt.Errorf("malformed --param %q (want <name>=<value>)", args[i])
			}
			params[k] = v
		default:
			return "", nil, fmt.Errorf("unknown argument %q", args[i])
		}
	}
	if contract == "" {
		return "", nil, fmt.Errorf("--contract is required")
	}
	return contract, params, nil
}

func readContract(path string, stdin io.Reader) ([]byte, error) {
	if path == "-" {
		data, err := io.ReadAll(stdin)
		if err != nil {
			return nil, fmt.Errorf("read contract from stdin: %w", err)
		}
		return data, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read contract %q: %w", path, err)
	}
	return data, nil
}

// receiptOutput is the CLI's structured receipt projection.
type receiptOutput struct {
	ContractID   string                    `json:"contract_id"`
	ProtocolID   string                    `json:"protocol_id"`
	Status       string                    `json:"status"`
	Decision     string                    `json:"decision"`
	Verdict      verdictOutput             `json:"verdict"`
	Params       map[string]string         `json:"params,omitempty"`
	AlphaSkills  []string                  `json:"alpha_skills"`
	BetaSkills   []string                  `json:"beta_skills"`
	Matter       string                    `json:"matter"`
	Review       reviewOutput              `json:"review"`
	EvidenceRefs []cellkernel.EvidenceRef  `json:"evidence_refs"`
	Repair       *cellkernel.RepairRequest `json:"repair,omitempty"`
}

type verdictOutput struct {
	Pass   bool     `json:"pass"`
	Failed []string `json:"failed,omitempty"`
}

type reviewOutput struct {
	Pass  bool   `json:"pass"`
	Notes string `json:"notes"`
}

func writeReceipt(w io.Writer, spec cellspec.CellSpec, r cellspec.Resolved, res cellkernel.EpisodeResult) error {
	out := receiptOutput{
		ContractID:   res.Contract.ID,
		ProtocolID:   spec.ProtocolID,
		Status:       string(res.Status),
		Decision:     string(res.Decision),
		Verdict:      verdictOutput{Pass: res.Verdict.Pass, Failed: res.Verdict.Failed},
		Params:       r.Params,
		AlphaSkills:  r.AlphaSkills,
		BetaSkills:   r.BetaSkills,
		Matter:       res.Matter.Data,
		Review:       reviewOutput{Pass: res.Review.Pass, Notes: res.Review.Notes},
		EvidenceRefs: res.Receipt.EvidenceRefs,
		Repair:       res.Repair,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("encode receipt: %w", err)
	}
	return nil
}

// CellRunExit carries the runner's exit code back to cmd/cn/main.go.
type CellRunExit struct{ Code int }

func (e *CellRunExit) Error() string {
	return fmt.Sprintf("cell run exited with status %d", e.Code)
}
