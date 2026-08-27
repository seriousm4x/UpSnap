//go:build !linux

package networking

func NmapScan(scanRange string) (Nmaprun, error) {
	return runNmap(scanRange)
}
