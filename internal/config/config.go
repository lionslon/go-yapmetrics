package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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
	parseClientEnv(c)
	return *cfg
}

func parseClientFlags(c *ClientConfig) {
	flag.StringVar(&c.Addr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&c.PollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()
}

func parseClientEnv(c *ClientConfig) {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		c.Addr = envRunAddr
	}
	if envRunAddr := os.Getenv("REPORT_INTERVAL"); envRunAddr != "" {
		c.ReportInterval, _ = strconv.Atoi(envRunAddr)
	}
	if envRunAddr := os.Getenv("POLL_INTERVAL"); envRunAddr != "" {
		c.PollInterval, _ = strconv.Atoi(envRunAddr)
	}
}

func (s *ServerConfig) New() ServerConfig {
	cfg := &ServerConfig{}
	parseServerFlags(s)
	parseServerEnv(s)
	return *cfg
}

func parseServerFlags(s *ServerConfig) {
	flag.StringVar(&s.Addr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&s.StoreInterval, "i", 300, "interval for saving metrics on the server")
	flag.StringVar(&s.FilePath, "f", "/tmp/metrics-db.json", "file storage path for saving data")
	flag.BoolVar(&s.Restore, "r", true, "need to load data at startup")
	flag.Parse()
}

func parseServerEnv(s *ServerConfig) {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		s.Addr = envRunAddr
	}
	if envValue := os.Getenv("STORE_INTERVAL"); envValue != "" {
		value, err := strconv.Atoi(envValue)
		if err != nil {
			fmt.Println(err)
		}
		s.StoreInterval = value
	}
	if envValue := os.Getenv("FILE_STORAGE_PATH"); envValue != "" {
		s.FilePath = envValue
	}
	if envValue := os.Getenv("RESTORE"); envValue != "" {
		value, err := strconv.ParseBool(envValue)
		if err != nil {
			fmt.Println(err)
		}
		s.Restore = value
	}
}
