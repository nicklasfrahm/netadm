package nsdp

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// RecordType describes which data a record contains.
type RecordType uint16

const (
	// Manufacturer's model name of the device.
	RecordModel RecordType = iota + 1
	// Unknown record 0x0002.
	Record0x0002
	// User-defined name of the device.
	RecordName
	// MAC address of the device.
	RecordMAC
	// Unknown record 0x0005.
	Record0x0005
	// IP address of the device.
	RecordIP
	// Netmask of the device.
	RecordNetmask
	// Gateway of the device.
	RecordGateway
	// Unknown record 0x0009.
	Record0x0009
	// Unknown record 0x000A.
	Record0x000A
	// DHCP status of the device.
	RecordDHCP
	// Unknown record 0x000C.
	Record0x000C
	// Version of the firmware currently running on the device.
	RecordFirmware
	// Unknown record 0x000E.
	Record0x000E
	// Unknown record 0x000F.
	Record0x000F
	// A special record type that identifies
	// the end of the message. Combined with
	// a length of 0, this forms the 4 magic
	// bytes that mark the end of the message
	// (0xFFFF0000).
	RecordEndOfMessage = RecordType(0xFFFF)
)

// OpCode describes the operation that a message is performing.
type OpCode uint8

const (
	// ReadRequest is the OpCode that identifies read
	// request messages sent by the host client.
	ReadRequest OpCode = iota + 1
	// ReadResponse is the OpCode that identifies read
	// response messages sent by the device server.
	ReadResponse
	// WriteRequest is the OpCode that identifies write
	// request messages sent by the host client.
	WriteRequest
	// WriteResponse is the OpCode that identifies write
	// response messages sent by the device server.
	WriteResponse
)

// Record defines the binary encoding of a
// type-length-value object, which makes it
// possible to encode variable length values
// in a binary format.
type Record struct {
	Type  RecordType
	Len   uint16
	Value []uint8
}

// Header defines the binary encoding of the
// UDP payload message header.
type Header struct {
	Version   uint8
	Operation OpCode
	Result    uint16
	_         [4]uint8
	ClientMAC [6]uint8
	ServerMAC [6]uint8
	_         [2]uint8
	Sequence  uint16
	Signature [4]uint8
	_         [4]uint8
}

// Message defines the binary encoding scheme of the
// UDP payload. The order of the fields determines
// how the data is encoded and decoded respectively.
type Message struct {
	Header  Header
	Records []Record
}

// Message creates a new message to the device with
// the default options.
func NewMessage() *Message {
	msg := Message{
		Header: Header{
			// The version of the protocol is always 1.
			Version: 1,
			// The signature of the protocol is always "NSDP".
			Signature: [4]uint8{'N', 'S', 'D', 'P'},
		},
		Records: make([]Record, 0),
	}

	return &msg
}

// UnmarshalBinary decodes the bytes of a message into the message structure.
func (m *Message) UnmarshalBinary(data []byte) error {
	// Decode message header.
	r := bytes.NewReader(data)
	if err := binary.Read(r, binary.BigEndian, &m.Header); err != nil {
		return err
	}

	// Decode message records.
	for r.Len() > 0 {
		var record Record
		// Decode record type.
		if err := binary.Read(r, binary.BigEndian, &record.Type); err != nil {
			return err
		}
		// Decode record length.
		if err := binary.Read(r, binary.BigEndian, &record.Len); err != nil {
			return err
		}

		// Check if magic bytes for end of message are reached.
		if record.Type == RecordEndOfMessage {
			// Check if the message is valid.
			if record.Len != 0 {
				return errors.New("invalid end of message")
			}
			return nil
		}

		// Decode record value.
		record.Value = make([]uint8, record.Len)
		if record.Len > 0 {
			if err := binary.Read(r, binary.BigEndian, &record.Value); err != nil {
				return err
			}
		}
		m.Records = append(m.Records, record)
	}

	return nil
}

// MarshalBinary encodes the message structure into a slice of bytes.
func (m *Message) MarshalBinary() ([]byte, error) {
	w := new(bytes.Buffer)

	// Encode message header.
	if err := binary.Write(w, binary.BigEndian, m.Header); err != nil {
		return nil, err
	}

	// Encode message records.
	for _, record := range m.Records {
		// Encode record type.
		if err := binary.Write(w, binary.BigEndian, record.Type); err != nil {
			return nil, err
		}

		// Encode record length.
		if err := binary.Write(w, binary.BigEndian, record.Len); err != nil {
			return nil, err
		}

		// Encode record value.
		if err := binary.Write(w, binary.BigEndian, record.Value); err != nil {
			return nil, err
		}
	}

	// Magic bytes that mark the end of the message.
	if err := binary.Write(w, binary.BigEndian, uint32(0xFFFF0000)); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
