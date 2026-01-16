package main

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

func genCommon(propTypes map[DirectiveIdentifier]string) string {
	var sb strings.Builder

	sb.WriteString("package configs\n\n")
	sb.WriteString(genSection(propTypes, "Install"))
	sb.WriteString("\n")
	sb.WriteString(genSection(propTypes, "Unit"))

	return sb.String()
}

func genSection(propTypes map[DirectiveIdentifier]string, section string) string {
	var sb strings.Builder

	typeName := toPascalCase(section) + "Block"
	sb.WriteString("type " + typeName + " struct {\n")

	ids := make([]DirectiveIdentifier, 0)
	for id := range propTypes {
		if id.Section == section {
			ids = append(ids, id)
		}
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Key < ids[j].Key
	})

	for _, id := range ids {
		prop := propTypes[id]
		fieldName := toPascalCase(id.Key)
		snakeName := toSnakeCase(id.Key)

		fmt.Fprintf(
			&sb,
			"\t%s %s `hcl:\"%s,optional\" systemd:\"%s\"`\n",
			fieldName,
			prop,
			snakeName,
			id.Key,
		)
	}

	sb.WriteString("}\n")
	return sb.String()
}

func gen(u *Unit, propTypes map[DirectiveIdentifier]string) string {
	var sb strings.Builder

	sb.WriteString("package configs\n\n")

	blockTypeName := toPascalCase(u.Name)
	subBlockTypeName := toPascalCase(u.Name) + "Block"
	hasBlockType := len(u.Options) > 0

	if u.Description != "" {
		if hasBlockType {
			sb.WriteString("// " + toPascalCase(u.Name) + "Block represents the [" + toPascalCase(u.Name) + "] section of a systemd unit.\n")
		} else {
			sb.WriteString("// " + toPascalCase(u.Name) + " represents the " + toPascalCase(u.Name) + " configuration file of a systemd unit.\n")
			sb.WriteString("// A separate [" + toPascalCase(u.Name) + "] section does not exist.\n")
		}
		sb.WriteString("//\n")
		paragraphs := strings.SplitSeq(u.Description, "\n\n")
		for paragraph := range paragraphs {
			paragraph = strings.TrimSpace(paragraph)
			if paragraph != "" {
				wrappedLines := wrapComment(paragraph, 100)
				for _, wrapped := range wrappedLines {
					sb.WriteString("// " + wrapped + "\n")
				}
				sb.WriteString("//\n")
			}
		}
	}

	if hasBlockType {
		sb.WriteString("type " + subBlockTypeName + " struct {\n")

		optionNames := make([]string, 0, len(u.Options))
		for optionName := range u.Options {
			optionNames = append(optionNames, optionName)
		}
		sort.Strings(optionNames)

		for _, optionName := range optionNames {
			optionDesc := u.Options[optionName]

			if optionDesc != "" {
				paragraphs := strings.Split(optionDesc, "\n\n")
				for _, paragraph := range paragraphs {
					paragraph = strings.TrimSpace(paragraph)
					if paragraph != "" {
						wrappedLines := wrapComment(paragraph, 100)
						for _, wrapped := range wrappedLines {
							sb.WriteString("\t// " + wrapped + "\n")
						}
						if len(paragraphs) > 1 {
							sb.WriteString("\t//\n")
						}
					}
				}
			}

			fieldName := toPascalCase(optionName)
			snakeName := toSnakeCase(optionName)
			systemdName := optionName

			propType, ok := propTypes[DirectiveIdentifier{
				Section: blockTypeName,
				Key:     systemdName,
			}]
			if !ok || propType == "" {
				propType = "string"
			}

			fmt.Fprintf(&sb, "\t%s %s `hcl:\"%s,optional\" systemd:\"%s\"`\n",
				fieldName, propType, snakeName, systemdName)
		}

		sb.WriteString("}\n\n")

	}
	sb.WriteString("type " + blockTypeName + " struct {\n")

	sb.WriteString("Name string `hcl:\"name,label\"`\n\n")
	sb.WriteString("Unit    UnitBlock    `hcl:\"unit,block\"`\n")
	if hasBlockType {
		fmt.Fprintf(&sb, "%s %s `hcl:\"%s,block\"`\n", blockTypeName, subBlockTypeName, toSnakeCase(blockTypeName))
	}
	sb.WriteString("Install InstallBlock `hcl:\"install,block\"`\n")

	sb.WriteString("}\n")

	return sb.String()
}

func toPascalCase(s string) string {
	if s == "" {
		return ""
	}

	if len(s) > 0 && unicode.IsUpper(rune(s[0])) {
		hasLower := false
		hasUpper := false
		for _, r := range s {
			if unicode.IsLower(r) {
				hasLower = true
			}
			if unicode.IsUpper(r) {
				hasUpper = true
			}
		}
		if hasLower && hasUpper {
			return s
		}
	}

	words := splitWords(s)
	var result strings.Builder

	for _, word := range words {
		if word == "" {
			continue
		}
		result.WriteString(strings.ToUpper(string(word[0])))
		if len(word) > 1 {
			result.WriteString(strings.ToLower(word[1:]))
		}
	}

	return result.String()
}

func toSnakeCase(s string) string {
	if s == "" {
		return ""
	}

	if strings.ContainsAny(s, "-_") {
		return strings.ToLower(strings.ReplaceAll(s, "-", "_"))
	}

	var out strings.Builder
	runes := []rune(s)

	for i := range runes {
		r := runes[i]

		if i > 0 {
			prev := runes[i-1]

			if unicode.IsLower(prev) && unicode.IsUpper(r) {
				out.WriteByte('_')

			} else if unicode.IsUpper(prev) &&
				unicode.IsUpper(r) &&
				i+1 < len(runes) &&
				unicode.IsLower(runes[i+1]) {
				out.WriteByte('_')
			}
		}

		out.WriteRune(unicode.ToLower(r))
	}

	return out.String()
}

func splitWords(s string) []string {
	var words []string
	var currentWord strings.Builder

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			currentWord.WriteRune(r)
		} else {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		}
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

func wrapComment(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}
