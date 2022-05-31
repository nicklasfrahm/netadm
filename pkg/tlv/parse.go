package tlv

// Record is a TLV record.
type Record struct {
	Type   int16
	Length int16
	Value  []byte
}
