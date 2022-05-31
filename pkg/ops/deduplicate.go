package ops

import (
	"github.com/nicklasfrahm/netadm/pkg/driver"
)

// DeduplicateDevices accepts a single slice and returns a new slice with all unique items.
func DeduplicateDevices(duplicates []driver.Device) []driver.Device {
	unique := make(map[string]*driver.Device, 0)

	// Add existing objects to map and
	// identify them by their MAC.
	for _, object := range duplicates {
		unique[object.System.MAC.String()] = &object
	}

	// Convert map to slice.
	objects := make([]driver.Device, len(unique))
	i := 0
	for _, object := range unique {
		objects[i] = *object
		i += 1
	}

	return objects
}
