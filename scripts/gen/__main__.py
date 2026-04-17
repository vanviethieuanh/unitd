#!/usr/bin/env python3
"""Generate Go config structs from systemd man pages and gperf data.

Configuration is loaded from scripts/gen/gen-config.yaml.

Usage:
    python -m scripts.gen
"""

from __future__ import annotations

import logging
import subprocess
import sys
from pathlib import Path

import yaml

from .codegen import gen_common, generate_known_units_code, generate_unit_code, render_unit_file
from .directive import Directive
from .gperf import extract_all_gperfs, load_all_directives, load_parser_map
from .parser import parse_applicable_types, parse_descriptions, parse_special_units, parse_unit, scan_shipped_units

log = logging.getLogger(__name__)

_CONFIG_PATH = Path(__file__).parent / "gen-config.yaml"


def main() -> None:
    logging.basicConfig(level=logging.INFO, format="%(levelname)s %(message)s")

    cfg = yaml.safe_load(_CONFIG_PATH.read_text())

    systemd_src = Path(cfg["systemd_src"])
    man_dir = systemd_src / "man"
    pkg = cfg["package"]
    output_dir = Path(cfg["output_dir"])
    man_files: list[str] = cfg["man_files"]
    shared_mans: list[str] = cfg.get("shared_man_files", [])

    output_dir.mkdir(parents=True, exist_ok=True)

    # --- Load directives (shared across all man pages) ---
    gperf_records = extract_all_gperfs(systemd_src, python=sys.executable)
    log.info("extracted %d gperf records", len(gperf_records))

    parser_map = _run_extract_parser_types(systemd_src, gperf_records)
    log.info("loaded %d parser type mappings", len(parser_map))

    directives = load_all_directives(gperf_records, parser_map)
    log.info("built %d directives", len(directives))

    # --- Generate unit config files ---
    for man in man_files:
        if man == "systemd.unit":
            code = gen_common(pkg, directives)
        else:
            code = _generate_unit(man, man_dir, pkg, shared_mans, directives)

        out_path = output_dir / f"{man}.go"
        _write_go_file(out_path, code)
        log.info("wrote %s", out_path)

    # --- Generate known units registry ---
    code = _generate_known_units(man_dir, systemd_src / "units", pkg)
    out_path = output_dir / "known_units.go"
    _write_go_file(out_path, code)
    log.info("wrote %s", out_path)


def _write_go_file(path: Path, code: str) -> None:
    """Write code to path, formatted with gofmt."""
    result = subprocess.run(
        ["gofmt"],
        input=code,
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        log.error("gofmt failed for %s: %s", path, result.stderr)
        path.write_text(code)
        return
    path.write_text(result.stdout)


def _run_extract_parser_types(
    systemd_src: Path,
    gperf_records: list,
) -> dict[str, str]:
    """Run extract_parser_types.py to get parser→type map."""
    import json
    import tempfile

    scripts_dir = Path(__file__).resolve().parent.parent
    extract_script = scripts_dir / "extract_parser_types.py"

    tmpdir = Path(tempfile.mkdtemp())
    output_path = tmpdir / "parser_type.jsonl"

    gperf_jsonl_path = tmpdir / "gperf_all.jsonl"
    with open(gperf_jsonl_path, "w") as f:
        for r in gperf_records:
            json.dump({
                "system": r.system,
                "section": r.section,
                "property": r.property,
                "parser": r.parser,
                "ltype": r.ltype,
                "offset": r.offset,
            }, f)
            f.write("\n")

    cmd = [sys.executable, str(extract_script), str(systemd_src), str(output_path), str(tmpdir)]
    subprocess.run(cmd, check=True, capture_output=True, text=True)

    parser_map = load_parser_map(output_path)

    import shutil
    shutil.rmtree(tmpdir, ignore_errors=True)

    return parser_map


def _generate_unit(
    man: str,
    man_dir: Path,
    pkg: str,
    shared_mans: list[str],
    directives: list[Directive],
) -> str:
    """Generate code for a single unit type."""
    xml_path = man_dir / f"{man}.xml"
    unit_type = man.removeprefix("systemd.")

    extra_descriptions: dict[str, str] = {}
    for extra in shared_mans:
        extra_path = man_dir / f"{extra}.xml"
        if not extra_path.exists():
            log.warning("shared man %s not found at %s", extra, extra_path)
            continue

        types = parse_applicable_types(extra_path)
        if unit_type not in types:
            continue

        log.debug("shared man page applies: %s -> %s", extra, unit_type)
        descs = parse_descriptions(extra_path)
        for k, v in descs.items():
            if k not in extra_descriptions:
                extra_descriptions[k] = v

    unit = parse_unit(xml_path, directives, extra_descriptions)

    unit_code, imports = generate_unit_code(unit)

    return render_unit_file(pkg, man, sorted(imports), unit_code)


def _generate_known_units(man_dir: Path, units_dir: Path, pkg: str) -> str:
    """Generate the known units registry Go file."""
    special_xml = man_dir / "systemd.special.xml"
    special_units = parse_special_units(special_xml)
    log.info("parsed %d special units from systemd.special.xml", len(special_units))

    shipped_units = scan_shipped_units(units_dir)
    log.info("scanned %d shipped units from %s", len(shipped_units), units_dir)

    return generate_known_units_code(pkg, special_units, shipped_units)


if __name__ == "__main__":
    main()
