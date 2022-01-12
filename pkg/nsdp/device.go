package nsdp

import (
	"net"
	"reflect"
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
		// Get record type based on type identifier.
		recordType := RecordTypeIDs[record.Type]
		if recordType == nil {
			continue
		}

		// Parse values into their according types.
		field := reflect.ValueOf(d).Elem().FieldByName(recordType.Name)
		switch recordType.Example.(type) {
		case string:
			field.SetString(string(record.Value))
		case bool:
			field.SetBool(record.Value[0] == 1)
		case net.HardwareAddr:
			field.Set(reflect.ValueOf(net.HardwareAddr(record.Value)))
		case net.IP:
			field.Set(reflect.ValueOf(net.IP(record.Value)))
		}
	}

	return nil
}
