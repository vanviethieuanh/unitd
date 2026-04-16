package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

var gperfDir = filepath.Join("tmp", "gperf")

func main() {
	man := flag.String("man", "", "systemd man page name (e.g. systemd.timer)")
	flag.Parse()

	if *man == "" {
		log.Fatal("--man is required (e.g. --man systemd.timer)")
	}

	directives := LoadAllDirectives(gperfDir)
	if *man == "systemd.unit" {
		fmt.Print(genCommon(directives))
		return
	}

	path := filepath.Join(
		".",
		"tmp",
		"systemd",
		"man",
		*man+".xml",
	)

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	unit := parseUnit(f, directives)

	var code strings.Builder

	fmt.Fprintf(&code, "// Generated based on man page %s of systemd\n\n", *man)
	code.WriteString("package configs\n\n")
	unitCode, imports := GenerateUnitCode(unit)
	writeImports(&code, imports.Sorted())
	code.WriteString(unitCode)

	fmt.Print(code.String())
}
