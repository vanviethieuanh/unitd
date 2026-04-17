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
	var sharedMans multiFlag
	flag.Var(&sharedMans, "shared-man", "shared man pages to check for applicability (e.g. systemd.exec); may be repeated")
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

	code, err := generateUnit(*man, *manDir, *pkg, sharedMans, directives)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(code)
}

// multiFlag allows a flag to be specified multiple times.
type multiFlag []string

func (f *multiFlag) String() string { return strings.Join(*f, ", ") }
func (f *multiFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func generateUnit(man, manDir, pkg string, sharedMans []string, directives []Directive) (string, error) {
	path := filepath.Join(manDir, man+".xml")

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Determine the unit type (e.g. "service" from "systemd.service").
	unitType := strings.TrimPrefix(man, "systemd.")

	// Load descriptions from shared man pages whose <refsynopsisdiv>
	// lists this unit type (e.g. systemd.exec applies to .service).
	extraDescriptions := make(map[string]string)
	for _, extra := range sharedMans {
		extraPath := filepath.Join(manDir, extra+".xml")
		ef, err := os.Open(extraPath)
		if err != nil {
			return "", fmt.Errorf("open shared man %s: %w", extra, err)
		}
		types, err := parseApplicableTypes(ef)
		ef.Close()
		if err != nil {
			return "", fmt.Errorf("parse shared man synopsis %s: %w", extra, err)
		}

		applies := false
		for _, t := range types {
			if t == unitType {
				applies = true
				break
			}
		}
		if !applies {
			continue
		}

		log.Debug("shared man page applies", "man", extra, "unit", unitType)

		ef2, err := os.Open(extraPath)
		if err != nil {
			return "", fmt.Errorf("open shared man %s: %w", extra, err)
		}
		descs, err := parseDescriptions(ef2)
		ef2.Close()
		if err != nil {
			return "", fmt.Errorf("parse shared man %s: %w", extra, err)
		}
		for k, v := range descs {
			if _, exists := extraDescriptions[k]; !exists {
				extraDescriptions[k] = v
			}
		}
	}

	unit, err := parseUnit(f, directives, extraDescriptions)
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
