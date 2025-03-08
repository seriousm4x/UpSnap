//go:build !windows

package networking

import "os"

func isRoot() bool {
	return os.Geteuid() == 0
}
