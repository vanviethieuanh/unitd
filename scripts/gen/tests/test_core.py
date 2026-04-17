"""Tests for scripts.gen — ported from cmd/gen/gen_test.go and cmd/gen/type_test.go."""

from __future__ import annotations

import pytest

from scripts.gen.types import TypeExpr, ValueType, parse_type_expr
from scripts.gen.util import to_pascal_case, to_snake_case


# ---------- toPascalCase tests (from gen_test.go) ----------


@pytest.mark.parametrize(
    "input_,expected",
    [
        ("timer", "Timer"),
        ("on-calendar", "OnCalendar"),
        ("OnActiveSec", "OnActiveSec"),
        ("AccuracySec", "AccuracySec"),
        ("on-boot-sec", "OnBootSec"),
    ],
)
def test_to_pascal_case(input_: str, expected: str) -> None:
    assert to_pascal_case(input_) == expected


# ---------- toSnakeCase tests (from gen_test.go) ----------


@pytest.mark.parametrize(
    "input_,expected",
    [
        ("OnActiveSec", "on_active_sec"),
        ("OnCalendar", "on_calendar"),
        ("AccuracySec", "accuracy_sec"),
        ("on-boot-sec", "on_boot_sec"),
    ],
)
def test_to_snake_case(input_: str, expected: str) -> None:
    assert to_snake_case(input_) == expected


# ---------- ParseTypeExpr tests (from type_test.go) ----------


def test_parse_simple() -> None:
    expr = parse_type_expr("PATH")
    assert expr == TypeExpr(base=ValueType.PATH)


def test_parse_simple_list() -> None:
    expr = parse_type_expr("PATH [...]")
    assert expr == TypeExpr(base=ValueType.PATH, repeated=True)


def test_parse_nested_list() -> None:
    expr = parse_type_expr("PATH [ARGUMENT [...]]")
    assert expr == TypeExpr(
        base=ValueType.PATH,
        inner=TypeExpr(base=ValueType.ARGUMENT, repeated=True),
    )


def test_parse_deep_nesting() -> None:
    expr = parse_type_expr("PATH [ARGUMENT [UNIT [STRING [...]]]]")
    assert expr == TypeExpr(
        base=ValueType.PATH,
        inner=TypeExpr(
            base=ValueType.ARGUMENT,
            inner=TypeExpr(
                base=ValueType.UNIT,
                inner=TypeExpr(base=ValueType.STRING, repeated=True),
            ),
        ),
    )


def test_parse_missing_closing_bracket() -> None:
    with pytest.raises(ValueError):
        parse_type_expr("PATH [ARGUMENT [...]")


def test_parse_unexpected_tokens() -> None:
    with pytest.raises(ValueError):
        parse_type_expr("PATH PATH")


def test_parse_empty() -> None:
    with pytest.raises(ValueError):
        parse_type_expr("")


# ---------- TypeExpr.toGoType tests (from type_test.go) ----------


def test_go_type_single_string() -> None:
    expr = TypeExpr(base=ValueType.STRING)
    assert expr.to_go_type()[0] == "string"


def test_go_type_repeated_string() -> None:
    expr = TypeExpr(base=ValueType.STRING, repeated=True)
    assert expr.to_go_type()[0] == "[]string"


def test_go_type_two_layer_nested() -> None:
    expr = TypeExpr(
        base=ValueType.STRING,
        inner=TypeExpr(base=ValueType.STRING, repeated=True),
    )
    assert expr.to_go_type()[0] == "[][]string"


def test_go_type_three_layer_nested() -> None:
    expr = TypeExpr(
        base=ValueType.INTEGER,
        inner=TypeExpr(
            base=ValueType.LEVEL,
            inner=TypeExpr(base=ValueType.SERVICEEXITTYPE, repeated=True),
        ),
    )
    assert expr.to_go_type()[0] == "[][][]int"


def test_go_type_mismatch_error() -> None:
    expr = TypeExpr(
        base=ValueType.STRING,
        inner=TypeExpr(base=ValueType.BOOLEAN),
    )
    with pytest.raises(ValueError, match="mismatch"):
        expr.to_go_type()
