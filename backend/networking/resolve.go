//go:build windows || darwin

package networking

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"time"
)

/*
Takes one parameter 'host' and returns its IP.  If 'host' is already formatted
as an IP, simply return it. If 'host' is formatted as an FQDN, attempt to
resolve it with mDNS (if .local suffix) or else use a DNS lookup with
net.DefaultResolver

Windows and MacOS can resolve mDNS/DNS with one call to net.DefaultResolver
*/
func ResolveToIPAddr(host string) (string, error) {
	// If host is already an IP, simply return it
	addr, err := netip.ParseAddr(host)
	if err == nil && addr.Is4() {
		return host, nil
	}
	if err := ValidateFQDN(host); err != nil {
		return "", fmt.Errorf("%q is the host", host)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip4", host)
	if err != nil {
		if strings.HasSuffix(host, ".local") {
			var dnsErr *net.DNSError
			isNotFound := errors.As(err, &dnsErr) && dnsErr.IsNotFound
			isTimeout := errors.Is(err, context.DeadlineExceeded) || ctx.Err() == context.DeadlineExceeded
			if isTimeout || isNotFound {
				return "", nil
			}
		}
		return "", err
	}
	if len(ips) > 0 {
		// only return one result
		return ips[0].String(), nil
	}
	// No results, but it's expected for a .local FQDN device that is offline. Don't return an error.
	if strings.HasSuffix(host, ".local") {
		return "", nil
	} else {
		// No results, but FQDN should have resolved since it's not .local. Return an error.
		return "", fmt.Errorf("No IP address found for %s", host)
	}
}
