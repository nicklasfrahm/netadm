package ops

import (
	"fmt"
	"os"

	"github.com/nicklasfrahm/netadm/pkg/driver"
	"github.com/nicklasfrahm/netadm/pkg/netgear"
)

var drivers = map[string]driver.Driver{
	"netgear": netgear.NewDriver(),
}

func Scan(options ...driver.Option) ([]driver.Device, error) {
	// Get options for the current operation.
	opts, err := driver.GetDefaultOptions().Apply(options...)
	if err != nil {
		return nil, err
	}

	// TODO: Use options.
	_ = opts

	// TODO: Disable drivers based on deny-list.

	devices := make([]driver.Device, 0)
	// TODO: Scan for devices using different drivers.

	for _, driver := range drivers {
		devs, err := driver.Scan()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}

		devices = append(devices, devs...)
	}

	return devices, nil
}
