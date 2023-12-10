package config

import (
	"flag"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
)

type ClientConfig struct {
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	Addr           string `env:"ADDRESS"`
}

type ServerConfig struct {
	Addr          string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	Restore       bool   `env:"RESTORE"`
}

func (c *ClientConfig) New() ClientConfig {
	cfg := &ClientConfig{}
	parseClientFlags(c)
	err := env.Parse(c)

	if err != nil {
		zap.S().Error(err)
	}
	return *cfg
}

func parseClientFlags(c *ClientConfig) {
	flag.StringVar(&c.Addr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&c.PollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()
}

func (s *ServerConfig) New() ServerConfig {
	cfg := &ServerConfig{}
	parseServerFlags(s)
	err := env.Parse(s)

	if err != nil {
		zap.S().Error(err)
	}
	return *cfg
}

func parseServerFlags(s *ServerConfig) {
	flag.StringVar(&s.Addr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&s.StoreInterval, "i", 300, "interval for saving metrics on the server")
	flag.StringVar(&s.FilePath, "f", "/tmp/metrics-db.json", "file storage path for saving data")
	flag.BoolVar(&s.Restore, "r", true, "need to load data at startup")
	flag.Parse()
}
