package networking

import (
	"testing"
)

func TestDeviceSubnet(t *testing.T) {
	testCases := []struct {
		name      string
		ip        string
		netmask   string
		wantCidr  string
		wantError bool
	}{
		// Valid case: address is masked down to the network
		{
			name:      "Valid Case",
			ip:        "192.168.1.100",
			netmask:   "255.255.255.0",
			wantCidr:  "192.168.1.0/24",
			wantError: false,
		},
		// Valid case: wider mask
		{
			name:      "Valid /16",
			ip:        "10.1.10.142",
			netmask:   "255.255.0.0",
			wantCidr:  "10.1.0.0/16",
			wantError: false,
		},
		// Valid case: host mask
		{
			name:      "Valid /32",
			ip:        "192.168.1.5",
			netmask:   "255.255.255.255",
			wantCidr:  "192.168.1.5/32",
			wantError: false,
		},
		// Invalid IP address
		{
			name:      "Invalid IP",
			ip:        "256.1.1.1", // invalid
			netmask:   "255.255.255.0",
			wantCidr:  "",
			wantError: true,
		},
		// Invalid netmask
		{
			name:      "Invalid Netmask",
			ip:        "192.168.1.100",
			netmask:   "300.255.255.0", // invalid
			wantCidr:  "",
			wantError: true,
		},
		// IPv6 address is not a valid device ip
		{
			name:      "IPv6 IP",
			ip:        "fe80::1",
			netmask:   "255.255.255.0",
			wantCidr:  "",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			subnet, err := DeviceSubnet(tc.ip, tc.netmask)
			if tc.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Got unexpected error: %v", err)
				} else if subnet.String() != tc.wantCidr {
					t.Errorf("Subnet mismatch: expected %s, got %s", tc.wantCidr, subnet.String())
				}
			}
		})
	}
}

func TestValidateScannableSubnet(t *testing.T) {
	testCases := []struct {
		name      string
		ip        string
		netmask   string
		wantError bool
	}{
		// A /32 has nothing to scan
		{
			name:      "Host Mask",
			ip:        "192.168.1.5",
			netmask:   "255.255.255.255",
			wantError: true,
		},
		// A non-contiguous mask has no cidr notation to scan
		{
			name:      "Non-contiguous Mask",
			ip:        "10.0.1.5",
			netmask:   "255.0.255.0",
			wantError: true,
		},
		// Anything larger than a /16 takes too long to scan
		{
			name:      "Too Large",
			ip:        "10.0.0.1",
			netmask:   "255.0.0.0",
			wantError: true,
		},
		// The loopback address as a /16 is within limits and always attached
		{
			name:      "Local /16",
			ip:        "127.0.0.1",
			netmask:   "255.255.0.0",
			wantError: false,
		},
		// A sub-block of the attached loopback /8 is on-link even though
		// it doesn't contain the host's own 127.0.0.1
		{
			name:      "Nested In Attached",
			ip:        "127.0.1.5",
			netmask:   "255.255.255.0",
			wantError: false,
		},
		// TEST-NET-1 (RFC 5737) is never assigned to an interface
		{
			name:      "No Overlap",
			ip:        "192.0.2.1",
			netmask:   "255.255.255.0",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			subnet, err := DeviceSubnet(tc.ip, tc.netmask)
			if err != nil {
				t.Fatalf("Got unexpected error: %v", err)
			}
			if err := ValidateScannableSubnet(subnet); (err != nil) != tc.wantError {
				t.Errorf("ValidateScannableSubnet(%s): expected error=%v, got %v", subnet.String(), tc.wantError, err)
			}
		})
	}
}
