package main

import (
	"github.com/lionslon/go-yapmetrics/internal/api"
	"github.com/lionslon/go-yapmetrics/internal/config"
	"github.com/lionslon/go-yapmetrics/internal/grpcserver"
	"log"
)

func main() {
	config.PrintBuildInfo()
	s := api.New()

	grpcSrv := grpcserver.NewServer()
	go func() {
		if err := grpcserver.StartGRPCServer(":9090", grpcSrv); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()
	log.Println("gRPC server is running on :9090")

	s.Start()
}
