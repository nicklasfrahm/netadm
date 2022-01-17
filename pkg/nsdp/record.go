package nsdp

import (
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"strings"
)

// LinkStatus defines the speed of a network link.
type LinkStatus uint8

const (
	// LinkDown marks a link to be down.
	LinkDown LinkStatus = iota
	// LinkSpeed10MbitHalfDuplex is the speed of a 10Mbit Ethernet link in half-duplex mode.
	LinkSpeed10MbitHalfDuplex
	// LinkSpeed10Mbit is the speed of a 10Mbit Ethernet link in full-duplex mode.
	LinkSpeed10Mbit
	// LinkSpeed100MbitHalfDuplex is the speed of a 100Mbit Ethernet link in half-duplex mode.
	LinkSpeed100MbitHalfDuplex
	// LinkSpeed100Mbit is the speed of a 100Mbit Ethernet link in full-duplex mode.
	LinkSpeed100Mbit
	// LinkSpeed1Gbit is the speed of a 1Gbit Ethernet link in full-duplex mode.
	LinkSpeed1Gbit
	// LinkSpeed10Gbit is the speed of a 10Gbit Ethernet link in full-duplex mode.
	LinkSpeed10Gbit
)

// String returns the string representation of the link status.
func (l LinkStatus) String() string {
	switch l {
	case LinkDown:
		return "Down"
	case LinkSpeed10MbitHalfDuplex:
		return "10M+HX"
	case LinkSpeed10Mbit:
		return "10M"
	case LinkSpeed100MbitHalfDuplex:
		return "100M+HX"
	case LinkSpeed100Mbit:
		return "100M"
	case LinkSpeed1Gbit:
		return "1G"
	case LinkSpeed10Gbit:
		return "10G"
	default:
		return "Unknown"
	}
}

// VLANEngine defines the VLAN engine.
type VLANEngine uint8

const (
	// VLANEngineDisabled is the VLAN engine if it is disabled.
	VLANEngineDisabled VLANEngine = iota
	// VLANEnginePortBasic is the VLAN engine if it uses basic port-based VLANs.
	VLANEnginePortBasic
	// VLANEnginePortAdvanced is the VLAN engine if it uses advanced port-based VLANs.
	VLANEnginePortAdvanced
	// VLANEngine802QBasic is the VLAN engine if it uses basic 802.1Q VLANs.
	VLANEngine802QBasic
	// VLANEngine802QAdvanced is the VLAN engine if it uses advanced 802.1Q VLANs.
	VLANEngine802QAdvanced
)

// String returns the string representation of a VLAN Engine
func (e VLANEngine) String() string {
	switch e {
	case VLANEngineDisabled:
		return "Disabled"
	case VLANEnginePortBasic:
		return "PortBasic"
	case VLANEnginePortAdvanced:
		return "PortAdvanced"
	case VLANEngine802QBasic:
		return "802.1QBasic"
	case VLANEngine802QAdvanced:
		return "802.1QAdvanced"
	default:
		return "Unknown"
	}
}

// QoSEngine defines the quality of service engine.
type QoSEngine uint8

const (
	// QoSPort is a port-based QoS engine.
	QoSPort QoSEngine = iota + 1
	// QoSDSCP is a DSCP-based QoS engine.
	QoSDSCP
)

// String returns the string representation of a QoS engine.
func (q QoSEngine) String() string {
	switch q {
	case QoSPort:
		return "Port"
	case QoSDSCP:
		return "DSCP"
	default:
		return "Unknown"
	}
}

// QoSPriority describes a port-based QoS priority.
type QoSPriority uint8

const (
	// QoSPriorityHigh assigns traffic the highest priority.
	QoSPriorityHigh QoSPriority = iota + 1
	// QoSPriorityMedium assigns traffic an above-normal priority.
	QoSPriorityMedium
	// QoSPriorityNormal assigns traffic a normal priority.
	QoSPriorityNormal
	// QoSPriorityLow is the lowest priority.
	QoSPriorityLow
)

// String returns the string representation of a QoS priority.
func (p QoSPriority) String() string {
	switch p {
	case QoSPriorityHigh:
		return "High"
	case QoSPriorityMedium:
		return "Medium"
	case QoSPriorityNormal:
		return "Normal"
	case QoSPriorityLow:
		return "Low"
	default:
		return "Unknown"
	}
}

// BandwidthLimit describes a bandwidth limit.
type BandwidthLimit uint8

const (
	// BandwidthLimitNone does not apply any bandwidth limits.
	BandwidthLimitNone BandwidthLimit = iota
	// BandwidthLimit512Kbps limits the bandwidth to 512Kbps.
	BandwidthLimit512Kbps
	// BandwidthLimit1Mbps limits the bandwidth to 1Mbps.
	BandwidthLimit1Mbps
	// BandwidthLimit2Mbps limits the bandwidth to 2Mbps.
	BandwidthLimit2Mbps
	// BandwidthLimit4Mbps limits the bandwidth to 4Mbps.
	BandwidthLimit4Mbps
	// BandwidthLimit8Mbps limits the bandwidth to 8Mbps.
	BandwidthLimit8Mbps
	// BandwidthLimit16Mbps limits the bandwidth to 16Mbps.
	BandwidthLimit16Mbps
	// BandwidthLimit32Mbps limits the bandwidth to 32Mbps.
	BandwidthLimit32Mbps
	// BandwidthLimit64Mbps limits the bandwidth to 64Mbps.
	BandwidthLimit64Mbps
	// BandwidthLimit128Mbps limits the bandwidth to 128Mbps.
	BandwidthLimit128Mbps
	// BandwidthLimit256Mbps limits the bandwidth to 256Mbps.
	BandwidthLimit256Mbps
	// BandwidthLimit512Mbps limits the bandwidth to 512Mbps.
	BandwidthLimit512Mbps
)

// String returns the string representation of a bandwidth limit.
func (b BandwidthLimit) String() string {
	switch b {
	case BandwidthLimitNone:
		return "None"
	case BandwidthLimit512Kbps:
		return "512Kbps"
	case BandwidthLimit1Mbps:
		return "1Mbps"
	case BandwidthLimit2Mbps:
		return "2Mbps"
	case BandwidthLimit4Mbps:
		return "4Mbps"
	case BandwidthLimit8Mbps:
		return "8Mbps"
	case BandwidthLimit16Mbps:
		return "16Mbps"
	case BandwidthLimit32Mbps:
		return "32Mbps"
	case BandwidthLimit64Mbps:
		return "64Mbps"
	case BandwidthLimit128Mbps:
		return "128Mbps"
	case BandwidthLimit256Mbps:
		return "256Mbps"
	case BandwidthLimit512Mbps:
		return "512Mbps"
	default:
		return "Unknown"
	}
}

// PortSpeed describes the speed of a port.
type PortSpeed struct {
	ID    uint8
	Speed LinkStatus
}

// String returns the string representation of the port speed.
func (p PortSpeed) String() string {
	return fmt.Sprintf("%d:%s", p.ID, p.Speed.String())
}

// PortMetric contains network traffic metrics of a port.
// TODO: Find out what the other metrics are.
type PortMetric struct {
	ID              uint8
	BytesReceived   uint64
	BytesSent       uint64
	ErrorsPacketCRC uint64
}

// String returns the string representation of a port metric.
func (p PortMetric) String() string {
	return fmt.Sprintf("%d:%d/%d/%d", p.ID, p.BytesReceived, p.BytesSent, p.ErrorsPacketCRC)
}

// PortMirroring describes the port mirroring configuration of all ports.
type PortMirroring struct {
	Destination uint8
	Sources     []uint8
}

// String returns the string representation of the port mirroring configuration.
func (p PortMirroring) String() string {
	if p.Destination == 0 {
		return "Disabled"
	}
	return fmt.Sprintf("%d:%s", p.Destination, joinInts(p.Sources, "+"))
}

// VLANPort describes the configuration of a port-based VLAN.
type VLANPort struct {
	ID    uint16
	Ports []uint8
}

// String returns the string representation of a VLAN.
func (v VLANPort) String() string {
	return fmt.Sprintf("%d:%s", v.ID, joinInts(v.Ports, "+"))
}

// VLAN802Q describes the configuration of an 802.1Q VLAN.
type VLAN802Q struct {
	ID       uint16
	Tagged   []uint8
	Untagged []uint8
}

// String returns the string representation of an 802.1Q VLAN.
func (v VLAN802Q) String() string {
	return fmt.Sprintf("%dt%su%s", v.ID, joinInts(v.Tagged, "+"), joinInts(v.Untagged, "+"))
}

// PVID describes the PVID assignment of a port.
type PVID struct {
	ID   uint8
	PVID uint16
}

// String returns the string representation of a PVID mapping.
func (p PVID) String() string {
	return fmt.Sprintf("%d:%d", p.ID, p.PVID)
}

// QoSPolicy describes the QoS policy of a port.
type QoSPolicy struct {
	ID       uint8
	Priority QoSPriority
}

// String returns the string representation of a QoS policy.
func (q QoSPolicy) String() string {
	return fmt.Sprintf("%d:%s", q.ID, q.Priority.String())
}

// BandwidthPolicy describes the bandwidth limit of a port.
type BandwidthPolicy struct {
	ID    uint8
	Limit BandwidthLimit
}

// String returns the string representation of a bandwidth policy.
func (b BandwidthPolicy) String() string {
	return fmt.Sprintf("%d:%s", b.ID, b.Limit.String())
}

// CableTestResult contains the results of a cable test.
type CableTestResult []uint8

// RecordTypeID is the ID of a RecordType.
type RecordTypeID uint16

// IGMPSnoopingVLAN describes the VLAN ID of the IGMP
// snooping VLAN. If this value is zero, it is disabled.
type IGMPSnoopingVLAN uint16

// RecordType describes which data a Record contains.
type RecordType struct {
	ID      RecordTypeID
	Name    string
	Example interface{}
	Slice   bool
}

// NewRecordType creates a new record type.
func NewRecordType(id uint16, name string, example interface{}) *RecordType {
	return &RecordType{
		ID:      RecordTypeID(id),
		Name:    name,
		Example: example,
	}
}

// SetSlice sets the slice flag on the record.
func (r *RecordType) SetSlice(slice bool) *RecordType {
	r.Slice = slice
	return r
}

// TODO: Define interface for record type that allows encoding and decoding into more semantic structs.
// TODO: Add missing records types once all operations are implemented.

var (
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
	// PasswordEncryption specifies whether the password is transmitted encrypted or plain-text.
	RecordPasswordEncryption = NewRecordType(0x0014, "PasswordEncryption", false)
	// RecordPortSpeeds contains the link status and the speed of a port.
	RecordPortSpeeds = NewRecordType(0x0C00, "PortSpeeds", []PortSpeed{{1, LinkSpeed1Gbit}, {2, LinkDown}}).SetSlice(true)
	// RecordPortMetrics contains network traffic metrics of a port.
	RecordPortMetrics = NewRecordType(0x1000, "PortMetrics", []PortMetric{{1, 64, 32, 0}}).SetSlice(true)
	// RecordCableTestResult contains the result of a cable test.
	RecordCableTestResult = NewRecordType(0x1C00, "CableTestResult", CableTestResult{0, 0, 0, 0, 0, 119, 30, 183, 118})
	// RecordVLANEngine contains the active VLAN engine.
	RecordVLANEngine = NewRecordType(0x2000, "VLANEngine", VLANEngineDisabled)
	// RecordVLANPort contains the configuration of a VLAN.
	RecordVLANPort = NewRecordType(0x2400, "VLANsPort", []VLANPort{{1, []uint8{1, 2, 3, 4, 5, 6, 7, 8}}}).SetSlice(true)
	// RecordVLAN802Q contains the configuration of a 802.1Q VLAN.
	RecordVLAN802Q = NewRecordType(0x2800, "VLANs802Q", []VLAN802Q{{1, []uint8{1, 2}, []uint8{3, 4, 5, 6, 7, 8}}}).SetSlice(true)
	// RecordPVIDs contains the 802.1Q VLAN IDs for each port often also referred to as PVIDs.
	RecordPVIDs = NewRecordType(0x3000, "PVIDs", []PVID{{1, 2}, {2, 2}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, {7, 1}, {8, 1}}).SetSlice(true)
	// RecordQoSEngine contains the QoS engine.
	RecordQoSEngine = NewRecordType(0x3400, "QoSEngine", QoSDSCP)
	// RecordQoSPolicies contains the QoS policy of a port.
	RecordQoSPolicies = NewRecordType(0x3800, "QoSPolicies", []QoSPolicy{{1, QoSPriorityNormal}, {2, QoSPriorityHigh}}).SetSlice(true)
	// RecordBandwidthLimitsIn contains the inbound bandwidth limit of a port.
	RecordBandwidthLimitsIn = NewRecordType(0x4C00, "BandwidthLimitsIn", []BandwidthPolicy{{1, BandwidthLimit256Mbps}, {2, BandwidthLimitNone}}).SetSlice(true)
	// RecordBandwidthLimitsOut contains the inbound bandwidth limit of a port.
	RecordBandwidthLimitsOut = NewRecordType(0x5000, "BandwidthLimitsOut", []BandwidthPolicy{{1, BandwidthLimit256Mbps}, {2, BandwidthLimitNone}}).SetSlice(true)
	// RecordBroadcastFilter defines whether broadcast storm control is enabled.
	RecordBroadcastFilter = NewRecordType(0x5400, "BroadcastFilter", false)
	// RecordBroadcastLimits contains the broadcast filter configuration of a port.
	RecordBroadcastLimits = NewRecordType(0x5800, "BroadcastLimits", []BandwidthPolicy{{1, BandwidthLimit256Mbps}, {2, BandwidthLimitNone}}).SetSlice(true)
	// RecordPortMirroring contains the mirroring configuration of all ports.
	RecordPortMirroring = NewRecordType(0x5C00, "PortMirroring", PortMirroring{1, []uint8{2, 3}})
	// RecordPortCount contains the number of ports on the device.
	RecordPortCount = NewRecordType(0x6000, "PortCount", uint8(5))
	// RecordIGMPSnoopingVLAN contains the VLAN ID used for IGMP snooping.
	RecordIGMPSnoopingVLAN = NewRecordType(0x6800, "IGMPSnoopingVLAN", IGMPSnoopingVLAN(1))
	// RecordMulticastFilter defines whether the device is configured to filter unknown multicast addresses.
	RecordMulticastFilter = NewRecordType(0x6C00, "MulticastFilter", false)
	// RecordIGMPHeaderValidation contains the IGMPv3 header validation status of the device.
	RecordIGMPHeaderValidation = NewRecordType(0x7000, "IGMPHeaderValidation", false)
	// RecordLoopDetection contains the loop detection status of the device.
	RecordLoopDetection = NewRecordType(0x9000, "LoopDetection", false)
	// RecordEndOfMessage special record type that identifies the end
	// of the message. Combined with a length of 0, this forms the 4
	// magic bytes that mark the end of the message (0xFFFF0000).
	RecordEndOfMessage = NewRecordType(0xFFFF, "EndOfMessage", nil)
)

// RecordTypeByID maps the ID of a record to a record type.
var RecordTypeByID = map[RecordTypeID]*RecordType{
	RecordModel.ID:                RecordModel,
	RecordName.ID:                 RecordName,
	RecordMAC.ID:                  RecordMAC,
	RecordIP.ID:                   RecordIP,
	RecordNetmask.ID:              RecordNetmask,
	RecordGateway.ID:              RecordGateway,
	RecordDHCP.ID:                 RecordDHCP,
	RecordFirmware.ID:             RecordFirmware,
	RecordPasswordEncryption.ID:   RecordPasswordEncryption,
	RecordPortSpeeds.ID:           RecordPortSpeeds,
	RecordPortMetrics.ID:          RecordPortMetrics,
	RecordCableTestResult.ID:      RecordCableTestResult,
	RecordVLANEngine.ID:           RecordVLANEngine,
	RecordVLANPort.ID:             RecordVLANPort,
	RecordVLAN802Q.ID:             RecordVLAN802Q,
	RecordPVIDs.ID:                RecordPVIDs,
	RecordQoSEngine.ID:            RecordQoSEngine,
	RecordQoSPolicies.ID:          RecordQoSPolicies,
	RecordBandwidthLimitsIn.ID:    RecordBandwidthLimitsIn,
	RecordBandwidthLimitsOut.ID:   RecordBandwidthLimitsOut,
	RecordBroadcastFilter.ID:      RecordBroadcastFilter,
	RecordBroadcastLimits.ID:      RecordBroadcastLimits,
	RecordPortMirroring.ID:        RecordPortMirroring,
	RecordPortCount.ID:            RecordPortCount,
	RecordIGMPSnoopingVLAN.ID:     RecordIGMPSnoopingVLAN,
	RecordMulticastFilter.ID:      RecordMulticastFilter,
	RecordIGMPHeaderValidation.ID: RecordIGMPHeaderValidation,
	RecordLoopDetection.ID:        RecordLoopDetection,
	RecordEndOfMessage.ID:         RecordEndOfMessage,
}

// RecordTypeNames maps the name of a record to a record type.
var RecordTypeByName = indexRecordTypeNames()

// indexRecordTypeNames builds an index of the record names.
func indexRecordTypeNames() map[string]*RecordType {
	recordNames := make(map[string]*RecordType, len(RecordTypeByID))

	for _, record := range RecordTypeByID {
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
	ID    RecordTypeID
	Len   uint16
	Value []uint8
}

// Type returns the type of the record.
func (r Record) Type() *RecordType {
	return RecordTypeByID[r.ID]
}

// Reflect returns a reflect.Value of the record's value.
func (r Record) Reflect() reflect.Value {
	rt := r.Type()
	if rt == nil {
		return reflect.ValueOf((*byte)(nil))
	}

	switch rt.Example.(type) {
	case string:
		return reflect.ValueOf(string(r.Value))
	case uint8:
		return reflect.ValueOf(uint8(r.Value[0]))
	case bool:
		return reflect.ValueOf(bool(r.Value[0] > 0))
	case net.HardwareAddr:
		return reflect.ValueOf(net.HardwareAddr(r.Value))
	case net.IP:
		return reflect.ValueOf(net.IP(r.Value))
	case []PortSpeed:
		return reflect.ValueOf(PortSpeed{
			ID:    r.Value[0],
			Speed: LinkStatus(r.Value[1]),
		})
	case []PortMetric:
		return reflect.ValueOf(PortMetric{
			ID:              r.Value[0],
			BytesReceived:   binary.BigEndian.Uint64(r.Value[1:9]),
			BytesSent:       binary.BigEndian.Uint64(r.Value[9:17]),
			ErrorsPacketCRC: binary.BigEndian.Uint64(r.Value[41:49]),
		})
	case PortMirroring:
		// I can for sure make out that uint8[0] is the destination
		// port. The other bits seem to be a bitmask. On my 8-port
		// switch, GS308E, uint8[2] seems to map to port 1 through 8,
		// where the most significant bit (7) corresponds to port 1
		// and the least significant bit (0) corresponds to port 8.
		// My educated guess is therefore that uint8[2] maps to port
		// 9 through 16, where the most significant bit (7) corresponds
		// to port 9 and the least significant bit (0) corresponds to
		// port 16. Below you can find a few examples that I configured
		// via the web UI to figure out the bitmask.
		//
		//   1. Port mirroring disabled: [0, 0, 0]
		//   2. Mirror ports 5 to port 6: [6, 0, 8]
		//   3. Mirror ports 3 and 7 to port 4: [4, 0, 34]
		//   4. Mirror ports 1 and 8 to port 2: [2, 0, 129]
		return reflect.ValueOf(PortMirroring{
			Destination: r.Value[0],
			Sources:     decodePortBitmask(r.Value[1:]),
		})
	case IGMPSnoopingVLAN:
		// If the value is 1, the IGMP snooping is enabled.
		if binary.BigEndian.Uint16(r.Value[0:2]) == 0x0001 {
			return reflect.ValueOf(IGMPSnoopingVLAN(binary.BigEndian.Uint16(r.Value[2:4])))
		}
		return reflect.ValueOf(IGMPSnoopingVLAN(0))
	case VLANEngine:
		return reflect.ValueOf(VLANEngine(r.Value[0]))
	case []VLANPort:
		return reflect.ValueOf(VLANPort{
			ID:    binary.BigEndian.Uint16(r.Value[0:2]),
			Ports: decodePortBitmask(r.Value[2:]),
		})
	case []VLAN802Q:
		portGroups := (r.Len - 2) / 2
		return reflect.ValueOf(VLAN802Q{
			ID:       binary.BigEndian.Uint16(r.Value[0:2]),
			Untagged: decodePortBitmask(r.Value[2 : 2+portGroups]),
			Tagged:   decodePortBitmask(r.Value[2+portGroups:]),
		})
	case []PVID:
		return reflect.ValueOf(PVID{
			ID:   r.Value[0],
			PVID: binary.BigEndian.Uint16(r.Value[1:3]),
		})
	case QoSEngine:
		return reflect.ValueOf(QoSEngine(r.Value[0]))
	case []QoSPolicy:
		return reflect.ValueOf(QoSPolicy{
			ID:       r.Value[0],
			Priority: QoSPriority(r.Value[1]),
		})
	case []BandwidthPolicy:
		return reflect.ValueOf(BandwidthPolicy{
			ID:    r.Value[0],
			Limit: BandwidthLimit(r.Value[4]),
		})
	default:
		// TODO: Parse CableTestResult.
		return reflect.ValueOf(r.Value)
	}
}

// decodePortBitmask takes a bitmask and returns a slice of ports.
func decodePortBitmask(portGroups []uint8) []uint8 {
	ports := make([]uint8, 0)
	// As previously outlined, each byte describes
	// a port group of 8 switch ports as each byte
	// has 8 bits.
	portGroupSize := 8
	portGroupCount := len(portGroups)
	for pg, portGroup := range portGroups {
		// The port groups are in reverse order, such that
		// port group containing the highest port number is
		// the mapped to the first byte. Thus, we find the
		// port offset by subtracting the current port group
		// from the total number of port groups. We also need
		// to subtract 1 because the last port group byte has
		// no offset and is therefore 0-indexed.
		portOffset := (portGroupCount - pg - 1) * 8
		for bit := portGroupSize - 1; bit >= 0; bit-- {
			if portGroup&(1<<bit) != 0 {
				// As bits are numbered starting from the least
				// significant bit and our ports are mapped in
				// reverse order, we need to subtract the bit
				// number port group size. Here, we do not need
				// to subtract 1 because the ports are 1-indexed.
				port := portOffset + (portGroupSize - bit)
				ports = append(ports, uint8(port))
			}
		}
	}
	return ports
}

// joinInts converts a slice of integers to a string.
func joinInts(ints []uint8, delimiter string) string {
	return strings.Trim(strings.ReplaceAll(fmt.Sprint(ints), " ", delimiter), "[]")
}
