# 🔍 Nexa - Enterprise Network Connectivity Monitor

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Release](https://img.shields.io/github/v/release/ferchd/nexa)](https://github.com/ferchd/nexa/releases)
[![Build Status](https://github.com/ferchd/nexa/workflows/CI/badge.svg)](https://github.com/ferchd/nexa/actions)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/ferchd/nexa.svg)](https://hub.docker.com/r/ferchd/nexa)
[![Go Report Card](https://goreportcard.com/badge/github.com/ferchd/nexa)](https://goreportcard.com/report/github.com/ferchd/nexa)
[![Go Reference](https://pkg.go.dev/badge/github.com/ferchd/nexa.svg)](https://pkg.go.dev/github.com/ferchd/nexa)
[![GitHub All Releases](https://img.shields.io/github/downloads/ferchd/nexa/total.svg)](https://github.com/ferchd/nexa/releases)
[![Docker Image Size](https://img.shields.io/docker/image-size/tuusuario/netcheck/latest)](https://hub.docker.com/r/tuusuario/netcheck)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/m/ferchd/nexa)](https://github.com/ferchd/nexa/pulse)

**Nexa** es una herramienta CLI empresarial escrita en Go para monitoreo de conectividad de red. Verifica simultáneamente conectividad a Internet y redes corporativas con soporte multi-protocolo.

**Características principales:** 
- ✅ **Verificaciones multi-protocolo**: TCP, HTTP, DNS, ICMP Ping real
- ✅ **Ejecución concurrente**: Hasta 8 verificaciones simultáneas
- ✅ **Métricas Prometheus**: Exportador integrado para monitoring
- ✅ **Configuración flexible**: Archivos YAML, variables entorno, CLI
- ✅ **Enterprise-ready**: Systemd, Docker, Logging rotativo
- ✅ **Cross-platform**: Linux, Windows, macOS, contenedores

## 🚀 Quick Start

```bash
curl -sSL https://raw.githubusercontent.com/ferchd/nexa/main/scripts/install.sh | sudo bash

nexa --external 8.8.8.8:53 --corp fileserver.corp.local:445 --stdout-json
```