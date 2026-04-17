package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type RefEntry struct {
	XMLName     xml.Name       `xml:"refentry"`
	Meta        RefMeta        `xml:"refmeta"`
	NameDiv     RefNameDiv     `xml:"refnamediv"`
	SynopsisDiv RefSynopsisDiv `xml:"refsynopsisdiv"`
	Sections    []RefSect1     `xml:"refsect1"`
}

type RefSynopsisDiv struct {
	Para SynopsisPara `xml:"para"`
}

type SynopsisPara struct {
	FileNames []string `xml:"filename"`
}

type RefMeta struct {
	Title  string `xml:"refentrytitle"`
	Volnum string `xml:"manvolnum"`
}

type RefNameDiv struct {
	Name    string `xml:"refname"`
	Purpose string `xml:"refpurpose"`
}

type RefSect1 struct {
	Title        string       `xml:"title"`
	Para         []string     `xml:"para"`
	VariableList VariableList `xml:"variablelist"`
}

type VariableList struct {
	Entries []VarListEntry `xml:"varlistentry"`
}

type VarListEntry struct {
	Terms    []string `xml:"term>varname"`
	ListItem ListItem `xml:"listitem"`
}

type ListItem struct {
	Para  []Para `xml:"para"`
	Table *Table `xml:"table"`
}

type Para struct {
	Content string `xml:",innerxml"`
}

type Table struct {
	TGroup TGroup `xml:"tgroup"`
}

type TGroup struct {
	Body TBody `xml:"tbody"`
}

type TBody struct {
	Rows []Row `xml:"row"`
}

type Row struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Content string `xml:",innerxml"`
}

func parseUnit(r io.Reader, directivesList []Directive, extraDescriptions map[string]string) (*Unit, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	var refEntry RefEntry
	if err := xml.Unmarshal(data, &refEntry); err != nil {
		return nil, fmt.Errorf("parsing XML: %w", err)
	}

	unitName := parseName(refEntry)

	var options []Directive
	parsedOptions := parseOptions(refEntry)
	seen := make(map[string]struct{})
	for _, d := range directivesList {
		if !strings.EqualFold(d.Identifier.Section, unitName) {
			continue
		}

		key := d.Identifier.Key
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		if desc, ok := parsedOptions[key]; ok {
			d.Description = desc
		} else if desc, ok := extraDescriptions[key]; ok {
			d.Description = desc
		}
		options = append(options, d)
	}

	return &Unit{
		Name:        parseName(refEntry),
		Title:       parseTitle(refEntry),
		Purpose:     parsePurpose(refEntry),
		Description: parseDescription(refEntry),
		Options:     options,
	}, nil
}

func parseName(refEntry RefEntry) string {
	name := refEntry.NameDiv.Name
	if after, ok := strings.CutPrefix(name, "systemd."); ok {
		name = after
	}
	return name
}

func parseTitle(refEntry RefEntry) string {
	return refEntry.Meta.Title
}

func parsePurpose(refEntry RefEntry) string {
	return refEntry.NameDiv.Purpose
}

func parseDescription(refEntry RefEntry) string {
	for _, section := range refEntry.Sections {
		if section.Title == "Description" {
			return joinCleanedTexts(section.Para)
		}
	}
	return ""
}

func parseOptions(refEntry RefEntry) map[string]string {
	options := make(map[string]string)

	for _, section := range refEntry.Sections {
		if section.Title == "Description" || len(section.VariableList.Entries) == 0 {
			continue
		}
		for _, entry := range section.VariableList.Entries {
			if entry.ListItem.Table != nil {
				for _, term := range entry.Terms {
					cleanTerm := strings.TrimSpace(strings.TrimSuffix(term, "="))
					desc := extractTableDescription(&entry, cleanTerm)
					if desc != "" {
						options[cleanTerm] = desc
					}
				}
			} else {
				for _, term := range entry.Terms {
					cleanTerm := strings.TrimSpace(strings.TrimSuffix(term, "="))
					desc := extractFullDescription(&entry)
					if desc != "" {
						options[cleanTerm] = desc
					}
				}
			}
		}
	}

	return options
}

func extractTableDescription(entry *VarListEntry, term string) string {
	if entry.ListItem.Table == nil {
		return ""
	}

	cleanTerm := strings.TrimSuffix(term, "=")

	for _, row := range entry.ListItem.Table.TGroup.Body.Rows {
		if len(row.Entries) >= 2 {
			firstEntry := cleanText(row.Entries[0].Content)
			firstEntry = strings.TrimPrefix(firstEntry, "<varname>")
			firstEntry = strings.TrimSuffix(firstEntry, "</varname>")
			firstEntry = strings.TrimSuffix(firstEntry, "=")

			if strings.Contains(firstEntry, cleanTerm) {
				return cleanText(row.Entries[1].Content)
			}
		}
	}

	return ""
}

func extractFullDescription(entry *VarListEntry) string {
	texts := make([]string, len(entry.ListItem.Para))
	for i, p := range entry.ListItem.Para {
		texts[i] = p.Content
	}
	return joinCleanedTexts(texts)
}

// joinCleanedTexts cleans each text and joins non-empty results with double newlines.
func joinCleanedTexts(texts []string) string {
	var parts []string
	for _, t := range texts {
		cleaned := cleanText(t)
		if cleaned != "" {
			parts = append(parts, cleaned)
		}
	}
	return strings.Join(parts, "\n\n")
}

var xmlTagRe = regexp.MustCompile(`</?(?:literal|filename|varname|replaceable|option|constant|command|emphasis)>`)

// parseApplicableTypes extracts unit type names from <refsynopsisdiv>.
// For example, systemd.kill.xml lists .service, .socket, .mount, .swap, .scope
// and this returns ["service", "socket", "mount", "swap", "scope"].
func parseApplicableTypes(r io.Reader) ([]string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	var refEntry RefEntry
	if err := xml.Unmarshal(data, &refEntry); err != nil {
		return nil, fmt.Errorf("parsing XML: %w", err)
	}

	var types []string
	for _, fn := range refEntry.SynopsisDiv.Para.FileNames {
		// fn is like "service.service" — extract extension after the dot
		if idx := strings.LastIndex(fn, "."); idx >= 0 {
			types = append(types, fn[idx+1:])
		}
	}
	return types, nil
}

// parseDescriptions extracts option key->description mappings from an XML man page.
func parseDescriptions(r io.Reader) (map[string]string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	var refEntry RefEntry
	if err := xml.Unmarshal(data, &refEntry); err != nil {
		return nil, fmt.Errorf("parsing XML: %w", err)
	}

	return parseOptions(refEntry), nil
}

func cleanText(text string) string {
	text = strings.TrimSpace(text)
	text = xmlTagRe.ReplaceAllString(text, "")

	// Normalize whitespace: convert all whitespace sequences to single spaces
	text = strings.Join(strings.Fields(text), " ")

	return text
}
