"""Directive model and Go struct field code generation."""

from __future__ import annotations

from dataclasses import dataclass, field

from .types import parse_type_expr
from .util import to_pascal_case, to_snake_case, wrap_comment


@dataclass
class DirectiveIdentifier:
    section: str
    key: str


@dataclass
class Directive:
    identifier: DirectiveIdentifier
    type: str  # Go type string
    system: str
    description: str = ""
    deps: list[str] = field(default_factory=list)
    native_type: bool = True


def new_directive(section: str, key: str, system: str, type_str: str) -> Directive:
    """Create a Directive by parsing a type expression string."""
    parsed = parse_type_expr(type_str)
    go_type, deps = parsed.to_go_type()

    return Directive(
        identifier=DirectiveIdentifier(section=section, key=key),
        type=go_type,
        system=system,
        deps=deps,
        native_type=len(deps) == 0,
    )


def write_directive(d: Directive) -> str:
    """Generate the Go struct field code for a single directive."""
    lines: list[str] = []

    if d.description:
        paragraphs = d.description.split("\n\n")
        for paragraph in paragraphs:
            paragraph = paragraph.strip()
            if paragraph:
                wrapped = wrap_comment(paragraph, 100)
                for line in wrapped:
                    lines.append(f"\t// {line}")
                if len(paragraphs) > 1:
                    lines.append("\t//")

    field_name = to_pascal_case(d.identifier.key)
    snake_name = to_snake_case(d.identifier.key)
    systemd_name = d.identifier.key

    if d.native_type:
        lines.append(
            f'\t{field_name} {d.type} `hcl:"{snake_name},optional" systemd:"{systemd_name}"`'
        )
    else:
        lines.append(
            f'\t{field_name} {d.type} `unitd:"{snake_name},optional" systemd:"{systemd_name}"`'
        )

    return "\n".join(lines) + "\n"


def generate_block_struct(
    directives: list[Directive],
    section: str,
    system: str,
) -> tuple[str, set[str]]:
    """Generate a Go struct for a config section block.

    Returns (code_string, import_set).
    """
    imports: set[str] = set()

    type_name = to_pascal_case(section) + "Block"

    filtered = [
        d
        for d in directives
        if d.identifier.section.lower() == section.lower()
        and d.system.lower() == system.lower()
    ]
    filtered.sort(key=lambda d: to_pascal_case(d.identifier.key))

    out = f"type {type_name} struct {{\n"
    for d in filtered:
        out += write_directive(d)
        imports.update(d.deps)
    out += "}\n"

    return out, imports
