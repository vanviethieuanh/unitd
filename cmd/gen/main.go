package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func main() {
	man := flag.String("man", "", "systemd man page name (e.g. systemd.timer)")
	outputJSON := flag.Bool("json", false, "output JSON instead of Go code")
	flag.Parse()

	if *man == "" {
		log.Fatal("--man is required (e.g. --man systemd.timer)")
	}

	propTypes := getDirectives()
	if *man == "systemd.unit" {
		fmt.Print(genCommon(propTypes))
		return
	}

	path := filepath.Join(
		".",
		"systemd-man",
		"man",
		*man+".xml",
	)

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	unit := parseUnit(f)

	if *outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(unit); err != nil {
			log.Fatal(err)
		}
	} else {
		code := gen(&unit, propTypes)
		fmt.Print(code)
	}
}
