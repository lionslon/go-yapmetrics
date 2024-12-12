package main

import (
	"github.com/lionslon/go-yapmetrics/internal/api"
	"github.com/lionslon/go-yapmetrics/internal/config"
)

func main() {
	config.PrintBuildInfo()
	s := api.New()
	if err := s.Start(); err != nil {
		panic(err)
	}
}
