```markdown
# Nexa Architecture

## Overview

Nexa is designed as a modular, concurrent network connectivity checker with clear separation of concerns.

## Project Structure

```

```
nexa/
├── cmd/
│   └── nexa/
│       └── main.go              # Application entry point
├── internal/
│   ├── checker/                 # Core checking logic
│   │   ├── checker.go          # Main orchestrator
│   │   ├── tcp.go              # TCP connectivity checks
│   │   ├── http.go             # HTTP connectivity checks
│   │   ├── dns.go              # DNS resolution checks
│   │   ├── ping.go             # ICMP ping checks
│   │   └── *_test.go           # Unit tests
│   ├── config/                  # Configuration management
│   │   ├── config.go           # Config structures
│   │   └── flags.go            # CLI flag parsing
│   └── metrics/                 # Metrics exporters
│       └── prometheus.go       # Prometheus exporter
├── pkg/
│   └── utils/                   # Utility functions
│       └── retry.go            # Retry logic
├── scripts/                     # Build and deployment scripts
├── examples/                    # Example configurations
└── docs/                        # Documentation
```

## Component Diagram

```
┌─────────────────────────────────────────────────────────┐
│                        main.go                          │
│                   (Entry Point)                         │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                    config.Load()                        │
│              (Configuration Loading)                    │
│  • CLI Flags  • Env Vars  • YAML Files  • Defaults    │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                  checker.NewNexa()                      │
│                (Checker Initialization)                 │
│  • Initialize logger                                    │
│  • Start Prometheus metrics (if enabled)               │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                    nexa.Run()                           │
│              (Main Execution Loop)                      │
└────────────────────┬────────────────────────────────────┘
                     │
         ┌───────────┴───────────┐
         ▼                       ▼
┌──────────────────┐   ┌──────────────────┐
│  checkExternal() │   │  checkCorporate()│
│   (Goroutines)   │   │   (Goroutines)  │
└────────┬─────────┘   └────────┬─────────┘
         │                      │
    ┌────┴────┐            ┌────┴────┐
    ▼         ▼            ▼         ▼
┌─────┐   ┌─────┐      ┌─────┐   ┌─────┐
│ TCP │   │HTTP │      │ TCP │   │ DNS │
└─────┘   └─────┘      └─────┘   └─────┘
    │         │            │         │
    └────┬────┘            └────┬────┘
         ▼                      ▼
┌────────────────────────────────────────┐
│       Aggregate Results                │
│  • Internet Status                     │
│  • Corporate Status                    │
│  • Detailed Check Results              │
│  • Summary Statistics                  │
└────────┬───────────────────────────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌─────────┐ ┌──────────────┐
│ Metrics │ │ Output       │
│ Update  │ │ (JSON/Human) │
└─────────┘ └──────────────┘
```

## Data Flow

### 1. Configuration Phase
```
CLI Args → Environment → YAML → Defaults → Merged Config
```

### 2. Check Execution Phase
```
Config → Worker Pool → Concurrent Checks → Results Channel → Aggregation
```

### 3. Output Phase
```
Aggregated Results → Metrics Update → Format (JSON/Human) → stdout
```

## Concurrency Model

Nexa uses a worker pool pattern for concurrent checks:

```go
// Simplified concurrency flow
results := make(chan CheckResult, totalChecks)
var wg sync.WaitGroup

for _, host := range config.ExternalHosts {
    wg.Add(1)
    go func(host HostPort) {
        defer wg.Done()
        result := performCheck(host)
        results <- result
    }(host)
}

wg.Wait()
close(results)

// Aggregate results
for result := range results {
    aggregateResult(result)
}
```

## Key Design Decisions

### 1. **Separation of Check Types**
- Each protocol (TCP, HTTP, DNS, ICMP) has its own file
- Easy to add new check types
- Testable in isolation

### 2. **Concurrent Execution**
- Uses goroutines for parallel checks
- Configurable worker count
- Non-blocking channel communication

### 3. **Retry Logic**
- Exponential backoff
- Configurable attempts
- Isolated in utils package

### 4. **Configuration Flexibility**
- Multiple sources (CLI, env, file)
- Clear precedence order
- Type-safe with viper

### 5. **Metrics Export**
- Optional Prometheus integration
- Minimal overhead when disabled
- Standard metric types

## Error Handling

```
Check Failure → Retry (with backoff) → Final Result
                                             │
                                             ├→ Success: true
                                             └→ Success: false + error details
```

## Testing Strategy

1. **Unit Tests**: Individual check functions with mocks
2. **Integration Tests**: Full check cycles with test servers
3. **Table-Driven Tests**: Multiple scenarios per test
4. **Mock Servers**: Local TCP/HTTP servers for testing

## Performance Considerations

- **Memory**: ~50MB typical usage
- **CPU**: Minimal, spikes during checks
- **Network**: Dependent on check frequency
- **Goroutines**: Max = workers + background tasks

## Security Considerations

1. **Privileged Operations**: ICMP requires root or capabilities
2. **Timeouts**: All network operations have timeouts
3. **Input Validation**: Config values are validated
4. **Secrets**: No sensitive data in logs or metrics

## Extensibility Points

### Adding New Check Types

```go
// 1. Create new file: internal/checker/newcheck.go
func CheckNew(params) bool {
    // Implementation
}

// 2. Update checker.go to include new check
result.Details["new_check"] = CheckNew(...)

// 3. Add tests
func TestCheckNew(t *testing.T) { ... }
```

### Adding New Metrics

```go
// 1. Add metric in prometheus.go
newMetric: prometheus.NewGauge(...)

// 2. Register metric
prometheus.MustRegister(metrics.newMetric)

// 3. Update method
func (m *PrometheusMetrics) UpdateNewMetric(value float64) {
    m.newMetric.Set(value)
}
```

## Future Architecture Improvements

1. **Plugin System**: Load check types dynamically
2. **gRPC API**: Remote check orchestration
3. **Distributed Mode**: Multiple nexa instances
4. **State Persistence**: Historical data storage
5. **Advanced Scheduling**: Cron-like check scheduling