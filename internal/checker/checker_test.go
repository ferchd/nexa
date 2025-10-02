package checker

import (
	"net"
	"testing"
	"time"

	"github.com/ferchd/nexa/internal/config"
)

func TestCheckTCP_Success(t *testing.T) {
	addr, cleanup := MockTCPServer(t)
	defer cleanup()
	
	host, portStr, _ := net.SplitHostPort(addr)
	var port int
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
		t.Fatalf("Failed to parse port: %v", err)
	}
	
	result := CheckTCP(host, port, 2*time.Second)
	if !result {
		t.Errorf("Expected TCP check to pass for mock server")
	}
}

func TestCheckTCP_Failure(t *testing.T) {
	result := CheckTCP("127.0.0.1", 9999, 1*time.Second)
	if result {
		t.Errorf("Expected TCP check to fail for closed port")
	}
}

func TestCheckTCP_Timeout(t *testing.T) {
	// Test with very short timeout
	result := CheckTCP("1.1.1.1", 81, 1*time.Nanosecond)
	if result {
		t.Errorf("Expected TCP check to timeout")
	}
}

func TestCheckDNS_Success(t *testing.T) {
	result := CheckDNS("localhost")
	if !result {
		t.Errorf("Expected DNS resolution to work for localhost")
	}
}

func TestCheckDNS_Failure(t *testing.T) {
	result := CheckDNS("this-domain-should-never-exist-12345.invalid")
	if result {
		t.Errorf("Expected DNS resolution to fail for invalid domain")
	}
}

func TestCheckHTTP_Success(t *testing.T) {
	server, url := MockHTTPServer(t, 200)
	defer server.Close()
	
	result := CheckHTTP(url, 5*time.Second)
	if !result {
		t.Errorf("Expected HTTP check to pass for 200 status")
	}
}

func TestCheckHTTP_Failure(t *testing.T) {
	server, url := MockHTTPServer(t, 404)
	defer server.Close()
	
	result := CheckHTTP(url, 5*time.Second)
	if result {
		t.Errorf("Expected HTTP check to fail for 404 status")
	}
}

func TestCheckHTTP_Timeout(t *testing.T) {
	result := CheckHTTP("http://10.255.255.1", 100*time.Millisecond)
	if result {
		t.Errorf("Expected HTTP check to timeout")
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
			got := result.ExitCode()
			if got != tc.expected {
				t.Errorf("Expected exit code %d, got %d", tc.expected, got)
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
	
	if checker.config == nil {
		t.Error("Expected config to be set")
	}
}

func TestNexaInitialization_WithPrometheus(t *testing.T) {
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
		Prometheus:  true,
		PromPort:    9090,
	}

	checker, err := NewNexa(cfg)
	if err != nil {
		t.Fatalf("Failed to create Nexa with Prometheus: %v", err)
	}

	if checker.metrics == nil {
		t.Error("Expected metrics to be initialized")
	}
}

func TestCheckResult_JSON(t *testing.T) {
	result := &GlobalResult{
		InternetOK:  true,
		CorporateOK: false,
		Timestamp:   time.Now(),
		InternetDetails: make(map[string]CheckResult),
		CorporateDetails: make(map[string]CheckResult),
	}
	
	// Should not panic
	result.PrintJSON()
}

func TestCheckResult_Human(t *testing.T) {
	result := &GlobalResult{
		InternetOK:  true,
		CorporateOK: true,
		Timestamp:   time.Now(),
		ElapsedSeconds: 1.234,
		InternetDetails: make(map[string]CheckResult),
		CorporateDetails: make(map[string]CheckResult),
		Summary: SummaryStats{
			TotalChecks: 5,
			Successful: 4,
			Failed: 1,
			ExternalChecks: 3,
			CorporateChecks: 2,
		},
	}
	
	// Should not panic
	result.PrintHuman()
}