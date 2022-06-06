package netgear

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/nicklasfrahm/netadm/pkg/tlv"
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

// Header defines the binary encoding of the
// UDP payload message header.
type Header struct {
	Version   uint8
	Operation OpCode
	Result    uint16
	_         [4]byte
	ClientMAC [6]byte
	ServerMAC [6]byte
	_         [2]byte
	Sequence  uint16
	Signature [4]byte
	_         [4]byte
}

// HeaderSize denotes the size of the message header in bytes.
const HeaderSize = 32

// Message defines the binary encoding scheme of the
// UDP payload. The order of the fields determines
// how the data is encoded and decoded respectively.
type Message struct {
	Header  Header
	Records tlv.RecordList
}

// NewMessage creates a new message to the device with
// the default options.
func NewMessage(operation OpCode) *Message {
	msg := Message{
		Header: Header{
			// The version of the protocol is always 1.
			Version: 1,
			// The signature of the protocol is always "NSDP".
			Signature: [4]byte{'N', 'S', 'D', 'P'},
			// Configure the operation based on the provided OpCode.
			Operation: operation,
		},
		Records: make(tlv.RecordList, 0),
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
	r := bytes.NewReader(data[:HeaderSize])
	if err := binary.Read(r, binary.BigEndian, &m.Header); err != nil {
		return err
	}

	if err := m.Records.UnmarshalBinary(data[HeaderSize:]); err != nil {
		return err
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

	// Append end of message record.
	m.Records = append(m.Records, tlv.RecordEndOfMessage)

	// Encode message records.
	recordsBinary, err := m.Records.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(recordsBinary); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
