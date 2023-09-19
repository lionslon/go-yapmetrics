package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	Addr           string `env:"ADDRESS"`
}

func GetParameters() Config {
	var cfg Config
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&cfg.PollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.Addr = envRunAddr
	}
	if envRunAddr := os.Getenv("REPORT_INTERVAL"); envRunAddr != "" {
		cfg.ReportInterval, _ = strconv.Atoi(envRunAddr)
	}
	if envRunAddr := os.Getenv("POLL_INTERVAL"); envRunAddr != "" {
		cfg.PollInterval, _ = strconv.Atoi(envRunAddr)
	}
	return cfg
}
