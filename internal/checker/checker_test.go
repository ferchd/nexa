package checker

import (
	"testing"
	"time"

	"github.com/ferchd/nexa/internal/config"
)

func TestCheckTCP(t *testing.T) {
	result := CheckTCP("google.com", 80, 5*time.Second)
	if !result {
		t.Errorf("Expected TCP check to pass for google.com:80")
	}

	result = CheckTCP("invalid-host-that-should-not-exist.local", 9999, 1*time.Second)
	if result {
		t.Errorf("Expected TCP check to fail for invalid host")
	}
}

func TestCheckDNS(t *testing.T) {
	result := CheckDNS("google.com")
	if !result {
		t.Errorf("Expected DNS resolution to work for google.com")
	}

	result = CheckDNS("invalid-domain-that-should-not-exist-12345.local")
	if result {
		t.Errorf("Expected DNS resolution to fail for invalid domain")
	}
}

func TestCheckHTTP(t *testing.T) {
	result := CheckHTTP("https://httpbin.org/status/200", 5*time.Second)
	if !result {
		t.Errorf("Expected HTTP check to pass for valid URL")
	}

	result = CheckHTTP("https://httpbin.org/status/404", 5*time.Second)
	if result {
		t.Errorf("Expected HTTP check to fail for 404")
	}
}

func TestGlobalResultExitCodes(t *testing.T) {
	testCases := []struct {
		name      string
		internet  bool
		corporate bool
		expected  int
	}{
		{"Both OK", true, true, 0},
		{"Internet fail", false, true, 1},
		{"Corporate fail", true, false, 2},
		{"Both fail", false, false, 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := &GlobalResult{
				InternetOK:  tc.internet,
				CorporateOK: tc.corporate,
			}
			if result.ExitCode() != tc.expected {
				t.Errorf("Expected exit code %d, got %d", tc.expected, result.ExitCode())
			}
		})
	}
}

func TestNexaInitialization(t *testing.T) {
	cfg := &config.Config{
		ExternalHosts: []config.HostPort{
			{Host: "8.8.8.8", Port: 53},
		},
		TCPTimeout:  2 * time.Second,
		HTTPTimeout: 5 * time.Second,
		PingTimeout: 3 * time.Second,
		Attempts:    2,
		Backoff:     1 * time.Second,
		Workers:     4,
	}

	checker, err := NewNexa(cfg)
	if err != nil {
		t.Fatalf("Failed to create Nexa: %v", err)
	}

	if checker == nil {
		t.Error("Expected Nexa to be initialized")
	}
}