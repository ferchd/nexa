package checker

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ferchd/nexa/internal/config"
	"github.com/ferchd/nexa/internal/metrics"
	"github.com/ferchd/nexa/pkg/utils"
)

type CheckType string

const (
	CheckTypeExternal CheckType = "external"
	CheckTypeCorporate CheckType = "corporate"
)

type CheckResult struct {
	Type      CheckType               `json:"type"`
	Host      string                  `json:"host"`
	Port      int                     `json:"port,omitempty"`
	Success   bool                    `json:"success"`
	Error     string                  `json:"error,omitempty"`
	Details   map[string]interface{}  `json:"details"`
	Duration  time.Duration           `json:"duration_ms"`
	Timestamp time.Time               `json:"timestamp"`
}

type GlobalResult struct {
	InternetOK       bool                   `json:"internet"`
	CorporateOK      bool                   `json:"corporate"`
	Timestamp        time.Time              `json:"timestamp"`
	ElapsedSeconds   float64                `json:"elapsed_s"`
	InternetDetails  map[string]CheckResult `json:"internet_details"`
	CorporateDetails map[string]CheckResult `json:"corporate_details"`
	Summary          SummaryStats           `json:"summary"`
}

type SummaryStats struct {
	TotalChecks    int `json:"total_checks"`
	Successful     int `json:"successful"`
	Failed         int `json:"failed"`
	ExternalChecks int `json:"external_checks"`
	CorporateChecks int `json:"corporate_checks"`
}

type Nexa struct {
	config  *config.Config
	metrics *metrics.PrometheusMetrics
	logger  *log.Logger
}

func NewNexa(cfg *config.Config) (*Nexa, error) {
	var promMetrics *metrics.PrometheusMetrics
	if cfg.Prometheus {
		var err error
		promMetrics, err = metrics.NewPrometheusMetrics(cfg.PromPort)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize metrics: %v", err)
		}
	}

	return &Nexa{
		config:  cfg,
		metrics: promMetrics,
		logger:  log.New(log.Writer(), "[nexa] ", log.LstdFlags),
	}, nil
}

func (nc *Nexa) Run() *GlobalResult {
	startTime := time.Now()
	result := &GlobalResult{
		Timestamp:        startTime,
		InternetDetails:  make(map[string]CheckResult),
		CorporateDetails: make(map[string]CheckResult),
	}

	var wg sync.WaitGroup
	results := make(chan CheckResult, len(nc.config.ExternalHosts)+len(nc.config.CorpHosts))

	for _, hp := range nc.config.ExternalHosts {
		wg.Add(1)
		go func(hp config.HostPort) {
			defer wg.Done()
			checkResult := nc.checkExternal(hp)
			results <- checkResult
		}(hp)
	}

	for _, hp := range nc.config.CorpHosts {
		wg.Add(1)
		go func(hp config.HostPort) {
			defer wg.Done()
			checkResult := nc.checkCorporate(hp)
			results <- checkResult
		}(hp)
	}

	wg.Wait()
	close(results)

	for checkResult := range results {
		key := fmt.Sprintf("%s:%s:%d", checkResult.Type, checkResult.Host, checkResult.Port)
		if checkResult.Type == CheckTypeExternal {
			result.InternetDetails[key] = checkResult
		} else {
			result.CorporateDetails[key] = checkResult
		}
	}

	result.InternetOK = nc.determineInternetStatus(result)
	result.CorporateOK = nc.determineCorporateStatus(result)
	result.ElapsedSeconds = time.Since(startTime).Seconds()
	result.Summary = nc.calculateSummary(result)

	if nc.metrics != nil {
		nc.metrics.UpdateInternetStatus(result.InternetOK)
		nc.metrics.UpdateCorporateStatus(result.CorporateOK)
		nc.metrics.UpdateCheckSummary(result.Summary)
	}

	return result
}

func (nc *Nexa) checkExternal(hp config.HostPort) CheckResult {
	startTime := time.Now()
	result := CheckResult{
		Type:      CheckTypeExternal,
		Host:      hp.Host,
		Port:      hp.Port,
		Details:   make(map[string]interface{}),
		Timestamp: startTime,
	}

	if hp.Port > 0 {
		tcpOK := utils.Retry(nc.config.Attempts, nc.config.Backoff, func() bool {
			return CheckTCP(hp.Host, hp.Port, nc.config.TCPTimeout)
		})
		result.Details["tcp"] = tcpOK
	}

	pingOK := utils.Retry(nc.config.Attempts, nc.config.Backoff, func() bool {
		return CheckPing(hp.Host, nc.config.PingTimeout, nc.config.Attempts)
	})
	result.Details["ping"] = pingOK

	if nc.config.HTTPURL != "" {
		httpOK := utils.Retry(nc.config.Attempts, nc.config.Backoff, func() bool {
			return CheckHTTP(nc.config.HTTPURL, nc.config.HTTPTimeout)
		})
		result.Details["http"] = httpOK
		result.Details["http_url"] = nc.config.HTTPURL
	}

	result.Success = result.Details["tcp"] == true || 
	                 result.Details["ping"] == true || 
	                 result.Details["http"] == true
	result.Duration = time.Since(startTime)

	return result
}

func (nc *Nexa) checkCorporate(hp config.HostPort) CheckResult {
	startTime := time.Now()
	result := CheckResult{
		Type:      CheckTypeCorporate,
		Host:      hp.Host,
		Port:      hp.Port,
		Details:   make(map[string]interface{}),
		Timestamp: startTime,
	}

	if hp.Port > 0 {
		tcpOK := utils.Retry(nc.config.Attempts, nc.config.Backoff, func() bool {
			return CheckTCP(hp.Host, hp.Port, nc.config.TCPTimeout)
		})
		result.Details["tcp"] = tcpOK
	}

	if nc.config.DNSProbe != "" {
		dnsOK := utils.Retry(nc.config.Attempts, nc.config.Backoff, func() bool {
			return CheckDNS(nc.config.DNSProbe)
		})
		result.Details["dns"] = dnsOK
		result.Details["dns_probe"] = nc.config.DNSProbe
	}

	result.Success = result.Details["tcp"] == true || result.Details["dns"] == true
	result.Duration = time.Since(startTime)

	return result
}

func (nc *Nexa) determineInternetStatus(result *GlobalResult) bool {
	for _, check := range result.InternetDetails {
		if check.Success {
			return true
		}
	}

	if nc.config.HTTPURL != "" {
		httpOK := CheckHTTP(nc.config.HTTPURL, nc.config.HTTPTimeout)
		if httpOK {
			fallbackResult := CheckResult{
				Type:    CheckTypeExternal,
				Host:    nc.config.HTTPURL,
				Success: true,
				Details: map[string]interface{}{
					"http":          true,
					"http_fallback": true,
				},
				Timestamp: time.Now(),
			}
			result.InternetDetails["http_fallback"] = fallbackResult
			return true
		}
	}

	return false
}

func (nc *Nexa) determineCorporateStatus(result *GlobalResult) bool {
	for _, check := range result.CorporateDetails {
		if check.Success {
			return true
		}
	}
	return false
}

func (nc *Nexa) calculateSummary(result *GlobalResult) SummaryStats {
	stats := SummaryStats{}
	
	for _, check := range result.InternetDetails {
		stats.TotalChecks++
		stats.ExternalChecks++
		if check.Success {
			stats.Successful++
		} else {
			stats.Failed++
		}
	}
	
	for _, check := range result.CorporateDetails {
		stats.TotalChecks++
		stats.CorporateChecks++
		if check.Success {
			stats.Successful++
		} else {
			stats.Failed++
		}
	}
	
	return stats
}

func (r *GlobalResult) ExitCode() int {
	if r.InternetOK && r.CorporateOK {
		return 0
	} else if !r.InternetOK && r.CorporateOK {
		return 1
	} else if r.InternetOK && !r.CorporateOK {
		return 2
	} else {
		return 3
	}
}

func (r *GlobalResult) PrintJSON() {
	jsonData, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}
	fmt.Println(string(jsonData))
}

func (r *GlobalResult) PrintHuman() {
	status := "✅"
	if !r.InternetOK || !r.CorporateOK {
		status = "❌"
	}
	
	fmt.Printf("NetCheck Results %s\n", status)
	fmt.Printf("Internet:  %v\n", r.InternetOK)
	fmt.Printf("Corporate: %v\n", r.CorporateOK)
	fmt.Printf("Duration:  %.3fs\n", r.ElapsedSeconds)
	fmt.Printf("Checks:    %d total (%d external, %d corporate)\n", 
		r.Summary.TotalChecks, r.Summary.ExternalChecks, r.Summary.CorporateChecks)
	fmt.Printf("Success:   %d/%d\n", r.Summary.Successful, r.Summary.TotalChecks)
}