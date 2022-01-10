package nsdp

import (
	"errors"
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

	// Check if the provided interface has a valid configuration.
	iface, ip, err := GetValidInterface(ifaceName)
	if err != nil {
		return nil, err
	}

	// Create a UDP socket to listen for incoming packets.
	socketAddr := net.UDPAddr{
		IP:   *ip,
		Port: ClientPort,
	}
	socket, err := net.ListenUDP("udp", &socketAddr)
	if err != nil {
		return nil, err
	}
	defer socket.Close()

	devices := make([]Device, 0)
	errs := make(chan error, 1)

	// Create a goroutine to listen for incoming packets.
	go func() {
		for {
			select {
			case <-opts.Context.Done():
				return
			default:
				buf := make([]byte, 1500)

				n, err := socket.Read(buf)
				if err != nil {
					errs <- err
					return
				}

				msg := new(Message)
				if err := msg.UnmarshalBinary(buf[:n]); err != nil {
					errs <- errors.New("malformed response message")
					return
				}

				device := new(Device)
				if err := device.UnmarshalMessage(msg); err != nil {
					errs <- err
					return
				}
				devices = append(devices, *device)
			}
		}
	}()

	// Create discovery message and encode it into its binary form.
	msg, err := NewDiscoveryMessage(iface).MarshalBinary()
	if err != nil {
		return nil, err
	}

	// Send the message to the broadcast address.
	deviceAddr := net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: ServerPort,
	}
	if _, err := socket.WriteToUDP(msg, &deviceAddr); err != nil {
		return nil, err
	}

	select {
	case <-opts.Context.Done():
		return devices, nil
	case err := <-errs:
		return nil, err
	}
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
	msg.Header.ServerMAC = [6]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
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
		{Type: RecordModel},
		{Type: RecordName},
		{Type: RecordMAC},
		{Type: RecordIP},
		{Type: RecordNetmask},
		{Type: RecordGateway},
		{Type: RecordDHCP},
		{Type: RecordFirmware},
	}
	msg.Records = append(msg.Records, scanRecords...)

	return msg
}

// GetValidInterface fetches the interface based on the provided
// interface name if it is up and has an IPv4 address.
func GetValidInterface(ifaceName string) (*net.Interface, *net.IP, error) {
	// Fetch specified interface by name.
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, nil, errors.New("interface unknown")
	}

	// Check if interface is up.
	if iface.Flags&net.FlagUp == 0 {
		return nil, nil, errors.New("interface is down")
	}

	// Check if interface has addresses.
	addresses, err := iface.Addrs()
	if err != nil {
		return nil, nil, err
	}
	if len(addresses) == 0 {
		return nil, nil, errors.New("interface has no address")
	}

	// Select IPv4 interface address.
	var ip *net.IP
	var address net.Addr
	for _, addr := range addresses {
		// Check if address is IPv4.
		ipNet, ok := addr.(*net.IPNet)
		if ok && ipNet.IP.To4() != nil {
			ip = &ipNet.IP
			address = addr
			break
		}
	}

	// Check if interface has a valid IPv4 address.
	if address == nil {
		return nil, nil, errors.New("interface has no valid IPv4 address")
	}

	return iface, ip, nil
}
