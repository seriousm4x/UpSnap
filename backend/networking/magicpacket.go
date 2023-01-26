package networking

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/mdlayher/wol"
	"github.com/pocketbase/pocketbase/models"
)

func SendMagicPacket(device *models.Record) error {
	ip := device.GetString("ip")
	mac := device.GetString("mac")
	netmask := device.GetString("netmask")
	password := device.GetString("password")

	// parse inputs
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("error: invalid ip address")
	}
	parsedMac, err := net.ParseMAC(mac)
	if err != nil {
		return err
	}
	parsedNetmask := net.ParseIP(netmask)
	if parsedNetmask == nil {
		return fmt.Errorf("error: invalid netmask")
	}
	var bytePassword []byte
	if len(password) == 0 || len(password) == 4 || len(password) == 6 {
		bytePassword = []byte(password)
	} else {
		return fmt.Errorf("error: password must be 0, 4 or 6 characters long")
	}

	// get target addr
	ipNet := &net.IPNet{
		IP:   parsedIP,
		Mask: parsedNetmask.DefaultMask(),
	}
	broadcastIp, err := getBroadcastIp(ipNet)
	if err != nil {
		return err
	}
	targetAddr := fmt.Sprintf("%s:%d", broadcastIp, 9)

	// send wake via udp port 9
	if err := wakeUDP(targetAddr, parsedMac, bytePassword); err != nil {
		return err
	}

	return nil
}

func wakeUDP(addr string, target net.HardwareAddr, password []byte) error {
	c, err := wol.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()
	return c.WakePassword(addr, target, password)
}

func getBroadcastIp(n *net.IPNet) (net.IP, error) {
	if n.IP.To4() == nil {
		return net.IP{}, errors.New("does not support IPv6 addresses")
	}
	ip := make(net.IP, len(n.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return ip, nil
}
