package main

import (
	"fmt"
	"strings"
)

type ValueType string

const (
	TypeAccess             ValueType = "ACCESS"
	TypeAction             ValueType = "ACTION"
	TypeArgument           ValueType = "ARGUMENT"
	TypeBoolean            ValueType = "BOOLEAN"
	TypeCondition          ValueType = "CONDITION"
	TypeInteger            ValueType = "INTEGER"
	TypeLong               ValueType = "LONG"
	TypeMode               ValueType = "MODE"
	TypeNetworkInterface   ValueType = "NETWORKINTERFACE"
	TypeNode               ValueType = "NODE"
	TypeNotSupported       ValueType = "NOTSUPPORTED"
	TypePath               ValueType = "PATH"
	TypeSeconds            ValueType = "SECONDS"
	TypeService            ValueType = "SERVICE"
	TypeServiceExitType    ValueType = "SERVICEEXITTYPE"
	TypeServiceRestart     ValueType = "SERVICERESTART"
	TypeServiceRestartMode ValueType = "SERVICERESTARTMODE"
	TypeServiceType        ValueType = "SERVICETYPE"
	TypeSignal             ValueType = "SIGNAL"
	TypeSize               ValueType = "SIZE"
	TypeSocket             ValueType = "SOCKET"
	TypeSocketBind         ValueType = "SOCKETBIND"
	TypeSockets            ValueType = "SOCKETS"
	TypeStatus             ValueType = "STATUS"
	TypeString             ValueType = "STRING"
	TypeTimeoutMode        ValueType = "TIMEOUTMODE"
	TypeTimer              ValueType = "TIMER"
	TypeTOS                ValueType = "TOS"
	TypeUnit               ValueType = "UNIT"
	TypeUnknown            ValueType = "UNKNOWN"
	TypeLevel              ValueType = "LEVEL"
	TypeUnsigned           ValueType = "UNSIGNED"
	TypeURL                ValueType = "URL"
)

func (v ValueType) toGoType() (goType string, imports []string) {
	switch v {
	case TypeAccess,
		TypeAction,
		TypeArgument,
		TypeCondition,
		TypeNetworkInterface,
		TypeNode,
		TypePath,
		TypeService,
		TypeServiceRestart,
		TypeServiceRestartMode,
		TypeServiceType,
		TypeStatus,
		TypeString,
		TypeTimeoutMode,
		TypeURL:
		return "string", nil

	case TypeUnit:
		return "hcl.Traversal", []string{"github.com/hashicorp/hcl/v2"}

	case TypeBoolean:
		return "bool", nil

	case TypeInteger, TypeSeconds:
		return "int", nil

	case TypeLong, TypeSize:
		return "int64", nil

	case TypeUnsigned:
		return "uint64", nil

	case TypeMode:
		return "os.FileMode", []string{"os"}

	case TypeTimer:
		return "time.Duration", []string{"time"}

	case TypeServiceExitType:
		return "int", nil

	case TypeSignal:
		return "syscall.Signal", []string{"syscall"}

	case TypeSockets:
		return "[]string", nil

	case TypeLevel:
		return "int", nil

	case TypeUnknown, TypeNotSupported:
		return "string", nil

	default:
		return "string", nil
	}
}

type tokenKind int

const (
	tokenIdent tokenKind = iota
	tokenLBracket
	tokenRBracket
	tokenEllipsis
)

type token struct {
	Kind tokenKind
	Text string
}

func tokenizeType(s string) ([]token, error) {
	var tokens []token

	for i := 0; i < len(s); {
		switch {
		case s[i] == ' ' || s[i] == '\t' || s[i] == '\n':
			i++

		case s[i] == '[':
			tokens = append(tokens, token{Kind: tokenLBracket, Text: "["})
			i++

		case s[i] == ']':
			tokens = append(tokens, token{Kind: tokenRBracket, Text: "]"})
			i++

		case strings.HasPrefix(s[i:], "..."):
			tokens = append(tokens, token{Kind: tokenEllipsis, Text: "..."})
			i += 3

		case isIdentStart(s[i]):
			start := i
			i++
			for i < len(s) && isIdentPart(s[i]) {
				i++
			}
			tokens = append(tokens, token{
				Kind: tokenIdent,
				Text: s[start:i],
			})

		default:
			return nil, fmt.Errorf("unexpected character %q at %d", s[i], i)
		}
	}

	return tokens, nil
}

func isIdentStart(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

func isIdentPart(b byte) bool {
	return (b >= 'A' && b <= 'Z') || b == '_'
}

type TypeExpr struct {
	Base     ValueType
	Inner    *TypeExpr
	Repeated bool
}

func (e *TypeExpr) toGoType() (goType string, deps []string) {
	if e == nil {
		return "any", nil
	}

	var (
		baseType  string
		depth     int
		importSet = make(map[string]struct{})
	)

	for cur := e; cur != nil; cur = cur.Inner {
		bt, imps := cur.Base.toGoType()
		for _, p := range imps {
			importSet[p] = struct{}{}
		}

		if baseType == "" {
			baseType = bt
		} else if baseType != bt {
			panic(fmt.Sprintf("nested type mismatch: %q vs %q", baseType, bt))
		}

		if e.Inner != nil {
			depth++
		}
	}

	if e.Repeated {
		depth++
	}

	goType = strings.Repeat("[]", depth) + baseType

	deps = make([]string, 0, len(importSet))
	for p := range importSet {
		deps = append(deps, p)
	}

	return goType, deps
}

func parseType(tokens []token) (*TypeExpr, int, error) {
	if len(tokens) == 0 {
		return nil, 0, fmt.Errorf("empty token stream")
	}

	if tokens[0].Kind != tokenIdent {
		return nil, 0, fmt.Errorf("expected identifier")
	}

	expr := &TypeExpr{
		Base: ValueType(tokens[0].Text),
	}

	pos := 1

	if pos < len(tokens) && tokens[pos].Kind == tokenLBracket {
		pos++

		if pos < len(tokens) && tokens[pos].Kind == tokenEllipsis {
			expr.Repeated = true
			pos++
		} else {
			inner, consumed, err := parseType(tokens[pos:])
			if err != nil {
				return nil, 0, err
			}
			expr.Inner = inner
			pos += consumed
		}

		if pos >= len(tokens) || tokens[pos].Kind != tokenRBracket {
			return nil, 0, fmt.Errorf("expected ']'")
		}
		pos++
	}

	return expr, pos, nil
}

func ParseTypeExpr(s string) (*TypeExpr, error) {
	tokens, err := tokenizeType(s)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize expression %q: %w", s, err)
	}

	expr, consumed, err := parseType(tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression %q: %w", s, err)
	}

	if consumed != len(tokens) {
		return nil, fmt.Errorf(
			"unexpected trailing, %d tokens left in %q",
			len(tokens)-consumed,
			s,
		)
	}

	return expr, nil
}
