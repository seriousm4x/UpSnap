package networking

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

func TestPingDevice(t *testing.T) {
	testCases := []struct {
		name      string
		ip        string
		ping_cmd  string
		wantError bool
	}{
		// Valid case: all inputs correct
		{
			name:      "Valid Case",
			ip:        "8.8.8.8",
			ping_cmd:  "exit 0",
			wantError: false,
		},
		// Invalid ping command
		{
			name:      "Invalid Ping Command",
			ip:        "8.8.8.8",
			ping_cmd:  "exit 1",
			wantError: true,
		},
	}

	collection := &core.Collection{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			device := core.NewRecord(collection)
			device.Set("ip", tc.ip)
			device.Set("ping_cmd", tc.ping_cmd)

			_, err := PingDevice(device)
			if err == nil && tc.wantError {
				t.Errorf("Expected error but got none")
			} else if err != nil && !tc.wantError {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}

func TestCheckPort(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		port      string
		wantError bool
	}{
		// Valid case: all inputs correct
		{
			name:      "Valid Case",
			host:      "8.8.8.8",
			port:      "443",
			wantError: false,
		},
		// Invalid IP address
		{
			name:      "Invalid IP",
			host:      "256.1.1.1", // invalid
			port:      "443",
			wantError: true,
		},
		// Invalid port number
		{
			name:      "Invalid Port",
			host:      "8.8.8.8",
			port:      "70000", // invalid
			wantError: true,
		},
		// Port zero
		{
			name:      "Port Zero",
			host:      "8.8.8.8",
			port:      "0",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CheckPort(tc.host, tc.port)
			if err == nil && tc.wantError {
				t.Errorf("Expected error but got none")
			} else if err != nil && !tc.wantError {
				t.Errorf("Got unexpected error: %v", err)
			}
		})
	}
}
