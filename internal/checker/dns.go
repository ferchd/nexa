package checker

import "net"

func CheckDNS(hostname string) bool {
	_, err := net.LookupHost(hostname)
	return err == nil
}