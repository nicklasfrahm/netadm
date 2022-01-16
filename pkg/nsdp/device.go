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
	PortSpeeds         []PortSpeed
	CableTestResult    CableTestResult
	PortMetrics        []PortMetric
	PortMirroring      PortMirroring
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
			// This is a minor hack as using reflect.Kind() == reflect.Slice
			// will give false positives for MAC and IP addresses.
			if rt.Slice {
				// Initialize slice if it is nil.
				if field.IsZero() {
					field.Set(reflect.MakeSlice(field.Type(), 0, 0))
				}
				field.Set(reflect.Append(field, value))
			} else {
				field.Set(value)
			}
		}
	}

	return nil
}
