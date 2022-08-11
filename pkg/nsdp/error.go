package nsdp

import "errors"

var (
	// ErrRecordTypeUnknown is returned if the record type is not supported.
	ErrRecordTypeUnknown = errors.New("record type unknown")
	// ErrNoDevicesFound is returned if no devices responded within the timeout period.
	ErrNoDevicesFound = errors.New("no devices found")
	// ErrInvalidDeviceIdentifier is returned if the device identifier is not a MAC or an IP.
	ErrInvalidDeviceIdentifier = errors.New("device identifier must be a MAC address or an IP address")
	// ErrInvalidEncryptionMode is returned if the encryption mode is not supported.
	ErrInvalidEncryptionMode = errors.New("invalid encryption mode")
	// ErrInvalidEndOfMessage is returned if the end of message is invalid.
	ErrInvalidEndOfMessage = errors.New("invalid end of message")
	// ErrInterfaceDown is returned if the interface is not connected and up.
	ErrInterfaceDown = errors.New("interface is down")
	// ErrInvalidInterfaceAddress is returned if the interface has no addresses.
	ErrInvalidInterfaceAddress = errors.New("invalid interface address")
	// ErrInvalidSelector is returned if the selector is invalid.
	ErrInvalidSelector = errors.New("invalid selector")
	// ErrInvalidResponse is returned if the response is invalid.
	ErrInvalidResponse = errors.New("invalid response")
	// ErrInvalidRecordLength is returned if the record length is invalid.
	ErrInvalidRecordLength = errors.New("invalid record length")
	// ErrInvalidPassword is returned if the password is invalid.
	ErrInvalidPassword = errors.New("invalid password")
	// ErrInvalidPasswordLockdown is returned if the password is invalid 3 times in a row.
	ErrInvalidPasswordLockdown = errors.New("device locked due to too many invalid password attempts")
	// ErrFailedNonceRetrieval is returned if the nonce retrieval failed.
	ErrFailedNonceRetrieval = errors.New("failed to retrieve password encryption nonce")
)

// ResponseCode describes the response code of a NSDP message.
type ResponseCode uint16

const (
	// ResponseCodeInvalidRecordLength is returned when the record length is invalid.
	ResponseCodeInvalidRecordLength ResponseCode = 0x0400
	// ResponseCodeInvalidPassword is returned when the password is invalid.
	ResponseCodeInvalidPassword ResponseCode = 0x0D00
	// ResponseCodeInvalidPasswordLockdown is returned if the user provides the wrong
	// password during 3 separate attempts. The devices enters a lockdown state for 30
	// minutes, during which it will not respond to any requests.
	ResponseCodeInvalidPasswordLockdown ResponseCode = 0x0E00
)
