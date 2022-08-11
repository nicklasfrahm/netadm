package nsdp

import (
	"context"
	"fmt"
	"net"
)

// Set provides a simplified way to set configuration keys on devices.
func Set(id string, values map[string]string, options ...Option) ([]Device, error) {
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
				return nil, ErrInvalidDeviceIdentifier
			}
			selector.SetMAC(&mac)
		} else {
			selector.SetIP(&ip)
		}
	}

	// Check if all keys are valid.
	for key := range values {
		// Check if key is valid.
		if RecordTypeByName[key] == nil {
			return nil, fmt.Errorf(`unknown configuration key "%s"`, key)
		}
	}

	// Prepare password for authentication.
	devices, err := Get(id, []string{"mac", "ip", "passwordencryption"}, options...)
	if err != nil {
		return nil, err
	}
	encryptionMode := devices[0].PasswordEncryption
	id = devices[0].IP.String()

	nonce := make([]byte, 4)
	if encryptionMode == EncryptionModeHash32 || encryptionMode == EncryptionModeHash64 {
		devs, err := Get(id, []string{"mac", "passwordnonce"}, options...)
		if err != nil {
			return nil, err
		}

		if len(nonce) == 0 {
			return nil, ErrFailedNonceRetrieval
		}

		copy(nonce, devs[0].PasswordNonce)
	}

	encryptedPassword, err := EncryptPassword(encryptionMode, devices[0].MAC, nonce, []byte(opts.Password))
	if err != nil {
		return nil, err
	}

	// Create new message.
	request := NewMessage(WriteRequest)

	if encryptionMode == EncryptionModeNone || encryptionMode == EncryptionModeSimple {
		request.Records = append(request.Records, Record{
			ID:    RecordPassword.ID,
			Value: encryptedPassword,
			Len:   uint16(len(encryptedPassword)),
		})
	}

	if encryptionMode == EncryptionModeHash32 || encryptionMode == EncryptionModeHash64 {
		request.Records = append(request.Records,
			Record{
				ID:    RecordPasswordHash.ID,
				Value: encryptedPassword,
				Len:   uint16(len(encryptedPassword)),
			},
		)
	}

	// Add request records.
	for key, value := range values {
		encodedValue := []byte(value)

		request.Records = append(request.Records, Record{
			ID:    RecordTypeByName[key].ID,
			Value: encodedValue,
			Len:   uint16(len(encodedValue)),
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

	// Check if any devices were found.
	if len(devices) == 0 {
		return nil, ErrNoDevicesFound
	}

	return devices, nil
}
