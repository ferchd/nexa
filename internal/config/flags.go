package config

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func parseFlags() {
	pflag.StringSlice("external", []string{}, 
		"External host or host:port to probe (can repeat). Example: --external 8.8.8.8:53 --external 1.1.1.1")
	pflag.StringSlice("corp", []string{},
		"Corporate host or host:port to probe (can repeat). Example: --corp fileserver.corp.local:445")
	pflag.String("http-url", "https://www.google.com/generate_204", 
		"HTTP URL for captive-portal detection")
	pflag.String("dns-probe", "", 
		"Internal DNS name for corporate indicator")

	pflag.Duration("tcp-timeout", 2*time.Second, "TCP connect timeout")
	pflag.Duration("http-timeout", 5*time.Second, "HTTP timeout")
	pflag.Duration("ping-timeout", 3*time.Second, "Ping timeout")

	pflag.Int("attempts", 2, "Retry attempts per check")
	pflag.Duration("backoff", 1500*time.Millisecond, "Backoff between retries")
	pflag.Int("workers", 8, "Worker count for concurrent checks")

	pflag.Bool("stdout-json", false, "Print JSON result to stdout")

	pflag.Bool("prometheus", false, "Enable Prometheus exporter")
	pflag.Int("prom-port", 9000, "Prometheus exporter port")

	pflag.String("log-file", "/var/log/nexa.log", "Log file path")
	pflag.String("log-level", "info", "Log level (debug, info, warn, error)")
	pflag.Int("log-max-size-mb", 10, "Maximum log file size in MB")
	pflag.Int("log-max-backups", 3, "Maximum number of old log files to retain")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if externalHosts := viper.GetStringSlice("external"); len(externalHosts) > 0 {
		parsed := parseHostStrings(externalHosts)
		viper.Set("external_hosts", parsed)
	}

	if corpHosts := viper.GetStringSlice("corp"); len(corpHosts) > 0 {
		parsed := parseHostStrings(corpHosts)
		viper.Set("corp_hosts", parsed)
	}
}

func parseHostStrings(hostStrings []string) []map[string]interface{} {
	var hosts []map[string]interface{}
	for _, s := range hostStrings {
		host, port := parseHostPort(s)
		hosts = append(hosts, map[string]interface{}{
			"host": host,
			"port": port,
		})
	}
	return hosts
}

func parseHostPort(s string) (string, int) {
	parts := strings.Split(s, ":")
	if len(parts) == 1 {
		return strings.TrimSpace(parts[0]), 0
	}

	host := strings.TrimSpace(strings.Join(parts[:len(parts)-1], ":"))
	portStr := parts[len(parts)-1]

	var port int
	fmt.Sscanf(portStr, "%d", &port)
	return host, port
}