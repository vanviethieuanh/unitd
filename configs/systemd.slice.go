// Generated based on man page systemd.slice of systemd

package configs

// Slice is for Slice systemd unit file
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
type Slice struct {
	Name string `hcl:"name,label"`

	Unit    UnitBlock    `hcl:"unit,block"`
	Install InstallBlock `hcl:"install,block"`
}
