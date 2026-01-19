// Generated based on man page systemd.mount of systemd

package configs

import (
	"os"
)

// MountBlock is for [Mount] systemd unit block
//
// A unit configuration file whose name ends in encodes information about a file system mount point
// controlled and supervised by systemd. This man page lists the configuration options specific to this
// unit type. See for the common options of all unit configuration files. The common configuration
// items are configured in the generic [Unit] and [Install] sections. The mount specific configuration
// options are configured in the [Mount] section. Additional options are listed in , which define the
// execution environment the program is executed in, and in , which define the way the processes are
// terminated, and in , which configure resource control settings for the processes of the service.
// Note that the options and are not useful for mount units. systemd passes two parameters to ; the
// values of and . When invoked in this way, does not read any options from , and must be run as UID 0.
// Mount units must be named after the mount point directories they control. Example: the mount point
// must be configured in a unit file . For details about the escaping logic used to convert a file
// system path to a unit name, see . Note that mount units cannot be templated, nor is possible to add
// multiple names to a mount unit by creating symlinks to its unit file. Optionally, a mount unit may
// be accompanied by an automount unit, to allow on-demand or parallelized mounting. See . Mount points
// created at runtime (independently of unit files or ) will be monitored by systemd and appear like
// any other mount unit in systemd. See description in . Some file systems have special semantics as
// API file systems for kernel-to-userspace and userspace-to-userspace interfaces. Some of them may not
// be changed via mount units, and cannot be disabled. For a longer discussion see . The command allows
// creating and units dynamically and transiently from the command line.
type MountBlock struct {
	// Directories of mount points (and any parent directories) are automatically created if needed. This
	// option specifies the file system access mode used when creating these directories. Takes an access
	// mode in octal notation. Defaults to 0755.
	DirectoryMode os.FileMode `unitd:"directory_mode,optional" systemd:"DirectoryMode"`
	// Takes a boolean argument. If true, force an unmount (in case of an unreachable NFS system). This
	// corresponds with <citerefentry
	// project='man-pages'><refentrytitle>umount</refentrytitle><manvolnum>8</manvolnum></citerefentry>'s
	// <parameter>-f</parameter> switch. Defaults to off.
	ForceUnmount bool `hcl:"force_unmount,optional" systemd:"ForceUnmount"`
	// Takes a boolean argument. If true, detach the filesystem from the filesystem hierarchy at time of
	// the unmount operation, and clean up all references to the filesystem as soon as they are not busy
	// anymore. This corresponds with <citerefentry
	// project='man-pages'><refentrytitle>umount</refentrytitle><manvolnum>8</manvolnum></citerefentry>'s
	// <parameter>-l</parameter> switch. Defaults to off.
	LazyUnmount bool `hcl:"lazy_unmount,optional" systemd:"LazyUnmount"`
	// Mount options to use when mounting. This takes a comma-separated list of options. This setting is
	// optional. Note that the usual specifier expansion is applied to this setting, literal percent
	// characters should hence be written as <literal class='specifiers'>%%.
	Options string `hcl:"options,optional" systemd:"Options"`
	// Takes a boolean argument. If false, a mount point that shall be mounted read-write but cannot be
	// mounted so is retried to be mounted read-only. If true the operation will fail immediately after the
	// read-write mount attempt did not succeed. This corresponds with <citerefentry
	// project='man-pages'><refentrytitle>mount</refentrytitle><manvolnum>8</manvolnum></citerefentry>'s
	// <parameter>-w</parameter> switch. Defaults to off.
	ReadWriteOnly bool `hcl:"read_write_only,optional" systemd:"ReadWriteOnly"`
	// Takes a boolean argument. If true, parsing of the options specified in Options= is relaxed, and
	// unknown mount options are tolerated. This corresponds with <citerefentry
	// project='man-pages'><refentrytitle>mount</refentrytitle><manvolnum>8</manvolnum></citerefentry>'s
	// <parameter>-s</parameter> switch. Defaults to off.
	SloppyOptions bool `hcl:"sloppy_options,optional" systemd:"SloppyOptions"`
	// Configures the time to wait for the mount command to finish. If a command does not exit within the
	// configured time, the mount will be considered failed and be shut down again. All commands still
	// running will be terminated forcibly via SIGTERM, and after another delay of this time with SIGKILL.
	// (See KillMode= in
	// <citerefentry><refentrytitle>systemd.kill</refentrytitle><manvolnum>5</manvolnum></citerefentry>.)
	// Takes a unit-less value in seconds, or a time span value such as "5min 20s". Pass 0 to disable the
	// timeout logic. The default value is set from DefaultTimeoutStartSec= option in
	// <citerefentry><refentrytitle>systemd-system.conf</refentrytitle><manvolnum>5</manvolnum></citerefentry>.
	TimeoutSec string `hcl:"timeout_sec,optional" systemd:"TimeoutSec"`
	// Takes a string for the file system type. See <citerefentry
	// project='man-pages'><refentrytitle>mount</refentrytitle><manvolnum>8</manvolnum></citerefentry> for
	// details. This setting is optional.
	//
	// If the type is overlay, and upperdir= or workdir= are specified as options and the directories do
	// not exist, they will be created.
	//
	Type string `hcl:"type,optional" systemd:"Type"`
	// Takes an absolute path or a fstab-style identifier of a device node, file or other resource to
	// mount. See <citerefentry
	// project='man-pages'><refentrytitle>mount</refentrytitle><manvolnum>8</manvolnum></citerefentry> for
	// details. If this refers to a device node, a dependency on the respective device unit is
	// automatically created. (See
	// <citerefentry><refentrytitle>systemd.device</refentrytitle><manvolnum>5</manvolnum></citerefentry>
	// for more information.) This option is mandatory. Note that the usual specifier expansion is applied
	// to this setting, literal percent characters should hence be written as <literal
	// class='specifiers'>%%. If this mount is a bind mount and the specified path does not exist yet it is
	// created as directory.
	What string `hcl:"what,optional" systemd:"What"`
	// Takes an absolute path of a file or directory for the mount point; in particular, the destination
	// cannot be a symbolic link. If the mount point does not exist at the time of mounting, it is created
	// as either a directory or a file. The former is the usual case; the latter is done only if this mount
	// is a bind mount and the source (What=) is not a directory. This string must be reflected in the unit
	// filename. (See above.) This option is mandatory.
	Where string `hcl:"where,optional" systemd:"Where"`
}

type Mount struct {
	Name string `hcl:"name,label"`

	Unit    UnitBlock    `hcl:"unit,block"`
	Mount   MountBlock   `hcl:"mount,block"`
	Install InstallBlock `hcl:"install,block"`
}
