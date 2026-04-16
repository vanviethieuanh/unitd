// Generated based on man page systemd.device of systemd

package configs

// Device is for Device systemd unit file
// A unit configuration file whose name ends in encodes information about a device unit as exposed in
// the sysfs/ device tree. This may be used to define dependencies between devices and other units.
// This unit type has no specific options. See for the common options of all unit configuration files.
// The common configuration items are configured in the generic [Unit] and [Install] sections. A
// separate [Device] section does not exist, since no device-specific options may be configured.
// systemd will dynamically create device units for all kernel devices that are marked with the udev
// tag (by default all block and network devices, and a few others). Note that . Device units are named
// after the and paths they control. Example: the device is exposed in systemd as . For details about
// the escaping logic used to convert a file system path to a unit name see . To tag a udev device, use
// in the udev rules file, see for details. Device units will be reloaded by systemd whenever the
// corresponding device generates a event. Other units can use to react to that event.
type Device struct {
	Name string `hcl:"name,label"`

	Unit    UnitBlock    `hcl:"unit,block"`
	Install InstallBlock `hcl:"install,block"`
}
