"""Load gperf data and parser type mappings."""

from __future__ import annotations

import json
import logging
import re
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path

from .directive import Directive, new_directive

log = logging.getLogger(__name__)

def load_parser_map(path: Path) -> dict[str, str]:
    """Load parser→type mappings from a JSONL file."""
    result: dict[str, str] = {}
    for line in path.read_text().splitlines():
        if not line.strip():
            continue
        rec = json.loads(line)
        result[rec["parser"]] = rec["type"]
    return result


_GPERF_LINE_RE = re.compile(r"^[A-Za-z0-9_.]+,")

_EXPAND_JINJA2 = Path(__file__).resolve().parent.parent / "expand_jinja2.py"


@dataclass
class GperfRecord:
    system: str
    section: str
    property: str
    parser: str
    ltype: str
    offset: str


def _extract_gperf_file(
    item: Path,
    python: str,
) -> tuple[str, list[GperfRecord]]:
    """Extract gperf records from a single .gperf or .gperf.in file."""
    stem = item.name
    stem = re.sub(r"-gperf\.gperf(\.in)?$", "", stem)
    source_name = stem
    system = "core" if stem == "load-fragment" else stem

    if item.suffix == ".in":
        result = subprocess.run(
            [python, str(_EXPAND_JINJA2), str(item)],
            capture_output=True,
            text=True,
            check=True,
        )
        content = result.stdout
    else:
        content = item.read_text()

    records: list[GperfRecord] = []
    for line in content.splitlines():
        if not _GPERF_LINE_RE.match(line):
            continue

        parts = line.split(",")
        if len(parts) < 4:
            continue

        col1 = parts[0].strip()
        col2 = parts[1].strip()
        col3 = parts[2].strip()
        col4 = ",".join(parts[3:]).strip()

        dot_parts = col1.split(".", 1)
        if len(dot_parts) != 2:
            continue
        section, key = dot_parts

        records.append(
            GperfRecord(
                system=system,
                section=section,
                property=key,
                parser=col2,
                ltype=col3,
                offset=col4,
            )
        )

    return source_name, records


def extract_all_gperfs(
    systemd_src: Path,
    python: str = sys.executable,
    debug_dir: Path | None = None,
) -> list[GperfRecord]:
    """Find and extract all gperf files from the systemd source tree.

    If debug_dir is set, write JSONL files there.
    """
    gperf_files = sorted(
        list(systemd_src.glob("src/**/*.gperf"))
        + list(systemd_src.glob("src/**/*.gperf.in"))
    )

    all_records: list[GperfRecord] = []
    by_source: dict[str, list[GperfRecord]] = {}

    for item in gperf_files:
        source_name, records = _extract_gperf_file(item, python)
        all_records.extend(records)
        by_source.setdefault(source_name, []).extend(records)

    if debug_dir is not None:
        debug_dir.mkdir(parents=True, exist_ok=True)
        for source_name, records in sorted(by_source.items()):
            out_path = debug_dir / f"gperf_{source_name}.jsonl"
            with open(out_path, "w") as f:
                for r in records:
                    json.dump(
                        {
                            "system": r.system,
                            "section": r.section,
                            "property": r.property,
                            "parser": r.parser,
                            "ltype": r.ltype,
                            "offset": r.offset,
                        },
                        f,
                    )
                    f.write("\n")
            log.debug("wrote %s (%d records)", out_path, len(records))

    return all_records


_NULL_PARSER_TYPES: dict[str, str] = {
    "Alias": "STRING",
    "WantedBy": "UNIT [...]",
    "RequiredBy": "UNIT [...]",
    "UpheldBy": "UNIT [...]",
    "Also": "UNIT [...]",
    "DefaultInstance": "STRING",
}


def load_all_directives(
    gperf_records: list[GperfRecord],
    parser_map: dict[str, str],
) -> list[Directive]:
    """Convert gperf records to Directives using the parser map."""
    result: list[Directive] = []

    for gr in gperf_records:
        type_str = ""
        found = False

        if gr.parser == "NULL":
            if gr.property in _NULL_PARSER_TYPES:
                type_str = _NULL_PARSER_TYPES[gr.property]
                found = True
        if not found:
            if gr.parser in parser_map:
                type_str = parser_map[gr.parser]
                found = True

        if not found:
            log.warning("unknown parser %r for %s.%s — skipping", gr.parser, gr.section, gr.property)
            continue

        type_str = type_str.strip()
        if not type_str or type_str == "NOTSUPPORTED":
            continue

        try:
            d = new_directive(gr.section, gr.property, gr.system, type_str)
        except ValueError as e:
            log.warning("parse type %r (parser %s): %s — skipping", type_str, gr.parser, e)
            continue

        result.append(d)

    return result
