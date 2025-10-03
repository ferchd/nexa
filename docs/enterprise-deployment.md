# Enterprise Deployment Guide

## High Availability Setup
```yaml
# Kubernetes Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nexa
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: nexa
        image: ghcr.io/ferchd/nexa:latest
        ports:
        - containerPort: 9000