package nsdp

import (
	"math"
	"net"
)

// FixedLengthXOR creates a fixed-length XOR encryption of the given data.
// The length of the hash will be the minimum length of the data and the key.
func FixedLengthXOR(data []byte, key []byte) []byte {
	// Suprisingly, Golang does not have a built-in `min` function for integers,
	// so we have to resort to this ugliness.
	length := int(math.Min(float64(len(data)), float64(len(key))))
	hash := make([]byte, length)

	for i := 0; i < length; i++ {
		hash[i] = data[i] ^ key[i]
	}

	return hash
}

// EncryptPassword encrypts a password for transmission. I have no idea what kind of dark magic this is, see below.
// Reference: https://github.com/rhulme/ProSafeLinux/blob/bdd8f55a67f2b5930377ed3a055bd29a9a6fc019/psl_class.py#L538
func EncryptPassword(encryptionMode EncryptionMode, macAddr net.HardwareAddr, nonce []byte, password []byte) ([]byte, error) {
	switch encryptionMode {
	case EncryptionModeNone:
		return []byte(password), nil

	case EncryptionModeSimple:
		return FixedLengthXOR([]byte(password), []byte("NtgrSmartSwitchRock")), nil

	case EncryptionModeHash32:
	case EncryptionModeHash64:
		// This ensures that the algorithm also works
		// if the MAC address is longer than 6 bytes.
		mac := make([]byte, 6)
		copy(mac, macAddr)

		// Seed hash with device MAC address.
		hash := []byte{
			mac[1] ^ mac[5],
			mac[0] ^ mac[4],
			mac[2] ^ mac[3],
			mac[4] ^ mac[5],
		}

		// Seed the hash with the nonce.
		hash[0] ^= nonce[3] ^ nonce[2]
		hash[1] ^= nonce[3] ^ nonce[1]
		hash[2] ^= nonce[0] ^ nonce[2]
		hash[3] ^= nonce[0] ^ nonce[1]

		if encryptionMode == EncryptionModeHash32 {
			for i := 0; i < int(math.Min(float64(len(password)), 16)); i++ {
				j := 0
				if i < 4 || i > 7 {
					j = ((i + 3) % 4)
					j ^= (j / 2)
				} else {
					j = 3 - (i % 4)
				}

				hash[j] ^= password[i]
			}
		} else {
			hash = append(hash, hash[0], hash[1], hash[2], hash[3])

			hash[6] ^= password[0]

			for i := 0; i < len(password); i++ {
				hash[i/3] ^= password[i]

				if i < 6 && i%2 != 0 {
					hash[7] ^= password[i]
				}
			}
		}

		return hash, nil
	}

	return nil, ErrInvalidEncryptionMode
}
