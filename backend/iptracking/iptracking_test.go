package iptracking

import (
	"encoding/xml"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/seriousm4x/upsnap/networking"

	// register the app migrations so test apps get the real schema
	_ "github.com/seriousm4x/upsnap/migrations"
)

// newTestApp returns a throwaway app with all migrations applied.
func newTestApp(t *testing.T) *tests.TestApp {
	t.Helper()
	app, err := tests.NewTestApp(t.TempDir())
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	t.Cleanup(app.Cleanup)
	return app
}

func buildDevice(t *testing.T, app core.App, name, ip, netmask, mac string, trackIp bool) *core.Record {
	t.Helper()
	collection, err := app.FindCollectionByNameOrId("devices")
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	device := core.NewRecord(collection)
	device.Set("name", name)
	device.Set("ip", ip)
	device.Set("netmask", netmask)
	device.Set("mac", mac)
	device.Set("track_ip", trackIp)
	return device
}

func newDevice(t *testing.T, app core.App, name, ip, netmask, mac string, trackIp bool) *core.Record {
	t.Helper()
	device := buildDevice(t, app, name, ip, netmask, mac, trackIp)
	if err := app.Save(device); err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	return device
}

// newLegacyDevice saves a device bypassing field validation, like rows
// created before the ip and mac format rules existed.
func newLegacyDevice(t *testing.T, app core.App, name, ip, netmask, mac string, trackIp bool) *core.Record {
	t.Helper()
	device := buildDevice(t, app, name, ip, netmask, mac, trackIp)
	if err := app.SaveNoValidate(device); err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	return device
}

// stubScan replaces the nmap seam with a canned scan reporting the given
// mac addresses at the given ips, and returns the list of scanned ranges.
func stubScan(t *testing.T, macToHostIp map[string]string) *[]string {
	t.Helper()
	var sample strings.Builder
	sample.WriteString("<nmaprun>")
	for mac, ip := range macToHostIp {
		fmt.Fprintf(&sample, `<host><address addr=%q addrtype="ipv4"/><address addr=%q addrtype="mac"/></host>`, ip, mac)
	}
	sample.WriteString("</nmaprun>")
	scan := networking.Nmaprun{}
	if err := xml.Unmarshal([]byte(sample.String()), &scan); err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	var scanned []string
	orig := nmapScan
	nmapScan = func(scanRange string) (networking.Nmaprun, error) {
		scanned = append(scanned, scanRange)
		return scan, nil
	}
	t.Cleanup(func() { nmapScan = orig })
	return &scanned
}

func mustSubnet(t *testing.T, cidr string) *net.IPNet {
	t.Helper()
	_, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	return subnet
}

// The loopback /16 passes the scannable checks on any host, making it the
// only subnet a test can scan without depending on the environment.
const (
	testSubnet  = "127.0.0.0/16"
	testNetmask = "255.255.0.0"
)

func TestTrackOneSubnetSkipsUntrackedDevices(t *testing.T) {
	app := newTestApp(t)
	device := newDevice(t, app, "untracked", "127.0.0.50", testNetmask, "AA:BB:CC:DD:EE:01", false)
	stubScan(t, map[string]string{"AA:BB:CC:DD:EE:01": "127.0.0.99"})

	if err := TrackOneSubnet(app, mustSubnet(t, testSubnet)); err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	fresh, err := app.FindRecordById("devices", device.Id)
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	if ip := fresh.GetString("ip"); ip != "127.0.0.50" {
		t.Errorf("Ip of device without track_ip changed to %s", ip)
	}
}

func TestTrackOneSubnetSkipPaths(t *testing.T) {
	app := newTestApp(t)
	saves := 0
	app.OnRecordUpdate("devices").BindFunc(func(e *core.RecordEvent) error {
		saves++
		return e.Next()
	})

	testCases := []struct {
		name   string
		ip     string
		mac    string
		legacy bool
	}{
		// A device mac that can't be parsed never matches
		{
			name:   "Invalid Mac",
			ip:     "127.0.0.50",
			mac:    "not-a-mac",
			legacy: true,
		},
		// A device mac missing from the scan results is left alone
		{
			name: "Mac Not In Scan",
			ip:   "127.0.0.50",
			mac:  "AA:BB:CC:DD:EE:02",
		},
		// A device ip without a computable subnet is never moved
		{
			name:   "Invalid Device Ip",
			ip:     "999.999.999.999",
			mac:    "AA:BB:CC:DD:EE:01",
			legacy: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			create := newDevice
			if tc.legacy {
				create = newLegacyDevice
			}
			device := create(t, app, tc.name, tc.ip, testNetmask, tc.mac, true)
			stubScan(t, map[string]string{"AA:BB:CC:DD:EE:01": "127.0.0.99"})

			if err := TrackOneSubnet(app, mustSubnet(t, testSubnet)); err != nil {
				t.Fatalf("Got unexpected error: %v", err)
			}

			fresh, err := app.FindRecordById("devices", device.Id)
			if err != nil {
				t.Fatalf("Got unexpected error: %v", err)
			}
			if ip := fresh.GetString("ip"); ip != tc.ip {
				t.Errorf("Ip mismatch: expected %s, got %s", tc.ip, ip)
			}
			if saves != 0 {
				t.Errorf("Expected no record updates, got %d", saves)
			}
		})
	}
}

func TestTrackOneSubnetUnscannable(t *testing.T) {
	app := newTestApp(t)
	// TEST-NET-1 (RFC 5737) is never assigned to an interface
	scanned := stubScan(t, nil)

	if err := TrackOneSubnet(app, mustSubnet(t, "192.0.2.0/24")); err == nil {
		t.Error("Expected error but got none")
	}
	if len(*scanned) != 0 {
		t.Errorf("Expected no scans of an unscannable subnet, got %v", *scanned)
	}
}

func TestTrackOneSubnetSameSubnetGuard(t *testing.T) {
	app := newTestApp(t)

	testCases := []struct {
		name    string
		ip      string
		netmask string
		mac     string
		scanMac string
		foundIp string
		wantIp  string
	}{
		// A move within the device's own subnet is applied, with the
		// stored dashed lower case mac matching the scanned
		// colon-separated upper case one
		{
			name:    "Within Own Subnet",
			ip:      "127.0.0.50",
			netmask: testNetmask,
			mac:     "aa-bb-cc-dd-01-01",
			scanMac: "AA:BB:CC:DD:01:01",
			foundIp: "127.0.0.99",
			wantIp:  "127.0.0.99",
		},
		// The device's own subnet bounds the move, not the scanned one
		{
			name:    "Narrower Subnet Within Scan",
			ip:      "127.0.5.5",
			netmask: "255.255.255.0",
			mac:     "AA:BB:CC:DD:01:02",
			foundIp: "127.0.5.9",
			wantIp:  "127.0.5.9",
		},
		// A found ip outside the device's narrower subnet is not applied,
		// even though it is inside the scanned subnet
		{
			name:    "Outside Own Subnet",
			ip:      "127.0.5.5",
			netmask: "255.255.255.0",
			mac:     "AA:BB:CC:DD:01:03",
			foundIp: "127.0.6.9",
			wantIp:  "127.0.5.5",
		},
		// A device from a disjoint subnet is never moved into the scanned one
		{
			name:    "Disjoint Subnet",
			ip:      "10.9.9.9",
			netmask: "255.255.255.0",
			mac:     "AA:BB:CC:DD:01:04",
			foundIp: "127.0.0.99",
			wantIp:  "10.9.9.9",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			device := newDevice(t, app, tc.name, tc.ip, tc.netmask, tc.mac, true)
			scanMac := tc.scanMac
			if scanMac == "" {
				scanMac = tc.mac
			}
			stubScan(t, map[string]string{scanMac: tc.foundIp})

			if err := TrackOneSubnet(app, mustSubnet(t, testSubnet)); err != nil {
				t.Fatalf("Got unexpected error: %v", err)
			}

			fresh, err := app.FindRecordById("devices", device.Id)
			if err != nil {
				t.Fatalf("Got unexpected error: %v", err)
			}
			if ip := fresh.GetString("ip"); ip != tc.wantIp {
				t.Errorf("Ip mismatch: expected %s, got %s", tc.wantIp, ip)
			}
		})
	}
}

// TrackAllSubnets scans each unique scannable subnet once, silently skips the
// rest, and updates any tracked device found at a new ip within its own
// subnet — even one whose own subnet couldn't be scanned.
func TestTrackAllSubnets(t *testing.T) {
	app := newTestApp(t)
	first := newDevice(t, app, "first", "127.0.0.10", testNetmask, "AA:BB:CC:DD:02:01", true)
	second := newDevice(t, app, "second", "127.0.0.20", testNetmask, "AA:BB:CC:DD:02:02", true)
	// TEST-NET-1 (RFC 5737) is never assigned to an interface
	unscannable := newDevice(t, app, "unscannable", "192.0.2.5", "255.255.255.0", "AA:BB:CC:DD:02:03", true)
	// a /8 subnet is too large to scan but contains the found ip
	wide := newDevice(t, app, "wide", "127.0.0.60", "255.0.0.0", "AA:BB:CC:DD:02:04", true)
	scanned := stubScan(t, map[string]string{
		"AA:BB:CC:DD:02:01": "127.0.0.11",
		"AA:BB:CC:DD:02:02": "127.0.0.21",
		"AA:BB:CC:DD:02:04": "127.0.0.61",
	})

	TrackAllSubnets(app)

	if len(*scanned) != 1 || (*scanned)[0] != testSubnet {
		t.Errorf("Expected a single scan of %s, got %v", testSubnet, *scanned)
	}
	wantIps := map[string]string{
		first.Id:       "127.0.0.11",
		second.Id:      "127.0.0.21",
		unscannable.Id: "192.0.2.5",
		wide.Id:        "127.0.0.61",
	}
	for id, wantIp := range wantIps {
		fresh, err := app.FindRecordById("devices", id)
		if err != nil {
			t.Fatalf("Got unexpected error: %v", err)
		}
		if ip := fresh.GetString("ip"); ip != wantIp {
			t.Errorf("Ip mismatch for %s: expected %s, got %s", fresh.GetString("name"), wantIp, ip)
		}
	}
}
