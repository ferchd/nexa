# 🔍 Nexa - Enterprise Network Connectivity Monitor

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Release](https://img.shields.io/github/v/release/ferchd/nexa)](https://github.com/ferchd/nexa/releases)
[![Build Status](https://github.com/ferchd/nexa/workflows/Build%20and%20Test/badge.svg)](https://github.com/ferchd/nexa/actions)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ferchd/nexa)](https://goreportcard.com/report/github.com/ferchd/nexa)
[![Coverage](https://img.shields.io/codecov/c/github/ferchd/nexa)](https://codecov.io/gh/ferchd/nexa)

**Nexa** is an enterprise-grade CLI tool written in Go for comprehensive network connectivity monitoring. It simultaneously verifies Internet and corporate network connectivity with multi-protocol support.

---

## ✨ Features

- ✅ **Multi-Protocol Checks**: TCP, HTTP, DNS, ICMP Ping
- ✅ **Concurrent Execution**: Up to 8 simultaneous checks
- ✅ **Prometheus Metrics**: Built-in metrics exporter for monitoring
- ✅ **Flexible Configuration**: YAML files, environment variables, CLI flags
- ✅ **Enterprise-Ready**: Systemd service, Docker support, log rotation
- ✅ **Cross-Platform**: Linux, Windows, macOS, containers
- ✅ **Smart Retry Logic**: Configurable attempts with exponential backoff
- ✅ **Structured Output**: JSON and human-readable formats

---

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Usage](#-usage)
- [Configuration](#-configuration)
- [Exit Codes](#-exit-codes)
- [Prometheus Metrics](#-prometheus-metrics)
- [Examples](#-examples)
- [Deployment](#-deployment)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)

---

## 🚀 Quick Start

### One-line Installation

```bash
curl -sSL https://raw.githubusercontent.com/ferchd/nexa/main/scripts/install.sh | sudo bash
```

### Basic Usage

```bash
# Check internet connectivity
nexa --external 8.8.8.8:53 --stdout-json

# Check corporate network
nexa --corp fileserver.corp.local:445 --corp dc01.corp.local:389

# Use configuration file
nexa --config /etc/nexa/config.yaml

# Enable Prometheus metrics
nexa --prometheus --prom-port 9000
```

---

## 📦 Installation

### Method 1: Binary Release (Recommended)

```bash
# Linux
wget https://github.com/ferchd/nexa/releases/latest/download/nexa-linux-amd64
chmod +x nexa-linux-amd64
sudo mv nexa-linux-amd64 /usr/local/bin/nexa

# macOS
brew install ferchd/tap/nexa

# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/ferchd/nexa/releases/latest/download/nexa-windows-amd64.exe -OutFile nexa.exe
```

### Method 2: From Source

```bash
git clone https://github.com/ferchd/nexa.git
cd nexa
make build
sudo make install
```

### Method 3: Docker

```bash
docker pull ghcr.io/ferchd/nexa:latest
docker run --rm ghcr.io/ferchd/nexa:latest --help
```

### Method 4: Go Install

```bash
go install github.com/ferchd/nexa/cmd/nexa@latest
```

---

## 💻 Usage

### Command-Line Flags

```bash
nexa [flags]

Flags:
  --external strings         External host:port to probe (repeatable)
  --corp strings            Corporate host:port to probe (repeatable)
  --http-url string         HTTP URL for connectivity check (default "https://www.google.com/generate_204")
  --dns-probe string        Internal DNS name for corporate detection
  
  --tcp-timeout duration    TCP connection timeout (default 2s)
  --http-timeout duration   HTTP request timeout (default 5s)
  --ping-timeout duration   ICMP ping timeout (default 3s)
  
  --attempts int           Retry attempts per check (default 2)
  --backoff duration       Backoff between retries (default 1.5s)
  --workers int            Concurrent worker count (default 8)
  
  --stdout-json            Output results as JSON
  --config string          Configuration file path
  
  --prometheus             Enable Prometheus metrics exporter
  --prom-port int          Prometheus port (default 9000)
  
  --log-file string        Log file path (default "/var/log/nexa.log")
  --log-level string       Log level: debug|info|warn|error (default "info")
  
  -v, --version            Show version information
  -h, --help               Show help
```

### Configuration Precedence

1. CLI Flags (highest priority)
2. Environment Variables (prefix: NEXA_)
3. Config File (YAML)
4. Defaults (lowest priority)

---

## ⚙️ Configuration

### Configuration File (config.yaml)

```yaml
# External connectivity checks
external_hosts:
  - host: "8.8.8.8"
    port: 53
  - host: "1.1.1.1"
    port: 53
  - host: "cloudflare.com"
    port: 443

# Corporate network checks
corp_hosts:
  - host: "fileserver.corp.local"
    port: 445
  - host: "dc01.corp.local"
    port: 389
  - host: "exchange.corp.local"
    port: 443

# HTTP connectivity test
http_url: "https://www.google.com/generate_204"

# DNS probe for corporate network
dns_probe: "internal.corp.local"

# Timeouts
tcp_timeout: "2s"
http_timeout: "5s"
ping_timeout: "3s"

# Retry settings
attempts: 2
backoff: "1500ms"

# Performance
workers: 8

# Output
stdout_json: false

# Prometheus metrics
prometheus: true
prom_port: 9000

# Logging
log_file: "/var/log/nexa/nexa.log"
log_level: "info"
log_max_size_mb: 10
log_max_backups: 3
```

### Environment Variables

```bash
export NEXA_EXTERNAL='8.8.8.8:53,1.1.1.1:53'
export NEXA_CORP='fileserver.corp.local:445'
export NEXA_HTTP_URL='https://www.google.com/generate_204'
export NEXA_TCP_TIMEOUT='2s'
export NEXA_ATTEMPTS=3
export NEXA_PROMETHEUS=true
export NEXA_PROM_PORT=9000
```

---

## 🚦 Exit Codes

Nexa uses specific exit codes to indicate connectivity status:

| Exit Code | Meaning | Internet | Corporate |
|-----------|---------|----------|-----------|
| **0** | ✅ Both OK | ✅ Up | ✅ Up |
| **1** | ⚠️ Internet Down | ❌ Down | ✅ Up |
| **2** | ⚠️ Corporate Down | ✅ Up | ❌ Down |
| **3** | ❌ Both Down | ❌ Down | ❌ Down |

### Usage in Scripts

```bash
#!/bin/bash
nexa --config /etc/nexa/config.yaml

case $? in
  0)
    echo "All systems operational"
    ;;
  1)
    echo "Internet connectivity lost"
    send_alert "Internet down"
    ;;
  2)
    echo "Corporate network unreachable"
    send_alert "VPN/Corporate network down"
    ;;
  3)
    echo "Complete network failure"
    send_alert "CRITICAL: All networks down"
    ;;
esac
```

---

## 📊 Prometheus Metrics

When `--prometheus` is enabled, Nexa exposes metrics on `/metrics`:

### Available Metrics

```
# Internet connectivity status (1=up, 0=down)
nexa_internet_up

# Corporate network status (1=up, 0=down)
nexa_corporate_up

# Last check duration in seconds
nexa_check_duration_seconds

# Total checks performed by type
nexa_checks_total{type="external|corporate"}

# Successful checks by type
nexa_checks_success_total{type="external|corporate"}

# Failed checks by type
nexa_checks_failed_total{type="external|corporate"}
```

### Prometheus Configuration

```yaml
scrape_configs:
  - job_name: 'nexa'
    static_configs:
      - targets: ['localhost:9000']
    scrape_interval: 30s
```

### Grafana Dashboard

Import dashboard ID: `TBD` (coming soon)

---

## 📚 Examples

### Example 1: Basic Internet Check

```bash
nexa --external 8.8.8.8:53 --external 1.1.1.1:53
```

**Output:**
```
NetCheck Results ✅
Internet:  true
Corporate: false
Duration:  0.245s
Checks:    2 total (2 external, 0 corporate)
Success:   2/2
```

### Example 2: Corporate Network Detection

```bash
nexa \
  --external 8.8.8.8:53 \
  --corp fileserver.corp.local:445 \
  --corp dc01.corp.local:389 \
  --dns-probe internal.corp.local
```

### Example 3: JSON Output for Automation

```bash
nexa --external 8.8.8.8:53 --stdout-json | jq .
```

**Output:**
```json
{
  "internet": true,
  "corporate": false,
  "timestamp": "2025-10-02T10:30:45Z",
  "elapsed_s": 0.234,
  "internet_details": {
    "external:8.8.8.8:53": {
      "type": "external",
      "host": "8.8.8.8",
      "port": 53,
      "success": true,
      "details": {
        "tcp": true,
        "ping": true
      },
      "duration_ms": 45
    }
  },
  "corporate_details": {},
  "summary": {
    "total_checks": 1,
    "successful": 1,
    "failed": 0,
    "external_checks": 1,
    "corporate_checks": 0
  }
}
```

### Example 4: Monitoring Loop

```bash
#!/bin/bash
while true; do
  nexa --config /etc/nexa/config.yaml
  STATUS=$?
  
  if [ $STATUS -ne 0 ]; then
    logger -t nexa "Connectivity issue detected: exit code $STATUS"
  fi
  
  sleep 60
done
```

### Example 5: Docker Compose

```yaml
version: '3.8'
services:
  nexa:
    image: ghcr.io/ferchd/nexa:latest
    container_name: nexa
    restart: always
    ports:
      - "9000:9000"
    volumes:
      - ./config.yaml:/etc/nexa/config.yaml:ro
    command: ["--config", "/etc/nexa/config.yaml"]
    healthcheck:
      test: ["CMD", "/root/nexa", "--external", "8.8.8.8:53"]
      interval: 30s
      timeout: 5s
      retries: 3
```

---

## 🚀 Deployment

### Systemd Service

```bash
# Install service
sudo cp examples/systemd/nexa.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable nexa
sudo systemctl start nexa

# Check status
sudo systemctl status nexa

# View logs
sudo journalctl -u nexa -f
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nexa
  labels:
    app: nexa
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nexa
  template:
    metadata:
      labels:
        app: nexa
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9000"
    spec:
      containers:
      - name: nexa
        image: ghcr.io/ferchd/nexa:latest
        ports:
        - containerPort: 9000
          name: metrics
        volumeMounts:
        - name: config
          mountPath: /etc/nexa
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "100m"
        livenessProbe:
          exec:
            command:
            - /root/nexa
            - --external
            - 8.8.8.8:53
          initialDelaySeconds: 30
          periodSeconds: 60
      volumes:
      - name: config
        configMap:
          name: nexa-config
---
apiVersion: v1
kind: Service
metadata:
  name: nexa
  labels:
    app: nexa
spec:
  ports:
  - port: 9000
    name: metrics
  selector:
    app: nexa
```

### Ansible Playbook

```yaml
---
- name: Deploy Nexa
  hosts: all
  become: yes
  tasks:
    - name: Download Nexa binary
      get_url:
        url: https://github.com/ferchd/nexa/releases/latest/download/nexa-linux-amd64
        dest: /usr/local/bin/nexa
        mode: '0755'
    
    - name: Create config directory
      file:
        path: /etc/nexa
        state: directory
        mode: '0755'
    
    - name: Copy configuration
      template:
        src: config.yaml.j2
        dest: /etc/nexa/config.yaml
        mode: '0644'
    
    - name: Install systemd service
      copy:
        src: nexa.service
        dest: /etc/systemd/system/nexa.service
        mode: '0644'
    
    - name: Enable and start service
      systemd:
        name: nexa
        enabled: yes
        state: started
        daemon_reload: yes
```

---

## 🔧 Troubleshooting

### Common Issues

#### 1. Permission Denied for ICMP Ping

**Problem:**
```
Failed to create ICMP socket: operation not permitted
```

**Solution:**
```bash
# Option 1: Run as root
sudo nexa --config /etc/nexa/config.yaml

# Option 2: Set capabilities (Linux only)
sudo setcap cap_net_raw+ep /usr/local/bin/nexa

# Option 3: Disable ICMP in config
# Remove ping checks or set enable_icmp: false
```

#### 2. DNS Resolution Fails

**Problem:**
```
DNS probe failed for internal.corp.local
```

**Solution:**
```bash
# Check DNS configuration
cat /etc/resolv.conf

# Test DNS manually
nslookup internal.corp.local

# Verify corporate DNS is accessible
dig @10.0.0.1 internal.corp.local
```

#### 3. Port Already in Use

**Problem:**
```
Failed to start Prometheus server: address already in use
```

**Solution:**
```bash
# Check what's using the port
sudo lsof -i :9000

# Use different port
nexa --prometheus --prom-port 9001
```

#### 4. Timeout Errors

**Problem:**
```
All checks timing out
```

**Solution:**
```yaml
# Increase timeouts in config
tcp_timeout: "5s"
http_timeout: "10s"
ping_timeout: "5s"
attempts: 3
backoff: "2s"
```

### Debug Mode

```bash
# Enable debug logging
nexa --log-level debug --config /etc/nexa/config.yaml

# View detailed logs
tail -f /var/log/nexa/nexa.log
```

### Health Checks

```bash
# Test individual components
nexa --external 8.8.8.8:53 --stdout-json | jq .internet_details

# Check Prometheus metrics
curl http://localhost:9000/metrics

# Verify systemd service
systemctl status nexa
journalctl -u nexa --since "1 hour ago"
```

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone repository
git clone https://github.com/ferchd/nexa.git
cd nexa

# Install dependencies
go mod download

# Run tests
make test

# Run linter
make lint

# Run security scan
make security

# Build locally
make build
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
make coverage

# Run specific test
go test -v -run TestCheckTCP ./internal/checker/

# Run tests with race detection
go test -race ./...
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [go-ping](https://github.com/go-ping/ping) for ICMP
- Uses [Prometheus](https://prometheus.io/) for metrics
- Uses [Viper](https://github.com/spf13/viper) for configuration

---

## 📞 Support

- 📧 Email: ferchd@outlook.com
- 🐛 Issues: [GitHub Issues](https://github.com/ferchd/nexa/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/ferchd/nexa/discussions)
- 📖 Documentation: [Wiki](https://github.com/ferchd/nexa/wiki)

---

## 🗺️ Roadmap

- [ ] Web UI dashboard
- [ ] Slack/Teams notifications
- [ ] InfluxDB support
- [ ] Advanced alerting rules
- [ ] Multi-region checks
- [ ] Historical data storage
- [ ] REST API endpoint
- [ ] Custom check plugins

---

**Made with ❤️ by [Fernando Duarte](https://github.com/ferchd)**
```

---