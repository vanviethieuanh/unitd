// Package configs define all the blocks and syntax for this DSL.
package configs

import "fmt"

type Config struct {
	Services []Service `hcl:"service,block"`
}

func (c *Config) Validate() error {
	seen := make(map[string]struct{})

	for _, svc := range c.Services {
		if _, ok := seen[svc.Name]; ok {
			return fmt.Errorf(
				"duplicate service %q defined more than once",
				svc.Name,
			)
		}
		seen[svc.Name] = struct{}{}
	}

	return nil
}
