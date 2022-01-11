package nsdp

import (
	"net"
	"time"
)

// Scan performs a discovery for devices on the network by sending
// an NSDP message to the broadcast address.
func Scan(ifaceName string, options ...Option) ([]Device, error) {
	opts, err := GetDefaultOptions().Apply(options...)
	if err != nil {
		return nil, err
	}

	// Fetch network interface.
	iface, err := GetInterface(ifaceName)
	if err != nil {
		return nil, err
	}

	// Create new discovery message.
	request := NewDiscoveryMessage(iface)

	// Send message to broadcast address.
	responses, err := Send(opts.Context, iface, request)
	if err != nil {
		return nil, err
	}

	// Convert responses to devices.
	devices := make([]Device, len(responses))
	for i, response := range responses {
		// This is safe because we previously allocated the slice.
		if err := devices[i].UnmarshalMessage(&response); err != nil {
			return nil, err
		}
	}

	return devices, nil
}

// NewDiscoveryMessage creates a new message that can be
// broadcasted to discover other devices on the network.
func NewDiscoveryMessage(iface *net.Interface) *Message {
	// Fetch MAC address from client interface.
	clientMAC := [6]uint8{}
	for i := 0; i < len(clientMAC); i++ {
		clientMAC[i] = iface.HardwareAddr[i]
	}

	// Create discovery message.
	msg := NewMessage()
	msg.Header.Operation = ReadRequest
	msg.Header.ClientMAC = clientMAC
	// The server MAC during discovery should be all-zero
	// as this will be interpreted as a multicast address
	// and cause all devices to respond to the message.
	msg.Header.ServerMAC = SelectorAll.MACMarshalBinary()
	// HACK: Because we want the CLI to be stateless we can't
	// keep track of a sequence number between subsequent calls.
	// But if we use the remainder when dividing the current
	// timestamp by our maximum sequence number we can get a
	// number that is very likely to be incrementing between
	// subsequent calls. If it is not incrementing the previous
	// call is ASSUMED to be so much in the past that the sequence
	// number is ASSUMED to be valid again. This SHOULD guarantee
	// a response from the device on every call.
	msg.Header.Sequence = uint16(time.Now().UnixNano()/1e6) % 0xFFFF

	// Define the information we would like to receive during
	// discovery. The list of records is limited to the most
	// common ones and therefore NOT the same as used by the
	// original tool provided by the manufacturer.
	scanRecords := []Record{
		{Type: RecordModel.ID},
		{Type: RecordName.ID},
		{Type: RecordMAC.ID},
		{Type: RecordIP.ID},
		{Type: RecordNetmask.ID},
		{Type: RecordGateway.ID},
		{Type: RecordDHCP.ID},
		{Type: RecordFirmware.ID},
	}
	msg.Records = append(msg.Records, scanRecords...)

	return msg
}
