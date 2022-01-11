package nsdp

import (
	"errors"
	"net"
)

// GetInterface fetches the interface based on the provided
// interface name if it up.
func GetInterface(ifaceName string) (*net.Interface, error) {
	// Fetch specified interface by name.
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, errors.New("interface unknown")
	}

	// Check if interface is up.
	if iface.Flags&net.FlagUp == 0 {
		return nil, errors.New("interface is down")
	}

	return iface, nil
}

// GetInterfaceIP fetches the IPv4 address of the interface.
func GetInterfaceIP(iface *net.Interface) (*net.IP, error) {
	// Check if interface has addresses.
	addresses, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	if len(addresses) == 0 {
		return nil, errors.New("interface has no address")
	}

	// Select IPv4 interface address.
	var ip *net.IP
	var address net.Addr
	for _, addr := range addresses {
		// Check if address is IPv4.
		ipNet, ok := addr.(*net.IPNet)
		if ok && ipNet.IP.To4() != nil {
			ip = &ipNet.IP
			address = addr
			break
		}
	}

	// Check if interface has a valid IPv4 address.
	if address == nil {
		return nil, errors.New("interface has no valid IPv4 address")
	}

	return ip, nil
}
