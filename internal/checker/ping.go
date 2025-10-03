package checker

import (
	"time"

	"github.com/go-ping/ping"
)

func CheckPing(host string, timeout time.Duration, count int) bool {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return false
	}

	pinger.Count = count
	pinger.Timeout = timeout
	
	pinger.SetPrivileged(false)

	err = pinger.Run()
	if err != nil {
		return false
	}

	stats := pinger.Statistics()
	return stats.PacketsRecv > 0
}