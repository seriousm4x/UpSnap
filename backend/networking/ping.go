package networking

import (
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/pocketbase/pocketbase/core"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/seriousm4x/upsnap/logger"
)

func PingDevice(device *core.Record) bool {
	ping_cmd := device.GetString("ping_cmd")
	if ping_cmd == "" {
		pinger, err := probing.NewPinger(device.GetString("ip"))
		if err != nil {
			logger.Error.Println(err)
			return false
		}
		pinger.Count = 1
		pinger.Timeout = 500 * time.Millisecond
		privileged, err := strconv.ParseBool(os.Getenv("UPSNAP_PING_PRIVILEGED"))
		if err != nil {
			privileged = true
		}
		pinger.SetPrivileged(privileged)
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
	} else {
		var shell string
		var shell_arg string
		if runtime.GOOS == "windows" {
			shell = "cmd"
			shell_arg = "/C"
		} else {
			shell = "/bin/sh"
			shell_arg = "-c"
		}

		cmd := exec.Command(shell, shell_arg, ping_cmd)
		err := cmd.Run()

		return err == nil
	}
}

func CheckPort(host string, port string) bool {
	timeout := 500 * time.Millisecond
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return conn != nil
}
