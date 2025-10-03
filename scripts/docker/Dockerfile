FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=${VERSION:-dev}" -o nexa ./cmd/nexa

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder --chown=appuser:appgroup /app/nexa .
COPY --chown=appuser:appgroup examples/config.yaml /etc/nexa/config.yaml

RUN mkdir -p /var/log/nexa && \
    chown appuser:appgroup /var/log/nexa

USER appuser

EXPOSE 9000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ./nexa --external 8.8.8.8:53 --stdout-json > /dev/null || exit 1

ENTRYPOINT ["./nexa"]
CMD ["--config", "/etc/nexa/config.yaml"]