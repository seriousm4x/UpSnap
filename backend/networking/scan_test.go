package networking

import (
	"encoding/xml"
	"testing"
)

// Trimmed from real `nmap -sn -oX -` output: one host with a mac (scanned
// with raw socket privileges), one without (the scanning host itself).
const nmapSample = `<?xml version="1.0" encoding="UTF-8"?>
<nmaprun scanner="nmap" args="nmap -sn -oX - 10.1.10.0/24" version="7.94">
<host><status state="up" reason="arp-response"/>
<address addr="10.1.10.1" addrtype="ipv4"/>
<address addr="AA:BB:CC:DD:EE:FF" addrtype="mac" vendor="Ubiquiti Networks"/>
<hostnames></hostnames>
</host>
<host><status state="up" reason="localhost-response"/>
<address addr="10.1.10.9" addrtype="ipv4"/>
<hostnames></hostnames>
</host>
<runstats><finished time="1754000000" exit="success"/></runstats>
</nmaprun>`

func TestNmaprunUnmarshal(t *testing.T) {
	nmapOutput := Nmaprun{}
	if err := xml.Unmarshal([]byte(nmapSample), &nmapOutput); err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	if len(nmapOutput.Host) != 2 {
		t.Fatalf("Host count mismatch: expected 2, got %d", len(nmapOutput.Host))
	}

	first := nmapOutput.Host[0].Address
	if len(first) != 2 {
		t.Fatalf("Address count mismatch: expected 2, got %d", len(first))
	}
	if first[0].Addrtype != "ipv4" || first[0].Addr != "10.1.10.1" {
		t.Errorf("Unexpected ipv4 address: %+v", first[0])
	}
	if first[1].Addrtype != "mac" || first[1].Addr != "AA:BB:CC:DD:EE:FF" || first[1].Vendor != "Ubiquiti Networks" {
		t.Errorf("Unexpected mac address: %+v", first[1])
	}

	second := nmapOutput.Host[1].Address
	if len(second) != 1 || second[0].Addrtype != "ipv4" {
		t.Errorf("Unexpected addresses for mac-less host: %+v", second)
	}
}

func TestNmaprunMacToIP(t *testing.T) {
	nmapOutput := Nmaprun{}
	if err := xml.Unmarshal([]byte(nmapSample), &nmapOutput); err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	macToIp := nmapOutput.MacToIP()
	if len(macToIp) != 1 {
		t.Fatalf("Entry count mismatch: expected 1, got %d", len(macToIp))
	}
	// the mac address is normalized to lower case, the mac-less host is skipped
	if ip := macToIp["aa:bb:cc:dd:ee:ff"]; ip != "10.1.10.1" {
		t.Errorf("Unexpected ip for aa:bb:cc:dd:ee:ff: %q", ip)
	}
}
