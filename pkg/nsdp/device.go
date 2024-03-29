package nsdp

import (
	"net"
	"reflect"
)

// TODO: Rework this model structure. It's bad. Properties
// should be group by port and not by feature.

// Device represents a switch network device.
type Device struct {
	Model                string
	Name                 string
	MAC                  net.HardwareAddr
	IP                   net.IP
	Netmask              net.IP
	Gateway              net.IP
	DHCP                 bool
	Firmware             string
	PasswordEncryption   EncryptionMode
	PasswordNonce        []byte
	PortSpeeds           []PortSpeed
	CableTestResult      CableTestResult
	VLANEngine           VLANEngine
	VLANsPort            []VLANPort
	VLANs802Q            []VLAN802Q
	PVIDs                []PVID
	QoSEngine            QoSEngine
	QoSPolicies          []QoSPolicy
	BandwidthLimitsIn    []BandwidthPolicy
	BandwidthLimitsOut   []BandwidthPolicy
	BroadcastFilter      bool
	BroadcastLimits      []BandwidthPolicy
	PortMetrics          []PortMetric
	PortMirroring        PortMirroring
	PortCount            uint8
	LoopDetection        bool
	IGMPSnoopingVLAN     IGMPSnoopingVLAN
	MulticastFilter      bool
	IGMPHeaderValidation bool
}

// UnmarshalMessage decodes a message into a Device.
func (d *Device) UnmarshalMessage(msg *Message) error {
	for _, record := range msg.Records {
		// Fetch record type, because it tells us which field to map it to.
		rt := record.Type()
		if rt == nil {
			return ErrRecordTypeUnknown
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
