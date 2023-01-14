package networking

import (
	"encoding/xml"
	"net"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
)

type Nmaprun struct {
	Host []struct {
		Address []struct {
			Addr     string `xml:"addr,attr" binding:"required"`
			Addrtype string `xml:"addrtype,attr" binding:"required"`
			Vendor   string `xml:"vendor,attr"`
		} `xml:"address"`
	} `xml:"host"`
}

func ScanNetwork(c echo.Context) error {
	scanRange := os.Getenv("UPSNAP_SCAN_RANGE")
	_, _, err := net.ParseCIDR(scanRange)
	if err != nil {
		return apis.NewBadRequestError(err.Error(), nil)
	}

	cmdOutput, err := exec.Command(
		"sudo", "nmap", "-sP", "-oX", "-", scanRange,
		"--host-timeout", "500ms").Output()
	if err != nil {
		return err
	}

	nmapOutput := Nmaprun{}
	if err := xml.Unmarshal(cmdOutput, &nmapOutput); err != nil {
		return err
	}

	type Device struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
		MAC  string `json:"mac"`
	}
	data := []Device{}

	for _, host := range nmapOutput.Host {
		dev := Device{}
		for _, addr := range host.Address {
			if addr.Addrtype == "ipv4" {
				dev.IP = addr.Addr
			} else if addr.Addrtype == "mac" {
				dev.MAC = addr.Addr
			}
			if addr.Vendor != "" {
				dev.Name = addr.Vendor
			}
		}
		if dev.IP == "" || dev.MAC == "" {
			continue
		}
		if dev.Name == "" {
			dev.Name = "Unknown"
		}
		data = append(data, dev)
	}

	return c.JSON(http.StatusOK, data)
}
