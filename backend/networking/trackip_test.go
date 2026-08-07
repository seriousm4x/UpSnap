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

func TestIsLocalSubnet(t *testing.T) {
	testCases := []struct {
		name    string
		ip      string
		netmask string
		want    bool
	}{
		// The loopback address is configured on every host
		{
			name:    "Loopback",
			ip:      "127.0.0.1",
			netmask: "255.0.0.0",
			want:    true,
		},
		// TEST-NET-1 (RFC 5737) is reserved for documentation and never assigned
		{
			name:    "Reserved Test Net",
			ip:      "192.0.2.1",
			netmask: "255.255.255.0",
			want:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			subnet, err := DeviceSubnet(tc.ip, tc.netmask)
			if err != nil {
				t.Fatalf("Got unexpected error: %v", err)
			}
			if got := IsLocalSubnet(subnet); got != tc.want {
				t.Errorf("IsLocalSubnet(%s): expected %v, got %v", subnet.String(), tc.want, got)
			}
		})
	}
}
