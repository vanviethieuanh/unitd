// Generated based on man page systemd.scope of systemd

package configs

// ScopeBlock is for [Scope] systemd unit block
//
// Scope units are not configured via unit configuration files, but are only created programmatically
// using the bus interfaces of systemd. They are named similar to filenames. A unit whose name ends in
// refers to a scope unit. Scopes units manage a set of system processes. Unlike service units, scope
// units manage externally created processes, and do not fork off processes on its own. The main
// purpose of scope units is grouping worker processes of a system service for organization and for
// managing resources. may be used to easily launch a command in a new scope unit from the command
// line. See the for an introduction on how to make use of scope units from programs. Note that, unlike
// service units, scope units have no "main" process: all processes in the scope are equivalent. The
// lifecycle of the scope unit is thus not bound to the lifetime of one specific process, but to the
// existence of at least one process in the scope. This also means that the exit statuses of these
// processes are not relevant for the scope unit failure state. Scope units may still enter a failure
// state, for example due to resource exhaustion or stop timeouts being reached, but not due to
// programs inside of them terminating uncleanly. Since processes managed as scope units generally
// remain children of the original process that forked them off, it is also the job of that process to
// collect their exit statuses and act on them as needed.
type ScopeBlock struct {
	OOMPolicy string `hcl:"oom_policy,optional" systemd:"OOMPolicy"`
	// Configures a maximum time for the scope to run. If this is used and the scope has been active for
	// longer than the specified time it is terminated and put into a failure state. Pass infinity (the
	// default) to configure no runtime limit.
	RuntimeMaxSec int `hcl:"runtime_max_sec,optional" systemd:"RuntimeMaxSec"`
	// This option modifies RuntimeMaxSec= by increasing the maximum runtime by an evenly distributed
	// duration between 0 and the specified value (in seconds). If RuntimeMaxSec= is unspecified, then this
	// feature will be disabled.
	RuntimeRandomizedExtraSec int `hcl:"runtime_randomized_extra_sec,optional" systemd:"RuntimeRandomizedExtraSec"`
	TimeoutStopSec            int `hcl:"timeout_stop_sec,optional" systemd:"TimeoutStopSec"`
}

type Scope struct {
	Name string `hcl:"name,label"`

	Unit    UnitBlock    `hcl:"unit,block"`
	Scope   ScopeBlock   `hcl:"scope,block"`
	Install InstallBlock `hcl:"install,block"`
}
