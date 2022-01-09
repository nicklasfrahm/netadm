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
		switch record.Type {
		case RecordModel:
			d.Model = string(record.Value)
		case RecordName:
			d.Name = string(record.Value)
		case RecordMAC:
			d.MAC = net.HardwareAddr(record.Value)
		case RecordIP:
			d.IP = net.IP(record.Value)
		case RecordNetmask:
			d.Netmask = net.IP(record.Value)
		case RecordGateway:
			d.Gateway = net.IP(record.Value)
		case RecordDHCP:
			d.DHCP = record.Value[0] == 1
		case RecordFirmware:
			d.Firmware = string(record.Value)
		}
	}

	return nil
}
