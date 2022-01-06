package nsdp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
)

func Scan(ifaceName string, options ...Option) (*[]Device, error) {
	opts := GetDefaultOptions()

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

	opts.Address = addresses[0]
	clientAddr := net.UDPAddr{
		IP:   net.ParseIP(opts.Address.String()),
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
	// TODO: Use proper mechanism to set sequence number.
	msg.Header.Sequence = 1

	// Define the information we would like to receive during discovery.
	scanRecords := []Record{
		{Type: RecordModel},
		{Type: RecordName},
		{Type: RecordMAC},
		{Type: RecordIP},
	}
	msg.Records = append(msg.Records, scanRecords...)

	buf := bytes.NewBuffer([]byte{})
	if err := msg.Write(buf); err != nil {
		return nil, err
	}

	// Split the string into chunks of 4 characters.
	hex := fmt.Sprintf("%X", buf)
	for i := 0; i < len(hex); i += 4 {
		fmt.Printf("%s ", []byte(hex[i:i+4]))
	}

	fmt.Printf("\t%d\n", buf.Len())

	// TODO: Remove this line.
	os.Exit(0)

	// Write message to connection.
	err = msg.Write(conn)
	if err != nil {
		return nil, err
	}
	fmt.Println("Message sent")

	res := make([]byte, 10)
	// devices := make([]Device, 0)
	_, err = bufio.NewReader(conn).Read(res)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%X\n", res)
	conn.Close()

	os.Exit(0)

	return nil, nil
}
