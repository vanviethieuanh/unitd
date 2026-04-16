package main

import (
	"fmt"
	"strings"
)

type Unit struct {
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	Purpose     string      `json:"purpose"`
	Description string      `json:"description"`
	Options     []Directive `json:"options"`
}

func GenerateUnitCode(u *Unit) (string, ImportSet) {
	var out strings.Builder
	var importSet = NewImportSet()

	writeUnitDoc(&out, u)

	blockType := toPascalCase(u.Name)
	subBlockType := blockType + "Block"
	hasBlockType := len(u.Options) > 0

	if hasBlockType {
		code, deps := generateBlockStruct(u.Options, u.Name, "core")
		importSet.Merge(deps)
		out.WriteString(code)
		out.WriteString("\n")
	}

	fmt.Fprintf(&out, "type %s struct {\n", blockType)
	out.WriteString("\tName string `hcl:\"name,label\"`\n\n")
	out.WriteString("\tUnit    UnitBlock    `hcl:\"unit,block\"`\n")

	if hasBlockType {
		fmt.Fprintf(
			&out,
			"\t%s %s `hcl:\"%s,block\"`\n",
			blockType,
			subBlockType,
			toSnakeCase(blockType),
		)
	}

	out.WriteString("\tInstall InstallBlock `hcl:\"install,block\"`\n")
	out.WriteString("}\n")

	return out.String(), importSet
}

func writeUnitDoc(out *strings.Builder, u *Unit) {
	if u.Description == "" {
		return
	}

	name := toPascalCase(u.Name)

	if len(u.Options) > 0 {
		fmt.Fprintf(out, "// %sBlock is for [%s] systemd unit block\n", name, name)
		fmt.Fprint(out, "//\n")
	} else {

		fmt.Fprintf(out, "// %s is for %s systemd unit file\n", name, name)
	}

	for _, line := range wrapComment(u.Description, 100) {
		out.WriteString("// " + line + "\n")
	}
}
