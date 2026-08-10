// Package iptracking updates the stored ip addresses of opted-in devices by
// scanning their local subnets for their mac addresses.
package iptracking

import (
	"net"
	"sync/atomic"

	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/upsnap/logger"
	"github.com/seriousm4x/upsnap/networking"
)

var sweepRunning atomic.Bool

// nmapScan is a seam for tests to stub out the privileged nmap invocation
var nmapScan = networking.NmapScan

// TrackAllSubnets scans the local subnets of devices with ip tracking enabled
// and updates their ip address if their mac address is found at a different
// one. Non-scannable subnets are silently skipped, but their devices can
// still be updated by another subnet's scan that finds their mac address,
// as long as the new ip stays within the device's own subnet.
func TrackAllSubnets(app core.App) {
	// skip if the previous sweep is still running
	if !sweepRunning.CompareAndSwap(false, true) {
		return
	}
	defer sweepRunning.Store(false)

	devices, err := app.FindRecordsByFilter("devices", "track_ip = true", "", 0, 0)
	if err != nil {
		logger.Error.Println(err)
		return
	}

	// collect the unique scannable subnets of the tracked devices
	subnets := make(map[string]*net.IPNet)
	for _, device := range devices {
		subnet, err := networking.DeviceSubnet(device.GetString("ip"), device.GetString("netmask"))
		if err != nil {
			logger.Error.Println("Ip tracking for", device.GetString("name")+":", err)
			continue
		}
		if networking.ValidateScannableSubnet(subnet) != nil {
			continue
		}
		subnets[subnet.String()] = subnet
	}

	for cidr, subnet := range subnets {
		if err := TrackOneSubnet(app, subnet); err != nil {
			logger.Error.Println("Ip tracking scan for", cidr+":", err)
		}
	}
}

// TrackOneSubnet runs an nmap scan of the given subnet and updates the ip
// address of any ip-tracked device whose mac address is found at a new
// address within its own subnet. Returns an error for a non-scannable subnet.
func TrackOneSubnet(app core.App, subnet *net.IPNet) error {
	if err := networking.ValidateScannableSubnet(subnet); err != nil {
		return err
	}

	scan, err := nmapScan(subnet.String())
	if err != nil {
		return err
	}
	macToIp := scan.MacToIP()

	devices, err := app.FindRecordsByFilter("devices", "track_ip = true", "", 0, 0)
	if err != nil {
		return err
	}

	for _, device := range devices {
		parsedMac, err := net.ParseMAC(device.GetString("mac"))
		if err != nil {
			continue
		}
		newIp, ok := macToIp[parsedMac.String()]
		if !ok || newIp == device.GetString("ip") {
			continue
		}
		// only move a device within its own subnet, so that a scan can't
		// relocate devices tracked on other subnets that see the same mac
		deviceSubnet, err := networking.DeviceSubnet(device.GetString("ip"), device.GetString("netmask"))
		if err != nil || !deviceSubnet.Contains(net.ParseIP(newIp)) {
			continue
		}
		logger.Info.Println("Ip tracking: updating", device.GetString("name"), "from", device.GetString("ip"), "to", newIp)
		device.Set("ip", newIp)
		// only write the changed ip field to avoid clobbering concurrent
		// status updates from the ping and wake/shutdown cronjobs
		device.IgnoreUnchangedFields(true)
		if err := app.Save(device); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}
	}
	return nil
}
