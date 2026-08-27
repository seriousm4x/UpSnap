package networking

import (
	"errors"
	"net"
	"os"
	"syscall"
	"time"
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
