package configs

// PathBlock represents the [Path] section of a systemd unit.
//
// A unit configuration file whose name ends in encodes information about a path monitored by systemd,
// for path-based activation.
//
// This man page lists the configuration options specific to this unit type. See for the common options
// of all unit configuration files. The common configuration items are configured in the generic [Unit]
// and [Install] sections. The path specific configuration options are configured in the [Path]
// section.
//
// For each path file, a matching unit file must exist, describing the unit to activate when the path
// changes. By default, a service by the same name as the path (except for the suffix) is activated.
// Example: a path file activates a matching service . The unit to activate may be controlled by (see
// below).
//
// Internally, path units use the API to monitor file systems. Due to that, it suffers by the same
// limitations as inotify, and for example cannot be used to monitor files or directories changed by
// other machines on remote NFS file systems.
//
// When a service unit triggered by a path unit terminates (regardless whether it exited successfully
// or failed), monitored paths are checked immediately again, and the service accordingly restarted
// instantly. As protection against busy looping in this trigger/start cycle, a start rate limit is
// enforced on the service unit, see and in . Unlike other service failures, the error condition that
// the start rate limit is hit is propagated from the service unit to the path unit and causes the path
// unit to fail as well, thus ending the loop.
type PathBlock struct {
	// If MakeDirectory= is enabled, use the mode specified here to create the directories in question.
	// Takes an access mode in octal notation. Defaults to 0755.
	DirectoryMode string `hcl:"directory_mode,optional" systemd:"DirectoryMode"`
	// Defines paths to monitor for certain changes: PathExists= may be used to watch the mere existence of
	// a file or directory. If the file specified exists, the configured unit is activated. PathExistsGlob=
	// works similarly, but checks for the existence of at least one file matching the globbing pattern
	// specified. PathChanged= may be used to watch a file or directory and activate the configured unit
	// whenever it changes. It is not activated on every write to the watched file but it is activated if
	// the file which was open for writing gets closed. PathModified= is similar, but additionally it is
	// activated also on simple writes to the watched file. DirectoryNotEmpty= may be used to watch a
	// directory and activate the configured unit whenever it contains at least one file.
	//
	// The arguments of these directives must be absolute file system paths.
	//
	// Multiple directives may be combined, of the same and of different types, to watch multiple paths. If
	// the empty string is assigned to any of these options, the list of paths to watch is reset, and any
	// prior assignments of these options will not have any effect.
	//
	// If a path already exists (in case of PathExists= and PathExistsGlob=) or a directory already is not
	// empty (in case of DirectoryNotEmpty=) at the time the path unit is activated, then the configured
	// unit is immediately activated as well. Something similar does not apply to PathChanged= and
	// PathModified=.
	//
	// If the path itself or any of the containing directories are not accessible, systemd will watch for
	// permission changes and notice that conditions are satisfied when permissions allow that.
	//
	// Note that files whose name starts with a dot (i.e. hidden files) are generally ignored when
	// monitoring these paths.
	//
	DirectoryNotEmpty string `hcl:"directory_not_empty,optional" systemd:"DirectoryNotEmpty"`
	// Takes a boolean argument. If true, the directories to watch are created before watching. This option
	// is ignored for PathExists= settings. Defaults to false.
	MakeDirectory bool `hcl:"make_directory,optional" systemd:"MakeDirectory"`
	// Defines paths to monitor for certain changes: PathExists= may be used to watch the mere existence of
	// a file or directory. If the file specified exists, the configured unit is activated. PathExistsGlob=
	// works similarly, but checks for the existence of at least one file matching the globbing pattern
	// specified. PathChanged= may be used to watch a file or directory and activate the configured unit
	// whenever it changes. It is not activated on every write to the watched file but it is activated if
	// the file which was open for writing gets closed. PathModified= is similar, but additionally it is
	// activated also on simple writes to the watched file. DirectoryNotEmpty= may be used to watch a
	// directory and activate the configured unit whenever it contains at least one file.
	//
	// The arguments of these directives must be absolute file system paths.
	//
	// Multiple directives may be combined, of the same and of different types, to watch multiple paths. If
	// the empty string is assigned to any of these options, the list of paths to watch is reset, and any
	// prior assignments of these options will not have any effect.
	//
	// If a path already exists (in case of PathExists= and PathExistsGlob=) or a directory already is not
	// empty (in case of DirectoryNotEmpty=) at the time the path unit is activated, then the configured
	// unit is immediately activated as well. Something similar does not apply to PathChanged= and
	// PathModified=.
	//
	// If the path itself or any of the containing directories are not accessible, systemd will watch for
	// permission changes and notice that conditions are satisfied when permissions allow that.
	//
	// Note that files whose name starts with a dot (i.e. hidden files) are generally ignored when
	// monitoring these paths.
	//
	PathChanged string `hcl:"path_changed,optional" systemd:"PathChanged"`
	// Defines paths to monitor for certain changes: PathExists= may be used to watch the mere existence of
	// a file or directory. If the file specified exists, the configured unit is activated. PathExistsGlob=
	// works similarly, but checks for the existence of at least one file matching the globbing pattern
	// specified. PathChanged= may be used to watch a file or directory and activate the configured unit
	// whenever it changes. It is not activated on every write to the watched file but it is activated if
	// the file which was open for writing gets closed. PathModified= is similar, but additionally it is
	// activated also on simple writes to the watched file. DirectoryNotEmpty= may be used to watch a
	// directory and activate the configured unit whenever it contains at least one file.
	//
	// The arguments of these directives must be absolute file system paths.
	//
	// Multiple directives may be combined, of the same and of different types, to watch multiple paths. If
	// the empty string is assigned to any of these options, the list of paths to watch is reset, and any
	// prior assignments of these options will not have any effect.
	//
	// If a path already exists (in case of PathExists= and PathExistsGlob=) or a directory already is not
	// empty (in case of DirectoryNotEmpty=) at the time the path unit is activated, then the configured
	// unit is immediately activated as well. Something similar does not apply to PathChanged= and
	// PathModified=.
	//
	// If the path itself or any of the containing directories are not accessible, systemd will watch for
	// permission changes and notice that conditions are satisfied when permissions allow that.
	//
	// Note that files whose name starts with a dot (i.e. hidden files) are generally ignored when
	// monitoring these paths.
	//
	PathExists string `hcl:"path_exists,optional" systemd:"PathExists"`
	// Defines paths to monitor for certain changes: PathExists= may be used to watch the mere existence of
	// a file or directory. If the file specified exists, the configured unit is activated. PathExistsGlob=
	// works similarly, but checks for the existence of at least one file matching the globbing pattern
	// specified. PathChanged= may be used to watch a file or directory and activate the configured unit
	// whenever it changes. It is not activated on every write to the watched file but it is activated if
	// the file which was open for writing gets closed. PathModified= is similar, but additionally it is
	// activated also on simple writes to the watched file. DirectoryNotEmpty= may be used to watch a
	// directory and activate the configured unit whenever it contains at least one file.
	//
	// The arguments of these directives must be absolute file system paths.
	//
	// Multiple directives may be combined, of the same and of different types, to watch multiple paths. If
	// the empty string is assigned to any of these options, the list of paths to watch is reset, and any
	// prior assignments of these options will not have any effect.
	//
	// If a path already exists (in case of PathExists= and PathExistsGlob=) or a directory already is not
	// empty (in case of DirectoryNotEmpty=) at the time the path unit is activated, then the configured
	// unit is immediately activated as well. Something similar does not apply to PathChanged= and
	// PathModified=.
	//
	// If the path itself or any of the containing directories are not accessible, systemd will watch for
	// permission changes and notice that conditions are satisfied when permissions allow that.
	//
	// Note that files whose name starts with a dot (i.e. hidden files) are generally ignored when
	// monitoring these paths.
	//
	PathExistsGlob string `hcl:"path_exists_glob,optional" systemd:"PathExistsGlob"`
	// Defines paths to monitor for certain changes: PathExists= may be used to watch the mere existence of
	// a file or directory. If the file specified exists, the configured unit is activated. PathExistsGlob=
	// works similarly, but checks for the existence of at least one file matching the globbing pattern
	// specified. PathChanged= may be used to watch a file or directory and activate the configured unit
	// whenever it changes. It is not activated on every write to the watched file but it is activated if
	// the file which was open for writing gets closed. PathModified= is similar, but additionally it is
	// activated also on simple writes to the watched file. DirectoryNotEmpty= may be used to watch a
	// directory and activate the configured unit whenever it contains at least one file.
	//
	// The arguments of these directives must be absolute file system paths.
	//
	// Multiple directives may be combined, of the same and of different types, to watch multiple paths. If
	// the empty string is assigned to any of these options, the list of paths to watch is reset, and any
	// prior assignments of these options will not have any effect.
	//
	// If a path already exists (in case of PathExists= and PathExistsGlob=) or a directory already is not
	// empty (in case of DirectoryNotEmpty=) at the time the path unit is activated, then the configured
	// unit is immediately activated as well. Something similar does not apply to PathChanged= and
	// PathModified=.
	//
	// If the path itself or any of the containing directories are not accessible, systemd will watch for
	// permission changes and notice that conditions are satisfied when permissions allow that.
	//
	// Note that files whose name starts with a dot (i.e. hidden files) are generally ignored when
	// monitoring these paths.
	//
	PathModified string `hcl:"path_modified,optional" systemd:"PathModified"`
	// Configures a limit on how often this path unit may be activated within a specific time interval. The
	// TriggerLimitIntervalSec= may be used to configure the length of the time interval in the usual time
	// units us, ms, s, min, h, … and defaults to 2s. See
	// <citerefentry><refentrytitle>systemd.time</refentrytitle><manvolnum>7</manvolnum></citerefentry> for
	// details on the various time units understood. The TriggerLimitBurst= setting takes a positive
	// integer value and specifies the number of permitted activations per time interval, and defaults to
	// 200. Set either to 0 to disable any form of trigger rate limiting. If the limit is hit, the unit is
	// placed into a failure mode, and will not watch the paths anymore until restarted. Note that this
	// limit is enforced before the service activation is enqueued.
	TriggerLimitBurst int `hcl:"trigger_limit_burst,optional" systemd:"TriggerLimitBurst"`
	// Configures a limit on how often this path unit may be activated within a specific time interval. The
	// TriggerLimitIntervalSec= may be used to configure the length of the time interval in the usual time
	// units us, ms, s, min, h, … and defaults to 2s. See
	// <citerefentry><refentrytitle>systemd.time</refentrytitle><manvolnum>7</manvolnum></citerefentry> for
	// details on the various time units understood. The TriggerLimitBurst= setting takes a positive
	// integer value and specifies the number of permitted activations per time interval, and defaults to
	// 200. Set either to 0 to disable any form of trigger rate limiting. If the limit is hit, the unit is
	// placed into a failure mode, and will not watch the paths anymore until restarted. Note that this
	// limit is enforced before the service activation is enqueued.
	TriggerLimitIntervalSec int `hcl:"trigger_limit_interval_sec,optional" systemd:"TriggerLimitIntervalSec"`
	// The unit to activate when any of the configured paths changes. The argument is a unit name, whose
	// suffix is not .path. If not specified, this value defaults to a service that has the same name as
	// the path unit, except for the suffix. (See above.) It is recommended that the unit name that is
	// activated and the unit name of the path unit are named identical, except for the suffix.
	Unit string `hcl:"unit,optional" systemd:"Unit"`
}

type Path struct {
	Name string `hcl:"name,label"`

	Unit    UnitBlock    `hcl:"unit,block"`
	Path    PathBlock    `hcl:"path,block"`
	Install InstallBlock `hcl:"install,block"`
}
