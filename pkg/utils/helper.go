package utils

import (
	"fmt"
	"strings"
	"time"
)

func ParseHostPort(s string) (string, int, error) {
	if !strings.Contains(s, ":") {
		return strings.TrimSpace(s), 0, nil
	}

	parts := strings.Split(s, ":")
	if len(parts) < 2 {
		return "", 0, fmt.Errorf("invalid host:port format: %s", s)
	}

	host := strings.TrimSpace(strings.Join(parts[:len(parts)-1], ":"))
	portStr := parts[len(parts)-1]

	var port int
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port in '%s': %v", s, err)
	}

	return host, port, nil
}

func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return d.String()
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}