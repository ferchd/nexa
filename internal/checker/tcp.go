package checker

import (
	"net"
	"time"
)

func CheckTCP(host string, port int, timeout time.Duration) bool {
	address := net.JoinHostPort(host, string(port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}