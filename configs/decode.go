package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// DecodeFile parses an HCL file and decodes it into a Config.
//
// It pre-scans block labels to build user-defined unit references,
// then decodes with an EvalContext that resolves both builtin and
// user-defined unit names.
func DecodeFile(path string) (*Config, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse %s: %s", path, diags.Error())
	}

	// Pre-scan: extract all block types and labels.
	userBlocks, err := extractBlockLabels(file.Body)
	if err != nil {
		return nil, fmt.Errorf("extract labels from %s: %w", path, err)
	}

	ctx := BuildEvalContext(DefaultKnownUnits, userBlocks)

	var config Config
	diags = gohcl.DecodeBody(file.Body, ctx, &config)
	if diags.HasErrors() {
		return nil, fmt.Errorf("decode %s: %s", path, diags.Error())
	}

	return &config, nil
}

// extractBlockLabels scans an HCL body for top-level blocks and collects
// their labels, grouped by block type. This is used to build the
// user-defined unit reference namespace (e.g. service.nginx → "nginx.service").
func extractBlockLabels(body hcl.Body) (map[string][]string, error) {
	// Define the schema of known block types.
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "service", LabelNames: []string{"name"}},
			{Type: "timer", LabelNames: []string{"name"}},
			{Type: "socket", LabelNames: []string{"name"}},
			{Type: "path", LabelNames: []string{"name"}},
			{Type: "mount", LabelNames: []string{"name"}},
			{Type: "automount", LabelNames: []string{"name"}},
			{Type: "swap", LabelNames: []string{"name"}},
			{Type: "target", LabelNames: []string{"name"}},
			{Type: "slice", LabelNames: []string{"name"}},
			{Type: "scope", LabelNames: []string{"name"}},
			{Type: "device", LabelNames: []string{"name"}},
		},
	}

	content, _, diags := body.PartialContent(schema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("partial content: %s", diags.Error())
	}

	result := make(map[string][]string)
	for _, block := range content.Blocks {
		if len(block.Labels) > 0 {
			result[block.Type] = append(result[block.Type], block.Labels[0])
		}
	}

	return result, nil
}
