package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

type gperfRecord struct {
	System   string `json:"system"`
	Section  string `json:"section"`
	Property string `json:"property"`
	Parser   string `json:"parser"`
}

// nullParserTypes provides type strings for properties whose parser is "NULL"
// in gperf. These are Install section entries that systemd handles outside of
// the normal config_parse_* machinery.
var nullParserTypes = map[string]string{
	"Alias":           "STRING",
	"WantedBy":        "UNIT [...]",
	"RequiredBy":      "UNIT [...]",
	"UpheldBy":        "UNIT [...]",
	"Also":            "UNIT [...]",
	"DefaultInstance": "STRING",
}

// loadGperfDirectives reads all gperf_*.jsonl files in dir (skipping
// parser_type.jsonl) and converts each record to a Directive using parserMap
// to resolve the type string from the parser name.
func loadGperfDirectives(dir string, parserMap map[string]string) []Directive {
	files, err := filepath.Glob(filepath.Join(dir, "gperf_*.jsonl"))
	if err != nil {
		log.Fatal(err)
	}

	var result []Directive

	for _, path := range files {
		f, err := os.Open(path)
		if err != nil {
			log.Fatalf("open %s: %v", path, err)
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Bytes()

			var gr gperfRecord
			if err := json.Unmarshal(line, &gr); err != nil {
				log.Fatalf("unmarshal %s: %v", path, err)
			}

			typeStr, ok := parserMap[gr.Parser]
			if !ok {
				if gr.Parser == "NULL" {
					typeStr, ok = nullParserTypes[gr.Property]
				}
				if !ok {
					log.Warnf("unknown parser %q in %s — skipping", gr.Parser, filepath.Base(path))
					continue
				}
			}

			// parsers mapped to empty type string are intentionally unsupported
			typeStr = strings.TrimSpace(typeStr)
			if typeStr == "" || typeStr == "NOTSUPPORTED" {
				continue
			}

			parsedType, err := ParseTypeExpr(typeStr)
			if err != nil {
				log.Warnf("parse type %q (parser %s): %v — skipping", typeStr, gr.Parser, err)
				continue
			}

			goType, deps := parsedType.toGoType()

			result = append(result, Directive{
				Identifier: DirectiveIdentifier{
					Section: gr.Section,
					Key:     gr.Property,
				},
				Type:       goType,
				System:     gr.System,
				Deps:       deps,
				NativeType: len(deps) == 0,
			})
		}

		f.Close()

		if err := scanner.Err(); err != nil {
			log.Fatalf("scan %s: %v", path, err)
		}
	}

	return result
}

// loadParserMapDefault loads the parser_type.jsonl from gperfDir.
func loadParserMapDefault(gperfDir string) map[string]string {
	pm, err := loadParserMap(filepath.Join(gperfDir, "parser_type.jsonl"))
	if err != nil {
		log.Fatalf("load parser map: %v", err)
	}
	return pm
}

// LoadAllDirectives is the single entry-point: load parser map and all gperf
// records from gperfDir, returning the merged []Directive slice.
func LoadAllDirectives(gperfDir string) []Directive {
	pm := loadParserMapDefault(gperfDir)
	return loadGperfDirectives(gperfDir, pm)
}
