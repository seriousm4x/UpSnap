package networking

import (
	"errors"
	"net"
)

// DeviceSubnet returns the device's IPv4 subnet computed from its ip and netmask.
func DeviceSubnet(ipStr, maskStr string) (*net.IPNet, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil || ip.To4() == nil {
		return nil, errors.New("ip not a valid ipv4 address")
	}
	ip = ip.To4()

	mask := net.ParseIP(maskStr)
	if mask == nil || mask.To4() == nil {
		return nil, errors.New("subnet mask not a valid ipv4 address")
	}
	ipMask := net.IPMask(mask.To4())

	return &net.IPNet{IP: ip.Mask(ipMask), Mask: ipMask}, nil
}

// IsLocalSubnet reports whether one of the host's interface addresses is
// inside the given subnet, i.e. the subnet is directly attached and not routed.
func IsLocalSubnet(subnet *net.IPNet) bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipNet.IP.To4() != nil && subnet.Contains(ipNet.IP) {
			return true
		}
	}
	return false
}
