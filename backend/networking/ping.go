package networking

import (
	"context"
	"errors"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// pingCmdTimeout is how long a custom ping command may run. A command that
// takes longer is killed with its child processes and the device counts as
// offline, so a stalled command such as a hung ssh cannot block the next ping.
var pingCmdTimeout = 10 * time.Second

// runPingCmd runs a custom ping command and reports whether it exited with 0.
func runPingCmd(shell string, shellArg string, pingCmd string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pingCmdTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, shell, shellArg, pingCmd)
	SetProcessAttributes(cmd)
	cmd.Cancel = func() error {
		return killProcessGroup(cmd.Process)
	}
	cmd.WaitDelay = time.Second

	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		if ctx.Err() != nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

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

func CheckPort(host string, port string) (bool, error) {
	timeout := 500 * time.Millisecond
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		// treat "host unreachable", "connection refused" and "timeout" as no error
		if netErr, ok := errors.AsType[*net.OpError](err); ok {
			if errors.Is(netErr.Err, syscall.EHOSTUNREACH) ||
				errors.Is(netErr.Err, syscall.ECONNREFUSED) ||
				netErr.Timeout() {
				return false, nil
			}
		}
		return false, err
	}
	defer conn.Close()
	return conn != nil, nil
}
