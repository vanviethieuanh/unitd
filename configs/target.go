package configs

// Target represents the Target configuration file of a systemd unit.
// A separate [Target] section does not exist.
//
// A unit configuration file whose name ends in encodes information about a target unit of systemd.
// Target units are used to group units and to set synchronization points for ordering dependencies
// with other unit files.
//
// This unit type has no specific options. See for the common options of all unit configuration files.
// The common configuration items are configured in the generic [Unit] and [Install] sections. A
// separate [Target] section does not exist, since no target-specific options may be configured.
//
// Target units do not offer any additional functionality on top of the generic functionality provided
// by units. They merely group units, allowing a single target name to be used in and settings to
// establish a dependency on a set of units defined by the target, and in and settings to establish
// ordering. Targets establish standardized names for synchronization points during boot and shutdown.
// Importantly, see for examples and descriptions of standard systemd targets.
//
// Target units provide a more flexible replacement for SysV runlevels in the classic SysV init system.
//
// Note that a target unit file must not be empty, lest it be considered a masked unit. It is
// recommended to provide a [Unit] section which includes informative and options.
type Target struct {
	Name string `hcl:"name,label"`

	Unit    UnitBlock    `hcl:"unit,block"`
	Install InstallBlock `hcl:"install,block"`
}
