package nsdp

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

func Scan(ifaceName string, options ...Option) (*[]Device, error) {
	opts, err := GetDefaultOptions().Apply(options...)
	if err != nil {
		return nil, err
	}

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, errors.New("unknown interface")
	}

	addresses, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	if len(addresses) == 0 {
		return nil, errors.New("interface has no address")
	}

	// Select IPv4 interface address.
	var ip net.IP
	for _, address := range addresses {
		// Check if address is IPv4.
		ipNet, ok := address.(*net.IPNet)
		if ok && ipNet.IP.To4() != nil {
			opts.Address = address
			ip = ipNet.IP
			break
		}
	}

	// Check if interface has a valid IPv4 address.
	if opts.Address == nil {
		return nil, errors.New("interface has no valid IPv4 address")
	}

	clientAddr := net.UDPAddr{
		IP:   ip,
		Port: ClientPort,
	}
	serverAddr := net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: ServerPort,
	}
	conn, err := net.DialUDP("udp", &clientAddr, &serverAddr)
	if err != nil {
		return nil, err
	}
	deadline, hasDeadline := opts.Context.Deadline()
	if hasDeadline {
		conn.SetDeadline(deadline)
		conn.SetReadDeadline(deadline)
	}

	// Close connection on context cancel.
	// go func() {
	// 	<-opts.Context.Done()
	// 	conn.Close()
	// 	return
	// }()

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
	msg.Header.Sequence = uint16(time.Now().UnixNano()/1e9) % 0xFFFF

	// Define the information we would like to receive during discovery.
	scanRecords := []Record{
		{Type: RecordModel},
		{Type: RecordName},
		{Type: RecordMAC},
		{Type: RecordIP},
	}
	msg.Records = append(msg.Records, scanRecords...)

	// Buffer message to ensure it is sent in a single datagram.
	buf := bytes.NewBuffer([]byte{})
	if err := msg.Write(buf); err != nil {
		return nil, err
	}

	// TODO: Open seperate UDP listener as traffic
	// is not related to broadcast address.

	// DEBUG: Split the string into chunks of 4 characters.
	hex := fmt.Sprintf("%X", buf)
	for i := 0; i < len(hex); i += 4 {
		fmt.Printf("%s ", []byte(hex[i:i+4]))
	}

	fmt.Printf("\t%d\n", buf.Len())

	// Write buffered message to connection.
	res := make([]byte, 1500)
	n, err := conn.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}
	fmt.Printf("Wrote %d bytes\n", n)
	// devices := make([]Device, 0)
	_, _, err = conn.ReadFromUDP(res)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%X\n", res)

	return nil, nil
}
