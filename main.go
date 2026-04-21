package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/vanviethieuanh/unitd/configs"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: unitd <src.hcl> <outdir>\n")
		os.Exit(1)
	}

	srcFile := os.Args[1]
	outDir := os.Args[2]

	_ = os.MkdirAll(outDir, 0o755)

	config, err := configs.DecodeFile(srcFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	for _, svc := range config.Services {
		content, err := svc.Encode()
		if err != nil {
			log.Fatalf("Failed to encode %s: %s", svc.Name, err)
		}

		for _, filename := range svc.UnitFilenames() {
			path := filepath.Join(outDir, filename)
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				log.Fatalf("Failed to write %s: %s", path, err)
			}
			fmt.Println("wrote", path)
		}
	}
}
