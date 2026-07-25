//go:build windows || darwin

package networking

import (
	"fmt"
	"net"
  "net/netip"
	"context"
  "regexp"
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
func ResolveToIPAddr(host string) (string, error){
  var fqdnRegex = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,63}\.?$`)
  ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()
  // If host is already an IP, simply return it
  _, err := netip.ParseAddr(host)
  if err == nil {
    return host, nil
  }
  // Validate FQDN format
  if ! fqdnRegex.MatchString(host) {
    return "", fmt.Errorf("%s is not a properly formatted FQDN",host)
  }
  ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip4", host)
  // return results
  if len(ips) > 0 {
    return ips[0].String(), nil
  }
	// no results, but .local FQDN (expected when mDNS device is offline), don't return an error with empty result
  if strings.HasSuffix(host, ".local") {
    return "", nil
  } else {
    // no results, but FQDN is not .local. Should have resolved, return an error so the user can fix the FQDN or DNS entries
    return "", err
  }
}
