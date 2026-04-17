"""Type expression parser and Go type mapping."""

from __future__ import annotations

from dataclasses import dataclass, field
from enum import StrEnum


class ValueType(StrEnum):
    ACCESS = "ACCESS"
    ACTION = "ACTION"
    ARGUMENT = "ARGUMENT"
    BOOLEAN = "BOOLEAN"
    CONDITION = "CONDITION"
    INTEGER = "INTEGER"
    LEVEL = "LEVEL"
    LONG = "LONG"
    MODE = "MODE"
    NETWORKINTERFACE = "NETWORKINTERFACE"
    NODE = "NODE"
    NOTSUPPORTED = "NOTSUPPORTED"
    PATH = "PATH"
    SECONDS = "SECONDS"
    SERVICE = "SERVICE"
    SERVICEEXITTYPE = "SERVICEEXITTYPE"
    SERVICERESTART = "SERVICERESTART"
    SERVICERESTARTMODE = "SERVICERESTARTMODE"
    SERVICETYPE = "SERVICETYPE"
    SIGNAL = "SIGNAL"
    SIZE = "SIZE"
    SOCKET = "SOCKET"
    SOCKETBIND = "SOCKETBIND"
    SOCKETS = "SOCKETS"
    STATUS = "STATUS"
    STRING = "STRING"
    TIMEOUTMODE = "TIMEOUTMODE"
    TIMER = "TIMER"
    TOS = "TOS"
    UNIT = "UNIT"
    UNKNOWN = "UNKNOWN"
    UNSIGNED = "UNSIGNED"
    URL = "URL"
    IOCLASS = "IOCLASS"


_GO_TYPE_MAP: dict[ValueType, tuple[str, list[str]]] = {
    ValueType.ACCESS: ("string", []),
    ValueType.ACTION: ("string", []),
    ValueType.ARGUMENT: ("string", []),
    ValueType.CONDITION: ("string", []),
    ValueType.NETWORKINTERFACE: ("string", []),
    ValueType.NODE: ("string", []),
    ValueType.PATH: ("string", []),
    ValueType.SERVICE: ("string", []),
    ValueType.SERVICERESTART: ("string", []),
    ValueType.SERVICERESTARTMODE: ("string", []),
    ValueType.SERVICETYPE: ("string", []),
    ValueType.STATUS: ("string", []),
    ValueType.STRING: ("string", []),
    ValueType.TIMEOUTMODE: ("string", []),
    ValueType.URL: ("string", []),
    ValueType.UNIT: ("hcl.Traversal", ["github.com/hashicorp/hcl/v2"]),
    ValueType.BOOLEAN: ("bool", []),
    ValueType.INTEGER: ("int", []),
    ValueType.SECONDS: ("int", []),
    ValueType.LONG: ("int64", []),
    ValueType.SIZE: ("int64", []),
    ValueType.UNSIGNED: ("uint64", []),
    ValueType.MODE: ("os.FileMode", ["os"]),
    ValueType.TIMER: ("time.Duration", ["time"]),
    ValueType.SERVICEEXITTYPE: ("int", []),
    ValueType.SIGNAL: ("syscall.Signal", ["syscall"]),
    ValueType.SOCKETS: ("[]string", []),
    ValueType.LEVEL: ("int", []),
    ValueType.UNKNOWN: ("string", []),
    ValueType.NOTSUPPORTED: ("string", []),
    ValueType.IOCLASS: ("string", []),
}


def value_type_to_go(vt: ValueType) -> tuple[str, list[str]]:
    """Return (go_type, imports) for a base ValueType."""
    return _GO_TYPE_MAP.get(vt, ("string", []))


class _TokenKind:
    IDENT = 0
    LBRACKET = 1
    RBRACKET = 2
    ELLIPSIS = 3


@dataclass
class _Token:
    kind: int
    text: str


def _tokenize(s: str) -> list[_Token]:
    tokens: list[_Token] = []
    i = 0
    while i < len(s):
        ch = s[i]
        if ch in (" ", "\t", "\n"):
            i += 1
        elif ch == "[":
            tokens.append(_Token(_TokenKind.LBRACKET, "["))
            i += 1
        elif ch == "]":
            tokens.append(_Token(_TokenKind.RBRACKET, "]"))
            i += 1
        elif s[i:i + 3] == "...":
            tokens.append(_Token(_TokenKind.ELLIPSIS, "..."))
            i += 3
        elif ch.isupper():
            start = i
            i += 1
            while i < len(s) and (s[i].isupper() or s[i] == "_"):
                i += 1
            tokens.append(_Token(_TokenKind.IDENT, s[start:i]))
        else:
            raise ValueError(f"unexpected character {ch!r} at {i}")
    return tokens


@dataclass
class TypeExpr:
    base: ValueType
    inner: TypeExpr | None = field(default=None)
    repeated: bool = field(default=False)

    def to_go_type(self) -> tuple[str, list[str]]:
        """Return (go_type_string, import_paths)."""
        if self is None:
            return "any", []

        base_type = ""
        depth = 0
        import_set: set[str] = set()

        cur: TypeExpr | None = self
        while cur is not None:
            bt, imps = value_type_to_go(cur.base)
            import_set.update(imps)

            if not base_type:
                base_type = bt
            elif base_type != bt:
                raise ValueError(f"nested type mismatch: {base_type!r} vs {bt!r}")

            if self.inner is not None:
                depth += 1

            cur = cur.inner

        if self.repeated:
            depth += 1

        go_type = "[]" * depth + base_type
        return go_type, sorted(import_set)


def _parse_type(tokens: list[_Token], pos: int = 0) -> tuple[TypeExpr, int]:
    if pos >= len(tokens):
        raise ValueError("empty token stream")

    if tokens[pos].kind != _TokenKind.IDENT:
        raise ValueError("expected identifier")

    expr = TypeExpr(base=ValueType(tokens[pos].text))
    pos += 1

    if pos < len(tokens) and tokens[pos].kind == _TokenKind.LBRACKET:
        pos += 1  # consume '['

        if pos < len(tokens) and tokens[pos].kind == _TokenKind.ELLIPSIS:
            expr.repeated = True
            pos += 1
        else:
            inner, pos = _parse_type(tokens, pos)
            expr.inner = inner

        if pos >= len(tokens) or tokens[pos].kind != _TokenKind.RBRACKET:
            raise ValueError("expected ']'")
        pos += 1

    return expr, pos


def parse_type_expr(s: str) -> TypeExpr:
    """Parse a type expression string like 'STRING', 'UNIT [...]', etc."""
    tokens = _tokenize(s)
    expr, consumed = _parse_type(tokens)

    if consumed != len(tokens):
        raise ValueError(
            f"unexpected trailing, {len(tokens) - consumed} tokens left in {s!r}"
        )

    return expr
