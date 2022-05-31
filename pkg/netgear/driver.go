package netgear

import (
	"github.com/nicklasfrahm/netadm/pkg/driver"
)

// Driver is a struct that implements the driver.Driver interface.
type Driver struct {
}

func NewDriver() driver.Driver {
	return &Driver{}
}

func (d *Driver) Scan(options ...driver.Option) ([]driver.Device, error) {
	// TODO: Implement this.

	return nil, nil
}

func (d *Driver) Service(service string) (*driver.Service, error) {
	// TODO: Implement this.

	return nil, nil
}
