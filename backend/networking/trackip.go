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

// ValidateScannableSubnet returns nil if the subnet is suitable for an nmap
// scan, or an error explaining why it isn't.
func ValidateScannableSubnet(subnet *net.IPNet) error {
	ones, bits := subnet.Mask.Size()
	if ones == 0 && bits == 0 {
		// Size() returns 0,0 for a non-contiguous mask, which has no cidr
		// notation to hand to nmap
		return errors.New("netmask is not contiguous")
	}
	if ones == bits {
		return errors.New("a single address has nothing to scan")
	}
	if ones < 16 {
		return errors.New("subnets larger than a /16 take too long to scan")
	}

	// arp only answers on-link, so the subnet must overlap one of the
	// host's attached ipv4 networks for a scan to find mac addresses
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		attached, ok := addr.(*net.IPNet)
		if !ok || attached.IP.To4() == nil {
			continue
		}
		// subnets are either disjoint or nested, so checking containment
		// in both directions covers overlap
		if attached.Contains(subnet.IP) || subnet.Contains(attached.IP) {
			return nil
		}
	}
	return errors.New("subnet does not overlap a directly attached network")
}
