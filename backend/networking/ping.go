package networking

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/pocketbase/pocketbase/core"
	probing "github.com/prometheus-community/pro-bing"
)

func PingDevice(device *core.Record) (bool, error) {
	ping_cmd := device.GetString("ping_cmd")
	if ping_cmd == "" {
		pinger, err := probing.NewPinger(device.GetString("ip"))
		if err != nil {
			return false, err
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
			return false, err
		}
		stats := pinger.Statistics()
		if stats.PacketLoss > 0 {
			return false, fmt.Errorf("packet loss is > 0: %v", stats.PacketLoss)
		} else {
			return true, nil
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

		return err == nil, err
	}
}

func CheckPort(host string, port string) (bool, error) {
	timeout := 500 * time.Millisecond
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return conn != nil, nil
}
