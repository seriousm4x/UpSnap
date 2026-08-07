package networking

import (
	"encoding/xml"
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
