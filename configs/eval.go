package configs

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// sanitizeHCLIdent converts a systemd unit name stem to a valid HCL identifier.
// Hyphens become underscores, "@" becomes "_at_", leading "-" becomes "root".
func sanitizeHCLIdent(name string) string {
	// Strip the type suffix (e.g. ".target", ".service")
	dot := strings.LastIndex(name, ".")
	if dot < 0 {
		return name
	}
	stem := name[:dot]

	// Special case: "-.slice" / "-.mount" → "root"
	if stem == "-" {
		return "root"
	}

	stem = strings.ReplaceAll(stem, "-", "_")
	stem = strings.ReplaceAll(stem, "@", "_at_")

	return stem
}

// BuildEvalContext creates an HCL EvalContext that resolves unit references.
//
// Builtin units (from DefaultKnownUnits) are accessible as:
//
//	builtin.<type>.<sanitized_name>  →  "original-name.type"
//
// User-defined blocks are accessible as:
//
//	<block_type>.<label>  →  "label.block_type"
func BuildEvalContext(knownUnits []KnownUnit, userBlocks map[string][]string) *hcl.EvalContext {
	vars := make(map[string]cty.Value)

	// Build builtin namespace: builtin.<type>.<name> = "unit-name.type"
	byType := make(map[string]map[string]cty.Value)
	for _, u := range knownUnits {
		if u.IsTemplate {
			continue
		}
		if _, ok := byType[u.UnitType]; !ok {
			byType[u.UnitType] = make(map[string]cty.Value)
		}
		ident := sanitizeHCLIdent(u.Name)
		byType[u.UnitType][ident] = cty.StringVal(u.Name)
	}

	builtinObj := make(map[string]cty.Value, len(byType))
	for unitType, names := range byType {
		builtinObj[unitType] = cty.ObjectVal(names)
	}
	vars["builtin"] = cty.ObjectVal(builtinObj)

	// Build user-defined block namespaces: <type>.<label> = "label.type"
	for blockType, labels := range userBlocks {
		labelMap := make(map[string]cty.Value, len(labels))
		for _, label := range labels {
			labelMap[label] = cty.StringVal(label + "." + blockType)
		}
		vars[blockType] = cty.ObjectVal(labelMap)
	}

	return &hcl.EvalContext{
		Variables: vars,
	}
}
