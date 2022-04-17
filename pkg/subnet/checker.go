package subnet

import (
	"net"
)

// CheckSubnet checks is subnet contains IP.
func CheckSubnet(ip, subnet string) bool {
	if subnet == "" || ip == "" {
		return false
	}
	_, sNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return false
	}

	return sNet.Contains(net.ParseIP(ip))
}
