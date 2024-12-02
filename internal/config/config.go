package config

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/lionslon/go-yapmetrics/internal/storage"
	"go.uber.org/zap"
)

// ClientConfig конфиг агента
type ClientConfig struct {
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
	Addr           string `env:"ADDRESS"`
	SignPass       string `env:"KEY"`
}

// ServerConfig конфиг сервера
type ServerConfig struct {
	Addr            string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FilePath        string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	SignPass        string `env:"KEY"`
	EnableProfiling bool   `env:"ENABLE_PROFILING"`
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
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&c.RateLimit, "l", 10, "rate limit")
	flag.IntVar(&c.PollInterval, "p", 2, "poll interval in seconds")
	flag.StringVar(&c.SignPass, "k", "", "signature for HashSHA256")
	flag.Parse()
}

// NewServer парсит флаги и env + инициализирует конфиг сервера
func NewServer() *ServerConfig {
	cfg := &ServerConfig{}
	parseServerFlags(cfg)
	err := env.Parse(cfg)

	if err != nil {
		zap.S().Error(err)
	}
	return cfg
}

func parseServerFlags(s *ServerConfig) {
	flag.StringVar(&s.Addr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&s.StoreInterval, "i", 300, "interval for saving metrics on the server")
	flag.StringVar(&s.FilePath, "f", "/tmp/metrics-db.json", "file storage path for saving data")
	flag.BoolVar(&s.Restore, "r", true, "need to load data at startup")
	flag.StringVar(&s.DatabaseDSN, "d", "", "Database Data Source Name")
	flag.StringVar(&s.SignPass, "k", "", "signature for HashSHA256")
	flag.BoolVar(&s.EnableProfiling, "p", false, "run pprof server")

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
