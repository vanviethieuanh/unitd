package main

import (
	"strings"
	"unicode"
)

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
