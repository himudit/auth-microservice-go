package utils

import (
	"net"
	"net/http"
	"strings"
)

// GetIP extracts the client's IP address from an *http.Request.
// Priority:
// 1. X-Forwarded-For (first non-empty, comma-separated entry)
// 2. X-Real-IP
// 3. r.RemoteAddr (host part)
// Returns empty string if no valid IP found.
func GetIP(r *http.Request) string {
	// 1) X-Forwarded-For (may contain comma-separated list)
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		// split and use the first non-empty element
		parts := strings.Split(xff, ",")
		for _, p := range parts {
			ip := strings.TrimSpace(p)
			if ip != "" {
				if parsed := parseIP(ip); parsed != "" {
					return parsed
				}
			}
		}
	}

	// 2) X-Real-IP
	if xr := strings.TrimSpace(r.Header.Get("X-Real-Ip")); xr != "" {
		if parsed := parseIP(xr); parsed != "" {
			return parsed
		}
	}

	// 3) RemoteAddr (may be "ip:port" or "[ip]:port")
	if ra := strings.TrimSpace(r.RemoteAddr); ra != "" {
		// try SplitHostPort first
		if host, _, err := net.SplitHostPort(ra); err == nil {
			if parsed := parseIP(host); parsed != "" {
				return parsed
			}
		} else {
			// fallback: maybe RemoteAddr is just an IP
			if parsed := parseIP(ra); parsed != "" {
				return parsed
			}
		}
	}

	return ""
}

// parseIP validates and returns the canonical IP string or empty if invalid.
func parseIP(s string) string {
	// strip possible surrounding brackets for IPv6 like "[::1]"
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	if ip := net.ParseIP(s); ip != nil {
		// return the string form without zone (if any)
		if i := ip.String(); i != "" {
			return i
		}
	}
	return ""
}
