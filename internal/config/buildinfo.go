package config

import (
	"fmt"
)

// Используйте go run -ldflags "-X main.buildVersion=v1.0.1 -X main.buildDate=11.11.2024 -X main.buildCommit=caferacer"
// Чтобы установить переменные для сборки
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func PrintBuildInfo() {
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
