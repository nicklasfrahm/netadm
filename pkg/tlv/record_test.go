package tlv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalBinary(t *testing.T) {
	assert := assert.New(t)

	// Arrange.
	data := []byte{
		0x00, 0x99, 0x00, 0x04, 't', 'e', 's', 't',
		0x99, 0x00, 0x00, 0x02, 'h', 'i',
		0xFF, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x07, 's', 'k', 'i', 'p', 'p', 'e', 'd',
	}
	records := new(RecordList)

	// Act.
	err := records.UnmarshalBinary(data)

	// Assert.
	assert.NoError(err, "should not return an error")
	assert.Equal(2, len(*records), "should stop decoding after end of message record")
	assert.Equal(RecordList{
		{
			Type:   0x0099,
			Length: 0x0004,
			Value:  []byte{'t', 'e', 's', 't'},
		},
		{
			Type:   0x9900,
			Length: 0x0002,
			Value:  []byte{'h', 'i'},
		},
	}, *records, "should decode data and omit end of message record")
}

func TestUnmarshalBinaryInvalidEndOfMessage(t *testing.T) {
	assert := assert.New(t)

	// Arrange.
	data := []byte{
		0x00, 0x99, 0x00, 0x04, 't', 'e', 's', 't',
		0xFF, 0xFF, 0x11, 0x11,
	}
	records := new(RecordList)

	// Act.
	err := records.UnmarshalBinary(data)

	// Assert.
	assert.Equal(ErrInvalidEndOfMessage, err, "should return correct error")
	assert.Equal(1, len(*records), "should partially decode data")
	assert.Equal(RecordList{
		{
			Type:   0x0099,
			Length: 0x0004,
			Value:  []byte{'t', 'e', 's', 't'},
		},
	}, *records, "should decode data until invalid end of message record")
}

func TestMarshalBinary(t *testing.T) {
	assert := assert.New(t)

	// Arrange.
	records := &RecordList{
		{
			Type:   0x0099,
			Length: 0x0004,
			Value:  []byte{'t', 'e', 's', 't'},
		},
		{
			Type:   0x9900,
			Length: 0x0002,
			Value:  []byte{'h', 'i'},
		},
		RecordEndOfMessage,
	}

	// Act.
	bytes, err := records.MarshalBinary()

	// Assert.
	assert.NoError(err, "should not return an error")
	assert.Equal(18, len(bytes), "should encode the correct number of bytes")
	assert.Equal([]byte{
		0x00, 0x99, 0x00, 0x04, 't', 'e', 's', 't',
		0x99, 0x00, 0x00, 0x02, 'h', 'i',
		0xFF, 0xFF, 0x00, 0x00,
	}, bytes, "should encode all data including the end of message record")
}
