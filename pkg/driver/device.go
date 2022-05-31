package driver

import (
	"net"
)

// System contains system-wide information about the device.
type System struct {
	Model    string           `yaml:"model"`
	Name     string           `yaml:"name"`
	MAC      net.HardwareAddr `yaml:"mac"`
	IP       net.IP           `yaml:"ip"`
	Netmask  net.IPMask       `yaml:"netmask"`
	Gateway  net.IP           `yaml:"gateway"`
	DHCP     bool             `yaml:"dhcp"`
	Firmware string           `yaml:"firmware"`
	Ports    uint8            `yaml:"ports"`
}

// TODO: Status contains status information and metrics about the device.
// type Status struct {
// 	Links
// 	Metrics
// 	Cables
// }

// Device describes the structure of a network device.
type Device struct {
	System *System `yaml:"system"`
}

// Driver returns the driver of the device.
func (d *Device) Driver() *Driver {
	// TODO: Implement this.

	return nil
}
