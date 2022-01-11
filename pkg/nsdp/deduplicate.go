package nsdp

func Deduplicate(existing []Device, current []Device) []Device {
	unique := make(map[string]*Device, 0)

	// Add existing devices to map and
	// identify them by their MAC.
	for _, device := range existing {
		unique[device.MAC.String()] = &device
	}

	// Overwrite existing devices that
	// have been updated by the current
	// operation.
	for _, device := range current {
		unique[device.MAC.String()] = &device
	}

	// Convert map to slice.
	devices := make([]Device, len(unique))
	i := 0
	for _, device := range unique {
		devices[i] = *device
	}

	return devices
}
