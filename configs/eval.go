package configs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// BuiltinVars builds the cty variables for the builtin namespace.
func BuiltinVars(knownUnits []KnownUnit) map[string]cty.Value {
	byType := make(map[string]map[string]cty.Value)
	for _, u := range knownUnits {
		if _, ok := byType[u.UnitType]; !ok {
			byType[u.UnitType] = make(map[string]cty.Value)
		}
		byType[u.UnitType][sanitizeHCLIdent(u.Name)] = cty.StringVal(u.Name)
	}

	builtinObj := make(map[string]cty.Value, len(byType))
	for unitType, names := range byType {
		builtinObj[unitType] = cty.ObjectVal(names)
	}
	return builtinObj
}

// SelfVars returns the cty variables for the self namespace (template specifiers).
func SelfVars() map[string]cty.Value {
	return map[string]cty.Value{
		"instance":           cty.StringVal("%i"),
		"instance_unescaped": cty.StringVal("%I"),
	}
}

// EachVars returns the cty variables for the each namespace (for_each iteration).
func EachVars(key, value string) map[string]cty.Value {
	return map[string]cty.Value{
		"key":   cty.StringVal(key),
		"value": cty.StringVal(value),
	}
}

// BuildEvalContext assembles an HCL EvalContext from per-block-type variable maps.
func BuildEvalContext(
	knownUnits []KnownUnit,
	services []ServiceMeta,
	instances []InstanceResolved,
) *hcl.EvalContext {
	vars := make(map[string]cty.Value)

	vars["builtin"] = cty.ObjectVal(BuiltinVars(knownUnits))

	if svcVars := ServiceVars(services); len(svcVars) > 0 {
		vars["service"] = cty.ObjectVal(svcVars)
	}

	if instVars := InstanceVars(instances); len(instVars) > 0 {
		vars["instance"] = cty.ObjectVal(instVars)
	}

	vars["self"] = cty.ObjectVal(SelfVars())

	return &hcl.EvalContext{
		Variables: vars,
	}
}

// WithEachVars returns a copy of the EvalContext with the each namespace added.
func WithEachVars(base *hcl.EvalContext, key, value string) *hcl.EvalContext {
	vars := make(map[string]cty.Value, len(base.Variables)+1)
	for k, v := range base.Variables {
		vars[k] = v
	}
	vars["each"] = cty.ObjectVal(EachVars(key, value))
	return &hcl.EvalContext{
		Variables: vars,
	}
}
