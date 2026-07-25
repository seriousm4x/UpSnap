package networking

import (
	"errors"
	"fmt"
	"net"

	"github.com/mdlayher/wol"
	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/upsnap/logger"
)

func SendMagicPacket(device *core.Record) error {
  ip, err := ResolveToIPAddr(device.GetString("ip"))
  if err != nil {
    return err
  }
  mac := device.GetString("mac")
	netmask := device.GetString("netmask")
	password := device.GetString("password")

	// Validate MAC
	parsedMac, err := net.ParseMAC(mac)
	if err != nil {
		return err
	}
	var bytePassword []byte
  var broadcastIp string
	if len(password) == 0 || len(password) == 4 || len(password) == 6 {
		bytePassword = []byte(password)
	} else {
		return errors.New("error: password must be 0, 4 or 6 characters long")
	}

	// an IP was provided or FQDN resolved to an IP so calculate broadcast destination based on that
  if ip != "" {
    broadcastIp, err = getBroadcastIp(ip, netmask)
    if err != nil {
      return err
    }
  } else {
    // No IP available, so no broadcast IP can be calculated. This occurs if FQDN is .local and the device is powered off
    broadcastIp = ""
    logger.Debug.Println(fmt.Sprintf("%s did not resolve to an IP, broadcasting WOL packet to all interfaces", device.GetString("ip")))
  }

	// send wake
	if err := wakeUDP(broadcastIp, ip, parsedMac, bytePassword); err != nil {
		return err
	}
	return nil
}

func wakeUDP(broadcastIp string, deviceIp string, target net.HardwareAddr, password []byte) error {
	c, err := wol.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	var destinations []string
  
  // Broadcast IP available, add to destinations
  if broadcastIp != "" {
    // user-calculated broadcast packet to port 9 and alternative port 7
    destinations = append(destinations,
      fmt.Sprintf("%s:9", broadcastIp),
      fmt.Sprintf("%s:7", broadcastIp),
    )
  }

  if deviceIp != "" {
    // For routed network WOL, send a unicast packet to device IP on port 9 and alternative port 7
    destinations = append(destinations,
      fmt.Sprintf("%s:9", deviceIp),
      fmt.Sprintf("%s:7", deviceIp),
    )
  }
  
  // No IP for the device is available, broadcast WOL packet to all interfaces
  if deviceIp == "" {
    allBroadcastDests, err := getAllBroadcastDests()
    if err != nil {
      return err
    }
    destinations = append(destinations, allBroadcastDests...)
  }
  logger.Debug.Println(fmt.Sprintf("Broadcast destinations: %q", destinations))
	for _, dest := range destinations {
		if err := c.WakePassword(dest, target, password); err != nil {
			return err
		}
	}
	return nil
}

func getBroadcastIp(ipStr, maskStr string) (string, error) {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return "", errors.New("ip not a valid ipv4 address")
	}

	mask := net.ParseIP(maskStr).To4()
	if mask == nil {
		return "", errors.New("subnet mask not a valid ipv4 address")
	}

	broadcast := make(net.IP, 4)
	for i := range ip {
		broadcast[i] = ip[i] | ^mask[i]
	}

	return broadcast.String(), nil
}

// Returns a slice of all broadcast addresses for ports 9 and 7 using local interfaces
func getAllBroadcastDests() ([]string, error) {
  var allIps []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

  // Iterate through all the interfaces
	for _, iface := range interfaces {
		isRunning := iface.Flags&net.FlagRunning != 0
		isLoopback := iface.Flags&net.FlagLoopback != 0
		isNotBroadcast := iface.Flags&net.FlagBroadcast == 0
    // exclude interfaces that are not running, are loopback, or can't broadcast
		if !isRunning || isLoopback || isNotBroadcast{
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
    // Iterate through each interface's address
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				ip := ipNet.IP
        // Only use IPv4
				if ip.To4() != nil {
          bcastip,err := getBroadcastIp(ip.String(),net.IP(ipNet.Mask).String())
 					if err != nil {
            return nil, err
          }
          allIps = append(allIps, fmt.Sprintf("%s:9",bcastip), fmt.Sprintf("%s:7",bcastip))
				}
			}
		}
	}
  return allIps, nil
}
