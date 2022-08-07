package nsdp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
)

// Set provides a simplified way to set configuration keys on devices.
// TOOD: Implement this!
func Set(id string, keys []string, options ...Option) ([]Device, error) {
	// Get operation options.
	opts, err := GetDefaultOptions().Apply(options...)
	if err != nil {
		return nil, err
	}

	selector := NewSelector()
	// Allow usage of keyword "all" to select all devices.
	if id != "all" {
		// Check if the device is identified via its IP address.
		ip := net.ParseIP(id)
		if ip == nil {
			// Fall back to MAC address device identification.
			mac, err := net.ParseMAC(id)
			if err != nil {
				return nil, errors.New("device identifier must be a MAC address or an IP address")
			}
			selector.SetMAC(&mac)
		} else {
			selector.SetIP(&ip)
		}
	}

	// Check if all keys are valid.
	for i, key := range keys {
		// Normalize key name.
		keys[i] = strings.ToLower(key)

		// Check if key is valid.
		if RecordTypeByName[key] == nil {
			return nil, fmt.Errorf(`unknown configuration key "%s"`, key)
		}
	}

	// Create slice to hold results.
	devices := make([]Device, 0)

	// Retry operation if retries is greater than 0.
	for i := uint(0); i <= opts.Retries; i++ {
		// Create new message.
		request := NewMessage(ReadRequest)

		// Add request records.
		for _, key := range keys {
			request.Records = append(request.Records, Record{
				ID: RecordTypeByName[key].ID,
			})
		}

		// Create context to handle timeout.
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()

		// Run scan for devices.
		devs, err := RequestDevices(opts.InterfaceName, request,
			WithContext(ctx),
			WithSelector(selector),
		)
		if err != nil {
			return nil, err
		}

		// Deduplicate results from all attempts.
		devices = DeduplicateDevices(devices, devs)
	}

	// Check if any devices were found.
	if len(devices) == 0 {
		return nil, errors.New("no switches found")
	}

	return devices, nil
}
