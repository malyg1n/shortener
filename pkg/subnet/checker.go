package subnet

import (
	"net"
)

// CheckSubnet checks is subnet contains IP.
func CheckSubnet(ip, subnet string) bool {
	if subnet == "" {
		return false
	}

	if ip == "" {
		return false
	}

	_, sNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return false
	}

	sIP := net.ParseIP(ip)

	return sNet.Contains(sIP)
}
