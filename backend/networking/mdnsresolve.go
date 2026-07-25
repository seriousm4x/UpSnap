//go:build linux || freebsd

package networking

import (
	"fmt"
	"net"
  "net/netip"
	"context"
	"errors"
	"strings"
  "regexp"
  "time"
	"github.com/godbus/dbus/v5"
	"github.com/holoplot/go-avahi"
)

/*
Takes one parameter 'host' and returns its IP.  If 'host' is already formatted 
as an IP, simply return it. If 'host' is formatted as an FQDN, attempt to 
resolve it with mDNS (if .local suffix) or else use a DNS lookup with 
net.DefaultResolver 

Linux and FreeBSD require dbus and avahi for mDNS lookups
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
	// attempt mDNS query first
	if strings.HasSuffix(host, ".local") {
    conn, err := dbus.SystemBus()
    if err != nil {
      return "", err
    }
    server, err := avahi.ServerNew(conn)
    if err != nil {
      return "", err
    }
    resolved, err := server.ResolveHostName(
      avahi.InterfaceUnspec,
      avahi.ProtoInet,
      host,
      avahi.ProtoInet,
      0,
    )
    if err != nil {
      var dbusErr dbus.Error
      if errors.As(err, &dbusErr) {
        switch dbusErr.Name {
        // Timeout error means no device replied to the query, which is normal when a device is powered off, don't treat it as an error
        case "org.freedesktop.Avahi.TimeoutError":
          break
        default:
          return "", fmt.Errorf("%s: %w", host, err)
        }
      } else {
          return "", fmt.Errorf("%s: %w", host, err)
      }
    } else {
      return resolved.Address, nil
    }
  }
  // mDNS did not return a result, attempt DNS lookup
	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip4", host)
  // return result from DNS lookup
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
