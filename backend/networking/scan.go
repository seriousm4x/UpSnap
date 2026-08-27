package networking

import (
	"encoding/xml"
	"net"
	"os"
	"os/exec"
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

// MacToIP maps each scanned host's normalized mac address to its ipv4
// address, skipping hosts that lack either.
func (n Nmaprun) MacToIP() map[string]string {
	macToIp := make(map[string]string)
	for _, host := range n.Host {
		var hostIp, hostMac string
		for _, addr := range host.Address {
			switch addr.Addrtype {
			case "ipv4":
				hostIp = addr.Addr
			case "mac":
				hostMac = addr.Addr
			}
		}
		if parsedMac, err := net.ParseMAC(hostMac); hostIp != "" && err == nil {
			macToIp[parsedMac.String()] = hostIp
		}
	}
	return macToIp
}

func runNmap(scanRange string) (Nmaprun, error) {
	nmapOutput := Nmaprun{}

	nmap, err := exec.LookPath("nmap")
	if err != nil {
		return nmapOutput, err
	}

	timeout := os.Getenv("UPSNAP_SCAN_TIMEOUT")
	if timeout == "" {
		timeout = "500ms"
	}

	cmd := exec.Command(nmap, "-sn", "-oX", "-", scanRange, "--host-timeout", timeout, "--privileged")
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nmapOutput, err
	}

	if err := xml.Unmarshal(cmdOutput, &nmapOutput); err != nil {
		return nmapOutput, err
	}

	return nmapOutput, nil
}
