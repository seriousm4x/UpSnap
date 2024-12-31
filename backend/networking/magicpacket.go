package networking

import (
	"errors"
	"fmt"
	"net"

	"github.com/mdlayher/wol"
	"github.com/pocketbase/pocketbase/core"
)

func SendMagicPacket(device *core.Record) error {
	ip := device.GetString("ip")
	mac := device.GetString("mac")
	netmask := device.GetString("netmask")
	password := device.GetString("password")

	// parse inputs
	parsedMac, err := net.ParseMAC(mac)
	if err != nil {
		return err
	}
	var bytePassword []byte
	if len(password) == 0 || len(password) == 4 || len(password) == 6 {
		bytePassword = []byte(password)
	} else {
		return fmt.Errorf("error: password must be 0, 4 or 6 characters long")
	}

	// get target addr
	broadcastIp, err := getBroadcastIp(ip, netmask)
	if err != nil {
		return err
	}

	// send wake via udp port 9
	if err := wakeUDP(broadcastIp, parsedMac, bytePassword); err != nil {
		return err
	}

	return nil
}

func wakeUDP(broadcastIp string, target net.HardwareAddr, password []byte) error {
	c, err := wol.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	// send 4 magic packets to different addresses to enhance change of waking up
	destinations := []struct {
		Addr     string
		Target   net.HardwareAddr
		Password []byte
	}{
		{
			// default user-calculated broadcast to port 9
			Addr:     fmt.Sprintf("%s:%d", broadcastIp, 9),
			Target:   target,
			Password: password,
		},
		{
			// user-calculated broadcast to port alternative port 7
			Addr:     fmt.Sprintf("%s:%d", broadcastIp, 7),
			Target:   target,
			Password: password,
		},
		{
			// broadcast to port 9
			Addr:     "255.255.255.255:9",
			Target:   target,
			Password: password,
		},
		{
			// broadcast to alternative port 7
			Addr:     "255.255.255.255:7",
			Target:   target,
			Password: password,
		},
	}

	for _, dest := range destinations {
		if err := c.WakePassword(dest.Addr, dest.Target, dest.Password); err != nil {
			return err
		}
	}

	// send packet to default broadcast address
	return c.WakePassword("255.255.255.255:9", target, password)
}

func getBroadcastIp(ipStr, maskStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", errors.New("ip not a valid ipv4 address")
	}
	ip = ip.To4()
	if ip == nil {
		return "", errors.New("ip not a valid ipv4 address")
	}

	mask := net.ParseIP(maskStr)
	if mask == nil {
		return "", errors.New("subnet mask not a valid ipv4 address")
	}
	mask = mask.To4()
	ip = ip.To4()
	if ip == nil {
		return "", errors.New("subnet mask not a valid ipv4 address")
	}

	broadcast := make(net.IP, 4)
	for i := range ip {
		broadcast[i] = ip[i] | ^mask[i]
	}

	return broadcast.String(), nil
}
