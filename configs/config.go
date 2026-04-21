// Package configs define all the blocks and syntax for this DSL.
package configs

import "fmt"

type Config struct {
	Services  []Service  `hcl:"service,block"`
	Instances []Instance `hcl:"instance,block"`
}

func (c *Config) Validate() error {
	type svcKey struct {
		name    string
		variant string
	}
	seen := make(map[svcKey]struct{})

	for _, svc := range c.Services {
		variant := ""
		for k := range svc.ForEach {
			variant = k
			break // expanded services have exactly one entry
		}
		key := svcKey{svc.Name, variant}
		if _, ok := seen[key]; ok {
			if variant != "" {
				return fmt.Errorf("duplicate service %q variant %q", svc.Name, variant)
			}
			return fmt.Errorf(
				"duplicate service %q defined more than once",
				svc.Name,
			)
		}
		seen[key] = struct{}{}
	}

	for _, inst := range c.Instances {
		if inst.Template == "" {
			return fmt.Errorf("instance %q is missing template", inst.Name)
		}
		if len(inst.Instances) == 0 {
			return fmt.Errorf("instance %q has no instances", inst.Name)
		}
	}

	return nil
}
