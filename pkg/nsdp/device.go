package nsdp

import (
	"net"
)

// Device represents a switch network device.
type Device struct {
	Model    string
	Name     string
	MAC      net.HardwareAddr
	IP       net.IP
	Netmask  net.IP
	Gateway  net.IP
	DHCP     bool
	Firmware string
}

// UnmarshalMessage decodes a message into a Device.
func (d *Device) UnmarshalMessage(msg *Message) error {
	for _, record := range msg.Records {
		// TODO: Can we use the reflect package for this
		// without making typing implicit for the Device?
		switch record.Type {
		case RecordModel.ID:
			d.Model = string(record.Value)
		case RecordName.ID:
			d.Name = string(record.Value)
		case RecordMAC.ID:
			d.MAC = net.HardwareAddr(record.Value)
		case RecordIP.ID:
			d.IP = net.IP(record.Value)
		case RecordNetmask.ID:
			d.Netmask = net.IP(record.Value)
		case RecordGateway.ID:
			d.Gateway = net.IP(record.Value)
		case RecordDHCP.ID:
			d.DHCP = record.Value[0] == 1
		case RecordFirmware.ID:
			d.Firmware = string(record.Value)
		}
	}

	return nil
}
