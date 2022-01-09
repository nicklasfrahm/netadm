package nsdp

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// RecordType describes which data a record contains.
type RecordType uint16

const (
	RecordModel  RecordType = iota + 1
	Record0x0002            // Is this the serial number?
	RecordName
	RecordMAC
	Record0x0005
	RecordIP
	RecordNetmask
	RecordGateway
	Record0x0009
	Record0x000A
	RecordDHCP
	Record0x000C
	RecordFirmware
	Record0x000E
	Record0x000F
	RecordEndOfMessage = RecordType(0xFFFF)
)

// OpCode describes the operation that a message is performing.
type OpCode uint8

const (
	ReadRequest OpCode = iota + 1
	ReadResponse
	WriteRequest
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
