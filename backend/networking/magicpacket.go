package networking

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

func SendMagicPacket(mac, netmask string) error {
	macAddr, err := hex.DecodeString(strings.ReplaceAll(mac, ":", ""))
	if err != nil {
		return fmt.Errorf("error decoding mac address: %v", err)
	}

	// Create an address to listen on
	address := net.UDPAddr{
		IP:   net.ParseIP(netmask),
		Port: 7,
	}

	// create a new UDP packet conn
	conn, err := net.DialUDP("udp", nil, &address)
	if err != nil {
		return fmt.Errorf("error creating UDP packet conn: %v", err)
	}
	defer conn.Close()

	// create a magic packet
	var magicPacket bytes.Buffer
	for i := 0; i < 6; i++ {
		magicPacket.WriteByte(0xff)
	}
	for i := 0; i < 16; i++ {
		magicPacket.Write(macAddr)
	}

	// send the magic packet
	_, err = conn.Write(magicPacket.Bytes())
	if err != nil {
		return fmt.Errorf("error sending magic packet: %v", err)
	}
	return nil
}
