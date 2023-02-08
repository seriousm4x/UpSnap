package networking

import (
	"errors"
	"fmt"
	"net"

	"github.com/mdlayher/wol"
	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/upsnap/logger"
)

func SendMagicPacket(device *models.Record) error {
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
	targetAddr := fmt.Sprintf("%s:%d", broadcastIp, 9)
	logger.Debug.Println(targetAddr)

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

func getBroadcastIp(ipStr, maskStr string) (string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", errors.New("invalid IP address")
	}
	mask := net.ParseIP(maskStr)
	if mask == nil {
		return "", errors.New("invalid subnet mask")
	}

	ip = ip.To4()
	mask = mask.To4()

	broadcast := make(net.IP, 4)
	for i := range ip {
		broadcast[i] = ip[i] | ^mask[i]
	}

	return broadcast.String(), nil
}
