package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// Instance declares instantiations of a template unit.
type Instance struct {
	Name      string   `hcl:"name,label"`
	Template  string   `hcl:"template"`
	Instances []string `hcl:"instances"`
}

// InstanceMeta holds pre-scanned metadata about an instance block.
type InstanceMeta struct {
	Name         string
	TemplateExpr hcl.Expression
	Instances    []string
}

// InstanceResolved holds a fully resolved instance block.
type InstanceResolved struct {
	Name         string
	TemplateName string // e.g. "worker-queue@.service"
	Instances    []string
}

// ExtractInstanceMeta pre-scans instance blocks for their expressions.
func ExtractInstanceMeta(body hcl.Body) ([]InstanceMeta, error) {
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "instance", LabelNames: []string{"name"}},
		},
	}

	content, _, diags := body.PartialContent(schema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("partial content: %s", diags.Error())
	}

	innerSchema := &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "template"},
			{Name: "instances"},
		},
	}

	var result []InstanceMeta
	for _, block := range content.Blocks {
		if len(block.Labels) == 0 {
			continue
		}

		meta := InstanceMeta{Name: block.Labels[0]}

		inner, _, diags := block.Body.PartialContent(innerSchema)
		if diags.HasErrors() {
			return nil, fmt.Errorf("instance %q: %s", meta.Name, diags.Error())
		}

		if attr, ok := inner.Attributes["template"]; ok {
			meta.TemplateExpr = attr.Expr
		}

		if attr, ok := inner.Attributes["instances"]; ok {
			val, diags := attr.Expr.Value(nil)
			if !diags.HasErrors() && val.CanIterateElements() {
				for it := val.ElementIterator(); it.Next(); {
					_, v := it.Element()
					if v.Type() == cty.String {
						meta.Instances = append(meta.Instances, v.AsString())
					}
				}
			}
		}

		result = append(result, meta)
	}

	return result, nil
}

// ResolveInstances evaluates instance template expressions using a partial
// EvalContext and produces resolved instance metadata.
func ResolveInstances(ctx *hcl.EvalContext, metas []InstanceMeta) ([]InstanceResolved, error) {
	var result []InstanceResolved
	for _, m := range metas {
		if m.TemplateExpr == nil {
			return nil, fmt.Errorf("instance %q: missing template attribute", m.Name)
		}

		val, diags := m.TemplateExpr.Value(ctx)
		if diags.HasErrors() {
			return nil, fmt.Errorf("instance %q template: %s", m.Name, diags.Error())
		}
		if val.Type() != cty.String {
			return nil, fmt.Errorf("instance %q template: expected string, got %s", m.Name, val.Type().FriendlyName())
		}

		result = append(result, InstanceResolved{
			Name:         m.Name,
			TemplateName: val.AsString(),
			Instances:    m.Instances,
		})
	}
	return result, nil
}

// InstanceVars builds the cty variables for the instance namespace.
func InstanceVars(instances []InstanceResolved) map[string]cty.Value {
	instMap := make(map[string]cty.Value, len(instances))
	for _, inst := range instances {
		idMap := make(map[string]cty.Value, len(inst.Instances))
		for _, id := range inst.Instances {
			idMap[id] = cty.StringVal(InstanceUnitName(inst.TemplateName, id))
		}
		instMap[inst.Name] = cty.ObjectVal(idMap)
	}
	return instMap
}
