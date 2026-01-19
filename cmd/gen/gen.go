package main

import (
	"fmt"
	"strings"
)

func genCommon(directives []Directive) string {
	var out strings.Builder

	installCode, installDeps := generateBlockStruct(directives, "Install", "core")
	unitCode, unitDeps := generateBlockStruct(directives, "Unit", "core")

	imports := NewImportSet()
	imports.Merge(installDeps)
	imports.Merge(unitDeps)

	out.WriteString("package configs\n\n")
	out.WriteString("\n")
	writeImports(&out, imports.Sorted())
	out.WriteString("\n")
	out.WriteString(installCode)
	out.WriteString("\n")
	out.WriteString(unitCode)

	return out.String()
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

	for _, d := range directives {
		if !strings.EqualFold(d.Identifier.Section, section) || !strings.EqualFold(d.System, system) {
			continue
		}

		WriteDirective(&out, &d)
		imports.AddAll(d.Deps)
	}

	out.WriteString("}\n")
	return out.String(), imports
}
