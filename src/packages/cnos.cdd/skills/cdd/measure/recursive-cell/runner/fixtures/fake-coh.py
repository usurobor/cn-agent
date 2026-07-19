#!/usr/bin/env python3
import argparse
import hashlib
import json
import math
from pathlib import Path
import sys

if "--version" in sys.argv:
    print("coh fixture 0.1 (deterministic)")
    raise SystemExit(0)

parser = argparse.ArgumentParser()
parser.add_argument("--target", required=True)
parser.add_argument("--registry")
parser.add_argument("--instruction")
parser.add_argument("--root")
parser.add_argument("--mode")
parser.add_argument("--emit-prompt")
parser.add_argument("--llm-response")
parser.add_argument("--output")
args = parser.parse_args()

root = Path(args.root) if args.root else Path(".")
registry = Path(args.registry) if args.registry else Path("")
instruction = Path(args.instruction) if args.instruction else Path("")
if registry.is_absolute() or instruction.is_absolute():
    print("fixture coh authority paths must be relative to --root", file=sys.stderr)
    raise SystemExit(2)
registry = root / registry
instruction = root / instruction
if not args.registry or not registry.is_file():
    print("fixture coh requires an existing registry", file=sys.stderr)
    raise SystemExit(2)
if not args.instruction or not instruction.is_file():
    print("fixture coh requires an existing instruction", file=sys.stderr)
    raise SystemExit(2)
if not args.root or not root.is_dir():
    print("fixture coh requires an existing root", file=sys.stderr)
    raise SystemExit(2)

if args.emit_prompt:
    if args.mode != "llm" or args.llm_response or args.output:
        print("fixture coh emit route requires only --mode llm and --emit-prompt", file=sys.stderr)
        raise SystemExit(2)
    path = Path(args.emit_prompt)
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(
        "fixture-prompt-v2\n"
        f"target={args.target}\n"
        f"registry={args.registry}\n"
        f"registry_sha256={hashlib.sha256(registry.read_bytes()).hexdigest()}\n"
        f"instruction={args.instruction}\n"
        f"instruction_sha256={hashlib.sha256(instruction.read_bytes()).hexdigest()}\n",
        encoding="utf-8",
    )
    raise SystemExit(0)

if args.mode != "hybrid" or not args.llm_response or not args.output:
    print("fixture coh ingest route requires --mode hybrid, --llm-response, and --output", file=sys.stderr)
    raise SystemExit(2)

with Path(args.llm_response).open(encoding="utf-8") as stream:
    witness = json.load(stream)
scores = [witness["alpha"], witness["beta"], witness["gamma"]]
aggregate = math.prod(scores) ** (1.0 / 3.0)
report = {
    "target": args.target,
    "mode": "hybrid",
    "alpha": scores[0],
    "beta": scores[1],
    "gamma": scores[2],
    "bottleneck_axis": min(["alpha", "beta", "gamma"], key=lambda axis: (witness[axis], ["alpha", "beta", "gamma"].index(axis))),
    "provenance": {
        "aggregate_math": {"C_sigma_math": aggregate},
        "aggregate_numeric": {"C_sigma_num": aggregate, "epsilon": 0.00001},
    },
}
output = Path(args.output)
output.mkdir(parents=True, exist_ok=True)
(output / "report.json").write_text(json.dumps(report, indent=2, sort_keys=True) + "\n", encoding="utf-8")
