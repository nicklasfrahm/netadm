package nsdp

import (
	"context"
	"errors"
	"net"
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
	SelectorAll = &Selector{MAC: &net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}
)

// Selector is a type that allows it to
// target specific devices with a message.
type Selector struct {
	MAC *net.HardwareAddr
}

// Options defines the configuration of an operation of this library.
type Options struct {
	Context  context.Context
	Selector *Selector
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

// WithMAC specifies which devices to send
// the message to.
func WithMAC(mac *net.HardwareAddr) Option {
	return func(o *Options) error {
		if mac == nil {
			return errors.New("MAC must not be empty")
		}

		o.Selector = &Selector{MAC: mac}
		return nil
	}
}
