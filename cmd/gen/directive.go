package main

import (
	"fmt"
	"sort"
	"strings"
)

// IncludeExec lists sections that should include Exec directives.
// TODO: wire this up (see TODO.md)
var IncludeExec = map[string]struct{}{
	"Service": {},
	"Socket":  {},
	"Mount":   {},
	"Swap":    {},
}

type DirectiveIdentifier struct {
	Section string
	Key     string
}

type Directive struct {
	Identifier  DirectiveIdentifier
	Type        string
	System      string
	Description string
	Deps        []string

	// Indicate if this type is natively to parse by hcl tag
	NativeType bool
}

// NewDirective creates a Directive by parsing a type expression string
// (e.g. "STRING", "UNIT [...]") into the corresponding Go type and imports.
func NewDirective(section, key, system, typeStr string) (Directive, error) {
	parsedType, err := ParseTypeExpr(typeStr)
	if err != nil {
		return Directive{}, fmt.Errorf("failed to parse type expression %q: %w", typeStr, err)
	}

	goType, deps := parsedType.toGoType()

	return Directive{
		Identifier: DirectiveIdentifier{
			Section: section,
			Key:     key,
		},
		Type:       goType,
		System:     system,
		Deps:       deps,
		NativeType: len(deps) == 0,
	}, nil
}

func WriteDirective(builder *strings.Builder, d *Directive) {
	if d.Description != "" {
		paragraphs := strings.Split(d.Description, "\n\n")
		for _, paragraph := range paragraphs {
			paragraph = strings.TrimSpace(paragraph)
			if paragraph != "" {
				wrappedLines := wrapComment(paragraph, 100)
				for _, wrapped := range wrappedLines {
					builder.WriteString("\t// " + wrapped + "\n")
				}
				if len(paragraphs) > 1 {
					builder.WriteString("\t//\n")
				}
			}
		}
	}

	fieldName := toPascalCase(d.Identifier.Key)
	snakeName := toSnakeCase(d.Identifier.Key)
	systemdName := d.Identifier.Key

	if d.NativeType {
		fmt.Fprintf(builder, "\t%s %s `hcl:\"%s,optional\" systemd:\"%s\"`\n",
			fieldName, d.Type, snakeName, systemdName)
	} else {
		fmt.Fprintf(builder, "\t%s %s `unitd:\"%s,optional\" systemd:\"%s\"`\n",
			fieldName, d.Type, snakeName, systemdName)
	}
}

func generateBlockStruct(
	directives []Directive,
	section string,
	system string,
) (string, ImportSet) {
	var out strings.Builder
	imports := NewImportSet()

	typeName := toPascalCase(section) + "Block"
	fmt.Fprintf(&out, "type %s struct {\n", typeName)

	var filtered []Directive
	for _, d := range directives {
		if !strings.EqualFold(d.Identifier.Section, section) || !strings.EqualFold(d.System, system) {
			continue
		}
		filtered = append(filtered, d)
	}
	sort.Slice(filtered, func(i, j int) bool {
		return toPascalCase(filtered[i].Identifier.Key) < toPascalCase(filtered[j].Identifier.Key)
	})

	for _, d := range filtered {
		WriteDirective(&out, &d)
		imports.AddAll(d.Deps)
	}

	out.WriteString("}\n")
	return out.String(), imports
}
