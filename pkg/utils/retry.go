package utils

import (
	"time"
)

func Retry(attempts int, sleep time.Duration, fn func() bool) bool {
	for i := 0; i < attempts; i++ {
		if fn() {
			return true
		}
		if i < attempts-1 {
			time.Sleep(sleep)
		}
	}
	return false
}

func RetryWithBackoff(attempts int, initialSleep time.Duration, fn func() bool) bool {
	sleep := initialSleep
	for i := 0; i < attempts; i++ {
		if fn() {
			return true
		}
		if i < attempts-1 {
			time.Sleep(sleep)
			sleep = sleep * 2 // Exponential backoff
		}
	}
	return false
}