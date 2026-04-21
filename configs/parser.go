package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// configFileSchema is the top-level schema for a unitd configuration file.
var configFileSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "service", LabelNames: []string{"name"}},
		{Type: "instance", LabelNames: []string{"name"}},
	},
}

// DecodeFile parses an HCL file and decodes it into a Config.
//
// Multi-phase decode:
//  1. Pre-scan service blocks for labels, template, for_each
//  2. Pre-scan instance blocks for labels, template expr, instances
//  3. Build partial EvalContext (builtins + services)
//  4. Resolve instance template expressions
//  5. Build full EvalContext (+ instances)
//  6. Decode blocks individually — services with for_each are expanded
//     (one Service per variant, each decoded with its own each.key/each.value)
func DecodeFile(path string) (*Config, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse %s: %s", path, diags.Error())
	}

	body := file.Body

	// Phase 1: Pre-scan service blocks.
	serviceMetas, err := ExtractServiceMeta(body)
	if err != nil {
		return nil, fmt.Errorf("extract services from %s: %w", path, err)
	}

	// Phase 2: Pre-scan instance blocks.
	instanceMetas, err := ExtractInstanceMeta(body)
	if err != nil {
		return nil, fmt.Errorf("extract instances from %s: %w", path, err)
	}

	// Phase 3: Build partial context (builtins + services, no instances yet).
	partialCtx := BuildEvalContext(DefaultKnownUnits, serviceMetas, nil)

	// Phase 4: Resolve instance template expressions.
	resolved, err := ResolveInstances(partialCtx, instanceMetas)
	if err != nil {
		return nil, fmt.Errorf("resolve instances in %s: %w", path, err)
	}

	// Phase 5: Build full base context.
	baseCtx := BuildEvalContext(DefaultKnownUnits, serviceMetas, resolved)

	// Phase 6: Decode blocks individually.
	content, _, diags := body.PartialContent(configFileSchema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("read blocks from %s: %s", path, diags.Error())
	}

	metaIndex := make(map[string]ServiceMeta, len(serviceMetas))
	for _, m := range serviceMetas {
		metaIndex[m.Name] = m
	}

	var config Config
	for _, block := range content.Blocks {
		switch block.Type {
		case "service":
			name := block.Labels[0]
			meta := metaIndex[name]

			if len(meta.ForEach) > 0 {
				// Expand: decode once per variant with each.key/each.value
				for key, value := range meta.ForEach {
					ctx := WithEachVars(baseCtx, key, value)
					var svc Service
					diags := gohcl.DecodeBody(block.Body, ctx, &svc)
					if diags.HasErrors() {
						return nil, fmt.Errorf("decode service %q variant %q in %s: %s", name, key, path, diags.Error())
					}
					svc.Name = name
					svc.ForEach = map[string]string{key: value}
					config.Services = append(config.Services, svc)
				}
			} else {
				var svc Service
				diags := gohcl.DecodeBody(block.Body, baseCtx, &svc)
				if diags.HasErrors() {
					return nil, fmt.Errorf("decode service %q in %s: %s", name, path, diags.Error())
				}
				svc.Name = name
				config.Services = append(config.Services, svc)
			}

		case "instance":
			var inst Instance
			diags := gohcl.DecodeBody(block.Body, baseCtx, &inst)
			if diags.HasErrors() {
				return nil, fmt.Errorf("decode instance %q in %s: %s", block.Labels[0], path, diags.Error())
			}
			inst.Name = block.Labels[0]
			config.Instances = append(config.Instances, inst)
		}
	}

	return &config, nil
}
