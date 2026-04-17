"""String utilities: case conversion and comment wrapping."""

from __future__ import annotations

import re
import textwrap
import unicodedata


def _split_words(s: str) -> list[str]:
    """Split a string into words on non-alphanumeric boundaries."""
    words: list[str] = []
    current: list[str] = []
    for ch in s:
        if ch.isalnum():
            current.append(ch)
        else:
            if current:
                words.append("".join(current))
                current = []
    if current:
        words.append("".join(current))
    return words


def to_pascal_case(s: str) -> str:
    if not s:
        return ""

    # If already mixed-case (has both upper and lower), return as-is.
    if s[0].isupper():
        has_lower = any(c.islower() for c in s)
        has_upper = any(c.isupper() for c in s)
        if has_lower and has_upper:
            return s

    words = _split_words(s)
    result: list[str] = []
    for word in words:
        if not word:
            continue
        result.append(word[0].upper() + word[1:].lower())
    return "".join(result)


def to_snake_case(s: str) -> str:
    if not s:
        return ""

    if "-" in s or "_" in s:
        return s.lower().replace("-", "_")

    out: list[str] = []
    runes = list(s)
    for i, r in enumerate(runes):
        if i > 0:
            prev = runes[i - 1]
            if prev.islower() and r.isupper():
                out.append("_")
            elif (
                prev.isupper()
                and r.isupper()
                and i + 1 < len(runes)
                and runes[i + 1].islower()
            ):
                out.append("_")
        out.append(r.lower())
    return "".join(out)


def wrap_comment(text: str, width: int = 100) -> list[str]:
    """Word-wrap *text* into lines of at most *width* characters."""
    if len(text) <= width:
        return [text]

    lines: list[str] = []
    current: list[str] = []
    current_len = 0

    for word in text.split():
        if current and current_len + 1 + len(word) > width:
            lines.append(" ".join(current))
            current = []
            current_len = 0

        current.append(word)
        current_len += (1 if current_len > 0 else 0) + len(word)

    if current:
        lines.append(" ".join(current))

    return lines
