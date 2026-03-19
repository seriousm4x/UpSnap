//go:build linux

package pb

import (
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"kernel.org/pub/linux/libs/security/libcap/cap"
)

func HandlerScan(e *core.RequestEvent) error {
	// check if nmap installed
	nmap, err := exec.LookPath("nmap")
	if err != nil {
		return apis.NewBadRequestError(err.Error(), nil)
	}

	// check if scan range is valid
	allPrivateSettings, err := e.App.FindAllRecords("settings_private")
	if err != nil {
		return err
	}
	settingsPrivate := allPrivateSettings[0]
	scanRange := settingsPrivate.GetString("scan_range")
	_, ipNet, err := net.ParseCIDR(scanRange)
	if err != nil {
		return apis.NewBadRequestError(err.Error(), nil)
	}

	// run nmap
	timeout := os.Getenv("UPSNAP_SCAN_TIMEOUT")
	if timeout == "" {
		timeout = "500ms"
	}
	orig := cap.GetProc()
	defer orig.SetProc() // restore original caps on exit.

	c, err := orig.Dup()
	if err != nil {
		return fmt.Errorf("Failed to dup existing capabilities: %v", err)
	}

	if on, _ := c.GetFlag(cap.Permitted, cap.NET_RAW); !on {
		return fmt.Errorf("unable to get NET_RAW permissions")
	}

	if err := c.SetFlag(cap.Effective, true, cap.NET_RAW); err != nil {
		return fmt.Errorf("unable to set NET_RAW capability effective")
	}

	if err := c.SetFlag(cap.Inheritable, true, cap.NET_RAW); err != nil {
		return fmt.Errorf("unable to set NET_RAW capability inheritable")
	}

	if err := c.SetProc(); err != nil {
		return fmt.Errorf("unable to raise NET_RAW capability")
	}

	if err := cap.SetAmbient(true, cap.NET_RAW); err != nil {
		return fmt.Errorf("unable to set NET_RAW capability ambient")
	}

	cmd := exec.Command(nmap, "-sn", "-oX", "-", scanRange, "--host-timeout", timeout, "--privileged")
	cmdOutput, err := cmd.Output()
	if err != nil {
		return err
	}

	// unmarshal xml
	nmapOutput := Nmaprun{}
	if err := xml.Unmarshal(cmdOutput, &nmapOutput); err != nil {
		return err
	}

	type Device struct {
		Name      string `json:"name"`
		IP        string `json:"ip"`
		MAC       string `json:"mac"`
		MACVendor string `json:"mac_vendor"`
	}

	// extract info from struct into data
	type Response struct {
		Netmask string   `json:"netmask"`
		Devices []Device `json:"devices"`
	}
	res := Response{}
	var nm []string
	for _, octet := range ipNet.Mask {
		nm = append(nm, strconv.Itoa(int(octet)))
	}
	res.Netmask = strings.Join(nm, ".")

	for _, host := range nmapOutput.Host {
		dev := Device{}
		for _, addr := range host.Address {
			if addr.Addrtype == "ipv4" {
				dev.IP = addr.Addr
			} else if addr.Addrtype == "mac" {
				dev.MAC = addr.Addr
			}
			if addr.Vendor != "" {
				dev.MACVendor = addr.Vendor
			} else {
				dev.MACVendor = "Unknown"
			}
		}

		if dev.IP == "" || dev.MAC == "" {
			continue
		}

		names, err := net.LookupAddr(dev.IP)
		if err != nil || len(names) == 0 {
			dev.Name = dev.MACVendor
		} else {
			dev.Name = strings.TrimSuffix(names[0], ".")
		}

		if dev.Name == "" && dev.MACVendor == "" {
			continue
		}

		res.Devices = append(res.Devices, dev)
	}

	return e.JSON(http.StatusOK, res)
}
