package main

import (
	"fmt"
	"strings"
)

func genCommon(pkg string, directives []Directive) string {
	var out strings.Builder

	installCode, installDeps := generateBlockStruct(directives, "Install", "core")
	unitCode, unitDeps := generateBlockStruct(directives, "Unit", "core")

	imports := NewImportSet()
	imports.Merge(installDeps)
	imports.Merge(unitDeps)

	fmt.Fprintf(&out, "package %s\n\n", pkg)
	out.WriteString("\n")
	writeImports(&out, imports.Sorted())
	out.WriteString("\n")
	out.WriteString(installCode)
	out.WriteString("\n")
	out.WriteString(unitCode)

	return out.String()
}
