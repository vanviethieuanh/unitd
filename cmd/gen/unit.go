package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/log"
)

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

func (d *Directive) UnmarshalJSON(b []byte) error {
	type alias struct {
		System   string `json:"system"`
		Section  string `json:"section"`
		Property string `json:"property"`
		Type     string `json:"type"`
	}

	var tmp alias
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	parsedType, err := ParseTypeExpr(tmp.Type)
	if err != nil {
		return fmt.Errorf("failed to parse type expression %q: %w", tmp.Type, err)
	}

	goType, deps := parsedType.toGoType()

	d.System = tmp.System
	d.Type = goType
	d.Deps = deps
	d.NativeType = len(deps) == 0
	d.Identifier = DirectiveIdentifier{
		Section: tmp.Section,
		Key:     tmp.Property,
	}

	return nil
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

func getDirectives(f io.Reader) []Directive {

	scanner := bufio.NewScanner(f)

	var result []Directive

	for scanner.Scan() {
		var d Directive
		if err := json.Unmarshal(scanner.Bytes(), &d); err != nil {
			log.Fatalf("decode error: %v", err)
		}

		result = append(result, d)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

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
