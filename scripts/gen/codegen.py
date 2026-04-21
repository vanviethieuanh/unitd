"""Go code generation for systemd unit config structs."""

from __future__ import annotations

from pathlib import Path

import jinja2

from .directive import Directive, generate_block_struct
from .parser import KnownUnit, Unit
from .util import to_pascal_case, to_snake_case, wrap_comment

_TEMPLATES_DIR = Path(__file__).parent / "templates"

_JINJA_ENV = jinja2.Environment(
    loader=jinja2.FileSystemLoader(str(_TEMPLATES_DIR)),
    keep_trailing_newline=True,
    lstrip_blocks=True,
    trim_blocks=True,
)

_COMMON_TEMPLATE = _JINJA_ENV.get_template("common.go.j2")
_UNIT_TEMPLATE = _JINJA_ENV.get_template("unit.go.j2")
_KNOWN_UNITS_TEMPLATE = _JINJA_ENV.get_template("known_units.go.j2")


def gen_common(pkg: str, directives: list[Directive]) -> str:
    """Generate the shared Unit + Install blocks (systemd.unit)."""
    install_code, install_deps = generate_block_struct(directives, "Install", "core")
    unit_code, unit_deps = generate_block_struct(directives, "Unit", "core")

    imports = sorted(install_deps | unit_deps)

    return _COMMON_TEMPLATE.render(
        pkg=pkg,
        imports=imports,
        install_block=install_code,
        unit_block=unit_code,
    )


def generate_unit_code(u: Unit) -> tuple[str, set[str]]:
    """Generate Go code for a unit type.

    Returns (code, import_set).
    """
    import_set: set[str] = set()
    parts: list[str] = []

    doc = _format_unit_doc(u)
    if doc:
        parts.append(doc.rstrip("\n"))

    block_type = to_pascal_case(u.name)
    sub_block_type = block_type + "Block"
    has_block_type = len(u.options) > 0

    if has_block_type:
        code, deps = generate_block_struct(u.options, u.name, "core")
        import_set |= deps
        parts.append(code)

    snake = to_snake_case(block_type) if has_block_type else ""
    struct_lines = [f"type {block_type} struct {{"]
    struct_lines.append(f'\tName string `hcl:"name,label"`')
    struct_lines.append("")
    struct_lines.append(f'\tTemplate bool              `hcl:"template,optional"`')
    struct_lines.append(f'\tForEach  map[string]string `hcl:"for_each,optional"`')
    struct_lines.append("")
    struct_lines.append(f'\tUnit    UnitBlock    `hcl:"unit,block"`')
    if has_block_type:
        struct_lines.append(f'\t{block_type} {sub_block_type} `hcl:"{snake},block"`')
    struct_lines.append(f'\tInstall InstallBlock `hcl:"install,block"`')
    struct_lines.append("}")
    parts.append("\n".join(struct_lines) + "\n")

    return "\n".join(parts), import_set


def _format_unit_doc(u: Unit) -> str:
    if not u.description:
        return ""

    name = to_pascal_case(u.name)
    out = ""

    if u.options:
        out += f"// {name}Block is for [{name}] systemd unit block\n"
        out += "//\n"
    else:
        out += f"// {name} is for {name} systemd unit file\n"

    for line in wrap_comment(u.description, 100):
        out += f"// {line}\n"

    return out


def format_imports(packages: list[str]) -> str:
    """Format an import block. Public for use by __main__.py."""
    return _format_imports(packages)


def _format_imports(packages: list[str]) -> str:
    if not packages:
        return ""
    out = "import (\n"
    for p in packages:
        out += f'\t"{p}"\n'
    out += ")\n"
    return out


def render_unit_file(pkg: str, man: str, imports: list[str], unit_code: str) -> str:
    """Render a complete unit Go file from pre-generated unit code."""
    return _UNIT_TEMPLATE.render(
        pkg=pkg,
        man=man,
        imports=imports,
        unit_code=unit_code,
    )


_KNOWN_UNITS_TEMPLATE = _JINJA_ENV.get_template("known_units.go.j2")


def generate_known_units_code(
    pkg: str,
    special_units: list[KnownUnit],
    shipped_units: list[KnownUnit],
) -> str:
    """Generate Go code for the known units registry."""
    # Merge and deduplicate: special takes precedence over shipped
    seen: set[str] = set()
    all_units: list[KnownUnit] = []
    for u in special_units:
        if u.name not in seen:
            seen.add(u.name)
            all_units.append(u)
    for u in shipped_units:
        if u.name not in seen:
            seen.add(u.name)
            all_units.append(u)

    all_units.sort(key=lambda u: u.name)

    return _KNOWN_UNITS_TEMPLATE.render(pkg=pkg, units=all_units)
