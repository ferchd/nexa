package checker

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MockTCPServer(t *testing.T) (string, func()) {
	t.Helper()
	
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create mock TCP server: %v", err)
	}
	
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	
	cleanup := func() {
		listener.Close()
	}
	
	return listener.Addr().String(), cleanup
}

func MockHTTPServer(t *testing.T, statusCode int) (*httptest.Server, string) {
	t.Helper()
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))
	
	return server, server.URL
}