package configs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// ServiceMeta holds pre-scanned metadata about a service block.
type ServiceMeta struct {
	Name     string
	Template bool
	ForEach  map[string]string // key → value from for_each map
}

// UnitFilenames returns the systemd unit filenames this service produces.
//
//	no template, no for_each  →  ["nginx.service"]
//	no template, for_each     →  ["worker-queue.service", "worker-email.service"]
//	template, no for_each     →  ["worker@.service"]
//	template + for_each       →  ["worker-queue@.service", "worker-email@.service"]
func (s *Service) UnitFilenames() []string {
	switch {
	case s.Template && len(s.ForEach) > 0:
		names := make([]string, 0, len(s.ForEach))
		for v := range s.ForEach {
			names = append(names, TemplateUnitName(s.Name, v, "service"))
		}
		return names
	case len(s.ForEach) > 0:
		names := make([]string, 0, len(s.ForEach))
		for v := range s.ForEach {
			names = append(names, s.Name+"-"+v+".service")
		}
		return names
	case s.Template:
		return []string{TemplateUnitName(s.Name, "", "service")}
	default:
		return []string{s.Name + ".service"}
	}
}

// Encode converts a Service to a systemd .service unit file string.
func (s *Service) Encode() (string, error) {
	var b strings.Builder

	sections := []struct {
		name string
		data any
	}{
		{"Unit", s.Unit},
		{"Service", s.Service},
		{"Install", s.Install},
	}

	for _, sec := range sections {
		entries, err := EncodeSystemdSection(sec.data)
		if err != nil {
			return "", err
		}
		if len(entries) > 0 {
			writeSection(&b, sec.name, entries)
		}
	}

	return strings.TrimSpace(b.String()) + "\n", nil
}

// ExtractServiceMeta pre-scans service blocks for template/for_each metadata.
func ExtractServiceMeta(body hcl.Body) ([]ServiceMeta, error) {
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "service", LabelNames: []string{"name"}},
		},
	}

	content, _, diags := body.PartialContent(schema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("partial content: %s", diags.Error())
	}

	innerSchema := &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "template"},
			{Name: "for_each"},
		},
	}

	var result []ServiceMeta
	for _, block := range content.Blocks {
		if len(block.Labels) == 0 {
			continue
		}

		meta := ServiceMeta{Name: block.Labels[0]}

		inner, _, diags := block.Body.PartialContent(innerSchema)
		if diags.HasErrors() {
			return nil, fmt.Errorf("service %q: %s", meta.Name, diags.Error())
		}

		if attr, ok := inner.Attributes["template"]; ok {
			val, diags := attr.Expr.Value(nil)
			if !diags.HasErrors() && val.Type() == cty.Bool {
				meta.Template = val.True()
			}
		}

		if attr, ok := inner.Attributes["for_each"]; ok {
			val, diags := attr.Expr.Value(nil)
			if !diags.HasErrors() && (val.Type().IsObjectType() || val.Type().IsMapType()) {
				meta.ForEach = make(map[string]string)
				for it := val.ElementIterator(); it.Next(); {
					k, v := it.Element()
					if k.Type() == cty.String {
						vs := ""
						if v.Type() == cty.String {
							vs = v.AsString()
						}
						meta.ForEach[k.AsString()] = vs
					}
				}
			}
		}

		result = append(result, meta)
	}

	return result, nil
}

// ServiceVars builds the cty variables for the service namespace.
func ServiceVars(services []ServiceMeta) map[string]cty.Value {
	svcMap := make(map[string]cty.Value, len(services))
	for _, s := range services {
		switch {
		case s.Template && len(s.ForEach) > 0:
			variantMap := make(map[string]cty.Value, len(s.ForEach))
			for k := range s.ForEach {
				variantMap[k] = cty.StringVal(TemplateUnitName(s.Name, k, "service"))
			}
			svcMap[s.Name] = cty.ObjectVal(variantMap)
		case len(s.ForEach) > 0:
			variantMap := make(map[string]cty.Value, len(s.ForEach))
			for k := range s.ForEach {
				variantMap[k] = cty.StringVal(fmt.Sprintf("%s-%s.service", s.Name, k))
			}
			svcMap[s.Name] = cty.ObjectVal(variantMap)
		case s.Template:
			svcMap[s.Name] = cty.StringVal(TemplateUnitName(s.Name, "", "service"))
		default:
			svcMap[s.Name] = cty.StringVal(s.Name + ".service")
		}
	}
	return svcMap
}
