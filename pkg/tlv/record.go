package tlv

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	// TypeEndOfMessage is the type identifier for
	// the record that marks the end of the message.
	TypeEndOfMessage = 0xFFFF
)

// Record is a TLV record.
type Record struct {
	Type   uint16
	Length uint16
	Value  []byte
}

// RecordList is a slice of TLV records.
type RecordList []Record

// MarshalBinary encodes the RecordList into a slice of bytes.
func (records *RecordList) MarshalBinary() ([]byte, error) {
	w := new(bytes.Buffer)

	for _, record := range *records {
		// Encode record type identifier.
		if err := binary.Write(w, binary.BigEndian, record.Type); err != nil {
			return nil, err
		}

		// Encode record length.
		if err := binary.Write(w, binary.BigEndian, record.Length); err != nil {
			return nil, err
		}

		// Encode record value.
		if err := binary.Write(w, binary.BigEndian, record.Value); err != nil {
			return nil, err
		}
	}

	return w.Bytes(), nil
}

// UnmarshalBinary decodes the slice of bytes into a RecordList.
func (records *RecordList) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	// Decode message records.
	for r.Len() > 0 {
		var record Record

		// Decode record type identifier.
		if err := binary.Read(r, binary.BigEndian, &record.Type); err != nil {
			return err
		}
		// Decode record length.
		if err := binary.Read(r, binary.BigEndian, &record.Length); err != nil {
			return err
		}

		// Check if magic bytes for end of message are reached.
		if record.Type == TypeEndOfMessage {
			// Check if the message is valid.
			if record.Length != 0 {
				return errors.New("invalid end of message")
			}
			return nil
		}

		// Decode record value.
		record.Value = make([]byte, record.Length)
		if record.Length > 0 {
			if err := binary.Read(r, binary.BigEndian, &record.Value); err != nil {
				return err
			}
		}
		*records = append(*records, record)
	}

	return nil
}
