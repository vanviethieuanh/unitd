package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

func main() {
	man := flag.String("man", "", "systemd man page name (e.g. systemd.timer)")
	gperfDir := flag.String("gperf-dir", filepath.Join("tmp", "gperf"), "directory containing gperf JSONL files")
	manDir := flag.String("man-dir", filepath.Join("tmp", "systemd", "man"), "directory containing systemd XML man pages")
	pkg := flag.String("package", "configs", "Go package name for generated files")
	logLevel := flag.String("log-level", "info", "log level (debug, info, warn, error, fatal)")
	flag.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal("invalid log level", "level", *logLevel)
	}
	log.SetLevel(level)

	if *man == "" {
		log.Fatal("--man is required (e.g. --man systemd.timer)")
	}

	directives, err := LoadAllDirectives(*gperfDir)
	if err != nil {
		log.Fatal(err)
	}

	if *man == "systemd.unit" {
		fmt.Print(genCommon(*pkg, directives))
		return
	}

	code, err := generateUnit(*man, *manDir, *pkg, directives)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(code)
}

func generateUnit(man, manDir, pkg string, directives []Directive) (string, error) {
	path := filepath.Join(manDir, man+".xml")

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	unit, err := parseUnit(f, directives)
	if err != nil {
		return "", err
	}

	var code strings.Builder
	fmt.Fprintf(&code, "// Generated based on man page %s of systemd\n\n", man)
	fmt.Fprintf(&code, "package %s\n\n", pkg)

	unitCode, imports := GenerateUnitCode(unit)
	writeImports(&code, imports.Sorted())
	code.WriteString(unitCode)

	return code.String(), nil
}
