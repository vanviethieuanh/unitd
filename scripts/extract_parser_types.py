#!/usr/bin/env python3
"""Extract parser→type mappings from systemd C source code.

Replaces the manual extract_dump_config.sh + supplement_parser_types.sh pipeline
with automated extraction from three sources:

1. The unit_dump_config_items() table in load-fragment.c   (explicit mappings)
2. DEFINE_CONFIG_PARSE* macro invocations                  (enum/ptr/string macros)
3. Function-name heuristics for remaining parsers          (fallback patterns)

Usage:
    python3 extract_parser_types.py <systemd_src_dir> <output.jsonl> [<gperf_dir>]

If <gperf_dir> is provided, the script also checks which parsers are actually
used in gperf JSONL files and warns about any that remain unmapped.
"""

import json
import os
import re
import sys
from pathlib import Path


def extract_dump_table(load_fragment_c: str) -> dict[str, str]:
    """Extract {parser, "TYPE"} pairs from the unit_dump_config_items() table."""
    result = {}
    for m in re.finditer(
        r'\{\s*(config_parse_\w+)\s*,\s*"([^"]*)"\s*\}', load_fragment_c
    ):
        result[m.group(1)] = m.group(2)
    return result


# Normalize systemd's specialized dump-table type names into the canonical
# types recognized by our Go code generator (see cmd/gen/type.go).
# Types not listed here pass through unchanged.
_TYPE_NORMALIZE: dict[str, str] = {
    # List types → STRING [...]
    "ARCHS": "STRING [...]",
    "BOUNDINGSET": "STRING [...]",
    "BPF_DELEGATE_ATTACHMENTS": "STRING [...]",
    "BPF_DELEGATE_COMMANDS": "STRING [...]",
    "BPF_DELEGATE_MAPS": "STRING [...]",
    "BPF_DELEGATE_PROGRAMS": "STRING [...]",
    "DEVICE": "STRING [...]",
    "DEVICELATENCY": "STRING [...]",
    "DEVICEWEIGHT": "STRING [...]",
    "ENVIRON": "STRING [...]",
    "FAMILIES": "STRING [...]",
    "FILESYSTEMS": "STRING [...]",
    "LIMIT": "STRING",
    "NAMESPACES": "STRING [...]",
    "REGEX": "STRING [...]",
    "SYSCALLS": "STRING [...]",
    # Scalar types → STRING
    "CPUAFFINITY": "STRING",
    "CPUSCHEDPOLICY": "STRING",
    "CPUSCHEDPRIO": "STRING",
    "ERRNO": "STRING",
    "FACILITY": "STRING",
    "FILE": "STRING",
    "INPUT": "STRING",
    "KILLMODE": "STRING",
    "LABEL": "STRING",
    "MOUNTFLAG": "STRING",
    "NICE": "STRING",
    "OOMSCOREADJUST": "STRING",
    "OUTPUT": "STRING",
    "PERSONALITY": "STRING",
    "POLICY": "STRING",
    "SECUREBITS": "STRING",
    "SLICE": "STRING",
    # Numeric types
    "CPUWEIGHT": "UNSIGNED",
    "NANOSECONDS": "SECONDS",
    "WEIGHT": "UNSIGNED",
    # Complex compound types → simplified
    "PATH [ARGUMENT [...]]": "STRING",
    "PATH[:PATH[:OPTIONS]] [...]": "STRING [...]",
}


def normalize_type(type_str: str) -> str:
    """Normalize a systemd dump-table type into a canonical type."""
    return _TYPE_NORMALIZE.get(type_str, type_str)

def extract_macros(load_fragment_c: str, conf_parser_c: str) -> dict[str, str]:
    """Extract parser types from DEFINE_CONFIG_PARSE* macro invocations.

    Only ENUM and plain DEFINE_CONFIG_PARSE are used here. DEFINE_CONFIG_PARSE_PTR
    is intentionally skipped because its C storage type (e.g. uint64_t for bitmasks)
    does not reflect the config-file value format.
    """
    result = {}
    combined = load_fragment_c + "\n" + conf_parser_c

    # DEFINE_CONFIG_PARSE_ENUM      (name, ..., EnumType)     → STRING
    # DEFINE_CONFIG_PARSE_ENUM_FULL (name, ..., EnumType)     → STRING
    # DEFINE_CONFIG_PARSE_ENUM_WITH_DEFAULT(name, ..., ctype, default) → STRING
    for m in re.finditer(
        r"DEFINE_CONFIG_PARSE_ENUM\w*\(\s*(config_parse_\w+)", combined
    ):
        result[m.group(1)] = "STRING"

    # DEFINE_CONFIG_PARSE(name, conv_func)  — always parses from string
    for m in re.finditer(
        r"DEFINE_CONFIG_PARSE\(\s*(config_parse_\w+)\s*,", combined
    ):
        result.setdefault(m.group(1), "STRING")

    return result


# Ordered patterns: first match wins.
_NAME_PATTERNS: list[tuple[re.Pattern, str]] = [
    (re.compile(r"_sec_|_timeout_|_duration_"), "SECONDS"),
    (re.compile(r"_path_strv|_paths$"), "PATH [...]"),
    (re.compile(r"_path$|_pid_file|_working_directory"), "PATH"),
    (re.compile(r"_strv|_environ|_families|_filter_patterns|_filesystems|_interfaces$"
                r"|_images$|_options$|_fields|_credential$"
                r"|_bind_paths|_temporary_filesystems|_device_allow"
                r"|_io_device_|_io_limit|_socket_bind|_nft_set"
                r"|_bpf_progs|_open_file|_exit_status"
                r"|_directories$|_prefixes$|_foreign_program$"
                r"|_disable_controllers$|_bind_network_interface$"
                r"|_syscall_filter$|_syscall_log$"), "STRING [...]"),
    (re.compile(r"_bool$|_tristate$"), "BOOLEAN"),
    (re.compile(r"_tty_size$"), "UNSIGNED"),
    (re.compile(r"_priority$"), "INTEGER"),
    (re.compile(r"_obsolete_"), "NOTSUPPORTED"),
]


def infer_from_name(parser: str) -> str | None:
    """Infer type from parser function name patterns."""
    for pattern, type_str in _NAME_PATTERNS:
        if pattern.search(parser):
            return type_str
    return None


def collect_used_parsers(gperf_dir: str) -> set[str]:
    """Collect all parser names referenced in gperf JSONL files."""
    parsers = set()
    for path in Path(gperf_dir).glob("gperf_*.jsonl"):
        for line in path.read_text().splitlines():
            rec = json.loads(line)
            parsers.add(rec["parser"])
    return parsers


def main():
    if len(sys.argv) < 3:
        print(f"Usage: {sys.argv[0]} <systemd_src_dir> <output.jsonl> [<gperf_dir>]",
              file=sys.stderr)
        sys.exit(1)

    src_dir = sys.argv[1]
    output_path = sys.argv[2]
    gperf_dir = sys.argv[3] if len(sys.argv) > 3 else None

    load_fragment_path = os.path.join(src_dir, "src", "core", "load-fragment.c")
    conf_parser_path = os.path.join(src_dir, "src", "shared", "conf-parser.c")

    load_fragment_c = open(load_fragment_path).read()
    conf_parser_c = open(conf_parser_path).read()

    # Phase 1: explicit table
    mappings = extract_dump_table(load_fragment_c)

    # Phase 2: macro definitions (only for parsers not already in the dump table)
    macro_mappings = extract_macros(load_fragment_c, conf_parser_c)
    for parser, type_str in macro_mappings.items():
        mappings.setdefault(parser, type_str)

    # Phase 3: name-pattern overrides for parsers whose dump-table type
    # is known to be ambiguous (e.g. LIMIT used for both single-value memory
    # settings and multi-value IO device settings). These take precedence.
    for parser in list(mappings):
        override = infer_from_name(parser)
        if override is not None:
            mappings[parser] = override

    # Phase 4: name-based heuristics for parsers used in gperf but still unmapped
    if gperf_dir:
        used = collect_used_parsers(gperf_dir)
        used.discard("NULL")
        unmapped = used - set(mappings)

        for parser in sorted(unmapped):
            inferred = infer_from_name(parser)
            if inferred:
                mappings[parser] = inferred
            else:
                # Default: most custom parsers accept a string value
                mappings[parser] = "STRING"
                print(f"WARN: no pattern for {parser}, defaulting to STRING",
                      file=sys.stderr)

    # Normalize all types to canonical forms recognized by the Go code generator.
    mappings = {p: normalize_type(t) for p, t in mappings.items()}

    # Always include NULL → "" (handled specially by the Go code)
    mappings["NULL"] = ""

    # Write output
    with open(output_path, "w") as f:
        for parser in sorted(mappings):
            json.dump({"parser": parser, "type": mappings[parser]}, f)
            f.write("\n")

    print(f"Wrote {len(mappings)} parser mappings to {output_path}", file=sys.stderr)


if __name__ == "__main__":
    main()
