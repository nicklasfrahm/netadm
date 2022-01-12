package nsdp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"strings"
	"time"
)

// RecordTypeID is the ID of a RecordType.
type RecordTypeID uint16

// RecordType describes which data a Record contains.
type RecordType struct {
	ID      RecordTypeID
	Name    string
	Example interface{}
}

// NewRecordType creates a new record type.
func NewRecordType(id uint16, name string, example interface{}) *RecordType {
	return &RecordType{
		ID:      RecordTypeID(id),
		Name:    name,
		Example: example,
	}
}

// TODO: Add missing records types once all operations are implemented.

var (
	// RecordNone is a placeholder for an invalid or empty record.
	RecordNone = NewRecordType(0x0000, "None", nil)
	// RecordModel contains the device's manufacturer-provided model name.
	RecordModel = NewRecordType(0x0001, "Model", "GS308E")
	// RecordName contains the device's user-defined name.
	RecordName = NewRecordType(0x0003, "Name", "switch-0")
	// RecordMAC contains the device's MAC address.
	RecordMAC = NewRecordType(0x0004, "MAC", net.HardwareAddr{0x33, 0x0B, 0xC9, 0x5E, 0x51, 0x3A})
	// RecordIP contains the device's IP address.
	RecordIP = NewRecordType(0x0006, "IP", net.IP{192, 168, 0, 253})
	// RecordNetmask contains the device's netmask.
	RecordNetmask = NewRecordType(0x0007, "Netmask", net.IP{255, 255, 255, 0})
	// RecordGateway contains the device's gateway.
	RecordGateway = NewRecordType(0x0008, "Gateway", net.IP{192, 168, 0, 254})
	// RecordDHCP contains the device's DHCP status.
	RecordDHCP = NewRecordType(0x000B, "DHCP", false)
	// RecordFirmware contains the device's firmware version.
	RecordFirmware = NewRecordType(0x000D, "Firmware", "1.00.10")
	// RecordEndOfMessage special record type that identifies the end
	// of the message. Combined with a length of 0, this forms the 4
	// magic bytes that mark the end of the message (0xFFFF0000).
	RecordEndOfMessage = NewRecordType(0xFFFF, "EndOfMessage", nil)
)

// RecordTypeIDs maps the ID of a record to a record type.
var RecordTypeIDs = map[RecordTypeID]*RecordType{
	RecordNone.ID:         RecordNone,
	RecordModel.ID:        RecordModel,
	RecordName.ID:         RecordName,
	RecordMAC.ID:          RecordMAC,
	RecordIP.ID:           RecordIP,
	RecordNetmask.ID:      RecordNetmask,
	RecordGateway.ID:      RecordGateway,
	RecordDHCP.ID:         RecordDHCP,
	RecordFirmware.ID:     RecordFirmware,
	RecordEndOfMessage.ID: RecordEndOfMessage,
}

// RecordTypeNames maps the name of a record to a record type.
var RecordTypeNames = indexRecordTypeNames()

// indexRecordTypeNames builds an index of the record names.
func indexRecordTypeNames() map[string]*RecordType {
	recordNames := make(map[string]*RecordType, len(RecordTypeIDs))

	for _, record := range RecordTypeIDs {
		// Exclude the None and the EndOfMessage record types.
		if record.Example != nil {
			recordNames[strings.ToLower(record.Name)] = record
		}
	}

	return recordNames
}

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
	Type  RecordTypeID
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
		// Decode record type.
		if err := binary.Read(r, binary.BigEndian, &record.Type); err != nil {
			return err
		}
		// Decode record length.
		if err := binary.Read(r, binary.BigEndian, &record.Len); err != nil {
			return err
		}

		// Check if magic bytes for end of message are reached.
		if record.Type == RecordEndOfMessage.ID {
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

// NewDiscoveryMessage creates a new message that can be
// broadcasted to discover other devices on the network.
func NewDiscoveryMessage() *Message {
	// Create discovery message.
	msg := NewMessage(ReadRequest)

	// The server MAC during discovery should be all-zero
	// as this will be interpreted as a multicast address
	// and cause all devices to respond to the message.
	msg.Header.ServerMAC = MACMarshalBinary(SelectorAll.MAC)

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
