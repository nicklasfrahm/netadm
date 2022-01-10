package nsdp

import (
	"context"
)

const (
	// ClientPort is the port the client host
	// uses to send and receive messages.
	ClientPort = 63321
	// ServerPort is the port the device server
	// uses to send and receive messages.
	ServerPort = 63322
)

// Options defines the configuration of an operation of this library.
type Options struct {
	Context context.Context
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
