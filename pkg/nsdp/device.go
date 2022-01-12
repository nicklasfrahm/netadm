package nsdp

import (
	"errors"
	"net"
	"reflect"
)

// Device represents a switch network device.
type Device struct {
	Model              string
	Name               string
	MAC                net.HardwareAddr
	IP                 net.IP
	Netmask            net.IP
	Gateway            net.IP
	DHCP               bool
	Firmware           string
	PasswordEncryption bool
	PortCount          uint8
}

// UnmarshalMessage decodes a message into a Device.
func (d *Device) UnmarshalMessage(msg *Message) error {
	for _, record := range msg.Records {
		// Fetch record type, because it tells us which field to map it to.
		rt := record.Type()
		if rt == nil {
			return errors.New("record type unknown")
		}

		// Dynamically decode the record type.
		value := record.Reflect()

		// Set the value of the field.
		field := reflect.ValueOf(d).Elem().FieldByName(rt.Name)
		if field.IsValid() {
			field.Set(value)
		}
	}

	return nil
}
