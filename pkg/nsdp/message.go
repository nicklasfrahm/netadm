package nsdp

import (
	"bytes"
	"encoding/binary"
	"time"
)

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

// NewMessage creates a new message to the device with
// the default options.
func NewMessage(operation OpCode) *Message {
	msg := Message{
		Header: Header{
			// The version of the protocol is always 1.
			Version: 1,
			// The signature of the protocol is always "NSDP".
			Signature: [4]uint8{'N', 'S', 'D', 'P'},
			// Configure the operation based on the provided OpCode.
			Operation: operation,
			// Call all devices by default.
			ServerMAC: MACMarshalBinary(SelectorAll.MAC),
		},
		Records: make([]Record, 0),
	}

	// HACK: Because we want the CLI to be stateless we can't
	// keep track of a sequence number between subsequent calls.
	// But if we use the remainder when dividing the current
	// timestamp by our maximum sequence number we can get a
	// number that is very likely to be incrementing between
	// subsequent calls. If it is not incrementing the previous
	// call is ASSUMED to be so much in the past that the sequence
	// number is ASSUMED to be valid again. This SHOULD maximize
	// the chance to get a response from the device on every call.
	msg.Header.Sequence = uint16(time.Now().UnixNano()/1e6) % 0xFFFF

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
		// Decode record type identifier.
		if err := binary.Read(r, binary.BigEndian, &record.ID); err != nil {
			return err
		}
		// Decode record length.
		if err := binary.Read(r, binary.BigEndian, &record.Len); err != nil {
			return err
		}

		// Check if magic bytes for end of message are reached.
		if record.ID == RecordEndOfMessage.ID {
			// Check if the message is valid.
			if record.Len != 0 {
				return ErrInvalidEndOfMessage
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
		// Encode record type identifier.
		if err := binary.Write(w, binary.BigEndian, record.ID); err != nil {
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
