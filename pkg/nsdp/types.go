package nsdp

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"time"
)

const (
	ClientPort = 63321
	ServerPort = 63322
)

// RecordType describes what data a record contains.
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
	Record0x000B
	Record0x000C
	Record0x000D
	Record0x000E
	Record0x000F
)

// OpCode describes the operation
// that a message is performing.
type OpCode uint8

const (
	ReadRequest OpCode = iota + 1
	ReadResponse
	WriteRequest
	WriteResponse
)

// Options defines the configureation of an
// operations of this library.
type Options struct {
	Context context.Context
	Timeout time.Duration
}

// Apply applies the option functions to the current set of options.
func (o *Options) Apply(options ...Option) (*Options, error) {
	for _, option := range options {
		if err := option(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

// Option defines the function signature to set
// an option for the operations of this library.
type Option func(*Options) error

// GetDefaultOptions returns the default options
// for all operations of this library.
func GetDefaultOptions() *Options {
	return &Options{
		Timeout: time.Second,
		Context: context.Background(),
	}
}

// WithContext supplies a custom context the
// operations of this library. This makes it
// possible to cancel the operations of this
// library by using a timeout for example.
func WithContext(ctx context.Context) Option {
	return func(o *Options) error {
		o.Context = ctx
		return nil
	}
}

// Device represents a switch network device.
type Device struct {
	Model    string
	Name     string
	MAC      net.HardwareAddr
	IP       net.IP
	Netmask  net.IP
	Gateway  net.IP
	DHCP     bool
	Firmware string
}

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
	}

	return &msg
}

// Write encodes the message into bytes via an io.Writer and
// uses the big-endian byte order like most network protocols.
func (m *Message) Write(w io.Writer) error {
	// Encode message header.
	binary.Write(w, binary.BigEndian, m.Header)

	// Encode message records.
	for _, record := range m.Records {
		binary.Write(w, binary.BigEndian, record.Type)
		binary.Write(w, binary.BigEndian, record.Len)
		binary.Write(w, binary.BigEndian, record.Value)
	}

	// Magic bytes that mark the end of the message.
	binary.Write(w, binary.BigEndian, uint32(0xFFFF0000))

	return nil
}
