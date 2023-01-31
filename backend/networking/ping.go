package networking

import (
	"net"
	"time"

	"github.com/pocketbase/pocketbase/models"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/seriousm4x/upsnap/logger"
)

func PingDevice(device *models.Record) bool {
	pinger, err := probing.NewPinger(device.GetString("ip"))
	if err != nil {
		logger.Error.Println(err)
		return false
	}
	pinger.Count = 1
	pinger.Timeout = 500 * 1000000 // 500ms
	pinger.SetPrivileged(true)
	err = pinger.Run()
	if err != nil {
		logger.Error.Println(err)
		return false
	}
	stats := pinger.Statistics()
	if stats.PacketLoss > 0 {
		return false
	} else {
		return true
	}
}

func CheckPort(host string, port string) bool {
	timeout := 500 * time.Millisecond
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}
