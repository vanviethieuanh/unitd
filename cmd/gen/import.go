package main

import (
	"fmt"
	"sort"
	"strings"
)

type ImportSet map[string]struct{}

func NewImportSet() ImportSet {
	return make(map[string]struct{})
}

func (s ImportSet) Add(pkg string) {
	if pkg != "" {
		s[pkg] = struct{}{}
	}
}

func (s ImportSet) AddAll(pkgs []string) {
	for _, p := range pkgs {
		s.Add(p)
	}
}

func (s ImportSet) Merge(other ImportSet) {
	for p := range other {
		s[p] = struct{}{}
	}
}

func (s ImportSet) Sorted() []string {
	out := make([]string, 0, len(s))
	for p := range s {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

func writeImports(out *strings.Builder, packages []string) {
	if len(packages) == 0 {
		return
	}

	out.WriteString("import (\n")
	for _, p := range packages {
		fmt.Fprintf(out, "\t%q\n", p)
	}
	out.WriteString(")\n")
}
