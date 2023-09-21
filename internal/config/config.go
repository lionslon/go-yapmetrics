package config

import (
	"flag"
	"os"
	"strconv"
)

type ClientConfig struct {
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	Addr           string `env:"ADDRESS"`
}

type ServerConfig struct {
	Addr string `env:"ADDRESS"`
}

func (c *ClientConfig) New() ClientConfig {
	cfg := &ClientConfig{}
	c.parseFlags()
	c.parseEnv()
	return *cfg
}

func (c *ClientConfig) parseFlags() {
	flag.StringVar(&c.Addr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&c.PollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()
}

func (c *ClientConfig) parseEnv() {
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
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		s.Addr = envRunAddr
	} else {
		flag.StringVar(&s.Addr, "a", "localhost:8080", "address and port to run server")
		flag.Parse()
	}
	return *cfg
}
