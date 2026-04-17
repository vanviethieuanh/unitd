"""XML man page parser for systemd documentation."""

from __future__ import annotations

import re
import xml.etree.ElementTree as ET
from dataclasses import dataclass, field
from pathlib import Path

from .directive import Directive


@dataclass
class Unit:
    name: str
    title: str
    purpose: str
    description: str
    options: list[Directive] = field(default_factory=list)


def parse_unit(
    path: Path,
    directives_list: list[Directive],
    extra_descriptions: dict[str, str] | None = None,
) -> Unit:
    """Parse a systemd XML man page and enrich directives with descriptions."""
    if extra_descriptions is None:
        extra_descriptions = {}

    tree = ET.parse(path)
    root = tree.getroot()

    unit_name = _parse_name(root)
    parsed_options = _parse_options(root)

    options: list[Directive] = []
    seen: set[str] = set()
    for d in directives_list:
        if d.identifier.section.lower() != unit_name.lower():
            continue

        key = d.identifier.key
        if key in seen:
            continue
        seen.add(key)

        if key in parsed_options:
            d.description = parsed_options[key]
        elif key in extra_descriptions:
            d.description = extra_descriptions[key]

        options.append(d)

    return Unit(
        name=unit_name,
        title=_parse_title(root),
        purpose=_parse_purpose(root),
        description=_parse_description(root),
        options=options,
    )


def parse_applicable_types(path: Path) -> list[str]:
    """Extract unit type names from <refsynopsisdiv>.

    For example, systemd.kill.xml lists .service, .socket, .mount, .swap, .scope
    and this returns ["service", "socket", "mount", "swap", "scope"].
    """
    tree = ET.parse(path)
    root = tree.getroot()

    types: list[str] = []
    synopsis_div = root.find("refsynopsisdiv")
    if synopsis_div is not None:
        para = synopsis_div.find("para")
        if para is not None:
            for fn in para.iter("filename"):
                text = "".join(fn.itertext())
                idx = text.rfind(".")
                if idx >= 0:
                    types.append(text[idx + 1 :])

    return types


def parse_descriptions(path: Path) -> dict[str, str]:
    """Extract option key→description mappings from an XML man page."""
    tree = ET.parse(path)
    root = tree.getroot()
    return _parse_options(root)


def _parse_name(root: ET.Element) -> str:
    el = root.find("refnamediv/refname")
    name = (el.text or "") if el is not None else ""
    if name.startswith("systemd."):
        name = name[len("systemd.") :]
    return name


def _parse_title(root: ET.Element) -> str:
    el = root.find("refmeta/refentrytitle")
    return (el.text or "") if el is not None else ""


def _parse_purpose(root: ET.Element) -> str:
    el = root.find("refnamediv/refpurpose")
    return (el.text or "") if el is not None else ""


def _parse_description(root: ET.Element) -> str:
    for sect in root.iter("refsect1"):
        title_el = sect.find("title")
        if title_el is not None and title_el.text == "Description":
            paras = [_text_only(p) for p in sect.findall("para")]
            return _join_cleaned(paras)
    return ""


def _parse_options(root: ET.Element) -> dict[str, str]:
    options: dict[str, str] = {}

    for sect in root.iter("refsect1"):
        title_el = sect.find("title")
        if title_el is not None and title_el.text == "Description":
            continue

        varlist = sect.find("variablelist")
        if varlist is None:
            continue

        for entry in varlist.findall("varlistentry"):
                terms: list[str] = []
                for term_el in entry.findall("term/varname"):
                    t = (term_el.text or "").strip().rstrip("=")
                    if t:
                        terms.append(t)

                listitem = entry.find("listitem")
                if listitem is None:
                    continue

                table = listitem.find("table")
                if table is not None:
                    for term in terms:
                        desc = _extract_table_description(table, term)
                        if desc:
                            options[term] = desc
                else:
                    desc = _extract_full_description(listitem)
                    for term in terms:
                        if desc:
                            options[term] = desc

    return options


def _extract_table_description(table: ET.Element, term: str) -> str:
    clean_term = term.rstrip("=")
    tbody = table.find("tgroup/tbody")
    if tbody is None:
        return ""

    for row in tbody.findall("row"):
        entries = row.findall("entry")
        if len(entries) >= 2:
            first = _clean_text(_inner_xml(entries[0]))
            first = first.removeprefix("<varname>").removesuffix("</varname>")
            first = first.rstrip("=")
            if clean_term in first:
                return _clean_text(_inner_xml(entries[1]))
    return ""


def _extract_full_description(listitem: ET.Element) -> str:
    texts = [_inner_xml(p) for p in listitem.findall("para")]
    return _join_cleaned(texts)


def _join_cleaned(texts: list[str]) -> str:
    parts = [_clean_text(t) for t in texts if t]
    parts = [p for p in parts if p]
    return "\n\n".join(parts)


_XML_TAG_RE = re.compile(
    r"</?(?:literal|filename|varname|replaceable|option|constant|command|emphasis)>"
)


def _clean_text(text: str) -> str:
    text = text.strip()
    text = _XML_TAG_RE.sub("", text)
    # Normalize whitespace
    text = " ".join(text.split())
    return text


def _inner_xml(el: ET.Element) -> str:
    """Get the inner XML (including sub-tags) of an element as a string."""
    # ET.tostring gives us <tag>...</tag>, we want what's between
    raw = ET.tostring(el, encoding="unicode", method="xml")
    # Strip outer tag
    # Find first > and last <
    start = raw.find(">")
    end = raw.rfind("<")
    if start >= 0 and end > start:
        return raw[start + 1 : end]
    return ""


def _text_only(el: ET.Element) -> str:
    """Extract only direct text nodes, discarding child elements."""
    parts: list[str] = []
    if el.text:
        parts.append(el.text)
    for child in el:
        # Skip child element content entirely — only grab tail text
        if child.tail:
            parts.append(child.tail)
    return "".join(parts)
