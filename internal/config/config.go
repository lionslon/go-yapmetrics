package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/lionslon/go-yapmetrics/internal/storage"
	"go.uber.org/zap"
	"os"
	"time"
)

// ClientConfig конфиг агента
type ClientConfig struct {
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
	Addr           string `env:"ADDRESS"`
	SignPass       string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY"`
	ConfigJson     string `env:"CONFIG"`
}

// ServerConfig конфиг сервера
type ServerConfig struct {
	Addr            string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FilePath        string `env:"FILE_STORAGE_PATH"`
	CryptoKey       string `env:"CRYPTO_KEY"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	SignPass        string `env:"KEY"`
	EnableProfiling bool   `env:"ENABLE_PROFILING"`
	ConfigJson      string `env:"CONFIG"`
}

// NewClient парсит флаги и env + инициализирует конфиг агента
func NewClient() *ClientConfig {
	cfg := &ClientConfig{}
	parseClientFlags(cfg)
	err := env.Parse(cfg)

	if err != nil {
		zap.S().Error(err)
	}
	return cfg
}

func parseClientFlags(c *ClientConfig) {
	flag.StringVar(&c.Addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.CryptoKey, "c", "", "public crypto-key path")
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&c.RateLimit, "l", 10, "rate limit")
	flag.IntVar(&c.PollInterval, "p", 2, "poll interval in seconds")
	flag.StringVar(&c.SignPass, "k", "", "signature for HashSHA256")
	flag.StringVar(&c.ConfigJson, "config", "", "json config")
	flag.Parse()
}

// NewServer парсит флаги и env + инициализирует конфиг сервера
func NewServer() *ServerConfig {
	cfg := &ServerConfig{}
	if cfg.ConfigJson != `` {
		err := cfg.formServerJson()
		if err != nil {
			zap.S().Error(err)
		}
	}
	parseServerFlags(cfg)
	err := env.Parse(cfg)

	if err != nil {
		zap.S().Error(err)
	}
	return cfg
}

func parseServerFlags(s *ServerConfig) {
	flag.StringVar(&s.Addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&s.CryptoKey, "c", "", "private crypto-key path")
	flag.IntVar(&s.StoreInterval, "i", 300, "interval for saving metrics on the server")
	flag.StringVar(&s.FilePath, "f", "/tmp/metrics-db.json", "file storage path for saving data")
	flag.BoolVar(&s.Restore, "r", true, "need to load data at startup")
	flag.StringVar(&s.DatabaseDSN, "d", "", "Database Data Source Name")
	flag.StringVar(&s.SignPass, "k", "", "signature for HashSHA256")
	flag.BoolVar(&s.EnableProfiling, "p", false, "run pprof server")
	flag.StringVar(&s.ConfigJson, "config", "", "json config")

	flag.Parse()
}

func (s *ServerConfig) StoreIntervalNotZero() bool {
	return s.StoreInterval != 0
}

func (s *ServerConfig) GetProvider() storage.StorageProvider {
	if s.DatabaseDSN != "" {
		return storage.DBProvider
	}
	if s.FilePath != "" {
		return storage.FileProvider
	}
	return 0
}

// formServerJson дополняет отсутствующие параметры сервера из json
func (s *ServerConfig) formServerJson() error {

	data, err := os.ReadFile(s.ConfigJson)
	if err != nil {
		return fmt.Errorf("cannot read json config: %w", err)
	}

	var settings map[string]interface{}

	err = json.Unmarshal(data, &settings)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json settings: %w", err)
	}

	for stype, value := range settings {
		switch stype {
		case "address":
			if s.Addr == `` {
				s.Addr = value.(string)
			}
		case "restore":
			if !s.Restore {
				s.Restore = value.(bool)
			}
		case "store_interval":
			if s.StoreInterval == 0 {
				duration, err := time.ParseDuration(value.(string))
				if err != nil {
					return fmt.Errorf("bad json param 'store_interval': %w", err)
				}
				s.StoreInterval = int(duration.Seconds())
			}
		case "store_file":
			if s.FilePath == `` {
				s.FilePath = value.(string)
			}
		case "database_dsn":
			if s.DatabaseDSN == `` {
				s.DatabaseDSN = value.(string)
			}
		case "sign_key":
			if s.SignPass == `` {
				s.SignPass = value.(string)
			}
		case "crypto_key":
			if s.CryptoKey == `` {
				s.CryptoKey = value.(string)
			}
		}
	}

	return nil
}

// formServerJson дополняет отсутствующие параметры агента из json
func (c *ClientConfig) formClientJson() error {

	data, err := os.ReadFile(c.ConfigJson)
	if err != nil {
		return fmt.Errorf("cannot read json config: %w", err)
	}

	var settings map[string]interface{}

	err = json.Unmarshal(data, &settings)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json settings: %w", err)
	}

	for stype, value := range settings {
		switch stype {
		case "address":
			if c.Addr == `` {
				c.Addr = value.(string)
			}
		case "report_interval":
			if c.ReportInterval == 0 {
				duration, err := time.ParseDuration(value.(string))
				if err != nil {
					return fmt.Errorf("bad json param 'report_interval': %w", err)
				}
				c.ReportInterval = int(duration.Seconds())
			}
		case "poll_interval":
			if c.PollInterval == 0 {
				duration, err := time.ParseDuration(value.(string))
				if err != nil {
					return fmt.Errorf("bad json param 'poll_interval': %w", err)
				}
				c.PollInterval = int(duration.Seconds())
			}
		case "sign_key":
			if c.SignPass == `` {
				c.SignPass = value.(string)
			}
		case "crypto_key":
			if c.CryptoKey == `` {
				c.CryptoKey = value.(string)
			}
		}
	}

	return nil
}
