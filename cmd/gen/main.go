package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

func main() {
	man := flag.String("man", "", "systemd man page name (e.g. systemd.timer)")
	outputJSON := flag.Bool("json", false, "output JSON instead of Go code")
	flag.Parse()

	if *man == "" {
		log.Fatal("--man is required (e.g. --man systemd.timer)")
	}

	fragPath := filepath.Join("out", "directives.jsonl")
	f, err := os.Open(fragPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	directives := getDirectives(f)
	if *man == "systemd.unit" {
		fmt.Print(genCommon(directives))
		return
	}

	path := filepath.Join(
		".",
		"systemd",
		"man",
		*man+".xml",
	)

	f, err = os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	unit := parseUnit(f, directives)

	if *outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(unit); err != nil {
			log.Fatal(err)
		}
	} else {
		var code strings.Builder

		fmt.Fprintf(&code, "// Generated based on man page %s of systemd\n\n", *man)
		code.WriteString("package configs\n\n")
		unitCode, imports := GenerateUnitCode(unit)
		writeImports(&code, imports.Sorted())
		code.WriteString(unitCode)

		fmt.Print(code.String())
	}
}
