package networking

import (
	"errors"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/pocketbase/pocketbase/core"
	probing "github.com/prometheus-community/pro-bing"
)

func isNoRouteOrDownError(err error) bool {
	opErr, ok := err.(*net.OpError)
	if !ok {
		return false
	}
	syscallErr, ok := opErr.Err.(*os.SyscallError)
	if !ok {
		return false
	}
	return syscallErr.Err == syscall.EHOSTUNREACH || syscallErr.Err == syscall.EHOSTDOWN
}

func PingDevice(device *core.Record) (bool, error) {
	ping_cmd := device.GetString("ping_cmd")
	if ping_cmd == "" {
		pinger, err := probing.NewPinger(device.GetString("ip"))
		if err != nil {
			return false, err
		}
		pinger.Count = 1
		pinger.Timeout = 500 * time.Millisecond

		privileged := isRoot()
		privilegedEnv := os.Getenv("UPSNAP_PING_PRIVILEGED")
		if privilegedEnv != "" {
			privileged, err = strconv.ParseBool(privilegedEnv)
			if err != nil {
				privileged = false
			}
		}
		pinger.SetPrivileged(privileged)

		err = pinger.Run()
		if err != nil {
			if isNoRouteOrDownError(err) {
				return false, nil
			}
			return false, err
		}
		stats := pinger.Statistics()
		return stats.PacketLoss == 0, nil
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
		// treat "host unreachable" and "timeout" as no error
		var netErr *net.OpError
		if errors.As(err, &netErr) {
			if errors.Is(netErr.Err, syscall.EHOSTUNREACH) || netErr.Timeout() {
				return false, nil
			}
		}
		return false, err
	}
	defer conn.Close()
	return conn != nil, nil
}
