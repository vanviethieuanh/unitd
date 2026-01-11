package configs

type UnitBlock struct {
	Description string   `hcl:"description,optional" systemd:"Description"`
	After       []string `hcl:"after,optional" systemd:"After"`
	Before      []string `hcl:"before,optional" systemd:"Before"`
	Wants       []string `hcl:"wants,optional" systemd:"Wants"`
	Requires    []string `hcl:"requires,optional" systemd:"Requires"`
}

type InstallBlock struct {
	WantedBy   []string `hcl:"wanted_by,optional" systemd:"WantedBy"`
	RequiredBy []string `hcl:"required_by,optional" systemd:"RequiredBy"`
}
