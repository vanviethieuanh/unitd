package main

import (
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
