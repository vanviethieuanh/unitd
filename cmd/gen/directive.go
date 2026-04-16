package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type record struct {
	Parser string `json:"parser"`
	Type   string `json:"type"`
}

func loadParserMap(f string) (map[string]string, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()

		var r record
		if err := json.Unmarshal(line, &r); err != nil {
			return result, fmt.Errorf("failed to unmarshal line %s: %w", line, err)
		}

		result[r.Parser] = r.Type
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}
