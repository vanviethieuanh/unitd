#!/usr/bin/env python3
"""Expand a Jinja2 .gperf.in template to a plain .gperf file.

Outputs the expanded content to stdout.  All HAVE_* / ENABLE_* feature
flags default to True so we get the full superset of directives.
"""

import sys
from pathlib import Path

import jinja2


def main() -> None:
    if len(sys.argv) != 2:
        print(f"usage: {sys.argv[0]} <template.gperf.in>", file=sys.stderr)
        sys.exit(1)

    template_path = Path(sys.argv[1])
    template_dir = template_path.parent
    template_name = template_path.name

    env = jinja2.Environment(
        loader=jinja2.FileSystemLoader(str(template_dir)),
        undefined=jinja2.Undefined,
        keep_trailing_newline=True,
    )

    template = env.get_template(template_name)

    # Default all feature flags to True for maximum directive coverage.
    context = {
        "HAVE_SECCOMP": True,
        "HAVE_PAM": True,
        "HAVE_SELINUX": True,
        "HAVE_APPARMOR": True,
        "ENABLE_SMACK": True,
    }

    sys.stdout.write(template.render(context))


if __name__ == "__main__":
    main()
