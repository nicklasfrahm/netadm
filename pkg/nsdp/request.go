package nsdp

// RequestMessage is a high-level API that sends messages via the
// low-level Send API and returns the results as a slice of Messages.
func RequestMessage(ifaceName string, request *Message, options ...Option) ([]Message, error) {
	// Get operation options.
	opts, err := GetDefaultOptions().Apply(options...)
	if err != nil {
		return nil, err
	}

	// Fetch network interface.
	iface, err := GetInterface(ifaceName)
	if err != nil {
		return nil, err
	}

	// Inject the client MAC address into the request.
	request.Header.ClientMAC = MACMarshalBinary(&iface.HardwareAddr)

	// Overwrite server MAC address if the provided Selector
	// is different than the default SelectorAll.
	if opts.Selector != SelectorAll {
		request.Header.ServerMAC = MACMarshalBinary(opts.Selector.MAC)
	}

	// Send message to broadcast address.
	// TODO: Can we do this also via unicast?
	responses, err := Send(opts.Context, iface, request)
	if err != nil {
		return nil, err
	}

	return responses, nil
}

// RequestDevice is a high-level API that sends messages via the high-level
// RequestMessage API and returns the results as a slice of Devices.
func RequestDevice(ifaceName string, request *Message, options ...Option) ([]Device, error) {
	// Use the RequestMessage function to get the responses.
	responses, err := RequestMessage(ifaceName, request, options...)
	if err != nil {
		return nil, err
	}

	// Convert responses to devices.
	devices := make([]Device, len(responses))
	for i, response := range responses {
		// This is safe because we previously allocated the slice.
		if err := devices[i].UnmarshalMessage(&response); err != nil {
			return nil, err
		}
	}

	return devices, nil
}
