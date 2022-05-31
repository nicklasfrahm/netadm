package driver

import (
	"context"
)

// Attr is a key-value attribute pair that
// can be read from or written to the device.
type Attr struct {
	ID       int16
	Value    interface{}
	IsSlice  bool
	JSONPath string
}

// Service provides an RPC interface to a network API.
type Service interface {
	// Read reads one or more attributes from the device.
	Read(ctx context.Context, attr ...Attr) error
	// Write writes one or more attributes to the device.
	Write(ctx context.Context, attr ...Attr) error
}

// Driver is an abstract interface to interact with network devices and exposes
// RPC services to read or write configuration attributes. The RPC services are
// grouped by functionality as not every device may provide all of the services.
type Driver interface {
	// Scan scans for devices on the local network.
	Scan(options ...Option) ([]Device, error)

	// Service returns the named service or an error if it does not exist.
	Service(service string) (*Service, error)

	// VLAN returns an RPC service to manage virtual local area network settings.
	// VLAN() (*Service, error)

	// QoS returns an RPC service to manage Quality of Service settings.
	// QoS() (*Service, error)

	// TODO: Create link status service.
	// TODO: Create bandwidth limit service.
	// TODO: Create metric service.
	// TODO: Create port mirroring service.
}
