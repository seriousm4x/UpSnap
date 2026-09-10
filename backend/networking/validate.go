package networking

import (
	"fmt"
	"net/netip"
	"regexp"
	"strings"

	"golang.org/x/net/idna"
)

var fqdnRegex = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z0-9-]{2,63}$`)

/*
Validate's FQDN by converting it to Punycode and checking it against several criteria
*/
func ValidateFQDN(host string) error {
	// Remove trailing . just in case, but DB and form validation should already prevent that
	trimmedHost := strings.TrimSuffix(host, ".")
	if len(trimmedHost) == 0 {
		return fmt.Errorf("%q: invalid length of 0.  original %q", trimmedHost, host)
	}

	// convert to punycode
	punycodeHost, err := idna.Lookup.ToASCII(trimmedHost)
	if err != nil {
		return fmt.Errorf("%q: failed to convert to punycode - %w", trimmedHost, err)
	}

	if len(punycodeHost) < 1 || len(punycodeHost) > 253 {
		return fmt.Errorf("%q -> %q: invalid length %d bytes. Must be 1-253 bytes (chars) after conversion.", trimmedHost, punycodeHost, len(punycodeHost))
	}

	if !fqdnRegex.MatchString(punycodeHost) {
		return fmt.Errorf("%q -> %q: Failed validation.", trimmedHost, punycodeHost)
	}
	return nil
}

func ValidateSubnetMask(mask string) bool {
	addr, err := netip.ParseAddr(mask)
	if err != nil || !addr.Is4() {
		return false
	}

	byteArray := addr.As4()
	maskInt := uint32(byteArray[0])<<24 | uint32(byteArray[1])<<16 | uint32(byteArray[2])<<8 | uint32(byteArray[3])
	inverted := ^maskInt
	return (inverted & (inverted + 1)) == 0
}
