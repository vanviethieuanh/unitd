// Generated based on man page systemd.slice of systemd

package configs

// SliceBlock is for [Slice] systemd unit block
//
// A unit configuration file whose name ends in encodes information about a slice unit. A slice unit is
// a concept for hierarchically managing resources of a group of processes. This management is
// performed by creating a node in the Linux Control Group (cgroup) tree. Units that manage processes
// (primarily scope and service units) may be assigned to a specific slice. For each slice, certain
// resource limits may be set that apply to all processes of all units contained in that slice. Slices
// are organized hierarchically in a tree. The name of the slice encodes the location in the tree. The
// name consists of a dash-separated series of names, which describes the path to the slice from the
// root slice. The root slice is named . Example: is a slice that is located within , which in turn is
// located in the root slice . Note that slice units cannot be templated, nor is possible to add
// multiple names to a slice unit by creating additional symlinks to its unit file. By default, service
// and scope units are placed in , virtual machines and containers registered with are found in , and
// user sessions handled by in . See for more information. See for the common options of all unit
// configuration files. The common configuration items are configured in the generic [Unit] and
// [Install] sections. The slice specific configuration options are configured in the [Slice] section.
// Currently, only generic resource control settings as described in are allowed. See the for an
// introduction on how to make use of slice units from programs.
type SliceBlock struct {
	AllowedCPUs             string   `hcl:"allowed_cp_us,optional" systemd:"AllowedCPUs"`
	AllowedMemoryNodes      string   `hcl:"allowed_memory_nodes,optional" systemd:"AllowedMemoryNodes"`
	BPFProgram              []string `hcl:"bpf_program,optional" systemd:"BPFProgram"`
	BindNetworkInterface    []string `hcl:"bind_network_interface,optional" systemd:"BindNetworkInterface"`
	CPUPressureThresholdSec int      `hcl:"cpu_pressure_threshold_sec,optional" systemd:"CPUPressureThresholdSec"`
	CPUPressureWatch        string   `hcl:"cpu_pressure_watch,optional" systemd:"CPUPressureWatch"`
	CPUQuota                string   `hcl:"cpu_quota,optional" systemd:"CPUQuota"`
	CPUQuotaPeriodSec       int      `hcl:"cpu_quota_period_sec,optional" systemd:"CPUQuotaPeriodSec"`
	CPUWeight               uint64   `hcl:"cpu_weight,optional" systemd:"CPUWeight"`
	// Configures a hard and a soft limit on the maximum number of units assigned to this slice (or any
	// descendent slices) that may be active at the same time. If the hard limit is reached no further
	// units associated with the slice may be activated, and their activation will fail with an error. If
	// the soft limit is reached any further requested activation of units will be queued, but no immediate
	// error is generated. The queued activation job will remain queued until the number of concurrent
	// active units within the slice is below the limit again.
	//
	// If the special value infinity is specified, no concurrency limit is enforced. This is the default.
	//
	// Note that if multiple start jobs are queued for units, and all their dependencies are fulfilled
	// they'll be processed in an order that is dependent on the unit type, the CPU weight (for unit types
	// that know the concept, such as services), the nice level (similar), and finally in alphabetical
	// order by the unit name. This may be used to influence dispatching order when using
	// ConcurrencySoftMax= to pace concurrency within a slice unit.
	//
	// Note that these options have a hierarchial effect: a limit set for a slice unit will apply to both
	// the units immediately within the slice, but also all units further down the slice tree. Also note
	// that each sub-slice unit counts as one unit each too, and thus when choosing a limit for a slice
	// hierarchy the limit must provide room for both the payload units (i.e. services, mounts, …) and
	// structural units (i.e. slice units), if any are defined.
	//
	ConcurrencyHardMax string `hcl:"concurrency_hard_max,optional" systemd:"ConcurrencyHardMax"`
	// Configures a hard and a soft limit on the maximum number of units assigned to this slice (or any
	// descendent slices) that may be active at the same time. If the hard limit is reached no further
	// units associated with the slice may be activated, and their activation will fail with an error. If
	// the soft limit is reached any further requested activation of units will be queued, but no immediate
	// error is generated. The queued activation job will remain queued until the number of concurrent
	// active units within the slice is below the limit again.
	//
	// If the special value infinity is specified, no concurrency limit is enforced. This is the default.
	//
	// Note that if multiple start jobs are queued for units, and all their dependencies are fulfilled
	// they'll be processed in an order that is dependent on the unit type, the CPU weight (for unit types
	// that know the concept, such as services), the nice level (similar), and finally in alphabetical
	// order by the unit name. This may be used to influence dispatching order when using
	// ConcurrencySoftMax= to pace concurrency within a slice unit.
	//
	// Note that these options have a hierarchial effect: a limit set for a slice unit will apply to both
	// the units immediately within the slice, but also all units further down the slice tree. Also note
	// that each sub-slice unit counts as one unit each too, and thus when choosing a limit for a slice
	// hierarchy the limit must provide room for both the payload units (i.e. services, mounts, …) and
	// structural units (i.e. slice units), if any are defined.
	//
	ConcurrencySoftMax                  string   `hcl:"concurrency_soft_max,optional" systemd:"ConcurrencySoftMax"`
	CoredumpReceive                     bool     `hcl:"coredump_receive,optional" systemd:"CoredumpReceive"`
	Delegate                            string   `hcl:"delegate,optional" systemd:"Delegate"`
	DelegateSubgroup                    string   `hcl:"delegate_subgroup,optional" systemd:"DelegateSubgroup"`
	DeviceAllow                         []string `hcl:"device_allow,optional" systemd:"DeviceAllow"`
	DevicePolicy                        string   `hcl:"device_policy,optional" systemd:"DevicePolicy"`
	DisableControllers                  []string `hcl:"disable_controllers,optional" systemd:"DisableControllers"`
	IOAccounting                        bool     `hcl:"io_accounting,optional" systemd:"IOAccounting"`
	IODeviceLatencyTargetSec            []string `hcl:"io_device_latency_target_sec,optional" systemd:"IODeviceLatencyTargetSec"`
	IODeviceWeight                      []string `hcl:"io_device_weight,optional" systemd:"IODeviceWeight"`
	IOPressureThresholdSec              int      `hcl:"io_pressure_threshold_sec,optional" systemd:"IOPressureThresholdSec"`
	IOPressureWatch                     string   `hcl:"io_pressure_watch,optional" systemd:"IOPressureWatch"`
	IOReadBandwidthMax                  []string `hcl:"io_read_bandwidth_max,optional" systemd:"IOReadBandwidthMax"`
	IOReadIOPSMax                       []string `hcl:"io_read_iops_max,optional" systemd:"IOReadIOPSMax"`
	IOWeight                            uint64   `hcl:"io_weight,optional" systemd:"IOWeight"`
	IOWriteBandwidthMax                 []string `hcl:"io_write_bandwidth_max,optional" systemd:"IOWriteBandwidthMax"`
	IOWriteIOPSMax                      []string `hcl:"io_write_iops_max,optional" systemd:"IOWriteIOPSMax"`
	IPAccounting                        bool     `hcl:"ip_accounting,optional" systemd:"IPAccounting"`
	IPAddressAllow                      []string `hcl:"ip_address_allow,optional" systemd:"IPAddressAllow"`
	IPAddressDeny                       []string `hcl:"ip_address_deny,optional" systemd:"IPAddressDeny"`
	IPEgressFilterPath                  []string `hcl:"ip_egress_filter_path,optional" systemd:"IPEgressFilterPath"`
	IPIngressFilterPath                 []string `hcl:"ip_ingress_filter_path,optional" systemd:"IPIngressFilterPath"`
	ManagedOOMMemoryPressure            string   `hcl:"managed_oom_memory_pressure,optional" systemd:"ManagedOOMMemoryPressure"`
	ManagedOOMMemoryPressureDurationSec int      `hcl:"managed_oom_memory_pressure_duration_sec,optional" systemd:"ManagedOOMMemoryPressureDurationSec"`
	ManagedOOMMemoryPressureLimit       string   `hcl:"managed_oom_memory_pressure_limit,optional" systemd:"ManagedOOMMemoryPressureLimit"`
	ManagedOOMPreference                string   `hcl:"managed_oom_preference,optional" systemd:"ManagedOOMPreference"`
	ManagedOOMSwap                      string   `hcl:"managed_oom_swap,optional" systemd:"ManagedOOMSwap"`
	MemoryAccounting                    bool     `hcl:"memory_accounting,optional" systemd:"MemoryAccounting"`
	MemoryHigh                          string   `hcl:"memory_high,optional" systemd:"MemoryHigh"`
	MemoryLow                           string   `hcl:"memory_low,optional" systemd:"MemoryLow"`
	MemoryMax                           string   `hcl:"memory_max,optional" systemd:"MemoryMax"`
	MemoryMin                           string   `hcl:"memory_min,optional" systemd:"MemoryMin"`
	MemoryPressureThresholdSec          int      `hcl:"memory_pressure_threshold_sec,optional" systemd:"MemoryPressureThresholdSec"`
	MemoryPressureWatch                 string   `hcl:"memory_pressure_watch,optional" systemd:"MemoryPressureWatch"`
	MemorySwapMax                       string   `hcl:"memory_swap_max,optional" systemd:"MemorySwapMax"`
	MemoryZSwapMax                      string   `hcl:"memory_z_swap_max,optional" systemd:"MemoryZSwapMax"`
	MemoryZSwapWriteback                bool     `hcl:"memory_z_swap_writeback,optional" systemd:"MemoryZSwapWriteback"`
	NFTSet                              []string `hcl:"nft_set,optional" systemd:"NFTSet"`
	RestrictNetworkInterfaces           []string `hcl:"restrict_network_interfaces,optional" systemd:"RestrictNetworkInterfaces"`
	Slice                               string   `hcl:"slice,optional" systemd:"Slice"`
	SocketBindAllow                     []string `hcl:"socket_bind_allow,optional" systemd:"SocketBindAllow"`
	SocketBindDeny                      []string `hcl:"socket_bind_deny,optional" systemd:"SocketBindDeny"`
	StartupAllowedCPUs                  string   `hcl:"startup_allowed_cp_us,optional" systemd:"StartupAllowedCPUs"`
	StartupAllowedMemoryNodes           string   `hcl:"startup_allowed_memory_nodes,optional" systemd:"StartupAllowedMemoryNodes"`
	StartupCPUWeight                    uint64   `hcl:"startup_cpu_weight,optional" systemd:"StartupCPUWeight"`
	StartupIOWeight                     uint64   `hcl:"startup_io_weight,optional" systemd:"StartupIOWeight"`
	StartupMemoryHigh                   string   `hcl:"startup_memory_high,optional" systemd:"StartupMemoryHigh"`
	StartupMemoryLow                    string   `hcl:"startup_memory_low,optional" systemd:"StartupMemoryLow"`
	StartupMemoryMax                    string   `hcl:"startup_memory_max,optional" systemd:"StartupMemoryMax"`
	StartupMemorySwapMax                string   `hcl:"startup_memory_swap_max,optional" systemd:"StartupMemorySwapMax"`
	StartupMemoryZSwapMax               string   `hcl:"startup_memory_z_swap_max,optional" systemd:"StartupMemoryZSwapMax"`
	TasksAccounting                     bool     `hcl:"tasks_accounting,optional" systemd:"TasksAccounting"`
	TasksMax                            string   `hcl:"tasks_max,optional" systemd:"TasksMax"`
}

type Slice struct {
	Name string `hcl:"name,label"`

	Unit    UnitBlock    `hcl:"unit,block"`
	Slice   SliceBlock   `hcl:"slice,block"`
	Install InstallBlock `hcl:"install,block"`
}
