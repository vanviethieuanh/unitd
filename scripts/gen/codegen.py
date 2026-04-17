"""Go code generation for systemd unit config structs."""

from __future__ import annotations

from .directive import Directive, generate_block_struct
from .parser import Unit
from .util import to_pascal_case, to_snake_case, wrap_comment


def gen_common(pkg: str, directives: list[Directive]) -> str:
    """Generate the shared Unit + Install blocks (systemd.unit)."""
    install_code, install_deps = generate_block_struct(directives, "Install", "core")
    unit_code, unit_deps = generate_block_struct(directives, "Unit", "core")

    imports = install_deps | unit_deps

    out = f"package {pkg}\n\n"
    out += "\n"
    out += _format_imports(sorted(imports))
    out += "\n"
    out += install_code
    out += "\n"
    out += unit_code

    return out


def generate_unit_code(u: Unit) -> tuple[str, set[str]]:
    """Generate Go code for a unit type.

    Returns (code, import_set).
    """
    import_set: set[str] = set()
    out = ""

    out += _format_unit_doc(u)

    block_type = to_pascal_case(u.name)
    sub_block_type = block_type + "Block"
    has_block_type = len(u.options) > 0

    if has_block_type:
        code, deps = generate_block_struct(u.options, u.name, "core")
        import_set |= deps
        out += code
        out += "\n"

    out += f"type {block_type} struct {{\n"
    out += '\tName string `hcl:"name,label"`\n\n'
    out += '\tUnit    UnitBlock    `hcl:"unit,block"`\n'

    if has_block_type:
        snake = to_snake_case(block_type)
        out += f'\t{block_type} {sub_block_type} `hcl:"{snake},block"`\n'

    out += '\tInstall InstallBlock `hcl:"install,block"`\n'
    out += "}\n"

    return out, import_set


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
