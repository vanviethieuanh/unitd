package configs

// SwapBlock represents the [Swap] section of a systemd unit.
//
// A unit configuration file whose name ends in encodes information about a swap device or file for
// memory paging controlled and supervised by systemd.
//
// This man page lists the configuration options specific to this unit type. See for the common options
// of all unit configuration files. The common configuration items are configured in the generic [Unit]
// and [Install] sections. The swap specific configuration options are configured in the [Swap]
// section.
//
// Additional options are listed in , which define the execution environment the program is executed
// in, in , which define the way these processes are terminated, and in , which configure resource
// control settings for these processes of the unit.
//
// Swap units must be named after the devices or files they control. Example: the swap device must be
// configured in a unit file . For details about the escaping logic used to convert a file system path
// to a unit name, see . Note that swap units cannot be templated, nor is possible to add multiple
// names to a swap unit by creating additional symlinks to it.
//
// Note that swap support on Linux is privileged, swap units are hence only available in the system
// service manager (and root's user service manager), but not in unprivileged user's service manager.
type SwapBlock struct {
	// May contain an option string for the swap device. This may be used for controlling discard options
	// among other functionality, if the swap backing device supports the discard or trim operation. (See
	// <citerefentry
	// project='man-pages'><refentrytitle>swapon</refentrytitle><manvolnum>8</manvolnum></citerefentry> for
	// more information.) Note that the usual specifier expansion is applied to this setting, literal
	// percent characters should hence be written as %%.
	Options string `hcl:"options,optional" systemd:"Options"`
	// Swap priority to use when activating the swap device or file. This takes an integer. This setting is
	// optional and ignored when the priority is set by pri= in the Options= key.
	Priority int `hcl:"priority,optional" systemd:"Priority"`
	// Configures the time to wait for the swapon command to finish. If a command does not exit within the
	// configured time, the swap will be considered failed and be shut down again. All commands still
	// running will be terminated forcibly via SIGTERM, and after another delay of this time with SIGKILL.
	// (See KillMode= in
	// <citerefentry><refentrytitle>systemd.kill</refentrytitle><manvolnum>5</manvolnum></citerefentry>.)
	// Takes a unit-less value in seconds, or a time span value such as "5min 20s". Pass 0 to disable the
	// timeout logic. Defaults to DefaultTimeoutStartSec= from the manager configuration file (see
	// <citerefentry><refentrytitle>systemd-system.conf</refentrytitle><manvolnum>5</manvolnum></citerefentry>).
	TimeoutSec int `hcl:"timeout_sec,optional" systemd:"TimeoutSec"`
	// Takes an absolute path or a fstab-style identifier of a device node or file to use for paging. See
	// <citerefentry
	// project='man-pages'><refentrytitle>swapon</refentrytitle><manvolnum>8</manvolnum></citerefentry> for
	// details. If this refers to a device node, a dependency on the respective device unit is
	// automatically created. (See
	// <citerefentry><refentrytitle>systemd.device</refentrytitle><manvolnum>5</manvolnum></citerefentry>
	// for more information.) If this refers to a file, a dependency on the respective mount unit is
	// automatically created. (See
	// <citerefentry><refentrytitle>systemd.mount</refentrytitle><manvolnum>5</manvolnum></citerefentry>
	// for more information.) This option is mandatory. Note that the usual specifier expansion is applied
	// to this setting, literal percent characters should hence be written as <literal
	// class='specifiers'>%%.
	What string `hcl:"what,optional" systemd:"What"`
}

type Swap struct {
	Name string `hcl:"name,label"`

	Unit    UnitBlock    `hcl:"unit,block"`
	Swap    SwapBlock    `hcl:"swap,block"`
	Install InstallBlock `hcl:"install,block"`
}
