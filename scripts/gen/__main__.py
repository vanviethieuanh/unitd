#!/usr/bin/env python3
"""Generate Go config structs from systemd man pages and gperf data.

Usage:
    python -m scripts.gen --man systemd.timer --systemd-src tmp/systemd --package configs
    python -m scripts.gen --man systemd.unit  --systemd-src tmp/systemd --package configs
"""

from __future__ import annotations

import argparse
import logging
import subprocess
import sys
from pathlib import Path

from .codegen import gen_common, generate_known_units_code, generate_unit_code, format_imports
from .directive import Directive
from .gperf import extract_all_gperfs, load_all_directives, load_parser_map
from .parser import parse_applicable_types, parse_descriptions, parse_special_units, parse_unit, scan_shipped_units

log = logging.getLogger(__name__)


def main(argv: list[str] | None = None) -> None:
    p = argparse.ArgumentParser(description="Generate Go config structs from systemd data")
    p.add_argument("--man", default=None, help="systemd man page name (e.g. systemd.timer)")
    p.add_argument("--known-units", action="store_true", help="generate known units registry instead of config structs")
    p.add_argument(
        "--shared-man",
        action="append",
        default=[],
        help="shared man pages to check for applicability (repeatable)",
    )
    p.add_argument("--systemd-src", default="tmp/systemd", help="path to systemd source tree")
    p.add_argument("--man-dir", default=None, help="directory containing XML man pages (default: <systemd-src>/man)")
    p.add_argument("--gperf-dir", default=None, help="use pre-extracted gperf JSONL files from this directory instead of extracting from source")
    p.add_argument("--package", default="configs", help="Go package name")
    p.add_argument("--parser-types", default=None, help="path to parser_type.jsonl (default: generate in-memory)")
    p.add_argument("--debug", action="store_true", help="write JSONL intermediates to tmp/gperf/")
    p.add_argument("--log-level", default="info", help="log level")
    args = p.parse_args(argv)

    logging.basicConfig(
        level=getattr(logging, args.log_level.upper(), logging.INFO),
        format="%(levelname)s %(message)s",
    )

    systemd_src = Path(args.systemd_src)
    man_dir = Path(args.man_dir) if args.man_dir else systemd_src / "man"

    if args.known_units:
        code = _generate_known_units(man_dir, systemd_src / "units", args.package)
        print(code, end="")
        return

    if not args.man:
        p.error("--man is required when not using --known-units")

    debug_dir = Path("tmp/gperf") if args.debug else None

    python = sys.executable
    if args.gperf_dir:
        gperf_records = _load_gperf_from_jsonl(Path(args.gperf_dir))
    else:
        gperf_records = extract_all_gperfs(systemd_src, python=python, debug_dir=debug_dir)
    log.info("extracted %d gperf records", len(gperf_records))

    if args.parser_types:
        parser_map = load_parser_map(Path(args.parser_types))
    else:
        parser_map = _run_extract_parser_types(systemd_src, gperf_records, debug_dir)
    log.info("loaded %d parser type mappings", len(parser_map))

    directives = load_all_directives(gperf_records, parser_map)
    log.info("built %d directives", len(directives))

    if args.man == "systemd.unit":
        print(gen_common(args.package, directives), end="")
    else:
        code = _generate_unit(args.man, man_dir, args.package, args.shared_man, directives)
        print(code, end="")


def _load_gperf_from_jsonl(gperf_dir: Path) -> list:
    """Load pre-extracted gperf records from JSONL files."""
    import json
    from .gperf import GperfRecord

    records: list[GperfRecord] = []
    for path in sorted(gperf_dir.glob("gperf_*.jsonl")):
        for line in path.read_text().splitlines():
            if not line.strip():
                continue
            rec = json.loads(line)
            records.append(
                GperfRecord(
                    system=rec["system"],
                    section=rec["section"],
                    property=rec["property"],
                    parser=rec["parser"],
                    ltype=rec.get("ltype", ""),
                    offset=rec.get("offset", ""),
                )
            )
    return records


def _run_extract_parser_types(
    systemd_src: Path,
    gperf_records: list,
    debug_dir: Path | None,
) -> dict[str, str]:
    """Run extract_parser_types.py in-process to get parser→type map.

    If debug_dir is set, also writes the JSONL file.
    """
    import json
    import tempfile

    scripts_dir = Path(__file__).resolve().parent.parent
    extract_script = scripts_dir / "extract_parser_types.py"

    if debug_dir:
        output_path = debug_dir / "parser_type.jsonl"
        gperf_dir = debug_dir
        debug_dir.mkdir(parents=True, exist_ok=True)
        cleanup = False
    else:
        tmpdir = Path(tempfile.mkdtemp())
        output_path = tmpdir / "parser_type.jsonl"
        gperf_dir = tmpdir
        cleanup = True

    gperf_jsonl_path = gperf_dir / "gperf_all.jsonl"
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

    cmd = [sys.executable, str(extract_script), str(systemd_src), str(output_path), str(gperf_dir)]
    subprocess.run(cmd, check=True, capture_output=True, text=True)

    parser_map = load_parser_map(output_path)

    if cleanup:
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

    code = f"// Generated based on man page {man} of systemd\n\n"
    code += f"package {pkg}\n\n"

    unit_code, imports = generate_unit_code(unit)
    code += format_imports(sorted(imports))
    code += unit_code

    return code


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
