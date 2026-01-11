package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/vanviethieuanh/unitd/configs"
)

func main() {
	outDir := "./example/transplied"
	_ = os.MkdirAll(outDir, 0o755)

	var config configs.Config

	err := hclsimple.DecodeFile("./example/src/services.hcl", nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	for _, service := range config.Services {
		u, err := service.Encode()
		if err != nil {
			panic(err)
		}

		s := u.ToString()

		path := filepath.Join(outDir, service.Name+".service")

		if err := os.WriteFile(path, []byte(s), 0o644); err != nil {
			panic(err)
		}

		fmt.Println("wrote", path)
	}
}
