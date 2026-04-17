// Generated based on man page systemd.scope of systemd

package configs

import (
	"syscall"
)

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
	AllowedCPUs             string   `hcl:"allowed_cp_us,optional" systemd:"AllowedCPUs"`
	AllowedMemoryNodes      string   `hcl:"allowed_memory_nodes,optional" systemd:"AllowedMemoryNodes"`
	BPFProgram              []string `hcl:"bpf_program,optional" systemd:"BPFProgram"`
	BindNetworkInterface    []string `hcl:"bind_network_interface,optional" systemd:"BindNetworkInterface"`
	CPUPressureThresholdSec int      `hcl:"cpu_pressure_threshold_sec,optional" systemd:"CPUPressureThresholdSec"`
	CPUPressureWatch        string   `hcl:"cpu_pressure_watch,optional" systemd:"CPUPressureWatch"`
	CPUQuota                string   `hcl:"cpu_quota,optional" systemd:"CPUQuota"`
	CPUQuotaPeriodSec       int      `hcl:"cpu_quota_period_sec,optional" systemd:"CPUQuotaPeriodSec"`
	CPUWeight               uint64   `hcl:"cpu_weight,optional" systemd:"CPUWeight"`
	CoredumpReceive         bool     `hcl:"coredump_receive,optional" systemd:"CoredumpReceive"`
	Delegate                string   `hcl:"delegate,optional" systemd:"Delegate"`
	DelegateSubgroup        string   `hcl:"delegate_subgroup,optional" systemd:"DelegateSubgroup"`
	DeviceAllow             []string `hcl:"device_allow,optional" systemd:"DeviceAllow"`
	DevicePolicy            string   `hcl:"device_policy,optional" systemd:"DevicePolicy"`
	DisableControllers      []string `hcl:"disable_controllers,optional" systemd:"DisableControllers"`
	// Specifies which signal to send to remaining processes after a timeout if SendSIGKILL= is enabled.
	// The signal configured here should be one that is not typically caught and processed by services
	// (SIGTERM is not suitable). Developers can find it useful to use this to generate a coredump to
	// troubleshoot why a service did not terminate upon receiving the initial SIGTERM signal. This can be
	// achieved by configuring LimitCORE= and setting FinalKillSignal= to either SIGQUIT or SIGABRT.
	// Defaults to SIGKILL.
	FinalKillSignal          syscall.Signal `unitd:"final_kill_signal,optional" systemd:"FinalKillSignal"`
	IOAccounting             bool           `hcl:"io_accounting,optional" systemd:"IOAccounting"`
	IODeviceLatencyTargetSec []string       `hcl:"io_device_latency_target_sec,optional" systemd:"IODeviceLatencyTargetSec"`
	IODeviceWeight           []string       `hcl:"io_device_weight,optional" systemd:"IODeviceWeight"`
	IOPressureThresholdSec   int            `hcl:"io_pressure_threshold_sec,optional" systemd:"IOPressureThresholdSec"`
	IOPressureWatch          string         `hcl:"io_pressure_watch,optional" systemd:"IOPressureWatch"`
	IOReadBandwidthMax       []string       `hcl:"io_read_bandwidth_max,optional" systemd:"IOReadBandwidthMax"`
	IOReadIOPSMax            []string       `hcl:"io_read_iops_max,optional" systemd:"IOReadIOPSMax"`
	IOWeight                 uint64         `hcl:"io_weight,optional" systemd:"IOWeight"`
	IOWriteBandwidthMax      []string       `hcl:"io_write_bandwidth_max,optional" systemd:"IOWriteBandwidthMax"`
	IOWriteIOPSMax           []string       `hcl:"io_write_iops_max,optional" systemd:"IOWriteIOPSMax"`
	IPAccounting             bool           `hcl:"ip_accounting,optional" systemd:"IPAccounting"`
	IPAddressAllow           []string       `hcl:"ip_address_allow,optional" systemd:"IPAddressAllow"`
	IPAddressDeny            []string       `hcl:"ip_address_deny,optional" systemd:"IPAddressDeny"`
	IPEgressFilterPath       []string       `hcl:"ip_egress_filter_path,optional" systemd:"IPEgressFilterPath"`
	IPIngressFilterPath      []string       `hcl:"ip_ingress_filter_path,optional" systemd:"IPIngressFilterPath"`
	// Specifies how processes of this unit shall be killed. One of control-group, mixed, process, none.
	//
	// If set to control-group, all remaining processes in the control group of this unit will be killed on
	// unit stop (for services: after the stop command is executed, as configured with ExecStop=). If set
	// to mixed, the SIGTERM signal (see below) is sent to the main process while the subsequent SIGKILL
	// signal (see below) is sent to all remaining processes of the unit's control group. If set to
	// process, only the main process itself is killed (not recommended!). If set to none, no process is
	// killed (strongly recommended against!). In this case, only the stop command will be executed on unit
	// stop, but no process will be killed otherwise. Processes remaining alive after stop are left in
	// their control group and the control group continues to exist after stop unless empty.
	//
	// Note that it is not recommended to set KillMode= to process or even none, as this allows processes
	// to escape the service manager's lifecycle and resource management, and to remain running even while
	// their service is considered stopped and is assumed to not consume any resources.
	//
	// Processes will first be terminated via SIGTERM (unless the signal to send is changed via KillSignal=
	// or RestartKillSignal=). Optionally, this is immediately followed by a SIGHUP (if enabled with
	// SendSIGHUP=). If processes still remain after: <itemizedlist> <listitem><para>the main process of a
	// unit has exited (applies to KillMode=: mixed)</para></listitem> <listitem><para>the delay configured
	// via the TimeoutStopSec= has passed (applies to KillMode=: control-group, mixed,
	// process)</para></listitem> </itemizedlist> the termination request is repeated with the SIGKILL
	// signal or the signal specified via FinalKillSignal= (unless this is disabled via the SendSIGKILL=
	// option). See
	// <citerefentry><refentrytitle>kill</refentrytitle><manvolnum>2</manvolnum></citerefentry> for more
	// information.
	//
	// Defaults to control-group.
	//
	KillMode string `hcl:"kill_mode,optional" systemd:"KillMode"`
	// Specifies which signal to use when stopping a service. This controls the signal that is sent as
	// first step of shutting down a unit (see above), and is usually followed by SIGKILL (see above and
	// below). For a list of valid signals, see <citerefentry
	// project='man-pages'><refentrytitle>signal</refentrytitle><manvolnum>7</manvolnum></citerefentry>.
	// Defaults to SIGTERM.
	//
	// Note that, right after sending the signal specified in this setting, systemd will always send
	// SIGCONT, to ensure that even suspended tasks can be terminated cleanly.
	//
	KillSignal                          syscall.Signal `unitd:"kill_signal,optional" systemd:"KillSignal"`
	ManagedOOMMemoryPressure            string         `hcl:"managed_oom_memory_pressure,optional" systemd:"ManagedOOMMemoryPressure"`
	ManagedOOMMemoryPressureDurationSec int            `hcl:"managed_oom_memory_pressure_duration_sec,optional" systemd:"ManagedOOMMemoryPressureDurationSec"`
	ManagedOOMMemoryPressureLimit       string         `hcl:"managed_oom_memory_pressure_limit,optional" systemd:"ManagedOOMMemoryPressureLimit"`
	ManagedOOMPreference                string         `hcl:"managed_oom_preference,optional" systemd:"ManagedOOMPreference"`
	ManagedOOMSwap                      string         `hcl:"managed_oom_swap,optional" systemd:"ManagedOOMSwap"`
	MemoryAccounting                    bool           `hcl:"memory_accounting,optional" systemd:"MemoryAccounting"`
	MemoryHigh                          string         `hcl:"memory_high,optional" systemd:"MemoryHigh"`
	MemoryLow                           string         `hcl:"memory_low,optional" systemd:"MemoryLow"`
	MemoryMax                           string         `hcl:"memory_max,optional" systemd:"MemoryMax"`
	MemoryMin                           string         `hcl:"memory_min,optional" systemd:"MemoryMin"`
	MemoryPressureThresholdSec          int            `hcl:"memory_pressure_threshold_sec,optional" systemd:"MemoryPressureThresholdSec"`
	MemoryPressureWatch                 string         `hcl:"memory_pressure_watch,optional" systemd:"MemoryPressureWatch"`
	MemorySwapMax                       string         `hcl:"memory_swap_max,optional" systemd:"MemorySwapMax"`
	MemoryZSwapMax                      string         `hcl:"memory_z_swap_max,optional" systemd:"MemoryZSwapMax"`
	MemoryZSwapWriteback                bool           `hcl:"memory_z_swap_writeback,optional" systemd:"MemoryZSwapWriteback"`
	NFTSet                              []string       `hcl:"nft_set,optional" systemd:"NFTSet"`
	OOMPolicy                           string         `hcl:"oom_policy,optional" systemd:"OOMPolicy"`
	// Specifies which signal to use when restarting a service. The same as KillSignal= described above,
	// with the exception that this setting is used in a restart job. Not set by default, and the value of
	// KillSignal= is used.
	RestartKillSignal         syscall.Signal `unitd:"restart_kill_signal,optional" systemd:"RestartKillSignal"`
	RestrictNetworkInterfaces []string       `hcl:"restrict_network_interfaces,optional" systemd:"RestrictNetworkInterfaces"`
	// Configures a maximum time for the scope to run. If this is used and the scope has been active for
	// longer than the specified time it is terminated and put into a failure state. Pass infinity (the
	// default) to configure no runtime limit.
	RuntimeMaxSec int `hcl:"runtime_max_sec,optional" systemd:"RuntimeMaxSec"`
	// This option modifies RuntimeMaxSec= by increasing the maximum runtime by an evenly distributed
	// duration between 0 and the specified value (in seconds). If RuntimeMaxSec= is unspecified, then this
	// feature will be disabled.
	RuntimeRandomizedExtraSec int `hcl:"runtime_randomized_extra_sec,optional" systemd:"RuntimeRandomizedExtraSec"`
	// Specifies whether to send SIGHUP to remaining processes immediately after sending the signal
	// configured with KillSignal=. This is useful to indicate to shells and shell-like programs that their
	// connection has been severed. Takes a boolean value. Defaults to no.
	SendSIGHUP bool `hcl:"send_sighup,optional" systemd:"SendSIGHUP"`
	// Specifies whether to send SIGKILL (or the signal specified by FinalKillSignal=) to remaining
	// processes after a timeout, if the normal shutdown procedure left processes of the service around.
	// When disabled, a KillMode= of control-group or mixed service will not restart if processes from
	// prior services exist within the control group. Takes a boolean value. Defaults to yes.
	SendSIGKILL               bool     `hcl:"send_sigkill,optional" systemd:"SendSIGKILL"`
	Slice                     string   `hcl:"slice,optional" systemd:"Slice"`
	SocketBindAllow           []string `hcl:"socket_bind_allow,optional" systemd:"SocketBindAllow"`
	SocketBindDeny            []string `hcl:"socket_bind_deny,optional" systemd:"SocketBindDeny"`
	StartupAllowedCPUs        string   `hcl:"startup_allowed_cp_us,optional" systemd:"StartupAllowedCPUs"`
	StartupAllowedMemoryNodes string   `hcl:"startup_allowed_memory_nodes,optional" systemd:"StartupAllowedMemoryNodes"`
	StartupCPUWeight          uint64   `hcl:"startup_cpu_weight,optional" systemd:"StartupCPUWeight"`
	StartupIOWeight           uint64   `hcl:"startup_io_weight,optional" systemd:"StartupIOWeight"`
	StartupMemoryHigh         string   `hcl:"startup_memory_high,optional" systemd:"StartupMemoryHigh"`
	StartupMemoryLow          string   `hcl:"startup_memory_low,optional" systemd:"StartupMemoryLow"`
	StartupMemoryMax          string   `hcl:"startup_memory_max,optional" systemd:"StartupMemoryMax"`
	StartupMemorySwapMax      string   `hcl:"startup_memory_swap_max,optional" systemd:"StartupMemorySwapMax"`
	StartupMemoryZSwapMax     string   `hcl:"startup_memory_z_swap_max,optional" systemd:"StartupMemoryZSwapMax"`
	TasksAccounting           bool     `hcl:"tasks_accounting,optional" systemd:"TasksAccounting"`
	TasksMax                  string   `hcl:"tasks_max,optional" systemd:"TasksMax"`
	TimeoutStopSec            int      `hcl:"timeout_stop_sec,optional" systemd:"TimeoutStopSec"`
	// Specifies which signal to use to terminate the service when the watchdog timeout expires (enabled
	// through WatchdogSec=). Defaults to SIGABRT.
	WatchdogSignal syscall.Signal `unitd:"watchdog_signal,optional" systemd:"WatchdogSignal"`
}

type Scope struct {
	Name string `hcl:"name,label"`

	Unit    UnitBlock    `hcl:"unit,block"`
	Scope   ScopeBlock   `hcl:"scope,block"`
	Install InstallBlock `hcl:"install,block"`
}
