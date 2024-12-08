package main

import (
	"fmt"
	"github.com/lionslon/go-yapmetrics/internal/api"
)

// Используйте go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildDate=11.11.2024 -X main.buildCommit=caferacer"
// Чтобы установить переменные для сборки
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildInfo()
	s := api.New()
	if err := s.Start(); err != nil {
		panic(err)
	}
}

func printBuildInfo() {
	printOrPanic(fmt.Sprintf("Build version: %v", buildVersion))
	printOrPanic(fmt.Sprintf("Build date: %v", buildDate))
	printOrPanic(fmt.Sprintf("Build commit: %v", buildCommit))
}

func printOrPanic(data string) {
	_, err := fmt.Println(data)
	if err != nil {
		panic(err)
	}
}
