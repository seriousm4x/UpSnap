//go:build linux || freebsd

package networking

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/holoplot/go-avahi"
)

var (
	avahiServer *avahi.Server
	avahiErr    error
	avahiOnce   sync.Once
)

// getAvahiServer initializes and returns a shared, long-lived Avahi client.
// It handles high throughput gracefully by using a single D-Bus socket connection.
func getAvahiServer() (*avahi.Server, error) {
	avahiOnce.Do(func() {
		// Using the cached system bus here is safe because we never call server.Close() on it.
		conn, err := dbus.SystemBus()
		if err != nil {
			avahiErr = fmt.Errorf("failed to connect to dbus: %w", err)
			return
		}
		server, err := avahi.ServerNew(conn)
		if err != nil {
			avahiErr = fmt.Errorf("failed to create avahi server: %w", err)
			return
		}
		avahiServer = server
	})
	return avahiServer, avahiErr
}

/*
Takes one parameter 'host' and returns its IP.  If 'host' is already formatted
as an IP, simply return it. If 'host' is formatted as an FQDN, attempt to
resolve it with mDNS (if .local suffix) or else use a DNS lookup with
net.DefaultResolver

Linux and FreeBSD require dbus and avahi for mDNS lookups
*/
func ResolveToIPAddr(host string) (string, error) {
	// If host is already an IP, simply return it
	addr, err := netip.ParseAddr(host)
	if err == nil && addr.Is4() {
		return host, nil
	}
	if err := ValidateFQDN(host); err != nil {
		return "", err
	}
	// attempt mDNS query first
	if strings.HasSuffix(host, ".local") {
		//		conn, err := dbus.Dial("unix:path=/var/run/dbus/system_bus_socket")
		server, err := getAvahiServer()
		if err != nil {
			return "", fmt.Errorf("Failed to create connection: %w", err)
		}
		// defer conn.Close()

		// server, err := avahi.ServerNew(conn)
		// if err != nil {
		// 	return "", fmt.Errorf("Failed to create avahi server: %w", err)
		// }
		// defer server.Close()

		// set a context with 2 secs and only wait that long for an mDNS response.  otherwise, avahi takes 5 secs
		ctx1, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		type result struct {
			address string
			err     error
		}
		ch := make(chan result, 1)

		go func() {
			resolved, err := server.ResolveHostName(
				avahi.InterfaceUnspec,
				avahi.ProtoInet,
				host,
				avahi.ProtoInet,
				0,
			)
			if err != nil {
				ch <- result{"", err}
				return
			}
			// mDNS result received
			ch <- result{resolved.Address, nil}
		}()

		select {
		case <-ctx1.Done():
			// 2 second timeout reached,
			return "", nil

		case res := <-ch:
			if res.err != nil {
				var dbusErr *dbus.Error
				// in theory we shouldn't encounter this error since the 2 second context always beats the race
				if errors.As(res.err, &dbusErr) && dbusErr.Name == "org.freedesktop.Avahi.TimeoutError" {
					return "", nil
				}
				return "", fmt.Errorf("%s: %w", host, res.err)
			}
			return res.address, nil
		}
	}
	// Not .local, attempt unicast (regular) DNS lookup
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	ips, err := net.DefaultResolver.LookupNetIP(ctx2, "ip4", host)
	if err != nil {
		return "", err
	}
	if len(ips) > 0 {
		// only return one result
		return ips[0].String(), nil
	}
	// No results, but FQDN should have resolved since it's not .local. Return an error.
	return "", err
}
