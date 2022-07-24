package nsdp

// DeduplicateDevices merges two slices by retaining only all unique items.
func DeduplicateDevices(existing []Device, current []Device) []Device {
	unique := make(map[string]*Device, 0)

	// Add existing objects to map and identify them by their MAC.
	// Use indexing instead of "range" to obtain immutable pointers.
	for i := 0; i < len(existing); i++ {
		unique[existing[i].MAC.String()] = &existing[i]
	}

	// Overwrite existing objects that have been updated by the current operation.
	// Use indexing instead of "range" to obtain immutable pointers.
	for i := 0; i < len(current); i++ {
		unique[current[i].MAC.String()] = &current[i]
	}

	// Convert map to slice.
	objects := make([]Device, len(unique))
	i := 0
	for _, object := range unique {
		objects[i] = *object
		i += 1
	}

	return objects
}

// DeduplicateMessages merges two slices by retaining only all unique items.
func DeduplicateMessages(existing []Message, current []Message) []Message {
	unique := make(map[[6]uint8]*Message, 0)

	// Add existing objects to map and identify them by their MAC.
	// Use indexing instead of "range" to obtain immutable pointers.
	for i := 0; i < len(existing); i++ {
		unique[existing[i].Header.ServerMAC] = &existing[i]
	}

	// Overwrite existing objects that have been updated by the current operation.
	// Use indexing instead of "range" to obtain immutable pointers.
	for i := 0; i < len(current); i++ {
		unique[current[i].Header.ServerMAC] = &current[i]
	}

	// Convert map to slice.
	objects := make([]Message, len(unique))
	i := 0
	for _, object := range unique {
		objects[i] = *object
		i += 1
	}

	return objects
}
