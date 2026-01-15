package utils

import (
	"net"
	"strings"
)

// NormalizeIPAddress normalizes an IP address (optionally with prefix length) to match NetBox's canonical form.
//
// This is important for IPv6 addresses where NetBox removes leading zeros from segments.
// For example, "fd00:0f13:e5bf::ced6/128" becomes "fd00:f13:e5bf::ced6/128".
//
// Normalization is best-effort: if parsing fails, the input is returned unchanged.
func NormalizeIPAddress(address string) string {
	address = strings.TrimSpace(address)
	if address == "" {
		return address
	}

	ipStr, prefix, hasPrefix := strings.Cut(address, "/")

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return address
	}

	// Only normalize IPv6. IPv4 is returned as-is to avoid any unexpected formatting changes.
	if ip.To4() != nil {
		return address
	}

	if hasPrefix {
		return ip.String() + "/" + prefix
	}

	return ip.String()
}
