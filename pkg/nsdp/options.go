package nsdp

import (
	"context"
	"net"
	"time"
)

const (
	// ClientPort is the port the client host
	// uses to send and receive messages.
	ClientPort = 63321
	// ServerPort is the port the device server
	// uses to send and receive messages.
	ServerPort = 63322
)

var (
	// SelectorAll is a selector that will cause
	// all devices to respond to the mesage.
	SelectorAll = &Selector{MAC: &net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, IP: &net.IPv4bcast}
)

// Selector is a type that allows it to
// target specific devices with a message.
type Selector struct {
	MAC *net.HardwareAddr
	IP  *net.IP
}

// NewSelector returns a new Selector matching all devices.
func NewSelector() *Selector {
	return &Selector{
		MAC: &net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		IP:  &net.IPv4bcast,
	}
}

// SetMAC sets the MAC address of the selector and returns the selector.
func (s *Selector) SetMAC(mac *net.HardwareAddr) *Selector {
	s.MAC = mac
	return s
}

// SetIP sets the IP address of the selector and returns the selector.
func (s *Selector) SetIP(ip *net.IP) *Selector {
	s.IP = ip
	return s
}

// Options defines the configuration of an operation of this library.
type Options struct {
	Context       context.Context
	Selector      *Selector
	InterfaceName string
	Timeout       time.Duration
	Retries       uint
	Password      string
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
		Context:  context.Background(),
		Selector: SelectorAll,
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

// WithSelector allows to select devices either
// by their MAC address or by their IP address.
func WithSelector(selector *Selector) Option {
	return func(o *Options) error {
		if selector == nil {
			return ErrInvalidSelector
		}
		o.Selector = selector
		return nil
	}
}

// WithTimeout supplies a custom timeout for the operation.
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) error {
		o.Timeout = timeout
		return nil
	}
}

// WithRetries supplies the number of retries for the operation.
func WithRetries(retries uint) Option {
	return func(o *Options) error {
		o.Retries = retries
		return nil
	}
}

// WithInterfaceName supplies the name of the network interface
// to use for the operation.
func WithInterfaceName(name string) Option {
	return func(o *Options) error {
		o.InterfaceName = name
		return nil
	}
}

// WithPassword supplies the password for the operation.
func WithPassword(password string) Option {
	return func(o *Options) error {
		o.Password = password
		return nil
	}
}
