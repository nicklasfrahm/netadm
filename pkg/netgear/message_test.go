package netgear

import (
	"testing"

	"github.com/nicklasfrahm/netadm/pkg/tlv"
	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	assert := assert.New(t)

	// Arrange.
	opcode := WriteRequest

	// Act.
	msg := NewMessage(opcode)

	// Assert.
	assert.Equal(uint8(1), msg.Header.Version, "should set version")
	assert.Equal([4]byte{'N', 'S', 'D', 'P'}, msg.Header.Signature, "should set signature")
	assert.Equal(opcode, msg.Header.Operation, "should set provided opcode")
	assert.Equal([6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, msg.Header.ClientMAC, "should set client MAC")
	assert.Equal([6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, msg.Header.ServerMAC, "should set server MAC")
}

func TestMarshalBinary(t *testing.T) {
	assert := assert.New(t)

	// Arrange.
	msg := NewMessage(WriteRequest)
	msg.Header.Sequence = 1
	msg.Records = append(
		msg.Records,
		tlv.Record{
			Type:   0x01,
			Length: 0x04,
			Value:  []byte{'t', 'e', 's', 't'},
		},
		tlv.Record{
			Type:   0x02,
			Length: 0x02,
			Value:  []byte{'h', 'i'},
		},
	)
	msg.Header.ClientMAC = [6]byte{0x42, 0x00, 0x42, 0x00, 0x42, 0x00}
	msg.Header.ServerMAC = [6]byte{0x84, 0x00, 0x84, 0x00, 0x84, 0x00}

	// Act.
	data, err := msg.MarshalBinary()

	// Assert.
	assert.NoError(err, "should not return an error")
	assert.Equal([]byte{0x01}, data[0:1], "should encode version")
	assert.Equal([]byte{0x03}, data[1:2], "should encode operation")
	assert.Equal([]byte{0x00, 0x00}, data[2:4], "should encode result")
	assert.Equal([]byte{0x00, 0x00, 0x00, 0x00}, data[4:8], "should encode reserved bytes as zeros")
	assert.Equal([]byte{0x42, 0x00, 0x42, 0x00, 0x42, 0x00}, data[8:14], "should encode client MAC")
	assert.Equal([]byte{0x84, 0x00, 0x84, 0x00, 0x84, 0x00}, data[14:20], "should encode server MAC")
	assert.Equal([]byte{0x00, 0x00}, data[20:22], "should encode reserved bytes as zeros")
	assert.Equal([]byte{0x00, 0x01}, data[22:24], "should encode sequence number")
	assert.Equal([]byte{'N', 'S', 'D', 'P'}, data[24:28], "should encode protocol signature")
	assert.Equal([]byte{0x00, 0x00, 0x00, 0x00}, data[28:32], "should encode reserved bytes as zeros")
	assert.Equal([]byte{0x00, 0x01, 0x00, 0x04, 't', 'e', 's', 't'}, data[32:40], "should encode first TLV record")
	assert.Equal([]byte{0x00, 0x02, 0x00, 0x02, 'h', 'i'}, data[40:46], "should encode second TLV record")
	assert.Equal([]byte{0xFF, 0xFF, 0x00, 0x00}, data[46:50], "should append and encode end of message TLV record")
}
