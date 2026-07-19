#!/usr/bin/env python3
"""Deterministic State-A runner for the frozen cnos#662 recursive-cell CM."""

from __future__ import annotations

import argparse
import hashlib
import json
import math
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile
from datetime import datetime, timezone

TARGETS = [
    "cc662-system",
    "cc662-l0",
    "cc662-l1",
    "cc662-l2",
    "cc662-l3",
    "cc662-l4",
]
LEVEL_TARGETS = TARGETS[1:]
LEVEL_NAMES = ["L0", "L1", "L2", "L3", "L4"]
AXES = ["alpha", "beta", "gamma"]
INVARIANTS = [f"H{i:02d}" for i in range(1, 14)]
SHA40 = re.compile(r"^[0-9a-f]{40}$")
RFC3339_UTC = re.compile(r"^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?Z$")
SEVERITY_ORDER = {"none": 0, "cosmetic": 1, "isolated": 2, "systemic": 3}
CHECKLISTS = {
    "alpha": ["naming-drift", "duplicate-definition", "internal-contradiction", "unstable-boundary"],
    "beta": ["broken-reference", "authority-conflict", "fact-drift", "undeclared-relationship"],
    "gamma": ["unowned-change-path", "generated-canonical-confusion", "missing-migration-rule", "stale-transitional-marker"],
}
CARD_CATEGORIES = {category: axis for axis, categories in CHECKLISTS.items() for category in categories}
PROHIBITED_COMPUTED_FIELDS = {
    "coh", "Coh", "c_sigma", "C_sigma", "c_sigma_num", "C_sigma_num",
    "c_sigma_math", "C_sigma_math", "coherence", "coherence_score",
}


class Refusal(RuntimeError):
    pass


def sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for block in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(block)
    return digest.hexdigest()


def path_framed_sha256(base: Path, paths: list[Path]) -> str:
    """Hash sorted relative path, byte length, and bytes with NUL framing."""
    digest = hashlib.sha256()
    for path in sorted(paths, key=lambda item: item.relative_to(base).as_posix()):
        relative = path.relative_to(base).as_posix().encode("utf-8")
        content = path.read_bytes()
        digest.update(relative)
        digest.update(b"\0")
        digest.update(str(len(content)).encode("ascii"))
        digest.update(b"\0")
        digest.update(content)
        digest.update(b"\0")
    return digest.hexdigest()


def load_json(path: Path):
    try:
        with path.open(encoding="utf-8") as stream:
            return json.load(stream)
    except (OSError, json.JSONDecodeError) as error:
        raise Refusal(f"invalid JSON {path}: {error}") from error


def write_json(path: Path, value) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def snapshot_file(source: Path, destination: Path, label: str) -> None:
    try:
        destination.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(source, destination)
    except OSError as error:
        raise Refusal(f"could not snapshot {label}: {error}") from error


def run(command, *, cwd: Path | None = None, capture: bool = False) -> str:
    try:
        result = subprocess.run(
            [str(part) for part in command], cwd=cwd, check=True,
            text=True, stdout=subprocess.PIPE if capture else None,
            stderr=subprocess.PIPE if capture else None,
        )
    except subprocess.CalledProcessError as error:
        detail = (error.stderr or error.stdout or "").strip()
        raise Refusal(f"command failed ({' '.join(map(str, command))}): {detail}") from error
    return (result.stdout or "").strip()


def require_unit_interval(value, field: str) -> None:
    if isinstance(value, bool) or not isinstance(value, (int, float)) or not 0.0 <= value <= 1.0:
        raise Refusal(f"witness field {field} must be a number in [0,1]")


def reject_computed_fields(value, path="$") -> None:
    if isinstance(value, dict):
        for key, nested in value.items():
            if key in PROHIBITED_COMPUTED_FIELDS:
                raise Refusal(f"witness contains engine-owned computed field {path}.{key}")
            reject_computed_fields(nested, f"{path}.{key}")
    elif isinstance(value, list):
        for index, nested in enumerate(value):
            reject_computed_fields(nested, f"{path}[{index}]")


def validate_witness(witness, target: str) -> None:
    required = {
        "target", "alpha", "beta", "gamma", "delta_alpha_beta",
        "delta_beta_gamma", "delta_gamma_alpha", "bottleneck_axis",
        "confidence", "summary", "axis_evidence", "defect_cards",
        "unresolved_ambiguity", "next_fixes",
    }
    if not isinstance(witness, dict) or set(witness) != required:
        raise Refusal(f"witness {target} must have the exact TSC v3.2.4 top-level field set")
    if witness["target"] != target:
        raise Refusal(f"witness target mismatch: expected {target}, got {witness['target']!r}")
    for field in ["alpha", "beta", "gamma", "delta_alpha_beta", "delta_beta_gamma", "delta_gamma_alpha", "confidence"]:
        require_unit_interval(witness[field], field)
    if witness["bottleneck_axis"] not in AXES:
        raise Refusal(f"witness {target} has invalid bottleneck_axis")
    if not isinstance(witness["summary"], str) or not witness["summary"].strip():
        raise Refusal(f"witness {target} has empty summary")
    if not isinstance(witness["axis_evidence"], dict) or set(witness["axis_evidence"]) != set(AXES):
        raise Refusal(f"witness {target} axis_evidence must contain alpha/beta/gamma exactly")
    if not isinstance(witness["defect_cards"], list):
        raise Refusal(f"witness {target} defect_cards must be a list")
    aggregates = {axis: {category: [] for category in CHECKLISTS[axis]} for axis in AXES}
    for card in witness["defect_cards"]:
        keys = {"id", "primary_axis", "category", "severity", "evidence", "summary", "secondary_axes"}
        if not isinstance(card, dict) or set(card) != keys:
            raise Refusal(f"witness {target} has malformed defect card")
        axis = card["primary_axis"]
        category = card["category"]
        severity = card["severity"]
        if (
            not isinstance(card["id"], str) or not card["id"].strip()
            or axis not in AXES
            or CARD_CATEGORIES.get(category) != axis
            or severity not in SEVERITY_ORDER
            or severity == "none"
        ):
            raise Refusal(f"witness {target} has invalid defect-card axis/category/severity")
        if not isinstance(card["evidence"], str) or not card["evidence"].strip():
            raise Refusal(f"witness {target} defect card lacks evidence")
        if not isinstance(card["summary"], str) or not card["summary"].strip():
            raise Refusal(f"witness {target} defect card lacks summary")
        if (
            not isinstance(card["secondary_axes"], list)
            or len(card["secondary_axes"]) != len(set(card["secondary_axes"]))
            or any(a not in AXES or a == axis for a in card["secondary_axes"])
        ):
            raise Refusal(f"witness {target} has invalid secondary_axes")
        aggregates[axis][category].append(severity)
    for axis in AXES:
        evidence = witness["axis_evidence"][axis]
        if not isinstance(evidence, dict) or set(evidence) != {"positive", "negative", "reason", "checklist"}:
            raise Refusal(f"witness {target} has malformed {axis} axis evidence")
        if (
            not isinstance(evidence["positive"], list)
            or not isinstance(evidence["negative"], list)
            or any(not isinstance(item, str) or not item.strip() for item in evidence["positive"] + evidence["negative"])
            or not isinstance(evidence["reason"], str)
            or not evidence["reason"].strip()
        ):
            raise Refusal(f"witness {target} has malformed {axis} evidence values")
        if not isinstance(evidence["checklist"], dict) or set(evidence["checklist"]) != set(CHECKLISTS[axis]):
            raise Refusal(f"witness {target} has incorrect {axis} checklist categories")
        for category in CHECKLISTS[axis]:
            entry = evidence["checklist"][category]
            if not isinstance(entry, dict) or set(entry) != {"count", "severity"}:
                raise Refusal(f"witness {target} checklist entry is malformed")
            if isinstance(entry["count"], bool) or not isinstance(entry["count"], int) or entry["count"] < 0:
                raise Refusal(f"witness {target} checklist count must be a nonnegative integer")
            severities = aggregates[axis][category]
            expected_severity = max(severities, key=SEVERITY_ORDER.get) if severities else "none"
            if entry["count"] != len(severities) or entry["severity"] != expected_severity:
                raise Refusal(f"witness {target} checklist/card mismatch for {axis}/{category}")
    if (
        not isinstance(witness["unresolved_ambiguity"], list)
        or any(not isinstance(item, str) for item in witness["unresolved_ambiguity"])
    ):
        raise Refusal(f"witness {target} unresolved_ambiguity must be a string list")
    if not isinstance(witness["next_fixes"], list):
        raise Refusal(f"witness {target} next_fixes must be a list")
    for fix in witness["next_fixes"]:
        if not isinstance(fix, dict) or set(fix) != {"axis", "fix"}:
            raise Refusal(f"witness {target} next_fixes entries must be exact axis/fix objects")
        if fix["axis"] not in AXES or not isinstance(fix["fix"], str) or not fix["fix"].strip():
            raise Refusal(f"witness {target} has an invalid next_fixes entry")
    reject_computed_fields(witness)


def validate_invariants(
    document,
    assessment_prompt_sha: str,
    prompt_digests: dict,
    response_digests: dict,
    witnesses: dict,
) -> list:
    required = {
        "schema", "target", "assessment_prompt_sha256", "prompt_sha256",
        "response_sha256", "items",
    }
    if not isinstance(document, dict) or set(document) != required:
        raise Refusal("invariant assessment has incorrect top-level shape")
    if document["schema"] != "cnos-recursive-cell-invariants/0.2" or document["target"] != "cc662-system":
        raise Refusal("invariant assessment schema/target mismatch")
    if document["assessment_prompt_sha256"] != assessment_prompt_sha:
        raise Refusal("invariant assessment does not bind the emitted assessment prompt")
    if document["prompt_sha256"] != prompt_digests:
        raise Refusal("invariant assessment does not bind all six emitted TSC prompts")
    if document["response_sha256"] != response_digests:
        raise Refusal("invariant assessment does not bind all six TSC response witnesses")
    items = document["items"]
    if not isinstance(items, list) or [item.get("id") for item in items if isinstance(item, dict)] != INVARIANTS:
        raise Refusal("invariant assessment must contain H01-H13 exactly once in order")
    for item in items:
        if set(item) != {"id", "status", "evidence", "level", "primary_axis", "next_mca"}:
            raise Refusal(f"invariant {item.get('id')} has incorrect shape")
        if item["status"] not in {"pass", "fail", "unknown"}:
            raise Refusal(f"invariant {item['id']} has invalid status")
        if not isinstance(item["evidence"], str) or not item["evidence"].strip():
            raise Refusal(f"invariant {item['id']} lacks evidence")
        if item["level"] not in {*LEVEL_NAMES, "system"} or item["primary_axis"] not in AXES:
            raise Refusal(f"invariant {item['id']} has invalid level/axis")
        if item["status"] in {"fail", "unknown"} and (not isinstance(item["next_mca"], str) or not item["next_mca"].strip()):
            raise Refusal(f"invariant {item['id']} requires next_mca")
        if item["status"] == "pass" and item["next_mca"] is not None:
            raise Refusal(f"passing invariant {item['id']} must have null next_mca")
        if item["status"] == "fail":
            marker = f"[{item['id']}]"
            target = "cc662-system" if item["level"] == "system" else f"cc662-{item['level'].lower()}"
            cards = witnesses[target]["defect_cards"]
            if not any(
                card["severity"] == "systemic"
                and card["primary_axis"] == item["primary_axis"]
                and marker in (card["summary"] + " " + card["evidence"])
                for card in cards
            ):
                raise Refusal(
                    f"failed invariant {item['id']} lacks a systemic defect card "
                    f"on {target}/{item['primary_axis']}"
                )
    return items


def report_facts(report, target: str):
    if not isinstance(report, dict) or report.get("target") != target:
        raise Refusal(f"TSC report target mismatch for {target}")
    for axis in AXES:
        require_unit_interval(report.get(axis), f"report.{target}.{axis}")
    bottleneck_axis = report.get("bottleneck_axis")
    if bottleneck_axis not in AXES:
        raise Refusal(f"TSC report for {target} has invalid bottleneck_axis")
    try:
        c_math = report["provenance"]["aggregate_math"]["C_sigma_math"]
        c_num = report["provenance"]["aggregate_numeric"]["C_sigma_num"]
    except (KeyError, TypeError) as error:
        raise Refusal(f"TSC report for {target} lacks canonical aggregate provenance") from error
    require_unit_interval(c_math, f"report.{target}.C_sigma_math")
    require_unit_interval(c_num, f"report.{target}.C_sigma_num")
    return {
        "alpha": report["alpha"],
        "beta": report["beta"],
        "gamma": report["gamma"],
        "bottleneck_axis": bottleneck_axis,
        "C_sigma_math": c_math,
        "C_sigma_num": c_num,
    }


def prepare(args):
    root = args.root.resolve()
    cm = root / "src/packages/cnos.cdd/skills/cdd/measure/recursive-cell"
    registry = cm / "calibration/662/registry.tsc"
    instruction = cm / "INSTRUCTION.md"
    skill = cm / "SKILL.md"
    preflight = cm / "calibration/662/verify-target.sh"
    assembler = cm / "instruction/assemble-instruction.sh"
    schema = cm / "runner/recursive-cell-run.schema.cue"
    assessment_template = cm / "runner/invariant-assessment-template.md"
    active_runner = Path(__file__).resolve()
    cm_runner = cm / "runner/recursive-cell-runner.py"
    manifests = [
        cm / "calibration/662/targets/system.tsc",
        cm / "calibration/662/targets/l0-three-cell.tsc",
        cm / "calibration/662/targets/l1-kernel.tsc",
        cm / "calibration/662/targets/l2-contracts.tsc",
        cm / "calibration/662/targets/l3-fsm.tsc",
        cm / "calibration/662/targets/l4-instance.tsc",
    ]
    authority_paths = [
        skill, instruction, cm_runner, schema, assessment_template, registry,
        preflight, assembler,
        *manifests,
    ]
    coh_executable = args.coh.resolve()
    for path in [*authority_paths, active_runner, coh_executable]:
        if not path.is_file():
            raise Refusal(f"required State-A authority file missing: {path}")
    if sha256(active_runner) != sha256(cm_runner):
        raise Refusal("active runner bytes differ from the target-root CM authority runner")
    if not SHA40.fullmatch(args.cm_revision):
        raise Refusal("cm_revision must be a lowercase 40-hex Git object ID")
    if not SHA40.fullmatch(args.engine_revision):
        raise Refusal("engine_revision must be a lowercase 40-hex Git object ID")
    if not RFC3339_UTC.fullmatch(args.timestamp):
        raise Refusal("timestamp must be an RFC3339 UTC timestamp ending in Z")
    try:
        parsed_time = datetime.fromisoformat(args.timestamp.replace("Z", "+00:00"))
    except ValueError as error:
        raise Refusal("timestamp must be an RFC3339 UTC timestamp") from error
    if parsed_time.tzinfo is None or parsed_time.utcoffset() != timezone.utc.utcoffset(parsed_time):
        raise Refusal("timestamp must be an RFC3339 UTC timestamp")
    run([preflight, root])
    output = args.output.resolve()
    output.mkdir(parents=True, exist_ok=True)
    return (
        root, cm, registry, instruction, skill, schema, assessment_template,
        active_runner, authority_paths, coh_executable, output,
    )


def coh_base(args, target, registry, instruction, root):
    return [
        args.coh,
        "--target", target,
        "--registry", registry.relative_to(root),
        "--instruction", instruction.relative_to(root),
        "--root", root,
    ]


def require_unpublished(output: Path) -> None:
    publication = output / "publication"
    if publication.exists() or publication.is_symlink():
        raise Refusal(f"canonical publication already exists and is immutable: {publication}")
    lock = output / ".publication-lock"
    if lock.exists() or lock.is_symlink():
        raise Refusal(f"publication is already in progress or recovery is required: {lock}")
    emit_lock = output / ".emission-lock"
    if emit_lock.exists() or emit_lock.is_symlink():
        raise Refusal(f"prompt emission is already in progress or recovery is required: {emit_lock}")


def acquire_publication_lock(output: Path) -> Path:
    lock = output / ".publication-lock"
    try:
        lock.mkdir()
    except OSError as error:
        raise Refusal(f"could not acquire atomic publication lock: {error}") from error
    publication = output / "publication"
    if publication.exists() or publication.is_symlink():
        lock.rmdir()
        raise Refusal(f"canonical publication already exists and is immutable: {publication}")
    return lock


def acquire_emission_lock(output: Path) -> Path:
    lock = output / ".emission-lock"
    try:
        lock.mkdir()
    except OSError as error:
        raise Refusal(f"could not acquire atomic emission lock: {error}") from error
    unexpected = [path for path in output.iterdir() if path != lock]
    if unexpected:
        lock.rmdir()
        raise Refusal(f"emit requires a fresh run root; found existing path: {unexpected[0]}")
    return lock


def emit(args) -> None:
    root, _, registry, instruction, _, _, assessment_template, _, _, _, output = prepare(args)
    require_unpublished(output)
    lock = acquire_emission_lock(output)
    stage = None
    emitted = False
    try:
        stage = Path(tempfile.mkdtemp(prefix=".emission-stage-", dir=output))
        prompt_dir = stage / "prompts"
        prompt_dir.mkdir(parents=True)
        prompts = []
        for target in TARGETS:
            path = prompt_dir / f"{target}.md"
            run(coh_base(args, target, registry, instruction, root) + ["--mode", "llm", "--emit-prompt", path])
            prompts.append({
                "target": target,
                "path": str(path.relative_to(stage)),
                "sha256": sha256(path),
                "bytes": path.stat().st_size,
            })
        assessment_prompt = stage / "invariant-assessment-prompt.md"
        digest_lines = "\n".join(f"- `{row['target']}`: `{row['sha256']}`" for row in prompts)
        assessment_prompt.write_text(
            assessment_template.read_text(encoding="utf-8").rstrip()
            + "\n\n## Bound TSC prompt digests\n\n"
            + digest_lines
            + "\n",
            encoding="utf-8",
        )
        write_json(stage / "prompt-digests.json", {
            "schema": "cnos-recursive-cell-prompts/0.1",
            "targets": prompts,
            "invariant_assessment": {
                "path": str(assessment_prompt.relative_to(stage)),
                "sha256": sha256(assessment_prompt),
                "bytes": assessment_prompt.stat().st_size,
            },
        })
        emission = output / "emission"
        if emission.exists() or emission.is_symlink():
            raise Refusal(f"canonical emission appeared while the lock was held: {emission}")
        try:
            stage.rename(emission)
        except OSError as error:
            raise Refusal(f"atomic emission failed: {error}") from error
        emitted = True
        print(f"PASS emitted {len(prompts)} recursive-cell prompts to {emission / 'prompts'}")
    finally:
        if not emitted and stage is not None and stage.exists():
            shutil.rmtree(stage)
        lock.rmdir()


def select_report(directory: Path) -> Path:
    candidates = sorted(path for path in directory.glob("*.json") if path.is_file())
    if len(candidates) != 1:
        raise Refusal(f"expected exactly one TSC JSON report in {directory}, found {len(candidates)}")
    return candidates[0]


def disposition_for(items):
    failed = [item["id"] for item in items if item["status"] == "fail"]
    unknown = [item["id"] for item in items if item["status"] == "unknown"]
    if any(identifier in {"H09", "H10", "H11"} for identifier in failed):
        return "request_planning", "hard-invariant-failure"
    if "H12" in failed:
        return "request_working", "hard-invariant-failure"
    if failed or unknown:
        return "hold", "hard-invariant-failure-or-unknown"
    return "hold", "zero-standing-calibration-source"


def ingest(args) -> None:
    (
        root, cm, registry, instruction, skill, schema, assessment_template,
        active_runner, authority_paths, coh_executable, output,
    ) = prepare(args)
    require_unpublished(output)
    emission_source = output / "emission"
    if not emission_source.is_dir() or emission_source.is_symlink():
        raise Refusal("canonical emission missing or malformed; run emit in a fresh root")
    lock = acquire_publication_lock(output)
    stage = None
    published = False
    try:
        stage = Path(tempfile.mkdtemp(prefix=".publication-stage-", dir=output))
        emission = stage / "emission"
        try:
            shutil.copytree(emission_source, emission, symlinks=True)
        except OSError as error:
            raise Refusal(f"could not snapshot canonical emission: {error}") from error
        if any(path.is_symlink() for path in emission.rglob("*")):
            raise Refusal("canonical emission contains a symlink")

        response_paths = {}
        for target in TARGETS:
            response_path = stage / "inputs" / "responses" / f"{target}.json"
            snapshot_file(args.responses / f"{target}.json", response_path, f"response {target}")
            response_paths[target] = response_path
        invariant_path = stage / "inputs" / "invariant-assessment.json"
        snapshot_file(args.invariants, invariant_path, "invariant assessment")

        prompt_manifest_path = emission / "prompt-digests.json"
        prompt_manifest = load_json(prompt_manifest_path)
        if not isinstance(prompt_manifest, dict) or prompt_manifest.get("schema") != "cnos-recursive-cell-prompts/0.1":
            raise Refusal("prompt digest manifest missing or malformed; run emit first")
        prompt_rows = prompt_manifest.get("targets")
        row_keys = {"target", "path", "sha256", "bytes"}
        if (
            not isinstance(prompt_rows, list)
            or any(not isinstance(row, dict) or set(row) != row_keys for row in prompt_rows)
            or [row["target"] for row in prompt_rows] != TARGETS
        ):
            raise Refusal("prompt digest manifest target order mismatch")
        prompt_by_target = {}
        for row in prompt_rows:
            if (
                not isinstance(row["path"], str)
                or not isinstance(row["sha256"], str)
                or not re.fullmatch(r"[0-9a-f]{64}", row["sha256"])
                or isinstance(row["bytes"], bool)
                or not isinstance(row["bytes"], int)
                or row["bytes"] < 1
            ):
                raise Refusal(f"prompt digest row is malformed: {row['target']}")
            path = (emission / row["path"]).resolve()
            if emission not in path.parents:
                raise Refusal(f"prompt path escapes canonical emission: {row['target']}")
            if not path.is_file() or sha256(path) != row["sha256"] or path.stat().st_size != row["bytes"]:
                raise Refusal(f"emitted prompt changed before ingestion: {row['target']}")
            prompt_by_target[row["target"]] = row
        assessment_row = prompt_manifest.get("invariant_assessment")
        if not isinstance(assessment_row, dict) or set(assessment_row) != {"path", "sha256", "bytes"}:
            raise Refusal("prompt digest manifest lacks the invariant-assessment prompt")
        if (
            not isinstance(assessment_row["path"], str)
            or not isinstance(assessment_row["sha256"], str)
            or not re.fullmatch(r"[0-9a-f]{64}", assessment_row["sha256"])
            or isinstance(assessment_row["bytes"], bool)
            or not isinstance(assessment_row["bytes"], int)
            or assessment_row["bytes"] < 1
        ):
            raise Refusal("invariant-assessment prompt digest row is malformed")
        assessment_prompt = (emission / assessment_row["path"]).resolve()
        if emission not in assessment_prompt.parents:
            raise Refusal("invariant-assessment prompt path escapes canonical emission")
        if (
            not assessment_prompt.is_file()
            or sha256(assessment_prompt) != assessment_row["sha256"]
            or assessment_prompt.stat().st_size != assessment_row["bytes"]
        ):
            raise Refusal("emitted invariant-assessment prompt changed before ingestion")

        responses = {}
        response_digests = {}
        for target in TARGETS:
            witness = load_json(response_paths[target])
            validate_witness(witness, target)
            responses[target] = witness
            response_digests[target] = sha256(response_paths[target])
        invariant_document = load_json(invariant_path)
        prompt_digests = {target: prompt_by_target[target]["sha256"] for target in TARGETS}
        invariant_items = validate_invariants(
            invariant_document,
            assessment_row["sha256"],
            prompt_digests,
            response_digests,
            responses,
        )

        report_dir = stage / "reports"
        report_dir.mkdir(parents=True)
        report_rows = []
        facts = {}
        for target in TARGETS:
            scratch = stage / ".coh" / target
            scratch.mkdir(parents=True)
            run(coh_base(args, target, registry, instruction, root) + [
                "--mode", "hybrid", "--llm-response", response_paths[target], "--output", scratch,
            ])
            generated = select_report(scratch)
            canonical = report_dir / f"{target}.json"
            shutil.copyfile(generated, canonical)
            report = load_json(canonical)
            facts[target] = report_facts(report, target)
            report_rows.append({
                "target": target,
                "prompt_sha256": prompt_by_target[target]["sha256"],
                "response_sha256": response_digests[target],
                "report_sha256": sha256(canonical),
                "report_path": str(canonical.relative_to(stage)),
            })
        shutil.rmtree(stage / ".coh")

        level_facts = [facts[target] for target in LEVEL_TARGETS]
        c_math = 0.0 if any(row["C_sigma_math"] == 0 for row in level_facts) else math.prod(row["C_sigma_math"] for row in level_facts) ** (1.0 / 5.0)
        c_num = math.prod(row["C_sigma_num"] for row in level_facts) ** (1.0 / 5.0)
        level_index = min(range(5), key=lambda index: (level_facts[index]["C_sigma_num"], index))
        axis = level_facts[level_index]["bottleneck_axis"]
        gate_passed = all(item["status"] == "pass" for item in invariant_items)
        disposition, reason = disposition_for(invariant_items)
        engine_version = run([args.coh, "--version"], capture=True)
        result = {
            "schema": "cnos-recursive-cell-run/0.2",
            "target": "cc662-r7-recursive",
            "status": "calibration-source-zero-standing",
            "targets": report_rows,
            "cross_level": {
                "levels": [
                    {"level": LEVEL_NAMES[index], "target": LEVEL_TARGETS[index], **level_facts[index]}
                    for index in range(5)
                ],
                "C_sigma_cross_math": c_math,
                "C_sigma_cross_num": c_num,
                "semantics": "unweighted-geometric-mean-of-L0-L4-canonical-TSC-C_sigma",
            },
            "hard_invariant_gate": {"passed": gate_passed, "items": invariant_items},
            "bottleneck": {
                "level": LEVEL_NAMES[level_index],
                "axis": axis,
                "tie_break": "lowest-level-index-then-TSC-bottleneck-axis",
            },
            "disposition": {"value": disposition, "reason": reason, "standing": "none"},
            "provenance": {
                "declared_cm_revision": args.cm_revision,
                "target_revision": run(["git", "-C", root, "rev-parse", "HEAD"], capture=True),
                "cm_authority_bundle_sha256": path_framed_sha256(cm, authority_paths),
                "skill_sha256": sha256(skill),
                "registry_sha256": sha256(registry),
                "instruction_sha256": sha256(instruction),
                "runner_sha256": sha256(active_runner),
                "output_schema_sha256": sha256(schema),
                "invariant_assessment_template_sha256": sha256(assessment_template),
                "invariant_assessment_prompt_sha256": assessment_row["sha256"],
                "invariant_assessment_sha256": sha256(invariant_path),
                "response_sha256": response_digests,
                "declared_engine_revision": args.engine_revision,
                "coh_executable_sha256": sha256(coh_executable),
                "engine_version": engine_version,
                "run_time_utc": args.timestamp,
                "mode": "State-A external-response hybrid ingestion",
            },
        }
        result_path = stage / "recursive-cell-run.json"
        write_json(result_path, result)
        run(["cue", "vet", "-d", "#RecursiveCellRun", schema, result_path])

        success = {
            "schema": "cnos-recursive-cell-publication/0.2",
            "target": "cc662-r7-recursive",
            "status": "complete",
            "result": {"path": "recursive-cell-run.json", "sha256": sha256(result_path)},
            "reports": [
                {"target": row["target"], "path": row["report_path"], "sha256": row["report_sha256"]}
                for row in report_rows
            ],
            "emission": {
                "prompt_manifest": {
                    "path": str(prompt_manifest_path.relative_to(stage)),
                    "sha256": sha256(prompt_manifest_path),
                },
                "invariant_assessment_prompt": {
                    "path": str(assessment_prompt.relative_to(stage)),
                    "sha256": assessment_row["sha256"],
                },
                "prompts": [
                    {
                        "target": row["target"],
                        "path": str((emission / row["path"]).relative_to(stage)),
                        "sha256": row["sha256"],
                    }
                    for row in prompt_rows
                ],
            },
            "inputs": {
                "responses": [
                    {
                        "target": target,
                        "path": str(response_paths[target].relative_to(stage)),
                        "sha256": response_digests[target],
                    }
                    for target in TARGETS
                ],
                "invariant_assessment": {
                    "path": str(invariant_path.relative_to(stage)),
                    "sha256": sha256(invariant_path),
                },
            },
        }
        success_path = stage / "publication-success.json"
        write_json(success_path, success)
        run(["cue", "vet", "-d", "#PublicationSuccess", schema, success_path])

        publication = output / "publication"
        if publication.exists() or publication.is_symlink():
            raise Refusal(f"canonical publication appeared while the lock was held: {publication}")
        try:
            stage.rename(publication)
        except OSError as error:
            raise Refusal(f"atomic publication failed: {error}") from error
        published = True
        print(
            f"PASS recursive-cell publication: {publication} "
            f"disposition={disposition} bottleneck={LEVEL_NAMES[level_index]}/{axis}"
        )
    finally:
        if not published and stage is not None and stage.exists():
            shutil.rmtree(stage)
        lock.rmdir()


def parser() -> argparse.ArgumentParser:
    result = argparse.ArgumentParser(description=__doc__)
    subparsers = result.add_subparsers(dest="command", required=True)
    for command in ["emit", "ingest"]:
        child = subparsers.add_parser(command)
        child.add_argument("--coh", type=Path, required=True)
        child.add_argument("--root", type=Path, required=True)
        child.add_argument("--output", type=Path, required=True)
        child.add_argument("--engine-revision", required=True)
        child.add_argument("--cm-revision", required=True)
        child.add_argument("--timestamp", required=True, help="caller-supplied RFC3339 UTC provenance value")
        if command == "ingest":
            child.add_argument("--responses", type=Path, required=True)
            child.add_argument("--invariants", type=Path, required=True)
    return result


def main() -> int:
    args = parser().parse_args()
    try:
        emit(args) if args.command == "emit" else ingest(args)
        return 0
    except Refusal as error:
        print(f"REFUSE: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
