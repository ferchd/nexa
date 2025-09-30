#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

VERSION="1.0.0"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/nexa"
SERVICE_DIR="/etc/systemd/system"
LOG_DIR="/var/log/nexa"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="armv7"
        ;;
    i386|i686)
        ARCH="386"
        ;;
esac

BINARY_NAME="nexa"
DOWNLOAD_URL="https://github.com/ferchd/nexa/releases/download/v${VERSION}/nexa-${OS}-${ARCH}"

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_root() {
    if [[ $EUID -eq 0 ]]; then
        print_status "Running as root"
    else
        print_error "This script requires root privileges. Please run with sudo."
        exit 1
    fi
}

download_binary() {
    print_status "Downloading nexa binary..."
    
    if command -v curl &> /dev/null; then
        curl -L -o "/tmp/${BINARY_NAME}" "${DOWNLOAD_URL}"
    elif command -v wget &> /dev/null; then
        wget -O "/tmp/${BINARY_NAME}" "${DOWNLOAD_URL}"
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
    
    if [[ ! -f "/tmp/${BINARY_NAME}" ]]; then
        print_error "Failed to download binary"
        exit 1
    fi
    
    chmod +x "/tmp/${BINARY_NAME}"
}

install_binary() {
    print_status "Installing binary to ${INSTALL_DIR}..."
    
    mkdir -p "${INSTALL_DIR}"
    
    cp "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    
    if command -v "${BINARY_NAME}" &> /dev/null; then
        print_status "Binary installed successfully"
    else
        print_error "Binary installation failed"
        exit 1
    fi
}

create_directories() {
    print_status "Creating directories..."
    
    mkdir -p "${CONFIG_DIR}"
    mkdir -p "${LOG_DIR}"
    
    chmod 755 "${CONFIG_DIR}"
    chmod 755 "${LOG_DIR}"
}

create_config() {
    print_status "Creating configuration file..."
    
    if [[ ! -f "${CONFIG_DIR}/config.yaml" ]]; then
        cat > "${CONFIG_DIR}/config.yaml" << 'EOF'
# Nexa Configuration
external_hosts:
  - host: 8.8.8.8
    port: 53
  - host: 1.1.1.1
    port: 53

corp_hosts:
  - host: fileserver.corp.local
    port: 445

http_url: "https://www.google.com/generate_204"
dns_probe: "internal.corp.local"

tcp_timeout: 2s
http_timeout: 5s
ping_timeout: 3s

attempts: 2
backoff: 1500ms

workers: 8
stdout_json: false

prometheus: true
prom_port: 9000

log_file: "/var/log/nexa/nexa.log"
log_level: "info"
log_max_size_mb: 10
log_max_backups: 3
EOF
        print_status "Configuration file created at ${CONFIG_DIR}/config.yaml"
    else
        print_warning "Configuration file already exists at ${CONFIG_DIR}/config.yaml"
    fi
}

setup_systemd_service() {
    if [[ ! -d "/run/systemd/system" ]]; then
        print_warning "Systemd not detected, skipping service installation"
        return 0
    fi
    
    print_status "Setting up systemd service..."
    
    cat > "${SERVICE_DIR}/nexa.service" << EOF
[Unit]
Description=Nexa Network Connectivity Monitor
Documentation=https://github.com/ferchd/nexa
After=network.target

[Service]
Type=simple
User=root
ExecStart=${INSTALL_DIR}/nexa --config ${CONFIG_DIR}/config.yaml
Restart=always
RestartSec=30
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable nexa.service
    
    print_status "Systemd service installed and enabled"
    
    read -p "Do you want to start the nexa service now? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        systemctl start nexa.service
        print_status "Nexa service started"
    fi
}

setup_log_rotation() {
    print_status "Setting up log rotation..."
    
    if command -v logrotate &> /dev/null; then
        cat > "/etc/logrotate.d/nexa" << EOF
/var/log/nexa/*.log {
    daily
    missingok
    rotate 7
    compress
    delaycompress
    notifempty
    copytruncate
}
EOF
        print_status "Log rotation configured"
    else
        print_warning "logrotate not found, skipping log rotation setup"
    fi
}

cleanup() {
    rm -f "/tmp/${BINARY_NAME}"
    print_status "Installation cleanup completed"
}

main() {
    print_status "Starting Nexa installation..."
    
    check_root
    download_binary
    install_binary
    create_directories
    create_config
    setup_systemd_service
    setup_log_rotation
    cleanup
    
    print_status "Nexa installation completed successfully!"
    echo
    print_status "Quick start:"
    echo "  nexa --help"
    echo "  systemctl status nexa"
    echo
    print_status "Configuration: ${CONFIG_DIR}/config.yaml"
    print_status "Logs: ${LOG_DIR}/"
}

main "$@"