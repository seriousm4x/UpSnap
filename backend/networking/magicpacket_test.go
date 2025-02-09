package networking

import (
	"net"
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

type Device struct {
	IP                 string
	MAC                string
	Netmask            string
	Password           string
	ExpectFailPassword bool
}

func TestSendMagicPacket(t *testing.T) {
	testCases := []struct {
		name      string
		ip        string
		mac       string
		netmask   string
		password  string
		wantError bool
	}{
		// Valid case: all inputs correct
		{
			name:      "Valid Case",
			ip:        "192.168.1.100",
			mac:       "00:11:22:33:44:55",
			netmask:   "255.255.255.0",
			password:  "secret",
			wantError: false,
		},
		// Invalid MAC address
		{
			name:      "Invalid MAC",
			ip:        "192.168.1.100",
			mac:       "invalid",
			netmask:   "255.255.255.0",
			password:  "secret",
			wantError: true,
		},
		// Password too short
		{
			name:      "Password Too Short",
			ip:        "192.168.1.100",
			mac:       "00:11:22:33:44:55",
			netmask:   "255.255.255.0",
			password:  "s", // length 1
			wantError: true,
		},
		// Password too long
		{
			name:      "Password Too Long",
			ip:        "192.168.1.100",
			mac:       "00:11:22:33:44:55",
			netmask:   "255.255.255.0",
			password:  "password", // length 9
			wantError: true,
		},
		// Invalid IP address
		{
			name:      "Invalid IP",
			ip:        "256.1.1.1", // invalid
			mac:       "00:11:22:33:44:55",
			netmask:   "255.255.255.0",
			password:  "secret",
			wantError: true,
		},
		// Invalid netmask
		{
			name:      "Invalid Netmask",
			ip:        "192.168.1.100",
			mac:       "00:11:22:33:44:55",
			netmask:   "300.255.255.0", // invalid
			password:  "secret",
			wantError: true,
		},
	}

	collection := &core.Collection{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			device := core.NewRecord(collection)
			device.Set("ip", tc.ip)
			device.Set("mac", tc.mac)
			device.Set("netmask", tc.netmask)
			device.Set("password", tc.password)

			err := SendMagicPacket(device)
			if tc.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Got unexpected error: %v", err)
				}
			}
		})
	}
}

func TestWakeUDP(t *testing.T) {
	testCases := []struct {
		name        string
		broadcastIp string
		targetMac   string
		password    string
		wantError   bool
	}{
		// Valid case: all inputs correct
		{
			name:        "Valid Case",
			broadcastIp: "192.168.1.255",
			targetMac:   "00:11:22:33:44:55",
			password:    "secret",
			wantError:   false,
		},
		// Invalid MAC address
		{
			name:        "Invalid MAC",
			broadcastIp: "192.168.1.255",
			targetMac:   "invalid",
			password:    "secret",
			wantError:   true,
		},
		// Password too short
		{
			name:        "Password Too Short",
			broadcastIp: "192.168.1.255",
			targetMac:   "00:11:22:33:44:55",
			password:    "s", // length 1
			wantError:   true,
		},
		// Password too long
		{
			name:        "Password Too Long",
			broadcastIp: "192.168.1.255",
			targetMac:   "00:11:22:33:44:55",
			password:    "password", // length 9
			wantError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mac, err := net.ParseMAC(tc.targetMac)
			if err != nil && !tc.wantError {
				t.Errorf("Unexpected error parsing MAC: %v", err)
			} else if err != nil {
				return
			}

			// Now test the wakeUDP function
			err = wakeUDP(tc.broadcastIp, mac, []byte(tc.password))
			if tc.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			}
		})
	}
}

func TestGetBroadcastIp(t *testing.T) {
	testCases := []struct {
		name      string
		ip        string
		netmask   string
		wantIp    string
		wantError bool
	}{
		// Valid case: all inputs correct
		{
			name:      "Valid Case",
			ip:        "192.168.1.100",
			netmask:   "255.255.255.0",
			wantIp:    "192.168.1.255",
			wantError: false,
		},
		// Invalid IP address
		{
			name:      "Invalid IP",
			ip:        "256.1.1.1", // invalid
			netmask:   "255.255.255.0",
			wantIp:    "",
			wantError: true,
		},
		// Invalid netmask
		{
			name:      "Invalid Netmask",
			ip:        "192.168.1.100",
			netmask:   "300.255.255.0", // invalid
			wantIp:    "",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ip := net.ParseIP(tc.ip)
			if ip == nil && !tc.wantError {
				t.Errorf("Unexpected error parsing IP: %v", ip)
			} else if ip == nil {
				return
			}

			mask := net.ParseIP(tc.netmask)
			if mask == nil && !tc.wantError {
				t.Errorf("Unexpected error parsing netmask: %v", mask)
			} else if mask == nil {
				return
			}

			broadcast, err := getBroadcastIp(tc.ip, tc.netmask)
			if tc.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if broadcast != tc.wantIp {
					t.Errorf("Broadcast IP mismatch: expected %s, got %s", tc.wantIp, broadcast)
				}
			}
		})
	}
}
