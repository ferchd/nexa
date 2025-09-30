package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ExternalHosts []HostPort `mapstructure:"external_hosts"`
	CorpHosts     []HostPort `mapstructure:"corp_hosts"`
	HTTPURL       string     `mapstructure:"http_url"`
	DNSProbe      string     `mapstructure:"dns_probe"`
	
	TCPTimeout  time.Duration `mapstructure:"tcp_timeout"`
	HTTPTimeout time.Duration `mapstructure:"http_timeout"`
	PingTimeout time.Duration `mapstructure:"ping_timeout"`
	
	Attempts int           `mapstructure:"attempts"`
	Backoff  time.Duration `mapstructure:"backoff"`
	
	Workers    int  `mapstructure:"workers"`
	StdoutJSON bool `mapstructure:"stdout_json"`
	
	Prometheus bool `mapstructure:"prometheus"`
	PromPort   int  `mapstructure:"prom_port"`
	
	LogFile      string `mapstructure:"log_file"`
	LogLevel     string `mapstructure:"log_level"`
	LogMaxSizeMB int    `mapstructure:"log_max_size_mb"`
	LogMaxBackups int   `mapstructure:"log_max_backups"`
}

type HostPort struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func Load() (*Config, error) {
	setDefaults()

	viper.SetConfigName("nexa")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/nexa/")
	viper.AddConfigPath("$HOME/.nexa")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}

	viper.SetEnvPrefix("NEXA")
	viper.AutomaticEnv()

	parseFlags()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %v", err)
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("external_hosts", []map[string]interface{}{
		{"host": "8.8.8.8", "port": 53},
		{"host": "1.1.1.1", "port": 53},
	})
	viper.SetDefault("corp_hosts", []map[string]interface{}{})
	viper.SetDefault("http_url", "https://www.google.com/generate_204")
	viper.SetDefault("tcp_timeout", 2*time.Second)
	viper.SetDefault("http_timeout", 5*time.Second)
	viper.SetDefault("ping_timeout", 3*time.Second)
	viper.SetDefault("attempts", 2)
	viper.SetDefault("backoff", 1500*time.Millisecond)
	viper.SetDefault("workers", 8)
	viper.SetDefault("prometheus", false)
	viper.SetDefault("prom_port", 9000)
	viper.SetDefault("log_file", "/var/log/nexa.log")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("log_max_size_mb", 10)
	viper.SetDefault("log_max_backups", 3)
}