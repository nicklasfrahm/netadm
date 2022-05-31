package ops

// DeduplicateDevices merges two slices by retaining only all unique items.
func DeduplicateDevices(existing []Device, current []Device) []Device {
	unique := make(map[string]*Device, 0)

	// Add existing objects to map and
	// identify them by their MAC.
	for _, object := range existing {
		unique[object.MAC.String()] = &object
	}

	// Overwrite existing objects that
	// have been updated by the current
	// operation.
	for _, object := range current {
		unique[object.MAC.String()] = &object
	}

	// Convert map to slice.
	objects := make([]Device, len(unique))
	i := 0
	for _, object := range unique {
		objects[i] = *object
	}

	return objects
}

// DeduplicateMessages merges two slices by retaining only all unique items.
func DeduplicateMessages(existing []Message, current []Message) []Message {
	unique := make(map[string]*Message, 0)

	// Add existing objects to map and
	// identify them by their MAC.
	for _, object := range existing {
		unique[string(object.Header.ClientMAC[:])] = &object
	}

	// Overwrite existing objects that
	// have been updated by the current
	// operation.
	for _, object := range current {
		unique[string(object.Header.ClientMAC[:])] = &object
	}

	// Convert map to slice.
	objects := make([]Message, len(unique))
	i := 0
	for _, object := range unique {
		objects[i] = *object
	}

	return objects
}
