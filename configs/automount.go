// Generated based on man page systemd.automount of systemd

package configs

import (
	"os"
)

// AutomountBlock is for [Automount] systemd unit block
//
// A unit configuration file whose name ends in encodes information about a file system automount point
// controlled and supervised by systemd. Automount units may be used to implement on-demand mounting as
// well as parallelized mounting of file systems. This man page lists the configuration options
// specific to this unit type. See for the common options of all unit configuration files. The common
// configuration items are configured in the generic [Unit] and [Install] sections. The automount
// specific configuration options are configured in the [Automount] section. Automount units must be
// named after the automount directories they control. Example: the automount point must be configured
// in a unit file . For details about the escaping logic used to convert a file system path to a unit
// name see . Note that automount units cannot be templated, nor is it possible to add multiple names
// to an automount unit by creating symlinks to its unit file. For each automount unit file a matching
// mount unit file (see for details) must exist which is activated when the automount path is accessed.
// Example: if an automount unit is active and the user accesses the mount unit will be activated. Note
// that automount units are separate from the mount itself, so you should not set or for mount
// dependencies here. For example, you should not set or similar on network filesystems. Doing so may
// result in an ordering cycle. Note that automount support on Linux is privileged, automount units are
// hence only available in the system service manager (and root's user service manager), but not in
// unprivileged users' service managers. Note that automount units should not be nested. (The
// establishment of the inner automount point would unconditionally pin the outer mount point,
// defeating its purpose.)
type AutomountBlock struct {
	// Directories of automount points (and any parent directories) are automatically created if needed.
	// This option specifies the file system access mode used when creating these directories. Takes an
	// access mode in octal notation. Defaults to 0755.
	DirectoryMode os.FileMode `unitd:"directory_mode,optional" systemd:"DirectoryMode"`
	// Extra mount options to use when creating the autofs mountpoint. This takes a comma-separated list of
	// options. This setting is optional. Note that the usual specifier expansion is applied to this
	// setting, literal percent characters should hence be written as <literal class='specifiers'>%%.
	ExtraOptions string `hcl:"extra_options,optional" systemd:"ExtraOptions"`
	// Configures an idle timeout. Once the mount has been idle for the specified time, systemd will
	// attempt to unmount. Takes a unit-less value in seconds, or a time span value such as "5min 20s".
	// Pass 0 to disable the timeout logic. The timeout is disabled by default.
	TimeoutIdleSec string `hcl:"timeout_idle_sec,optional" systemd:"TimeoutIdleSec"`
	// Takes an absolute path of a directory of the automount point. If the automount point does not exist
	// at time that the automount point is installed, it is created. This string must be reflected in the
	// unit filename. (See above.) This option is mandatory.
	Where string `hcl:"where,optional" systemd:"Where"`
}

type Automount struct {
	Name string `hcl:"name,label"`

	Unit      UnitBlock      `hcl:"unit,block"`
	Automount AutomountBlock `hcl:"automount,block"`
	Install   InstallBlock   `hcl:"install,block"`
}
